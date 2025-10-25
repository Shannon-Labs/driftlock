package models

import (
	"time"

	"github.com/google/uuid"
)

// StreamType represents the type of telemetry stream
type StreamType string

const (
	StreamTypeLogs    StreamType = "logs"
	StreamTypeMetrics StreamType = "metrics"
	StreamTypeTraces  StreamType = "traces"
	StreamTypeLLM     StreamType = "llm"
)

// AnomalyStatus represents the current status of an anomaly
type AnomalyStatus string

const (
	StatusPending       AnomalyStatus = "pending"
	StatusAcknowledged  AnomalyStatus = "acknowledged"
	StatusDismissed     AnomalyStatus = "dismissed"
	StatusInvestigating AnomalyStatus = "investigating"
)

// Anomaly represents a detected anomaly with full CBAD metrics
type Anomaly struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`

	// Stream identification
	StreamType StreamType `json:"stream_type" db:"stream_type"`

	// Core CBAD metrics
	NCDScore float64 `json:"ncd_score" db:"ncd_score"`
	PValue   float64 `json:"p_value" db:"p_value"`

	// Status
	Status AnomalyStatus `json:"status" db:"status"`

	// Explanations
	GlassBoxExplanation string  `json:"glass_box_explanation" db:"glass_box_explanation"`
	DetailedExplanation *string `json:"detailed_explanation,omitempty" db:"detailed_explanation"`

	// Compression metrics
	CompressionBaseline     float64 `json:"compression_baseline" db:"compression_baseline"`
	CompressionWindow       float64 `json:"compression_window" db:"compression_window"`
	CompressionCombined     float64 `json:"compression_combined" db:"compression_combined"`
	CompressionRatioChange  float64 `json:"compression_ratio_change" db:"compression_ratio_change"`

	// Entropy metrics
	BaselineEntropy *float64 `json:"baseline_entropy,omitempty" db:"baseline_entropy"`
	WindowEntropy   *float64 `json:"window_entropy,omitempty" db:"window_entropy"`
	EntropyChange   *float64 `json:"entropy_change,omitempty" db:"entropy_change"`

	// Statistical significance
	ConfidenceLevel            float64 `json:"confidence_level" db:"confidence_level"`
	IsStatisticallySignificant bool    `json:"is_statistically_significant" db:"is_statistically_significant"`

	// Data payloads (JSON)
	BaselineData map[string]interface{} `json:"baseline_data,omitempty" db:"baseline_data"`
	WindowData   map[string]interface{} `json:"window_data,omitempty" db:"window_data"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// Tags
	Tags []string `json:"tags,omitempty" db:"tags"`

	// User interaction
	AcknowledgedBy *string    `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	DismissedBy    *string    `json:"dismissed_by,omitempty" db:"dismissed_by"`
	DismissedAt    *time.Time `json:"dismissed_at,omitempty" db:"dismissed_at"`
	Notes          *string    `json:"notes,omitempty" db:"notes"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AnomalyCreate represents the payload for creating a new anomaly
type AnomalyCreate struct {
	Timestamp               time.Time              `json:"timestamp" validate:"required"`
	StreamType              StreamType             `json:"stream_type" validate:"required"`
	NCDScore                float64                `json:"ncd_score" validate:"required,gte=0,lte=1"`
	PValue                  float64                `json:"p_value" validate:"required,gte=0,lte=1"`
	GlassBoxExplanation     string                 `json:"glass_box_explanation" validate:"required"`
	CompressionBaseline     float64                `json:"compression_baseline" validate:"required"`
	CompressionWindow       float64                `json:"compression_window" validate:"required"`
	CompressionCombined     float64                `json:"compression_combined" validate:"required"`
	ConfidenceLevel         float64                `json:"confidence_level" validate:"required,gte=0,lte=1"`
	BaselineData            map[string]interface{} `json:"baseline_data,omitempty"`
	WindowData              map[string]interface{} `json:"window_data,omitempty"`
	Metadata                map[string]interface{} `json:"metadata,omitempty"`
	Tags                    []string               `json:"tags,omitempty"`
}

// AnomalyUpdate represents the payload for updating an anomaly
type AnomalyUpdate struct {
	Status AnomalyStatus `json:"status" validate:"required,oneof=pending acknowledged dismissed investigating"`
	Notes  *string       `json:"notes,omitempty"`
}

// AnomalyFilter represents filter parameters for querying anomalies
type AnomalyFilter struct {
	StreamType  *StreamType    `json:"stream_type,omitempty"`
	Status      *AnomalyStatus `json:"status,omitempty"`
	MinNCDScore *float64       `json:"min_ncd_score,omitempty"`
	MaxPValue   *float64       `json:"max_p_value,omitempty"`
	StartTime   *time.Time     `json:"start_time,omitempty"`
	EndTime     *time.Time     `json:"end_time,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	OnlySignificant bool       `json:"only_significant,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// AnomalyListResponse represents a paginated list of anomalies
type AnomalyListResponse struct {
	Anomalies  []Anomaly `json:"anomalies"`
	Total      int       `json:"total"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
	HasMore    bool      `json:"has_more"`
}

// EvidenceBundle represents an exportable evidence package
type EvidenceBundle struct {
	Anomaly            Anomaly                `json:"anomaly"`
	ExportedAt         time.Time              `json:"exported_at"`
	ExportedBy         string                 `json:"exported_by"`
	Version            string                 `json:"version"`
	Signature          *string                `json:"signature,omitempty"`
	AdditionalMetadata map[string]interface{} `json:"additional_metadata,omitempty"`
}

// IsAnomaly returns true if the metrics indicate an anomaly
func (a *Anomaly) IsAnomaly() bool {
	return a.IsStatisticallySignificant && a.PValue < 0.05
}

// GetSeverity returns a severity level based on metrics
func (a *Anomaly) GetSeverity() string {
	if a.PValue < 0.001 && a.NCDScore > 0.5 {
		return "critical"
	}
	if a.PValue < 0.01 && a.NCDScore > 0.3 {
		return "high"
	}
	if a.PValue < 0.05 {
		return "medium"
	}
	return "low"
}
