# Quickstart Guide (5 minutes)

Run your first Driftlock detection in minutes. This guide covers signup, API key setup, and first requests with cURL and Python.

## 1) Get an API key

- Sign up at [https://driftlock.net/#signup](https://driftlock.net/#signup) (Pilot plan: 10k events/month free).
- Open **Dashboard → API Keys → Create key** and copy it. Treat it like a secret.
- No signup? Use the public demo endpoint: [POST /v1/demo/detect](../api/endpoints/demo.md) (rate-limited, no persistence).

## 2) Set your environment

```bash
export DRIFTLOCK_API_URL="https://api.driftlock.net/v1"
export DRIFTLOCK_API_KEY="YOUR_API_KEY"
```

## 3) Send your first detection (cURL)

```bash
curl -X POST "$DRIFTLOCK_API_URL/detect" \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"latency_ms": 120}},
      {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"latency_ms": 125}},
      {"timestamp": "2025-01-01T10:02:00Z", "type": "metric", "body": {"latency_ms": 118}},
      {"timestamp": "2025-01-01T10:03:00Z", "type": "metric", "body": {"latency_ms": 950}}  /* Anomaly */
    ]
  }'
```

**Expected response (truncated):**

```json
{
  "success": true,
  "anomaly_count": 1,
  "anomalies": [
    {
      "index": 3,
      "metrics": {"ncd": 0.72, "p_value": 0.004, "confidence": 0.996},
      "why": "Significant latency spike detected"
    }
  ],
  "request_id": "req_abc123"
}
```

## 4) Python example

```python
import os
import requests

api_url = os.getenv("DRIFTLOCK_API_URL", "https://api.driftlock.net/v1")
api_key = os.environ["DRIFTLOCK_API_KEY"]

payload = {
    "stream_id": "default",
    "events": [
        {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"latency_ms": 120}},
        {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"latency_ms": 125}},
        {"timestamp": "2025-01-01T10:02:00Z", "type": "metric", "body": {"latency_ms": 118}},
        {"timestamp": "2025-01-01T10:03:00Z", "type": "metric", "body": {"latency_ms": 950}},
    ],
}

resp = requests.post(
    f"{api_url}/detect",
    headers={"X-Api-Key": api_key, "Content-Type": "application/json"},
    json=payload,
    timeout=15,
)
resp.raise_for_status()
data = resp.json()

print(f"Found {data['anomaly_count']} anomalies")
for anom in data.get("anomalies", []):
    print(f"Index {anom['index']}: {anom['why']}")
```

## 5) View anomalies and iterate

- Query history: `curl "$DRIFTLOCK_API_URL/anomalies" -H "X-Api-Key: $DRIFTLOCK_API_KEY"`
- Tweak sensitivity: use `config_override` on `/v1/detect` or set a profile (see [Detection Profiles](../guides/detection-profiles.md)).
- Send feedback: mark false positives/confirmed via `/v1/anomalies/{id}/feedback` when available.

## Troubleshooting

| Symptom | Fix |
| --- | --- |
| `unauthorized` | Check `X-Api-Key` header, ensure key is active, no extra spaces. |
| `rate_limit_exceeded` | Respect `retry_after_seconds`, batch events (up to 256), or upgrade plan. |
| Slow response | Keep batches ≤256 events and payloads <10 MB; prefer `zstd` (default). |

## Next steps

- **Concepts:** [Baselines, windows, NCD, p-values](./concepts.md)
- **API reference:** [REST API overview](../api/rest-api.md)
- **Examples:** [cURL](../api/examples/curl-examples.md), [Python](../api/examples/python-examples.md)
- **Demo without signup:** [POST /v1/demo/detect](../api/endpoints/demo.md)
