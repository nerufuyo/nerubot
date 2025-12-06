# NeruBot Project Structure

This document provides a detailed overview of NeruBot's file and directory organization.

---

## 📁 Root Directory Structure

```
nerubot/
├── cmd/                    # Application entry points
├── internal/               # Private application code
├── data/                   # Data storage (JSON files)
├── deploy/                 # Deployment configurations
├── docs/                   # Documentation
├── .env.example            # Environment template
├── .gitignore             # Git ignore rules
├── ARCHITECTURE.md        # Architecture overview
├── CHANGELOG.md           # Version history
├── CONTRIBUTING.md        # Contribution guidelines
├── docker-compose.yml     # Docker Compose config
├── Dockerfile             # Docker build file
├── go.mod                 # Go module definition
├── go.sum                 # Go dependency checksums
├── LICENSE                # MIT License
├── Makefile               # Build automation
└── README.md              # Project documentation
```

---

## 🏗️ Detailed Structure

### `/cmd/` - Application Entry Points

```
cmd/
└── nerubot/
    └── main.go             # Bot application entry point
```

**Purpose:**
- Contains executable applications
- Minimal logic (initialization only)
- Calls internal packages

**main.go responsibilities:**
- Load configuration
- Initialize logger
- Create bot instance
- Start bot
- Handle graceful shutdown

---

### `/internal/` - Private Application Code

The core of the application, organized by Clean Architecture layers.

```
internal/
├── config/
│   ├── config.go           # Main configuration structure
│   ├── constants.go        # Constants and defaults
│   └── messages.go         # Response messages
├── delivery/
│   └── discord/
│       ├── bot.go          # Bot initialization
│       └── handlers.go     # Command handlers
├── entity/
│   ├── confession.go       # Confession domain model
│   ├── music.go            # Music domain model
│   ├── news.go             # News domain model
│   ├── roast.go            # Roast domain model
│   └── whale.go            # Whale domain model
├── pkg/
│   ├── ai/
│   │   ├── deepseek.go     # DeepSeek AI integration
│   │   └── interface.go    # AI provider interface
│   ├── ffmpeg/
│   │   └── ffmpeg.go       # FFmpeg wrapper
│   ├── logger/
│   │   └── logger.go       # Logging utility
│   └── ytdlp/
│       └── ytdlp.go        # yt-dlp wrapper
├── repository/
│   ├── confession_repository.go  # Confession data access
│   ├── repository.go             # Base repository
│   └── roast_repository.go       # Roast data access
└── usecase/
    ├── chatbot/
    │   └── chatbot_service.go    # AI chatbot logic
    ├── confession/
    │   └── confession_service.go # Confession management
    ├── music/
    │   └── music_service.go      # Music streaming
    ├── news/
    │   └── news_service.go       # News aggregation
    ├── roast/
    │   └── roast_service.go      # Roast generation
    └── whale/
        └── whale_service.go      # Whale alerts
```

#### `/internal/config/` - Configuration Management

| File | Purpose | Key Components |
|------|---------|----------------|
| `config.go` | Main configuration | Config struct, Load(), Validate() |
| `constants.go` | Constants | Default values, limits, timeouts |
| `messages.go` | Response messages | Success/error messages, embeds |

#### `/internal/entity/` - Domain Models

| File | Purpose | Key Structs |
|------|---------|-------------|
| `confession.go` | Confession models | Confession, ConfessionSettings, Reply |
| `music.go` | Music models | Song, Queue, PlaybackState |
| `roast.go` | Roast models | UserProfile, RoastPattern, RoastHistory |

**Principles:**
- Pure Go structs
- No external dependencies
- Business logic only
- JSON/YAML tags for serialization

#### `/internal/repository/` - Data Persistence

| File | Purpose | Responsibilities |
|------|---------|------------------|
| `repository.go` | Base interfaces | Common CRUD operations |
| `confession_repository.go` | Confession storage | Save, find, update confessions |
| `roast_repository.go` | Roast storage | Save, find user profiles |

