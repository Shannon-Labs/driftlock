# Driftlock

Deterministic, compression-based anomaly detection with a small, reproducible demo.

This repository ships a **proof-of-concept CLI demo** plus an **experimental HTTP API**.  
For an authoritative description of what is actually implemented, see `FINAL-STATUS.md`.

---

## What this repo contains

- **Rust core** (`cbad-core/`): compression-based anomaly detection library.
- **Go CLI demo** (`cmd/demo/`): reads synthetic payment data and produces an HTML report.
- **Synthetic data** (`test-data/financial-demo.json`): 5,000 payment-like events with injected anomalies.
- **Optional HTTP API prototype** (`collector-processor/cmd/driftlock-http`): JSON `/v1/detect` endpoint backed by Postgres.

The only path we rely on for verification and CI is the **CLI demo** described below.

---

## Quickstart: CLI HTML demo (single binary)

This is the simplest, fully-supported way to see Driftlock work end-to-end.

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html  # use xdg-open on Linux
```

What you should see:

- First ~400 events build the baseline.
- The next 1,600 events are scanned for anomalies.
- The HTML report highlights ~10–30 anomalies with compression-based metrics and explanations.

You can re-run the demo as many times as you like; outputs are deterministic for the same input.

---

## Experimental HTTP API prototype

There is an in-repo HTTP service that exposes the same core algorithm over JSON.  
It is useful for local experiments but **not a hardened product**.

High-level steps:

1. Build the Rust core (needed for all Go binaries):

   ```bash
   cd cbad-core
   cargo build --release
   cd ..
   ```

2. Start Postgres with Docker Compose and run the API demo script:

   ```bash
   export DRIFTLOCK_DEV_MODE=true  # dev-only, bypasses licensing
   ./scripts/run-api-demo.sh
   ```

The script:

- Builds the `driftlock-http` binary.
- Starts a local Postgres instance.
- Applies migrations and creates a demo tenant + API key.
- Calls `/v1/detect` with synthetic data and prints follow-up `curl` and `psql` commands.

For the manual, step-by-step version see `docs/API-DEMO-WALKTHROUGH.md`.  
For the HTTP API schema, see `docs/API.md`.

**OpenZL note:** OpenZL integration is optional and experimental. All demos and default builds use generic compressors (zstd, lz4, gzip). See `docs/OPENZL_ANALYSIS.md` if you want to opt in.

---

## How the demo works (conceptually)

At a high level, both the CLI and HTTP flows do the same thing:

1. Build a baseline from normal events.
2. Compare new events to that baseline using compression distance (NCD).
3. Run permutation testing to estimate p-values / confidence.
4. Emit anomalies with:
   - NCD
   - compression ratios
   - entropy change
   - p-value and confidence
   - a short explanation string.

The math and implementation details are documented in `docs/ALGORITHMS.md`.

---

## Project status

See `FINAL-STATUS.md` for the current repository status. As of that file’s last update:

- ✅ Rust + Go CLI demo is stable and exercised in CI via `./verify-yc-ready.sh`.
- ✅ Synthetic dataset and HTML report are suitable for screenshots and quick demos.
- ⚠️ HTTP API, multi-tenant flows, and pricing/ROI language are **prototype only** and may change.

If you are evaluating Driftlock for anything beyond local experiments, treat this repo as an engine prototype rather than a finished product.

---

## License

Apache 2.0 for the open-source portions of this repository.  
See `LICENSE` and `LICENSE-COMMERCIAL.md` for details.

