# K6 Load Testing Quick Start

Get up and running with load tests in 5 minutes.

## 1. Install K6

### macOS
```bash
brew install k6
```

### Linux (Debian/Ubuntu)
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

### Windows
```bash
choco install k6
```

### Verify Installation
```bash
k6 version
# Expected: k6 v0.48.0 or higher
```

## 2. Start Driftlock API

```bash
cd /Volumes/VIXinSSD/driftlock
cargo run -p driftlock-api --release
```

Wait for:
```
INFO  Server running on http://0.0.0.0:8080
```

## 3. Run Your First Test (Smoke Test)

```bash
cd /Volumes/VIXinSSD/driftlock
k6 run scripts/load-test/smoke.js
```

Expected output:
```
running (1m01.0s), 0/5 VUs, 300 complete and 0 interrupted iterations

✓ health check returns 200
✓ readiness check returns 200
✓ demo detect returns 200
✓ demo detect has processed count
✓ demo detect has anomalies array

checks.........................: 100.00% ✓ 1500      ✗ 0
http_req_duration..............: avg=45ms min=5ms med=38ms max=120ms p(95)=85ms p(99)=110ms
```

If you see this, you're ready to go!

## 4. Run Load Test

```bash
k6 run scripts/load-test/load.js
```

This simulates production traffic for 15 minutes.

## Troubleshooting

### Error: Connection Refused
```
ERRO[0000] GoError: Get "http://localhost:8080/healthz": dial tcp 127.0.0.1:8080: connect: connection refused
```

**Fix:** Start the API server first:
```bash
cargo run -p driftlock-api --release
```

### Error: Module Not Found
```
ERRO[0000] module './helpers.js' not found
```

**Fix:** Run k6 from the project root:
```bash
cd /Volumes/VIXinSSD/driftlock
k6 run scripts/load-test/smoke.js
```

## What to Run

| Scenario | Test | Time |
|----------|------|------|
| Quick check | `smoke.js` | 1 min |
| Production readiness | `load.js` | 15 min |
| Find limits | `stress.js` | 10 min |
| Memory leaks | `soak.js` | 30 min |
| Stream capacity | `detector-capacity.js` | 10 min |
| Rate limiting | `rate-limit-validation.js` | 3 min |

## Next Steps

See full documentation: `README.md`