**Current Implementation:** JSON file-based storage  
**Future:** PostgreSQL/MongoDB migration ready

#### `/internal/usecase/` - Business Logic

Each feature has its own service:

| Directory | Service | Responsibilities |
|-----------|---------|------------------|
| `chatbot/` | ChatbotService | AI conversation management |
| `confession/` | ConfessionService | Confession workflow, moderation |
| `music/` | MusicService | Playback, queue, audio streaming |
| `news/` | NewsService | RSS aggregation, publishing |
| `roast/` | RoastService | Activity tracking, roast generation |
| `whale/` | WhaleService | Transaction monitoring, alerts |

**Design Pattern:** Service pattern with dependency injection

#### `/internal/delivery/discord/` - Discord Interface

| File | Purpose | Key Functions |
|------|---------|---------------|
| `bot.go` | Bot lifecycle | New(), Start(), Stop(), registerHandlers() |
| `handlers.go` | Command handlers | handlePlayCommand(), handleConfessCommand(), etc. |

**Responsibilities:**
- Discord API interaction
- Command registration
- Event handling
- Response formatting

#### `/internal/pkg/` - Shared Utilities

| Package | Purpose | Key Types |
|---------|---------|-----------|
| `ai/` | AI provider abstraction | AIProvider interface, DeepSeekProvider |
| `ffmpeg/` | Audio processing | FFmpeg struct, Convert(), Stream() |
| `logger/` | Structured logging | Logger struct, Info(), Error(), Debug() |
| `ytdlp/` | YouTube extraction | YtDlp struct, ExtractInfo(), GetStreamURL() |

---

### `/data/` - Data Storage

```
data/
├── confessions/
│   ├── confessions.json    # Active confessions
│   ├── queue.json          # Pending moderation
│   ├── replies.json        # Confession replies
│   └── settings.json       # Per-guild settings
└── roasts/
    ├── activities.json     # User activity tracking
    ├── patterns.json       # Roast pattern templates
    ├── profiles.json       # User profiles
    └── stats.json          # Global statistics
```

**Storage Format:** JSON  
**Access:** Through repository layer only  
**Backup:** Recommended daily backups

---

### `/deploy/` - Deployment Configurations

```
deploy/
├── setup.sh                # VPS setup script
├── monitor.sh              # Health monitoring script
├── status.sh               # Status check script
├── update.sh               # Update script
├── README.md               # Deployment documentation
├── cron/
│   └── nerubot-crontab     # Cron jobs
├── logrotate/
│   └── nerubot             # Log rotation config
├── nginx/
│   └── nerubot.conf        # Nginx configuration
└── systemd/
    └── nerubot.service     # Systemd service file
```

**Usage:**
- Automated VPS deployment
- Service management
- Monitoring and maintenance

---

### `/docs/` - Documentation

```
docs/
├── ARCHITECTURE.md         # System architecture guide
├── DEPLOYMENT.md           # Deployment instructions
└── PROJECT_STRUCTURE.md    # This file
```

**Purpose:**
- Technical documentation
- Deployment guides
- Architecture explanations

---

## 📦 Key Files

### Configuration Files

#### `.env` (Created from `.env.example`)

```env
# Discord
DISCORD_TOKEN=...

# AI
DEEPSEEK_API_KEY=...

# Features
ENABLE_MUSIC=true
ENABLE_CONFESSION=true
ENABLE_ROAST=true
```

**Security:** Never commit to git (listed in `.gitignore`)

#### `docker-compose.yml`

```yaml
version: '3.8'
services:
  nerubot:
    build: .
    env_file: .env
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    restart: unless-stopped
```

**Purpose:** Local development with Docker

#### `Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
# ... build stage ...

FROM alpine:latest
# ... runtime stage ...
```

**Features:**
- Multi-stage build
- Alpine Linux (minimal size)
- Non-root user
- Health check

#### `Makefile`

```makefile
build:
    go build -o build/nerubot cmd/nerubot/main.go

run: build
    ./build/nerubot

