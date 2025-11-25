#!/bin/bash
# Automatic screen recording on anomaly detection
# Supports virtual X11 display capture (default) or macOS avfoundation fallback

set -uo pipefail

# Configuration
SESSION_DIR="${SESSION_DIR:-logs/session_$(date +%s)}"
ALERT_LOG="${ALERT_LOG:-logs/anomaly_alerts.log}"
RECORDING_DIR="${RECORDING_DIR:-$SESSION_DIR/recordings}"
RECORDING_DURATION="${RECORDING_DURATION:-60}"  # 60 seconds per anomaly
SCREEN_RESOLUTION="${SCREEN_RESOLUTION:-1920x1080}"
FRAMERATE="${FRAMERATE:-30}"
CAPTURE_BACKEND="${CAPTURE_BACKEND:-x11}" # x11 (virtual) or avfoundation (macOS)
VIRTUAL_DISPLAY="${VIRTUAL_DISPLAY:-:99}"
PIX_FMT_X11="${PIX_FMT_X11:-yuv420p}"
PIX_FMT_AVFOUNDATION="${PIX_FMT_AVFOUNDATION:-uyvy422}"
# macOS fallback device (kept for compatibility)
DEFAULT_DEVICE="${SCREEN_DEVICE:-3:none}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

# Check if ffmpeg is available
if ! command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${RED}❌ ffmpeg not found. Install with: brew install ffmpeg${NC}"
    exit 1
fi

mkdir -p "$RECORDING_DIR"
mkdir -p "$(dirname "$ALERT_LOG")"
touch "$ALERT_LOG"

detect_screen_device() {
    # Only used for macOS avfoundation fallback
    echo "🔍 Detecting screen devices..." >&2
    DEVICES_OUTPUT=$(ffmpeg -f avfoundation -list_devices true -i "" 2>&1 || true)
    
    if echo "$DEVICES_OUTPUT" | grep -q "\\[3\\]"; then
        echo -e "${GREEN}✅ Found preferred device 3${NC}" >&2
        echo "3:none"
        return
    fi

    DETECTED=$(echo "$DEVICES_OUTPUT" | grep "Capture screen" | head -1 | awk -F'[][]' '{print $2}')
    
    if [ -n "$DETECTED" ]; then
        echo -e "${YELLOW}⚠️  Device 3 not found. Falling back to detected device $DETECTED${NC}" >&2
        echo "${DETECTED}:none"
    else
        echo -e "${RED}❌ No suitable screen device found. Defaulting to $DEFAULT_DEVICE${NC}" >&2
        echo "$DEFAULT_DEVICE"
    fi
}

# Build ffmpeg input arguments for the selected backend
build_capture_args() {
    local args=()
    case "$CAPTURE_BACKEND" in
        x11)
            local x_display="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
            args=(-f x11grab -framerate "$FRAMERATE" -video_size "$SCREEN_RESOLUTION" -i "${x_display}+0,0")
            ;;
        avfoundation)
            args=(-f avfoundation -framerate "$FRAMERATE" -i "$SCREEN_DEVICE" -s "$SCREEN_RESOLUTION")
            ;;
        *)
            echo -e "${YELLOW}⚠️  Unknown CAPTURE_BACKEND=$CAPTURE_BACKEND, falling back to x11${NC}" >&2
            local x_display="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
            args=(-f x11grab -framerate "$FRAMERATE" -video_size "$SCREEN_RESOLUTION" -i "${x_display}+0,0")
            ;;
    esac
    printf '%s\n' "${args[@]}"
}

# Resolve the device for the chosen backend
case "$CAPTURE_BACKEND" in
    x11)
        SCREEN_DEVICE="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
        ;;
    avfoundation)
        SCREEN_DEVICE=$(detect_screen_device | tail -n 1)
        ;;
    *)
        CAPTURE_BACKEND="x11"
        SCREEN_DEVICE="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
        ;;
esac

if [ "$CAPTURE_BACKEND" = "x11" ]; then
    export DISPLAY="$SCREEN_DEVICE"
    socket="/tmp/.X11-unix/X${SCREEN_DEVICE#:}"
    if [ ! -S "$socket" ]; then
        echo -e "${YELLOW}⚠️  X11 display $SCREEN_DEVICE not found. Start it with scripts/start_virtual_display.sh${NC}" >&2
    fi
fi

echo "🎬 Auto Screen Recorder Started - $(date '+%Y-%m-%d %H:%M:%S')"
echo "   Alert Log: $ALERT_LOG"
echo "   Recording Dir: $RECORDING_DIR"
echo "   Duration per anomaly: ${RECORDING_DURATION}s"
echo "   Backend: ${CAPTURE_BACKEND}"
echo "   Screen device: ${SCREEN_DEVICE}"
echo ""

# Simple audible alert (macOS-friendly). Non-blocking and best-effort.
play_alert() {
    if command -v say >/dev/null 2>&1; then
        say "Anomaly detected" >/dev/null 2>&1 || true
    fi
    if command -v afplay >/dev/null 2>&1; then
        afplay /System/Library/Sounds/Glass.aiff >/dev/null 2>&1 || true
    fi
}

# Track last processed line
LAST_LINE=$(wc -l < "$ALERT_LOG" 2>/dev/null | tr -d ' ' || echo 0)
RECORDING_PID=""

