#!/bin/bash
# Synthetic anomaly injection
# Monitors for 30-minute no-anomaly windows and injects test anomalies

set -euo pipefail

# Configuration
SESSION_DIR="${SESSION_DIR:-logs/session_$(date +%s)}"
ANOMALY_LOG="${ANOMALY_LOG:-$SESSION_DIR/kraken_anomalies.ndjson}"
SYNTHETIC_LOG="${SYNTHETIC_LOG:-$SESSION_DIR/synthetic.log}"
STREAMER_PID_FILE="${STREAMER_PID_FILE:-$SESSION_DIR/streamer.pid}"
NO_ANOMALY_THRESHOLD="${NO_ANOMALY_THRESHOLD:-1800}" # 30 minutes in seconds
FORCE_INJECT="${FORCE_INJECT:-false}"

# Colors
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --force-inject)
            FORCE_INJECT=true
            shift
            ;;
        *)
            shift
            ;;
    esac
done

mkdir -p "$SESSION_DIR"

echo "🧪 Synthetic Anomaly Injector Started - $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$SYNTHETIC_LOG"
echo "   Threshold: ${NO_ANOMALY_THRESHOLD}s ($(($NO_ANOMALY_THRESHOLD / 60)) minutes)"
echo "   Anomaly Log: $ANOMALY_LOG"
echo ""

# Function to get last anomaly timestamp
get_last_anomaly_time() {
    if [ ! -f "$ANOMALY_LOG" ]; then
        echo "0"
        return
    fi
    
    local last_anomaly=$(tac "$ANOMALY_LOG" | grep -m 1 '"anomaly".*true\|"is_anomaly".*true\|"detected".*true' 2>/dev/null || echo "")
    if [ -n "$last_anomaly" ]; then
        echo "$last_anomaly" | jq -r '.ts // .timestamp // "0"' 2>/dev/null || echo "0"
    else
        echo "0"
    fi
}

# Function to inject synthetic anomaly via pair switch
inject_pair_switch() {
    local reason="$1"
    
    echo -e "${YELLOW}🧪 INJECTING SYNTHETIC ANOMALY${NC}" | tee -a "$SYNTHETIC_LOG"
    echo "   Reason: $reason" | tee -a "$SYNTHETIC_LOG"
    echo "   Method: Switch BTC/USD → ETH/USD → BTC/USD" | tee -a "$SYNTHETIC_LOG"
    echo "   Time: $(date '+%Y-%m-%d %H:%M:%S')" | tee -a "$SYNTHETIC_LOG"
    
    # Get streamer PID
    if [ -f "$STREAMER_PID_FILE" ]; then
        local streamer_pid=$(cat "$STREAMER_PID_FILE")
        
        # Kill existing BTC/USD streamer
        if kill -0 "$streamer_pid" 2>/dev/null; then
            echo "   Stopping BTC/USD streamer (PID: $streamer_pid)..." | tee -a "$SYNTHETIC_LOG"
            kill "$streamer_pid" 2>/dev/null || true
            sleep 2
        fi
    fi
    
    # Start ETH/USD streamer
    echo "   Starting ETH/USD streamer for 120 seconds..." | tee -a "$SYNTHETIC_LOG"
    export KRAKEN_PAIR="ETH/USD"
    
    # Start streamer in background, capture PID
    python3 -u scripts/stream_kraken_ws.py 2>>"$SESSION_DIR/streamer.stderr.log" | \
        tee -a "$SESSION_DIR/kraken_raw.ndjson" > /tmp/kraken_synthetic.pipe &
    
    local eth_pid=$!
    echo "$eth_pid" > "$STREAMER_PID_FILE"
    
    # Wait 120 seconds
    sleep 120
    
    # Kill ETH/USD streamer
    if kill -0 "$eth_pid" 2>/dev/null; then
        echo "   Stopping ETH/USD streamer (PID: $eth_pid)..." | tee -a "$SYNTHETIC_LOG"
        kill "$eth_pid" 2>/dev/null || true
        sleep 2
    fi
    
    # Restart BTC/USD streamer
    echo "   Restarting BTC/USD streamer..." | tee -a "$SYNTHETIC_LOG"
    export KRAKEN_PAIR="BTC/USD"
    
    python3 -u scripts/stream_kraken_ws.py 2>>"$SESSION_DIR/streamer.stderr.log" | \
        tee -a "$SESSION_DIR/kraken_raw.ndjson" > /tmp/kraken.pipe &
    
    local btc_pid=$!
    echo "$btc_pid" > "$STREAMER_PID_FILE"
    
    echo -e "${GREEN}✅ Synthetic injection complete - reverted to BTC/USD (PID: $btc_pid)${NC}" | tee -a "$SYNTHETIC_LOG"
    echo ""
}

# Force inject if requested
if [ "$FORCE_INJECT" = true ]; then
    inject_pair_switch "Force inject requested via --force-inject"
    exit 0
fi

# Main monitoring loop
CHECK_INTERVAL=300 # Check every 5 minutes
LAST_INJECTION_TIME=0

while true; do
    sleep "$CHECK_INTERVAL"
    
    NOW=$(date +%s)
    LAST_ANOMALY_TS=$(get_last_anomaly_time)
    
    # Convert to integer timestamp
    if [[ "$LAST_ANOMALY_TS" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        LAST_ANOMALY_SEC=${LAST_ANOMALY_TS%.*}
    else
        LAST_ANOMALY_SEC=0
    fi
    
    # Calculate time since last anomaly
    if [ "$LAST_ANOMALY_SEC" -eq 0 ]; then
        # No anomalies yet - use session start as reference
        TIME_SINCE_ANOMALY=$((NOW - $(stat -f%B "$SYNTHETIC_LOG" 2>/dev/null || echo "$NOW")))
    else
        TIME_SINCE_ANOMALY=$((NOW - LAST_ANOMALY_SEC))
    fi
    
    # Check if we should inject
    if [ "$TIME_SINCE_ANOMALY" -gt "$NO_ANOMALY_THRESHOLD" ]; then
        # Avoid injecting too frequently (at least 30 min between injections)
        TIME_SINCE_INJECTION=$((NOW - LAST_INJECTION_TIME))
        
        if [ "$TIME_SINCE_INJECTION" -gt 1800 ] || [ "$LAST_INJECTION_TIME" -eq 0 ]; then
            inject_pair_switch "No anomalies detected for $((TIME_SINCE_ANOMALY / 60)) minutes"
            LAST_INJECTION_TIME=$NOW
        else
            echo "   Skipping injection - too soon since last injection" | tee -a "$SYNTHETIC_LOG"
        fi
    fi
done