clean:
    rm -rf build/
```

**Commands:** `make build`, `make run`, `make clean`

### Go Module Files

#### `go.mod`

```go
module github.com/nerufuyo/nerubot

go 1.21

require (
    github.com/bwmarrin/discordgo v0.29.0
    // ... other dependencies
)
```

**Purpose:** Define module and dependencies

#### `go.sum`

**Purpose:** Dependency checksums for security

---

## 🔍 File Naming Conventions

### General Rules

1. **Lowercase with underscores:** `user_profile.go`
2. **Descriptive names:** `confession_repository.go` not `repo.go`
3. **Suffix for type:** `_service.go`, `_repository.go`, `_test.go`

### Examples

**Good:**
```
✅ confession_service.go
✅ music_repository.go
✅ deepseek_provider.go
✅ config_test.go
```

**Avoid:**
```
❌ svc.go (too abbreviated)
❌ ConfessionService.go (wrong case)
❌ confessionService.go (wrong case)
❌ cs.go (not descriptive)
```

---

## 📝 Code Organization Principles

### 1. Package by Feature

**Good:**
```
usecase/
├── music/
│   └── music_service.go
├── confession/
│   └── confession_service.go
```

**Bad:**
```
usecase/
├── services/
│   ├── music.go
│   └── confession.go
```

### 2. Minimal Package Dependencies

```
entity/      → No dependencies
repository/  → Depends on entity/
usecase/     → Depends on entity/, repository/
delivery/    → Depends on usecase/
```

### 3. Interface in User Package

```go
// usecase/music/service.go
type Repository interface {
    Save(song *entity.Song) error
}

// repository/music_repository.go
type MusicRepository struct {}

func (r *MusicRepository) Save(song *entity.Song) error {
    // Implementation
}
```

---

## 🧪 Test Files

### Location

Tests are co-located with source files:

```
usecase/music/
├── music_service.go
└── music_service_test.go
```

### Naming

- Test file: `*_test.go`
- Test function: `Test<Function>` (e.g., `TestPlay`)
- Benchmark: `Benchmark<Function>`

### Example

```go
// music_service_test.go
package music_test

func TestMusicService_Play(t *testing.T) {
    // Test implementation
}
```

---

## 📚 Import Organization

### Order

1. Standard library
2. External packages
3. Internal packages

### Example

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // External packages
    "github.com/bwmarrin/discordgo"

    // Internal packages
    "github.com/nerufuyo/nerubot/internal/config"
    "github.com/nerufuyo/nerubot/internal/entity"
)
```

---

## 🔧 Build Artifacts

### `/build/` Directory (Created at build time)

```
build/
└── nerubot              # Compiled binary
```

**Note:** Listed in `.gitignore`, not committed to repository

### Generated Files

- `build/nerubot` - Main executable
- `logs/*.log` - Log files
- Coverage reports (during testing)

---

## 🚫 What's NOT in the Repository

Files excluded by `.gitignore`:

```
.env                # Environment variables (secrets)
build/              # Compiled binaries
logs/               # Log files
*.log               # Any log files
data/backups/       # Data backups
.DS_Store           # macOS files
*.swp               # Vim swap files
```

---

## 📖 Related Documentation

- [Architecture Guide](ARCHITECTURE.md) - System design principles
- [Deployment Guide](DEPLOYMENT.md) - Production deployment
- [Contributing Guide](../CONTRIBUTING.md) - Development workflow

---

## 🔗 Quick Navigation

**Find a file:**
```bash
# Find by name
find . -name "music_service.go"

# Find by pattern
find . -name "*_repository.go"

# List all Go files
find . -name "*.go" -not -path "*/vendor/*"
```

**Count lines of code:**
```bash
# Total Go code
find . -name "*.go" -not -path "*/vendor/*" | xargs wc -l

# By directory
find internal/usecase -name "*.go" | xargs wc -l
```

---

**Last Updated:** December 6, 2025  
**Version:** 3.0.0  
**Author:** @nerufuyo
