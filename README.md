# Driftlock

Driftlock is a **developer-first anomaly detection platform**: we use deterministic compression math (cbad-core) to flag weird behaviour in milliseconds, then optionally hand the anomaly off to Gemini to explain it in plain English. The repo contains everything needed to ship that experience as a SaaS product: core engine, HTTP service, landing page + dashboard, Firebase Auth, Stripe hooks, and deployment runbooks.

For an authoritative description of what ships today, see `.archive/reports/FINAL-STATUS.md` (snapshot) or `NEXT_STEPS.md` (current). The short version:

- ✅ Math engine + CLI demo for deterministic verification
- ✅ Multi-tenant HTTP API (`/v1/detect`) with Firebase Auth
- ✅ Landing page + dashboard (Vue 3 + Tailwind + Pinia)
- ✅ Pricing tiers (Developer free, Starter $25, Pro custom) wired to Stripe and the dashboard usage endpoints
- ✅ Gemini-based explainability layer ready to be toggled on per anomaly

**🚀 Ready for SaaS Launch:** Full Cloud SQL + Firebase Auth deployment setup available in `docs/deployment/CLOUDSQL_FIREBASE_SETUP_GUIDE.md`.

---

## What this repo contains

- **Rust core** (`cbad-core/`): compression-based anomaly detection library (FFI boundary to Go). The math is deterministic and seeded for reproducibility.
- **Go HTTP API service** (`collector-processor/cmd/driftlock-http`): multi-tenant `/v1/detect`, `/v1/anomalies`, `/v1/me/*` endpoints with Firebase Auth + Cloud SQL.
- **Usage tracker + pricing hooks** (`collector-processor/cmd/driftlock-http/dashboard.go`): exposes current plan limits so the dashboard and billing flows stay in sync.
- **Vue landing page + dashboard** (`landing-page/`): marketing site, signup flow, pricing section (Developer/Starter/Pro), dashboard with API keys + usage, docs viewer.
- **CLI demo** (`cmd/demo/`): reproducible HTML report for quick verification and CI.
- **Synthetic data** (`test-data/financial-demo.json`): 5,000 payment-like events with injected anomalies (used by both CLI + HTTP scripts).
- **Docs and runbooks** (`docs/`): organized by `architecture`, `deployment`, `launch`, `compliance`, and `development`.
- **SDKs (beta)** (`sdks/typescript`, `sdks/python`): thin clients for `/v1` detect/anomalies/healthz.

The CLI demo remains the fastest path to verify the engine locally, but the **primary product surface is the hosted HTTP API + dashboard**.

Running everything end-to-end? See `docs/deployment/QUICKSTART_GOD_MODE.md` for a guided “god mode” checklist covering the CLI demo, API, landing page, and deploy.

---

## Product Overview

Driftlock is purpose-built for teams who need **provable anomaly detection** without training data:

- **Math detects it**: cbad-core builds a baseline from your first ~400 events, then scores every new event using Normalized Compression Distance (NCD), permutation tests, and entropy deltas.
- **AI explains it**: when an anomaly survives the math, we optionally send the evidence payload + metrics to Gemini Flash. The response becomes the \"plain English\" field in dashboards, email alerts, or Slack posts.
- **Dev-first ergonomics**: deterministic CLI demo, REST API, Vue dashboard, Firebase Auth, Stripe billing, Cloud SQL deployment scripts.
- **Use cases**: financial compliance (EU DORA, FFIEC), DDOS/API abuse, AI agent monitoring, IoT/smart-home telemetry. See `docs/launch/USE_CASES.md`.

## Quickstart: CLI HTML demo (deterministic)

Still the fastest way to convince yourself the math works.

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

The HTTP layer is the canonical product surface. It wraps the Rust detector with tenant-aware auth, usage tracking, and pricing guardrails.

### Local Development

```bash
# Start Postgres and run API locally
export DRIFTLOCK_DEV_MODE=true  # dev-only, bypasses licensing
./scripts/run-api-demo.sh
```

The script:

- Builds the `driftlock-http` binary.
- Starts a local Postgres instance (matching Cloud SQL schema).
- Applies migrations and creates a demo tenant + API key.
- Calls `/v1/detect` with synthetic data, prints follow-up `curl` + `psql` commands, and shows how Gemini explainability would be triggered.

For the manual, step-by-step version see `docs/development/API_DEMO_WALKTHROUGH.md`.  
For the HTTP API schema, see `docs/architecture/API.md`.

