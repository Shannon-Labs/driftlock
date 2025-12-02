package ai

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Router handles intelligent routing of events to AI models based on cost controls
type Router struct {
	limiter  *CostLimiter
	models   map[string]*ModelClient
	trackers map[string]*UsageTracker
}

// ModelClient represents an AI model client
type ModelClient struct {
	Name          string
	CostPerInput  float64 // Cost per 1M input tokens
	CostPerOutput float64 // Cost per 1M output tokens
	Speed         int     // Relative speed (1-10)
	Quality       int     // Relative quality (1-10)
	BatchDiscount float64 // Discount factor for batch API (0.5 = 50% off)
}

// UsageTracker tracks usage statistics for optimization
type UsageTracker struct {
	TenantID    string
	Model       string
	TotalCost   float64
	TotalCalls  int
	AvgLatency  time.Duration
	SuccessRate float64
}

// RoutingDecision contains the routing decision for an event
type RoutingDecision struct {
	ShouldAnalyze   bool    `json:"should_analyze"`
	Model           string  `json:"model,omitempty"`
	BatchWithOthers bool    `json:"batch_with_others"`
	Priority        int     `json:"priority"`
	Reason          string  `json:"reason"`
	EstimatedCost   float64 `json:"estimated_cost"`
}

// Event represents an anomaly detection event
type Event struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	StreamID     string                 `json:"stream_id"`
	AnomalyScore float64                `json:"anomaly_score"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
}

// NewRouter creates a new smart router
func NewRouter(limiter *CostLimiter) *Router {
	return &Router{
		limiter: limiter,
		models: map[string]*ModelClient{
			// Direct Anthropic API models
			"claude-haiku-4-5-20251001": {
				Name:          "haiku-4.5",
				CostPerInput:  1.00,
				CostPerOutput: 5.00,
				Speed:         10,
				Quality:       7,
				BatchDiscount: 0.5, // 50% discount with batch
			},
			"claude-sonnet-4-5-20250929": {
				Name:          "sonnet-4.5",
				CostPerInput:  3.00,
				CostPerOutput: 15.00,
				Speed:         7,
				Quality:       9,
				BatchDiscount: 0.5, // 50% discount with batch
			},
			"claude-opus-4-5-20251101": {
				Name:          "opus-4.5",
				CostPerInput:  5.00,
				CostPerOutput: 25.00,
				Speed:         5,
				Quality:       10,
				BatchDiscount: 0.5, // 50% discount with batch
			},
			// Vertex AI Claude models (same models, GCP billing)
			"claude-haiku-4-5@20251001": {
				Name:          "vertex-haiku-4.5",
				CostPerInput:  1.00,
				CostPerOutput: 5.00,
				Speed:         10,
				Quality:       7,
				BatchDiscount: 0.5,
			},
			"claude-sonnet-4-5@20250929": {
				Name:          "vertex-sonnet-4.5",
				CostPerInput:  3.00,
				CostPerOutput: 15.00,
				Speed:         7,
				Quality:       9,
				BatchDiscount: 0.5,
			},
			"claude-opus-4-5@20251101": {
				Name:          "vertex-opus-4.5",
				CostPerInput:  15.00,
				CostPerOutput: 75.00,
				Speed:         5,
				Quality:       10,
				BatchDiscount: 0.5,
			},
		},
		trackers: make(map[string]*UsageTracker),
	}
}

// RouteEvent decides how to handle an event
func (r *Router) RouteEvent(ctx context.Context, event Event) (*RoutingDecision, error) {
	// Get tenant's cost control config
	config, err := r.limiter.repo.GetConfig(ctx, event.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	decision := &RoutingDecision{
		ShouldAnalyze: false,
		Reason:        "",
	}

	// Check if AI is enabled
	if !config.Enabled {
		decision.Reason = "AI analysis is disabled"
		return decision, nil
	}

	// Check if anomaly meets threshold
	if event.AnomalyScore < config.AnalysisThreshold {
		decision.Reason = fmt.Sprintf("Anomaly score %.2f below threshold %.2f",
			event.AnomalyScore, config.AnalysisThreshold)
		return decision, nil
	}

	// Select best model based on configuration and optimization target
	model, reason := r.selectModel(config, event)
	if model == "" {
		decision.Reason = reason
		return decision, nil
	}

	// Estimate cost (assume batch processing by default)
	estimatedCost := r.estimateCost(model, event, true, config.Plan)

	// Check limits
	limitResult, err := r.limiter.CheckLimits(ctx, event.TenantID, model, estimatedCost)
	if err != nil {
		return nil, err
	}

	if !limitResult.Allowed {
		decision.Reason = limitResult.Reason
		return decision, nil
	}

	// Build routing decision
	decision.ShouldAnalyze = true
	decision.Model = model
	decision.EstimatedCost = estimatedCost
	decision.Reason = reason

	// Determine batching and priority based on optimization target
	switch config.OptimizeFor {
	case "cost":
		decision.BatchWithOthers = true
		decision.Priority = 1 // Low priority
	case "speed":
		decision.BatchWithOthers = false
		decision.Priority = 10 // High priority
	case "accuracy":
		// Use best available model, don't batch if high confidence
		if event.AnomalyScore > 0.9 {
			decision.BatchWithOthers = false
			decision.Priority = 8
		} else {
			decision.BatchWithOthers = true
			decision.Priority = 5
		}
	}

	// Add notification if needed
	if limitResult.ShouldNotify {
		decision.Reason += fmt.Sprintf(" | %s", limitResult.NotifyMessage)
	}

	return decision, nil
}

// selectModel chooses the best model based on configuration and optimization
func (r *Router) selectModel(config *CostControlConfig, event Event) (string, string) {
	// If only one model is allowed, use it
	if len(config.Models) == 1 {
		return config.Models[0], fmt.Sprintf("Using configured model: %s", config.Models[0])
	}

	// Select based on optimization target
	switch config.OptimizeFor {
	case "cost":
		// Always use cheapest available model
		return r.getCheapestModel(config.Models), "Using cheapest model for cost optimization"

	case "speed":
		// Use fastest available model
		return r.getFastestModel(config.Models), "Using fastest model for speed optimization"

	case "accuracy":
		// Use best quality model for high confidence anomalies
		if event.AnomalyScore > 0.9 {
			return r.getBestQualityModel(config.Models), "Using best quality model for high confidence anomaly"
		}
		// Otherwise use balanced approach
		return r.getBalancedModel(config.Models), "Using balanced model for moderate confidence anomaly"

	default:
		return r.getCheapestModel(config.Models), "Using default cheapest model"
	}
}

// getCheapestModel returns the model with lowest cost
func (r *Router) getCheapestModel(allowedModels []string) string {
	if len(allowedModels) == 0 {
		return ""
	}

	cheapest := allowedModels[0]
	minCost := r.models[cheapest].CostPerInput + r.models[cheapest].CostPerOutput

	for _, model := range allowedModels[1:] {
		cost := r.models[model].CostPerInput + r.models[model].CostPerOutput
		if cost < minCost {
			cheapest = model
			minCost = cost
		}
	}

	return cheapest
}

// getFastestModel returns the model with highest speed rating
func (r *Router) getFastestModel(allowedModels []string) string {
	if len(allowedModels) == 0 {
		return ""
	}

	fastest := allowedModels[0]
	maxSpeed := r.models[fastest].Speed

	for _, model := range allowedModels[1:] {
		speed := r.models[model].Speed
		if speed > maxSpeed {
			fastest = model
			maxSpeed = speed
		}
	}

	return fastest
}

// getBestQualityModel returns the model with highest quality rating
func (r *Router) getBestQualityModel(allowedModels []string) string {
	if len(allowedModels) == 0 {
		return ""
	}

	best := allowedModels[0]
	maxQuality := r.models[best].Quality

	for _, model := range allowedModels[1:] {
		quality := r.models[model].Quality
		if quality > maxQuality {
			best = model
			maxQuality = quality
		}
	}

	return best
}

// getBalancedModel returns a model with good balance of speed, quality, and cost
func (r *Router) getBalancedModel(allowedModels []string) string {
	if len(allowedModels) == 0 {
		return ""
	}

	// Calculate balanced score (normalized)
	bestModel := allowedModels[0]
	bestScore := r.calculateBalancedScore(allowedModels[0])

	for _, model := range allowedModels[1:] {
		score := r.calculateBalancedScore(model)
		if score > bestScore {
			bestModel = model
			bestScore = score
		}
	}

	return bestModel
}

// calculateBalancedScore computes a balanced score for model selection
func (r *Router) calculateBalancedScore(model string) float64 {
	m := r.models[model]

	// Normalize factors (inverse for cost)
	costScore := 1.0 / (m.CostPerInput + m.CostPerOutput)
	speedScore := float64(m.Speed) / 10.0
	qualityScore := float64(m.Quality) / 10.0

	// Weighted average (favor balanced approach)
	return (costScore*0.4 + speedScore*0.3 + qualityScore*0.3)
}

// estimateCost estimates the cost of analyzing an event with a model
func (r *Router) estimateCost(model string, event Event, isBatch bool, plan string) float64 {
	// Typical analysis: 500 input tokens, 200 output tokens
	const (
		avgInputTokens  = 500
		avgOutputTokens = 200
		margin          = 1.15 // 15% margin
	)

	if m, ok := r.models[model]; ok {
		cost := (float64(avgInputTokens)/1000000.0*m.CostPerInput +
			float64(avgOutputTokens)/1000000.0*m.CostPerOutput)

		// Apply batch discount if applicable AND tenant is on enterprise plan
		if isBatch && plan == "orbit" {
			cost *= m.BatchDiscount
		}

		cost *= margin
		return math.Round(cost*1000000) / 1000000 // Round to 6 decimal places
	}

	return 0
}

// CalculateCost calculates the actual cost for a given model and token usage
func (r *Router) CalculateCost(model string, inputTokens, outputTokens int64) float64 {
	if m, ok := r.models[model]; ok {
		cost := (float64(inputTokens)/1000000.0*m.CostPerInput +
			float64(outputTokens)/1000000.0*m.CostPerOutput)
		return math.Round(cost*1000000) / 1000000 // Round to 6 decimal places
	}
	return 0
}

// OptimizeBatch creates optimized batches for processing
func (r *Router) OptimizeBatch(events []Event, config *CostControlConfig) [][]Event {
	if config.BatchSize <= 1 || len(events) <= 1 {
		// No batching
		batches := make([][]Event, len(events))
		for i, event := range events {
			batches[i] = []Event{event}
		}
		return batches
	}

	// Sort by priority (higher first)
	sorted := make([]Event, len(events))
	copy(sorted, events)

	// Group by priority and model
	batches := make([][]Event, 0)
	currentBatch := make([]Event, 0, config.BatchSize)

	for _, event := range sorted {
		decision, err := r.RouteEvent(context.Background(), event)
		if err != nil || !decision.ShouldAnalyze {
			continue
		}

		// Start new batch if model changes or batch is full
		if len(currentBatch) > 0 &&
			(len(currentBatch) >= config.BatchSize ||
				currentBatch[0].TenantID != event.TenantID) {
			batches = append(batches, currentBatch)
			currentBatch = make([]Event, 0, config.BatchSize)
		}

		currentBatch = append(currentBatch, event)
	}

	// Add last batch if not empty
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}

// UpdateUsageTracker updates usage statistics for optimization
func (r *Router) UpdateUsageTracker(tenantID, model string, cost float64, latency time.Duration, success bool) {
	key := fmt.Sprintf("%s:%s", tenantID, model)

	tracker, exists := r.trackers[key]
	if !exists {
		tracker = &UsageTracker{
			TenantID: tenantID,
			Model:    model,
		}
		r.trackers[key] = tracker
	}

	tracker.TotalCost += cost
	tracker.TotalCalls++

	// Update exponential moving average for latency
	if tracker.TotalCalls == 1 {
		tracker.AvgLatency = latency
	} else {
		alpha := 0.1 // Smoothing factor
		tracker.AvgLatency = time.Duration(alpha*float64(latency) + (1-alpha)*float64(tracker.AvgLatency))
	}

	// Update success rate
	if success {
		tracker.SuccessRate = (tracker.SuccessRate*float64(tracker.TotalCalls-1) + 1.0) / float64(tracker.TotalCalls)
	} else {
		tracker.SuccessRate = (tracker.SuccessRate * float64(tracker.TotalCalls-1)) / float64(tracker.TotalCalls)
	}
}
