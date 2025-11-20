# Python SDK

The official Python client for Driftlock. Designed for data science and backend applications with full async support.

## Installation

Install the package via pip:

```bash
pip install driftlock
```

Or using poetry:

```bash
poetry add driftlock
```

## Quick Start

Initialize the client and start detecting anomalies.

```python
import os
import asyncio
from driftlock import DriftlockClient

# Initialize the client
client = DriftlockClient(api_key=os.getenv("DRIFTLOCK_API_KEY"))

async def main():
    # Detect anomalies
    result = await client.detect(
        stream_id="sensor-data",
        events=[
            {
                "timestamp": "2025-01-01T12:00:00Z",
                "type": "temperature",
                "body": {"value": 25.5, "unit": "C"}
            }
        ]
    )

    if result.anomalies:
        print(f"Detected {len(result.anomalies)} anomalies:")
        for anomaly in result.anomalies:
            print(f"- {anomaly.why} (Confidence: {anomaly.metrics.confidence})")
    else:
        print("No anomalies detected.")

if __name__ == "__main__":
    asyncio.run(main())
```

## Synchronous Usage

If you prefer synchronous code, use `DriftlockSyncClient`.

```python
from driftlock import DriftlockSyncClient

client = DriftlockSyncClient(api_key="your-api-key")

result = client.detect(
    stream_id="sync-stream",
    events=[...]
)
```

## Configuration

The client accepts several configuration options:

```python
client = DriftlockClient(
    api_key="your-api-key",
    base_url="https://driftlock-api-o6kjgrsowq-uc.a.run.app",
    timeout=10.0,  # seconds
    max_retries=3
)
```

## Pandas Integration

The Python SDK integrates seamlessly with Pandas DataFrames.

```python
import pandas as pd

df = pd.DataFrame({
    'timestamp': [...],
    'value': [10, 12, 11, 100, 10],
    'category': ['A', 'A', 'A', 'A', 'A']
})

# Convert DataFrame to Driftlock events
events = client.utils.from_dataframe(
    df, 
    timestamp_col='timestamp',
    body_cols=['value', 'category']
)

result = await client.detect(stream_id="pandas-stream", events=events)
```

## Error Handling

Driftlock exceptions are available in `driftlock.exceptions`.

```python
from driftlock.exceptions import DriftlockAPIError, RateLimitError

try:
    await client.detect(...)
except RateLimitError as e:
    print(f"Rate limit exceeded. Retry in {e.retry_after} seconds.")
except DriftlockAPIError as e:
    print(f"API Error: {e.message}")
```

## Logging

Enable debug logging to see request details.

```python
import logging

logging.basicConfig(level=logging.DEBUG)
logging.getLogger("driftlock").setLevel(logging.DEBUG)
```

## Support

For issues and feature requests, visit our [GitHub repository](https://github.com/Shannon-Labs/driftlock-python) or email [support@driftlock.io](mailto:support@driftlock.io).
