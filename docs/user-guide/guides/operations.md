# Operations & Runbooks

Practical guidance for running Driftlock in production.

## Health & readiness
- `GET /healthz` — liveness (no auth)
- `GET /readyz` — readiness (if exposed)
- Watch error rates, p99 latency, and 5xxs on `/v1/detect` and `/v1/anomalies`.

## Deployment
- Preferred: Cloud Run + managed Postgres + Redis.
- Build once, deploy with immutable images; set `DATABASE_URL`, `STRIPE_*`, `SENDGRID_API_KEY`, `FIREBASE_PROJECT_ID`.
- For Firebase Hosting proxy, ensure `X-Api-Key` headers pass through when used.

## Scaling
- Batch events (1–256) to reduce QPS and rate-limit hits.
- Use `lz4` for very high throughput; `zstd` is the default balance.
- Horizontal scale API service; keep Redis/Postgres sized for write bursts.

## Observability
- Enable structured logs; include `request_id` and `stream_id`.
- Metrics to watch: request rate, latency, anomaly_count, rate_limit_exceeded, fallback_from_algo, error_rate.
- Alert on sustained 5xx, rising latency, or spikes in rate-limit responses.

## Backups & retention
- Daily Postgres backups; test restores.
- Retain anomalies/history per compliance requirements; document retention in your org policy.

## Runbooks
- 429s: confirm batching, inspect `Retry-After`, consider plan upgrade.
- Elevated 5xx: capture request IDs; check backend health; roll back recent deploys if correlated.
- Fallback from `openzl`: ensure CPU/memory headroom; consider locking to `zstd`/`lz4`.

## Security
- Rotate API keys every 90 days; bind stream-role keys where possible.
- Store secrets in Secret Manager/Secrets Manager; never ship in frontend code.

## Related
- [Cloud Run setup](../../deployment/cloud-run-setup.md)
- [Deployment runbooks](../../deployment/RUNBOOKS.md)
- [Authentication](../getting-started/authentication.md)
- [Error Codes](../api/errors.md)
