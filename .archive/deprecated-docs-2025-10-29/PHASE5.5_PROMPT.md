# Phase 5.5 / 6 Prompt - Testing Completion & Production Deployment

## Context

Phase 5 has been successfully completed with core production features implemented:
- ✅ Comprehensive handler tests (13 tests, all passing)
- ✅ Production hardening (validation, rate limiting, error handling)
- ✅ Structured logging with slog
- ✅ Database migrations
- ✅ API and deployment documentation
- ✅ Performance benchmarks

**Current State**:
- Test coverage: 29.6% overall, >70% for critical handlers
- All code compiles and runs
- CBAD Rust library still blocked (documented)

See `PHASE5_SUMMARY.md` for complete Phase 5 details.

## Your Mission

Complete the testing pyramid and prepare for production deployment.

---

## Priority 1: Complete Test Coverage (Aim for 70%+)

### 1.1 Storage Layer Integration Tests

**File**: `api-server/internal/storage/postgres_test.go`

**Requirements**:
- Use Docker for test PostgreSQL instance
- Test all CRUD operations
- Test connection pooling
- Test transaction handling
- Test error recovery

**Example Setup**:
```go
func setupTestDB(t *testing.T) *Storage {
    // Start test PostgreSQL via testcontainers-go
    // Run migrations
    // Return storage instance
}

func TestCreateAnomaly_Integration(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    // Test actual database operations
}
```

### 1.2 Additional Handler Tests

**Files to Test**:
- `api-server/internal/handlers/analytics.go` (0% coverage)
- `api-server/internal/handlers/config.go` (0% coverage)
- `api-server/internal/handlers/export.go` (0% coverage)

**Create**: 
- `analytics_test.go`
- `config_test.go`
- `export_test.go`

**Goal**: Achieve 70%+ coverage for each handler

### 1.3 Middleware Tests

**File**: `api-server/internal/middleware/middleware_test.go`

**Test**:
- Rate limiting behavior
- Request ID generation
- Logging middleware
- CORS handling
- Recovery from panics

### 1.4 Validation Tests

**File**: `api-server/internal/validation/validator_test.go`

**Test**:
- Valid and invalid payloads
- Edge cases (empty, null, overflow)
- Type validation
- Size limits

---

## Priority 2: End-to-End Testing

### 2.1 E2E Test Suite

**File**: `tests/e2e/e2e_test.go`

**Test Flow**:
1. Start API server
2. Start PostgreSQL
3. Generate synthetic OTLP data
4. POST to `/v1/events`
5. Verify anomaly creation via `/v1/anomalies`
6. Test SSE stream `/v1/stream/anomalies`
7. Export anomaly
8. Update status

**Example**:
```go
func TestE2E_AnomalyDetection(t *testing.T) {
    // Setup
    server := startTestServer(t)
    defer server.Close()
    
    // Generate data
    data := generateSyntheticOTLP()
    
    // Post events
    resp := postEvents(server.URL, data)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verify anomalies
    anomalies := getAnomalies(server.URL)
    assert.Greater(t, len(anomalies), 0)
}
```

### 2.2 SSE Stream Testing

**Test**:
- Client connection/disconnection
- Heartbeat delivery
- Anomaly broadcast
- Client limit enforcement
- Reconnection handling

---

## Priority 3: CI/CD Pipeline

### 3.1 GitHub Actions Workflow

**File**: `.github/workflows/ci.yml`

**Jobs**:
1. **Lint**: golangci-lint, ESLint
2. **Test**: Run all tests with coverage
3. **Build**: Build API server and UI
4. **Integration**: Run E2E tests
5. **Security**: Dependency scanning

**Example**:
```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - run: go test -v -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: actions/upload-artifact@v3
        with:
          name: coverage
          path: coverage.html
```

### 3.2 Docker Build

**Files**:
- `api-server/Dockerfile`
- `ui/Dockerfile`
- `docker-compose.yml`

**Build**:
```dockerfile
# Multi-stage build for API server
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o driftlock-api ./cmd/driftlock-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/driftlock-api /usr/local/bin/
EXPOSE 8080
CMD ["driftlock-api"]
```

---

## Priority 4: Monitoring & Observability

### 4.1 Prometheus Metrics

**File**: `api-server/internal/metrics/prometheus.go`

**Metrics to Implement**:
```go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    anomaliesDetected = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "anomalies_detected_total",
            Help: "Total anomalies detected",
        },
        []string{"stream_type", "severity"},
    )
)
```

### 4.2 Grafana Dashboard

**File**: `monitoring/grafana/dashboard.json`

