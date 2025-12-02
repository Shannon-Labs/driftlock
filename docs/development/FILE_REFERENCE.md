# Automated Crypto Anomaly Detection - File Reference

## Core Scripts (in `scripts/`)

Created for automated 4-hour anomaly detection sessions:

| Script | Size | Purpose |
|--------|------|---------|
| `auto_crypto_monitor.sh` | 8.4KB | **Main orchestrator** - runs entire 4-hour session |
| `alert_on_anomaly.sh` | 3.8KB | Real-time anomaly alerting (visual + audio + notifications) |
| `status_reporter.sh` | 4.6KB | 15-minute status updates (trade counts, health checks) |
| `synthetic_anomaly.sh` | 4.5KB | Synthetic anomaly injection (if no real anomalies for 30 min) |
| `loom_controller.sh` | 2.8KB | Loom recording state tracking |
| `auto_screen_recorder.sh` | 2.7KB | **Auto screen recording** on anomaly detection (ffmpeg) |
| `stream_kraken_ws.py` | 2.9KB | Enhanced Kraken WebSocket streamer (connection logging) |

## Documentation

| File | Purpose |
|------|---------|
| `CRYPTO_ANOMALY_AUTOMATION.md` | **Main guide** - Quick start, usage, configuration |
| `README.md` | Driftlock project README |

## Quick Start

```bash
# Full 4-hour session with auto screen recording
./scripts/auto_crypto_monitor.sh

# 2-minute test
./scripts/auto_crypto_monitor.sh --test-mode --duration 120
```

## What Was Cleaned Up

Removed old test files:
- ❌ `README_CRYPTO_TEST.md` (superseded by CRYPTO_ANOMALY_AUTOMATION.md)
- ❌ `README_CRYPTO_TEST_LOOM.md` (superseded by CRYPTO_ANOMALY_AUTOMATION.md)
- ❌ `QUICK_START_LOOM.md` (superseded by CRYPTO_ANOMALY_AUTOMATION.md)
- ❌ `AUTO_RECORDING_SOLUTION.md` (merged into CRYPTO_ANOMALY_AUTOMATION.md)
- ❌ All test session directories (`logs/session_*`)
- ❌ All old test logs (`logs/*-2025*.log`)
- ❌ All old PID files

## What Remains

**Active logs:**
- `logs/anomaly_alerts.log` - Live anomaly alert history (shared across sessions)

**Next session will create:**
- `logs/session_<timestamp>/` - New session directory
  - `kraken_raw.ndjson` - Raw trades
  - `kraken_anomalies.ndjson` - Driftlock output
  - `recordings/` - Auto-recorded screen captures (if ffmpeg)
  - Various logs and state files

Clean and ready! 🚀
