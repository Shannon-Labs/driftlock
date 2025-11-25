# Error Codes Reference

All API errors follow a consistent format with an error code, message, and request ID for debugging.

## Error Response Format

```json
{
  "error": {
    "code": "error_code",
    "message": "Human-readable description",
    "request_id": "req_abc123",
    "retry_after_seconds": 30
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `code` | string | Machine-readable error code |
| `message` | string | Human-readable error description |
| `request_id` | string | Unique ID for support/debugging |
| `retry_after_seconds` | integer | (Rate limits only) Seconds to wait before retry |

---

## Authentication Errors (4xx)

### `unauthorized`
**HTTP Status:** 401

The request is missing or has an invalid API key.

```json
{
  "error": {
    "code": "unauthorized",
    "message": "Invalid or missing API key",
    "request_id": "req_abc123"
  }
}
```

**Causes:**
- Missing `X-Api-Key` header
- Invalid or expired API key
- API key has been revoked

**Resolution:**
1. Ensure `X-Api-Key` header is included
2. Verify your API key in the [Dashboard](https://driftlock.net/dashboard)
3. Generate a new key if yours was revoked

---

### `forbidden`
**HTTP Status:** 403

The API key doesn't have permission for this operation.

```json
{
  "error": {
    "code": "forbidden",
    "message": "API key does not have permission for this operation",
    "request_id": "req_abc123"
  }
}
```

**Causes:**
- API key role doesn't permit the action
- Account is suspended
- Trial has expired without subscription

**Resolution:**
1. Check your API key's role/permissions
2. Upgrade your plan if needed
3. Contact support if account is suspended

---

## Rate Limiting Errors

### `rate_limit_exceeded`
**HTTP Status:** 429

You've exceeded your rate limit.

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "request_id": "req_abc123",
    "retry_after_seconds": 30
  }
}
```

**Rate Limits by Plan:**

| Plan | Requests/minute | Events/month |
|------|-----------------|--------------|
| Demo (no auth) | 10 | N/A |
| Pilot (Free) | 60 | 10,000 |
| Radar | 300 | 500,000 |
| Lock | 1,000+ | 5,000,000 |

**Resolution:**
1. Wait for `retry_after_seconds` before retrying
2. Implement exponential backoff
3. Batch events to reduce request count
4. Upgrade your plan for higher limits

**Example retry logic (Python):**

```python
import time
import requests

def detect_with_retry(events, max_retries=3):
    for attempt in range(max_retries):
        response = requests.post(
            "https://api.driftlock.net/v1/detect",
            headers={"X-Api-Key": API_KEY},
            json={"events": events}
        )

        if response.status_code == 429:
            retry_after = response.json()["error"].get("retry_after_seconds", 60)
            time.sleep(retry_after)
            continue

        return response.json()

    raise Exception("Max retries exceeded")
```

---

### `usage_limit_exceeded`
**HTTP Status:** 429

You've exceeded your monthly event quota.

```json
{
  "error": {
    "code": "usage_limit_exceeded",
    "message": "Monthly event limit exceeded. Upgrade your plan for more capacity.",
    "request_id": "req_abc123"
  }
}
```

