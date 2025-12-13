pub mod claude;

use crate::error::Result;
use async_trait::async_trait;

/// Response from analyzing an anomaly
#[derive(Debug, Clone)]
pub struct AnalysisResponse {
    /// The explanation text
    pub text: String,
    /// Number of input tokens used
    pub input_tokens: i64,
    /// Number of output tokens used
    pub output_tokens: i64,
    /// Estimated cost in USD
    pub cost_usd: f64,
}

/// Trait for AI clients that can analyze anomalies
#[async_trait]
pub trait AiClient: Send + Sync {
    /// Get the provider name (e.g., "anthropic")
    fn provider(&self) -> &str;

    /// Get the default model name
    fn default_model(&self) -> &str;

    /// Analyze an anomaly with the given prompt
    ///
    /// # Arguments
    /// * `model` - The model to use (e.g., "claude-haiku-4-5-20251001")
    /// * `prompt` - The prompt describing the anomaly
    ///
    /// # Returns
    /// Analysis response with text, token usage, and cost
    async fn analyze_anomaly(&self, model: &str, prompt: &str) -> Result<AnalysisResponse>;
}
