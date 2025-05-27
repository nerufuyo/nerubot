# NeruBot - Discord Music Bot

A fast, efficient Discord music bot with high-quality audio streaming and smart queue management.

## ✨ Features

- 🎵 **Multi-Source Music** - YouTube, Spotify, SoundCloud support
- 🔄 **Smart Queue Management** - Loop modes, shuffle, auto-queue
- 🎛️ **Intuitive Controls** - Play, pause, skip, volume control
- 📱 **Interactive Help** - Paginated menus with button navigation
- 🌐 **24/7 Operation** - Continuous playback mode
- ⚡ **High Performance** - Optimized for low latency and resource usage

## 🚀 Quick Start

### Local Development
```bash
# Clone and setup
git clone https://github.com/your-username/nerubot.git
cd nerubot
./run.sh  # Automated setup and run
```

### VPS Deployment
```bash
# One-command VPS setup (Ubuntu/Debian)
curl -fsSL https://raw.githubusercontent.com/your-username/nerubot/main/deploy/setup.sh | sudo bash
```

## 🎵 Commands

| Command | Description |
|---------|-------------|
| `/play <song>` | Play music from any supported source |
| `/queue` | Show current music queue |
| `/skip` | Skip to next song |
| `/pause` / `/resume` | Control playbook |
| `/loop [mode]` | Set loop mode (off/single/queue) |
| `/247` | Toggle 24/7 continuous mode |
| `/help` | Interactive help system |

## 🛠️ Configuration

Create `.env` file:
```env
DISCORD_TOKEN=your_bot_token_here
LOG_LEVEL=INFO
COMMAND_PREFIX=!
```

Get your Discord bot token: [Discord Developer Portal](https://discord.com/developers/applications)

## 🔧 Management Scripts

| Script | Purpose |
|--------|---------|
| `./run.sh` | Setup and start the bot |
| `./clean.sh` | Clean temporary files |
| `deploy/setup.sh` | VPS environment setup |
| `deploy/monitor.sh` | Health monitoring |
| `deploy/status.sh` | Service status dashboard |

## 📖 Documentation

- **[Deployment Guide](deploy/README.md)** - VPS setup and configuration
- **[Contributing](CONTRIBUTING.md)** - Development guidelines
- **[Architecture](ARCHITECTURE.md)** - Technical overview
- **[Changelog](CHANGELOG.md)** - Version history

## 🏗️ Architecture

```
src/
├── main.py              # Bot entry point
├── config/              # Settings and configuration
├── features/            # Feature modules (music, help, news)
│   ├── music/          # Music streaming functionality
│   ├── help/           # Interactive help system
│   └── news/           # News broadcasting
└── interfaces/         # Discord interface layer
```

## 📊 Requirements

- **Python 3.8+**
- **FFmpeg** (for audio processing)
- **Discord Bot Token**
- **2GB+ RAM** (recommended for VPS)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Need help?** Join our Discord server or open an issue on GitHub!