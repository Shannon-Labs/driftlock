#!/bin/bash
# Generate synthetic transaction data for Driftlock testing

set -e

echo "🧪 Generating Driftlock test data..."

# Configuration
NORMAL_COUNT=500
ANOMALOUS_COUNT=100
MIXED_COUNT=1000
ANOMALY_RATE=0.05  # 5% anomalies in mixed data

# Generate normal transactions
echo "📊 Generating $NORMAL_COUNT normal transactions..."
cat > /Volumes/VIXinSSD/driftlock/test-data/normal-transactions.jsonl << 'EOF'
{"timestamp":"2025-11-10T10:00:00Z","user_id":"user_001","amount":45.99,"merchant":"Starbucks","location":"US","payment_method":"credit","device":"mobile","transaction_type":"purchase","metadata":{"category":"food_and_beverage","currency":"USD"}}
{"timestamp":"2025-11-10T10:00:01Z","user_id":"user_002","amount":125.50,"merchant":"Amazon","location":"US","payment_method":"credit","device":"web","transaction_type":"purchase","metadata":{"category":"electronics","currency":"USD"}}
{"timestamp":"2025-11-10T10:00:02Z","user_id":"user_003","amount":23.75,"merchant":"Uber","location":"US","payment_method":"debit","device":"mobile","transaction_type":"ride","metadata":{"category":"transportation","currency":"USD"}}
{"timestamp":"2025-11-10T10:00:03Z","user_id":"user_004","amount":89.99,"merchant":"Netflix","location":"US","payment_method":"credit","device":"smart_tv","transaction_type":"subscription","metadata":{"category":"entertainment","currency":"USD"}}
{"timestamp":"2025-11-10T10:00:04Z","user_id":"user_005","amount":12.50,"merchant":"Spotify","location":"US","payment_method":"digital","device":"mobile","transaction_type":"subscription","metadata":{"category":"entertainment","currency":"USD"}}
EOF

# Generate 495 more normal transactions
for i in {006..500}; do
  amount=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 500 + 10" | bc)")
  merchant=$(shuf -n 1 -e "Starbucks" "Amazon" "Uber" "Netflix" "Spotify" "Walmart" "Target" "Apple" "Google" "Microsoft")
  location=$(shuf -n 1 -e "US" "CA" "UK" "DE" "FR")
  payment_method=$(shuf -n 1 -e "credit" "debit" "digital")
  device=$(shuf -n 1 -e "mobile" "web" "pos" "smart_tv")
  category=$(shuf -n 1 -e "food_and_beverage" "electronics" "transportation" "entertainment" "shopping" "software")
  
  cat >> /Volumes/VIXinSSD/driftlock/test-data/normal-transactions.jsonl << EOF
{"timestamp":"2025-11-10T10:00:${i}Z","user_id":"user_${i}","amount":${amount},"merchant":"${merchant}","location":"${location}","payment_method":"${payment_method}","device":"${device}","transaction_type":"purchase","metadata":{"category":"${category}","currency":"USD"}}
EOF
done

# Generate anomalous transactions
echo "🚨 Generating $ANOMALOUS_COUNT anomalous transactions..."
cat > /Volumes/VIXinSSD/driftlock/test-data/anomalous-transactions.jsonl << 'EOF'
{"timestamp":"2025-11-10T10:01:00Z","user_id":"user_suspicious_001","amount":125000.00,"merchant":"Offshore Casino","location":"RU","payment_method":"crypto","device":"tor_browser","transaction_type":"transfer","metadata":{"category":"gambling","currency":"BTC","risk_score":0.95}}
{"timestamp":"2025-11-10T10:01:01Z","user_id":"user_suspicious_002","amount":99999.99,"merchant":"Darknet Marketplace","location":"??","payment_method":"anonymous","device":"vpn","transaction_type":"purchase","metadata":{"category":"illegal_goods","currency":"XMR","risk_score":0.99}}
{"timestamp":"2025-11-10T10:01:02Z","user_id":"user_suspicious_003","amount":75000.00,"merchant":"Shell Company A","location":"KY","payment_method":"wire","device":"unknown","transaction_type":"transfer","metadata":{"category":"money_laundering","currency":"USD","risk_score":0.92}}
{"timestamp":"2025-11-10T10:01:03Z","user_id":"user_suspicious_004","amount":50000.00,"merchant":"Unregistered Exchange","location":"CN","payment_method":"crypto","device":"mobile","transaction_type":"exchange","metadata":{"category":"unregulated_crypto","currency":"ETH","risk_score":0.88}}
{"timestamp":"2025-11-10T10:01:04Z","user_id":"user_suspicious_005","amount":250000.00,"merchant":"High Risk Jurisdiction","location":"IR","payment_method":"cash","device":"atm","transaction_type":"withdrawal","metadata":{"category":"sanctions_violation","currency":"USD","risk_score":0.97}}
EOF

