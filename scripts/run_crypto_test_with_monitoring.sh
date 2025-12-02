#!/bin/bash
# Run 4-hour crypto test with real-time anomaly monitoring
# This script starts both the test and the monitor in separate terminals/tabs
# Perfect for capturing anomalies on Loom!

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🚀 Starting 4-Hour Crypto Test with Anomaly Monitoring"
echo ""

# Check if we're on macOS (for terminal splitting)
if [[ "$OSTYPE" == "darwin"* ]]; then
    HAS_ITERM=$(command -v osascript >/dev/null 2>&1 && echo "yes" || echo "no")
else
    HAS_ITERM="no"
fi

# Load .env if it exists
if [ -f "$REPO_ROOT/.env" ]; then
    set -a
    source "$REPO_ROOT/.env"
    set +a
fi

# Check for API key
if [ -z "$DRIFTLOCK_API_KEY" ]; then
    echo "❌ DRIFTLOCK_API_KEY not set. Please:"
    echo "   1. Sign up at https://driftlock.web.app"
    echo "   2. Export your API key: export DRIFTLOCK_API_KEY='dlk_...'"
    exit 1
fi

export DRIFTLOCK_API_URL="${DRIFTLOCK_API_URL:-https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1}"

echo "✅ API Key: ${DRIFTLOCK_API_KEY:0:20}..."
echo "✅ API URL: $DRIFTLOCK_API_URL"
echo ""

# Start the 4-hour test in background
echo "📊 Starting crypto test (will run for 4 hours)..."
"$SCRIPT_DIR/start_crypto_test.sh" &
TEST_PID=$!

# Wait a moment for log file to be created
sleep 3

# Start the anomaly monitor
echo "🔍 Starting anomaly monitor..."
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "📹 SETUP YOUR LOOM RECORDING NOW!"
echo ""
echo "The monitor will alert you when anomalies are detected."
echo "When you see the alert, your Loom should already be recording!"
echo ""
echo "To monitor manually in another terminal:"
echo "   ./scripts/monitor_anomalies.sh"
echo ""
echo "To view the full log:"
echo "   tail -f logs/crypto-api-test-*.log"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Run the monitor in foreground (user can Ctrl+C to stop)
"$SCRIPT_DIR/monitor_anomalies.sh"

# Cleanup
echo ""
echo "Stopping test (PID: $TEST_PID)..."
kill $TEST_PID 2>/dev/null || true









