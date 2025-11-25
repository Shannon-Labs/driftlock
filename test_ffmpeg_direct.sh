#!/bin/bash
# Test ffmpeg directly with full error output

OUTPUT_FILE="logs/session_1763871492/recordings/manual_test_$(date +%s).mp4"
mkdir -p logs/session_1763871492/recordings

echo "Testing ffmpeg screen recording..."
echo "Output: $OUTPUT_FILE"
echo ""

ffmpeg -f avfoundation -framerate 30 -i "4:none" \
    -t 5 \
    -s 1920x1080 \
    -c:v libx264 \
    -preset ultrafast \
    -crf 23 \
    -pix_fmt uyvy422 \
    "$OUTPUT_FILE" 2>&1

echo ""
echo "Exit code: $?"
echo "File exists:"
ls -lah "$OUTPUT_FILE" 2>&1
