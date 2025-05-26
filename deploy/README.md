# NeruBot VPS Deployment Guide

This guide will help you deploy your NeruBot Discord music bot to a Virtual Private Server (VPS) for 24/7 operation.

## 🚀 Quick Start

### Option 1: One-Command Deployment (Easiest)
```bash
# Complete deployment with domain and SSL
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/quick_deploy.sh | sudo bash -s -- --domain bot.yourdomain.com --ssl

# Basic deployment
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/quick_deploy.sh | sudo bash
```

### Option 2: Traditional VPS Deployment (Recommended)
```bash
# On your VPS (as root)
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/vps_setup.sh | sudo bash
```

### Option 3: Docker Deployment
```bash
# Clone your repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot
chmod +x deploy/docker_setup.sh
./deploy/docker_setup.sh
```

## 📋 Prerequisites

### VPS Requirements
- **OS**: Ubuntu 20.04+ or Debian 11+ (recommended)
- **RAM**: Minimum 1GB, recommended 2GB+
- **Storage**: Minimum 10GB
- **Network**: Stable internet connection

### VPS Providers
Popular and reliable options:
- **DigitalOcean** - $5/month droplet
- **Linode** - $5/month VPS
- **Vultr** - $3.50/month VPS
- **AWS EC2** - t3.micro (free tier)
- **Google Cloud** - e2-micro (free tier)

## 🛠️ Manual Setup Instructions

### Step 1: Connect to Your VPS
```bash
ssh root@your_vps_ip
```

### Step 2: Run the Setup Script
```bash
# Download and run the setup script
wget https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/vps_setup.sh
chmod +x vps_setup.sh
sudo ./vps_setup.sh
```

### Step 3: Deploy Your Bot
```bash
# Switch to bot user
sudo su - nerubot

# Clone your repository
git clone https://github.com/nerufuyo/nerubot.git nerubot
cd nerubot

# Run setup (dependencies only)
./run_nerubot.sh --setup-only

# Configure environment
nano .env
# Add your Discord token and other configuration
```

### Step 4: Start the Service
```bash
# Start the bot service
sudo systemctl start nerubot

# Check status
sudo systemctl status nerubot

# Enable auto-start on boot
sudo systemctl enable nerubot
```

## 🔧 Configuration

### Environment Variables (.env)
```bash
# Required
DISCORD_TOKEN=your_discord_bot_token_here

# Optional
LOG_LEVEL=INFO
BOT_PREFIX=!
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
```

### Discord Bot Token Setup
1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to "Bot" section
4. Create a bot and copy the token
5. Add the token to your `.env` file

## 📊 Monitoring and Management

### Service Management
```bash
# Check bot status
sudo systemctl status nerubot

# View logs
sudo journalctl -u nerubot -f

# Restart bot
sudo systemctl restart nerubot

# Stop bot
sudo systemctl stop nerubot
```

### Using the Monitor Script
```bash
# Quick status overview
/home/nerubot/monitor.sh
```

### Log Management
- Logs are automatically rotated daily
- Last 30 days of logs are kept
- Location: `/home/nerubot/logs/`

## 🛠️ Management Scripts

Your deployment includes several management scripts for easy maintenance:

### Health Monitoring
```bash
# Check bot health
/home/nerubot/nerubot/deploy/scripts/health_check.sh

# View performance metrics
/home/nerubot/nerubot/deploy/scripts/performance_monitor.sh

# Save metrics to file
/home/nerubot/nerubot/deploy/scripts/performance_monitor.sh save
```

### Updates and Maintenance
```bash
# Update bot safely with rollback capability
/home/nerubot/nerubot/deploy/scripts/update.sh

# Check for updates without applying
/home/nerubot/nerubot/deploy/scripts/update.sh --check-only

# Rollback to previous version
/home/nerubot/nerubot/deploy/scripts/update.sh --rollback
```

### Automated Tasks
The deployment automatically sets up cron jobs for:
- Health checks every 5 minutes
- Performance monitoring every 15 minutes  
- Daily backups at 2 AM
- Weekly log cleanup
- Daily update checks

## 🔄 Updates and Deployment

### Automatic Updates
Use the deployment script for easy updates:
```bash
/home/nerubot/deploy.sh
```

### Manual Updates
```bash
# Switch to bot user
sudo su - nerubot
cd nerubot

# Pull latest changes
git pull origin main

# Update dependencies
source nerubot_env/bin/activate
pip install -r requirements.txt

# Restart service
sudo systemctl restart nerubot
```

## 🐳 Docker Deployment

