#!/bin/bash
# NeruBot Deployment Dashboard
# Provides a comprehensive overview of the bot's status

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

SERVICE_NAME="nerubot"
PROJECT_DIR="/home/nerubot/nerubot"

# Function to get colored status
get_status_color() {
    if [ "$1" = "active" ] || [ "$1" = "running" ] || [ "$1" = "healthy" ]; then
        echo -e "${GREEN}$1${NC}"
    elif [ "$1" = "inactive" ] || [ "$1" = "stopped" ] || [ "$1" = "failed" ]; then
        echo -e "${RED}$1${NC}"
    else
        echo -e "${YELLOW}$1${NC}"
    fi
}

# Get service status
get_service_status() {
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        echo "running"
    else
        echo "stopped"
    fi
}

# Get uptime
get_uptime() {
    local start_time=$(systemctl show "$SERVICE_NAME" --property=ActiveEnterTimestamp --value)
    if [ -n "$start_time" ] && [ "$start_time" != "n/a" ]; then
        local start_epoch=$(date -d "$start_time" +%s 2>/dev/null)
        local current_epoch=$(date +%s)
        if [ -n "$start_epoch" ]; then
            local uptime_seconds=$((current_epoch - start_epoch))
            local days=$((uptime_seconds / 86400))
            local hours=$(((uptime_seconds % 86400) / 3600))
            local minutes=$(((uptime_seconds % 3600) / 60))
            
            if [ $days -gt 0 ]; then
                echo "${days}d ${hours}h ${minutes}m"
            elif [ $hours -gt 0 ]; then
                echo "${hours}h ${minutes}m"
            else
                echo "${minutes}m"
            fi
        else
            echo "unknown"
        fi
    else
        echo "not running"
    fi
}

# Get memory usage
get_memory_usage() {
    local pid=$(systemctl show "$SERVICE_NAME" --property=MainPID --value)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        local mem_kb=$(ps -p "$pid" -o rss= 2>/dev/null | xargs)
        if [ -n "$mem_kb" ]; then
            echo "$((mem_kb / 1024))MB"
        else
            echo "unknown"
        fi
    else
        echo "not running"
    fi
}

# Get CPU usage
get_cpu_usage() {
    local pid=$(systemctl show "$SERVICE_NAME" --property=MainPID --value)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        local cpu=$(ps -p "$pid" -o %cpu= 2>/dev/null | xargs)
        if [ -n "$cpu" ]; then
            echo "${cpu}%"
        else
            echo "unknown"
        fi
    else
        echo "not running"
    fi
}

# Get recent errors
get_recent_errors() {
    journalctl -u "$SERVICE_NAME" --since "1 hour ago" --no-pager -q | grep -i "error\|exception\|failed" | wc -l
}

# Get version info
get_version_info() {
    if [ -d "$PROJECT_DIR/.git" ]; then
        cd "$PROJECT_DIR"
        local commit=$(git rev-parse --short HEAD 2>/dev/null)
        local branch=$(git branch --show-current 2>/dev/null)
        if [ -n "$commit" ] && [ -n "$branch" ]; then
            echo "$branch@$commit"
        else
            echo "unknown"
        fi
    else
        echo "not a git repository"
    fi
}

# Check for available updates
check_updates() {
    if [ -d "$PROJECT_DIR/.git" ]; then
        cd "$PROJECT_DIR"
        git fetch origin 2>/dev/null
        local commits_behind=$(git rev-list --count HEAD..origin/main 2>/dev/null)
        if [ "$commits_behind" -gt 0 ]; then
            echo -e "${YELLOW}$commits_behind updates available${NC}"
        else
            echo -e "${GREEN}up to date${NC}"
        fi
    else
        echo "unknown"
    fi
}

# Main dashboard
show_dashboard() {
    clear
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║                    ${PURPLE}NeruBot Dashboard${CYAN}                        ║${NC}"
    echo -e "${CYAN}╠══════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${CYAN}║ Generated: ${BLUE}$(date '+%Y-%m-%d %H:%M:%S')${CYAN}                              ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo
    
    # Service Status
    local service_status=$(get_service_status)
    echo -e "${YELLOW}🤖 Service Status${NC}"
    echo -e "   Status: $(get_status_color "$service_status")"
    echo -e "   Uptime: ${BLUE}$(get_uptime)${NC}"
    echo
    
    # Resource Usage
    echo -e "${YELLOW}📊 Resource Usage${NC}"
    echo -e "   Memory: ${BLUE}$(get_memory_usage)${NC}"
    echo -e "   CPU: ${BLUE}$(get_cpu_usage)${NC}"
    echo
    
    # Health Status
    local recent_errors=$(get_recent_errors)
    echo -e "${YELLOW}🏥 Health Status${NC}"
    if [ "$recent_errors" -eq 0 ]; then
        echo -e "   Recent Errors (1h): $(get_status_color "0 - healthy")"
    elif [ "$recent_errors" -lt 5 ]; then
        echo -e "   Recent Errors (1h): ${YELLOW}$recent_errors - warning${NC}"
    else
        echo -e "   Recent Errors (1h): ${RED}$recent_errors - critical${NC}"
    fi
    echo
    
    # Version Information
    echo -e "${YELLOW}📦 Version Information${NC}"
    echo -e "   Current: ${BLUE}$(get_version_info)${NC}"
    echo -e "   Updates: $(check_updates)"
    echo
    
    # System Information
    echo -e "${YELLOW}🖥️  System Information${NC}"
    echo -e "   Load: ${BLUE}$(uptime | awk -F'load average:' '{print $2}' | sed 's/^ *//')${NC}"
    echo -e "   Disk: ${BLUE}$(df -h /home/nerubot | tail -1 | awk '{print $5}' | sed 's/%//')% used${NC}"
    echo
    
    # Quick Actions
    echo -e "${YELLOW}⚡ Quick Actions${NC}"
    echo -e "   View Logs:     ${CYAN}sudo journalctl -u $SERVICE_NAME -f${NC}"
    echo -e "   Restart Bot:   ${CYAN}sudo systemctl restart $SERVICE_NAME${NC}"
    echo -e "   Health Check:  ${CYAN}/home/nerubot/nerubot/deploy/scripts/health_check.sh${NC}"
    echo -e "   Update Bot:    ${CYAN}/home/nerubot/nerubot/deploy/scripts/update.sh${NC}"
    echo
    
    # Recent Log Entries
    echo -e "${YELLOW}📋 Recent Log Entries (Last 5)${NC}"
    echo -e "${BLUE}────────────────────────────────────────────────────────────${NC}"
    journalctl -u "$SERVICE_NAME" --no-pager -n 5 --output=short-iso | sed 's/^/   /'
    echo
}

# Watch mode
watch_dashboard() {
    while true; do
        show_dashboard
        echo -e "${CYAN}Press Ctrl+C to exit watch mode...${NC}"
        sleep 5
    done
}

# Main function
main() {
    case "${1:-}" in
        "watch"|"-w")
            watch_dashboard
            ;;
        "help"|"-h"|"--help")
            echo "NeruBot Dashboard"
            echo
            echo "Usage: $0 [option]"
            echo
            echo "Options:"
            echo "  (none)    Show dashboard once"
            echo "  watch     Continuous update mode"
            echo "  help      Show this help"
            ;;
        "")
            show_dashboard
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
}

main "$@"
