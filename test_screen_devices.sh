#!/bin/bash
echo "Testing screen capture devices..."
for i in 0 1 2 3 4; do
    echo -e "\n=== Device $i ==="
    timeout 3 ffmpeg -f avfoundation -i "${i}:none" -t 1 -f null - 2>&1 | head -8
done
