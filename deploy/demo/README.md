# Driftlock Demo Environment

This directory contains the configuration and deployment files for the Driftlock demo environment at demo.driftlock.com.

## Overview

The demo environment showcases Driftlock's anomaly detection capabilities with realistic synthetic data. It demonstrates:

- Real-time anomaly detection using CBAD (Compression-Based Anomaly Detection)
- Multi-tenant architecture with tenant isolation
- Streaming data processing with Kafka
- Analytics with ClickHouse
- Monitoring with Prometheus and Grafana
- Interactive web UI

## Architecture

```
┌─────────────────┐
│  Web UI        │
│  (Next.js)     │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  API Server     │
│  (Go)          │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  OTel Collector  │
│                 │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  Kafka          │
│  (Streaming)    │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  ClickHouse      │
│  (Analytics)    │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  PostgreSQL      │
│  (Database)     │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  Redis          │
│  (Cache)        │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  Prometheus     │
│  (Monitoring)   │
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  Grafana       │
│  (Dashboards)  │
└─────────────────┘
```

## Quick Start

1. **Start the demo environment:**
   ```bash
   cd deploy/demo
   docker-compose up -d
   ```

2. **Access the demo:**
   - Web UI: http://localhost:3000 (admin/admin)
   - API: http://localhost:8080
   - Grafana: http://localhost:3000 (admin/admin)
   - Prometheus: http://localhost:9090

3. **Generate synthetic data:**
   ```bash
   docker-compose up -d data-generator
   # Or run with custom parameters:
   docker-compose run --rm data-generator \
     -e API_URL=http://api:8080 \
     -e EVENTS_PER_SECOND=20 \
     -e ANOMALY_RATE=0.1
   ```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|----------|-------------|
| API_URL | http://api:8080 | API server endpoint |
| EVENTS_PER_SECOND | 10 | Events to generate per second |
| ANOMALY_RATE | 0.05 | Probability of anomalies (0.0-1.0) |
| TENANT_ID | demo | Tenant ID for multi-tenant setup |

### Synthetic Data

The synthetic data generator creates realistic OTLP events with the following characteristics:

- **Services**: API Gateway, User Service, Payment Service, etc.
- **Operations**: GET, POST, PUT, DELETE, PATCH
- **Status**: Success, Error, Timeout
- **Anomalies**: Statistical, Behavioral, Contextual, Collective, Temporal, Spatial, Graph
- **Indicators**: Request rate, Response time, Error rate, CPU usage, Memory usage, etc.

### Demo Scenarios

1. **Normal Operations**: Regular traffic patterns with normal response times
2. **Anomaly Spike**: Sudden increase in request rate or response time
3. **Error Storm**: High error rate with specific error patterns
4. **Geographic Anomaly**: Unusual access from different geographic regions
5. **Behavioral Anomaly**: Atypical user behavior patterns

## Monitoring

### Grafana Dashboards

- **Overview Dashboard**: System health, request rates, anomaly detection
- **Tenant Dashboard**: Per-tenant resource usage and quotas
- **Anomaly Dashboard**: Detailed anomaly analysis and trends
- **Performance Dashboard**: System performance metrics

### Prometheus Metrics

- **System Metrics**: CPU, memory, disk, network
- **Application Metrics**: Request rate, response time, error rate
- **Anomaly Metrics**: Detection rate, types, confidence scores
- **Tenant Metrics**: Resource usage per tenant

## Troubleshooting

### Common Issues

1. **Services not starting**: Check network connectivity and dependencies
2. **High memory usage**: Reduce synthetic data generation rate
3. **No anomalies detected**: Increase ANOMALY_RATE environment variable
4. **Grafana not showing data**: Verify Prometheus datasource configuration

### Logs

- **API Server**: `docker-compose logs api`
- **Data Generator**: `docker-compose logs data-generator`
- **Kafka**: `docker-compose logs kafka`
- **ClickHouse**: `docker-compose logs clickhouse`

## Security

- **Authentication**: Demo uses basic authentication (admin/admin)
- **Network**: All services communicate within a Docker network
- **Data**: All synthetic data is randomly generated and contains no real user information

## Scaling

The demo environment is configured for small-scale demonstrations:

- **API Server**: 1 replica, 512MB memory limit
- **Kafka**: 1 broker, 1GB disk limit
- **ClickHouse**: 1 node, 2GB disk limit
- **PostgreSQL**: 1 instance, 1GB disk limit
- **Redis**: 1 instance, 256MB memory limit

For production deployments, adjust these values based on your requirements.
