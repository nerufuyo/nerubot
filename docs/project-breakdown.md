# NeruBot - Python to Golang Migration Project Breakdown

## Project Overview

**Project Name:** NeruBot Discord Bot Migration  
**Current Stack:** Python 3.8+ with discord.py  
**Target Stack:** Golang with DiscordGo  
**Migration Type:** Complete rewrite using Clean Architecture  
**Project Owner:** @nerufuyo  
**Created:** November 9, 2025

## Executive Summary

This document outlines the complete breakdown of migrating NeruBot, a feature-rich Discord bot, from Python to Golang. The migration will follow Clean Architecture principles with proper separation of concerns across Delivery, Use Case, Entity, and Repository layers.

## Current System Analysis

### Technology Stack (Python)
- **Language:** Python 3.8+
- **Discord Library:** discord.py 2.3.0+
- **Audio Processing:** FFmpeg, yt-dlp, PyNaCl
- **AI Services:** OpenAI, Anthropic Claude, Google Gemini
- **Data Persistence:** JSON files
- **Deployment:** Docker, systemd, nginx

### Project Structure
```
src/
├── main.py                    # Bot entry point
├── config/                    # Configuration and messages
│   ├── settings.py           # Bot settings, limits, features
│   └── messages.py           # User-facing messages
├── core/                      # Core utilities
│   ├── constants.py          # Shared constants
│   └── utils/                # Logging, file utilities
├── features/                  # Feature modules (8 features)
│   ├── music/                # Multi-source music streaming
│   ├── confession/           # Anonymous confession system
│   ├── chatbot/              # AI-powered chatbot
│   ├── roast/                # User roasting system
│   ├── news/                 # News broadcasting
│   ├── whale_alerts/         # Crypto whale alerts
│   └── help/                 # Help system (4 cogs)
└── interfaces/               # External interfaces
    └── discord/              # Discord bot implementation
```

### Features Inventory

#### 1. Music System (Priority: HIGH)
**Complexity:** High  
**Components:**
- Multi-source support (YouTube, Spotify, SoundCloud)
- Queue management with loop modes (off/single/queue)
- 24/7 continuous mode
- Voice channel management
- FFmpeg audio streaming
- Playlist support

**Key Services:**
- `MusicService` - Queue, playback control
- `SourceManager` - Multi-source search and extraction
- YouTube, Spotify, SoundCloud source handlers

**Models:**
- `Song` - Song metadata (title, artist, duration, URL, source)

**Data Storage:** In-memory queues, no persistence

#### 2. Confession System (Priority: HIGH)
**Complexity:** Medium  
**Components:**
- Anonymous confession submission
- Image attachment support
- Reply system with threading
- Queue-based processing
- Cooldown management
- Guild-specific settings

**Key Services:**
- `ConfessionService` - Confession management
- `QueueService` - Async queue processing

**Models:**
- `Confession` - ID, content, author, guild, status, timestamps
- `ConfessionReply` - Reply to confession with threading
- `GuildConfessionSettings` - Channel settings, cooldowns

**Data Storage:** JSON files (confessions.json, replies.json, settings.json, queue.json)

#### 3. AI Chatbot (Priority: MEDIUM)
**Complexity:** Medium  
**Components:**
- Multi-AI support (Claude, Gemini, OpenAI)
- Smart fallback system
- Session management with timeout
- Welcome and thank you messages
- Personality: Fun, witty gaming/anime character

**Key Services:**
- `ChatbotService` - Multi-provider AI integration
- Session tracking and timeout management

**Models:**
- `ChatSession` - User sessions with message history

**Data Storage:** In-memory sessions

**External APIs:**
- Anthropic Claude API
- Google Gemini API
- OpenAI API

#### 4. Roast System (Priority: MEDIUM)
**Complexity:** Medium-High  
**Components:**
- User behavior tracking (messages, voice time, commands)
- AI-powered roast generation
- 8 roast categories (night owl, spammer, lurker, etc.)
- Activity pattern analysis
- Statistics and insights
- Cooldown system

**Key Services:**
- `RoastService` - Roast generation and behavior analysis
- Activity tracking across Discord events

**Models:**
- `UserProfile` - User activity data
- `RoastPattern` - Roast templates and categories
- `ActivityStats` - Detailed statistics

**Data Storage:** JSON files (profiles.json, patterns.json, stats.json, activities.json)

