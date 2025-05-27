# Simple VPS Deployment Guide for NeruBot

This is a simplified guide to deploy your NeruBot on a VPS without monitoring features.

## 🚀 Quick Deployment

### One-Command Setup
```bash
# On your VPS (as root)
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/simple_vps_setup.sh | sudo bash
```

### Manual Setup
```bash
# Download the script
wget https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/simple_vps_setup.sh
chmod +x simple_vps_setup.sh
sudo ./simple_vps_setup.sh
```

## 📋 VPS Requirements

- **OS**: Ubuntu 20.04+ or Debian 11+
- **RAM**: Minimum 1GB
- **Storage**: Minimum 5GB
- **Network**: Stable internet connection

## 🛠️ Step-by-Step Setup

### 1. Connect to Your VPS
```bash
ssh root@your_vps_ip
```

### 2. Run the Setup Script
```bash
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/simple_vps_setup.sh | sudo bash
```

### 3. Switch to Bot User and Clone Repository
```bash
# Switch to bot user
sudo su - nerubot

# Clone your repository
git clone https://github.com/nerufuyo/nerubot.git nerubot
cd nerubot

# Setup dependencies
./run_nerubot.sh --setup-only
```

### 4. Configure Your Bot
```bash
# Create environment file
nano .env

# Add your Discord token:
DISCORD_TOKEN=your_discord_bot_token_here
```

### 5. Start the Bot
```bash
# Start the bot service
sudo systemctl start nerubot

# Check if it's running
sudo systemctl status nerubot
```

## 🔧 Bot Management

### Basic Commands
```bash
# Start bot
sudo systemctl start nerubot

# Stop bot
sudo systemctl stop nerubot

# Restart bot
sudo systemctl restart nerubot

# Check status
sudo systemctl status nerubot

# View logs (real-time)
sudo journalctl -u nerubot -f

# View recent logs
sudo journalctl -u nerubot -n 50
```

### Updates
```bash
# Simple update script (created automatically)
/home/nerubot/deploy.sh
```

## 🔑 Getting Your Discord Token

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to "Bot" section
4. Create a bot and copy the token
5. Add the token to your `.env` file

## ⚠️ Troubleshooting

### Bot Won't Start
```bash
# Check detailed status
sudo systemctl status nerubot -l

# Check logs for errors
sudo journalctl -u nerubot -n 50

# Verify files exist
sudo -u nerubot ls -la /home/nerubot/nerubot/

# Test manually
sudo su - nerubot
cd nerubot
source nerubot_env/bin/activate
python src/main.py
```

### Permission Issues
```bash
# Fix ownership
sudo chown -R nerubot:nerubot /home/nerubot/nerubot
```

### Audio Issues
```bash
# Test FFmpeg
ffmpeg -version

# Test in Python
python3 -c "import discord; print('Opus loaded:', discord.opus.is_loaded())"
```

## 🔄 What This Setup Includes

- ✅ Python 3 and required dependencies
- ✅ FFmpeg for audio processing
- ✅ Basic firewall (UFW) with SSH access
- ✅ Systemd service for auto-start
- ✅ Simple deployment script for updates
- ✅ Dedicated `nerubot` user for security

## 🚫 What This Setup Does NOT Include

- ❌ Monitoring scripts
- ❌ Performance metrics
- ❌ Automated backups
- ❌ Health checks
- ❌ Fail2ban protection
- ❌ Nginx web server
- ❌ SSL certificates
- ❌ Log rotation

## 📝 Quick Reference

### File Locations
- Bot files: `/home/nerubot/nerubot/`
- Environment file: `/home/nerubot/nerubot/.env`
- Logs: View with `sudo journalctl -u nerubot -f`
- Service file: `/etc/systemd/system/nerubot.service`

### Useful Commands
```bash
# Switch to bot user
sudo su - nerubot

# Go to bot directory
cd /home/nerubot/nerubot

# Update bot
/home/nerubot/deploy.sh

# Check bot status
sudo systemctl status nerubot
```

---

**That's it! Your bot should now be running 24/7 on your VPS.** 🎵
