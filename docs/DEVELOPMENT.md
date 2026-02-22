# Development Guide

How to set up, build, and extend NeruBot locally.

---

## Prerequisites

- **Go 1.21+** — [install](https://go.dev/dl/)
- **Git**
- **Discord bot token** — from the [Developer Portal](https://discord.com/developers/applications)

---

## Local Setup

```bash
# Clone
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Configure
cp .env.example .env
# Edit .env — at minimum set DISCORD_TOKEN

# Run
go run ./cmd/nerubot
```

The bot will start, register slash commands, and connect to Discord.

---

## Build

```bash
make build          # Compile to ./build/nerubot
make run            # Build and execute
make clean          # Remove build artifacts
```

Or directly:
```bash
go build -o nerubot ./cmd/nerubot
```

---

## Code Quality

```bash
go fmt ./...        # Format all files
go vet ./...        # Static analysis
go test ./...       # Run tests
```

Or via Make:
```bash
make fmt
make vet
make test
```

---

## Architecture

The project follows **Clean Architecture** with clear dependency direction:

```
delivery (Discord handlers)
    ↓
usecase (business logic)
    ↓
entity (domain models)
    ↓
repository (data persistence)
```

### Directory Layout

| Directory | Purpose |
|-----------|---------|
| `cmd/nerubot/` | Entry point (`main.go`) |
| `internal/config/` | Configuration loading, constants, message templates |
| `internal/delivery/discord/` | Discord session, command routing, handler files |
| `internal/entity/` | Domain types (no dependencies on other packages) |
| `internal/usecase/` | One package per feature with business logic |
| `internal/repository/` | JSON file storage |
| `internal/pkg/` | Shared utilities: AI client, logger |
| `data/` | Runtime JSON files (gitignored) |
| `deploy/` | Production configs (systemd, nginx, cron) |

### Key Files

| File | Role |
|------|------|
| `cmd/nerubot/main.go` | Loads config, creates bot, handles shutdown signals |
| `internal/delivery/discord/bot.go` | Bot lifecycle: session, services init, command routing, slash command registration |
| `internal/delivery/discord/handlers.go` | Help command and shared response helpers |
| `internal/delivery/discord/handler_*.go` | One file per feature (chatbot, confession, roast, news, whale, analytics, reminder) |

---

## Adding a New Feature

### 1. Define Domain Types

Create `internal/entity/<feature>.go` with your structs and types:

```go
package entity

type MyFeature struct {
    ID    string
    Value string
}
```

### 2. Create Service

Create `internal/usecase/<feature>/<feature>_service.go`:

```go
package feature

type FeatureService struct {
    // dependencies
}

func NewFeatureService() *FeatureService {
    return &FeatureService{}
}

func (s *FeatureService) DoSomething() (string, error) {
    return "result", nil
}
```

### 3. Add Handler

Create `internal/delivery/discord/handler_<feature>.go`:

```go
package discord

import "github.com/bwmarrin/discordgo"

func (b *Bot) handleMyFeature(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if err := b.deferResponse(s, i); err != nil {
        return
    }
    // call service, build response
    b.followUp(s, i, "response text")
}
```

### 4. Wire in bot.go

1. Add a struct field: `myFeatureService *feature.FeatureService`
2. Initialize in `New()`: `b.myFeatureService = feature.NewFeatureService()`
3. Add command routing in `handleCommand()`: `case "myfeature": b.handleMyFeature(s, i)`
4. Register slash command in `registerCommands()`

### 5. Update Help

Add a field entry in the help embed in `handlers.go`.

### 6. Add Config (optional)

If the feature needs environment variables:
- Add a config struct and field in `internal/config/config.go`
- Read the env var in `Load()`
- Add to `.env.example`

---

## Configuration Reference

All configuration is via environment variables (loaded from `.env`).

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DISCORD_TOKEN` | Yes | — | Bot token from Discord Developer Portal |
| `DEEPSEEK_API_KEY` | No | — | DeepSeek AI API key for `/chat` |
| `ENABLE_CONFESSION` | No | `true` | Enable anonymous confessions |
| `ENABLE_ROAST` | No | `true` | Enable roast command |
| `ENABLE_REMINDER` | No | `true` | Enable holiday/Ramadan reminders |
| `WHALE_ALERT_API_KEY` | No | — | Whale Alert API key |
| `REMINDER_CHANNEL_ID` | No | — | Channel ID for automatic reminders |
| `LOG_LEVEL` | No | `INFO` | Log level: DEBUG, INFO, WARN, ERROR |
| `ENVIRONMENT` | No | `development` | Runtime environment |

---

## Docker

### Build and Run

```bash
docker compose up -d --build
```

### Logs

```bash
docker compose logs -f
```

### Rebuild

```bash
docker compose down
docker compose up -d --build
```

---

## Deployment

### Railway

1. Push to GitHub
2. Connect repo on [Railway](https://railway.app)
3. Set environment variables in Railway dashboard
4. Deploy — `railway.toml` handles the rest

### VPS

Use the scripts in `deploy/`:

```bash
cd deploy
chmod +x setup.sh
./setup.sh
```

This sets up a systemd service, nginx reverse proxy, logrotate, and cron monitoring.
