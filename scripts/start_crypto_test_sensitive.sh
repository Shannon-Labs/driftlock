#!/bin/bash
# Start 4-hour crypto test in SENSITIVE MODE
# This increases the likelihood of detecting anomalies for demos/recordings
# Use this when you want to capture anomalies on Loom!

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🚀 Starting 4-Hour Crypto Test (SENSITIVE MODE)"
echo "   This mode uses more sensitive detection settings to increase"
echo "   the likelihood of detecting anomalies for demos/recordings."
echo ""

# Load .env if it exists
if [ -f "$REPO_ROOT/.env" ]; then
    echo "📥 Loading .env file..."
    set -a
    source "$REPO_ROOT/.env"
    set +a
fi

# Check for API key
if [ -z "$DRIFTLOCK_API_KEY" ]; then
    echo "⚠️  No API key found. Creating one automatically..."
    echo ""
    
    if "$SCRIPT_DIR/create-test-api-key-cloudrun.sh"; then
        if [ -f "$REPO_ROOT/.env" ]; then
            set -a
            source "$REPO_ROOT/.env"
            set +a
        fi
        echo "✅ API key created and loaded!"
    else
        echo ""
        echo "❌ Failed to create API key automatically."
        echo ""
        echo "Please either:"
        echo "  1. Run manually: ./scripts/create-test-api-key-cloudrun.sh"
        echo "  2. Or sign up at https://driftlock.web.app and set:"
        echo "     export DRIFTLOCK_API_KEY='dlk_...'"
        exit 1
    fi
fi

export DRIFTLOCK_API_URL="${DRIFTLOCK_API_URL:-https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1}"

echo "✅ API Key: ${DRIFTLOCK_API_KEY:0:20}..."
echo "✅ API URL: $DRIFTLOCK_API_URL"
echo ""
echo "⚠️  SENSITIVE MODE ENABLED"
echo "   Settings: window_size=20, baseline_lines=40, ncd_threshold=0.25"
echo "   This will detect more anomalies than normal mode"
echo ""
echo "📊 Starting 4-hour crypto test..."
echo "   Source: ${CRYPTO_SOURCE:-coingecko}"
echo "   Logs will be saved to: logs/crypto-api-test-*.log"
echo ""
echo "💡 TIP: Run in another terminal to monitor for anomalies:"
echo "   ./scripts/monitor_anomalies.sh"
echo ""

# Create logs directory
mkdir -p "$REPO_ROOT/logs"

# Timestamp for this run
TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")
LOG_FILE="$REPO_ROOT/logs/crypto-api-test-${TIMESTAMP}.log"
PID_FILE="$REPO_ROOT/logs/crypto-api-test-${TIMESTAMP}.pid"

# Function to cleanup on exit
cleanup() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            echo ""
            echo "Stopping test (PID: $PID)..."
            kill "$PID" 2>/dev/null || true
            wait "$PID" 2>/dev/null || true
        fi
        rm -f "$PID_FILE"
    fi
    echo "✅ Test completed at $(date)"
}

trap cleanup EXIT INT TERM

# Start the sensitive test script in background
echo "📊 Starting crypto data stream (SENSITIVE MODE)..."
python3 "$SCRIPT_DIR/api_crypto_test_sensitive.py" \
    --api-key "$DRIFTLOCK_API_KEY" \
    --api-url "$DRIFTLOCK_API_URL" \
    > "$LOG_FILE" 2>&1 &

TEST_PID=$!
echo $TEST_PID > "$PID_FILE"

echo "✅ Test started (PID: $TEST_PID)"
echo "   Log file: $LOG_FILE"
echo "   To monitor: tail -f $LOG_FILE"
echo "   To monitor anomalies: ./scripts/monitor_anomalies.sh"
echo "   To stop: kill $TEST_PID or Ctrl+C"
echo ""

# Wait for 4 hours (14400 seconds)
START_TIME=$(date +%s)
END_TIME=$((START_TIME + 14400))

echo "⏳ Running for 4 hours..."
echo "   Will stop at: $(date -v+4H 2>/dev/null || date -d '+4 hours' 2>/dev/null || echo '4 hours from now')"
echo "   Press Ctrl+C to stop early"
echo ""

while [ $(date +%s) -lt $END_TIME ]; do
    sleep 60
    if ! ps -p "$TEST_PID" > /dev/null 2>&1; then
        echo "⚠️  Test process ended early (PID: $TEST_PID)"
        break
    fi
done

echo ""
echo "⏰ 4 hours elapsed. Stopping test..."
cleanup

echo ""
echo "📊 Test Summary:"
echo "   Log file: $LOG_FILE"
if [ -f "$LOG_FILE" ]; then
    echo "   Total batches: $(grep -c "Batch sent" "$LOG_FILE" 2>/dev/null || echo "0")"
    echo "   Anomalies detected: $(grep -c "ANOMALY DETECTED" "$LOG_FILE" 2>/dev/null || echo "0")"
fi
echo ""









