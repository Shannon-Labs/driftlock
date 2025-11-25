# Automated Crypto Anomaly Detection - Quick Start Guide

## 🚀 Overview

This system provides a fully automated 4-hour crypto anomaly detection session with:
- **Live Kraken WebSocket streaming** (BTC/USD trades)
- **Driftlock entropy detector** (baseline 120 lines, zstd compression)
- **Real-time anomaly alerts** with "🚨 START LOOM NOW!" notifications
- **15-minute status updates** (trade counts, anomaly stats, stream health)
- **Auto-reconnection** if WebSocket drops
- **Synthetic anomaly injection** if no real anomalies for 30 minutes
- **Loom recording tracking** to correlate anomalies with recordings

## 📋 Prerequisites

✅ Python 3 with `websocket` module installed  
✅ `driftlock` CLI available at `./bin/driftlock`  
✅ `jq` command-line JSON processor  
✅ **`ffmpeg`** (for automatic screen recording) - **RECOMMENDED**  
✅ Loom (optional, manual fallback if no ffmpeg)

Check prerequisites:
```bash
python3 -c "import websocket; print('websocket OK')"
./bin/driftlock scan --help
jq --version
ffmpeg -version  # Install with: brew install ffmpeg
```

### Installing ffmpeg (for Auto Recording)

**macOS:**
```bash
brew install ffmpeg
```

**First Run Permission:**
On macOS, the first time ffmpeg tries to record the screen, you'll need to grant screen recording permission:
1. System Settings → Privacy & Security → Screen Recording
2. Enable Terminal (or your terminal app)

## 🎬 Quick Start

### 1. Start a 4-Hour Session

```bash
cd /Volumes/VIXinSSD/driftlock
./scripts/auto_crypto_monitor.sh
```

This will:
- Stream live BTC/USD trades from Kraken
- Detect anomalies using Driftlock entropy analysis
- Alert you immediately when anomalies are detected
- Report status every 15 minutes
- Run for 4 hours total

### 2. Start a Test Session (2 minutes)

```bash
./scripts/auto_crypto_monitor.sh --test-mode --duration 120
```

### 3. Monitor Different Pair

```bash
export KRAKEN_PAIR="ETH/USD"
./scripts/auto_crypto_monitor.sh
```

## 📊 What You'll See

### Startup
```
═══════════════════════════════════════════════════════════
  🚀 AUTOMATED CRYPTO ANOMALY DETECTION SESSION
═══════════════════════════════════════════════════════════

   Pair: BTC/USD
   Duration: 4h 0m
   Baseline: 120 lines
   Compression: zstd
   Session Dir: logs/session_1763870312
   Started: 2025-11-22 21:58:32

   📹 Loom will be prompted on anomaly detection

═══════════════════════════════════════════════════════════

[1/6] Starting Kraken WebSocket streamer...
   PID: 3759
   ✅ Streamer running

[2/6] Starting Driftlock entropy detector...
   PID: 3893
   ✅ Detector running
   Alerter PID: 3894
   ✅ Alerter running

[3/6] Starting status reporter (15-minute intervals)...
   PID: 4050
   ✅ Status reporter running

[4/6] Starting synthetic anomaly injector (30-min threshold)...
   PID: 4051
   ✅ Synthetic injector running

[5/6] Starting Loom controller...
   PID: 4059
   ✅ Loom controller running

[6/6] Session active - monitoring for 14400s...
   Press Ctrl+C to stop early
```

### Anomaly Alert
```
═══════════════════════════════════════════════════════════
🚨 START LOOM NOW! 🚨
═══════════════════════════════════════════════════════════
Alert #1
Time: 2025-11-22 22:15:43
Trade Time: 2025-11-22 22:15:42
Anomaly ID: kraken-1763871342123-5
Pair: BTC/USD
Price: 87456.2
Score: 0.87

Full Payload:
{"sequence":247,"line":"{...}","entropy":5.42,"score":0.87,"is_anomaly":true,...}
═══════════════════════════════════════════════════════════

🎬 RECORDING STARTED
   Time: 2025-11-22 22:15:43
   Triggered by Anomaly: kraken-1763871342123-5
```

