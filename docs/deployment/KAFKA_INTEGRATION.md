# Kafka Integration with DriftLock

This document explains how Kafka is integrated with the DriftLock project for streaming OTLP events and anomaly detection results.

## Architecture Overview

DriftLock uses Apache Kafka as a streaming backbone for:

1. **OTLP Events Ingestion**: OpenTelemetry logs and metrics are published to the `otlp-events` topic
2. **Anomaly Events Distribution**: Detected anomalies are published to the `anomaly-events` topic
3. **Distributed Processing**: Multiple collector instances can process events in parallel

## Components

### 1. Kafka Publisher

The `collector-processor/driftlockcbad/kafka/publisher.go` component handles publishing OTLP events to Kafka:

- Serializes log and metric data as JSON
- Adds metadata headers for event routing
- Supports batch publishing for performance
- Configurable TLS support for secure connections

### 2. Collector Processor Integration

The `driftlock_cbad` processor can be configured to publish events to Kafka:

```yaml
processors:
  driftlock_cbad:
    window_size: 1024
    hop_size: 256
    threshold: 0.9
    determinism: true
    kafka:
      enabled: true
      brokers: ["localhost:9092"]
      client_id: "driftlock-collector"
      events_topic: "otlp-events"
      tls_enabled: false
      batch_size: 100
      batch_timeout_ms: 5
```

### 3. Kafka Topics

Two primary topics are used:

- `otlp-events`: For incoming OpenTelemetry data
  - Partitions: 3 (default)
  - Replication factor: 1 (default)
  
- `anomaly-events`: For anomaly detection results
  - Partitions: 3 (default)
  - Replication factor: 1 (default)

## Setup Options

### Option 1: Native Installation with Colima (Recommended for Development)

Use native Homebrew installation with Colima for better performance on macOS:

```bash
# Install Kafka and Zookeeper
./scripts/kafka-native-setup.sh install

# Start Kafka
./scripts/kafka-native-setup.sh start
```

### Option 2: Docker Installation with Colima

Use Docker Compose with Colima for a self-contained environment:

```bash
# Start complete stack with Kafka
./scripts/start-stack.sh start
```

## Configuration

### Collector Configuration

Update `deploy/collector-config/config.yaml` to enable Kafka publishing:

```yaml
processors:
  driftlock_cbad:
    # ... other settings ...
    kafka:
      enabled: true
      brokers: ["localhost:9092"]  # Use "kafka:29092" for Docker
      client_id: "driftlock-collector"
      events_topic: "otlp-events"
      tls_enabled: false
      batch_size: 100
      batch_timeout_ms: 5
```

### Docker Compose Configuration

The `deploy/docker-compose.yml` includes Kafka and Zookeeper services:

```yaml
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    # ... configuration ...

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    depends_on:
      - zookeeper
    # ... configuration ...
```

## Event Format

### OTLP Log Event

```json
{
  "type": "log",
  "data": {
    "timestamp": "2025-10-26T13:14:22.123456789Z",
    "severity": "INFO",
    "body": "Request completed successfully",
    "severity_number": 9,
    "flags": 0,
    "attributes": {
      "service.name": "api-gateway",
      "http.method": "GET",
      "http.status_code": "200"
    }
  },
  "timestamp": "2025-10-26T13:14:22.123Z",
  "source": "collector-processor"
}
```

### OTLP Metric Event

```json
{
  "type": "metric",
  "data": {
    "name": "http_request_duration",
    "description": "HTTP request duration",
    "unit": "ms",
    "type": "histogram",
    "histogram": [
      {
        "timestamp": "2025-10-26T13:14:22.123456789Z",
        "count": 100,
        "sum": 5234.5,
        "min": 12.3,
        "max": 234.5,
        "attributes": {
          "service.name": "api-gateway",
          "http.method": "GET"
        }
      }
    ]
  },
  "timestamp": "2025-10-26T13:14:22.123Z",
  "source": "collector-processor"
}
```

