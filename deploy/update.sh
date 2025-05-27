#!/bin/bash
# NeruBot Update Script - Simple and safe updates
set -e

# Colors
G='\033[0;32m'; R='\033[0;31m'; Y='\033[1;33m'; B='\033[0;34m'; NC='\033[0m'
log() { echo -e "${B}[Update]${NC} $1"; }
success() { echo -e "${G}✓${NC} $1"; }
error() { echo -e "${R}✗${NC} $1"; exit 1; }

SERVICE_NAME="nerubot"
PROJECT_DIR="/home/nerubot/nerubot"

# Check if running as nerubot user
[[ "$(whoami)" == "nerubot" ]] || error "Run as nerubot user: sudo su - nerubot"

# Navigate to project directory
cd "$PROJECT_DIR" || error "Project directory not found"

log "Starting NeruBot update..."

# Stop service
log "Stopping service..."
sudo systemctl stop "$SERVICE_NAME"

# Update code
log "Pulling latest changes..."
git fetch origin
COMMITS_BEHIND=$(git rev-list --count HEAD..origin/main)

if [[ $COMMITS_BEHIND -eq 0 ]]; then
    log "Already up to date"
    sudo systemctl start "$SERVICE_NAME"
    success "No updates needed"
    exit 0
fi

log "Found $COMMITS_BEHIND new commits"
git pull origin main

# Update dependencies
log "Updating dependencies..."
source nerubot_env/bin/activate
pip install -r requirements.txt -q

# Start service
log "Starting service..."
sudo systemctl start "$SERVICE_NAME"

# Verify
sleep 5
if sudo systemctl is-active --quiet "$SERVICE_NAME"; then
    success "Update completed successfully!"
    log "View logs: sudo journalctl -u $SERVICE_NAME -f"
else
    error "Service failed to start - check logs"
fi
