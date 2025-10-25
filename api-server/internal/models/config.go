package models

import "time"

// DetectionConfig represents the global CBAD detection configuration
type DetectionConfig struct {
	ID              int                    `json:"id" db:"id"`
	NCDThreshold    float64                `json:"ncd_threshold" db:"ncd_threshold"`
	PValueThreshold float64                `json:"p_value_threshold" db:"p_value_threshold"`
	BaselineSize    int                    `json:"baseline_size" db:"baseline_size"`
	WindowSize      int                    `json:"window_size" db:"window_size"`
	HopSize         int                    `json:"hop_size" db:"hop_size"`
	StreamOverrides map[string]interface{} `json:"stream_overrides,omitempty" db:"stream_overrides"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy       *string                `json:"created_by,omitempty" db:"created_by"`
	Notes           *string                `json:"notes,omitempty" db:"notes"`
	IsActive        bool                   `json:"is_active" db:"is_active"`
}

// DetectionConfigUpdate represents an update to detection configuration
type DetectionConfigUpdate struct {
	NCDThreshold    *float64               `json:"ncd_threshold,omitempty" validate:"omitempty,gte=0,lte=1"`
	PValueThreshold *float64               `json:"p_value_threshold,omitempty" validate:"omitempty,gte=0,lte=1"`
	BaselineSize    *int                   `json:"baseline_size,omitempty" validate:"omitempty,gte=1"`
	WindowSize      *int                   `json:"window_size,omitempty" validate:"omitempty,gte=1"`
	HopSize         *int                   `json:"hop_size,omitempty" validate:"omitempty,gte=1"`
	StreamOverrides map[string]interface{} `json:"stream_overrides,omitempty"`
	Notes           *string                `json:"notes,omitempty"`
}

// PerformanceMetric represents a performance measurement
type PerformanceMetric struct {
	ID           int64                  `json:"id" db:"id"`
	Timestamp    time.Time              `json:"timestamp" db:"timestamp"`
	MetricType   string                 `json:"metric_type" db:"metric_type"`
	Endpoint     *string                `json:"endpoint,omitempty" db:"endpoint"`
	DurationMs   float64                `json:"duration_ms" db:"duration_ms"`
	Success      bool                   `json:"success" db:"success"`
	ErrorMessage *string                `json:"error_message,omitempty" db:"error_message"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// APIKey represents an API authentication key
type APIKey struct {
	ID                 string     `json:"id" db:"id"`
	KeyHash            string     `json:"-" db:"key_hash"` // Never expose the hash
	Name               string     `json:"name" db:"name"`
	Description        *string    `json:"description,omitempty" db:"description"`
	Role               string     `json:"role" db:"role"`
	Scopes             []string   `json:"scopes" db:"scopes"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	CreatedBy          *string    `json:"created_by,omitempty" db:"created_by"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	LastUsedAt         *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	RateLimitPerMinute int        `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
}

// APIKeyCreate represents a request to create a new API key
type APIKeyCreate struct {
	Name               string   `json:"name" validate:"required"`
	Description        *string  `json:"description,omitempty"`
	Role               string   `json:"role" validate:"required,oneof=admin analyst viewer"`
	Scopes             []string `json:"scopes,omitempty"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	RateLimitPerMinute int      `json:"rate_limit_per_minute,omitempty"`
}

// APIKeyResponse includes the actual key (only returned once at creation)
type APIKeyResponse struct {
	APIKey
	Key string `json:"key"` // Only included at creation time
}