## Management Scripts

### Kafka Native Setup Script

`scripts/kafka-native-setup.sh` provides commands for managing a native Kafka installation:

```bash
./scripts/kafka-native-setup.sh install  # Install Kafka and Zookeeper
./scripts/kafka-native-setup.sh start    # Start services and create topics
./scripts/kafka-native-setup.sh stop     # Stop services
./scripts/kafka-native-setup.sh status   # Check service status
./scripts/kafka-native-setup.sh topics   # List topics
./scripts/kafka-native-setup.sh test     # Test connectivity
```

### Stack Management Script

`scripts/start-stack.sh` provides commands for managing the complete DriftLock stack:

```bash
./scripts/start-stack.sh start    # Start the complete stack
./scripts/start-stack.sh stop     # Stop the complete stack
./scripts/start-stack.sh restart  # Restart the complete stack
./scripts/start-stack.sh status   # Check status of all services
./scripts/start-stack.sh logs     # Show logs from all services
```

## Performance Considerations

1. **Batch Size**: Adjust `batch_size` based on throughput requirements
2. **Batch Timeout**: Tune `batch_timeout_ms` for latency vs. throughput balance
3. **Partitions**: Increase topic partitions for higher parallelism
4. **Compression**: Consider enabling compression for large payloads

## Troubleshooting

### Connection Issues with Colima

1. Verify Colima is running:
   ```bash
   colima list
   ```

2. Check if Kafka is running:
   ```bash
   ./scripts/kafka-native-setup.sh status
   ```

3. Check Kafka logs:
   ```bash
   ./scripts/kafka-native-setup.sh logs
   ```

4. Test connectivity:
   ```bash
   ./scripts/kafka-native-setup.sh test
   ```

### Container Access Issues

When using Colima, container networking works differently:

1. Use container IPs instead of localhost:
   ```bash
   # Get container IP
   docker inspect <container-name> | grep IPAddress
   ```

2. Use `colima exec` instead of `docker exec`:
   ```bash
   colima exec <container-name> <command>
   ```

3. Port forwarding is handled automatically by Colima

### API Container Issues

If API container is running but not responding:

1. Check container logs:
   ```bash
   docker logs <container-name>
   ```

2. Access API using container IP:
   ```bash
   # Get container IP
   IP=$(docker inspect <container-name> | grep IPAddress | awk '{print $2}')
   curl -s http://$IP:8080/healthz
   ```

3. The distroless/static image doesn't include curl - use a different base image if needed:
   ```dockerfile
   FROM alpine:latest
   RUN apk add --no-cache curl
   # ... rest of your Dockerfile
   ```

## Current Status

✅ **Kafka**: Running with native Homebrew installation
✅ **Zookeeper**: Running with native Homebrew installation  
✅ **Kafka Topics**: `otlp-events` and `anomaly-events` created
✅ **Docker Compose**: Running with Colima (Zookeeper, Kafka, Collector, API)
⚠️ **API Container**: Running but not responding to requests (networking issue with Colima)

## Next Steps

1. **Fix API Container Networking Issue**:
   - The API container is using a minimal distroless/static image without curl
   - Need to either:
     a) Add curl to the Dockerfile, or
     b) Use a different base image that includes common utilities
   - Update the start-stack.sh script to handle container IP addressing

2. **Test Kafka Integration**:
   - Send test OTLP events to the collector
   - Verify events are published to Kafka topics
   - Check if anomaly detection is working

3. **Implement Kafka Consumer**:
   - Create a consumer service to read from `anomaly-events` topic
   - Build a simple web UI to display detected anomalies
   - Add real-time streaming of anomaly events

4. **Documentation Updates**:
   - Update API documentation with correct endpoints
   - Add examples of sending OTLP data to the collector
   - Document the complete data flow from collector → Kafka → consumer → UI
