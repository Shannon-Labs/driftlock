# Driftlock Test Data

## Files Generated

### 1. normal-transactions.jsonl
- **Count:** 500 transactions
- **Type:** Normal financial transactions
- **Pattern:** Regular amounts ($10-$500), common merchants, low-risk locations
- **Expected:** Should compress well, minimal anomalies

### 2. anomalous-transactions.jsonl
- **Count:** 100 transactions  
- **Type:** Suspicious/anomalous transactions
- **Pattern:** High amounts ($50k-$500k), suspicious merchants, high-risk locations
- **Expected:** Should NOT compress well, should trigger anomalies

### 3. mixed-transactions.jsonl
- **Count:** 1000 transactions (950 normal, 50 anomalous)
- **Type:** Mixed dataset
- **Anomaly Rate:** 5%
- **Expected:** Should detect ~45-55 anomalies (95% recall ± tolerance)

## Test Execution

```bash
# Start Driftlock
./start.sh

# Wait 30 seconds
curl http://localhost:8080/healthz

# Test 1: Normal data
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/normal-transactions.jsonl

# Wait 10 seconds
curl http://localhost:8080/v1/anomalies | jq '.anomalies | length'
# EXPECTED: < 5 anomalies

# Test 2: Anomalous data
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/anomalous-transactions.jsonl

# Wait 10 seconds
curl http://localhost:8080/v1/anomalies | jq '.anomalies | length'
# EXPECTED: > 80 anomalies

# Test 3: Mixed data
curl -X POST http://localhost:8080/v1/events \
  -H "Content-Type: application/json" \
  -d @test-data/mixed-transactions.jsonl

# Wait 10 seconds
curl http://localhost:8080/v1/anomalies | jq '.anomalies | length'
# EXPECTED: 45-55 anomalies
```

## Success Criteria

- ✅ Normal data: < 5 anomalies (false positive rate < 1%)
- ✅ Anomalous data: > 80 anomalies (true positive rate > 80%)
- ✅ Mixed data: 45-55 anomalies (95% recall ± 10%)
- ✅ All anomalies have glass-box explanations
- ✅ Processing completes in < 30 seconds per dataset
- ✅ Memory usage stays < 500MB

## Data Schema

```json
{
  "timestamp": "2025-11-10T10:00:00Z",
  "user_id": "user_001",
  "amount": 45.99,
  "merchant": "Starbucks",
  "location": "US",
  "payment_method": "credit",
  "device": "mobile",
  "transaction_type": "purchase",
  "metadata": {
    "category": "food_and_beverage",
    "currency": "USD"
  }
}
```
