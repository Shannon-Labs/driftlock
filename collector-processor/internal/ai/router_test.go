package ai

import (
	"testing"
	"time"
)

func TestCalculateCost(t *testing.T) {
	limiter := &CostLimiter{repo: nil} // Repo not needed for CalculateCost
	router := NewRouter(limiter)

	tests := []struct {
		name         string
		model        string
		inputTokens  int64
		outputTokens int64
		expectedCost float64
	}{
		{
			name:         "Haiku 4.5 Small",
			model:        "claude-haiku-4-5-20251001",
			inputTokens:  1000,
			outputTokens: 1000,
			expectedCost: 0.006000, // (1000/1M * 1.00) + (1000/1M * 5.00) = 0.001 + 0.005 = 0.006
		},
		{
			name:         "Sonnet 4.5 Medium",
			model:        "claude-sonnet-4-5-20250929",
			inputTokens:  1000000,
			outputTokens: 1000000,
			expectedCost: 18.000000, // (1 * 3.00) + (1 * 15.00) = 18.00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := router.CalculateCost(tt.model, tt.inputTokens, tt.outputTokens)
			if cost != tt.expectedCost {
				t.Errorf("CalculateCost() = %v, want %v", cost, tt.expectedCost)
			}
		})
	}
}

func TestEstimateCost_Batching(t *testing.T) {
	limiter := &CostLimiter{repo: nil}
	router := NewRouter(limiter)

	event := Event{
		ID:           "test-event",
		TenantID:     "test-tenant",
		StreamID:     "test-stream",
		AnomalyScore: 0.8,
		CreatedAt:    time.Now(),
	}

	// Haiku: Input $1/M, Output $5/M
	// Avg: 500 input, 200 output
	// Base Cost: (500/1M * 1) + (200/1M * 5) = 0.0005 + 0.001 = 0.0015
	// Margin 15%: 0.0015 * 1.15 = 0.001725

	// Batch Discount 50%: 0.0015 * 0.5 * 1.15 = 0.0008625 -> 0.000863

	tests := []struct {
		name         string
		plan         string
		isBatch      bool
		expectedCost float64
	}{
		{
			name:         "Orbit Plan with Batching",
			plan:         "orbit",
			isBatch:      true,
			expectedCost: 0.000863, // Discount applied
		},
		{
			name:         "Radar Plan with Batching",
			plan:         "radar",
			isBatch:      true,
			expectedCost: 0.001725, // No discount
		},
		{
			name:         "Orbit Plan without Batching",
			plan:         "orbit",
			isBatch:      false,
			expectedCost: 0.001725, // No discount
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := router.estimateCost("claude-haiku-4-5-20251001", event, tt.isBatch, tt.plan)
			if cost != tt.expectedCost {
				t.Errorf("estimateCost() = %v, want %v", cost, tt.expectedCost)
			}
		})
	}
}
