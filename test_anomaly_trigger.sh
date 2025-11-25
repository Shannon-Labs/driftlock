#!/bin/bash
# Quick test to trigger screen recording

echo "🧪 Injecting test anomaly to trigger screen recorder..."

# Create a realistic anomaly alert entry
cat >> logs/anomaly_alerts.log << 'EOF'
═══════════════════════════════════════════════════════════
🚨 START LOOM NOW! 🚨
═══════════════════════════════════════════════════════════
Alert #999
Time: 2025-11-22 22:20:00
Trade Time: 2025-11-22 22:20:00
Anomaly ID: TEST-RECORDING-001
Pair: BTC/USD
Price: 86443.6
Score: 0.85

Full Payload:
{"ts":1763871600,"price":86443.6,"qty":1.0,"side":"sell","pair":"BTC/USD","id":"TEST-RECORDING-001","source":"synthetic","is_anomaly":true,"score":0.85,"reason":"Unusual pattern detected"}
═══════════════════════════════════════════════════════════
EOF

echo "✅ Test anomaly injected! Waiting for screen recorder to detect..."
