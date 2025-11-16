#!/usr/bin/env bash
# Friendly wrapper around scripts/test-integration.sh that guides users through
# the HTTP API + Postgres demo flow.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INTEGRATION_SCRIPT="${ROOT_DIR}/scripts/test-integration.sh"
INFO_PREFIX="✨"
WARN_PREFIX="⚠️"
SUCCESS_PREFIX="✅"

KEEP_POSTGRES="false"
CUSTOM_PORT="${INTEGRATION_API_PORT:-}"

usage() {
  cat <<'EOF'
Usage: ./scripts/run-api-demo.sh [options]

Spin up Postgres, run the Driftlock HTTP API demo, and print follow-up commands
for exploring anomalies.

Options:
  --keep-postgres     Leave the dockerized Postgres instance running after the
                      demo completes (helpful if you want to keep querying).
  --port <number>     Force the API server to bind to a specific port.
  -h, --help          Show this message.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --keep-postgres)
      KEEP_POSTGRES="true"
      shift
      ;;
    --port)
      CUSTOM_PORT="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "${WARN_PREFIX} Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ ! -x "${INTEGRATION_SCRIPT}" ]]; then
  echo "${WARN_PREFIX} Unable to find ${INTEGRATION_SCRIPT}"
  exit 1
fi

SUMMARY_FILE="$(mktemp)"
trap 'rm -f "${SUMMARY_FILE}"' EXIT

if [[ -z "${DRIFTLOCK_LICENSE_KEY:-}" && -z "${DRIFTLOCK_DEV_MODE:-}" ]]; then
  export DRIFTLOCK_DEV_MODE="true"
  DEV_MODE_MESSAGE=" (dev mode enabled automatically)"
else
  DEV_MODE_MESSAGE=""
fi

cat <<EOF
${INFO_PREFIX} This script will:
   1. Build driftlock-http (Go + cbad-core FFI)
   2. Start Postgres via docker compose (and reuse it if already running)
   3. Apply migrations, create a tenant/key, and hit /v1/detect
   4. Show you commands to inspect the anomalies in the API + database${DEV_MODE_MESSAGE}

Dependencies: docker, go, jq, curl, psql, base64

EOF

ENV_VARS=(
  "INTEGRATION_SUMMARY_FILE=${SUMMARY_FILE}"
)

if [[ "${KEEP_POSTGRES}" == "true" ]]; then
  ENV_VARS+=("INTEGRATION_PRESERVE_POSTGRES=true")
fi

if [[ -n "${CUSTOM_PORT}" ]]; then
  ENV_VARS+=("INTEGRATION_API_PORT=${CUSTOM_PORT}")
fi

echo "${INFO_PREFIX} Running integration workflow..."
(cd "${ROOT_DIR}" && env USE_OPENZL=false "${ENV_VARS[@]}" "${INTEGRATION_SCRIPT}")

if [[ ! -s "${SUMMARY_FILE}" ]]; then
  echo "${WARN_PREFIX} Demo completed but summary file was empty."
  exit 1
fi

# shellcheck source=/dev/null
source "${SUMMARY_FILE}"

if [[ -z "${API_KEY:-}" || -z "${STREAM_ID:-}" ]]; then
  echo "${WARN_PREFIX} Missing API_KEY or STREAM_ID in summary output."
  exit 1
fi

DETECT_RESPONSE_JSON="$(printf '%s' "${DETECT_RESPONSE_B64}" | base64 --decode | jq '.')"

cat <<EOF

${SUCCESS_PREFIX} Driftlock API demo is ready.

Tenant stream: ${STREAM_ID}
API key (admin scope): ${API_KEY}
API base URL: ${BASE_URL}
Database URL: ${DATABASE_URL}
Detected anomalies: ${ANOMALY_COUNT}

Replay /v1/detect in another terminal:

  jq '.[0:600]' test-data/financial-demo.json \\
    | curl -sS -X POST \\
        -H "X-Api-Key: ${API_KEY}" \\
        -H 'Content-Type: application/json' \\
        --data @- \\
        ${BASE_URL}/v1/detect | jq '.anomalies[0]'

Inspect the anomaly we fetched during the run:

  curl -sS -H "X-Api-Key: ${API_KEY}" \\
    ${BASE_URL}/v1/anomalies/${FIRST_ANOMALY} | jq '{id, ncd, p_value, explanation}'

View persisted rows via psql:

  psql "${DATABASE_URL}" \\
    -c "SELECT id, stream_id, ROUND(ncd::numeric, 4) AS ncd, ROUND(p_value::numeric, 4) AS p_value FROM anomalies ORDER BY detected_at DESC LIMIT 5;"

Health probe (expects license + DB + queue status):

  curl -sS ${BASE_URL}/healthz | jq

To keep the Postgres container running for manual explorations, rerun this
script with --keep-postgres. Otherwise, docker compose will clean up any
containers it started once the script exits.

EOF

echo "${INFO_PREFIX} Sample detect response excerpt:"
echo "${DETECT_RESPONSE_JSON}" | jq '{anomaly_count, anomalies: [.anomalies[0]]}'
