#!/usr/bin/env bash
# End-to-end integration test for the Driftlock HTTP API.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${ROOT_DIR}/bin"
API_BIN="${BIN_DIR}/driftlock-http"
LOG_FILE=""
SERVER_PID=""
POSTGRES_SERVICE="driftlock-postgres"
POSTGRES_ALREADY_RUNNING="false"
SUMMARY_FILE="${INTEGRATION_SUMMARY_FILE:-}"
TEMP_FILES=()
INFO_PREFIX="🔧"
SUCCESS_PREFIX="✅"
FAIL_PREFIX="❌"

cleanup() {
  if [[ -n "${SERVER_PID}" ]]; then
    kill "${SERVER_PID}" >/dev/null 2>&1 || true
    wait "${SERVER_PID}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${LOG_FILE}" && -f "${LOG_FILE}" ]]; then
    rm -f "${LOG_FILE}"
  fi
  for file in "${TEMP_FILES[@]}"; do
    if [[ -f "${file}" ]]; then
      rm -f "${file}"
    fi
  done
  if [[ "${POSTGRES_ALREADY_RUNNING}" != "true" && "${INTEGRATION_PRESERVE_POSTGRES:-}" != "true" ]]; then
    docker compose -f "${ROOT_DIR}/docker-compose.yml" rm -sf "${POSTGRES_SERVICE}" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "${FAIL_PREFIX} Missing dependency: $1"
    exit 1
  fi
}

echo "${INFO_PREFIX} Checking prerequisites..."
for cmd in docker jq curl go psql base64; do
  require_command "${cmd}"
done

# Use dev mode if no license key is provided
if [[ -z "${DRIFTLOCK_LICENSE_KEY:-}" ]]; then
  if [[ -z "${DRIFTLOCK_DEV_MODE:-}" ]]; then
    echo "${INFO_PREFIX} No DRIFTLOCK_LICENSE_KEY set, enabling development mode..."
    export DRIFTLOCK_DEV_MODE="true"
  fi
fi

CBAD_LIB_DIR="${ROOT_DIR}/cbad-core/target/release"
if [[ ! -d "${CBAD_LIB_DIR}" ]] || [[ -z "$(ls "${CBAD_LIB_DIR}"/libcbad_core.* 2>/dev/null)" ]]; then
  echo "${FAIL_PREFIX} cbad-core artifacts not found. Run 'cargo build --release' first."
  exit 1
fi

export LD_LIBRARY_PATH="${CBAD_LIB_DIR}:${LD_LIBRARY_PATH:-}"
export DYLD_LIBRARY_PATH="${CBAD_LIB_DIR}:${DYLD_LIBRARY_PATH:-}"
export USE_OPENZL=false

DB_PORT="${POSTGRES_PORT:-7543}"
export DATABASE_URL="${DATABASE_URL:-postgres://driftlock:driftlock@localhost:${DB_PORT}/driftlock?sslmode=disable}"
if [[ -z "${INTEGRATION_API_PORT:-}" ]]; then
  if command -v python3 >/dev/null 2>&1; then
    API_PORT=$(python3 - <<'PY'
import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind(("127.0.0.1", 0))
port = s.getsockname()[1]
s.close()
print(port)
PY
)
  else
    API_PORT=18080
  fi
else
  API_PORT="${INTEGRATION_API_PORT}"
fi
WRITE_TIMEOUT_SEC="${INTEGRATION_WRITE_TIMEOUT_SEC:-180}"
BASE_URL="http://localhost:${API_PORT}"

echo "${INFO_PREFIX} Building driftlock-http..."
mkdir -p "${BIN_DIR}"
pushd "${ROOT_DIR}/collector-processor/cmd/driftlock-http" >/dev/null
go build -o "${API_BIN}" .
popd >/dev/null

EXISTING_POSTGRES_ID=$(docker compose -f "${ROOT_DIR}/docker-compose.yml" ps -q "${POSTGRES_SERVICE}" 2>/dev/null || true)
if [[ -n "${EXISTING_POSTGRES_ID}" ]]; then
  POSTGRES_RUNNING_STATE=$(docker inspect -f '{{.State.Running}}' "${EXISTING_POSTGRES_ID}" 2>/dev/null || echo "false")
  if [[ "${POSTGRES_RUNNING_STATE}" == "true" ]]; then
    POSTGRES_ALREADY_RUNNING="true"
    echo "${INFO_PREFIX} Reusing running Postgres container (${EXISTING_POSTGRES_ID})"
  fi
fi

