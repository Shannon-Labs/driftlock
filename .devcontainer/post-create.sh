#!/usr/bin/env bash
set -euo pipefail

echo "Installing workspace tooling (just, Google Cloud SDK, Firebase CLI, Docker CLI)..."

sudo apt-get update
sudo apt-get install -y \
  apt-transport-https \
  ca-certificates \
  curl \
  gnupg \
  lsb-release \
  build-essential \
  pkg-config \
  docker.io \
  just

# Google Cloud SDK (idempotent install)
if [ ! -f /etc/apt/sources.list.d/google-cloud-sdk.list ]; then
  echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list >/dev/null
  curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo tee /usr/share/keyrings/cloud.google.gpg >/dev/null
fi

sudo apt-get update
sudo apt-get install -y google-cloud-cli

# Ensure Rust tooling used by clippy/rustfmt is present
if command -v rustup >/dev/null 2>&1; then
  rustup component add clippy rustfmt >/dev/null
fi

# Firebase CLI via npm (after Node feature install)
if ! command -v firebase >/dev/null 2>&1; then
  npm install -g firebase-tools
fi

echo "Workspace tooling ready. 'just setup' has been run; rerun it if dependencies change."
