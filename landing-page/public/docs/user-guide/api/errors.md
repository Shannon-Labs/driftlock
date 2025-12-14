# Error Codes

Standard error payload:

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

| Code | HTTP | When it happens | How to fix |
| --- | --- | --- | --- |
| `unauthorized` | 401 | Missing or invalid `X-Api-Key` | Send a valid API key; rotate if leaked |
| `forbidden` | 403 | Key lacks required role (e.g., stream vs admin) | Use an admin key or adjust key role |
| `invalid_argument` | 400 | Payload/schema error, missing required fields | Validate request body and types |
| `not_found` | 404 | Resource does not exist (stream/anomaly) | Verify IDs and environment |
| `conflict` | 409 | Duplicate idempotency key or state conflict | Use new idempotency key or resolve race |
| `rate_limit_exceeded` | 429 | Plan limit exceeded | Back off, respect `retry_after_seconds`, or upgrade plan |
| `internal` | 500 | Server error | Retry with jitter; contact support if persistent |

## Troubleshooting

- Always log `request_id` from responses; include it in support requests.
- 429s: implement exponential backoff and honor `retry_after_seconds`.
- 4xx validation: ensure `events` array is present and within limits (1–256 events).
- 401: confirm header casing (`X-Api-Key`) and that the key is active.

## Related
- [REST API Reference](./rest-api.md)
- [POST /v1/detect](./endpoints/detect.md)
- [GET /v1/anomalies](./endpoints/anomalies.md)
