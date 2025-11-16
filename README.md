# Driftlock: Explainable Anomaly Detection for Global Banks

Driftlock ships today as a deterministic, compression-based anomaly engine with a guided API demo, CLI report, and ROI story that make sense to YC reviewers, senior engineers, and developers alike.

- **YC reviewers / investors** see a live pipeline: Rust core + Go API + Postgres + HTML evidence, all reproducible via scripts and CI.
- **Senior engineers** get a math-first detector that exposes NCD, compression ratios, entropy, p-value, confidence, and explanations for every anomaly.
- **Developers** run a single script to mint a tenant, hit `/v1/detect`, and inspect persisted anomalies—no black-box ML or vendor lock-in.

**EU DORA**: Up to 2% of annual turnover in fines | **US Regulations**: CFPB Fair Lending, NYDFS Part 500, FFIEC Guidelines

![Driftlock API demo – terminal session](screenshots/api-demo-demo.gif)

## Quickstart — HTTP API Demo (Docker + Postgres)

Run the same flow we show in YC and pilot calls. It builds the Rust core, the Go HTTP API, and a dockerized Postgres, then walks you through `/v1/detect`, `/v1/anomalies/{id}`, and `/healthz`.

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
git submodule update --init --recursive
cargo build --release
DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
```

The script prints ready-to-run `curl`, `psql`, and `/healthz` commands so you can keep exploring immediately. Prefer a manual walkthrough? See [docs/API-DEMO-WALKTHROUGH.md](docs/API-DEMO-WALKTHROUGH.md) or watch the landing-page video (`landing-page/public/api-demo-demo.mp4`) for the same flow.

## Try Driftlock via the HTTP API (Docker + Postgres)

If you prefer to step through what `./scripts/run-api-demo.sh` automates, run each command yourself:

1. Clone the repo and pull submodules:

   ```bash
   git clone https://github.com/Shannon-Labs/driftlock.git
   cd driftlock
   git submodule update --init --recursive
   ```

2. Build the cbad-core artifacts (required for all Go FFI binaries):

   ```bash
   cargo build --release
   ```

3. Run the end-to-end API demo. The script builds `driftlock-http`, starts Postgres via Docker Compose, runs migrations, creates a tenant + API key, and hits `/v1/detect`, `/v1/anomalies/{id}`, and the export stubs. Dev mode is enabled automatically if no license key is present.

   ```bash
   DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
   ```

   The script prints ready-to-run `curl`, `psql`, and `/healthz` commands so you can continue exploring the API immediately.

🎥 **Watch the full run** — the terminal session below shows `./scripts/run-api-demo.sh` provisioning Postgres, creating a tenant, calling `/v1/detect`, and surfacing the follow-up curl/psql commands. The responses in the recording use the latest anomaly metrics (NCD, compression ratios, entropy change, p-value, confidence, and explanations):

![Driftlock API demo – terminal session](screenshots/api-demo-demo.gif)

1. Prefer to run the commands manually? Follow [docs/API-DEMO-WALKTHROUGH.md](docs/API-DEMO-WALKTHROUGH.md) for the step-by-step version (docker compose up, `driftlock-http migrate up`, `create-tenant`, `curl /v1/detect`, `psql` queries, etc.).

## Manual HTTP API Walkthrough

The Go HTTP service (`cmd/driftlock-http`) is still the canonical API for integrating pilot workloads. The commands below mirror what `scripts/run-api-demo.sh` automates when you need to drive each step yourself:

```bash
# Required for all API flows
export DRIFTLOCK_LICENSE_KEY="<signed license from Shannon Labs>"
export DATABASE_URL="postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable"

# Or use dev mode locally (bypasses license validation, not for production)
export DRIFTLOCK_DEV_MODE=true

docker compose up -d driftlock-postgres    # managed Postgres container
./bin/driftlock-http migrate up            # apply goose migrations
./bin/driftlock-http create-tenant --name "Demo" --key-role admin --json
docker compose up --build driftlock-http   # HTTP API

