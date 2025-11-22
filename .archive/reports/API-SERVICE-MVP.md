# Driftlock API Service MVP – Phase 0 Decisions

This page captures the Phase 0 architecture decisions for bringing the current CLI/demo pipeline to a multi-tenant, launchable HTTP API. It narrows scope for the first production-ready milestone and lists the questions that must be answered before implementation begins.

## Scope & Goals

- Promote the deterministic CBAD core into a **tenant-scoped API service** without altering math behaviour proven by the CLI demo.
- Stick with shipping artifacts that already exist in the repository (Rust `cbad-core`, Go HTTP binary) and constrain the MVP to **one deployable Go service** plus managed dependencies (Postgres, optional Redis, Prometheus/Grafana).
- Provide enough infrastructure (migrations, CLI bootstrap, docs) so a pilot customer can ingest events, trigger detection, fetch anomalies, and export evidence inside Docker Compose, with a clear path to Helm/Kubernetes.
- Make the API-first experience the *public* story: README/DEMO landing pages point to `scripts/run-api-demo.sh` and [docs/API-DEMO-WALKTHROUGH.md](./API-DEMO-WALKTHROUGH.md) so anyone can reproduce the Postgres-backed flow.

## Platform & Deployment Decisions

| Area | Decision | Notes |
|------|----------|-------|
| Service language | **Go** (same repo, extend `collector-processor/cmd/driftlock-http`) | Keeps FFI surface with `cbad-core` untouched and reuses existing HTTP wiring. |
| Persistence | **Postgres 15+ (Supabase-hosted by default)** | Default `docker-compose.yml` gains a managed Postgres container; Supabase (us-east-1) hosts dev/prod Postgres/pgBouncer. Atlas/goose migrations tracked under `api/migrations/`. Deterministic ordering on all reads (`ORDER BY created_at, id`). |
| Cache / rate limit / short queues | **Redis (optional)** | Used for token-bucket rate limiting per tenant and as a lightweight job queue. Must remain optional; when absent, fall back to in-memory buckets + Go channels (implemented now for Phase 0 integration tests). |
| Metrics / dashboards | **Prometheus + Grafana** | `/metrics` exposes Prometheus format; Grafana dashboards live under `deploy/`. |
| Deployment targets | **Docker Compose (pilot) + Helm/K8s (enterprise)** | Compose stack runs API + Postgres (+ optional Redis, Grafana). Helm chart added under `deploy/helm/driftlock-api/` targeting managed Postgres/Redis services. |
| Migrations & CLI | **Goose migrations** + `./driftlock-api migrate up` | Reuse Go binary with Cobra CLI to expose `migrate`, `create-tenant`, `rotate-api-key` commands. |
| Licensing | **Rust core (`cbad-core`) remains Apache 2.0; new API binary/container ships under a source-available Driftlock Commercial License** | Keeps math core OSS while gating hosted binaries; evaluation builds stay available for pilots under time-limited keys validated via `DRIFTLOCK_LICENSE_KEY`. Dev-mode bypass (`DRIFTLOCK_DEV_MODE=true`) exists for local demos only and is surfaced via `/healthz`. |

## API Contract Overview

The API surface mirrors the CLI behaviour plus new multi-tenant capabilities. All endpoints live under `/v1/*` and require an API key unless noted. Pagination is cursor-based (`page_token`) with deterministic order guarantees.

