#!/bin/bash
# NeruBot Update Script
# Handles safe updates with rollback capability

set -e

# Configuration
SERVICE_NAME="nerubot"
PROJECT_DIR="/home/nerubot/nerubot"
BACKUP_DIR="/home/nerubot/backups"
VENV_DIR="$PROJECT_DIR/nerubot_env"
LOG_FILE="/home/nerubot/logs/update.log"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Logging function
log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Print colored messages
print_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
    log_message "SUCCESS: $1"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
    log_message "ERROR: $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
    log_message "WARNING: $1"
}

print_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
    log_message "INFO: $1"
}

# Check if running as nerubot user
check_user() {
    if [ "$(whoami)" != "nerubot" ]; then
        print_error "This script must be run as the nerubot user"
        print_info "Run: sudo su - nerubot"
        exit 1
    fi
}

# Create backup before update
create_backup() {
    local backup_name="pre_update_$(date +%Y%m%d_%H%M%S)"
    local backup_path="$BACKUP_DIR/$backup_name.tar.gz"
    
    print_info "Creating backup: $backup_name"
    
    mkdir -p "$BACKUP_DIR"
    
    cd "$PROJECT_DIR"
    tar -czf "$backup_path" \
        --exclude='nerubot_env' \
        --exclude='__pycache__' \
        --exclude='*.pyc' \
        --exclude='.git' \
        --exclude='logs' \
        .
    
    print_success "Backup created: $backup_path"
    echo "$backup_path" > "$BACKUP_DIR/latest_backup.txt"
}

# Stop the service
stop_service() {
    print_info "Stopping $SERVICE_NAME service..."
    sudo systemctl stop "$SERVICE_NAME"
    
    # Wait for service to stop
    sleep 3
    
    if ! sudo systemctl is-active --quiet "$SERVICE_NAME"; then
        print_success "Service stopped successfully"
    else
        print_error "Failed to stop service"
        exit 1
    fi
}

# Start the service
start_service() {
    print_info "Starting $SERVICE_NAME service..."
    sudo systemctl start "$SERVICE_NAME"
    
    # Wait for service to start
    sleep 5
    
    if sudo systemctl is-active --quiet "$SERVICE_NAME"; then
        print_success "Service started successfully"
    else
        print_error "Failed to start service"
        return 1
    fi
}

# Check service health
check_service_health() {
    print_info "Checking service health..."
    
    # Wait for service to fully initialize
    sleep 10
    
    # Check if service is running
    if ! sudo systemctl is-active --quiet "$SERVICE_NAME"; then
        print_error "Service is not running"
        return 1
    fi
    
    # Check for recent errors in logs
    local recent_errors=$(sudo journalctl -u "$SERVICE_NAME" --since "2 minutes ago" --no-pager -q | grep -i "error\|exception\|failed" | wc -l)
    
    if [ "$recent_errors" -gt 0 ]; then
        print_warning "Found $recent_errors recent errors in logs"
        print_info "Recent logs:"
        sudo journalctl -u "$SERVICE_NAME" --since "2 minutes ago" --no-pager -n 10
        return 1
    fi
    
    print_success "Service health check passed"
    return 0
}

# Git operations
update_code() {
    print_info "Updating code from repository..."
    
    cd "$PROJECT_DIR"
    
    # Stash any local changes
    if git status --porcelain | grep -q .; then
        print_warning "Local changes detected, stashing..."
        git stash
    fi
    
    # Fetch latest changes
    git fetch origin
    
    # Show what will be updated
    local commits_behind=$(git rev-list --count HEAD..origin/main)
    if [ "$commits_behind" -eq 0 ]; then
        print_info "Already up to date"
        return 0
    fi
    
    print_info "Found $commits_behind new commits"
    git log --oneline HEAD..origin/main
    
    # Pull changes
    git pull origin main
    
    print_success "Code updated successfully"
}

