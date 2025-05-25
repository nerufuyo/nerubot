# NeruBot - Discord Music Bot

A clean, efficient Discord music bot with high-quality audio streaming and advanced queue management.

## ✨ Features

### 🎵 Music Streaming
- Play music from YouTube, Spotify, and SoundCloud
- Advanced playback controls (pause, resume, skip, stop)
- Loop modes: off, single song, or entire queue
- 24/7 mode for continuous playback
- Auto-disconnect after 5 minutes of inactivity
- High-quality audio streaming with FFmpeg
- Server-specific configuration

### 🤖 Interactive Help System
- Paginated help menu with category navigation
- Compact command reference card
- Detailed feature showcase
- Interactive UI with buttons for navigation

## 🚀 Quick Start

1. **Install Dependencies**
   ```bash
   pip install -r requirements.txt
   ```

2. **Configure Environment**
   ```bash
   # Create .env file
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

3. **Run the Bot**
   ```bash
   python3 -m src.main
   ```

## 📁 New Improved Architecture

The bot now uses a feature-based modular architecture following DRY and KISS principles:

```
src/
├── main.py                    # Entry point
├── interfaces/               # Discord interface layer
│   └── discord/
│       ├── bot.py           # Main bot class
│       └── music_cog.py     # Music commands
├── features/                # Feature modules (NEW!)
│   ├── news/               # News feature
│   │   ├── cogs/           # Discord commands
│   │   ├── services/       # Business logic
│   │   └── models/         # Data models
│   ├── quotes/             # AI quotes feature
## 🎵 Music Commands

### Basic Commands
- `/play <song>` - Play a song from YouTube, Spotify, or SoundCloud
- `/queue` - Show current queue
- `/skip` - Skip current song
- `/stop` - Stop music and clear queue
- `/pause` / `/resume` - Control playback
- `/join` / `/leave` - Voice channel management

### Advanced Features
- `/loop [mode]` - Set loop mode (off/single/queue)
- `/247` - Toggle 24/7 mode (no auto-disconnect)
- `/nowplaying` - Show currently playing song with details
- `/sources` - Display available music sources

## 🤖 Help Commands

### General Help
- `/help` - Interactive paginated help system with categories
- `/commands` - Compact command reference card
- `/about` - Bot information and statistics
- `/features` - Show current and upcoming features

## 📁 Project Structure

```
src/
├── main.py                 # Entry point
├── interfaces/            # Discord interface layer
│   └── discord/
│       ├── bot.py         # Main bot class
│       └── help_cog.py    # Help system
├── features/              # Feature modules
│   └── music/             # Music feature
│       ├── cogs/          # Discord commands
│       ├── services/      # Business logic
│       └── models/        # Data structures
└── core/                  # Core utilities
    └── utils/             # Shared utilities
```

## 🛠️ Development

### Architecture Benefits

- **🧩 Modular:** Each feature is independent
- **🔧 Maintainable:** Clear separation of concerns
- **📈 Scalable:** Easy to add/remove features
- **🧪 Testable:** Services can be tested independently
- **♻️ DRY:** Shared utilities prevent code duplication
- **💋 KISS:** Simple, clean interfaces

## 🔧 Configuration

Environment variables in `.env`:

```bash
# Required
DISCORD_TOKEN=your_bot_token

# Optional - For Spotify Support
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# Optional
COMMAND_PREFIX=!
LOG_LEVEL=INFO
```

## 📊 Testing

Test the news feature:
```bash
python3 test_news.py
```

## 🤝 Contributing

1. Follow the feature-based architecture
2. Use the shared utilities for common functionality
3. Write tests for new services
4. Follow Python best practices and type hints

## 📝 License

MIT License - see LICENSE file for details.

---

**NeruBot v2.0** - Now with advanced features and improved architecture! 🚀
