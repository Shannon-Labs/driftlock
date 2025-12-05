# AI Testing Prompt for Driftlock Production Readiness

Use this prompt with an AI assistant that has Docker access to comprehensively test the Driftlock production deployment.

## Context

You are testing a production-ready Docker deployment of Driftlock, an anomaly detection API service. The project includes:
- HTTP API server (`driftlock-http`)
- Kafka collector processor (`driftlock-collector`) - optional
- Kafka + Zookeeper infrastructure - optional
- Web playground UI - optional

## Your Task

Comprehensively test the entire dockerized Driftlock package to ensure it's production-ready and demo-ready. Follow this systematic testing plan:

## Phase 1: Docker Build Testing

1. **Test all Docker builds:**
   ```bash
   cd /Volumes/VIXinSSD/driftlock
   ./scripts/test-docker-build.sh
   ```
   
   Expected: All 4 builds should pass (http, http-openzl, collector, collector-openzl)
   - Verify builds complete without errors
   - Check image sizes are reasonable (< 500MB)
   - Note any build warnings or issues

2. **Manual build verification:**
   ```bash
   # Test HTTP API build
   docker build -t driftlock-http:test --build-arg USE_OPENZL=false -f collector-processor/cmd/driftlock-http/Dockerfile .
   
   # Test collector build
   docker build -t driftlock-collector:test --build-arg USE_OPENZL=false -f collector-processor/cmd/driftlock-collector/Dockerfile .
   ```

## Phase 2: Service Startup Testing

1. **Start services:**
   ```bash
   docker compose up -d driftlock-http
   sleep 5
   ```

2. **Run service tests:**
   ```bash
   ./scripts/test-services.sh
   ```

3. **Manual health check:**
   ```bash
   curl -s http://localhost:8080/healthz | jq '.'
   ```
   
   Expected response structure:
   ```json
   {
    "success": true,
    "request_id": "...",
    "library_status": "healthy",
    "version": "1.0.0",
    "available_algos": ["zlab", "zstd", "lz4", "gzip"]
  }
  ```

4. **Verify service logs:**
   ```bash
   docker compose logs driftlock-http --tail 50
   ```
   - Check for startup errors
   - Verify Prometheus metrics registration
   - Check for any warnings

## Phase 3: API Endpoint Testing

1. **Run comprehensive API tests:**
   ```bash
   ./scripts/test-api.sh
   ```

2. **Manual endpoint verification:**

   **Health endpoint:**
   ```bash
   curl -I http://localhost:8080/healthz
   ```

   **Metrics endpoint:**
   ```bash
   curl -s http://localhost:8080/metrics | grep driftlock_http
   ```
   - Verify `driftlock_http_requests_total` exists
   - Verify `driftlock_http_request_duration_seconds` exists

   **Detection endpoint (NDJSON):**
   ```bash
   echo '{"test":"data1"}
   {"test":"data2"}' | curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @- | jq '.'
   ```

   **Detection endpoint (JSON array):**
   ```bash
   echo '[{"test":"data1"},{"test":"data2"}]' | curl -X POST "http://localhost:8080/v1/detect?format=json" \
     -H "Content-Type: application/json" \
     --data-binary @- | jq '.'
   ```

3. **Test error handling:**
   ```bash
   # Invalid JSON
   curl -X POST http://localhost:8080/v1/detect \
     -H "Content-Type: application/json" \
     --data "invalid json"
   
   # Empty body
   curl -X POST http://localhost:8080/v1/detect \
     -H "Content-Type: application/json"
   ```

4. **Verify security headers:**
   ```bash
   curl -I http://localhost:8080/healthz | grep -i "x-frame-options\|x-xss-protection\|content-security-policy\|referrer-policy"
   ```
   
   Expected headers:
   - `X-Frame-Options: DENY`
   - `X-XSS-Protection: 1; mode=block`
   - `Content-Security-Policy: default-src 'self'`
   - `Referrer-Policy: strict-origin-when-cross-origin`

## Phase 4: Integration Testing with Real Data

1. **Verify test data exists:**
   ```bash
   ls -lh test-data/*.jsonl
   ```

2. **Run integration tests:**
   ```bash
   ./scripts/test-integration.sh
   ```

3. **Manual anomaly detection tests:**

   **Normal transactions (should find < 5 anomalies):**
   ```bash
   curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/normal-transactions.jsonl | jq '.anomaly_count'
   ```

   **Anomalous transactions (should find > 80 anomalies):**
   ```bash
   curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/anomalous-transactions.jsonl | jq '.anomaly_count'
   ```

   **Mixed transactions (should find 45-55 anomalies):**
   ```bash
   curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/mixed-transactions.jsonl | jq '.'
   ```

4. **Verify response structure:**
   ```bash
   curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/mixed-transactions.jsonl | jq 'keys'
   ```
   
   Expected keys: `success`, `request_id`, `total_events`, `anomaly_count`, `anomalies`, `processing_time`, `compression_algo`

5. **Test Prometheus metrics increment:**
   ```bash
   # Get initial count
   INITIAL=$(curl -s http://localhost:8080/metrics | grep "driftlock_http_requests_total" | awk '{print $2}' | head -1)
   
   # Make requests
   for i in {1..5}; do
     curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
       -H "Content-Type: application/json" \
       --data-binary @test-data/normal-transactions.jsonl > /dev/null
   done
   
   # Get new count
   NEW=$(curl -s http://localhost:8080/metrics | grep "driftlock_http_requests_total" | awk '{print $2}' | head -1)
   
   echo "Initial: $INITIAL, New: $NEW"
   # Should show increment
   ```

