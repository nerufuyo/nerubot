# Railway Deployment Guide for NeruBot

## Prerequisites

- Railway CLI installed ✅
- Railway account created
- Discord Bot Token
- (Optional) DeepSeek API Key for AI chatbot

## Step 1: Login to Railway

```bash
railway login
```

This will open your browser to authenticate.

## Step 2: Initialize Railway Project

In your project directory:

```bash
cd c:\Projects\nerubot
railway init
```

Choose:
- Create a new project: **NeruBot**
- Or link to existing project if you already created one on railway.app

## Step 3: Add Environment Variables

Add your environment variables to Railway:

```bash
# Required - Discord Bot Token
railway variables set DISCORD_TOKEN=your_discord_token_here

# Optional - AI Chatbot
railway variables set DEEPSEEK_API_KEY=your_deepseek_api_key

# Optional - Spotify Integration
railway variables set SPOTIFY_CLIENT_ID=your_spotify_id
railway variables set SPOTIFY_CLIENT_SECRET=your_spotify_secret

# Bot Configuration
railway variables set COMMAND_PREFIX=!
railway variables set LOG_LEVEL=INFO

# Feature Toggles
railway variables set ENABLE_MUSIC=true
railway variables set ENABLE_CHATBOT=true
railway variables set ENABLE_CONFESSION=true
railway variables set ENABLE_ROAST=true
railway variables set ENABLE_NEWS=true

# Music Settings
railway variables set MAX_QUEUE_SIZE=100
railway variables set AUTO_DISCONNECT_TIME=300
```

Alternatively, you can set variables via Railway Dashboard:
1. Go to https://railway.app
2. Select your project
3. Go to Variables tab
4. Add all environment variables

## Step 4: Deploy to Railway

```bash
# Deploy the bot
railway up
```

This will:
1. Build the Docker image using Dockerfile.go
2. Push to Railway
3. Deploy and start the bot

## Step 5: View Logs

```bash
# Watch logs in real-time
railway logs
```

## Step 6: Check Deployment Status

```bash
# Check if bot is running
railway status
```

Or visit the Railway dashboard at https://railway.app to see your deployment status.

## Troubleshooting

### Build Fails

If the build fails, check:
- Go version in Dockerfile matches your go.mod
- All dependencies are available
- Dockerfile.go has correct syntax

### Bot Won't Start

Check logs:
```bash
railway logs
```

Common issues:
- Missing DISCORD_TOKEN
- Invalid token
- Network connectivity issues

### Bot Crashes

Check Railway logs for error messages:
```bash
railway logs --follow
```

## Managing Your Deployment

### Update Environment Variables
```bash
railway variables set VARIABLE_NAME=new_value
```

### Redeploy After Code Changes
```bash
railway up
```

### Open Railway Dashboard
```bash
railway open
```

### Link to Different Project
```bash
railway link
```

## Production Tips

1. **Always set required environment variables** before deploying
2. **Monitor logs** regularly for errors
3. **Use Railway's built-in metrics** to track performance
4. **Set up alerts** in Railway dashboard for downtime
5. **Keep your dependencies updated** in go.mod

## Cost Optimization

Railway offers:
- $5/month free tier (500 hours)
- Good for hobby projects
- Scales based on usage

For 24/7 operation:
- Consider Railway's Hobby plan ($5/month)
- Or upgrade to Pro plan for larger servers

## Next Steps

After deployment:
1. Verify bot is online in Discord
2. Test slash commands
3. Monitor resource usage in Railway dashboard
4. Set up error notifications
5. Configure automatic backups if needed

## Support

- Railway Docs: https://docs.railway.app
- NeruBot Issues: https://github.com/nerufuyo/nerubot/issues
- Railway Discord: https://discord.gg/railway
