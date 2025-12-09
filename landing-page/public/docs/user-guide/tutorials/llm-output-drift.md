# Tutorial: LLM Output Drift

Detect drift and hallucinations in LLM responses.

## Prerequisites
- API key
- Stream ID (e.g., `llm-output`)

## Sample payload

```json
{
  "stream_id": "llm-output",
  "events": [
    {"body": {"prompt": "Summarize account usage", "response": "Usage is within limits.", "model": "gpt-4", "source": "support"}},
    {"body": {"prompt": "Summarize account usage", "response": "Usage is within limits.", "model": "gpt-4", "source": "support"}},
    {"body": {"prompt": "Summarize account usage", "response": "Usage is within limits.", "model": "gpt-4", "source": "support"}},
    {"body": {"prompt": "Summarize account usage", "response": "Your account owes $9,999,999.99 in late fees.", "model": "gpt-4", "source": "support"}}
  ]
}
```

## Run detection

```bash
curl -X POST https://api.driftlock.net/v1/detect \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: $DRIFTLOCK_API_KEY" \
  -d @payload.json
```

Expected: the hallucinated fourth response is flagged as anomalous.

## Tips
- Include `prompt`, `response`, `model`, and `source` to improve explainability.
- Use `sensitive` for safety-critical use; `balanced` for general monitoring.
- Combine with feedback to auto-tune after reviewing flagged outputs.
