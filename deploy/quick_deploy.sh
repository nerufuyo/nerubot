#!/bin/bash
# NeruBot Quick Deploy Script
# One-command deployment for VPS

set -e

# Configuration
REPO_URL="${REPO_URL:-https://github.com/your-username/nerubot.git}"
BRANCH="${BRANCH:-main}"
DOMAIN="${DOMAIN:-}"
ENABLE_SSL="${ENABLE_SSL:-false}"
INSTALL_NGINX="${INSTALL_NGINX:-false}"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_step() {
    echo -e "${BLUE}[STEP] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root"
        print_step "Usage: curl -fsSL https://raw.githubusercontent.com/your-username/nerubot/main/deploy/quick_deploy.sh | sudo bash"
        exit 1
    fi
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --repo)
                REPO_URL="$2"
                shift 2
                ;;
            --branch)
                BRANCH="$2"
                shift 2
                ;;
            --domain)
                DOMAIN="$2"
                INSTALL_NGINX="true"
                shift 2
                ;;
            --ssl)
                ENABLE_SSL="true"
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    echo "NeruBot Quick Deploy Script"
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  --repo URL        Git repository URL"
    echo "  --branch NAME     Git branch to deploy (default: main)"
    echo "  --domain NAME     Domain name (enables Nginx)"
    echo "  --ssl             Enable SSL with Let's Encrypt"
    echo "  --help            Show this help"
    echo
    echo "Environment variables:"
    echo "  DISCORD_TOKEN     Discord bot token"
    echo "  REPO_URL          Git repository URL"
    echo
    echo "Example:"
    echo "  $0 --repo https://github.com/user/nerubot.git --domain bot.example.com --ssl"
}

# Main deployment
deploy_nerubot() {
    print_step "Starting NeruBot quick deployment..."
    
    # Run main setup script
    if command -v wget > /dev/null; then
        wget -O vps_setup.sh https://raw.githubusercontent.com/your-username/nerubot/main/deploy/vps_setup.sh
    else
        curl -fsSL https://raw.githubusercontent.com/your-username/nerubot/main/deploy/vps_setup.sh -o vps_setup.sh
    fi
    
    chmod +x vps_setup.sh
    ./vps_setup.sh
    
    # Switch to bot user and clone repository
    print_step "Cloning repository..."
    sudo -u nerubot git clone "$REPO_URL" /home/nerubot/nerubot
    cd /home/nerubot/nerubot
    
    if [ "$BRANCH" != "main" ]; then
        sudo -u nerubot git checkout "$BRANCH"
    fi
    
    # Setup Python environment
    print_step "Setting up Python environment..."
    sudo -u nerubot bash -c "cd /home/nerubot/nerubot && ./run_nerubot.sh --setup-only"
    
    # Configure environment
    print_step "Configuring environment..."
    if [ -n "$DISCORD_TOKEN" ]; then
        sudo -u nerubot bash -c "cd /home/nerubot/nerubot && echo 'DISCORD_TOKEN=$DISCORD_TOKEN' > .env"
        sudo -u nerubot bash -c "cd /home/nerubot/nerubot && echo 'LOG_LEVEL=INFO' >> .env"
    else
        print_warning "DISCORD_TOKEN not provided. You'll need to configure it manually."
        sudo -u nerubot cp /home/nerubot/nerubot/.env.production /home/nerubot/nerubot/.env
    fi
    
    # Setup Nginx if domain provided
    if [ "$INSTALL_NGINX" = "true" ] && [ -n "$DOMAIN" ]; then
        setup_nginx
    fi
    
    # Start the service
    print_step "Starting NeruBot service..."
    systemctl start nerubot
    systemctl enable nerubot
    
    # Verify deployment
    sleep 10
    if systemctl is-active --quiet nerubot; then
        print_success "NeruBot deployed successfully!"
    else
        print_error "Deployment failed. Check logs with: journalctl -u nerubot -f"
        exit 1
    fi
    
    # Show final information
    show_deployment_info
}

# Setup Nginx
setup_nginx() {
    print_step "Setting up Nginx for domain: $DOMAIN"
    
    # Copy Nginx configuration
    cp /home/nerubot/nerubot/deploy/nginx/nerubot.conf /etc/nginx/sites-available/nerubot
    sed -i "s/your-domain.com/$DOMAIN/g" /etc/nginx/sites-available/nerubot
    
    # Enable site
    ln -sf /etc/nginx/sites-available/nerubot /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default
    
    # Create web directory
    mkdir -p /var/www/nerubot
    cat > /var/www/nerubot/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>NeruBot Status</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .status { padding: 15px; border-radius: 5px; margin: 20px 0; }
        .online { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .offline { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        h1 { color: #333; text-align: center; }
        .info { color: #666; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🤖 NeruBot Status</h1>
        <div class="status online">
            <strong>✅ Bot is Online</strong><br>
            NeruBot is running and ready to play music!
        </div>
        <div class="info">
            <h3>🎵 Features</h3>
            <ul>
                <li>Play music from YouTube, Spotify, and SoundCloud</li>
                <li>Advanced queue management</li>
                <li>24/7 mode support</li>
                <li>Interactive help system</li>
            </ul>
            <h3>🚀 Commands</h3>
            <ul>
                <li><code>/play &lt;song&gt;</code> - Play music</li>
                <li><code>/help</code> - Show help menu</li>
                <li><code>/queue</code> - Show current queue</li>
            </ul>
        </div>
    </div>
</body>
</html>
EOF
    
    # Test Nginx configuration
    nginx -t
    systemctl reload nginx
    
    # Setup SSL if requested
    if [ "$ENABLE_SSL" = "true" ]; then
        print_step "Setting up SSL certificate..."
        certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email admin@"$DOMAIN"
    fi
    
    print_success "Nginx configured for $DOMAIN"
}

# Show deployment information
show_deployment_info() {
    echo
    print_success "=== Deployment Complete! ==="
    echo
    print_step "Service Status:"
    systemctl status nerubot --no-pager -l
    echo
    
    print_step "Useful Commands:"
    echo "  Status:     sudo systemctl status nerubot"
    echo "  Logs:       sudo journalctl -u nerubot -f"
    echo "  Restart:    sudo systemctl restart nerubot"
    echo "  Update:     /home/nerubot/nerubot/deploy/scripts/update.sh"
    echo "  Monitor:    /home/nerubot/monitor.sh"
    echo
    
    if [ -n "$DOMAIN" ]; then
        if [ "$ENABLE_SSL" = "true" ]; then
            print_step "Web Interface: https://$DOMAIN"
        else
            print_step "Web Interface: http://$DOMAIN"
        fi
    fi
    
    print_step "Configuration:"
    echo "  Config:     /home/nerubot/nerubot/.env"
    echo "  Logs:       /home/nerubot/logs/"
    echo "  Backups:    /home/nerubot/backups/"
    echo
    
    if [ -z "$DISCORD_TOKEN" ]; then
        print_warning "Don't forget to configure your Discord token:"
        echo "  sudo -u nerubot nano /home/nerubot/nerubot/.env"
        echo "  sudo systemctl restart nerubot"
    fi
    
    print_success "NeruBot is ready! 🎵"
}

# Main execution
main() {
    parse_args "$@"
    check_root
    deploy_nerubot
}

main "$@"
