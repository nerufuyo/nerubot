# Lavalink Deployment Guide for Railway

## Quick Deploy to Railway

### Step 1: Create New Service for Lavalink

1. Go to your Railway project: https://railway.app/project/bf89d675-5758-4c73-8952-3c6ffc093bfb
2. Click **"+ New"** → **"Empty Service"**
3. Name it: **"lavalink"**

### Step 2: Configure Lavalink Service

1. Click on the new **lavalink** service
2. Go to **Settings** → **Source**
3. Connect to this GitHub repo: `nerufuyo/nerubot`
4. Set **Root Directory**: Leave blank (uses repo root)
5. Set **Build Command**: Leave default
6. Set **Start Command**: Leave default (uses Dockerfile)

### Step 3: Add Lavalink Configuration File

In the lavalink service settings:
1. Go to **Settings** → **Variables**
2. Add these environment variables:

```
SERVER_PORT=2333
LAVALINK_SERVER_PASSWORD=youshallnotpass
```

### Step 4: Deploy Lavalink

1. Go to **Settings** → **Deploy**
2. Click **"Deploy"**
3. Wait 2-3 minutes for build to complete
4. Check **Deployments** tab for success

### Step 5: Get Lavalink Internal URL

Once deployed:
1. Go to **Settings** → **Networking**
2. Copy the **Private Network URL** (e.g., `lavalink.railway.internal`)
3. Note: This URL only works within Railway's network

### Step 6: Update NeruBot Configuration

Run these commands locally:

```bash
# Update bot to use Lavalink
railway service nerubot
railway variables --set "LAVALINK_ENABLED=true"
railway variables --set "LAVALINK_HOST=lavalink.railway.internal"
railway variables --set "LAVALINK_PORT=2333"
railway variables --set "LAVALINK_PASSWORD=youshallnotpass"

# Redeploy bot
railway up
```

---

## Alternative: Manual Railway CLI Setup

If you prefer using Railway CLI to create the service:

```bash
# Link to Railway project
railway link

# The Lavalink service must be created via Railway dashboard
# Then you can configure it with CLI
```

---

## Testing Lavalink

Once deployed, test the connection:

```bash
# From your local machine (won't work - internal only)
curl http://lavalink.railway.internal:2333/version

# Check Railway logs instead
railway logs --service lavalink
```

---

## Files Created for Lavalink

- `Dockerfile.lavalink` - Docker build configuration
- `application.yml` - Lavalink server configuration
- `railway.lavalink.toml` - Railway deployment settings

---

## Troubleshooting

### Lavalink Won't Start
- Check Railway logs: `railway logs --service lavalink`
- Verify environment variables are set
- Ensure port 2333 is available

### Bot Can't Connect to Lavalink
- Verify `LAVALINK_HOST=lavalink.railway.internal` (not localhost)
- Check both services are in same Railway project
- Verify password matches: `youshallnotpass`

### Audio Still Not Playing
- Check bot logs: `railway logs --service nerubot`
- Verify `LAVALINK_ENABLED=true`
- Restart bot after Lavalink is running

---

## Next Steps

After deployment:
1. Test `/play` command in Discord
2. Bot should join voice channel
3. Music should play successfully
4. Bot auto-disconnects after 3 minutes of inactivity

