#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
LANDING_DIR="$ROOT_DIR/landing-page"

if [[ ! -d "$LANDING_DIR" ]]; then
  echo "Landing page directory not found at $LANDING_DIR" >&2
  exit 1
fi

if [[ -z "${FIREBASE_TOKEN:-}" ]]; then
  echo "FIREBASE_TOKEN must be set in the environment." >&2
  exit 1
fi

pushd "$LANDING_DIR" >/dev/null

echo "Installing landing page dependencies..."
bun install --frozen-lockfile

echo "Building landing page..."
bun run build

DEPLOY_ARGS=(--only hosting)
if [[ -n "${FIREBASE_PROJECT:-}" ]]; then
  DEPLOY_ARGS+=("--project" "$FIREBASE_PROJECT")
fi

echo "Deploying to Firebase Hosting..."
npx firebase deploy "${DEPLOY_ARGS[@]}"

popd >/dev/null
echo "Landing page deployed successfully."
