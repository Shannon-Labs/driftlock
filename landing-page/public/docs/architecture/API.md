# Driftlock API (Launchable MVP Spec)

Base URL defaults to `http://localhost:8080` when running Docker Compose; production deployments front the service with HTTPS. Versioned under `/v1/*`.

## Demo Quickstart

- `cargo build --release` → produce the `cbad-core` shared library for Go FFI.
- `DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh` → guided script for humans. Builds the Go binary, starts Postgres via Docker Compose, runs migrations, creates a tenant/key, exercises `/v1/detect`, `/v1/anomalies/{id}`, and the export stubs, then prints follow-up commands.
- See [docs/API-DEMO-WALKTHROUGH.md](./API-DEMO-WALKTHROUGH.md) when you need to run each command manually (migrate, create-tenant, start server, curl `/v1/detect`, and query Postgres) for screen recordings or compliance evidence.

`DRIFTLOCK_DEV_MODE=true` is for **local demos only**. Any pilot or production deployment must set a signed `DRIFTLOCK_LICENSE_KEY` issued by Shannon Labs. `/healthz` reports the current license status so compliance teams can audit every environment.

## Authentication

- Every request (except `/healthz` and `/metrics` when exposed internally) must send `X-Api-Key: <secret>`.
- Keys are tenant-scoped secrets stored hashed (Argon2id + salt). Each key has:
  - `role=admin` → manage configs, API keys, exports.
  - `role=stream` → ingest events (`/v1/streams/*`), run synchronous detect, fetch anomaly listings.
  - Optional `stream_id` binding for least-privilege ingestion keys.
- Requests missing or presenting invalid keys return `401` with `error.code="unauthorized"`.
- Supabase Auth / Ory Kratos / SSO integrations are future work (Phase 3); document Okta/Azure AD as planned.
- **Admin CLI**: use `./driftlock-http create-tenant --name ... --key-role ...` to bootstrap tenants and keys. `migrate up` / `migrate status` operate against the configured `DATABASE_URL` using the bundled goose migrations under `api/migrations/`.

## Rate Limiting & Quotas

- Deterministic token-bucket per tenant enforced in-process (`RATE_LIMIT_RPS` env, default 60 req/sec). Admins can override per-tenant rates; API keys may optionally set lower overrides.
- Authenticated responses now include `X-RateLimit-Limit` and `X-RateLimit-Remaining` (minimum of tenant + key limits). Throttled responses add `Retry-After`, `X-RateLimit-Reset`, and a JSON envelope with `error.code="rate_limit_exceeded"` plus `retry_after_seconds`.
- Redis-backed buckets land in Phase 4; today’s build keeps the bucket in memory and is suitable for single-instance pilots.

## Endpoint Summary

| Endpoint | Method | Description | Auth |
|----------|--------|-------------|------|
| `/v1/streams/{streamId}/events` | POST | Batch ingest OTLP/log/metric events for async detection. | Stream/Admin |
| `/v1/detect` | POST | Run synchronous detect against a payload (CLI parity). | Stream/Admin |
| `/v1/anomalies` | GET | Paginated anomaly listing, filters, cursor pagination. | Stream/Admin |
| `/v1/anomalies/{id}` | GET | Fetch single anomaly with evidence metadata. | Stream/Admin |
| `/v1/anomalies/export` | POST | Spawn export job for filtered anomalies. | Admin |
| `/v1/anomalies/{id}/export` | POST | Generate evidence bundle for a single anomaly. | Stream/Admin |
| `/v1/config` | GET | Fetch tenant + stream configuration snapshot. | Admin |
| `/v1/config` | PUT | Update tenant + stream configuration (If-Match). | Admin |
| `/healthz` | GET | Health probe with CBAD build info, license status, queue backend, and database connectivity. | Optional |
| `/metrics` | GET | Prometheus metrics (restrict via network or token). | Optional |
| `/v1/stream/anomalies` | GET | SSE feed for real-time anomaly notifications. | Stream/Admin |

Future: OTLP HTTP/gRPC receiver compatibility (`/otlp/v1/*`) and Kafka ingestion (feature-flagged) arrive in Phase 4.

## Endpoint Details

### `POST /v1/streams/{streamId}/events`

Queues batches for worker processing. JSON schema:

