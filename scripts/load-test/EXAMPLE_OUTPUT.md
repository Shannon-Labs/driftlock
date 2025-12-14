# K6 Load Test - Example Output

This document shows what successful test runs look like.

## Smoke Test Output

```
          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
    script: scripts/load-test/smoke.js
    output: -

  scenarios: (100.00%) 1 scenario, 5 max VUs, 1m30s max duration (incl. graceful stop):
           * default: 5 looping VUs for 1m0s (gracefulStop: 30s)

INFO[0000] Starting smoke test against http://localhost:8080
INFO[0000] Testing basic endpoint functionality...

running (1m01.0s), 0/5 VUs, 300 complete and 0 interrupted iterations
default ✓ [======================================] 5 VUs  1m0s

     ✓ health check returns 200
     ✓ health check has ok status
     ✓ readiness check returns 200
     ✓ demo detect returns 200
     ✓ demo detect has processed count
     ✓ demo detect has anomalies array
     ✓ invalid endpoint returns 404

     checks.........................: 100.00% ✓ 2100      ✗ 0
     data_received..................: 2.1 MB  35 kB/s
     data_sent......................: 450 kB  7.4 kB/s
     demo_detect_success............: 100.00% ✓ 300       ✗ 0
     health_check_success...........: 100.00% ✓ 300       ✗ 0
     http_req_blocked...............: avg=12.34µs  min=2µs      med=6µs      max=1.2ms    p(95)=15µs     p(99)=45µs
     http_req_connecting............: avg=4.23µs   min=0s       med=0s       max=890µs    p(95)=0s       p(99)=0s
     http_req_duration..............: avg=45.67ms  min=5.12ms   med=38.23ms  max=125.45ms p(95)=85.34ms  p(99)=108.12ms
       { expected_response:true }...: avg=45.67ms  min=5.12ms   med=38.23ms  max=125.45ms p(95)=85.34ms  p(99)=108.12ms
     http_req_failed................: 0.00%   ✓ 0         ✗ 1200
     http_req_receiving.............: avg=156.45µs min=45µs     med=134µs    max=1.23ms   p(95)=245µs    p(99)=456µs
     http_req_sending...............: avg=45.23µs  min=12µs     med=38µs     max=234µs    p(95)=78µs     p(99)=134µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(95)=0s       p(99)=0s
     http_req_waiting...............: avg=45.47ms  min=5.01ms   med=38.03ms  max=125.12ms p(95)=85.12ms  p(99)=107.89ms
     http_reqs......................: 1200    19.67/s
     iteration_duration.............: avg=3.04s    min=3s       med=3.04s    max=3.18s    p(95)=3.12s    p(99)=3.15s
     iterations.....................: 300     4.917/s
     vus............................: 5       min=5       max=5
     vus_max........................: 5       min=5       max=5

INFO[0061] Smoke test complete!
INFO[0061] Tested against: http://localhost:8080

PASS: All checks passed
```

**Interpretation:**
- ✅ All 2100 checks passed (100%)
- ✅ P95 latency: 85ms (excellent, < 500ms target)
- ✅ P99 latency: 108ms (excellent, < 1000ms target)
- ✅ 0% error rate
- ✅ Ready for load testing

---

## Load Test Output

