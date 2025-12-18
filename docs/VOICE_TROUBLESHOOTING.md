# Voice Connection Troubleshooting Guide

## Known Issues & Solutions

### Issue: "Unknown encryption mode" Error

**Error Message:**
```
Voice websocket error: 4016 - Unknown encryption mode
timeout waiting for voice connection
```

**Root Cause:**
Discord's voice protocol uses encryption modes that `discordgo` (the Go Discord library) doesn't fully support. This is a library limitation, not a bug in NeruBot.

**Solution:**
Use **Lavalink** for voice audio - it handles all voice protocol complexity. See [MUSIC_SETUP.md](MUSIC_SETUP.md).

### Issue: "timeout waiting for voice"

**What happens:**
- Bot attempts to join voice channel
- Connection starts but fails during encryption negotiation
- Bot times out and disconnects

**Why:**
Discord changed voice encryption modes. `discordgo` is incompatible with newer modes.

**Fix:**
Configure Lavalink (see [MUSIC_SETUP.md](MUSIC_SETUP.md) - Quick Start with Docker section).

---

## Quick Fixes

### Temporary: Disable Voice Joining

If you don't need voice audio yet, the bot still works for:
- ✅ Song searching
- ✅ Queue management  
- ✅ /queue commands
- ✅ All other features

All voice-related code is safely skipped.

### Permanent: Set Up Lavalink

**Fastest option (local development):**
```bash
docker-compose -f docker-compose.music.yml up -d
```

This automatically:
- Starts Lavalink
- Configures the bot
- Enables all music features

**Result:** Full voice audio playback in under 2 minutes.

---

## What Works Without Lavalink

| Feature | Status |
|---------|--------|
| Music Search | ✅ Works |
| Queue Management | ✅ Works |
| `/play` command | ✅ Works (queues, no audio) |
| `/queue` command | ✅ Works |
| `/skip` command | ✅ Works (queue only) |
| `/stop` command | ✅ Works |
| Voice Channel Audio | ❌ Needs Lavalink |
| Skip While Playing | ❌ Needs Lavalink |
| Audio Playback | ❌ Needs Lavalink |

---

## Architecture Explanation

```
Without Lavalink:
┌─────────────┐
│   Discord   │
│     Bot     │ ← Can't send audio (discordgo limitation)
└──────┬──────┘
       │
    (Queue
     Only)
       │
       ▼
┌─────────────┐
│   Discord   │
│   Channel   │ ← Bot joins but can't play audio
└─────────────┘


With Lavalink (Correct Setup):
┌─────────────┐         ┌─────────────┐
│   Discord   │────────▶│  Lavalink   │
│     Bot     │ (API)   │   Server    │
└─────────────┘         └──────┬──────┘
                               │
                          (Handles
                           voice
                          protocol)
                               │
                               ▼
                        ┌─────────────┐
                        │   Discord   │
                        │   Channel   │ ← Audio streams properly
                        └─────────────┘
```

---

## Environment Variables

### Current (Without Lavalink)
```env
LAVALINK_ENABLED=false
LAVALINK_HOST=localhost
LAVALINK_PORT=2333
```

### Updated (With Lavalink)
```env
LAVALINK_ENABLED=true
LAVALINK_HOST=localhost
LAVALINK_PORT=2333
LAVALINK_PASSWORD=youshallnotpass
```

After updating `.env`, restart the bot:
```bash
# Local development
go run cmd/nerubot/main.go

# Railway
railway up
```

---

## Testing Your Setup

### Test 1: Can bot join voice channel?

```bash
# Check in Discord - bot should attempt to join
# Look for any errors in logs
```

### Test 2: Check Lavalink connection

```bash
# If Lavalink is running
curl -H "Authorization: youshallnotpass" http://localhost:2333/v4/info

# Expected response:
# {"version":{"major":4,"minor":X,"patch":X},"buildTime":...}
```

### Test 3: Check bot can find songs

```bash
# Use /play command - should show search results
# Bot should queue songs successfully
```

### Test 4: Check audio playback (Lavalink only)

```bash
# Use /play command while in voice channel
# Should hear audio in Discord
```

---

## Common Misconceptions

❌ **"The bot is broken"**
✅ It's working! Just needs Lavalink for audio.

❌ **"yt-dlp is broken"**  
✅ yt-dlp works fine for searching. It's Discord voice that needs Lavalink.

❌ **"Fix the voice encryption error"**
✅ Can't be fixed in Go without using Lavalink. It's a library design issue.

---

## Getting Help

**Issue: Lavalink won't start**
- Check Java is installed: `java -version`
- Check port 2333 is available: `netstat -an | grep 2333`

**Issue: Bot can't connect to Lavalink**
- Verify host/port in `.env`
- Verify password is correct
- Check firewall allows connection

**Issue: Audio still plays in voice but sounds bad**
- Check bot permissions (Speak/Connect)
- Check voice channel quality settings
- Check bandwidth availability

---

## Next Steps

1. **Read:** [MUSIC_SETUP.md](MUSIC_SETUP.md) - Full setup guide
2. **Quick Start:** Use `docker-compose -f docker-compose.music.yml up`
3. **Test:** Try `/play command` in Discord
4. **Enjoy:** Full music bot functionality!

---

## References

- [Lavalink Official Docs](https://lavalink.dev/)
- [NeruBot Music Setup](MUSIC_SETUP.md)
- [NeruBot Architecture](ARCHITECTURE.md)
- [Deployment Guide](DEPLOYMENT.md)