```json
{
  "events": [
    {
      "timestamp": "2025-03-01T12:34:56.789Z",
      "type": "log|metric|trace|llm",
      "body": {},
      "attributes": {
        "service.name": "checkout",
        "region": "us-east-1"
      },
      "idempotency_key": "evt_123",
      "sequence": 1234567890
    }
  ],
  "ingest_options": {
    "async": true,
    "max_batch_bytes": 1048576
  }
}
```

- Limits: ≤256 events or 1 MiB per request (configurable per tenant).
- Response (202):

```json
{
  "queued": 128,
  "dropped": 0,
  "worker": "memory",
  "queue_depth": 32
}
```

- Validation errors → `400 invalid_argument`.

### `POST /v1/detect`

Synchronous detection for CLI-parity workloads.

```json
{
  "stream_id": "stream_logs",
  "events": [ ... same shape as ingest ... ],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd",
    "seed": 42
  }
}
```

Response mirrors CLI output:

```json
{
  "anomalies": [
    {
      "id": "anom_abc",
      "ncd": 0.72,
      "compression_ratio": 1.41,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "explanation": "Latency spike relative to baseline",
      "evidence_uri": "evidence/anom_abc/index.html"
    }
  ],
  "metrics": {
    "processed": 2000,
    "baseline": 400,
    "window": 50,
    "duration_ms": 4200
  }
}
```

Validation notes:
- Request bodies are capped at `MAX_BODY_MB` (default 10 MiB) and a maximum of `MAX_EVENTS` events (default 1,000). Empty or `null` events are rejected with `400 invalid_argument`.
- When callers request `openzl` but the binary was built without the OpenZL feature flag, the response falls back to `zstd` and includes `"fallback_from_algo": "openzl"`.

### `GET /v1/anomalies`

Query params:
- `limit` ≤ 200, default 50
- `page_token` (opaque cursor)
- `stream_id`, `min_ncd`, `max_p_value`, `status`, `since`, `until`, `has_evidence`

Response:

```json
{
  "anomalies": [
    {
      "id": "anom_abc",
      "tenant_id": "ten_bankalpha",
      "stream_id": "stream_logs",
      "ncd": 0.72,
      "compression_ratio": 1.41,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "status": "new",
      "explanation": "Latency spike relative to baseline",
      "detected_at": "2025-03-01T12:58:14Z"
    }
  ],
  "next_page_token": "eyJjcmVhdGVkX2F0IjoiMjAyNS0wMy0wMVQxMjo1ODoxNFoiLCJpZCI6ImFub21fYWJjIn0=",
  "total": 1532
}
```

### `GET /v1/anomalies/{id}`

Returns the metrics + evidence bundle for a single anomaly:

```json
{
  "id": "anom_abc",
  "stream_id": "stream_logs",
  "batch_id": "ingest_123",
  "status": "new",
  "detected_at": "2025-03-01T12:58:14Z",
  "explanation": "Latency spike relative to baseline",
  "metrics": {
    "ncd": 0.72,
    "compression_ratio": 1.41,
    "entropy_change": 0.13,
    "p_value": 0.004,
    "confidence": 0.96
  },
  "details": {
    "event": { "...raw payload..." },
    "metrics": { "...full cbad metrics..." }
  },
  "baseline_snapshot": {},
  "window_snapshot": {},
  "evidence": [
    {
      "format": "markdown",
      "uri": "local://evidence/anom_abc.md",
      "checksum": "ad54...9c",
      "size_bytes": 2048,
      "created_at": "2025-03-01T12:58:15Z"
    }
  ]
}
```

### `POST /v1/anomalies/export`

```json
{
  "format": "json|markdown|html",
  "filters": {
    "stream_id": "stream_logs",
    "status": "new",
    "date_range": {
      "start": "2025-02-01T00:00:00Z",
      "end": "2025-02-28T23:59:59Z"
    }
  },
  "delivery": {
    "type": "sync|webhook|s3",
    "webhook_url": "https://hooks.company.com/driftlock"
  }
}
```

Returns `202 Accepted` with a stub status (worker arrives in Phase 4):

```json
{
  "job_id": "exp_123",
  "status": "not_implemented",
  "message": "export worker queue stubbed; payload recorded for future worker"
}
```

`POST /v1/anomalies/{id}/export` behaves the same but pre-populates the filter with `{"anomaly_id": "<id>"}`.

### `GET/PUT /v1/config`

