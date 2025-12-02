package ai

import (
	"context"
	"time"
)

// MockAIClient is a mock implementation of AIClient for testing
type MockAIClient struct {
	AnalyzeAnomalyFunc func(ctx context.Context, model string, prompt string) (string, int64, int64, error)
}

func (m *MockAIClient) AnalyzeAnomaly(ctx context.Context, model string, prompt string) (string, int64, int64, error) {
	if m.AnalyzeAnomalyFunc != nil {
		return m.AnalyzeAnomalyFunc(ctx, model, prompt)
	}
	return "Mock explanation", 100, 50, nil
}

func (m *MockAIClient) Provider() string {
	return "mock"
}

// MockConfigRepository mocks ConfigRepository for testing
type MockConfigRepository struct {
	GetConfigFunc func(ctx context.Context, tenantID string) (*CostControlConfig, error)
}

func (m *MockConfigRepository) GetConfig(ctx context.Context, tenantID string) (*CostControlConfig, error) {
	if m.GetConfigFunc != nil {
		return m.GetConfigFunc(ctx, tenantID)
	}
	return &CostControlConfig{
		Enabled: true,
		Models:  []string{"claude-haiku-4-5-20251001"},
		Plan:    "radar",
	}, nil
}

func (m *MockConfigRepository) UpdateConfig(ctx context.Context, config *CostControlConfig) error {
	return nil
}

func (m *MockConfigRepository) GetCurrentUsage(ctx context.Context, tenantID string) (*CurrentUsage, error) {
	return &CurrentUsage{
		CallsToday:    0,
		CallsThisHour: 0,
		CostThisMonth: 0.0,
		LastResetTime: time.Now(),
	}, nil
}
