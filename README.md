# Driftlock

**Regulator-proof anomaly detection for high-compliance environments.**

Driftlock detects data drift and anomalies in milliseconds using deterministic compression mathematics (NCD), then explains findings in plain English. 
*   **No training data required.**
*   **Zero configuration.**
*   **Provably deterministic.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Production--Ready-green.svg)]()

---

## 🚀 Features

*   **Mathematical Certainty:** Uses Normalized Compression Distance (NCD) to measure how "surprising" new data is. If the compression ratio shifts, the data has changed.
*   **Zero Training:** Builds a baseline from your first ~400 events. No labeled datasets needed.
*   **AI Explainability:** Generates human-readable explanations for anomalies using LLMs.
*   **Compliance Ready:** Perfect for DORA, NIS2, and API abuse detection where auditability is key.
*   **Full Platform:** Includes multi-tenant API, real-time dashboard, and usage-based billing.

## 🛠️ Quick Start

### Run the Rust API locally
Bring up the API server with PostgreSQL in a few commands:

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Start Postgres (Docker)
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15

# Run the API (requires Rust 1.75+)
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
  cargo run -p driftlock-api

# Smoke test
curl http://localhost:8080/healthz
```

### Local Development

1. **Prerequisites:** Rust 1.75+, Node.js 18+, Docker (optional), PostgreSQL 15+.
2. **Configure env:** `cp .env.example .env` and set `DATABASE_URL`, `FIREBASE_PROJECT_ID`, Stripe keys (if billing).
3. **Build:** `cargo build --workspace`
4. **Run API:** `cargo run -p driftlock-api`
5. **Run dashboard (optional):**
    ```bash
    cd landing-page
    npm install
    npm run dev
    ```

## 🏗️ Architecture

*   **Core Engine (`cbad-core`):** High-performance Rust library for NCD/entropy detection.
*   **API (`driftlock-api`):** Rust Axum server for ingestion, billing, onboarding, and management.
*   **Database Layer (`driftlock-db`):** sqlx-powered repository layer shared across services.
*   **Dashboard (`landing-page`):** Vue 3 + Tailwind interface for onboarding and analysis.
*   **Infrastructure:** Runs on Cloud Run / Docker with PostgreSQL and Firebase Auth.

## 📚 Documentation

*   [User Guide](docs/user-guide/)
*   [Architecture & API](docs/architecture/)
*   [Development](docs/development/)
*   [Deployment](docs/deployment/)

## Pricing & Tiers

| Tier | Price | Events/Month | Streams | Features |
|------|-------|--------------|---------|----------|
| **Free** | $0 | 10,000 | 5 | Basic detection, 14-day retention |
| **Pro** | $99/mo | 500,000 | 20 | Email alerts, 90-day retention |
| **Team** | $199/mo | 5,000,000 | 100 | DORA/NIS2 evidence bundles, 1-year retention |
| **Enterprise** | Custom | Unlimited | 500+ | EU data residency, self-hosting, SLA |

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) (if available) or check out the `AGENTS.md` for current development priorities.

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) for details.
Commercial licenses available for enterprise deployment.
