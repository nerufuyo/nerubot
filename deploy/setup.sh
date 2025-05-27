#!/bin/bash
# NeruBot VPS Setup - Simple and efficient
set -e

# Colors
G='\033[0;32m'; R='\033[0;31m'; Y='\033[1;33m'; B='\033[0;34m'; NC='\033[0m'
log() { echo -e "${B}[Setup]${NC} $1"; }
success() { echo -e "${G}✓${NC} $1"; }
error() { echo -e "${R}✗${NC} $1"; exit 1; }

# Configuration
BOT_USER="nerubot"
BOT_HOME="/home/$BOT_USER"
PROJECT_DIR="$BOT_HOME/nerubot"
SERVICE_NAME="nerubot"

# Check root
[[ $EUID -eq 0 ]] || error "Run as root: sudo $0"

log "Setting up NeruBot VPS environment..."

# Update system
log "Updating system packages..."
apt update -y && apt install -y python3 python3-pip python3-venv git ffmpeg curl ufw

# Create user
if ! id "$BOT_USER" &>/dev/null; then
    useradd -m -s /bin/bash "$BOT_USER"
    success "User $BOT_USER created"
fi

# Setup firewall
ufw --force reset && ufw default deny incoming && ufw default allow outgoing
ufw allow ssh && ufw --force enable
success "Firewall configured"

# Create systemd service
cat > "/etc/systemd/system/$SERVICE_NAME.service" << EOF
[Unit]
Description=NeruBot Discord Music Bot
After=network.target

[Service]
Type=simple
User=$BOT_USER
WorkingDirectory=$PROJECT_DIR
Environment=PYTHONPATH=$PROJECT_DIR
ExecStart=/bin/bash -c 'source nerubot_env/bin/activate && python src/main.py'
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload && systemctl enable "$SERVICE_NAME"
success "Service created and enabled"

# Create directories
sudo -u "$BOT_USER" mkdir -p "$BOT_HOME/logs" "$BOT_HOME/backups"

# Create deployment script
cat > "$BOT_HOME/deploy.sh" << 'EOF'
#!/bin/bash
# Simple deployment script
set -e
cd /home/nerubot/nerubot
sudo systemctl stop nerubot
git pull origin main
source nerubot_env/bin/activate
pip install -r requirements.txt -q
sudo systemctl start nerubot
sleep 3
if sudo systemctl is-active --quiet nerubot; then
    echo "✓ Deployment successful"
else
    echo "✗ Deployment failed"
    exit 1
fi
EOF

chmod +x "$BOT_HOME/deploy.sh"
chown "$BOT_USER:$BOT_USER" "$BOT_HOME/deploy.sh"

success "VPS setup complete!"
echo
log "Next steps:"
echo "1. sudo su - $BOT_USER"
echo "2. git clone <your-repo-url> $PROJECT_DIR"
echo "3. cd $PROJECT_DIR && ./run.sh"
echo "4. Edit .env file with your Discord token"
echo "5. sudo systemctl start $SERVICE_NAME"
