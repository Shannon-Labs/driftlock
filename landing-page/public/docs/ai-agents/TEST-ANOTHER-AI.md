# TEST PROMPT FOR ANOTHER AI AGENT

## 🎯 **Your Mission: Verify Driftlock Actually Works**

You are an independent AI agent. Your job is to test whether Driftlock's compression-based anomaly detection (CBAD) actually works as advertised.

**Repository:** `/Volumes/VIXinSSD/driftlock`

## 📋 **Test Protocol**

### **Step 1: Environment Setup**
```bash
cd /Volumes/VIXinSSD/driftlock

# Check if Driftlock is already running
if curl -s http://localhost:8080/healthz > /dev/null; then
    echo "Driftlock is running"
else
    echo "Starting Driftlock..."
    ./start.sh &
    sleep 30
fi
```

### **Step 2: Verify System Health**
```bash
# Check API health
curl http://localhost:8080/healthz
# EXPECTED: "ok"

# Check version
curl http://localhost:8080/v1/version | jq .
# EXPECTED: {"version":"..."}

# Check dashboard
curl -s http://localhost:3000 > /dev/null && echo "Dashboard reachable"
# EXPECTED: "Dashboard reachable"
```

### **Step 3: Run Tests**

**Test A: Normal Data (Should NOT trigger anomalies)**
```bash
# Ingest normal transactions
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/normal-transactions.jsonl

# Wait 10 seconds
sleep 10

# Count anomalies
ANOMALY_COUNT=$(curl -s http://localhost:8080/v1/anomalies | jq '.anomalies | length')
echo "Anomalies detected in normal data: $ANOMALY_COUNT"

# PASS if: $ANOMALY_COUNT < 5 (false positive rate < 1%)
if [ "$ANOMALY_COUNT" -lt 5 ]; then
    echo "✅ TEST A PASSED"
else
    echo "❌ TEST A FAILED: Too many false positives"
fi
```

**Test B: Anomalous Data (SHOULD trigger anomalies)**
```bash
# Ingest anomalous transactions
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/anomalous-transactions.jsonl

# Wait 10 seconds
sleep 10

# Count anomalies
ANOMALY_COUNT=$(curl -s http://localhost:8080/v1/anomalies | jq '.anomalies | length')
echo "Anomalies detected in anomalous data: $ANOMALY_COUNT"

# PASS if: $ANOMALY_COUNT > 80 (true positive rate > 80%)
if [ "$ANOMALY_COUNT" -gt 80 ]; then
    echo "✅ TEST B PASSED"
else
    echo "❌ TEST B FAILED: Too many false negatives"
fi
```

**Test C: Mixed Data (Should detect ~5% anomalies)**
```bash
# Ingest mixed transactions
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/mixed-transactions.jsonl

# Wait 10 seconds
sleep 10

# Count anomalies
ANOMALY_COUNT=$(curl -s http://localhost:8080/v1/anomalies | jq '.anomalies | length')
echo "Anomalies detected in mixed data: $ANOMALY_COUNT"

# PASS if: 45 <= $ANOMALY_COUNT <= 55 (95% recall ± tolerance)
if [ "$ANOMALY_COUNT" -ge 45 ] && [ "$ANOMALY_COUNT" -le 55 ]; then
    echo "✅ TEST C PASSED"
else
    echo "❌ TEST C FAILED: Detection rate outside expected range"
fi
```

**Test D: Glass-Box Explanations**
```bash
# Get first anomaly
ANOMALY_ID=$(curl -s http://localhost:8080/v1/anomalies | jq -r '.anomalies[0].id')

if [ "$ANOMALY_ID" != "null" ] && [ -n "$ANOMALY_ID" ]; then
    # Get anomaly details
    curl http://localhost:8080/v1/anomalies/$ANOMALY_ID | jq .
    
    # Check for required fields
    HAS_NCD=$(curl -s http://localhost:8080/v1/anomalies/$ANOMALY_ID | jq '.ncd_score != null')
    HAS_EXPLANATION=$(curl -s http://localhost:8080/v1/anomalies/$ANOMALY_ID | jq '.glass_box_explanation != null')
    HAS_CONFIDENCE=$(curl -s http://localhost:8080/v1/anomalies/$ANOMALY_ID | jq '.confidence_level != null')
    
    if [ "$HAS_NCD" = "true" ] && [ "$HAS_EXPLANATION" = "true" ] && [ "$HAS_CONFIDENCE" = "true" ]; then
        echo "✅ TEST D PASSED: Glass-box explanations working"
    else
        echo "❌ TEST D FAILED: Missing explanation fields"
    fi
else
    echo "⚠️  No anomalies to test explanations"
fi
```

