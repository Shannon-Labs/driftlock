package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/your-org/driftlock/api-server/internal/export"
	"github.com/your-org/driftlock/api-server/internal/storage"
)

// ExportHandler handles evidence bundle exports
type ExportHandler struct {
	storage  *storage.Storage
	exporter *export.Exporter
}

// NewExportHandler creates a new export handler
func NewExportHandler(storage *storage.Storage, exporter *export.Exporter) *ExportHandler {
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
	if len(path) < len(pathPrefix)+len(pathSuffix)+36 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := path[len(pathPrefix) : len(path)-len(pathSuffix)]

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
