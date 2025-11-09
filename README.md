# NeruBot - Your Ultimate Discord Companion 🎵

<div align="center">

![NeruBot Banner](https://imgur.com/yh3j7PK.png)

[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/oauth2/authorize?client_id=yourid&permissions=8&scope=bot%20applications.commands)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-3.0.0-blue?style=for-the-badge)](CHANGELOG.md)

**A powerful, feature-rich Discord bot built with Go - bringing music, community engagement, and entertainment to your server**

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [📖 Documentation](#-documentation) • [🤝 Support](#-support)

</div>

---

## 🎯 About NeruBot

NeruBot is a comprehensive Discord companion created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. Built with Go for superior performance and reliability, NeruBot delivers high-quality multi-source music streaming, anonymous confession systems, personalized user roasting, and intuitive slash commands.

### 🏆 Why Choose NeruBot?

- **⚡ Lightning Fast** - Built with Go for exceptional performance and low resource usage
- **🎵 Premium Audio Quality** - Crystal-clear streaming from YouTube with yt-dlp integration
- **🛡️ Privacy-First Design** - Anonymous features with robust security measures
- **🏗️ Clean Architecture** - Maintainable, scalable codebase following industry best practices
- **🔒 Production Ready** - Enterprise-grade architecture with thread-safe operations
- **💰 Completely Free** - No premium features, everything included!

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎵 **Music System**
- **YouTube Support** - High-quality audio streaming via yt-dlp
- **Smart Queue Management** - Add, skip, stop, shuffle songs
- **Loop Modes** - None, single song, or entire queue
- **Voice State Detection** - Automatic channel validation
- **Rich Embeds** - Beautiful now playing displays
- **Thread-Safe Operations** - Concurrent queue management

### 📝 **Anonymous Confession System**
- **Complete Anonymity** - Secure, private confession sharing
- **Image Support** - Attach images to confessions
- **Reply System** - Anonymous community engagement
- **Moderation Queue** - Admin approval workflow
- **Settings Management** - Per-server configuration
- **Thread-Safe Storage** - JSON-based persistence

</td>
<td width="50%">

### � **User Roasting System**
- **Activity Tracking** - Monitor messages, reactions, voice time
- **Smart Patterns** - 8 roast categories (spammer, lurker, etc.)
- **Profile Analysis** - Detailed user behavior insights
- **Statistics** - Comprehensive roast metrics
- **Safety Systems** - Cooldowns and friendly content
- **Persistent Data** - JSON storage for long-term tracking

### ℹ️ **User-Friendly Interface**
- **Slash Commands** - Modern Discord command system
- **Rich Embeds** - Beautiful, consistent message formatting
- **Error Handling** - Comprehensive error messages
- **Interactive Components** - Buttons and modals
- **Multi-Feature Support** - Easy feature toggling via config

</td>
</tr>
</table>

**🚧 Coming Soon:**
- 🤖 AI-Powered Chatbot (Claude, Gemini, OpenAI)
- 📰 Real-Time News & Alerts
- 💰 Crypto Whale Alerts
- 📊 Advanced Analytics Dashboard

---

## 🚀 Quick Start

### 🖥️ Local Development

```bash
# Clone the repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Copy environment template
cp .env.example .env

# Edit .env with your Discord bot token and configuration
nano .env  # or use your preferred editor

# Build the bot
make build

# Run the bot
make run
```

### Prerequisites
- **Go 1.21+** - [Download](https://go.dev/dl/)
- **FFmpeg** - For audio processing
- **yt-dlp** - For YouTube downloads

**Install dependencies (macOS):**
```bash
brew install ffmpeg yt-dlp
```

**Install dependencies (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install -y ffmpeg
sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

### 🌐 VPS Deployment (Production)

```bash
# One-command VPS setup (Ubuntu/Debian)
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/setup.sh | sudo bash
```

**What this does:**
- 🔧 Installs Go, FFmpeg, yt-dlp, and dependencies
- 👤 Creates secure `nerubot` user
- 🛡️ Configures firewall (SSH only)
- ⚙️ Sets up systemd service
- 📊 Enables health monitoring

---

## 📋 Command Reference

### 🎵 Music Commands
| Command | Description | Status |
|---------|-------------|--------|
| `/play <song>` | Play music from YouTube | ✅ |
| `/skip` | Skip to the next song | ✅ |
| `/stop` | Stop playback and clear queue | ✅ |
| `/queue` | Display current music queue | ✅ |

### 📝 Confession Commands
| Command | Description | Status |
|---------|-------------|--------|
| `/confess` | Submit anonymous confession (modal) | ✅ |

### 🔥 Roast Commands
| Command | Description | Status |
|---------|-------------|--------|
| `/roast [user]` | Generate personalized roast | ✅ |

### ℹ️ Information Commands
| Command | Description | Status |
|---------|-------------|--------|
| `/help` | Display help information | ✅ |

**🚧 Additional commands will be added as features are completed**

---

## ⚙️ Configuration

### 🔑 Environment Setup

1. **Create `.env` file from template:**
```bash
cp .env.example .env
```

2. **Configure required settings:**
```env
# Required
BOT_TOKEN=your_discord_bot_token_here
BOT_PREFIX=!
BOT_NAME=NeruBot

# Discord
DISCORD_GUILD_ID=your_guild_id_here
DISCORD_OWNER_ID=your_user_id_here

# Feature Toggles
FEATURE_MUSIC=true
FEATURE_CONFESSION=true
FEATURE_ROAST=true
FEATURE_CHATBOT=false
FEATURE_NEWS=false
FEATURE_WHALE_ALERTS=false
```

3. **Get Discord Bot Token:**
   - Visit [Discord Developer Portal](https://discord.com/developers/applications)
   - Create new application → Bot → Copy token
   - Enable all necessary intents (Server Members, Message Content)

4. **Bot Permissions:**
   - Send Messages
   - Embed Links
   - Read Message History
   - Connect to Voice
   - Speak in Voice
   - Use Slash Commands

### 🎛️ Advanced Configuration

All configuration is managed through environment variables. See [`.env.example`](.env.example) for all available options:
- Bot behavior and status
- Feature toggles
- Resource limits
- Audio settings
- AI configuration
- Logging preferences

---

## 🛠️ Management & Monitoring

### 📊 Service Management
```bash
# Check bot status
sudo systemctl status nerubot

# View real-time logs
sudo journalctl -u nerubot -f

# Restart service
sudo systemctl restart nerubot
```

### 📈 Monitoring Tools
```bash
# Quick status dashboard
./deploy/status.sh

# Health monitoring
./deploy/monitor.sh

# Update bot to latest version
./deploy/update.sh
```

---

## 🏗️ Architecture

NeruBot follows **Clean Architecture** principles for maximum maintainability and testability:

```
internal/
├── config/                 # Configuration and constants
│   ├── config.go          # Environment configuration
│   ├── messages.go        # Bot messages and responses
│   └── constants.go       # Application constants
├── entity/                # Domain models (business entities)
│   ├── music.go           # Music domain models
│   ├── confession.go      # Confession domain models
│   ├── roast.go           # Roast domain models
│   ├── news.go            # News domain models
│   └── whale.go           # Whale alert domain models
├── repository/            # Data persistence layer
│   ├── repository.go      # Base JSON repository
│   ├── confession_repository.go
│   └── roast_repository.go
├── usecase/               # Business logic layer
│   ├── music/             # Music service
│   ├── confession/        # Confession service
│   ├── roast/             # Roast service
│   ├── chatbot/           # AI chatbot service
│   ├── news/              # News service
│   └── whale/             # Whale alerts service
├── delivery/              # Interface layer
│   └── discord/           # Discord bot implementation
│       ├── bot.go         # Bot lifecycle and setup
│       └── handlers.go    # Command handlers
└── pkg/                   # Shared utilities
    ├── logger/            # Structured logging
    ├── ffmpeg/            # FFmpeg wrapper
    └── ytdlp/             # yt-dlp wrapper
```

**Key Principles:**
- 🏛️ **Clean Architecture** - Clear separation of concerns
- 🧹 **SOLID Principles** - Well-designed, maintainable code
- 🔒 **Thread Safety** - Concurrent operations with sync.RWMutex
- 📈 **Scalable** - Ready for high-traffic servers
- 🧪 **Testable** - Dependency injection for easy testing

**Data Flow:**
```
Discord → Delivery → Use Case → Entity
                ↓
           Repository → JSON Files
```

---

## 📊 System Requirements

### Minimum Requirements
- **OS:** Ubuntu 20.04+ / Debian 11+ / Windows 10+ / macOS 10.15+
- **Go:** 1.21 or higher
- **RAM:** 512MB
- **Storage:** 2GB
- **Network:** Stable internet connection

### Recommended (VPS)
- **CPU:** 2+ cores
- **RAM:** 1GB+
- **Storage:** 5GB+
- **Bandwidth:** 500GB/month

### Dependencies
- **Go 1.21+** - Programming language
- **FFmpeg** - Audio processing
- **yt-dlp** - YouTube downloads
- **Git** - Version control

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| **[🚀 Deployment Guide](deploy/README.md)** | Complete VPS setup and management |
| **[🤝 Contributing Guide](CONTRIBUTING.md)** | Development guidelines and setup |
| **[🏗️ Architecture Overview](ARCHITECTURE.md)** | Technical architecture details |
| **[📝 Changelog](CHANGELOG.md)** | Version history and updates |
| **[📋 Feature Guides](src/features/)** | Individual feature documentation |

---

## 🤝 Contributing

We welcome contributions! NeruBot is built with ❤️ by the community.

### Quick Contribution Guide
1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/yourusername/nerubot.git`
3. **Create** feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes following our [coding standards](CONTRIBUTING.md)
5. **Test** thoroughly: `make test`
6. **Build** to verify: `make build`
7. **Submit** pull request

### Development Setup
```bash
# Install Go dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Code formatting
gofmt -s -w .
go vet ./...
```

**Contribution Areas:**
- 🎵 Music features and sources
- 🛡️ Security improvements
- 📱 Discord interaction enhancements
- 📚 Documentation
- 🧪 Testing coverage
- 🌐 Internationalization
- ⚡ Performance optimization

---

## 💫 Support & Community

### 🆘 Getting Help
- **[GitHub Issues](https://github.com/nerufuyo/nerubot/issues)** - Bug reports and feature requests
- **[Discussions](https://github.com/nerufuyo/nerubot/discussions)** - Questions and community chat
- **[Discord Server](https://discord.gg/yourserver)** - Real-time support and community
- **[Documentation](https://github.com/nerufuyo/nerubot/wiki)** - Comprehensive guides

### 🏷️ Project Status
- ✅ **Active Development** - Regular updates and improvements
- 🛡️ **Production Ready** - Used in 100+ Discord servers
- 🧪 **Well Tested** - Comprehensive test suite
- 📚 **Documented** - Complete documentation and guides

---

## 🙏 Acknowledgments

**Created with ❤️ by [@nerufuyo](https://github.com/nerufuyo)**

Special thanks to:
- **Discord.py Community** - Amazing framework and support
- **Contributors** - Everyone who helped improve NeruBot
- **Users** - Servers and communities using NeruBot
- **Open Source Projects** - Libraries and tools that make this possible

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

**TL;DR:** You can use, modify, and distribute this code freely, just keep the license notice.

---

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=nerufuyo/nerubot&type=Timeline)](https://star-history.com/#nerufuyo/nerubot&Timeline)

---

<div align="center">

**Made with 💖 by the NeruBot Team**

[⭐ Star on GitHub](https://github.com/nerufuyo/nerubot) • [🐛 Report Bug](https://github.com/nerufuyo/nerubot/issues) • [💡 Request Feature](https://github.com/nerufuyo/nerubot/issues)

</div>