```
          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
    script: scripts/load-test/load.js
    output: -

  scenarios: (100.00%) 1 scenario, 100 max VUs, 15m30s max duration (incl. graceful stop):
           * default: Up to 100 looping VUs for 15m0s over 6 stages (gracefulRampDown: 30s, gracefulStop: 30s)

INFO[0000] Starting load test against http://localhost:8080
INFO[0000] Simulating production traffic patterns...

running (15m00.2s), 000/100 VUs, 8523 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  15m0s

     ✓ detect status is 200
     ✓ detect response is valid
     ✓ large batch status is 200
     ✓ large batch processes all events
     ✓ health check is 200

     anomalies_detected.............: 3421    3.802/s
     checks.........................: 99.95%  ✓ 42510     ✗ 21
     data_received..................: 89 MB   99 kB/s
     data_sent......................: 125 MB  139 kB/s
     detect_success.................: 99.98%  ✓ 8521      ✗ 2
     events_processed...............: 523401  582.666/s
     http_req_blocked...............: avg=15.23µs  min=2µs      med=7µs      max=2.1ms    p(95)=24µs     p(99)=67µs
     http_req_connecting............: avg=5.12µs   min=0s       med=0s       max=1.2ms    p(95)=0s       p(99)=0s
     http_req_duration..............: avg=245.67ms min=12.34ms  med=198.45ms max=1234.56ms p(95)=456.78ms p(99)=789.12ms
       { expected_response:true }...: avg=245.34ms min=12.34ms  med=198.23ms max=1234.56ms p(95)=456.45ms p(99)=788.23ms
     http_req_failed................: 0.024%  ✓ 2         ✗ 8521
     http_req_receiving.............: avg=234.56µs min=67µs     med=189µs    max=3.45ms   p(95)=456µs    p(99)=789µs
     http_req_sending...............: avg=78.23µs  min=23µs     med=67µs     max=890µs    p(95)=156µs    p(99)=234µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(95)=0s       p(99)=0s
     http_req_waiting...............: avg=245.36ms min=12.12ms  med=198.12ms max=1234.23ms p(95)=456.34ms p(99)=788.67ms
     http_reqs......................: 8523    9.469/s
     iteration_duration.............: avg=1.54s    min=500ms    med=1.48s    max=3.89s    p(95)=2.12s    p(99)=2.67s
     iterations.....................: 8523    9.469/s
     vus............................: 1       min=1       max=100
     vus_max........................: 100     min=100     max=100

     detect_latency.................: avg=234.56ms min=10.23ms  med=189.34ms max=1189.23ms p(95)=423.45ms p(99)=756.78ms

INFO[0900] Load test complete! Duration: 900.02s
INFO[0900] Tested against: http://localhost:8080

PASS: Production traffic simulation successful
```

**Interpretation:**
- ✅ 99.95% checks passed
- ✅ P95 latency: 456ms (within 500ms target)
- ✅ P99 latency: 789ms (within 1000ms target)
- ✅ 0.024% error rate (well below 1% target)
- ✅ Throughput: 9.5 req/s (scales with VUs)
- ✅ Processed 523k events
- ✅ Detected 3421 anomalies
- ✅ System handled 100 concurrent VUs

---

## Stress Test Output (Finding Limits)

```
INFO[0000] Starting stress test against http://localhost:8080
INFO[0000] WARNING: This test will push the system to its limits!

running (10m00.3s), 000/200 VUs, 12456 complete and 0 interrupted iterations
default ✓ [======================================] 000/200 VUs  10m0s

     ✓ status is not 5xx
     ✓ status is not timeout
     ✓ response time is acceptable

     checks.........................: 97.85%  ✓ 36512     ✗ 856
     data_received..................: 156 MB  260 kB/s
     data_sent......................: 234 MB  390 kB/s
     detect_latency.................: avg=456.78ms min=23.45ms  med=389.23ms max=2345.67ms p(95)=923.45ms p(99)=1567.89ms
     error_rate.....................: 2.15%   ✓ 268       ✗ 12188
     events_processed...............: 678901  1131.5/s
     http_req_blocked...............: avg=23.45µs  min=3µs      med=12µs     max=4.5ms    p(95)=45µs     p(99)=123µs
     http_req_duration..............: avg=567.89ms min=45.67ms  med=478.23ms max=2678.45ms p(95)=1123.45ms p(99)=1789.23ms
     http_req_failed................: 2.15%   ✓ 268       ✗ 12188
     http_reqs......................: 12456   20.76/s
     server_errors..................: 145     errors started at ~180 VUs
     timeouts.......................: 123     timeouts at ~195 VUs
     iteration_duration.............: avg=1.23s    min=600ms    med=1.12s    max=4.56s    p(95)=2.34s    p(99)=3.12s
     vus............................: 1       min=1       max=200

INFO[0601] Stress test complete! Duration: 600.12s

=== Stress Test Summary ===
Total requests: 12456
Failed requests: 2.15%
P95 latency: 1123ms
P99 latency: 1789ms
Server errors: 145
Timeouts: 123

Breaking point identified:
- Errors start at ~180 VUs
- Timeouts at ~195 VUs
- System degrades gracefully under extreme load
```

