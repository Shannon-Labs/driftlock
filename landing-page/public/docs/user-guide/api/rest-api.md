# REST API Reference

Complete reference for Driftlock's REST API. All endpoints are versioned under `/v1/` and require authentication via API key.

## Base URL

```text
https://api.driftlock.net/v1
```

## Authentication

All API requests (except `/healthz`) require an API key in the `X-Api-Key` header:

```bash
curl -H "X-Api-Key: YOUR_API_KEY" \
  https://api.driftlock.net/v1/anomalies
```

See [Authentication Guide](../getting-started/authentication.md) for details on creating and managing API keys.

## Rate Limiting

API keys are rate limited based on your plan:

| Plan | Requests/Minute |
|------|----------------|
| Free | 60 |
| Standard ($15/mo) | 300 |
| Pro ($100/mo) | 600 |
| Enterprise ($299/mo) | 1,000+ |

### Rate Limit Headers

Every response includes:

```text
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1672531200
```

When rate limited, you'll receive HTTP 429:

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "retry_after_seconds": 30
  }
}
```

## Common Headers

All requests should include:

```text
Content-Type: application/json
X-Api-Key: YOUR_API_KEY
```

All responses include:

```text
Content-Type: application/json
X-Request-ID: req_abc123
```

## Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable error message",
    "request_id": "req_abc123",
    "retry_after_seconds": 30
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `unauthorized` | 401 | Missing or invalid API key |
| `forbidden` | 403 | Insufficient permissions |
| `invalid_argument` | 400 | Invalid request parameters |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource conflict |
| `rate_limit_exceeded` | 429 | Rate limit exceeded |
| `internal` | 500 | Server error |

See [complete error reference](errors.md) for all error codes.

## Endpoints

### Detection

| Endpoint | Method | Description |
|----------|--------|-------------|
| [/v1/detect](./endpoints/detect.md) | POST | Run synchronous anomaly detection |
| /v1/streams/{id}/events | POST | Queue events for async detection |

### Anomalies

| Endpoint | Method | Description |
|----------|--------|-------------|
| [/v1/anomalies](./endpoints/anomalies.md) | GET | List anomalies with filtering |
| [/v1/anomalies/{id}](./endpoints/anomaly-detail.md) | GET | Get anomaly details |
| /v1/anomalies/export | POST | Export anomalies (bulk) |
| /v1/anomalies/{id}/export | POST | Export single anomaly |

### Dashboard

| Endpoint | Method | Description |
|----------|--------|-------------|
| /v1/me/keys | GET | List your API keys (Firebase Auth required) |
| /v1/me/usage | GET | Get usage statistics (Firebase Auth required) |

### Billing

| Endpoint | Method | Description |
|----------|--------|-------------|
| /v1/billing/checkout | POST | Create Stripe checkout session |
| /v1/billing/portal | POST | Get Stripe customer portal link |
| /v1/billing/webhook | POST | Stripe webhook handler (internal) |

### Health & Monitoring

| Endpoint | Method | Description |
|----------|--------|-------------|
| /healthz | GET | Health check (no auth required) |
| /metrics | GET | Prometheus metrics |

## Pagination

List endpoints use cursor-based pagination:

```json
{
  "anomalies": [...],
  "next_page_token": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMy0wMVQxMjo1ODoxNFoifQ==",
  "total": 1532
}
```

To get the next page:

```bash
curl "https://api.driftlock.net/v1/anomalies?page_token=eyJj..." \
  -H "X-Api-Key: YOUR_API_KEY"
```

## Filtering

Many endpoints support filtering via query parameters:

```bash
# Filter anomalies by stream
GET /v1/anomalies?stream_id=stream_logs

# Filter by NCD threshold
GET /v1/anomalies?min_ncd=0.5

# Filter by date range
GET /v1/anomalies?since=2025-01-01T00:00:00Z&until=2025-01-31T23:59:59Z

# Combine filters
GET /v1/anomalies?stream_id=stream_logs&min_ncd=0.5&status=new
```

## Versioning

The API is versioned via the URL path (`/v1/`). We maintain backward compatibility within a major version.

**Current version**: v1  
**Status**: Stable

Breaking changes will be introduced in v2, v3, etc. with advance notice.

## SDK & Libraries

- **Node.js/TypeScript**: [Node.js SDK](/docs/sdk/nodejs.md)
- **Python**: [Python SDK](/docs/sdk/python.md)
- **REST API Examples**: [cURL examples](./examples/curl-examples.md)

## Quick Examples

### Detect Anomalies

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "X-Api-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "default",
    "events": [
      {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"value": 100}}
    ]
  }'
```

### List Anomalies

```bash
curl https://api.driftlock.net/v1/anomalies \
  -H "X-Api-Key: YOUR_API_KEY"
```

### Get Anomaly Details

```bash
curl https://api.driftlock.net/v1/anomalies/anom_abc123 \
  -H "X-Api-Key: YOUR_API_KEY"
```

## Next Steps

- **[POST /v1/detect](./endpoints/detect.md)** - Detailed detection endpoint documentation
- **[GET /v1/anomalies](./endpoints/anomalies.md)** - Query and filter anomalies
- **[Code Examples](./examples/curl-examples.md)** - Complete working examples
- **[Error Codes](errors.md)** - Complete error reference

---

**Need help?** Check out our [tutorials](../tutorials/) or contact [support@driftlock.io](mailto:support@driftlock.io)
