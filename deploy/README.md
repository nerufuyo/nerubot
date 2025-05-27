# NeruBot Deployment Guide

Simple and efficient deployment options for NeruBot.

## 🚀 Quick Deployment

### VPS Setup (Recommended)
```bash
# On your Ubuntu/Debian VPS (as root)
curl -fsSL https://raw.githubusercontent.com/your-username/nerubot/main/deploy/setup.sh | sudo bash
```

This script will:
- Install Python 3, FFmpeg, and dependencies
- Create `nerubot` user and secure environment
- Setup systemd service
- Configure firewall

### Manual Setup
```bash
# 1. Clone repository
sudo su - nerubot
git clone https://github.com/your-username/nerubot.git
cd nerubot

# 2. Setup and configure
./run.sh  # Follow prompts to configure .env

# 3. Start service
sudo systemctl start nerubot
sudo systemctl enable nerubot
```

## 🔧 Management

### Service Control
```bash
# Check status
sudo systemctl status nerubot

# View logs
sudo journalctl -u nerubot -f

# Restart service
sudo systemctl restart nerubot
```

### Monitoring
```bash
# Quick status dashboard
./deploy/status.sh

# Health check
./deploy/monitor.sh

# Update bot
cd /home/nerubot/nerubot
./deploy/update.sh
```

## 📊 Server Requirements

**Minimum:**
- 1 CPU core
- 1GB RAM
- 5GB storage
- Ubuntu 20.04+ or Debian 11+

**Recommended:**
- 2 CPU cores
- 2GB RAM
- 10GB storage

## 🛠️ Configuration

### Environment Variables
```bash
# Required
DISCORD_TOKEN=your_bot_token_here

# Optional
LOG_LEVEL=INFO
COMMAND_PREFIX=!
SPOTIFY_CLIENT_ID=optional_spotify_id
SPOTIFY_CLIENT_SECRET=optional_spotify_secret
```

### Security
The setup script automatically:
- Creates non-root user for the bot
- Configures UFW firewall (SSH only)
- Sets up proper file permissions
- Enables service auto-restart

## 🚨 Troubleshooting

### Common Issues

**Bot won't start:**
```bash
# Check logs
sudo journalctl -u nerubot -n 50

# Verify token
sudo -u nerubot cat /home/nerubot/nerubot/.env
```

**Audio issues:**
```bash
# Check FFmpeg
ffmpeg -version

# Reinstall if needed
sudo apt update && sudo apt install ffmpeg
```

**Permission errors:**
```bash
# Fix ownership
sudo chown -R nerubot:nerubot /home/nerubot/nerubot
```

### Service Recovery
```bash
# If service fails repeatedly
sudo systemctl stop nerubot
cd /home/nerubot/nerubot
./run.sh  # Test manually first
sudo systemctl start nerubot
```

## 🔄 Updates

### Automatic Updates (Recommended)
```bash
# Setup daily updates (run as nerubot user)
echo "0 4 * * * cd /home/nerubot/nerubot && ./deploy/update.sh" | crontab -
```

### Manual Updates
```bash
cd /home/nerubot/nerubot
git pull origin main
pip install -r requirements.txt
sudo systemctl restart nerubot
```

## 🐳 Docker Alternative

If you prefer Docker:
```bash
git clone https://github.com/your-username/nerubot.git
cd nerubot
docker-compose up -d
```

## 📝 Notes

- Bot runs as `nerubot` user for security
- Logs are available via `journalctl -u nerubot`
- Service auto-restarts on failure
- Uses systemd for process management
- Firewall blocks all ports except SSH

---

**Need support?** Open an issue or check the troubleshooting section above.