**Interpretation:**
- ✅ Handled 12,456 requests with 200 concurrent VUs
- ⚠️ 2.15% error rate (acceptable for stress test)
- ⚠️ P95 latency: 1123ms (higher than normal, expected)
- ✅ Breaking point identified: ~180-195 VUs
- ✅ System degrades gracefully (no crashes)
- 📊 Capacity limit established

---

## Soak Test Output (Stability)

```
INFO[0000] Starting soak test against http://localhost:8080
INFO[0000] Duration: 30 minutes
INFO[0000] Monitoring for memory leaks and performance degradation...

running (30m00.5s), 00/50 VUs, 36789 complete and 0 interrupted iterations
default ✓ [======================================] 50 VUs  30m0s

     ✓ status is 200
     ✓ response is valid
     ✓ no significant degradation
     ✓ health check still passing

     anomalies_detected.............: 14567   8.093/s
     checks.........................: 99.92%  ✓ 146945    ✗ 111
     data_received..................: 389 MB  216 kB/s
     data_sent......................: 567 MB  315 kB/s
     detect_latency.................: avg=267.89ms min=15.67ms  med=234.56ms max=1567.89ms p(95)=534.67ms p(99)=923.45ms
     detect_success.................: 99.94%  ✓ 36767     ✗ 22
     events_processed...............: 1234567 686.981/s
     http_req_duration..............: avg=289.45ms min=23.45ms  med=256.78ms max=1678.90ms p(95)=567.89ms p(99)=1012.34ms
     http_req_failed................: 0.06%   ✓ 22        ✗ 36767
     http_reqs......................: 36789   20.438/s
     memory_leak_indicator..........: 0.087   (8.7% degradation - PASS)
     response_time_trend............: avg=289.45ms p(95)=567.89ms p(99)=1012.34ms

INFO[1801] Soak test complete! Duration: 30.01 minutes
INFO[1801] Tested against: http://localhost:8080

=== Performance Over Time ===
Early average latency: 263.45ms
Late average latency: 286.34ms
Degradation: 8.70%

PASS: No significant degradation detected.

=== Soak Test Summary ===
Total requests: 36789
Success rate: 99.94%
P95 latency: 567ms
P99 latency: 1012ms
Events processed: 1234567
Anomalies detected: 14567
```

**Interpretation:**
- ✅ 99.92% success rate over 30 minutes
- ✅ Stable latency (P95: 567ms, P99: 1012ms)
- ✅ Only 8.7% degradation (well below 20% threshold)
- ✅ No memory leak detected
- ✅ Processed 1.2M events successfully
- ✅ System is stable under sustained load

---

## Capacity Test Output (LRU Eviction)

```
INFO[0000] Starting capacity test against http://localhost:8080
INFO[0000] Testing LRU eviction with 1000+ unique streams...

INFO[0060] VU 23: Created 120 unique streams
INFO[0120] VU 45: Created 245 unique streams
INFO[0180] Estimated LRU evictions: 78 (1078 total streams)
INFO[0240] VU 67: Created 389 unique streams
INFO[0300] Estimated LRU evictions: 234 (1234 total streams)
INFO[0360] VU 89: Created 512 unique streams

running (10m00.2s), 000/100 VUs, 15234 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  10m0s

     ✓ status is 200
     ✓ response has stream_id
     ✓ detector handles new stream

     checks.........................: 98.67%  ✓ 45021     ✗ 681
     data_received..................: 234 MB  390 kB/s
     data_sent......................: 345 MB  575 kB/s
     detect_latency.................: avg=345.67ms min=34.56ms  med=289.45ms max=1456.78ms p(95)=678.90ms p(99)=1023.45ms
     detect_success.................: 98.89%  ✓ 15065     ✗ 169
     events_processed...............: 789456  1315.76/s
     http_req_duration..............: avg=378.90ms min=56.78ms  med=321.45ms max=1567.89ms p(95)=723.45ms p(99)=1134.56ms
     http_req_failed................: 1.11%   ✓ 169       ✗ 15065
     http_reqs......................: 15234   25.39/s
     unique_streams_created.........: 1523    streams

INFO[0600] Capacity test complete! Duration: 10.00 minutes
INFO[0600] Tested against: http://localhost:8080
INFO[0600] Total unique streams created: 1523

PASS: Successfully created 1000+ unique streams
Estimated LRU evictions: 523

=== Capacity Test Summary ===
Unique streams created: 1523
Total requests: 15234
Success rate: 98.89%
P95 latency: 723ms
P99 latency: 1134ms
Detect P95 latency: 678ms
Detect P99 latency: 1023ms
```

