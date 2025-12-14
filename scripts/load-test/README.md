# Driftlock K6 Load Testing Suite

Comprehensive stress testing suite for the Driftlock Rust API to validate production readiness.

## Prerequisites

1. **Install K6:**
   ```bash
   # macOS
   brew install k6

   # Linux
   sudo gpg -k
   sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
   echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
   sudo apt-get update
   sudo apt-get install k6

   # Windows
   choco install k6
   ```

2. **Start Driftlock API:**
   ```bash
   cd /Volumes/VIXinSSD/driftlock
   cargo run -p driftlock-api --release
   ```

## Test Suite

| Test | Purpose | Duration | VUs | Command |
|------|---------|----------|-----|---------|
| **smoke.js** | Quick sanity check | 1 min | 5 | `k6 run scripts/load-test/smoke.js` |
| **load.js** | Production simulation | 15 min | 0→100 | `k6 run scripts/load-test/load.js` |
| **stress.js** | Breaking point test | 10 min | 0→200 | `k6 run scripts/load-test/stress.js` |
| **soak.js** | Memory leak detection | 30 min | 50 | `k6 run scripts/load-test/soak.js` |
| **detector-capacity.js** | LRU eviction test | 10 min | 100 | `k6 run scripts/load-test/detector-capacity.js` |
| **rate-limit-validation.js** | Rate limit verification | 3 min | 1 | `k6 run scripts/load-test/rate-limit-validation.js` |

## Quick Start

### 1. Smoke Test (Start Here)

Quick sanity check to verify all endpoints work:

```bash
k6 run scripts/load-test/smoke.js
```

**Expected output:**
- All health checks pass
- Demo detect endpoint returns 200
- Response structure is valid

### 2. Load Test (Production Simulation)

Simulate realistic production traffic:

```bash
k6 run scripts/load-test/load.js
```

**Traffic pattern:**
- 70% normal detection requests (10-110 events)
- 20% large batch requests (500-1000 events)
- 10% health checks

**Success criteria:**
- P95 latency < 500ms
- P99 latency < 1000ms
- Error rate < 1%
- Throughput > 100 req/s

### 3. Stress Test (Find Breaking Points)

Push the system to its limits:

```bash
k6 run scripts/load-test/stress.js
```

**Ramp pattern:**
- 0 → 50 → 100 → 150 → 200 VUs
- Aggressive load with 200-1000 events per request

**What to watch:**
- At what VU count do errors start?
- Where does latency spike?
- Are there 5xx errors or timeouts?

### 4. Soak Test (Memory Leak Detection)

Run sustained load for 30 minutes:

```bash
k6 run scripts/load-test/soak.js
```

**Success criteria:**
- No significant performance degradation over time
- Memory leak indicator < 50%
- P95 latency remains stable

**Warning signs:**
- Increasing latency over time
- Growing error rate
- Memory leak indicator > 50%

### 5. Capacity Test (LRU Eviction)

Test detector capacity with 1000+ unique streams:

```bash
k6 run scripts/load-test/detector-capacity.js
```

**Success criteria:**
- Create 1000+ unique streams
- No errors due to detector exhaustion
- LRU eviction works correctly

### 6. Rate Limit Validation

Verify rate limiting is enforced:

```bash
k6 run scripts/load-test/rate-limit-validation.js
```

**Expected behavior:**
- First 10 requests per minute: 200 OK
- Subsequent requests: 429 Too Many Requests
- Retry-After header is present

## Configuration

### Environment Variables

```bash
# API base URL (default: http://localhost:8080)
export BASE_URL=http://localhost:8080

# API key for authenticated endpoints (optional for demo tests)
export API_KEY=your-api-key-here
```

### Custom Test Run

```bash
# Run with custom base URL
BASE_URL=https://api.driftlock.net k6 run scripts/load-test/load.js

# Run with custom VUs and duration
k6 run --vus 50 --duration 5m scripts/load-test/load.js

# Save results to file
k6 run --out json=results.json scripts/load-test/load.js

# Run with CSV output
k6 run --out csv=results.csv scripts/load-test/load.js
```

## Performance Thresholds

Default thresholds configured in `config/thresholds.json`:

