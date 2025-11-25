#!/bin/bash
# Inject a full-formatted alert like the real ones

cat >> logs/anomaly_alerts.log << 'DELIM'

═══════════════════════════════════════════════════════════
🚨 AUTO RECORDING TEST 🚨
═══════════════════════════════════════════════════════════
Alert #888
Time: 2025-11-22 22:23:00
Trade Time: 2025-11-22 22:23:00
Anomaly ID: FULL-FORMAT-004
Pair: BTC/USD
Price: 86443.6
Score: 0.90

Full Payload:
{"ts":1763871780,"price":86443.6,"qty":0.5,"side":"buy","pair":"BTC/USD","id":"FULL-FORMAT-004","source":"driftlock","is_anomaly":true,"score":0.90,"reason":"Compression entropy spike detected"}
═══════════════════════════════════════════════════════════
DELIM

echo "✅ Full-format anomaly alert injected"
