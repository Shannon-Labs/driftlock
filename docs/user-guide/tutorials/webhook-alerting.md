# Tutorial: Webhook & Alerting

Send anomalies to a webhook (e.g., Slack/PagerDuty) after detection.

## Prerequisites
- API key
- A webhook endpoint URL (example: Slack incoming webhook)

## Steps

1. **Run detection** (example uses metrics spike):
   ```bash
   curl -X POST https://api.driftlock.net/v1/detect \
     -H "Content-Type: application/json" \
     -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
     -d '{"events": [{"body": {"latency_ms": 950}}]}' > detect.json
   ```

2. **Extract anomalies and post to webhook:**
   ```bash
   jq -c '.anomalies[] | {title: "Driftlock Anomaly", text: (.why // "Anomaly detected"), ncd: .metrics.ncd, p_value: .metrics.p_value}' detect.json \
     | while read -r payload; do
         curl -X POST "$WEBHOOK_URL" \
           -H "Content-Type: application/json" \
           -d "{\"text\": \"$payload\"}";
       done
   ```

3. **Optional:** filter by confidence before sending:
   ```bash
   jq -c '.anomalies[] | select(.metrics.confidence > 0.9)' detect.json | ...
   ```

## Tips
- Batch detections to reduce webhook noise.
- Include `stream_id`, `why`, and `request_id` in alerts for fast triage.
- Pair with the `strict` profile for paging channels; use `balanced` for chat notifications.
