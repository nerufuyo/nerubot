# NeruBot - Your Ultimate Discord Companion 🎵

<div align="center">

![NeruBot Banner](https://imgur.com/yh3j7PK.png)

[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/oauth2/authorize?client_id=yourid&permissions=8&scope=bot%20applications.commands)
[![Python](https://img.shields.io/badge/Python-3.8+-3776ab?style=for-the-badge&logo=python&logoColor=white)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-2.0.0-blue?style=for-the-badge)](CHANGELOG.md)

**A powerful, feature-rich Discord bot designed to bring music, community engagement, and entertainment to your server**

[🚀 Quick Start](#-quick-start) • [📋 Features](#-features) • [📖 Documentation](#-documentation) • [🤝 Support](#-support)

</div>

---

## 🎯 About NeruBot

NeruBot is a comprehensive Discord companion created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. With high-quality multi-source music streaming, anonymous confession systems, real-time news updates, whale alerts, and intuitive slash commands, NeruBot delivers a premium Discord experience.

### 🏆 Why Choose NeruBot?

- **🎵 Premium Audio Quality** - Crystal-clear streaming from YouTube, Spotify, and SoundCloud
- **🛡️ Privacy-First Design** - Anonymous features with robust security measures
- **⚡ Lightning Fast** - Optimized performance with minimal resource usage
- **🎨 Beautiful Interface** - Modern Discord UI with interactive components
- **🔒 Production Ready** - Enterprise-grade architecture and deployment tools
- **💰 Completely Free** - No premium features, everything included!

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎵 **Advanced Music System**
- **Multi-Platform Support** - YouTube, Spotify, SoundCloud
- **Smart Queue Management** - Loop modes, shuffle, auto-queue
- **High-Quality Audio** - Optimized streaming with minimal latency
- **24/7 Mode** - Continuous playbook in voice channels
- **Playlist Support** - Import and manage playlists seamlessly
- **Interactive Controls** - Volume, skip, pause, resume

### 🤖 **AI-Powered Chatbot**
- **Multi-AI Support** - Claude, Gemini, OpenAI with smart fallback
- **Unique Personality** - Fun, witty gaming/anime character
- **Smart Sessions** - Welcome messages & 5-min timeout thanks
- **Natural Conversations** - Responds to mentions and DMs
- **Global AI Service** - Available for all bot features

### 🔥 **User Roasting System**
- **Behavior Analysis** - AI-powered analysis of user Discord habits
- **Personalized Roasts** - Hilarious, custom roasts based on activity patterns
- **Activity Tracking** - Monitors messages, voice time, commands, and more
- **Smart Categories** - 8 different roast types (night owl, spammer, lurker, etc.)
- **Safety Systems** - Cooldowns and friendly community-appropriate content
- **Rich Statistics** - Detailed behavior insights and roasting analytics

</td>
<td width="50%">

### 📝 **Anonymous Confession System**
- **Complete Anonymity** - Secure, private confession sharing
- **Image Support** - Attach images to confessions and replies
- **Interactive Replies** - Anonymous community engagement
- **Smart Moderation** - Cooldown protection and content filtering
- **Server Isolation** - Confessions stay within your community
- **Beautiful UI** - Modern Discord modals and buttons

</td>
</tr>
<tr>
<td width="50%">

### 📰 **Real-Time News & Alerts**
- **Trusted Sources** - 12+ international and Indonesian news outlets
- **Crypto Intelligence** - Whale alerts and market updates
- **Guru Monitoring** - Track crypto influencer tweets with sentiment analysis
- **Auto-Publishing** - Scheduled updates every 10 minutes
- **Manual Control** - Start/stop updates with admin commands
- **Smart Formatting** - Clean, readable news presentation

</td>
<td width="50%">

### 🤖 **User-Friendly Interface**
- **Slash Commands** - Modern Discord command system
- **Interactive Help** - Paginated help with button navigation
- **Rich Embeds** - Beautiful, consistent message formatting
- **Error Handling** - Comprehensive error messages and recovery
- **Multi-Language Ready** - Architecture supports localization

</td>
</tr>
</table>

---

## 🚀 Quick Start

### 🖥️ Local Development

```bash
# Clone the repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Automated setup and run
./run.sh
```

The script will:
- ✅ Create virtual environment
- ✅ Install dependencies
- ✅ Generate `.env` template
- ✅ Guide you through Discord token setup
- ✅ Start the bot

### 🌐 VPS Deployment (Production)

```bash
# One-command VPS setup (Ubuntu/Debian)
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/setup.sh | sudo bash
```

**What this does:**
- 🔧 Installs Python 3, FFmpeg, and dependencies
- 👤 Creates secure `nerubot` user
- 🛡️ Configures firewall (SSH only)
- ⚙️ Sets up systemd service
- 📊 Enables health monitoring

---

## 📋 Command Reference

### 🎵 Music Commands
| Command | Description |
|---------|-------------|
| `/play <song>` | Play music from any supported platform |
| `/queue` | Display current music queue |
| `/skip` | Skip to the next song |
| `/pause` / `/resume` | Control playback |
| `/loop [mode]` | Set loop mode (off/single/queue) |
| `/247` | Toggle 24/7 continuous mode |
| `/volume <level>` | Adjust playback volume |
| `/nowplaying` | Show currently playing song |

### 📝 Confession Commands
| Command | Description |
|---------|-------------|
| `/confess [image]` | Submit anonymous confession |
| `/reply <id> [image]` | Reply to a confession anonymously |
| `/confession-setup <channel>` | Set confession channel (Admin) |
| `/confession-stats` | View confession statistics |

### � Roast Commands
| Command | Description |
|---------|-------------|
| `/roast [target] [custom]` | Generate personalized roast based on user behavior |
| `/roast-stats [user]` | View roasting statistics and insights |
| `/behavior-analysis [user]` | Detailed Discord behavior analysis |

### 🤖 Chatbot Commands
| Command | Description |
|---------|-------------|
| `/chat <message>` | Start a conversation with the AI |
| `/reset-chat` | Reset your conversation history |

### �📰 News & Crypto Commands
| Command | Description |
|---------|-------------|
| `/news latest [count]` | Get latest news updates |
| `/news sources` | List configured news sources |
| `/news set-channel <channel>` | Set news channel (Admin) |
| `/news start` / `/news stop` | Control auto-updates (Admin) |
| `/whale setup [channel]` | Enable whale alerts |
| `/whale recent [limit]` | Show recent whale transactions |
| `/guru setup [channel]` | Enable crypto guru tweets |
| `/guru accounts` | List monitored crypto influencers |

### ℹ️ Information Commands
| Command | Description |
|---------|-------------|
| `/help` | Interactive help system with navigation |
| `/about` | Bot information and creator details |
| `/features` | Showcase all available features |
| `/commands` | Quick command reference card |

---

## ⚙️ Configuration

### 🔑 Environment Setup

1. **Create `.env` file:**
```env
# Required
DISCORD_TOKEN=your_discord_bot_token_here

# Optional
LOG_LEVEL=INFO
COMMAND_PREFIX=!
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
```

2. **Get Discord Bot Token:**
   - Visit [Discord Developer Portal](https://discord.com/developers/applications)
   - Create new application → Bot → Copy token
   - Enable all necessary intents

3. **Bot Permissions:**
   - Send Messages
   - Embed Links
   - Read Message History
   - Connect to Voice
   - Speak in Voice
   - Use Slash Commands

### 🎛️ Advanced Configuration

See [`config/messages.py`](src/config/messages.py) for customizable:
- Bot responses and messages
- Help system content
- Error messages
- Feature descriptions

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

```
src/
├── main.py                 # Bot entry point
├── config/                 # Configuration and messages
├── core/                   # Shared utilities and helpers
├── features/               # Feature modules
│   ├── music/             # Music streaming system
│   ├── help/              # Interactive help system
│   ├── confession/        # Anonymous confession system
│   ├── news/              # News broadcasting system
│   └── whale_alerts/      # Crypto whale alerts
└── interfaces/            # Discord interface layer
```

**Key Principles:**
- 🏛️ **Modular Design** - Features can be easily added/removed
- 🧹 **Clean Code** - Well-structured, maintainable codebase
- 🔒 **Security First** - Production-grade security practices
- 📈 **Scalable** - Ready for high-traffic servers
- 🧪 **Testable** - Comprehensive testing infrastructure

---

## 📊 System Requirements

### Minimum Requirements
- **OS:** Ubuntu 20.04+ / Debian 11+ / Windows 10+ / macOS 10.15+
- **Python:** 3.8 or higher
- **RAM:** 1GB
- **Storage:** 5GB
- **Network:** Stable internet connection

### Recommended (VPS)
- **CPU:** 2+ cores
- **RAM:** 2GB+
- **Storage:** 10GB+
- **Bandwidth:** 1TB/month

### Dependencies
- **FFmpeg** - Audio processing
- **Git** - Version control
- **Discord.py 2.3+** - Discord API wrapper

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
5. **Test** thoroughly
6. **Submit** pull request

### Development Setup
```bash
# Setup development environment
./run.sh setup

# Install development dependencies
pip install -r requirements-dev.txt

# Run tests
python -m pytest

# Code formatting
black src/ && isort src/
```

**Contribution Areas:**
- 🎵 Music features and sources
- 🛡️ Security improvements
- 📱 UI/UX enhancements
- 📚 Documentation
- 🧪 Testing coverage
- 🌐 Internationalization

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