# NeruBot Testing Guide

## Prerequisites

Before testing, ensure you have:

1. **Discord Bot Token** - Set in `.env` file
   ```bash
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

2. **DeepSeek API Key** (Optional - for chatbot)
   ```bash
   DEEPSEEK_API_KEY=your_deepseek_api_key_here
   ```

3. **Bot invited to your Discord server** with proper permissions:
   - Send Messages
   - Embed Links
   - Attach Files
   - Use Slash Commands
   - Connect (for music)
   - Speak (for music)

## Running the Bot

### Option 1: Run Monolithic Bot (Recommended for Testing)
```powershell
# Build and run the main bot
go build -o build/bot.exe ./cmd/nerubot
.\build\bot.exe
```

### Option 2: Run Microservices Architecture
```powershell
# Terminal 1 - Gateway (Discord Bot)
go run ./services/gateway/cmd

# Terminal 2 - Music Service
go run ./services/music/cmd

# Terminal 3 - Confession Service
go run ./services/confession/cmd

# Terminal 4 - Roast Service
go run ./services/roast/cmd

# Terminal 5 - Chatbot Service
go run ./services/chatbot/cmd

# Terminal 6 - News Service
go run ./services/news/cmd

# Terminal 7 - Whale Service
go run ./services/whale/cmd
```

## Testing Checklist

### 1. Basic Bot Connection ✅
- [ ] Bot comes online in Discord
- [ ] Bot status shows as "Online"
- [ ] Bot custom status displays correctly

### 2. Slash Commands Registration ✅
Check if slash commands appear when typing `/` in Discord:
- [ ] `/play` - Music playback
- [ ] `/pause` - Pause music
- [ ] `/resume` - Resume music
- [ ] `/skip` - Skip track
- [ ] `/stop` - Stop playback
- [ ] `/queue` - Show queue
- [ ] `/nowplaying` - Current track
- [ ] `/confess` - Submit confession
- [ ] `/roast` - Get roasted
- [ ] `/profile` - View profile
- [ ] `/ping` - Bot status

### 3. Music Feature Tests 🎵

#### Test: Play Music
```
Command: /play song:Never Gonna Give You Up
Expected: Bot joins voice channel and plays music
Status: [ ]
```

#### Test: Queue Management
```
Command: /queue
Expected: Display current queue with song list
Status: [ ]
```

#### Test: Playback Controls
```
Commands:
  /pause   -> Should pause playback
  /resume  -> Should resume playback
  /skip    -> Should skip to next track
  /stop    -> Should stop and clear queue
Status: [ ]
```

#### Test: Now Playing
```
Command: /nowplaying
Expected: Display current track info with requester
Status: [ ]
```

### 4. Confession Feature Tests 🤐

#### Test: Submit Confession
```
Command: /confess message:This is a test confession
Expected: Confession submitted to moderation queue
Status: [ ]
```

**Note**: Confession feature requires:
- Admin approval system
- Anonymous submission
- Moderation queue

### 5. Roast Feature Tests 🔥

#### Test: Self Roast
```
Command: /roast
Expected: Bot generates a roast based on your activity
Status: [ ]
```

#### Test: Roast Another User
```
Command: /roast user:@username
Expected: Bot generates a roast for specified user
Status: [ ]
```

#### Test: View Profile
```
Command: /profile
Expected: Display activity statistics
  - Message count
  - Reaction count
  - Voice time
  - Command count
Status: [ ]
```

### 6. Chatbot Feature Tests 🤖

**Note**: Requires `DEEPSEEK_API_KEY` in `.env` file

#### Test: Basic Chat
```
Command: /chat message:Hello, how are you?
Expected: AI response from DeepSeek
Status: [ ]
```

#### Test: Context Retention
```
Commands:
  1. /chat message:My name is John
  2. /chat message:What is my name?
Expected: Second response should remember "John"
Status: [ ]
```

### 7. News Feature Tests 📰

#### Test: Fetch News
```
Command: /news
Expected: Display recent news articles
Status: [ ]
```

### 8. Whale Alerts Feature Tests 🐋

**Note**: Requires `WHALE_ALERT_API_KEY` in `.env` file

#### Test: Get Whale Transactions
```
Command: /whale
Expected: Display recent large crypto transactions
Status: [ ]
```

### 9. Ping/Health Check ✅

#### Test: Bot Status
```
Command: /ping
Expected: "🏓 Pong! API Gateway is running."
Status: [ ]
```

## Testing Microservices

### Health Check Endpoints

Test each service's health endpoint:

```powershell
# Gateway
curl http://localhost:8080/health

