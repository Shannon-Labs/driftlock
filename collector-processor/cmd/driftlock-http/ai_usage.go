package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/internal/ai"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AIUsageResponse represents the AI usage response
type AIUsageResponse struct {
	CallsUsed    int     `json:"calls_used"`
	CallsLimit   int     `json:"calls_limit"`
	CostUsed     float64 `json:"cost_used"`
	CostLimit    float64 `json:"cost_limit"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	ModelType    string  `json:"model_type"`
	Forecast     struct {
		EstimatedMonthly float64 `json:"estimated_monthly"`
		DaysRemaining    int     `json:"days_remaining"`
	} `json:"forecast"`
}

// AIConfigRequest represents AI configuration update request
type AIConfigRequest struct {
	AnalysisThreshold float64 `json:"analysis_threshold"`
	OptimizeFor       string  `json:"optimize_for"`
	MaxCost           float64 `json:"max_cost"`
}

// aiUsageHandler handles GET /v1/me/usage/ai
func aiUsageHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// Get AI config
		configRepo := ai.NewConfigRepository(store.pool)
		config, err := configRepo.GetConfig(r.Context(), tc.Tenant.ID.String())
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get config: %w", err))
			return
		}

		// Get current usage
		usage, err := getCurrentAIUsage(r.Context(), store.pool, tc.Tenant.ID.String())
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get usage: %w", err))
			return
		}

		// Build response
		response := AIUsageResponse{
			CallsUsed:    usage.CallsThisMonth,
			CallsLimit:   config.MaxCallsPerDay * 30, // Approximate monthly limit
			CostUsed:     usage.CostThisMonth,
			CostLimit:    config.MaxCostPerMonth,
			InputTokens:  usage.InputTokensThisMonth,
			OutputTokens: usage.OutputTokensThisMonth,
			ModelType:    getPrimaryModel(config.Models),
		}

		// Calculate forecast
		now := time.Now()
		daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		daysPassed := now.Day()
		dailyRate := response.CostUsed / float64(daysPassed)

		response.Forecast.EstimatedMonthly = dailyRate * float64(daysInMonth)
		response.Forecast.DaysRemaining = daysInMonth - daysPassed

		writeJSON(w, r, http.StatusOK, response)
	}
}

// aiConfigHandler handles POST /v1/me/ai/config
func aiConfigHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		var req AIConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
			return
		}

		// Validate request
		if req.AnalysisThreshold < 0 || req.AnalysisThreshold > 1 {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("analysis_threshold must be between 0 and 1"))
			return
		}

		if req.OptimizeFor != "" && req.OptimizeFor != "cost" && req.OptimizeFor != "speed" && req.OptimizeFor != "accuracy" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("optimize_for must be one of: cost, speed, accuracy"))
			return
		}

		// Get existing config
		configRepo := ai.NewConfigRepository(store.pool)
		config, err := configRepo.GetConfig(r.Context(), tc.Tenant.ID.String())
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get config: %w", err))
			return
		}

		// Update only allowed fields
		if req.AnalysisThreshold > 0 {
			config.AnalysisThreshold = req.AnalysisThreshold
		}
		if req.OptimizeFor != "" {
			config.OptimizeFor = req.OptimizeFor
		}
		if req.MaxCost > 0 && tc.Tenant.Plan == "orbit" { // Only allow for enterprise
			config.MaxCostPerMonth = req.MaxCost
		}

		config.TenantID = tc.Tenant.ID.String()
		config.UpdatedAt = time.Now()

		// Save updated config
		if err := configRepo.UpdateConfig(r.Context(), config); err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to update config: %w", err))
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]string{
			"status":  "success",
			"message": "AI configuration updated successfully",
		})
	}
}

// trackAIUsage tracks AI API usage
func trackAIUsage(store *store, tenantID, streamID, modelType string, inputTokens, outputTokens int64, cost float64) error {
	ctx := context.Background()

	// Calculate total charge with 15% margin
	marginPercent := 15.0
	totalCharge := cost * (1 + marginPercent/100)

	// Insert into ai_usage table
	query := `
		INSERT INTO ai_usage (
			tenant_id, stream_id, model_type, input_tokens, output_tokens,
			cost_usd, margin_percent, total_charge_usd, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`

	_, err := store.pool.Exec(ctx, query,
		tenantID, streamID, modelType, inputTokens, outputTokens,
		cost, marginPercent, totalCharge,
	)

	return err
}

// AIUsageStats represents usage statistics
type AIUsageStats struct {
	CallsThisMonth        int     `json:"calls_this_month"`
	CostThisMonth         float64 `json:"cost_this_month"`
	InputTokensThisMonth  int64   `json:"input_tokens_this_month"`
	OutputTokensThisMonth int64   `json:"output_tokens_this_month"`
}

func getCurrentAIUsage(ctx context.Context, pool *pgxpool.Pool, tenantID string) (*AIUsageStats, error) {
	stats := &AIUsageStats{}

	// Get current month usage
	query := `
		SELECT
			COUNT(*) as calls,
			COALESCE(SUM(cost_usd), 0) as cost,
			COALESCE(SUM(input_tokens), 0) as input_tokens,
			COALESCE(SUM(output_tokens), 0) as output_tokens
		FROM ai_usage
		WHERE tenant_id = $1
		  AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', NOW())
	`

	err := pool.QueryRow(ctx, query, tenantID).Scan(
		&stats.CallsThisMonth,
		&stats.CostThisMonth,
		&stats.InputTokensThisMonth,
		&stats.OutputTokensThisMonth,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func getPrimaryModel(models []string) string {
	if len(models) == 0 {
		return "none"
	}

	// Return the highest tier model available
	if contains(models, "claude-opus-4-5-20251101") {
		return "opus-4.5"
	}
	if contains(models, "claude-sonnet-4-5-20250929") {
		return "sonnet-4.5"
	}
	if contains(models, "claude-haiku-4-5-20251001") {
		return "haiku-4.5"
	}

	return models[0]
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
