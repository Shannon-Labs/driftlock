#!/bin/bash
# Launch readiness sweep for Driftlock SaaS (builds, configs, and live health).
# Keeps checks lightweight and fast; fails on missing hard requirements, warns on nice-to-haves.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PASS=0
WARN=0
FAIL=0

status() {
  local level="$1"; shift
  local message="$*"
  case "${level}" in
    OK)   PASS=$((PASS+1)); echo "[OK]   ${message}" ;;
    WARN) WARN=$((WARN+1)); echo "[WARN] ${message}" ;;
    FAIL) FAIL=$((FAIL+1)); echo "[FAIL] ${message}" ;;
  esac
}

need_cmd() {
  local cmd="$1"
  if command -v "${cmd}" >/dev/null 2>&1; then
    status OK "Command available: ${cmd}"
  else
    status FAIL "Missing required command: ${cmd}"
  fi
}

check_file() {
  local path="$1"
  local severity="${2:-FAIL}"
  if [[ -e "${path}" ]]; then
    status OK "Found ${path}"
  else
    status "${severity}" "Missing ${path}"
  fi
}

echo "Launch Readiness Check"
echo "----------------------"

need_cmd docker
need_cmd go
need_cmd node
need_cmd jq
need_cmd curl

if command -v firebase >/dev/null 2>&1; then
  status OK "Command available: firebase"
else
  status WARN "firebase CLI not found (needed for deploy)"
fi

check_file "${ROOT}/cloudbuild.yaml"
check_file "${ROOT}/firebase.json"
check_file "${ROOT}/api/migrations/20250302000000_onboarding.sql"
check_file "${ROOT}/docker-compose.yml"
check_file "${ROOT}/collector-processor/cmd/driftlock-http/main.go"
check_file "${ROOT}/cbad-core/Cargo.toml"

shopt -s nullglob
CBAD_ARTIFACTS=("${ROOT}/cbad-core/target/release"/libcbad_core.*)
shopt -u nullglob
if [[ ${#CBAD_ARTIFACTS[@]} -gt 0 ]]; then
  status OK "cbad-core release artifacts present"
else
  status WARN "cbad-core artifacts missing (run: cargo build --release)"
fi

if [[ -d "${ROOT}/landing-page/dist" ]]; then
  status OK "Frontend build present (landing-page/dist)"
else
  status WARN "Frontend build missing (run: cd landing-page && npm run build)"
fi

if [[ -f "${ROOT}/landing-page/.env.production" ]]; then
  status OK "Landing page env present (.env.production)"
else
  status WARN "Landing page env missing (.env.production)"
fi

if [[ -f "${ROOT}/functions/.env.local" ]]; then
  status OK "Functions env present (functions/.env.local)"
else
  status WARN "Functions env missing (functions/.env.local)"
fi

# Optional live health check (passive if API not running)
if curl -fsS http://localhost:8080/healthz >/dev/null 2>&1; then
  status OK "API /healthz reachable on localhost:8080"
else
  status WARN "API not reachable at http://localhost:8080/healthz (start via ./scripts/run-api-demo.sh)"
fi

echo "----------------------"
echo "Summary: ${PASS} OK, ${WARN} WARN, ${FAIL} FAIL"
if [[ "${FAIL}" -gt 0 ]]; then
  echo "Result: NOT READY (fix FAIL items above)."
  exit 1
fi
echo "Result: READY WITH WARNINGS=${WARN}"
