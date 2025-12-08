package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeClient handles interactions with the Claude API
type ClaudeClient struct {
	client *anthropic.Client
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ClaudeClient{
		client: &client,
	}, nil
}

// Provider returns the provider name
func (c *ClaudeClient) Provider() string {
	return "anthropic"
}

// Model returns the default model name
func (c *ClaudeClient) Model() string {
	return "claude-haiku-4-5-20251001" // Default Claude model
}

// AnalyzeAnomaly analyzes an anomaly using Claude
func (c *ClaudeClient) AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error) {
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(1024),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to call Claude API: %w", err)
	}

	if len(message.Content) == 0 {
		return "", 0, 0, fmt.Errorf("empty response from Claude")
	}

	response := message.Content[0].Text
	inputTokens := message.Usage.InputTokens
	outputTokens := message.Usage.OutputTokens

	return response, inputTokens, outputTokens, nil
}
