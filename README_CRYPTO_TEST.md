# 🚀 Start 4-Hour Crypto Test - Quick Guide

## Step 1: Get Your API Key

Visit **https://driftlock.web.app** and sign up. You'll get an API key immediately.

## Step 2: Start the Test

```bash
# Set your API key
export DRIFTLOCK_API_KEY="dlk_..."

# Start the 4-hour test
./scripts/start_crypto_test.sh
```

That's it! The test will:
- ✅ Stream live Binance crypto data for 4 hours
- ✅ Send batches to your API every 5 seconds
- ✅ Detect anomalies in real-time
- ✅ Log everything to `logs/crypto-api-test-*.log`

## Monitor Progress

In another terminal:
```bash
tail -f logs/crypto-api-test-*.log
```

## Stop Early

Press `Ctrl+C` or:
```bash
kill $(cat logs/crypto-api-test-*.pid)
```

## Run in Background (Overnight)

```bash
# Using nohup
nohup ./scripts/start_crypto_test.sh > /tmp/crypto-test.out 2>&1 &

# Or using tmux (recommended)
tmux new-session -d -s crypto-test './scripts/start_crypto_test.sh'
```

## What You'll See

- Real-time anomaly detection on live crypto trades
- Batch processing (10 events per batch)
- Anomaly explanations and scores
- Summary statistics at the end

Perfect for validating the system works! 🎯

