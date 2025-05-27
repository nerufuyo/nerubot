#!/bin/bash
# Upload NeruBot to VPS Script
# This script helps you upload your NeruBot files to a VPS when the repository is not public

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_message() {
    echo -e "${2}[NeruBot Upload] $1${NC}"
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

# Check if we're in the right directory
if [[ ! -f "deploy/simple_vps_setup.sh" ]]; then
    print_error "Please run this script from the NeruBot root directory"
    print_info "Usage: ./deploy/upload_to_vps.sh your_vps_ip"
    exit 1
fi

# Check for VPS IP argument
if [[ $# -eq 0 ]]; then
    print_error "Please provide your VPS IP address"
    print_info "Usage: ./deploy/upload_to_vps.sh your_vps_ip"
    exit 1
fi

VPS_IP="$1"

print_info "Starting upload to VPS: $VPS_IP"

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
