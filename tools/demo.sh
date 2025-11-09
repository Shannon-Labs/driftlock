#!/bin/bash
# AUTHENTIC DEMO SETUP - 3 commands, no theater

# Step 1: Generate realistic financial transaction data with anomalies
cat > demo-data.json << 'JSONEOF'
{
  "timestamp": "2025-11-10T10:00:00Z",
  "stream_type": "logs", 
  "ncd_score": 0.88,
  "p_value": 0.001,
  "glass_box_explanation": "Transaction amount exceeds 3-sigma threshold",
  "compression_baseline": 0.45,
  "compression_window": 0.92,
  "compression_combined": 0.88,
  "confidence_level": 0.999,
  "baseline_data": {"merchant": "Amazon", "amount": 150.50, "location": "DE"},
  "window_data": {"merchant": "Unknown", "amount": 12500.00, "location": "RU"},
  "metadata": {"transaction_id": "tx_002", "currency": "EUR"},
  "tags": ["high_amount", "unusual_location", "dora_compliance"]
}
JSONEOF

# Step 2: Send it to Driftlock via curl
echo "Sending demo data to Driftlock..."
curl -X POST http://localhost:8080/v1/anomalies \
  -H "Authorization: Bearer $DEMO_API_KEY" \
  -H "Content-Type: application/json" \
  -d @demo-data.json | jq .

# Step 3: Show the explanation
echo ""
echo "Now open http://localhost:3000 and click the anomaly to see the explanation."