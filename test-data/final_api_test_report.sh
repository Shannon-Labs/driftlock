#!/bin/bash

# Final Comprehensive API Test Report
API_URL="https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect"
RESULTS_DIR="/Volumes/VIXinSSD/driftlock/test-data/final_test_results"
mkdir -p "$RESULTS_DIR"

echo "============================================"
echo "DRIFTLOCK API DATASET TEST REPORT"
echo "============================================"
echo "API Endpoint: $API_URL"
echo "Test Date: $(date)"
echo ""

# Function to run test
run_test() {
    local name=$1
    local description=$2
    local payload_file=$3
    local expected=$4

    echo "TEST: $name"
    echo "Description: $description"
    echo "Expected: $expected"

    start_time=$(date +%s)
    response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d @"$payload_file")
    end_time=$(date +%s)

    echo "$response" > "$RESULTS_DIR/${name}_response.json"

    if echo "$response" | jq . >/dev/null 2>&1; then
        success=$(echo "$response" | jq -r '.success')
        total=$(echo "$response" | jq -r '.total_events')
        anomalies=$(echo "$response" | jq -r '.anomaly_count')
        proc_time=$(echo "$response" | jq -r '.processing_time')
        algo=$(echo "$response" | jq -r '.compression_algo')

        echo "  Result: SUCCESS=$success"
        echo "  Events Analyzed: $total"
        echo "  Anomalies Detected: $anomalies"
        echo "  Processing Time: $proc_time"
        echo "  Algorithm: $algo"
        echo "  Request Duration: $((end_time - start_time))s"

        if [ "$anomalies" -gt 0 ]; then
            echo "  Verdict: PASS - Anomalies detected"
            return 0
        else
            echo "  Verdict: WARN - No anomalies (expected: $expected)"
            return 1
        fi
    else
        echo "  Result: FAILED"
        echo "  Error: $response"
        echo "  Verdict: FAIL"
        return 2
    fi
}

# ==============================================
# TEST 1: NASA Turbofan Degradation
# ==============================================
echo ""
echo "============================================"
echo "1. NASA TURBOFAN SENSOR DEGRADATION"
echo "============================================"

cat > "$RESULTS_DIR/nasa_payload.json" << 'EOF'
{
  "events": [
    {"cycle":1,"sensor1":-0.0007,"sensor2":100,"sensor3":518.67,"sensor4":641.82},
    {"cycle":10,"sensor1":-0.0004,"sensor2":100,"sensor3":518.67,"sensor4":642.24},
    {"cycle":20,"sensor1":0.0005,"sensor2":100,"sensor3":518.67,"sensor4":642.42},
    {"cycle":30,"sensor1":-0.0002,"sensor2":100,"sensor3":518.67,"sensor4":642.61},
    {"cycle":40,"sensor1":-0.0004,"sensor2":100,"sensor3":518.67,"sensor4":642.24},
    {"cycle":50,"sensor1":-0.0001,"sensor2":100,"sensor3":518.67,"sensor4":642.77},
    {"cycle":60,"sensor1":0.0002,"sensor2":100,"sensor3":518.67,"sensor4":642.53},
    {"cycle":70,"sensor1":-0,"sensor2":100,"sensor3":518.67,"sensor4":642.48},
    {"cycle":80,"sensor1":-0.0004,"sensor2":100,"sensor3":518.67,"sensor4":642.35},
    {"cycle":90,"sensor1":0.0007,"sensor2":100,"sensor3":518.67,"sensor4":642.23},
    {"cycle":100,"sensor1":-0.0003,"sensor2":100,"sensor3":518.67,"sensor4":642.11}
  ]
}
EOF

run_test "nasa_turbofan" "Turbofan engine sensor degradation over time" \
    "$RESULTS_DIR/nasa_payload.json" "Sensor degradation anomalies"

echo ""

# ==============================================
# TEST 2: Terra Luna Crash
# ==============================================
echo "============================================"
echo "2. TERRA LUNA CRYPTO CRASH (May 2022)"
echo "============================================"

