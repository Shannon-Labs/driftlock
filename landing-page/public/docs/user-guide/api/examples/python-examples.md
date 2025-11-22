# Python Examples

Complete Python examples for integrating with Driftlock's API.

## Installation

```bash
pip install requests
```

For async support:
```bash
pip install httpx  # Modern async HTTP client
```

## Basic Client

### Simple Client Class

```python
import requests
from typing import Dict, List, Any, Optional
import time

class DriftlockClient:
    """Simple Driftlock API client"""
    
    def __init__(self, api_key: str, base_url: str = "https://driftlock-api-o6kjgrsowq-uc.a.run.app"):
        self.api_key = api_key
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({
            'X-Api-Key': api_key,
            'Content-Type': 'application/json'
        })
    
    def detect(self, events: List[Dict], stream_id: str = "default", 
               config_override: Optional[Dict] = None) -> Dict:
        """Run anomaly detection on events"""
        payload = {
            "stream_id": stream_id,
            "events": events
        }
        if config_override:
            payload["config_override"] = config_override
        
        response = self.session.post(f"{self.base_url}/v1/detect", json=payload)
        response.raise_for_status()
        return response.json()
    
    def list_anomalies(self, stream_id: Optional[str] = None, 
                      min_ncd: Optional[float] = None,
                      limit: int = 50) -> Dict:
        """List anomalies with optional filtering"""
        params = {"limit": limit}
        if stream_id:
            params["stream_id"] = stream_id
        if min_ncd is not None:
            params["min_ncd"] = min_ncd
        
        response = self.session.get(f"{self.base_url}/v1/anomalies", params=params)
        response.raise_for_status()
        return response.json()
    
    def get_anomaly(self, anomaly_id: str) -> Dict:
        """Get detailed information about a specific anomaly"""
        response = self.session.get(f"{self.base_url}/v1/anomalies/{anomaly_id}")
        response.raise_for_status()
        return response.json()
```

## Usage Examples

### Detect Anomalies

```python
from datetime import datetime, timedelta
import os

# Initialize client
api_key = os.getenv("DRIFTLOCK_API_KEY")
client = DriftlockClient(api_key)

# Prepare events
events = [
    {
        "timestamp": datetime.now().isoformat(),
        "type": "log",
        "body": {"message": "Normal login", "user": "alice"}
    },
    {
        "timestamp": (datetime.now() + timedelta(minutes=1)).isoformat(),
        "type": "log",
        "body": {"message": "Normal login", "user": "bob"}
    },
    {
        "timestamp": (datetime.now() + timedelta(minutes=2)).isoformat(),
        "type": "log",
        "body": {"message": "SQL INJECTION ATTEMPT!", "user": "hacker"}
    }
]

# Run detection
result = client.detect(events)

print(f"Processed {result['total_events']} events")
print(f"Found {result['anomaly_count']} anomalies")

# Print anomalies
for anomaly in result['anomalies']:
    print(f"\\nAnomaly {anomaly['id']}:")
    print(f"  NCD: {anomaly['metrics']['ncd']:.3f}")
    print(f"  Confidence: {anomaly['metrics']['confidence']:.3f}")
    print(f"  Explanation: {anomaly['why']}")
```

### Metric Monitoring

```python
# Monitor system metrics
def monitor_metrics(client: DriftlockClient):
    import psutil
    
    events = []
    for i in range(10):
        cpu = psutil.cpu_percent(interval=1)
        memory = psutil.virtual_memory().percent
        
        events.append({
            "timestamp": datetime.now().isoformat(),
            "type": "metric",
            "body": {
                "cpu_percent": cpu,
                "memory_percent": memory
            },
            "attributes": {
                "host": "server-01",
                "environment": "production"
            }
        })
        
        time.sleep(1)
    
    result = client.detect(events, stream_id="system-metrics")
    return result

# Run monitoring
result = monitor_metrics(client)
```

### List and Filter Anomalies

```python
# Get all recent anomalies
anomalies = client.list_anomalies(limit=100)
print(f"Total anomalies: {anomalies['total']}")

# Filter by stream
production_anomalies = client.list_anomalies(
    stream_id="production-logs",
    min_ncd=0.5
)

# Get detailed info
if anomalies['anomalies']:
    first_anomaly_id = anomalies['anomalies'][0]['id']
    details = client.get_anomaly(first_anomaly_id)
    print(f"Anomaly details: {details}")
```

## Advanced Client

### With Retry Logic

```python
from requests.adapters import HTTPAdapter
from requests.packages.urllib3.util.retry import Retry

class RobustDriftlockClient(DriftlockClient):
    """Driftlock client with automatic retries"""
    
    def __init__(self, api_key: str, **kwargs):
        super().__init__(api_key, **kwargs)
        
        # Configure retries
        retry_strategy = Retry(
            total=3,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["GET", "POST"],
            backoff_factor=1  # Wait 1, 2, 4 seconds
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("https://", adapter)
        self.session.mount("http://", adapter)
```

### With Rate Limit Handling

```python
class DriftlockClientWithRateLimit(DriftlockClient):
    """Client with rate limit handling"""
    
    def _make_request(self, method: str, endpoint: str, **kwargs) -> Dict:
        """Make request with rate limit handling"""
        max_retries = 5
        
        for attempt in range(max_retries):
            try:
                response = self.session.request(method, f"{self.base_url}{endpoint}", **kwargs)
                
                # Check rate limit
                if response.status_code == 429:
                    retry_after = int(response.headers.get('Retry-After', 60))
                    print(f"Rate limited. Retrying after {retry_after}s...")
                    time.sleep(retry_after)
                    continue
                
                response.raise_for_status()
                return response.json()
                
            except requests.exceptions.HTTPError as e:
                if attempt == max_retries - 1:
                    raise
                time.sleep(2 ** attempt)  # Exponential backoff
        
        raise Exception("Max retries exceeded")
    
    def detect(self, events: List[Dict], **kwargs) -> Dict:
        return self._make_request('POST', '/v1/detect', json={
            "events": events,
            **kwargs
        })
```

