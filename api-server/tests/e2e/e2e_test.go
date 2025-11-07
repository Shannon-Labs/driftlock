package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Shannon-Labs/driftlock/api-server/internal/export"
	"github.com/Shannon-Labs/driftlock/api-server/internal/handlers"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
)

// fakeStorage provides an in-memory implementation of storage interfaces used by the handlers.
type fakeStorage struct {
	mu        sync.RWMutex
	anomalies map[uuid.UUID]*models.Anomaly
	config    *models.DetectionConfig
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		anomalies: make(map[uuid.UUID]*models.Anomaly),
		config: &models.DetectionConfig{
			ID:              1,
			NCDThreshold:    0.25,
			PValueThreshold: 0.05,
			BaselineSize:    100,
			WindowSize:      100,
			HopSize:         50,
			CreatedAt:       time.Now().Add(-1 * time.Hour),
			UpdatedAt:       time.Now().Add(-1 * time.Hour),
			IsActive:        true,
		},
	}
}

func (s *fakeStorage) CreateAnomaly(_ context.Context, create *models.AnomalyCreate) (*models.Anomaly, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New()
	now := time.Now()
	compressionChange := 0.0
	if create.CompressionBaseline != 0 {
		compressionChange = ((create.CompressionWindow - create.CompressionBaseline) / create.CompressionBaseline) * 100
	}

	anomaly := &models.Anomaly{
		ID:                         id,
		Timestamp:                  create.Timestamp,
		StreamType:                 create.StreamType,
		NCDScore:                   create.NCDScore,
		PValue:                     create.PValue,
		Status:                     models.StatusPending,
		GlassBoxExplanation:        create.GlassBoxExplanation,
		CompressionBaseline:        create.CompressionBaseline,
		CompressionWindow:          create.CompressionWindow,
		CompressionCombined:        create.CompressionCombined,
		CompressionRatioChange:     compressionChange,
		ConfidenceLevel:            create.ConfidenceLevel,
		IsStatisticallySignificant: create.PValue < 0.05,
		BaselineData:               cloneMap(create.BaselineData),
		WindowData:                 cloneMap(create.WindowData),
		Metadata:                   cloneMap(create.Metadata),
		Tags:                       append([]string(nil), create.Tags...),
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	s.anomalies[id] = anomaly
	return cloneAnomaly(anomaly), nil
}

func (s *fakeStorage) GetAnomaly(_ context.Context, id uuid.UUID) (*models.Anomaly, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	anomaly, ok := s.anomalies[id]
	if !ok {
		return nil, fmt.Errorf("anomaly not found: %s", id)
	}
	return cloneAnomaly(anomaly), nil
}

func (s *fakeStorage) ListAnomalies(_ context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.Anomaly, 0, len(s.anomalies))
	for _, anomaly := range s.anomalies {
		if filter != nil {
			if filter.StreamType != nil && anomaly.StreamType != *filter.StreamType {
				continue
			}
			if filter.Status != nil && anomaly.Status != *filter.Status {
				continue
			}
			if filter.OnlySignificant && !anomaly.IsStatisticallySignificant {
				continue
			}
			if filter.StartTime != nil && anomaly.Timestamp.Before(*filter.StartTime) {
				continue
			}
			if filter.EndTime != nil && anomaly.Timestamp.After(*filter.EndTime) {
				continue
			}
		}
		items = append(items, *cloneAnomaly(anomaly))
	}

	limit := len(items)
	if filter != nil && filter.Limit > 0 && filter.Limit < limit {
		limit = filter.Limit
	}

	response := &models.AnomalyListResponse{
		Anomalies: items,
		Total:     len(items),
		Limit:     limit,
		Offset:    0,
		HasMore:   false,
	}

	return response, nil
}

func (s *fakeStorage) UpdateAnomalyStatus(_ context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	anomaly, ok := s.anomalies[id]
	if !ok {
		return fmt.Errorf("anomaly not found: %s", id)
	}

	anomaly.Status = update.Status
	anomaly.UpdatedAt = time.Now()
	anomaly.Notes = update.Notes
	if update.Status == models.StatusAcknowledged {
		anomaly.AcknowledgedBy = &username
		now := anomaly.UpdatedAt
		anomaly.AcknowledgedAt = &now
	}
	return nil
}

func (s *fakeStorage) GetActiveConfig(context.Context) (*models.DetectionConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneConfig(s.config), nil
}

