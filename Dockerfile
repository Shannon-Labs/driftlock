# Multi-stage build for Driftlock API
# Stage 1: Build Rust CBAD core
FROM rust:1.80 as rust-builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Copy Rust source
COPY cbad-core ./cbad-core
COPY openzl ./openzl

# Build CBAD core as cdylib (creates .so on Linux)
WORKDIR /build/cbad-core
RUN cargo build --release

# Verify the .so file exists and show what we have
RUN ls -lh /build/cbad-core/target/release/

# Stage 2: Build Go API server
FROM golang:1.24 as go-builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Copy entire project structure
COPY . .

# Copy Rust library from previous stage
COPY --from=rust-builder /build/cbad-core/target/release/*.so* /usr/local/lib/ || \
    COPY --from=rust-builder /build/cbad-core/target/release/*.dylib* /usr/local/lib/
RUN ldconfig /usr/local/lib || true

# Download dependencies
RUN go mod download

# Build API server with explicit CGO flags
WORKDIR /build/api-server
RUN CGO_ENABLED=1 CGO_CFLAGS="-I/usr/local/include" CGO_LDFLAGS="-L/usr/local/lib -lcbad_core" \
    go build -ldflags="-w -s" -o /build/driftlock-api ./cmd/api-server

# Stage 3: Final runtime image
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 65532 driftlock

WORKDIR /app

# Copy Go binary
COPY --from=go-builder /build/driftlock-api /app/driftlock-api

# Copy Rust library (try multiple possible names)
COPY --from=rust-builder /build/cbad-core/target/release/libcbad_core.so* /usr/local/lib/ || \
    COPY --from=rust-builder /build/cbad-core/target/release/libcbad_core.dylib* /usr/local/lib/ || \
    echo "No .so or .dylib found, checking what we have..."

# Update library cache
RUN ldconfig /usr/local/lib || true

# Create data directory
RUN mkdir -p /app/data && chown -R driftlock:driftlock /app

USER driftlock

ENV PORT=8080
ENV HOST=0.0.0.0
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

EXPOSE 8080

ENTRYPOINT ["/app/driftlock-api"]
