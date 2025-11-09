#!/bin/bash
set -euo pipefail

# Ensure the demo output exists
if [ ! -f "demo-output.html" ]; then
  echo "demo-output.html not found. Running demo..."
  ./driftlock-demo test-data/financial-demo.json >/dev/null 2>&1 || {
    echo "Failed to run demo. Build first with 'make demo'."; exit 1;
  }
fi

echo "Opening Safari and capturing first anomaly card..."
osascript scripts/capture-anomaly-card.applescript || {
  echo "AppleScript capture failed. Please follow docs/CAPTURE-ANOMALY-SCREENSHOT.md for manual capture."; exit 1;
}

if [ -f screenshots/demo-anomaly-card.png ]; then
  echo "✅ Saved: screenshots/demo-anomaly-card.png"
else
  echo "❌ Capture failed. See docs/CAPTURE-ANOMALY-SCREENSHOT.md"
  exit 1
fi