func (s *fakeStorage) UpdateConfig(_ context.Context, update *models.DetectionConfigUpdate) (*models.DetectionConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if update.NCDThreshold != nil {
		s.config.NCDThreshold = *update.NCDThreshold
	}
	if update.PValueThreshold != nil {
		s.config.PValueThreshold = *update.PValueThreshold
	}
	if update.BaselineSize != nil {
		s.config.BaselineSize = *update.BaselineSize
	}
	if update.WindowSize != nil {
		s.config.WindowSize = *update.WindowSize
	}
	if update.HopSize != nil {
		s.config.HopSize = *update.HopSize
	}
	if update.StreamOverrides != nil {
		s.config.StreamOverrides = cloneMap(update.StreamOverrides)
	}
	if update.Notes != nil {
		s.config.Notes = update.Notes
	}
	s.config.UpdatedAt = time.Now()

	return cloneConfig(s.config), nil
}

// Helpers --------------------------------------------------------------------

func cloneAnomaly(a *models.Anomaly) *models.Anomaly {
	if a == nil {
		return nil
	}
	clone := *a
	clone.BaselineData = cloneMap(a.BaselineData)
	clone.WindowData = cloneMap(a.WindowData)
	clone.Metadata = cloneMap(a.Metadata)
	if a.Tags != nil {
		clone.Tags = append([]string(nil), a.Tags...)
	}
	return &clone
}

