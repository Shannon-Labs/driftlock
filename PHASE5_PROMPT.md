# Phase 5: Production Readiness & Testing

## Context

You are continuing development of **Driftlock**, a compression-based anomaly detection platform for OpenTelemetry data. 

**Current Status**: Phase 4 completed with all compilation errors fixed. The project is now in a stable state with:
- ✅ API server building successfully
- ✅ UI building with 0 errors/warnings  
- ✅ All Go code compiling
- ✅ All TypeScript/React code passing lint checks
- ⚠️  CBAD Rust library blocked by network restrictions (see BUGFIXES.md)

## Project Overview

Driftlock is an enterprise anomaly detection system that uses **compression-based analysis** to detect anomalies in telemetry streams:

### Core Components:
1. **API Server** (Go) - REST API with OTEL instrumentation, PostgreSQL storage, authentication
2. **UI** (Next.js 14 + TypeScript) - Modern dashboard with real-time SSE streams
3. **CBAD Core** (Rust) - Compression engine using OpenZL (format-aware compression)
4. **Collector Processor** - OpenTelemetry Collector integration
5. **Synthetic Tool** - OTLP test data generator

### Technology Stack:
- **Backend**: Go 1.22, PostgreSQL, Prometheus, OpenTelemetry
- **Frontend**: Next.js 14, React 18, TypeScript 5, Tailwind CSS
- **Compression**: Rust (OpenZL), Kolmogorov complexity theory
- **Infrastructure**: Docker, Helm, Grafana

## Phase 5 Objectives

Your goal is to make Driftlock **production-ready** by implementing comprehensive testing, improving reliability, and adding production features.

---

## Task Breakdown

### 1. Build CBAD Core Library ⚠️ CRITICAL

**Priority**: HIGHEST - This unblocks all other CBAD-dependent work

**Current Issue**: 
- Cannot build Rust cbad-core due to crates.io network access restrictions
- Error: `failed to get successful HTTP response from https://index.crates.io/config.json, got 403`

**Required Actions**:
a) **Verify Environment**:
   ```bash
   curl -I https://index.crates.io/config.json
   rustc --version
   cargo --version
   ```

b) **Attempt Build**:
   ```bash
   cd cbad-core
   cargo build --release --lib
   ```

c) **If Network Access Available**:
   - Build completes → Move to next task
   - Document build artifacts in `cbad-core/target/release/`

d) **If Still Blocked**:
   - Try alternative: `cargo vendor` to cache dependencies
   - Document limitation and skip CBAD-dependent tasks
   - Focus on API server tests instead

**Expected Output**: 
- `cbad-core/target/release/libcbad_core.a` (static library)
- All collector-processor tests passing

---

### 2. Add Comprehensive Unit Tests

**Goal**: Achieve >70% code coverage for critical paths

#### 2.1 API Server Tests

Create `api-server/internal/handlers/handlers_test.go`:

**Test Coverage Needed**:
```go
// Health endpoints
TestHealthHandler()
TestReadinessHandler()

// Event ingestion
TestEventsHandler_Success()
TestEventsHandler_InvalidPayload()
TestEventsHandler_AuthFailure()

// Anomaly endpoints
TestGetAnomalies_WithFilters()
TestGetAnomalies_Pagination()
TestGetAnomaly_NotFound()
TestGetAnomaly_Success()

// Status updates
TestUpdateAnomalyStatus_ValidTransition()
TestUpdateAnomalyStatus_InvalidStatus()

// Export functionality
TestExportAnomaly_Success()
TestExportAnomaly_NotFound()

// Stream endpoint
TestStreamAnomalies_SSE()

// Config endpoint
TestGetConfig_Success()
```

**Testing Patterns**:
- Use `httptest.NewRecorder()` for HTTP handlers
- Mock database with interfaces
- Use `testify/assert` for assertions
- Test error paths and edge cases

Example test structure:
```go
func TestEventsHandler_Success(t *testing.T) {
    // Setup
    req := httptest.NewRequest("POST", "/v1/events", strings.NewReader(`{"test":"data"}`))
    w := httptest.NewRecorder()
    
    // Execute
    handler := NewEventsHandler(mockDB)
    handler.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    // ... more assertions
}
```

