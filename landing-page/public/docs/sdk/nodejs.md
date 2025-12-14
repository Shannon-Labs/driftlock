# Node.js SDK

Official Node.js SDK for Driftlock - Compression-Based Anomaly Detection for OpenTelemetry data.

[![npm version](https://img.shields.io/npm/v/@driftlock/node.svg)](https://www.npmjs.com/package/@driftlock/node)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Compression-Based Anomaly Detection (CBAD)**: Detect anomalies using information-theoretic compression
- **Streaming Support**: Automatic batching with event-driven architecture
- **OpenAI Integration**: Monitor LLM requests for unusual behavior
- **Full TypeScript**: Comprehensive type definitions with strict mode
- **Automatic Retries**: Exponential backoff with configurable retry logic
- **Rate Limit Handling**: Respects `Retry-After` headers
- **ESM & CommonJS**: Dual package format support

## Installation

```bash
npm install @driftlock/node
```

## Quick Start

```typescript
import { DriftlockClient } from '@driftlock/node';

const client = new DriftlockClient({
  apiKey: 'dlk_your-api-key'
});

// Detect anomalies in events
const result = await client.detect([
  { level: 'info', message: 'User logged in' },
  { level: 'error', message: 'CRITICAL: System failure' }
], { streamId: 'app-logs' });

console.log(`Found ${result.anomaly_count} anomalies`);
result.anomalies.forEach(anomaly => {
  console.log(`- ${anomaly.why}`);
});
```

## Getting Your API Key

1. Sign up at [driftlock.net](https://driftlock.net)
2. Verify your email
3. Copy your API key (format: `dlk_<uuid>.<secret>`)
4. Set as environment variable: `export DRIFTLOCK_API_KEY="dlk_..."`

Or try the demo endpoint (no signup required):

```typescript
const result = await client.detection.demo([
  { message: 'Test event 1' },
  { message: 'Test event 2' }
]);

console.log(result.demo?.remaining_calls); // Rate limit info
```

## Configuration

```typescript
const client = new DriftlockClient({
  apiKey: 'dlk_your-api-key',           // Required
  baseUrl: 'https://driftlock.net/api', // Optional
  timeout: 30000,                       // Request timeout (ms)
  maxRetries: 3,                        // Max retry attempts
  retryDelay: 1000,                     // Initial retry delay (ms)
  maxRetryDelay: 30000,                 // Max retry delay (ms)
});
```

## Detection API

### Simple Detection

```typescript
const result = await client.detect([
  { message: 'Normal log entry' },
  { message: 'Unusual event occurred' }
]);
```

### With Stream ID

```typescript
const result = await client.detect(events, {
  streamId: 'production-logs'
});
```

### With Configuration Override

```typescript
const result = await client.detect(events, {
  streamId: 'sensitive-data',
  config_override: {
    ncd_threshold: 0.25,        // Lower threshold = more sensitive
    p_value_threshold: 0.01,    // Stricter significance test
    compressor: 'zstd'          // Use specific compressor
  }
});
```

## Streaming Detection

Automatically batch events with configurable flush behavior:

```typescript
const stream = client.createStream({
  streamId: 'app-logs',
  batchSize: 100,      // Flush every 100 events
  flushInterval: 5000  // Or flush every 5 seconds
});

// Listen for anomalies
stream.on('anomaly', (anomaly) => {
  console.log(`Anomaly detected: ${anomaly.why}`);
  console.log(`  NCD: ${anomaly.metrics.ncd}`);
  console.log(`  Confidence: ${anomaly.metrics.confidence_level}`);
});

// Listen for batch completions
stream.on('batch', (response) => {
  console.log(`Processed ${response.total_events} events`);
});

// Listen for errors
stream.on('error', (error) => {
  console.error('Detection error:', error);
});

// Push events
stream.push({ message: 'Log entry 1' });
stream.push({ message: 'Log entry 2' });

// Manual flush
await stream.flush();

// Cleanup
await stream.destroy();
```

## OpenAI Integration

Monitor LLM requests for unusual behavior:

```typescript
import OpenAI from 'openai';
import { DriftlockClient, wrapOpenAI } from '@driftlock/node';

const openai = new OpenAI({ apiKey: process.env.OPENAI_API_KEY });
const driftlock = new DriftlockClient({ apiKey: process.env.DRIFTLOCK_API_KEY });

const wrapped = wrapOpenAI(openai, driftlock, {
  streamId: 'openai-prod',
  onAnomaly: (anomaly) => {
    console.error('Unusual LLM behavior:', anomaly.why);
    // Send alert, log to monitoring, etc.
  },
  throwOnAnomaly: false,  // Don't throw on detection
  includeData: true       // Include request/response metadata
});

// Use wrapped client normally
const response = await wrapped.chat.completions.create({
  model: 'gpt-4',
  messages: [{ role: 'user', content: 'Hello!' }]
});
```

## Anomaly Management

### List Anomalies

```typescript
const result = await client.anomalies.list({
  stream_id: 'app-logs',
  status: 'open',
  from: '2025-01-01T00:00:00Z',
  to: '2025-01-31T23:59:59Z',
  limit: 50,
  offset: 0
});

console.log(`Found ${result.total} anomalies`);
result.anomalies.forEach(a => {
  console.log(`- ${a.explanation} (NCD: ${a.ncd})`);
});
```

### Get Anomaly Details

```typescript
const anomaly = await client.anomalies.get('anomaly-id');
console.log(anomaly.metrics);
console.log(anomaly.event);
```

### Submit Feedback

Help improve detection with feedback:

```typescript
// Mark as false positive
await client.anomalies.feedback('anomaly-id', {
  feedback_type: 'false_positive',
  reason: 'Planned maintenance window'
});

// Confirm as real anomaly
await client.anomalies.feedback('anomaly-id', {
  feedback_type: 'confirmed'
});
```

### Export Anomalies

```typescript
// Export single anomaly
const job = await client.anomalies.export('anomaly-id', 'pdf');
console.log(`Export job: ${job.job_id}`);

// Export multiple anomalies
const bulkJob = await client.anomalies.exportBulk(
  ['id1', 'id2', 'id3'],
  'json'
);
```

## Stream Profiles and Tuning

### Get Detection Profile

```typescript
const profile = await client.streams.getProfile('app-logs');
console.log(`Profile: ${profile.profile}`);
console.log(`Auto-tune: ${profile.auto_tune_enabled}`);
console.log(`NCD threshold: ${profile.current_thresholds.ncd_threshold}`);
```

### Update Profile

```typescript
await client.streams.updateProfile('app-logs', {
  profile: 'sensitive',          // sensitive | balanced | strict | custom
  auto_tune_enabled: true,       // Learn from feedback
  adaptive_window_enabled: true  // Auto-adjust window sizes
});
```

### Get Tuning History

```typescript
const tuning = await client.streams.getTuning('app-logs');
console.log(`False positive rate: ${tuning.feedback_stats.false_positive_rate}`);

tuning.tune_history.forEach(t => {
  console.log(`${t.tune_type}: ${t.old_value} → ${t.new_value}`);
  console.log(`  Reason: ${t.reason}`);
});
```

### Anchor Management

Anchors enable drift detection:

```typescript
// Get anchor settings
const settings = await client.streams.getAnchor('app-logs');
console.log(`Has anchor: ${settings.has_active_anchor}`);

// Get anchor details
const details = await client.streams.getAnchorDetails('app-logs');
if (details.anchor) {
  console.log(`Event count: ${details.anchor.event_count}`);
  console.log(`Baseline entropy: ${details.anchor.baseline_entropy}`);
}

// Reset anchor with known-good events
await client.streams.resetAnchor('app-logs', {
  events: knownGoodEvents,
  force_reset: true
});

// Deactivate anchor
await client.streams.deactivateAnchor('app-logs');
```

## Error Handling

The SDK provides typed errors for different failure scenarios:

```typescript
import {
  DriftlockError,
  DriftlockAuthenticationError,
  DriftlockRateLimitError,
  DriftlockValidationError,
  DriftlockNetworkError
} from '@driftlock/node';

try {
  await client.detect(events);
} catch (error) {
  if (error instanceof DriftlockAuthenticationError) {
    console.error('Invalid API key');
  } else if (error instanceof DriftlockRateLimitError) {
    console.error(`Rate limited. Retry after ${error.retryAfter}s`);
  } else if (error instanceof DriftlockValidationError) {
    console.error('Invalid request:', error.message);
  } else if (error instanceof DriftlockNetworkError) {
    console.error('Network error:', error.message);
  } else if (error instanceof DriftlockError) {
    console.error(`API error [${error.code}]:`, error.message);
    console.error(`Request ID: ${error.requestId}`);
  }
}
```

## TypeScript Support

The SDK is written in TypeScript with strict mode enabled. All types are exported:

```typescript
import type {
  DetectRequest,
  DetectResponse,
  AnomalyOutput,
  AnomalyMetrics,
  ConfigOverride,
  DetectionProfile
} from '@driftlock/node';
```

## Health Checks

```typescript
// Liveness check
const health = await client.health();
console.log('Healthy:', health.ok);

// Readiness check (includes DB)
const ready = await client.readiness();
console.log('Ready:', ready.ready);
console.log('Database:', ready.database);
```

## Common Use Cases

### Application Logs

```typescript
const logEvents = [
  { timestamp: '2025-12-11T10:00:00Z', level: 'info', message: 'Request started' },
  { timestamp: '2025-12-11T10:00:01Z', level: 'info', message: 'Query executed' },
  { timestamp: '2025-12-11T10:00:02Z', level: 'error', message: 'Connection timeout' }
];

const result = await client.detect(logEvents, { streamId: 'app-logs' });
```

### API Metrics

```typescript
const metrics = [
  { endpoint: '/api/users', method: 'GET', duration_ms: 45, status: 200 },
  { endpoint: '/api/users', method: 'GET', duration_ms: 52, status: 200 },
  { endpoint: '/api/users', method: 'GET', duration_ms: 3000, status: 500 }
];

const result = await client.detect(metrics, { streamId: 'api-metrics' });
```

### Security Events

```typescript
const securityEvents = [
  { action: 'login', user: 'john', ip: '10.0.0.1', success: true },
  { action: 'login', user: 'jane', ip: '10.0.0.2', success: true },
  { action: 'login', user: 'admin', ip: '192.168.1.1', success: false, attempts: 50 }
];

const result = await client.detect(securityEvents, { streamId: 'security' });
```

### Transaction Monitoring

```typescript
const transactions = [
  { amount: 49.99, merchant: 'Amazon', country: 'US' },
  { amount: 125.00, merchant: 'Target', country: 'US' },
  { amount: 9999.99, merchant: 'Unknown', country: 'RU' }
];

const result = await client.detect(transactions, { streamId: 'transactions' });
```

## Understanding Detection Results

### Calibration Status

Streams need baseline events before detection works:

```typescript
if (!result.detection_ready) {
  console.log(`Calibrating: ${result.calibration.progress_percent}%`);
  console.log(`Need ${result.detection_events_needed} more events`);
} else {
  console.log('Detection is active');
}
```

### Anomaly Metrics

```typescript
result.anomalies.forEach(anomaly => {
  if (anomaly.detected) {
    console.log(`NCD: ${anomaly.metrics.ncd.toFixed(4)}`);
    console.log(`P-value: ${anomaly.metrics.p_value.toFixed(6)}`);
    console.log(`Confidence: ${(anomaly.metrics.confidence_level * 100).toFixed(2)}%`);
    console.log(`Statistically significant: ${anomaly.metrics.is_statistically_significant}`);
  }
});
```

## Detection Profiles

Choose sensitivity level for your stream:

```typescript
// More sensitive - catches more anomalies
await client.streams.updateProfile('logs', { profile: 'sensitive' });

// Default - balanced
await client.streams.updateProfile('logs', { profile: 'balanced' });

// Less sensitive - fewer false positives
await client.streams.updateProfile('logs', { profile: 'strict' });
```

## Auto-Tuning

Let Driftlock learn from feedback:

```typescript
await client.streams.updateProfile('logs', {
  auto_tune_enabled: true
});

// Submit feedback to improve detection
await client.anomalies.feedback(anomalyId, {
  feedback_type: 'false_positive',
  reason: 'Planned maintenance'
});
```

## API Reference

For complete API documentation and examples, see:

- [REST API Reference](../user-guide/api/rest-api.md)
- [OpenAPI Specification](https://driftlock.net/api/openapi.yaml)

## Support

- Documentation: https://driftlock.net/docs
- Issues: https://github.com/Shannon-Labs/driftlock/issues
- Email: support@driftlock.net

## License

MIT License - see [LICENSE](https://github.com/Shannon-Labs/driftlock/blob/main/LICENSE) for details.

---

Built with compression-based anomaly detection.
