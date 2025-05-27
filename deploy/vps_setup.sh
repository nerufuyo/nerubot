#!/bin/bash
# NeruBot VPS Deployment Script
# This script sets up a VPS for running NeruBot as a systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="NeruBot"
BOT_USER="nerubot"
BOT_HOME="/home/$BOT_USER"
PROJECT_DIR="$BOT_HOME/nerubot"
SERVICE_NAME="nerubot"
PYTHON_VERSION="3.11"

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

print_step() {
    print_message "$1" "$PURPLE"
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
    print_step "Updating system packages..."
    apt update && apt upgrade -y
    print_success "System updated"
}

# Install required packages
install_packages() {
    print_step "Installing required packages..."
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
        htop \
        ufw \
        fail2ban \
        supervisor \
        nginx \
        certbot \
        python3-certbot-nginx
    print_success "Required packages installed"
}

# Create bot user
create_bot_user() {
    print_step "Creating bot user..."
    if ! id "$BOT_USER" &>/dev/null; then
        useradd -m -s /bin/bash "$BOT_USER"
        print_success "User $BOT_USER created"
    else
        print_info "User $BOT_USER already exists"
    fi
}

# Setup firewall
setup_firewall() {
    print_step "Setting up firewall..."
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow ssh
    ufw allow 80
    ufw allow 443
    ufw --force enable
    print_success "Firewall configured"
}

# Setup fail2ban
setup_fail2ban() {
    print_step "Setting up fail2ban..."
    systemctl enable fail2ban
    systemctl start fail2ban
    print_success "Fail2ban configured"
}

# Create project directory structure
setup_project_structure() {
    print_step "Setting up project structure..."
    sudo -u "$BOT_USER" mkdir -p "$PROJECT_DIR"
    sudo -u "$BOT_USER" mkdir -p "$BOT_HOME/.config"
    sudo -u "$BOT_USER" mkdir -p "$BOT_HOME/logs"
    print_success "Project structure created"
}

# Create systemd service file
create_systemd_service() {
    print_step "Creating systemd service..."
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

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$PROJECT_DIR $BOT_HOME/logs
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
    print_success "Systemd service created and enabled"
}

# Create log rotation configuration
setup_log_rotation() {
    print_step "Setting up log rotation..."
    cat > "/etc/logrotate.d/$SERVICE_NAME" << EOF
$BOT_HOME/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 $BOT_USER $BOT_USER
    postrotate
        systemctl reload $SERVICE_NAME
    endscript
}
EOF
    print_success "Log rotation configured"
}

# Create deployment script for updates
create_deployment_script() {
    print_step "Creating deployment script..."
    cat > "$BOT_HOME/deploy.sh" << 'EOF'
#!/bin/bash
# NeruBot Deployment Script for VPS

set -e

PROJECT_DIR="/home/nerubot/nerubot"
SERVICE_NAME="nerubot"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}[NeruBot] Starting deployment...${NC}"

# Stop the service
echo -e "${YELLOW}[NeruBot] Stopping service...${NC}"
sudo systemctl stop $SERVICE_NAME

# Navigate to project directory
cd $PROJECT_DIR

# Pull latest changes
echo -e "${YELLOW}[NeruBot] Pulling latest changes...${NC}"
git pull origin main

# Activate virtual environment
source nerubot_env/bin/activate

# Update dependencies
echo -e "${YELLOW}[NeruBot] Updating dependencies...${NC}"
pip install --upgrade pip
pip install -r requirements.txt

# Start the service
echo -e "${YELLOW}[NeruBot] Starting service...${NC}"
sudo systemctl start $SERVICE_NAME

# Check status
sleep 3
if sudo systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${GREEN}[NeruBot] Deployment successful! Bot is running.${NC}"
else
    echo -e "\033[0;31m[NeruBot] Deployment failed! Check logs with: sudo journalctl -u $SERVICE_NAME -f${NC}"
    exit 1
