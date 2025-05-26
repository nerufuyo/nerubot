#!/bin/bash
# NeruBot Complete Setup and Run Script
# This script handles all setup, dependency installation, and running of the NeruBot Discord bot

set -e  # Exit on any error

# Debug mode - uncomment to enable
# set -x

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Project configuration
PROJECT_NAME="NeruBot"
PYTHON_CMD="python3"
VENV_NAME="nerubot_env"

# Function to print colored messages
print_message() {
    echo -e "${2}[${PROJECT_NAME}] $1${NC}"
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check Python version
check_python() {
    print_step "Checking Python installation..."
    
    if ! command_exists python3; then
        print_error "Python 3 is not installed. Please install Python 3.7+ first."
        exit 1
    fi
    
    PYTHON_VERSION=$($PYTHON_CMD --version 2>&1 | cut -d' ' -f2)
    print_success "Python $PYTHON_VERSION found"
    
    # Check if version is 3.7+
    PYTHON_MAJOR=$(echo $PYTHON_VERSION | cut -d'.' -f1)
    PYTHON_MINOR=$(echo $PYTHON_VERSION | cut -d'.' -f2)
    
    if [[ $PYTHON_MAJOR -lt 3 ]] || [[ $PYTHON_MAJOR -eq 3 && $PYTHON_MINOR -lt 7 ]]; then
        print_error "Python 3.7+ is required. Found: $PYTHON_VERSION"
        exit 1
    fi
}

# Function to create virtual environment
setup_virtual_env() {
    print_step "Setting up virtual environment..."
    
    if [[ ! -d "$VENV_NAME" ]]; then
        print_info "Creating virtual environment: $VENV_NAME"
        $PYTHON_CMD -m venv $VENV_NAME
        print_success "Virtual environment created"
    else
        print_info "Virtual environment already exists"
    fi
    
    # Activate virtual environment
    print_info "Activating virtual environment..."
    source $VENV_NAME/bin/activate
    print_success "Virtual environment activated"
}

# Function to install dependencies
install_dependencies() {
    print_step "Installing dependencies..."
    
    if [[ -f "requirements.txt" ]]; then
        print_info "Installing from requirements.txt..."
        pip install --upgrade pip
        pip install -r requirements.txt
        print_success "Dependencies installed successfully"
    else
        print_error "requirements.txt not found!"
        exit 1
    fi
}

# Function to check and setup .env file
setup_env_file() {
    print_step "Checking environment configuration..."
    
    if [[ ! -f ".env" ]]; then
        print_warning ".env file not found. Creating template..."
        cat > .env << EOF
# Discord Bot Configuration
DISCORD_TOKEN=your_discord_bot_token_here

# Optional: Set log level (DEBUG, INFO, WARNING, ERROR)
LOG_LEVEL=INFO

# Optional: Bot prefix (default is !)
BOT_PREFIX=!
EOF
        print_error "Please edit .env file and add your Discord bot token!"
        print_info "Open .env file and replace 'your_discord_bot_token_here' with your actual Discord bot token"
        print_info "You can get a bot token from: https://discord.com/developers/applications"
        
        # Ask if user wants to continue
        echo
        print_warning "Do you want to open .env file now? (y/n)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            if command_exists code; then
                code .env
            elif command_exists nano; then
                nano .env
            elif command_exists vim; then
                vim .env
            else
                print_info "Please edit .env manually with your preferred editor"
            fi
        fi
        
        print_warning "After setting up your token, run this script again."
        exit 0
    fi
    
    # Check if token is set
    if grep -q "DISCORD_TOKEN=your_discord_bot_token_here" .env || ! grep -q "DISCORD_TOKEN=" .env || grep -q "DISCORD_TOKEN=$" .env; then
        print_error "Discord token not properly set in .env file!"
        print_info "Please edit .env file and set your Discord bot token"
        exit 1
    fi
    
    print_success "Environment file configured"
}

# Function to check project structure
check_project_structure() {
    print_step "Verifying project structure..."
    
    REQUIRED_FILES=("src/main.py" "src/interfaces/discord/bot.py")
    MISSING_FILES=()
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [[ ! -f "$file" ]]; then
            MISSING_FILES+=("$file")
        fi
    done
    
    if [[ ${#MISSING_FILES[@]} -gt 0 ]]; then
        print_error "Missing required files:"
        for file in "${MISSING_FILES[@]}"; do
            print_error "  - $file"
        done
        exit 1
    fi
    
    print_success "Project structure verified"
}

# Function to run the bot
run_bot() {
    print_step "Starting $PROJECT_NAME..."
    print_info "Bot will start with prefix: !"
    print_info "Press Ctrl+C to stop the bot"
    echo
    
    # Ensure we're in the virtual environment
    source $VENV_NAME/bin/activate
    
    # Add current directory to Python path
    export PYTHONPATH="${PYTHONPATH}:$(pwd)"
    
    # Run the bot
    $PYTHON_CMD src/main.py
}

# Function to handle cleanup
cleanup() {
    print_info "Shutting down $PROJECT_NAME..."
    if [[ -n "$BOT_PID" ]]; then
        kill $BOT_PID 2>/dev/null || true
    fi
    print_success "Bot stopped"
    exit 0
}

# Function to show help
show_help() {
    echo
    print_info "NeruBot Setup and Run Script"
    echo
    print_info "Usage: ./run_nerubot.sh [option]"
    echo
    print_info "Options:"
    print_info "  --help, -h     Show this help message"
    print_info "  --setup-only   Only setup environment and dependencies, don't run"
    print_info "  --run-only     Only run the bot (assumes setup is complete)"
    print_info "  --clean        Remove virtual environment and start fresh"
    echo
    print_info "Default behavior: Complete setup and run the bot"
    echo
}

# Function to clean installation
clean_install() {
    print_step "Cleaning previous installation..."
    if [[ -d "$VENV_NAME" ]]; then
        rm -rf $VENV_NAME
        print_success "Virtual environment removed"
    fi
    if [[ -f "bot.log" ]]; then
        rm -f bot.log
        print_success "Log file removed"
    fi
}

# Main script execution
main() {
    clear
    echo
    print_success "=== $PROJECT_NAME Setup and Run Script ==="
    echo
    print_info "This script will:"
    print_info "1. Check Python installation"
    print_info "2. Setup virtual environment"
    print_info "3. Install dependencies"
    print_info "4. Configure environment"
    print_info "5. Run the Discord bot"
    echo
    
    # Parse command line arguments
    case "${1:-}" in
        --help|-h)
            show_help
            exit 0
            ;;
        --clean)
            clean_install
            exit 0
            ;;
        --setup-only)
            SETUP_ONLY=true
            ;;
        --run-only)
            RUN_ONLY=true
            ;;
    esac
    
    # Setup signal handlers
    trap cleanup SIGINT SIGTERM
    
    if [[ "$RUN_ONLY" != true ]]; then
        # Run setup steps
        check_python
        setup_virtual_env
        install_dependencies
        setup_env_file
        check_project_structure
        
        print_success "Setup completed successfully!"
        echo
    fi
    
    if [[ "$SETUP_ONLY" != true ]]; then
        # Run the bot
        run_bot
    fi
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