You'll also hear:
- 🔔 Terminal bell
- 🗣️ Voice alert: "Anomaly detected! Start Loom now! Alert number 1"
- 🔊 Sound effect (Glass.aiff)
- 💬 macOS notification

### Status Update (Every 15 Minutes)
```
⏰ STATUS UPDATE - 22:13:32
   Runtime: 0h 15m
   Total Trades: 1,247
   Total Anomalies: 2
   Last Anomaly: 2025-11-22 22:08:15 (ID: kraken-1763870895432-89)
   Raw Log: 156KB (logs/session_1763870312/kraken_raw.ndjson)
   Anomaly Log: 2KB (logs/session_1763870312/kraken_anomalies.ndjson)
   Stream Status: ✅ HEALTHY (last update 2s ago)
```

### Synthetic Anomaly Injection (If No Real Anomalies for 30 Min)
```
🧪 INJECTING SYNTHETIC ANOMALY
   Reason: No anomalies detected for 32 minutes
   Method: Switch BTC/USD → ETH/USD → BTC/USD
   Time: 2025-11-22 22:45:12
   
   Stopping BTC/USD streamer (PID: 3759)...
   Starting ETH/USD streamer for 120 seconds...
   [... 2 minutes pass ...]
   Stopping ETH/USD streamer (PID: 5123)...
   Restarting BTC/USD streamer...
   
✅ Synthetic injection complete - reverted to BTC/USD (PID: 5234)
```

### Session End
```
⏱️  Session duration reached - shutting down...

🛑 Shutting down session...
   Stopping process 3759...
   Stopping process 3893...
   Stopping process 4050...
   Stopping process 4051...
   Stopping process 4059...

📊 SESSION SUMMARY
   Duration: 14400s
   Session Dir: logs/session_1763870312
   Total Trades: 45,623
   Total Anomalies: 8

✅ Session complete
```

## 📁 Log Files

All logs are in `logs/session_<timestamp>/`:

```
logs/session_1763870312/
├── kraken_raw.ndjson          # Raw Kraken trades (NDJSON)
├── kraken_anomalies.ndjson    # Driftlock output with anomaly flags
├── master.log                 # Master orchestrator log
├── status.log                 # Status updates
├── synthetic.log              # Synthetic injection events
├── loom.log                   # Loom recording events
├── streamer.stderr.log        # Streamer connection events
├── detector.stderr.log        # Driftlock detector errors
└── loom.state                 # Current Loom state (idle/recording)
```

Plus: `logs/anomaly_alerts.log` - All anomaly alerts (shared across sessions)

## 🧪 Testing Individual Components

### Test Anomaly Alerter
```bash
echo '{"ts":1763870503,"price":100000,"is_anomaly":true,"score":0.95}' | \
  ./scripts/alert_on_anomaly.sh
```

### Test Status Reporter (Fast Mode)
```bash
export SESSION_DIR=logs/session_1763870312
./scripts/status_reporter.sh --test-interval 15
```

### Force Synthetic Injection
```bash
export SESSION_DIR=logs/session_1763870312
./scripts/synthetic_anomaly.sh --force-inject
```

### Test Loom Controller
```bash
export SESSION_DIR=logs/session_1763870312
./scripts/loom_controller.sh
# In another terminal:
echo "Anomaly ID: test-123" >> logs/anomaly_alerts.log
```

## 🎥 Automatic Screen Recording (Set & Forget!)

### How It Works

If `ffmpeg` is installed, the system **automatically records your screen** when anomalies are detected:

