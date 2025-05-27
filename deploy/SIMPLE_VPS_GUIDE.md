# Simple VPS Deployment Guide for NeruBot

This is a simplified guide to deploy your NeruBot on a VPS without monitoring features.

## 🚀 Quick Deployment

### Option 1: Upload Setup Script (Easy Method)
Use the provided upload helper script:

```bash
# On your local machine (from the nerubot directory)
./deploy/upload_to_vps.sh your_vps_ip

# Then follow the instructions displayed by the script
```

### Option 2: Manual Upload
Upload the setup script manually:

```bash
# On your local machine, upload the setup script to your VPS
scp deploy/simple_vps_setup.sh root@your_vps_ip:/tmp/

# On your VPS (as root)
chmod +x /tmp/simple_vps_setup.sh
sudo /tmp/simple_vps_setup.sh
```

### Option 2: Manual Setup (if script not available)
```bash
# Follow the step-by-step guide below
# This will manually set up everything the script would do
```

### Option 3: One-Command Setup (when repository is public)
```bash
# On your VPS (as root) - only works if repository is public
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/simple_vps_setup.sh | sudo bash
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
# If you uploaded the script manually
chmod +x /tmp/simple_vps_setup.sh
sudo /tmp/simple_vps_setup.sh

# OR if repository is public
curl -fsSL https://raw.githubusercontent.com/nerufuyo/nerubot/main/deploy/simple_vps_setup.sh | sudo bash
```

### 3. Upload and Setup Your Bot Code
```bash
# Switch to bot user
sudo su - nerubot

# If repository is public, clone it:
git clone https://github.com/nerufuyo/nerubot.git nerubot

# OR if repository is private/not available, upload your code manually:
# On your local machine:
# tar -czf nerubot.tar.gz --exclude=nerubot_env --exclude=.git --exclude=__pycache__ .
# scp nerubot.tar.gz nerubot@your_vps_ip:/home/nerubot/
# 
# Then on VPS as nerubot user:
# tar -xzf nerubot.tar.gz
# mv nerubot-main nerubot  # adjust if needed
# cd nerubot

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

### Repository Not Found (404 Error)
If you get a 404 error when trying to download the setup script:

1. **Repository is Private**: The GitHub repository may not be public yet
   - **Solution**: Upload files manually using `scp` as shown in Option 1 above

2. **Wrong Branch Name**: The URL might reference `master` instead of `main`
   - **Solution**: Try changing `master` to `main` in the URL

3. **Repository Doesn't Exist**: The repository might not be created yet
   - **Solution**: Use the manual upload method shown above

### Manual File Upload Process
```bash
# On your local machine (from the nerubot directory):
# 1. Upload the setup script
scp deploy/simple_vps_setup.sh root@your_vps_ip:/tmp/

# 2. Create a tarball of your project (excluding unnecessary files)
tar -czf nerubot.tar.gz --exclude=nerubot_env --exclude=.git --exclude=__pycache__ --exclude=*.log .

# 3. Upload the project files
scp nerubot.tar.gz root@your_vps_ip:/tmp/

# 4. On your VPS, run the setup
ssh root@your_vps_ip
chmod +x /tmp/simple_vps_setup.sh
sudo /tmp/simple_vps_setup.sh

# 5. Setup the bot files
sudo su - nerubot
cd /home/nerubot
tar -xzf /tmp/nerubot.tar.gz
mv nerubot-* nerubot  # adjust the directory name if needed
cd nerubot
./run_nerubot.sh --setup-only
```

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