# Verify health and available compressors (OpenZL is optional and reported separately)
curl -s http://localhost:8080/healthz | jq
# Run /v1/detect, fetch anomalies, and query Postgres as shown in docs/API-DEMO-WALKTHROUGH.md
```

### Demo Environment Variables

| Variable | Purpose | When to set |
|----------|---------|-------------|
| `DRIFTLOCK_DEV_MODE=true` | Bypasses license enforcement for local demos. `/healthz` reports `license.status="dev_mode"`. | Local development only. Production deployments **must** use `DRIFTLOCK_LICENSE_KEY`. |
| `DRIFTLOCK_LICENSE_KEY` | Signed evaluation or commercial key. Enables production mode licensing. | Required outside of dev mode (CI, pilots, prod). |
| `DATABASE_URL` | Connection string for Postgres (default `postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable`). | Always. Point to Compose, Supabase, or your managed Postgres. |
| `INTEGRATION_API_PORT` | Force the API port for `scripts/test-integration.sh` / `run-api-demo.sh` (defaults to a random free port). | Optional—use when firewall rules or demos require a fixed port. |
| `INTEGRATION_PRESERVE_POSTGRES=true` | Skip the automatic `docker compose rm -sf driftlock-postgres` cleanup after the integration script finishes. | Optional for long-running demos. |

### Legacy CLI HTML Demo (Still Available)

The CLI report is still shipped for backwards compatibility and offline screenshots. Treat it as the legacy/secondary path:

```bash
make demo                              # builds Rust core + Go CLI
./driftlock-demo test-data/financial-demo.json
open demo-output.html                  # macOS (use xdg-open on Linux)
```

`./verify-yc-ready.sh` exercises the same CLI pipeline end-to-end and remains part of CI, but all new onboarding now goes through the HTTP API.

Legacy HTML screenshots such as `screenshots/demo-anomaly-card.png` reflect this CLI report only; the canonical confidence and anomaly metrics are the API/terminal outputs shown in the demo GIF above.

If the license key is missing or expired, the server exits on startup and `/healthz` reports the invalid status.

### Admin CLI

`driftlock-http` now exposes management commands for migrations and tenant onboarding. Run them with the same binary used for the HTTP service:

```bash
# Apply pending goose migrations
./driftlock-http migrate up

# Inspect migration state
./driftlock-http migrate status

# Create a tenant, default stream, and API key
./driftlock-http create-tenant \
  --name "Bank Alpha" \
  --key-role admin

# List and revoke API keys
./driftlock-http list-keys --tenant bank-alpha
./driftlock-http revoke-key --id 8f97f374-...
```

All commands require `DATABASE_URL` (and `DRIFTLOCK_LICENSE_KEY` when running the server). Keys are printed once at creation; store them securely.

Pass `--json` to `create-tenant` when you need machine-readable output for scripts:

```bash
./driftlock-http create-tenant \
  --name "CI Tenant" \
  --key-role stream \
  --json | jq
