<div align="center">

# 🎵 NeruBot

### Your Ultimate Discord Companion

[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-3.0.0-blue?style=for-the-badge)](CHANGELOG.md)

**A powerful, feature-rich Discord bot built with Go - bringing music, community engagement, and entertainment to your server**

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [📖 Documentation](#-documentation) • [🤝 Contributing](#-contributing)

</div>

---

## 📖 Table of Contents

- [About](#-about-nerubot)
- [Features](#-features)
- [Quick Start](#-quick-start)
- [Commands](#-commands)
- [Configuration](#️-configuration)
- [Architecture](#️-architecture)
- [Deployment](#-deployment)
- [Documentation](#-documentation)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 About NeruBot

NeruBot is a comprehensive Discord companion created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. Built with Go for superior performance and reliability, NeruBot follows **Clean Architecture** principles for maintainability and scalability.

### 🏆 Why Choose NeruBot?

- **⚡ Lightning Fast** - Built with Go for exceptional performance
- **🎵 Premium Audio** - Crystal-clear YouTube streaming via yt-dlp
- **🛡️ Privacy-First** - Anonymous features with robust security
- **🏗️ Clean Architecture** - Maintainable, scalable codebase
- **🔒 Production Ready** - Thread-safe operations and error handling
- **💰 Completely Free** - No premium features, everything included!

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎵 **Music System**
- ✅ YouTube audio streaming (yt-dlp)
- ✅ Queue management & controls
- ✅ Loop modes (none/single/queue)
- ✅ Now playing with rich embeds
- ✅ Voice state detection
- ✅ Thread-safe operations

### 📝 **Confession System**
- ✅ Complete anonymity
- ✅ Image attachment support
- ✅ Moderation queue
- ✅ Reply system
- ✅ Per-guild settings
- ✅ Confession numbering

</td>
<td width="50%">

### 🔥 **Roast System**
- ✅ Activity tracking
- ✅ Smart pattern detection
- ✅ Profile analysis
- ✅ Leaderboards & stats
- ✅ Cooldown management
- ✅ 8 roast categories

### 🤖 **AI Chatbot** (Coming Soon)
- 🚧 Multi-provider support
- 🚧 DeepSeek integration
- 🚧 Context-aware conversations
- 🚧 Session management

### 📰 **Additional Features** (Planned)
- 🚧 RSS News aggregation
- 🚧 Crypto whale alerts
- 🚧 Advanced analytics

</td>
</tr>
</table>

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **FFmpeg** - For audio processing
- **yt-dlp** - For YouTube downloads
- **Discord Bot Token** - [Create a bot](https://discord.com/developers/applications)

### Installation

```bash
# Clone the repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your configuration
nano .env

# Edit .env with your configuration
nano .env

# Install system dependencies
# macOS
brew install ffmpeg
python -m pip install yt-dlp

# Ubuntu/Debian
sudo apt update
sudo apt install -y ffmpeg python3-pip
pip3 install yt-dlp

# Build and run
make build
make run
```

### Docker Setup

```bash
# Build and run with Docker
docker-compose up -d

# View logs
docker-compose logs -f

# Stop bot
docker-compose down
```

---

## 🎮 Commands

### 🎵 Music Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/play <song>` | Play music from YouTube | `/play never gonna give you up` |
| `/skip` | Skip to next song | `/skip` |
| `/stop` | Stop playback and clear queue | `/stop` |
| `/pause` | Pause current playback | `/pause` |
| `/resume` | Resume playback | `/resume` |
| `/queue` | Display current queue | `/queue` |
| `/nowplaying` | Show current song info | `/nowplaying` |

### 📝 Confession Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/confess` | Submit anonymous confession | Opens modal |
| `/confess-approve <id>` | Approve confession (Admin) | `/confess-approve 5` |
| `/confess-reject <id>` | Reject confession (Admin) | `/confess-reject 3` |
| `/confess-reply <id>` | Reply to confession (Admin) | Opens modal |

### 🔥 Roast Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/roast [@user]` | Generate personalized roast | `/roast @username` |
| `/profile [@user]` | View user activity profile | `/profile @username` |
| `/leaderboard` | Show roast leaderboard | `/leaderboard` |

### ℹ️ Utility Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/ping` | Check bot response time | `/ping` |
| `/help` | Display help information | `/help` |
| `/about` | Show bot information | `/about` |

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# === REQUIRED SETTINGS ===
# Discord Bot Token (Get from: https://discord.com/developers/applications)
DISCORD_TOKEN=your_discord_bot_token_here

# === AI CHATBOT SETTINGS ===
# DeepSeek API Key (Get from: https://platform.deepseek.com/)
DEEPSEEK_API_KEY=your_deepseek_api_key_here

# === OPTIONAL MUSIC SETTINGS ===
# Spotify Integration (Optional - for better music search)
# Get from: https://developer.spotify.com/dashboard/applications
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# === BOT CONFIGURATION ===
# Bot command prefix (default: !)
COMMAND_PREFIX=!

# Logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
LOG_LEVEL=INFO

# Enable/Disable Features
ENABLE_MUSIC=true
ENABLE_CHATBOT=true
ENABLE_CONFESSION=true
ENABLE_ROAST=true
ENABLE_NEWS=false
ENABLE_WHALE_ALERTS=false

# === MUSIC SETTINGS ===
# Maximum songs in queue per server
MAX_QUEUE_SIZE=100

# Auto-disconnect timeout in seconds (0 = disabled)
AUTO_DISCONNECT_TIME=300

# === ADVANCED SETTINGS ===
# Bot activity status
BOT_STATUS=🎵 Music for everyone!

# Database settings (if using database features)
DATABASE_URL=mongodb://localhost:27017

# Redis settings (if using caching)
REDIS_URL=redis://localhost:6379
```

### Feature Flags

Control which features are enabled:

```env
ENABLE_MUSIC=true          # Music streaming
ENABLE_CONFESSION=true     # Anonymous confessions
ENABLE_ROAST=true          # User roasting
ENABLE_CHATBOT=false       # AI chatbot (requires API key)
ENABLE_NEWS=false          # News aggregation
ENABLE_WHALE_ALERTS=false  # Crypto whale alerts
```

---

## 🏗️ Architecture

NeruBot follows **Clean Architecture** principles with clear separation of concerns:

```
nerubot/
├── cmd/
│   └── nerubot/              # Application entry point
│       └── main.go
├── internal/
│   ├── config/               # Configuration management
│   │   ├── config.go         # Main config structure
│   │   ├── constants.go      # Constants and defaults
│   │   └── messages.go       # Response messages
│   ├── entity/               # Domain models
│   │   ├── confession.go     # Confession entities
│   │   ├── music.go          # Music entities
│   │   ├── roast.go          # Roast entities
│   │   └── ...
│   ├── repository/           # Data persistence layer
│   │   ├── confession_repository.go
│   │   ├── roast_repository.go
│   │   └── repository.go
│   ├── usecase/              # Business logic layer
│   │   ├── chatbot/          # AI chatbot service
│   │   ├── confession/       # Confession management
│   │   ├── music/            # Music streaming
│   │   ├── news/             # News aggregation
│   │   ├── roast/            # Roast generation
│   │   └── whale/            # Whale alerts
│   ├── delivery/             # External interfaces
│   │   └── discord/          # Discord bot interface
│   │       ├── bot.go        # Bot initialization
│   │       └── handlers.go   # Command handlers
│   └── pkg/                  # Shared utilities
│       ├── ai/               # AI provider implementations
│       │   └── deepseek.go   # DeepSeek integration
│       ├── ffmpeg/           # FFmpeg wrapper
│       ├── logger/           # Logging utilities
│       └── ytdlp/            # yt-dlp wrapper
├── data/                     # Data storage (JSON files)
│   ├── confessions/          # Confession data
│   └── roasts/               # Roast data & patterns
├── deploy/                   # Deployment scripts
│   ├── setup.sh              # VPS setup script
│   ├── systemd/              # Systemd service files
│   └── docker/               # Docker configurations
├── docs/                     # Documentation
├── .env.example              # Environment template
├── docker-compose.yml        # Docker Compose config
├── Dockerfile                # Docker build file
├── Makefile                  # Build automation
└── go.mod                    # Go dependencies
```

### Layer Responsibilities

**1. Entity Layer** (`internal/entity/`)
- Pure business objects
- No external dependencies
- Defines core domain models

**2. Use Case Layer** (`internal/usecase/`)
- Business logic implementation
- Orchestrates data flow
- Independent of frameworks

**3. Repository Layer** (`internal/repository/`)
- Data persistence abstraction
- File/database operations
- Interface-based design

**4. Delivery Layer** (`internal/delivery/`)
- External interfaces (Discord, HTTP)
- Framework-specific code
- Converts external requests to use cases

**5. Infrastructure** (`internal/pkg/`)
- Shared utilities and tools
- External service wrappers
- Logging, AI providers, FFmpeg

### Design Principles

- ✅ **Dependency Inversion** - High-level modules don't depend on low-level modules
- ✅ **Single Responsibility** - Each module has one reason to change
- ✅ **Interface Segregation** - Clients depend on interfaces they use
- ✅ **Separation of Concerns** - Clear boundaries between layers
- ✅ **Testability** - Easy to test each component independently

For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## 🚀 Deployment

### Local Development

```bash
# Run directly
go run cmd/nerubot/main.go

# Or use Makefile
make run
```

### Docker Deployment

```bash
# Build image
docker build -t nerubot:latest .

# Run container
docker run -d \
  --name nerubot \
  --env-file .env \
  -v $(pwd)/data:/app/data \
  nerubot:latest

# View logs
docker logs -f nerubot
```

### VPS Deployment (Ubuntu/Debian)

```bash
# One-command setup
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/setup.sh | sudo bash

# Manual setup
sudo ./deploy/setup.sh

# Check status
sudo systemctl status nerubot

# View logs
sudo journalctl -u nerubot -f
```

### Production Checklist

- [ ] Set `LOG_LEVEL=INFO` or `WARNING`
- [ ] Configure proper `BOT_STATUS`
- [ ] Enable only required features
- [ ] Set up monitoring and alerts
- [ ] Configure log rotation
- [ ] Regular backups of `data/` directory
- [ ] Use strong Discord bot token
- [ ] Restrict file permissions (chmod 600 .env)

For detailed deployment instructions, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

---

## 📖 Documentation

### Available Documentation

- 📘 [Architecture Guide](docs/ARCHITECTURE.md) - System design and structure
- 🚀 [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment instructions
- 🔧 [Project Structure](docs/PROJECT_STRUCTURE.md) - Detailed file organization
- 🤝 [Contributing Guide](CONTRIBUTING.md) - How to contribute
- 📝 [Changelog](CHANGELOG.md) - Version history

### Additional Resources

- [Discord.js Guide](https://discordjs.guide/) - Discord bot development
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Architecture principles
- [Go Best Practices](https://golang.org/doc/effective_go.html) - Go programming guide

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'feat: Add amazing feature'`)
4. **Push to branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

**Examples:**
```
feat: Add playlist support to music system
fix: Resolve queue management race condition
docs: Update installation instructions
```

For more details, see [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📊 Project Status

### Current Version: 3.0.0

**Completed Features:**
- ✅ Music System (YouTube streaming)
- ✅ Confession System (Anonymous submissions)
- ✅ Roast System (Activity tracking & generation)
- ✅ Slash Commands (Modern Discord interface)
- ✅ Clean Architecture Implementation
- ✅ Docker Support

**In Development:**
- 🚧 AI Chatbot (DeepSeek integration)
- 🚧 News Aggregation System
- 🚧 Crypto Whale Alerts

**Planned:**
- 📋 Web Dashboard
- 📋 Database Migration (JSON → PostgreSQL)
- 📋 Microservices Architecture
- 📋 Advanced Analytics

---

## 🐛 Known Issues

- Music playback may have occasional buffering on slow connections
- Large confession images may take longer to process
- Roast cooldown is per-guild, not global

Report issues at: [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)

---

## 📜 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2024 nerufuyo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 🙏 Acknowledgments

- **[@bwmarrin](https://github.com/bwmarrin)** - DiscordGo library
- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** - YouTube download tool
- **[FFmpeg](https://ffmpeg.org/)** - Audio processing
- **Discord Community** - Feedback and support

---

## 📞 Support

### Get Help

- 📖 **Documentation:** Check [docs/](docs/) directory
- 💬 **Discord Server:** [Join our community](#) (Coming soon)
- 🐛 **Bug Reports:** [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)
- ✨ **Feature Requests:** [GitHub Discussions](https://github.com/nerufuyo/nerubot/discussions)

### Contact

- **Author:** [@nerufuyo](https://github.com/nerufuyo)
- **Email:** [Create an issue for contact](https://github.com/nerufuyo/nerubot/issues/new)
- **Website:** [Coming soon](#)

---

<div align="center">

**Made with ❤️ by [@nerufuyo](https://github.com/nerufuyo)**

⭐ Star this repository if you find it helpful!

[Report Bug](https://github.com/nerufuyo/nerubot/issues) · [Request Feature](https://github.com/nerufuyo/nerubot/issues) · [Documentation](docs/)

</div>

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