| Endpoint | Method | Purpose | Auth |
|----------|--------|---------|------|
| `/v1/streams/{streamId}/events` | POST | Ingest OTLP/log batches; queue them for async detection. | `X-Api-Key` (tenant scope) |
| `/v1/detect` | POST | Synchronous detect for smaller payloads; parity with CLI demo. | `X-Api-Key` |
| `/v1/anomalies` | GET | Paginated anomaly listing per tenant/stream with filters. | `X-Api-Key` |
| `/v1/anomalies/{id}` | GET | Fetch one anomaly plus evidence metadata. | `X-Api-Key` |
| `/v1/anomalies/export` | POST | Batch export by filter (JSON/HTML/Markdown). | `X-Api-Key` |
| `/v1/anomalies/{id}/export` | POST | Generate an immutable evidence bundle for a single anomaly (initially Markdown/HTML, PDF later). | `X-Api-Key` |
| `/v1/config` | GET/PUT | Retrieve or update tenant + stream configuration (thresholds, compressors, retention). | `X-Api-Key` (admin key) |
| `/healthz` | GET | Liveness + dependency probes (CBAD core, DB, Redis/queue). | No auth (optional `?full=1` requires key). |
| `/metrics` | GET | Prometheus metrics (requires pod/service network policies for exposure). | Optional auth (token or network restricted). |

### JSON Schemas (Simplified)

#### `POST /v1/streams/{streamId}/events`

```json
{
  "events": [
    {
      "timestamp": "RFC3339Nano",
      "type": "log|metric|trace|llm",
      "body": {},                     // OTLP JSON payload or flattened fields
      "attributes": {"service.name": "api-gateway", "region": "us-east-1"},
      "idempotency_key": "uuid",     // ensures dedupe per tenant
      "sequence": 1234567890
    }
  ],
  "ingest_options": {
    "max_batch_bytes": 1048576,
    "queue": "redis|memory",
    "async": true
  }
}
```

- Response: `202 Accepted` with `{ "queued": N, "dropped": 0, "worker": "redis" }`.
- Size limits: 256 events or 1 MiB per request (configurable).
- Validation errors return `400` with machine-readable codes (`invalid_event_type`, `batch_too_large`).

#### `POST /v1/detect`

Direct pass-through of the CLI payload:

```json
{
  "stream_id": "uuid",
  "events": [...same shape as above...],
  "config_override": {
    "baseline_size": 400,
    "window_size": 50,
    "ncd_threshold": 0.3,
    "p_value_threshold": 0.05,
    "compressor": "zstd|lz4|openzl"
  }
}
```

- Response: `200 OK` with the anomaly decisions and evidence identical to the CLI HTML data model.
- Deterministic seeds derived from `stream_id` + `tenant_id` + `config.seed_override`.

#### `GET /v1/anomalies`

Query parameters: `page_token`, `limit (<=200)`, `stream_id`, `min_ncd`, `max_p_value`, `status`, `evidence` (bool), `since`, `until`.

Response envelope:

```json
{
  "anomalies": [
    {
      "id": "anom_...",
      "tenant_id": "ten_...",
      "stream_id": "str_...",
      "ncd": 0.72,
      "compression_ratio": 1.43,
      "entropy_change": 0.13,
      "p_value": 0.004,
      "confidence": 0.96,
      "explanation": "Latency spike relative to baseline",
      "status": "new|ack|exported",
      "created_at": "2025-02-10T14:05:00Z"
    }
  ],
  "next_page_token": "opaque",
  "total": 1234
}
```

`GET /v1/anomalies/{id}` responds with the above plus full evidence references.

#### `POST /v1/anomalies/export`

```json
{
  "format": "json|markdown|html|pdf",
  "filters": {
    "stream_id": "str_...",
    "status": "new",
    "date_range": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-31T23:59:59Z"
    }
  },
  "delivery": {
    "type": "sync|webhook|s3",
    "webhook_url": "https://..."
  }
}
```

Returns `202 Accepted` with an export job id. Phase 0 queues the request in-memory and responds with `status="not_implemented"` plus a TODO message; Phase 4 workers will hydrate the job and stream artifacts to the requested destination.

#### `GET/PUT /v1/config`

- `GET`: returns tenant + stream configs with `etag`.
- `PUT`: accepts `If-Match` header and payload:

