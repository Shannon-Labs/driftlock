# Driftlock Development Guide

This document provides a comprehensive guide for developers working on the Driftlock project.

## Table of Contents
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Development Environment Setup](#development-environment-setup)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Contribution Workflow](#contribution-workflow)
- [Project Structure](#project-structure)

## Project Overview

Driftlock is an explainable, deterministic anomaly detection toolkit designed for regulated and audit-conscious teams. It leverages compression-based anomaly detection (CBAD) to identify deviations from expected patterns in telemetry data.

### Core Features
- **Deterministic**: 100% reproducible results across runs
- **Explainable**: Glass-box explanations for detected anomalies
- **Compliance-Ready**: Built-in DORA, NIS2, and Runtime AI compliance
- **Real-time**: Sub-second anomaly detection and alerting
- **Privacy-First**: On-premises deployment with configurable data redaction

### Technology Stack
- **CBAD Engine**: Rust with C FFI, WASM compilation target
- **API Layer**: Go with PostgreSQL, Redis caching
- **Frontend**: Next.js with TypeScript, real-time WebSocket integration
- **Infrastructure**: Kubernetes, Istio service mesh, Prometheus monitoring

## Architecture

### High-Level Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   OTel Collector│    │   Driftlock API │    │  UI Dashboard   │
│                 │    │                 │    │                 │
│  - Receives     │───▶│  - Anomaly      │───▶│  - Visualize    │
│    telemetry    │    │    detection    │    │    anomalies    │
│  - Routes to    │    │  - Storage      │    │  - Manage       │
│    CBAD         │    │    interface    │    │    alerts       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │                 │
                       │  - Anomaly      │
                       │    storage      │
                       │  - Detection    │
                       │    config       │
                       └─────────────────┘
```

### Data Flow
1. OTel Collector receives telemetry data (logs, metrics, traces)
2. Data is routed to CBAD engine for analysis
3. Compression ratios and statistical measures computed
4. Anomalies detected based on NCD and p-value thresholds
5. Anomaly details stored in PostgreSQL
6. Real-time notifications sent via SSE/webhook
7. UI dashboard displays anomalies with explanations

## Development Environment Setup

### Prerequisites
- Go 1.22+
- Rust 1.70+
- Docker and Docker Compose (or Colima for macOS)
- PostgreSQL 13+ (or Docker for testing)
- Node.js 18+ (for UI development)

### Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/Hmbown/driftlock.git
cd driftlock
```

2. **Set up environment variables:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Build the CBAD core (Rust):**
```bash
cd cbad-core
cargo build --release
```

4. **Build the API server:**
```bash
cd api-server
go build -o driftlock-api cmd/driftlock-api/main.go
```

5. **Start dependencies with Docker:**
```bash
docker-compose up -d postgres redis
```

6. **Run database migrations:**
```bash
go run cmd/migrate/main.go
```

7. **Start the API server:**
```bash
./driftlock-api
```

### Running with Docker Compose
For a complete local development setup:

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Stop services
docker-compose -f docker-compose.dev.yml down
```

## Running Tests

### API Server Tests
```bash
# Run all tests
cd api-server
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestAnomalyCRUD ./internal/handlers/
```

### Integration Tests
```bash
# Run integration tests (requires PostgreSQL)
go test -tags=integration ./internal/storage/...
```

### End-to-End Tests
```bash
# Run E2E tests (requires test environment)
cd tests/e2e
go test ./...
```

### Load Tests
```bash
# Run load tests using k6
cd tests/load
k6 run load_test.js
```

## Code Style

### Go Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofumpt` for formatting
- Use `golangci-lint` for linting
- Write comprehensive tests for all business logic

### Import Grouping
```go
import (
	// Standard library
	"context"
	"fmt"
	"net/http"
	
	// Third-party
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	
	// Internal
	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/storage"
)
```

### Error Handling
- Use the `errors` package for consistent error handling
- Always wrap errors with context using `fmt.Errorf("context: %w", err)`
- Return structured API errors using `errors.APIError`

### Naming Conventions
- Use descriptive names: `userID` instead of `uid`
- Be consistent with terminology throughout the codebase
- Prefix exported types, functions, and variables with meaningful names

## Contribution Workflow

### Branching Strategy
- Use feature branches off `develop`
- Name branches descriptively: `feature/anomaly-detection-improvements` or `bugfix/api-rate-limiting`
- Submit pull requests to `develop`

### Pull Request Process
1. Create an issue describing the feature or bug
2. Fork the repository and create a feature branch
3. Make your changes following the code style
4. Write or update tests as appropriate
5. Update documentation if needed
6. Submit a pull request with a clear description
7. Address review comments
8. Your PR will be merged after approval

### Commit Messages
- Use imperative mood: "Add feature" not "Added feature"
- First line should be 50 characters or less
- Include a body if needed to explain the "why"
- Link to issues when applicable

Example:
```
Add support for custom anomaly severity thresholds

Enables users to configure different severity levels based on
NCD scores and p-values. This allows for more granular alerting
strategies in different environments.

Resolves #123
```

## Project Structure

```
driftlock/
├── api-server/                 # Go API server
│   ├── cmd/                   # Main applications
│   ├── internal/              # Private application code
│   │   ├── api/              # HTTP handlers
│   │   ├── auth/             # Authentication
│   │   ├── errors/           # Error handling
│   │   ├── export/           # Evidence export logic
│   │   ├── logging/          # Logging utilities
│   │   ├── metrics/          # Prometheus metrics
│   │   ├── middleware/       # HTTP middleware
│   │   ├── models/           # Data models
│   │   └── storage/          # Database operations
│   ├── migrations/           # Database migrations
│   └── go.mod, go.sum        # Go module files
├── cbad-core/                # Rust CBAD engine
├── collector-processor/      # OTel collector processor
├── ui/                       # Frontend application
├── docs/                     # Documentation
├── tests/                    # Test files
├── k8s/                      # Kubernetes manifests
├── helm/                     # Helm charts
├── .github/                  # GitHub configuration
└── Dockerfile*               # Container definitions
```

### Key Directories
- `api-server/internal/api`: Contains the main HTTP router and handler registration
- `api-server/internal/storage`: Database access layer with PostgreSQL implementation
- `api-server/internal/auth`: Authentication and authorization logic
- `api-server/internal/export`: Evidence bundle generation and export functionality
- `tests/e2e`: End-to-end tests simulating real-world usage
- `tests/load`: Performance and load testing scripts
- `k8s/`: Production Kubernetes deployment manifests
- `helm/driftlock/`: Helm chart for easier deployment

## Performance Considerations

### API Server
- Use connection pooling for database access
- Implement proper request timeouts
- Cache frequently accessed configuration
- Use streaming responses for large datasets

### Database
- Create proper indexes for common query patterns
- Use prepared statements to prevent SQL injection
- Implement read replicas for read-heavy operations
- Set up proper connection pooling

### CBAD Engine
- Optimize compression algorithms for speed
- Implement efficient data buffering strategies
- Use memory-mapped files for large datasets when appropriate
- Consider SIMD instructions for mathematical operations

## Security Best Practices

### Input Validation
- Implement strict input validation on all endpoints
- Use parameterized queries to prevent SQL injection
- Sanitize and validate all user inputs
- Implement proper rate limiting

### Authentication & Authorization
- Use JWT tokens with proper expiration
- Implement role-based access control (RBAC)
- Use HTTPS in production
- Store secrets securely (not in code)

### Data Protection
- Implement data redaction for sensitive information
- Use encryption at rest for sensitive data
- Limit data retention according to compliance requirements
- Log access to sensitive data for audit purposes

## Troubleshooting

### Common Issues
- **Database connection failures**: Check connection string and ensure PostgreSQL is running
- **CBAD build failures**: Verify Rust installation and network access to crates.io
- **API rate limiting**: Check rate limiter configuration and adjust as needed
- **Memory issues**: Monitor memory usage and tune GC settings if needed

### Debugging Tips
- Set `LOG_LEVEL=debug` for detailed logging
- Use pprof for performance profiling: `go tool pprof http://localhost:8080/debug/pprof/profile`
- Check PostgreSQL logs for query performance issues
- Monitor metrics via the `/metrics` endpoint

## Contact and Support

- For questions about the codebase: [team@driftlock.example.com]
- For security issues: [security@driftlock.example.com]
- For general support: Open an issue in the GitHub repository

For more detailed information, see the individual README files in each major component directory.