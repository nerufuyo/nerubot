#!/bin/bash
# NeruBot Deployment Script
# Upload bot files to VPS

set -e

print_success() { echo -e "\033[0;32m✅ $1\033[0m"; }
print_error() { echo -e "\033[0;31m❌ $1\033[0m"; }

# Check arguments
if [[ $# -eq 0 ]]; then
    print_error "Please provide VPS IP address"
    echo "Usage: ./deploy/deploy.sh your_vps_ip"
    exit 1
fi

VPS_IP="$1"
VPS_USER="root"

# Check if we're in the right directory
if [[ ! -f "src/main.py" ]]; then
    print_error "Please run from NeruBot root directory"
    exit 1
fi

print_success "Deploying to VPS: $VPS_IP"

# Create deployment archive
echo "📦 Creating deployment package..."
tar --exclude='nerubot_env' \
    --exclude='.git' \
    --exclude='__pycache__' \
    --exclude='*.pyc' \
    --exclude='.env' \
    --exclude='bot.log' \
    -czf nerubot-deploy.tar.gz .

# Upload to VPS
echo "🚀 Uploading to VPS..."
scp nerubot-deploy.tar.gz $VPS_USER@$VPS_IP:/tmp/

# Extract and setup on VPS
echo "⚙️  Setting up on VPS..."
ssh $VPS_USER@$VPS_IP << 'EOF'
cd /home/nerubot
tar -xzf /tmp/nerubot-deploy.tar.gz
chown -R nerubot:nerubot /home/nerubot/nerubot
rm /tmp/nerubot-deploy.tar.gz
EOF

# Cleanup
rm nerubot-deploy.tar.gz

print_success "Deployment completed!"
echo "Next: SSH to VPS and run setup: sudo su - nerubot && cd nerubot && ./install.sh"

# Step 1: Upload setup script
print_info "1. Uploading setup script..."
scp deploy/simple_vps_setup.sh root@$VPS_IP:/tmp/
print_success "Setup script uploaded"

# Step 2: Create and upload project tarball
print_info "2. Creating project tarball..."
tar -czf nerubot.tar.gz \
    --exclude=nerubot_env \
    --exclude=.git \
    --exclude=__pycache__ \
    --exclude=*.log \
    --exclude="*.tar.gz" \
    .

print_info "3. Uploading project files..."
scp nerubot.tar.gz root@$VPS_IP:/tmp/
print_success "Project files uploaded"

# Cleanup local tarball
rm nerubot.tar.gz

print_success "Upload complete!"
print_info "Now SSH into your VPS and run:"
echo ""
echo "  ssh root@$VPS_IP"
echo "  chmod +x /tmp/simple_vps_setup.sh"
echo "  sudo /tmp/simple_vps_setup.sh"
echo ""
print_info "After the setup script completes, run:"
echo ""
echo "  sudo su - nerubot"
echo "  cd /home/nerubot"
echo "  tar -xzf /tmp/nerubot.tar.gz"
echo "  cd nerubot"
echo "  ./run_nerubot.sh --setup-only"
echo ""
print_warning "Don't forget to create your .env file with your Discord token!"
