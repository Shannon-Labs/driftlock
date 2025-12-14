# K6 Load Test Suite - Complete Index

Comprehensive K6 load testing suite for Driftlock Rust API production readiness validation.

## Overview

**Total Files:** 13
**Total Lines:** 2,629
**Test Coverage:** Health checks, detection endpoints, rate limiting, capacity, stability
**Ready to Run:** Yes (requires K6 installation)

---

## Quick Navigation

| File | Purpose | Lines |
|------|---------|-------|
| **QUICKSTART.md** | 5-minute setup guide | Start here |
| **README.md** | Full documentation | Complete reference |
| **TEST_SUMMARY.md** | Test suite overview | Understand tests |
| **EXAMPLE_OUTPUT.md** | Expected results | Interpret results |
| **INDEX.md** | This file | Navigation |

---

## Test Scripts (6 Total)

### 1. smoke.js (85 lines)
**Quick sanity check**
- Duration: 1 minute
- VUs: 5
- Tests: Health, readiness, demo detect
- Run: `k6 run scripts/load-test/smoke.js`

### 2. load.js (158 lines)
**Production simulation**
- Duration: 15 minutes
- VUs: 0 → 25 → 50 → 100
- Tests: Normal traffic, large batches, health checks
- Run: `k6 run scripts/load-test/load.js`

### 3. stress.js (131 lines)
**Breaking point test**
- Duration: 10 minutes
- VUs: 0 → 50 → 100 → 150 → 200
- Tests: Extreme load, system limits
- Run: `k6 run scripts/load-test/stress.js`

### 4. soak.js (208 lines)
**Memory leak detection**
- Duration: 30 minutes
- VUs: 50 (constant)
- Tests: Sustained load, stability, degradation
- Run: `k6 run scripts/load-test/soak.js`

### 5. detector-capacity.js (168 lines)
**LRU eviction test**
- Duration: 10 minutes
- VUs: 100
- Tests: 1000+ unique streams, LRU eviction
- Run: `k6 run scripts/load-test/detector-capacity.js`

### 6. rate-limit-validation.js (125 lines)
**Rate limit verification**
- Duration: 3 minutes
- VUs: 1 (single IP)
- Tests: 10 req/min rate limit, 429 responses
- Run: `k6 run scripts/load-test/rate-limit-validation.js`

---

## Helper Files

### helpers.js (229 lines)
**Utility functions and event generators**

**Event Generators:**
- `generateLogEvent()` - Realistic log events
- `generateAnomalousEvent()` - High-entropy anomalies
- `generateMetricEvent()` - Metric data points
- `generateEventBatch()` - Batch with anomaly rate
- `generateMixedBatch()` - Mixed event types
- `generateTransaction()` - Financial transactions
- `generateLLMEvent()` - LLM I/O events

**Utilities:**
- `authHeaders()` - Authentication headers
- `demoHeaders()` - Demo endpoint headers
- `getBaseUrl()` - API URL from env
- `getApiKey()` - API key from env
- `createDetectPayload()` - Build detect request
- `validateDetectResponse()` - Parse response
- `generateUniqueStreamId()` - Unique stream IDs
- `sleepWithJitter()` - Random sleep
- `recordDetectionMetrics()` - Track metrics
- `isRateLimited()` - Check 429 response

