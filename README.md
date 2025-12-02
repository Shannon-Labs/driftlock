# DRIFTLOCK

**Regulator-proof anomaly detection for high-compliance environments.**

Driftlock detects data drift and anomalies in milliseconds using deterministic compression mathematics (NCD), then optionally explains them in plain English using Gemini Flash.

No training data required. Zero configuration. Provably deterministic.

---

## CORE THESIS

Traditional anomaly detection relies on opaque ML models that hallucinate and require massive training sets. Driftlock is different:

1.  **MATHEMATICAL CERTAINTY**: We use **Normalized Compression Distance (NCD)** and Shannon Entropy deltas. If the compression ratio shifts, the data has changed. It is not a guess; it is a measurement.
2.  **ZERO TRAINING**: The system builds a baseline from your first ~400 events.
3.  **AI EXPLAINABILITY**: When the math flags an anomaly, the evidence is sent to Gemini Flash to generate a human-readable explanation for your dashboard or Slack alerts.

**Perfect for**: DORA compliance, FFIEC/NIS2 auditing, API abuse detection, and AI agent monitoring.

---

## QUICK DEMO (CLI)

The fastest way to verify the engine is to run the local deterministic demo. It processes 2,000 synthetic financial events and generates an HTML report.

```bash
# 1. Clone
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# 2. Run
make demo
./driftlock-demo test-data/financial-demo.json

# 3. View
open demo-output.html
```

**What you will see:**
- **Baseline Phase**: The first 400 events establish the compression dictionary.
- **Detection Phase**: The next 1,600 events are scored.
- **Anomalies**: ~10-30 events will be flagged with high NCD scores (e.g., `0.85+`), indicating they are mathematically distinct from the baseline.

---

## PRODUCT SURFACE

Driftlock ships as a complete SaaS platform.

### 1. The Engine (`cbad-core`)
A high-performance Rust library implementing the Compression-Based Anomaly Detection (CBAD) algorithms. It exposes a C-compatible FFI for integration into any system.

### 2. HTTP API (`driftlock-http`)
A multi-tenant Go service wrapping the engine.
- **Endpoint**: `POST /v1/detect`
- **Auth**: Firebase Auth (JWT) or API Key
- **Storage**: PostgreSQL (Cloud SQL) + TimeScale
- **Billing**: Built-in Stripe integration for usage-based metering

### 3. Dashboard (`landing-page`)
A Vue 3 + Tailwind interface for monitoring real-time streams, managing API keys, and viewing AI-generated explanations.

---

## DEPLOYMENT

The repository contains a complete "SaaS-in-a-box" deployment capability.

**Infrastructure Stack:**
- **Compute**: Google Cloud Run (Serverless)
- **Database**: Cloud SQL (Postgres)
- **Auth**: Firebase Auth
- **Secrets**: Google Secret Manager

**Deploy to Production:**
See `docs/deployment/CLOUDSQL_FIREBASE_SETUP_GUIDE.md` for the authoritative runbook.

```bash
# Quick infrastructure setup
./scripts/deployment/setup-gcp-cloudsql-firebase.sh
```

---

## DEVELOPMENT

**Local API Setup:**
```bash
export DRIFTLOCK_DEV_MODE=true
./scripts/run-api-demo.sh
```

**Directory Structure:**
- `cbad-core/`: Rust anomaly detection engine.
- `collector-processor/`: Go HTTP API and business logic.
- `landing-page/`: Frontend dashboard.
- `docs/`: Comprehensive architectural and deployment documentation.

---

## LICENSE

**Apache 2.0** for open-source components.
Commercial licenses available for enterprise deployment. See `LICENSE-COMMERCIAL.md`.