# Driftlock API Demo Walkthrough

This guide mirrors `./scripts/run-api-demo.sh` step-by-step so you can copy/paste the Docker + Postgres commands during a live call or in sandboxed environments.

## Prerequisites

- macOS/Linux shell with `git`, `docker`, `curl`, `jq`, `psql`, `go`, `cargo`, and `base64`
- `cbad-core` compiled (`cargo build --release`)
- Optional: signed `DRIFTLOCK_LICENSE_KEY`. Otherwise export `DRIFTLOCK_DEV_MODE=true` for local demos (dev-mode is not permitted in production).

## 1. Clone and prepare the repo

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
git submodule update --init --recursive
```

## 2. Build the core library

```bash
cargo build --release
export LD_LIBRARY_PATH="$(pwd)/cbad-core/target/release:${LD_LIBRARY_PATH:-}"
```

## 3. Start Postgres (Docker Compose)

```bash
docker compose up -d driftlock-postgres
```

The default connection string (overridable via `DATABASE_URL`) is:

```
postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable
```

## 4. Apply migrations

```bash
mkdir -p bin
(cd collector-processor/cmd/driftlock-http && go build -o ../../bin/driftlock-http .)
DATABASE_URL="postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable" \
  ./bin/driftlock-http migrate up
```

## 5. Create a tenant + API key

```bash
TENANT_JSON=$(DATABASE_URL=postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable \
  ./bin/driftlock-http create-tenant \
    --name "Demo Tenant" \
    --slug demo-tenant \
    --stream demo-stream \
    --key-role admin \
    --json)
API_KEY=$(echo "$TENANT_JSON" | jq -r '.api_key')
STREAM_ID=$(echo "$TENANT_JSON" | jq -r '.stream_id')
```

Store `API_KEY` securely—it is only printed once.

## 6. Launch the API server (dev mode bypass shown; provide a real license in production)

```bash
PORT=8080 \
DRIFTLOCK_DEV_MODE=true \
DATABASE_URL=postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable \
./bin/driftlock-http
```

Verify `/healthz`:

```bash
curl -s http://localhost:8080/healthz | jq
```

You should see `license.status` (`valid` or `dev_mode`), `database="connected"`, and `queue.status` describing the in-memory exporter.

## 7. Run `/v1/detect`

```bash
jq '.[0:600]' test-data/financial-demo.json \
  | curl -sS -X POST \
      -H "X-Api-Key: ${API_KEY}" \
      -H 'Content-Type: application/json' \
      --data @- \
      http://localhost:8080/v1/detect | jq '{anomaly_count, anomalies: [.anomalies[0]]}'
```

The response returns the same anomaly evidence as the legacy CLI demo, but now the batches/anomalies are persisted.

## 8. Inspect anomaly detail and Postgres rows

```bash
FIRST_ANOMALY=$(psql postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable -Atc \
  "SELECT id FROM anomalies ORDER BY detected_at DESC LIMIT 1")

curl -sS -H "X-Api-Key: ${API_KEY}" \
  "http://localhost:8080/v1/anomalies/${FIRST_ANOMALY}" | jq '{id, ncd, p_value, explanation}'
```

Query the persisted rows:

```bash
psql postgres://driftlock:driftlock@localhost:7543/driftlock?sslmode=disable \
  -c "SELECT id, stream_id, ROUND(ncd::numeric, 4) AS ncd, ROUND(p_value::numeric, 4) AS p_value FROM anomalies ORDER BY detected_at DESC LIMIT 5;"
```

## 9. Clean up

Press `Ctrl+C` to stop `driftlock-http`, then run:

```bash
docker compose rm -sf driftlock-postgres
```

The workflow above is identical to what `scripts/run-api-demo.sh` automates; keep this doc handy for compliance reviews or screen recordings where you need to show each command explicitly.