```json
{
  "tenant": {
    "name": "Bank Alpha",
    "retention_days": 30,
    "default_compressor": "zstd",
    "rate_limit_rps": 50,
    "max_parallel_exports": 2
  },
  "streams": [
    {
      "id": "str_logs",
      "type": "logs",
      "ncd_threshold": 0.32,
      "p_value_threshold": 0.04,
      "baseline_size": 400,
      "window_size": 50,
      "hop_size": 10,
      "seed": 42,
      "compressor": "openzl",
      "queue": "redis",
      "retention_days": 14
    }
  ]
}
```

`/healthz` and `/metrics` mirror the existing HTTP binary but add DB/queue status and Prometheus counters (`driftlock_queue_depth`, `driftlock_anomalies_detected_total`, etc.).

### Error Model & Pagination

- Standardized error envelope:

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Tenant exceeded POST /v1/detect quota",
    "retry_after_seconds": 30,
    "request_id": "req_..."
  }
}
```

- Error codes: `unauthorized`, `forbidden`, `invalid_argument`, `conflict`, `not_found`, `rate_limit_exceeded`, `queue_unavailable`, `cbad_failure`, `export_timeout`.
- Pagination tokens encode `tenant_id`, `stream_id`, and the last `(created_at, id)` tuple to guarantee deterministic scans.

## Authentication & Tenancy

- **Primary auth**: tenant-scoped `X-Api-Key` header. Each tenant receives at least one secret; secrets stored hashed (Argon2id + per-key salt) with a `role` flag.
- **Roles**: `role=admin` keys may manage configs/keys; `role=stream` keys limited to ingestion/detect operations. Middleware enforces scope before routing.
- **Key lifecycle**: CLI exposes `create-tenant`, `create-key`, `revoke-key`; all actions audited. Keys can optionally be tied to a specific `stream_id` for least privilege.
- **Future auth hooks**: document Supabase Auth / Ory Kratos / Okta / Azure AD as planned Phase 3+ enhancements, but do not implement UI login or SSO yet.
- **Rate limiting**: token bucket per tenant; Phase 1–2 uses in-process counters. Redis-backed buckets arrive in Phase 4, with responses including `X-RateLimit-*` headers and `429` codes when exceeded.
- **Request scoping**: middleware resolves `tenant_id` from key, preloads configs, and validates stream membership before handlers execute.
- **mTLS / SSO**: note as roadmap for enterprise clusters; not required for MVP.

## Persistence Model

Use goose migrations with idempotent, deterministic scripts. Core tables:

| Table | Purpose | Key Fields |
|-------|---------|------------|
| `tenants` | Tenant metadata, licensing tier, retention defaults | `id`, `name`, `status`, `retention_days`, `default_compressor`, `plan`, `created_at` |
| `api_keys` | Auth keys hashed + scoped | `id`, `tenant_id`, `name`, `hash`, `role`, `rate_limit_rps`, `created_at`, `last_used_at` |
| `streams` | Logical streams per tenant | `id`, `tenant_id`, `type`, `description`, `seed`, `compressor`, `queue_mode` |
| `stream_configs` | Versioned configs (baseline/window/hop thresholds) | `id`, `stream_id`, `version`, `config_json`, `created_at`, `created_by` |
| `ingest_batches` | Metadata for ingested batches | `id`, `tenant_id`, `stream_id`, `batch_hash`, `queued_at`, `status`, `worker` |
| `anomalies` | Detector outputs | `id`, `tenant_id`, `stream_id`, metric columns (NCD, ratios, entropy, p-value, confidence), `details_json`, `status`, `detected_at` |
| `anomaly_evidence` | Rendered HTML/Markdown plus attachments | `id`, `anomaly_id`, `format`, `uri` (S3/path), `checksum`, `size_bytes` |
| `export_jobs` | Async export tracking | `id`, `tenant_id`, `filters_json`, `format`, `status`, `dest`, `created_at`, `completed_at` |

The repository layer lives under `api/internal/storage/` with context-aware methods returning deterministic `[]Model` slices. All writes inside transactions; detectors append evidence + anomaly rows atomically.

## Runtime & Workers

- `POST /v1/streams/{streamId}/events` writes batch metadata + pushes IDs into a queue:
  - **Phase 1–2**: buffered Go channels (per-tenant + per-stream) to avoid external dependencies.
  - **Phase 4+**: optional Redis Streams backend (`driftlock:{tenant}:{stream}` consumer groups) and Redis-powered rate limiting; documented as recommended for production deployments.
- Worker goroutines read queue items, fetch payload from Postgres (or inline if memory queue), call CBAD detector with deterministic config, write anomalies/evidence, and update `queue_depth` gauge.
- Seeds and config snapshots stored alongside each detection to guarantee reproducibility.
- Kafka ingestion remains a stubbed interface (feature-flagged) until full implementation in Phase 4.

## Monitoring & Ops

- `/metrics` exports: request latency histograms, queue depth gauges, anomaly counters, detector durations, DB connection pool stats.
- `/healthz` returns per-dependency probes plus optional `diagnostics` array when `?full=1` (requires admin key).
- Structured logging via `zap` with fields: `tenant_id`, `stream_id`, `request_id`, `anomaly_id`, `config_version`.
- Alert hooks: webhook configuration per tenant (stored in `stream_configs`) invoked when anomaly confidence > threshold.

## Evidence & Export Strategy

- **Storage**: local disk (`./evidence/`) under Docker Compose for development; S3-compatible storage (customer S3 bucket or MinIO) for Helm/K8s deployments. Provide Helm `values.yaml` example with IAM/S3 credentials.
- **Formats**: render Markdown as the source of truth and ship HTML downloads immediately. PDF generation waits until Phase 3 (wkhtmltopdf/headless Chrome).
- Each anomaly row references immutable evidence versions; exports gather references and optionally package as ZIP.
- Retention jobs (Go cron + SQL) purge anomalies/evidence beyond `retention_days` while respecting legal hold flags.

## Open Questions & Follow-Ups

1. **SSO roadmap**: list Okta/Azure AD as “planned” in docs, but no implementation until Phase 3 when UI logins matter.
2. **Redis rollout**: ensure documentation clearly states Redis optional for OSS/demo and recommended for production once Phase 4 ships.
3. **Evidence backend**: supply Terraform/IaC examples for customer-managed S3 buckets and MinIO to accelerate production installs.
4. **License ops**: finalize signing tool + evaluation key issuance process backing `DRIFTLOCK_LICENSE_KEY` before first commercial pilot.
5. **Kafka consumer interface timing**: stub remains behind a build flag; revisit once ingestion workers stabilize in Phase 4.

## Supabase Defaults & Deployment Notes

- Managed Postgres lives in Supabase `us-east-1` by default; document data residency options for DORA/HIPAA customers and allow overriding connection strings per environment.
- Keep goose/Atlas migrations in-repo; Supabase is only the managed Postgres host. Avoid Supabase migration tooling to maintain portability.
- Supabase connection secrets feed the API via environment variables (`DATABASE_URL`, optional `PGBOUNCER_URL`). Compose defaults to local Postgres for offline development.

## Licensing & Distribution

- `cbad-core` remains Apache 2.0 (OSS) to keep the math transparent.
- The Go API binary/container is distributed under the Driftlock Commercial License. Deployments require a signed `DRIFTLOCK_LICENSE_KEY` environment variable checked on startup; the key encodes expiry + tier.
- Provide an automatic evaluation key path (e.g., `DRIFTLOCK_LICENSE_KEY=EVAL-<signature>`) that expires after N days. Server logs impending expiry and refuses to start once expired.
- `LICENSE-COMMERCIAL.md` documents the terms and enforcement expectations alongside activation instructions.

Answering these items unblocks Phase 1 work: updating `docs/API.md`, scaffolding migrations, and teaching the HTTP binary about Postgres-backed tenants/configuration without regressing the CLI demo or determinism guarantees.
