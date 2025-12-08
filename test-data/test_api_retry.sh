#!/bin/bash

# Retry testing with smaller payloads
API_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect"
RESULTS_DIR="/Volumes/VIXinSSD/driftlock/test-data/api_test_results"

echo "============================================"
echo "Retry Testing with Smaller Payloads"
echo "============================================"
echo ""

# Test Terra Luna with smaller payload (50 events)
echo "TEST 1: Terra Luna Crypto Crash (50 events)"
echo "----------------------------------------------"

terra_payload=$(awk -F, 'NR>1 && NR<=51 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"date\":\"%s\",\"price\":%s,\"asset\":\"LUNA\"},", $1, $2, $3
}' /Volumes/VIXinSSD/driftlock/test-data/terra_luna/terra-luna.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

response=$(echo "$terra_payload" | curl -s -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d @-)

echo "Response: $response" | jq '.' > "$RESULTS_DIR/terra_luna_retry_response.json"

if echo "$response" | jq . >/dev/null 2>&1; then
    anomaly_count=$(echo "$response" | jq '.anomaly_count // 0')
    total_events=$(echo "$response" | jq '.events_analyzed // 0')
    echo "  Events analyzed: $total_events"
    echo "  Anomalies detected: $anomaly_count"
else
    echo "  Error: $response"
fi
echo ""

# Test AWS CloudWatch with smaller payload (50 events)
echo "TEST 2: AWS CloudWatch CPU (50 events)"
echo "----------------------------------------------"

web_payload=$(awk -F, 'NR>1 && NR<=51 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"cpu_utilization\":%s,\"metric\":\"ec2_cpu\"},", $1, $2
}' /Volumes/VIXinSSD/driftlock/test-data/web_traffic/realAWSCloudwatch/realAWSCloudwatch/ec2_cpu_utilization_24ae8d.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

response=$(echo "$web_payload" | curl -s -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d @-)

echo "Response: $response" | jq '.' > "$RESULTS_DIR/aws_cloudwatch_retry_response.json"

if echo "$response" | jq . >/dev/null 2>&1; then
    anomaly_count=$(echo "$response" | jq '.anomaly_count // 0')
    total_events=$(echo "$response" | jq '.events_analyzed // 0')
    echo "  Events analyzed: $total_events"
    echo "  Anomalies detected: $anomaly_count"
else
    echo "  Error: $response"
fi
echo ""

echo "Testing complete!"
