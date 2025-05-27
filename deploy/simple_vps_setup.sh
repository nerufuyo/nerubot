#!/bin/bash
# Simple NeruBot VPS Setup Script
# Minimal setup without monitoring features

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="NeruBot"
BOT_USER="nerubot"
BOT_HOME="/home/$BOT_USER"
PROJECT_DIR="$BOT_HOME/nerubot"
SERVICE_NAME="nerubot"

# Function to print colored messages
print_message() {
    echo -e "${2}[$PROJECT_NAME] $1${NC}"
}

print_success() {
    print_message "$1" "$GREEN"
}

print_error() {
    print_message "$1" "$RED"
}

print_warning() {
    print_message "$1" "$YELLOW"
}

print_info() {
    print_message "$1" "$BLUE"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Update system packages
update_system() {
    print_info "Updating system packages..."
    apt update && apt upgrade -y
    print_success "System updated"
}

# Install required packages
install_packages() {
    print_info "Installing required packages..."
    apt install -y \
        python3 \
        python3-pip \
        python3-venv \
        python3-dev \
        build-essential \
        ffmpeg \
        libopus0 \
        libopus-dev \
        libffi-dev \
        libnacl-dev \
        git \
        curl \
        ufw
    print_success "Required packages installed"
}

# Create bot user
create_bot_user() {
    print_info "Creating bot user..."
    if ! id "$BOT_USER" &>/dev/null; then
        useradd -m -s /bin/bash "$BOT_USER"
        print_success "User $BOT_USER created"
    else
        print_info "User $BOT_USER already exists"
    fi
}

# Setup basic firewall
setup_firewall() {
    print_info "Setting up basic firewall..."
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow ssh
    ufw --force enable
    print_success "Firewall configured"
}

# Create project directory structure
setup_project_structure() {
    print_info "Setting up project structure..."
    sudo -u "$BOT_USER" mkdir -p "$PROJECT_DIR"
    sudo -u "$BOT_USER" mkdir -p "$BOT_HOME/logs"
    print_success "Project structure created"
}

# Create systemd service file
create_systemd_service() {
    print_info "Creating systemd service..."
    cat > "/etc/systemd/system/$SERVICE_NAME.service" << EOF
[Unit]
Description=NeruBot Discord Music Bot
After=network.target
Wants=network.target

[Service]
Type=simple
User=$BOT_USER
Group=$BOT_USER
WorkingDirectory=$PROJECT_DIR
Environment=PYTHONPATH=$PROJECT_DIR
Environment=PYTHONUNBUFFERED=1
ExecStart=/bin/bash -c 'cd $PROJECT_DIR && source nerubot_env/bin/activate && python src/main.py'
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$SERVICE_NAME

# Basic security settings
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    print_success "Systemd service created and enabled"
}

# Create simple deployment script
create_deployment_script() {
    print_info "Creating deployment script..."
    cat > "$BOT_HOME/deploy.sh" << 'EOF'
#!/bin/bash
# Simple NeruBot Deployment Script

set -e

PROJECT_DIR="/home/nerubot/nerubot"
SERVICE_NAME="nerubot"

echo "[NeruBot] Starting deployment..."

# Stop the service
echo "[NeruBot] Stopping service..."
sudo systemctl stop $SERVICE_NAME

# Navigate to project directory
cd $PROJECT_DIR

# Pull latest changes
echo "[NeruBot] Pulling latest changes..."
git pull origin main

# Activate virtual environment and update dependencies
echo "[NeruBot] Updating dependencies..."
source nerubot_env/bin/activate
pip install --upgrade pip
pip install -r requirements.txt

# Start the service
echo "[NeruBot] Starting service..."
sudo systemctl start $SERVICE_NAME

# Check status
sleep 3
if sudo systemctl is-active --quiet $SERVICE_NAME; then
    echo "[NeruBot] Deployment successful! Bot is running."
else
    echo "[NeruBot] Deployment failed! Check logs with: sudo journalctl -u $SERVICE_NAME -f"
    exit 1
fi
EOF

    chmod +x "$BOT_HOME/deploy.sh"
    chown "$BOT_USER:$BOT_USER" "$BOT_HOME/deploy.sh"
    print_success "Deployment script created"
}

# Main setup function
main() {
    print_success "=== Simple NeruBot VPS Setup ==="
    echo

    check_root
    update_system
    install_packages
    create_bot_user
    setup_firewall
    setup_project_structure
    create_systemd_service
    create_deployment_script

    print_success "=== VPS Setup Complete! ==="
    echo
    print_info "Next steps:"
    print_info "1. Switch to bot user: sudo su - $BOT_USER"
    print_info "2. Clone your repository: git clone https://github.com/nerufuyo/nerubot.git $PROJECT_DIR"
    print_info "3. Setup bot: cd $PROJECT_DIR && ./run_nerubot.sh --setup-only"
    print_info "4. Configure .env file: nano $PROJECT_DIR/.env"
    print_info "5. Add your Discord token to .env file"
    print_info "6. Start the service: sudo systemctl start $SERVICE_NAME"
    echo
    print_info "Useful commands:"
    print_info "- Check status: sudo systemctl status $SERVICE_NAME"
    print_info "- View logs: sudo journalctl -u $SERVICE_NAME -f"
    print_info "- Stop bot: sudo systemctl stop $SERVICE_NAME"
    print_info "- Restart bot: sudo systemctl restart $SERVICE_NAME"
    print_info "- Deploy updates: /home/$BOT_USER/deploy.sh"
}

# Run main function
main "$@"
