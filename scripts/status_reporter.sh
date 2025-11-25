#!/bin/bash
# 15-minute status reporter
# Monitors session statistics and stream health

set -euo pipefail

# Configuration
SESSION_DIR="${SESSION_DIR:-logs/session_$(date +%s)}"
RAW_LOG="${RAW_LOG:-$SESSION_DIR/kraken_raw.ndjson}"
ANOMALY_LOG="${ANOMALY_LOG:-$SESSION_DIR/kraken_anomalies.ndjson}"
STATUS_LOG="${STATUS_LOG:-$SESSION_DIR/status.log}"
REPORT_INTERVAL="${REPORT_INTERVAL:-900}" # 15 minutes in seconds
TEST_INTERVAL="${TEST_INTERVAL:-0}"        # For testing: override with --test-interval

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --test-interval)
            TEST_INTERVAL="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

# Use test interval if set
if [ "$TEST_INTERVAL" -gt 0 ]; then
    REPORT_INTERVAL="$TEST_INTERVAL"
fi

# Create logs directory
mkdir -p "$SESSION_DIR"

# Track session start
SESSION_START=$(date +%s)
echo "📊 Status Reporter Started - $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$STATUS_LOG"
echo "   Interval: ${REPORT_INTERVAL}s"
echo "   Session Dir: $SESSION_DIR"
echo ""

# Function to format large numbers with commas
format_number() {
    printf "%'d" "$1" 2>/dev/null || echo "$1"
}

# Function to format file size
format_size() {
    local bytes=$1
    if [ "$bytes" -lt 1024 ]; then
        echo "${bytes}B"
    elif [ "$bytes" -lt 1048576 ]; then
        echo "$(( bytes / 1024 ))KB"
    else
        echo "$(( bytes / 1048576 ))MB"
    fi
}

# Function to get last anomaly info
get_last_anomaly() {
    if [ -f "$ANOMALY_LOG" ]; then
        # Get last line with anomaly flag
        local last_anomaly=$(tac "$ANOMALY_LOG" | grep -m 1 '"anomaly".*true\|"is_anomaly".*true\|"detected".*true' || echo "")
        if [ -n "$last_anomaly" ]; then
            local ts=$(echo "$last_anomaly" | jq -r '.ts // .timestamp // "N/A"' 2>/dev/null || echo "N/A")
            local id=$(echo "$last_anomaly" | jq -r '.id // "N/A"' 2>/dev/null || echo "N/A")
            
            if [[ "$ts" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
                local readable=$(date -r "${ts%.*}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "$ts")
                echo "$readable|$id"
            else
                echo "$ts|$id"
            fi
        else
            echo "None|N/A"
        fi
    else
        echo "None|N/A"
    fi
}

# Main reporting loop
while true; do
    sleep "$REPORT_INTERVAL"
    
    # Calculate runtime
    NOW=$(date +%s)
    RUNTIME_SEC=$((NOW - SESSION_START))
    RUNTIME_MIN=$((RUNTIME_SEC / 60))
    RUNTIME_HOUR=$((RUNTIME_MIN / 60))
    RUNTIME_MIN_REMAINDER=$((RUNTIME_MIN % 60))
    
    # Count total trades
    TOTAL_TRADES=0
    if [ -f "$RAW_LOG" ]; then
        TOTAL_TRADES=$(wc -l < "$RAW_LOG" | tr -d ' ')
    fi
    
    # Count total anomalies
    TOTAL_ANOMALIES=0
    if [ -f "$ANOMALY_LOG" ]; then
        TOTAL_ANOMALIES=$(grep -c '"anomaly".*true\|"is_anomaly".*true\|"detected".*true' "$ANOMALY_LOG" 2>/dev/null || echo "0")
    fi
    
    # Get last anomaly
    LAST_ANOMALY_INFO=$(get_last_anomaly)
    LAST_ANOMALY_TIME=$(echo "$LAST_ANOMALY_INFO" | cut -d'|' -f1)
    LAST_ANOMALY_ID=$(echo "$LAST_ANOMALY_INFO" | cut -d'|' -f2)
    
    # Get file sizes
    RAW_SIZE=0
    ANOMALY_SIZE=0
    if [ -f "$RAW_LOG" ]; then
        RAW_SIZE=$(stat -f%z "$RAW_LOG" 2>/dev/null || stat -c%s "$RAW_LOG" 2>/dev/null || echo "0")
    fi
    if [ -f "$ANOMALY_LOG" ]; then
        ANOMALY_SIZE=$(stat -f%z "$ANOMALY_LOG" 2>/dev/null || stat -c%s "$ANOMALY_LOG" 2>/dev/null || echo "0")
    fi
    
    # Check stream health
    STREAM_STATUS="${RED}⚠️  UNHEALTHY${NC}"
    STREAM_DETAIL="no data"
    if [ -f "$RAW_LOG" ]; then
        # Check last update time
        LAST_UPDATE=$(stat -f%m "$RAW_LOG" 2>/dev/null || stat -c%Y "$RAW_LOG" 2>/dev/null || echo "0")
        SECONDS_SINCE=$((NOW - LAST_UPDATE))
        
        if [ "$SECONDS_SINCE" -lt 30 ]; then
            STREAM_STATUS="${GREEN}✅ HEALTHY${NC}"
            STREAM_DETAIL="last update ${SECONDS_SINCE}s ago"
        elif [ "$SECONDS_SINCE" -lt 120 ]; then
            STREAM_STATUS="${YELLOW}⚠️  SLOW${NC}"
            STREAM_DETAIL="last update ${SECONDS_SINCE}s ago"
        else
            STREAM_DETAIL="last update $((SECONDS_SINCE / 60))m ago"
        fi
    fi
    
    # Generate status report
    STATUS_REPORT="
⏰ STATUS UPDATE - $(date '+%H:%M:%S')
   Runtime: ${RUNTIME_HOUR}h ${RUNTIME_MIN_REMAINDER}m
   Total Trades: $(format_number $TOTAL_TRADES)
   Total Anomalies: $(format_number $TOTAL_ANOMALIES)
   Last Anomaly: ${LAST_ANOMALY_TIME} (ID: ${LAST_ANOMALY_ID})
   Raw Log: $(format_size $RAW_SIZE) ($RAW_LOG)
   Anomaly Log: $(format_size $ANOMALY_SIZE) ($ANOMALY_LOG)
   Stream Status: ${STREAM_STATUS} (${STREAM_DETAIL})
"
    
    # Print to terminal with colors
    echo -e "${CYAN}${STATUS_REPORT}${NC}"
    
    # Log to file (without colors)
    echo "$STATUS_REPORT" | sed 's/\x1b\[[0-9;]*m//g' >> "$STATUS_LOG"
done
