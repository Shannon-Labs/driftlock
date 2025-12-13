/// Claude model pricing information with 15% margin
#[derive(Debug, Clone)]
pub struct ModelPricing {
    /// Cost per million input tokens (USD)
    pub input_cost_per_million: f64,
    /// Cost per million output tokens (USD)
    pub output_cost_per_million: f64,
}

/// Get pricing for a Claude model (with 15% margin)
pub fn get_model_pricing(model: &str) -> Option<ModelPricing> {
    // Base pricing from Anthropic (as of 2025)
    let (base_input, base_output) = match model {
        // Haiku models - cheapest
        "claude-haiku-4-5-20251001" | "claude-3-5-haiku-20241022" => (1.0, 5.0),

        // Sonnet models - balanced
        "claude-sonnet-4-5-20250929" | "claude-3-5-sonnet-20241022" => (3.0, 15.0),

        // Opus models - best quality
        "claude-opus-4-5-20251101" | "claude-3-opus-20240229" => (5.0, 25.0),

        _ => return None,
    };

    // Apply 15% margin
    const MARGIN: f64 = 1.15;

    Some(ModelPricing {
        input_cost_per_million: base_input * MARGIN,
        output_cost_per_million: base_output * MARGIN,
    })
}

/// Calculate cost for token usage
pub fn calculate_cost(model: &str, input_tokens: i64, output_tokens: i64) -> Option<f64> {
    let pricing = get_model_pricing(model)?;

    let input_cost = (input_tokens as f64 / 1_000_000.0) * pricing.input_cost_per_million;
    let output_cost = (output_tokens as f64 / 1_000_000.0) * pricing.output_cost_per_million;

    Some(input_cost + output_cost)
}

/// Get model tier (haiku, sonnet, opus)
pub fn get_model_tier(model: &str) -> Option<&'static str> {
    if model.contains("haiku") {
        Some("haiku")
    } else if model.contains("sonnet") {
        Some("sonnet")
    } else if model.contains("opus") {
        Some("opus")
    } else {
        None
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_haiku_pricing() {
        let pricing = get_model_pricing("claude-haiku-4-5-20251001").unwrap();
        assert_eq!(pricing.input_cost_per_million, 1.15); // $1 * 1.15
        assert_eq!(pricing.output_cost_per_million, 5.75); // $5 * 1.15
    }

    #[test]
    fn test_sonnet_pricing() {
        let pricing = get_model_pricing("claude-sonnet-4-5-20250929").unwrap();
        assert!((pricing.input_cost_per_million - 3.45).abs() < 0.001); // $3 * 1.15
        assert!((pricing.output_cost_per_million - 17.25).abs() < 0.001); // $15 * 1.15
    }

    #[test]
    fn test_opus_pricing() {
        let pricing = get_model_pricing("claude-opus-4-5-20251101").unwrap();
        assert!((pricing.input_cost_per_million - 5.75).abs() < 0.001); // $5 * 1.15
        assert!((pricing.output_cost_per_million - 28.75).abs() < 0.001); // $25 * 1.15
    }

    #[test]
    fn test_calculate_cost() {
        // 1000 input tokens, 500 output tokens on Haiku
        let cost = calculate_cost("claude-haiku-4-5-20251001", 1000, 500).unwrap();
        // (1000/1M * 1.15) + (500/1M * 5.75) = 0.00115 + 0.002875 = 0.004025
        assert!((cost - 0.004025).abs() < 0.000001);
    }

    #[test]
    fn test_unknown_model() {
        assert!(get_model_pricing("unknown-model").is_none());
        assert!(calculate_cost("unknown-model", 1000, 500).is_none());
    }

    #[test]
    fn test_model_tier() {
        assert_eq!(get_model_tier("claude-haiku-4-5-20251001"), Some("haiku"));
        assert_eq!(get_model_tier("claude-sonnet-4-5-20250929"), Some("sonnet"));
        assert_eq!(get_model_tier("claude-opus-4-5-20251101"), Some("opus"));
        assert_eq!(get_model_tier("unknown-model"), None);
    }
}
