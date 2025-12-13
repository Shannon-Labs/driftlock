use serde::{Deserialize, Serialize};

/// Plan-based configuration for AI usage
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlanConfig {
    /// Plan name
    pub plan: String,
    /// Default model to use
    pub default_model: String,
    /// Allowed models for this plan
    pub allowed_models: Vec<String>,
    /// Maximum calls per day (0 = unlimited)
    pub max_calls_per_day: u32,
    /// Maximum calls per hour (0 = unlimited)
    pub max_calls_per_hour: u32,
    /// Maximum cost per month in USD (0.0 = unlimited)
    pub max_cost_per_month: f64,
    /// Analysis threshold (0.0-1.0), only analyze anomalies above this score
    pub analysis_threshold: f64,
    /// Batch size for processing multiple anomalies
    pub batch_size: u32,
    /// Optimization target: "speed", "cost", or "accuracy"
    pub optimize_for: String,
    /// Notification threshold (0.0-1.0), alert when using this % of limit
    pub notify_threshold: f64,
}

impl PlanConfig {
    /// Get default configuration for a plan
    pub fn for_plan(plan: &str) -> Self {
        match plan {
            "trial" | "pilot" => Self {
                plan: plan.to_string(),
                default_model: "claude-haiku-4-5-20251001".to_string(),
                allowed_models: vec!["claude-haiku-4-5-20251001".to_string()],
                max_calls_per_day: 100,
                max_calls_per_hour: 20,
                max_cost_per_month: 10.0,
                analysis_threshold: 0.7,
                batch_size: 50,
                optimize_for: "cost".to_string(),
                notify_threshold: 0.9,
            },
            "radar" => Self {
                plan: plan.to_string(),
                default_model: "claude-haiku-4-5-20251001".to_string(),
                allowed_models: vec!["claude-haiku-4-5-20251001".to_string()],
                max_calls_per_day: 200,
                max_calls_per_hour: 40,
                max_cost_per_month: 20.0,
                analysis_threshold: 0.6,
                batch_size: 50,
                optimize_for: "cost".to_string(),
                notify_threshold: 0.8,
            },
            "tensor" | "lock" => Self {
                plan: plan.to_string(),
                default_model: "claude-sonnet-4-5-20250929".to_string(),
                allowed_models: vec![
                    "claude-haiku-4-5-20251001".to_string(),
                    "claude-sonnet-4-5-20250929".to_string(),
                ],
                max_calls_per_day: 500,
                max_calls_per_hour: 100,
                max_cost_per_month: 150.0,
                analysis_threshold: 0.5,
                batch_size: 30,
                optimize_for: "accuracy".to_string(),
                notify_threshold: 0.7,
            },
            "orbit" => Self {
                plan: plan.to_string(),
                default_model: "claude-opus-4-5-20251101".to_string(),
                allowed_models: vec![
                    "claude-haiku-4-5-20251001".to_string(),
                    "claude-sonnet-4-5-20250929".to_string(),
                    "claude-opus-4-5-20251101".to_string(),
                ],
                max_calls_per_day: 0,    // Unlimited
                max_calls_per_hour: 0,   // Unlimited
                max_cost_per_month: 0.0, // Unlimited
                analysis_threshold: 0.4,
                batch_size: 20,
                optimize_for: "speed".to_string(),
                notify_threshold: 0.9,
            },
            _ => {
                // Default to trial configuration
                tracing::warn!("Unknown plan '{}', using trial defaults", plan);
                Self::for_plan("trial")
            }
        }
    }

    /// Check if a model is allowed for this plan
    pub fn is_model_allowed(&self, model: &str) -> bool {
        self.allowed_models.iter().any(|m| m == model)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_trial_plan() {
        let config = PlanConfig::for_plan("trial");
        assert_eq!(config.plan, "trial");
        assert_eq!(config.default_model, "claude-haiku-4-5-20251001");
        assert_eq!(config.max_calls_per_day, 100);
        assert_eq!(config.max_cost_per_month, 10.0);
        assert!(config.is_model_allowed("claude-haiku-4-5-20251001"));
        assert!(!config.is_model_allowed("claude-opus-4-5-20251101"));
    }

    #[test]
    fn test_radar_plan() {
        let config = PlanConfig::for_plan("radar");
        assert_eq!(config.plan, "radar");
        assert_eq!(config.max_calls_per_day, 200);
        assert_eq!(config.max_cost_per_month, 20.0);
    }

    #[test]
    fn test_tensor_plan() {
        let config = PlanConfig::for_plan("tensor");
        assert_eq!(config.plan, "tensor");
        assert_eq!(config.default_model, "claude-sonnet-4-5-20250929");
        assert!(config.is_model_allowed("claude-haiku-4-5-20251001"));
        assert!(config.is_model_allowed("claude-sonnet-4-5-20250929"));
        assert!(!config.is_model_allowed("claude-opus-4-5-20251101"));
    }

    #[test]
    fn test_orbit_plan() {
        let config = PlanConfig::for_plan("orbit");
        assert_eq!(config.plan, "orbit");
        assert_eq!(config.max_calls_per_day, 0); // Unlimited
        assert_eq!(config.max_cost_per_month, 0.0); // Unlimited
        assert!(config.is_model_allowed("claude-opus-4-5-20251101"));
    }

    #[test]
    fn test_unknown_plan_defaults_to_trial() {
        let config = PlanConfig::for_plan("unknown");
        assert_eq!(config.plan, "trial");
    }

    #[test]
    fn test_legacy_plan_names() {
        let pilot = PlanConfig::for_plan("pilot");
        assert_eq!(pilot.plan, "pilot");

        let lock = PlanConfig::for_plan("lock");
        assert_eq!(lock.plan, "lock");
        assert_eq!(lock.default_model, "claude-sonnet-4-5-20250929");
    }
}