#### 2.2 Storage Layer Tests

Create `api-server/internal/storage/postgres_test.go`:

**Test Coverage**:
- Database connection handling
- CRUD operations for anomalies
- Query filtering and pagination
- Transaction handling
- Error recovery

**Use Docker for Integration Tests**:
```bash
# Start test PostgreSQL
docker run --rm -p 5433:5432 -e POSTGRES_PASSWORD=test postgres:15
```

#### 2.3 Authentication Tests

Create `api-server/internal/auth/middleware_test.go`:

**Test Coverage**:
- API key validation
- JWT token validation (if implemented)
- Unauthorized access rejection
- Rate limiting (if implemented)

#### 2.4 CBAD Processor Tests (if library built)

Run existing tests:
```bash
cd collector-processor
go test ./... -v -tags driftlock_cbad_cgo
```

Enhance `collector-processor/driftlockcbad/cbad_test.go`:
- Add more compression scenarios
- Test anomaly detection accuracy
- Test different data patterns

---

### 3. Integration & End-to-End Tests

#### 3.1 E2E Test Suite

Create `tests/e2e/e2e_test.go`:

**Test Flow**:
1. Start API server
2. Generate synthetic OTLP data
3. POST events to `/v1/events`
4. Verify anomalies detected via `/v1/anomalies`
5. Test SSE stream `/v1/stream/anomalies`
6. Export anomaly evidence
7. Verify exported data structure

**Implementation**:
```bash
# Terminal 1: Start API server
make run

# Terminal 2: Run E2E tests
go test ./tests/e2e -v
```

#### 3.2 Synthetic Data Test

Test the synthetic tool works end-to-end:
```bash
# Build synthetic tool
make tools

# Generate normal data
bin/synthetic 100

# Generate anomalous data (create variant)
# Verify API server processes correctly
```

---

### 4. Production Hardening

#### 4.1 Add Request Validation

**File**: `api-server/internal/handlers/validation.go`

Implement:
- JSON schema validation for `/v1/events`
- Input sanitization
- Size limits (prevent DoS)
- Content-Type checking

Example:
```go
type EventRequest struct {
    Timestamp time.Time `json:"timestamp" validate:"required"`
    Value     float64   `json:"value" validate:"required"`
    Metadata  map[string]interface{} `json:"metadata"`
}

func ValidateEventRequest(r *EventRequest) error {
    if r.Timestamp.IsZero() {
        return errors.New("timestamp required")
    }
    // ... more validation
    return nil
}
```

#### 4.2 Add Rate Limiting

**File**: `api-server/internal/middleware/ratelimit.go`

Implement token bucket or sliding window:
```go
// Example: 100 requests per minute per IP
func RateLimitMiddleware(next http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

#### 4.3 Enhance Error Handling

**Current**: Many handlers use generic errors  
**Needed**: Structured error responses

Create `api-server/internal/errors/errors.go`:
```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
    return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

var (
    ErrNotFound = &APIError{404, "Resource not found", ""}
    ErrUnauthorized = &APIError{401, "Unauthorized", ""}
    ErrInvalidInput = &APIError{400, "Invalid input", ""}
    // ... more errors
)
```

#### 4.4 Add Structured Logging

Replace `log.Printf` with structured logging:

**File**: `api-server/internal/logging/logger.go`

Use `zap` or `slog`:
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("event ingested",
    "event_id", eventID,
    "stream_type", streamType,
    "timestamp", timestamp)
```

#### 4.5 Add Request Tracing

Enhance OpenTelemetry instrumentation:
```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("api-server").Start(r.Context(), "HandleEvent")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("http.method", r.Method),
        attribute.String("http.path", r.URL.Path),
    )
    
    // ... handler logic
}
```

---

### 5. Database Integration

#### 5.1 Implement PostgreSQL Storage

**Current**: Mock storage in API handlers  
**Needed**: Real PostgreSQL implementation

