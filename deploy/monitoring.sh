#!/bin/bash

# NeruBot Health Monitoring Script
# Checks bot status and reports health metrics

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="nerubot"
LOG_FILE="/var/log/nerubot/bot.log"
MAX_MEMORY_MB=512
MAX_CPU_PERCENT=80

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for some operations
check_permissions() {
    if [[ $EUID -eq 0 ]]; then
        log_warning "Running as root - consider using sudo for specific commands only"
    fi
}

# Check service status
check_service_status() {
    log_info "Checking service status..."
    
    if systemctl is-active --quiet $SERVICE_NAME; then
        log_success "Service is running"
        
        # Get service uptime
        uptime=$(systemctl show $SERVICE_NAME --property=ActiveEnterTimestamp --value)
        log_info "Service uptime: $uptime"
        
        return 0
    else
        log_error "Service is not running"
        
        # Show last few log entries
        log_info "Last service logs:"
        journalctl -u $SERVICE_NAME -n 5 --no-pager
        
        return 1
    fi
}

# Check system resources
check_resources() {
    log_info "Checking system resources..."
    
    # Memory usage
    if command -v systemctl &> /dev/null; then
        memory_bytes=$(systemctl show $SERVICE_NAME --property=MemoryCurrent --value 2>/dev/null)
        if [[ "$memory_bytes" != "[not set]" ]] && [[ -n "$memory_bytes" ]] && [[ "$memory_bytes" -gt 0 ]]; then
            memory_mb=$((memory_bytes / 1024 / 1024))
            log_info "Memory usage: ${memory_mb}MB"
            
            if [[ $memory_mb -gt $MAX_MEMORY_MB ]]; then
                log_warning "Memory usage is high: ${memory_mb}MB (limit: ${MAX_MEMORY_MB}MB)"
            fi
        else
            log_warning "Could not determine memory usage"
        fi
    fi
    
    # CPU usage (system-wide)
    cpu_usage=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1}')
    log_info "System CPU usage: ${cpu_usage}%"
    
    # Disk usage
    disk_usage=$(df -h / | awk 'NR==2{print $5}' | sed 's/%//')
    log_info "Disk usage: ${disk_usage}%"
    
    if [[ $disk_usage -gt 80 ]]; then
        log_warning "Disk usage is high: ${disk_usage}%"
    fi
}

# Check log file
check_logs() {
    log_info "Checking log file..."
    
    if [[ -f "$LOG_FILE" ]]; then
        log_size=$(du -h "$LOG_FILE" | cut -f1)
        log_info "Log file size: $log_size"
        
        # Check for recent errors
        recent_errors=$(grep -c "ERROR\|CRITICAL" "$LOG_FILE" 2>/dev/null || echo "0")
        if [[ $recent_errors -gt 0 ]]; then
            log_warning "Found $recent_errors recent errors in log file"
        else
            log_success "No recent errors found in log file"
        fi
    else
        log_warning "Log file not found: $LOG_FILE"
    fi
}

# Check Discord connectivity
check_discord_connection() {
    log_info "Checking Discord connectivity..."
    
    # Check if we can reach Discord API
    if curl -s --max-time 10 https://discord.com/api/v10/gateway > /dev/null; then
        log_success "Discord API is reachable"
    else
        log_error "Cannot reach Discord API"
    fi
}

# Check FFmpeg
check_ffmpeg() {
    log_info "Checking FFmpeg installation..."
    
    if command -v ffmpeg &> /dev/null; then
        ffmpeg_version=$(ffmpeg -version 2>&1 | head -n1 | cut -d' ' -f3)
        log_success "FFmpeg is installed: $ffmpeg_version"
    else
        log_error "FFmpeg is not installed"
    fi
}

# Main health check
main() {
    echo "=== NeruBot Health Check ==="
    echo "Timestamp: $(date)"
    echo ""
    
    check_permissions
    check_service_status
    service_running=$?
    
    check_resources
    check_logs
    check_discord_connection
    check_ffmpeg
    
    echo ""
    if [[ $service_running -eq 0 ]]; then
        log_success "Overall status: HEALTHY"
        exit 0
    else
        log_error "Overall status: UNHEALTHY"
        exit 1
    fi
}

# Parse command line arguments
case "${1:-}" in
    "status")
        check_service_status
        ;;
    "resources")
        check_resources
        ;;
    "logs")
        check_logs
        ;;
    "discord")
        check_discord_connection
        ;;
    "ffmpeg")
        check_ffmpeg
        ;;
    "help"|"-h"|"--help")
        echo "NeruBot Health Monitoring Script"
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  status    - Check service status only"
        echo "  resources - Check system resources only"
        echo "  logs      - Check log file only"
        echo "  discord   - Check Discord connectivity only"
        echo "  ffmpeg    - Check FFmpeg installation only"
        echo "  help      - Show this help message"
        echo ""
        echo "No command = Run full health check"
        ;;
    *)
        main
        ;;
esac