# Create payload with crash progression: $80 -> $60 -> $30 -> $1
cat > "$RESULTS_DIR/terra_payload.json" << 'EOF'
{
  "events": [
    {"timestamp":"1651852800","date":"2022-05-06T12:00:00","price":79.97},
    {"timestamp":"1651860000","date":"2022-05-06T14:00:00","price":79.98},
    {"timestamp":"1651870800","date":"2022-05-06T17:00:00","price":79.86},
    {"timestamp":"1651881600","date":"2022-05-06T20:00:00","price":81.32},
    {"timestamp":"1651894800","date":"2022-05-07T00:00:00","price":63.19},
    {"timestamp":"1651902000","date":"2022-05-07T02:00:00","price":60.76},
    {"timestamp":"1651913700","date":"2022-05-07T05:15:00","price":59.26},
    {"timestamp":"1651924500","date":"2022-05-07T08:15:00","price":60.77},
    {"timestamp":"1651935300","date":"2022-05-07T11:15:00","price":55.61},
    {"timestamp":"1651946100","date":"2022-05-07T14:15:00","price":52.73},
    {"timestamp":"1651981200","date":"2022-05-08T00:00:00","price":32.09},
    {"timestamp":"1651992000","date":"2022-05-08T03:00:00","price":26.78},
    {"timestamp":"1652004700","date":"2022-05-08T06:45:00","price":29.06},
    {"timestamp":"1652025300","date":"2022-05-08T12:45:00","price":23.97},
    {"timestamp":"1652068500","date":"2022-05-09T00:45:00","price":5.04},
    {"timestamp":"1652090100","date":"2022-05-09T06:45:00","price":0.8826},
    {"timestamp":"1652176500","date":"2022-05-10T06:45:00","price":0.1828},
    {"timestamp":"1652262900","date":"2022-05-11T06:45:00","price":0.8096},
    {"timestamp":"1652279400","date":"2022-05-11T10:30:00","price":7.1805},
    {"timestamp":"1652295600","date":"2022-05-11T15:00:00","price":1.3},
    {"timestamp":"1652349300","date":"2022-05-12T05:55:00","price":0.2044},
    {"timestamp":"1652435700","date":"2022-05-13T05:55:00","price":0.0002044}
  ]
}
EOF

run_test "terra_luna_crash" "Terra Luna collapse from $80 to near $0" \
    "$RESULTS_DIR/terra_payload.json" "High anomalies during crash"

echo ""

# ==============================================
# TEST 3: AWS CloudWatch CPU Spikes
# ==============================================
echo "============================================"
echo "3. AWS CLOUDWATCH EC2 CPU UTILIZATION"
echo "============================================"

cat > "$RESULTS_DIR/aws_payload.json" << 'EOF'
{
  "events": [
    {"timestamp":"2014-02-14 14:30:00","value":0.132},
    {"timestamp":"2014-02-14 14:35:00","value":0.134},
    {"timestamp":"2014-02-14 14:40:00","value":0.134},
    {"timestamp":"2014-02-14 14:45:00","value":0.134},
    {"timestamp":"2014-02-14 15:00:00","value":0.134},
    {"timestamp":"2014-02-14 15:05:00","value":0.134},
    {"timestamp":"2014-02-14 15:10:00","value":0.066},
    {"timestamp":"2014-02-14 15:15:00","value":0.066},
    {"timestamp":"2014-02-14 15:20:00","value":0.066},
    {"timestamp":"2014-02-14 15:25:00","value":0.132},
    {"timestamp":"2014-02-14 15:30:00","value":0.132},
    {"timestamp":"2014-02-14 15:35:00","value":0.132},
    {"timestamp":"2014-02-14 15:40:00","value":0.132},
    {"timestamp":"2014-02-14 15:45:00","value":0.132},
    {"timestamp":"2014-02-14 15:50:00","value":0.132},
    {"timestamp":"2014-02-14 15:55:00","value":0.132},
    {"timestamp":"2014-02-14 16:00:00","value":0.066},
    {"timestamp":"2014-02-14 16:05:00","value":0.066},
    {"timestamp":"2014-02-14 16:10:00","value":0.132},
    {"timestamp":"2014-02-14 16:15:00","value":0.132}
  ]
}
EOF

run_test "aws_cloudwatch" "EC2 CPU utilization with fluctuations" \
    "$RESULTS_DIR/aws_payload.json" "CPU spike anomalies"

echo ""

# ==============================================
# SUMMARY
# ==============================================
echo "============================================"
echo "TEST SUMMARY"
echo "============================================"
echo ""

total_tests=3
passed=0
failed=0

for result in "$RESULTS_DIR"/*_response.json; do
    if [ -f "$result" ]; then
        name=$(basename "$result" "_response.json")
        success=$(jq -r '.success // "false"' "$result" 2>/dev/null)
        anomalies=$(jq -r '.anomaly_count // 0' "$result" 2>/dev/null)

        if [ "$success" = "true" ]; then
            if [ "$anomalies" -gt 0 ]; then
                status="PASS"
                passed=$((passed + 1))
            else
                status="WARN"
            fi
        else
            status="FAIL"
            failed=$((failed + 1))
        fi

        printf "%-25s | Status: %-4s | Anomalies: %3s\n" "$name" "$status" "$anomalies"
    fi
done

echo ""
echo "Total Tests: $total_tests"
echo "Passed: $passed"
echo "Failed: $failed"
echo ""
echo "Results saved to: $RESULTS_DIR"
echo "============================================"
