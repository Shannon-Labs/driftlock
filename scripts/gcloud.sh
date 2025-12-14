#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Codex/CI sandboxes may not allow writing to ~/.config/gcloud.
# Keep gcloud state in-repo unless the caller overrides CLOUDSDK_CONFIG.
export CLOUDSDK_CONFIG="${CLOUDSDK_CONFIG:-"$ROOT_DIR/.gcloud"}"
mkdir -p "$CLOUDSDK_CONFIG"

exec gcloud "$@"

