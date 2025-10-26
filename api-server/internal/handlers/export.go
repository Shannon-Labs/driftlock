package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/storage"
)

// Exporter defines the interface for export operations
type Exporter interface {
	ExportAnomaly(anomaly *models.Anomaly, exportedBy string) (*models.EvidenceBundle, error)
	ExportJSON(bundle *models.EvidenceBundle) ([]byte, error)
	VerifySignature(bundle *models.EvidenceBundle) (bool, error)
}

// ExportHandler handles evidence bundle exports
type ExportHandler struct {
	storage  storage.AnomalyStorage
	exporter Exporter
}

// NewExportHandler creates a new export handler
func NewExportHandler(storage storage.AnomalyStorage, exporter Exporter) *ExportHandler {
	return &ExportHandler{
		storage:  storage,
		exporter: exporter,
	}
}

// ExportAnomaly handles GET /v1/anomalies/:id/export
func (h *ExportHandler) ExportAnomaly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL
	pathPrefix := "/v1/anomalies/"
	pathSuffix := "/export"
	path := r.URL.Path

	if !strings.HasPrefix(path, pathPrefix) || !strings.HasSuffix(path, pathSuffix) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSuffix(strings.TrimPrefix(path, pathPrefix), pathSuffix)
	idStr = strings.Trim(idStr, "/")

	if idStr == "" {
		http.Error(w, "Invalid anomaly ID", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid anomaly ID", http.StatusBadRequest)
		return
	}

	// Get anomaly
	anomaly, err := h.storage.GetAnomaly(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Anomaly not found: %v", err), http.StatusNotFound)
		return
	}

	// Get username from context
	username := getUsernameFromContext(ctx)

	// Generate evidence bundle
	bundle, err := h.exporter.ExportAnomaly(anomaly, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create bundle: %v", err), http.StatusInternalServerError)
		return
	}

	// Export as JSON
	data, err := h.exporter.ExportJSON(bundle)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export bundle: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"driftlock-evidence-%s.json\"", id))

	w.Write(data)
}