# Generate 95 more anomalous transactions
for i in {006..100}; do
  # High amount (50k-500k)
  amount=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 450000 + 50000" | bc)")
  
  # Suspicious merchant patterns
  merchant_patterns=("Offshore Casino" "Darknet Marketplace" "Shell Company" "Unregistered Exchange" "High Risk Jurisdiction" "Anonymous Transfer" "Crypto Mixer" "Sanctioned Entity")
  merchant=${merchant_patterns[$RANDOM % 8]}
  
  # High-risk locations
  locations=("RU" "CN" "IR" "KP" "SY" "??" "KY" "VG")
  location=${locations[$RANDOM % 8]}
  
  # Suspicious payment methods
  payment_methods=("crypto" "anonymous" "wire" "cash" "prepaid_card")
  payment_method=${payment_methods[$RANDOM % 5]}
  
  # Suspicious categories
  categories=("gambling" "illegal_goods" "money_laundering" "unregulated_crypto" "sanctions_violation" "fraud")
  category=${categories[$RANDOM % 6]}
  
  # High risk score
  risk_score=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 30 + 70" | bc | awk '{print $1/100}')")
  
  cat >> /Volumes/VIXinSSD/driftlock/test-data/anomalous-transactions.jsonl << EOF
{"timestamp":"2025-11-10T10:01:${i}Z","user_id":"user_suspicious_${i}","amount":${amount},"merchant":"${merchant}","location":"${location}","payment_method":"${payment_method}","device":"unknown","transaction_type":"transfer","metadata":{"category":"${category}","currency":"USD","risk_score":${risk_score}}}
EOF
done

# Generate mixed transactions
echo "🎲 Generating $MIXED_COUNT mixed transactions (5% anomalies)..."

# Start with 950 normal transactions
for i in {001..950}; do
  amount=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 500 + 10" | bc)")
  merchant=$(shuf -n 1 -e "Starbucks" "Amazon" "Uber" "Netflix" "Spotify" "Walmart" "Target" "Apple" "Google" "Microsoft")
  location=$(shuf -n 1 -e "US" "CA" "UK" "DE" "FR")
  payment_method=$(shuf -n 1 -e "credit" "debit" "digital")
  device=$(shuf -n 1 -e "mobile" "web" "pos" "smart_tv")
  category=$(shuf -n 1 -e "food_and_beverage" "electronics" "transportation" "entertainment" "shopping" "software")
  
  cat >> /Volumes/VIXinSSD/driftlock/test-data/mixed-transactions.jsonl << EOF
{"timestamp":"2025-11-10T10:02:${i}Z","user_id":"user_${i}","amount":${amount},"merchant":"${merchant}","location":"${location}","payment_method":"${payment_method}","device":"${device}","transaction_type":"purchase","metadata":{"category":"${category}","currency":"USD"}}
EOF
done

# Add 50 anomalous transactions (5% of 1000)
for i in {951..1000}; do
  amount=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 450000 + 50000" | bc)")
  merchant_patterns=("Offshore Casino" "Darknet Marketplace" "Shell Company" "Unregistered Exchange" "High Risk Jurisdiction")
  merchant=${merchant_patterns[$RANDOM % 5]}
  locations=("RU" "CN" "IR" "??" "KY")
  location=${locations[$RANDOM % 5]}
  payment_methods=("crypto" "anonymous" "wire" "cash")
  payment_method=${payment_methods[$RANDOM % 4]}
  categories=("gambling" "money_laundering" "unregulated_crypto" "fraud")
  category=${categories[$RANDOM % 4]}
  risk_score=$(printf "%.2f" "$(echo "scale=2; $RANDOM % 30 + 70" | bc | awk '{print $1/100}')")
  
  cat >> /Volumes/VIXinSSD/driftlock/test-data/mixed-transactions.jsonl << EOF
{"timestamp":"2025-11-10T10:02:${i}Z","user_id":"user_anomalous_${i}","amount":${amount},"merchant":"${merchant}","location":"${location}","payment_method":"${payment_method}","device":"unknown","transaction_type":"transfer","metadata":{"category":"${category}","currency":"USD","risk_score":${risk_score}}}
EOF
done

# Create test summary
cat > /Volumes/VIXinSSD/driftlock/test-data/README.md << 'EOF'
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
EOF

echo ""
echo "✅ Test data generation complete!"
echo ""
echo "Generated files:"
echo "  📄 test-data/normal-transactions.jsonl (500 transactions)"
echo "  🚨 test-data/anomalous-transactions.jsonl (100 transactions)"
echo "  🎲 test-data/mixed-transactions.jsonl (1000 transactions)"
echo "  📖 test-data/README.md (test instructions)"
echo ""
echo "Total size: $(du -sh /Volumes/VIXinSSD/driftlock/test-data | cut -f1)"
echo ""
echo "To run tests:"
echo "  1. ./start.sh"
echo "  2. Wait 30 seconds"
echo "  3. Follow instructions in test-data/README.md"
EOF

chmod +x /Volumes/VIXinSSD/driftlock/test-data/generate-test-data.sh