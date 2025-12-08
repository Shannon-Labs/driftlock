package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ========================================================================
// Anomaly Feedback Endpoint
// POST /v1/anomalies/{id}/feedback
// ========================================================================

type feedbackRequest struct {
	FeedbackType string `json:"feedback_type"` // "false_positive", "confirmed", "dismissed"
	Reason       string `json:"reason,omitempty"`
}

type feedbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func anomalyFeedbackHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// Parse anomaly ID from path
		// Path format: /v1/anomalies/{id}/feedback
		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/anomalies/"), "/")
		if len(pathParts) < 2 || pathParts[1] != "feedback" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid path"))
			return
		}
		anomalyIDStr := pathParts[0]
		anomalyID, err := uuid.Parse(anomalyIDStr)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid anomaly ID"))
			return
		}

		// Decode request
		var req feedbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Validate feedback type
		validTypes := map[string]bool{
			"false_positive": true,
			"confirmed":      true,
			"dismissed":      true,
		}
		if !validTypes[req.FeedbackType] {
			writeError(w, r, http.StatusBadRequest,
				fmt.Errorf("feedback_type must be: false_positive, confirmed, or dismissed"))
			return
		}

		// Get tenant context
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// Get anomaly details for metrics capture
		anomaly, err := store.getAnomalyByID(r.Context(), tc.Tenant.ID, anomalyID)
		if err != nil {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("anomaly not found"))
			return
		}

		// Record feedback
		var reason *string
		if req.Reason != "" {
			reason = &req.Reason
		}

		err = store.recordFeedback(r.Context(), FeedbackRecord{
			AnomalyID:             anomalyID,
			StreamID:              anomaly.StreamID,
			TenantID:              tc.Tenant.ID,
			FeedbackType:          req.FeedbackType,
			NCDAtDetection:        anomaly.NCD,
			PValueAtDetection:     anomaly.PValue,
			ConfidenceAtDetection: anomaly.Confidence,
			FeedbackReason:        reason,
			CreatedBy:             "api",
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Trigger async auto-tune check
		go func() {
			if err := store.applyAutoTune(context.Background(), anomaly.StreamID); err != nil {
				log.Printf("auto-tune failed for stream %s: %v", anomaly.StreamID, err)
			}
		}()

		writeJSON(w, r, http.StatusOK, feedbackResponse{
			Success: true,
			Message: "Feedback recorded",
		})
	}
}

// ========================================================================
// Stream Profile Endpoint
// GET/PATCH /v1/streams/{id}/profile
// ========================================================================

type profileRequest struct {
	Profile               string `json:"profile,omitempty"` // "sensitive", "balanced", "strict"
	AutoTuneEnabled       *bool  `json:"auto_tune_enabled,omitempty"`
	AdaptiveWindowEnabled *bool  `json:"adaptive_window_enabled,omitempty"`
}

type profileResponse struct {
	Profile               string                    `json:"profile"`
	AutoTuneEnabled       bool                      `json:"auto_tune_enabled"`
	AdaptiveWindowEnabled bool                      `json:"adaptive_window_enabled"`
	CurrentThresholds     map[string]interface{}    `json:"current_thresholds"`
	ProfileDescriptions   map[string]ProfileSummary `json:"profile_descriptions,omitempty"`
}

func streamProfileHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse stream ID from path
		// Path format: /v1/streams/{id}/profile
		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/streams/"), "/")
		if len(pathParts) < 2 || pathParts[1] != "profile" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid path"))
			return
		}
		streamIDStr := pathParts[0]
		streamID, err := uuid.Parse(streamIDStr)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid stream ID"))
			return
		}

		// Get tenant context
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Return current profile and settings
			settings, err := store.getStreamTuningSettings(r.Context(), streamID)
			if err != nil {
				writeError(w, r, http.StatusNotFound, err)
				return
			}

			// Verify tenant owns this stream
			if !store.verifyStreamOwnership(r.Context(), tc.Tenant.ID, streamID) {
				writeError(w, r, http.StatusForbidden, fmt.Errorf("access denied"))
				return
			}

			writeJSON(w, r, http.StatusOK, profileResponse{
				Profile:               settings.DetectionProfile,
				AutoTuneEnabled:       settings.AutoTuneEnabled,
				AdaptiveWindowEnabled: settings.AdaptiveWindowEnabled,
				CurrentThresholds: map[string]interface{}{
					"ncd_threshold":    settings.NCDThreshold,
					"pvalue_threshold": settings.PValueThreshold,
					"baseline_size":    settings.BaselineSize,
					"window_size":      settings.WindowSize,
				},
				ProfileDescriptions: GetProfileSummaries(),
			})

		case http.MethodPatch:
			var req profileRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, r, http.StatusBadRequest, err)
				return
			}

			// Validate profile
			if req.Profile != "" && !IsValidProfile(req.Profile) {
				writeError(w, r, http.StatusBadRequest,
					fmt.Errorf("profile must be: sensitive, balanced, strict, or custom"))
				return
			}

			// Verify tenant owns this stream
			if !store.verifyStreamOwnership(r.Context(), tc.Tenant.ID, streamID) {
				writeError(w, r, http.StatusForbidden, fmt.Errorf("access denied"))
				return
			}

			// Update settings
			err := store.updateStreamProfile(r.Context(), streamID, req.Profile, req.AutoTuneEnabled, req.AdaptiveWindowEnabled)
			if err != nil {
				writeError(w, r, http.StatusInternalServerError, err)
				return
			}

			writeJSON(w, r, http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "Profile updated",
			})

		default:
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		}
	}
}

