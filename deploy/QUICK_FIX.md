# Quick Fix for VPS Deployment Issue

## Problem
The command `wget https://raw.githubusercontent.com/nerufuyo/nerubot/master/deploy/simple_vps_setup.sh` returns a 404 error because the repository is either not public or doesn't exist at that URL.

## Solution

### Option 1: Use the Upload Helper Script (Recommended)
From your local machine in the nerubot directory:

```bash
./deploy/upload_to_vps.sh your_vps_ip
```

This will:
1. Upload the setup script to your VPS
2. Create and upload a tarball of your project
3. Give you step-by-step instructions to complete the setup

### Option 2: Manual Upload
```bash
# Upload setup script
scp deploy/simple_vps_setup.sh root@your_vps_ip:/tmp/

# Create project tarball
tar -czf nerubot.tar.gz --exclude=nerubot_env --exclude=.git --exclude=__pycache__ .

# Upload project files
scp nerubot.tar.gz root@your_vps_ip:/tmp/
```

### On Your VPS
```bash
# Run the setup script
chmod +x /tmp/simple_vps_setup.sh
sudo /tmp/simple_vps_setup.sh

# Setup the bot
sudo su - nerubot
cd /home/nerubot
tar -xzf /tmp/nerubot.tar.gz
cd nerubot
./run_nerubot.sh --setup-only

# Create .env file
nano .env
# Add: DISCORD_TOKEN=your_token_here

# Start the bot
exit  # exit from nerubot user
sudo systemctl start nerubot
sudo systemctl status nerubot
```

## Files Updated
- ✅ `deploy/SIMPLE_VPS_GUIDE.md` - Updated with alternative deployment methods
- ✅ `deploy/upload_to_vps.sh` - New helper script for easy uploading

The guide now provides multiple deployment options that don't rely on a public GitHub repository.
