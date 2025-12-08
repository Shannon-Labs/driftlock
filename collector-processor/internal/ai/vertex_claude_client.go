package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

// VertexClaudeClient uses Claude models via Google Vertex AI Model Garden
type VertexClaudeClient struct {
	project  string
	region   string
	model    string
	endpoint string
}

// NewVertexClaudeClient creates a new Vertex AI Claude client
func NewVertexClaudeClient(project, region, model string) (*VertexClaudeClient, error) {
	if project == "" {
		project = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if project == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT required for Vertex AI")
	}

	if region == "" {
		region = os.Getenv("VERTEX_AI_REGION") // Prefer VERTEX_AI_REGION for Claude
	}
	if region == "" {
		region = os.Getenv("GOOGLE_CLOUD_REGION")
	}
	if region == "" {
		region = "us-east5" // Claude models available in us-east5
	}

	if model == "" {
		model = os.Getenv("AI_MODEL")
	}
	if model == "" {
		model = "claude-haiku-4-5@20251001"
	}

	// Build the Vertex AI endpoint for Claude
	endpoint := fmt.Sprintf(
		"https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:rawPredict",
		region, project, region, model,
	)

	return &VertexClaudeClient{
		project:  project,
		region:   region,
		model:    model,
		endpoint: endpoint,
	}, nil
}

// AnalyzeAnomaly sends a prompt to Claude via Vertex AI and returns the response
func (v *VertexClaudeClient) AnalyzeAnomaly(ctx context.Context, model, prompt string) (string, int64, int64, error) {
	// Get GCP credentials using Application Default Credentials
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to get GCP credentials: %w", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to get access token: %w", err)
	}

	// Build Anthropic-style request body (Vertex AI uses same format)
	reqBody := map[string]interface{}{
		"anthropic_version": "vertex-2023-10-16",
		"max_tokens":        1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Use model override if provided
	endpoint := v.endpoint
	if model != "" && model != v.model {
		endpoint = fmt.Sprintf(
			"https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:rawPredict",
			v.region, v.project, v.region, model,
		)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("Vertex AI error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse Anthropic-style response
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int64 `json:"input_tokens"`
			OutputTokens int64 `json:"output_tokens"`
		} `json:"usage"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error.Message != "" {
		return "", 0, 0, fmt.Errorf("Claude error: %s - %s", result.Error.Type, result.Error.Message)
	}

	// Extract text from response
	text := ""
	for _, content := range result.Content {
		if content.Type == "text" {
			text += content.Text
		}
	}

	return text, result.Usage.InputTokens, result.Usage.OutputTokens, nil
}

// Provider returns the provider name
func (v *VertexClaudeClient) Provider() string {
	return "vertex-claude"
}

// Model returns the default model name
func (v *VertexClaudeClient) Model() string {
	return v.model
}
