ARG RUST_VERSION=1.82
ARG GO_VERSION=1.24
ARG GO_TARBALL=1.24.0
ARG DEBIAN_FRONTEND=noninteractive
ARG USE_OPENZL=false

FROM rust:${RUST_VERSION}-slim AS cbad
ARG USE_OPENZL
RUN apt-get update && apt-get install -y --no-install-recommends pkg-config clang make ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY cbad-core /src/cbad-core
COPY openzl /src/openzl
RUN if [ "${USE_OPENZL}" = "true" ]; then \
      echo "Building OpenZL library..."; \
      if [ ! -d "/src/openzl/deps/zstd" ]; then \
        echo "Warning: OpenZL nested submodule (zstd) not found. Ensure CI uses 'submodules: recursive' or run 'git submodule update --init --recursive' locally."; \
      fi; \
      make -C /src/openzl lib; \
    fi
WORKDIR /src/cbad-core
RUN if [ "${USE_OPENZL}" = "true" ]; then \
      cargo build --release --features openzl; \
    else \
      cargo build --release; \
    fi

FROM golang:${GO_VERSION}-bookworm AS builder
ARG GO_TARBALL
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-L/usr/local/lib -lcbad_core"
ENV LD_LIBRARY_PATH="/usr/local/lib:${LD_LIBRARY_PATH}"
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
COPY --from=cbad /src/cbad-core/target/release/libcbad_core.so /usr/local/lib/libcbad_core.so
COPY --from=cbad /src/cbad-core/target/release/libcbad_core.a /usr/local/lib/libcbad_core.a
COPY . /workspace
WORKDIR /workspace/collector-processor
RUN mkdir -p /workspace/bin && go build -o /workspace/bin/driftlock-http ./cmd/driftlock-http

FROM debian:bookworm-slim AS runtime
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && rm -rf /var/lib/apt/lists/*
ENV LD_LIBRARY_PATH="/usr/local/lib:${LD_LIBRARY_PATH}"
COPY --from=cbad /src/cbad-core/target/release/libcbad_core.so /usr/local/lib/libcbad_core.so
COPY --from=cbad /src/cbad-core/target/release/libcbad_core.a /usr/local/lib/libcbad_core.a
COPY --from=builder /workspace/bin/driftlock-http /usr/local/bin/driftlock-http
EXPOSE 8080
CMD ["/usr/local/bin/driftlock-http"]
