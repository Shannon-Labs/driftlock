# Driftlock: Explainable Anomaly Detection Platform

**Real-time Streaming Telemetry + Glass-box Anomaly Detection**

Driftlock provides explainable, deterministic anomaly detection for regulated industries. Built for financial services, healthcare, and critical infrastructure that need to explain every algorithmic decision to auditors.

## Two Deployment Options

### 1. 🖥️  CLI Demo & API Service (Production Ready)
The working demo described in FINAL-STATUS.md - perfect for pilots and partners.

**Try the Demo:**
```bash
# Build and run the CLI demo
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html

# Or run the verification script
./verify-yc-ready.sh
```

**Deploy the HTTP API:**
```bash
# Docker Compose (local development)
docker-compose up

# Cloud deployment  
./deploy.sh  # Firebase + Cloud Run
```

### 2. 🚀 SaaS Platform (In Development)
Modern web application with dashboard, real-time streaming, and API management.

**Preview:**
- Landing Page: `landing-page/` (Vue 3 + TypeScript)
- Dashboard: `/dashboard` (real-time anomaly monitoring) 
- API Docs: `/docs` (interactive API explorer)

## The Innovation: Compression-Based Anomaly Detection (CBAD)

Unlike black-box ML models, Driftlock uses **mathematical compression theory** to detect anomalies:

1. **Baseline Learning**: Compress normal data to learn patterns
2. **Anomaly Detection**: New data that compresses poorly = anomaly
3. **Glass-box Explanations**: NCD scores, p-values, compression ratios

**Why Math Works Better:**
- ✅ **Deterministic**: Same input = same output, always
- ✅ **Explainable**: Show mathematical proof to auditors  
- ✅ **No Training**: No ML models to train or retrain
- ✅ **Regulatory Ready**: Built-in compliance for DORA, NIS2, AI Act

## Quick Start Options

**For Developers:**
```bash
# Clone and run demo
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
```

**For Enterprises:**
```bash
# Deploy to your infrastructure  
./deploy.sh
# Configure your OTLP endpoints to point to Driftlock
```

**For SaaS Preview:**
```bash
# Run landing page locally
cd landing-page
npm install && npm run dev
# Visit localhost:5173
```

## Documentation

- [Demo Guide](FINAL-STATUS.md) - Working CLI demo
- [Roadmap](docs/ROADMAP_TO_LAUNCH.md) - Complete path to production
- [Architecture](docs/ARCHITECTURE.md) - Technical deep-dive
- [API Reference](docs/api/openapi.yaml) - OpenAPI specification
- [Deployment](docs/deployment/) - Cloud deployment guides

## Golden Invariants

These never change (per AGENTS.md):
- ✅ CLI demo remains working (`make demo`)
- ✅ Verification script passes (`./verify-yc-ready.sh`)  
- ✅ Deterministic outputs (same seed = same results)
- ✅ Glass-box explanations for every anomaly

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   OTLP Events   │───▶│  Driftlock API  │───▶│   Dashboard     │
│ (logs/metrics/  │    │  (Go + Rust)    │    │  (Vue 3 + TS)   │
│    traces)      │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │  (Anomalies +   │
                       │   Evidence)     │
                       └─────────────────┘
```

## License

- **Driftlock Core** (Rust): Apache 2.0
- **API Service & Dashboard**: Source-available (see LICENSE-COMMERCIAL.md)

---

**Next Steps:** See [ROADMAP_TO_LAUNCH.md](docs/ROADMAP_TO_LAUNCH.md) for the complete path from demo to $10M ARR.

