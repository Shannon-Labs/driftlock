# K6 Load Test Suite - Summary

Comprehensive stress testing suite for Driftlock Rust API production readiness validation.

## Files Created

```
scripts/load-test/
├── config/
│   └── thresholds.json           # Performance threshold definitions
├── helpers.js                    # Event generators, auth helpers, utilities
├── smoke.js                      # Quick sanity check (1 min, 5 VUs)
├── load.js                       # Production simulation (15 min, 0→100 VUs)
├── stress.js                     # Breaking point test (10 min, 0→200 VUs)
├── soak.js                       # Memory leak detection (30 min, 50 VUs)
├── detector-capacity.js          # LRU eviction test (10 min, 100 VUs)
├── rate-limit-validation.js      # Rate limit verification (3 min, 1 VU)
├── README.md                     # Full documentation
├── QUICKSTART.md                 # 5-minute setup guide
└── validate-scripts.sh           # Script validation tool
```

## Test Suite Overview

### 1. Smoke Test (`smoke.js`)
**Purpose:** Quick sanity check to verify all endpoints work

**What it tests:**
- Health check endpoint (`/healthz`)
- Readiness check endpoint (`/readyz`)
- Demo detect endpoint (`/v1/demo/detect`)
- Invalid endpoint handling (404s)
- Response structure validation

**Duration:** 1 minute
**VUs:** 5
**Run:** `k6 run scripts/load-test/smoke.js`

**Success criteria:**
- All health checks pass (>99%)
- Demo detect success rate >95%
- P95 latency <500ms

---

### 2. Load Test (`load.js`)
**Purpose:** Simulate realistic production traffic patterns

**What it tests:**
- Normal detection requests (10-110 events, 70% of traffic)
- Large batch requests (500-1000 events, 20% of traffic)
- Health checks (10% of traffic)
- Varied detection profiles (sensitive, balanced, strict)

**Duration:** 15 minutes
**VUs:** Ramps 0 → 25 → 50 → 100 → 50 → 0
**Run:** `k6 run scripts/load-test/load.js`

**Success criteria:**
- P95 latency <500ms
- P99 latency <1000ms
- Error rate <1%
- Throughput >100 req/s
- Process >10,000 events

**Custom metrics tracked:**
- `detect_latency`: Server-side detection time
- `events_processed`: Total events processed
- `anomalies_detected`: Total anomalies found
- `detect_success`: Detection success rate

---

### 3. Stress Test (`stress.js`)
**Purpose:** Push system to limits to find breaking points

**What it tests:**
- Aggressive load with large batches (200-1000 events)
- Rapid creation of unique streams
- System behavior under extreme load
- Error handling under stress

**Duration:** 10 minutes
**VUs:** Ramps 0 → 50 → 100 → 150 → 200 → 150 → 0
**Run:** `k6 run scripts/load-test/stress.js`

**Success criteria:**
- P95 latency <1000ms
- P99 latency <2000ms
- Error rate <5% (more lenient for stress)

**What to observe:**
- At what VU count do errors start?
- Where does latency spike?
- Are there 5xx errors or timeouts?
- Does system recover gracefully?

**Outputs:** `stress_test_results.json`

---

### 4. Soak Test (`soak.js`)
**Purpose:** Run sustained load to identify memory leaks and degradation

**What it tests:**
- Sustained 50 VU load for 30 minutes
- Performance stability over time
- Memory leak detection
- Resource exhaustion scenarios

**Duration:** 30 minutes
**VUs:** 50 (constant)
**Run:** `k6 run scripts/load-test/soak.js`

**Success criteria:**
- P95 latency <600ms (stable over time)
- P99 latency <1200ms (stable over time)
- Error rate <1%
- Memory leak indicator <50%

**Key metric:**
- `memory_leak_indicator`: Compares early vs. late performance
- `response_time_trend`: Tracks latency over time

**Warning signs:**
- Increasing latency over time
- Growing error rate
- Degradation >20% from start to finish

**Outputs:** `soak_test_results.json`

---

### 5. Detector Capacity Test (`detector-capacity.js`)
**Purpose:** Test LRU eviction with 1000+ unique streams

**What it tests:**
- Create 1000+ unique streams
- LRU cache eviction behavior
- Multi-stream handling
- Detector memory management

**Duration:** 10 minutes
**VUs:** 100
**Run:** `k6 run scripts/load-test/detector-capacity.js`

**Success criteria:**
- Create 1000+ unique streams
- No errors due to detector exhaustion
- P95 latency <1000ms
- Success rate >95%

**Key metric:**
- `unique_streams_created`: Must exceed 1000

**What it validates:**
- LRU eviction works correctly
- System handles stream churn
- No memory exhaustion with many streams

**Outputs:** `capacity_test_results.json`

---

### 6. Rate Limit Validation (`rate-limit-validation.js`)
**Purpose:** Verify rate limiting is enforced correctly

**What it tests:**
- Demo endpoint rate limit (10 req/min per IP)
- 429 response behavior
- Retry-After header presence
- Rate limit timing accuracy

**Duration:** 3 minutes
**VUs:** 1 (single IP simulation)
**Run:** `k6 run scripts/load-test/rate-limit-validation.js`

**Expected behavior:**
- First 10 requests per minute: 200 OK
- Subsequent requests: 429 Too Many Requests
- Retry-After header present

**Success criteria:**
- Rate limited responses >50%
- P95 latency <200ms (429s are fast)

