# Testing Strategies

Best practices for testing applications that integrate with Driftlock.

## Unit Testing

When running unit tests, you should mock the Driftlock client to avoid making real network requests and consuming your API quota.

### Node.js (Jest)

```javascript
jest.mock('@driftlock/client');
import { DriftlockClient } from '@driftlock/client';

test('should handle anomaly detection', async () => {
  const mockDetect = jest.fn().mockResolvedValue({
    success: true,
    anomalies: []
  });
  
  (DriftlockClient as jest.Mock).mockImplementation(() => ({
    detect: mockDetect
  }));

  // Run your code
  await processData();
  
  expect(mockDetect).toHaveBeenCalled();
});
```

### Python (Pytest)

```python
from unittest.mock import Mock, patch

@patch('driftlock.DriftlockClient')
def test_anomaly_handling(MockClient):
    client = MockClient.return_value
    client.detect.return_value = {'anomalies': []}
    
    # Run your code
    process_data()
    
    client.detect.assert_called_once()
```

## Integration Testing

For integration tests, you might want to verify that your application correctly communicates with Driftlock.

1. **Use a Test Stream**: Always use a dedicated `stream_id` (e.g., `ci-test`) for automated tests.
2. **Use a Test API Key**: Create a separate API key for CI/CD environments.

```javascript
// Integration test
test('should send event to driftlock', async () => {
  const client = new DriftlockClient({ apiKey: process.env.TEST_API_KEY });
  const result = await client.detect({
    streamId: 'ci-test',
    events: [{ type: 'test', body: { value: 1 } }]
  });
  
  expect(result.success).toBe(true);
});
```

## Load Testing

If you are load testing your application, you should consider how it interacts with Driftlock.

### Disabling Detection
You might want to disable Driftlock during load tests to isolate your application's performance, unless you are specifically testing the integration.

```javascript
const client = new DriftlockClient({
  apiKey: process.env.API_KEY,
  enabled: process.env.NODE_ENV !== 'load-test'
});
```

### Rate Limits
Remember that the Developer plan has a rate limit of 60 requests/minute. For load testing, ensure you are on a plan that supports your throughput, or mock the external call.

## Chaos Engineering

Test how your application behaves when Driftlock is unavailable or slow.

1. **Simulate Latency**: Configure your mock to delay responses.
2. **Simulate Errors**: Force your mock to return 500 errors or timeouts.

Ensure your application fails gracefully (e.g., logs the error but continues processing the user request).
