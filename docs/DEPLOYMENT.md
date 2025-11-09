# Deployment Guide

This guide covers different deployment methods for NeruBot.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Production Deployment](#production-deployment)
- [Systemd Service](#systemd-service)
- [Monitoring and Logging](#monitoring-and-logging)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **CPU:** 1 core minimum, 2 cores recommended
- **RAM:** 512MB minimum, 1GB recommended
- **Storage:** 500MB minimum
- **OS:** Linux (Ubuntu 20.04+), macOS, or Windows with WSL2

### Software Requirements

- Go 1.21 or higher
- FFmpeg
- yt-dlp
- Git
- (Optional) Docker and Docker Compose

## Environment Configuration

### 1. Create Environment File

```bash
cp .env.example .env
```

### 2. Configure Variables

Edit `.env` with your settings:

```env
# Discord Configuration (Required)
DISCORD_TOKEN=your_discord_bot_token_here
DISCORD_GUILD_ID=your_guild_id_here

# AI Providers (At least one required for chatbot)
ANTHROPIC_API_KEY=sk-ant-xxxxx
GEMINI_API_KEY=xxxxx
OPENAI_API_KEY=sk-xxxxx

# Optional Services
WHALE_ALERT_API_KEY=xxxxx

# Feature Flags
ENABLE_MUSIC=true
ENABLE_CONFESSION=true
ENABLE_ROAST=true
ENABLE_CHATBOT=true
ENABLE_NEWS=true
ENABLE_WHALE_ALERT=true

# Logging
LOG_LEVEL=INFO
LOG_FILE=logs/nerubot.log
```

### 3. Get Discord Bot Token

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to "Bot" section
4. Click "Reset Token" and copy it
5. Enable required intents:
   - Server Members Intent
   - Message Content Intent
   - Presence Intent

### 4. Get Guild ID

1. Enable Developer Mode in Discord (User Settings → Advanced)
2. Right-click your server
3. Click "Copy ID"

## Local Development

### Build from Source

```bash
# Clone repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Install dependencies
go mod download

# Build
make build

# Run
./build/nerubot
```

### Development Mode

```bash
# Run without building
go run cmd/nerubot/main.go

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air
```

## Docker Deployment

### Option 1: Docker (Simple)

```bash
# Build image
docker build -t nerubot:latest .

# Run container
docker run -d \
  --name nerubot \
  --restart unless-stopped \
  --env-file .env \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  nerubot:latest
```

### Option 2: Docker Compose (Recommended)

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  nerubot:
    build: .
    container_name: nerubot
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    environment:
      - TZ=America/New_York
```

**Deploy:**
```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild and restart
docker-compose up -d --build
```

### Docker Management

```bash
# View logs
docker logs -f nerubot

# Restart container
docker restart nerubot

# Stop container
docker stop nerubot

# Remove container
docker rm nerubot

# Enter container shell
docker exec -it nerubot sh
```

## Production Deployment

### 1. Server Setup (Ubuntu)

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y git build-essential ffmpeg

# Install yt-dlp
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 2. Deploy Application

```bash
# Create application directory
sudo mkdir -p /opt/nerubot
sudo chown $USER:$USER /opt/nerubot

# Clone and build
cd /opt/nerubot
git clone https://github.com/nerufuyo/nerubot.git .
go mod download
make build

# Create data and log directories
mkdir -p data logs

# Setup environment
cp .env.example .env
nano .env  # Configure your settings
```

### 3. Create Systemd Service

Create `/etc/systemd/system/nerubot.service`:

```ini
[Unit]
Description=NeruBot Discord Bot
After=network.target

[Service]
Type=simple
User=nerubot
WorkingDirectory=/opt/nerubot
ExecStart=/opt/nerubot/build/nerubot
Restart=always
RestartSec=10

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/nerubot/data /opt/nerubot/logs

[Install]
WantedBy=multi-user.target
```

### 4. Setup Service User

```bash
# Create service user
sudo useradd -r -s /bin/false nerubot

# Set permissions
sudo chown -R nerubot:nerubot /opt/nerubot
```

### 5. Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable nerubot

# Start service
sudo systemctl start nerubot

# Check status
sudo systemctl status nerubot
```

## Systemd Service

### Service Management

```bash
# Start
sudo systemctl start nerubot

# Stop
sudo systemctl stop nerubot

# Restart
sudo systemctl restart nerubot

# Status
sudo systemctl status nerubot

# Enable auto-start
sudo systemctl enable nerubot

# Disable auto-start
sudo systemctl disable nerubot

# View logs
sudo journalctl -u nerubot -f

# View last 100 lines
sudo journalctl -u nerubot -n 100
```

### Configuration Files

The systemd service files are located in `deploy/systemd/`:

- `nerubot.service` - Main service file
- Copy to `/etc/systemd/system/`

## Monitoring and Logging

### Log Rotation

Setup logrotate for NeruBot logs.

Create `/etc/logrotate.d/nerubot`:

```
/opt/nerubot/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 nerubot nerubot
    postrotate
        systemctl reload nerubot > /dev/null 2>&1 || true
    endscript
}
```

### View Logs

```bash
# Application logs
tail -f /opt/nerubot/logs/nerubot.log

# Systemd logs
sudo journalctl -u nerubot -f

# Docker logs
docker logs -f nerubot
```

### Monitoring

Monitor bot health:

```bash
# Check if running
systemctl is-active nerubot

# Check resource usage
ps aux | grep nerubot

# Memory usage
cat /proc/$(pgrep nerubot)/status | grep VmRSS
```

## Troubleshooting

### Bot Won't Start

**Check logs:**
```bash
sudo journalctl -u nerubot -n 50
```

**Common issues:**
- Missing Discord token in `.env`
- Invalid bot token
- Missing FFmpeg or yt-dlp
- Permission issues

**Solutions:**
```bash
# Verify environment
cat .env | grep DISCORD_TOKEN

# Check FFmpeg
which ffmpeg

# Check yt-dlp
which yt-dlp

# Fix permissions
sudo chown -R nerubot:nerubot /opt/nerubot
```

### Bot Crashes Repeatedly

**Check memory:**
```bash
free -h
```

**Check disk space:**
```bash
df -h
```

**Increase restart delay in systemd:**
```ini
[Service]
RestartSec=30
```

### Commands Not Working

**Verify bot permissions in Discord:**
- Administrator permission (or specific permissions)
- Slash commands enabled
- Bot has access to channels

**Re-register commands:**
Commands are registered on bot startup. Restart the bot:
```bash
sudo systemctl restart nerubot
```

### Audio Issues

**FFmpeg not found:**
```bash
sudo apt install ffmpeg
```

**yt-dlp not found:**
```bash
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

**Audio stuttering:**
- Increase server resources
- Check network connection
- Verify FFmpeg is working: `ffmpeg -version`

### Database/File Issues

**Permission denied errors:**
```bash
sudo chown -R nerubot:nerubot /opt/nerubot/data
sudo chmod -R 755 /opt/nerubot/data
```

**Corrupted JSON files:**
```bash
# Backup current data
cp -r data data.backup

# Reset data (WARNING: loses data)
rm -rf data/*
```

## Updates and Maintenance

### Update Bot

```bash
# Stop service
sudo systemctl stop nerubot

# Backup data
cp -r /opt/nerubot/data /opt/nerubot/data.backup

# Pull latest code
cd /opt/nerubot
git pull

# Rebuild
make build

# Start service
sudo systemctl start nerubot
```

### Rollback

```bash
# Stop service
sudo systemctl stop nerubot

# Restore backup
cp -r /opt/nerubot/data.backup /opt/nerubot/data

# Checkout previous version
git checkout <previous-commit-hash>
make build

# Start service
sudo systemctl start nerubot
```

## Security Best Practices

1. **Environment Variables**
   - Never commit `.env` to git
   - Use strong, unique tokens
   - Rotate tokens regularly

2. **File Permissions**
   - Restrict access to data directory
   - Run as dedicated user (not root)
   - Use systemd security features

3. **Network**
   - Use firewall rules
   - Only open necessary ports
   - Keep system updated

4. **Monitoring**
   - Enable logging
   - Set up alerts for crashes
   - Monitor resource usage

## Performance Optimization

### Memory

```bash
# Limit memory in systemd
[Service]
MemoryLimit=512M
```

### CPU

```bash
# Limit CPU in systemd
[Service]
CPUQuota=50%
```

### Storage

```bash
# Clean old logs
find /opt/nerubot/logs -name "*.log" -mtime +30 -delete

# Compress old data
tar -czf data-$(date +%Y%m%d).tar.gz data/
```

## Support

For additional help:
- Check [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)
- Read [Architecture Documentation](../ARCHITECTURE.md)
- Review [Contributing Guide](../CONTRIBUTING.md)

---

Last Updated: November 2025
