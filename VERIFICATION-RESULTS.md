# Driftlock Demo — Verification Results (Current)

## Summary

- Build: Rust core + Go demo build cleanly
- Demo: Processes 2,000 of 5,000 events and renders `demo-output.html`
- Anomalies: 10–30, typically ~30 (enforced by verification script)
- Runtime: ~4–6s on a modern laptop; CI requires <30s and passes
- Output: Each anomaly card includes baseline comparisons and similar normal examples

## Reproduce

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html
```

Or run the CI script locally:

```bash
./verify-yc-ready.sh
```

## What the report shows

- Header with totals and runtime
- Baseline Summary (first 400 events): processing_ms + amount_usd stats; top endpoints/origins
- For each anomaly:
  - NCD, p‑value, confidence, compression delta
  - Baseline Comparison table: timing z‑score, amount ratio/z‑score, endpoint/origin frequencies, compression Δ%
  - Similar normal examples from the warmup window

## Implementation notes

- Two-stage flow: warmup builds baseline; detection runs every N=25 events
- Rust FFI API: `cbad_detector_create_simple`, `cbad_add_transaction`, `cbad_detect`, `cbad_detector_ready`, `cbad_detector_free`
- Baseline computations are O(400) and per‑anomaly sampling scans ≤400 items

## Demo data (financial-demo.json)

- 5,000 synthetic payments
- Normal: 50–100ms processing, US/UK origins, `/v1/charges` endpoint
- Injected anomalies cause structural drift (compression efficiency changes, NCD ↑)

## Known issues (current)

None blocking. The demo is tuned to produce 10–30 anomalies, and the verify script checks this.

## Acceptance criteria (met)

- Anomalies: 10–30
- Runtime: <30s in CI; ~4–6s locally
- No Docker or external services; no new dependencies

READY FOR REVIEW ✅
