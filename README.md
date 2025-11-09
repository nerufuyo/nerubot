# NeruBot# NeruBot - Your Ultimate Discord Companion 🎵



A feature-rich Discord bot built with Golang, following Clean Architecture principles.<div align="center">



[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)![NeruBot Banner](https://imgur.com/yh3j7PK.png)

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[![Discord](https://img.shields.io/badge/Discord-Bot-7289DA?style=flat&logo=discord)](https://discord.com)[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/oauth2/authorize?client_id=yourid&permissions=8&scope=bot%20applications.commands)

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)

## 🎯 Features[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

[![Version](https://img.shields.io/badge/Version-3.0.0-blue?style=for-the-badge)](CHANGELOG.md)

### 🎵 Music System

- YouTube audio streaming with high quality**A powerful, feature-rich Discord bot built with Go - bringing music, community engagement, and entertainment to your server**

- Queue management (add, skip, pause, resume, stop)

- Now playing information[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [📖 Documentation](#-documentation) • [🤝 Support](#-support)

- Playback controls

</div>

### 🤐 Confession System

- Anonymous confession submission---

- Moderation queue with approval/rejection

- Reply system for confessions## 🎯 About NeruBot

- Guild-specific settings

NeruBot is a comprehensive Discord companion created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. Built with Go for superior performance and reliability, NeruBot delivers high-quality multi-source music streaming, anonymous confession systems, personalized user roasting, and intuitive slash commands.

### 🔥 Roast System

- User activity tracking (messages, voice time, commands)### 🏆 Why Choose NeruBot?

- AI-powered roast generation

- Statistics and leaderboards- **⚡ Lightning Fast** - Built with Go for exceptional performance and low resource usage

- Multiple roast categories- **🎵 Premium Audio Quality** - Crystal-clear streaming from YouTube with yt-dlp integration

- **🛡️ Privacy-First Design** - Anonymous features with robust security measures

### 🤖 AI Chatbot- **🏗️ Clean Architecture** - Maintainable, scalable codebase following industry best practices

- Multi-provider AI integration (Claude, Gemini, OpenAI)- **🔒 Production Ready** - Enterprise-grade architecture with thread-safe operations

- Automatic fallback between providers- **💰 Completely Free** - No premium features, everything included!

- Session management with 30-minute timeout

- Context-aware conversations---



### 📰 News System## ✨ Features

- RSS feed aggregation from multiple sources

- Concurrent news fetching<table>

- Customizable news sources<tr>

- Auto-publishing capability<td width="50%">



### 🐋 Whale Alerts### 🎵 **Music System**

- Cryptocurrency whale transaction monitoring- **YouTube Support** - High-quality audio streaming via yt-dlp

- Real-time alerts for large transactions- **Smart Queue Management** - Add, skip, stop, shuffle songs

- Configurable minimum threshold- **Loop Modes** - None, single song, or entire queue

- Support for multiple blockchains- **Voice State Detection** - Automatic channel validation

- **Rich Embeds** - Beautiful now playing displays

## 📋 Prerequisites- **Thread-Safe Operations** - Concurrent queue management



- Go 1.21 or higher### 📝 **Anonymous Confession System**

- FFmpeg (for audio processing)- **Complete Anonymity** - Secure, private confession sharing

- yt-dlp (for YouTube downloads)- **Image Support** - Attach images to confessions

- Discord Bot Token- **Reply System** - Anonymous community engagement

- (Optional) AI API keys for chatbot feature- **Moderation Queue** - Admin approval workflow

- **Settings Management** - Per-server configuration

## 🚀 Quick Start- **Thread-Safe Storage** - JSON-based persistence



### 1. Clone the Repository</td>

<td width="50%">

```bash

git clone https://github.com/nerufuyo/nerubot.git### � **User Roasting System**

cd nerubot- **Activity Tracking** - Monitor messages, reactions, voice time

```- **Smart Patterns** - 8 roast categories (spammer, lurker, etc.)

- **Profile Analysis** - Detailed user behavior insights

### 2. Install Dependencies- **Statistics** - Comprehensive roast metrics

- **Safety Systems** - Cooldowns and friendly content

```bash- **Persistent Data** - JSON storage for long-term tracking

# Install Go dependencies

go mod download### ℹ️ **User-Friendly Interface**

- **Slash Commands** - Modern Discord command system

# Install FFmpeg (macOS)- **Rich Embeds** - Beautiful, consistent message formatting

brew install ffmpeg- **Error Handling** - Comprehensive error messages

- **Interactive Components** - Buttons and modals

# Install FFmpeg (Ubuntu/Debian)- **Multi-Feature Support** - Easy feature toggling via config

sudo apt-get install ffmpeg

</td>

# Install yt-dlp</tr>

pip install yt-dlp</table>

# or

brew install yt-dlp**🚧 Coming Soon:**

```- 🤖 AI-Powered Chatbot (Claude, Gemini, OpenAI)

- 📰 Real-Time News & Alerts

### 3. Configure Environment- 💰 Crypto Whale Alerts

- 📊 Advanced Analytics Dashboard

```bash

cp .env.example .env---

```

## 🚀 Quick Start

Edit `.env` and add your configuration:

### 🖥️ Local Development

```env

# Required```bash

DISCORD_TOKEN=your_discord_bot_token# Clone the repository

DISCORD_GUILD_ID=your_guild_idgit clone https://github.com/nerufuyo/nerubot.git

cd nerubot

# AI Providers (at least one required for chatbot)

ANTHROPIC_API_KEY=your_claude_api_key# Copy environment template

GEMINI_API_KEY=your_gemini_api_keycp .env.example .env

OPENAI_API_KEY=your_openai_api_key

# Edit .env with your Discord bot token and configuration

# Optionalnano .env  # or use your preferred editor

WHALE_ALERT_API_KEY=your_whale_alert_api_key

# Build the bot

# Feature Flagsmake build

ENABLE_MUSIC=true

ENABLE_CONFESSION=true# Run the bot

ENABLE_ROAST=truemake run

ENABLE_CHATBOT=true```

ENABLE_NEWS=true

ENABLE_WHALE_ALERT=true### Prerequisites

```- **Go 1.21+** - [Download](https://go.dev/dl/)

- **FFmpeg** - For audio processing

### 4. Build and Run- **yt-dlp** - For YouTube downloads



```bash**Install dependencies (macOS):**

# Build the bot```bash

make buildbrew install ffmpeg yt-dlp

```

# Run the bot

./build/nerubot**Install dependencies (Ubuntu/Debian):**

```bash

# Or build and run in one stepsudo apt update

make runsudo apt install -y ffmpeg

```sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp

sudo chmod a+rx /usr/local/bin/yt-dlp

## 🎮 Commands```



### Music Commands### 🌐 VPS Deployment (Production)

- `/play <url>` - Play a YouTube video

- `/skip` - Skip the current song```bash

- `/pause` - Pause playback# One-command VPS setup (Ubuntu/Debian)

- `/resume` - Resume playbackcurl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/setup.sh | sudo bash

- `/stop` - Stop playback and clear queue```

- `/queue` - Show the current queue

- `/nowplaying` - Show currently playing song**What this does:**

- 🔧 Installs Go, FFmpeg, yt-dlp, and dependencies

### Confession Commands- 👤 Creates secure `nerubot` user

- `/confess <message>` - Submit an anonymous confession- 🛡️ Configures firewall (SSH only)

- `/confess-reply <id> <message>` - Reply to a confession (Admin)- ⚙️ Sets up systemd service

- `/confess-approve <id>` - Approve a confession (Admin)- 📊 Enables health monitoring

- `/confess-reject <id>` - Reject a confession (Admin)

---

### Roast Commands

- `/roast <user>` - Roast a user based on their activity## 📋 Command Reference

- `/roast-stats` - Show roast statistics

- `/roast-leaderboard` - Show roast leaderboard### 🎵 Music Commands

| Command | Description | Status |

### AI Chatbot Commands|---------|-------------|--------|

- `/chat <message>` - Chat with the AI bot| `/play <song>` | Play music from YouTube | ✅ |

- `/chat-reset` - Clear your chat session history| `/skip` | Skip to the next song | ✅ |

| `/stop` | Stop playback and clear queue | ✅ |

### News Commands| `/queue` | Display current music queue | ✅ |

- `/news [limit]` - Fetch latest news articles (default: 5)

### 📝 Confession Commands

### Whale Alert Commands| Command | Description | Status |

- `/whale [limit]` - Get recent crypto whale transactions (default: 5)|---------|-------------|--------|

| `/confess` | Submit anonymous confession (modal) | ✅ |

### Utility Commands

- `/help` - Show all available commands### 🔥 Roast Commands

| Command | Description | Status |

## 🏗️ Architecture|---------|-------------|--------|

| `/roast [user]` | Generate personalized roast | ✅ |

NeruBot follows Clean Architecture principles with clear separation of concerns:

### ℹ️ Information Commands

```| Command | Description | Status |

nerubot/|---------|-------------|--------|

├── cmd/nerubot/              # Application entry point| `/help` | Display help information | ✅ |

│   └── main.go

├── internal/**🚧 Additional commands will be added as features are completed**

│   ├── config/               # Configuration management

│   │   ├── config.go---

│   │   ├── constants.go

│   │   └── messages.go## ⚙️ Configuration

│   ├── entity/               # Domain models

│   │   ├── confession.go### 🔑 Environment Setup

│   │   ├── music.go

│   │   ├── news.go1. **Create `.env` file from template:**

│   │   ├── roast.go```bash

│   │   └── whale.gocp .env.example .env

│   ├── repository/           # Data persistence```

│   │   ├── confession_repository.go

│   │   ├── roast_repository.go2. **Configure required settings:**

│   │   └── repository.go```env

│   ├── usecase/              # Business logic# Required

│   │   ├── chatbot/BOT_TOKEN=your_discord_bot_token_here

│   │   ├── confession/BOT_PREFIX=!

│   │   ├── music/BOT_NAME=NeruBot

│   │   ├── news/

│   │   ├── roast/# Discord

│   │   └── whale/DISCORD_GUILD_ID=your_guild_id_here

│   ├── delivery/             # External interfacesDISCORD_OWNER_ID=your_user_id_here

│   │   └── discord/

│   │       ├── bot.go# Feature Toggles

│   │       └── handlers.goFEATURE_MUSIC=true

│   └── pkg/                  # Shared utilitiesFEATURE_CONFESSION=true

│       ├── ai/               # AI provider implementationsFEATURE_ROAST=true

│       ├── ffmpeg/           # FFmpeg wrapperFEATURE_CHATBOT=false

│       ├── logger/           # Logging utilitiesFEATURE_NEWS=false

│       └── ytdlp/            # yt-dlp wrapperFEATURE_WHALE_ALERTS=false

├── data/                     # Data storage```

│   ├── confessions/

│   └── roasts/3. **Get Discord Bot Token:**

└── deploy/                   # Deployment configurations   - Visit [Discord Developer Portal](https://discord.com/developers/applications)

    ├── systemd/   - Create new application → Bot → Copy token

    ├── nginx/   - Enable all necessary intents (Server Members, Message Content)

    ├── logrotate/

    └── cron/4. **Bot Permissions:**

```   - Send Messages

   - Embed Links

### Architecture Layers   - Read Message History

   - Connect to Voice

1. **Config Layer** - Environment configuration and settings   - Speak in Voice

2. **Delivery Layer** - Discord bot interface and command handlers   - Use Slash Commands

3. **Use Case Layer** - Business logic and service orchestration

4. **Entity Layer** - Domain models and data structures### 🎛️ Advanced Configuration

5. **Repository Layer** - Data persistence (JSON files)

6. **Pkg Layer** - Shared utilities and external integrationsAll configuration is managed through environment variables. See [`.env.example`](.env.example) for all available options:

- Bot behavior and status

## 🐳 Docker Deployment- Feature toggles

- Resource limits

### Build Docker Image- Audio settings

- AI configuration

```bash- Logging preferences

docker build -t nerubot:latest .

```---



### Run with Docker## 🛠️ Management & Monitoring



```bash### 📊 Service Management

docker run -d \```bash

  --name nerubot \# Check bot status

  --env-file .env \sudo systemctl status nerubot

  -v $(pwd)/data:/app/data \

  -v $(pwd)/logs:/app/logs \# View real-time logs

  nerubot:latestsudo journalctl -u nerubot -f

```

# Restart service

### Docker Composesudo systemctl restart nerubot

```

```bash

docker-compose up -d### 📈 Monitoring Tools

``````bash

# Quick status dashboard

## 🔧 Development./deploy/status.sh



### Building# Health monitoring

./deploy/monitor.sh

```bash

# Development build# Update bot to latest version

go build ./..../deploy/update.sh

```

# Production build (optimized)

make build---



# Clean build artifacts## 🏗️ Architecture

make clean

```NeruBot follows **Clean Architecture** principles for maximum maintainability and testability:



### Testing```

internal/

```bash├── config/                 # Configuration and constants

# Run all tests│   ├── config.go          # Environment configuration

go test ./...│   ├── messages.go        # Bot messages and responses

│   └── constants.go       # Application constants

# Run tests with coverage├── entity/                # Domain models (business entities)

go test -cover ./...│   ├── music.go           # Music domain models

│   ├── confession.go      # Confession domain models

# Run tests for specific package│   ├── roast.go           # Roast domain models

go test ./internal/usecase/music/│   ├── news.go            # News domain models

```│   └── whale.go           # Whale alert domain models

├── repository/            # Data persistence layer

### Code Quality│   ├── repository.go      # Base JSON repository

│   ├── confession_repository.go

```bash│   └── roast_repository.go

# Format code├── usecase/               # Business logic layer

go fmt ./...│   ├── music/             # Music service

│   ├── confession/        # Confession service

# Vet code│   ├── roast/             # Roast service

go vet ./...│   ├── chatbot/           # AI chatbot service

```│   ├── news/              # News service

│   └── whale/             # Whale alerts service

## 📊 Performance├── delivery/              # Interface layer

│   └── discord/           # Discord bot implementation

- **Binary Size:** ~8-10MB (optimized)│       ├── bot.go         # Bot lifecycle and setup

- **Memory Usage:** ~50-100MB (varies with features)│       └── handlers.go    # Command handlers

- **Startup Time:** <2 seconds└── pkg/                   # Shared utilities

- **Audio Latency:** <100ms    ├── logger/            # Structured logging

    ├── ffmpeg/            # FFmpeg wrapper

## 🔒 Security    └── ytdlp/             # yt-dlp wrapper

```

- Environment variables for sensitive data

- No hardcoded credentials**Key Principles:**

- Secure session management- 🏛️ **Clean Architecture** - Clear separation of concerns

- Rate limiting on commands- 🧹 **SOLID Principles** - Well-designed, maintainable code

- Admin-only commands for moderation- 🔒 **Thread Safety** - Concurrent operations with sync.RWMutex

- 📈 **Scalable** - Ready for high-traffic servers

## 🤝 Contributing- 🧪 **Testable** - Dependency injection for easy testing



Contributions are welcome! Please follow these guidelines:**Data Flow:**

```

1. Fork the repositoryDiscord → Delivery → Use Case → Entity

2. Create your feature branch (`git checkout -b feature/AmazingFeature`)                ↓

3. Commit your changes following [commit format guidelines](docs/format-commit.md)           Repository → JSON Files

4. Push to the branch (`git push origin feature/AmazingFeature`)```

5. Open a Pull Request

---

## 📝 License

## 📊 System Requirements

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Minimum Requirements

## 🙏 Acknowledgments- **OS:** Ubuntu 20.04+ / Debian 11+ / Windows 10+ / macOS 10.15+

- **Go:** 1.21 or higher

- [DiscordGo](https://github.com/bwmarrin/discordgo) - Discord API library for Go- **RAM:** 512MB

- [gofeed](https://github.com/mmcdole/gofeed) - RSS feed parser- **Storage:** 2GB

- [FFmpeg](https://ffmpeg.org/) - Audio processing- **Network:** Stable internet connection

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - YouTube download utility

### Recommended (VPS)

## 📞 Support- **CPU:** 2+ cores

- **RAM:** 1GB+

- **Issues:** [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)- **Storage:** 5GB+

- **Documentation:** [docs/](docs/)- **Bandwidth:** 500GB/month



## 🗺️ Roadmap### Dependencies

- **Go 1.21+** - Programming language

- [ ] Unit tests for all packages- **FFmpeg** - Audio processing

- [ ] Integration tests- **yt-dlp** - YouTube downloads

- [ ] CI/CD pipeline- **Git** - Version control

- [ ] Database support (PostgreSQL/MongoDB)

- [ ] Web dashboard---

- [ ] Metrics and monitoring

- [ ] Multi-guild support improvements## 📖 Documentation

- [ ] Additional music sources

| Document | Description |

---|----------|-------------|

| **[🚀 Deployment Guide](deploy/README.md)** | Complete VPS setup and management |

**Made with ❤️ by [@nerufuyo](https://github.com/nerufuyo)**| **[🤝 Contributing Guide](CONTRIBUTING.md)** | Development guidelines and setup |

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