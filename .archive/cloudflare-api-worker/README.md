# Cloudflare API Worker (Removed)

This directory previously contained a Cloudflare Worker that proxied Supabase API calls.
It has been removed as part of making DriftLock fully standalone without Supabase dependencies.

For OSS deployments, use the main API server directly. The API server supports API key authentication
and does not require Supabase.

If you need edge deployment, consider:
- Deploying the API server to Cloudflare Workers using the Go runtime
- Using Cloudflare Pages for the frontend
- Setting up a reverse proxy (nginx, Caddy) in front of the API server

