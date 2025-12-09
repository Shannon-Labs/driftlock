# Tutorial: Metrics Spike Detection

Catch sudden latency/CPU spikes in metrics data.

## Prerequisites
- API key
- Stream ID (e.g., `metrics`)

## Sample payload

```json
{
  "stream_id": "metrics",
  "events": [
    {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"cpu": 42, "latency_ms": 120}},
    {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"cpu": 44, "latency_ms": 118}},
    {"timestamp": "2025-01-01T10:02:00Z", "type": "metric", "body": {"cpu": 45, "latency_ms": 124}},
    {"timestamp": "2025-01-01T10:03:00Z", "type": "metric", "body": {"cpu": 92, "latency_ms": 950}}
  ]
}
```

## Run detection

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d @payload.json
```

Expected: the 4th event (CPU/latency spike) is flagged with high confidence.

## Tips
- Use `stream_id` per service to keep baselines tight.
- Switch to `strict` if normal traffic is noisy; use `sensitive` for SLO alerting.
- Add `region`/`host` attributes to speed investigations.
