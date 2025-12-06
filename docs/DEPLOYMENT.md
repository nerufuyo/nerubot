# NeruBot Deployment Guide

This guide covers various deployment methods for NeruBot, from local development to production environments.

---

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [VPS Deployment](#vps-deployment)
- [Production Best Practices](#production-best-practices)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

**Minimum:**
- CPU: 1 core
- RAM: 512 MB
- Storage: 500 MB
- OS: Linux (Ubuntu 20.04+), macOS, Windows 10+

**Recommended:**
- CPU: 2+ cores
- RAM: 1 GB
- Storage: 2 GB
- OS: Linux (Ubuntu 22.04 LTS)

### Software Requirements

**Required:**
- Go 1.21 or higher
- FFmpeg
- Python 3.8+ (for yt-dlp)
- Git

**Optional:**
- Docker & Docker Compose
- systemd (for service management)
- nginx (for reverse proxy)

---

## Environment Setup

### 1. Create Environment File

```bash
cp .env.example .env
```

### 2. Configure Essential Variables

Edit `.env` and set the following:

```env
# === REQUIRED ===
# Get from: https://discord.com/developers/applications
DISCORD_TOKEN=your_discord_bot_token_here

# === AI CHATBOT (Optional) ===
# Get from: https://platform.deepseek.com/
DEEPSEEK_API_KEY=your_deepseek_api_key

# === FEATURE TOGGLES ===
ENABLE_MUSIC=true
ENABLE_CONFESSION=true
ENABLE_ROAST=true
ENABLE_CHATBOT=false  # Set to true if you have API key

# === LOGGING ===
LOG_LEVEL=INFO  # DEBUG, INFO, WARNING, ERROR
```

### 3. Verify Configuration

```bash
# Check if token is set
grep DISCORD_TOKEN .env

# Validate .env syntax
cat .env | grep -v '^#' | grep -v '^$'
```

---

## Local Development

### Quick Start

```bash
# Install dependencies
go mod download

# Build the bot
go build -o build/nerubot cmd/nerubot/main.go

# Run the bot
./build/nerubot
```

### Using Makefile

```bash
# Build
make build

# Run
make run

# Build and run
make all

# Clean build artifacts
make clean
```

### Development Mode

For faster iteration during development:

```bash
# Run directly without building
go run cmd/nerubot/main.go

# Enable debug logging
LOG_LEVEL=DEBUG go run cmd/nerubot/main.go
```

### Hot Reload (Optional)

Install [air](https://github.com/cosmtrek/air) for hot reloading:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
air init

# Run with hot reload
air
```

---

## Docker Deployment

### Using Docker Compose (Recommended)

**1. Build and start:**
```bash
docker-compose up -d
```

**2. View logs:**
```bash
docker-compose logs -f
```

**3. Stop bot:**
```bash
docker-compose down
```

**4. Rebuild after changes:**
```bash
docker-compose up -d --build
```

### Using Docker Directly

**1. Build image:**
```bash
docker build -t nerubot:latest .
```

**2. Run container:**
```bash
docker run -d \
  --name nerubot \
  --restart unless-stopped \
  --env-file .env \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  nerubot:latest
```

**3. View logs:**
```bash
docker logs -f nerubot
```

**4. Stop container:**
```bash
docker stop nerubot
docker rm nerubot
```

### Docker Management Commands

```bash
# Restart container
docker restart nerubot

# Execute commands inside container
docker exec -it nerubot /bin/sh

# View resource usage
docker stats nerubot

# Inspect container
docker inspect nerubot
```

---

## VPS Deployment

### Automated Setup (Ubuntu/Debian)

**One-command installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/setup.sh | sudo bash
```

This script will:
- ✅ Install Go, FFmpeg, Python, yt-dlp
- ✅ Create dedicated `nerubot` user
- ✅ Clone repository
- ✅ Build bot
- ✅ Setup systemd service
- ✅ Configure firewall
- ✅ Enable auto-start on boot

### Manual Setup

#### 1. Install Dependencies

**Ubuntu/Debian:**
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install FFmpeg
sudo apt install -y ffmpeg

# Install Python and yt-dlp
sudo apt install -y python3 python3-pip
pip3 install yt-dlp

# Verify installations
go version
ffmpeg -version
yt-dlp --version
```

#### 2. Create Service User

```bash
# Create user without login shell
sudo useradd -r -s /bin/false nerubot

# Create directories
sudo mkdir -p /opt/nerubot
sudo chown nerubot:nerubot /opt/nerubot
```

#### 3. Deploy Bot

```bash
# Clone repository
cd /opt/nerubot
sudo -u nerubot git clone https://github.com/nerufuyo/nerubot.git .

# Configure environment
sudo -u nerubot cp .env.example .env
sudo nano .env  # Add your configuration

# Build bot
sudo -u nerubot go build -o build/nerubot cmd/nerubot/main.go

# Set permissions
sudo chmod 600 .env
sudo chown nerubot:nerubot .env
```

#### 4. Create Systemd Service

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

#### 5. Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (auto-start on boot)
sudo systemctl enable nerubot

# Start service
sudo systemctl start nerubot

# Check status
sudo systemctl status nerubot
```

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

# View logs
sudo journalctl -u nerubot -f

# View last 100 lines
sudo journalctl -u nerubot -n 100

# View logs from today
sudo journalctl -u nerubot --since today
```

---

## Production Best Practices

### Security

**1. Protect Secrets:**
```bash
# Restrict .env permissions
chmod 600 .env
chown nerubot:nerubot .env

# Never commit .env to git
echo ".env" >> .gitignore
```

**2. Use Firewall:**
```bash
# Enable UFW
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow custom ports if needed
sudo ufw allow 8080/tcp  # If running health check server

# Check status
sudo ufw status
```

**3. Regular Updates:**
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Update bot
cd /opt/nerubot
sudo -u nerubot git pull
sudo -u nerubot go build -o build/nerubot cmd/nerubot/main.go
sudo systemctl restart nerubot
```

### Logging

**1. Configure Log Rotation:**

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

**2. Monitor Logs:**
```bash
# Real-time logs
tail -f /opt/nerubot/logs/nerubot.log

# Search for errors
grep "ERROR" /opt/nerubot/logs/nerubot.log

# Count errors today
journalctl -u nerubot --since today | grep ERROR | wc -l
```

### Backups

**1. Data Backup Script:**

Create `/opt/nerubot/backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/opt/nerubot/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup data directory
tar -czf $BACKUP_DIR/data_$DATE.tar.gz data/

# Backup .env (encrypted)
gpg --symmetric --cipher-algo AES256 -o $BACKUP_DIR/env_$DATE.gpg .env

# Keep only last 7 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
find $BACKUP_DIR -name "*.gpg" -mtime +7 -delete

echo "Backup completed: $DATE"
```

**2. Schedule Backups:**
```bash
# Add to crontab
sudo crontab -e

# Daily backup at 2 AM
0 2 * * * /opt/nerubot/backup.sh >> /var/log/nerubot-backup.log 2>&1
```

### Resource Limits

**1. Configure systemd limits:**

Edit `/etc/systemd/system/nerubot.service`:

```ini
[Service]
# Memory limit
MemoryLimit=512M

# CPU limit (50% of one core)
CPUQuota=50%

# File descriptor limit
LimitNOFILE=4096
```

**2. Apply changes:**
```bash
sudo systemctl daemon-reload
sudo systemctl restart nerubot
```

---

## Monitoring & Maintenance

### Health Checks

**1. Check Bot Status:**
```bash
# Is service running?
sudo systemctl is-active nerubot

# Check process
ps aux | grep nerubot

# Check resource usage
top -p $(pgrep -f nerubot)
```

**2. Check Bot Connectivity:**
```bash
# View recent logs
sudo journalctl -u nerubot -n 50

# Check for Discord connection
sudo journalctl -u nerubot | grep "Discord bot connected"
```

### Performance Monitoring

**1. Resource Usage:**
```bash
# Memory usage
cat /proc/$(pgrep -f nerubot)/status | grep VmRSS

# CPU usage
top -bn1 -p $(pgrep -f nerubot) | tail -1

# Disk usage
du -sh /opt/nerubot/data
```

**2. Application Metrics:**
```bash
# Count log entries by level
journalctl -u nerubot --since today | grep INFO | wc -l
journalctl -u nerubot --since today | grep ERROR | wc -l

# Check uptime
systemctl show nerubot --property=ActiveEnterTimestamp
```

### Maintenance Tasks

**Daily:**
- ✅ Check service status
- ✅ Review error logs
- ✅ Monitor disk space

**Weekly:**
- ✅ Review all logs
- ✅ Check for updates
- ✅ Verify backups

**Monthly:**
- ✅ Update dependencies
- ✅ Security audit
- ✅ Performance review
- ✅ Clean old logs/backups

---

## Troubleshooting

### Common Issues

#### Bot Won't Start

**Symptoms:**
- Service fails to start
- Exits immediately

**Solutions:**
```bash
# Check logs
sudo journalctl -u nerubot -n 100

# Common issues:
# 1. Invalid Discord token
grep DISCORD_TOKEN .env

# 2. Missing dependencies
which ffmpeg
which yt-dlp

# 3. File permissions
ls -la .env
ls -la data/

# 4. Port conflicts (if using health check)
sudo netstat -tulpn | grep 8080
```

#### Music Not Working

**Symptoms:**
- Songs won't play
- Queue empty immediately

**Solutions:**
```bash
# Check yt-dlp installation
yt-dlp --version

# Test yt-dlp
yt-dlp --skip-download --print-json "https://youtube.com/watch?v=dQw4w9WgXcQ"

# Check FFmpeg
ffmpeg -version

# Check logs for errors
journalctl -u nerubot | grep "music" -i
```

#### High Memory Usage

**Symptoms:**
- OOM (Out of Memory) errors
- Service crashes

**Solutions:**
```bash
# Check current usage
free -h
ps aux | grep nerubot | awk '{print $6}'

# Set memory limits (see Resource Limits section)

# Restart service
sudo systemctl restart nerubot

# Monitor memory
watch -n 5 'ps aux | grep nerubot'
```

#### Permission Errors

**Symptoms:**
- Can't read/write files
- Access denied errors

**Solutions:**
```bash
# Fix ownership
sudo chown -R nerubot:nerubot /opt/nerubot

# Fix permissions
sudo chmod 755 /opt/nerubot
sudo chmod 600 /opt/nerubot/.env
sudo chmod -R 755 /opt/nerubot/data

# Verify
ls -la /opt/nerubot
```

### Debug Mode

Enable detailed logging:

```bash
# Edit .env
LOG_LEVEL=DEBUG

# Restart
sudo systemctl restart nerubot

# Watch logs
sudo journalctl -u nerubot -f
```

### Getting Help

1. **Check Logs:** Always start with logs
2. **Search Issues:** [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)
3. **Ask Community:** [GitHub Discussions](https://github.com/nerufuyo/nerubot/discussions)
4. **Report Bug:** Create detailed issue with logs

---

## Rollback Procedure

If an update causes issues:

```bash
# 1. Stop service
sudo systemctl stop nerubot

# 2. Restore from backup
cd /opt/nerubot
sudo -u nerubot tar -xzf backups/data_YYYYMMDD_HHMMSS.tar.gz

# 3. Checkout previous version
sudo -u nerubot git log --oneline  # Find previous commit
sudo -u nerubot git checkout <commit-hash>

# 4. Rebuild
sudo -u nerubot go build -o build/nerubot cmd/nerubot/main.go

# 5. Start service
sudo systemctl start nerubot

# 6. Verify
sudo systemctl status nerubot
```

---

## Additional Resources

- [Architecture Guide](ARCHITECTURE.md)
- [Project Structure](PROJECT_STRUCTURE.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [GitHub Issues](https://github.com/nerufuyo/nerubot/issues)

---

**Last Updated:** December 6, 2025  
**Version:** 3.0.0  
**Author:** @nerufuyo
