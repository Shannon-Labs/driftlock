package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Hmbown/driftlock/api-server/internal/errors"
	"github.com/Hmbown/driftlock/api-server/internal/models"
)

const (
	// MaxRequestBodySize is the maximum allowed request body size (1 MB)
	MaxRequestBodySize = 1 << 20
)

// Validator provides request validation functionality
type Validator struct {
	maxBodySize int64
}

// New creates a new validator with default settings
func New() *Validator {
	return &Validator{
		maxBodySize: MaxRequestBodySize,
	}
}

// NewWithMaxSize creates a validator with a custom max body size
func NewWithMaxSize(maxSize int64) *Validator {
	return &Validator{
		maxBodySize: maxSize,
	}
}

// DecodeAndValidate decodes JSON from request body and validates it
func (v *Validator) DecodeAndValidate(r *http.Request, dest interface{}) error {
	// Limit request body size
	r.Body = http.MaxBytesReader(nil, r.Body, v.maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		if err == io.EOF {
			return errors.ErrBadRequest.WithDetails("Request body cannot be empty")
		}
		return errors.ErrInvalidJSON.WithDetails(err.Error())
	}

	// Check for additional JSON content
	if decoder.More() {
		return errors.ErrInvalidJSON.WithDetails("Request body must contain only a single JSON value")
	}

	return nil
}

// ValidateAnomalyCreate validates an anomaly creation request
func ValidateAnomalyCreate(create *models.AnomalyCreate) error {
	fieldErrors := make(map[string]string)

	// Validate timestamp
	if create.Timestamp.IsZero() {
		fieldErrors["timestamp"] = "timestamp is required"
	} else if create.Timestamp.After(time.Now().Add(5 * time.Minute)) {
		fieldErrors["timestamp"] = "timestamp cannot be in the future"
	}

	// Validate stream type
	if create.StreamType == "" {
		fieldErrors["stream_type"] = "stream_type is required"
	} else {
		validTypes := map[models.StreamType]bool{
			models.StreamTypeLogs:    true,
			models.StreamTypeMetrics: true,
			models.StreamTypeTraces:  true,
			models.StreamTypeLLM:     true,
		}
		if !validTypes[create.StreamType] {
			fieldErrors["stream_type"] = "stream_type must be one of: logs, metrics, traces, llm"
		}
	}

	// Validate NCD score
	if create.NCDScore < 0 || create.NCDScore > 1 {
		fieldErrors["ncd_score"] = "ncd_score must be between 0 and 1"
	}

	// Validate p-value
	if create.PValue < 0 || create.PValue > 1 {
		fieldErrors["p_value"] = "p_value must be between 0 and 1"
	}

	// Validate confidence level
	if create.ConfidenceLevel < 0 || create.ConfidenceLevel > 1 {
		fieldErrors["confidence_level"] = "confidence_level must be between 0 and 1"
	}

	// Validate explanation
	if create.GlassBoxExplanation == "" {
		fieldErrors["glass_box_explanation"] = "glass_box_explanation is required"
	} else if len(create.GlassBoxExplanation) > 5000 {
		fieldErrors["glass_box_explanation"] = "glass_box_explanation must not exceed 5000 characters"
	}

	// Validate compression metrics
	if create.CompressionBaseline < 0 {
		fieldErrors["compression_baseline"] = "compression_baseline must be non-negative"
	}
	if create.CompressionWindow < 0 {
		fieldErrors["compression_window"] = "compression_window must be non-negative"
	}
	if create.CompressionCombined < 0 {
		fieldErrors["compression_combined"] = "compression_combined must be non-negative"
	}

	// Validate tags
	if len(create.Tags) > 50 {
		fieldErrors["tags"] = "cannot have more than 50 tags"
	}
	for i, tag := range create.Tags {
		if len(tag) > 100 {
			fieldErrors[fmt.Sprintf("tags[%d]", i)] = "tag must not exceed 100 characters"
		}
	}

	if len(fieldErrors) > 0 {
		return errors.ValidationError(fieldErrors)
	}

	return nil
}

// ValidateAnomalyUpdate validates an anomaly update request
func ValidateAnomalyUpdate(update *models.AnomalyUpdate) error {
	fieldErrors := make(map[string]string)

	// Validate status
	validStatuses := map[models.AnomalyStatus]bool{
		models.StatusPending:       true,
		models.StatusAcknowledged:  true,
		models.StatusDismissed:     true,
		models.StatusInvestigating: true,
	}
	if !validStatuses[update.Status] {
		fieldErrors["status"] = "status must be one of: pending, acknowledged, dismissed, investigating"
	}

	// Validate notes length
	if update.Notes != nil && len(*update.Notes) > 10000 {
		fieldErrors["notes"] = "notes must not exceed 10000 characters"
	}

	if len(fieldErrors) > 0 {
		return errors.ValidationError(fieldErrors)
	}

	return nil
}

// ValidateAnomalyFilter validates anomaly filter parameters
func ValidateAnomalyFilter(filter *models.AnomalyFilter) error {
	fieldErrors := make(map[string]string)

	// Validate pagination
	if filter.Limit < 0 {
		fieldErrors["limit"] = "limit must be non-negative"
	}
	if filter.Limit > 1000 {
		fieldErrors["limit"] = "limit must not exceed 1000"
	}
	if filter.Offset < 0 {
		fieldErrors["offset"] = "offset must be non-negative"
	}

	// Validate score ranges
	if filter.MinNCDScore != nil && (*filter.MinNCDScore < 0 || *filter.MinNCDScore > 1) {
		fieldErrors["min_ncd_score"] = "min_ncd_score must be between 0 and 1"
	}
	if filter.MaxPValue != nil && (*filter.MaxPValue < 0 || *filter.MaxPValue > 1) {
		fieldErrors["max_p_value"] = "max_p_value must be between 0 and 1"
	}

	// Validate time range
	if filter.StartTime != nil && filter.EndTime != nil {
		if filter.StartTime.After(*filter.EndTime) {
			fieldErrors["start_time"] = "start_time must be before end_time"
		}
	}

	if len(fieldErrors) > 0 {
		return errors.ValidationError(fieldErrors)
	}

	return nil
}

// ValidateContentType validates the request Content-Type header
func ValidateContentType(r *http.Request, expected string) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != expected {
		return errors.ErrBadRequest.WithDetails(
			fmt.Sprintf("Content-Type must be %s, got %s", expected, contentType),
		)
	}
	return nil
}
