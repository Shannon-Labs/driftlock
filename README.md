# **Driftlock: Regulator-Proof AI for DORA Compliance**

**Regulator-proof AI for DORA compliance.**

EU banks face €50M fines for black-box AI. LLM anomaly detectors can't explain themselves to auditors. Driftlock uses compression (NCD) to flag drift and generates math explanations regulators accept.

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## 🚀 Demo in 30 Seconds

```bash
git clone https://github.com/shannon-labs/driftlock.git
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

- **€50M fines** for EU banks using unexplainable AI under DORA regulations
- **Black-box LLMs** can't provide audit trails for financial regulators  
- **Driftlock** delivers glass-box anomaly detection with mathematical proof regulators accept

## How It Works

Driftlock analyzes payment gateway telemetry using compression-based anomaly detection:

1. **Builds a baseline** from normal transactions (first ~400 events in this demo)
2. **Detects anomalies** by comparing compression distance (NCD) of new vs. baseline data
3. **Generates explanations** with NCD scores, p-values, and statistical significance
4. **Outputs HTML report** with flagged anomalies and regulator-friendly math

## Project Status

**Alpha:** Core engine (Rust CBAD), Go CLI, and demo data implemented. Demo processes 2,000 transactions in ~5 seconds with full explainability. Not yet battle-tested—seeking early design partners.

## AI-Assisted Development

Built with AI coding assistants (Claude, Codex, Kimi CLI); see [docs/ai-agents/](docs/ai-agents/) for transparent prompts and verification.

## 📊 Demo Data

The demo uses `test-data/financial-demo.json` containing 5,000 synthetic payment transactions with:
- **Normal pattern**: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- **Anomalies**: Latency spikes up to 8000ms and a handful of malformed endpoints
- **Detection**: Demo tuned to flag ~30 anomalies (NCD + permutation test) from 2,000 processed events. The sample anomaly shown is a US transaction.

## 📚 Learn More

- **[DEMO.md](DEMO.md)** - 2-minute partner walkthrough with screenshots
- **[docs/](docs/)** - Full documentation and AI agent history

Visual proof (optional):
- Run: `./scripts/capture-anomaly-card.sh` (macOS Safari) to auto-capture the first anomaly card into `screenshots/demo-anomaly-card.png`. If it fails due to permissions, follow `docs/CAPTURE-ANOMALY-SCREENSHOT.md` for manual capture.

---

*Developed by Shannon Labs. Licensed under Apache 2.0.*

![Demo Anomaly](screenshots/demo-anomaly-card.png)
