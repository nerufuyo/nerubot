# User Guide

A complete guide to using NeruBot on your Discord server.

---

## Getting Started

Once NeruBot is added to your server, all commands are available as **slash commands**. Type `/` in any text channel to see the command list.

---

## AI Chat

Have a conversation with an AI assistant powered by DeepSeek.

| Command | What it does |
|---------|-------------|
| `/chat <message>` | Send a message to the AI. Each user has their own conversation history. |
| `/chat-reset` | Clear your conversation history and start fresh. |

The AI remembers your previous messages in the session, so you can have multi-turn conversations.

---

## Confessions

Send anonymous confessions to the server.

| Command | What it does |
|---------|-------------|
| `/confess <content>` | Submit an anonymous confession. It will be posted publicly but your identity stays hidden. |

No one — including server admins — can see who sent a confession.

---

## Roast

Get a humorous roast based on Discord activity.

| Command | What it does |
|---------|-------------|
| `/roast` | Roast yourself. |
| `/roast <user>` | Roast another user. |

Roasts are generated based on the user's activity patterns and profile. All in good fun.

---

## News

Get the latest headlines from multiple news sources.

| Command | What it does |
|---------|-------------|
| `/news` | Fetch and display recent news headlines. |

Headlines are pulled from RSS feeds and displayed in a clean embed.

---

## Whale Alerts

Monitor large cryptocurrency transactions.

| Command | What it does |
|---------|-------------|
| `/whale` | Show recent whale transactions (large crypto transfers). |

Requires `WHALE_ALERT_API_KEY` to be configured.

---

## Analytics

Track server and user activity.

| Command | What it does |
|---------|-------------|
| `/stats` | View server-wide statistics: messages sent, active users, command usage. |
| `/profile` | View your own activity profile. |
| `/profile <user>` | View another user's activity profile. |

---

## Reminders

Automatic reminders for Indonesian national holidays and Ramadan schedule.

### Automatic Notifications

When enabled, the bot sends messages to a configured channel:

- **Hari Libur Nasional** — At 07:00 WIB on every Indonesian national holiday, the bot posts a greeting with `@everyone`. Covers:
  - Fixed: Tahun Baru, Hari Buruh, Hari Lahir Pancasila, Hari Kemerdekaan, Hari Natal
  - Moving: Isra Mi'raj, Imlek, Nyepi, Idul Fitri, Wafat Isa, Waisak, Kenaikan Isa, Idul Adha, Tahun Baru Islam, Maulid Nabi

- **Sahoor** — Reminder at sahoor time (~03:50 WIB) during Ramadan with `@everyone`. The message uses a gentle, warm style in Indonesian.

- **Berbuka** — Reminder at Maghrib time (~17:57 WIB) during Ramadan with `@everyone`. The message uses a warm, encouraging style in Indonesian to celebrate breaking the fast.

### Manual Check

| Command | What it does |
|---------|-------------|
| `/reminder` | View the next upcoming national holidays and today's Ramadan sahoor/berbuka times (if applicable). |

### Setup

Set these in your `.env`:
```
ENABLE_REMINDER=true
REMINDER_CHANNEL_ID=123456789012345678
```

To get a channel ID: enable Developer Mode in Discord settings, right-click a channel, and click "Copy Channel ID".

---

## Help

| Command | What it does |
|---------|-------------|
| `/help` | Show a summary of all available commands. |

---

## Permissions

NeruBot needs these Discord permissions:
- **Send Messages** — respond to commands
- **Embed Links** — display rich embeds
- **Mention Everyone** — `@everyone` in reminders

When inviting the bot, ensure these permissions are granted via the OAuth2 URL in the Discord Developer Portal.

---

## Troubleshooting

| Problem | Solution |
|---------|---------|
| Commands don't appear | Wait a few minutes after bot starts — Discord caches slash commands. Try restarting Discord. |
| Reminders not sending | Verify `ENABLE_REMINDER=true` and `REMINDER_CHANNEL_ID` is set to a valid channel where the bot can send messages. |
| AI chat not responding | Check that `DEEPSEEK_API_KEY` is set and valid. |
| Whale alerts empty | Ensure `WHALE_ALERT_API_KEY` is configured. |
| Bot is offline | Check logs with `make run` or look at container logs with `docker compose logs`. |
