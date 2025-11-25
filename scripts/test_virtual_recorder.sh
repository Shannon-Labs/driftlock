#!/usr/bin/env bash
# Synthetic test to prove the virtual display and recorder pipeline produce a readable MP4.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

VIRTUAL_DISPLAY="${VIRTUAL_DISPLAY:-:99}"
SESSION_DIR="${SESSION_DIR:-$REPO_ROOT/logs/session_$(date +%s)_virtual_test}"
ALERT_LOG="${ALERT_LOG:-$SESSION_DIR/anomaly_alerts.log}"
RECORDING_DIR="${RECORDING_DIR:-$SESSION_DIR/recordings}"
STATE_DIR="$SESSION_DIR/.virtual_display"

mkdir -p "$SESSION_DIR" "$RECORDING_DIR" "$(dirname "$ALERT_LOG")"

export CAPTURE_BACKEND=x11
export VIRTUAL_DISPLAY
export SCREEN_DEVICE="$VIRTUAL_DISPLAY"
export SESSION_DIR ALERT_LOG RECORDING_DIR

cleanup() {
    if [ -n "${REC_PID:-}" ] && kill -0 "$REC_PID" 2>/dev/null; then
        kill "$REC_PID" 2>/dev/null || true
    fi
    VIRTUAL_DISPLAY_STATE_DIR="$STATE_DIR" "$SCRIPT_DIR/stop_virtual_display.sh" >/dev/null 2>&1 || true
}
trap cleanup EXIT

VIRTUAL_DISPLAY_STATE_DIR="$STATE_DIR" ALERT_LOG="$ALERT_LOG" "$SCRIPT_DIR/start_virtual_display.sh"

echo "Starting recorder against virtual display $VIRTUAL_DISPLAY"
"$SCRIPT_DIR/auto_screen_recorder.sh" >/tmp/virtual_recorder_test.log 2>&1 &
REC_PID=$!

sleep 3

echo "Triggering synthetic alert"
echo "$(date '+%F %T') Anomaly ID: synthetic-test virtual-capture-check" >> "$ALERT_LOG"

echo "Waiting for recording to be created in $RECORDING_DIR"
TARGET_FILE=""
for _ in $(seq 1 20); do
    TARGET_FILE=$(ls -t "$RECORDING_DIR"/anomaly_synthetic-test_* 2>/dev/null | head -n1 || true)
    if [ -n "$TARGET_FILE" ] && [ -s "$TARGET_FILE" ]; then
        break
    fi
    sleep 2
done

if [ -z "$TARGET_FILE" ] || [ ! -s "$TARGET_FILE" ]; then
    echo "ERROR: No recording produced. Check /tmp/virtual_recorder_test.log for ffmpeg output." >&2
    exit 1
fi

echo "Recording created: $TARGET_FILE (size: $(du -h "$TARGET_FILE" | cut -f1))"
exit 0

