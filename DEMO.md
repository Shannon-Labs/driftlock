```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
git submodule update --init --recursive
cargo build --release
DRIFTLOCK_DEV_MODE=true ./scripts/run-api-demo.sh
```

# Driftlock Demo Walkthrough

## The 2-Minute Partner Script (API-First)

### 1. The Problem

"Last year, EU banks paid €2.8B in algorithmic transparency fines. When your detection system flags a transaction as suspicious, GDPR Article 22 and Basel III require you to explain WHY in human terms. Black-box models can't. That's €50M-€200M per violation. Driftlock can."

### 2. The Demo (live API + Postgres)

1. Run `./scripts/run-api-demo.sh`. It builds the Go API, starts dockerized Postgres, runs migrations, creates a tenant/key, and calls `/v1/detect`.
2. Show `/healthz` returning `license`, `database`, and `queue` status.
3. Use the printed API key to `curl /v1/detect` or `curl /v1/anomalies/{id}` live.
4. Open a terminal tab with `psql ... -c "SELECT id, stream_id, ncd, p_value FROM anomalies ..."` to prove persistence.

![Driftlock API demo – terminal session](screenshots/api-demo-demo.gif)

Optional: run the manual commands from [docs/API-DEMO-WALKTHROUGH.md](docs/API-DEMO-WALKTHROUGH.md) for screen recordings.

### 3. What You'll Narrate

- **Multi-tenant**: Tenant + API key is minted in front of the user.
- **Explainable outputs**: `/v1/detect` returns NCD, compression ratios, entropy deltas, p-values, and human-language `explanation` strings.
- **Persistence**: PSQL query shows anomalies saved with ids, metrics, and timestamps.
- **Exports**: `/v1/anomalies/export` responds `202 Accepted` (stub queue today, job workers in roadmap).

### 4. Magic Moment

Walk through a single `/v1/anomalies/{id}` response. Highlight:

- Explanation text referencing compression deltas
- NCD, p-value, entropy change, confidence
- Raw event snippet + baseline/window snapshots

### 5. Close

"You just saw Driftlock provision a tenant, detect anomalies, and persist evidence in Postgres—all deterministic, no black-box ML. Export this audit trail, submit to regulators, avoid €50M fines."

### 6. Production Deployment Talking Points

- Drop-in with your existing payment gateway telemetry
- Deterministic streaming detection (same math as demo)
- `/healthz` surfaces license + queue + DB for compliance logging
- Docker Compose today, Helm/Supabase/SaaS on the roadmap

## Legacy CLI HTML Demo (Backup Flow)

Still need the standalone HTML for screenshots?

```bash
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html
```

![Legacy HTML anomaly card](screenshots/demo-anomaly-card.png)

Use it for static visuals, but lead with the HTTP API during partner conversations. The anomaly metrics and confidence levels in partner demos should match the API/terminal outputs shown in the GIF above.
