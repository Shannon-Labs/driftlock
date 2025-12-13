//! AI-powered anomaly explanations for Driftlock
//!
//! This crate provides integration with Claude API for generating human-readable
//! explanations of detected anomalies. It includes:
//!
//! - Cost calculation with 15% margin
//! - Plan-based model selection (Trial/Radar → Haiku, Tensor → Sonnet, Orbit → Opus)
//! - Prompt injection filtering for safe user data inclusion
//! - Token usage tracking and cost estimation
//!
//! # Example
//!
//! ```rust,no_run
//! use driftlock_ai::{ClaudeClient, AiClient, sanitize_for_prompt};
//!
//! #[tokio::main]
//! async fn main() -> Result<(), Box<dyn std::error::Error>> {
//!     let client = ClaudeClient::from_env()?;
//!
//!     let event_data = r#"{"error": "Connection timeout", "count": 50}"#;
//!     let sanitized = sanitize_for_prompt(event_data, 2048);
//!
//!     let prompt = format!(
//!         "Explain this anomaly:\n\nEvent data:\n{}\n\nProvide a brief explanation.",
//!         sanitized
//!     );
//!
//!     let response = client
//!         .analyze_anomaly("claude-haiku-4-5-20251001", &prompt)
//!         .await?;
//!
//!     println!("Explanation: {}", response.text);
//!     println!("Cost: ${:.6}", response.cost_usd);
//!
//!     Ok(())
//! }
//! ```

pub mod client;
pub mod config;
pub mod cost;
pub mod error;
pub mod sanitize;

// Re-export commonly used items
pub use client::claude::ClaudeClient;
pub use client::{AiClient, AnalysisResponse};
pub use config::PlanConfig;
pub use cost::{calculate_cost, get_model_pricing, get_model_tier, ModelPricing};
pub use error::{AiError, Result};
pub use sanitize::sanitize_for_prompt;
