#!/bin/bash
# Quick fix for NeruBot systemd service execution issue

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_message() {
    echo -e "${2}[NeruBot Fix] $1${NC}"
}

print_success() {
    print_message "$1" "$GREEN"
}

print_warning() {
    print_message "$1" "$YELLOW"
}

print_error() {
    print_message "$1" "$RED"
}

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   print_error "This script must be run as root (use sudo)"
   exit 1
fi

print_success "Fixing NeruBot systemd service..."

# Stop the service
systemctl stop nerubot 2>/dev/null || true

# Update systemd service file
cat > /etc/systemd/system/nerubot.service << 'EOF'
[Unit]
Description=NeruBot Discord Music Bot
After=network.target
Wants=network.target

[Service]
Type=simple
User=nerubot
Group=nerubot
WorkingDirectory=/home/nerubot/nerubot
Environment=PYTHONPATH=/home/nerubot/nerubot
Environment=PYTHONUNBUFFERED=1
ExecStart=/bin/bash -c 'cd /home/nerubot/nerubot && source nerubot_env/bin/activate && python src/main.py'
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=nerubot

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/home/nerubot/nerubot /home/nerubot/logs
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
systemctl daemon-reload

print_success "Systemd service updated!"

# Check if virtual environment exists
if [[ ! -d "/home/nerubot/nerubot/nerubot_env" ]]; then
    print_warning "Virtual environment not found. Creating it now..."
    
    # Switch to nerubot user and create virtual environment
    sudo -u nerubot bash -c '
        cd /home/nerubot/nerubot
        python3 -m venv nerubot_env
        source nerubot_env/bin/activate
        pip install --upgrade pip
        if [[ -f requirements.txt ]]; then
            pip install -r requirements.txt
        fi
    '
    
    print_success "Virtual environment created and dependencies installed!"
fi

# Check if .env file exists
if [[ ! -f "/home/nerubot/nerubot/.env" ]]; then
    print_warning ".env file not found. Creating template..."
    
    sudo -u nerubot bash -c '
        cd /home/nerubot/nerubot
        cat > .env << "EOF"
# Required
DISCORD_TOKEN=your_discord_bot_token_here

# Optional
LOG_LEVEL=INFO
BOT_PREFIX=!
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
EOF
    '
    
    print_warning "Please edit /home/nerubot/nerubot/.env and add your Discord token!"
    print_warning "Use: sudo -u nerubot nano /home/nerubot/nerubot/.env"
fi

print_success "Fix completed! You can now:"
print_success "1. Configure your .env file: sudo -u nerubot nano /home/nerubot/nerubot/.env"
print_success "2. Start the service: sudo systemctl start nerubot"
print_success "3. Check status: sudo systemctl status nerubot"
print_success "4. View logs: sudo journalctl -u nerubot -f"
