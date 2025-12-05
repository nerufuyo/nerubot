# NeruBot - Discord Microservices Bot 🎵

<div align="center">

![NeruBot Banner](https://imgur.com/yh3j7PK.png)

[![Discord Bot](https://img.shields.io/badge/Discord-Bot-7289da?style=for-the-badge&logo=discord&logoColor=white)](https://discord.com/oauth2/authorize?client_id=yourid&permissions=8&scope=bot%20applications.commands)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-3.0.0-blue?style=for-the-badge)](CHANGELOG.md)

**A powerful, scalable Discord bot built with Go microservices architecture - bringing music, community engagement, and entertainment to your server**

[🚀 Quick Start](#-quick-start) • [✨ Features](#-features) • [📖 Documentation](#-documentation) • [🏗️ Architecture](#️-architecture)

</div>

---

## 🎯 About NeruBot

NeruBot is a comprehensive Discord companion created by **[@nerufuyo](https://github.com/nerufuyo)** that transforms your server into an interactive entertainment hub. Built with **Go microservices architecture** for superior performance, scalability, and reliability.

### 🏆 Why Choose NeruBot?

- **⚡ Lightning Fast** - Built with Go for exceptional performance
- **🏗️ Microservices Architecture** - Scalable, maintainable, independent services
- **🎵 Premium Audio** - Crystal-clear YouTube streaming with yt-dlp
- **🛡️ Privacy-First** - Anonymous features with robust security
- **🔒 Production Ready** - Enterprise-grade with Railway deployment
- **💰 Completely Free** - No premium tiers, everything included!

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎵 **Music Service**
- **YouTube Support** - High-quality audio streaming
- **Queue Management** - Add, skip, stop, shuffle
- **Loop Modes** - None, single, or entire queue
- **Voice Detection** - Auto-validation
- **Rich Embeds** - Beautiful displays
- **Thread-Safe** - Concurrent operations

### 📝 **Confession Service**
- **Complete Anonymity** - Secure, private
- **Image Support** - Attach images
- **Reply System** - Community engagement
- **Moderation Queue** - Admin approval
- **Per-Server Settings** - Customizable
- **Thread-Safe Storage** - PostgreSQL

### 🔥 **Roast Service**
- **Activity Tracking** - Messages, reactions, voice
- **Smart Patterns** - 8 roast categories
- **Profile Analysis** - Behavior insights
- **Statistics** - Comprehensive metrics
- **Safety Systems** - Cooldowns & friendly
- **Persistent Data** - Database storage

</td>
<td width="50%">

### 🤖 **Chatbot Service** (Coming Soon)
- **Multi-Provider AI** - OpenAI, Claude, Gemini
- **Auto Fallback** - Provider redundancy
- **Session Management** - 30-min timeout
- **Context-Aware** - Smart conversations
- **Token Tracking** - Usage monitoring

### 📰 **News Service** (Coming Soon)
- **RSS Aggregation** - Multiple sources
- **Concurrent Fetching** - Fast updates
- **Auto-Publishing** - Scheduled posts
- **Custom Sources** - Configurable feeds
- **Deduplication** - No repeats

### 🐋 **Whale Service** (Coming Soon)
- **Crypto Monitoring** - Whale transactions
- **Real-Time Alerts** - Instant notifications
- **Multi-Blockchain** - Multiple networks
- **Configurable Thresholds** - Custom amounts
- **Transaction History** - Complete logs

</td>
</tr>
</table>

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│             Discord API                      │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│          API Gateway (8080)                  │
│  - Discord Bot Handler                       │
│  - Command Router                            │
│  - Response Aggregator                       │
└────┬────┬────┬────┬────┬────┬──────────────┘
     │    │    │    │    │    │
     ▼    ▼    ▼    ▼    ▼    ▼
   Music Conf Roast Chat News Whale
   :8081 :8082 :8083 :8084 :8085 :8086
     │    │    │    │    │    │
     └────┴────┴────┴────┴────┘
                   │
                   ▼
          ┌───────────────┐
          │  PostgreSQL   │
          │     Redis     │
          └───────────────┘
```

### Key Benefits

✅ **Independent Scaling** - Scale services independently  
✅ **Fault Isolation** - Service failures don't affect others  
✅ **Technology Flexibility** - Use best tool for each job  
✅ **Easier Maintenance** - Smaller, focused codebases  
✅ **Parallel Development** - Teams work independently  
✅ **Simplified Deployment** - Deploy services separately  

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Docker** - [Download](https://www.docker.com/get-started)
- **Discord Bot Token** - [Create Bot](https://discord.com/developers/applications)

### Option 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Copy environment template
cp .env.example .env

# Edit .env with your Discord token
nano .env

# Start all services
docker-compose -f docker-compose.microservices.yml up -d

# View logs
docker-compose -f docker-compose.microservices.yml logs -f
```

### Option 2: Local Development

```bash
# Clone repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Install dependencies
make -f Makefile.microservices deps

# Generate gRPC code
make -f Makefile.microservices proto

# Build services
make -f Makefile.microservices build

# Start infrastructure
docker-compose -f docker-compose.microservices.yml up -d postgres redis

# Run API Gateway
./build/gateway/gateway
```

### Option 3: Railway Deployment (Production)

See [Railway Deployment Guide](docs/RAILWAY_DEPLOYMENT.md) for complete instructions.

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Deploy
railway up
```

---

## 📖 Documentation

- **[Development Guide](docs/DEVELOPMENT.md)** - Setup, building, testing
- **[Railway Deployment](docs/RAILWAY_DEPLOYMENT.md)** - Production deployment
- **[Architecture](ARCHITECTURE.md)** - Clean architecture design
- **[Project Brief](projects/project_brief.md)** - Project overview
- **[Migration Plan](projects/project_plan.md)** - Microservices migration
- **[Contributing](CONTRIBUTING.md)** - How to contribute

---

## 🎮 Commands

### Music Commands
| Command | Description |
|---------|-------------|
| `/play <song>` | Play music from YouTube |
| `/pause` | Pause current song |
| `/resume` | Resume playback |
| `/skip` | Skip to next song |
| `/stop` | Stop and clear queue |
| `/queue` | Show current queue |
| `/nowplaying` | Show playing song info |

### Confession Commands
| Command | Description |
|---------|-------------|
| `/confess <message>` | Submit anonymous confession |

### Roast Commands
| Command | Description |
|---------|-------------|
| `/roast [user]` | Get roasted! |
| `/profile [user]` | View activity profile |

### General Commands
| Command | Description |
|---------|-------------|
| `/ping` | Check bot status |

---

## 🛠️ Development

### Project Structure

```
nerubot/
├── api/proto/          # gRPC protocol definitions
├── services/           # Microservices
│   ├── gateway/       # API Gateway
│   ├── music/         # Music Service
│   ├── confession/    # Confession Service
│   ├── roast/         # Roast Service
│   ├── chatbot/       # Chatbot Service
│   ├── news/          # News Service
│   └── whale/         # Whale Service
├── internal/          # Shared packages
│   ├── config/       # Configuration
│   ├── entity/       # Domain models
│   ├── pkg/          # Utilities
│   └── repository/   # Data access
├── docs/             # Documentation
└── scripts/          # Utility scripts
```

### Build Commands

```bash
# Build all services
make -f Makefile.microservices build

# Run tests
make -f Makefile.microservices test

# Start with Docker
make -f Makefile.microservices docker-up

# View logs
make -f Makefile.microservices docker-logs
```

See [Development Guide](docs/DEVELOPMENT.md) for detailed instructions.

---

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📊 Status

| Service | Status | Version | Railway |
|---------|--------|---------|---------|
| API Gateway | ✅ Ready | 1.0.0 | [Deploy](docs/RAILWAY_DEPLOYMENT.md) |
| Music | 🚧 In Progress | 0.9.0 | Pending |
| Confession | 🚧 In Progress | 0.9.0 | Pending |
| Roast | 🚧 In Progress | 0.9.0 | Pending |
| Chatbot | 📋 Planned | 0.0.0 | Planned |
| News | 📋 Planned | 0.0.0 | Planned |
| Whale | 📋 Planned | 0.0.0 | Planned |

---

## 📈 Roadmap

### Phase 1: Core Services (Current)
- [x] Project planning and documentation
- [x] gRPC protocol definitions
- [x] API Gateway service
- [ ] Music service migration
- [ ] Confession service migration
- [ ] Roast service migration

### Phase 2: Enhanced Features
- [ ] Chatbot service with AI providers
- [ ] News aggregation service
- [ ] Whale alert service
- [ ] Advanced monitoring & logging

### Phase 3: Scaling & Optimization
- [ ] Load balancing
- [ ] Service mesh
- [ ] Advanced caching strategies
- [ ] Performance optimization

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **[DiscordGo](https://github.com/bwmarrin/discordgo)** - Discord API wrapper
- **[gRPC](https://grpc.io/)** - High-performance RPC framework
- **[Railway](https://railway.app/)** - Deployment platform
- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** - YouTube downloader

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)
- **Discussions:** [GitHub Discussions](https://github.com/nerufuyo/nerubot/discussions)
- **Email:** nerufuyo@example.com

---

<div align="center">

**Made with ❤️ by [@nerufuyo](https://github.com/nerufuyo)**

⭐ Star this repository if you find it helpful!

</div>
