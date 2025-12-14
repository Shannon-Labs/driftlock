# Driftlock API Reference

Base URL: `http://localhost:8080` (development) or your production URL

All versioned endpoints under `/v1/*`.

## Quick Start

```bash
# Build and run
cargo build -p driftlock-api --release
cargo run -p driftlock-api --release

# Test demo endpoint (no auth required)
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["normal log entry", "ERROR: unexpected failure"]}'
```

## Authentication

### API Key Authentication
Most endpoints require an API key via `X-Api-Key` header:

```bash
curl -H "X-Api-Key: your-api-key" http://localhost:8080/v1/anomalies
```

### Firebase JWT Authentication
Onboarding endpoints use Firebase JWT via `Authorization: Bearer <token>`:

```bash
curl -H "Authorization: Bearer <firebase-jwt>" http://localhost:8080/v1/auth/me
```

## Endpoint Summary

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/healthz` | GET | None | Liveness probe |
| `/readyz` | GET | None | Readiness probe |
| `/v1/version` | GET | None | Version info |
| `/metrics` | GET | None | Prometheus metrics |
| `/v1/waitlist` | POST | IP Rate | Email waitlist signup |
| `/v1/demo/detect` | POST | IP Rate | Anonymous demo detection |
| `/v1/detect` | POST | API Key | Authenticated detection |
| `/v1/anomalies` | GET | API Key | List anomalies |
| `/v1/anomalies/:id` | GET | API Key | Get anomaly details |
| `/v1/anomalies/:id/feedback` | POST | API Key | Submit feedback |
| `/v1/streams` | GET/POST | API Key | List/create streams |
| `/v1/streams/:id` | GET | API Key | Get stream |
| `/v1/streams/:id/profile` | GET/PATCH | API Key | Detection profile |
| `/v1/streams/:id/tuning` | GET | API Key | Tuning history |
| `/v1/streams/:id/anchor` | GET/DELETE | API Key | Anchor management |
| `/v1/streams/:id/anchor/details` | GET | API Key | Full anchor data |
| `/v1/streams/:id/reset-anchor` | POST | API Key | Create new anchor |
| `/v1/profiles` | GET | API Key | List detection profiles |
| `/v1/account` | GET/PATCH | API Key | Account management |
| `/v1/account/usage` | GET | API Key | Usage summary |
| `/v1/api-keys` | GET/POST | API Key | API key management |
| `/v1/api-keys/:id` | DELETE | API Key | Revoke key |
| `/v1/api-keys/:id/regenerate` | POST | API Key | Regenerate key |
| `/v1/auth/signup` | POST | Firebase | Create tenant |
| `/v1/auth/me` | GET | Firebase | Get current user |
| `/v1/billing/checkout` | POST | API Key | Create checkout |
| `/v1/billing/portal` | POST | API Key | Customer portal |
| `/v1/billing/webhook` | POST | Stripe Sig | Stripe webhooks |
| `/v1/me/billing` | GET | API Key | Billing status |
| `/v1/me/usage/details` | GET | API Key | Daily usage |
| `/v1/me/usage/ai` | GET | API Key | AI usage (stub) |
| `/v1/me/ai/config` | GET | API Key | AI config (stub) |

## Core Endpoints

### POST /v1/detect

Synchronous anomaly detection for authenticated users.

**Request:**
```json
{
  "stream_id": "stream_logs",
  "events": ["log entry 1", "log entry 2", "..."],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd"
  }
}
```

**Response:**
```json
{
  "anomalies": [
    {
      "id": "anom_abc123",
      "ncd": 0.72,
      "compression_ratio": 1.41,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "explanation": "Significant deviation from baseline pattern"
    }
  ],
  "metrics": {
    "processed": 100,
    "baseline": 400,
    "window": 50,
    "duration_ms": 42
  }
}
```

### POST /v1/demo/detect

Anonymous demo endpoint with IP-based rate limiting (10 req/hour).

**Request:** Same as `/v1/detect` but `stream_id` optional.

### GET /v1/anomalies

List anomalies with pagination and filtering.

**Query Parameters:**
- `limit` (default: 50, max: 200)
- `offset` (default: 0)
- `stream_id` - Filter by stream
- `min_ncd` - Minimum NCD score
- `max_p_value` - Maximum p-value
- `since` - ISO timestamp
- `until` - ISO timestamp

**Response:**
```json
{
  "anomalies": [...],
  "total": 1532,
  "limit": 50,
  "offset": 0
}
```

### GET /v1/profiles

List available detection profiles.

**Response:**
```json
{
  "profiles": [
    {
      "name": "sensitive",
      "description": "Higher sensitivity for security-critical streams",
      "ncd_threshold": 0.20,
      "p_value_threshold": 0.10,
      "window_size": 30,
      "baseline_size": 300
    },
    {
      "name": "balanced",
      "description": "Balanced detection for general use",
      "ncd_threshold": 0.30,
      "p_value_threshold": 0.05,
      "window_size": 50,
      "baseline_size": 400
    },
    {
      "name": "strict",
      "description": "Lower false positive rate",
      "ncd_threshold": 0.45,
      "p_value_threshold": 0.01,
      "window_size": 100,
      "baseline_size": 500
    },
    {
      "name": "custom",
      "description": "User-defined or auto-tuned settings"
    }
  ]
}
```

## Stream Management

### POST /v1/streams

Create a new stream.

**Request:**
```json
{
  "slug": "my-app-logs",
  "stream_type": "logs",
  "detection_profile": "balanced"
}
```

### GET /v1/streams/:id/anchor

Get anchor settings for drift detection.

**Response:**
```json
{
  "has_active_anchor": true,
  "anchor_enabled": true,
  "drift_ncd_threshold": 0.4,
  "anchor_reset_on_drift": false
}
```

### POST /v1/streams/:id/reset-anchor

Create a new anchor from provided events.

**Request:**
```json
{
  "events": ["baseline event 1", "baseline event 2", "..."],
  "compressor": "zstd"
}
```

## Billing Endpoints

### POST /v1/billing/checkout

Create a Stripe checkout session.

**Request:**
```json
{
  "plan": "starter",
  "success_url": "https://app.driftlock.net/success",
  "cancel_url": "https://app.driftlock.net/cancel"
}
```

**Response:**
```json
{
  "url": "https://checkout.stripe.com/c/pay/cs_test_123",
  "session_id": "cs_test_123"
}
```

### GET /v1/me/billing

Get current billing status.

**Response:**
```json
{
  "plan": "starter",
  "status": "free",
  "current_period_end": "2025-01-15T00:00:00Z",
  "events_used": 12500,
  "events_limit": 250000
}
```

## Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "baseline_size must be >= window_size",
    "request_id": "req_abc123"
  }
}
```