**OpenZL note:** OpenZL integration is optional/experimental. Defaults stay on zstd/lz4/gzip for determinism and portability. To exercise OpenZL: `USE_OPENZL=true docker compose build driftlock-http && USE_OPENZL=true PREFER_OPENZL=true docker compose up driftlock-postgres driftlock-http`, then hit `/healthz` and look for `"openzl_available": true`. To prefer it when compiled in, set `PREFER_OPENZL=true` (requests can still override `compressor`). See `.archive/reports/OPENZL_ANALYSIS.md`.

### OpenZL Docker Workflow

To build and run with OpenZL enabled:

```bash
# Build with OpenZL support
USE_OPENZL=true docker compose build driftlock-http

# Run with OpenZL preference
USE_OPENZL=true PREFER_OPENZL=true docker compose up -d driftlock-postgres driftlock-http

# Verify OpenZL availability
curl -s http://localhost:8080/healthz | jq '.openzl_available'
```


---

## Production Deployment

**Driftlock is ready to deploy as a SaaS product.**

### Quick Deploy (30 minutes) - Cloud SQL + Firebase Auth

```bash
# 1. Set up Cloud SQL + Firebase Auth infrastructure
./scripts/setup-gcp-cloudsql-firebase.sh

# 2. Deploy to production
./scripts/deploy-production-cloudsql-firebase.sh

# 3. Test the deployment
./scripts/test-deployment-complete.sh
```

### Complete Launch Guide

For a **comprehensive step-by-step deployment plan**, see:
- **`docs/deployment/CLOUDSQL_FIREBASE_SETUP_GUIDE.md`** - Complete Cloud SQL + Firebase Auth setup
- **`docs/deployment/COMPLETE_DEPLOYMENT_PLAN.md`** - Legacy Supabase setup (deprecated)

### Launch Readiness

To launch Driftlock as a SaaS product, see:
- **`docs/launch/LAUNCH_SUMMARY.md`** - What remains to implement (90% complete)
- **`docs/launch/LAUNCH_CHECKLIST.md`** - Day-by-day launch plan
- **`scripts/test-deployment.sh`** - Validate deployment before launch

**Current Status:** Infrastructure is production-ready and fully wired to the new landing page/dashboard. Remaining work before GA:
- Harden onboarding verification emails (optional; Firebase handles most flows)
- Expand Gemini explainability presets (today it’s manual per deployment)
- Finish Stripe self-serve upgrade UX (backend endpoints exist; UI CTA shipping in Starter tier)

### Billing & Stripe Integration

Driftlock includes built-in Stripe support for subscription management.

**Required Secrets:**
- `STRIPE_SECRET_KEY`: Your Stripe Secret Key (`sk_...`)
- `STRIPE_PRICE_ID_PRO`: Price ID for the Pro plan (`price_...`)

**Setup:**
1. Create a Stripe account and a Product with a Price (Recurring).
2. Add secrets to Google Cloud Secret Manager (see `docs/deployment/GCP_SECRETS_CHECKLIST.md`).
3. Deploy using `cloudbuild.yaml`.


---

## Security

Driftlock is designed with security in mind. We follow best practices for handling secrets and API keys.

- **Backend Secrets:** All backend secrets, such as database credentials and third-party API keys, are managed through Google Secret Manager.
- **Frontend API Keys:** The Firebase API key for the frontend is also stored in Google Secret Manager and served securely to the client via a Cloud Function. This prevents the key from being exposed in the frontend code.
- **Vulnerability Reporting:** We take security vulnerabilities seriously. Please report any vulnerabilities to us privately.

For more details on our security policies and best practices, see `docs/architecture/SECURITY.md`.

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

The math and implementation details are documented in `docs/architecture/ALGORITHMS.md`.

---

## Project status

See `.archive/reports/FINAL-STATUS.md` for the snapshot repository status. As of that file's last update:

- ✅ Rust + Go CLI demo is stable and exercised in CI via `./scripts/verify-launch-readiness.sh`.
- ✅ Synthetic dataset and HTML report are suitable for screenshots and quick demos.
- ✅ HTTP API service (`driftlock-http`) is production-ready and deployable to Google Cloud Run.
- ✅ Complete deployment guide available: `docs/deployment/COMPLETE_DEPLOYMENT_PLAN.md`
- ✅ Launch materials created: `docs/launch/LAUNCH_SUMMARY.md`
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
cat docs/launch/LAUNCH_SUMMARY.md
```

Then follow the step-by-step guide:

```bash
cat docs/deployment/CLOUDSQL_FIREBASE_SETUP_GUIDE.md
```

Estimated time to first customer: **2-3 days** for MVP launch.