# Music Service
curl http://localhost:8081/health

# Confession Service
curl http://localhost:8082/health

# Roast Service
curl http://localhost:8083/health

# Chatbot Service
curl http://localhost:8084/health

# News Service
curl http://localhost:8085/health

# Whale Service
curl http://localhost:8086/health
```

Expected response for each:
```json
{"status":"healthy","service":"service-name","version":"1.0.0"}
```

### gRPC Communication Test

If running microservices, verify gRPC ports are listening:

```powershell
# Check if gRPC ports are open
netstat -an | Select-String "50051|50052|50053|50054|50055|50056"
```

Expected:
- 50051 - Music Service
- 50052 - Confession Service
- 50053 - Roast Service
- 50054 - Chatbot Service
- 50055 - News Service
- 50056 - Whale Service

## Troubleshooting

### Bot doesn't start
```powershell
# Check if Discord token is set
$env:DISCORD_TOKEN
# Should output your token, not "your_discord_bot_token_here"
```

### Music doesn't play
- Ensure you're in a voice channel before using `/play`
- Check FFmpeg is installed: `ffmpeg -version`
- Check yt-dlp is installed: `yt-dlp --version`

### Chatbot doesn't respond
- Verify DeepSeek API key is set and valid
- Check logs for API errors
- Test API key: `curl https://api.deepseek.com/v1/models -H "Authorization: Bearer YOUR_KEY"`

### Slash commands don't appear
- Wait 1-2 minutes after bot starts (Discord sync delay)
- Try running `/` in a different channel
- Reinvite bot with proper OAuth2 scopes

## Expected Log Output

When bot starts successfully:

```
[INFO] === NeruBot v3.0.0 (Golang Edition) ===
[INFO] Starting Discord bot...
[INFO] Configuration loaded successfully
[INFO] Features enabled
[INFO] Repositories initialized successfully
[INFO] Music service initialized
[INFO] Confession service initialized
[INFO] Roast service initialized
[INFO] Bot is running. Press CTRL+C to exit.
```

## Performance Tests

### Load Test: Multiple Commands
```
1. Send 10 /ping commands rapidly
   Expected: All respond within 2 seconds

2. Queue 5 songs with /play
   Expected: All added to queue successfully

3. Send 3 /roast commands
   Expected: Cooldown enforced after first roast
```

### Stress Test: Concurrent Users
```
Have 3-5 users simultaneously:
- Use different commands
- Play music
- Send chat messages
- Request roasts

Expected: All commands process successfully
```

## Error Handling Tests

### Test Invalid Input
```
1. /play song:invalid_url_here
   Expected: Error message displayed

2. /roast user:@nonexistent
   Expected: Graceful error handling

3. /chat message: (empty)
   Expected: "Please provide a message" error
```

### Test Rate Limiting
```
1. Send /roast 10 times quickly
   Expected: Cooldown message after first use

2. Search for songs repeatedly
   Expected: Rate limit enforced
```

## Success Criteria

✅ **All features working if:**
- [ ] Bot connects to Discord
- [ ] All slash commands registered
- [ ] At least 3 commands working:
  - `/ping` (100% required)
  - `/play` or `/confess` or `/roast`
  - One more feature of your choice
- [ ] No critical errors in logs
- [ ] Health endpoints respond

## Quick Test Script

Run this to test basic functionality:

```powershell
# Quick test commands (run in Discord after bot is online)
# Copy and paste one by one

/ping
/play song:test
/queue
/roast
/profile
```

## Next Steps After Testing

1. **If all tests pass:**
   - Deploy to production (Railway)
   - Monitor logs
   - Set up alerts

2. **If tests fail:**
   - Check error logs
   - Verify API keys
   - Check Discord permissions
   - Review troubleshooting section

3. **For production deployment:**
   - Set environment variables in Railway
   - Configure health checks
   - Set up monitoring
   - Enable auto-scaling

## Support

If you encounter issues:
1. Check logs for error messages
2. Verify all API keys are valid
3. Ensure bot has proper Discord permissions
4. Test health endpoints
5. Check if services are running

---

**Happy Testing! 🎉**
