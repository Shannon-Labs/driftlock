# Driftlock Repository Status — Ready for Review

This file reflects the current, simplified demo that ships in this repository.

## What ships today

- Rust core (`cbad-core/`) with FFI for Go
- Go CLI demo (`cmd/demo/main.go`) that reads a static JSON file
- Synthetic demo data (`test-data/financial-demo.json`)
- One minimal CI workflow that builds, runs, and verifies the demo
- HTML output (`demo-output.html`) with anomaly cards, baseline comparisons, and similar normal examples

## Quick start

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html  # macOS (or xdg-open on Linux)
```

Or run the single verification script used in CI:

```bash
./verify-yc-ready.sh
```

## Current expected results

- Processes the first 2,000 events from a 5,000‑row dataset
- Warmup: first 400 events build the baseline
- Anomalies: 10–30 (typically ~30)
- Runtime: ~4–6s on a modern laptop; <30s in CI
- Detection rate in the report: ~0.6% (30/5000 total events)

## Notes on removed/archived components

Earlier iterations included a Docker/React dashboard (docker‑compose, web frontend, API, DB). That stack is not part of this repository anymore. The current demo is a self‑contained Rust+Go CLI that renders HTML — no Docker, DB, or external services.

If you find any remaining references to docker-compose or a web frontend in older docs, treat them as historical. The authoritative instructions are in `README.md`, `DEMO.md`, and `verify-yc-ready.sh`.

---

Prepared for partner review.
