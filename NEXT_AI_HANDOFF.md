Next AI Handoff

Purpose: Provide a focused, execution-ready handoff for the next AI/agent to finish deployment, integrate payments, and validate production readiness.

Repo State Summary
- API builds and runs. Endpoints:
  - GET `/healthz`, GET `/readyz`, GET `/v1/version`
  - POST `/v1/events` (ingestion; optional API key auth)
  - GET/POST `/v1/anomalies`, GET `/v1/anomalies/{id}`, PATCH `/v1/anomalies/{id}/status`
  - GET `/v1/stream/anomalies` (SSE)
- Detection path: ingestion → CBAD detector → Postgres → SSE → Supabase sync → meter-usage Edge Function.
- Supabase client is minimal HTTP-based and supports REST + Edge Functions. Readiness checks Supabase health (best-effort).
- Cloudflare Worker has security headers + basic in-memory rate limiting.

High-Impact Next Steps
1) Deploy Cloudflare Worker (API gateway)
   - cd `cloudflare-api-worker`
   - `npm install && npx wrangler login`
   - `npx wrangler secret put SUPABASE_SERVICE_ROLE_KEY`
   - `npx wrangler secret put STRIPE_WEBHOOK_SECRET`
   - `npx wrangler deploy --env staging` then `--env production`
   - Update Stripe webhook URL → `https://driftlock-api.<account>.workers.dev/stripe-webhook`

2) Deploy Web Frontend (Cloudflare Pages)
   - cd `web-frontend`
   - `npm install && npm run build`
   - `npx wrangler pages project create driftlock-web` (first time)
   - `npx wrangler pages deploy dist --project-name driftlock-web`
   - Configure env vars in Pages: `VITE_SUPABASE_URL`, `VITE_SUPABASE_ANON_KEY`, `VITE_STRIPE_PUBLISHABLE_KEY`

3) Stripe Setup (test mode)
   - Create Standard/Growth products and metered overage prices
   - Update Supabase `plan_price_map` with actual Stripe IDs
   - Add `STRIPE_WEBHOOK_SECRET` to Supabase Edge Function secrets if needed

4) Supabase
   - Confirm tables: `anomalies`, `usage_counters`
   - Edge Functions: `health`, `meter-usage`, `stripe-webhook`, `send-alert-email`
   - Confirm RLS and schema match frontend queries

5) API auth + org scoping
   - For production, set `DEFAULT_API_KEY` and `DEFAULT_ORG_ID` (or implement DB-backed API keys); enforce required auth on `/v1/events`.
   - Ensure org id is derived from key or auth context — payload override should be disallowed unless trusted.

6) Tests and load
   - Go tests: `make test` (some test packages outside the API path need module fixes)
   - Rust tests in `cbad-core`: `cargo test`
   - Load tests: install `k6`, then run `k6 run tests/load/driftlock-load-test.js`

Local Runbook
- Env: `cp .env.example .env` and fill values.
- Migrate DB: `make migrate` (requires `DATABASE_URL` or DB_* vars)
- Start API: `make run`; check `/healthz` and `/readyz`
- Ingest sample:
  ```bash
  curl -X POST http://localhost:8080/v1/events \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${DEFAULT_API_KEY:-testkey}" \
    -d '{"organization_id":"org_123","event_type":"log","data":{"message":"test","level":"INFO"}}'
  ```
- Stream anomalies: `curl -N http://localhost:8080/v1/stream/anomalies`

Secrets Checklist
- `.env`: `SUPABASE_BASE_URL`, `SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_ROLE_KEY`, `DEFAULT_API_KEY`, `DEFAULT_ORG_ID`, Stripe keys (test mode)
- Worker Secrets: `SUPABASE_SERVICE_ROLE_KEY`, `STRIPE_WEBHOOK_SECRET`
- Pages Vars: `VITE_SUPABASE_URL`, `VITE_SUPABASE_ANON_KEY`, `VITE_STRIPE_PUBLISHABLE_KEY`

Validation Script
- Run `./test-integration-complete.sh` once API is up and Supabase vars are set.

References
- `README.md` (Quick Start + Endpoints)
- `CLOUDFLARE_DEPLOYMENT.md` and `README_CLOUDFLARE.md`
- `DEPLOYMENT_SUMMARY.md`, `REMAINING_ISSUES.md`, `OVERNIGHT_WORK_LOG.md`

