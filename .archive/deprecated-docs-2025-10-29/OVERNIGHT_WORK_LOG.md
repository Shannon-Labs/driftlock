# Overnight Work Log

Start Time: TBD

## Completed

- Fixed Go workspace `go.work` by removing missing `./productized` entry.
- Added Supabase usage metering method (`api-server/internal/supabase/client.go:MeterUsage`).
- Integrated best-effort usage metering into anomaly creation handler (`api-server/internal/handlers/anomalies.go`).
- Hardened PostgreSQL connection pooling (`api-server/internal/storage/postgres.go`).
- Added security headers middleware to Cloudflare Worker (`cloudflare-api-worker/src/index.ts`).
- Updated `.env.example` with Stripe, JWT, and OTEL variables.
- Added end-to-end integration script `test-integration-complete.sh`.
 - Wired anomalies REST and SSE routes; added readiness DB ping (`api-server/internal/api/server.go`).
 - Bootstrapped Postgres + Supabase + SSE in server main; enabled CBAD detector when available (`api-server/cmd/driftlock-api/main.go`).
 - Added contextual org_id/event_type propagation from `/v1/events` to detector (`api-server/internal/api/server.go`, `api-server/internal/ctxutil`).
 - Synced anomalies + metered usage from detector via Supabase on true anomalies (`api-server/internal/cbad/integration.go`).
- Added lightweight SQL migration runner (`tools/migrate`) and `make migrate`.
- Added optional API key auth for `/v1/events` with org derivation via `DEFAULT_API_KEY`/`DEFAULT_ORG_ID` (`auth` + server wiring).

## Notes

- Existing Supabase client already present; reused and extended instead of duplicating.
- Anomalies handler attempted to sync to Supabase; fields updated to align with current model.
- Usage metering uses `TENANT_DEFAULT_TENANT` or `SUPABASE_PROJECT_ID` env as fallback for `organization_id`.
  Detector path now prefers organization_id from `/v1/events` body when provided.

## Pending / External

- Build and deploy (requires network access and credentials).
- Stripe product/price IDs mapping in Supabase.
- Cloudflare Worker deployment + custom domain routing.
- End-to-end load testing and webhook validation.
