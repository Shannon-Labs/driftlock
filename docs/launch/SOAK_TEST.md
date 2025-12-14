# Driftlock Live Crypto Soak Test

This guide explains how to run the overnight soak test using live Cryptocurrency trade data from Binance.

## Overview

The soak test pipeline consists of:
1. **Bridge**: `scripts/crypto_bridge.py` connects to the Binance WebSocket API (free, public) and streams high-volatility trades (DOGE, PEPE, XRP, SOL, BTC).
2. **API**: Driftlock Rust API server processes the stream for anomalies via `/v1/demo/detect` endpoint.
3. **Runner**: `scripts/soak_runner.py` orchestrates the process, logging all output and sending detected anomalies to analysis.

## Prerequisites

- Python 3.9+ with `websockets` and `requests` installed:
  ```bash
  pip install websockets requests
  ```
- Running Driftlock API server:
  ```bash
  cargo build -p driftlock-api --release
  DATABASE_URL="postgres://..." ./target/release/driftlock-api
  ```

## Running the Test

It is recommended to run this in a `tmux` or `screen` session to ensure it persists overnight.

### 1. Start the API Server

```bash
# Build and run the Rust API
cargo build -p driftlock-api --release

# Start the server (ensure DATABASE_URL is set)
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  ./target/release/driftlock-api
```

### 2. Start the Runner

```bash
# From the repository root
./scripts/soak_runner.py
```

You should see output indicating the bridge is connecting and data is being processed.

### 3. Verify Operation

Check the `logs/` directory:

```bash
# Watch the stream log growing
tail -f logs/live-crypto.ndjson

# Watch for analysis responses (only appears when anomalies are found)
tail -f logs/live-gemini.ndjson
```

### Alternative: Direct API Testing

You can also test detection directly via curl:

```bash
# Send sample crypto events to the demo endpoint
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "DOGE/USDT: 0.08234 @ 1000000 volume",
      "DOGE/USDT: 0.08235 @ 1500000 volume",
      "DOGE/USDT: 0.15000 @ 50000000 volume"
    ]
  }'
```

## Verification (Next Morning)

To verify the test was successful:

1. **Check Log Volume**: Ensure `logs/live-crypto.ndjson` has grown significantly.
2. **Count Anomalies**:
   ```bash
   grep -c '"anomaly":true' logs/live-crypto.ndjson
   ```
3. **Review API Metrics**:
   ```bash
   curl http://localhost:8080/metrics | grep driftlock_anomalies
   ```
4. **Check Health**:
   ```bash
   curl http://localhost:8080/healthz
   ```

## Troubleshooting

- **No Data**: Check internet connection. The Binance API `wss://stream.binance.com:9443` must be accessible.
- **API Not Running**: Ensure the Rust API server is running on port 8080.
- **Rate Limited**: The demo endpoint has a rate limit of 10 requests/hour per IP. For extended testing, use authenticated endpoints with an API key.
- **Database Connection**: Ensure PostgreSQL is running and DATABASE_URL is correct.

## Driftlog (Debug Logging)

Enable detailed logging to troubleshoot issues:

```bash
# Run API with debug logging
RUST_LOG=debug cargo run -p driftlock-api

# Filter to detection module
RUST_LOG=driftlock_api::routes::detection=debug cargo run -p driftlock-api
```

## Production Soak Test

For production soak testing with authenticated endpoints:

```bash
# Get your API key from the dashboard
export DRIFTLOCK_API_KEY="dlk_..."

# Create a dedicated stream for crypto testing
curl -X POST http://localhost:8080/v1/streams \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "crypto-soak-test",
    "stream_type": "metrics",
    "detection_profile": "sensitive"
  }'

# Use authenticated detection
curl -X POST http://localhost:8080/v1/detect \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "crypto-soak-test",
    "events": ["..."]
  }'
```
