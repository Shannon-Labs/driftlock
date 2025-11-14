# Driftlock: Explainable Anomaly Detection for Global Banks

Stop multi-million dollar fines with math-based fraud detection that auditors love.

Banks worldwide face massive fines for black-box algorithms. When your system flags a transaction as suspicious, regulators demand: "Show your work." Driftlock uses compression math (not ML) to detect fraud with explanations auditors can verify.

**EU DORA**: Up to 2% of annual turnover in fines | **US Regulations**: CFPB Fair Lending, NYDFS Part 500, FFIEC Guidelines

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## Try It: Detect Payment Fraud in 30 Seconds

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Initialize git submodules (OpenZL and its nested dependencies)
git submodule update --init --recursive

# Build CBAD core (Rust) and the demo (Go)
make demo

# Or build the demo directly (no Makefile)
go build -o driftlock-demo cmd/demo/main.go

# Run the demo
./driftlock-demo test-data/financial-demo.json

# Open the results
open demo-output.html  # macOS
# xdg-open demo-output.html  # Linux
```

## Run the HTTP API with Docker

The Go HTTP service (`cmd/driftlock-http`) is the canonical API for integrating pilot workloads. Build it locally or use Docker:

```bash
# Build the Rust core + Go binary into a container
make docker-http

# Or run the full compose stack (HTTP API only)
docker compose up --build driftlock-http

# Verify health and available compressors (OpenZL is optional and reported separately)
curl -s http://localhost:8080/healthz | jq
```

The GitHub Actions `CI` workflow now runs `scripts/test-docker-build.sh` to guarantee the Dockerfiles stay in sync with `cbad-core`. OpenZL-enhanced images remain opt-in; set `USE_OPENZL=true` (and provide the private `openzl/` artifacts) to enable the feature flag.

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
- `LD_LIBRARY_PATH=cbad-core/target/release go test ./...` (from `collector-processor/`) — run Go unit + FFI integration tests after `cargo build --release`.

## 📊 Demo Data

The demo uses `test-data/financial-demo.json` containing 5,000 synthetic payment transactions with:
- **Normal pattern**: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- **Anomalies**: Latency spikes up to 8000ms and a handful of malformed endpoints
- **Detection**: Demo tuned to flag ~30 anomalies from 2,000 processed events (~1.5% detection rate).

## 📚 Learn More

- **[DEMO.md](DEMO.md)** - 2-minute partner walkthrough with screenshots
- **[docs/](docs/)** - Full documentation and agent automation history

Visual proof (optional):
- Run: `./scripts/capture-anomaly-card.sh` (macOS Safari) to auto‑capture the first anomaly card into `screenshots/demo-anomaly-card.png`. If it fails due to permissions, follow `docs/CAPTURE-ANOMALY-SCREENSHOT.md` for manual capture.

---

*Developed by Shannon Labs. Licensed under Apache 2.0.*
