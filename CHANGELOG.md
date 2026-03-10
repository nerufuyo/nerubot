# Changelog

All notable changes to NeruBot will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.0.0] - 2026-03-10

### Added
- **Core Commands** — `/ping`, `/botinfo`, `/serverinfo`, `/userinfo`, `/avatar` for essential bot & server information
- **Moderation System** — `/kick`, `/ban`, `/timeout`, `/purge`, `/warn`, `/warnings`, `/clearwarnings` with proper permission checks
- **Fun Commands** — `/coinflip`, `/8ball` (magic 8-ball) for casual server entertainment
- **Utility Commands** — `/poll` (reaction-based polls up to 5 options), `/calc` (math calculator with 6 operations)
- **Warning System** — Persistent warning records stored in MongoDB with per-guild, per-user tracking
- **Security: Spam Detection** — Automatic spam detection (5+ messages in 5 seconds) with message deletion
- **Security: Malicious Link Filter** — Detects and removes phishing Discord gift links, IP loggers, and suspicious domains
- **Paginated Help Menu** — Interactive `/help` with Prev/Next navigation buttons across 6 categorized pages
- **Moderation Repository** — `moderation_repository.go` for warning persistence, `moderation.go` entity

### Changed
- **Version bump**: 4.0.0 → 5.0.0
- **Help command**: Completely rewritten with Discord button-based pagination (6 pages: Overview, Fun, Utility, Moderation, Scheduling, Music)
- **Help translations**: All 5 languages (EN, ID, JP, KR, ZH) updated with new command entries
- **Content filter**: Enhanced with regex-based malicious link patterns and spam tracking
- **Command routing**: Reorganized with clear sections (core, features, fun, utility, moderation, music)
- **README.md**: Updated with comprehensive command tables organized by category

## [4.0.0] - 2025-07-21

### Added
- **Indonesian Holiday Reminders** — automatic `@everyone` notifications at 07:00 WIB for national holidays (Tahun Baru, Idul Fitri, Kemerdekaan, Natal, Nyepi, Waisak, and more) covering 2025-2027
- **Ramadan Sahoor & Berbuka Reminders** — automatic `@everyone` notifications at sahoor and Maghrib times during Ramadan, Jakarta timezone, with warm Indonesian-language messages
- `/reminder` command to view upcoming holidays and today's Ramadan schedule
- `.env.example` with clean configuration template
- `.gitignore` (proper Go-focused ignore rules)
- `docs/USER_GUIDE.md` — complete user-facing documentation
- `docs/DEVELOPMENT.md` — developer setup and architecture guide

### Changed
- **Code Structure Overhaul**
  - Split monolithic `handlers.go` (732 lines) into feature-specific handler files:
    `handler_music.go`, `handler_confession.go`, `handler_roast.go`, `handler_chatbot.go`,
    `handler_news.go`, `handler_whale.go`, `handler_analytics.go`, `handler_reminder.go`
  - Moved shared response helpers (`respond`, `respondEmbed`, `followUp`, etc.) into `handlers.go`
  - Simplified `main.go` — removed unused imports and redundant service initialization
  - Fixed bot lifecycle: `Start()` is non-blocking; signal handling is in `main.go` only
- **Emoji Cleanup** — replaced all decorative Unicode emojis with text indicators (`[OK]`, `[ERR]`, `>>`, `||`, etc.)
- **Documentation** — rewrote README.md, added user guide and development guide

### Removed
- Unused gRPC microservices (`services/`, `api/proto/`)
- Lavalink Docker/Railway configs (`Dockerfile.lavalink`, `application.yml`, `railway.lavalink.toml`)
- Extra Docker Compose files (`docker-compose.microservices.yml`, `docker-compose.music.yml`)
- Proto generation scripts, init-db.sql, test-services.ps1
- Outdated documentation files (ARCHITECTURE.md, CONTRIBUTING.md, docs/*.md)
- gRPC/protobuf dependencies from go.mod

### Fixed
- **Logger `sprintf`**: Was returning the format string unformatted; now uses `fmt.Sprintf`
- **Music shuffle**: `shuffleSongs` used biased `time.Now().UnixNano() % n` — replaced with `math/rand`
- **Reuters RSS URL**: Fixed invalid URL (`reedsnews`) — pointed to a working Google News / Reuters feed
- **Sorting**: Replaced O(n²) bubble sort with `sort.Slice` in analytics and news services

### Updated
- All Go dependencies upgraded to latest versions
- Version bumped from 3.0.0 to 4.0.0

---

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