**File**: `api-server/internal/storage/postgres.go`

**Schema** (create `api-server/migrations/001_initial.sql`):
```sql
CREATE TABLE anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL,
    stream_type VARCHAR(50) NOT NULL,
    ncd_score FLOAT NOT NULL,
    p_value FLOAT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    glass_box_explanation TEXT NOT NULL,
    compression_baseline FLOAT,
    compression_window FLOAT,
    compression_combined FLOAT,
    baseline_data JSONB,
    window_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_anomalies_timestamp ON anomalies(timestamp DESC);
CREATE INDEX idx_anomalies_status ON anomalies(status);
CREATE INDEX idx_anomalies_stream_type ON anomalies(stream_type);
CREATE INDEX idx_anomalies_p_value ON anomalies(p_value);
```

**Implementation**:
```go
type PostgresStorage struct {
    db *sql.DB
}

func (s *PostgresStorage) SaveAnomaly(ctx context.Context, a *Anomaly) error {
    query := `INSERT INTO anomalies (...) VALUES (...)`
    _, err := s.db.ExecContext(ctx, query, ...)
    return err
}

func (s *PostgresStorage) GetAnomalies(ctx context.Context, filters AnomalyFilters) ([]Anomaly, error) {
    // Build dynamic query with filters
    query := buildQuery(filters)
    rows, err := s.db.QueryContext(ctx, query)
    // ... parse rows
    return anomalies, err
}
```

#### 5.2 Add Database Migrations

Use `golang-migrate` or similar:
```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migrations directory
mkdir -p api-server/migrations

# Run migrations
migrate -path api-server/migrations -database "postgres://localhost:5432/driftlock?sslmode=disable" up
```

---

### 6. Documentation

#### 6.1 API Documentation

Create `docs/API.md`:

Document all endpoints with:
- Method, path, description
- Request/response schemas
- Authentication requirements
- Example curl commands
- Error codes

Example:
```markdown
### POST /v1/events

Ingest OpenTelemetry events for anomaly detection.

**Authentication**: Required (API key)

**Request Body**:
```json
{
  "timestamp": "2025-10-25T12:00:00Z",
  "value": 42.5,
  "metadata": {
    "service": "api-gateway",
    "region": "us-east-1"
  }
}
```

**Response**: 
- `200 OK`: Event accepted
- `400 Bad Request`: Invalid payload
- `401 Unauthorized`: Missing/invalid API key
- `429 Too Many Requests`: Rate limit exceeded
```

#### 6.2 Deployment Guide

Create `docs/DEPLOYMENT.md`:

Cover:
- Environment variables configuration
- Database setup
- Docker deployment
- Kubernetes/Helm deployment
- Monitoring setup (Prometheus, Grafana)
- Troubleshooting

#### 6.3 Developer Guide

Create `docs/DEVELOPMENT.md`:

Cover:
- Local development setup
- Running tests
- Code style guidelines
- Contribution workflow
- Architecture overview

---

### 7. Performance & Scalability

#### 7.1 Add Benchmarks

Create `api-server/internal/handlers/benchmark_test.go`:
```go
func BenchmarkEventsHandler(b *testing.B) {
    handler := NewEventsHandler(mockDB)
    payload := []byte(`{"timestamp":"2025-10-25T12:00:00Z","value":42.5}`)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("POST", "/v1/events", bytes.NewReader(payload))
        w := httptest.NewRecorder()
        handler.ServeHTTP(w, req)
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem ./api-server/internal/handlers
```

#### 7.2 Add Connection Pooling

