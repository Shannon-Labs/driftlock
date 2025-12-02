package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CostControlConfig represents tenant-specific AI cost control settings
type CostControlConfig struct {
	TenantID          string    `json:"tenant_id"`
	Plan              string    `json:"plan"` // "trial", "radar", "lock", "orbit"
	Enabled           bool      `json:"enabled"`
	Models            []string  `json:"models"` // Allowed models: ["haiku-4.5", "sonnet-4.5", "opus-4.5"]
	MaxCallsPerDay    int       `json:"max_calls_per_day"`
	MaxCallsPerHour   int       `json:"max_calls_per_hour"`
	MaxCostPerMonth   float64   `json:"max_cost_per_month"` // Hard limit in USD
	AnalysisThreshold float64   `json:"analysis_threshold"` // 0.0-1.0, only analyze anomalies above this
	BatchSize         int       `json:"batch_size"`         // Number of requests to batch
	OptimizeFor       string    `json:"optimize_for"`       // "speed", "cost", or "accuracy"
	NotifyThreshold   float64   `json:"notify_threshold"`   // Alert when using this % of limit
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Default cost control configurations by plan
var DefaultConfigs = map[string]CostControlConfig{
	"trial": {
		Enabled:           true,  // Trial uses Z.AI/GLM (cheap) via AI_PROVIDER=zai
		Models:            []string{"glm-4"}, // Z.AI model - actual model set by AI_MODEL env
		MaxCallsPerDay:    100,
		MaxCallsPerHour:   20,
		MaxCostPerMonth:   0, // Z.AI cost tracked separately, not Claude pricing
		AnalysisThreshold: 0.7,
		BatchSize:         50,
		OptimizeFor:       "cost",
		NotifyThreshold:   0.9,
	},
	"radar": {
		Enabled:           true,
		Models:            []string{"claude-haiku-4-5-20251001"}, // Haiku only for lowest paid tier
		MaxCallsPerDay:    200,
		MaxCallsPerHour:   40,
		MaxCostPerMonth:   20.00, // ~$20/month AI budget (Haiku is cheap)
		AnalysisThreshold: 0.6,
		BatchSize:         50,
		OptimizeFor:       "cost",
		NotifyThreshold:   0.8,
	},
	"tensor": { // Renamed from "lock"
		Enabled:           true,
		Models:            []string{"claude-opus-4-5-20251101"}, // Opus for pro tier (skip Sonnet)
		MaxCallsPerDay:    500,
		MaxCallsPerHour:   100,
		MaxCostPerMonth:   150.00, // $150/month AI budget for Opus
		AnalysisThreshold: 0.5,
		BatchSize:         30,
		OptimizeFor:       "accuracy",
		NotifyThreshold:   0.7,
	},
	"orbit": {
		Enabled:           true,
		Models:            []string{"claude-opus-4-5-20251101"}, // Opus for enterprise
		MaxCallsPerDay:    0, // Unlimited
		MaxCallsPerHour:   0, // Unlimited
		MaxCostPerMonth:   0, // No limit (configurable per contract)
		AnalysisThreshold: 0.4,
		BatchSize:         20,
		OptimizeFor:       "speed",
		NotifyThreshold:   0.9,
	},
	// Backward compatibility aliases
	"lock": { // Legacy name for tensor
		Enabled:           true,
		Models:            []string{"claude-opus-4-5-20251101"},
		MaxCallsPerDay:    500,
		MaxCallsPerHour:   100,
		MaxCostPerMonth:   150.00,
		AnalysisThreshold: 0.5,
		BatchSize:         30,
		OptimizeFor:       "accuracy",
		NotifyThreshold:   0.7,
	},
	"pilot": { // Legacy name for trial/free
		Enabled:           true,
		Models:            []string{"glm-4"},
		MaxCallsPerDay:    100,
		MaxCallsPerHour:   20,
		MaxCostPerMonth:   0,
		AnalysisThreshold: 0.7,
		BatchSize:         50,
		OptimizeFor:       "cost",
		NotifyThreshold:   0.9,
	},
}

// ConfigRepo defines the interface for configuration persistence
type ConfigRepo interface {
	GetConfig(ctx context.Context, tenantID string) (*CostControlConfig, error)
	UpdateConfig(ctx context.Context, config *CostControlConfig) error
	GetCurrentUsage(ctx context.Context, tenantID string) (*CurrentUsage, error)
}

// ConfigRepository handles persistence of cost control configs
type ConfigRepository struct {
	pool *pgxpool.Pool
}

func NewConfigRepository(pool *pgxpool.Pool) *ConfigRepository {
	return &ConfigRepository{pool: pool}
}

// GetConfig retrieves cost control config for a tenant
func (r *ConfigRepository) GetConfig(ctx context.Context, tenantID string) (*CostControlConfig, error) {
	var config CostControlConfig
	var modelsJSON string

	query := `
		SELECT c.tenant_id, t.plan, c.enabled, c.models, c.max_calls_per_day, c.max_calls_per_hour,
			   c.max_cost_per_month, c.analysis_threshold, c.batch_size, c.optimize_for,
			   c.notify_threshold, c.created_at, c.updated_at
		FROM ai_cost_control_configs c
		JOIN tenants t ON t.id = c.tenant_id
		WHERE c.tenant_id = $1
	`

	err := r.pool.QueryRow(ctx, query, tenantID).Scan(
		&config.TenantID,
		&config.Plan,
		&config.Enabled,
		&modelsJSON,
		&config.MaxCallsPerDay,
		&config.MaxCallsPerHour,
		&config.MaxCostPerMonth,
		&config.AnalysisThreshold,
		&config.BatchSize,
		&config.OptimizeFor,
		&config.NotifyThreshold,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		// Return default config based on plan
		plan, err := r.getTenantPlan(ctx, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenant plan: %w", err)
		}

		defaultConfig := DefaultConfigs[plan]
		defaultConfig.TenantID = tenantID
		defaultConfig.Plan = plan
		return &defaultConfig, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if err := json.Unmarshal([]byte(modelsJSON), &config.Models); err != nil {
		return nil, fmt.Errorf("failed to unmarshal models: %w", err)
	}

	return &config, nil
}

// UpdateConfig saves or updates cost control config for a tenant
func (r *ConfigRepository) UpdateConfig(ctx context.Context, config *CostControlConfig) error {
	modelsJSON, err := json.Marshal(config.Models)
	if err != nil {
		return fmt.Errorf("failed to marshal models: %w", err)
	}

	query := `
		INSERT INTO ai_cost_control_configs (
			tenant_id, enabled, models, max_calls_per_day, max_calls_per_hour,
			max_cost_per_month, analysis_threshold, batch_size, optimize_for,
			notify_threshold, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (tenant_id) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			models = EXCLUDED.models,
			max_calls_per_day = EXCLUDED.max_calls_per_day,
			max_calls_per_hour = EXCLUDED.max_calls_per_hour,
			max_cost_per_month = EXCLUDED.max_cost_per_month,
			analysis_threshold = EXCLUDED.analysis_threshold,
			batch_size = EXCLUDED.batch_size,
			optimize_for = EXCLUDED.optimize_for,
			notify_threshold = EXCLUDED.notify_threshold,
			updated_at = NOW()
	`

	_, err = r.pool.Exec(ctx, query,
		config.TenantID,
		config.Enabled,
		modelsJSON,
		config.MaxCallsPerDay,
		config.MaxCallsPerHour,
		config.MaxCostPerMonth,
		config.AnalysisThreshold,
		config.BatchSize,
		config.OptimizeFor,
		config.NotifyThreshold,
	)

	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func (r *ConfigRepository) getTenantPlan(ctx context.Context, tenantID string) (string, error) {
	var plan string
	err := r.pool.QueryRow(ctx,
		"SELECT COALESCE(plan, 'trial') FROM tenants WHERE id = $1",
		tenantID,
	).Scan(&plan)

	if err != nil {
		return "", err
	}

	return plan, nil
}

// CostLimiter enforces cost control limits
type CostLimiter struct {
	repo ConfigRepo
}

func NewCostLimiter(repo ConfigRepo) *CostLimiter {
	return &CostLimiter{repo: repo}
}

// CheckLimits verifies if a request is within limits
func (l *CostLimiter) CheckLimits(ctx context.Context, tenantID string, model string, estimatedCost float64) (*LimitCheckResult, error) {
	config, err := l.repo.GetConfig(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	result := &LimitCheckResult{
		Allowed: true,
		Reason:  "",
		Config:  config,
	}

	// Check if AI is enabled
	if !config.Enabled {
		result.Allowed = false
		result.Reason = "AI analysis is disabled for this tenant"
		return result, nil
	}

	// Check if model is allowed
	modelAllowed := false
	for _, allowedModel := range config.Models {
		if allowedModel == model {
			modelAllowed = true
			break
		}
	}
	if !modelAllowed {
		result.Allowed = false
		result.Reason = fmt.Sprintf("Model %s is not allowed for this plan", model)
		return result, nil
	}

	// Check usage limits
	usage, err := l.repo.GetCurrentUsage(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current usage: %w", err)
	}

	// Check daily calls
	if config.MaxCallsPerDay > 0 && usage.CallsToday >= config.MaxCallsPerDay {
		result.Allowed = false
		result.Reason = "Daily call limit exceeded"
		return result, nil
	}

	// Check hourly calls
	if config.MaxCallsPerHour > 0 && usage.CallsThisHour >= config.MaxCallsPerHour {
		result.Allowed = false
		result.Reason = "Hourly call limit exceeded"
		return result, nil
	}

	// Check monthly cost
	if config.MaxCostPerMonth > 0 && (usage.CostThisMonth+estimatedCost) > config.MaxCostPerMonth {
		result.Allowed = false
		result.Reason = "Monthly cost limit would be exceeded"
		return result, nil
	}

	// Check if approaching notification threshold
	if config.MaxCostPerMonth > 0 {
		costPercent := (usage.CostThisMonth + estimatedCost) / config.MaxCostPerMonth
		if costPercent >= config.NotifyThreshold {
			result.ShouldNotify = true
			result.NotifyMessage = fmt.Sprintf(
				"You've used %.1f%% of your monthly AI budget",
				costPercent*100,
			)
		}
	}

	return result, nil
}

// LimitCheckResult contains the result of a limit check
type LimitCheckResult struct {
	Allowed       bool               `json:"allowed"`
	Reason        string             `json:"reason"`
	ShouldNotify  bool               `json:"should_notify"`
	NotifyMessage string             `json:"notify_message"`
	Config        *CostControlConfig `json:"config"`
}

// CurrentUsage represents current usage metrics
type CurrentUsage struct {
	CallsToday    int       `json:"calls_today"`
	CallsThisHour int       `json:"calls_this_hour"`
	CostThisMonth float64   `json:"cost_this_month"`
	LastResetTime time.Time `json:"last_reset_time"`
}

func (r *ConfigRepository) GetCurrentUsage(ctx context.Context, tenantID string) (*CurrentUsage, error) {
	usage := &CurrentUsage{
		LastResetTime: time.Now(),
	}

	// Get calls today
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_usage
		 WHERE tenant_id = $1 AND DATE(created_at) = CURRENT_DATE`,
		tenantID,
	).Scan(&usage.CallsToday)
	if err != nil {
		return nil, err
	}

	// Get calls this hour
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ai_usage
		 WHERE tenant_id = $1 AND created_at >= date_trunc('hour', NOW())`,
		tenantID,
	).Scan(&usage.CallsThisHour)
	if err != nil {
		return nil, err
	}

	// Get cost this month
	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(cost_usd), 0) FROM ai_usage
		 WHERE tenant_id = $1 AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', NOW())`,
		tenantID,
	).Scan(&usage.CostThisMonth)
	if err != nil {
		return nil, err
	}

	return usage, nil
}
