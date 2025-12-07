package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GeminiClient handles interactions with Google's Gemini API
// Uses the OpenAI-compatible endpoint for simplicity
type GeminiClient struct {
	apiKey       string
	defaultModel string
	httpClient   *http.Client
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.5-flash" // Default to Flash 2.5
	}

	return &GeminiClient{
		apiKey:       apiKey,
		defaultModel: model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Provider returns the provider name
func (c *GeminiClient) Provider() string {
	return "gemini"
}

// AnalyzeAnomaly analyzes an anomaly using Gemini API
func (c *GeminiClient) AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error) {
	if model == "" {
		model = c.defaultModel
	}

	// Use Google's OpenAI-compatible endpoint
	// https://ai.google.dev/gemini-api/docs/openai
	reqBody := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   1024,
		Temperature: 0.3,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Gemini's OpenAI-compatible endpoint
	url := "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result chatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	if result.Error != nil {
		return "", 0, 0, fmt.Errorf("Gemini error: %s (type: %s)", result.Error.Message, result.Error.Type)
	}

	if len(result.Choices) == 0 {
		return "", 0, 0, fmt.Errorf("empty response from Gemini")
	}

	response := result.Choices[0].Message.Content
	inputTokens := result.Usage.PromptTokens
	outputTokens := result.Usage.CompletionTokens

	return response, inputTokens, outputTokens, nil
}