# Update dependencies
update_dependencies() {
    print_info "Updating Python dependencies..."
    
    cd "$PROJECT_DIR"
    source "$VENV_DIR/bin/activate"
    
    # Backup current requirements
    pip freeze > "$BACKUP_DIR/requirements_before_update_$(date +%Y%m%d_%H%M%S).txt"
    
    # Update pip
    pip install --upgrade pip
    
    # Install/update requirements
    pip install -r requirements.txt --upgrade
    
    print_success "Dependencies updated successfully"
}

# Rollback to previous version
rollback() {
    print_error "Rolling back to previous version..."
    
    local latest_backup_file="$BACKUP_DIR/latest_backup.txt"
    if [ ! -f "$latest_backup_file" ]; then
        print_error "No backup found for rollback"
        exit 1
    fi
    
    local backup_path=$(cat "$latest_backup_file")
    if [ ! -f "$backup_path" ]; then
        print_error "Backup file not found: $backup_path"
        exit 1
    fi
    
    print_info "Restoring from backup: $backup_path"
    
    # Stop service
    stop_service
    
    # Restore backup
    cd "$PROJECT_DIR"
    tar -xzf "$backup_path"
    
    # Update dependencies from backup
    if [ -f "$BACKUP_DIR"/requirements_before_update_*.txt ]; then
        local backup_reqs=$(ls -t "$BACKUP_DIR"/requirements_before_update_*.txt | head -1)
        source "$VENV_DIR/bin/activate"
        pip install -r "$backup_reqs"
    fi
    
    # Start service
    start_service
    
    if check_service_health; then
        print_success "Rollback completed successfully"
    else
        print_error "Rollback failed"
        exit 1
    fi
}

# Show help
show_help() {
    echo "NeruBot Update Script"
    echo
    echo "Usage: $0 [options]"
    echo
    echo "Options:"
    echo "  --help, -h        Show this help message"
    echo "  --rollback        Rollback to previous version"
    echo "  --check-only      Check for updates without applying"
    echo "  --force           Force update even if no changes"
    echo "  --no-backup       Skip backup creation (not recommended)"
    echo
    echo "Default behavior: Create backup, update code and dependencies, restart service"
}

# Check for updates without applying
check_updates() {
    print_info "Checking for updates..."
    
    cd "$PROJECT_DIR"
    git fetch origin
    
    local commits_behind=$(git rev-list --count HEAD..origin/main)
    
    if [ "$commits_behind" -eq 0 ]; then
        print_info "No updates available"
        return 0
    else
        print_info "Updates available: $commits_behind commits"
        print_info "New commits:"
        git log --oneline HEAD..origin/main
        return 1
    fi
}

# Main update process
perform_update() {
    local skip_backup=${1:-false}
    local force=${2:-false}
    
    print_info "Starting NeruBot update process..."
    
    # Check for updates first (unless forced)
    if [ "$force" != "true" ]; then
        if check_updates; then
            print_info "No updates needed"
            return 0
        fi
    fi
    
    # Create backup
    if [ "$skip_backup" != "true" ]; then
        create_backup
    fi
    
    # Stop service
    stop_service
    
    # Update code
    update_code
    
    # Update dependencies
    update_dependencies
    
    # Start service
    start_service
    
    # Check health
    if check_service_health; then
        print_success "Update completed successfully!"
        print_info "Check status with: sudo systemctl status $SERVICE_NAME"
        print_info "View logs with: sudo journalctl -u $SERVICE_NAME -f"
    else
        print_error "Update failed health check"
        print_warning "Consider rolling back with: $0 --rollback"
        exit 1
    fi
}

# Main function
main() {
    check_user
    
    case "${1:-}" in
        --help|-h)
            show_help
            exit 0
            ;;
        --rollback)
            rollback
            exit 0
            ;;
        --check-only)
            check_updates
            exit $?
            ;;
        --force)
            perform_update false true
            ;;
        --no-backup)
            if [ "$2" = "--force" ]; then
                perform_update true true
            else
                perform_update true false
            fi
            ;;
        "")
            perform_update false false
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
