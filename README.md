# 🎵 NeruBot - Advanced Discord Music Bot

> **A modern, feature-rich Discord music bot with clean architecture and professional-grade audio streaming.**

[![Python 3.8+](https://img.shields.io/badge/python-3.8+-blue.svg)](https://www.python.org/downloads/)
[![Discord.py](https://img.shields.io/badge/discord.py-2.0+-blue.svg)](https://github.com/Rapptz/discord.py)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-enabled-blue.svg)](https://www.docker.com/)

---

## ✨ Features

### 🎵 **Multi-Source Music Streaming**
- **YouTube** - Direct video playback with high-quality audio
- **Spotify** - Track and playlist support (searches equivalent on YouTube)
- **SoundCloud** - Direct streaming from SoundCloud tracks
- **Direct Links** - Support for MP3, MP4, and other audio formats

### 🎛️ **Advanced Playback Controls**
- Smart queue management with unlimited songs
- Loop modes: Off, Single Song, or Entire Queue
- 24/7 mode for continuous server presence
- Auto-disconnect after 5 minutes of inactivity
- High-quality audio with FFmpeg optimization
- Server-specific configurations

### 🤖 **Interactive Interface**
- Modern slash command system (`/play`, `/queue`, etc.)
- Paginated help system with category navigation
- Real-time now-playing information with thumbnails
- Interactive command reference cards
- Rich embed responses with music source indicators

### 🛡️ **Enterprise-Grade Architecture**
- Clean, modular codebase following SOLID principles
- Feature-based organization for easy maintenance
- Comprehensive error handling and logging
- Docker support for containerized deployment
- VPS deployment scripts for production use

---

## 🚀 Quick Start

### **Prerequisites**
- Python 3.8 or higher
- FFmpeg installed on your system
- Discord Bot Token ([Get one here](https://discord.com/developers/applications))

### **1. Clone & Setup**
```bash
# Clone the repository
git clone https://github.com/yourusername/nerubot.git
cd nerubot

# Quick setup and run
./start.sh setup
```

### **2. Configure Your Bot**
The script will create a `.env` file template. Edit it with your Discord bot token:
```bash
# Edit the environment file
nano .env

# Add your Discord bot token
DISCORD_TOKEN=your_actual_discord_bot_token_here
```

### **3. Run the Bot**
```bash
# Start the bot
./start.sh run

# Or run in debug mode
./start.sh debug
```

**That's it!** Your bot is now running and ready to play music! 🎉

---

## 📖 Usage Guide

### **Basic Commands**
| Command | Description | Example |
|---------|-------------|---------|
| `/play <song>` | Play a song or add to queue | `/play Bohemian Rhapsody` |
| `/queue` | Show current music queue | `/queue` |
| `/skip` | Skip the current song | `/skip` |
| `/pause` / `/resume` | Pause/resume playback | `/pause` |
| `/stop` | Stop music and clear queue | `/stop` |

### **Advanced Commands**
| Command | Description | Options |
|---------|-------------|---------|
| `/loop [mode]` | Set loop mode | `off`, `single`, `queue` |
| `/247` | Toggle 24/7 mode | - |
| `/nowplaying` | Show current song details | - |
| `/sources` | Show available music sources | - |
| `/help` | Interactive help system | - |

### **Voice Channel Management**
| Command | Description |
|---------|-------------|
| `/join` | Join your voice channel |
| `/leave` | Leave the voice channel |

### **Help Commands**
| Command | Description |
|---------|-------------|
| `/help` | Interactive paginated help system |
| `/commands` | Compact command reference card |
| `/about` | Bot information and statistics |
| `/features` | Show current and upcoming features |

---

## 🏗️ Project Architecture

NeruBot follows a clean, modular architecture that makes it easy to maintain and extend:

```
nerubot/
├── src/                     # Source code
│   ├── main.py             # Application entry point
│   ├── config/             # Configuration management
│   │   ├── settings.py     # Bot settings and limits
│   │   └── messages.py     # User-facing messages
│   ├── core/               # Core utilities
│   │   └── utils/          # Logging, file operations
│   ├── features/           # Feature modules
│   │   ├── music/          # Music streaming feature
│   │   │   ├── cogs/       # Discord commands
│   │   │   ├── services/   # Business logic
│   │   │   └── models/     # Data models
│   │   ├── help/           # Help system
│   │   └── news/           # News updates (optional)
│   └── interfaces/         # External interfaces
│       └── discord/        # Discord bot interface
├── deploy/                 # Deployment scripts
├── requirements.txt        # Python dependencies
├── start.sh               # Main startup script
├── Dockerfile             # Docker configuration
└── README.md              # This file
```

### **Key Benefits**
- **🧩 Modular**: Each feature is independent and can be easily modified
- **🔧 Maintainable**: Clear separation of concerns and single responsibility
- **📈 Scalable**: Easy to add new features without affecting existing ones
- **🧪 Testable**: Services can be unit tested independently
- **♻️ DRY**: Shared utilities prevent code duplication
- **💋 KISS**: Simple, clean interfaces

---

## 🌐 Production Deployment

### **Docker Deployment (Recommended)**
```bash
# Clone and build
git clone https://github.com/yourusername/nerubot.git
cd nerubot

# Build and run with Docker
docker build -t nerubot .
docker run -d --name nerubot --env-file .env nerubot
```

### **VPS Deployment**
```bash
# On your VPS (Ubuntu/Debian)
curl -fsSL https://raw.githubusercontent.com/yourusername/nerubot/main/deploy/vps_setup.sh | sudo bash
```

### **Manual VPS Setup**
```bash
# Install dependencies
sudo apt update && sudo apt install -y python3 python3-pip ffmpeg git

# Clone and setup
git clone https://github.com/yourusername/nerubot.git
cd nerubot
./start.sh setup

# Run with systemd (persistent)
sudo cp deploy/nerubot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable nerubot
sudo systemctl start nerubot
```

---

## 🔧 Configuration

### **Environment Variables**
Create a `.env` file in the project root:

```bash
# Required
DISCORD_TOKEN=your_discord_bot_token_here

# Optional - Enhanced Spotify Support
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# Optional - Bot Configuration
COMMAND_PREFIX=!
LOG_LEVEL=INFO
MAX_QUEUE_SIZE=100
ENABLE_24_7=true
AUTO_DISCONNECT_TIME=300
```

### **Advanced Configuration**
- **Queue Limits**: Modify `MAX_QUEUE_SIZE` in settings
- **Audio Quality**: FFmpeg options in `src/config/settings.py`
- **Logging**: Adjust `LOG_LEVEL` (DEBUG, INFO, WARNING, ERROR)
- **Timeouts**: Configure auto-disconnect timing

---

## 🛠️ Development

### **Starting Script Options**
```bash
./start.sh setup    # Install dependencies and create .env
./start.sh run      # Start the bot normally
./start.sh debug    # Start with debug logging
./start.sh clean    # Clean cache and temporary files
./start.sh help     # Show all available options
```

### **Testing**
```bash
# Run unit tests
python3 -m unittest discover -s src/features/music/tests

# Test specific feature
python3 -m unittest src.features.music.tests.test_music_service
```

### **Development Setup**
```bash
# Install development dependencies
pip install -r requirements-dev.txt

# Run linting
flake8 src/
black src/

# Type checking
mypy src/
```

---

## 🚨 Troubleshooting

### **Common Issues**

| Issue | Solution |
|-------|----------|
| Bot not responding | Check Discord token in `.env` file |
| Audio quality poor | Ensure FFmpeg is properly installed |
| Commands not showing | Bot needs proper Discord permissions |
| Memory issues | Restart bot: `./start.sh clean && ./start.sh run` |

### **FFmpeg Installation**
```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# CentOS/RHEL
sudo yum install ffmpeg

# macOS
brew install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
```

### **Discord Permissions**
Your bot needs these permissions:
- Send Messages
- Use Slash Commands
- Connect to Voice Channels
- Speak in Voice Channels
- Embed Links

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### **Development Guidelines**
1. Follow the feature-based architecture
2. Use shared utilities for common functionality
3. Write tests for new services
4. Follow Python best practices and type hints
5. Use meaningful commit messages

### **Contribution Process**
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the architecture
4. Add tests for new functionality
5. Submit a pull request

### **Code Standards**
- Use Black for code formatting
- Follow PEP 8 guidelines
- Add type hints to all functions
- Document complex functions with docstrings
- Keep functions small and focused

---

## 📊 Monitoring & Logs

### **Log Files**
```bash
# View current logs
tail -f logs/bot.log

# View error logs only
grep -E "(ERROR|CRITICAL)" logs/bot.log

# Rotate logs (automated daily)
./start.sh clean
```

### **Health Monitoring**
- Bot status: Check Discord presence
- Memory usage: Monitor with `htop` or `ps`
- Queue performance: Check `/queue` response times

---

## 📈 Roadmap

### **Upcoming Features**
- [ ] Web dashboard for bot management
- [ ] Multi-server playlist sharing
- [ ] Advanced audio effects and filters
- [ ] Integration with more music services
- [ ] Machine learning-based music recommendations
- [ ] REST API for external integrations

### **Recently Added**
- [x] Unified startup script
- [x] Docker containerization
- [x] Enhanced help system
- [x] Improved error handling
- [x] Modular architecture refactor

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💝 Support

If you find NeruBot helpful:
- ⭐ Star this repository
- 🐛 Report bugs via [Issues](https://github.com/yourusername/nerubot/issues)
- 💡 Suggest features via [Discussions](https://github.com/yourusername/nerubot/discussions)
- 🤝 Contribute via Pull Requests

---

<div align="center">

**NeruBot v2.0** - Professional Discord Music Bot 🚀

*Built with ❤️ for the Discord community*

</div>