1. System monitors for anomaly alerts
2. When anomaly detected → **automatically starts 60-second screen recording**
3. Saves to `logs/session_*/recordings/anomaly_<id>_<timestamp>.mp4`
4. No manual clicking required!

### What Gets Recorded

- **Duration:** 60 seconds per anomaly (configurable)
- **Resolution:** 1920x1080 (configurable)
- **Format:** MP4 (H.264)
- **Frame Rate:** 30 FPS
- **Location:** `logs/session_<timestamp>/recordings/`

### Example Output

```bash
logs/session_1763870312/recordings/
├── anomaly_kraken-1763871342123-5_20251122_221543.mp4  (15MB)
├── anomaly_kraken-1763872891456-12_20251122_224131.mp4 (14MB)
└── anomaly_test-anomaly-1_20251122_225917.mp4           (16MB)
```

### Configuration

```bash
# Change recording duration (default: 60s)
export RECORDING_DURATION=90

# Change resolution (default: 1920x1080)
export SCREEN_RESOLUTION=2560x1440

# Then start session
./scripts/auto_crypto_monitor.sh
```

### What You'll See

```
[5/7] Starting auto screen recorder...
   PID: 4123
   ✅ Auto recorder running (60s per anomaly)
```

When anomaly detected:
```
🎥 RECORDING STARTED
   Anomaly ID: kraken-1763871342123-5
   Output: logs/session_1763870312/recordings/anomaly_kraken-1763871342123-5_20251122_221543.mp4
   Duration: 60s
   PID: 4567

[... 60 seconds later ...]
✅ Recording complete: anomaly_kraken-1763871342123-5_20251122_221543.mp4 (15MB)
```

### MacOS Permissions

**First time only:** Grant screen recording permission to your terminal:

1. Run the session once
2. macOS will prompt forpermission
3. Go to: **System Settings → Privacy & Security → Screen Recording**
4. Enable your terminal app (Terminal, iTerm2, etc.)
5. Restart the session

### Manual Loom Fallback

If you don't have `ffmpeg` or prefer Loom:

1. System will show: `📹 Loom will be prompted on anomaly detection`
2. Keep Loom ready
3. When alert fires ("🚨 START LOOM NOW!"), manually click Loom to record
4. Optionally mark state:
   ```bash
   echo 'recording' > logs/session_*/loom.state
   ```

### Standalone Screen Recorder

You can also run the recorder independently:

```bash
# Set session directory
export SESSION_DIR=logs/session_1763870312

# Start recorder
./scripts/auto_screen_recorder.sh
```

It will watch `logs/anomaly_alerts.log` and auto-record when new anomalies appear.

## 🔧 Configuration

### Environment Variables

```bash
# Change trading pair
export KRAKEN_PAIR="ETH/USD"

# Adjust baseline window (default: 120)
export BASELINE_LINES=200

# Change compression algorithm (default: zstd)
export COMPRESSION_ALGO="gzip"

# Adjust synthetic injection threshold (default: 1800s = 30 min)
export NO_ANOMALY_THRESHOLD=3600  # 1 hour
```

### Command-Line Flags

```bash
# Short test run (2 minutes)
./scripts/auto_crypto_monitor.sh --test-mode --duration 120

# Longer session (8 hours)
./scripts/auto_crypto_monitor.sh --duration 28800

# Different pair
./scripts/auto_crypto_monitor.sh --pair ETH/USD
```

## 🛑 Stopping Early

Press `Ctrl+C` to gracefully shutdown. The script will:
- Stop all processes
- Clean up named pipes
- Show session summary
- Preserve all logs

## 📝 Notes

- **Anomalies are rare** in normal market conditions. The 4-hour duration should catch some, but synthetic injection ensures testing works.
- **Baseline warmup**: Driftlock needs 120 trades before it can detect anomalies (`ready:false` until then).
- **Status updates**: Default 15-minute interval. Use `--test-interval` to speed up for testing.
- **Loom has no API**: The controller tracks state but cannot programmatically control Loom. You must start/stop manually.