# Function to start recording
start_recording() {
    local anomaly_id="$1"
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local output_file="$RECORDING_DIR/anomaly_${anomaly_id}_${timestamp}.mp4"
    local log_file="/tmp/ffmpeg_record_${anomaly_id}.log"
    
    echo -e "${GREEN}🎥 RECORDING STARTED${NC}"
    echo "   Anomaly ID: $anomaly_id"
    echo "   Output: $output_file"
    echo "   Duration: ${RECORDING_DURATION}s"
    play_alert
    
    # Retry logic for recording start
    local max_retries=3
    local retry_count=0
    local success=false

    while [ $retry_count -lt $max_retries ]; do
        local pix_fmt="$PIX_FMT_X11"
        if [ "$CAPTURE_BACKEND" = "avfoundation" ]; then
            pix_fmt="$PIX_FMT_AVFOUNDATION"
        fi

        # Build capture args inline since mapfile is not available
        case "$CAPTURE_BACKEND" in
            x11)
                local x_display="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
                ffmpeg -f x11grab -framerate "$FRAMERATE" -video_size "$SCREEN_RESOLUTION" -i "${x_display}+0,0" \
                    -t "$RECORDING_DURATION" \
                    -c:v libx264 \
                    -preset ultrafast \
                    -crf 23 \
                    -pix_fmt "$pix_fmt" \
                    -y \
                    "$output_file" \
                    > "$log_file" 2>&1 &
                ;;
            avfoundation)
                ffmpeg -f avfoundation -framerate "$FRAMERATE" -i "$SCREEN_DEVICE" -s "$SCREEN_RESOLUTION" \
                    -t "$RECORDING_DURATION" \
                    -c:v libx264 \
                    -preset ultrafast \
                    -crf 23 \
                    -pix_fmt "$pix_fmt" \
                    -y \
                    "$output_file" \
                    > "$log_file" 2>&1 &
                ;;
            *)
                # Fallback to x11
                local x_display="${SCREEN_DEVICE:-$VIRTUAL_DISPLAY}"
                ffmpeg -f x11grab -framerate "$FRAMERATE" -video_size "$SCREEN_RESOLUTION" -i "${x_display}+0,0" \
                    -t "$RECORDING_DURATION" \
                    -c:v libx264 \
                    -preset ultrafast \
                    -crf 23 \
                    -pix_fmt "$pix_fmt" \
                    -y \
                    "$output_file" \
                    > "$log_file" 2>&1 &
                ;;
        esac
        
        RECORDING_PID=$!
        
        # Give it a moment to fail if it's going to fail immediately
        sleep 2
        if kill -0 "$RECORDING_PID" 2>/dev/null; then
            echo "   PID: $RECORDING_PID"
            success=true
            break
        else
            echo -e "${YELLOW}⚠️  Recording failed to start (Attempt $((retry_count+1))/$max_retries)${NC}"
            if [ -f "$log_file" ]; then
                cat "$log_file"
            else
                echo "   Log file not found: $log_file"
            fi
            retry_count=$((retry_count+1))
        fi
    done

    if [ "$success" = false ]; then
        echo -e "${RED}❌ Failed to start recording after $max_retries attempts.${NC}"
        return
    fi
    
    echo ""
    
    # Wait for recording to finish, then report
    (
        wait "$RECORDING_PID" 2>/dev/null
        local size=$(du -h "$output_file" 2>/dev/null | cut -f1 || echo "N/A")
        if [ -f "$output_file" ] && [ -s "$output_file" ]; then
            echo -e "${CYAN}✅ Recording complete: $output_file (${size})${NC}"
        else
            echo -e "${RED}❌ Recording failed (file empty or missing): $output_file${NC}"
            echo "   See $log_file for details."
        fi
    ) &
}

# Main monitoring loop
echo "Monitoring for anomaly alerts..."
echo "Press Ctrl+C to exit"
echo ""

while true; do
    sleep 2
    
    if [ ! -f "$ALERT_LOG" ]; then
        continue
    fi
    
    # Get current line count
    CURRENT_LINES=$(wc -l < "$ALERT_LOG" | tr -d ' ')
    
    # Check if there are new alerts
    if [ "$CURRENT_LINES" -gt "$LAST_LINE" ]; then
        # Process new lines - look for Anomaly IDs
        NEW_ALERTS=$(tail -n +$((LAST_LINE + 1)) "$ALERT_LOG" | grep "Anomaly ID:" || true)
        
        if [ -n "$NEW_ALERTS" ]; then
            # Extract anomaly ID from alert
            # Handle multiple alerts by taking the last one or iterating (taking last for now to avoid overlap)
            ANOMALY_ID=$(echo "$NEW_ALERTS" | tail -1 | sed 's/.*Anomaly ID: //' | awk '{print $1}')
            
            # Check if already recording
            if [ -n "$RECORDING_PID" ] && kill -0 "$RECORDING_PID" 2>/dev/null; then
                echo -e "${YELLOW}⚠️  Already recording (PID: $RECORDING_PID), skipping anomaly $ANOMALY_ID${NC}"
            else
                start_recording "$ANOMALY_ID"
            fi
        fi
        
        LAST_LINE=$CURRENT_LINES
    fi
done
