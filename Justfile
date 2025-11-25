# Justfile for Driftlock

set shell := ["bash", "-c"]

# List available commands
default:
    @just --list

# --- Setup ---

# Install dependencies
setup:
    npm install
    cd landing-page && npm install
    cd functions && npm install
    cd extensions/vscode-driftlock && npm install
    @echo "Dependencies installed."

# --- Build ---

# Build all components
build: build-core build-demo build-landing build-functions build-extension

# Build Rust core
build-core:
    cd cbad-core && cargo build --release

# Build Rust core with OpenZL (requires OpenZL library built first)
build-core-openzl: build-openzl-lib
    cd cbad-core && OPENZL_LIB_DIR="$(pwd)/../openzl" cargo build --release --features openzl

# Build Go demo
build-demo: build-core
    go build -o driftlock-demo cmd/demo/main.go

# Build Landing Page
build-landing:
    cd landing-page && npm run build

# Build Functions
build-functions:
    cd functions && npm run build

# Build VS Code Extension
build-extension:
    cd extensions/vscode-driftlock && npm run compile

# Build Docker image
docker-build:
    docker build -t driftlock-http:dev -f collector-processor/cmd/driftlock-http/Dockerfile .

# Build Docker image with OpenZL
docker-build-openzl:
    docker build -t driftlock-http:openzl -f collector-processor/cmd/driftlock-http/Dockerfile --build-arg USE_OPENZL=true .

# Smoke test Docker builds
docker-test:
    ./scripts/test-docker-build.sh

# --- Test ---

# Run all tests
test: test-core test-go test-landing test-functions test-extension

# Test Rust core
test-core:
    cd cbad-core && cargo test

# Test Go components
test-go: build-core
    for mod in $(find . -name go.mod -not -path "*/vendor/*"); do \
        dir=$(dirname $mod); \
        echo "Testing $dir..."; \
        (cd $dir && go test ./...); \
    done

# Test Landing Page
test-landing:
    cd landing-page && npm run type-check

# Test Functions
test-functions:
    cd functions && npm run build

# Test VS Code Extension
test-extension:
    cd extensions/vscode-driftlock && npm test

# Run verification script
verify: build-demo
    ./scripts/verify-launch-readiness.sh

# --- Lint ---

# Lint all code
lint: lint-go lint-rust lint-landing lint-functions lint-extension

# Lint Go
lint-go: build-core
    for mod in $(find . -name go.mod -not -path "*/vendor/*"); do \
        dir=$(dirname $mod); \
        echo "Linting $dir..."; \
        (cd $dir && go vet ./...); \
    done

# Lint Rust
lint-rust:
    cd cbad-core && cargo clippy

# Lint Landing Page
lint-landing:
    cd landing-page && npm run lint

# Lint Functions
lint-functions:
    cd functions && npm run lint

# Lint Extension
lint-extension:
    cd extensions/vscode-driftlock && npm run lint

# --- Security ---

# Run security scans (govulncheck, cargo audit, npm audit)
security: security-go security-rust security-npm
    @echo "✅ All security scans complete."

# Go vulnerability check
security-go:
    @echo "🔍 Running govulncheck on Go modules..."
    @command -v govulncheck >/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
    for mod in $(find . -name go.mod -not -path "*/vendor/*" -not -path "*/.archive/*"); do \
        dir=$$(dirname $$mod); \
        echo "  Scanning $$dir..."; \
        (cd $$dir && govulncheck ./...) || true; \
    done

# Rust vulnerability check
security-rust:
    @echo "🔍 Running cargo audit on Rust crates..."
    @command -v cargo-audit >/dev/null || cargo install cargo-audit
    cd cbad-core && cargo audit

# npm audit across all JS/TS projects
security-npm:
    @echo "🔍 Running npm audit..."
    cd landing-page && npm audit --audit-level=high || true
    cd functions && npm audit --audit-level=high || true
    cd extensions/vscode-driftlock && npm audit --audit-level=high || true

# --- Deploy ---

# Deploy to Firebase
deploy:
    firebase deploy

# Deploy Landing Page only
deploy-landing:
    ./scripts/deploy-landing.sh

# --- Dev ---

# Run Landing Page dev server
dev-landing:
    cd landing-page && npm run dev

# --- Utilities ---

# Clean build artifacts
clean:
    rm -rf driftlock-demo
    cd cbad-core && cargo clean
    rm -rf landing-page/dist
    rm -rf functions/lib
    rm -rf extensions/vscode-driftlock/out
    @echo "🧹 Build artifacts cleaned."

# Format all code
fmt: fmt-rust fmt-go
    @echo "✅ Code formatted."

# Format Rust code
fmt-rust:
    cd cbad-core && cargo fmt

# Format Go code
fmt-go:
    for mod in $(find . -name go.mod -not -path "*/vendor/*"); do \
        dir=$$(dirname $$mod); \
        echo "Formatting $$dir..."; \
        (cd $$dir && go fmt ./...); \
    done

# --- OpenZL (Optional Format-Aware Compression) ---

# Build OpenZL static library (prerequisites: C/C++ toolchain, make)
build-openzl-lib:
    @echo "🔧 Building OpenZL library..."
    @if [ ! -d "openzl/src" ]; then \
        echo "❌ OpenZL submodule not initialized. Run: git submodule update --init --recursive"; \
        exit 1; \
    fi
    cd openzl && CFLAGS="-fPIC $${CFLAGS:-}" CXXFLAGS="-fPIC $${CXXFLAGS:-}" make lib
    @echo "✅ OpenZL library built: openzl/libopenzl.a"

# Test cbad-core with OpenZL feature enabled
test-core-openzl: build-openzl-lib
    @echo "🧪 Testing cbad-core with OpenZL..."
    cd cbad-core && OPENZL_LIB_DIR="$(pwd)/../openzl" cargo test --features openzl

# Full OpenZL build and test (Rust + Go integration)
test-openzl: build-openzl-lib test-core-openzl
    @echo "🧪 Testing Go with OpenZL-enabled cbad-core..."
    cd collector-processor && \
        LD_LIBRARY_PATH="$(pwd)/../cbad-core/target/release" \
        CGO_LDFLAGS="-L$(pwd)/../cbad-core/target/release -lcbad_core" \
        go test ./...
    @echo "✅ OpenZL integration tests passed!"

# Clean OpenZL build artifacts
clean-openzl:
    @echo "🧹 Cleaning OpenZL artifacts..."
    cd openzl && make clean 2>/dev/null || true
    @echo "✅ OpenZL cleaned."

# Show OpenZL status
openzl-status:
    @echo "OpenZL Integration Status:"
    @echo "=========================="
    @if [ -f "openzl/libopenzl.a" ]; then \
        echo "✅ OpenZL library: BUILT"; \
        ls -lh openzl/libopenzl.a; \
    else \
        echo "❌ OpenZL library: NOT BUILT (run 'just build-openzl-lib')"; \
    fi
    @if [ -f "cbad-core/target/release/libcbad_core.so" ] || [ -f "cbad-core/target/release/libcbad_core.dylib" ]; then \
        echo "✅ cbad-core library: BUILT"; \
    else \
        echo "❌ cbad-core library: NOT BUILT (run 'just build-core')"; \
    fi
    @echo ""
    @echo "To build with OpenZL: just build-core-openzl"
    @echo "To test with OpenZL:  just test-openzl"