## 🎯 What Makes This Novel

This is a **demo automation agent** that:
1. ✅ Keeps a long-running session active (4 hours)
2. ✅ Monitors live streaming data (Kraken WebSocket)
3. ✅ Detects anomalies deterministically (Driftlock entropy)
4. ✅ Fires immediate alerts (visual + audio + notifications)
5. ✅ Provides periodic status updates (15 minutes)
6. ✅ Auto-recovers from failures (WebSocket reconnect)
7. ✅ Ensures testability (synthetic anomaly injection)
8. ✅ Integrates with human workflow (Loom prompts)

No expensive AI calls. All deterministic. All scriptable. All loggable.

Perfect for demos! 🎬
## ✅ Solution: Automatic Screen Recording Added!

### What Changed

Added **automatic screen recording** using `ffmpeg` so you can truly "set and forget":

**New Script:** [`auto_screen_recorder.sh`](file:///Volumes/VIXinSSD/driftlock/scripts/auto_screen_recorder.sh)
- Monitors anomaly alerts automatically
- When anomaly detected → **starts 60-second screen recording**  
- Saves to `logs/session_*/recordings/anomaly_<id>_<timestamp>.mp4`
- **No manual clicks needed!**

### How To Use

Since you already have `ffmpeg` installed (`/opt/homebrew/bin/ffmpeg`), just run:

```bash
cd /Volumes/VIXinSSD/driftlock
./scripts/auto_crypto_monitor.sh
```

**That's it!** The system will:
1. Stream live BTC/USD trades from Kraken ✅
2. Detect anomalies with Driftlock entropy ✅
3. **Automatically record 60s video when anomaly found** ✅
4. Save videos to `logs/session_*/recordings/` ✅

### First-Time Setup (macOS Permission)

The first time ffmpeg tries to record, macOS will ask for permission:

1. System Settings → Privacy & Security → Screen Recording
2. Enable your terminal app (Terminal, iTerm2, etc.)
3. Restart the session

### What You Asked About

> "is there any way to like auto trigger loom to start or something"

**Answer:** Loom has no API, BUT we can use **ffmpeg to auto-record the screen** instead! This is actually better because:
- ✅ Fully automatic (no clicking)
- ✅ Works headlessly
- ✅ Saves directly to MP4 files with anomaly IDs
- ✅ More reliable than trying to control Loom
- ✅ You already have ffmpeg installed!

> "i'm trying to figure out how to set this and forget it with the real data you know?"

**Answer:** Now you can! Just run `./scripts/auto_crypto_monitor.sh` and walk away. When anomalies are detected:
- 🎥 Screen automatically records for 60 seconds
- 📊 Status updates every 15 minutes
- 🔄 Auto-reconnects if stream drops
- 🧪 Synthetic anomaly injection if nothing happens for 30 min

Truly set-and-forget!

### Test It Quick

```bash
# 2-minute test session
./scripts/auto_crypto_monitor.sh --test-mode --duration 120

# Then inject a fake anomaly to test recording
echo '{"ts":1763870503,"price":100000,"is_anomaly":true,"id":"test-1"}' >> logs/anomaly_alerts.log

# Check for video in logs/session_*/recordings/
```

### Full Documentation

- **Quick Start:** [`CRYPTO_ANOMALY_AUTOMATION.md`](file:///Volumes/VIXinSSD/driftlock/CRYPTO_ANOMALY_AUTOMATION.md)
- **Walkthrough:** [walkthrough.md](file:///Users/hunterbown/.gemini/antigravity/brain/f897386b-d847-4a52-9da5-ef7ad7dae175/walkthrough.md)

Ready to run! 🚀

---

_NOTE: This replaces the older QUICK_START_LOOM.md, README_CRYPTO_TEST.md, and README_CRYPTO_TEST_LOOM.md files._
