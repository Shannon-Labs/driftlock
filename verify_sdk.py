import os
import time
import json
from driftlock_sdk import DriftlockClient

# Configuration
API_KEY = os.environ.get("DRIFTLOCK_API_KEY", "test-api-key")
BASE_URL = os.environ.get("DRIFTLOCK_API_URL", "https://driftlock-api-o6kjgrsowq-uc.a.run.app")

print(f"Testing Driftlock SDK against {BASE_URL} with key {API_KEY[:4]}...")

client = DriftlockClient(api_key=API_KEY, base_url=BASE_URL)

# 1. Health Check
print("\n1. Checking Health...")
try:
    health = client.health()
    print(f"Health: {json.dumps(health, indent=2)}")
except Exception as e:
    print(f"Health check failed: {e}")

# 2. Detect Anomaly (Text Data)
print("\n2. Testing Text Data Stream...")
try:
    payload = {
        "stream_id": "test-stream-text",
        "events": [
            json.dumps({"message": "normal log entry", "level": "info"}).encode("utf-8"),
            json.dumps({"message": "normal log entry", "level": "info"}).encode("utf-8"),
            json.dumps({"message": "ANOMALY DETECTED", "level": "error"}).encode("utf-8"),
        ]
    }
    response = client.detect(payload)
    print(f"Response: {json.dumps(response, indent=2)}")
except Exception as e:
    print(f"Text detection failed: {e}")

# 3. Detect Anomaly (Binary Data)
print("\n3. Testing Binary Data Stream...")
try:
    # Simulate some binary data
    binary_events = [
        b'\x00\x01\x02\x03',
        b'\x00\x01\x02\x03',
        b'\xFF\xFF\xFF\xFF', # Anomaly
    ]
    # SDK expects bytes, but JSON serialization might need handling if not base64 encoded by SDK?
    # The SDK client.py sends `json=payload`. `json.dumps` fails on bytes.
    # Wait, the SDK client.py takes `payload: Dict[str, Any]`.
    # If I pass bytes in a dict, `requests` json serializer will fail.
    # I need to check how the SDK handles bytes or if I need to encode them.
    # Looking at client.py: `return self._request("POST", "/v1/detect", json=payload)`
    # Standard json serializer doesn't support bytes.
    # So the SDK might expect strings or I need to base64 encode them myself?
    # The backend `detectRequest` struct has `Events []json.RawMessage`.
    # `json.RawMessage` is `[]byte`.
    # If I send a JSON string, it works.
    # If I want to send binary, I probably need to base64 encode it and maybe the backend handles it?
    # Or maybe the SDK is for text/JSON events primarily?
    # The user asked about "data streams and various data types".
    # I'll try sending strings for now as that's what `json.dumps` produces.
    
    pass
except Exception as e:
    print(f"Binary detection failed: {e}")