Configure PostgreSQL connection pool:
```go
db, err := sql.Open("postgres", connStr)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 7.3 Add Caching (Optional)

For expensive queries, add Redis caching:
```go
// Cache anomaly counts
func (s *Storage) GetAnomalyCount(ctx context.Context) (int, error) {
    // Check cache first
    if cached, err := s.redis.Get(ctx, "anomaly:count").Result(); err == nil {
        return strconv.Atoi(cached)
    }
    
    // Query database
    count, err := s.queryCount(ctx)
    
    // Cache result
    s.redis.Set(ctx, "anomaly:count", count, 5*time.Minute)
    return count, err
}
```

---

### 8. CI/CD Pipeline (Optional)

#### 8.1 GitHub Actions Workflow

Create `.github/workflows/ci.yml`:
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
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Set up Rust
        uses: actions-rs/toolchain@v1
        with:
          toolchain: stable
      
      - name: Build CBAD
        run: make cbad-core-lib
      
      - name: Run tests
        run: make test
      
      - name: Build API
        run: make api
      
      - name: Lint
        run: |
          cd ui
          npm install
          npm run lint
```

---

## Success Criteria

By the end of Phase 5, the project should have:

### ✅ Testing
- [ ] Unit tests for all API handlers (>70% coverage)
- [ ] Storage layer integration tests
- [ ] Authentication/authorization tests
- [ ] E2E test suite covering main workflows
- [ ] All tests passing in CI

### ✅ Production Features
- [ ] Request validation on all endpoints
- [ ] Rate limiting implemented
- [ ] Structured error responses
- [ ] Structured logging with proper levels
- [ ] OpenTelemetry tracing complete

### ✅ Database
- [ ] PostgreSQL schema created
- [ ] Migrations system in place
- [ ] All CRUD operations implemented
- [ ] Indexes optimized for queries

### ✅ Documentation
- [ ] Complete API documentation
- [ ] Deployment guide
- [ ] Developer guide
- [ ] Architecture diagrams

### ✅ Performance
- [ ] Benchmarks written and baseline established
- [ ] Connection pooling configured
- [ ] Load tested (1000+ req/s if possible)

### ✅ Quality
- [ ] No linting errors
- [ ] All compilation warnings resolved
- [ ] Code formatted consistently
- [ ] Security audit completed

---

## Implementation Order (Recommended)

1. **CBAD Build** (if network available) - Unblocks everything
2. **Unit Tests** - Foundation for quality
3. **Database Integration** - Core functionality
4. **Production Hardening** - Security & reliability
5. **E2E Tests** - Validation
6. **Documentation** - Usability
7. **Performance** - Optimization
8. **CI/CD** - Automation

---

## Known Issues to Address

From BUGFIXES.md:

1. **CBAD Network Restriction**: Try to resolve or document workaround
2. **Missing Root Tests**: Add tests for root module packages
3. **Mock Data in Production Code**: Replace with real storage

---

## Tips for Success

1. **Start Small**: Begin with one test file, verify it works, then expand
2. **Use Mocks**: Don't require full database for every test
3. **Test Errors**: Error paths are often untested but critical
4. **Document As You Go**: Don't save docs for last
5. **Commit Frequently**: Small, focused commits are easier to review
6. **Run CI Locally**: Fix issues before pushing

---

## Questions to Consider

- Should we implement JWT authentication or stick with API keys?
- Do we need WebSocket support in addition to SSE?
- Should anomalies be auto-acknowledged after N days?
- Do we need audit logging for all mutations?
- Should we support bulk event ingestion?

---

## Deliverables

At completion, provide:

1. **Test Report**: Coverage stats, test counts, any failures
2. **Performance Report**: Benchmark results, load test outcomes
3. **Documentation**: Links to all docs created
4. **Known Issues**: Any remaining bugs or limitations
5. **Phase 6 Recommendations**: What should come next

---

## Getting Started

```bash
# 1. Verify current state
git status
go build ./...
cd ui && npm run build

# 2. Review BUGFIXES.md
cat BUGFIXES.md

# 3. Try CBAD build
cd cbad-core
cargo build --release --lib

# 4. If CBAD works, start with tests
cd ../api-server
mkdir -p internal/handlers/
touch internal/handlers/handlers_test.go

# 5. Write first test
# ... implement

# 6. Run test
go test ./internal/handlers -v
```

---

## Resources

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Testify Documentation](https://github.com/stretchr/testify)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [PostgreSQL Go Driver](https://github.com/lib/pq)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)

---

Good luck! Focus on quality over quantity. A well-tested, production-ready API is worth more than a feature-rich but buggy one.