```

The GitHub Actions `CI` workflow now runs `scripts/test-docker-build.sh` to guarantee the Dockerfiles stay in sync with `cbad-core`. OpenZL-enhanced images remain opt-in; set `USE_OPENZL=true` (and provide the private `openzl/` artifacts) to enable the feature flag.

### Integration Tests (API + Postgres)

- `./scripts/run-api-demo.sh` — human-friendly onboarding script described above. Prints the generated tenant/key, curl snippets, and Postgres queries.
- `./scripts/test-integration.sh` — CI-grade helper invoked by the wrapper. Provision Postgres via Docker Compose, run goose migrations, create a tenant/key, and exercise `/v1/detect`, `/v1/anomalies/{id}`, plus the export stubs. Accepts `INTEGRATION_*` env vars for automation.

Prerequisites (for both scripts):

1. `cargo build --release` (installs `cbad-core` artifacts under `cbad-core/target/release`)
2. Either:
   - `export DRIFTLOCK_LICENSE_KEY="<signed key from Shannon Labs>"` (production/pilot)
   - OR `export DRIFTLOCK_DEV_MODE=true` (development - bypasses license validation)
3. `export LD_LIBRARY_PATH=cbad-core/target/release:$LD_LIBRARY_PATH`
4. Docker, `jq`, `psql`, `curl`, `base64`

Run it manually when you need raw logs:

```bash
USE_OPENZL=false ./scripts/test-integration.sh
```

`INTEGRATION_PRESERVE_POSTGRES=true` keeps the Compose container running for deeper analysis; otherwise the container is removed on exit. `/healthz` now reports license status, queue backend, and database reachability for compliance logging.

## Multi-Million Dollar Fines Start January 2025

Regulators worldwide are auditing bank mathematical systems for explainability compliance. Black-box fraud detection = automatic failure.

- **EU DORA**: Requires explainable algorithms for all automated fraud detection (up to 2% annual turnover in fines)
- **US Regulations**: CFPB Fair Lending, NYDFS Part 500, and FFIEC guidelines mandate explainability
- Black-box models = automatic audit failure + massive fines
- Driftlock provides compression-based analysis with human-readable explanations

## How It Works

Driftlock analyzes payment gateway telemetry using compression-based anomaly detection:

1. **Builds a baseline** from normal transactions (first ~400 events in this demo)
2. **Detects anomalies** by comparing compression distance (NCD) of new vs. baseline data
3. **Generates explanations** with NCD scores, p-values, and statistical significance
4. **Outputs HTML report** with flagged anomalies and regulator-friendly math

The HTML includes a baseline comparison panel and similar normal examples for each anomaly to make the "why" obvious to auditors and compliance teams.

## What Banks Get

1. **Drop-in replacement** for black-box fraud detection
2. **Mathematical proof** for every flagged transaction
3. **Audit-ready reports** in 1 click
4. **Works with existing data** - no infrastructure changes

## Developer-First Anomaly API & Pricing Direction

Driftlock is designed to be a **developer-first anomaly detection API**: drop a deterministic `/v1/detect` endpoint into payment gateways, risk engines, or AI monitoring pipelines and get glass-box evidence back.

Formal pricing will be published with pilot feedback, but per [docs/SALES-MODELING.md](docs/SALES-MODELING.md) the direction is:

- **Generous free tier** so individuals and CI pipelines can experiment without a sales call.
- **Usage-based API** targeting roughly **$1 per million anomaly checks** with automatic volume discounts as event counts grow.
- **Data-based options** priced **meaningfully below legacy observability per-GB rates**, because Driftlock delivers math-backed anomaly signals instead of storing every log forever.
- **Enterprise compliance plans** that start in the **low thousands per month** so regulated teams stay under the combined cost of fines, legacy tooling, and manual audit prep.

All numbers remain illustrative, not a binding quote. The ROI calculator on the landing page uses the same assumptions so YC reviewers, engineers, and finance partners can reason about value while we continue to tune costs with pilot data. The goal stays constant: beat black-box tooling on price-per-proof while keeping enough margin to support mission-critical deployments.

## What This Demo Proves

- ✅ Core algorithm works (compression distance + statistical testing)
- ✅ Generates explainable, audit-friendly output
- ✅ Processes 2,000 transactions in 4-6 seconds
- ✅ Glass-box explanations (compression ratios, p-values, statistical significance)

**Next step**: Pilot integration with bank payment gateways to validate real-world performance.

## Current Status

- ✅ Working prototype detects fraud in synthetic payment data
- ✅ 0% false negatives in demo (tuned for ~1.5% detection rate)
- 🎯 **Next**: 3 pilot banks for Q1 2025

**Target customers**: Regulated banks and PSPs worldwide replacing black-box model-based anomaly detectors (EU DORA, US CFPB/NYDFS, global compliance).

**We need**: Introductions to bank compliance teams facing regulatory audits in EU, US, and other jurisdictions.

## Development

Built with modern tooling; transparent development process documented in [docs/ai-agents/](docs/ai-agents/).

Helpful commands:

- `make demo` — build Rust core + Go CLI demo.
- `make docker-http` — build the containerized HTTP engine locally.
- `make docker-test` — run Docker smoke tests (generic compressors by default; set `ENABLE_OPENZL_BUILD=true` to cover the optional OpenZL images when libraries are present).
- `./scripts/run-api-demo.sh` — guided HTTP API onboarding (build + Postgres + /v1/detect + anomaly detail + psql checks).
- `USE_OPENZL=false LD_LIBRARY_PATH=cbad-core/target/release go test ./...` (from `collector-processor/`) — run Go unit + FFI integration tests after rebuilding the Rust core with generic compressors only:

```bash
cd cbad-core
USE_OPENZL=false cargo clean
USE_OPENZL=false cargo build --release
cd ..
USE_OPENZL=false LD_LIBRARY_PATH=cbad-core/target/release go test ./collector-processor/...
```

## 📊 Demo Data

The demo uses `test-data/financial-demo.json` containing 5,000 synthetic payment transactions with:

- **Normal pattern**: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- **Anomalies**: Latency spikes up to 8000ms and a handful of malformed endpoints
- **Detection**: Demo tuned to flag ~30 anomalies from 2,000 processed events (~1.5% detection rate).

## 📚 Learn More

- **[DEMO.md](DEMO.md)** - 2-minute partner walkthrough with screenshots
- **[docs/API-DEMO-WALKTHROUGH.md](docs/API-DEMO-WALKTHROUGH.md)** - Manual API commands mirroring the onboarding script
- **[docs/](docs/)** - Full documentation and agent automation history

Visual proof (optional):

- Run: `./scripts/capture-anomaly-card.sh` (macOS Safari) to auto‑capture the first anomaly card into `screenshots/demo-anomaly-card.png`. If it fails due to permissions, follow `docs/CAPTURE-ANOMALY-SCREENSHOT.md` for manual capture.

---

*Developed by Shannon Labs. Licensed under Apache 2.0.*
