# DRIFTLOCK

**Regulator-proof anomaly detection for high-compliance environments.**

Driftlock detects data drift and anomalies in milliseconds using deterministic compression mathematics, then explains findings in plain English. No training data required. Zero configuration. Provably deterministic.

[![CI](https://github.com/Shannon-Labs/driftlock/actions/workflows/ci.yml/badge.svg)](https://github.com/Shannon-Labs/driftlock/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

---

## Status

**Launch hardening** — Core platform built; distribution, monitoring, and compliance items tracked in [docs/launch/ROADMAP.md](docs/launch/ROADMAP.md#-remaining-critical-path-target-ga-by-dec-5-2025).

| Component | Status |
|-----------|--------|
| Detection Engine (Rust) | Production-ready |
| HTTP API (Go) | Production-ready |
| Dashboard (Vue 3) | Ready; production cutover pending custom domain/SSL (see roadmap) |
| Stripe Billing | Integrated |
| Firebase Auth | Integrated |
| Cloud Run Deployment | Staged; monitoring/alerting and domain selection in progress (see roadmap) |
| VS Code Extension | Marketplace publish blocked on publisher cert (see roadmap) |
| CLI Releases | Release binaries WIP for macOS/Linux (see roadmap) |
| Analytics & Alerting | Pending Firebase Analytics + Cloud Monitoring + PagerDuty (see roadmap) |
| Security & Policy | Pentest and privacy policy updates scheduled (see roadmap) |

---

## The Problem

Traditional anomaly detection relies on opaque ML models that:
- Require massive labeled training sets
- Produce unexplainable results ("the model said so")
- Fail regulatory audits due to non-determinism

**Regulated industries need provable, auditable detection.**

---

## Our Approach

Driftlock uses **Normalized Compression Distance (NCD)** — information-theoretic mathematics that measures how "surprising" new data is relative to a baseline.

1. **Mathematical Certainty**: If the compression ratio shifts, the data has changed. Not a guess — a measurement.
2. **Zero Training**: The system builds a baseline from your first ~400 events. No labeled data required.
3. **AI Explainability**: When math flags an anomaly, we generate human-readable explanations for dashboards and alerts.

**Perfect for**: DORA compliance, NIS2 auditing, API abuse detection, AI agent monitoring.

---

## Quick Demo

```bash
# Clone and run the local deterministic demo
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html
```

Processes 2,000 synthetic events. Flags ~10-30 anomalies with NCD scores indicating mathematical distinctness from baseline.

---

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Your Data     │────▶│  Driftlock API  │────▶│   Dashboard     │
│  (OTLP/JSON)    │     │   (Go + Rust)   │     │    (Vue 3)      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │
                    ┌──────────┴──────────┐
                    ▼                     ▼
             ┌───────────┐         ┌───────────┐
             │  cbad-core │         │  Gemini   │
             │   (Rust)   │         │  (Explain) │
             └───────────┘         └───────────┘
```

### Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| `cbad-core` | Rust | High-performance NCD/entropy detection engine |
| `collector-processor` | Go | Multi-tenant HTTP API with billing |
| `landing-page` | Vue 3 + Tailwind | Real-time dashboard |

### Infrastructure

- **Compute**: Google Cloud Run (serverless, auto-scaling)
- **Database**: Cloud SQL (PostgreSQL)
- **Auth**: Firebase Authentication
- **Billing**: Stripe (usage-based metering)
- **Secrets**: Google Secret Manager

---

## Pricing

| Tier | Price | Events/Month | Features |
|------|-------|--------------|----------|
| Pulse | Free | 10,000 | Basic detection, 14-day retention |
| Radar | $15/mo | 500,000 | Email alerts, 30-day retention |
| Tensor | $100/mo | 5,000,000 | DORA/NIS2 evidence bundles |
| Orbit | $499/mo | Unlimited | Dedicated support, SLA |

---

## Development

```bash
# Local API
export DRIFTLOCK_DEV_MODE=true
./scripts/run-api-demo.sh

# Frontend
cd landing-page && bun install && bun run dev

# Tests
make test
```

**Directory Structure:**
- `cbad-core/` — Rust detection engine with C FFI
- `collector-processor/` — Go HTTP API and business logic
- `landing-page/` — Vue 3 dashboard
- `docs/` — Architecture and deployment documentation

---

## License

**Apache 2.0** for open-source components.

Commercial licenses available for enterprise deployment. Contact: hello@shannon-labs.com
