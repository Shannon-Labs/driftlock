# Driftlock Demo Guide

This guide helps you prepare and execute a successful Driftlock demonstration.

## Pre-Demo Setup

### 1. Prerequisites Check

```bash
# Verify Docker is running
docker --version
docker compose version

# Verify test data exists
ls test-data/*.jsonl
```

### 2. Quick Setup

Run the demo script:

```bash
./scripts/demo.sh
```

This will:
- Start all required services
- Verify service health
- Run a sample detection
- Provide next steps

### 3. Manual Setup (Alternative)

```bash
# Start API server
docker compose up -d driftlock-http

# Wait for service to be ready
sleep 5

# Verify health
curl http://localhost:8080/healthz
```

## Demo Flow

### Step 1: Show Service Health

```bash
# Health check
curl http://localhost:8080/healthz | jq '.'
```

**Key Points:**
- Service is running and healthy
- Library validation passed
- Available algorithms listed

### Step 2: Demonstrate Anomaly Detection

#### Normal Data (Low Anomaly Rate)

```bash
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/normal-transactions.jsonl | jq '.anomaly_count'
```

**Expected:** < 5 anomalies

**Talking Points:**
- Normal transaction patterns compress well
- CBAD algorithm identifies expected patterns
- Low false positive rate

#### Anomalous Data (High Anomaly Rate)

```bash
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/anomalous-transactions.jsonl | jq '.anomaly_count'
```

**Expected:** > 80 anomalies

**Talking Points:**
- Suspicious transactions don't compress well
- High anomaly detection rate
- Algorithm identifies unusual patterns

#### Mixed Data (Balanced)

```bash
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq '.'
```

**Expected:** 45-55 anomalies (5% of 1000 events)

**Talking Points:**
- Real-world scenario with mixed data
- Accurate detection rate
- Detailed anomaly explanations

### Step 3: Show Detailed Anomaly Information

```bash
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq '.anomalies[0]'
```

**Key Points:**
- Glass-box explanations
- NCD (Normalized Compression Distance) values
- P-values for statistical significance
- Detailed "why" explanations

### Step 4: Demonstrate Prometheus Metrics

```bash
# Show metrics endpoint
curl http://localhost:8080/metrics | grep driftlock_http

# Make requests and show metrics increment
for i in {1..5}; do
  curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
    -H "Content-Type: application/json" \
    --data-binary @test-data/normal-transactions.jsonl > /dev/null
done

curl -s http://localhost:8080/metrics | grep driftlock_http_requests_total
```

**Key Points:**
- Production-ready observability
- Request counting
- Duration tracking
- Integration with Prometheus/Grafana

### Step 5: Show Playground UI (Optional)

```bash
# Start playground
cd playground
npm install
cp .env.example .env
npm run dev
```

**Open:** http://localhost:5174

**Features to Demonstrate:**
- Upload JSON/NDJSON files
- Adjust detection parameters
- View results in table format
- Download results as JSON
- API health indicator

### Step 6: Show Security Features

```bash
# Security headers
curl -I http://localhost:8080/healthz | grep -i "x-frame-options\|x-xss-protection\|content-security-policy"
```

**Key Points:**
- Enterprise-grade security headers
- XSS protection
- Frame options
- Content Security Policy

## Demo Scripts

### Quick Demo (5 minutes)

```bash
# 1. Start services
./scripts/demo.sh

# 2. Show health
curl http://localhost:8080/healthz | jq '.'

# 3. Run detection
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq '.anomaly_count'

# 4. Show metrics
curl -s http://localhost:8080/metrics | grep driftlock_http_requests_total
```

### Full Demo (15 minutes)

1. **Introduction** (2 min)
   - What is Driftlock
   - Use cases (fraud detection, anomaly detection)
   - Key differentiator: glass-box explanations

2. **Technical Overview** (3 min)
   - CBAD algorithm (Compression-Based Anomaly Detection)
   - How it works
   - Why compression distance matters

3. **Live Demo** (7 min)
   - Start services
   - Show health check
   - Run detection on normal data
   - Run detection on anomalous data
   - Run detection on mixed data
   - Show detailed anomaly explanations

4. **Production Features** (2 min)
   - Prometheus metrics
   - Security headers
   - Structured logging
   - Docker deployment

5. **Q&A** (1 min)

## Common Demo Scenarios

### Scenario 1: Financial Fraud Detection

**Setup:**
```bash
# Use transaction data
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq '.anomalies[] | select(.metrics.ncd > 0.5)'
```

**Talking Points:**
- Detects unusual transaction patterns
- Identifies high-value suspicious transactions
- Provides explanations for compliance

### Scenario 2: Log Anomaly Detection

**Setup:**
```bash
# Create sample log data
cat > /tmp/logs.jsonl << EOF
{"timestamp":"2025-01-01T00:00:00Z","level":"INFO","message":"User login successful"}
{"timestamp":"2025-01-01T00:00:01Z","level":"INFO","message":"User login successful"}
{"timestamp":"2025-01-01T00:00:02Z","level":"ERROR","message":"Unauthorized access attempt from suspicious IP"}
EOF

curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @/tmp/logs.jsonl | jq '.'
```

**Talking Points:**
- Detects unusual log patterns
- Identifies security events
- Helps with incident response

## Troubleshooting During Demo

### Service Not Responding

```bash
# Quick restart
docker compose restart driftlock-http

# Check logs
docker compose logs driftlock-http --tail 20
```

### Slow Response Times

- Reduce test data size
- Use smaller baseline/window parameters
- Check system resources

### Unexpected Results

- Verify test data format
- Check algorithm parameters
- Review anomaly thresholds

## Post-Demo

### Cleanup

```bash
# Stop services
docker compose down

# Remove test images (optional)
docker rmi driftlock-http:test driftlock-collector:test
```

### Follow-up Materials

- [API Documentation](./api/README.md)
- [Architecture Overview](./ARCHITECTURE.md)
- [Deployment Guide](./DEPLOYMENT.md)

## Tips for Success

1. **Prepare Test Data:** Have test files ready before demo
2. **Practice Flow:** Run through demo once before presentation
3. **Have Backup:** Keep terminal history for quick commands
4. **Explain Concepts:** Don't just show, explain the "why"
5. **Handle Questions:** Be ready to dive deeper into algorithm details

## Key Metrics to Highlight

- **Accuracy:** 95%+ recall on mixed data
- **Performance:** < 5 seconds for 1000 events
- **Explainability:** Glass-box explanations for every anomaly
- **Production Ready:** Prometheus metrics, security headers, structured logging