func cloneConfig(cfg *models.DetectionConfig) *models.DetectionConfig {
	if cfg == nil {
		return nil
	}
	clone := *cfg
	clone.StreamOverrides = cloneMap(cfg.StreamOverrides)
	return &clone
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	clone := make(map[string]interface{}, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

func withTestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "username", "e2e-user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setupTestServer(t *testing.T) (*httptest.Server, *fakeStorage) {
	store := newFakeStorage()
	exporter := export.NewExporter(true)

	anomaliesHandler := handlers.NewAnomaliesHandler(store, nil, nil)
	configHandler := handlers.NewConfigHandler(store)
	analyticsHandler := handlers.NewAnalyticsHandler(store)
	exportHandler := handlers.NewExportHandler(store, exporter)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.Handle("/v1/anomalies", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			anomaliesHandler.ListAnomalies(w, r)
		case http.MethodPost:
			anomaliesHandler.CreateAnomaly(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/v1/anomalies/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/status") && r.Method == http.MethodPatch:
			anomaliesHandler.UpdateAnomalyStatus(w, r)
		case strings.HasSuffix(path, "/export") && r.Method == http.MethodGet:
			exportHandler.ExportAnomaly(w, r)
		case r.Method == http.MethodGet:
			anomaliesHandler.GetAnomaly(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))

	mux.Handle("/v1/config", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			configHandler.GetConfig(w, r)
		case http.MethodPatch:
			configHandler.UpdateConfig(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle("/v1/analytics/summary", http.HandlerFunc(analyticsHandler.GetSummary))
	mux.Handle("/v1/analytics/compression-timeline", http.HandlerFunc(analyticsHandler.GetCompressionTimeline))
	mux.Handle("/v1/analytics/ncd-heatmap", http.HandlerFunc(analyticsHandler.GetNCDHeatmap))

	server := httptest.NewServer(withTestContext(mux))
	t.Cleanup(server.Close)

	return server, store
}

func createAnomalyPayload() *models.AnomalyCreate {
	return &models.AnomalyCreate{
		Timestamp:           time.Now().UTC().Truncate(time.Second),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.42,
		PValue:              0.01,
		GlassBoxExplanation: "Compression ratio increased significantly",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
		BaselineData: map[string]interface{}{
			"service": "api",
		},
		WindowData: map[string]interface{}{
			"service":    "api",
			"error_rate": 0.21,
		},
		Metadata: map[string]interface{}{
			"source": "e2e-test",
		},
		Tags: []string{"e2e", "critical"},
	}
}

func TestE2E_APIFlow(t *testing.T) {
	server, _ := setupTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Create anomaly via POST /v1/anomalies
	payload := createAnomalyPayload()
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := client.Post(server.URL+"/v1/anomalies", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Anomaly
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.Equal(t, payload.StreamType, created.StreamType)
	assert.True(t, created.IsStatisticallySignificant)

	// Retrieve anomaly via GET /v1/anomalies/:id
	getResp, err := client.Get(server.URL + "/v1/anomalies/" + created.ID.String())
	require.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var fetched models.Anomaly
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&fetched))
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.GlassBoxExplanation, fetched.GlassBoxExplanation)

	// List anomalies
	listResp, err := client.Get(server.URL + "/v1/anomalies")
	require.NoError(t, err)
	defer listResp.Body.Close()
	assert.Equal(t, http.StatusOK, listResp.StatusCode)

	var list models.AnomalyListResponse
	require.NoError(t, json.NewDecoder(listResp.Body).Decode(&list))
	assert.Equal(t, 1, list.Total)
	assert.Len(t, list.Anomalies, 1)

	// Update anomaly status
	updateBody, err := json.Marshal(models.AnomalyUpdate{Status: models.StatusAcknowledged})
	require.NoError(t, err)

	updateReq, err := http.NewRequest(http.MethodPatch, server.URL+"/v1/anomalies/"+created.ID.String()+"/status", bytes.NewReader(updateBody))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	require.NoError(t, err)
	defer updateResp.Body.Close()
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var updated models.Anomaly
	require.NoError(t, json.NewDecoder(updateResp.Body).Decode(&updated))
	assert.Equal(t, models.StatusAcknowledged, updated.Status)
	assert.NotNil(t, updated.AcknowledgedBy)
	assert.Equal(t, "e2e-user", *updated.AcknowledgedBy)

	// Analytics summary should reflect a single anomaly
	summaryResp, err := client.Get(server.URL + "/v1/analytics/summary")
	require.NoError(t, err)
	defer summaryResp.Body.Close()
	assert.Equal(t, http.StatusOK, summaryResp.StatusCode)

	var summary struct {
		TotalAnomalies       int                  `json:"total_anomalies"`
		SignificantAnomalies int                  `json:"significant_anomalies"`
		AnomaliesByStatus    map[string]int       `json:"anomalies_by_status"`
		AnomaliesByStream    map[string]int       `json:"anomalies_by_stream_type"`
		AverageNCDScore      float64              `json:"average_ncd_score"`
		AveragePValue        float64              `json:"average_p_value"`
		TimeRange            map[string]time.Time `json:"time_range"`
	}
	require.NoError(t, json.NewDecoder(summaryResp.Body).Decode(&summary))
	assert.Equal(t, 1, summary.TotalAnomalies)
	assert.Equal(t, 1, summary.SignificantAnomalies)
	assert.Equal(t, 0.42, summary.AverageNCDScore)
	assert.InDelta(t, 0.01, summary.AveragePValue, 1e-6)
	assert.Equal(t, 1, summary.AnomaliesByStatus[string(models.StatusAcknowledged)])

	// Export anomaly evidence bundle
	exportResp, err := client.Get(server.URL + "/v1/anomalies/" + created.ID.String() + "/export")
	require.NoError(t, err)
	defer exportResp.Body.Close()
	assert.Equal(t, http.StatusOK, exportResp.StatusCode)
	assert.Equal(t, "application/json", exportResp.Header.Get("Content-Type"))

	var bundle models.EvidenceBundle
	require.NoError(t, json.NewDecoder(exportResp.Body).Decode(&bundle))
	assert.Equal(t, created.ID, bundle.Anomaly.ID)
	assert.Equal(t, "e2e-user", bundle.ExportedBy)
	assert.NotEmpty(t, bundle.Version)

	// Fetch and update config
	configResp, err := client.Get(server.URL + "/v1/config")
	require.NoError(t, err)
	defer configResp.Body.Close()
	assert.Equal(t, http.StatusOK, configResp.StatusCode)

	var cfg models.DetectionConfig
	require.NoError(t, json.NewDecoder(configResp.Body).Decode(&cfg))
	assert.True(t, cfg.IsActive)

	newThreshold := 0.3
	configUpdate := models.DetectionConfigUpdate{NCDThreshold: &newThreshold}
	configUpdateBody, err := json.Marshal(configUpdate)
	require.NoError(t, err)

	patchReq, err := http.NewRequest(http.MethodPatch, server.URL+"/v1/config", bytes.NewReader(configUpdateBody))
	require.NoError(t, err)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := client.Do(patchReq)
	require.NoError(t, err)
	defer patchResp.Body.Close()
	assert.Equal(t, http.StatusOK, patchResp.StatusCode)

	var updatedCfg models.DetectionConfig
	require.NoError(t, json.NewDecoder(patchResp.Body).Decode(&updatedCfg))
	assert.Equal(t, newThreshold, updatedCfg.NCDThreshold)
}
