package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/driftlock/api-server/internal/models"
	"github.com/your-org/driftlock/api-server/internal/storage"
)

// AnalyticsHandler handles analytics and statistics endpoints
type AnalyticsHandler struct {
	storage storage.AnomalyStorage
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(storage storage.AnomalyStorage) *AnalyticsHandler {
	return &AnalyticsHandler{storage: storage}
}

// Summary represents statistical summary of anomalies
type Summary struct {
	TotalAnomalies           int                        `json:"total_anomalies"`
	AnomaliesByStreamType    map[string]int             `json:"anomalies_by_stream_type"`
	AnomaliesByStatus        map[string]int             `json:"anomalies_by_status"`
	SignificantAnomalies     int                        `json:"significant_anomalies"`
	AverageNCDScore          float64                    `json:"average_ncd_score"`
	AveragePValue            float64                    `json:"average_p_value"`
	AverageCompressionChange float64                    `json:"average_compression_change"`
	TimeRange                map[string]time.Time       `json:"time_range"`
}

// GetSummary handles GET /v1/analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse time range from query params
	var startTime, endTime *time.Time
	if st := r.URL.Query().Get("start_time"); st != "" {
		if t, err := time.Parse(time.RFC3339, st); err == nil {
			startTime = &t
		}
	}
	if et := r.URL.Query().Get("end_time"); et != "" {
		if t, err := time.Parse(time.RFC3339, et); err == nil {
			endTime = &t
		}
	}

	// Get all anomalies in time range
	filter := &models.AnomalyFilter{
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     10000, // High limit for analytics
	}

	response, err := h.storage.ListAnomalies(ctx, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get anomalies: %v", err), http.StatusInternalServerError)
		return
	}

	// Compute summary statistics
	summary := Summary{
		TotalAnomalies:        len(response.Anomalies),
		AnomaliesByStreamType: make(map[string]int),
		AnomaliesByStatus:     make(map[string]int),
		TimeRange:             make(map[string]time.Time),
	}

	var totalNCD, totalPValue, totalCompChange float64
	var minTime, maxTime time.Time

	for i, a := range response.Anomalies {
		// Count by stream type
		summary.AnomaliesByStreamType[string(a.StreamType)]++

		// Count by status
		summary.AnomaliesByStatus[string(a.Status)]++

		// Count significant anomalies
		if a.IsStatisticallySignificant {
			summary.SignificantAnomalies++
		}

		// Sum for averages
		totalNCD += a.NCDScore
		totalPValue += a.PValue
		totalCompChange += a.CompressionRatioChange

		// Track time range
		if i == 0 || a.Timestamp.Before(minTime) {
			minTime = a.Timestamp
		}
		if i == 0 || a.Timestamp.After(maxTime) {
			maxTime = a.Timestamp
		}
	}

	// Calculate averages
	if len(response.Anomalies) > 0 {
		summary.AverageNCDScore = totalNCD / float64(len(response.Anomalies))
		summary.AveragePValue = totalPValue / float64(len(response.Anomalies))
		summary.AverageCompressionChange = totalCompChange / float64(len(response.Anomalies))
		summary.TimeRange["start"] = minTime
		summary.TimeRange["end"] = maxTime
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// CompressionTimePoint represents compression ratio at a point in time
type CompressionTimePoint struct {
	Timestamp          time.Time `json:"timestamp"`
	BaselineRatio      float64   `json:"baseline_ratio"`
	WindowRatio        float64   `json:"window_ratio"`
	CompressionChange  float64   `json:"compression_change"`
}

// GetCompressionTimeline handles GET /v1/analytics/compression-timeline
func (h *AnalyticsHandler) GetCompressionTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse parameters
	var startTime, endTime *time.Time
	if st := r.URL.Query().Get("start_time"); st != "" {
		if t, err := time.Parse(time.RFC3339, st); err == nil {
			startTime = &t
		}
	}
	if et := r.URL.Query().Get("end_time"); et != "" {
		if t, err := time.Parse(time.RFC3339, et); err == nil {
			endTime = &t
		}
	}

	var streamType *models.StreamType
	if st := r.URL.Query().Get("stream_type"); st != "" {
		s := models.StreamType(st)
		streamType = &s
	}

	filter := &models.AnomalyFilter{
		StreamType: streamType,
		StartTime:  startTime,
		EndTime:    endTime,
		Limit:      1000,
	}

	response, err := h.storage.ListAnomalies(ctx, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get anomalies: %v", err), http.StatusInternalServerError)
		return
	}

	// Build timeline
	timeline := make([]CompressionTimePoint, len(response.Anomalies))
	for i, a := range response.Anomalies {
		timeline[i] = CompressionTimePoint{
			Timestamp:         a.Timestamp,
			BaselineRatio:     a.CompressionBaseline,
			WindowRatio:       a.CompressionWindow,
			CompressionChange: a.CompressionRatioChange,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timeline": timeline,
		"total":    len(timeline),
	})
}

// NCDHeatmapCell represents a cell in the NCD heatmap
type NCDHeatmapCell struct {
	StreamType string  `json:"stream_type"`
	Hour       int     `json:"hour"`
	AverageNCD float64 `json:"average_ncd"`
	Count      int     `json:"count"`
}

// GetNCDHeatmap handles GET /v1/analytics/ncd-heatmap
func (h *AnalyticsHandler) GetNCDHeatmap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get anomalies for last 24 hours by default
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	filter := &models.AnomalyFilter{
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     10000,
	}

	response, err := h.storage.ListAnomalies(ctx, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get anomalies: %v", err), http.StatusInternalServerError)
		return
	}

	// Build heatmap: stream_type x hour -> average NCD
	type cellKey struct {
		streamType string
		hour       int
	}
	cellData := make(map[cellKey]struct {
		totalNCD float64
		count    int
	})

	for _, a := range response.Anomalies {
		hour := a.Timestamp.Hour()
		key := cellKey{streamType: string(a.StreamType), hour: hour}
		cell := cellData[key]
		cell.totalNCD += a.NCDScore
		cell.count++
		cellData[key] = cell
	}

	// Convert to response format
	heatmap := []NCDHeatmapCell{}
	for key, cell := range cellData {
		heatmap = append(heatmap, NCDHeatmapCell{
			StreamType: key.streamType,
			Hour:       key.hour,
			AverageNCD: cell.totalNCD / float64(cell.count),
			Count:      cell.count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"heatmap":    heatmap,
		"start_time": startTime,
		"end_time":   endTime,
	})
}
