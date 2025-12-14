# Driftlock Testing Guide

Comprehensive testing procedures for the Driftlock Rust API.

## Quick Start

```bash
# Run all tests
cargo test --workspace

# Run API tests only
cargo test -p driftlock-api

# Run with output
cargo test -p driftlock-api -- --nocapture

# Run specific test
cargo test -p driftlock-api test_health_check
```

## Test Structure

```
driftlock/
├── cbad-core/
│   └── src/           # Unit tests in source files
├── crates/
│   ├── driftlock-api/
│   │   └── src/       # API integration tests
│   ├── driftlock-db/
│   │   └── src/       # Database layer tests
│   └── driftlock-auth/
│       └── src/       # Auth tests
└── tests/             # E2E tests (if any)
```

## Unit Tests

### CBAD Core Tests

```bash
# Run CBAD algorithm tests
cargo test -p cbad-core

# Run with release optimizations (faster)
cargo test -p cbad-core --release

# Run specific module
cargo test -p cbad-core anomaly
```

### Database Layer Tests

```bash
# Run database tests (requires DATABASE_URL)
DATABASE_URL="postgres://..." cargo test -p driftlock-db
```

### Auth Tests

```bash
# Run auth tests
cargo test -p driftlock-auth
```

## Integration Tests

### API Server Tests

```bash
# Full API test suite
cargo test -p driftlock-api

# Health endpoint tests
cargo test -p driftlock-api health

# Detection tests
cargo test -p driftlock-api detection
```

## Manual Testing

### Health Checks

```bash
# Start the server
cargo run -p driftlock-api --release

# Test liveness
curl http://localhost:8080/healthz

# Test readiness
curl http://localhost:8080/readyz

# Expected response
{
  "status": "ok",
  "version": "1.0.0"
}
```

### Demo Detection (No Auth)

```bash
curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "normal log entry 1",
      "normal log entry 2",
      "ERROR: unexpected failure in module X"
    ]
  }'
```

### Authenticated Detection

```bash
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: your-api-key" \
  -d '{
    "stream_id": "my-stream",
    "events": ["event1", "event2", "event3"]
  }'
```

### Prometheus Metrics

```bash
curl http://localhost:8080/metrics

# Check specific metrics
curl -s http://localhost:8080/metrics | grep driftlock_http_requests_total
curl -s http://localhost:8080/metrics | grep driftlock_anomalies_detected_total
```

## Load Testing

### Basic Performance Test

```bash
# Time a request
time curl -X POST http://localhost:8080/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["test event 1", "test event 2"]}'

# Expected: < 100ms for small payloads
```

### Concurrent Requests

```bash
# 10 concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost:8080/v1/demo/detect \
    -H "Content-Type: application/json" \
    -d '{"events": ["event'$i'"]}' &
done
wait
```

### K6 Load Test (Optional)

```javascript
// load-test.js
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  vus: 10,
  duration: '30s',
};

export default function() {
  let res = http.post('http://localhost:8080/v1/demo/detect',
    JSON.stringify({ events: ['test event'] }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, { 'status is 200': (r) => r.status === 200 });
}
```

```bash
k6 run load-test.js
```

## Error Handling Tests

### Invalid JSON

```bash
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  -d "invalid json"

# Expected: 400 Bad Request
```

### Missing Auth

```bash
curl -X POST http://localhost:8080/v1/detect \
  -H "Content-Type: application/json" \
  -d '{"events": ["test"]}'

# Expected: 401 Unauthorized
```

### Rate Limiting

```bash
# Demo endpoint: 10 req/hour per IP
for i in {1..15}; do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -X POST http://localhost:8080/v1/demo/detect \
    -H "Content-Type: application/json" \
    -d '{"events": ["test"]}'
done

# Expected: First 10 return 200, then 429 Too Many Requests
```

## Database Testing

### Migration Testing

```bash
# Run migrations
sqlx migrate run --database-url "$DATABASE_URL"

# Check migration status
sqlx migrate info --database-url "$DATABASE_URL"

# Revert last migration (if needed)
sqlx migrate revert --database-url "$DATABASE_URL"
```

### Query Testing

```bash
# Test database connection
psql "$DATABASE_URL" -c "SELECT 1"

# Check tables
psql "$DATABASE_URL" -c "\dt"
```

## CI/CD Testing

### GitHub Actions

The CI workflow runs:
1. `cargo fmt --check` - Code formatting
2. `cargo clippy` - Linting
3. `cargo build --release` - Build verification
4. `cargo test` - All tests

### Local CI Simulation

```bash
# Run full CI check locally
cargo fmt --check && \
cargo clippy -- -D warnings && \
cargo build --release && \
cargo test
```

## Test Data

### Sample Events

```json
{
  "events": [
    "2025-01-15T10:00:00Z INFO User logged in",
    "2025-01-15T10:00:01Z INFO Page viewed",
    "2025-01-15T10:00:02Z ERROR Database connection failed",
    "2025-01-15T10:00:03Z WARN High memory usage detected"
  ]
}
```

### Anomalous Events

```json
{
  "events": [
    "CRITICAL: System breach detected",
    "ERROR: Unauthorized access attempt from 192.168.1.100",
    "ALERT: Multiple failed login attempts"
  ]
}
```

## Success Criteria

### Critical (Must Pass)
- [x] All unit tests pass
- [x] API server starts and responds to health checks
- [x] `/v1/demo/detect` endpoint works
- [x] Authentication works (API key + Firebase)
- [x] Database operations work

### Important (Should Pass)
- [x] Rate limiting enforced
- [x] Error responses follow format
- [x] Prometheus metrics exposed
- [x] Performance < 100ms p99 for small payloads

### Optional (Nice to Have)
- [ ] Load testing passes
- [ ] 100% code coverage
- [ ] All edge cases handled

## Troubleshooting

### Tests Failing

```bash
# Run with verbose output
cargo test -- --nocapture

# Run single test with backtrace
RUST_BACKTRACE=1 cargo test test_name -- --nocapture
```

### Database Connection Issues

```bash
# Check DATABASE_URL is set
echo $DATABASE_URL

# Test connection
psql "$DATABASE_URL" -c "SELECT 1"
```

### Build Issues

```bash
# Clean and rebuild
cargo clean
cargo build --release
```

## Driftlog

For detailed debugging, check the driftlog output:

```bash
# Run with debug logging
RUST_LOG=debug cargo run -p driftlock-api

# Filter to specific modules
RUST_LOG=driftlock_api=debug,driftlock_db=info cargo run -p driftlock-api
```
