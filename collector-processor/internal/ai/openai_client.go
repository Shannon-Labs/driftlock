package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient handles interactions with OpenAI-compatible APIs (OpenAI, Z.AI, etc.)
type OpenAIClient struct {
	apiKey       string
	baseURL      string
	defaultModel string
	httpClient   *http.Client
}

// OpenAI API request/response structures
type chatCompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []chatMessage   `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// NewOpenAIClient creates a new OpenAI-compatible client
func NewOpenAIClient(apiKey, baseURL, defaultModel string) (*OpenAIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Default to OpenAI API if no base URL provided
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// Default model
	if defaultModel == "" {
		defaultModel = "gpt-4o-mini"
	}

	return &OpenAIClient{
		apiKey:       apiKey,
		baseURL:      baseURL,
		defaultModel: defaultModel,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Provider returns the provider name
func (c *OpenAIClient) Provider() string {
	// Detect provider from base URL
	if c.baseURL == "https://api.z.ai/api/coding/paas/v4" {
		return "zai"
	}
	if c.baseURL == "https://api.openai.com/v1" {
		return "openai"
	}
	return "openai-compatible"
}

// AnalyzeAnomaly analyzes an anomaly using the OpenAI-compatible API
func (c *OpenAIClient) AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error) {
	// Use provided model or fall back to default
	if model == "" {
		model = c.defaultModel
	}

	// Build request
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

	// Create HTTP request
	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var result chatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	// Check for API error
	if result.Error != nil {
		return "", 0, 0, fmt.Errorf("API error: %s (type: %s, code: %s)",
			result.Error.Message, result.Error.Type, result.Error.Code)
	}

	// Check for valid response
	if len(result.Choices) == 0 {
		return "", 0, 0, fmt.Errorf("empty response from API")
	}

	response := result.Choices[0].Message.Content
	inputTokens := result.Usage.PromptTokens
	outputTokens := result.Usage.CompletionTokens

	return response, inputTokens, outputTokens, nil
}
