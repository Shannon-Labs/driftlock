# Build stage
FROM rust:1.83-slim-bookworm AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy Cargo files for dependency caching
COPY Cargo.toml Cargo.lock ./
COPY cbad-core/Cargo.toml cbad-core/
COPY crates/driftlock-api/Cargo.toml crates/driftlock-api/
COPY crates/driftlock-db/Cargo.toml crates/driftlock-db/
COPY crates/driftlock-auth/Cargo.toml crates/driftlock-auth/
COPY crates/driftlock-billing/Cargo.toml crates/driftlock-billing/
COPY crates/driftlock-email/Cargo.toml crates/driftlock-email/

# Create dummy source files for dependency caching
RUN mkdir -p cbad-core/src && echo "fn main() {}" > cbad-core/src/lib.rs
RUN mkdir -p crates/driftlock-api/src && echo "fn main() {}" > crates/driftlock-api/src/main.rs
RUN mkdir -p crates/driftlock-db/src && echo "" > crates/driftlock-db/src/lib.rs
RUN mkdir -p crates/driftlock-auth/src && echo "" > crates/driftlock-auth/src/lib.rs
RUN mkdir -p crates/driftlock-billing/src && echo "" > crates/driftlock-billing/src/lib.rs
RUN mkdir -p crates/driftlock-email/src && echo "" > crates/driftlock-email/src/lib.rs

# Build dependencies (cached layer)
RUN cargo build --release -p driftlock-api || true

# Copy actual source code
COPY cbad-core cbad-core/
COPY crates crates/

# Touch files to ensure rebuild
RUN touch cbad-core/src/lib.rs crates/driftlock-api/src/main.rs

# Build the actual binary
RUN cargo build --release -p driftlock-api

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary
COPY --from=builder /app/target/release/driftlock-api /usr/local/bin/

# Create non-root user
RUN useradd -r -s /bin/false driftlock
USER driftlock

# Default port
ENV PORT=8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

CMD ["driftlock-api"]
