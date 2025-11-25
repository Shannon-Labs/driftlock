#!/bin/bash
# Real-time anomaly alerting script
# Monitors Driftlock NDJSON output stream and fires alerts on anomaly detection

set -euo pipefail

# Configuration
ANOMALY_LOG="${ANOMALY_LOG:-logs/anomaly_alerts.log}"
ALERT_SOUND="${ALERT_SOUND:-/System/Library/Sounds/Glass.aiff}"

# Colors for terminal output
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Create logs directory if needed
mkdir -p "$(dirname "$ANOMALY_LOG")"

# Initialize alert counter
ALERT_COUNT=0

echo "🔍 Anomaly Alert Monitor Started - $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$ANOMALY_LOG"
echo "   Reading from: stdin"
echo "   Logging to: $ANOMALY_LOG"
echo ""

# Read NDJSON from stdin
while IFS= read -r line; do
    # Check if this is a valid JSON line
    if ! echo "$line" | jq -e . >/dev/null 2>&1; then
        continue
    fi
    
    # Check for anomaly flag
    # Driftlock can mark anomalies with different fields, check common ones
    IS_ANOMALY=$(echo "$line" | jq -r '
        if .anomaly == true or .is_anomaly == true or .detected == true then
            "true"
        else
            "false"
        end
    ' 2>/dev/null || echo "false")
    
    if [ "$IS_ANOMALY" == "true" ]; then
        ALERT_COUNT=$((ALERT_COUNT + 1))
        
        # Extract key fields
        TIMESTAMP=$(echo "$line" | jq -r '.ts // .timestamp // (((.line|fromjson?) // {}) | .ts // .timestamp) // "N/A"' 2>/dev/null || echo "N/A")
        ANOMALY_ID=$(echo "$line" | jq -r '.id // (((.line|fromjson?) // {}) | .id) // "N/A"' 2>/dev/null || echo "N/A")
        PAIR=$(echo "$line" | jq -r '.pair // (((.line|fromjson?) // {}) | .pair) // "N/A"' 2>/dev/null || echo "N/A")
        PRICE=$(echo "$line" | jq -r '.price // (((.line|fromjson?) // {}) | .price) // "N/A"' 2>/dev/null || echo "N/A")
        SCORE=$(echo "$line" | jq -r '.ncd_score // .score // .entropy // (((.line|fromjson?) // {}) | .score // .entropy) // "N/A"' 2>/dev/null || echo "N/A")
        
        # Convert Unix timestamp to readable date
        if [[ "$TIMESTAMP" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
            READABLE_TIME=$(date -r "${TIMESTAMP%.*}" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "$TIMESTAMP")
        else
            READABLE_TIME="$TIMESTAMP"
        fi
        
        # Fire the alert!
        ALERT_MESSAGE="
═══════════════════════════════════════════════════════════
🚨 START LOOM NOW! 🚨
═══════════════════════════════════════════════════════════
Alert #${ALERT_COUNT}
Time: $(date '+%Y-%m-%d %H:%M:%S')
Trade Time: ${READABLE_TIME}
Anomaly ID: ${ANOMALY_ID}
Pair: ${PAIR}
Price: ${PRICE}
Score: ${SCORE}

Full Payload:
${line}
═══════════════════════════════════════════════════════════
"
        
        # Print to terminal with color
        echo -e "${RED}${ALERT_MESSAGE}${NC}"
        
        # Log to file
        echo "$ALERT_MESSAGE" >> "$ANOMALY_LOG"
        
        # Visual + audio alerts
        # macOS terminal bell
        printf '\a'
        
        # macOS say command (non-blocking)
        if command -v say >/dev/null 2>&1; then
            say "Anomaly detected! Start Loom now! Alert number ${ALERT_COUNT}" &
        fi
        
        # Sound effect
        if [ -f "$ALERT_SOUND" ] && command -v afplay >/dev/null 2>&1; then
            afplay "$ALERT_SOUND" &
        fi
        
        # Desktop notification (if osascript available)
        if command -v osascript >/dev/null 2>&1; then
            osascript -e "display notification \"Anomaly #${ALERT_COUNT} detected! ID: ${ANOMALY_ID}\" with title \"🚨 START LOOM NOW!\" sound name \"Glass\"" 2>/dev/null &
        fi
    fi
done

echo ""
echo "🔍 Anomaly Alert Monitor Stopped - $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$ANOMALY_LOG"
echo "   Total Alerts: $ALERT_COUNT"
