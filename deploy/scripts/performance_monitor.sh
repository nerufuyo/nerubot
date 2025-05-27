#!/bin/bash
# NeruBot Performance Monitor
# Provides detailed performance metrics

# Configuration
SERVICE_NAME="nerubot"
LOG_FILE="/home/nerubot/logs/performance.log"
METRICS_FILE="/home/nerubot/logs/metrics.json"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Get service PID
get_service_pid() {
    systemctl show "$SERVICE_NAME" --property=MainPID --value 2>/dev/null
}

# Get memory usage
get_memory_usage() {
    local pid=$(get_service_pid)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        ps -p "$pid" -o rss= 2>/dev/null | xargs
    else
        echo "0"
    fi
}

# Get CPU usage
get_cpu_usage() {
    local pid=$(get_service_pid)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        ps -p "$pid" -o %cpu= 2>/dev/null | xargs
    else
        echo "0.0"
    fi
}

# Get service uptime
get_service_uptime() {
    local start_time=$(systemctl show "$SERVICE_NAME" --property=ActiveEnterTimestamp --value)
    if [ -n "$start_time" ] && [ "$start_time" != "n/a" ]; then
        local start_epoch=$(date -d "$start_time" +%s 2>/dev/null)
        local current_epoch=$(date +%s)
        if [ -n "$start_epoch" ]; then
            echo $((current_epoch - start_epoch))
        else
            echo "0"
        fi
    else
        echo "0"
    fi
}

# Format uptime
format_uptime() {
    local total_seconds=$1
    local days=$((total_seconds / 86400))
    local hours=$(((total_seconds % 86400) / 3600))
    local minutes=$(((total_seconds % 3600) / 60))
    local seconds=$((total_seconds % 60))
    
    if [ $days -gt 0 ]; then
        echo "${days}d ${hours}h ${minutes}m ${seconds}s"
    elif [ $hours -gt 0 ]; then
        echo "${hours}h ${minutes}m ${seconds}s"
    elif [ $minutes -gt 0 ]; then
        echo "${minutes}m ${seconds}s"
    else
        echo "${seconds}s"
    fi
}

# Get system load
get_system_load() {
    uptime | awk -F'load average:' '{print $2}' | sed 's/^ *//'
}

# Get disk usage
get_disk_usage() {
    df -h /home/nerubot | tail -1 | awk '{print $5}' | sed 's/%//'
}

# Get network connections
get_network_connections() {
    local pid=$(get_service_pid)
    if [ "$pid" != "0" ] && [ -n "$pid" ]; then
        netstat -tulpn 2>/dev/null | grep "$pid" | wc -l
    else
        echo "0"
    fi
}

# Get recent error count
get_recent_errors() {
    journalctl -u "$SERVICE_NAME" --since "1 hour ago" --no-pager -q | grep -i "error\|exception\|failed" | wc -l
}

# Generate JSON metrics
generate_metrics_json() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local memory_kb=$(get_memory_usage)
    local memory_mb=$((memory_kb / 1024))
    local cpu_usage=$(get_cpu_usage)
    local uptime_seconds=$(get_service_uptime)
    local system_load=$(get_system_load)
    local disk_usage=$(get_disk_usage)
    local network_connections=$(get_network_connections)
    local recent_errors=$(get_recent_errors)
    local service_active=$(systemctl is-active "$SERVICE_NAME")
    
    cat << EOF
{
  "timestamp": "$timestamp",
  "service": {
    "name": "$SERVICE_NAME",
    "status": "$service_active",
    "uptime_seconds": $uptime_seconds,
    "pid": $(get_service_pid)
  },
  "resources": {
    "memory_kb": $memory_kb,
    "memory_mb": $memory_mb,
    "cpu_percent": $cpu_usage,
    "network_connections": $network_connections
  },
  "system": {
    "load_average": "$system_load",
    "disk_usage_percent": $disk_usage
  },
  "health": {
    "recent_errors": $recent_errors
  }
}
EOF
}

# Display performance metrics
display_metrics() {
    local memory_kb=$(get_memory_usage)
    local memory_mb=$((memory_kb / 1024))
    local cpu_usage=$(get_cpu_usage)
    local uptime_seconds=$(get_service_uptime)
    local uptime_formatted=$(format_uptime $uptime_seconds)
    local system_load=$(get_system_load)
    local disk_usage=$(get_disk_usage)
    local network_connections=$(get_network_connections)
    local recent_errors=$(get_recent_errors)
    local service_status=$(systemctl is-active "$SERVICE_NAME")
    
    echo -e "${BLUE}=== NeruBot Performance Metrics ===${NC}"
    echo -e "${GREEN}Generated: $(date)${NC}"
    echo
    
    # Service Status
    echo -e "${YELLOW}📊 Service Status${NC}"
    if [ "$service_status" = "active" ]; then
        echo -e "  Status: ${GREEN}✅ Running${NC}"
    else
        echo -e "  Status: ${RED}❌ Stopped${NC}"
    fi
    echo -e "  Uptime: $uptime_formatted"
    echo -e "  PID: $(get_service_pid)"
    echo
    
    # Resource Usage
    echo -e "${YELLOW}💾 Resource Usage${NC}"
    echo -e "  Memory: ${memory_mb}MB (${memory_kb}KB)"
    echo -e "  CPU: ${cpu_usage}%"
    echo -e "  Network Connections: $network_connections"
    echo
    
    # System Metrics
    echo -e "${YELLOW}🖥️  System Metrics${NC}"
    echo -e "  Load Average: $system_load"
    echo -e "  Disk Usage: ${disk_usage}%"
    echo
    
    # Health Indicators
    echo -e "${YELLOW}🏥 Health Indicators${NC}"
    if [ "$recent_errors" -eq 0 ]; then
        echo -e "  Recent Errors (1h): ${GREEN}$recent_errors${NC}"
    elif [ "$recent_errors" -lt 5 ]; then
        echo -e "  Recent Errors (1h): ${YELLOW}$recent_errors${NC}"
    else
        echo -e "  Recent Errors (1h): ${RED}$recent_errors${NC}"
    fi
    echo
    
    # Performance Tips
    if [ "$memory_mb" -gt 300 ]; then
        echo -e "${RED}⚠️  Warning: High memory usage detected${NC}"
    fi
    
    if [ "$disk_usage" -gt 80 ]; then
        echo -e "${RED}⚠️  Warning: High disk usage detected${NC}"
    fi
    
    if [ "$recent_errors" -gt 10 ]; then
        echo -e "${RED}⚠️  Warning: High error rate detected${NC}"
    fi
}

# Save metrics to file
save_metrics() {
    local json_metrics=$(generate_metrics_json)
    echo "$json_metrics" > "$METRICS_FILE"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Metrics saved" >> "$LOG_FILE"
}

# Main function
main() {
    case "${1:-display}" in
        "json")
            generate_metrics_json
            ;;
        "save")
            save_metrics
            echo "Metrics saved to $METRICS_FILE"
            ;;
        "display"|"")
            display_metrics
            ;;
        "log")
            display_metrics | tee -a "$LOG_FILE"
            ;;
        *)
            echo "Usage: $0 [display|json|save|log]"
            echo "  display - Show metrics in terminal (default)"
            echo "  json    - Output metrics as JSON"
            echo "  save    - Save metrics to file"
            echo "  log     - Display and log metrics"
            exit 1
            ;;
    esac
}

main "$@"
