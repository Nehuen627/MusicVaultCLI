package main

import (
    "bytes"
    "bufio"
    "encoding/json" 
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/exec"
	"runtime"
    "strings"
    "github.com/joho/godotenv"
)
type SpotifyTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    RefreshToken string `json:"refresh_token"`
}
type SpotifyPlaylistResponse struct {
    Items []SpotifyPlaylist `json:"items"`
}

type SpotifyPlaylist struct {
    Name   string `json:"name"`
    ID     string `json:"id"`
    Tracks struct {
        Total int `json:"total"`
    } `json:"tracks"`
}

type TracksPage struct {
    Next  string `json:"next"` 
    Items []struct {
        Track struct {
            Name    string `json:"name"`
            Artists []struct {
                Name string `json:"name"`
            } `json:"artists"`
        } `json:"track"`
    } `json:"items"`
}

type Song struct {
    Name   string `json:"name"`
    Artist string `json:"artist"`
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

    fmt.Println("Welcome to Music Vault")
	menu()
}

func menu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose an action: [A] Backup playlist, [B] See backed up playlist, [C] Quit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
		case "a":
			fmt.Println("Only spotify is available")
			userToken := fetch_spotify_token()
			playlists := get_playlists(userToken)
			fetch_and_backup_playlist(playlists, userToken)
		case "b":
			fmt.Println("Reading backup")
		case "c":
			fmt.Println("Quitting...")
		default:
			fmt.Println("Invalid choice, please try again.")
			menu()
	}
}

func fetch_spotify_token() string {
	fmt.Println("Go to your browser and log in to your spotify account...")
	codeChannel := make(chan string)
	loginURL := construct_URL()
	fetch_code(codeChannel)
	go http.ListenAndServe(":8000", nil)
	os := runtime.GOOS
	switch os {
		case "windows" :
			exec.Command("rundll32", "url.dll,FileProtocolHandler", loginURL).Start()
		case "linux" :
			exec.Command("xdg-open", loginURL).Start()
		case "darwin" :
			exec.Command("open", loginURL).Start()
	default :
		fmt.Println("Os not supported...")
	}
	authToken := <- codeChannel
	userToken := change_token_to_user_code(authToken)
	return userToken
}

func construct_URL() string {
	data := url.Values{}
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("response_type", "code")
    data.Set("redirect_uri", "http://127.0.0.1:8000/callback")
	data.Set("scope", "playlist-read-private")

	loginURL := "https://accounts.spotify.com/authorize?" + data.Encode()
	return loginURL
}
func change_token_to_user_code(authToken string) string {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authToken)
	data.Set("redirect_uri", "http://127.0.0.1:8000/callback")
	data.Set("client_id", os.Getenv("CLIENT_ID"))
    data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
    if err != nil {
        panic(err)
    }
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
    resp, err := client.Do(req) 
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var tokenResponse SpotifyTokenResponse
    json.NewDecoder(resp.Body).Decode(&tokenResponse)
    return tokenResponse.AccessToken
}

func fetch_code(pipe chan string) {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		pipe <- string(code)
		fmt.Fprintf(w, "You can close this tab now!")
	})
}
func get_playlists(userToken string) []SpotifyPlaylist {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)
    if err != nil {
        panic(err)
    }
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + userToken)
	client := &http.Client{}
    resp, err := client.Do(req) 
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

	var spotifyPlaylistResponse SpotifyPlaylistResponse 
	json.NewDecoder(resp.Body).Decode(&spotifyPlaylistResponse)
	return spotifyPlaylistResponse.Items
}

func fetch_and_backup_playlist(playlists []SpotifyPlaylist, userToken string) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("\n--- Starting playlist backup section ---")
    for _, playlists:= range playlists {
        fmt.Printf("Do you want to back up '%s'? (y/n): ", playlists.Name)
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(strings.ToLower(input))

        if input == "y" {
            songs := fetch_playlist_songs(playlists.ID, userToken)
            backup_playlist(playlists.Name, songs)
        }
    }
}

func fetch_playlist_songs(playlistID string, userToken string) []Song {
    url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
    
    var allSongs []Song

    fmt.Print("Fetching playlist songs")
    
    for url != "" {
        req, _ := http.NewRequest("GET", url, nil)
        req.Header.Set("Authorization", "Bearer " + userToken)
        
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println("Error:", err)
            break
        }
        
        var page TracksPage
        json.NewDecoder(resp.Body).Decode(&page)
        resp.Body.Close()

        for _, item := range page.Items {
            if item.Track.Name == "" { continue }
            
            artistName := "Unknown"
            if len(item.Track.Artists) > 0 {
                artistName = item.Track.Artists[0].Name
            }

            allSongs = append(allSongs, Song{
                Name:   item.Track.Name,
                Artist: artistName,
            })
        }
        
        fmt.Print(".")
        
        url = page.Next 
    }
    fmt.Println("\nAll pages done!")
    return allSongs
}

func backup_playlist(name string, songs []Song) {
    folderName := "playlists"

    if _, err := os.Stat(folderName); os.IsNotExist(err) {
        os.Mkdir(folderName, 0755)
    }

    fileData, _ := json.MarshalIndent(songs, "", "  ")
    
    safeName := strings.ReplaceAll(name, "/", "-") + ".json"
    
    fullPath := folderName + string(os.PathSeparator) + safeName

    err := os.WriteFile(fullPath, fileData, 0644) 
    
    if err == nil {
        fmt.Printf("Succesfully saved %d songs to '%s'\n", len(songs), fullPath)
    } else {
        fmt.Println("Error saving file:", err)
    }
}