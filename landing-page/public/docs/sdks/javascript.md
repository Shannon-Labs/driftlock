# JavaScript SDK

The official browser-based client for Driftlock. Detect anomalies directly from your frontend applications.

## Installation

Install via npm:

```bash
npm install @driftlock/browser
```

Or use via CDN:

```html
<script src="https://cdn.jsdelivr.net/npm/@driftlock/browser@latest/dist/driftlock.min.js"></script>
```

## Quick Start

Initialize the client and start detecting anomalies.

```javascript
import { DriftlockBrowser } from '@driftlock/browser';

// Initialize the client
const client = new DriftlockBrowser({
  apiKey: 'your-public-api-key' // Use a restricted API key for frontend usage
});

// Detect anomalies
async function checkUserBehavior() {
  const result = await client.detect({
    streamId: 'user-interactions',
    events: [
      {
        timestamp: new Date().toISOString(),
        type: 'click',
        body: { x: 100, y: 200, element: 'button' }
      }
    ]
  });

  if (result.anomalies.length > 0) {
    console.warn('Anomalous behavior detected:', result.anomalies);
  }
}
```

## Security Note

> [!WARNING]
> **Never expose your secret API key in frontend code.**
> Always use a **restricted API key** that only has permissions to `POST /detect` and is restricted to your domain. You can create restricted keys in the [Driftlock Dashboard](https://driftlock.web.app/dashboard).

## Configuration

```javascript
const client = new DriftlockBrowser({
  apiKey: 'your-public-key',
  endpoint: 'https://driftlock-api-o6kjgrsowq-uc.a.run.app',
  autoRetry: true
});
```

## Framework Integration

### React

```jsx
import { useEffect } from 'react';
import { useDriftlock } from '@driftlock/react';

function App() {
  const { detect } = useDriftlock();

  const handleClick = async () => {
    await detect({
      streamId: 'button-clicks',
      events: [{ type: 'click', body: { timestamp: Date.now() } }]
    });
  };

  return <button onClick={handleClick}>Click Me</button>;
}
```

### Vue.js

```vue
<script setup>
import { inject } from 'vue';

const driftlock = inject('driftlock');

const trackEvent = async () => {
  await driftlock.detect({
    streamId: 'vue-events',
    events: [{ type: 'interaction', body: { page: 'home' } }]
  });
};
</script>
```

## Error Handling

```javascript
try {
  await client.detect(...);
} catch (error) {
  console.error('Driftlock error:', error);
}
```

## Support

For issues, please [open an issue on GitHub](https://github.com/Shannon-Labs/driftlock-browser) or contact [support@driftlock.io](mailto:support@driftlock.io).
