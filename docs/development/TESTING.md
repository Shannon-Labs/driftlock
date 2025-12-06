# Production Testing Guide

This guide covers comprehensive testing procedures for the Driftlock production deployment.

## Quick Start

Run all tests with the provided scripts:

```bash
# Test Docker builds
./scripts/test-docker-build.sh

# Test service startup
./scripts/test-services.sh

# Test API endpoints
./scripts/test-api.sh

# Full integration test
./scripts/test-integration.sh

# Test Kafka integration (optional)
./scripts/test-kafka.sh

# One-command demo setup
./scripts/demo.sh
```

## Test Scripts

### 1. Docker Build Testing (`test-docker-build.sh`)

Tests that all Docker images build successfully:

- `driftlock-http` (with and without OpenZL)
- `driftlock-collector` (with and without OpenZL)
- Verifies image sizes are reasonable (< 500MB)
- Checks for build errors

**Usage:**
```bash
./scripts/test-docker-build.sh
```

**Expected Output:**
- All builds complete successfully
- Image sizes reported
- Build logs saved to `/tmp/docker-build-*.log`

### 2. Service Testing (`test-services.sh`)

Tests Docker Compose service startup and health:

- Starts `driftlock-http` service
- Tests `/healthz` endpoint
- Verifies health response structure
- Tests `/metrics` endpoint
- Tests graceful shutdown

**Usage:**
```bash
./scripts/test-services.sh
```

**Prerequisites:**
- Docker daemon running
- Port 8080 available

### 3. API Testing (`test-api.sh`)

Tests all API endpoints with various inputs:

- Health endpoint (`/healthz`)
- Metrics endpoint (`/metrics`)
- Detection endpoint (`/v1/detect`)
- Error handling
- Security headers verification

**Usage:**
```bash
# With default API URL (http://localhost:8080)
./scripts/test-api.sh

# With custom API URL
API_URL=https://api.driftlock.net ./scripts/test-api.sh
```

**Test Cases:**
- NDJSON format detection
- JSON array format detection
- Auto-detection (no format parameter)
- Query parameters (baseline, window, hop, algo)
- Invalid JSON handling
- Empty body handling
- Security headers presence

### 4. Integration Testing (`test-integration.sh`)

Full integration test with real test data:

- Tests with `normal-transactions.jsonl` (< 5 anomalies expected)
- Tests with `anomalous-transactions.jsonl` (> 80 anomalies expected)
- Tests with `mixed-transactions.jsonl` (45-55 anomalies expected)
- Verifies Prometheus metrics increment
- Validates response structure

**Usage:**
```bash
./scripts/test-integration.sh
```

**Prerequisites:**
- API server running
- Test data files in `test-data/` directory

### 5. Kafka Testing (`test-kafka.sh`)

Tests Kafka integration (optional):

- Verifies Kafka broker connectivity
- Verifies Zookeeper connectivity
- Checks collector service status
- Validates Kafka publisher initialization

**Usage:**
```bash
./scripts/test-kafka.sh
```

**Prerequisites:**
- Kafka services started (`docker compose --profile kafka up -d`)
- Collector configured with Kafka enabled

### 6. Demo Script (`demo.sh`)

One-command demo setup:

- Starts all required services
- Verifies service health
- Runs sample detection
- Provides next steps

**Usage:**
```bash
./scripts/demo.sh
```

## Manual Testing Procedures

### Health Check Testing

```bash
# Basic health check
curl http://localhost:8080/healthz

# Verify response structure
curl -s http://localhost:8080/healthz | jq '.'
```

Expected response:
```json
{
  "success": true,
  "request_id": "...",
  "library_status": "healthy",
  "version": "1.0.0",
  "available_algos": ["zlab", "zstd", "lz4", "gzip"]
}
```

### Prometheus Metrics Testing

```bash
# Get metrics
curl http://localhost:8080/metrics

# Verify specific metrics
curl -s http://localhost:8080/metrics | grep driftlock_http_requests_total
curl -s http://localhost:8080/metrics | grep driftlock_http_request_duration_seconds
```

### Anomaly Detection Testing

```bash
# Test with NDJSON
curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl | jq '.'

# Test with JSON array
curl -X POST "http://localhost:8080/v1/detect?format=json" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/test-demo.json | jq '.'
```

### Security Headers Testing

```bash
# Check headers
curl -I http://localhost:8080/healthz

# Verify specific headers
curl -I http://localhost:8080/healthz | grep -i "x-frame-options"
curl -I http://localhost:8080/healthz | grep -i "x-xss-protection"
curl -I http://localhost:8080/healthz | grep -i "content-security-policy"
```

### Error Handling Testing

```bash
# Invalid JSON
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  --data "invalid json"

# Empty body
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json"
```

## Performance Testing

### Basic Performance

```bash
# Time a request
time curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
  -H "Content-Type: application/json" \
  --data-binary @test-data/mixed-transactions.jsonl
```

### Load Testing

```bash
# Concurrent requests
for i in {1..10}; do
  curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
    -H "Content-Type: application/json" \
    --data-binary @test-data/mixed-transactions.jsonl &
done
wait
```

## Success Criteria

### Critical (Must Pass)
- ✅ All Docker images build successfully
- ✅ HTTP API server starts and responds to health checks
- ✅ `/v1/detect` endpoint works with test data
- ✅ Anomaly detection produces correct results
- ✅ Prometheus metrics are exposed and working
- ✅ Security headers are present

### Important (Should Pass)
- ✅ Kafka integration works (if enabled)
- ✅ Structured logging works correctly
- ✅ Error handling is robust
- ✅ Performance is acceptable (< 5s for 1000 events)

### Optional (Nice to Have)
- ✅ Load testing passes
- ✅ All edge cases handled

## Troubleshooting

### Docker Build Failures

```bash
# Check build logs
cat /tmp/docker-build-*.log

# Clean up and rebuild
docker compose down
docker system prune -f
docker compose build --no-cache
```

### Service Not Starting

```bash
# Check logs
docker compose logs driftlock-http

# Check service status
docker compose ps

# Restart service
docker compose restart driftlock-http
```

### API Not Responding

```bash
# Check if service is running
docker compose ps driftlock-http

# Check health endpoint
curl -v http://localhost:8080/healthz

# Check port availability
lsof -i :8080
```

### Kafka Connection Issues

```bash
# Check Kafka logs
docker compose logs kafka

# Test Kafka connectivity
docker compose exec kafka kafka-broker-api-versions --bootstrap-server localhost:9092

# Check collector logs
docker compose logs driftlock-collector
```

## Continuous Integration

These test scripts can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Test Docker Builds
  run: ./scripts/test-docker-build.sh

- name: Test Services
  run: ./scripts/test-services.sh

- name: Test API
  run: ./scripts/test-api.sh
```

## Next Steps

After passing all tests:

1. Review test results
2. Fix any failures
3. Run demo script: `./scripts/demo.sh`
4. Follow [DEMO_GUIDE.md](./DEMO_GUIDE.md) for presentation