#### 5. News System (Priority: LOW)
**Complexity:** Medium  
**Components:**
- Multi-source news aggregation (12+ sources)
- Auto-publishing every 10 minutes
- Manual start/stop controls
- Rich formatting
- Source management

**Key Services:**
- `NewsService` - News fetching and publishing
- RSS feed parsing
- Background task scheduler

**Models:**
- `NewsItem` - Title, description, link, source, published date

**Data Storage:** In-memory, no persistence

**External APIs:**
- RSS feeds from 12+ news sources

#### 6. Whale Alerts (Priority: LOW)
**Complexity:** Medium  
**Components:**
- Crypto whale transaction monitoring
- Guru tweet tracking with sentiment
- Real-time alerts
- Channel configuration

**Key Services:**
- `WhaleAlertService` - Transaction monitoring
- `GuruService` - Tweet monitoring

**Models:**
- `WhaleTransaction` - Amount, blockchain, timestamp
- `GuruTweet` - Tweet content, sentiment, author

**Data Storage:** In-memory

**External APIs:**
- Whale Alert API
- Twitter/X API

#### 7. Help System (Priority: HIGH)
**Complexity:** Low  
**Components:**
- Interactive paginated help
- Feature showcase
- Command reference
- About information
- Button navigation

**Key Services:**
- Help embed generation
- Navigation handling

**Models:**
- Help content from config

**Data Storage:** Static configuration

#### 8. General Bot Features (Priority: HIGH)
**Complexity:** Low  
**Components:**
- Slash command system
- Error handling
- Logging system
- Status/presence management
- Command sync

## Target Architecture (Golang)

### Technology Stack
- **Language:** Go 1.21+
- **Discord Library:** DiscordGo
- **Audio Processing:** FFmpeg (exec), yt-dlp (exec)
- **AI Services:** Native HTTP clients for APIs
- **Data Persistence:** JSON files with encoding/json
- **HTTP Client:** net/http, context for timeouts
- **Concurrency:** Goroutines, channels
- **Build:** Go modules
- **Deployment:** Docker multi-stage, systemd

### Clean Architecture Layers

```
nerubot-go/
├── cmd/
│   └── nerubot/              # Application entry point
│       └── main.go
├── internal/
│   ├── config/               # Configuration (Model layer)
│   │   ├── config.go        # Bot settings
│   │   ├── messages.go      # User messages
│   │   └── constants.go     # Constants
│   ├── entity/               # Domain entities (Model layer)
│   │   ├── song.go
│   │   ├── confession.go
│   │   ├── news.go
│   │   ├── roast.go
│   │   └── chat.go
│   ├── usecase/              # Business logic (Use Case layer)
│   │   ├── music/
│   │   │   ├── music_service.go
│   │   │   └── source_manager.go
│   │   ├── confession/
│   │   │   ├── confession_service.go
│   │   │   └── queue_service.go
│   │   ├── chatbot/
│   │   │   └── chatbot_service.go
│   │   ├── roast/
│   │   │   └── roast_service.go
│   │   ├── news/
│   │   │   └── news_service.go
│   │   └── whale/
│   │       └── whale_service.go
│   ├── repository/           # Data access (Repository layer)
│   │   ├── confession_repo.go
│   │   ├── roast_repo.go
│   │   └── settings_repo.go
│   ├── delivery/             # External interfaces (Delivery/Gateway layer)
│   │   └── discord/
│   │       ├── bot.go       # Bot setup
│   │       ├── handlers/    # Command handlers
│   │       │   ├── music_handler.go
│   │       │   ├── confession_handler.go
│   │       │   ├── chatbot_handler.go
│   │       │   ├── roast_handler.go
│   │       │   ├── news_handler.go
│   │       │   ├── whale_handler.go
│   │       │   └── help_handler.go
│   │       └── events/      # Event listeners
│   │           └── event_handlers.go
│   └── pkg/                  # Shared packages
│       ├── logger/          # Logging utilities
│       ├── ffmpeg/          # FFmpeg wrapper
│       ├── ytdlp/           # yt-dlp wrapper
│       └── ai/              # AI client interfaces
│           ├── claude.go
│           ├── gemini.go
│           └── openai.go
├── data/                     # Data storage (same as Python)
│   ├── confessions/
│   ├── roasts/
│   └── settings/
├── deploy/                   # Deployment configs
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### Dependency Flow (Clean Architecture)
```
Delivery (Discord Handlers)
    ↓ depends on
Use Case (Services)
    ↓ depends on
Entity (Domain Models)
    ↑ implements
