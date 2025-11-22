# Driftlock God Mode Quickstart

Fastest path to prove the math, exercise the HTTP API, and see the SaaS surface. This is a **hands-on checklist**—follow it in order and you’ll have the demo, API, and landing page running locally in under an hour.

## Prereqs (local)

- macOS/Linux with Docker running
- Git, bash/zsh
- Go 1.21+ and Rust stable (for the cbad core + Go FFI)
- Node 18+ and npm (for the landing page/dev server)
- `firebase-tools` if you plan to deploy (`npm i -g firebase-tools`)

## 1) Clone and orient

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
```

Review the launch invariants: `AGENTS.md`, `FIREBASE_SAAS_COMPLETE.md`, and `FINAL-STATUS.md` (core demo must stay deterministic).

## 2) Deterministic CLI demo (math proof)

```bash
make demo
./driftlock-demo test-data/financial-demo.json
open demo-output.html  # use xdg-open on Linux
```

What you’re proving: NCD, compression ratios, entropy deltas, p-values—all deterministic for the same input.

## 3) HTTP API + dev tenant

```bash
export DRIFTLOCK_DEV_MODE=true  # dev-only, bypasses licensing
./scripts/run-api-demo.sh
```

The script builds the Go binary, spins up Postgres, applies migrations, seeds a demo tenant + API key, and exercises `/v1/detect`. While it runs, open another shell and hit health + detect:

```bash
curl -s http://localhost:8080/healthz | jq .
curl -s -X POST http://localhost:8080/v1/detect \
  -H "X-Api-Key: <API_KEY_FROM_SCRIPT_OUTPUT>" \
  -H "Content-Type: application/json" \
  --data-binary @<(jq '{stream_id:"demo-stream", events:.}' test-data/financial-demo.json) | jq .
```

The script prints the `API key`, `stream_id`, and sample follow-up commands; if you set `INTEGRATION_SUMMARY_FILE`, it also writes them to disk.

If you prefer Docker Compose: `docker compose up driftlock-http`.

## 4) Landing page + dashboard

```bash
cd landing-page
npm install
npm run dev   # visit http://localhost:5173
```

Use the seeded demo API key to view dashboard states; Firebase Auth can be toggled on with your project config via `.env.local`.

## 5) Launch readiness sweep

Run the consolidated checks before recording or deploying:

```bash
./scripts/verify-launch-readiness.sh
```

See the script output for any TODOs (env files, firebase setup, migrations).

## 6) Deploy (Firebase Hosting + Cloud Run)

When you’re ready to push a pilot:

```bash
./scripts/setup-gcp-cloudsql-firebase.sh   # infra bootstrap
./scripts/deploy-production-cloudsql-firebase.sh
./scripts/test-deployment-complete.sh
firebase deploy   # hosts landing page + functions
```

Secrets live in Google Secret Manager (no hard-coded keys). AI is opt-in; default flows stay math-only for cost control.

## 7) Upgrade ideas (pressure-tested path)

- Turn on Gemini for explainability per anomaly (premium)
- Capture real-time anomaly feed via SSE `/v1/stream/anomalies`
- Wire Stripe price IDs to the dashboard upgrade CTA

You now have the core math proof, API surface, UI, and deployment path validated. Iterate from here. 