## Phase 5: Kafka Integration Testing (Optional)

1. **Start Kafka services:**
   ```bash
   docker compose --profile kafka up -d
   sleep 10
   ```

2. **Run Kafka tests:**
   ```bash
   ./scripts/test-kafka.sh
   ```

3. **Manual Kafka verification:**
   ```bash
   # Check Kafka broker
   docker compose exec kafka kafka-broker-api-versions --bootstrap-server localhost:9092
   
   # Check collector logs
   docker compose logs driftlock-collector | grep -i kafka
   ```

## Phase 6: Structured Logging Verification

1. **Check logs format:**
   ```bash
   docker compose logs driftlock-http --tail 20 | head -5
   ```

2. **Make a request and check logs:**
   ```bash
   curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/normal-transactions.jsonl > /dev/null
   
   docker compose logs driftlock-http --tail 10
   ```

3. **Verify JSON log format:**
   - Logs should be in JSON format
   - Should include `request_id`, `method`, `path`, `event` fields
   - Should have request_start and request_complete events

## Phase 7: Performance Testing

1. **Test response time:**
   ```bash
   time curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
     -H "Content-Type: application/json" \
     --data-binary @test-data/mixed-transactions.jsonl > /dev/null
   ```
   - Should complete in < 5 seconds for 1000 events

2. **Test concurrent requests:**
   ```bash
   for i in {1..10}; do
     curl -X POST "http://localhost:8080/v1/detect?format=ndjson" \
       -H "Content-Type: application/json" \
       --data-binary @test-data/normal-transactions.jsonl > /dev/null &
   done
   wait
   echo "Concurrent requests completed"
   ```

3. **Check memory usage:**
   ```bash
   docker stats driftlock-http --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
   ```

## Phase 8: Graceful Shutdown Testing

1. **Test graceful shutdown:**
   ```bash
   docker compose stop driftlock-http
   sleep 2
   docker compose logs driftlock-http --tail 10 | grep -i shutdown
   ```

2. **Verify shutdown message:**
   - Should see "Server received signal" message
   - Should see "Server gracefully shutdown" message
   - Should complete within 30 seconds

## Phase 9: Playground Testing (Optional)

1. **Check playground setup:**
   ```bash
   cd playground
   ls -la .env.example
   cat .env.example
   ```

2. **Test playground build:**
   ```bash
   npm install
   npm run build
   ```

3. **Note:** Playground requires manual browser testing, but verify build works

## Phase 10: Demo Script Testing

1. **Run demo script:**
   ```bash
   cd /Volumes/VIXinSSD/driftlock
   ./scripts/demo.sh
   ```

2. **Verify demo script:**
   - Should start services
   - Should verify health
   - Should run sample detection
   - Should provide next steps

## Expected Results Summary

### Must Pass (Critical):
- ✅ All Docker images build successfully
- ✅ HTTP API server starts and responds to health checks
- ✅ `/v1/detect` endpoint works with test data
- ✅ Anomaly detection produces correct results:
  - Normal data: < 5 anomalies
  - Anomalous data: > 80 anomalies
  - Mixed data: 45-55 anomalies
- ✅ Prometheus metrics are exposed and increment correctly
- ✅ Security headers are present
- ✅ Response structure is correct

### Should Pass (Important):
- ✅ Structured logging works correctly
- ✅ Error handling is robust
- ✅ Performance is acceptable (< 5s for 1000 events)
- ✅ Graceful shutdown works

### Nice to Have (Optional):
- ✅ Kafka integration works (if enabled)
- ✅ Playground builds successfully
- ✅ Demo script works flawlessly

## Reporting

After completing all tests, provide:

1. **Test Results Summary:**
   - Total tests run
   - Tests passed
   - Tests failed
   - Tests skipped

2. **Issues Found:**
   - List any failures or errors
   - Include error messages and logs
   - Note any warnings

3. **Performance Metrics:**
   - Build times
   - Startup times
   - Response times
   - Memory usage

4. **Recommendations:**
   - Any fixes needed
   - Any improvements suggested
   - Production readiness assessment

## Files to Check

- `scripts/test-*.sh` - All test scripts
- `docker-compose.yml` - Service configuration
- `collector-processor/cmd/driftlock-http/Dockerfile` - HTTP API Dockerfile
- `collector-processor/cmd/driftlock-collector/Dockerfile` - Collector Dockerfile
- `test-data/*.jsonl` - Test data files
- `playground/.env.example` - Playground configuration template

## Quick Start Command

For a quick comprehensive test run:

```bash
cd /Volumes/VIXinSSD/driftlock

# 1. Test builds
./scripts/test-docker-build.sh

# 2. Start services and test
docker compose up -d driftlock-http
sleep 5
./scripts/test-services.sh
./scripts/test-api.sh
./scripts/test-integration.sh

# 3. Test Kafka (optional)
docker compose --profile kafka up -d
sleep 10
./scripts/test-kafka.sh

# 4. Run demo
./scripts/demo.sh
```

## Notes

- All scripts are executable and ready to run
- Test data files should exist in `test-data/` directory
- Docker Compose should be version 2.0+
- Port 8080 should be available
- For Kafka tests, ports 9092 and 2181 should be available

Good luck! Test thoroughly and report all findings.

