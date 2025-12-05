# NeruBot Quick Start

## 🚀 Setup (First Time Only)

### 1. Get Discord Bot Token
1. Go to https://discord.com/developers/applications
2. Create new application → Bot section
3. Copy the bot token

### 2. Configure Environment
```powershell
# Edit .env file and set:
DISCORD_TOKEN=your_actual_token_here
DEEPSEEK_API_KEY=your_deepseek_key_here  # Optional, for chatbot
```

### 3. Invite Bot to Server
1. Discord Developer Portal → OAuth2 → URL Generator
2. Select: `bot` + `applications.commands`
3. Permissions: Send Messages, Embed Links, Connect, Speak, Use Slash Commands
4. Copy URL and open in browser

## 🎮 Running & Testing

### Option 1: Quick Test (Recommended)
```powershell
.\test-services.ps1
# Choose option 1 - Run main bot
```

### Option 2: Manual Run
```powershell
go build -o build/bot.exe ./cmd/nerubot
.\build\bot.exe
```

### Option 3: Microservices
```powershell
.\test-services.ps1
# Choose option 2 - Health check test
```

## 📝 Discord Commands to Test

Once bot is online, try these in Discord:

```
/ping                          → Bot status
/play song:Never Gonna Give You Up  → Play music
/queue                         → Show queue
/nowplaying                    → Current song
/pause                         → Pause music
/resume                        → Resume music
/skip                          → Skip song
/stop                          → Stop & clear

/roast                         → Get roasted
/roast user:@someone           → Roast someone
/profile                       → View activity stats

/confess message:your confession  → Submit anonymous confession

/chat message:Hello            → AI chatbot (requires DEEPSEEK_API_KEY)
```

## ✅ Success Checklist

- [ ] Bot appears online in Discord
- [ ] Slash commands appear when typing `/`
- [ ] `/ping` responds with "Pong!"
- [ ] At least one feature works (music, roast, or confess)

## 🔧 Troubleshooting

### Bot doesn't start
- Check DISCORD_TOKEN in .env
- Verify token is valid (not expired)

### Slash commands don't appear
- Wait 1-2 minutes after bot starts
- Try in different channel
- Reinvite bot with correct permissions

### Music doesn't play
- Join voice channel before using /play
- Install FFmpeg: `winget install FFmpeg`
- Install yt-dlp: `pip install yt-dlp`

### Chatbot doesn't work
- Add DEEPSEEK_API_KEY to .env
- Get key from https://platform.deepseek.com/

## 📊 Expected Output

When bot starts successfully:
```
[INFO] === NeruBot v3.0.0 (Golang Edition) ===
[INFO] Starting Discord bot...
[INFO] Configuration loaded successfully
[INFO] Features enabled
[INFO] Bot is running. Press CTRL+C to exit.
```

## 🎯 Quick Test Sequence

1. Start bot: `.\test-services.ps1` → Option 1
2. Wait for "Bot is ready" message
3. In Discord: `/ping`
4. Should respond: "🏓 Pong! API Gateway is running."
5. Test more commands from list above

## 📚 More Info

- Full testing guide: `TESTING_GUIDE.md`
- Architecture docs: `projects/MICROSERVICES_COMPLETE.md`
- Deployment guide: `docs/RAILWAY_DEPLOYMENT.md`

## 🆘 Need Help?

1. Check logs for error messages
2. Verify .env configuration
3. Test health endpoints
4. Review TESTING_GUIDE.md for detailed procedures
