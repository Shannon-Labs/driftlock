#!/bin/bash
# CI verification script for CBAD build and integration

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prereqs() {
    echo_info "Checking prerequisites..."
    
    # Check Go version
    if ! command -v go &> /dev/null; then
        echo_error "Go is not installed"
        exit 1
    fi
    
    local go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
    echo_info "Go version: $go_version"
    
    # Check Rust/Cargo
    if ! command -v cargo &> /dev/null; then
        echo_error "Rust/Cargo is not installed"
        exit 1
    fi
    
    local rust_version=$(cargo --version | cut -d' ' -f2)
    echo_info "Rust version: $rust_version"
}

# Build CBAD core static library
build_cbad_core() {
    echo_info "Building CBAD core static library..."
    
    if [ ! -f "cbad-core/Cargo.toml" ]; then
        echo_error "cbad-core/Cargo.toml not found"
        exit 1
    fi
    
    # Build static library
    make cbad-core-lib
    
    # Verify library exists
    local lib_path="cbad-core/target/release/libcbad_core.a"
    if [ ! -f "$lib_path" ]; then
        echo_error "Static library not found at $lib_path"
        exit 1
    fi
    
    echo_info "CBAD core library built successfully"
    ls -la "$lib_path"
}

# Test Go integration with CBAD
test_go_integration() {
    echo_info "Testing Go integration with CBAD..."
    
    # Build collector with CBAD integration
    if ! make collector; then
        echo_error "Failed to build collector with CBAD integration"
        exit 1
    fi
    
    echo_info "Collector built successfully with CBAD integration"
}

# Run linting and formatting checks
run_linting() {
    echo_info "Running linting and formatting checks..."
    
    # Go formatting
    if ! gofmt -l . | grep -q '^$'; then
        echo_error "Go files are not properly formatted. Run 'gofmt -w .'"
        gofmt -l .
        exit 1
    fi
    
    # Go linting (if staticcheck is available)
    if command -v staticcheck &> /dev/null; then
        if ! staticcheck ./...; then
            echo_error "Go linting failed"
            exit 1
        fi
    else
        echo_warn "staticcheck not available, skipping Go linting"
    fi
    
    # Rust formatting
    if ! (cd cbad-core && cargo fmt -- --check); then
        echo_error "Rust files are not properly formatted. Run 'cargo fmt'"
        exit 1
    fi
    
    # Rust linting
    if ! (cd cbad-core && cargo clippy --all-targets -- -D warnings); then
        echo_error "Rust linting failed"
        exit 1
    fi
    
    echo_info "All linting checks passed"
}

# Run test suites
run_tests() {
    echo_info "Running test suites..."
    
    # Go tests
    if ! go test ./...; then
        echo_error "Go tests failed"
        exit 1
    fi
    
    # Rust tests
    if ! (cd cbad-core && cargo test); then
        echo_error "Rust tests failed"
        exit 1
    fi
    
    echo_info "All tests passed"
}

# Verify build tags and FFI integration
verify_ffi_integration() {
    echo_info "Verifying FFI integration..."
    
    # Default build must link against the Rust CBAD library
    if ! go build ./collector-processor/...; then
        echo_error "Failed to build collector with CBAD integration (CGO enabled build is required)"
        exit 1
    fi
    
    # Verify that the optional driftlock_no_cbad tag falls back to the stub
    if ! go build -tags driftlock_no_cbad ./collector-processor/...; then
        echo_error "Failed to build collector with driftlock_no_cbad stub tag"
        exit 1
    fi
    
    echo_info "FFI integration verified"
}

# Main execution
main() {
    echo_info "Starting CI verification for Driftlock CBAD build..."
    
    check_prereqs
    run_linting
    build_cbad_core
    test_go_integration
    verify_ffi_integration
    run_tests
    
    echo_info "All CI checks passed successfully!"
}

# Run main function
main "$@"
