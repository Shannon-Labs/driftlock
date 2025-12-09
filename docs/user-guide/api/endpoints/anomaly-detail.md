# GET /v1/anomalies/{id}

Retrieve detailed information about a specific anomaly, including metrics, evidence, and the original event.

## Endpoint

```
GET https://api.driftlock.net/v1/anomalies/{id}
```

## Authentication

`X-Api-Key` header with `stream` or `admin` role.

## Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | The anomaly ID (e.g., `anom_xyz789`) |

## Response

```json
{
  "id": "anom_xyz789",
  "stream_id": "c4e6f7a8-...",
  "timestamp": "2025-01-01T10:42:00Z",
  "detected_at": "2025-01-01T10:42:05Z",
  "metrics": {
    "ncd": 0.72,
    "p_value": 0.004,
    "confidence": 0.996,
    "compression_ratio": 1.41,
    "entropy_change": 0.13,
    "baseline_size": 400,
    "window_size": 50
  },
  "event": {
    "timestamp": "2025-01-01T10:42:00Z",
    "type": "metric",
    "body": {"latency": 950},
    "attributes": {"host": "server-01"}
  },
  "explanation": {
    "summary": "Significant latency spike detected",
    "details": "Event compressed poorly relative to baseline (NCD 0.72). Entropy increased by 13%."
  },
  "evidence": {
    "baseline_sample": [...],
    "window_sample": [...]
  },
  "status": "new"
}
```

## Examples

```bash
curl "https://api.driftlock.net/v1/anomalies/anom_xyz789" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Error responses

- **404 not_found** — anomaly does not exist or belongs to another stream
- **403 forbidden** — key lacks permission for this anomaly
