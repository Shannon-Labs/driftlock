# Express.js Integration

Integrate Driftlock anomaly detection into your Express.js application using our middleware.

## Installation

```bash
npm install @driftlock/express
```

## Basic Usage

The middleware automatically captures request metrics (latency, status code, payload size) and sends them to Driftlock for anomaly detection.

```javascript
const express = require('express');
const driftlock = require('@driftlock/express');

const app = express();

// Initialize middleware
app.use(driftlock({
  apiKey: process.env.DRIFTLOCK_API_KEY,
  streamId: 'api-requests'
}));

app.get('/', (req, res) => {
  res.send('Hello World!');
});

app.listen(3000, () => {
  console.log('Server started on port 3000');
});
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `apiKey` | `string` | Required | Your Driftlock API key |
| `streamId` | `string` | `'default'` | The stream ID to send events to |
| `sampleRate` | `number` | `1.0` | Percentage of requests to sample (0.0 - 1.0) |
| `excludePaths` | `string[]` | `[]` | List of paths to exclude from monitoring |

```javascript
app.use(driftlock({
  apiKey: '...',
  streamId: 'production-api',
  sampleRate: 0.5, // Sample 50% of requests
  excludePaths: ['/health', '/metrics']
}));
```

## Custom Event Data

You can attach custom data to the Driftlock event within your route handlers.

```javascript
app.post('/checkout', (req, res) => {
  // Add custom context to the current request's anomaly detection event
  req.driftlock.addContext({
    userId: req.user.id,
    cartValue: req.body.total
  });

  // Process checkout...
  res.json({ success: true });
});
```

## Handling Anomalies

By default, the middleware logs anomalies to the console. You can provide a custom callback to handle them (e.g., send to Slack, block request).

**Note**: Anomaly detection happens *asynchronously* after the response is sent to minimize latency impact, unless you configure it to block.

```javascript
app.use(driftlock({
  apiKey: '...',
  onAnomaly: (anomaly, req) => {
    console.warn(`Anomaly detected on ${req.path}:`, anomaly.why);
    // Send alert to Slack/PagerDuty
  }
}));
```

## Blocking Anomalous Requests

If you want to block requests *before* they are processed based on anomaly detection (e.g., for WAF-like behavior), use the `blockAnomalies` option. **Warning**: This adds latency to every request.

```javascript
app.use(driftlock({
  apiKey: '...',
  blockAnomalies: true, // Wait for detection before proceeding
  onAnomaly: (anomaly, req, res) => {
    res.status(403).json({ error: 'Request blocked due to anomaly detection' });
  }
}));
```

## Error Handling

The middleware will not crash your application if Driftlock is unreachable. Errors are logged via `console.error` by default.

```javascript
app.use(driftlock({
  apiKey: '...',
  onError: (err) => {
    // Custom error logging
    logger.error('Driftlock middleware error:', err);
  }
}));
```
