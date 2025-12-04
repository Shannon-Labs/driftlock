package main

// SHA-140: Anchor Baseline Handlers for Drift Detection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

// anchorResponse is the API response for anchor endpoints
type anchorResponse struct {
	ID                       string   `json:"id"`
	StreamID                 string   `json:"stream_id"`
	Compressor               string   `json:"compressor"`
	EventCount               int      `json:"event_count"`
	CalibrationCompletedAt   string   `json:"calibration_completed_at"`
	IsActive                 bool     `json:"is_active"`
	BaselineEntropy          *float64 `json:"baseline_entropy,omitempty"`
	BaselineCompressionRatio *float64 `json:"baseline_compression_ratio,omitempty"`
	DriftNCDThreshold        float64  `json:"drift_ncd_threshold"`
	CreatedAt                string   `json:"created_at"`
}

// anchorSettingsResponse is the API response for anchor settings
type anchorSettingsResponse struct {
	AnchorEnabled      bool    `json:"anchor_enabled"`
	DriftNCDThreshold  float64 `json:"drift_ncd_threshold"`
	AnchorResetOnDrift bool    `json:"anchor_reset_on_drift"`
	HasActiveAnchor    bool    `json:"has_active_anchor"`
	AnchorID           *string `json:"anchor_id,omitempty"`
}

// getAnchorHandler handles GET /v1/streams/{id}/anchor
// Returns the current active anchor for a stream, or null if none exists
func getAnchorHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		streamID := r.PathValue("id")
		if streamID == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("stream ID required"))
			return
		}

		// Verify stream belongs to tenant
		stream, streamCfg, ok := store.streamBySlugOrID(tc.Tenant.ID, streamID)
		if !ok {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("stream not found"))
			return
		}

		// Get anchor settings
		anchorSettings, err := store.getAnchorSettings(r.Context(), stream.ID)
		if err != nil && err != errNotFound {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Get active anchor
		anchor, err := store.getActiveAnchor(r.Context(), stream.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		response := anchorSettingsResponse{
			AnchorEnabled:      anchorSettings != nil && anchorSettings.AnchorEnabled,
			DriftNCDThreshold:  streamCfg.NCDThreshold,
			AnchorResetOnDrift: false,
			HasActiveAnchor:    anchor != nil,
		}
		if anchorSettings != nil {
			response.DriftNCDThreshold = anchorSettings.DriftNCDThreshold
			response.AnchorResetOnDrift = anchorSettings.AnchorResetOnDrift
		}
		if anchor != nil {
			id := anchor.ID.String()
			response.AnchorID = &id
		}

		writeJSON(w, r, http.StatusOK, response)
	}
}

// getAnchorDetailsHandler handles GET /v1/streams/{id}/anchor/details
// Returns full anchor data including metadata
func getAnchorDetailsHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		streamID := r.PathValue("id")
		if streamID == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("stream ID required"))
			return
		}

		// Verify stream belongs to tenant
		stream, _, ok := store.streamBySlugOrID(tc.Tenant.ID, streamID)
		if !ok {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("stream not found"))
			return
		}

		// Get active anchor
		anchor, err := store.getActiveAnchor(r.Context(), stream.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		if anchor == nil {
			writeJSON(w, r, http.StatusOK, map[string]interface{}{
				"anchor":  nil,
				"message": "No active anchor. Stream may still be calibrating.",
			})
			return
		}

		response := anchorResponse{
			ID:                       anchor.ID.String(),
			StreamID:                 anchor.StreamID.String(),
			Compressor:               anchor.Compressor,
			EventCount:               anchor.EventCount,
			CalibrationCompletedAt:   anchor.CalibrationCompletedAt.Format(time.RFC3339),
			IsActive:                 anchor.IsActive,
			BaselineEntropy:          anchor.BaselineEntropy,
			BaselineCompressionRatio: anchor.BaselineCompressionRatio,
			DriftNCDThreshold:        anchor.DriftNCDThreshold,
			CreatedAt:                anchor.CreatedAt.Format(time.RFC3339),
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"anchor": response,
		})
	}
}

// resetAnchorRequest is the request body for POST /v1/streams/{id}/reset-anchor
type resetAnchorRequest struct {
	// ForceReset allows resetting even if not calibrated (admin override)
	ForceReset bool              `json:"force_reset,omitempty"`
	Events     []json.RawMessage `json:"events,omitempty"`
}

