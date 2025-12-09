# Tutorial: Choosing Detection Profiles

Compare `sensitive`, `balanced`, and `strict` profiles to pick the right sensitivity.

## Prerequisites

- API key
- A stream with data (or use the sample events below)

## 1) Check current profile

```bash
curl "https://api.driftlock.net/v1/streams/logs-production/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## 2) Run detection with the default (balanced)

```bash
curl -X POST "https://api.driftlock.net/v1/detect" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
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

Note the anomaly count.

## 3) Try the sensitive profile

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/logs-production/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"profile": "sensitive"}'
```

Re-run detection; expect more anomalies.

## 4) Try the strict profile

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/logs-production/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"profile": "strict"}'
```

Re-run detection; expect fewer, higher-confidence anomalies.

## 5) Enable adaptive features

```bash
curl -X PATCH "https://api.driftlock.net/v1/streams/logs-production/profile" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"auto_tune_enabled": true, "adaptive_window_enabled": true}'
```

This shifts the stream to `custom` when thresholds are tuned from feedback.

## Profile selection matrix

| Situation | Recommended profile |
| --- | --- |
| Starting out | balanced |
| Missing anomalies | sensitive |
| Too many false alerts | strict |
| Need manual thresholds | custom (set overrides) |

**Next:** [Detection Profiles](../guides/detection-profiles.md), [Feedback Loop](./feedback-loop.md)
