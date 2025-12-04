# Driftlock Python SDK (beta)

Thin wrapper around Driftlock’s `/v1` API. Keeps responses math-first (NCD, p-value, entropy) for auditability and leaves AI commentary server-side.

## Install

```bash
cd sdks/python
python -m venv .venv
source .venv/bin/activate
pip install -e .
```

## Usage

```python
from driftlock_sdk import DriftlockClient
import os

client = DriftlockClient(
    api_key=os.environ["DRIFTLOCK_API_KEY"],
    base_url=os.getenv("DRIFTLOCK_BASE_URL", "http://localhost:8080"),
)

# Health probe
health = client.health()
print(health.get("license"), health.get("database"))

# Detect
resp = client.detect({
    "stream_id": "integration-stream",
    "events": [{"timestamp": "2025-03-04T12:00:00Z", "type": "log", "body": {"msg": "demo"}}]
})
print("anomalies:", len(resp.get("anomalies", [])))

# Fetch anomaly detail
if resp.get("anomalies"):
    anomaly = client.get_anomaly(resp["anomalies"][0]["id"])
    print(anomaly.get("metrics"))
```

## Methods

- `health()` → `/healthz`
- `detect(payload: dict)` → `POST /v1/detect`
- `list_anomalies(params: dict)` → `GET /v1/anomalies`
- `get_anomaly(anomaly_id: str)` → `GET /v1/anomalies/{id}`

Raises `DriftlockError` on non-2xx responses or parse errors.

## Notes

- Default timeout: 10s. Override with `timeout`.
- Uses a shared `requests.Session` for connection reuse.
- No AI calls are made by default; this SDK only wraps the math-first endpoints. 
