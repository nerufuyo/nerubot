# NeruBot Deployment Guide

This directory contains deployment scripts and configurations for running NeruBot in production environments.

## 🚀 Quick Deployment

### VPS Deployment (Ubuntu/Debian)
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/nerubot/main/deploy/vps_setup.sh | sudo bash
```

### Docker Deployment
```bash
docker build -t nerubot .
docker run -d --name nerubot --env-file .env nerubot
```

## 📁 Files Overview

| File | Description |
|------|-------------|
| `vps_setup.sh` | Automated VPS deployment script |
| `nerubot.service` | Systemd service configuration |
| `docker-compose.yml` | Docker Compose configuration |
| `monitoring.sh` | Health monitoring script |

## 🛠️ Manual VPS Setup

### 1. System Prerequisites
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y python3 python3-pip python3-venv ffmpeg git
```

### 2. Create User and Directory
```bash
# Create dedicated user
sudo useradd -r -m -s /bin/bash nerubot

# Create application directory
sudo mkdir -p /opt/nerubot
sudo chown nerubot:nerubot /opt/nerubot
```

### 3. Setup Application
```bash
# Clone repository
cd /opt/nerubot
sudo -u nerubot git clone https://github.com/yourusername/nerubot.git .

# Setup virtual environment
sudo -u nerubot python3 -m venv venv
sudo -u nerubot ./venv/bin/pip install -r requirements.txt
```

### 4. Configure Environment
```bash
# Copy environment template
sudo -u nerubot cp .env.example .env

# Edit configuration
sudo -u nerubot nano .env
```

### 5. Install Systemd Service
```bash
# Copy service file
sudo cp deploy/nerubot.service /etc/systemd/system/

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable nerubot
sudo systemctl start nerubot
```

## 🐳 Docker Setup

### Using Docker Compose
```yaml
# Create docker-compose.yml
version: '3.8'
services:
  nerubot:
    build: .
    container_name: nerubot
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
```

### Using Docker CLI
```bash
# Build image
docker build -t nerubot .

# Run container
docker run -d \
  --name nerubot \
  --restart unless-stopped \
  --env-file .env \
  -v $(pwd)/logs:/app/logs \
  nerubot
```

## 📊 Monitoring

### Service Status
```bash
# Check service status
sudo systemctl status nerubot

# View logs
sudo journalctl -u nerubot -f

# Check resource usage
sudo systemctl show nerubot --property=MemoryCurrent,CPUUsageNSec
```

### Log Management
```bash
# View recent logs
tail -f /var/log/nerubot/bot.log

# Search error logs
grep -E "(ERROR|CRITICAL)" /var/log/nerubot/bot.log

# Rotate logs manually
sudo logrotate -f /etc/logrotate.d/nerubot
```

## 🔧 Management Commands

### Service Control
```bash
# Start service
sudo systemctl start nerubot

# Stop service
sudo systemctl stop nerubot

# Restart service
sudo systemctl restart nerubot

# Reload configuration
sudo systemctl reload nerubot
```

### Application Updates
```bash
# Stop service
sudo systemctl stop nerubot

# Update code
cd /opt/nerubot
sudo -u nerubot git pull

# Update dependencies
sudo -u nerubot ./venv/bin/pip install -r requirements.txt

# Start service
sudo systemctl start nerubot
```

## 🔒 Security

### Firewall Configuration
```bash
# Install UFW
sudo apt install ufw

# Allow SSH
sudo ufw allow ssh

# Enable firewall
sudo ufw enable
```

### Service Security
The systemd service includes security hardening:
- Runs as non-root user
- Private temporary directory
- Protected home directory
- Read-only system directories
- No new privileges

### Regular Maintenance
```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Clean old logs
sudo find /var/log/nerubot -name "*.log.*" -mtime +30 -delete

# Monitor disk usage
df -h
```

## 🚨 Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| Service fails to start | Check logs with `journalctl -u nerubot` |
| FFmpeg not found | Install with `sudo apt install ffmpeg` |
| Permission denied | Check file ownership and permissions |
| Bot not responding | Verify Discord token in `.env` |

### Debug Mode
```bash
# Stop service
sudo systemctl stop nerubot

# Run manually with debug
cd /opt/nerubot
sudo -u nerubot ./venv/bin/python src/main.py --debug
```

### Health Check
```bash
# Check if bot is online
./deploy/monitoring.sh

# Verify Discord connection
curl -H "Authorization: Bot YOUR_TOKEN" \
  https://discord.com/api/v10/users/@me
```

## 📈 Performance Tuning

### Memory Optimization
```bash
# Set memory limits in systemd service
echo "MemoryMax=512M" >> /etc/systemd/system/nerubot.service
sudo systemctl daemon-reload
```

### Audio Quality
```bash
# Optimize FFmpeg settings in config
echo "FFMPEG_OPTIONS=-vn -f opus" >> .env
```

### Queue Performance
```bash
# Limit queue size
echo "MAX_QUEUE_SIZE=50" >> .env
```