**Panels**:
- Request rate (req/s)
- Error rate (%)
- Response time (p50, p95, p99)
- Anomalies by stream type
- Database connection pool usage
- Memory usage

---

## Priority 5: Performance Testing

### 5.1 Load Testing

**Tool**: k6 or wrk

**File**: `tests/load/load_test.js`

**Scenarios**:
1. Sustained load: 1000 req/s for 10 minutes
2. Spike test: 0 → 5000 req/s in 30 seconds
3. Stress test: Gradually increase until failure

**Example (k6)**:
```javascript
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 1000 },
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 0 },
  ],
};

export default function() {
  let response = http.get('http://localhost:8080/v1/anomalies');
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
}
```

### 5.2 Database Performance

**Optimize**:
- Add missing indexes based on slow query log
- Tune connection pool settings
- Implement query result caching
- Add database monitoring

---

## Priority 6: Security Hardening

### 6.1 Security Audit

**Checklist**:
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (prepared statements)
- [ ] XSS prevention (Content-Security-Policy)
- [ ] CSRF protection (if cookies used)
- [ ] Rate limiting enabled
- [ ] HTTPS enforced in production
- [ ] Secrets not in code/git
- [ ] Dependencies scanned for vulnerabilities

### 6.2 Authentication Implementation

**Options**:
1. API Key authentication (simple)
2. JWT tokens (more features)
3. OAuth2 (enterprise)

**File**: `api-server/internal/auth/auth.go`

**Implement**:
- Token validation
- User context injection
- Permission checking
- Token refresh (if JWT)

---

## Priority 7: Production Deployment

### 7.1 Kubernetes Deployment

**Files**:
- `k8s/deployment.yaml`
- `k8s/service.yaml`
- `k8s/configmap.yaml`
- `k8s/secrets.yaml`

**Example**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: driftlock-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: driftlock-api
  template:
    metadata:
      labels:
        app: driftlock-api
    spec:
      containers:
      - name: api
        image: driftlock-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: driftlock-secrets
              key: database-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8080
```

### 7.2 Helm Chart

**File**: `helm/driftlock/Chart.yaml`

**Values**:
- Replica count
- Resource limits
- Database config
- Secrets
- Ingress rules

---

## Priority 8: Documentation Completion

### 8.1 Developer Guide

**File**: `docs/DEVELOPMENT.md`

**Contents**:
- Local setup
- Running tests
- Code style
- Contribution workflow
- Architecture overview

### 8.2 Architecture Diagrams

**Tool**: Mermaid or PlantUML

**Diagrams**:
- System architecture
- Data flow
- Database schema
- Deployment topology

### 8.3 Runbooks

**File**: `docs/RUNBOOKS.md`

**Scenarios**:
- High CPU usage
- Database connection exhaustion
- Memory leaks
- Slow queries
- Deployment rollback

---

## Success Criteria

By the end of Phase 5.5/6, you should have:

### ✅ Testing
- [ ] 70%+ test coverage overall
- [ ] All handlers tested
- [ ] Storage layer integration tests
- [ ] E2E test suite
- [ ] Load testing completed

### ✅ CI/CD
- [ ] GitHub Actions workflow
- [ ] Automated testing
- [ ] Docker builds
- [ ] Deployment automation

### ✅ Monitoring
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Alerting rules
- [ ] Log aggregation

### ✅ Security
- [ ] Security audit complete
- [ ] Authentication implemented
- [ ] Secrets management
- [ ] TLS enabled

### ✅ Deployment
- [ ] Kubernetes manifests
- [ ] Helm chart
- [ ] Production checklist
- [ ] Disaster recovery plan

### ✅ Documentation
- [ ] Developer guide
- [ ] Architecture diagrams
- [ ] Runbooks
- [ ] API versioning strategy

---

## Implementation Order

1. **Week 1**: Testing completion (storage, handlers, middleware)
2. **Week 2**: E2E tests + CI/CD pipeline
3. **Week 3**: Monitoring, metrics, load testing
4. **Week 4**: Security, authentication, production prep

---

## Known Blockers

1. **CBAD Build**: Still blocked by network restrictions
   - **Workaround**: Use mock CBAD for testing
   - **Long-term**: Build in environment with crates.io access

---

## Getting Started

```bash
# 1. Review Phase 5 work
cat PHASE5_SUMMARY.md

# 2. Check current test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# 3. Start with storage tests
cd api-server/internal/storage
touch postgres_test.go

# 4. Run tests
go test -v ./...
```

---

## Resources

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [k6 Load Testing](https://k6.io/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)

---

Good luck! Focus on quality testing and production readiness. The foundation is solid - now make it bulletproof.
