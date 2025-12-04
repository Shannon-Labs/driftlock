# AI Agent Integration Guide

Driftlock provides compression-based anomaly detection that's easy for AI agents to use.

## Quick Reference

| What | Value |
|------|-------|
| Base URL | `https://driftlock.net/api` |
| OpenAPI Spec | `https://driftlock.net/docs/architecture/api/openapi.yaml` |
| Demo (no auth) | `POST /v1/demo/detect` |
| Full API | `POST /v1/detect` with `X-Api-Key` header |

## Try Without Signup

Test Driftlock instantly using the demo endpoint:

```bash
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {"level": "info", "message": "User logged in"},
      {"level": "info", "message": "Request processed"},
      {"level": "error", "message": "CRITICAL: Database connection failed"}
    ]
  }'
```

**Limits**: 10 requests/minute per IP, max 50 events, results not persisted.

## Full API Access

### 1. Get an API Key

```bash
# Sign up
curl -X POST https://driftlock.net/api/v1/onboard/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "company_name": "Your Company"}'

# Check email for verification link
# After verification, you'll receive your API key
```

### 2. Detect Anomalies

```bash
curl -X POST https://driftlock.net/api/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: dlk_your-key-here.secret" \
  -d '{
    "stream_id": "logs-production",
    "events": [
      {"timestamp": "2024-01-15T10:00:00Z", "latency_ms": 45},
      {"timestamp": "2024-01-15T10:00:01Z", "latency_ms": 52},
      {"timestamp": "2024-01-15T10:00:02Z", "latency_ms": 30000}
    ]
  }'
```

## Understanding Results

```json
{
  "success": true,
  "anomaly_count": 1,
  "anomalies": [
    {
      "index": 2,
      "detected": true,
      "why": "High NCD (0.72) indicates significant deviation from baseline",
      "metrics": {
        "ncd": 0.72,
        "p_value": 0.004,
        "confidence_level": 0.996,
        "is_statistically_significant": true
      }
    }
  ]
}
```

### Key Metrics

| Metric | Range | Meaning |
|--------|-------|---------|
| `ncd` | 0-1 | Normalized Compression Distance. Higher = more anomalous |
| `p_value` | 0-1 | Statistical significance. Lower = more significant |
| `confidence_level` | 0-1 | 1 - p_value. Higher = more confident |
| `is_statistically_significant` | bool | True if p_value < threshold |

### Decision Logic

An event is flagged as anomalous when:
- `ncd > ncd_threshold` (default 0.3)
- `p_value < p_value_threshold` (default 0.05)

## Configuration Override

Fine-tune detection per request:

```json
{
  "events": [...],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd"
  }
}
```

## Rate Limits

| Endpoint | Limit | Notes |
|----------|-------|-------|
| Demo | 10/min per IP | Max 50 events |
| Authenticated | Based on plan | See /v1/me/usage |

Rate limit headers:
- `X-RateLimit-Limit`: Max requests
- `X-RateLimit-Remaining`: Remaining requests
- `Retry-After`: Seconds to wait (on 429)

## Error Handling

All errors return consistent JSON:

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Too many requests",
    "request_id": "req_abc123",
    "retry_after_seconds": 60
  }
}
```

## Python Client Example

```python
import requests

class DriftlockClient:
    def __init__(self, api_key: str):
        self.base_url = "https://driftlock.net/api"
        self.session = requests.Session()
        self.session.headers["X-Api-Key"] = api_key
        self.session.headers["Content-Type"] = "application/json"

    def detect(self, events: list, stream_id: str = None, config: dict = None):
        payload = {"events": events}
        if stream_id:
            payload["stream_id"] = stream_id
        if config:
            payload["config_override"] = config

        resp = self.session.post(f"{self.base_url}/v1/detect", json=payload)

        if resp.status_code == 429:
            retry_after = resp.json()["error"].get("retry_after_seconds", 60)
            raise RateLimitError(f"Rate limited. Retry after {retry_after}s")

        resp.raise_for_status()
        return resp.json()

# Usage
client = DriftlockClient("dlk_your-key.secret")
result = client.detect([
    {"message": "normal event"},
    {"message": "CRITICAL ERROR"}
])
print(f"Found {result['anomaly_count']} anomalies")
```

## MCP Server Integration (Coming Soon)

For Claude Code and other AI agents that support MCP:

```json
{
  "mcpServers": {
    "driftlock": {
      "url": "https://driftlock.net/api/mcp",
      "auth": {
        "type": "apiKey",
        "header": "X-Api-Key"
      }
    }
  }
}
```

## Use Cases

### Log Anomaly Detection
Send structured logs, detect unusual patterns:
```json
{"events": [
  {"level": "info", "service": "api", "latency": 45},
  {"level": "error", "service": "api", "latency": 30000}
]}
```

### Metric Spike Detection
Detect if a metric change is statistically significant:
```json
{"events": [
  {"cpu": 45}, {"cpu": 48}, {"cpu": 47}, {"cpu": 95}
]}
```

### API Traffic Analysis
Find unusual request patterns:
```json
{"events": [
  {"endpoint": "/api/users", "method": "GET", "status": 200},
  {"endpoint": "/api/admin/delete-all", "method": "DELETE", "status": 403}
]}
```

## Support

- OpenAPI Spec: https://driftlock.net/docs/architecture/api/openapi.yaml
- Documentation: https://driftlock.net/docs
- Email: support@driftlock.net
