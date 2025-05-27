#!/bin/bash

# NeruBot VPS Setup Script
# Automated deployment script for Ubuntu/Debian servers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NERUBOT_USER="nerubot"
NERUBOT_DIR="/opt/nerubot"
SERVICE_NAME="nerubot"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   log_error "This script must be run as root (use sudo)"
   exit 1
fi

log_info "Starting NeruBot VPS deployment..."

# Update system
log_info "Updating system packages..."
apt update && apt upgrade -y

# Install dependencies
log_info "Installing dependencies..."
apt install -y python3 python3-pip python3-venv ffmpeg git curl wget htop

# Create nerubot user
if ! id "$NERUBOT_USER" &>/dev/null; then
    log_info "Creating nerubot user..."
    useradd -r -m -s /bin/bash $NERUBOT_USER
else
    log_warning "User $NERUBOT_USER already exists"
fi

# Create nerubot directory
log_info "Setting up directory structure..."
mkdir -p $NERUBOT_DIR
chown $NERUBOT_USER:$NERUBOT_USER $NERUBOT_DIR

# Clone repository
log_info "Cloning NeruBot repository..."
cd $NERUBOT_DIR
if [ -d ".git" ]; then
    log_warning "Repository already exists, pulling latest changes..."
    sudo -u $NERUBOT_USER git pull
else
    sudo -u $NERUBOT_USER git clone https://github.com/yourusername/nerubot.git .
fi

# Setup Python environment
log_info "Setting up Python virtual environment..."
sudo -u $NERUBOT_USER python3 -m venv venv
sudo -u $NERUBOT_USER ./venv/bin/pip install --upgrade pip
sudo -u $NERUBOT_USER ./venv/bin/pip install -r requirements.txt

# Create .env file template if it doesn't exist
if [ ! -f "$NERUBOT_DIR/.env" ]; then
    log_info "Creating .env template..."
    sudo -u $NERUBOT_USER cat > $NERUBOT_DIR/.env << EOF
# Discord Bot Configuration
DISCORD_TOKEN=your_discord_bot_token_here

# Optional - Spotify Integration
SPOTIFY_CLIENT_ID=
SPOTIFY_CLIENT_SECRET=

# Bot Settings
COMMAND_PREFIX=!
LOG_LEVEL=INFO
MAX_QUEUE_SIZE=100
ENABLE_24_7=true
AUTO_DISCONNECT_TIME=300
EOF
    log_warning "Please edit $NERUBOT_DIR/.env with your Discord bot token"
fi

# Create systemd service
log_info "Creating systemd service..."
cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=NeruBot Discord Music Bot
After=network.target

[Service]
Type=simple
User=$NERUBOT_USER
Group=$NERUBOT_USER
WorkingDirectory=$NERUBOT_DIR
Environment=PATH=$NERUBOT_DIR/venv/bin
ExecStart=$NERUBOT_DIR/venv/bin/python src/main.py
Restart=always
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=$SERVICE_NAME

[Install]
WantedBy=multi-user.target
EOF

# Create log directory
log_info "Setting up logging..."
mkdir -p /var/log/nerubot
chown $NERUBOT_USER:$NERUBOT_USER /var/log/nerubot

# Configure logrotate
cat > /etc/logrotate.d/nerubot << EOF
/var/log/nerubot/*.log {
    daily
    missingok
    rotate 30
    compress
    notifempty
    create 644 $NERUBOT_USER $NERUBOT_USER
    postrotate
        systemctl reload $SERVICE_NAME || true
    endscript
}
EOF

# Enable and start service
log_info "Enabling systemd service..."
systemctl daemon-reload
systemctl enable $SERVICE_NAME

# Setup firewall (UFW)
if command -v ufw &> /dev/null; then
    log_info "Configuring firewall..."
    ufw allow ssh
    ufw --force enable
fi

# Create management scripts
log_info "Creating management scripts..."
cat > $NERUBOT_DIR/manage.sh << 'EOF'
#!/bin/bash

SERVICE_NAME="nerubot"

case "$1" in
    start)
        sudo systemctl start $SERVICE_NAME
        echo "NeruBot started"
        ;;
    stop)
        sudo systemctl stop $SERVICE_NAME
        echo "NeruBot stopped"
        ;;
    restart)
        sudo systemctl restart $SERVICE_NAME
        echo "NeruBot restarted"
        ;;
    status)
        sudo systemctl status $SERVICE_NAME
        ;;
    logs)
        sudo journalctl -u $SERVICE_NAME -f
        ;;
    update)
        sudo systemctl stop $SERVICE_NAME
        git pull
        ./venv/bin/pip install -r requirements.txt
        sudo systemctl start $SERVICE_NAME
        echo "NeruBot updated and restarted"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs|update}"
        exit 1
        ;;
esac
EOF

chmod +x $NERUBOT_DIR/manage.sh
chown $NERUBOT_USER:$NERUBOT_USER $NERUBOT_DIR/manage.sh

# Final setup
log_success "NeruBot VPS deployment completed!"
echo ""
log_info "Next steps:"
echo "1. Edit the Discord token: sudo nano $NERUBOT_DIR/.env"
echo "2. Start the bot: sudo systemctl start $SERVICE_NAME"
echo "3. Check status: sudo systemctl status $SERVICE_NAME"
echo "4. View logs: sudo journalctl -u $SERVICE_NAME -f"
echo ""
log_info "Management commands:"
echo "- Start: $NERUBOT_DIR/manage.sh start"
echo "- Stop: $NERUBOT_DIR/manage.sh stop"
echo "- Restart: $NERUBOT_DIR/manage.sh restart"
echo "- Status: $NERUBOT_DIR/manage.sh status"
echo "- Logs: $NERUBOT_DIR/manage.sh logs"
echo "- Update: $NERUBOT_DIR/manage.sh update"
echo ""
log_warning "Remember to configure your Discord bot token in $NERUBOT_DIR/.env"
