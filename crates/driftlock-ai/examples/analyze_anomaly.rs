use driftlock_ai::{sanitize_for_prompt, AiClient, ClaudeClient, PlanConfig};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    tracing_subscriber::fmt::init();

    // Show plan configurations
    println!("=== Plan Configurations ===\n");
    for plan in &["trial", "radar", "tensor", "orbit"] {
        let config = PlanConfig::for_plan(plan);
        println!("Plan: {}", config.plan);
        println!("  Default model: {}", config.default_model);
        println!("  Max calls/day: {}", config.max_calls_per_day);
        println!(
            "  Max cost/month: ${}",
            if config.max_cost_per_month == 0.0 {
                "unlimited".to_string()
            } else {
                format!("{:.2}", config.max_cost_per_month)
            }
        );
        println!();
    }

    // Create client from environment
    println!("=== Creating Claude Client ===\n");
    let client = match ClaudeClient::from_env() {
        Ok(client) => {
            println!("Client created successfully!");
            println!("Provider: {}", client.provider());
            println!("Default model: {}", client.default_model());
            client
        }
        Err(e) => {
            eprintln!("Error: {}", e);
            eprintln!("\nTo use this example, set ANTHROPIC_API_KEY environment variable:");
            eprintln!("  export ANTHROPIC_API_KEY='your-key-here'");
            return Ok(());
        }
    };

    // Example anomaly data
    println!("\n=== Analyzing Anomaly ===\n");
    let event_data = r#"{
        "timestamp": "2025-12-11T10:30:00Z",
        "service": "api-gateway",
        "error": "Connection timeout",
        "count": 50,
        "duration_ms": 30000,
        "endpoint": "/api/v1/users"
    }"#;

    // Sanitize the data
    let sanitized = sanitize_for_prompt(event_data, 2048);
    println!("Sanitized event data:\n{}\n", sanitized);

    // Build the prompt
    let prompt = format!(
        r#"You are an expert in analyzing system anomalies. Analyze this anomaly and provide:
1. A brief explanation of what happened
2. The likely cause
3. Potential impact
4. Recommended actions

Event data:
{}

Keep your response concise (2-3 sentences per point)."#,
        sanitized
    );

    // Use Haiku model (cheapest)
    let model = "claude-haiku-4-5-20251001";
    println!("Using model: {}\n", model);

    // Analyze the anomaly
    println!("Calling Claude API...\n");
    let response = client.analyze_anomaly(model, &prompt).await?;

    // Display results
    println!("=== Analysis Results ===\n");
    println!("{}\n", response.text);
    println!("=== Usage Statistics ===\n");
    println!("Input tokens: {}", response.input_tokens);
    println!("Output tokens: {}", response.output_tokens);
    println!(
        "Total tokens: {}",
        response.input_tokens + response.output_tokens
    );
    println!("Cost: ${:.6}", response.cost_usd);

    // Estimate costs for different models
    println!("\n=== Cost Comparison ===\n");
    let models = vec![
        ("claude-haiku-4-5-20251001", "Haiku (cheapest)"),
        ("claude-sonnet-4-5-20250929", "Sonnet (balanced)"),
        ("claude-opus-4-5-20251101", "Opus (best quality)"),
    ];

    for (model, desc) in models {
        if let Some(cost) =
            driftlock_ai::calculate_cost(model, response.input_tokens, response.output_tokens)
        {
            println!("{}: ${:.6}", desc, cost);
        }
    }

    Ok(())
}