fi
EOF

    chmod +x "$BOT_HOME/deploy.sh"
    chown "$BOT_USER:$BOT_USER" "$BOT_HOME/deploy.sh"
    print_success "Deployment script created"
}

# Create monitoring script
create_monitoring_script() {
    print_step "Creating monitoring script..."
    cat > "$BOT_HOME/monitor.sh" << 'EOF'
#!/bin/bash
# NeruBot Monitoring Script

SERVICE_NAME="nerubot"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== NeruBot Status ==="

# Service status
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "Service Status: ${GREEN}RUNNING${NC}"
else
    echo -e "Service Status: ${RED}STOPPED${NC}"
fi

# Uptime
UPTIME=$(systemctl show $SERVICE_NAME --property=ActiveEnterTimestamp --value)
if [ -n "$UPTIME" ] && [ "$UPTIME" != "n/a" ]; then
    echo "Started: $UPTIME"
fi

# Memory usage
PID=$(systemctl show $SERVICE_NAME --property=MainPID --value)
if [ "$PID" != "0" ] && [ -n "$PID" ]; then
    MEM=$(ps -p $PID -o rss= 2>/dev/null | xargs)
    if [ -n "$MEM" ]; then
        MEM_MB=$((MEM / 1024))
        echo "Memory Usage: ${MEM_MB}MB"
    fi
fi

# Recent logs
echo -e "\n${YELLOW}Recent Logs:${NC}"
journalctl -u $SERVICE_NAME --no-pager -n 10 --output=short-iso
EOF

    chmod +x "$BOT_HOME/monitor.sh"
    chown "$BOT_USER:$BOT_USER" "$BOT_HOME/monitor.sh"
    print_success "Monitoring script created"
}

# Create backup script
create_backup_script() {
    print_step "Creating backup script..."
    cat > "$BOT_HOME/backup.sh" << 'EOF'
#!/bin/bash
# NeruBot Backup Script

BACKUP_DIR="/home/nerubot/backups"
PROJECT_DIR="/home/nerubot/nerubot"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup configuration and logs
tar -czf "$BACKUP_DIR/nerubot_backup_$DATE.tar.gz" \
    -C "/home/nerubot" \
    nerubot/.env \
    nerubot/bot.log \
    logs/

echo "Backup created: nerubot_backup_$DATE.tar.gz"

# Keep only last 7 backups
find $BACKUP_DIR -name "nerubot_backup_*.tar.gz" -type f -mtime +7 -delete
EOF

    chmod +x "$BOT_HOME/backup.sh"
    chown "$BOT_USER:$BOT_USER" "$BOT_HOME/backup.sh"
    
    # Add to crontab for daily backups
    (crontab -u $BOT_USER -l 2>/dev/null; echo "0 2 * * * $BOT_HOME/backup.sh") | crontab -u $BOT_USER -
    print_success "Backup script created and scheduled"
}

# Main setup function
main() {
    print_success "=== NeruBot VPS Setup Script ==="
    echo

    check_root
    update_system
    install_packages
    create_bot_user
    setup_firewall
    setup_fail2ban
    setup_project_structure
    create_systemd_service
    setup_log_rotation
    create_deployment_script
    create_monitoring_script
    create_backup_script

    print_success "=== VPS Setup Complete! ==="
    echo
    print_info "Next steps:"
    print_info "1. Switch to bot user: sudo su - $BOT_USER"
    print_info "2. Clone your repository: git clone <your-repo-url> $PROJECT_DIR"
    print_info "3. Run the setup script: cd $PROJECT_DIR && ./run_nerubot.sh --setup-only"
    print_info "4. Configure .env file with your Discord token"
    print_info "5. Start the service: sudo systemctl start $SERVICE_NAME"
    echo
    print_info "Useful commands:"
    print_info "- Check status: sudo systemctl status $SERVICE_NAME"
    print_info "- View logs: sudo journalctl -u $SERVICE_NAME -f"
    print_info "- Monitor bot: /home/$BOT_USER/monitor.sh"
    print_info "- Deploy updates: /home/$BOT_USER/deploy.sh"
}

# Run main function
main "$@"
