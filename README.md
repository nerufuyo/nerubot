# NeruBot v5.0.1

A feature-rich Discord bot built with **Rust** 🦀

## Features

| Category | Commands | Description |
|----------|----------|-------------|
| **AI Chat** | `/chat`, `/chat-reset` | Chat with DeepSeek AI, per-user history |
| **Roast** | `/roast` | Humorous roasts based on activity |
| **Moderation** | `/kick`, `/ban`, `/timeout`, `/warn`, `/purge` | Server moderation |
| **Fun** | `/coinflip`, `/8ball`, `/meme`, `/dad-joke` | Fun commands |
| **Analytics** | `/stats`, `/profile` | Server & user statistics |
| **Reminders** | `/reminder` | Indonesian holidays & Ramadan |
| **Utility** | `/calc`, `/poll` | Calculator & polls |
| **Help** | `/help` | Show all commands |

## Stack

- **Language:** Rust (2024 Edition)
- **Discord:** Serenity
- **Database:** PostgreSQL (shared with Neruwork)
- **Cache:** Redis (shared with Neruwork)
- **AI:** DeepSeek API

## Quick Start

```bash
# Clone
git clone https://github.com/nerufuyo/nerubot-rs.git
cd nerubot-rs

# Configure
cp .env.example .env
# Edit .env with your DISCORD_TOKEN

# Run
cargo run --release
```

## Docker

```bash
docker compose up -d --build
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DISCORD_TOKEN` | ✅ | Discord bot token |
| `DATABASE_URL` | ✅ | PostgreSQL connection string |
| `REDIS_URL` | ✅ | Redis connection string |
| `DEEPSEEK_API_KEY` | ⚡ | For AI chat feature |

## Database

Uses the **same PostgreSQL database** as Neruwork. Migrations run automatically on startup.

Tables:
- `guild_config` — Server settings
- `chat_history` — AI chat sessions
- `roasts` — Roast history
- `mod_logs` — Moderation logs
- `warnings` — User warnings
- `message_stats` — Message analytics
- `command_usage` — Command analytics
- `reminders` — User reminders
- `polls` — Active polls

## License

MIT