- `GET` returns tenant + stream configuration snapshot with `etag` for optimistic locking.

```json
{
  "etag": "W/\"config-42\"",
  "tenant": {
    "name": "Bank Alpha",
    "retention_days": 30,
    "default_compressor": "zstd",
    "rate_limit_rps": 60
  },
  "streams": [
    {
      "id": "stream_logs",
      "type": "logs",
      "ncd_threshold": 0.32,
      "p_value_threshold": 0.04,
      "baseline_size": 400,
      "window_size": 50,
      "hop_size": 10,
      "seed": 42,
      "compressor": "openzl",
      "retention_days": 14
    }
  ]
}
```

- `PUT` requires `If-Match` header matching the `etag`. On success returns `204 No Content` with new `etag` header.

### `GET /v1/stream/anomalies`

Server-Sent Events feed for live anomaly notifications.

- `Last-Event-ID` supported for resume.
- Events: `heartbeat` every 15 s, `anomaly` payload identical to `/v1/anomalies/{id}` summary.

### `GET /healthz`

Returns 200 when core dependencies are healthy and exposes compliance-friendly diagnostics. Fields of note:

- `openzl_available` (bool) — true when the binary includes OpenZL symbols.
- `available_algos` — includes `"openzl"` only when compiled in; requests for `openzl` will fall back to `zstd` when unavailable and surface `fallback_from_algo` in responses.
- `license.status` — must be `"valid"` outside dev; dev mode is for local demos only.

```json
{
  "build": {
    "revision": "5b93f09",
    "cbad_core": "0.3.2",
    "compressors": ["zstd", "lz4", "gzip", "openzl"]
  },
  "license": {
    "status": "dev_mode",
    "tier": "EVAL",
    "expires_at": "2025-06-01T00:00:00Z"
  },
  "database": "connected",
  "queue": {
    "backend": "memory",
    "pending": 0
  },
  "timestamp": "2025-03-01T13:00:00Z"
}
```

`/healthz?full=1` (admin key required) adds goose migration status, Redis connection info (when enabled), queue depth metrics, and the currently loaded OpenZL plans. Treat dev mode as local-only; production clusters must report `license.status="valid"`.

### `GET /metrics`

Prometheus exposition format with metrics:
- `driftlock_http_requests_total{path,method,status}`
- `driftlock_http_request_duration_seconds_bucket`
- `driftlock_queue_depth{tenant,stream}`
- `driftlock_anomalies_detected_total{tenant,stream}`
- `driftlock_detector_duration_seconds`

Restrict via network policy or require admin API key via reverse proxy.

## Pagination & Error Model

- Cursor-based pagination: `page_token` encodes `(tenant_id, stream_id, created_at, id)`.
- Errors share a common envelope:

```json
{
  "error": {
    "code": "invalid_argument",
    "message": "baseline_size must be >= window_size",
    "request_id": "req_abc123",
    "retry_after_seconds": 30
  }
}
```

- Codes: `unauthorized`, `forbidden`, `invalid_argument`, `not_found`, `conflict`, `rate_limit_exceeded`, `queue_unavailable`, `cbad_failure`, `export_timeout`, `internal`.

## Determinism & Explainability

- Seeds derive from `tenant_id`, `stream_id`, and config seed override ensuring reproducible detection.
- Every anomaly response includes: NCD, compression ratios, entropy delta, p-value, confidence, explanation string, baseline/window snapshots, evidence URIs.

## Licensing & Deployment Notes

- Running the API binary/container requires `DRIFTLOCK_LICENSE_KEY` (signed). The only exception is `DRIFTLOCK_DEV_MODE=true`, which is reserved for local demos and causes `/healthz` to report `license.status="dev_mode"`.
- `cbad-core` remains Apache 2.0; only the API binary/container is governed by the Driftlock Commercial License (see `LICENSE-COMMERCIAL.md`).
- Migrations (goose/Atlas) live in-repo; Supabase supplies managed Postgres/pgBouncer (default region `us-east-1`, override per-customer for residency).

## Roadmap Callouts

- OTLP HTTP receiver compatibility and Kafka ingestion are Phase 4 features; the interfaces are stubbed but disabled until the queue/Redis work ships.
- PDF exports (wkhtmltopdf/headless Chrome) land in Phase 3.
- Redis-backed rate limiting and queues become the production default once Phase 4 completes; OSS demo keeps in-process channels for simplicity.
