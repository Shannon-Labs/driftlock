package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shannon-labs/driftlock/api-server/internal/models"
	"github.com/shannon-labs/driftlock/api-server/internal/storage"
)

// ConfigHandler handles configuration-related HTTP requests
type ConfigHandler struct {
	storage storage.ConfigStorage
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler(storage storage.ConfigStorage) *ConfigHandler {
	return &ConfigHandler{storage: storage}
}

// GetConfig handles GET /v1/config
func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config, err := h.storage.GetActiveConfig(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateConfig handles PATCH /v1/config
func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var update models.DetectionConfigUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate thresholds
	if update.NCDThreshold != nil && (*update.NCDThreshold < 0 || *update.NCDThreshold > 1) {
		http.Error(w, "NCD threshold must be between 0 and 1", http.StatusBadRequest)
		return
	}

	if update.PValueThreshold != nil && (*update.PValueThreshold < 0 || *update.PValueThreshold > 1) {
		http.Error(w, "P-value threshold must be between 0 and 1", http.StatusBadRequest)
		return
	}

	config, err := h.storage.UpdateConfig(ctx, &update)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update config: %v", err), http.StatusInternalServerError)
		return
	}

	// TODO: Notify OTel Collector processor of config change for live reload

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
