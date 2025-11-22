# Driftlock Live Crypto Soak Test

This guide explains how to run the overnight soak test using live Cryptocurrency trade data from Binance.

## Overview

The soak test pipeline consists of:
1.  **Bridge**: `scripts/crypto_bridge.py` connects to the Binance WebSocket API (free, public) and streams high-volatility trades (DOGE, PEPE, XRP, SOL, BTC).
2.  **Scanner**: `bin/driftlock` scans the stream for anomalies using the entropy window algorithm.
3.  **Runner**: `scripts/soak_runner.py` orchestrates the process, logging all output and sending detected anomalies to the Firebase AI analysis function.

## Prerequisites

- Python 3.9+ with `websockets` installed (`pip install websockets`)
- Built `driftlock` binary in `bin/driftlock` (Run `make build` or `go build -o bin/driftlock cmd/driftlock-cli/main.go`)

## Running the Test

It is recommended to run this in a `tmux` or `screen` session to ensure it persists overnight.

### 1. Start the Runner

```bash
# From the repository root
./scripts/soak_runner.py
```

You should see output indicating the bridge is connecting and the worker is started.

### 2. Verify Operation

Check the `logs/` directory:

```bash
# Watch the stream log growing
tail -f logs/live-crypto.ndjson

# Watch for AI analysis responses (only appears when anomalies are found and sent)
tail -f logs/live-gemini.ndjson
```

## Verification (Next Morning)

To verify the test was successful:

1.  **Check Log Volume**: Ensure `logs/live-crypto.ndjson` has grown significantly.
2.  **Count Anomalies**:
    ```bash
    grep -c '"anomaly":true' logs/live-crypto.ndjson
    ```
3.  **Review AI Summaries**:
    Check `logs/live-gemini.ndjson` for successful responses from the Firebase function.

## Troubleshooting

- **No Data**: Check internet connection. The Binance API `wss://stream.binance.com:9443` must be accessible.
- **Binary Not Found**: Ensure `bin/driftlock` exists.
- **Firebase Errors**: Check `live-gemini.ndjson` for error messages.
