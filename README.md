# Driftlock

**Regulator-proof anomaly detection for high-compliance environments.**

Driftlock detects data drift and anomalies in milliseconds using deterministic compression mathematics (NCD), then explains findings in plain English. 
*   **No training data required.**
*   **Zero configuration.**
*   **Provably deterministic.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Production--Ready-green.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/shannon-labs/driftlock)](https://goreportcard.com/report/github.com/shannon-labs/driftlock)

---

## 🚀 Features

*   **Mathematical Certainty:** Uses Normalized Compression Distance (NCD) to measure how "surprising" new data is. If the compression ratio shifts, the data has changed.
*   **Zero Training:** Builds a baseline from your first ~400 events. No labeled datasets needed.
*   **AI Explainability:** Generates human-readable explanations for anomalies using LLMs.
*   **Compliance Ready:** Perfect for DORA, NIS2, and API abuse detection where auditability is key.
*   **Full Platform:** Includes multi-tenant API, real-time dashboard, and usage-based billing.

## 🛠️ Quick Start

### Run the Deterministic Demo
Process synthetic events and see the detection engine in action:

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Build and run demo
make demo
./driftlock-demo test-data/financial-demo.json

# View results
open demo-output.html
```

### Local Development

1.  **Prerequisites:** Go 1.21+, Rust 1.75+, Node.js 20+, Docker.
2.  **Start Infrastructure:**
    ```bash
    docker compose up -d
    ```
3.  **Run API:**
    ```bash
    export DRIFTLOCK_DEV_MODE=true
    ./scripts/run-api-demo.sh
    ```
4.  **Run Dashboard:**
    ```bash
    cd landing-page
    npm install
    npm run dev
    ```

## 🏗️ Architecture

*   **Core Engine (`cbad-core`):** High-performance Rust library for NCD/entropy detection.
*   **API (`collector-processor`):** Go HTTP API handling ingestion, multi-tenancy, and billing.
*   **Dashboard (`landing-page`):** Vue 3 + Tailwind interface for real-time monitoring.
*   **Infrastructure:** Designed for Google Cloud Run (Serverless) + PostgreSQL + Firebase Auth.

## 📚 Documentation

*   [API Documentation](docs/architecture/api/openapi.yaml)
*   [Architecture Overview](docs/architecture/)
*   [Use Cases](landing-page/public/docs/use-cases/)

## 📦 Pricing & Tiers

| Tier | Events/Month | Features |
|------|--------------|----------|
| **Pulse** (Free) | 10,000 | Basic detection, 14-day retention |
| **Radar** ($15/mo) | 500,000 | Email alerts, 30-day retention |
| **Tensor** ($100/mo) | 5,000,000 | DORA/NIS2 evidence bundles |
| **Orbit** (Custom) | Unlimited | Dedicated support, SLA |

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) (if available) or check out the `AGENTS.md` for current development priorities.

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) for details.
Commercial licenses available for enterprise deployment.
