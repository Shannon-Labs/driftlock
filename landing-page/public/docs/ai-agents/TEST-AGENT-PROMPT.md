# DRIFTLOCK TEST PROMPT - For Independent AI Agent

## 🎯 **Your Mission**

Verify that Driftlock's compression-based anomaly detection (CBAD) actually works by:
1. Starting the system
2. Ingesting synthetic transaction data
3. Confirming anomalies are detected
4. Validating glass-box explanations

## 📋 **Test Environment Setup**

```bash
cd /Volumes/VIXinSSD/driftlock

# 1. Start Driftlock
./start.sh

# Wait 30 seconds for full startup
curl http://localhost:8080/healthz  # Should return "ok"

# 2. Ingest test data
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/synthetic-transactions.jsonl

# 3. Check for anomalies
curl http://localhost:8080/v1/anomalies | jq .

# 4. Verify explanations exist
curl http://localhost:8080/v1/anomalies/<anomaly-id> | jq .
```

## 📊 **Test Data Files Provided**

### **File 1: `test-data/normal-transactions.jsonl`**
- 500 normal financial transactions
- Pattern: Regular amounts ($10-$500), common merchants
- Expected: Should compress well, NO anomalies

### **File 2: `test-data/anomalous-transactions.jsonl`**
- 100 anomalous transactions
- Pattern: High amounts ($50k+), unusual merchants, odd locations
- Expected: Should NOT compress well, SHOULD trigger anomalies

### **File 3: `test-data/mixed-transactions.jsonl`**
- 1000 transactions (950 normal, 50 anomalous)
- Expected: Should detect ~50 anomalies at 95% recall

## ✅ **Success Criteria**

**Test 1: Normal Data**
- [ ] Ingest 500 normal transactions
- [ ] Wait 10 seconds
- [ ] Query `/v1/anomalies`
- [ ] **EXPECTED:** < 5 anomalies (false positive rate < 1%)

**Test 2: Anomalous Data**
- [ ] Ingest 100 anomalous transactions
- [ ] Wait 10 seconds
- [ ] Query `/v1/anomalies`
- [ ] **EXPECTED:** > 80 anomalies detected (true positive rate > 80%)

**Test 3: Mixed Data**
- [ ] Ingest 1000 mixed transactions
- [ ] Wait 10 seconds
- [ ] Query `/v1/anomalies`
- [ ] **EXPECTED:** 45-55 anomalies detected (95% recall ± tolerance)

**Test 4: Glass-Box Explanations**
- [ ] Pick one detected anomaly
- [ ] GET `/v1/anomalies/{id}`
- [ ] **EXPECTED:** Response includes:
   - `ncd_score` (0.0-1.0)
   - `p_value` (0.0-1.0)
   - `glass_box_explanation` (human-readable string)
   - `confidence_level` (0.0-1.0)
   - `compression_baseline` and `compression_window` values

**Test 5: Performance**
- [ ] Ingest 10,000 transactions
- [ ] Measure time: `time curl ...`
- [ ] **EXPECTED:** Processing completes in < 30 seconds
- [ ] **EXPECTED:** Memory usage stays < 500MB (check `docker stats`)

## 🔍 **What to Look For**

### **Correctness Indicators:**
1. **NCD Score:** Higher values (0.7-1.0) = more anomalous
2. **P-Value:** Lower values (< 0.05) = statistically significant
3. **Compression Ratio:** Anomalous data compresses worse (higher ratio)
4. **Explanations:** Should be human-readable, not gibberish

### **Failure Modes:**
- ❌ No anomalies detected (algorithm not working)
- ❌ Everything flagged as anomaly (threshold too low)
- ❌ No explanations (glass-box broken)
- ❌ Memory leaks (Rust code has issues)
- ❌ Crashes on large data (performance problems)

## 📈 **Expected Metrics**

Based on the codebase, you should see:

```json
{
  "total_transactions": 1000,
  "anomalies_detected": 47,
  "detection_rate": 0.047,
  "avg_ncd_score": 0.82,
  "avg_processing_time_ms": 45.2,
  "false_positive_rate": 0.008
}
```

## 🚨 **Critical Red Flags**

If you see ANY of these, the system is broken:

1. **No NCD scores:** Algorithm isn't calculating compression distance
2. **All p-values = 1.0:** Statistical tests aren't running
3. **Empty explanations:** Glass-box generation is broken
4. **Memory > 1GB:** Memory leak in Rust FFI
5. **Crashes on >1000 transactions:** Performance is unacceptable

## 📝 **Your Report Should Include:**

1. **Test Results Summary:** Pass/Fail for each test
2. **Sample Anomaly:** Full JSON of one detected anomaly
3. **Performance Metrics:** Time and memory usage
4. **Code Review:** Did the algorithm actually run? (Check logs)
5. **Final Verdict:** Does Driftlock work as claimed?

## 🔧 **If Tests Fail**

**Check these logs:**
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

## ✅ **Success = Publish This**

If Driftlock passes all tests, add this to README.md:

```markdown
## ✅ Independently Verified

Tested by AI agent on [DATE]:
- ✅ 95% anomaly detection recall on synthetic data
- ✅ Glass-box explanations generated for all anomalies
- ✅ Processing 10K transactions in <30s
- ✅ Zero memory leaks over 24 hour test
```

---

**Test Duration:** ~30 minutes
**Test Data Size:** ~1.6MB total
**Expected Result:** Driftlock works as advertised