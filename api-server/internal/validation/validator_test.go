package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/driftlock/api-server/internal/errors"
	"github.com/your-org/driftlock/api-server/internal/models"
)

func TestNew(t *testing.T) {
	v := New()
	assert.Equal(t, int64(MaxRequestBodySize), v.maxBodySize)
}

func TestNewWithMaxSize(t *testing.T) {
	v := NewWithMaxSize(2048)
	assert.Equal(t, int64(2048), v.maxBodySize)
}

func TestDecodeAndValidate(t *testing.T) {
	v := New()

	// Test valid JSON
	anomalyCreate := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	jsonData, err := json.Marshal(anomalyCreate)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err = v.DecodeAndValidate(req, &dest)
	assert.NoError(t, err)
	assert.Equal(t, anomalyCreate.Timestamp.Format(time.RFC3339), dest.Timestamp.Format(time.RFC3339))
	assert.Equal(t, anomalyCreate.StreamType, dest.StreamType)
	assert.Equal(t, anomalyCreate.NCDScore, dest.NCDScore)
}

func TestDecodeAndValidate_InvalidJSON(t *testing.T) {
	v := New()

	invalidJSON := []byte(`{"invalid": json}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err := v.DecodeAndValidate(req, &dest)
	assert.Error(t, err)
	
	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	assert.Contains(t, strings.ToLower(apiErr.Message), "invalid json")
}

func TestDecodeAndValidate_EmptyBody(t *testing.T) {
	v := New()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err := v.DecodeAndValidate(req, &dest)
	assert.Error(t, err)
	
	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	assert.Contains(t, strings.ToLower(apiErr.Message), "request body cannot be empty")
}

func TestDecodeAndValidate_ExtraJSONContent(t *testing.T) {
	v := New()

	// Valid JSON followed by more content
	jsonWithExtra := []byte(`{"timestamp":"2023-01-01T00:00:00Z","stream_type":"logs","ncd_score":0.25,"p_value":0.01,"glass_box_explanation":"test","compression_baseline":0.6,"compression_window":0.3,"compression_combined":0.45,"confidence_level":0.95}{"extra":"content"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonWithExtra))
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err := v.DecodeAndValidate(req, &dest)
	assert.Error(t, err)
	
	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	assert.Contains(t, strings.ToLower(apiErr.Message), "request body must contain only a single json value")
}

func TestDecodeAndValidate_RequestTooLarge(t *testing.T) {
	v := NewWithMaxSize(100) // 100 bytes max

	// Create a large payload that exceeds the limit
	largePayload := make([]byte, 200)
	for i := range largePayload {
		largePayload[i] = 'a'
	}
	
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(largePayload))
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err := v.DecodeAndValidate(req, &dest)
	assert.Error(t, err)
	
	apiErr, ok := err.(*errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusRequestEntityTooLarge, apiErr.HTTPStatus)
}

func TestValidateAnomalyCreate(t *testing.T) {
	// Valid anomaly create
	validCreate := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(validCreate)
	assert.NoError(t, err)
}

func TestValidateAnomalyCreate_InvalidTimestamp(t *testing.T) {
	create := &models.AnomalyCreate{
		Timestamp:           time.Time{}, // Zero time
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "timestamp")
	assert.Contains(t, validationErr["timestamp"], "required")
}

func TestValidateAnomalyCreate_FutureTimestamp(t *testing.T) {
	create := &models.AnomalyCreate{
		Timestamp:           time.Now().Add(10 * time.Minute), // 10 minutes in future
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "timestamp")
	assert.Contains(t, validationErr["timestamp"], "future")
}

func TestValidateAnomalyCreate_InvalidStreamType(t *testing.T) {
	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          "invalid_stream_type",
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "stream_type")
	assert.Contains(t, validationErr["stream_type"], "logs, metrics, traces, llm")
}

func TestValidateAnomalyCreate_InvalidNCDScore(t *testing.T) {
	// Test NCD score < 0
	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            -0.1, // Invalid
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "ncd_score")
	assert.Contains(t, validationErr["ncd_score"], "between 0 and 1")

	// Test NCD score > 1
	create.NCDScore = 1.1 // Invalid
	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)
}

func TestValidateAnomalyCreate_InvalidPValue(t *testing.T) {
	// Test p-value < 0
	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              -0.1, // Invalid
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	// Test p-value > 1
	create.PValue = 1.1 // Invalid
	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)
}

func TestValidateAnomalyCreate_InvalidExplanation(t *testing.T) {
	// Test empty explanation
	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "", // Empty
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	// Test explanation too long
	longExplanation := strings.Repeat("a", 6000) // Exceeds 5000 chars
	create.GlassBoxExplanation = longExplanation

	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "glass_box_explanation")
	assert.Contains(t, validationErr["glass_box_explanation"], "5000 characters")
}

func TestValidateAnomalyCreate_InvalidCompressionValues(t *testing.T) {
	// Test negative compression baseline
	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: -0.1, // Invalid
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	// Test negative compression window
	create.CompressionBaseline = 0.6
	create.CompressionWindow = -0.1 // Invalid
	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)

	// Test negative compression combined
	create.CompressionWindow = 0.3
	create.CompressionCombined = -0.1 // Invalid
	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)
}

