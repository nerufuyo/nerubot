<div align="center">

# NeruBot

### A Production-Ready Discord Bot Built with Go

[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-3.0.0-blue?style=for-the-badge)](CHANGELOG.md)

**A powerful, feature-rich Discord bot featuring music streaming, anonymous confessions, and community engagement tools**

[Quick Start](#quick-start) • [Features](#features) • [Documentation](#documentation) • [Contributing](#contributing)

</div>

---

## Table of Contents

- [About](#about-nerubot)
- [Features](#features)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Deployment](#deployment)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

---

## About NeruBot

NeruBot is a comprehensive Discord bot created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. Built with Go for superior performance and reliability, NeruBot follows **Clean Architecture** principles for maintainability and scalability.

### Why Choose NeruBot?

- **Lightning Fast** - Built with Go for exceptional performance
- **Premium Audio** - Crystal-clear YouTube streaming via yt-dlp
- **Privacy-First** - Anonymous features with robust security
- **Clean Architecture** - Maintainable, scalable codebase
- **Production Ready** - Thread-safe operations and error handling
- **Completely Free** - No premium features, everything included

---

## Features

<table>
<tr>
<td width="50%">

### Music System
- YouTube audio streaming (yt-dlp)
- Queue management & controls
- Loop modes (none/single/queue)
- Now playing with rich embeds
- Voice state detection
- Thread-safe operations

### Confession System
- Complete anonymity
- Image attachment support
- Moderation queue
- Reply system
- Per-guild settings
- Confession numbering

### Roast System
- Activity tracking
- Smart pattern detection
- Profile analysis
- Leaderboards & stats
- Cooldown management
- 8 roast categories

</td>
<td width="50%">

### AI Chatbot (Neru)
- DeepSeek integration
- Friendly personality
- Context-aware conversations
- Session management
- Natural, helpful responses
- Chat history tracking

### News Aggregation
- RSS feed integration
- Multiple news sources
- Latest headlines
- Rich embed formatting
- TechCrunch, The Verge, CNN, etc.

### Crypto Whale Alerts
- Real-time transaction tracking
- Large transaction monitoring
- Blockchain analysis
- USD value formatting
- Multiple chain support

### Analytics System
- Server statistics tracking
- User activity profiles
- Command usage analytics
- Top users and commands
- Historical data

</td>
</tr>
</table>

---

## Quick Start

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

## Commands

### Music Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/play <song>` | Play music from YouTube | `/play never gonna give you up` |
| `/skip` | Skip to next song | `/skip` |
| `/stop` | Stop playback and clear queue | `/stop` |
| `/pause` | Pause current playback | `/pause` |
| `/resume` | Resume playback | `/resume` |
| `/queue` | Display current queue | `/queue` |
| `/nowplaying` | Show current song info | `/nowplaying` |

### Confession Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/confess` | Submit anonymous confession | Opens modal |
| `/confess-approve <id>` | Approve confession (Admin) | `/confess-approve 5` |
| `/confess-reject <id>` | Reject confession (Admin) | `/confess-reject 3` |
| `/confess-reply <id>` | Reply to confession (Admin) | Opens modal |

### Roast Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/roast [@user]` | Generate personalized roast | `/roast @username` |

### AI Chatbot Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/chat <message>` | Chat with Neru (AI assistant) | `/chat what's your favorite music?` |
| `/chat-reset` | Clear your chat history | `/chat-reset` |

### News & Crypto Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/news` | Get latest news headlines | `/news` |
| `/whale` | View recent whale transactions | `/whale` |

### Analytics Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/stats` | View server statistics | `/stats` |
| `/profile [@user]` | View user activity profile | `/profile @username` |

### Utility Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/ping` | Check bot response time | `/ping` |
| `/help` | Display help information | `/help` |
| `/about` | Show bot information | `/about` |

---

## Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# === REQUIRED SETTINGS ===
# Discord Bot Token (Get from: https://discord.com/developers/applications)
DISCORD_TOKEN=your_discord_bot_token_here

# === AI CHATBOT SETTINGS ===
# DeepSeek API Key (Get from: https://platform.deepseek.com/)
# Required for /chat commands with Neru AI personality
DEEPSEEK_API_KEY=your_deepseek_api_key_here

# === CRYPTO WHALE ALERTS ===
# Whale Alert API Key (Get from: https://whale-alert.io/)
# Required for /whale commands to track large crypto transactions
WHALE_ALERT_API_KEY=your_whale_alert_api_key

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
ENABLE_CONFESSION=true
ENABLE_ROAST=true

# === MUSIC SETTINGS ===
# Maximum songs in queue per server
MAX_QUEUE_SIZE=100

# Auto-disconnect timeout in seconds (0 = disabled)
AUTO_DISCONNECT_TIME=300

# === ADVANCED SETTINGS ===
# Bot activity status
BOT_STATUS=Your friendly Discord companion

# Database settings (if using database features)
DATABASE_URL=mongodb://localhost:27017

# Redis settings (if using caching)
REDIS_URL=redis://localhost:6379
```

### Feature Flags

Features are automatically enabled/disabled based on API key availability:

- **Music System** - Controlled by `ENABLE_MUSIC` environment variable
- **AI Chatbot** - Auto-enabled when `DEEPSEEK_API_KEY` is configured
- **Confession System** - Controlled by `ENABLE_CONFESSION` (default: true)
- **Roast System** - Controlled by `ENABLE_ROAST` (default: true)
- **News Aggregation** - Always available (uses free RSS feeds)
- **Whale Alerts** - Auto-enabled when `WHALE_ALERT_API_KEY` is configured
- **Analytics** - Always enabled (no API key required)

---

## Architecture

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

- **Dependency Inversion** - High-level modules don't depend on low-level modules
- **Single Responsibility** - Each module has one reason to change
- **Interface Segregation** - Clients depend on interfaces they use
- **Separation of Concerns** - Clear boundaries between layers
- **Testability** - Easy to test each component independently

For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## Deployment

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

## Documentation

### Available Documentation

- [Architecture Guide](docs/ARCHITECTURE.md) - System design and structure
- [Deployment Guide](docs/DEPLOYMENT.md) - Production deployment instructions
- [Project Structure](docs/PROJECT_STRUCTURE.md) - Detailed file organization
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Changelog](CHANGELOG.md) - Version history

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

## Project Status

### Current Version: 3.0.0

**Completed Features:**
- Music System (YouTube streaming with yt-dlp)
- Confession System (Anonymous submissions with moderation)
- Roast System (Activity tracking & personalized generation)
- AI Chatbot (Neru personality with DeepSeek)
- News Aggregation (RSS feeds from multiple sources)
- Crypto Whale Alerts (Real-time transaction monitoring)
- Analytics System (Server & user statistics)
- Slash Commands (Modern Discord interface)
- Clean Architecture Implementation
- Docker Support

**Planned:**
- Web Dashboard
- Database Migration (JSON → PostgreSQL)
- Microservices Architecture
- Advanced visualization
- Mobile app companion

---

## Known Issues

- Music playback may have occasional buffering on slow connections
- Large confession images may take longer to process
- Roast cooldown is per-guild, not global

Report issues at: [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)

---

## License

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

## Support

### Get Help

- **Documentation:** Check [docs/](docs/) directory
- **Discord Server:** [Join our community](#) (Coming soon)
- **Bug Reports:** [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)
- **Feature Requests:** [GitHub Discussions](https://github.com/nerufuyo/nerubot/discussions)

### Contact

- **Author:** [@nerufuyo](https://github.com/nerufuyo)
- **Email:** [Create an issue for contact](https://github.com/nerufuyo/nerubot/issues/new)
- **Website:** [Coming soon](#)

---

<div align="center">

**Made by [@nerufuyo](https://github.com/nerufuyo)**

[Report Bug](https://github.com/nerufuyo/nerubot/issues) · [Request Feature](https://github.com/nerufuyo/nerubot/issues) · [Documentation](docs/)

</div>

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