Repository (Data Access)
```

**Key Principles:**
1. **Dependency Rule:** Inner layers never depend on outer layers
2. **Interface Segregation:** Use cases define repository interfaces
3. **Single Responsibility:** Each layer has one clear purpose
4. **KISS:** Keep implementations simple and straightforward

## Migration Breakdown by Component

### 1. Core Infrastructure

#### 1.1 Configuration System
**Python:**
- `config/settings.py` - Bot settings, limits, feature flags
- `config/messages.py` - User-facing messages
- `core/constants.py` - Shared constants

**Golang:**
- `internal/config/config.go` - Struct-based configuration with defaults
- `internal/config/messages.go` - Message templates
- `internal/config/constants.go` - Constant definitions
- Load from environment variables and `.env` file using `godotenv`

**Changes:**
- Use structs instead of dictionaries
- Implement config validation
- Support YAML/JSON config files (optional)

#### 1.2 Logging System
**Python:**
- `core/utils/logging_utils.py` - Custom logger setup
- Rotating file handler

**Golang:**
- `internal/pkg/logger/logger.go` - Structured logging with levels
- Use `log/slog` (Go 1.21+) or `logrus`
- File rotation with `lumberjack`

**Changes:**
- Structured logging with fields
- Context-aware logging
- Better performance

#### 1.3 File Utilities
**Python:**
- `core/utils/file_utils.py` - FFmpeg detection, file operations

**Golang:**
- `internal/pkg/ffmpeg/ffmpeg.go` - FFmpeg wrapper
- `internal/pkg/ytdlp/ytdlp.go` - yt-dlp wrapper
- Use `os/exec` for external commands

**Changes:**
- Better error handling
- Context-based timeouts
- Process management

### 2. Discord Bot Layer (Delivery)

#### 2.1 Bot Setup
**Python:**
- `interfaces/discord/bot.py` - Bot initialization, cog loading
- `main.py` - Entry point

**Golang:**
- `cmd/nerubot/main.go` - Application entry
- `internal/delivery/discord/bot.go` - Bot setup and lifecycle
- Use DiscordGo session management

**Changes:**
- Explicit dependency injection
- Graceful shutdown with context
- Health check endpoint (optional)

#### 2.2 Command Handlers
**Python:**
- Cogs with `@app_commands.command` decorators
- Automatic command registration

**Golang:**
- Handler functions in `internal/delivery/discord/handlers/`
- Manual command registration with DiscordGo
- Slash command support with `ApplicationCommand`

**Changes:**
- Struct-based handlers with dependencies injected
- Middleware pattern for common functionality
- Better type safety

#### 2.3 Event Handling
**Python:**
- `@commands.Cog.listener()` decorators
- Event handlers in cogs

**Golang:**
- `internal/delivery/discord/events/event_handlers.go`
- DiscordGo event handlers with callbacks

**Changes:**
- Centralized event routing
- Goroutines for concurrent event processing
- Better error recovery

### 3. Feature Migration

#### 3.1 Music System
**Complexity:** HIGH  
**Effort:** 40 hours

**Components to migrate:**
1. **Music Service** (`usecase/music/music_service.go`)
   - Queue management with `container/list` or custom queue
   - Playback control with goroutines
   - Loop mode handling
   - 24/7 mode state management

2. **Source Manager** (`usecase/music/source_manager.go`)
   - YouTube source with yt-dlp
   - Spotify API integration with YouTube fallback
   - SoundCloud support
   - Concurrent search with `errgroup`

3. **Voice Handler**
   - Voice channel connection with DiscordGo
   - Audio streaming with DCA (Discord Compatible Audio)
   - FFmpeg piping
   - Voice state management

4. **Models** (`entity/song.go`)
   - Song struct with metadata
   - Queue item struct

**Key Challenges:**
- Audio encoding to DCA format
- Voice connection management in DiscordGo
- Concurrent playback handling
- FFmpeg process management

**Go Libraries:**
- `github.com/bwmarrin/discordgo` - Discord API
- `github.com/bwmarrin/dca` - Audio encoding
- `os/exec` - FFmpeg/yt-dlp execution
- `zmb3/spotify` - Spotify API

#### 3.2 Confession System
**Complexity:** MEDIUM  
**Effort:** 24 hours

**Components to migrate:**
1. **Confession Service** (`usecase/confession/confession_service.go`)
   - Confession creation with ID generation
   - Reply handling
   - Status management
   - Cooldown tracking

2. **Queue Service** (`usecase/confession/queue_service.go`)
   - Async queue with channels
   - Worker pool pattern
   - Queue persistence

3. **Repository** (`repository/confession_repo.go`)
   - JSON file operations with `encoding/json`
   - Concurrent access with `sync.RWMutex`
   - Atomic writes

4. **Models** (`entity/confession.go`)
   - Confession struct
   - Reply struct
   - Settings struct
   - Status enum

**Key Challenges:**
- Channel-based queue processing
- Thread-safe data access
- JSON marshaling/unmarshaling with custom types

**Go Libraries:**
- `encoding/json` - JSON operations
- `sync` - Mutexes for thread safety

#### 3.3 AI Chatbot
**Complexity:** MEDIUM  
**Effort:** 20 hours

**Components to migrate:**
1. **Chatbot Service** (`usecase/chatbot/chatbot_service.go`)
   - Multi-provider AI client
   - Fallback logic
   - Session management
   - Timeout handling

2. **AI Clients** (`internal/pkg/ai/`)
   - `claude.go` - Anthropic Claude client
   - `gemini.go` - Google Gemini client
   - `openai.go` - OpenAI client
   - Common interface for all providers

3. **Session Management**
   - In-memory sessions with `sync.Map`
   - TTL-based expiration
   - Message history tracking

4. **Models** (`entity/chat.go`)
   - ChatSession struct
   - Message struct

**Key Challenges:**
- HTTP client with proper timeout/retry
- Context propagation
- Efficient session cleanup

**Go Libraries:**
- `net/http` - HTTP clients
- Standard library for JSON

#### 3.4 Roast System
**Complexity:** MEDIUM-HIGH  
**Effort:** 28 hours

**Components to migrate:**
1. **Roast Service** (`usecase/roast/roast_service.go`)
   - Behavior analysis
   - Roast generation with AI
   - Pattern matching
   - Statistics calculation

2. **Activity Tracker**
   - Event listeners for Discord activities
   - Concurrent activity recording
   - Aggregation logic

3. **Repository** (`repository/roast_repo.go`)
   - User profile persistence
   - Activity log storage
   - Statistics tracking

4. **Models** (`entity/roast.go`)
   - UserProfile struct
   - RoastPattern struct
   - ActivityStats struct

**Key Challenges:**
- Concurrent activity tracking
- Efficient data aggregation
- AI integration for roast generation

#### 3.5 News System
**Complexity:** MEDIUM  
**Effort:** 16 hours

**Components to migrate:**
1. **News Service** (`usecase/news/news_service.go`)
   - RSS feed fetching
   - News aggregation
   - Background scheduler with `time.Ticker`
   - Channel management

2. **Models** (`entity/news.go`)
   - NewsItem struct
   - Source struct

**Key Challenges:**
- RSS feed parsing
- Background task management
- Graceful scheduler shutdown

**Go Libraries:**
- `mmcdole/gofeed` - RSS parsing
- `time.Ticker` - Scheduling

#### 3.6 Whale Alerts
**Complexity:** MEDIUM  
**Effort:** 16 hours

**Components to migrate:**
1. **Whale Service** (`usecase/whale/whale_service.go`)
   - Transaction monitoring
   - Alert generation
   - API integration

2. **Guru Service**
   - Tweet monitoring
   - Sentiment analysis (optional)

3. **Models** (`entity/whale.go`)
   - WhaleTransaction struct
   - GuruTweet struct

**Key Challenges:**
- External API integration
- Rate limiting
- Real-time monitoring

**Go Libraries:**
- `net/http` - API clients
- `golang.org/x/time/rate` - Rate limiting

#### 3.7 Help System
**Complexity:** LOW  
**Effort:** 12 hours

**Components to migrate:**
1. **Help Handlers** (`delivery/discord/handlers/help_handler.go`)
   - Embed generation
   - Pagination with buttons
   - Command reference

2. **Content Management**
   - Static help content from config
   - Template-based formatting

**Key Challenges:**
- Discord embed formatting
- Button interaction handling

### 4. Data Persistence

#### 4.1 Repository Pattern
**Implementation:**
```go
type ConfessionRepository interface {
    Save(confession *entity.Confession) error
    FindByID(id int) (*entity.Confession, error)
    FindByGuild(guildID string) ([]*entity.Confession, error)
    Update(confession *entity.Confession) error
    Delete(id int) error
}

