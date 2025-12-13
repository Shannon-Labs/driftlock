# Integration Guide: driftlock-ai

This guide shows how to integrate the `driftlock-ai` crate into the Driftlock API.

## Overview

The `driftlock-ai` crate provides AI-powered anomaly explanations using Claude API. It includes:

- Cost calculation with 15% margin
- Plan-based model selection
- Prompt injection filtering
- Token usage tracking

## Integration Steps

### 1. Add Dependency

In `crates/driftlock-api/Cargo.toml`:

```toml
[dependencies]
driftlock-ai = { path = "../driftlock-ai" }
```

### 2. Add Environment Variables

```bash
# Required
ANTHROPIC_API_KEY=your_api_key_here

# Optional (override default model)
ANTHROPIC_MODEL=claude-haiku-4-5-20251001
```

### 3. Usage in API Routes

```rust
use driftlock_ai::{ClaudeClient, AiClient, PlanConfig, sanitize_for_prompt};

// In your detection handler
pub async fn analyze_anomaly(
    anomaly_id: Uuid,
    tenant_plan: String,
    pool: Pool<Postgres>,
) -> Result<String, AppError> {
    // Get plan configuration
    let config = PlanConfig::for_plan(&tenant_plan);

    // Create client
    let client = ClaudeClient::from_env()
        .map_err(|e| AppError::Internal(format!("AI client error: {}", e)))?;

    // Fetch anomaly data from DB
    let anomaly = sqlx::query!(
        "SELECT event_data, score FROM anomalies WHERE id = $1",
        anomaly_id
    )
    .fetch_one(&pool)
    .await?;

    // Sanitize event data
    let sanitized = sanitize_for_prompt(&anomaly.event_data, 2048);

    // Build prompt
    let prompt = format!(
        "Analyze this anomaly (score: {}):\n\n{}\n\nProvide a brief explanation.",
        anomaly.score,
        sanitized
    );

    // Call AI with plan's default model
    let response = client
        .analyze_anomaly(&config.default_model, &prompt)
        .await
        .map_err(|e| AppError::Internal(format!("AI analysis failed: {}", e)))?;

    // Store usage metrics
    sqlx::query!(
        r#"
        INSERT INTO ai_usage (
            tenant_id, anomaly_id, model, input_tokens,
            output_tokens, cost_usd, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        "#,
        tenant_id,
        anomaly_id,
        config.default_model,
        response.input_tokens,
        response.output_tokens,
        response.cost_usd as f64
    )
    .execute(&pool)
    .await?;

    Ok(response.text)
}
```

### 4. Database Schema

Add these tables for AI usage tracking:

```sql
-- AI usage tracking
CREATE TABLE IF NOT EXISTS ai_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    anomaly_id UUID REFERENCES anomalies(id),
    model VARCHAR(255) NOT NULL,
    input_tokens BIGINT NOT NULL,
    output_tokens BIGINT NOT NULL,
    cost_usd DECIMAL(10, 6) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ai_usage_tenant ON ai_usage(tenant_id, created_at DESC);
CREATE INDEX idx_ai_usage_anomaly ON ai_usage(anomaly_id);
```

### 5. Add API Endpoints

Add these endpoints to `crates/driftlock-api/src/routes/`:

- `POST /v1/anomalies/:id/analyze` - Generate AI explanation
- `GET /v1/usage/ai` - Get AI usage statistics
- `GET /v1/config/ai` - Get/update AI configuration

### 6. Cost Control

Implement rate limiting and cost controls:

```rust
use driftlock_ai::PlanConfig;

pub async fn check_ai_quota(
    tenant_id: Uuid,
    plan: &str,
    pool: &Pool<Postgres>,
) -> Result<bool, AppError> {
    let config = PlanConfig::for_plan(plan);

    // Check daily limit
    let today_count: i64 = sqlx::query_scalar!(
        "SELECT COUNT(*) FROM ai_usage
         WHERE tenant_id = $1 AND DATE(created_at) = CURRENT_DATE",
        tenant_id
    )
    .fetch_one(pool)
    .await?;

    if config.max_calls_per_day > 0 && today_count >= config.max_calls_per_day as i64 {
        return Ok(false);
    }

    // Check monthly cost
    let month_cost: Option<f64> = sqlx::query_scalar!(
        "SELECT COALESCE(SUM(cost_usd), 0) FROM ai_usage
         WHERE tenant_id = $1 AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', NOW())",
        tenant_id
    )
    .fetch_one(pool)
    .await?;

    if config.max_cost_per_month > 0.0 && month_cost.unwrap_or(0.0) >= config.max_cost_per_month {
        return Ok(false);
    }

    Ok(true)
}
```

## Testing

### Unit Tests

```bash
# Run all tests
cargo test -p driftlock-ai

# Run with integration tests (requires ANTHROPIC_API_KEY)
ANTHROPIC_API_KEY=your_key cargo test -p driftlock-ai -- --ignored
```

### Example

```bash
# Run the example
ANTHROPIC_API_KEY=your_key cargo run -p driftlock-ai --example analyze_anomaly
```

## Plan Configuration Reference

| Plan | Model | Max Calls/Day | Max Cost/Month | Use Case |
|------|-------|---------------|----------------|----------|
| Trial | Haiku | 100 | $10 | Testing, demos |
| Radar | Haiku | 200 | $20 | Small teams |
| Tensor | Sonnet | 500 | $150 | Growing teams |
| Orbit | Opus | Unlimited | Unlimited | Enterprise |

## Cost Optimization Tips

1. **Use Haiku for simple explanations** - 5x cheaper than Opus
2. **Cache explanations** - Store AI responses to avoid repeated calls
3. **Batch processing** - Analyze multiple anomalies in one prompt
4. **Threshold filtering** - Only analyze anomalies above config.analysis_threshold
5. **Smart routing** - Use cheaper models for low-severity anomalies

## Error Handling

Common errors and how to handle them:

```rust
match client.analyze_anomaly(model, prompt).await {
    Ok(response) => Ok(response),
    Err(AiError::MissingApiKey(_)) => {
        // API key not configured
        Err(AppError::Configuration("AI service not configured".to_string()))
    }
    Err(AiError::RateLimit(msg)) => {
        // Rate limited by Anthropic
        Err(AppError::RateLimit(msg))
    }
    Err(AiError::CostLimit(msg)) => {
        // Tenant exceeded their cost limit
        Err(AppError::QuotaExceeded(msg))
    }
    Err(e) => {
        // Other errors
        Err(AppError::Internal(format!("AI error: {}", e)))
    }
}
```

## Monitoring

Track these metrics:

- AI calls per tenant per day
- Monthly AI costs per tenant
- Average tokens per request
- Error rates
- Response times
- Cost per anomaly type

## Security Considerations

1. **Prompt Injection** - Always use `sanitize_for_prompt()` on user data
2. **API Key** - Store in environment variables, never in code
3. **Rate Limiting** - Implement per-tenant limits
4. **Cost Controls** - Monitor and alert on unexpected cost increases
5. **Data Privacy** - Ensure user data handling complies with policies

## Next Steps

1. Add `driftlock-ai` dependency to API
2. Add database migrations for AI usage tracking
3. Implement AI analysis endpoints
4. Add cost control middleware
5. Update dashboard to show AI explanations
6. Add usage/cost monitoring

## Support

For issues or questions, see:
- `crates/driftlock-ai/README.md` - Usage documentation
- `crates/driftlock-ai/examples/` - Example code
- Anthropic API docs: https://docs.anthropic.com
