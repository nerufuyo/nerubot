#!/bin/bash
# NeruBot Status Dashboard - Compact and informative
SERVICE_NAME="nerubot"

# Colors
G='\033[0;32m'; R='\033[0;31m'; Y='\033[1;33m'; B='\033[0;34m'; C='\033[0;36m'; NC='\033[0m'

# Status check
if systemctl is-active --quiet "$SERVICE_NAME"; then
    STATUS="${G}RUNNING${NC}"
else
    STATUS="${R}STOPPED${NC}"
fi

# Get uptime
UPTIME=$(systemctl show "$SERVICE_NAME" --property=ActiveEnterTimestamp --value 2>/dev/null)
if [[ -n "$UPTIME" && "$UPTIME" != "n/a" ]]; then
    START_TIME=$(date -d "$UPTIME" +%s 2>/dev/null)
    if [[ -n "$START_TIME" ]]; then
        CURRENT_TIME=$(date +%s)
        UPTIME_SEC=$((CURRENT_TIME - START_TIME))
        UPTIME_STR="${UPTIME_SEC}s"
        
        if [[ $UPTIME_SEC -gt 86400 ]]; then
            DAYS=$((UPTIME_SEC / 86400))
            HOURS=$(((UPTIME_SEC % 86400) / 3600))
            UPTIME_STR="${DAYS}d ${HOURS}h"
        elif [[ $UPTIME_SEC -gt 3600 ]]; then
            HOURS=$((UPTIME_SEC / 3600))
            MINS=$(((UPTIME_SEC % 3600) / 60))
            UPTIME_STR="${HOURS}h ${MINS}m"
        elif [[ $UPTIME_SEC -gt 60 ]]; then
            MINS=$((UPTIME_SEC / 60))
            UPTIME_STR="${MINS}m"
        fi
    fi
else
    UPTIME_STR="N/A"
fi

# Get memory usage
PID=$(systemctl show "$SERVICE_NAME" --property=MainPID --value 2>/dev/null)
if [[ "$PID" != "0" && -n "$PID" ]]; then
    MEM_KB=$(ps -p "$PID" -o rss= 2>/dev/null | xargs)
    if [[ -n "$MEM_KB" ]]; then
        MEM_MB=$((MEM_KB / 1024))
        MEMORY="${MEM_MB}MB"
    else
        MEMORY="N/A"
    fi
else
    MEMORY="N/A"
fi

# Display dashboard
echo -e "${C}╔══════════════════════════════════════╗${NC}"
echo -e "${C}║${NC}             ${B}NeruBot Status${NC}             ${C}║${NC}"
echo -e "${C}╠══════════════════════════════════════╣${NC}"
echo -e "${C}║${NC} Status:  $STATUS                    ${C}║${NC}"
echo -e "${C}║${NC} Uptime:  ${Y}$UPTIME_STR${NC}                     ${C}║${NC}"
echo -e "${C}║${NC} Memory:  ${Y}$MEMORY${NC}                      ${C}║${NC}"
echo -e "${C}╚══════════════════════════════════════╝${NC}"
echo
echo -e "${Y}Quick Commands:${NC}"
echo -e "  Status:   ${B}systemctl status $SERVICE_NAME${NC}"
echo -e "  Logs:     ${B}journalctl -u $SERVICE_NAME -f${NC}"
echo -e "  Restart:  ${B}systemctl restart $SERVICE_NAME${NC}"
