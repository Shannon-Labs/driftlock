#!/bin/bash
# Loom recording state controller
# Tracks Loom recording state and logs anomaly events for recording correlation

set -euo pipefail

# Configuration
SESSION_DIR="${SESSION_DIR:-logs/session_$(date +%s)}"
LOOM_LOG="${LOOM_LOG:-$SESSION_DIR/loom.log}"
LOOM_STATE_FILE="${LOOM_STATE_FILE:-$SESSION_DIR/loom.state}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

mkdir -p "$SESSION_DIR"

# Initialize state
if [ ! -f "$LOOM_STATE_FILE" ]; then
    echo "idle" > "$LOOM_STATE_FILE"
fi

echo "🎬 Loom Controller Started - $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$LOOM_LOG"
echo "   State File: $LOOM_STATE_FILE"
echo "   Current State: $(cat "$LOOM_STATE_FILE")"
echo ""
echo "   📝 Note: Loom has no public API for automation."
echo "   This controller tracks recording state and logs anomaly events."
echo "   Please start/stop Loom manually when prompted."
echo ""

# Function to get current state
get_state() {
    cat "$LOOM_STATE_FILE" 2>/dev/null || echo "idle"
}

# Function to set state
set_state() {
    local new_state="$1"
    echo "$new_state" > "$LOOM_STATE_FILE"
}

# Function to log anomaly event
log_anomaly_event() {
    local anomaly_id="$1"
    local timestamp="$2"
    local state=$(get_state)
    
    local event="
📹 ANOMALY EVENT LOGGED
   Time: $(date '+%Y-%m-%d %H:%M:%S')
   Anomaly ID: ${anomaly_id}
   Trade Timestamp: ${timestamp}
   Loom State: ${state}
"
    
    echo "$event" | tee -a "$LOOM_LOG"
    
    if [ "$state" == "idle" ]; then
        echo -e "${YELLOW}⚠️  LOOM IS IDLE - Please start recording now!${NC}" | tee -a "$LOOM_LOG"
        echo ""
        echo "   To mark as recording, run:"
        echo "   echo 'recording' > $LOOM_STATE_FILE"
        echo ""
    else
        echo -e "${GREEN}✅ LOOM IS RECORDING - Anomaly captured${NC}" | tee -a "$LOOM_LOG"
    fi
}

# Function to start recording (user notification)
start_recording() {
    local anomaly_id="$1"
    local timestamp="$2"
    
    set_state "recording"
    
    echo -e "${GREEN}🎬 RECORDING STARTED${NC}" | tee -a "$LOOM_LOG"
    echo "   Time: $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$LOOM_LOG"
    echo "   Triggered by Anomaly: ${anomaly_id}" | tee -a "$LOOM_LOG"
    echo ""
}

# Function to stop recording
stop_recording() {
    set_state "idle"
    
    echo -e "${CYAN}⏹️  RECORDING STOPPED${NC}" | tee -a "$LOOM_LOG"
    echo "   Time: $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$LOOM_LOG"
    echo ""
}

# Main monitoring loop - watch for anomaly alerts
# This reads from the anomaly alert log
ALERT_LOG="${SESSION_DIR}/../anomaly_alerts.log"

echo "Monitoring for anomaly alerts..."
echo "Press Ctrl+C to exit"
echo ""

# Track last processed line
LAST_LINE=0

while true; do
    sleep 5
    
    if [ ! -f "$ALERT_LOG" ]; then
        continue
    fi
    
    # Get current line count
    CURRENT_LINES=$(wc -l < "$ALERT_LOG" | tr -d ' ')
    
    # Check if there are new alerts
    if [ "$CURRENT_LINES" -gt "$LAST_LINE" ]; then
        # Process new lines
        NEW_ALERTS=$(tail -n +$((LAST_LINE + 1)) "$ALERT_LOG" | grep "Anomaly ID:" || true)
        
        if [ -n "$NEW_ALERTS" ]; then
            # Extract anomaly ID from alert
            ANOMALY_ID=$(echo "$NEW_ALERTS" | head -1 | sed 's/.*Anomaly ID: //' | awk '{print $1}')
            TRADE_TIME=$(echo "$NEW_ALERTS" | head -1 | sed 's/.*Trade Time: //' | awk '{print $1, $2}')
            
            log_anomaly_event "$ANOMALY_ID" "$TRADE_TIME"
            
            # If idle, prompt to start recording
            if [ "$(get_state)" == "idle" ]; then
                start_recording "$ANOMALY_ID" "$TRADE_TIME"
            fi
        fi
        
        LAST_LINE=$CURRENT_LINES
    fi
done
