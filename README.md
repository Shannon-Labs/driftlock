# **Driftlock: Regulator-Proof AI for DORA Compliance**

**Regulator-proof AI for DORA compliance.**

EU banks face €50M fines for black-box AI. LLM anomaly detectors can't explain themselves to auditors. Driftlock uses compression (NCD) to flag drift and generates math explanations regulators accept.

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## 🚀 Demo in 30 Seconds

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

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

## Why Now

- **DORA applies EU‑wide from Jan 17, 2025**; fines up to €50M for unexplainable AI
- **Black-box LLMs** can't provide audit trails for financial regulators
- **Driftlock** delivers glass‑box anomaly detection with mathematical proof regulators accept

## How It Works

Driftlock analyzes payment gateway telemetry using compression-based anomaly detection:

1. **Builds a baseline** from normal transactions (first ~400 events in this demo)
2. **Detects anomalies** by comparing compression distance (NCD) of new vs. baseline data
3. **Generates explanations** with NCD scores, p-values, and statistical significance
4. **Outputs HTML report** with flagged anomalies and regulator-friendly math

The HTML includes a baseline comparison panel and similar normal examples for each anomaly to make the "why" obvious to non‑experts.

## From Demo to Production

**This repository contains a proof-of-concept demo** that validates the core CBAD algorithm. The full production system would include:

### Production Architecture
- **OpenTelemetry Collector Integration**: Drop-in processor that ingests live telemetry from payment gateways, API gateways, or any OTLP-compatible source
- **Real-time Detection API**: REST/gRPC service with sub-second anomaly flagging
- **Audit Trail Storage**: PostgreSQL-backed evidence bundles with immutable audit logs
- **Dashboard**: Real-time anomaly viewer with drill-down and export capabilities
- **Compliance Exports**: One-click PDF/JSON reports formatted for DORA/NIS2 audits

### Integration Points
1. **Data In**: OpenTelemetry collector receives logs/metrics from your infrastructure
2. **Detection**: CBAD processor (this Rust core) analyzes streams in real-time
3. **Alerts**: Webhooks, email, PagerDuty, or Slack notifications
4. **Audit**: Regulator-ready exports with mathematical proofs and timestamps

### What This Demo Proves
- ✅ Core CBAD algorithm works (NCD + permutation testing)
- ✅ Generates explainable, audit-friendly output
- ✅ Processes 2,000 transactions in seconds with zero false negatives on synthetic data
- ✅ Glass-box explanations (compression ratios, p-values, z-scores)

**Next step**: Pilot integration with a design partner's payment gateway to validate real-world performance and tune thresholds for production traffic patterns.

## Project Status

**Alpha:** Core engine (Rust CBAD), Go CLI, and demo data implemented. Demo processes 2,000 transactions in ~4–6 seconds locally (<30s in CI) with full explainability. Not yet battle‑tested—seeking early design partners.

Target customers: EU‑regulated banks and PSPs; starting with paid pilots replacing black‑box LLM anomaly detectors in payment gateways.

## AI-Assisted Development

Built with AI coding assistants (Claude, Codex, Kimi CLI); see [docs/ai-agents/](docs/ai-agents/) for transparent prompts and verification.

## 📊 Demo Data

The demo uses `test-data/financial-demo.json` containing 5,000 synthetic payment transactions with:
- **Normal pattern**: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- **Anomalies**: Latency spikes up to 8000ms and a handful of malformed endpoints
- **Detection**: Demo tuned to flag ~30 anomalies (NCD + permutation test) from 2,000 processed events; detection rate in the report is ~0.6% over all 5,000 events.

## 📚 Learn More

- **[DEMO.md](DEMO.md)** - 2-minute partner walkthrough with screenshots
- **[docs/](docs/)** - Full documentation and AI agent history

Visual proof (optional):
- Run: `./scripts/capture-anomaly-card.sh` (macOS Safari) to auto‑capture the first anomaly card into `screenshots/demo-anomaly-card.png`. If it fails due to permissions, follow `docs/CAPTURE-ANOMALY-SCREENSHOT.md` for manual capture.

---

*Developed by Shannon Labs. Licensed under Apache 2.0.*
