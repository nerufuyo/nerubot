# NeruBot

A feature-rich Discord bot built with Rust. AI chat, moderation, utility tools, fun commands, confessions, roasts, news, whale alerts, analytics, music, and scheduled reminders.

**v5.0.1** | Rust 2024 Edition | MIT License

---

## Features

| Feature | Description |
|---------|-------------|
| **Core** | Ping, bot info, server info, user info, avatar |
| **AI Chat** | Chat with DeepSeek AI, per-user history |
| **Moderation** | Kick, ban, timeout, purge messages, warn system |
| **Utility** | Calculator, polls with reactions |
| **Fun** | Coinflip, 8ball, memes, dad jokes, confessions, roasts |
| **Music** | Full music player with queue, filters, playlists, radio |
| **News** | Latest headlines from RSS feeds |
| **Whale Alerts** | Large cryptocurrency transactions |
| **Analytics** | Server and user statistics |
| **Mental Health** | Scheduled mental health tips and reminders |
| **Reminders** | Indonesian national holidays and Ramadan schedule |
| **Security** | Spam detection, malicious link filter, content filter |
| **Multi-language** | EN, ID, JP, KR, ZH support |

---

## Quick Start

### Prerequisites

- Rust (2024 edition)
- A Discord bot token ([Discord Developer Portal](https://discord.com/developers/applications))

### Setup

```bash
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot
cp .env.example .env
# Edit .env with your values
```

### Run

```bash
# Direct
cargo run

# Or build first
cargo build --release
./target/release/nerubot
```

### Docker

```bash
docker compose up -d
```

---

## Commands

All commands use Discord slash commands (`/`). Use `/help` for paginated interactive help with navigation buttons.

### Core Commands
| Command | Description |
|---------|-------------|
| `/ping` | Check bot latency |
| `/botinfo` | Show bot info, uptime, developer |
| `/serverinfo` | Show server details |
| `/userinfo [user]` | Show user information |
| `/avatar [user]` | Display a user's avatar |
| `/help [lang]` | Interactive paginated help menu |

### Fun Commands
| Command | Description |
|---------|-------------|
| `/coinflip` | Flip a coin |
| `/8ball <question>` | Ask the magic 8-ball |
| `/dadjoke` | Get a random dad joke |
| `/meme` | Get a random meme |
| `/confess <content>` | Submit anonymous confession |
| `/roast [user]` | Roast a user based on activity |

### Utility Commands
| Command | Description |
|---------|-------------|
| `/chat <message>` | Chat with AI |
| `/chat-reset` | Clear your chat history |
| `/calc <a> <op> <b>` | Simple math calculator |
| `/poll <question> <options>` | Create a poll with reactions |
| `/news [lang]` | Latest news headlines |
| `/whale` | Recent whale crypto transactions |
| `/stats` | Server statistics |
| `/profile [user]` | User activity profile |

### Moderation Commands (Admin/Mod)
| Command | Description |
|---------|-------------|
| `/kick <user> [reason]` | Kick a user from the server |
| `/ban <user> [reason]` | Ban a user from the server |
| `/timeout <user> [duration] [reason]` | Temporarily mute a user |
| `/purge <amount>` | Delete multiple messages (1-100) |
| `/warn <user> [reason]` | Warn a user and record it |
| `/warnings <user>` | View warnings for a user |
| `/clearwarnings <user>` | Clear all warnings for a user |

### Scheduling Commands
| Command | Description |
|---------|-------------|
| `/reminder` | Upcoming holidays and Ramadan schedule |
| `/reminder-set <channel>` | Set reminder channel (admin) |
| `/mentalhealth [lang]` | Get a mental health tip |
| `/mentalhealth-setup` | Schedule mental health reminders |
| `/dadjoke-setup <channel> <interval>` | Schedule dad jokes |
| `/meme-setup <channel> <interval>` | Schedule memes |

### Music Commands
| Command | Description |
|---------|-------------|
| `/play <query>` | Play a song from YouTube, Spotify, SoundCloud |
| `/pause` `/resume` `/stop` | Playback controls |
| `/skip` `/previous` `/seek` | Track navigation |
| `/queue` `/shuffle` `/volume` | Queue management |
| `/loop` `/filter` `/playlist` | Advanced features |
| `/lyrics` `/recommend` `/radio` | Discovery & extras |

---

## Configuration

Copy `.env.example` to `.env` and set your values:

```dotenv
# Required
DISCORD_TOKEN=your_bot_token

# AI (optional)
DEEPSEEK_API_KEY=your_key

# Feature toggles (true/false)
ENABLE_CONFESSION=true
ENABLE_ROAST=true
ENABLE_REMINDER=true

# Whale alerts (optional)
WHALE_ALERT_API_KEY=

# Reminder channel
REMINDER_CHANNEL_ID=your_channel_id

# Runtime
LOG_LEVEL=INFO
ENVIRONMENT=development
```

### Feature Flags

Each feature can be toggled independently. If a feature is disabled, its commands still register but respond with a "not enabled" message. The reminder feature requires `REMINDER_CHANNEL_ID` to send scheduled messages.

---

## Reminder Feature

Sends automatic reminders to a configured Discord channel:

**Indonesian National Holidays** — Posted at 07:00 WIB on the holiday date with `@everyone` tag. Covers fixed holidays (Tahun Baru, Hari Buruh, Pancasila, Kemerdekaan, Natal) and moving holidays (Idul Fitri, Idul Adha, Nyepi, Imlek, Isra Mi'raj, Waisak, Maulid Nabi, etc.) for 2025-2027.

**Ramadan Sahoor** — Reminder at sahoor time (around 03:50 WIB) with `@everyone` tag. Uses a warm, romantic Indonesian style: *"Hai sayang... bangun dong, jangan ketiduran..."*

**Ramadan Berbuka** — Reminder at Maghrib time (around 17:57 WIB) with `@everyone` tag. Romantic style: *"Alhamdulillah... buka dengan yang manis ya — seperti senyummu..."*

Use `/reminder` to view upcoming holidays and today's Ramadan schedule.

---

## Project Structure

```
nerubot/
├── src/main.rs                  # Entry point
├── src/
│   ├── commands/                # Slash commands (chat, roast, fun, moderation, etc.)
│   ├── handlers/                # Interaction and event handlers
│   ├── config/                  # Environment configuration
│   ├── db/                      # Database connection and migrations
│   ├── models/                  # Domain models
│   ├── utils/                   # AI client, reminder utilities
│   ├── api/                     # Admin API + Lark webhook
│   ├── lark/                    # Lark bot integration
│   └── workers/                 # Background workers (reminders)
├── static/                      # Static assets
├── migrations/                  # SQL migration files
├── Cargo.toml                   # Rust dependencies
├── Dockerfile                   # Container build
├── docker-compose.yml           # Docker orchestration
└── Makefile                     # Build tasks
```

The project follows **Clean Architecture**: `delivery` -> `usecase` -> `entity` -> `repository`. Each feature is isolated in its own usecase package.

---

## Development

### Build

```bash
make build          # Build release binary
make run            # Build and run
make clean          # Remove build artifacts
```

### Code Quality

```bash
make fmt            # Format code (cargo fmt)
make check          # Check code (cargo check)
make test           # Run tests (cargo test)
make lint           # Run clippy linter
```

### Adding a Feature

1. Create module in `src/commands/<feature>.rs`
2. Add handler function with standard signature
3. Register in `src/commands/mod.rs`
4. Add route in `src/handlers/mod.rs` match block
5. Register slash command in `src/main.rs` `ready()` handler
6. Update help text in `show_help()`

---

## Deployment

### Railway

Push to GitHub and connect to [Railway](https://railway.app). Set environment variables in Railway dashboard. The `railway.toml` is pre-configured.

### Docker

```bash
docker compose up -d
```

### VPS / Bare Metal

See `deploy/` for systemd service, nginx, logrotate, and cron configurations.

```bash
cd deploy
chmod +x setup.sh
./setup.sh
```

---

## License

MIT — see [LICENSE](LICENSE).

Built by [@nerufuyo](https://github.com/nerufuyo).

## GitHub About suggestion

See `GITHUB_ABOUT.md` for a ready-to-paste repository description and topics.
