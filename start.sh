#!/bin/bash
# NeruBot - Discord Music Bot
# Unified startup script for development and production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Project settings
PROJECT_NAME="NeruBot"
PYTHON_CMD="python3"
VENV_NAME="nerubot_env"

# Print colored messages
print_message() {
    echo -e "${2}${1}${NC}"
}

print_success() { print_message "✅ $1" "$GREEN"; }
print_error() { print_message "❌ $1" "$RED"; }
print_warning() { print_message "⚠️  $1" "$YELLOW"; }
print_info() { print_message "ℹ️  $1" "$BLUE"; }
print_step() { print_message "🔄 $1" "$PURPLE"; }

# Help function
show_help() {
    echo
    print_info "$PROJECT_NAME - Discord Music Bot"
    echo
    print_info "Usage: ./start.sh [OPTION]"
    echo
    print_info "Options:"
    print_info "  help, -h, --help    Show this help message"
    print_info "  setup               Setup environment and dependencies only"
    print_info "  run                 Run the bot (assumes setup is complete)"
    print_info "  debug               Run in debug mode with verbose logging"
    print_info "  clean               Clean install (remove virtual environment)"
    echo
    print_info "Default: Complete setup and run"
    echo
}

# Check if Python 3 is available
check_python() {
    print_step "Checking Python installation..."
    
    if ! command -v $PYTHON_CMD &> /dev/null; then
        print_error "Python 3 is not installed"
        print_info "Please install Python 3.8+ from https://python.org"
        exit 1
    fi
    
    # Check Python version
    PYTHON_VERSION=$($PYTHON_CMD --version 2>&1 | cut -d' ' -f2)
    print_success "Python $PYTHON_VERSION found"
    
    # Verify it's 3.8+
    MAJOR=$(echo $PYTHON_VERSION | cut -d'.' -f1)
    MINOR=$(echo $PYTHON_VERSION | cut -d'.' -f2)
    
    if [[ $MAJOR -lt 3 ]] || [[ $MAJOR -eq 3 && $MINOR -lt 8 ]]; then
        print_error "Python 3.8+ required, found $PYTHON_VERSION"
        exit 1
    fi
}

# Setup virtual environment
setup_venv() {
    print_step "Setting up virtual environment..."
    
    if [[ ! -d "$VENV_NAME" ]]; then
        print_info "Creating virtual environment..."
        $PYTHON_CMD -m venv $VENV_NAME
        print_success "Virtual environment created"
    else
        print_info "Virtual environment exists"
    fi
    
    # Activate virtual environment
    source $VENV_NAME/bin/activate
    print_success "Virtual environment activated"
}

# Install dependencies
install_deps() {
    print_step "Installing dependencies..."
    
    if [[ ! -f "requirements.txt" ]]; then
        print_error "requirements.txt not found"
        exit 1
    fi
    
    # Upgrade pip first
    pip install --upgrade pip --quiet
    
    # Install requirements
    pip install -r requirements.txt --quiet
    print_success "Dependencies installed"
}

# Setup environment file
setup_env() {
    print_step "Checking environment configuration..."
    
    if [[ ! -f ".env" ]]; then
        print_info "Creating .env template..."
        cat > .env << 'EOF'
# Discord Bot Token (Required)
DISCORD_TOKEN=your_discord_bot_token_here

# Optional: Spotify API credentials
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# Optional: Bot configuration
LOG_LEVEL=INFO
COMMAND_PREFIX=!
EOF
        print_warning ".env file created with template"
        print_warning "Please edit .env and add your Discord bot token"
        print_info "Get your token from: https://discord.com/developers/applications"
        return 1
    fi
    
    # Check if token is set
    if grep -q "your_discord_bot_token_here" .env; then
        print_error "Discord token not configured in .env"
        print_info "Please edit .env and set DISCORD_TOKEN"
        return 1
    fi
    
    print_success "Environment configured"
    return 0
}

# Run the bot
run_bot() {
    print_step "Starting $PROJECT_NAME..."
    
    # Ensure virtual environment is activated
    if [[ "$VIRTUAL_ENV" == "" ]]; then
        source $VENV_NAME/bin/activate
    fi
    
    # Set Python path
    export PYTHONPATH="$(pwd):$PYTHONPATH"
    
    # Set debug mode if requested
    if [[ "$1" == "debug" ]]; then
        export LOG_LEVEL=DEBUG
        print_info "Debug mode enabled"
    fi
    
    print_success "$PROJECT_NAME is starting..."
    print_info "Use Ctrl+C to stop the bot"
    echo
    
    # Run the bot
    $PYTHON_CMD src/main.py
}

# Clean installation
clean_install() {
    print_step "Cleaning installation..."
    
    if [[ -d "$VENV_NAME" ]]; then
        rm -rf $VENV_NAME
        print_success "Virtual environment removed"
    fi
    
    # Clean cache files
    find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
    find . -name "*.pyc" -delete 2>/dev/null || true
    find . -name "*.log" -delete 2>/dev/null || true
    
    print_success "Cleanup complete"
}

# Main function
main() {
    # Change to script directory
    cd "$(dirname "$0")"
    
    # Handle command line arguments
    case "${1:-}" in
        "help"|"-h"|"--help")
            show_help
            exit 0
            ;;
        "clean")
            clean_install
            exit 0
            ;;
        "setup")
            print_success "=== $PROJECT_NAME Setup ==="
            echo
            check_python
            setup_venv
            install_deps
            if setup_env; then
                print_success "Setup complete! Run './start.sh' to start the bot"
            else
                print_warning "Setup complete, but environment needs configuration"
            fi
            exit 0
            ;;
        "run")
            run_bot
            exit 0
            ;;
        "debug")
            setup_venv
            run_bot debug
            exit 0
            ;;
        "")
            # Default: full setup and run
            print_success "=== $PROJECT_NAME - Discord Music Bot ==="
            echo
            check_python
            setup_venv
            install_deps
            
            if ! setup_env; then
                exit 1
            fi
            
            echo
            run_bot
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Trap signals for clean shutdown
trap 'print_info "Shutting down..."; exit 0' SIGINT SIGTERM

# Run main function
main "$@"
