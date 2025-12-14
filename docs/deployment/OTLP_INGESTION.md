# OTLP gRPC Ingestion

This document explains how to use OpenTelemetry Protocol (OTLP) gRPC ingestion with the DriftLock Rust API.

## Overview

DriftLock's Rust API includes an optional OTLP gRPC server that accepts OpenTelemetry logs, metrics, and traces. When enabled, the API listens for OTLP data and processes it through the CBAD anomaly detection engine.

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  OTel Collector │────▶│  Driftlock API  │────▶│   PostgreSQL    │
│  or SDK         │     │  (gRPC :4317)   │     │   Database      │
│  (OTLP gRPC)    │     └────────┬────────┘     └─────────────────┘
└─────────────────┘              │
                                 ▼
                        ┌─────────────────┐
                        │   CBAD Core     │
                        │   (Detection)   │
                        └─────────────────┘
```

## Building with OTLP Support

The OTLP server is feature-gated. Build with the `otlp` feature enabled:

```bash
# Build with OTLP support
cargo build -p driftlock-api --features otlp --release

# Or run directly
cargo run -p driftlock-api --features otlp --release

# Enable both OTLP and Kafka
cargo build -p driftlock-api --features otlp,kafka --release
```

## Configuration

Configure OTLP via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `OTLP_ENABLED` | Enable OTLP gRPC server | `false` |
| `OTLP_ADDR` | gRPC bind address | `0.0.0.0:4317` |
| `OTLP_STREAM_ID_ATTR` | Attribute key for stream ID | `driftlock.stream_id` |
| `OTLP_DEFAULT_STREAM_ID` | Fallback stream UUID | (none) |

### Example Configuration

```bash
export OTLP_ENABLED=true
export OTLP_ADDR=0.0.0.0:4317
export OTLP_STREAM_ID_ATTR=driftlock.stream_id
export OTLP_DEFAULT_STREAM_ID=550e8400-e29b-41d4-a716-446655440000
```

## Supported Signal Types

The OTLP server accepts all three OpenTelemetry signal types:

### Logs

Log records are serialized and processed through CBAD. The body and attributes are combined for compression-based analysis.

```
LogsService.Export() → /opentelemetry.proto.collector.logs.v1.LogsService/Export
```

### Metrics

Metrics are serialized (debug format) and processed. Stream ID is extracted from resource attributes.

```
MetricsService.Export() → /opentelemetry.proto.collector.metrics.v1.MetricsService/Export
```

### Traces (Spans)

Spans are serialized and processed. Span attributes are checked for stream ID.

```
TraceService.Export() → /opentelemetry.proto.collector.trace.v1.TraceService/Export
```

## Stream ID Resolution

Each OTLP signal must map to a DriftLock stream. Resolution order:

1. **Signal Attributes**: Check log/span attributes for `OTLP_STREAM_ID_ATTR` key
2. **Resource Attributes**: Check resource attributes (for metrics)
3. **Default Stream**: Fall back to `OTLP_DEFAULT_STREAM_ID` if set
4. **Error**: Reject if no stream ID found

### Setting Stream ID in OTel SDK

**Go:**
```go
import "go.opentelemetry.io/otel/attribute"

logger.Info("message",
    attribute.String("driftlock.stream_id", "550e8400-e29b-41d4-a716-446655440000"),
)
```

**Python:**
```python
logger.info("message", extra={
    "driftlock.stream_id": "550e8400-e29b-41d4-a716-446655440000"
})
```

**Node.js:**
```javascript
logger.info('message', {
    'driftlock.stream_id': '550e8400-e29b-41d4-a716-446655440000'
});
```

## OpenTelemetry Collector Configuration

Configure the OTel Collector to export to DriftLock:

```yaml
exporters:
  otlp/driftlock:
    endpoint: "driftlock-api:4317"
    tls:
      insecure: true  # Set to false in production

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/driftlock]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/driftlock]
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/driftlock]
```

## Metrics

When OTLP is enabled, these Prometheus metrics are emitted:

| Metric | Description |
|--------|-------------|
| `driftlock_events_processed_total` | Total events processed |
| `driftlock_anomalies_detected_total` | Total anomalies detected |
| `driftlock_stream_events_total{stream_id}` | Events per stream |
| `driftlock_stream_anomalies_total{stream_id}` | Anomalies per stream |

## Local Development

### Testing with grpcurl

```bash
# Start API with OTLP enabled
OTLP_ENABLED=true \
OTLP_DEFAULT_STREAM_ID=550e8400-e29b-41d4-a716-446655440000 \
cargo run -p driftlock-api --features otlp

# Test with grpcurl (requires stream in database)
grpcurl -plaintext localhost:4317 list
```

### Testing with OTel Collector

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4318

exporters:
  otlp/driftlock:
    endpoint: "localhost:4317"
    tls:
      insecure: true

service:
  pipelines:
    logs:
      receivers: [otlp]
      exporters: [otlp/driftlock]
```

```bash
docker run -v $(pwd)/otel-collector-config.yaml:/etc/otel/config.yaml \
  otel/opentelemetry-collector:latest \
  --config /etc/otel/config.yaml
```

## Production Considerations

### Security

- Enable TLS for gRPC connections (configure certificate paths)
- Use mTLS for client authentication
- Restrict access via network policies

### Performance

- OTLP processes signals synchronously
- Each signal type triggers CBAD detection
- Consider batching at the collector level

### High Availability

- Deploy multiple API instances behind a load balancer
- gRPC supports client-side load balancing
- Stateless processing allows horizontal scaling

## Troubleshooting

### Server Not Starting

Check logs for binding errors:
```bash
cargo run -p driftlock-api --features otlp 2>&1 | grep -i otlp
```

Common issues:
- Port 4317 already in use
- Invalid bind address format
- Feature not enabled at build time

### Signals Not Processing

1. Verify stream exists in database
2. Check stream ID attribute key matches configuration
3. Ensure `OTLP_DEFAULT_STREAM_ID` is set if not using attributes
4. Check resource/log attributes are set correctly

### Connection Refused

- Verify `OTLP_ENABLED=true`
- Check firewall allows port 4317
- Ensure gRPC client uses correct protocol (not HTTP)

## Implementation Details

Source: `crates/driftlock-api/src/otlp.rs`

Key components:
- `spawn_otlp_server()`: Spawns gRPC server task
- `LogsService::export()`: Handles log ingestion
- `MetricsService::export()`: Handles metric ingestion
- `TraceService::export()`: Handles span ingestion
- `resolve_stream_id()`: Extract stream ID from attributes