type jsonConfessionRepository struct {
    filePath string
    mu       sync.RWMutex
}
```

**Features:**
- Thread-safe operations
- Atomic writes with temp file + rename
- Automatic backup
- Error recovery

#### 4.2 Data Migration
**Files to migrate:**
- `data/confessions/` - Direct copy
- `data/roasts/` - Direct copy
- Ensure JSON compatibility

### 5. Deployment

#### 5.1 Docker
**Changes:**
- Multi-stage build (builder + runtime)
- Smaller image size with alpine
- Go binary is single executable
- No Python dependencies

**New Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o nerubot cmd/nerubot/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /root/
COPY --from=builder /app/nerubot .
COPY --from=builder /app/data ./data
CMD ["./nerubot"]
```

#### 5.2 Systemd Service
**Changes:**
- Update ExecStart to Go binary
- Remove Python-specific settings
- Add resource limits

#### 5.3 Deployment Scripts
**Updates:**
- `deploy/setup.sh` - Install Go instead of Python
- `deploy/update.sh` - Build Go binary
- `deploy/monitor.sh` - Update for Go binary

## Dependencies Mapping

### Python → Golang Libraries

| Python | Golang | Purpose |
|--------|--------|---------|
| discord.py | github.com/bwmarrin/discordgo | Discord API |
| yt-dlp | os/exec + yt-dlp binary | YouTube download |
| spotipy | github.com/zmb3/spotify | Spotify API |
| feedparser | github.com/mmcdole/gofeed | RSS parsing |
| tweepy | net/http + Twitter API | Twitter integration |
| openai | net/http + OpenAI API | OpenAI client |
| anthropic | net/http + Anthropic API | Claude client |
| google-generativeai | net/http + Gemini API | Gemini client |
| python-dotenv | github.com/joho/godotenv | Environment variables |
| asyncio | goroutines + channels | Concurrency |
| aiohttp | net/http | HTTP client |

