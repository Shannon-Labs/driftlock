make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html

# Driftlock Demo Walkthrough

## The 2-Minute Partner Script

### 1. The Problem

"EU banks process millions of payments daily. When AI flags a transaction as anomalous, regulators demand: **'Show your work.'** Black-box LLMs can't. Driftlock can."

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

**HTML Report:**
- Synthetic payment data processed (≈2,000 transactions)
- Multiple flagged anomalies in the list with badges
- Explanation panel per anomaly with:
  - NCD score and p-value
  - Compression ratio/entropy deltas
  - A natural-language summary for auditors

### 4. The Magic Moment

Click any flagged anomaly → see the **explanation panel** with:
- Mathematical compression distance calculation
- Reference to historical patterns
- Audit trail ID for regulator export

### 5. The Close

"Export this audit log, hand it to regulators → fine avoided. Glass-box AI that ships."

## Screenshot Reference

![Demo Anomaly](screenshots/demo-anomaly-card.png)

## What To Look For

- The HTML shows total processed transactions and anomalies detected.
- Each anomaly has a red badge and an explanation panel with key metrics.
- Use this page as the talking point; no extra setup required.

## Next Steps

1. **Integration**: Connect your OpenTelemetry pipeline
2. **Customization**: Adjust thresholds for your data patterns
3. **Scale**: Stream logs/OTLP into the detector (collector component)
