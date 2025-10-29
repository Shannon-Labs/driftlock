package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/storage"
	"github.com/Hmbown/driftlock/api-server/internal/stream"
	"github.com/Hmbown/driftlock/api-server/internal/streaming"
	"github.com/Hmbown/driftlock/api-server/internal/supabase"
	"github.com/google/uuid"
)

// AnomaliesHandler handles anomaly-related HTTP requests
type AnomaliesHandler struct {
	storage  storage.AnomalyStorage
	streamer *stream.Streamer
	events   streaming.EventPublisher
	supabase *supabase.Client
}

// NewAnomaliesHandler creates a new anomalies handler
func NewAnomaliesHandler(storage storage.AnomalyStorage, streamer *stream.Streamer, events streaming.EventPublisher) *AnomaliesHandler {
	return &AnomaliesHandler{
		storage:  storage,
		streamer: streamer,
		events:   events,
	}
}

// NewAnomaliesHandlerWithSupabase creates a new anomalies handler with Supabase integration
func NewAnomaliesHandlerWithSupabase(storage storage.AnomalyStorage, streamer *stream.Streamer, events streaming.EventPublisher, supabaseClient *supabase.Client) *AnomaliesHandler {
	return &AnomaliesHandler{
		storage:  storage,
		streamer: streamer,
		events:   events,
		supabase: supabaseClient,
	}
}

// ListAnomalies handles GET /v1/anomalies
func (h *AnomaliesHandler) ListAnomalies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	filter := &models.AnomalyFilter{
		Limit:  parseIntParam(r, "limit", 50),
		Offset: parseIntParam(r, "offset", 0),
	}

	if streamType := r.URL.Query().Get("stream_type"); streamType != "" {
		st := models.StreamType(streamType)
		filter.StreamType = &st
	}

	if status := r.URL.Query().Get("status"); status != "" {
		s := models.AnomalyStatus(status)
		filter.Status = &s
	}

	if minNCD := r.URL.Query().Get("min_ncd_score"); minNCD != "" {
		if val, err := strconv.ParseFloat(minNCD, 64); err == nil {
			filter.MinNCDScore = &val
		}
	}

	if maxPValue := r.URL.Query().Get("max_p_value"); maxPValue != "" {
		if val, err := strconv.ParseFloat(maxPValue, 64); err == nil {
			filter.MaxPValue = &val
		}
	}

	if startTime := r.URL.Query().Get("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = &t
		}
	}

	if endTime := r.URL.Query().Get("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filter.EndTime = &t
		}
	}

	if r.URL.Query().Get("only_significant") == "true" {
		filter.OnlySignificant = true
	}

	// Query anomalies
	response, err := h.storage.ListAnomalies(ctx, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list anomalies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAnomaly handles GET /v1/anomalies/:id
func (h *AnomaliesHandler) GetAnomaly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL path
	idStr := r.URL.Path[len("/v1/anomalies/"):]
	if idx := len(idStr); idx > 36 {
		idStr = idStr[:36]
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid anomaly ID", http.StatusBadRequest)
		return
	}

	anomaly, err := h.storage.GetAnomaly(ctx, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Anomaly not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anomaly)
}

// UpdateAnomalyStatus handles PATCH /v1/anomalies/:id/status
func (h *AnomaliesHandler) UpdateAnomalyStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL
	pathPrefix := "/v1/anomalies/"
	pathSuffix := "/status"
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

	// Parse request body
	var update models.AnomalyUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get username from context (set by auth middleware)
	username := getUsernameFromContext(ctx)

	// Update anomaly
	if err := h.storage.UpdateAnomalyStatus(ctx, id, &update, username); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update anomaly: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch updated anomaly
	anomaly, err := h.storage.GetAnomaly(ctx, id)
	if err != nil {
		http.Error(w, "Failed to fetch updated anomaly", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anomaly)
}

// CreateAnomaly handles POST /v1/anomalies (internal use by CBAD processor)
func (h *AnomaliesHandler) CreateAnomaly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var create models.AnomalyCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create anomaly in database
	anomaly, err := h.storage.CreateAnomaly(ctx, &create)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create anomaly: %v", err), http.StatusInternalServerError)
		return
	}

	// Also create in Supabase if client is available
	if h.supabase != nil {
		supabaseAnomaly := map[string]interface{}{
			"id":          anomaly.ID.String(),
			"event_type":  string(create.StreamType),
			"severity":    anomaly.GetSeverity(),
			"status":      "open",
			"created_at":  time.Now().Format(time.RFC3339),
			"metadata":    create.Metadata,
			"ncd_score":   anomaly.NCDScore,
			"p_value":     anomaly.PValue,
		}

		if err := h.supabase.CreateAnomaly(ctx, supabaseAnomaly); err != nil {
			log.Printf("anomalies: failed to create anomaly in Supabase: %v", err)
		} else {
			log.Printf("anomalies: created anomaly in Supabase: %s", anomaly.ID.String())
		}

		// Best-effort usage metering (only when anomaly is detected)
		orgID := os.Getenv("TENANT_DEFAULT_TENANT")
		if orgID == "" {
			orgID = os.Getenv("SUPABASE_PROJECT_ID")
		}
		if orgID != "" {
			if err := h.supabase.MeterUsage(ctx, orgID, true, 1); err != nil {
				log.Printf("anomalies: failed to meter usage: %v", err)
			}
		}
	}

	// Broadcast to SSE clients
	if h.streamer != nil {
		h.streamer.BroadcastAnomaly(anomaly)
	}

	// Emit anomaly-created event
	if h.events != nil {
		if err := h.events.AnomalyCreated(ctx, anomaly); err != nil {
			log.Printf("anomalies: failed to publish anomaly-created event: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(anomaly)
}

// Helper functions

func parseIntParam(r *http.Request, param string, defaultValue int) int {
	if val := r.URL.Query().Get(param); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

func getUsernameFromContext(ctx context.Context) string {
	if ctx == nil {
		return "system"
	}

	if username, ok := ctx.Value("username").(string); ok && username != "" {
		return username
	}

	return "system"
}