### **Step 4: Performance Test**
```bash
# Measure time for 10K transactions
echo "Testing performance with 10K transactions..."

time_start=$(date +%s)

# Generate 10K transactions
for i in {1..10000}; do
    echo '{"timestamp":"2025-11-10T10:00:00Z","user_id":"user_'$i'","amount":'$((RANDOM % 500 + 10))'.99,"merchant":"Test","location":"US","payment_method":"credit","device":"mobile","transaction_type":"purchase","metadata":{"category":"test","currency":"USD"}}' >> /tmp/bulk-test.jsonl
done

# Ingest
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @/tmp/bulk-test.jsonl

time_end=$(date +%s)
time_taken=$((time_end - time_start))

echo "Time taken: ${time_taken} seconds"

# PASS if: time_taken < 30 seconds
if [ "$time_taken" -lt 30 ]; then
    echo "✅ PERFORMANCE TEST PASSED"
else
    echo "❌ PERFORMANCE TEST FAILED: Too slow"
fi

# Cleanup
rm /tmp/bulk-test.jsonl
```

### **Step 5: Memory Check**
```bash
# Check for memory leaks
echo "Checking memory usage..."

# Get container ID
CONTAINER_ID=$(docker ps | grep driftlock-api | awk '{print $1}')

if [ -n "$CONTAINER_ID" ]; then
    # Check memory usage
    MEMORY_USAGE=$(docker stats --no-stream --format "table {{.MemUsage}}" $CONTAINER_ID | tail -1 | awk '{print $1}' | sed 's/MiB//')
    
    echo "Memory usage: ${MEMORY_USAGE}MiB"
    
    # PASS if: MEMORY_USAGE < 500MiB
    if (( $(echo "$MEMORY_USAGE < 500" | bc -l) )); then
        echo "✅ MEMORY TEST PASSED"
    else
        echo "❌ MEMORY TEST FAILED: Too high"
    fi
else
    echo "⚠️  Could not find container to check memory"
fi
```

## 📊 **Test Data Provided**

All test data is in `/Volumes/VIXinSSD/driftlock/test-data/`:

- **normal-transactions.jsonl** (500 transactions): Should NOT trigger anomalies
- **anomalous-transactions.jsonl** (100 transactions): SHOULD trigger anomalies  
- **mixed-transactions.jsonl** (1000 transactions, 5% anomalous): Should detect ~50 anomalies
- **README.md**: Detailed test instructions

## ✅ **Success Criteria**

**All tests must pass:**
1. ✅ Normal data: < 5 anomalies detected
2. ✅ Anomalous data: > 80 anomalies detected
3. ✅ Mixed data: 45-55 anomalies detected
4. ✅ Glass-box explanations present for all anomalies
5. ✅ Performance: 10K transactions in < 30 seconds
6. ✅ Memory: < 500MiB usage

## 📝 **Your Report Should Include:**

1. **Test Results Summary:** Pass/Fail for each test
2. **Sample Anomaly:** Full JSON output from one anomaly
3. **Performance Metrics:** Time and memory measurements
4. **Code Review:** Did CBAD algorithm actually run? (check logs)
5. **Final Verdict:** Does Driftlock work as claimed?

## 🚨 **Critical Red Flags**

If ANY of these happen, Driftlock is broken:

- ❌ No anomalies detected (algorithm not working)
- ❌ All transactions flagged as anomalies (threshold broken)
- ❌ No NCD scores in output (compression not calculated)
- ❌ No glass-box explanations (explanation generation broken)
- ❌ Memory > 1GB (memory leak)
- ❌ Crashes on >1000 transactions (performance unacceptable)

## 🔧 **If Tests Fail**

**Check logs:**
```bash
docker logs driftlock-api              # API server logs
docker logs driftlock-collector        # OTEL collector logs
docker logs postgres                   # Database logs
```

**Common issues:**
- Database not ready: Wait 30 seconds after `./start.sh`
- API key wrong: Check `.env` file
- Port conflicts: `lsof -i :8080`, `lsof -i :3000`
- Rust library not built: `cd cbad-core && cargo build --release`

## 🎯 **Final Verdict**

**If all tests pass:**
- ✅ Driftlock works as advertised
- ✅ CBAD algorithm is functional
- ✅ Glass-box explanations work
- ✅ Performance is acceptable
- ✅ Ready for YC application

**If tests fail:**
- ❌ Document exactly what failed
- ❌ Check if it's a configuration issue or algorithm bug
- ❌ Report findings to repository maintainers

---

**Test Duration:** ~30 minutes
**Test Data Size:** ~400KB total
**Expected Result:** Driftlock detects anomalies with 95% recall