**Resolution:**
1. Check your usage in the [Dashboard](https://driftlock.net/dashboard)
2. Wait for your billing period to reset
3. Upgrade your plan for more events

---

## Validation Errors (400)

### `invalid_argument`
**HTTP Status:** 400

The request contains invalid parameters.

**Common messages:**

| Message | Cause | Fix |
|---------|-------|-----|
| `events required` | Empty events array | Provide at least 1 event |
| `event N is empty` | Null or empty event at index N | Remove or fix the event |
| `too many events: max 256 per request` | Too many events | Split into multiple requests |
| `invalid json: ...` | Malformed JSON | Validate JSON syntax |

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "events required",
    "request_id": "req_abc123"
  }
}
```

**Resolution:**
1. Validate your JSON structure
2. Ensure events array is not empty
3. Keep events under 256 per request (50 for demo)

---

### `invalid_stream`
**HTTP Status:** 400

The specified stream doesn't exist or you don't have access.

```json
{
  "error": {
    "code": "invalid_stream",
    "message": "Stream not found or access denied",
    "request_id": "req_abc123"
  }
}
```

**Resolution:**
1. Verify the stream_id exists
2. Check you have access to the stream
3. Omit stream_id to use the default stream

---

## Not Found Errors (404)

### `not_found`
**HTTP Status:** 404

The requested resource doesn't exist.

```json
{
  "error": {
    "code": "not_found",
    "message": "anomaly not found",
    "request_id": "req_abc123"
  }
}
```

**Causes:**
- Invalid anomaly/resource ID
- Resource was deleted
- Resource belongs to a different tenant

**Resolution:**
1. Verify the resource ID
2. Check if it was deleted
3. Ensure you're using the correct API key

---

## Server Errors (5xx)

### `internal_error`
**HTTP Status:** 500

An unexpected error occurred on the server.

```json
{
  "error": {
    "code": "internal_error",
    "message": "Internal server error",
    "request_id": "req_abc123"
  }
}
```

**Resolution:**
1. Note the `request_id` for support
2. Retry after a brief delay
3. If persistent, contact support@driftlock.io

---

### `service_unavailable`
**HTTP Status:** 503

The service is temporarily unavailable.

```json
{
  "error": {
    "code": "service_unavailable",
    "message": "Service temporarily unavailable",
    "request_id": "req_abc123"
  }
}
```

**Resolution:**
1. Retry with exponential backoff
2. Check [status.driftlock.net](https://status.driftlock.net) for outages
3. Contact support if persists

---

## Method Errors

### `method_not_allowed`
**HTTP Status:** 405

Wrong HTTP method for the endpoint.

```json
{
  "error": {
    "code": "method_not_allowed",
    "message": "method not allowed",
    "request_id": "req_abc123"
  }
}
```

**Correct Methods:**

| Endpoint | Method |
|----------|--------|
| `/v1/detect` | POST |
| `/v1/demo/detect` | POST |
| `/v1/anomalies` | GET |
| `/v1/anomalies/{id}` | GET |

---

## Best Practices

### 1. Always Check for Errors

```python
response = requests.post(url, json=payload, headers=headers)

if response.status_code >= 400:
    error = response.json().get("error", {})
    print(f"Error {error.get('code')}: {error.get('message')}")
    print(f"Request ID: {error.get('request_id')}")
else:
    data = response.json()
    # Process success response
```

### 2. Implement Retry Logic

```python
from tenacity import retry, stop_after_attempt, wait_exponential

@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=1, min=1, max=60)
)
def detect_events(events):
    response = requests.post(url, json={"events": events}, headers=headers)

    if response.status_code == 429:
        retry_after = response.json()["error"].get("retry_after_seconds", 30)
        raise Exception(f"Rate limited, retry after {retry_after}s")

    response.raise_for_status()
    return response.json()
```

### 3. Log Request IDs

Always log the `request_id` for debugging:

```python
import logging

logger = logging.getLogger(__name__)

response = requests.post(url, json=payload, headers=headers)

if response.status_code >= 400:
    error = response.json().get("error", {})
    logger.error(
        "Driftlock API error",
        extra={
            "error_code": error.get("code"),
            "message": error.get("message"),
            "request_id": error.get("request_id"),
        }
    )
```

### 4. Handle Common Scenarios

```python
def handle_driftlock_response(response):
    if response.status_code == 200:
        return response.json()

    error = response.json().get("error", {})
    code = error.get("code")

    if code == "rate_limit_exceeded":
        # Implement backoff
        raise RateLimitError(error)
    elif code == "unauthorized":
        # Re-authenticate or refresh key
        raise AuthError(error)
    elif code == "usage_limit_exceeded":
        # Alert user to upgrade
        raise QuotaError(error)
    elif response.status_code >= 500:
        # Retry with backoff
        raise ServerError(error)
    else:
        # Log and raise
        raise DriftlockError(error)
```

---

## Need Help?

If you encounter persistent errors:

1. **Check your request** against our [API Reference](./endpoints/detect.md)
2. **Search our docs** for specific error messages
3. **Contact support** at support@driftlock.io with:
   - Your `request_id`
   - The full error response
   - What you were trying to do

---

**Next**: [Quickstart Guide →](../getting-started/quickstart.md)
