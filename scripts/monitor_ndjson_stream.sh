#!/bin/bash
# Monitor NDJSON stream with driftlock scan
# Launches the streamer and driftlock CLI together

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Create logs directory
mkdir -p "$REPO_ROOT/logs"

# Timestamp for this run
TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")
LOG_FILE="$REPO_ROOT/logs/ndjson-driftlock-$TIMESTAMP.log"

# Default settings
INTERVAL=5
SYNTHETIC_EVERY=10  # Inject synthetic anomaly every 10 batches
BASELINE_LINES=120

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║        LIVE CRYPTO MONITORING WITH DRIFTLOCK CLI            ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "📄 Log file: $LOG_FILE"
echo "📊 Source: CoinGecko REST (5 coins)"
echo "🎯 Synthetic spikes: every $SYNTHETIC_EVERY batches"
echo "🎛️  Driftlock: baseline-lines=$BASELINE_LINES, algo=entropy, show-all"
echo ""
echo "🚀 Starting pipeline..."
echo ""

# Launch the pipeline: streamer | driftlock | tee log
python3 -u "$SCRIPT_DIR/stream_ndjson_simple.py" \
  --interval "$INTERVAL" \
  --synthetic-every "$SYNTHETIC_EVERY" \
  2>"$LOG_FILE.streamer" | \
 "$REPO_ROOT/driftlock-demo" scan \
  --stdin \
  --follow \
  --format ndjson \
  --baseline-lines "$BASELINE_LINES" \
  --algo entropy \
  --show-all \
  2>&1 | tee "$LOG_FILE"

# Cleanup on exit
echo ""
echo "✅ Monitoring complete"
echo "📊 Statistics:"
echo "   Streamer log: $LOG_FILE.streamer"
echo "   Driftlock log: $LOG_FILE"
if [ -f "$LOG_FILE.streamer" ]; then
  TRADES=$(grep -c "trade" "$LOG_FILE.streamer" 2>/dev/null || echo "0")
  SYNTHETIC=$(grep -c "synthetic" "$LOG_FILE.streamer" 2>/dev/null || echo "0")
  echo "   Total trades: $TRADES"
  echo "   Synthetic spikes: $SYNTHETIC"
fi
if [ -f "$LOG_FILE" ]; then
  ANOMALIES=$(grep -c "ANOMALY\|anomaly" "$LOG_FILE" 2>/dev/null || echo "0")
  echo "   Anomalies detected: $ANOMALIES"
fi