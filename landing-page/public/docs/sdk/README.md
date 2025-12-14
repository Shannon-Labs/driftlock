# SDKs

Official Driftlock SDKs provide high-velocity integration with automatic batching, backpressure handling, and authentication.

## Available SDKs

| SDK | Status | Installation | Features |
|-----|--------|--------------|----------|
| [Node.js/TypeScript](#nodejs) | Available | `npm install @driftlock/node` | Full TypeScript, streaming, OpenAI wrapper |
| [Python](#python) | Available | `pip install driftlock` | Async/await, decorators, context managers |
| [REST API](#rest-api) | Always Available | HTTP client | Direct API access, universal compatibility |

## Node.js

The official Node.js SDK with full TypeScript support.

### Installation

```bash
npm install @driftlock/node
```

### Quick Start

```typescript
import { DriftlockClient } from '@driftlock/node';

const client = new DriftlockClient({
  apiKey: process.env.DRIFTLOCK_API_KEY
});

const result = await client.detect([
  { level: 'info', message: 'Normal operation' },
  { level: 'error', message: 'CRITICAL: System failure' }
], { streamId: 'app-logs' });

console.log(`Found ${result.anomaly_count} anomalies`);
```

### Features

- **Full TypeScript**: Comprehensive types with strict mode
- **Streaming Support**: Automatic batching with event emitters
- **OpenAI Integration**: Monitor LLM requests out-of-the-box
- **Error Handling**: Typed exceptions for all failure scenarios
- **Dual Format**: ESM and CommonJS support
- **Auto Retry**: Exponential backoff with rate limit handling

### Documentation

See [Node.js SDK documentation](./nodejs.md) for complete reference including:
- Configuration options
- Detection API
- Stream management
- Profile tuning
- Error handling
- Examples for logs, metrics, security, and transactions

## Python

The official Python SDK with async/await support.

### Installation

```bash
pip install driftlock
```

### Quick Start

```python
import asyncio
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient(api_key='dlk_your-api-key')

    result = await client.detect([
        {'level': 'info', 'message': 'Normal operation'},
        {'level': 'error', 'message': 'CRITICAL: System failure'}
    ], stream_id='app-logs')

    print(f"Found {result.anomaly_count} anomalies")

asyncio.run(main())
```

### Features

- **Async/Await**: Full asyncio support for high-performance applications
- **Decorators**: Simple `@client.monitor` for existing functions
- **Context Managers**: Clean resource management with async with
- **Streaming**: Automatic batching and event handling
- **Error Handling**: Typed exceptions for all scenarios
- **Type Hints**: Full type annotations for IDE support

### Documentation

See [Python SDK documentation](./python.md) for complete reference including:
- Configuration options
- Async detection API
- Streaming with context managers
- Decorator integration
- Profile tuning
- Error handling
- Examples for logs, metrics, security, and transactions

## REST API

Direct HTTP API access for maximum flexibility.

### Quick Start

```bash
# Try demo endpoint (no auth required)
curl -X POST https://driftlock.net/api/v1/demo/detect \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      "normal log 1",
      "normal log 2",
      "ERROR: anomalous event"
    ]
  }'
```

### Authentication

```bash
curl -X POST https://driftlock.net/api/v1/detect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer dlk_your-api-key" \
  -d '{
    "stream_id": "my-logs",
    "events": ["event 1", "event 2"]
  }'
```

### Documentation

See [REST API documentation](../user-guide/api/rest-api.md) for:
- Full endpoint reference
- Request/response formats
- Error codes
- Rate limiting
- Examples in cURL, JavaScript, Python

## SDK Comparison

| Feature | Node.js | Python | REST API |
|---------|---------|--------|----------|
| Language | TypeScript/JavaScript | Python 3.8+ | Any (HTTP) |
| Async Support | Yes (Promises) | Yes (asyncio) | Manual |
| Streaming | Yes (EventEmitter) | Yes (async context) | Manual batching |
| Auto Retry | Yes | Yes | Manual |
| Error Types | Typed exceptions | Typed exceptions | HTTP status codes |
| Learning Curve | Medium | Low | Low |
| Performance | High | High | Depends on client |

## Common Use Cases

### Application Logs

Monitor logs from your application:

```typescript
// Node.js
const result = await client.detect(logEvents, { streamId: 'app-logs' });
```

```python
# Python
result = await client.detect(log_events, stream_id='app-logs')
```

### API Metrics

Detect performance anomalies:

```typescript
// Node.js
const result = await client.detect(metrics, { streamId: 'api-metrics' });
```

```python
# Python
result = await client.detect(metrics, stream_id='api-metrics')
```

### Security Events

Monitor authentication and access:

```typescript
// Node.js
const result = await client.detect(securityEvents, { streamId: 'security' });
```

```python
# Python
result = await client.detect(security_events, stream_id='security')
```

### LLM Monitoring

Monitor OpenAI requests (Node.js only):

```typescript
import OpenAI from 'openai';
import { wrapOpenAI } from '@driftlock/node';

const wrapped = wrapOpenAI(openai, client, { streamId: 'openai-prod' });
const response = await wrapped.chat.completions.create({...});
```

### Transaction Fraud

Detect fraudulent transactions:

```typescript
// Node.js
const result = await client.detect(transactions, { streamId: 'transactions' });
```

```python
# Python
result = await client.detect(transactions, stream_id='transactions')
```

## Getting Started

### 1. Sign Up

Visit [driftlock.net](https://driftlock.net) to create an account and get your API key.

### 2. Choose an SDK

- **Node.js/TypeScript projects**: Use [Node.js SDK](./nodejs.md)
- **Python projects**: Use [Python SDK](./python.md)
- **Other languages**: Use [REST API](../user-guide/api/rest-api.md)

### 3. Install

```bash
# Node.js
npm install @driftlock/node

# Python
pip install driftlock
```

### 4. Detect Anomalies

See quick start examples above for your chosen SDK.

### 5. Monitor in Production

- Set up stream profiles for your use case
- Enable auto-tuning to learn from feedback
- Configure alerts based on anomalies
- Export data for compliance

## Advanced Features

### Detection Profiles

Choose sensitivity level:

```typescript
// Node.js
await client.streams.updateProfile('logs', { profile: 'sensitive' });
```

```python
# Python
await client.streams.update_profile('logs', {'profile': 'sensitive'})
```

Options: `sensitive`, `balanced`, `strict`, or `custom`

### Auto-Tuning

Improve accuracy with feedback:

```typescript
// Node.js
await client.anomalies.feedback(anomalyId, {
  feedback_type: 'false_positive',
  reason: 'Planned maintenance'
});
```

```python
# Python
await client.anomalies.feedback(anomaly_id,
  feedback_type='false_positive',
  reason='Planned maintenance'
)
```

### Stream Management

Access detection history and metrics:

```typescript
// Node.js
const tuning = await client.streams.getTuning('app-logs');
```

```python
# Python
tuning = await client.streams.get_tuning('app-logs')
```

## Troubleshooting

### Authentication Error

```
DriftlockAuthenticationError: Invalid API key
```

Check that:
1. Your API key is correct (format: `dlk_<uuid>.<secret>`)
2. It's set in `DRIFTLOCK_API_KEY` environment variable or passed to client
3. Your account is active and verified

### Rate Limit

```
DriftlockRateLimitError: Rate limited
```

The SDK automatically handles retries with exponential backoff. Check your plan limits on [driftlock.net](https://driftlock.net).

### Detection Not Ready

If `detection_ready` is false, the stream needs more baseline events. Continue sending normal events until calibration completes.

## Support

- Documentation: https://driftlock.net/docs
- Issues: https://github.com/Shannon-Labs/driftlock/issues
- Email: support@driftlock.net
- Status: https://status.driftlock.net

## License

MIT License - All SDKs are open source and available on GitHub.

---

Ready to get started? [Choose an SDK](#available-sdks) and see the [getting started guide](../user-guide/getting-started/GETTING_STARTED.md).
