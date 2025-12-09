# Python Examples

Minimal Python examples for Driftlock's REST API using `requests`.

## Install

```bash
pip install requests
```

## Constants

```python
import os
import requests

API_URL = os.getenv("DRIFTLOCK_API_URL", "https://api.driftlock.net/v1")
API_KEY = os.environ["DRIFTLOCK_API_KEY"]
HEADERS = {"X-Api-Key": API_KEY, "Content-Type": "application/json"}
```

## Detect anomalies

```python
def detect(events, stream_id="default", config_override=None):
    payload = {"events": events, "stream_id": stream_id}
    if config_override:
        payload["config_override"] = config_override

    resp = requests.post(f"{API_URL}/detect", headers=HEADERS, json=payload, timeout=15)
    resp.raise_for_status()
    return resp.json()

result = detect([
    {"timestamp": "2025-01-01T10:00:00Z", "type": "metric", "body": {"latency_ms": 120}},
    {"timestamp": "2025-01-01T10:01:00Z", "type": "metric", "body": {"latency_ms": 950}},
])

print(f"Anomalies: {result['anomaly_count']}")
for anom in result.get("anomalies", []):
    print(f"Index {anom['index']}: {anom['why']} (NCD={anom['metrics']['ncd']})")
```

## Detection with overrides

```python
result = detect(
    events=[{"type": "log", "body": {"message": "SQL injection attempt"}}],
    stream_id="prod-logs",
    config_override={"ncd_threshold": 0.35, "p_value_threshold": 0.02, "compressor": "lz4"},
)
```

## List anomalies

```python
def list_anomalies(stream_id=None, min_ncd=None, limit=50):
    params = {"limit": limit}
    if stream_id:
        params["stream_id"] = stream_id
    if min_ncd is not None:
        params["min_ncd"] = min_ncd

    resp = requests.get(f"{API_URL}/anomalies", headers=HEADERS, params=params, timeout=10)
    resp.raise_for_status()
    return resp.json()

anomalies = list_anomalies(stream_id="prod-logs", min_ncd=0.5)
print(f"Total anomalies: {anomalies['total']}")
```

## Basic retry on 429

```python
import time

def detect_with_retry(events, retries=3):
    for attempt in range(retries):
        resp = requests.post(f"{API_URL}/detect", headers=HEADERS, json={"events": events}, timeout=15)
        if resp.status_code == 429:
            retry_after = int(resp.headers.get("Retry-After", 60))
            time.sleep(retry_after)
            continue
        resp.raise_for_status()
        return resp.json()
    raise RuntimeError("Rate limit: max retries exceeded")
```

## Demo endpoint (no auth)

```python
demo_resp = requests.post(
    f"{API_URL}/demo/detect",
    headers={"Content-Type": "application/json"},
    json={"events": [{"body": {"cpu": 45}}, {"body": {"cpu": 99}}]},
)
print(demo_resp.json())
```

**Next:** [REST API reference](../rest-api.md)
