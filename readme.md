# Music Vault

A command-line tool to backup your music playlists locally. Music Vault allows you to create JSON backups of your playlists, ensuring you never lose your carefully curated music collections.

## Introduction

Music Vault is designed to help you preserve your music playlists by creating local backups in a simple JSON format. Whether you're worried about losing access to your account, want to keep a historical record of your playlists, or simply want to have your music data in a portable format, Music Vault has you covered.

## Current Status

**Currently Available:**
- ✅ Spotify playlist backup

**In Development:**
- 🚧 Additional streaming platform support
- 🚧 Cross-platform playlist transfer

## Prerequisites

- Go 1.16 or higher
- A Spotify account
- A Spotify Developer App

## Spotify Setup

Before you can use Music Vault with Spotify, you need to create a Spotify Developer application:

### 1. Create a Spotify App

1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Log in with your Spotify account
3. Click **"Create app"**
4. Fill in the app details:
   - **App name:** Music Vault (or any name you prefer)
   - **App description:** Personal playlist backup tool
   - **Redirect URI:** `http://127.0.0.1:8000/callback`
   - Check the agreement box
5. Click **"Save"**

### 2. Get Your Credentials

1. In your app dashboard, click on **"Settings"**
2. Copy your **Client ID**
3. Click **"View client secret"** and copy your **Client Secret**

### 3. Configure Environment Variables

1. Create a `.env` file in the project root directory
2. Add your Spotify credentials:

```env
CLIENT_ID=your_client_id_here
CLIENT_SECRET=your_client_secret_here
```

**Important:** Never share or commit your `.env` file to version control!

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd music-vault
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your `.env` file as described above

## Usage

Run the application:
```bash
go run main.go
```

### Menu Options

- **[A] Backup playlist** - Authenticate with Spotify and select playlists to backup
- **[B] See backed up playlist** - View previously backed up playlists (coming soon)
- **[C] Quit** - Exit the application

### Backing Up Playlists

1. Select option **A** from the menu
2. Your browser will automatically open for Spotify authentication
3. Log in and authorize the application
4. You'll be prompted for each playlist whether you want to back it up
5. Type `y` to backup a playlist or `n` to skip
6. Backups are saved in the `playlists/` directory as JSON files

## Backup Format

Playlists are saved as JSON files with the following structure:

```json
[
  {
    "name": "Song Title",
    "artist": "Artist Name"
  },
  ...
]
```

## Troubleshooting

**Browser doesn't open automatically:**
- Manually copy the URL from the terminal and paste it into your browser

**Authentication fails:**
- Verify your `CLIENT_ID` and `CLIENT_SECRET` are correct in `.env`
- Ensure the redirect URI in your Spotify app settings is exactly `http://127.0.0.1:8000/callback`

**Port 8000 already in use:**
- Close any applications using port 8000 or modify the port in the code

## Roadmap

- [ ] Add support for additional streaming platforms
- [ ] Implement playlist viewing functionality
- [ ] Enable cross-platform playlist transfer

## License

This project is for personal use. Please respect the terms of service of any streaming platforms you use with this tool.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---
**Created by:** Nehuen
**Note:** This tool is not affiliated with Spotify or any other music streaming service.
