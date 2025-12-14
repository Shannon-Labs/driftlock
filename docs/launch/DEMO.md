```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock
cargo build -p driftlock-api --release
docker compose up -d postgres
DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" ./target/release/driftlock-api
```

# Driftlock Demo Walkthrough

## The 2-Minute Partner Script (API-First)

### 1. The Problem

"Last year, EU banks paid €2.8B in algorithmic transparency fines. When your detection system flags a transaction as suspicious, GDPR Article 22 and Basel III require you to explain WHY in human terms. Black-box models can't. That's €50M-€200M per violation. Driftlock can."

### 2. The Demo (Live API + Postgres)

1. Build and start the Rust API server:
   ```bash
   cargo build -p driftlock-api --release
   docker compose up -d postgres
   DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock" \
     ./target/release/driftlock-api
   ```
2. Show `/healthz` returning status and version:
   ```bash
   curl http://localhost:8080/healthz
   ```
3. Test the demo detection endpoint (no auth needed):
   ```bash
   curl -X POST http://localhost:8080/v1/demo/detect \
     -H "Content-Type: application/json" \
     -d '{
       "events": [
         "2025-01-01T10:00:00Z Transaction $150 account A123",
         "2025-01-01T10:00:01Z Transaction $155 account A123",
         "2025-01-01T10:00:02Z Transaction $9999999 account UNKNOWN"
       ]
     }'
   ```
4. Open a terminal tab with `psql` to show persistence:
   ```bash
   psql "$DATABASE_URL" -c "SELECT id, stream_id, ncd, p_value FROM anomalies LIMIT 5"
   ```

### 3. What You'll Narrate

- **Multi-tenant**: Create tenants via Firebase auth, get API keys from dashboard.
- **Explainable outputs**: `/v1/detect` returns NCD, compression ratios, entropy deltas, p-values, and human-language `explanation` strings.
- **Persistence**: PostgreSQL query shows anomalies saved with IDs, metrics, and timestamps.
- **Exports**: `/v1/anomalies` endpoint for listing and filtering anomalies.

### 4. Magic Moment

Walk through the detection response. Highlight:

- Explanation text referencing compression deltas
- NCD, p-value, entropy change, confidence
- Raw event data and detection metrics

### 5. Close

"You just saw Driftlock detect anomalies and persist evidence in Postgres—all deterministic, no black-box ML. Export this audit trail, submit to regulators, avoid €50M fines."

### 6. Production Deployment Talking Points

- Drop-in with your existing telemetry
- Deterministic streaming detection (same math as demo)
- `/healthz` and `/metrics` for compliance logging and monitoring
- Replit deployment today, Cloud Run/Kubernetes on the roadmap

## Demo Detection Examples

### Basic Demo Request

```bash
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "normal log entry 1",
      "normal log entry 2",
      "ERROR: unusual event detected"
    ]
  }'
```

### With Config Overrides

```bash
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": ["event1", "event2", "anomalous"],
    "config_override": {
      "ncd_threshold": 0.25,
      "p_value_threshold": 0.05,
      "compressor": "zstd"
    }
  }'
```

### Authenticated Production Request

```bash
# Get API key from dashboard first
curl -X POST http://localhost:8080/v1/detect \
  -H "X-Api-Key: dlk_your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "stream_id": "production-logs",
    "events": ["..."]
  }'
```

## Metrics & Monitoring

```bash
# Check Prometheus metrics
curl http://localhost:8080/metrics

# Key metrics to show:
# - driftlock_http_requests_total
# - driftlock_events_processed_total
# - driftlock_anomalies_detected_total
```