**Interpretation:**
- ✅ Created 1523 unique streams (exceeds 1000 target)
- ✅ LRU eviction triggered and working (523 evictions)
- ✅ 98.89% success rate
- ✅ System handles stream churn gracefully
- ✅ No memory exhaustion with many streams

---

## Rate Limit Test Output

```
INFO[0000] Starting rate limit validation test against http://localhost:8080
INFO[0000] Expected limit: 10 requests per minute per IP
INFO[0000] Testing with single VU to simulate single IP...

INFO[0002] Request 1: SUCCESS (200)
INFO[0004] Request 2: SUCCESS (200)
INFO[0006] Request 3: SUCCESS (200)
INFO[0008] Request 4: SUCCESS (200)
INFO[0010] Request 5: SUCCESS (200)
INFO[0012] Request 6: SUCCESS (200)
INFO[0014] Request 7: SUCCESS (200)
INFO[0016] Request 8: SUCCESS (200)
INFO[0018] Request 9: SUCCESS (200)
INFO[0020] Request 10: SUCCESS (200)
INFO[0022] Request 11: RATE LIMITED (429) - Total rate limits: 1
INFO[0024] Request 12: RATE LIMITED (429) - Total rate limits: 2
INFO[0026] Request 13: RATE LIMITED (429) - Total rate limits: 3

running (3m00.1s), 0/1 VUs, 90 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  3m0s

     ✓ rate limit returns 429
     ✓ rate limit has retry-after header
     ✓ rate limit is enforced correctly

     checks.........................: 100.00% ✓ 270       ✗ 0
     http_req_duration..............: avg=45.67ms  min=12.34ms  med=39.23ms  max=123.45ms p(95)=78.90ms  p(99)=98.76ms
     http_req_failed................: 66.67%  ✓ 60        ✗ 30
     http_reqs......................: 90      0.5/s
     rate_limit_hit_correctly.......: 100.00% ✓ 90        ✗ 0
     rate_limited_responses.........: 60      66.67%
     successful_requests............: 30      33.33%

=== Rate Limit Validation Summary ===
Test duration: 3.00 minutes
Total requests sent: 90
Successful requests: 30
Rate limited (429): 60
Expected max success: ~30

PASS: Rate limiting is working correctly!
```

**Interpretation:**
- ✅ Rate limit enforced at 10 req/min
- ✅ First 10 requests succeed per minute
- ✅ Subsequent requests return 429
- ✅ Retry-After header present
- ✅ 30 successful requests over 3 minutes (10/min)
- ✅ Rate limiter working as designed

---

## What to Look For

### Green Flags (Good)
- ✅ Check pass rate >99%
- ✅ P95 latency within thresholds
- ✅ Error rate <1% (normal tests) or <5% (stress)
- ✅ Stable performance over time (soak)
- ✅ Graceful degradation under stress
- ✅ No crashes or panics

### Yellow Flags (Investigate)
- ⚠️ Check pass rate 95-99%
- ⚠️ P95 latency slightly over threshold
- ⚠️ Error rate 1-3%
- ⚠️ Performance degradation 10-20%
- ⚠️ Occasional timeouts

### Red Flags (Issues)
- ❌ Check pass rate <95%
- ❌ P95 latency 2x+ over threshold
- ❌ Error rate >5%
- ❌ Performance degradation >20%
- ❌ Server crashes
- ❌ Memory leaks detected
- ❌ Consistent timeouts

---

## Saving Results

### JSON Output
```bash
k6 run --out json=results.json scripts/load-test/load.js
```

### CSV Output
```bash
k6 run --out csv=results.csv scripts/load-test/load.js
```

### Cloud Output
```bash
k6 login cloud
k6 cloud scripts/load-test/load.js
```

---

## Next Steps After Testing

1. Document baseline performance in `docs/PERFORMANCE.md`
2. Set up automated testing in CI/CD
3. Create performance regression tests
4. Monitor production against test results
5. Iterate on improvements based on findings
