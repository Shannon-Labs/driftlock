#!/usr/bin/env bash
# Stop the virtual X11 display and helper processes started by start_virtual_display.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

STATE_DIR="${VIRTUAL_DISPLAY_STATE_DIR:-$REPO_ROOT/logs/.virtual_display}"

stop_pid() {
    local name="$1"
    local file="$STATE_DIR/$2"
    if [ -f "$file" ]; then
        local pid
        pid=$(cat "$file")
        if kill -0 "$pid" 2>/dev/null; then
            echo "Stopping $name (PID $pid)"
            kill "$pid" 2>/dev/null || true
        fi
        rm -f "$file"
    fi
}

stop_pid "Playwright dashboard" "playwright.pid"
stop_pid "Feed window" "feed.pid"
stop_pid "Xvfb" "xvfb.pid"

echo "Virtual display processes stopped (state dir: $STATE_DIR)"

