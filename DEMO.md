```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
```

# Driftlock Demo Walkthrough

## The 2-Minute Partner Script

### 1. The Problem

"Last year, EU banks paid €2.8B in algorithmic transparency fines. When your detection system flags a transaction as suspicious, GDPR Article 22 and Basel III require you to explain WHY in human terms. Black-box models can't. That's €50M-€200M per violation. Driftlock can."

### 2. The Demo

No Docker needed. Build and run the local demo:

```bash
make demo               # builds Rust core + Go demo
./driftlock-demo test-data/financial-demo.json

# Open the HTML report
open demo-output.html   # macOS
# xdg-open demo-output.html  # Linux
```

### 3. What You'll See

**What Compliance Officers See:**
- **Risk level**: High/Medium/Low with clear reasoning
- **Regulatory explanation**: Plain English why this transaction is suspicious
- **Audit trail**: Timestamped decision log with mathematical proof
- **Similar cases**: 2-3 legitimate transactions for comparison
- **Export ready**: One-click PDF for regulator submission

### 4. The Magic Moment

Click any flagged transaction → see the **regulatory explanation** with:
- **Plain English reasoning**: "This payment is 47x larger than customer's typical transactions"
- **Audit trail ID**: Unique identifier for regulator review
- **Mathematical proof**: Compression-based analysis (industry standard method)
- **Export button**: Generate compliance report in 3 seconds

### 5. The Close

"Export this audit trail, submit to regulators, avoid €50M fines. That's explainable fraud detection that actually works. No black boxes. No regulatory risk. Just compliance."

### 6. Production Deployment

"This demo proves our algorithm works. In production, Driftlock:
- **Integrates without code changes** (protects your existing payment systems)
- **Detects fraud in real-time** (prevents losses before they happen)
- **Alerts your team immediately** (no compliance surprises)
- **Stores immutable audit trails** (regulator-ready documentation)
- **Generates compliance reports** (saves weeks of manual work)

Ready to pilot on your real payment data?"

## Screenshot Reference

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## What To Look For

- The HTML shows total processed transactions and anomalies detected.
- Each anomaly has a red badge and an explanation panel with key metrics.
- Use this page as the talking point; no extra setup required.

## Next Steps

1. **Pilot**: Run on your historical payment data
2. **Validate**: Compare our detection vs. your current system
3. **Deploy**: Integrate with your payment gateway

## Need the API Instead of the CLI?

Spin up the HTTP engine (same CBAD core) in Docker:

```bash
# Build + run the HTTP API locally
docker compose up --build driftlock-http

# Hit the detector
curl -X POST http://localhost:8080/v1/detect \
  -H 'Content-Type: application/json' \
  -d @test-data/financial-demo.json
```

`/healthz` reports the available compressors (zstd/lz4/gzip are always present; the optional OpenZL adapter appears when the proprietary library is mounted under `openzl/` or `OPENZL_LIB_DIR`).
