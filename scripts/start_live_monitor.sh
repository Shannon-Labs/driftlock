#!/bin/bash
# Start live crypto monitoring with driftlock CLI
# Launches streamer + driftlock with anomaly alerts

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Create logs directory
mkdir -p "$REPO_ROOT/logs"

# Timestamp
TIMESTAMP=$(date -u +"%Y%m%dT%H%M%SZ")
STREAM_LOG="$REPO_ROOT/logs/streamer-$TIMESTAMP.log"
ANOMALY_LOG="$REPO_ROOT/logs/anomalies-$TIMESTAMP.log"
PID_FILE="$REPO_ROOT/logs/monitor-$TIMESTAMP.pid"

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║        LIVE CRYPTO ANOMALY MONITORING WITH DRIFTLOCK        ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "📄 Logs:"
echo "   Streamer: $STREAM_LOG"
echo "   Anomalies: $ANOMALY_LOG"
echo "   PID: $PID_FILE"
echo ""
echo "📊 Configuration:"
echo "   • Source: CoinGecko (5 coins)"
echo "   • Synthetic spike: every 10 batches"
echo "   • Driftlock: baseline-lines=40, algo=zstd, entropy detection"
echo "   • Alert: 🚨 on anomaly detection + Loom prompt"
echo ""

# Launch pipeline in background
echo "🚀 Starting monitoring pipeline..."
{
  python3 -u "$SCRIPT_DIR/stream_ndjson_simple.py" \
    --interval 5 \
    --synthetic-every 10 \
    2>"$STREAM_LOG" | \
  "$REPO_ROOT/bin/driftlock" scan \
    --stdin \
    --follow \
    --format ndjson \
    --baseline-lines 40 \
    --algo zstd \
    --show-all \
    2>&1 | \
  while IFS= read -r line; do
    echo "$line"
    # Check for anomaly patterns
    if echo "$line" | grep -qE "(ANOMALY|anomaly|NCD|compression|entropy)"; then
      echo "$line" >> "$ANOMALY_LOG"
      # Alert for actual anomalies (not just stats)
      if echo "$line" | grep -qE "(NCD:|compression ratio:|entropy:)"; then
        echo "═══════════════════════════════════════════════════════════"
        echo "🚨 ANOMALY DETECTED IN CRYPTO STREAM! 🚨"
        echo "Time: $(date '+%Y-%m-%d %H:%M:%S')"
        echo "Log: $ANOMALY_LOG"
        echo ""
        echo "📹 START YOUR LOOM RECORDING NOW!"
        echo "═══════════════════════════════════════════════════════════"
        # Play audio alert
        say "Anomaly detected in crypto stream! Start Loom now!" 2>/dev/null || true
        afplay /System/Library/Sounds/Glass.aiff 2>/dev/null || true
      fi
    fi
  done
} &

MONITOR_PID=$!
echo $MONITOR_PID > "$PID_FILE"

sleep 2
if ps -p "$MONITOR_PID" > /dev/null 2>&1; then
  echo "✅ Pipeline started successfully!"
  echo "   Monitor PID: $MONITOR_PID"
  echo ""
  echo "🔍 To watch for anomalies:"
  echo "   tail -f $ANOMALY_LOG"
  echo ""
  echo "📊 To watch full stream:"
  echo "   tail -f $STREAM_LOG"
  echo ""
  echo "⏹️  To stop:"
  echo "   kill $MONITOR_PID"
  echo "   or: kill $(cat $PID_FILE)"
else
  echo "❌ Failed to start pipeline"
  exit 1
fi

# Keep script running
trap "kill $MONITOR_PID 2>/dev/null; rm -f $PID_FILE; echo 'Stopped.'; exit 0" INT TERM
wait $MONITOR_PID