// resetAnchorHandler handles POST /v1/streams/{id}/reset-anchor
// Creates a new anchor from the current baseline window, deactivating the old one
func resetAnchorHandler(store *store, cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		streamID := r.PathValue("id")
		if streamID == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("stream ID required"))
			return
		}

		// Parse optional request body
		var req resetAnchorRequest
		if r.Body != nil && r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
				return
			}
		}

		// Verify stream belongs to tenant
		stream, settings, ok := store.streamBySlugOrID(tc.Tenant.ID, streamID)
		if !ok {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("stream not found"))
			return
		}

		// Check calibration status
		calibStatus, err := store.getStreamCalibrationStatus(r.Context(), stream.ID)
		if err != nil {
			log.Printf("Warning: failed to get calibration status: %v", err)
			// Continue anyway - might be legacy stream without calibration columns
		}

		if calibStatus != nil && !calibStatus.IsCalibrated && !req.ForceReset {
			writeJSON(w, r, http.StatusPreconditionFailed, map[string]interface{}{
				"success": false,
				"error":   "stream_not_calibrated",
				"message": fmt.Sprintf("Stream is still calibrating (%d/%d events). Cannot create anchor until calibration completes.",
					calibStatus.EventsIngested, calibStatus.MinBaselineSize),
				"calibration": map[string]interface{}{
					"events_ingested":   calibStatus.EventsIngested,
					"min_baseline_size": calibStatus.MinBaselineSize,
					"progress_percent":  calibStatus.ProgressPercent,
				},
			})
			return
		}

		// Get anchor settings
		anchorSettings, _ := store.getAnchorSettings(r.Context(), stream.ID)
		driftThreshold := 0.35
		if anchorSettings != nil {
			driftThreshold = anchorSettings.DriftNCDThreshold
		}

		if len(req.Events) == 0 {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("events required to create anchor"))
			return
		}

		plan := buildDetectionSettings(cfg, stream, settings, nil)
		usedAlgo := plan.CompressionAlgorithm
		if usedAlgo == "openzl" && !driftlockcbad.HasOpenZL() {
			usedAlgo = "zstd"
		}

		var tokenizer *driftlockcbad.Tokenizer
		if settings.TokenizerEnabled {
			tokenizer = driftlockcbad.GetTokenizer(driftlockcbad.TokenizerConfig{
				EnableUUID:   settings.TokenizeUUID,
				EnableHash:   settings.TokenizeHash,
				EnableBase64: settings.TokenizeBase64,
				EnableJWT:    settings.TokenizeJWT,
			})
		}

		snapshotLimit := minInt(len(req.Events), 200)
		snapshot, eventCount, err := buildSnapshot(req.Events, tokenizer, snapshotLimit, false)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid events for anchor: %w", err))
			return
		}
		if eventCount == 0 || len(bytes.TrimSpace(snapshot)) == 0 {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("events required to create anchor"))
			return
		}

		anchor := newAnchorRecord(stream.ID, usedAlgo, snapshot, eventCount, time.Now().UTC(), driftThreshold, plan.Seed, plan.PermutationCount)
		if anchor == nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("failed to build anchor snapshot"))
			return
		}

		if err := store.createStreamAnchor(r.Context(), anchor); err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create anchor: %w", err))
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"success":   true,
			"anchor_id": anchor.ID.String(),
			"message":   "Anchor created successfully. Drift detection will now compare against this baseline.",
		})
	}
}

// deleteAnchorHandler handles DELETE /v1/streams/{id}/anchor
// Deactivates the current anchor (disables drift detection until new anchor is created)
func deleteAnchorHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		streamID := r.PathValue("id")
		if streamID == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("stream ID required"))
			return
		}

		// Verify stream belongs to tenant
		stream, _, ok := store.streamBySlugOrID(tc.Tenant.ID, streamID)
		if !ok {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("stream not found"))
			return
		}

		// Deactivate anchor
		if err := store.deactivateAnchor(r.Context(), stream.ID); err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to deactivate anchor: %w", err))
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Anchor deactivated. Drift detection is now disabled until a new anchor is created.",
		})
	}
}
