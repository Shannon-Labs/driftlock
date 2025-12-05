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

// OllamaClient handles interactions with local Ollama instance
// Ollama provides an OpenAI-compatible API at /v1/chat/completions
type OllamaClient struct {
	baseURL      string
	defaultModel string
	httpClient   *http.Client
}

// NewOllamaClient creates a new Ollama client for local LLM inference
func NewOllamaClient() (*OllamaClient, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "mistral" // Default to Mistral, change to ministral:3b if available
	}

	return &OllamaClient{
		baseURL:      baseURL,
		defaultModel: model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for local inference
		},
	}, nil
}

// Provider returns the provider name
func (c *OllamaClient) Provider() string {
	return "ollama"
}

// AnalyzeAnomaly analyzes an anomaly using local Ollama instance
func (c *OllamaClient) AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error) {
	if model == "" {
		model = c.defaultModel
	}

	// Use Ollama's OpenAI-compatible API
	reqBody := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Ollama's OpenAI-compatible endpoint
	url := c.baseURL + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// No Authorization header needed for local Ollama

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to call Ollama API (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result chatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	if result.Error != nil {
		return "", 0, 0, fmt.Errorf("Ollama error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", 0, 0, fmt.Errorf("empty response from Ollama")
	}

	response := result.Choices[0].Message.Content
	inputTokens := result.Usage.PromptTokens
	outputTokens := result.Usage.CompletionTokens

	return response, inputTokens, outputTokens, nil
}

// CheckHealth verifies Ollama is running and the model is available
func (c *OllamaClient) CheckHealth(ctx context.Context) error {
	url := c.baseURL + "/api/tags"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama not reachable at %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check failed: status %d", resp.StatusCode)
	}

	return nil
}
