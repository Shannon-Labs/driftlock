# Remaining Issues / Blockers

## Credentials & Access
- Supabase dashboard access needed to set service role key and confirm schema (anomalies, usage_counters) constraints.
- Stripe dashboard access needed to create products/prices and set webhook secrets.
- Cloudflare account access required to deploy Workers/Pages and set secrets.

## Go API Integration
- Wired `/v1/anomalies` routes and SSE. Remaining: expand auth and make org scoping explicit across handlers.
- `/v1/events` now propagates `organization_id` and `event_type` to the detector; detector syncs to Supabase and meters usage on true anomalies. Remaining: richer event schemas and auth-derived org context.

## Supabase Sync
- `supabase.CreateAnomaly` call may fail if the table requires `organization_id`. Provide a reliable org/tenant ID from auth or request context.
- Usage metering currently uses fallback org ID; switch to authenticated tenant context or event payload field.

## Rate Limiting
- Workers rate limiting requires KV/Durable Object or `@cloudflare/workers-rate-limit` dependency. Add infra and wire once available.

## Testing & Load
- Run `make test` (Go) and `cargo test` (Rust) once network + toolchain available.
- Execute `k6` load tests targeting Worker/API to validate p95 < 200ms at 1k rps.