**Environment Variables:**
- `BASE_URL` - API base URL (default: http://localhost:8080)
- `API_KEY` - API key for authenticated endpoints

---

## Configuration

### config/thresholds.json (38 lines)
**Performance threshold definitions**

```json
{
  "smoke": {
    "http_req_duration": ["p(95)<500", "p(99)<1000"],
    "http_req_failed": ["rate<0.01"]
  },
  "load": {
    "http_req_duration": ["p(95)<500", "p(99)<1000"],
    "detect_latency": ["p(95)<400", "p(99)<800"],
    "events_processed": ["count>10000"]
  },
  "stress": {
    "http_req_duration": ["p(95)<1000", "p(99)<2000"],
    "http_req_failed": ["rate<0.05"]
  }
  // ... more thresholds
}
```

---

## Documentation (4 Files)

### README.md (332 lines)
**Complete documentation**

**Sections:**
- Prerequisites and installation
- Test suite overview
- Quick start guide
- Configuration options
- Performance thresholds
- Understanding results
- Troubleshooting
- Advanced usage
- CI/CD integration
- Best practices

### QUICKSTART.md (71 lines)
**5-minute setup guide**

**Sections:**
1. Install K6
2. Start Driftlock API
3. Run first test
4. Troubleshooting
5. What to run next

### TEST_SUMMARY.md (497 lines)
**Comprehensive test suite overview**

**Sections:**
- Files created
- Test suite overview (all 6 tests)
- Helper functions reference
- Performance thresholds
- Quick start
- Recommended testing sequence
- Expected results
- Troubleshooting
- CI/CD integration

### EXAMPLE_OUTPUT.md (489 lines)
**Expected test results**

**Sections:**
- Smoke test output (passing)
- Load test output (passing)
- Stress test output (finding limits)
- Soak test output (stable)
- Capacity test output (LRU working)
- Rate limit test output (enforced)
- What to look for (green/yellow/red flags)
- Saving results

### INDEX.md (This file)
**Navigation and overview**

---

## Validation

### validate-scripts.sh (Executable)
**Script validation tool**

**Checks:**
- K6 installation
- Script syntax
- Required files exist
- Configuration valid

**Usage:**
```bash
./scripts/load-test/validate-scripts.sh
```

---

## File Structure

```
scripts/load-test/
├── config/
│   └── thresholds.json           # Performance thresholds
│
├── Test Scripts (6)
│   ├── smoke.js                  # 1 min, 5 VUs
│   ├── load.js                   # 15 min, 0→100 VUs
│   ├── stress.js                 # 10 min, 0→200 VUs
│   ├── soak.js                   # 30 min, 50 VUs
│   ├── detector-capacity.js      # 10 min, 100 VUs
│   └── rate-limit-validation.js  # 3 min, 1 VU
│
├── Helper Files
│   └── helpers.js                # Event generators & utilities
│
├── Documentation (5)
│   ├── README.md                 # Full documentation
│   ├── QUICKSTART.md             # 5-minute guide
│   ├── TEST_SUMMARY.md           # Test overview
│   ├── EXAMPLE_OUTPUT.md         # Expected results
│   └── INDEX.md                  # This file
│
└── Tools
    └── validate-scripts.sh       # Validation script
```

---

## Getting Started (3 Steps)

### 1. Install K6
```bash
# macOS
brew install k6

# Verify
k6 version
```

### 2. Start API
```bash
cd /Volumes/VIXinSSD/driftlock
cargo run -p driftlock-api --release
```

### 3. Run Tests
```bash
# Smoke test (1 min)
k6 run scripts/load-test/smoke.js

# Load test (15 min)
k6 run scripts/load-test/load.js

# Full suite (70+ min)
for test in smoke load stress detector-capacity rate-limit-validation; do
  k6 run scripts/load-test/${test}.js
done
```

---

## Test Progression

### Phase 1: Validation (2 min)
```bash
k6 run scripts/load-test/smoke.js
```
- Verify all endpoints work
- Establish baseline

### Phase 2: Performance (15 min)
```bash
k6 run scripts/load-test/load.js
```
- Production traffic simulation
- Document P95/P99 latencies

### Phase 3: Limits (23 min)
```bash
k6 run scripts/load-test/stress.js
k6 run scripts/load-test/detector-capacity.js
k6 run scripts/load-test/rate-limit-validation.js
```
- Find breaking points
- Validate LRU eviction
- Verify rate limits

### Phase 4: Stability (30 min)
```bash
k6 run scripts/load-test/soak.js
```
- Memory leak detection
- Long-term stability
- Performance degradation

---

## Performance Targets

| Metric | Target | Test |
|--------|--------|------|
| P95 latency | <500ms | Load |
| P99 latency | <1000ms | Load |
| Error rate | <1% | Load |
| Throughput | >100 req/s | Load |
| Degradation | <20% | Soak |
| Unique streams | >1000 | Capacity |
| Rate limit | 10/min | Rate limit |

---

## Custom Metrics Tracked

### Detection Metrics
- `detect_latency` - Server-side detection time (ms)
- `events_processed` - Total events processed (count)
- `anomalies_detected` - Total anomalies found (count)
- `detect_success` - Detection success rate (rate)

### Capacity Metrics
- `unique_streams_created` - Unique streams created (count)
- `lru_evictions_estimated` - Estimated LRU evictions (count)

### Stability Metrics
- `memory_leak_indicator` - Performance degradation ratio
- `response_time_trend` - Latency trend over time

### Rate Limiting Metrics
- `rate_limited_responses` - Count of 429 responses (count)
- `rate_limit_hit_correctly` - Rate limit enforcement (rate)

---

## Common Commands

### Basic Run
```bash
k6 run scripts/load-test/smoke.js
```

### Custom Base URL
```bash
BASE_URL=https://api.driftlock.net k6 run scripts/load-test/load.js
```

### Save Results
```bash
k6 run --out json=results.json scripts/load-test/load.js
```

### Custom VUs/Duration
```bash
k6 run --vus 50 --duration 5m scripts/load-test/load.js
```

### Cloud Run
```bash
k6 login cloud
k6 cloud scripts/load-test/load.js
```

---

## Output Files Generated

When tests complete, they may generate:
- `stress_test_results.json` - Stress test detailed results
- `soak_test_results.json` - Soak test detailed results
- `capacity_test_results.json` - Capacity test detailed results
- `rate_limit_test_results.json` - Rate limit test detailed results

---

## Troubleshooting Quick Reference

### Connection Refused
```bash
# Start API
cargo run -p driftlock-api --release
```

### Module Not Found
```bash
# Run from project root
cd /Volumes/VIXinSSD/driftlock
k6 run scripts/load-test/smoke.js
```

### High Error Rate
- Check server logs
- Verify database connection
- Monitor memory usage

### No Rate Limiting
- Verify rate limiter enabled
- Check demo endpoint config
- Ensure correct endpoint

---

## CI/CD Integration Example

```yaml
name: K6 Load Tests

on: [push, pull_request]

jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Start API
        run: |
          cargo run -p driftlock-api --release &
          sleep 10

      - name: Install K6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Run Smoke Test
        run: k6 run --out json=smoke-results.json scripts/load-test/smoke.js

      - name: Run Load Test
        run: k6 run --out json=load-results.json scripts/load-test/load.js

      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: k6-results
          path: |
            smoke-results.json
            load-results.json
```

---

## Next Steps

1. **Install K6:** See QUICKSTART.md
2. **Run smoke test:** Verify setup
3. **Run load test:** Establish baseline
4. **Run full suite:** Comprehensive validation
5. **Document results:** Save to `docs/PERFORMANCE.md`
6. **Automate:** Add to CI/CD pipeline

---

## Support & Resources

- **Quick Start:** `QUICKSTART.md`
- **Full Docs:** `README.md`
- **Test Details:** `TEST_SUMMARY.md`
- **Example Output:** `EXAMPLE_OUTPUT.md`
- **K6 Docs:** https://k6.io/docs/
- **Driftlock Docs:** `../../docs/TESTING.md`

---

## License

Part of the Driftlock project. See main LICENSE file.

---

**Last Updated:** 2025-12-11
**Status:** Ready for use
**Version:** 1.0.0