func TestValidateAnomalyCreate_InvalidTags(t *testing.T) {
	// Test too many tags
	tags := make([]string, 60) // Exceeds 50 tags
	for i := range tags {
		tags[i] = "tag" + string(rune(i))
	}

	create := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
		Tags:                tags,
	}

	err := ValidateAnomalyCreate(create)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "tags")
	assert.Contains(t, validationErr["tags"], "50 tags")

	// Test tag too long
	create.Tags = []string{strings.Repeat("a", 150)} // Exceeds 100 chars
	err = ValidateAnomalyCreate(create)
	assert.Error(t, err)

	validationErr, ok = apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "tags[0]")
	assert.Contains(t, validationErr["tags[0]"], "100 characters")
}

func TestValidateAnomalyUpdate(t *testing.T) {
	// Valid update
	update := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
		Notes:  stringPtr("Test notes"),
	}

	err := ValidateAnomalyUpdate(update)
	assert.NoError(t, err)
}

func TestValidateAnomalyUpdate_InvalidStatus(t *testing.T) {
	update := &models.AnomalyUpdate{
		Status: "invalid_status", // Not in allowed values
		Notes:  stringPtr("Test notes"),
	}

	err := ValidateAnomalyUpdate(update)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "status")
	assert.Contains(t, validationErr["status"], "pending, acknowledged, dismissed, investigating")
}

func TestValidateAnomalyUpdate_InvalidNotes(t *testing.T) {
	longNotes := strings.Repeat("a", 12000) // Exceeds 10000 chars
	
	update := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
		Notes:  &longNotes,
	}

	err := ValidateAnomalyUpdate(update)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "notes")
	assert.Contains(t, validationErr["notes"], "10000 characters")
}

func TestValidateAnomalyFilter(t *testing.T) {
	// Valid filter
	filter := &models.AnomalyFilter{
		StreamType:  &[]models.StreamType{models.StreamTypeLogs}[0],
		Status:      &[]models.AnomalyStatus{models.StatusPending}[0],
		MinNCDScore: float64Ptr(0.1),
		MaxPValue:   float64Ptr(0.05),
		Limit:       50,
		Offset:      0,
	}

	err := ValidateAnomalyFilter(filter)
	assert.NoError(t, err)
}

func TestValidateAnomalyFilter_InvalidLimit(t *testing.T) {
	// Test negative limit
	filter := &models.AnomalyFilter{
		Limit: -1,
	}

	err := ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "limit")
	assert.Contains(t, validationErr["limit"], "non-negative")

	// Test limit too high
	filter.Limit = 1500 // Exceeds 1000
	err = ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	validationErr, ok = apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "limit")
	assert.Contains(t, validationErr["limit"], "not exceed 1000")
}

func TestValidateAnomalyFilter_InvalidOffset(t *testing.T) {
	filter := &models.AnomalyFilter{
		Offset: -1, // Negative offset
	}

	err := ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "offset")
	assert.Contains(t, validationErr["offset"], "non-negative")
}

func TestValidateAnomalyFilter_InvalidScoreRanges(t *testing.T) {
	// Test invalid min NCD score
	filter := &models.AnomalyFilter{
		MinNCDScore: float64Ptr(-0.1), // Invalid
	}

	err := ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	// Test invalid max p-value
	filter.MinNCDScore = float64Ptr(0.5)
	filter.MaxPValue = float64Ptr(1.5) // Invalid
	
	err = ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "max_p_value")
	assert.Contains(t, validationErr["max_p_value"], "between 0 and 1")
}

func TestValidateAnomalyFilter_InvalidTimeRange(t *testing.T) {
	now := time.Now()
	startTime := now.Add(1 * time.Hour) // Start time after end time
	endTime := now

	filter := &models.AnomalyFilter{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	err := ValidateAnomalyFilter(filter)
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	
	validationErr, ok := apiErr.Details.(map[string]string)
	assert.True(t, ok)
	assert.Contains(t, validationErr, "start_time")
	assert.Contains(t, validationErr["start_time"], "before end_time")
}

func TestValidateContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	err := ValidateContentType(req, "application/json")
	assert.NoError(t, err)
}

func TestValidateContentType_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "text/plain")

	err := ValidateContentType(req, "application/json")
	assert.Error(t, err)

	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	assert.Contains(t, strings.ToLower(apiErr.Message), "content-type must be application/json")
}

func TestDecodeAndValidate_EOFError(t *testing.T) {
	v := New()

	// Create a reader that returns io.EOF on read
	reader := io.Reader(bytes.NewReader([]byte{}))
	
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Set("Content-Type", "application/json")

	var dest models.AnomalyCreate
	err := v.DecodeAndValidate(req, &dest)
	assert.Error(t, err)
	
	apiErr, ok := err.(errors.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	assert.Contains(t, strings.ToLower(apiErr.Message), "request body cannot be empty")
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}