### Prerequisites
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Setup and Run
```bash
# Clone repository
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot

# Setup Docker files
chmod +x deploy/docker_setup.sh
./deploy/docker_setup.sh

# Configure environment
cp .env.docker .env
nano .env  # Add your Discord token

# Build and start
./docker-build.sh
./docker-start.sh
```

### Docker Management
```bash
# View logs
./docker-logs.sh

# Stop bot
./docker-stop.sh

# Update bot
./docker-update.sh
```

## 🔒 Security Best Practices

### Firewall Configuration
The setup script automatically configures UFW:
- SSH (port 22) - allowed
- HTTP (port 80) - allowed
- HTTPS (port 443) - allowed
- All other ports - denied

### Fail2Ban
Automatic protection against brute-force attacks is enabled.

### User Security
- Bot runs as non-root user `nerubot`
- Limited file system access
- Systemd security features enabled

### Additional Recommendations
```bash
# Change default SSH port (optional)
sudo nano /etc/ssh/sshd_config
# Change: Port 22 to Port 2222
sudo systemctl restart ssh

# Disable root SSH login
sudo nano /etc/ssh/sshd_config
# Add: PermitRootLogin no
sudo systemctl restart ssh

# Use SSH keys instead of passwords
ssh-copy-id nerubot@your_vps_ip
```

## 📝 Backup and Recovery

### Automatic Backups
Daily backups are automatically created:
- Location: `/home/nerubot/backups/`
- Retention: 7 days
- Includes: configuration, logs, data

### Manual Backup
```bash
/home/nerubot/backup.sh
```

### Restore from Backup
```bash
# Extract backup
cd /home/nerubot
tar -xzf backups/nerubot_backup_YYYYMMDD_HHMMSS.tar.gz

# Restart service
sudo systemctl restart nerubot
```

## 🎵 Music Dependencies

The setup script automatically installs:
- **FFmpeg** - Audio processing
- **libopus** - Audio encoding for Discord
- **yt-dlp** - YouTube audio extraction

### Troubleshooting Audio Issues
```bash
# Test FFmpeg
ffmpeg -version

# Test Opus
python3 -c "import discord; print(discord.opus.is_loaded())"

# Check bot logs for audio errors
sudo journalctl -u nerubot -f | grep -i audio
```

## 🌐 Domain and SSL (Optional)

### Setup Domain
1. Point your domain to your VPS IP
2. Configure Nginx (already installed):
```bash
sudo nano /etc/nginx/sites-available/nerubot
```

### SSL Certificate
```bash
# Get free SSL certificate
sudo certbot --nginx -d your-domain.com
```

## 🚨 Troubleshooting

### Common Issues

#### Bot Not Starting
```bash
# Check service status
sudo systemctl status nerubot

# Check logs
sudo journalctl -u nerubot -n 50

# Verify configuration
sudo -u nerubot cat /home/nerubot/nerubot/.env
```

#### Permission Errors
```bash
# Fix ownership
sudo chown -R nerubot:nerubot /home/nerubot/nerubot

# Check file permissions
ls -la /home/nerubot/nerubot/
```

#### Audio Not Working
```bash
# Install audio dependencies
sudo apt install ffmpeg libopus0 libopus-dev

# Test in Python
python3 -c "
import discord
print('Opus loaded:', discord.opus.is_loaded())
"
```

#### Memory Issues
```bash
# Check memory usage
free -h
htop

# Restart bot if needed
sudo systemctl restart nerubot
```

### Getting Help
- Check logs: `sudo journalctl -u nerubot -f`
- Monitor script: `/home/nerubot/monitor.sh`
- Discord.py docs: https://discordpy.readthedocs.io/

## 📊 Performance Monitoring

### System Resources
```bash
# CPU and memory usage
htop

# Disk usage
df -h

# Network usage
iftop
```

### Bot Metrics
```bash
# Bot uptime and status
/home/nerubot/monitor.sh

# Recent errors
sudo journalctl -u nerubot --since "1 hour ago" | grep -i error
```

## 🔄 Maintenance Schedule

### Daily
- Automatic log rotation
- Automatic backups
- Service health checks

### Weekly
- Update system packages: `sudo apt update && sudo apt upgrade`
- Check disk space: `df -h`
- Review error logs

### Monthly
- Update bot dependencies
- Review and clean old backups
- Security updates

## 📞 Support

If you encounter issues:
1. Check the troubleshooting section
2. Review bot logs
3. Check Discord.py documentation
4. Open an issue on GitHub

---

**Happy Deploying! 🎵**
