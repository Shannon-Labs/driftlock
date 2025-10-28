#!/bin/bash

# DriftLock Test Data Generator
# Generates realistic test data for anomaly detection

set -e

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
JWT_TOKEN="${JWT_TOKEN:-}"

if [ -z "$JWT_TOKEN" ]; then
    echo "Error: JWT_TOKEN environment variable must be set"
    echo "Usage: JWT_TOKEN=your_token_here ./generate-test-data.sh"
    exit 1
fi

echo "========================================"
echo "DriftLock Test Data Generator"
echo "========================================"
echo "API URL: $API_BASE_URL"
echo ""

# Function to generate random log message
generate_log_message() {
    local severity=$1
    local messages=(
        "User authentication successful"
        "Database query executed"
        "API request processed"
        "Cache hit for key"
        "Background job completed"
        "File uploaded successfully"
        "Session created"
        "Email sent"
        "Payment processed"
        "Configuration updated"
    )

    local error_messages=(
        "Connection timeout"
        "Database query failed"
        "Authentication failed"
        "Rate limit exceeded"
        "Invalid request payload"
        "Service unavailable"
        "Memory allocation failed"
        "Disk space low"
    )

    if [ "$severity" == "error" ]; then
        echo "${error_messages[$RANDOM % ${#error_messages[@]}]}"
    else
        echo "${messages[$RANDOM % ${#messages[@]}]}"
    fi
}

# Function to generate metric value
generate_metric_value() {
    local metric_type=$1
    case $metric_type in
        "cpu")
            echo "scale=2; $RANDOM % 100" | bc
            ;;
        "memory")
            echo "scale=2; $RANDOM % 16000" | bc
            ;;
        "latency")
            echo "scale=2; $RANDOM % 1000" | bc
            ;;
        "requests")
            echo $((RANDOM % 10000))
            ;;
    esac
}

# Generate normal log events
echo "Generating normal log events..."
for i in {1..50}; do
    timestamp=$(date -u -d "-$((RANDOM % 3600)) seconds" +%Y-%m-%dT%H:%M:%SZ)
    message=$(generate_log_message "info")

    event_data='[{
        "timestamp": "'$timestamp'",
        "stream_type": "logs",
        "data": "'$message'",
        "metadata": {
            "source": "test-generator",
            "severity": "info",
            "service": "api-server"
        }
    }]'

    curl -s -X POST "$API_BASE_URL/api/v1/events/ingest" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "$event_data" > /dev/null

    echo -n "."
done
echo " ✓ Generated 50 normal log events"

# Generate some anomalous log events (errors)
echo "Generating anomalous log events..."
for i in {1..10}; do
    timestamp=$(date -u -d "-$((RANDOM % 1800)) seconds" +%Y-%m-%dT%H:%M:%SZ)
    message=$(generate_log_message "error")

    event_data='[{
        "timestamp": "'$timestamp'",
        "stream_type": "logs",
        "data": "'$message'",
        "metadata": {
            "source": "test-generator",
            "severity": "error",
            "service": "api-server"
        }
    }]'

    curl -s -X POST "$API_BASE_URL/api/v1/events/ingest" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "$event_data" > /dev/null

    echo -n "."
done
echo " ✓ Generated 10 anomalous log events"

# Generate metric events
echo "Generating metric events..."
for i in {1..30}; do
    timestamp=$(date -u -d "-$((RANDOM % 3600)) seconds" +%Y-%m-%dT%H:%M:%SZ)

    for metric_type in "cpu" "memory" "latency"; do
        value=$(generate_metric_value "$metric_type")

        event_data='[{
            "timestamp": "'$timestamp'",
            "stream_type": "metrics",
            "data": "{'$metric_type': '$value'}",
            "metadata": {
                "source": "test-generator",
                "metric_type": "'$metric_type'",
                "unit": "percent"
            }
        }]'

        curl -s -X POST "$API_BASE_URL/api/v1/events/ingest" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d "$event_data" > /dev/null
    done

    echo -n "."
done
echo " ✓ Generated 90 metric events"

# Generate trace events
echo "Generating trace events..."
for i in {1..20}; do
    timestamp=$(date -u -d "-$((RANDOM % 1800)) seconds" +%Y-%m-%dT%H:%M:%SZ)
    duration=$((RANDOM % 1000))

    event_data='[{
        "timestamp": "'$timestamp'",
        "stream_type": "traces",
        "data": "span_id='$(uuidgen)' duration='$duration'ms",
        "metadata": {
            "source": "test-generator",
            "trace_id": "'$(uuidgen)'",
            "span_id": "'$(uuidgen)'",
            "duration_ms": '$duration',
            "service": "api-server"
        }
    }]'

    curl -s -X POST "$API_BASE_URL/api/v1/events/ingest" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -d "$event_data" > /dev/null

    echo -n "."
done
echo " ✓ Generated 20 trace events"

echo ""
echo "========================================"
echo "Test Data Generation Complete"
echo "========================================"
echo "Total events generated: 170"
echo "  - Normal logs: 50"
echo "  - Anomalous logs: 10"
echo "  - Metrics: 90"
echo "  - Traces: 20"
echo ""
echo "Check the dashboard for anomalies!"
echo "Visit: http://localhost:3000/dashboard"
