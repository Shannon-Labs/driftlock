# Performance Optimization

Maximize the performance of your application while using Driftlock.

## Asynchronous Detection

The single most important performance tip is to **never block your user's request** while waiting for anomaly detection, unless you strictly require it (e.g., for security blocking).

### Good Pattern (Fire and Forget)

```javascript
// Node.js / Express
app.post('/api/data', async (req, res) => {
  // 1. Process user request immediately
  const data = await db.save(req.body);
  res.json(data);

  // 2. Send to Driftlock in background
  driftlock.detect({ ... }).catch(console.error);
});
```

### Bad Pattern (Blocking)

```javascript
app.post('/api/data', async (req, res) => {
  // ⚠️ Adds network latency to every request!
  await driftlock.detect({ ... });
  
  const data = await db.save(req.body);
  res.json(data);
});
```

## Batching

If you have a high-throughput stream (e.g., > 100 events/second), sending an HTTP request for every event is inefficient. Use batching.

### Manual Batching

Accumulate events in a buffer and flush them periodically or when the buffer is full.

```javascript
const buffer = [];
const FLUSH_SIZE = 50;
const FLUSH_INTERVAL = 1000; // 1 second

function track(event) {
  buffer.push(event);
  if (buffer.length >= FLUSH_SIZE) flush();
}

setInterval(flush, FLUSH_INTERVAL);

async function flush() {
  if (buffer.length === 0) return;
  const events = [...buffer];
  buffer.length = 0;
  
  await driftlock.detect({ events });
}
```

## Connection Reuse

Ensure your HTTP client reuses TCP connections (Keep-Alive). The official SDKs handle this automatically, but if you are using a raw HTTP client, verify your configuration.

## Compression

Driftlock accepts compressed request bodies. If you are sending large batches, enable GZIP compression in your HTTP client to reduce bandwidth usage.

```bash
curl -H "Content-Encoding: gzip" --data-binary @events.json.gz ...
```

## Latency Expectations

- **P95 Latency**: ~150ms for single events
- **P99 Latency**: ~300ms

If you observe higher latencies:
1. Check your network connection to Google Cloud (us-central1).
2. Verify you aren't sending excessively large payloads (> 1MB).
3. Ensure you aren't hitting rate limits (which return 429s).
