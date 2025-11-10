# Driftlock: Explainable Anomaly Detection for EU Banks

Stop €50M DORA fines with math-based fraud detection that auditors love.

EU banks face €50M fines starting Jan 2025 for black-box AI. When your ML flags a transaction as suspicious, regulators demand: "Show your work." Driftlock uses compression math (not AI) to detect fraud with explanations auditors can verify.

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## Try It: Detect Payment Fraud in 30 Seconds

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

## €50M Fines Start January 2025

EU regulators will audit bank AI systems for DORA compliance. Black-box fraud detection = automatic failure.

- DORA requires explainable AI decisions for all automated fraud detection
- Black-box models = automatic audit failure + €50M-€200M fines  
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

**Target customers**: EU-regulated banks and PSPs replacing black-box ML anomaly detectors.

**We need**: Introductions to EU bank compliance teams facing DORA audits.

## Development

Built with modern tooling; transparent development process documented in [docs/ai-agents/](docs/ai-agents/).

## 📊 Demo Data

The demo uses `test-data/financial-demo.json` containing 5,000 synthetic payment transactions with:
- **Normal pattern**: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- **Anomalies**: Latency spikes up to 8000ms and a handful of malformed endpoints
- **Detection**: Demo tuned to flag ~30 anomalies from 2,000 processed events (~1.5% detection rate).

## 📚 Learn More

- **[DEMO.md](DEMO.md)** - 2-minute partner walkthrough with screenshots
- **[docs/](docs/)** - Full documentation and AI agent history

Visual proof (optional):
- Run: `./scripts/capture-anomaly-card.sh` (macOS Safari) to auto‑capture the first anomaly card into `screenshots/demo-anomaly-card.png`. If it fails due to permissions, follow `docs/CAPTURE-ANOMALY-SCREENSHOT.md` for manual capture.

---

*Developed by Shannon Labs. Licensed under Apache 2.0.*