## Async Client

### Using httpx

```python
import httpx
import asyncio
from typing import List, Dict

class AsyncDriftlockClient:
    """Async Driftlock API client"""
    
    def __init__(self, api_key: str, base_url: str = "https://driftlock-api-o6kjgrsowq-uc.a.run.app"):
        self.api_key = api_key
        self.base_url = base_url
        self.headers = {
            'X-Api-Key': api_key,
            'Content-Type': 'application/json'
        }
    
    async def detect(self, events: List[Dict], stream_id: str = "default") -> Dict:
        """Async anomaly detection"""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/v1/detect",
                json={"stream_id": stream_id, "events": events},
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()
    
    async def list_anomalies(self, **params) -> Dict:
        """Async list anomalies"""
        async with httpx.AsyncClient() as client:
            response = await client.get(
                f"{self.base_url}/v1/anomalies",
                params=params,
                headers=self.headers
            )
            response.raise_for_status()
            return response.json()

# Usage
async def main():
    client = AsyncDriftlockClient(os.getenv("DRIFTLOCK_API_KEY"))
    
    # Run multiple detections concurrently
    tasks = [
        client.detect(events1, stream_id="stream1"),
        client.detect(events2, stream_id="stream2"),
        client.detect(events3, stream_id="stream3")
    ]
    
    results = await asyncio.gather(*tasks)
    for result in results:
        print(f"Found {result['anomaly_count']} anomalies")

# Run
asyncio.run(main())
```

## Real-World Examples

### Log File Monitoring

```python
import json

def monitor_log_file(client: DriftlockClient, log_file: str, batch_size: int = 50):
    """Monitor log file for anomalies"""
    events = []
    
    with open(log_file, 'r') as f:
        for line in f:
            try:
                log_entry = json.loads(line)
                events.append({
                    "timestamp": log_entry.get("timestamp", datetime.now().isoformat()),
                    "type": "log",
                    "body": log_entry,
                    "attributes": {
                        "source": log_file
                    }
                })
                
                # Batch detection
                if len(events) >= batch_size:
                    result = client.detect(events, stream_id="log-monitoring")
                    if result['anomaly_count'] > 0:
                        print(f"🚨 Found {result['anomaly_count']} anomalies!")
                        for anomaly in result['anomalies']:
                            print(f"  - {anomaly['why']}")
                    events = []
                    
            except json.JSONDecodeError:
                continue
    
    # Process remaining events
    if events:
        result = client.detect(events, stream_id="log-monitoring")

# Usage
monitor_log_file(client, "/var/log/application.log")
```

### Financial Transaction Monitoring

```python
def monitor_transactions(client: DriftlockClient, transactions: List[Dict]):
    """Monitor financial transactions for anomalies"""
    events = [
        {
            "timestamp": tx['timestamp'],
            "type": "llm",  # Use llm type for structured financial data
            "body": {
                "amount": tx['amount'],
                "currency": tx['currency'],
                "merchant": tx['merchant'],
                "card_type": tx['card_type']
            },
            "attributes": {
                "user_id": tx['user_id'],
                "country": tx['country']
            },
            "idempotency_key": f"tx_{tx['id']}"
        }
        for tx in transactions
    ]
    
    # Use stricter thresholds for financial data
    result = client.detect(
        events,
        stream_id="financial-transactions",
        config_override={
            "ncd_threshold": 0.25,  # More sensitive
            "p_value_threshold": 0.01  # Higher significance
        }
    )
    
    # Alert on anomalies
    if result['anomaly_count'] > 0:
        for anomaly in result['anomalies']:
            send_alert(f"Suspicious transaction detected", anomaly)
    
    return result
```

## Error Handling

```python
from requests.exceptions import HTTPError, RequestException

try:
    result = client.detect(events)
except HTTPError as e:
    if e.response.status_code == 401:
        print("Invalid API key")
    elif e.response.status_code == 429:
        error_data = e.response.json()
        retry_after = error_data['error'].get('retry_after_seconds', 60)
        print(f"Rate limited. Wait {retry_after} seconds")
    else:
        print(f"HTTP error: {e}")
except RequestException as e:
    print(f"Network error: {e}")
```

## Testing

```python
import unittest
from unittest.mock import Mock, patch

class TestDriftlockClient(unittest.TestCase):
    
    def setUp(self):
        self.client = DriftlockClient("test_api_key")
    
    @patch('requests.Session.post')
    def test_detect(self, mock_post):
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "success": True,
            "anomaly_count": 1
        }
        mock_post.return_value = mock_response
        
        # Test
        result = self.client.detect([{"type": "log", "body": {}}])
        self.assertEqual(result['anomaly_count'], 1)

if __name__ == '__main__':
    unittest.main()
```

## Next Steps

- **[Node.js Examples](./node-examples.md)** - JavaScript/TypeScript client
- **[cURL Examples](./curl-examples.md)** - Command-line examples
- **[API Reference](../rest-api.md)** - Complete API documentation

---

**Tip**: Check out the [official Python SDK](https://github.com/Shannon-Labs/driftlock-python) (coming soon) for a fully-featured client!
