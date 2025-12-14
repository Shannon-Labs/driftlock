# Python SDK

Official Python SDK for Driftlock - Compression-Based Anomaly Detection for OpenTelemetry data.

**Status**: Available now with full async/await support and decorator-based integration.

## Installation

```bash
pip install driftlock
```

Or install from source:

```bash
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock/packages/driftlock-python
pip install -e .
```

## Quick Start

```python
import asyncio
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient(api_key='dlk_your-api-key')

    # Detect anomalies in events
    result = await client.detect([
        {'level': 'info', 'message': 'User logged in'},
        {'level': 'error', 'message': 'CRITICAL: System failure'}
    ], stream_id='app-logs')

    print(f"Found {result.anomaly_count} anomalies")
    for anomaly in result.anomalies:
        print(f"- {anomaly.why}")

asyncio.run(main())
```

## Getting Your API Key

1. Sign up at [driftlock.net](https://driftlock.net)
2. Verify your email
3. Copy your API key (format: `dlk_<uuid>.<secret>`)
4. Set as environment variable: `export DRIFTLOCK_API_KEY="dlk_..."`

Or try the demo endpoint (no signup required):

```python
import asyncio
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient()

    result = await client.detection.demo([
        {'message': 'Test event 1'},
        {'message': 'Test event 2'}
    ])

    print(f"Remaining calls: {result.demo.remaining_calls}")

asyncio.run(main())
```

## Configuration

```python
from driftlock import DriftlockClient

client = DriftlockClient(
    api_key='dlk_your-api-key',           # Required
    base_url='https://driftlock.net/api', # Optional
    timeout=30,                           # Request timeout (seconds)
    max_retries=3,                        # Max retry attempts
    retry_delay=1.0,                      # Initial retry delay (seconds)
    max_retry_delay=30.0,                 # Max retry delay (seconds)
)
```

## Detection API

### Simple Detection

```python
result = await client.detect([
    {'message': 'Normal log entry'},
    {'message': 'Unusual event occurred'}
])
```

### With Stream ID

```python
result = await client.detect(events, stream_id='production-logs')
```

### With Configuration Override

```python
result = await client.detect(events, stream_id='sensitive-data',
    config_override={
        'ncd_threshold': 0.25,        # Lower threshold = more sensitive
        'p_value_threshold': 0.01,    # Stricter significance test
        'compressor': 'zstd'          # Use specific compressor
    }
)
```

## Streaming Detection

Automatically batch events with configurable flush behavior:

```python
import asyncio
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient(api_key='dlk_your-api-key')

    stream = client.create_stream(
        stream_id='app-logs',
        batch_size=100,      # Flush every 100 events
        flush_interval=5     # Or flush every 5 seconds
    )

    @stream.on('anomaly')
    def on_anomaly(anomaly):
        print(f"Anomaly detected: {anomaly.why}")
        print(f"  NCD: {anomaly.metrics.ncd}")
        print(f"  Confidence: {anomaly.metrics.confidence_level}")

    @stream.on('batch')
    def on_batch(response):
        print(f"Processed {response.total_events} events")

    @stream.on('error')
    def on_error(error):
        print(f"Detection error: {error}")

    # Push events
    await stream.push({'message': 'Log entry 1'})
    await stream.push({'message': 'Log entry 2'})

    # Manual flush
    await stream.flush()

    # Cleanup
    await stream.close()

asyncio.run(main())
```

## Anomaly Management

### List Anomalies

```python
result = await client.anomalies.list(
    stream_id='app-logs',
    status='open',
    from_date='2025-01-01T00:00:00Z',
    to_date='2025-01-31T23:59:59Z',
    limit=50,
    offset=0
)

print(f"Found {result.total} anomalies")
for anomaly in result.anomalies:
    print(f"- {anomaly.explanation} (NCD: {anomaly.ncd})")
```

### Get Anomaly Details

```python
anomaly = await client.anomalies.get('anomaly-id')
print(anomaly.metrics)
print(anomaly.event)
```

### Submit Feedback

Help improve detection with feedback:

```python
# Mark as false positive
await client.anomalies.feedback('anomaly-id',
    feedback_type='false_positive',
    reason='Planned maintenance window'
)

# Confirm as real anomaly
await client.anomalies.feedback('anomaly-id',
    feedback_type='confirmed'
)
```

### Export Anomalies

```python
# Export single anomaly
job = await client.anomalies.export('anomaly-id', 'pdf')
print(f"Export job: {job.job_id}")

# Export multiple anomalies
bulk_job = await client.anomalies.export_bulk(
    ['id1', 'id2', 'id3'],
    'json'
)
```

## Stream Profiles and Tuning

### Get Detection Profile

```python
profile = await client.streams.get_profile('app-logs')
print(f"Profile: {profile.profile}")
print(f"Auto-tune: {profile.auto_tune_enabled}")
print(f"NCD threshold: {profile.current_thresholds.ncd_threshold}")
```

### Update Profile

```python
await client.streams.update_profile('app-logs', {
    'profile': 'sensitive',          # sensitive | balanced | strict | custom
    'auto_tune_enabled': True,       # Learn from feedback
    'adaptive_window_enabled': True  # Auto-adjust window sizes
})
```

### Get Tuning History

```python
tuning = await client.streams.get_tuning('app-logs')
print(f"False positive rate: {tuning.feedback_stats.false_positive_rate}")

for t in tuning.tune_history:
    print(f"{t.tune_type}: {t.old_value} → {t.new_value}")
    print(f"  Reason: {t.reason}")
```

### Anchor Management

Anchors enable drift detection:

```python
# Get anchor settings
settings = await client.streams.get_anchor('app-logs')
print(f"Has anchor: {settings.has_active_anchor}")

# Get anchor details
details = await client.streams.get_anchor_details('app-logs')
if details.anchor:
    print(f"Event count: {details.anchor.event_count}")
    print(f"Baseline entropy: {details.anchor.baseline_entropy}")

# Reset anchor with known-good events
await client.streams.reset_anchor('app-logs', {
    'events': known_good_events,
    'force_reset': True
})

# Deactivate anchor
await client.streams.deactivate_anchor('app-logs')
```

## Error Handling

The SDK provides typed exceptions for different failure scenarios:

```python
from driftlock.exceptions import (
    DriftlockError,
    DriftlockAuthenticationError,
    DriftlockRateLimitError,
    DriftlockValidationError,
    DriftlockNetworkError
)

try:
    await client.detect(events)
except DriftlockAuthenticationError:
    print("Invalid API key")
except DriftlockRateLimitError as e:
    print(f"Rate limited. Retry after {e.retry_after}s")
except DriftlockValidationError as e:
    print(f"Invalid request: {e}")
except DriftlockNetworkError as e:
    print(f"Network error: {e}")
except DriftlockError as e:
    print(f"API error [{e.code}]: {e.message}")
    print(f"Request ID: {e.request_id}")
```

## Decorator-Based Integration

Use decorators for simple integration with existing functions:

```python
from driftlock import DriftlockClient

client = DriftlockClient(api_key='dlk_your-api-key')

@client.monitor(stream_id='api-calls')
async def process_request(data):
    # Function is automatically monitored
    return await handle_data(data)

# Use normally - monitoring happens automatically
result = await process_request({'key': 'value'})
```

## Context Manager Support

Manage streams with context managers:

```python
from driftlock import DriftlockClient

async def main():
    client = DriftlockClient(api_key='dlk_your-api-key')

    async with client.create_stream(stream_id='app-logs') as stream:
        # Push events
        await stream.push({'message': 'Log 1'})
        await stream.push({'message': 'Log 2'})

    # Automatically flushed and closed

asyncio.run(main())
```

## Common Use Cases

### Application Logs

```python
log_events = [
    {'timestamp': '2025-12-11T10:00:00Z', 'level': 'info', 'message': 'Request started'},
    {'timestamp': '2025-12-11T10:00:01Z', 'level': 'info', 'message': 'Query executed'},
    {'timestamp': '2025-12-11T10:00:02Z', 'level': 'error', 'message': 'Connection timeout'}
]

result = await client.detect(log_events, stream_id='app-logs')
```

### API Metrics

```python
metrics = [
    {'endpoint': '/api/users', 'method': 'GET', 'duration_ms': 45, 'status': 200},
    {'endpoint': '/api/users', 'method': 'GET', 'duration_ms': 52, 'status': 200},
    {'endpoint': '/api/users', 'method': 'GET', 'duration_ms': 3000, 'status': 500}
]

result = await client.detect(metrics, stream_id='api-metrics')
```

### Security Events

```python
security_events = [
    {'action': 'login', 'user': 'john', 'ip': '10.0.0.1', 'success': True},
    {'action': 'login', 'user': 'jane', 'ip': '10.0.0.2', 'success': True},
    {'action': 'login', 'user': 'admin', 'ip': '192.168.1.1', 'success': False, 'attempts': 50}
]

result = await client.detect(security_events, stream_id='security')
```

### Transaction Monitoring

```python
transactions = [
    {'amount': 49.99, 'merchant': 'Amazon', 'country': 'US'},
    {'amount': 125.00, 'merchant': 'Target', 'country': 'US'},
    {'amount': 9999.99, 'merchant': 'Unknown', 'country': 'RU'}
]

result = await client.detect(transactions, stream_id='transactions')
```

## Understanding Detection Results

### Calibration Status

Streams need baseline events before detection works:

```python
if not result.detection_ready:
    print(f"Calibrating: {result.calibration.progress_percent}%")
    print(f"Need {result.detection_events_needed} more events")
else:
    print("Detection is active")
```

### Anomaly Metrics

```python
for anomaly in result.anomalies:
    if anomaly.detected:
        print(f"NCD: {anomaly.metrics.ncd:.4f}")
        print(f"P-value: {anomaly.metrics.p_value:.6f}")
        print(f"Confidence: {anomaly.metrics.confidence_level * 100:.2f}%")
        print(f"Statistically significant: {anomaly.metrics.is_statistically_significant}")
```

## Detection Profiles

Choose sensitivity level for your stream:

```python
# More sensitive - catches more anomalies
await client.streams.update_profile('logs', {'profile': 'sensitive'})

# Default - balanced
await client.streams.update_profile('logs', {'profile': 'balanced'})

# Less sensitive - fewer false positives
await client.streams.update_profile('logs', {'profile': 'strict'})
```

## Auto-Tuning

Let Driftlock learn from feedback:

```python
await client.streams.update_profile('logs', {
    'auto_tune_enabled': True
})

# Submit feedback to improve detection
await client.anomalies.feedback(anomaly_id,
    feedback_type='false_positive',
    reason='Planned maintenance'
)
```

## Health Checks

```python
# Liveness check
health = await client.health()
print(f"Healthy: {health.ok}")

# Readiness check (includes DB)
ready = await client.readiness()
print(f"Ready: {ready.ready}")
print(f"Database: {ready.database}")
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