## Testing Strategy

### Unit Tests
- Test each use case independently
- Mock repositories and external services
- Use `testify` for assertions

### Integration Tests
- Test repository implementations
- Test Discord command flow
- Use test Discord bot token

### Load Tests
- Concurrent queue operations
- Voice connection stability
- Memory leak detection

## Risk Assessment

### High Risk
1. **Audio Streaming** - Voice encoding complexity
   - Mitigation: Use proven DCA library, extensive testing
   
2. **Data Loss** - Migration of JSON files
   - Mitigation: Backup before migration, validation scripts

### Medium Risk
1. **DiscordGo API Differences** - Different from discord.py
   - Mitigation: Study DiscordGo documentation, prototype first
   
2. **Performance Issues** - Concurrent operations
   - Mitigation: Profiling, benchmarking, gradual rollout

### Low Risk
1. **Configuration Changes** - Different format
   - Mitigation: Conversion scripts, documentation

## Success Metrics

1. **Feature Parity** - All Python features working in Go
2. **Performance** - 50% reduction in memory usage
3. **Reliability** - 99.9% uptime
4. **Response Time** - <100ms for slash commands
5. **Code Quality** - >80% test coverage

## Timeline Overview

### Phase 1: Foundation (Week 1-2)
- Setup Go project structure
- Implement core packages (config, logger, utils)
- Create entity models
- Setup CI/CD

### Phase 2: High Priority Features (Week 3-5)
- Music system
- Confession system
- Help system
- Basic Discord bot

### Phase 3: Medium Priority Features (Week 6-7)
- AI Chatbot
- Roast system

### Phase 4: Low Priority Features (Week 8)
- News system
- Whale alerts

### Phase 5: Testing & Deployment (Week 9-10)
- Integration testing
- Performance testing
- Documentation
- Production deployment

**Total Estimated Time:** 10 weeks (200 hours)

## Resources Required

### Development
- 1 Senior Go Developer (you)
- Go 1.21+ environment
- Discord test server
- API keys for testing

### Infrastructure
- VPS for testing (2GB RAM minimum)
- Docker environment
- CI/CD pipeline (GitHub Actions)

### External Services
- Discord Bot Token
- Spotify API credentials
- OpenAI/Claude/Gemini API keys
- Whale Alert API key
- Twitter API credentials

## Next Steps

1. Review and approve this breakdown
2. Create detailed project plan
3. Generate implementation tickets
4. Setup Go development environment
5. Begin Phase 1 implementation

---

**Document Version:** 1.0  
**Last Updated:** November 9, 2025  
**Status:** Draft - Pending Review
