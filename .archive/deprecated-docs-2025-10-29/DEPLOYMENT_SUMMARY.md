# Deployment Summary

This file tracks key endpoints, configuration, and secrets required for production.

## Endpoints
- API (Go, local dev): `http://localhost:8080`
- Ingestion: `POST /v1/events`
- Anomalies API: `GET/POST /v1/anomalies`, `GET /v1/anomalies/{id}`, `PATCH /v1/anomalies/{id}/status`
- SSE Stream: `GET /v1/stream/anomalies`
- API Gateway (Workers): `https://driftlock-api.<account>.workers.dev`
- Frontend (Pages): `https://driftlock-web.pages.dev`
- Supabase project: `https://nfkdeeunyvnntvpvwpwh.supabase.co`

## Required Secrets
- Supabase: `SUPABASE_SERVICE_ROLE_KEY`, `SUPABASE_ANON_KEY`
- Stripe: `STRIPE_SECRET_KEY`, `STRIPE_PUBLISHABLE_KEY`, `STRIPE_WEBHOOK_SECRET`
- Worker (wrangler secrets): `SUPABASE_SERVICE_ROLE_KEY`, `STRIPE_WEBHOOK_SECRET`
 - API (optional dev): `DEFAULT_API_KEY`, `DEFAULT_ORG_ID`

## Env Vars
- `.env.example` updated; copy to `.env` per environment with real values.

## Deploy Steps (High-Level)
1. Cloudflare Workers: `cd cloudflare-api-worker && npx wrangler deploy` (set secrets before deploy).
2. Cloudflare Pages: build frontend and `wrangler pages deploy` to project `driftlock-web`.
3. Configure custom domains and SSL in Cloudflare.
4. Update Stripe webhook to point to Worker: `/stripe-webhook`.
