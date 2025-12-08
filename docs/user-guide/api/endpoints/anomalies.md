# GET /v1/anomalies

List detected anomalies with support for filtering, pagination, and sorting.

## Endpoint

```
GET https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies
```

## Authentication

Requires API key with `stream` or `admin` role:
```
X-Api-Key: YOUR_API_KEY
```

## Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `limit` | integer | No | Number of items to return (default: 50, max: 100) |
| `page_token` | string | No | Token for pagination (from previous response) |
| `stream_id` | string | No | Filter by stream ID |
| `min_ncd` | float | No | Filter by minimum NCD score (0.0-1.0) |
| `max_p_value` | float | No | Filter by maximum p-value (0.0-1.0) |
| `status` | string | No | Filter by status (`new`, `viewed`, `archived`) |
| `since` | string | No | Filter by start date (ISO 8601) |
| `until` | string | No | Filter by end date (ISO 8601) |
| `has_evidence` | boolean | No | Filter for anomalies with saved evidence |

## Response

### Success (HTTP 200)

```json
{
  "anomalies": [
    {
      "id": "anom_xyz789",
      "stream_id": "c4e6f7a8-...",
      "timestamp": "2025-01-01T10:42:00Z",
      "metrics": {
        "ncd": 0.72,
        "p_value": 0.004,
        "confidence": 0.996
      },
      "event_summary": {
        "type": "metric",
        "preview": "{\"latency\": 950}"
      },
      "why": "Significant latency spike detected",
      "status": "new"
    }
  ],
  "next_page_token": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMy0wMVQxMjo1ODoxNFoifQ==",
  "total": 1532
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `anomalies` | array | List of anomaly summaries |
| `next_page_token` | string | Token to fetch the next page (null if no more pages) |
| `total` | integer | Total count of anomalies matching filters |

## Examples

### Basic List
```bash
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies" \
  -H "X-Api-Key: YOUR_API_KEY"
```

### Filter by Stream and Score
```bash
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies?stream_id=prod-logs&min_ncd=0.5" \
  -H "X-Api-Key: YOUR_API_KEY"
```

### Date Range Filtering
```bash
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies?since=2025-01-01T00:00:00Z&until=2025-01-02T00:00:00Z" \
  -H "X-Api-Key: YOUR_API_KEY"
```

### Pagination
```bash
# Get first page
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies?limit=10" \
  -H "X-Api-Key: YOUR_API_KEY"

# Get next page using token from response
curl "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies?page_token=eyJ..." \
  -H "X-Api-Key: YOUR_API_KEY"
```

## Error Responses

### 400 Bad Request
```json
{
  "error": {
    "code": "invalid_argument",
    "message": "Invalid date format for 'since' parameter"
  }
}
```

### 401 Unauthorized
```json
{
  "error": {
    "code": "unauthorized",
    "message": "Invalid API key"
  }
}
```

# POST /v1/anomalies/{id}/feedback

Record feedback for a specific anomaly to improve auto-tuning.

## Endpoint

```
POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies/{id}/feedback
```

## Authentication

Requires Firebase bearer token or API key associated with the anomaly’s tenant.

## Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `feedback_type` | string | Yes | One of `false_positive`, `confirmed`, `dismissed` |
| `reason` | string | No | Optional context on why you marked it |

Example:

```bash
curl -X POST "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/anomalies/123/feedback" \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"feedback_type":"confirmed","reason":"Clear production incident"}'
```

## Response

`200 OK` with `{ "success": true, "message": "Feedback recorded" }`
