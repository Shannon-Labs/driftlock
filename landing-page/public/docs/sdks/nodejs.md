# Node.js SDK

The official Node.js client for Driftlock. Integrate anomaly detection into your Node.js applications with full TypeScript support.

## Installation

Install the package via npm:

```bash
npm install @driftlock/client
```

Or using yarn:

```bash
yarn add @driftlock/client
```

## Quick Start

Initialize the client with your API key and start detecting anomalies.

```typescript
import { DriftlockClient } from '@driftlock/client';

// Initialize the client
const client = new DriftlockClient({
  apiKey: process.env.DRIFTLOCK_API_KEY,
});

// Detect anomalies in a stream of events
async function detectAnomalies() {
  try {
    const result = await client.detect({
      streamId: 'payment-processing',
      events: [
        {
          timestamp: new Date().toISOString(),
          type: 'transaction',
          body: { amount: 500, currency: 'USD', latency: 120 }
        }
      ]
    });

    if (result.anomalies.length > 0) {
      console.log('Anomalies detected:', result.anomalies);
    } else {
      console.log('No anomalies detected.');
    }
  } catch (error) {
    console.error('Error detecting anomalies:', error);
  }
}

detectAnomalies();
```

## Configuration

The `DriftlockClient` constructor accepts the following options:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `apiKey` | `string` | Required | Your Driftlock API key |
| `baseUrl` | `string` | `https://driftlock-api-o6kjgrsowq-uc.a.run.app` | API base URL |
| `timeout` | `number` | `10000` | Request timeout in milliseconds |
| `retries` | `number` | `3` | Number of retries for failed requests |

```typescript
const client = new DriftlockClient({
  apiKey: 'your-api-key',
  timeout: 5000,
  retries: 5
});
```

## Advanced Usage

### Batch Processing

For high-throughput applications, you can send events in batches.

```typescript
const events = [
  { timestamp: '...', body: { ... } },
  { timestamp: '...', body: { ... } },
  // ... more events
];

const result = await client.detect({
  streamId: 'high-volume-stream',
  events: events
});
```

### Error Handling

The SDK throws typed errors that you can catch and handle.

```typescript
import { DriftlockError, RateLimitError } from '@driftlock/client';

try {
  await client.detect({ ... });
} catch (error) {
  if (error instanceof RateLimitError) {
    console.error('Rate limit exceeded. Retry after:', error.retryAfter);
  } else if (error instanceof DriftlockError) {
    console.error('Driftlock API error:', error.message);
  } else {
    console.error('Unexpected error:', error);
  }
}
```

## Types

The SDK exports TypeScript definitions for all request and response objects.

```typescript
import type { DetectionRequest, DetectionResponse, Anomaly } from '@driftlock/client';

const request: DetectionRequest = {
  streamId: 'test',
  events: []
};
```

## Testing

You can mock the Driftlock client in your tests.

```typescript
// jest.config.js
jest.mock('@driftlock/client');

// my-service.test.ts
import { DriftlockClient } from '@driftlock/client';

test('should handle anomalies', async () => {
  const mockDetect = jest.fn().mockResolvedValue({
    success: true,
    anomalies: [{ id: '1', confidence: 0.9 }]
  });
  
  (DriftlockClient as jest.Mock).mockImplementation(() => ({
    detect: mockDetect
  }));

  // Run your code
});
```

## Support

If you encounter any issues, please [open an issue on GitHub](https://github.com/Shannon-Labs/driftlock-node) or contact [support@driftlock.io](mailto:support@driftlock.io).