if [[ "${POSTGRES_ALREADY_RUNNING}" != "true" ]]; then
  echo "${INFO_PREFIX} Starting Postgres container..."
  docker compose -f "${ROOT_DIR}/docker-compose.yml" up -d "${POSTGRES_SERVICE}"
fi

echo "${INFO_PREFIX} Waiting for Postgres to accept connections..."
for _ in {1..30}; do
  if PGPASSWORD=driftlock psql "${DATABASE_URL}" -c '\q' >/dev/null 2>&1; then
    break
  fi
  sleep 1
done
if ! PGPASSWORD=driftlock psql "${DATABASE_URL}" -c '\q' >/dev/null 2>&1; then
  echo "${FAIL_PREFIX} Postgres did not become ready in time."
  exit 1
fi

echo "${INFO_PREFIX} Running goose migrations..."
"${API_BIN}" migrate up

echo "${INFO_PREFIX} Creating integration tenant via CLI..."
SUFFIX=$(date +%s)
TENANT_JSON=$("${API_BIN}" create-tenant \
  --name "Integration Tenant ${SUFFIX}" \
  --slug "integration-${SUFFIX}" \
  --stream "integration-stream-${SUFFIX}" \
  --key-role "admin" \
  --key-name "integration-cli" \
  --json)

API_KEY=$(echo "${TENANT_JSON}" | jq -r '.api_key')
STREAM_ID=$(echo "${TENANT_JSON}" | jq -r '.stream_id')
if [[ -z "${API_KEY}" || -z "${STREAM_ID}" || "${API_KEY}" == "null" ]]; then
  echo "${FAIL_PREFIX} Failed to parse tenant creation response:"
  echo "${TENANT_JSON}"
  exit 1
fi

echo "${INFO_PREFIX} Starting API server on ${BASE_URL}..."
LOG_FILE="$(mktemp)"
(
  cd "${ROOT_DIR}"
  PORT="${API_PORT}" \
  WRITE_TIMEOUT_SEC="${WRITE_TIMEOUT_SEC}" \
  DATABASE_URL="${DATABASE_URL}" \
  DRIFTLOCK_LICENSE_KEY="${DRIFTLOCK_LICENSE_KEY:-}" \
  DRIFTLOCK_DEV_MODE="${DRIFTLOCK_DEV_MODE:-}" \
  LD_LIBRARY_PATH="${LD_LIBRARY_PATH}" \
  DYLD_LIBRARY_PATH="${DYLD_LIBRARY_PATH}" \
  "${API_BIN}" >"${LOG_FILE}" 2>&1
) &
SERVER_PID=$!
sleep 1

echo "${INFO_PREFIX} Waiting for /healthz..."
for _ in {1..30}; do
  if curl -sSf "${BASE_URL}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done
if ! curl -sSf "${BASE_URL}/healthz" >/dev/null 2>&1; then
  echo "${FAIL_PREFIX} API server failed to start. Logs:"
  cat "${LOG_FILE}"
  exit 1
fi

HEALTH=$(curl -s "${BASE_URL}/healthz")
LICENSE_STATUS=$(echo "${HEALTH}" | jq -r '.license.status')
DB_STATUS=$(echo "${HEALTH}" | jq -r '.database')
if [[ "${LICENSE_STATUS}" != "valid" || "${DB_STATUS}" != "connected" ]]; then
  echo "${FAIL_PREFIX} Unexpected /healthz response:"
  echo "${HEALTH}"
  exit 1
fi
echo "${SUCCESS_PREFIX} Health check reports license=${LICENSE_STATUS}, db=${DB_STATUS}"

EVENT_FILE=$(mktemp)
TEMP_FILES+=("${EVENT_FILE}")
jq '.[0:600]' "${ROOT_DIR}/test-data/financial-demo.json" > "${EVENT_FILE}"
DETECT_PAYLOAD=$(mktemp)
TEMP_FILES+=("${DETECT_PAYLOAD}")
REQUEST_ID=$(uuidgen 2>/dev/null || date +%s)
jq --arg stream "${STREAM_ID}" --arg req "${REQUEST_ID}" '{stream_id:$stream, request_id:$req, events:.}' "${EVENT_FILE}" > "${DETECT_PAYLOAD}"
DETECT_RESPONSE_FILE=$(mktemp)
TEMP_FILES+=("${DETECT_RESPONSE_FILE}")

