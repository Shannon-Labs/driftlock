package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Hmbown/driftlock/api-server/internal/models"
)

func TestAnalyticsHandler_GetSummary(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Create test anomalies
	now := time.Now()
	anomalies := []models.Anomaly{
		{
			ID:                         uuid.New(),
			Timestamp:                  now,
			StreamType:                 models.StreamTypeLogs,
			Status:                     models.StatusPending,
			NCDScore:                   0.25,
			PValue:                     0.01,
			CompressionRatioChange:     50.0,
			IsStatisticallySignificant: true,
		},
		{
			ID:                         uuid.New(),
			Timestamp:                  now.Add(-1 * time.Hour),
			StreamType:                 models.StreamTypeMetrics,
			Status:                     models.StatusAcknowledged,
			NCDScore:                   0.30,
			PValue:                     0.02,
			CompressionRatioChange:     60.0,
			IsStatisticallySignificant: false,
		},
	}

	response := &models.AnomalyListResponse{
		Anomalies: anomalies,
		Total:     2,
		Limit:     10000,
		Offset:    0,
		HasMore:   false,
	}

	// Set up mock expectations
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return(response, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/summary", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetSummary(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var summary map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &summary)
	assert.NoError(t, err)

	// Verify summary fields
	assert.Equal(t, float64(2), summary["total_anomalies"])
	assert.Equal(t, float64(1), summary["significant_anomalies"])

	anomaliesByStreamType := summary["anomalies_by_stream_type"].(map[string]interface{})
	assert.Equal(t, float64(1), anomaliesByStreamType["logs"])
	assert.Equal(t, float64(1), anomaliesByStreamType["metrics"])

	anomaliesByStatus := summary["anomalies_by_status"].(map[string]interface{})
	assert.Equal(t, float64(1), anomaliesByStatus["pending"])
	assert.Equal(t, float64(1), anomaliesByStatus["acknowledged"])

	// Verify averages
	assert.Equal(t, 0.275, summary["average_ncd_score"])         // (0.25 + 0.30) / 2
	assert.Equal(t, 0.015, summary["average_p_value"])           // (0.01 + 0.02) / 2
	assert.Equal(t, 55.0, summary["average_compression_change"]) // (50.0 + 60.0) / 2

	// Verify time range
	timeRange := summary["time_range"].(map[string]interface{})
	assert.NotEmpty(t, timeRange["start"])
	assert.NotEmpty(t, timeRange["end"])

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestAnalyticsHandler_GetSummary_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Set up mock to return an error
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return((*models.AnomalyListResponse)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/summary", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetSummary(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get anomalies")

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestAnalyticsHandler_GetCompressionTimeline(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Create test anomalies
	now := time.Now()
	anomalies := []models.Anomaly{
		{
			ID:                     uuid.New(),
			Timestamp:              now,
			StreamType:             models.StreamTypeLogs,
			Status:                 models.StatusPending,
			CompressionBaseline:    0.6,
			CompressionWindow:      0.3,
			CompressionRatioChange: 50.0,
		},
		{
			ID:                     uuid.New(),
			Timestamp:              now.Add(-1 * time.Hour),
			StreamType:             models.StreamTypeMetrics,
			Status:                 models.StatusAcknowledged,
			CompressionBaseline:    0.7,
			CompressionWindow:      0.4,
			CompressionRatioChange: 42.9,
		},
	}

	response := &models.AnomalyListResponse{
		Anomalies: anomalies,
		Total:     2,
		Limit:     1000,
		Offset:    0,
		HasMore:   false,
	}

	// Set up mock expectations
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return(response, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/compression-timeline", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetCompressionTimeline(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Verify timeline structure
	assert.Equal(t, float64(2), result["total"])
	timeline := result["timeline"].([]interface{})
	assert.Len(t, timeline, 2)

	// Verify first point
	firstPoint := timeline[0].(map[string]interface{})
	assert.Equal(t, 0.6, firstPoint["baseline_ratio"])
	assert.Equal(t, 0.3, firstPoint["window_ratio"])
	assert.Equal(t, 50.0, firstPoint["compression_change"])

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestAnalyticsHandler_GetCompressionTimeline_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Set up mock to return an error
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return((*models.AnomalyListResponse)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/compression-timeline", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetCompressionTimeline(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get anomalies")

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestAnalyticsHandler_GetNCDHeatmap(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Create test anomalies
	now := time.Now()
	anomalies := []models.Anomaly{
		{
			ID:         uuid.New(),
			Timestamp:  now,
			StreamType: models.StreamTypeLogs,
			Status:     models.StatusPending,
			NCDScore:   0.25,
		},
		{
			ID:         uuid.New(),
			Timestamp:  now.Add(-1 * time.Hour),
			StreamType: models.StreamTypeMetrics,
			Status:     models.StatusAcknowledged,
			NCDScore:   0.30,
		},
	}

	response := &models.AnomalyListResponse{
		Anomalies: anomalies,
		Total:     2,
		Limit:     10000,
		Offset:    0,
		HasMore:   false,
	}

	// Set up mock expectations
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return(response, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/ncd-heatmap", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetNCDHeatmap(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Verify heatmap structure
	assert.NotNil(t, result["heatmap"])
	assert.NotNil(t, result["start_time"])
	assert.NotNil(t, result["end_time"])

	heatmap := result["heatmap"].([]interface{})
	assert.Greater(t, len(heatmap), 0)

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestAnalyticsHandler_GetNCDHeatmap_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewAnalyticsHandler(mockStorage)

	// Set up mock to return an error
	mockStorage.On("ListAnomalies", mock.Anything, mock.AnythingOfType("*models.AnomalyFilter")).Return((*models.AnomalyListResponse)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/analytics/ncd-heatmap", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetNCDHeatmap(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get anomalies")

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}
