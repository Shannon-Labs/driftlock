# GET /v1/anomalies

List detected anomalies with filtering and pagination.

## Endpoint

```
GET https://api.driftlock.net/v1/anomalies
```

## Authentication

`X-Api-Key` header with `stream` or `admin` role.

## Query parameters

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| `limit` | integer | 50 (max 100) | Page size |
| `page_token` | string | — | Cursor from previous page |
| `stream_id` | string | — | Filter by stream |
| `min_ncd` | float | — | Minimum NCD |
| `max_p_value` | float | — | Maximum p-value |
| `status` | string | — | `new`, `viewed`, `archived` |
| `since` | ISO-8601 | — | Start timestamp filter |
| `until` | ISO-8601 | — | End timestamp filter |
| `has_evidence` | boolean | — | Only anomalies with saved evidence |

## Response

```json
{
  "anomalies": [
    {
      "id": "anom_xyz789",
      "stream_id": "c4e6f7a8-...",
      "timestamp": "2025-01-01T10:42:00Z",
      "metrics": {"ncd": 0.72, "p_value": 0.004, "confidence": 0.996},
      "event_summary": {"type": "metric", "preview": "{\"latency\":950}"},
      "why": "Significant latency spike detected",
      "status": "new"
    }
  ],
  "next_page_token": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMy0wMVQxMjo1ODoxNFoifQ==",
  "total": 1532
}
```

### Response fields

| Field | Description |
| --- | --- |
| `anomalies` | Array of anomaly summaries |
| `next_page_token` | Cursor for the next page (null if none) |
| `total` | Total anomalies matching filters |

## Examples

### Basic list
```bash
curl "https://api.driftlock.net/v1/anomalies" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Filter by stream and score
```bash
curl "https://api.driftlock.net/v1/anomalies?stream_id=prod-logs&min_ncd=0.5" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

### Date range + pagination
```bash
curl "https://api.driftlock.net/v1/anomalies?since=2025-01-01T00:00:00Z&until=2025-01-02T00:00:00Z&limit=20" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"

# Next page
curl "https://api.driftlock.net/v1/anomalies?page_token=eyJ..." \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY"
```

## Error responses

- **400 invalid_argument** — bad date format or invalid limit
- **401 unauthorized** — missing/invalid API key

## Feedback (optional)

If enabled, record feedback to improve auto-tuning:

```bash
curl -X POST "https://api.driftlock.net/v1/anomalies/{id}/feedback" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type":"confirmed","reason":"Clear production incident"}'
```

## Related

- [POST /v1/detect](./detect.md)
- [Error codes](../errors.md)