| Metric | Target |
|--------|--------|
| P95 latency | < 500ms |
| P99 latency | < 1000ms |
| Error rate | < 1% |
| Success rate | > 99% |
| Throughput | > 100 req/s (load test) |

## Understanding Results

### Key Metrics

1. **http_req_duration**: Total request duration (including network)
   - P95: 95% of requests complete within this time
   - P99: 99% of requests complete within this time

2. **http_req_failed**: Percentage of failed requests
   - Target: < 1% for production

3. **detect_latency**: Server-side detection time
   - Excludes network latency
   - Pure algorithm performance

4. **events_processed**: Total events processed
   - Validates throughput capacity

### Reading K6 Output

```
✓ detect status is 200
✓ detect response is valid

checks.........................: 100.00% ✓ 1234      ✗ 0
data_received..................: 1.2 MB  40 kB/s
http_req_duration..............: avg=250ms min=50ms med=200ms max=1000ms p(95)=450ms p(99)=800ms
  { expected_response:true }...: avg=250ms min=50ms med=200ms max=1000ms p(95)=450ms p(99)=800ms
http_reqs......................: 1234    41.13/s
```

**Interpretation:**
- All checks passed (100%)
- Average latency: 250ms
- P95 latency: 450ms (good!)
- Throughput: 41 req/s

## Troubleshooting

### Test Fails with Connection Refused

**Problem:** `dial: connection refused`

**Solution:**
1. Verify API is running: `curl http://localhost:8080/healthz`
2. Check correct port (default: 8080)
3. Set BASE_URL if using different port

### High Error Rate

**Problem:** Error rate > 5%

**Possible causes:**
1. Server overloaded (reduce VUs)
2. Database connection pool exhausted
3. Memory leak (check soak test results)

**Debug:**
```bash
# Run with verbose logging
k6 run --http-debug scripts/load-test/load.js

# Check server logs for errors
```

### Rate Limit Test Shows No Rate Limiting

**Problem:** All requests succeed, no 429s

**Solution:**
1. Verify rate limiter is enabled in API config
2. Check rate limit is set to 10/min for demo endpoint
3. Ensure using demo endpoint (not authenticated endpoint)

### Soak Test Shows Degradation

**Problem:** Memory leak indicator > 50%

**Solution:**
1. Check server memory usage during test
2. Review detector cleanup logic
3. Verify LRU eviction is working
4. Check for connection/resource leaks

## Advanced Usage

### Cloud Execution

Run tests on K6 Cloud for distributed load:

```bash
# Login to K6 cloud
k6 login cloud

# Run on cloud
k6 cloud scripts/load-test/load.js
```

### Custom Scenarios

Create custom test by copying a template:

```bash
cp scripts/load-test/load.js scripts/load-test/custom.js
# Edit custom.js with your scenario
k6 run scripts/load-test/custom.js
```

### Grafana Integration

Stream metrics to Grafana:

```bash
# Install K6 Prometheus exporter
docker run -p 9090:9090 prom/prometheus

# Run with Prometheus output
k6 run --out statsd scripts/load-test/load.js
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run K6 Load Tests
  run: |
    k6 run --out json=results.json scripts/load-test/smoke.js
    k6 run --out json=results.json scripts/load-test/load.js
```

### Performance Regression Detection

Compare results over time:

```bash
# Baseline
k6 run --out json=baseline.json scripts/load-test/load.js

# After changes
k6 run --out json=current.json scripts/load-test/load.js

# Compare (manual or automated)
```

## Best Practices

1. **Start small:** Always run smoke test first
2. **Ramp gradually:** Don't jump straight to max load
3. **Monitor server:** Watch CPU, memory, connections during tests
4. **Realistic data:** Use helpers to generate realistic event payloads
5. **Clean state:** Reset database between major test runs
6. **Document findings:** Save results and note any issues

## Next Steps

1. Run smoke test to verify setup
2. Run load test to establish baseline
3. Run stress test to find limits
4. Run soak test overnight to catch leaks
5. Document results in `docs/PERFORMANCE.md`

## Support

Questions or issues? Check:
- API logs: `cargo run -p driftlock-api --release`
- K6 docs: https://k6.io/docs/
- Driftlock docs: `docs/TESTING.md`