**Outputs:** `rate_limit_test_results.json`

---

## Helper Functions (`helpers.js`)

### Event Generators
- `generateLogEvent()`: Realistic log events
- `generateAnomalousEvent()`: High-entropy anomalous events
- `generateMetricEvent()`: Metric data points
- `generateEventBatch()`: Batch of events with configurable anomaly rate
- `generateMixedBatch()`: Mixed logs, metrics, traces
- `generateTransaction()`: Financial transaction data
- `generateLLMEvent()`: LLM I/O events

### Utilities
- `authHeaders()`: API key authentication headers
- `demoHeaders()`: Demo endpoint headers
- `getBaseUrl()`: Get API URL from env or default
- `getApiKey()`: Get API key from env or test key
- `createDetectPayload()`: Build detect request
- `validateDetectResponse()`: Parse and validate response
- `generateUniqueStreamId()`: Create unique stream IDs
- `sleepWithJitter()`: Sleep with randomization
- `recordDetectionMetrics()`: Track custom metrics
- `isRateLimited()`: Check for 429 response

### Environment Variables
- `BASE_URL`: API base URL (default: `http://localhost:8080`)
- `API_KEY`: API key for authenticated endpoints

---

## Performance Thresholds (`config/thresholds.json`)

Pre-configured thresholds for each test type:

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
  },
  "soak": {
    "http_req_duration": ["p(95)<600", "p(99)<1200"],
    "detect_latency": ["p(95)<500", "p(99)<1000"]
  },
  "capacity": {
    "unique_streams_created": ["value>1000"]
  },
  "rate_limit": {
    "rate_limited_responses": ["rate>0.5"]
  }
}
```

---

## Quick Start

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

### 3. Run Smoke Test
```bash
k6 run scripts/load-test/smoke.js
```

### 4. Run Full Suite
```bash
# Smoke (1 min)
k6 run scripts/load-test/smoke.js

# Load (15 min)
k6 run scripts/load-test/load.js

# Stress (10 min)
k6 run scripts/load-test/stress.js

# Capacity (10 min)
k6 run scripts/load-test/detector-capacity.js

# Rate limit (3 min)
k6 run scripts/load-test/rate-limit-validation.js

# Soak (30 min) - run last or overnight
k6 run scripts/load-test/soak.js
```

---

## Recommended Testing Sequence

### Phase 1: Validation (2 minutes)
1. Run `smoke.js` - Verify everything works
2. Fix any issues before proceeding

### Phase 2: Baseline (15 minutes)
1. Run `load.js` - Establish baseline performance
2. Document P95/P99 latencies
3. Note throughput capacity

### Phase 3: Limits (13 minutes)
1. Run `stress.js` - Find breaking points
2. Run `detector-capacity.js` - Validate LRU eviction
3. Run `rate-limit-validation.js` - Verify rate limits

### Phase 4: Stability (30 minutes+)
1. Run `soak.js` - Check for memory leaks
2. Monitor server resources during test
3. Compare early vs. late performance

---

## Expected Results (Healthy System)

### Smoke Test
```
checks.........................: 100.00% ✓ 1500
http_req_duration..............: avg=45ms p(95)=85ms p(99)=110ms
http_req_failed................: 0.00%
```

### Load Test
```
checks.........................: 99.95%
http_req_duration..............: avg=250ms p(95)=450ms p(99)=800ms
http_reqs......................: 120/s
detect_latency.................: p(95)=350ms p(99)=650ms
events_processed...............: 50000+
```

### Stress Test
```
checks.........................: 98.00%
http_req_duration..............: avg=500ms p(95)=900ms p(99)=1500ms
http_req_failed................: 2.00%
errors start around 180-200 VUs
```

### Soak Test
```
checks.........................: 99.90%
http_req_duration..............: avg=280ms p(95)=550ms p(99)=950ms
memory_leak_indicator..........: <0.15 (15% degradation)
stable performance over 30 min
```

### Capacity Test
```
unique_streams_created.........: 1200+
http_req_failed................: 3.00%
LRU eviction working correctly
```

### Rate Limit Test
```
rate_limited_responses.........: 65%
successful_requests............: ~30 (10/min × 3 min)
429 responses after limit hit
```

---

## Troubleshooting

### Connection Refused
**Problem:** `dial: connection refused`

**Fix:**
```bash
# Start API
cargo run -p driftlock-api --release

# Verify
curl http://localhost:8080/healthz
```

### High Error Rate
**Problem:** Error rate >5% in load test

**Check:**
1. Server logs for errors
2. Database connection pool size
3. Memory usage during test

### No Rate Limiting
**Problem:** All requests succeed in rate limit test

**Fix:**
1. Verify rate limiter is enabled
2. Check demo endpoint configuration
3. Ensure using correct endpoint

---

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run K6 Load Tests
  run: |
    brew install k6
    k6 run --out json=smoke-results.json scripts/load-test/smoke.js
    k6 run --out json=load-results.json scripts/load-test/load.js
```

---

## Next Steps

1. Run smoke test to verify setup
2. Run load test to establish baseline
3. Document results in `docs/PERFORMANCE.md`
4. Set up automated testing in CI/CD
5. Run soak test overnight before production deploy

---

## Support

- Full docs: `README.md`
- Quick start: `QUICKSTART.md`
- K6 docs: https://k6.io/docs/
- Driftlock docs: `docs/TESTING.md`
