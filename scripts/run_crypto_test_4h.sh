#!/bin/bash
# Run Driftlock Crypto API Test for 4 Hours
# Streams live crypto market data to the Driftlock API for 4 hours and logs
# all activity for analysis. CoinGecko (no API key) is the default source;
# set CRYPTO_SOURCE=binance to use Binance WebSocket when available.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG_DIR="$REPO_ROOT/logs"
TEST_SCRIPT="$SCRIPT_DIR/api_crypto_test.py"

# Create logs directory if it doesn't exist
mkdir -p "$LOG_DIR"

# Timestamp for this run
TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")
LOG_FILE="$LOG_DIR/crypto-api-test-${TIMESTAMP}.log"
PID_FILE="$LOG_DIR/crypto-api-test-${TIMESTAMP}.pid"

# Check for API key
if [ -z "$DRIFTLOCK_API_KEY" ]; then
    echo "❌ Error: DRIFTLOCK_API_KEY environment variable not set"
    echo ""
    echo "To get an API key:"
    echo "  1. Sign up at https://driftlock.web.app"
    echo "  2. Or use an existing API key"
    echo ""
    echo "Then run:"
    echo "  export DRIFTLOCK_API_KEY='dlk_...'"
    echo "  export DRIFTLOCK_API_URL='https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1'"
    echo "  $0"
    exit 1
fi

# Default API URL if not set
export DRIFTLOCK_API_URL="${DRIFTLOCK_API_URL:-https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1}"

echo "🚀 Starting 4-Hour Driftlock Crypto API Test"
echo "   Started: $(date)"
echo "   Will run until: $(date -v+4H 2>/dev/null || date -d '+4 hours' 2>/dev/null || echo '4 hours from now')"
echo "   API: $DRIFTLOCK_API_URL"
echo "   Source: ${CRYPTO_SOURCE:-coingecko} (CRYPTO_SOURCE=binance for WS)"
echo "   Coins (CoinGecko): ${COINGECKO_IDS:-bitcoin,ethereum,solana,chainlink,avalanche-2,dogecoin,litecoin}"
if [ -n "$COINGECKO_API_KEY" ]; then
    echo "   CoinGecko API key: set"
fi
echo "   Log: $LOG_FILE"
echo "   PID: $PID_FILE"
echo ""

# Check if Python script exists
if [ ! -f "$TEST_SCRIPT" ]; then
    echo "❌ Error: Test script not found at $TEST_SCRIPT"
    exit 1
fi

# Check if Python dependencies are installed
if ! python3 -c "import websockets, requests" 2>/dev/null; then
    echo "⚠️  Warning: Missing Python dependencies"
    echo "   Installing: pip install websockets requests"
    pip3 install websockets requests certifi || {
        echo "❌ Failed to install dependencies"
        exit 1
    }
fi

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

# Start the test script in background, logging to file
echo "📊 Starting crypto data stream..."
python3 "$TEST_SCRIPT" \
    --api-key "$DRIFTLOCK_API_KEY" \
    --api-url "$DRIFTLOCK_API_URL" \
    > "$LOG_FILE" 2>&1 &

TEST_PID=$!
echo $TEST_PID > "$PID_FILE"

echo "✅ Test started (PID: $TEST_PID)"
echo "   Monitoring log: tail -f $LOG_FILE"
echo "   To stop: kill $TEST_PID or Ctrl+C"
echo ""

# Wait for 4 hours (14400 seconds) or until interrupted
echo "⏳ Running for 4 hours (will stop at $(date -v+4H 2>/dev/null || date -d '+4 hours' 2>/dev/null || echo '4 hours from now'))..."
echo "   You can monitor progress with: tail -f $LOG_FILE"
echo "   Press Ctrl+C to stop early"
echo ""

# Wait for 4 hours, but allow early termination
START_TIME=$(date +%s)
END_TIME=$((START_TIME + 14400))

while [ $(date +%s) -lt $END_TIME ]; do
    sleep 60  # Check every minute
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
    echo "   Total events: $(grep -c "Batch sent" "$LOG_FILE" 2>/dev/null || echo "0") batches"
    echo "   Anomalies detected: $(grep -c "anomalies detected" "$LOG_FILE" 2>/dev/null || echo "0")"
fi
echo ""
