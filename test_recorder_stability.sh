#!/bin/bash
set -euo pipefail

echo "Starting test at $(date)"

# Simulate what auto_screen_recorder does
SESSION_DIR=logs/session_debug
RECORDING_DIR=$SESSION_DIR/recordings
mkdir -p $RECORDING_DIR

# Simulate recording
ffmpeg -f avfoundation -framerate 30 -i "2:none" \
    -t 10 \
    -s 1920x1080 \
    -c:v libx264 \
    -preset ultrafast \
    -crf 23 \
    -pix_fmt uyvy422 \
    "$RECORDING_DIR/test_debug.mp4" \
    >/dev/null 2>&1 &

FFMPEG_PID=$!
echo "FFmpeg PID: $FFMPEG_PID"

# The script should not exit here, but continue monitoring
sleep 2
if kill -0 $FFMPEG_PID 2>/dev/null; then
    echo "FFmpeg still running after 2s: YES"
else
    echo "FFmpeg still running after 2s: NO"
fi

# Wait for it to finish
wait $FFMPEG_PID
echo "FFmpeg finished"
echo "Test file size: $(du -h $RECORDING_DIR/test_debug.mp4 | cut -f1)"
