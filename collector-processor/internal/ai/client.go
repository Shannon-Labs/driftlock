package ai

import (
	"context"
	"fmt"
	"os"
)

// AIClient defines the interface for AI model interactions
type AIClient interface {
	AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error)
	Provider() string
}

// ProviderConfig holds configuration for an AI provider
type ProviderConfig struct {
	Provider string // "anthropic", "openai", "zai"
	APIKey   string
	BaseURL  string // For OpenAI-compatible APIs (Z.AI, etc.)
	Model    string // Default model to use
}

// NewAIClientFromEnv creates an AI client from environment variables
func NewAIClientFromEnv() (AIClient, error) {
	provider := os.Getenv("AI_PROVIDER")
	if provider == "" {
		// Default to anthropic if ANTHROPIC_API_KEY is set
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			provider = "anthropic"
		} else if os.Getenv("AI_API_KEY") != "" {
			provider = "openai"
		} else {
			return nil, fmt.Errorf("no AI provider configured: set AI_PROVIDER and AI_API_KEY or ANTHROPIC_API_KEY")
		}
	}

	switch provider {
	case "anthropic", "claude":
		return NewClaudeClient()
	case "openai", "zai", "glm":
		apiKey := os.Getenv("AI_API_KEY")
		baseURL := os.Getenv("AI_BASE_URL")
		model := os.Getenv("AI_MODEL")
		if apiKey == "" {
			return nil, fmt.Errorf("AI_API_KEY required for %s provider", provider)
		}
		return NewOpenAIClient(apiKey, baseURL, model)
	case "vertex", "vertex-claude", "gcp-claude":
		// Claude models via Google Vertex AI Model Garden
		project := os.Getenv("GOOGLE_CLOUD_PROJECT")
		region := os.Getenv("GOOGLE_CLOUD_REGION")
		model := os.Getenv("AI_MODEL")
		return NewVertexClaudeClient(project, region, model)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", provider)
	}
}
