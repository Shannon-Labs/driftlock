# Virtual Screen Capture for Driftlock

This replaces the macOS `avfoundation 3:none` flow with a headless X11 screen that is always safe to record. The existing recorder pipeline stays intact: alerts → `scripts/auto_screen_recorder.sh` → MP4 in `logs/session_<ts>/recordings`.

## Dependencies
- ffmpeg
- Xvfb
- xterm
- Optional: node + `playwright` (for the Chromium dashboard view)

Install hints:
- macOS (Homebrew): `brew install ffmpeg xorg-server xterm`
- Debian/Ubuntu: `sudo apt-get install -y ffmpeg xvfb x11-apps xterm`

## Start the virtual display
```
VIRTUAL_DISPLAY=:99 \
ALERT_LOG=logs/anomaly_alerts.log \
./scripts/start_virtual_display.sh
```

This starts Xvfb with a tailing xterm window bound to `DISPLAY=:99`. State/pids live in `logs/.virtual_display` by default. Stop it with:
```
VIRTUAL_DISPLAY_STATE_DIR=logs/.virtual_display ./scripts/stop_virtual_display.sh
```

## Run the recorder against the virtual screen
```
CAPTURE_BACKEND=x11 \
SCREEN_DEVICE=:99 \
DISPLAY=:99 \
SESSION_DIR=logs/session_$(date +%s) \
./scripts/auto_screen_recorder.sh
```

Trigger an alert to generate a file:
```
echo "$(date) Anomaly ID: demo-virtual" >> logs/anomaly_alerts.log
```
Result: `logs/session_<ts>/recordings/anomaly_demo-virtual_<timestamp>.mp4`.

## Quick end-to-end test
```
./scripts/test_virtual_recorder.sh
```
This starts the virtual display, runs the recorder in `CAPTURE_BACKEND=x11`, appends a synthetic alert, and asserts that a non-empty MP4 exists under `logs/session_*_virtual_test/recordings`.

## Optional Playwright dashboard
To render a Chromium dashboard instead of the xterm feed, install Playwright (`npm install --save-dev playwright`) and run:
```
USE_PLAYWRIGHT_DASHBOARD=1 DISPLAY=:99 ./scripts/start_virtual_display.sh
```
The dashboard reads the alert log, refreshes every few seconds, and remains confined to the virtual display.

