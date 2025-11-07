package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/shannon-labs/driftlock/api-server/internal/models"
)

// MockExporter is a mock implementation of the exporter
type MockExporter struct {
	mock.Mock
}

func (m *MockExporter) ExportAnomaly(anomaly *models.Anomaly, exportedBy string) (*models.EvidenceBundle, error) {
	args := m.Called(anomaly, exportedBy)
	return args.Get(0).(*models.EvidenceBundle), args.Error(1)
}

func (m *MockExporter) ExportJSON(bundle *models.EvidenceBundle) ([]byte, error) {
	args := m.Called(bundle)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockExporter) VerifySignature(bundle *models.EvidenceBundle) (bool, error) {
	args := m.Called(bundle)
	return args.Bool(0), args.Error(1)
}

func (m *MockExporter) signBundle(bundle *models.EvidenceBundle) (string, error) {
	args := m.Called(bundle)
	return args.String(0), args.Error(1)
}

func TestExportHandler_ExportAnomaly(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Create test anomaly
	anomalyID := uuid.New()
	anomaly := &models.Anomaly{
		ID:                  anomalyID,
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		Status:              models.StatusPending,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Test explanation",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Create test evidence bundle
	bundle := &models.EvidenceBundle{
		Anomaly:    *anomaly,
		ExportedAt: time.Now(),
		ExportedBy: "test-user",
		Version:    "1.0.0",
		AdditionalMetadata: map[string]interface{}{
			"export_format": "driftlock-evidence-v1",
		},
	}

	// Create test JSON export
	jsonData := []byte(`{"test": "data"}`)

	// Set up mock expectations
	mockStorage.On("GetAnomaly", mock.Anything, anomalyID).Return(anomaly, nil)
	mockExporter.On("ExportAnomaly", anomaly, "test-user").Return(bundle, nil)
	mockExporter.On("ExportJSON", bundle).Return(jsonData, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/anomalies/%s/export", anomalyID), nil)
	w := httptest.NewRecorder()

	// Call handler
	// Note: We need to add the username to the context, but for testing purposes
	// we'll use a custom context
	ctx := context.WithValue(req.Context(), "username", "test-user")
	req = req.WithContext(ctx)

	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonData, w.Body.Bytes())
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf("attachment; filename=\"driftlock-evidence-%s.json\"", anomalyID), w.Header().Get("Content-Disposition"))

	// Verify mocks were called
	mockStorage.AssertExpectations(t)
	mockExporter.AssertExpectations(t)
}

func TestExportHandler_ExportAnomaly_InvalidID(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Create request with invalid UUID
	req := httptest.NewRequest(http.MethodGet, "/v1/anomalies/invalid-uuid/export", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid anomaly ID")

	// Verify mocks were not called
	mockStorage.AssertNotCalled(t, "GetAnomaly")
	mockExporter.AssertNotCalled(t, "ExportAnomaly")
}

func TestExportHandler_ExportAnomaly_AnomalyNotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Use a valid UUID but one that doesn't exist
	anomalyID := uuid.New()

	// Set up mock to return not found error
	mockStorage.On("GetAnomaly", mock.Anything, anomalyID).Return((*models.Anomaly)(nil), fmt.Errorf("anomaly not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/anomalies/%s/export", anomalyID), nil)
	w := httptest.NewRecorder()

	// Inject username context for exporter expectation
	req = req.WithContext(context.WithValue(req.Context(), "username", "test-user"))

	// Call handler
	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Anomaly not found")

	// Verify mocks were called
	mockStorage.AssertExpectations(t)
	mockExporter.AssertNotCalled(t, "ExportAnomaly")
}

func TestExportHandler_ExportAnomaly_ExporterError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Create test anomaly
	anomalyID := uuid.New()
	anomaly := &models.Anomaly{
		ID:         anomalyID,
		Timestamp:  time.Now(),
		StreamType: models.StreamTypeLogs,
		Status:     models.StatusPending,
		NCDScore:   0.25,
		PValue:     0.01,
	}

	// Set up mock expectations
	mockStorage.On("GetAnomaly", mock.Anything, anomalyID).Return(anomaly, nil)
	mockExporter.On("ExportAnomaly", anomaly, "test-user").Return((*models.EvidenceBundle)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/anomalies/%s/export", anomalyID), nil)
	w := httptest.NewRecorder()

	// Inject username context for exporter expectation
	req = req.WithContext(context.WithValue(req.Context(), "username", "test-user"))

	// Call handler
	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create bundle")

	// Verify mocks were called
	mockStorage.AssertExpectations(t)
	mockExporter.AssertExpectations(t)
}

func TestExportHandler_ExportAnomaly_ExportJSONError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Create test anomaly
	anomalyID := uuid.New()
	anomaly := &models.Anomaly{
		ID:         anomalyID,
		Timestamp:  time.Now(),
		StreamType: models.StreamTypeLogs,
		Status:     models.StatusPending,
		NCDScore:   0.25,
		PValue:     0.01,
	}

	// Create test evidence bundle
	bundle := &models.EvidenceBundle{
		Anomaly:    *anomaly,
		ExportedAt: time.Now(),
		ExportedBy: "test-user",
		Version:    "1.0.0",
	}

	// Set up mock expectations
	mockStorage.On("GetAnomaly", mock.Anything, anomalyID).Return(anomaly, nil)
	mockExporter.On("ExportAnomaly", anomaly, "test-user").Return(bundle, nil)
	mockExporter.On("ExportJSON", bundle).Return([]byte(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/anomalies/%s/export", anomalyID), nil)
	req = req.WithContext(context.WithValue(req.Context(), "username", "test-user"))
	w := httptest.NewRecorder()

	// Call handler
	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to export bundle")

	// Verify mocks were called
	mockStorage.AssertExpectations(t)
	mockExporter.AssertExpectations(t)
}

func TestExportHandler_ExportAnomaly_InvalidURL(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)
	handler := NewExportHandler(mockStorage, mockExporter)

	// Create request with invalid URL
	req := httptest.NewRequest(http.MethodGet, "/v1/anomalies/abc/export", nil) // Invalid UUID format
	w := httptest.NewRecorder()

	// Call handler
	handler.ExportAnomaly(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid anomaly ID")

	// Verify mocks were not called
	mockStorage.AssertNotCalled(t, "GetAnomaly")
	mockExporter.AssertNotCalled(t, "ExportAnomaly")
}

func TestNewExportHandler(t *testing.T) {
	mockStorage := new(MockStorage)
	mockExporter := new(MockExporter)

	handler := NewExportHandler(mockStorage, mockExporter)

	assert.NotNil(t, handler)
	assert.Equal(t, mockStorage, handler.storage)
	assert.Equal(t, mockExporter, handler.exporter)
}