echo "${INFO_PREFIX} Posting /v1/detect payload..."
DETECT_RESPONSE=$(curl -s \
  -H "X-Api-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  --data-binary "@${DETECT_PAYLOAD}" \
  "${BASE_URL}/v1/detect" || true)
printf '%s' "${DETECT_RESPONSE}" > "${DETECT_RESPONSE_FILE}"

SUCCESS=$(echo "${DETECT_RESPONSE}" | jq -r '.success')
ANOMALY_COUNT=$(echo "${DETECT_RESPONSE}" | jq -r '.anomaly_count')
if [[ "${SUCCESS}" != "true" ]]; then
  echo "${FAIL_PREFIX} /v1/detect failed:"
  echo "${DETECT_RESPONSE}"
  if [[ -f "${LOG_FILE}" ]]; then
    echo "--- driftlock-http log ---"
    cat "${LOG_FILE}"
    echo "-------------------------"
  fi
  exit 1
fi
if [[ "${ANOMALY_COUNT}" -lt 1 ]]; then
  echo "${FAIL_PREFIX} Expected anomalies but found ${ANOMALY_COUNT}"
  exit 1
fi
echo "${SUCCESS_PREFIX} /v1/detect produced ${ANOMALY_COUNT} anomalies"

echo "${INFO_PREFIX} Verifying database persistence..."
BATCHES=$(psql "${DATABASE_URL}" -Atc "SELECT COUNT(*) FROM ingest_batches")
ANOMALIES=$(psql "${DATABASE_URL}" -Atc "SELECT COUNT(*) FROM anomalies")
if [[ "${BATCHES}" -lt 1 || "${ANOMALIES}" -lt 1 ]]; then
  echo "${FAIL_PREFIX} Expected persisted rows (batches=${BATCHES}, anomalies=${ANOMALIES})"
  exit 1
fi
echo "${SUCCESS_PREFIX} Database rows present (batches=${BATCHES}, anomalies=${ANOMALIES})"

FIRST_ANOMALY=$(echo "${DETECT_RESPONSE}" | jq -r '.anomalies[0].id // empty')
if [[ -n "${FIRST_ANOMALY}" ]]; then
  echo "${INFO_PREFIX} Fetching anomaly detail..."
  DETAIL=$(curl -s -H "X-Api-Key: ${API_KEY}" "${BASE_URL}/v1/anomalies/${FIRST_ANOMALY}" || true)
  DETAIL_STATUS=$(echo "${DETAIL}" | jq -r '.status // empty')
  if [[ -z "${DETAIL_STATUS}" ]]; then
    echo "${FAIL_PREFIX} Detail endpoint did not return expected payload."
    echo "${DETAIL}"
    exit 1
  fi
  echo "${SUCCESS_PREFIX} Anomaly detail retrieved (status=${DETAIL_STATUS})"
fi

echo "${INFO_PREFIX} Exercising export stubs..."
EXPORT_BULK=$(curl -s -o /dev/null -w "%{http_code}" \
  -X POST -H "X-Api-Key: ${API_KEY}" \
  -H "Content-Type: application/json" \
  --data '{"format":"json"}' \
  "${BASE_URL}/v1/anomalies/export")
if [[ "${EXPORT_BULK}" != "202" ]]; then
  echo "${FAIL_PREFIX} Bulk export endpoint returned HTTP ${EXPORT_BULK}"
  exit 1
fi
if [[ -n "${FIRST_ANOMALY}" ]]; then
  EXPORT_SINGLE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST -H "X-Api-Key: ${API_KEY}" \
    -H "Content-Type: application/json" \
    "${BASE_URL}/v1/anomalies/${FIRST_ANOMALY}/export")
  if [[ "${EXPORT_SINGLE}" != "202" ]]; then
    echo "${FAIL_PREFIX} Single export endpoint returned HTTP ${EXPORT_SINGLE}"
    exit 1
  fi
fi
echo "${SUCCESS_PREFIX} Export jobs accepted (stubbed)"

echo "${SUCCESS_PREFIX} Integration test completed successfully."

if [[ -n "${SUMMARY_FILE}" ]]; then
  DETECT_RESPONSE_B64=$(printf '%s' "${DETECT_RESPONSE}" | base64 | tr -d '\n')
  cat <<EOF > "${SUMMARY_FILE}"
API_KEY=${API_KEY}
STREAM_ID=${STREAM_ID}
BASE_URL=${BASE_URL}
DATABASE_URL=${DATABASE_URL}
FIRST_ANOMALY=${FIRST_ANOMALY}
ANOMALY_COUNT=${ANOMALY_COUNT}
REQUEST_ID=${REQUEST_ID}
API_PORT=${API_PORT}
WRITE_TIMEOUT_SEC=${WRITE_TIMEOUT_SEC}
DETECT_RESPONSE_B64=${DETECT_RESPONSE_B64}
EOF
fi
