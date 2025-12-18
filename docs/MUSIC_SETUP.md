# Music Playback Setup Guide

## Overview

NeruBot uses **Lavalink**, a standalone audio server, to handle Discord voice audio streaming. This is necessary because the Go Discord library (`discordgo`) doesn't have native Opus codec support for voice channels.

## Why Lavalink?

- **Reliable**: Widely used in production Discord bots
- **Feature-Rich**: Supports playlists, filters, advanced search
- **Scalable**: Can handle multiple voice connections
- **Open Source**: Free to use and deploy

## Quick Start with Docker (Recommended)

The easiest way to get music working is using Docker Compose:

```bash
# Start both Lavalink and the bot with music support
docker-compose -f docker-compose.music.yml up -d

# View logs
docker-compose -f docker-compose.music.yml logs -f

# Stop
docker-compose -f docker-compose.music.yml down
```

This will:
1. Start a Lavalink server on `localhost:2333`
2. Start the bot with Lavalink enabled
3. Set up all required environment variables automatically

## Manual Setup

### Option 1: Docker (Lavalink only)

```bash
docker run -d \
  --name lavalink \
  -e _JAVA_OPTIONS="-Xmx6G" \
  -e SERVER_PORT=2333 \
  -e LAVALINK_SERVER_PASSWORD=youshallnotpass \
  -p 2333:2333 \
  ghcr.io/lavalink-devs/lavalink:latest
```

### Option 2: Direct Installation (Linux/macOS)

1. **Install Java**
   ```bash
   # Ubuntu/Debian
   sudo apt install openjdk-17-jre-headless
   
   # macOS
   brew install openjdk@17
   ```

2. **Download Lavalink**
   ```bash
   # Download latest release
   curl -L https://github.com/lavalink-devs/Lavalink/releases/download/v4.0.0/Lavalink.jar -o Lavalink.jar
   ```

3. **Create application.yml**
   ```yaml
   server:
     port: 2333
   
   lavalink:
     server:
       password: "youshallnotpass"
       sources:
         http:
           playlist:
             enabled: true
       filters:
         enabled: true
   ```

4. **Run Lavalink**
   ```bash
   java -Xmx6G -jar Lavalink.jar
   ```

### Option 3: Windows Setup

1. Download Java 17+ from [oracle.com](https://www.oracle.com/java/technologies/downloads/)
2. Download Lavalink.jar from [GitHub releases](https://github.com/lavalink-devs/Lavalink/releases)
3. Create `application.yml` in the same directory
4. Run: `java -Xmx6G -jar Lavalink.jar`

## Configuration

### Environment Variables

```env
# Enable Lavalink music support
LAVALINK_ENABLED=true

# Lavalink server connection details
LAVALINK_HOST=localhost
LAVALINK_PORT=2333
LAVALINK_PASSWORD=youshallnotpass
```

### Railway Deployment

For Railway, set these environment variables in your Railway dashboard:

```
LAVALINK_ENABLED=true
LAVALINK_HOST=lavalink.railway.internal
LAVALINK_PORT=2333
LAVALINK_PASSWORD=youshallnotpass
```

Then deploy a separate Lavalink service on Railway and link them in the same project.

## Testing

Once Lavalink is running and the bot is connected:

1. **Join a voice channel** in Discord
2. **Use the music commands**:
   ```
   /play <song name or URL>
   /skip
   /pause
   /resume
   /stop
   /queue
   /nowplaying
   ```

3. **Check bot logs** for connection confirmations:
   ```
   Lavalink connected to: localhost:2333
   Playing audio successfully
   ```

## Troubleshooting

### Bot can't connect to Lavalink

**Error**: `Failed to connect to Lavalink`

**Solutions**:
- Verify Lavalink is running: `curl http://localhost:2333/loadbalance`
- Check firewall allows port 2333
- Verify `LAVALINK_HOST` and `LAVALINK_PORT` are correct
- Ensure same network if using Docker

### Music commands work but no audio

**Error**: Bot joins voice but produces no sound

**Solutions**:
- Check bot has permission to speak in voice channel
- Verify Lavalink is actually connected
- Try different songs (some may be region-blocked)
- Check bot volume isn't set to 0

### Lavalink crashes on startup

**Error**: `OutOfMemoryError` or crash

**Solutions**:
- Increase Java heap: `-Xmx8G` instead of `-Xmx6G`
- Check system has enough free RAM
- Run on a machine with at least 2GB RAM

### Lavalink connection refused

**Error**: `Connection refused` on port 2333

**Solutions**:
- Is Lavalink running? Check: `netstat -ln | grep 2333`
- Change port in `application.yml` if 2333 is in use
- Check firewall rules: `sudo ufw allow 2333`

## Advanced Configuration

### Lavalink Filters

Edit `application.yml` to enable audio filters:

```yaml
lavalink:
  server:
    filters:
      enabled: true
      karaoke:
        enabled: true
      timescale:
        enabled: true
      distortion:
        enabled: true
      equalizer:
        enabled: true
```

### Performance Tuning

```yaml
lavalink:
  server:
    buffer:
      duration: 400
    frameBufferDuration: 5000
    trackStuckThresholdMs: 10000
```

### Multiple Nodes (Advanced)

For high load, run multiple Lavalink nodes:

```bash
# Node 1
java -Xmx6G -jar Lavalink.jar --server.port 2333

# Node 2
java -Xmx6G -jar Lavalink.jar --server.port 2334

# Node 3
java -Xmx6G -jar Lavalink.jar --server.port 2335
```

## Resources

- **Lavalink GitHub**: https://github.com/lavalink-devs/Lavalink
- **Lavalink Docs**: https://lavalink.dev/
- **Docker Image**: `ghcr.io/lavalink-devs/lavalink:latest`
- **Community Support**: https://discord.gg/lavalink

## Performance Requirements

| Component | Minimum | Recommended |
|-----------|---------|------------|
| CPU | 1 core | 2+ cores |
| RAM | 2 GB | 4-8 GB |
| Network | 1 Mbps | 10+ Mbps |
| Java | 11+ | 17+ LTS |

## Next Steps

1. Deploy Lavalink (Docker recommended)
2. Update environment variables
3. Restart bot
4. Test music commands in Discord
5. Enjoy your music bot! 🎵
