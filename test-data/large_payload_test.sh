#!/bin/bash

# Test with larger payloads (100+ events)
API_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect"
RESULTS_DIR="/Volumes/VIXinSSD/driftlock/test-data/final_test_results"

echo "============================================"
echo "LARGE PAYLOAD TESTS (100+ events)"
echo "============================================"
echo ""

# ==============================================
# TEST 1: Terra Luna with 100 events covering crash
# ==============================================
echo "TEST 1: Terra Luna Crash (100 events)"
echo "----------------------------------------------"

terra_payload=$(awk -F, 'NR>1 && NR<=101 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"price\":%s},", $1, $3
}' /Volumes/VIXinSSD/driftlock/test-data/terra_luna/terra-luna.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

response=$(echo "$terra_payload" | curl -s -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d @- 2>&1)

echo "$response" > "$RESULTS_DIR/terra_100_response.json"

if echo "$response" | jq . >/dev/null 2>&1; then
    success=$(echo "$response" | jq -r '.success')
    total=$(echo "$response" | jq -r '.total_events')
    anomalies=$(echo "$response" | jq -r '.anomaly_count')
    proc_time=$(echo "$response" | jq -r '.processing_time')

    echo "  Success: $success"
    echo "  Events: $total"
    echo "  Anomalies: $anomalies"
    echo "  Time: $proc_time"

    if [ "$anomalies" -gt 0 ]; then
        echo "  Verdict: PASS"
        # Show sample anomalies
        echo "$response" | jq -r '.anomalies[0:3] | .[] | "    - Event \(.index): price=\(.event.price) (NCD=\(.metrics.NCD | tostring | .[0:5]))"' 2>/dev/null
    else
        echo "  Verdict: WARN - No anomalies"
    fi
else
    echo "  ERROR: $response"
fi
echo ""

# ==============================================
# TEST 2: Terra Luna focusing on May 7-11 (crash period)
# ==============================================
echo "TEST 2: Terra Luna May 7-11 Crash Period (rows 50-250)"
echo "----------------------------------------------"

terra_crash_payload=$(awk -F, 'NR>50 && NR<=150 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"price\":%s},", $1, $3
}' /Volumes/VIXinSSD/driftlock/test-data/terra_luna/terra-luna.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

response=$(echo "$terra_crash_payload" | curl -s -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d @- 2>&1)

echo "$response" > "$RESULTS_DIR/terra_crash_period_response.json"

if echo "$response" | jq . >/dev/null 2>&1; then
    success=$(echo "$response" | jq -r '.success')
    total=$(echo "$response" | jq -r '.total_events')
    anomalies=$(echo "$response" | jq -r '.anomaly_count')
    proc_time=$(echo "$response" | jq -r '.processing_time')

    echo "  Success: $success"
    echo "  Events: $total"
    echo "  Anomalies: $anomalies"
    echo "  Time: $proc_time"

    if [ "$anomalies" -gt 0 ]; then
        echo "  Verdict: PASS"
        # Show sample anomalies
        echo "$response" | jq -r '.anomalies[0:3] | .[] | "    - Event \(.index): price=\(.event.price) (NCD=\(.metrics.NCD | tostring | .[0:5]))"' 2>/dev/null
    else
        echo "  Verdict: WARN - No anomalies"
    fi
else
    echo "  ERROR: $response"
fi
echo ""

# ==============================================
# TEST 3: AWS CloudWatch with 100 events
# ==============================================
echo "TEST 3: AWS CloudWatch CPU (100 events)"
echo "----------------------------------------------"

aws_payload=$(awk -F, 'NR>1 && NR<=101 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"value\":%s},", $1, $2
}' /Volumes/VIXinSSD/driftlock/test-data/web_traffic/realAWSCloudwatch/realAWSCloudwatch/ec2_cpu_utilization_24ae8d.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

response=$(echo "$aws_payload" | curl -s -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d @- 2>&1)

echo "$response" > "$RESULTS_DIR/aws_100_response.json"

if echo "$response" | jq . >/dev/null 2>&1; then
    success=$(echo "$response" | jq -r '.success')
    total=$(echo "$response" | jq -r '.total_events')
    anomalies=$(echo "$response" | jq -r '.anomaly_count')
    proc_time=$(echo "$response" | jq -r '.processing_time')

    echo "  Success: $success"
    echo "  Events: $total"
    echo "  Anomalies: $anomalies"
    echo "  Time: $proc_time"

    if [ "$anomalies" -gt 0 ]; then
        echo "  Verdict: PASS"
        # Show sample anomalies
        echo "$response" | jq -r '.anomalies[0:3] | .[] | "    - Event \(.index): value=\(.event.value) (NCD=\(.metrics.NCD | tostring | .[0:5]))"' 2>/dev/null
    else
        echo "  Verdict: WARN - No anomalies"
    fi
else
    echo "  ERROR: $response"
fi
echo ""

# ==============================================
# SUMMARY
# ==============================================
echo "============================================"
echo "SUMMARY"
echo "============================================"

for file in terra_100 terra_crash_period aws_100; do
    if [ -f "$RESULTS_DIR/${file}_response.json" ]; then
        anomalies=$(jq -r '.anomaly_count // 0' "$RESULTS_DIR/${file}_response.json" 2>/dev/null)
        events=$(jq -r '.total_events // 0' "$RESULTS_DIR/${file}_response.json" 2>/dev/null)
        printf "%-25s | Events: %3s | Anomalies: %3s\n" "$file" "$events" "$anomalies"
    fi
done

echo ""
echo "Complete!"
