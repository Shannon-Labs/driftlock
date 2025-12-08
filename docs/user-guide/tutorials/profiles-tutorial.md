# Tutorial: Choosing Detection Profiles

This tutorial walks through selecting and comparing detection profiles for your stream.

## Prerequisites

- A Driftlock API key
- A stream with some data (or use the examples below)

## Step 1: Check Current Profile

First, see your stream's current settings:

```bash
curl https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here"
```

Response:

```json
{
  "stream_id": "logs-production",
  "profile": "balanced",
  "auto_tune_enabled": false,
  "current_thresholds": {
    "ncd_threshold": 0.30,
    "pvalue_threshold": 0.05,
    "baseline_size": 400,
    "window_size": 50
  }
}
```

## Step 2: Run Detection with Default Profile

Send some test events:

```bash
curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "logs-production",
    "events": [
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 45}},
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 52}},
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 48}},
      {"body": {"level": "error", "message": "Connection timeout", "latency_ms": 30000}}
    ]
  }'
```

Note the anomaly count in the response.

## Step 3: Switch to Sensitive Profile

If you're missing anomalies, try the sensitive profile:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"profile": "sensitive"}'
```

Response:

```json
{
  "stream_id": "logs-production",
  "profile": "sensitive",
  "current_thresholds": {
    "ncd_threshold": 0.20,
    "pvalue_threshold": 0.10,
    "baseline_size": 200,
    "window_size": 30
  }
}
```

## Step 4: Re-run Detection

Send the same events again and compare results:

```bash
curl -X POST https://driftlock.net/api/v1/detect \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "logs-production",
    "events": [
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 45}},
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 52}},
      {"body": {"level": "info", "message": "Request processed", "latency_ms": 48}},
      {"body": {"level": "error", "message": "Connection timeout", "latency_ms": 30000}}
    ]
  }'
```

The sensitive profile should flag more anomalies.

## Step 5: Try Strict Profile (If Too Many Alerts)

If you're getting too many false positives:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{"profile": "strict"}'
```

Re-run detection to see fewer, higher-confidence anomalies.

## Step 6: Enable Adaptive Features

Once you've found a baseline, enable learning:

```bash
curl -X PATCH https://driftlock.net/api/v1/streams/logs-production/profile \
  -H "X-Api-Key: dlk_your_key_here" \
  -H "Content-Type: application/json" \
  -d '{
    "auto_tune_enabled": true,
    "adaptive_window_enabled": true
  }'
```

Now your stream will automatically refine itself based on feedback.

## Profile Selection Matrix

| Your Situation | Recommended Profile |
|---------------|---------------------|
| Starting out, unsure | balanced |
| Missing real anomalies | sensitive |
| Too many false alerts | strict |
| Have specific requirements | custom |

## What's Next?

- [Feedback Loop Tutorial](./feedback-loop.md) - Learn to refine detection with feedback
- [Detection Profiles Guide](../guides/detection-profiles.md) - Deep dive into profiles
