# Changelog

All notable changes to NeruBot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-10

### Added
- **Music System**
  - YouTube audio streaming with yt-dlp
  - Queue management (add, skip, pause, resume, stop)
  - Now playing information
  - FFmpeg audio processing
  - DCA encoding for Discord voice

- **Confession System**
  - Anonymous confession submission
  - Moderation queue with approval/rejection
  - Reply system with threading
  - Guild-specific settings
  - JSON-based persistence

- **Roast System**
  - User activity tracking (messages, voice time, commands)
  - AI-powered roast generation
  - Statistics and leaderboards
  - 8 roast categories (night owl, spammer, lurker, etc.)
  - Pattern-based roast templates

- **AI Chatbot**
  - Multi-provider support (Claude, Gemini, OpenAI)
  - Automatic fallback between providers
  - Session management with 30-minute timeout
  - Background session cleanup
  - Context-aware conversations

- **News System**
  - RSS feed aggregation from 5 sources
  - Concurrent fetching with goroutines
  - Auto-publishing capability
  - Customizable news sources

- **Whale Alerts**
  - Cryptocurrency whale transaction monitoring
  - Whale Alert API integration
  - Configurable minimum threshold ($1M default)
  - Real-time alert capability

- **Core Infrastructure**
  - Clean Architecture implementation (5 layers)
  - Configuration management with environment variables
  - Structured logging with file rotation
  - FFmpeg wrapper for audio processing
  - yt-dlp wrapper for YouTube downloads
  - Discord bot integration with DiscordGo

- **Commands (19 total)**
  - Music: `/play`, `/skip`, `/pause`, `/resume`, `/stop`, `/queue`, `/nowplaying`
  - Confession: `/confess`, `/confess-reply`, `/confess-approve`, `/confess-reject`
  - Roast: `/roast`, `/roast-stats`, `/roast-leaderboard`
  - Chatbot: `/chat`, `/chat-reset`
  - News: `/news`
  - Whale: `/whale`
  - Utility: `/help`

### Technical Details
- Golang 1.21+ with Clean Architecture
- DiscordGo for Discord API integration
- JSON-based data persistence
- Multi-provider AI integration
- Concurrent operations with goroutines
- Thread-safe data access with sync.RWMutex
- Docker deployment support
- Systemd service configuration

### Dependencies
- `github.com/bwmarrin/discordgo` v0.28.1 - Discord API
- `github.com/joho/godotenv` v1.5.1 - Environment variables
- `github.com/mmcdole/gofeed` v1.3.0 - RSS parsing
- External: FFmpeg, yt-dlp

### Performance
- Binary size: ~8.8MB (optimized)
- Memory usage: ~50-100MB
- Startup time: <2 seconds
- Audio latency: <100ms

---

## Release Notes

### v1.0.0 - Initial Release
This is the first production release of NeruBot, a complete rewrite in Golang with Clean Architecture principles. All core features and optional features are fully implemented and production-ready.

**Highlights:**
- 6 major feature systems
- 19 slash commands
- Multi-provider AI chatbot
- Real-time crypto whale alerts
- High-quality audio streaming
- Anonymous confession system
- AI-powered user roasting

**Migration Notes:**
This version represents a complete migration from Python (discord.py) to Golang (DiscordGo). All features from the Python version have been reimplemented with improved performance and architecture.

---

For detailed documentation, see [README.md](README.md) and [ARCHITECTURE.md](ARCHITECTURE.md).
