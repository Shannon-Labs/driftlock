use super::{AiClient, AnalysisResponse};
use crate::cost;
use crate::error::{AiError, Result};
use async_trait::async_trait;
use reqwest::Client;
use serde::{Deserialize, Serialize};

const ANTHROPIC_API_URL: &str = "https://api.anthropic.com/v1/messages";
const ANTHROPIC_VERSION: &str = "2023-06-01";

/// Claude API client using reqwest
#[derive(Debug)]
pub struct ClaudeClient {
    client: Client,
    api_key: String,
    default_model: String,
}

impl ClaudeClient {
    /// Create a new Claude client
    ///
    /// # Arguments
    /// * `api_key` - Anthropic API key
    /// * `default_model` - Default model to use (e.g., "claude-haiku-4-5-20251001")
    pub fn new(api_key: String, default_model: Option<String>) -> Result<Self> {
        if api_key.is_empty() {
            return Err(AiError::MissingApiKey(
                "ANTHROPIC_API_KEY is required".to_string(),
            ));
        }

        let client = Client::builder()
            .timeout(std::time::Duration::from_secs(60))
            .build()
            .map_err(|e| AiError::Other(format!("Failed to create HTTP client: {}", e)))?;

        Ok(Self {
            client,
            api_key,
            default_model: default_model.unwrap_or_else(|| "claude-haiku-4-5-20251001".to_string()),
        })
    }

    /// Create a Claude client from environment variable
    pub fn from_env() -> Result<Self> {
        let api_key = std::env::var("ANTHROPIC_API_KEY")
            .map_err(|_| AiError::MissingApiKey("ANTHROPIC_API_KEY not set".to_string()))?;

        let default_model = std::env::var("ANTHROPIC_MODEL").ok();

        Self::new(api_key, default_model)
    }
}

#[async_trait]
impl AiClient for ClaudeClient {
    fn provider(&self) -> &str {
        "anthropic"
    }

    fn default_model(&self) -> &str {
        &self.default_model
    }

    async fn analyze_anomaly(&self, model: &str, prompt: &str) -> Result<AnalysisResponse> {
        let request = ClaudeRequest {
            model: model.to_string(),
            max_tokens: 1024,
            messages: vec![ClaudeMessage {
                role: "user".to_string(),
                content: prompt.to_string(),
            }],
        };

        let response = self
            .client
            .post(ANTHROPIC_API_URL)
            .header("x-api-key", &self.api_key)
            .header("anthropic-version", ANTHROPIC_VERSION)
            .header("content-type", "application/json")
            .json(&request)
            .send()
            .await?;

        if !response.status().is_success() {
            let status = response.status();
            let body = response.text().await.unwrap_or_default();
            return Err(AiError::Api(format!(
                "API request failed with status {}: {}",
                status, body
            )));
        }

        let claude_response: ClaudeResponse = response.json().await?;

        if claude_response.content.is_empty() {
            return Err(AiError::EmptyResponse);
        }

        let text = claude_response
            .content
            .first()
            .map(|c| c.text.clone())
            .unwrap_or_default();

        let input_tokens = claude_response.usage.input_tokens;
        let output_tokens = claude_response.usage.output_tokens;

        // Calculate cost with 15% margin
        let cost_usd =
            cost::calculate_cost(model, input_tokens, output_tokens).unwrap_or_else(|| {
                tracing::warn!("Unknown model '{}', cost calculation unavailable", model);
                0.0
            });

        Ok(AnalysisResponse {
            text,
            input_tokens,
            output_tokens,
            cost_usd,
        })
    }
}

// Claude API request/response types

#[derive(Debug, Serialize)]
struct ClaudeRequest {
    model: String,
    max_tokens: u32,
    messages: Vec<ClaudeMessage>,
}

#[derive(Debug, Serialize)]
struct ClaudeMessage {
    role: String,
    content: String,
}

#[derive(Debug, Deserialize)]
struct ClaudeResponse {
    content: Vec<ClaudeContent>,
    usage: ClaudeUsage,
}

#[derive(Debug, Deserialize)]
struct ClaudeContent {
    #[serde(rename = "type")]
    #[allow(dead_code)]
    content_type: String,
    text: String,
}

#[derive(Debug, Deserialize)]
struct ClaudeUsage {
    input_tokens: i64,
    output_tokens: i64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_client_creation_empty_key() {
        let result = ClaudeClient::new("".to_string(), None);
        assert!(result.is_err());
        assert!(matches!(result.unwrap_err(), AiError::MissingApiKey(_)));
    }

    // Integration tests require ANTHROPIC_API_KEY
    #[tokio::test]
    #[ignore] // Run with: cargo test -- --ignored
    async fn test_analyze_anomaly_integration() {
        let client = ClaudeClient::from_env().expect("ANTHROPIC_API_KEY must be set");

        let prompt = "Analyze this anomaly: Sudden spike in error rate from 0.1% to 5%";
        let response = client
            .analyze_anomaly("claude-haiku-4-5-20251001", prompt)
            .await
            .expect("API call failed");

        assert!(!response.text.is_empty());
        assert!(response.input_tokens > 0);
        assert!(response.output_tokens > 0);
        assert!(response.cost_usd > 0.0);

        println!("Response: {}", response.text);
        println!(
            "Tokens: {} in, {} out",
            response.input_tokens, response.output_tokens
        );
        println!("Cost: ${:.6}", response.cost_usd);
    }
}
