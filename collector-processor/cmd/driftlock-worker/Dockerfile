# Driftlock Worker Dockerfile
# Batch processing service for Anthropic Batch API integration

ARG GO_VERSION=1.24
ARG GO_TARBALL=1.24.0
ARG DEBIAN_FRONTEND=noninteractive

FROM golang:${GO_VERSION}-bookworm AS builder
ARG GO_TARBALL

WORKDIR /workspace

# Fallback: if Go binary appears broken/empty, install official Go toolchain
RUN if [ ! -s /usr/local/go/bin/go ]; then \
      set -eux; \
      echo "Go binary missing or empty, installing official Go ${GO_TARBALL}..."; \
      apt-get update && apt-get install -y --no-install-recommends ca-certificates curl tar && \
      rm -rf /var/lib/apt/lists/*; \
      arch="$(uname -m)"; \
      case "$arch" in \
        aarch64|arm64) GOARCH=arm64 ;; \
        x86_64|amd64) GOARCH=amd64 ;; \
        *) echo "Unsupported arch: $arch"; exit 1 ;; \
      esac; \
      curl -fsSL "https://dl.google.com/go/go${GO_TARBALL}.linux-${GOARCH}.tar.gz" -o /tmp/go.tgz; \
      rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tgz; \
      ln -sf /usr/local/go/bin/go /usr/local/bin/go; \
      ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt; \
      /usr/local/go/bin/go version; \
    fi

# Copy source code
COPY . /workspace

# Build the worker binary
WORKDIR /workspace/collector-processor
RUN mkdir -p /workspace/bin && \
    go build -o /workspace/bin/driftlock-worker ./cmd/driftlock-worker

# Runtime stage - minimal image
FROM debian:bookworm-slim AS runtime
ARG DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary
COPY --from=builder /workspace/bin/driftlock-worker /usr/local/bin/driftlock-worker

# Run as non-root user for security
RUN useradd -r -s /bin/false driftlock && \
    chown driftlock:driftlock /usr/local/bin/driftlock-worker

USER driftlock

# No exposed ports - this is a background worker
CMD ["/usr/local/bin/driftlock-worker"]