// ========================================================================
// Stream Tuning History Endpoint
// GET /v1/streams/{id}/tuning
// ========================================================================

type tuningResponse struct {
	CurrentSettings StreamTuningSettings `json:"current_settings"`
	TuneHistory     []TuneHistoryRecord  `json:"tune_history"`
	FeedbackStats   *FeedbackStats       `json:"feedback_stats,omitempty"`
}

func streamTuningHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// Parse stream ID from path
		pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/streams/"), "/")
		if len(pathParts) < 2 || pathParts[1] != "tuning" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid path"))
			return
		}
		streamIDStr := pathParts[0]
		streamID, err := uuid.Parse(streamIDStr)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid stream ID"))
			return
		}

		// Get tenant context
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// Verify tenant owns this stream
		if !store.verifyStreamOwnership(r.Context(), tc.Tenant.ID, streamID) {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("access denied"))
			return
		}

		// Get current settings
		settings, err := store.getStreamTuningSettings(r.Context(), streamID)
		if err != nil {
			writeError(w, r, http.StatusNotFound, err)
			return
		}

		// Get tune history
		history, err := store.getTuneHistory(r.Context(), streamID, 20)
		if err != nil {
			// Non-fatal, continue without history
			history = []TuneHistoryRecord{}
		}

		// Get feedback stats (last 30 days)
		var stats *FeedbackStats
		if settings.AutoTuneEnabled {
			stats, _ = store.getFeedbackStats(r.Context(), streamID, time.Now().AddDate(0, 0, -30))
		}

		writeJSON(w, r, http.StatusOK, tuningResponse{
			CurrentSettings: *settings,
			TuneHistory:     history,
			FeedbackStats:   stats,
		})
	}
}

// ========================================================================
// Profiles Listing Endpoint (public)
// GET /v1/profiles
// ========================================================================

func profilesListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"profiles": GetProfileSummaries(),
		})
	}
}

// ========================================================================
// Helper functions for store operations
// ========================================================================

// AnomalyRecord holds minimal anomaly data for feedback
type AnomalyRecord struct {
	ID         uuid.UUID
	StreamID   uuid.UUID
	TenantID   uuid.UUID
	NCD        float64
	PValue     float64
	Confidence float64
}

// getAnomalyByID retrieves anomaly metrics for feedback recording
func (s *store) getAnomalyByID(ctx context.Context, tenantID, anomalyID uuid.UUID) (*AnomalyRecord, error) {
	rec := &AnomalyRecord{ID: anomalyID}
	err := s.pool.QueryRow(ctx, `
		SELECT stream_id, tenant_id,
		       COALESCE(ncd, 0), COALESCE(p_value, 0), COALESCE(confidence_level, 0)
		FROM anomalies
		WHERE id = $1 AND tenant_id = $2`,
		anomalyID, tenantID).Scan(&rec.StreamID, &rec.TenantID, &rec.NCD, &rec.PValue, &rec.Confidence)
	return rec, err
}

// verifyStreamOwnership checks if a tenant owns a stream
func (s *store) verifyStreamOwnership(ctx context.Context, tenantID, streamID uuid.UUID) bool {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM streams WHERE id = $1 AND tenant_id = $2)`,
		streamID, tenantID).Scan(&exists)
	return err == nil && exists
}

// updateStreamProfile updates the detection profile for a stream
func (s *store) updateStreamProfile(ctx context.Context, streamID uuid.UUID, profile string, autoTune, adaptiveWindow *bool) error {
	// Build dynamic update query
	updates := []string{}
	args := []interface{}{streamID}
	argIdx := 2

	if profile != "" {
		updates = append(updates, fmt.Sprintf("detection_profile = $%d", argIdx))
		args = append(args, profile)
		argIdx++
	}
	if autoTune != nil {
		updates = append(updates, fmt.Sprintf("auto_tune_enabled = $%d", argIdx))
		args = append(args, *autoTune)
		argIdx++
	}
	if adaptiveWindow != nil {
		updates = append(updates, fmt.Sprintf("adaptive_window_enabled = $%d", argIdx))
		args = append(args, *adaptiveWindow)
		argIdx++
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	query := fmt.Sprintf("UPDATE streams SET %s WHERE id = $1", strings.Join(updates, ", "))
	_, err := s.pool.Exec(ctx, query, args...)
	return err
}
