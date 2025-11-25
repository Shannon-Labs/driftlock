#!/usr/bin/env bash
# Start a virtual X11 display for headless screen capture
# This keeps recordings off the physical desktop and feeds the recorder with a stable view.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

VIRTUAL_DISPLAY="${VIRTUAL_DISPLAY:-:99}"
VIRTUAL_DISPLAY_GEOMETRY="${VIRTUAL_DISPLAY_GEOMETRY:-1920x1080x24}"
STATE_DIR="${VIRTUAL_DISPLAY_STATE_DIR:-$REPO_ROOT/logs/.virtual_display}"
ALERT_LOG="${ALERT_LOG:-$REPO_ROOT/logs/anomaly_alerts.log}"
FEED_CMD="${VIRTUAL_FEED_CMD:-$SCRIPT_DIR/virtual_display_feed.sh}"
FEED_TITLE="${VIRTUAL_FEED_TITLE:-Driftlock Virtual Feed}"
USE_PLAYWRIGHT_DASHBOARD="${USE_PLAYWRIGHT_DASHBOARD:-0}"

mkdir -p "$STATE_DIR" "$(dirname "$ALERT_LOG")"

require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "ERROR: Missing dependency: $1" >&2
        echo "   Install via brew (macOS): brew install ffmpeg xorg-server xterm" >&2
        echo "   Install via apt (Linux): sudo apt-get install -y ffmpeg xvfb x11-apps xterm" >&2
        exit 1
    fi
}

require_cmd ffmpeg
require_cmd Xvfb
require_cmd xterm

export DISPLAY="$VIRTUAL_DISPLAY"
SOCKET="/tmp/.X11-unix/X${VIRTUAL_DISPLAY#:}"

start_xvfb() {
    if [ -S "$SOCKET" ]; then
        echo "INFO: Xvfb already running on $VIRTUAL_DISPLAY (socket $SOCKET)"
        return
    fi

    echo "Starting Xvfb on $VIRTUAL_DISPLAY with geometry $VIRTUAL_DISPLAY_GEOMETRY"
    Xvfb "$VIRTUAL_DISPLAY" -screen 0 "$VIRTUAL_DISPLAY_GEOMETRY" -ac +extension RANDR \
        >"$STATE_DIR/xvfb.log" 2>&1 &
    echo $! >"$STATE_DIR/xvfb.pid"
    sleep 1
}

start_feed_window() {
    if [ ! -x "$FEED_CMD" ]; then
        echo "ERROR: Feed command not executable: $FEED_CMD" >&2
        exit 1
    fi

    if [ -f "$STATE_DIR/feed.pid" ] && kill -0 "$(cat "$STATE_DIR/feed.pid")" 2>/dev/null; then
        echo "INFO: Feed already running (PID $(cat "$STATE_DIR/feed.pid"))"
        return
    fi

    echo "Launching virtual feed window ($FEED_TITLE)"
    ALERT_LOG="$ALERT_LOG" DISPLAY="$VIRTUAL_DISPLAY" xterm -geometry 160x48 -bg black -fg green \
        -fa "Menlo" -fs 12 -T "$FEED_TITLE" -e "$FEED_CMD" &
    echo $! >"$STATE_DIR/feed.pid"
}

start_playwright_dashboard() {
    if [ "$USE_PLAYWRIGHT_DASHBOARD" != "1" ]; then
        return
    fi

    if ! command -v node >/dev/null 2>&1; then
        echo "WARNING: Node not found; skipping Playwright dashboard" >&2
        return
    fi

    if ! node -e "require('playwright')" >/dev/null 2>&1; then
        echo "WARNING: Playwright not installed. Install with: npm install --save-dev playwright" >&2
        return
    fi

    echo "Starting Playwright dashboard view"
    DISPLAY="$VIRTUAL_DISPLAY" node "$SCRIPT_DIR/virtual_dashboard.js" &
    echo $! >"$STATE_DIR/playwright.pid"
}

start_xvfb
start_feed_window
start_playwright_dashboard

echo ""
echo "Virtual display ready"
echo "   DISPLAY=$VIRTUAL_DISPLAY"
echo "   Geometry=$VIRTUAL_DISPLAY_GEOMETRY"
echo "   State dir: $STATE_DIR"
echo "   Feed: $FEED_CMD"
echo "   Alert log: $ALERT_LOG"
echo ""
echo "Use these to record: CAPTURE_BACKEND=x11 SCREEN_DEVICE=$VIRTUAL_DISPLAY DISPLAY=$VIRTUAL_DISPLAY"
