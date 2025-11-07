package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
)

func TestConfigHandler_GetConfig(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create test config
	config := &models.DetectionConfig{
		ID:              1,
		NCDThreshold:    0.2,
		PValueThreshold: 0.05,
		BaselineSize:    100,
		WindowSize:      100,
		HopSize:         50,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
	}

	// Set up mock expectations
	mockStorage.On("GetActiveConfig", mock.Anything).Return(config, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/config", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseConfig map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseConfig)
	assert.NoError(t, err)

	// Verify config fields
	assert.Equal(t, float64(1), responseConfig["id"])
	assert.Equal(t, 0.2, responseConfig["ncd_threshold"])
	assert.Equal(t, 0.05, responseConfig["p_value_threshold"])
	assert.Equal(t, float64(100), responseConfig["baseline_size"])
	assert.Equal(t, float64(100), responseConfig["window_size"])
	assert.Equal(t, float64(50), responseConfig["hop_size"])
	assert.Equal(t, true, responseConfig["is_active"])

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestConfigHandler_GetConfig_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Set up mock to return an error
	mockStorage.On("GetActiveConfig", mock.Anything).Return((*models.DetectionConfig)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/v1/config", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.GetConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get config")

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestConfigHandler_UpdateConfig(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create test config update
	updateData := models.DetectionConfigUpdate{
		NCDThreshold:    float64Ptr(0.3),
		PValueThreshold: float64Ptr(0.01),
	}

	updateJSON, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// Create updated config to return
	updatedConfig := &models.DetectionConfig{
		ID:              1,
		NCDThreshold:    0.3,
		PValueThreshold: 0.01,
		BaselineSize:    100,
		WindowSize:      100,
		HopSize:         50,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
	}

	// Set up mock expectations
	mockStorage.On("UpdateConfig", mock.Anything, mock.AnythingOfType("*models.DetectionConfigUpdate")).Return(updatedConfig, nil)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/v1/config", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.UpdateConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseConfig map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseConfig)
	assert.NoError(t, err)

	// Verify updated config fields
	assert.Equal(t, 0.3, responseConfig["ncd_threshold"])
	assert.Equal(t, 0.01, responseConfig["p_value_threshold"])

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func TestConfigHandler_UpdateConfig_InvalidJSON(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create invalid JSON
	invalidJSON := []byte(`{"invalid": json}`)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/v1/config", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.UpdateConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")

	// Verify mock was not called
	mockStorage.AssertNotCalled(t, "UpdateConfig")
}

func TestConfigHandler_UpdateConfig_InvalidNCDThreshold(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create update with invalid NCD threshold
	updateData := models.DetectionConfigUpdate{
		NCDThreshold: float64Ptr(-0.1), // Invalid
	}

	updateJSON, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/v1/config", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.UpdateConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "NCD threshold must be between 0 and 1")

	// Verify mock was not called
	mockStorage.AssertNotCalled(t, "UpdateConfig")
}

func TestConfigHandler_UpdateConfig_InvalidPValueThreshold(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create update with invalid p-value threshold
	updateData := models.DetectionConfigUpdate{
		PValueThreshold: float64Ptr(1.5), // Invalid
	}

	updateJSON, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/v1/config", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.UpdateConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "P-value threshold must be between 0 and 1")

	// Verify mock was not called
	mockStorage.AssertNotCalled(t, "UpdateConfig")
}

func TestConfigHandler_UpdateConfig_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := NewConfigHandler(mockStorage)

	// Create test config update
	updateData := models.DetectionConfigUpdate{
		NCDThreshold: float64Ptr(0.25),
	}

	updateJSON, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// Set up mock to return an error
	mockStorage.On("UpdateConfig", mock.Anything, mock.AnythingOfType("*models.DetectionConfigUpdate")).Return((*models.DetectionConfig)(nil), assert.AnError)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/v1/config", bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.UpdateConfig(w, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to update config")

	// Verify mock was called
	mockStorage.AssertExpectations(t)
}

func float64Ptr(f float64) *float64 {
	return &f
}