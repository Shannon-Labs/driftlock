# Driftlock API Demo Endpoint Test Results

**Test Date:** 2025-12-07
**API URL:** https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect
**Compression Algorithm:** zstd

---

## Executive Summary

Tested 5 different datasets against the deployed Driftlock API demo endpoint. All tests completed successfully, but anomaly detection rates show concerning patterns that differ from expected behavior.

**Overall Status:** PARTIALLY PASSING (3/5 tests failed expected criteria)

---

## Test Results

### Test 1: Financial Demo Dataset
**Dataset:** /Volumes/VIXinSSD/driftlock/test-data/financial-demo.json
**Events Sent:** 100
**Anomalies Detected:** 61
**Anomaly Rate:** 61%
**Processing Time:** 14.03 seconds
**Expected:** Mixed dataset (baseline test)
**Status:** BASELINE (no specific criteria)

**Analysis:**
- High anomaly rate indicates the demo dataset contains significant pattern variations
- All anomalies detected with moderate confidence (44-48%)
- NCD scores ranging from 0.58-0.63 (moderate dissimilarity)
- Most anomalies show 43-47% compression efficiency drop
- None of the anomalies were statistically significant (high p-values: 0.05-0.98)

---

### Test 2: Normal Transactions
**Dataset:** /Volumes/VIXinSSD/driftlock/test-data/normal-transactions.jsonl
**Events Sent:** 50
**Anomalies Detected:** 11
**Anomaly Rate:** 22%
**Processing Time:** 1.82 seconds
**Expected:** < 5% anomaly rate (< 2.5 anomalies out of 50)
**Status:** FAIL

**Analysis:**
- EXPECTED: < 2.5 anomalies (< 5% false positive rate)
- ACTUAL: 11 anomalies (22% false positive rate)
- FALSE POSITIVE RATE: 22% (4.4x higher than expected)
- This indicates the system is flagging normal transactions as anomalous at an unacceptably high rate

---

### Test 3: Anomalous Transactions
**Dataset:** /Volumes/VIXinSSD/driftlock/test-data/anomalous-transactions.jsonl
**Events Sent:** 50
**Anomalies Detected:** 11
**Anomaly Rate:** 22%
**Processing Time:** 8.66 seconds
**Expected:** > 50% anomaly rate (> 25 anomalies out of 50)
**Status:** FAIL

**Analysis:**
- EXPECTED: > 25 anomalies (> 50% true positive rate, ideally > 80%)
- ACTUAL: 11 anomalies (22% true positive rate)
- TRUE POSITIVE RATE: 22% (significantly below expected 80%)
- This indicates the system is MISSING 78% of known anomalous transactions
- CRITICAL ISSUE: The system should detect at least 40 anomalies (80% TP rate)

---

### Test 4: Mixed Transactions
**Dataset:** /Volumes/VIXinSSD/driftlock/test-data/mixed-transactions.jsonl
**Events Sent:** 100 (95 normal + 5 anomalous based on 5% anomaly rate)
**Anomalies Detected:** 61
**Anomaly Rate:** 61%
**Processing Time:** 14.09 seconds
**Expected:** ~5% anomaly rate (4-6 anomalies out of 100)
**Status:** FAIL

**Analysis:**
- EXPECTED: 4-6 anomalies (5% ± tolerance)
- ACTUAL: 61 anomalies (61% rate)
- DEVIATION: 55% over expected (12.2x higher than expected)
- This suggests the system is severely over-flagging on mixed datasets

---

### Test 5: Fraud Dataset
**Dataset:** /Volumes/VIXinSSD/driftlock/test-data/fraud/driftlock_ready.json
**Events Sent:** 50 (all flagged as "review_flagged" in source data)
**Anomalies Detected:** 11
**Anomaly Rate:** 22%
**Processing Time:** 4.81 seconds
**Expected:** High detection rate (these are all fraud-flagged transactions)
**Status:** FAIL (if all events should be detected)

**Analysis:**
- All 50 events in the fraud dataset have status "review_flagged"
- Only 11/50 (22%) were detected as anomalies
- If all events are fraudulent, TRUE POSITIVE RATE: 22%
- MISSED FRAUD CASES: 78%
- This is concerning for fraud detection use cases

---

## Success Criteria Analysis

### From test-data/README.md:

