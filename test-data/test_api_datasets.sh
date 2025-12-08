#!/bin/bash

# Driftlock API Dataset Testing Script
# Tests Terra Luna, NASA Turbofan, and Web Traffic datasets

API_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect"
RESULTS_DIR="/Volumes/VIXinSSD/driftlock/test-data/api_test_results"
mkdir -p "$RESULTS_DIR"

echo "============================================"
echo "Driftlock API Dataset Testing"
echo "API: $API_URL"
echo "Results: $RESULTS_DIR"
echo "Start Time: $(date)"
echo "============================================"
echo ""

# Function to test API endpoint
test_dataset() {
    local name=$1
    local payload=$2
    local expected_anomalies=$3

    echo "Testing: $name"
    echo "Expected anomalies: $expected_anomalies"

    # Send request and save response
    response=$(echo "$payload" | curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d @- 2>&1)

    # Save raw response
    echo "$response" > "$RESULTS_DIR/${name}_response.json"

    # Parse results
    if echo "$response" | jq . >/dev/null 2>&1; then
        anomaly_count=$(echo "$response" | jq '.anomaly_count // 0' 2>/dev/null || echo "0")
        total_events=$(echo "$response" | jq '.events_analyzed // 0' 2>/dev/null || echo "0")

        echo "  Events analyzed: $total_events"
        echo "  Anomalies detected: $anomaly_count"

        # Save summary
        echo "{\"dataset\":\"$name\",\"events\":$total_events,\"anomalies\":$anomaly_count,\"expected\":\"$expected_anomalies\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" \
            > "$RESULTS_DIR/${name}_summary.json"

        # Verdict
        if [ "$anomaly_count" -gt 0 ]; then
            echo "  Status: PASS - Anomalies detected"
            return 0
        else
            echo "  Status: WARN - No anomalies detected (expected: $expected_anomalies)"
            return 1
        fi
    else
        echo "  Status: FAIL - Invalid JSON response"
        echo "  Response: $response"
        return 2
    fi
}

# ==============================================
# TEST 1: Terra Luna Crash Data (May 2022)
# ==============================================
echo ""
echo "TEST 1: Terra Luna Crypto Crash"
echo "----------------------------------------------"

# Read Terra Luna CSV and convert to JSON events
# Focus on the crash period (May 7-13, 2022)
terra_payload=$(awk -F, 'NR>1 && NR<=200 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"date\":\"%s\",\"price\":%s,\"asset\":\"LUNA\"},", $1, $2, $3
}' /Volumes/VIXinSSD/driftlock/test-data/terra_luna/terra-luna.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

test_dataset "terra_luna_crash" "$terra_payload" "High anomalies during May 7-13 crash"
echo ""

# ==============================================
# TEST 2: NASA Turbofan Sensor Data
# ==============================================
echo "TEST 2: NASA Turbofan Sensor Degradation"
echo "----------------------------------------------"

# Convert NASA turbofan data to JSON
# Each row has 26 columns: unit_id, cycle, 3 operational settings, 21 sensor readings
nasa_payload=$(awk 'NR<=100 {
    cycle=$2;
    sensor1=$4; sensor2=$5; sensor3=$6; sensor4=$7;
    printf "{\"cycle\":%d,\"sensor1\":%s,\"sensor2\":%s,\"sensor3\":%s,\"sensor4\":%s},", cycle, sensor1, sensor2, sensor3, sensor4
}' /Volumes/VIXinSSD/driftlock/test-data/nasa_turbofan/CMaps/train_FD001.txt | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

test_dataset "nasa_turbofan_degradation" "$nasa_payload" "Anomalies as sensors degrade"
echo ""

# ==============================================
# TEST 3: Web Traffic - AWS CloudWatch
# ==============================================
echo "TEST 3: AWS CloudWatch CPU Utilization"
echo "----------------------------------------------"

# Convert AWS CloudWatch CSV to JSON
web_payload=$(awk -F, 'NR>1 && NR<=150 {
    gsub(/"/, "", $0);
    printf "{\"timestamp\":\"%s\",\"cpu_utilization\":%s,\"metric\":\"ec2_cpu\"},", $1, $2
}' /Volumes/VIXinSSD/driftlock/test-data/web_traffic/realAWSCloudwatch/realAWSCloudwatch/ec2_cpu_utilization_24ae8d.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

test_dataset "aws_cloudwatch_cpu" "$web_payload" "CPU spike anomalies"
echo ""

# ==============================================
# TEST 4: Known Anomaly - NYC Taxi
# ==============================================
echo "TEST 4: NYC Taxi (Known Anomaly Dataset)"
echo "----------------------------------------------"

if [ -f "/Volumes/VIXinSSD/driftlock/test-data/web_traffic/realKnownCause/realKnownCause/nyc_taxi.csv" ]; then
    nyc_payload=$(awk -F, 'NR>1 && NR<=200 {
        gsub(/"/, "", $0);
        printf "{\"timestamp\":\"%s\",\"value\":%s,\"metric\":\"taxi_count\"},", $1, $2
    }' /Volumes/VIXinSSD/driftlock/test-data/web_traffic/realKnownCause/realKnownCause/nyc_taxi.csv | sed 's/,$//' | awk '{print "{\"events\":[" $0 "]}"}')

    test_dataset "nyc_taxi_known_anomaly" "$nyc_payload" "Known anomaly events"
else
    echo "  NYC Taxi dataset not found - SKIPPED"
fi
echo ""

# ==============================================
# Summary Report
# ==============================================
echo "============================================"
echo "Test Summary"
echo "============================================"
echo ""

for summary in "$RESULTS_DIR"/*_summary.json; do
    if [ -f "$summary" ]; then
        dataset=$(jq -r '.dataset' "$summary")
        events=$(jq -r '.events' "$summary")
        anomalies=$(jq -r '.anomalies' "$summary")
        expected=$(jq -r '.expected' "$summary")

        printf "%-30s | Events: %4s | Anomalies: %4s | Expected: %s\n" \
            "$dataset" "$events" "$anomalies" "$expected"
    fi
done

echo ""
echo "Detailed results saved to: $RESULTS_DIR"
echo "End Time: $(date)"
echo "============================================"