**Error Codes:**
- `unauthorized` - Missing or invalid authentication
- `forbidden` - Insufficient permissions
- `invalid_argument` - Invalid request parameters
- `not_found` - Resource not found
- `conflict` - Resource conflict
- `rate_limit_exceeded` - Too many requests
- `internal` - Server error

## Rate Limiting

Public endpoints (demo, waitlist) use IP-based rate limiting:
- `/v1/demo/detect`: 10 requests per hour per IP
- `/v1/waitlist`: 5 requests per hour per IP

Authenticated endpoints use tenant-based limits based on plan tier.

Rate limit headers included in responses:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`
- `Retry-After` (when rate limited)

## Prometheus Metrics

Available at `/metrics`:

```
# HELP driftlock_http_requests_total Total HTTP requests
# TYPE driftlock_http_requests_total counter
driftlock_http_requests_total{method="POST",path="/v1/detect",status="200"} 1234

# HELP driftlock_events_processed_total Total events processed
# TYPE driftlock_events_processed_total counter
driftlock_events_processed_total 567890

# HELP driftlock_anomalies_detected_total Total anomalies detected
# TYPE driftlock_anomalies_detected_total counter
driftlock_anomalies_detected_total 123

# HELP driftlock_detectors_active Active detector count
# TYPE driftlock_detectors_active gauge
driftlock_detectors_active 45
```
