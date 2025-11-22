# Driftlock API Documentation

## Overview

The Driftlock API provides endpoints for anomaly detection, event processing, and system management. This documentation covers all available endpoints, authentication methods, and integration examples.

## Base URL

```
https://api.driftlock.com/v1
```

## Authentication

Driftlock uses API key authentication for secure access to the API.

### Headers

Include the following headers in all API requests:

```
Authorization: Bearer YOUR_API_KEY
X-Tenant-ID: your-tenant-id
Content-Type: application/json
```

### Getting an API Key

1. Sign in to the Driftlock dashboard
2. Navigate to Settings > API Keys
3. Click "Generate New Key"
4. Copy the key and include it in the `Authorization` header

## Endpoints

### Anomalies

#### Get Anomalies

Retrieve detected anomalies with filtering and pagination.

```http
GET /anomalies
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|-------------|-------------|
| limit | integer | No | Maximum number of anomalies to return (default: 20) |
| offset | integer | No | Number of anomalies to skip (default: 0) |
| start_time | string | No | ISO 8601 timestamp to start from |
| end_time | string | No | ISO 8601 timestamp to end at |
| severity | string | No | Filter by severity (low, medium, high, critical) |
| status | string | No | Filter by status (open, investigating, resolved) |

**Response:**

```json
{
  "anomalies": [
    {
      "id": "anomaly_123456",
      "timestamp": "2023-10-25T14:30:00Z",
      "severity": "high",
      "status": "open",
      "confidence": 0.95,
      "type": "statistical",
      "description": "Unusual spike in request rate",
      "indicators": [
        {
          "name": "request_rate",
          "value": 150,
          "threshold": 100,
          "unit": "requests/sec"
        }
      ],
      "evidence": {
        "baseline_period": "2023-10-24T00:00:00Z/2023-10-25T14:00:00Z",
        "compression_ratio": 0.85,
        "p_value": 0.001
      },
      "metadata": {
        "service": "api-gateway",
        "operation": "GET",
        "user_id": "user-12345",
        "session_id": "session-67890"
      }
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

#### Get Anomaly by ID

Retrieve details of a specific anomaly.

```http
GET /anomalies/{id}
```

**Response:**

```json
{
  "id": "anomaly_123456",
  "timestamp": "2023-10-25T14:30:00Z",
  "severity": "high",
  "status": "open",
  "confidence": 0.95,
  "type": "statistical",
  "description": "Unusual spike in request rate",
  "indicators": [
    {
      "name": "request_rate",
      "value": 150,
      "threshold": 100,
      "unit": "requests/sec"
    }
  ],
  "evidence": {
    "baseline_period": "2023-10-24T00:00:00Z/2023-10-25T14:00:00Z",
    "compression_ratio": 0.85,
    "p_value": 0.001
  },
  "metadata": {
    "service": "api-gateway",
    "operation": "GET",
    "user_id": "user-12345",
    "session_id": "session-67890"
  },
  "timeline": [
    {
      "timestamp": "2023-10-25T14:25:00Z",
      "event": "anomaly_detected",
      "description": "Statistical anomaly detected"
    },
    {
      "timestamp": "2023-10-25T14:30:00Z",
      "event": "investigation_started",
      "description": "Investigation initiated"
    }
  ]
}
```

### Events

#### Ingest Events

Submit events for anomaly detection.

```http
POST /events
```

**Request Body:**

```json
{
  "events": [
    {
      "timestamp": "2023-10-25T14:30:00Z",
      "service": "api-gateway",
      "operation": "GET",
      "duration": 150,
      "status": "success",
      "user_id": "user-12345",
      "session_id": "session-67890",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
      "custom_attributes": {
        "region": "us-west-1",
        "environment": "production"
      }
    }
  ]
}
```

**Response:**

```json
{
  "processed": 10,
  "failed": 0,
  "message": "Events processed successfully"
}
```

### Tenants

#### Create Tenant

Create a new tenant for multi-tenant deployments.

```http
POST /tenants
```

**Request Body:**

```json
{
  "name": "Acme Corporation",
  "domain": "acme.driftlock.com",
  "plan": "pro",
  "industry": "finance",
  "company_size": "medium"
}
```

**Response:**

```json
{
  "id": "tenant_123456",
  "name": "Acme Corporation",
  "domain": "acme.driftlock.com",
  "status": "trial",
  "plan": "pro",
  "industry": "finance",
  "company_size": "medium",
  "created_at": "2023-10-25T14:30:00Z",
  "updated_at": "2023-10-25T14:30:00Z",
  "quotas": {
    "max_anomalies_per_day": 2000,
    "max_events_per_day": 20000,
    "max_storage_gb": 200,
    "max_api_requests_per_min": 1200
  },
  "usage": {
    "anomalies_detected": 0,
    "events_processed": 0,
    "storage_used_gb": 0,
    "api_requests_today": 0,
    "last_reset": "2023-10-25T00:00:00Z"
  }
}
```

#### Get Tenant

Retrieve tenant information.

```http
GET /tenants/{id}
```

**Response:**

```json
{
  "id": "tenant_123456",
  "name": "Acme Corporation",
  "domain": "acme.driftlock.com",
  "status": "active",
  "plan": "pro",
  "industry": "finance",
  "company_size": "medium",
  "created_at": "2023-10-25T14:30:00Z",
  "updated_at": "2023-10-25T14:30:00Z",
  "quotas": {
    "max_anomalies_per_day": 2000,
    "max_events_per_day": 20000,
    "max_storage_gb": 200,
    "max_api_requests_per_min": 1200
  },
  "usage": {
    "anomalies_detected": 150,
    "events_processed": 15000,
    "storage_used_gb": 45,
    "api_requests_today": 800,
    "last_reset": "2023-10-25T00:00:00Z"
  }
}
```

### System

#### Health Check

Check system health and status.

```http
GET /healthz
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2023-10-25T14:30:00Z",
  "version": "1.0.0",
  "uptime": "99.9%",
  "checks": {
    "database": "healthy",
    "kafka": "healthy",
    "redis": "healthy",
    "clickhouse": "healthy"
  }
}
```

#### Metrics

Get system metrics for monitoring.

```http
GET /metrics
```

**Response:**

Prometheus-compatible metrics format.

```
# HELP
http_requests_total{method="GET",status="200"} 12345
http_request_duration_seconds_bucket{le="0.1"} 100
http_request_duration_seconds_sum 456.7
http_request_duration_seconds_count 12345
driftlock_anomalies_detected_total 150
driftlock_anomalies_processed_total 148
driftlock_events_processed_total 15000
```

## SDKs and Libraries

### OpenTelemetry

Send events to Driftlock using OpenTelemetry:

```javascript
import { trace } from '@opentelemetry/api';

// Configure OTLP exporter
const provider = new BasicTracerProvider({
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'X-Tenant-ID': 'your-tenant-id'
  }
});

// Send events
const tracer = trace.getTracer('your-service-name');
tracer.startActiveSpan('operation').setAttributes({
  'service.name': 'api-gateway',
  'operation.name': 'GET',
  'user.id': 'user-123'
}).end();
```

### Python Client

```python
import requests

headers = {
    'Authorization': 'Bearer YOUR_API_KEY',
    'X-Tenant-ID': 'your-tenant-id',
    'Content-Type': 'application/json'
}

response = requests.post(
    'https://api.driftlock.com/v1/events',
    json={'events': events},
    headers=headers
)
```

### Go Client

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Event struct {
    Timestamp string `json:"timestamp"`
    Service   string `json:"service"`
    Operation string `json:"operation"`
    Duration  int    `json:"duration"`
    Status    string `json:"status"`
    UserID    string `json:"user_id"`
    SessionID string `json:"session_id"`
}

func SendEvents(events []Event) error {
    jsonData, err := json.Marshal(map[string]interface{}{
        "events": events,
    })
    if err != nil {
        return err
    }

    req, err := http.NewRequest(
        "POST",
        "https://api.driftlock.com/v1/events",
        bytes.NewReader(jsonData),
    )
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer YOUR_API_KEY")
    req.Header.Set("X-Tenant-ID", "your-tenant-id")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err = client.Do(req)
    return err
}
```

## Error Codes

| Code | Description |
|-------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 429 | Too Many Requests |
| 500 | Internal Server Error |

## Rate Limits

| Plan | Events/Day | Anomalies/Day | API Requests/Min |
|-------|-------------|---------------|-------------------|
| Trial | 1,000 | 100 | 60 |
| Starter | 5,000 | 500 | 300 |
| Pro | 20,000 | 2,000 | 1,200 |
| Enterprise | Unlimited | Unlimited | 6,000 |

## Integration Guides

### OpenTelemetry Collector

Configure the OpenTelemetry Collector to send events to Driftlock:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: api.driftlock.com:4317
      http:
        endpoint: api.driftlock.com:4318

processors:
  batch:

exporters:
  otlp:
    endpoint: api.driftlock.com:4317
    headers:
      Authorization: "Bearer YOUR_API_KEY"
      X-Tenant-ID: "your-tenant-id"

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

### Kafka Producer

Send events to Driftlock via Kafka:

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/Shopify/sarama"
)

type Event struct {
    Timestamp string `json:"timestamp"`
    Service   string `json:"service"`
    Operation string `json:"operation"`
    Duration  int    `json:"duration"`
    Status    string `json:"status"`
    UserID    string `json:"user_id"`
    SessionID string `json:"session_id"`
}

func main() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Compression = sarama.CompressionSnappy
    config.Producer.Flush.Frequency = 500
    config.Producer.Flush.Messages = 100

    producer, err := sarama.NewSyncProducer([]string{"driftlock-events"}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    for i := 0; i < 1000; i++ {
        event := Event{
            Timestamp: time.Now().Format(time.RFC3339Nano),
            Service:   "api-gateway",
            Operation: "GET",
            Duration:  100 + (i % 200),
            Status:    "success",
            UserID:    "user-12345",
            SessionID: "session-67890",
        }

        jsonData, err := json.Marshal(event)
        if err != nil {
            log.Printf("Failed to marshal event: %v", err)
            continue
        }

        msg := &sarama.ProducerMessage{
            Topic: "driftlock-events",
            Value: jsonData,
            Headers: []sarama.RecordHeader{
                {Key: "X-Tenant-ID", Value: "your-tenant-id"},
                {Key: "Authorization", Value: "Bearer YOUR_API_KEY"},
            },
        }

        _, _, err := producer.SendMessage(msg)
        if err != nil {
            log.Printf("Failed to send message: %v", err)
        }

        time.Sleep(100 * time.Millisecond)
    }
}
```

## Support

For support and questions:

- Email: support@driftlock.com
- Documentation: https://docs.driftlock.com
- Status Page: https://status.driftlock.com
