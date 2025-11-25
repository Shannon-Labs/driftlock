#!/bin/bash
# Real-time anomaly monitor for crypto test
# Watches the log file and alerts when anomalies are detected
# Perfect for triggering Loom recordings!

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG_DIR="$REPO_ROOT/logs"

# Find the most recent log file
LATEST_LOG=$(ls -t "$LOG_DIR"/crypto-api-test-*.log 2>/dev/null | head -n1)

if [ -z "$LATEST_LOG" ]; then
    echo "❌ No log file found. Start the crypto test first:"
    echo "   ./scripts/start_crypto_test.sh"
    exit 1
fi

echo "🔍 Monitoring for anomalies in: $LATEST_LOG"
echo "   Press Ctrl+C to stop monitoring"
echo ""

# Track if we've seen anomalies
ANOMALY_COUNT=0
LAST_ANOMALY_TIME=""

# Function to play alert sound (macOS)
play_alert() {
    # Try different alert methods
    if command -v say >/dev/null 2>&1; then
        say "Anomaly detected!" 2>/dev/null || true
    fi
    if command -v afplay >/dev/null 2>&1; then
        afplay /System/Library/Sounds/Glass.aiff 2>/dev/null || true
    fi
    # Visual alert
    echo ""
    echo "🚨🚨🚨 ANOMALY DETECTED! 🚨🚨🚨"
    echo ""
}

# Monitor the log file in real-time
tail -f "$LATEST_LOG" | while IFS= read -r line; do
    # Check for anomaly detection (both standard and sensitive mode patterns)
    if echo "$line" | grep -qE "(anomalies detected|ANOMALY DETECTED)"; then
        ANOMALY_COUNT=$((ANOMALY_COUNT + 1))
        LAST_ANOMALY_TIME=$(date '+%Y-%m-%d %H:%M:%S')
        
        # Extract anomaly count from log line (try both patterns)
        ANOMALY_NUM=$(echo "$line" | grep -oE '[0-9]+ anomalies' | grep -oE '[0-9]+' | head -n1 || echo "1")
        if [ -z "$ANOMALY_NUM" ]; then
            ANOMALY_NUM=$(echo "$line" | grep -oE '[0-9]+ → [0-9]+' | grep -oE '[0-9]+$' | head -n1 || echo "1")
        fi
        
        echo ""
        echo "═══════════════════════════════════════════════════════════"
        echo "🚨 ANOMALY DETECTED! ($ANOMALY_COUNT total)"
        echo "   Time: $LAST_ANOMALY_TIME"
        echo "   Count: $ANOMALY_NUM anomalies in this batch"
        echo "   Log line: $line"
        echo ""
        echo "📹 START YOUR LOOM RECORDING NOW!"
        echo "   The anomaly details are in the log above this line."
        echo "═══════════════════════════════════════════════════════════"
        echo ""
        
        # Play alert
        play_alert
        
        # Show the last few lines for context
        echo "📋 Recent log context:"
        tail -n 20 "$LATEST_LOG" | grep -A 5 -B 5 -E "(anomalies detected|ANOMALY DETECTED)" | tail -n 15 || true
        echo ""
        
    # Also check for anomaly explanation lines
    elif echo "$line" | grep -qE "(Anomaly:|🎯 Anomaly)"; then
        echo "   $line"
    fi
done

