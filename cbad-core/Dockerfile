# Multi-stage build for cbad_core with optional OpenZL feature
# Usage:
#   docker build -t cbad-core:base -f cbad-core/Dockerfile .
#   docker build -t cbad-core:openzl --build-arg USE_OPENZL=true -f cbad-core/Dockerfile .

ARG RUST_VERSION=1.82
ARG DEBIAN_FRONTEND=noninteractive

FROM rust:${RUST_VERSION}-slim AS builder

ARG USE_OPENZL=false

RUN apt-get update && \
    apt-get install -y --no-install-recommends pkg-config clang make ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache deps first
COPY cbad-core/Cargo.toml cbad-core/Cargo.lock ./cbad-core/
RUN mkdir -p cbad-core/src && \
    echo "pub fn _placeholder() {}" > cbad-core/src/lib.rs
WORKDIR /app/cbad-core
RUN if [ "${USE_OPENZL}" = "true" ]; then \
      cargo build --release --features openzl || true; \
    else \
      cargo build --release || true; \
    fi

# Now copy real sources and build
WORKDIR /app
COPY cbad-core ./cbad-core
WORKDIR /app/cbad-core
RUN if [ "${USE_OPENZL}" = "true" ]; then \
      echo "Building with OpenZL feature enabled"; \
      cargo build --release --features openzl; \
    else \
      echo "Building without OpenZL (default)"; \
      cargo build --release; \
    fi

# Runtime stage
FROM debian:bookworm-slim AS runtime
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/cbad
COPY --from=builder /app/cbad-core/target/release/libcbad_core.so /usr/local/lib/libcbad_core.so
COPY --from=builder /app/cbad-core/target/release/libcbad_core.a /usr/local/lib/libcbad_core.a
COPY --from=builder /app/cbad-core/target/release/libcbad_core.rlib /usr/local/lib/libcbad_core.rlib

ENV LD_LIBRARY_PATH=/usr/local/lib:${LD_LIBRARY_PATH}

# Provide a basic health check for presence of the library
RUN test -f /usr/local/lib/libcbad_core.so || test -f /usr/local/lib/libcbad_core.a

CMD ["/bin/sh", "-c", "echo 'cbad_core runtime image ready'; sleep infinity"]

