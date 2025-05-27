#!/bin/bash
# NeruBot - Simple startup script
set -e

# Colors
G='\033[0;32m'; R='\033[0;31m'; Y='\033[1;33m'; B='\033[0;34m'; NC='\033[0m'

# Functions
log() { echo -e "${B}[NeruBot]${NC} $1"; }
success() { echo -e "${G}✓${NC} $1"; }
error() { echo -e "${R}✗${NC} $1"; exit 1; }
warn() { echo -e "${Y}!${NC} $1"; }

# Check Python
command -v python3 >/dev/null || error "Python 3 not found"

# Setup virtual environment
if [[ ! -d "nerubot_env" ]]; then
    log "Setting up virtual environment..."
    python3 -m venv nerubot_env
    success "Virtual environment created"
fi

# Activate virtual environment
source nerubot_env/bin/activate

# Install dependencies
if [[ ! -f "nerubot_env/installed" ]]; then
    log "Installing dependencies..."
    pip install --upgrade pip -q
    pip install -r requirements.txt -q
    touch nerubot_env/installed
    success "Dependencies installed"
fi

# Check environment file
if [[ ! -f ".env" ]]; then
    log "Creating .env template..."
    cat > .env << 'EOF'
DISCORD_TOKEN=your_discord_bot_token_here
LOG_LEVEL=INFO
COMMAND_PREFIX=!
EOF
    warn ".env created - Please add your Discord token"
    exit 1
fi

# Check token
if grep -q "your_discord_bot_token_here" .env; then
    error "Please configure your Discord token in .env"
fi

# Set Python path and run
export PYTHONPATH="$(pwd):$PYTHONPATH"

case "${1:-}" in
    "debug") export LOG_LEVEL=DEBUG; log "Starting in debug mode..." ;;
    "help"|"-h") 
        echo "Usage: $0 [debug|help]"
        echo "  debug  - Run with debug logging"
        echo "  help   - Show this help"
        exit 0 ;;
    *) log "Starting NeruBot..." ;;
esac

success "NeruBot is starting! Use Ctrl+C to stop."
python3 src/main.py
