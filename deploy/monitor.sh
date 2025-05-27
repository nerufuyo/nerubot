#!/bin/bash
# NeruBot Health Monitor - Simple and effective
SERVICE_NAME="nerubot"
LOG_FILE="/home/nerubot/logs/health.log"

# Simple logging
log() { echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"; }

# Create log directory
mkdir -p "$(dirname "$LOG_FILE")"

# Check service status
if ! systemctl is-active --quiet "$SERVICE_NAME"; then
    log "❌ Service down, restarting..."
    systemctl restart "$SERVICE_NAME"
    sleep 10
    
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        log "✅ Service restarted successfully"
    else
        log "❌ Restart failed - manual intervention required"
        exit 1
    fi
else
    log "✅ Service healthy"
fi

# Check memory usage (optional alert if > 500MB)
PID=$(systemctl show "$SERVICE_NAME" --property=MainPID --value 2>/dev/null)
if [[ "$PID" != "0" && -n "$PID" ]]; then
    MEM_KB=$(ps -p "$PID" -o rss= 2>/dev/null | xargs)
    if [[ -n "$MEM_KB" ]]; then
        MEM_MB=$((MEM_KB / 1024))
        if [[ $MEM_MB -gt 500 ]]; then
            log "⚠️ High memory usage: ${MEM_MB}MB"
        fi
    fi
fi
