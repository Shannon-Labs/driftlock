#!/bin/bash
set -e

# Build frontend if not built
if [ ! -d "landing-page/dist" ]; then
  echo "Building frontend..."
  cd landing-page && npm ci && npm run build && cd ..
fi

# Build and run Rust API (serves static files from landing-page/dist)
echo "Starting Rust API server..."
cargo run -p driftlock-api --release