| Criterion | Expected | Actual | Pass/Fail |
|-----------|----------|--------|-----------|
| Normal data: < 5 anomalies | < 2.5 out of 50 (< 5%) | 11 out of 50 (22%) | FAIL |
| Anomalous data: > 80% TP rate | > 40 out of 50 | 11 out of 50 (22%) | FAIL |
| Mixed data: ~5% anomaly rate | 4-6 out of 100 | 61 out of 100 (61%) | FAIL |

---

## Key Findings

### Concerning Patterns:

1. **Consistent 22% Detection Rate Across All Datasets:**
   - Normal transactions: 22% (11/50)
   - Anomalous transactions: 22% (11/50)
   - Fraud dataset: 22% (11/50)
   - This suggests the detection algorithm may have a fixed threshold or window size issue

2. **Low Statistical Significance:**
   - Most detected anomalies have high p-values (0.05-0.98)
   - This indicates the anomalies are not statistically significant
   - Confidence levels are moderate (44-48%) but not high

3. **NCD Score Patterns:**
   - NCD scores consistently in 0.58-0.63 range (moderate dissimilarity)
   - Compression ratio drops of 43-47% are consistent
   - Entropy changes are minimal (-0.4% to +0.0%)

4. **Processing Performance:**
   - 50 events: 1.8-8.7 seconds (good)
   - 100 events: 14.0-14.1 seconds (good)
   - Processing time is acceptable for demo use

### Potential Issues:

1. **Sliding Window Size:** The algorithm may be using a fixed window that triggers at event 39-40 (seen in Test 1)
2. **Threshold Calibration:** Detection thresholds may need adjustment for different data patterns
3. **Baseline Training:** The baseline compression ratio may not be adapting to normal vs anomalous patterns
4. **False Positive/Negative Balance:** Current settings favor precision over recall, missing many true anomalies

---

## Recommendations

### Immediate Actions:

1. **Investigate the 22% pattern** - Why are all datasets yielding the same detection rate?
2. **Review window size configuration** - Anomalies start at index 39 consistently
3. **Lower detection threshold** - To improve recall on anomalous and fraud datasets
4. **Add statistical significance filtering** - Consider only flagging anomalies with p < 0.05

### Algorithmic Improvements:

1. **Adaptive Baseline:** Train baseline on first N events to establish normal patterns
2. **Dynamic Thresholds:** Adjust NCD threshold based on data characteristics
3. **Confidence Weighting:** Require higher confidence for normal-looking data
4. **Multi-Window Analysis:** Use multiple window sizes to catch different anomaly types

### Testing Improvements:

1. **Add ground truth labels** - Label which events in mixed dataset are actually anomalous
2. **Calculate precision/recall metrics** - Not just anomaly count
3. **Test edge cases** - Single anomaly, bursts of anomalies, gradual drift
4. **Benchmark against other algorithms** - Compare to isolation forest, LOF, etc.

---

## Conclusion

The Driftlock API demo endpoint is **functionally operational** but the anomaly detection accuracy **does not meet the specified success criteria** from the test data README.

**Critical Issues:**
- False positive rate: 22% (expected < 5%)
- True positive rate: 22% (expected > 80%)
- Consistent 22% detection rate across all dataset types suggests systematic issue

**Next Steps:**
1. Debug why detection rate is fixed at ~22% regardless of data type
2. Investigate sliding window and threshold parameters
3. Re-calibrate algorithm for better true positive/false positive balance
4. Validate results with labeled ground truth data

---

## Test Commands Used

```bash
# Test 1: Financial demo
jq '.[0:100]' /Volumes/VIXinSSD/driftlock/test-data/financial-demo.json | jq -c '{events: .}' | \
  curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
  -H "Content-Type: application/json" -d @-

# Test 2: Normal transactions
head -50 /Volumes/VIXinSSD/driftlock/test-data/normal-transactions.jsonl | jq -s '{events: .}' | \
  curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
  -H "Content-Type: application/json" -d @-

# Test 3: Anomalous transactions
head -50 /Volumes/VIXinSSD/driftlock/test-data/anomalous-transactions.jsonl | jq -s '{events: .}' | \
  curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
  -H "Content-Type: application/json" -d @-

# Test 4: Mixed transactions
head -100 /Volumes/VIXinSSD/driftlock/test-data/mixed-transactions.jsonl | jq -s '{events: .}' | \
  curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
  -H "Content-Type: application/json" -d @-

# Test 5: Fraud dataset
jq -c '{events: .[0:50]}' /Volumes/VIXinSSD/driftlock/test-data/fraud/driftlock_ready.json | \
  curl -s -X POST https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1/demo/detect \
  -H "Content-Type: application/json" -d @-
```
