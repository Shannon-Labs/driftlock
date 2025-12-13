# driftlock-ai

AI-powered anomaly explanations for Driftlock using Claude API.

## Features

- **Claude API Integration**: Direct integration with Anthropic's Claude API using reqwest
- **Cost Tracking**: Automatic cost calculation with 15% margin for all Claude models
- **Plan-Based Model Selection**: Default model selection based on subscription tier
- **Prompt Injection Protection**: Sanitization filters to prevent prompt injection attacks
- **Token Usage Tracking**: Full tracking of input/output tokens and costs

## Installation

Add to your `Cargo.toml`:

```toml
[dependencies]
driftlock-ai = { path = "../driftlock-ai" }
```

## Usage

### Basic Example

```rust
use driftlock_ai::{ClaudeClient, AiClient, sanitize_for_prompt};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Create client from environment variable ANTHROPIC_API_KEY
    let client = ClaudeClient::from_env()?;

    // Sanitize user data for safe inclusion in prompt
    let event_data = r#"{"error": "Connection timeout", "count": 50}"#;
    let sanitized = sanitize_for_prompt(event_data, 2048);

    // Build prompt
    let prompt = format!(
        "Explain this anomaly:\n\nEvent data:\n{}\n\nProvide a brief explanation.",
        sanitized
    );

    // Analyze anomaly
    let response = client
        .analyze_anomaly("claude-haiku-4-5-20251001", &prompt)
        .await?;

    println!("Explanation: {}", response.text);
    println!("Input tokens: {}", response.input_tokens);
    println!("Output tokens: {}", response.output_tokens);
    println!("Cost: ${:.6}", response.cost_usd);

    Ok(())
}
```

### Plan-Based Configuration

```rust
use driftlock_ai::PlanConfig;

// Get default configuration for a plan
let config = PlanConfig::for_plan("tensor");

println!("Default model: {}", config.default_model);
println!("Max calls per day: {}", config.max_calls_per_day);
println!("Max cost per month: ${}", config.max_cost_per_month);

// Check if a model is allowed
if config.is_model_allowed("claude-sonnet-4-5-20250929") {
    println!("Sonnet is allowed for Tensor plan");
}
```

## Models and Pricing

All pricing includes 15% margin:

| Model | Input (per 1M tokens) | Output (per 1M tokens) | Tier |
|-------|----------------------|------------------------|------|
| claude-haiku-4-5-20251001 | $1.15 | $5.75 | Haiku |
| claude-sonnet-4-5-20250929 | $3.45 | $17.25 | Sonnet |
| claude-opus-4-5-20251101 | $5.75 | $28.75 | Opus |

## Plan Defaults

| Plan | Default Model | Max Calls/Day | Max Cost/Month |
|------|--------------|---------------|----------------|
| Trial/Pilot | Haiku | 100 | $10 |
| Radar | Haiku | 200 | $20 |
| Tensor/Lock | Sonnet | 500 | $150 |
| Orbit | Opus | Unlimited | Unlimited |

## Environment Variables

- `ANTHROPIC_API_KEY` - Your Anthropic API key (required)
- `ANTHROPIC_MODEL` - Override default model (optional)

## Security

The crate includes built-in prompt injection protection via `sanitize_for_prompt()`:

- Filters known injection patterns (e.g., "ignore previous instructions")
- Removes special tokens (e.g., `</s>`, `<|im_end|>`)
- Adds clear data boundaries
- Limits data length to prevent context stuffing

## Testing

```bash
# Run all tests
cargo test -p driftlock-ai

# Run integration tests (requires ANTHROPIC_API_KEY)
cargo test -p driftlock-ai -- --ignored
```

## License

Apache-2.0
