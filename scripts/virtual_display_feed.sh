#!/usr/bin/env bash
# Simple terminal UI shown on the virtual display
# Continuously tails the anomaly alert log so the capture has meaningful content.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

ALERT_LOG="${ALERT_LOG:-$REPO_ROOT/logs/anomaly_alerts.log}"
TAIL_LINES="${TAIL_LINES:-30}"
REFRESH_SECONDS="${REFRESH_SECONDS:-1}"
RECENT_TRADES_LINES="${RECENT_TRADES_LINES:-5}"

mkdir -p "$(dirname "$ALERT_LOG")"
touch "$ALERT_LOG"

while true; do
    clear
    echo "[Driftlock Virtual Feed]"
    echo "Log: $ALERT_LOG"
    echo "Display: ${DISPLAY:-unset}"
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "------------------------------------------------------------"

    if [ -s "$ALERT_LOG" ]; then
        tail -n "$TAIL_LINES" "$ALERT_LOG"
    else
        echo "Waiting for anomalies... (write 'Anomaly ID: <id>' lines to trigger)"
    fi

    # Show latest live trades from active sessions (ETH/SOL) to keep the feed moving
    ETH_RAW=$(ls -t "$REPO_ROOT"/logs/session_kraken_eth_*/kraken_raw.ndjson 2>/dev/null | head -1 || true)
    SOL_RAW=$(ls -t "$REPO_ROOT"/logs/session_kraken_sol_*/kraken_raw.ndjson 2>/dev/null | head -1 || true)

    echo ""
    echo "--- Recent ETH/USD trades ---"
    if [ -n "$ETH_RAW" ]; then
        tail -n "$RECENT_TRADES_LINES" "$ETH_RAW"
    else
        echo "(no ETH feed detected)"
    fi

    echo ""
    echo "--- Recent SOL/USD trades ---"
    if [ -n "$SOL_RAW" ]; then
        tail -n "$RECENT_TRADES_LINES" "$SOL_RAW"
    else
        echo "(no SOL feed detected)"
    fi

    sleep "$REFRESH_SECONDS"
done
