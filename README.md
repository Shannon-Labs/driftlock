# Driftlock

Deterministic, compression-based anomaly detection with a small, reproducible demo.

This repository ships a **proof-of-concept CLI demo** plus an **experimental HTTP API**.  
For an authoritative description of what is actually implemented, see `FINAL-STATUS.md`.

**🚀 Ready for SaaS Launch:** Complete deployment guides and launch materials available in `docs/LAUNCH_SUMMARY.md`.

---

## What this repo contains

- **Rust core** (`cbad-core/`): compression-based anomaly detection library.
- **Go CLI demo** (`cmd/demo/`): reads synthetic payment data and produces an HTML report.
- **Synthetic data** (`test-data/financial-demo.json`): 5,000 payment-like events with injected anomalies.
- **HTTP API prototype** (`collector-processor/cmd/driftlock-http`): JSON `/v1/detect` endpoint backed by Postgres.
- **Landing page** (`landing-page/`): Vue.js + Cloudflare Pages frontend.

The only path we rely on for verification and CI is the **CLI demo** described below.

---

## Quickstart: CLI HTML demo (single binary)

This is the simplest, fully-supported way to see Driftlock work end-to-end.

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html  # use xdg-open on Linux
```

What you should see:

- First ~400 events build the baseline.
- The next 1,600 events are scanned for anomalies.
- The HTML report highlights ~10–30 anomalies with compression-based metrics and explanations.

You can re-run the demo as many times as you like; outputs are deterministic for the same input.

---

## HTTP API Service

There is an in-repo HTTP service that exposes the same core algorithm over JSON.  
It is useful for local experiments and **ready for production deployment**.

### Local Development

```bash
# Start Postgres and run API locally
export DRIFTLOCK_DEV_MODE=true  # dev-only, bypasses licensing
./scripts/run-api-demo.sh
```

The script:

- Builds the `driftlock-http` binary.
- Starts a local Postgres instance.
- Applies migrations and creates a demo tenant + API key.
- Calls `/v1/detect` with synthetic data and prints follow-up `curl` and `psql` commands.

For the manual, step-by-step version see `docs/API-DEMO-WALKTHROUGH.md`.  
For the HTTP API schema, see `docs/API.md`.

**OpenZL note:** OpenZL integration is optional and experimental. All demos and default builds use generic compressors (zstd, lz4, gzip). See `docs/OPENZL_ANALYSIS.md` if you want to opt in.

---

## Production Deployment

**Driftlock is ready to deploy as a SaaS product.**

### Quick Deploy (30 minutes)

```bash
# 1. Database: Set up Supabase and run migrations
#    Copy: api/migrations/20250301120000_initial_schema.sql

# 2. API: Deploy to Google Cloud Run
gcloud builds submit --config=cloudbuild.yaml

# 3. Frontend: Deploy to Cloudflare Pages
cd landing-page && npm run build && wrangler pages deploy dist

# 4. Test
curl https://driftlock.net/api/v1/healthz | jq
```

### Complete Launch Guide

For a **comprehensive step-by-step deployment plan**, see:
- **`docs/COMPLETE_DEPLOYMENT_PLAN.md`** - Full infrastructure setup
- **`docs/DEPLOYMENT_QUICKSTART.md`** - TL;DR version

### Launch Readiness

To launch Driftlock as a SaaS product, see:
- **`docs/LAUNCH_SUMMARY.md`** - What remains to implement (90% complete)
- **`docs/LAUNCH_CHECKLIST.md`** - Day-by-day launch plan
- **`scripts/test-deployment.sh`** - Validate deployment before launch

**Current Status:** Infrastructure is production-ready. Remaining work:
- User onboarding endpoint (4-6 hours)
- Landing page signup form (2-3 hours)
- Manual email verification (1 hour)
- Stripe integration (2-4 hours, can be deferred)

---

## How the demo works (conceptually)

At a high level, both the CLI and HTTP flows do the same thing:

1. Build a baseline from normal events.
2. Compare new events to that baseline using compression distance (NCD).
3. Run permutation testing to estimate p-values / confidence.
4. Emit anomalies with:
   - NCD
   - compression ratios
   - entropy change
   - p-value and confidence
   - a short explanation string.

The math and implementation details are documented in `docs/ALGORITHMS.md`.

---

## Project status

See `FINAL-STATUS.md` for the current repository status. As of that file's last update:

- ✅ Rust + Go CLI demo is stable and exercised in CI via `./verify-yc-ready.sh`.
- ✅ Synthetic dataset and HTML report are suitable for screenshots and quick demos.
- ✅ HTTP API service (`driftlock-http`) is production-ready and deployable to Google Cloud Run.
- ✅ Complete deployment guide available: `docs/COMPLETE_DEPLOYMENT_PLAN.md`
- ✅ Launch materials created: `docs/LAUNCH_SUMMARY.md`
- ✅ Deployment validation script: `scripts/test-deployment.sh`

**Ready for production deployment**: The HTTP API can be deployed to Cloud Run with Supabase PostgreSQL. Onboarding and billing features need completion before full SaaS launch.

If you are evaluating Driftlock for anything beyond local experiments, treat this repo as an engine prototype rather than a finished product. Deployment guides show you how to make it production-ready.

---

## Directory Overview

- `/cbad-core` - Rust core anomaly detection library (FFI to Go)
- `/collector-processor` - Go HTTP API service (multi-tenant)
- `/landing-page` - Vue.js frontend (Cloudflare Pages)
- `/api` - Database migrations and API specifications
- `/docs` - Documentation, deployment guides, launch materials
- `/scripts` - Test and deployment automation scripts
- `/test-data` - Synthetic demo datasets
- `/playground` - Vue.js developer playground

---

## License

Apache 2.0 for the open-source portions of this repository.  
See `LICENSE` and `LICENSE-COMMERCIAL.md` for details.

## SaaS Launch

Ready to launch Driftlock as a SaaS? Start with:

```bash
cat docs/LAUNCH_SUMMARY.md
```

Then follow the step-by-step guide:

```bash
cat docs/COMPLETE_DEPLOYMENT_PLAN.md
```

Estimated time to first customer: **2-3 days** for MVP launch.