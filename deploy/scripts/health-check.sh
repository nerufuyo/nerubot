#!/bin/bash
# NeruBot Simple Health Check
# Basic monitoring and restart if needed

SERVICE_NAME="nerubot"
LOG_FILE="/home/nerubot/logs/health.log"

# Simple logging
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
}

# Check if service is running
if ! systemctl is-active --quiet "$SERVICE_NAME"; then
    log "Service is down, attempting restart"
    systemctl restart "$SERVICE_NAME"
    
    # Wait and check again
    sleep 10
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log "Service restarted successfully"
        echo "✅ Service restarted"
    else
        log "Service restart failed"
        echo "❌ Service restart failed"
        exit 1
    fi
else
    log "Service is running normally"
    echo "✅ Service is healthy"
fi
    if echo "$recent_logs" | grep -qi "connection.*error\|websocket.*error\|login.*failed"; then
        return 1
    fi
    
    # Check for successful activity
    if echo "$recent_logs" | grep -qi "logged in as\|ready\|connected"; then
        return 0
    fi
    
    # If no recent activity, check older logs
    local older_logs=$(journalctl -u "$SERVICE_NAME" --since "1 hour ago" --no-pager -q)
    if echo "$older_logs" | grep -qi "logged in as\|ready"; then
        return 0
    fi
    
    return 1
}

# Restart service with backoff
restart_service() {
    local attempt=$1
    local wait_time=$((attempt * 30))
    
    log_message "⚠️  Attempting to restart $SERVICE_NAME (attempt $attempt)"
    
    if [ $attempt -gt 1 ]; then
        log_message "⏳ Waiting ${wait_time}s before restart..."
        sleep $wait_time
    fi
    
    systemctl restart "$SERVICE_NAME"
    sleep 10
    
    if check_service_status; then
        log_message "✅ Service restarted successfully"
        return 0
    else
        log_message "❌ Service restart failed"
        return 1
    fi
}

# Send notification (placeholder for webhook/email)
send_notification() {
    local message="$1"
    local severity="$2"
    
    # Add your notification logic here
    # Examples: Discord webhook, email, Slack, etc.
    log_message "📢 NOTIFICATION [$severity]: $message"
}

# Main health check function
perform_health_check() {
    local issues=0
    local restart_needed=false
    
    # Check if service is running
    if ! check_service_status; then
        log_message "❌ Service is not running"
        restart_needed=true
        issues=$((issues + 1))
    else
        log_message "✅ Service is running"
    fi
    
    # Check memory usage
    local memory_mb=$(check_memory_usage)
    if [ "$memory_mb" -gt "$MAX_MEMORY_MB" ]; then
        log_message "⚠️  High memory usage: ${memory_mb}MB (limit: ${MAX_MEMORY_MB}MB)"
        restart_needed=true
        issues=$((issues + 1))
    elif [ "$memory_mb" -gt 0 ]; then
        log_message "📊 Memory usage: ${memory_mb}MB"
    fi
    
    # Check Discord connection
    if ! check_discord_connection; then
        log_message "❌ Discord connection issues detected"
        restart_needed=true
        issues=$((issues + 1))
    else
        log_message "✅ Discord connection healthy"
    fi
    
    # Restart if needed
    if [ "$restart_needed" = true ]; then
        send_notification "NeruBot health check failed ($issues issues). Attempting restart." "WARNING"
        
        for attempt in 1 2 3; do
            if restart_service $attempt; then
                send_notification "NeruBot successfully restarted after health check failure." "INFO"
                break
            elif [ $attempt -eq 3 ]; then
                send_notification "NeruBot restart failed after 3 attempts. Manual intervention required." "CRITICAL"
                log_message "🚨 CRITICAL: Failed to restart service after 3 attempts"
            fi
        done
    else
        log_message "✅ All health checks passed"
    fi
}

# Main execution
main() {
    # Create log directory if it doesn't exist
    mkdir -p "$(dirname "$LOG_FILE")"
    
    log_message "🔍 Starting health check for $SERVICE_NAME"
    perform_health_check
    log_message "🏁 Health check completed"
}

# Run main function
main "$@"
