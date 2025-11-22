# Live Crypto Test - 4 Hour Run Guide

## Overview

This guide explains how to run the 4-hour live cryptocurrency anomaly detection test using real Binance trade data.

## Prerequisites

1. **API Key**: Sign up at https://driftlock.web.app to get your API key
2. **Python 3.9+** with dependencies:
   ```bash
   pip3 install websockets requests certifi
   ```

## Quick Start

### 1. Get Your API Key

Visit https://driftlock.web.app and sign up. You'll receive an API key immediately.

### 2. Set Environment Variables

```bash
export DRIFTLOCK_API_KEY="dlk_..."
export DRIFTLOCK_API_URL="https://driftlock.web.app/api/v1"
```

### 3. Start the 4-Hour Test

```bash
./scripts/start_crypto_test.sh
```

Or run directly:

```bash
./scripts/run_crypto_test_4h.sh
```

## What It Does

1. **Connects to Binance WebSocket** - Streams live crypto trade data for:
   - BTC/USDT
   - ETH/USDT
   - SOL/USDT
   - LINK/USDT
   - AVAX/USDT
   - DOGE/USDT
   - LTC/USDT

2. **Sends Batches to Driftlock API** - Every 5 seconds or 10 events (whichever comes first)

3. **Detects Anomalies** - Real-time compression-based anomaly detection

4. **Logs Everything** - All activity saved to `logs/crypto-api-test-*.log`

## Monitoring

### Watch Live Progress

```bash
# In another terminal
tail -f logs/crypto-api-test-*.log
```

### Check for Anomalies

```bash
grep "anomalies detected" logs/crypto-api-test-*.log
```

### View Summary

The script will print a summary when it completes:
- Total events processed
- Total anomalies detected
- Anomaly rate percentage

## Running in Background

To run in the background (e.g., overnight):

```bash
# Using nohup
nohup ./scripts/start_crypto_test.sh > /tmp/crypto-test.out 2>&1 &

# Or using tmux (recommended)
tmux new-session -d -s crypto-test './scripts/start_crypto_test.sh'
tmux attach -t crypto-test  # To view
```

## Stopping Early

Press `Ctrl+C` or:

```bash
# Find the PID
cat logs/crypto-api-test-*.pid

# Kill the process
kill $(cat logs/crypto-api-test-*.pid)
```

## Expected Results

- **Event Rate**: ~100-500 events per minute (depends on market activity)
- **Anomaly Rate**: Typically 1-5% of events (volatility spikes, unusual trade patterns)
- **Batch Size**: 10 events per batch
- **API Calls**: ~6-30 calls per minute

## Troubleshooting

### Connection Issues

If Binance WebSocket fails:
- Check internet connection
- Try global Binance: `export BINANCE_WS_URL="wss://stream.binance.com:9443/ws"`

### API Errors

If you see 401 errors:
- Verify your API key is correct
- Check that the key hasn't been revoked

If you see 429 errors:
- You're hitting rate limits
- The script will retry automatically

### Python Dependencies

```bash
pip3 install --upgrade websockets requests certifi
```

## Use Cases

This test demonstrates:
- ✅ Real-time anomaly detection on live data
- ✅ High-throughput event processing
- ✅ Compression-based detection accuracy
- ✅ API reliability and performance
- ✅ Perfect for demos and validation

## Next Steps

After the test completes:
1. Review the log file for detected anomalies
2. Check the anomaly explanations
3. Use the results to validate detection accuracy
4. Share the results as a use case example

