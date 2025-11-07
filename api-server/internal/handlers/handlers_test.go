package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
	"github.com/Shannon-Labs/driftlock/api-server/internal/storage"
	"github.com/Shannon-Labs/driftlock/api-server/internal/stream"
)

// mockStorage implements a mock storage layer for testing
type mockStorage struct {
	anomalies         map[uuid.UUID]*models.Anomaly
	createAnomalyFunc func(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error)
	getAnomalyFunc    func(ctx context.Context, id uuid.UUID) (*models.Anomaly, error)
	listAnomaliesFunc func(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error)
	updateStatusFunc  func(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		anomalies: make(map[uuid.UUID]*models.Anomaly),
	}
}

func (m *mockStorage) CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error) {
	if m.createAnomalyFunc != nil {
		return m.createAnomalyFunc(ctx, create)
	}

	id := uuid.New()
	now := time.Now()

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
		CompressionRatioChange:     ((create.CompressionWindow - create.CompressionBaseline) / create.CompressionBaseline) * 100,
		ConfidenceLevel:            create.ConfidenceLevel,
		IsStatisticallySignificant: create.PValue < 0.05,
		BaselineData:               create.BaselineData,
		WindowData:                 create.WindowData,
		Metadata:                   create.Metadata,
		Tags:                       create.Tags,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	m.anomalies[id] = anomaly
	return anomaly, nil
}

func (m *mockStorage) GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error) {
	if m.getAnomalyFunc != nil {
		return m.getAnomalyFunc(ctx, id)
	}

	anomaly, ok := m.anomalies[id]
	if !ok {
		return nil, assert.AnError
	}
	return anomaly, nil
}

func (m *mockStorage) ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error) {
	if m.listAnomaliesFunc != nil {
		return m.listAnomaliesFunc(ctx, filter)
	}

	anomalies := []models.Anomaly{}
	for _, a := range m.anomalies {
		// Apply filters
		if filter.StreamType != nil && a.StreamType != *filter.StreamType {
			continue
		}
		if filter.Status != nil && a.Status != *filter.Status {
			continue
		}
		if filter.MinNCDScore != nil && a.NCDScore < *filter.MinNCDScore {
			continue
		}
		if filter.MaxPValue != nil && a.PValue > *filter.MaxPValue {
			continue
		}
		if filter.OnlySignificant && !a.IsStatisticallySignificant {
			continue
		}
		anomalies = append(anomalies, *a)
	}

	total := len(anomalies)
	limit := filter.Limit
	if limit == 0 {
		limit = 50
	}
	offset := filter.Offset

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(anomalies) {
		start = len(anomalies)
	}
	if end > len(anomalies) {
		end = len(anomalies)
	}

	paginatedAnomalies := anomalies[start:end]

	return &models.AnomalyListResponse{
		Anomalies: paginatedAnomalies,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
		HasMore:   offset+limit < total,
	}, nil
}

func (m *mockStorage) UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, update, username)
	}

	anomaly, ok := m.anomalies[id]
	if !ok {
		return assert.AnError
	}

	anomaly.Status = update.Status
	anomaly.UpdatedAt = time.Now()

	if update.Notes != nil {
		anomaly.Notes = update.Notes
	}

	if update.Status == models.StatusAcknowledged {
		now := time.Now()
		anomaly.AcknowledgedAt = &now
		anomaly.AcknowledgedBy = &username
	}

	if update.Status == models.StatusDismissed {
		now := time.Now()
		anomaly.DismissedAt = &now
		anomaly.DismissedBy = &username
	}

	return nil
}

// Verify mockStorage implements the storage interface methods
var _ storage.AnomalyStorage = (*mockStorage)(nil)

type stubEventPublisher struct {
	called  bool
	anomaly *models.Anomaly
	err     error
}

func (s *stubEventPublisher) AnomalyCreated(_ context.Context, anomaly *models.Anomaly) error {
	s.called = true
	s.anomaly = anomaly
	return s.err
}

// TestCreateAnomaly tests the CreateAnomaly handler
func TestCreateAnomaly_Success(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	createReq := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.85,
		PValue:              0.001,
		GlassBoxExplanation: "High compression ratio change detected",
		CompressionBaseline: 1000.0,
		CompressionWindow:   1500.0,
		CompressionCombined: 2500.0,
		ConfidenceLevel:     0.99,
		BaselineData:        map[string]interface{}{"test": "data"},
		WindowData:          map[string]interface{}{"test": "window"},
		Metadata:            map[string]interface{}{"source": "test"},
		Tags:                []string{"critical", "production"},
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/v1/anomalies", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateAnomaly(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Anomaly
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, response.ID)
	assert.Equal(t, createReq.NCDScore, response.NCDScore)
	assert.Equal(t, createReq.PValue, response.PValue)
	assert.Equal(t, models.StatusPending, response.Status)
	assert.True(t, response.IsStatisticallySignificant)
}

func TestCreateAnomaly_PublishesEvent(t *testing.T) {
	mockStore := newMockStorage()
	publisher := &stubEventPublisher{}
	handler := NewAnomaliesHandler(mockStore, nil, publisher)

	createReq := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.5,
		PValue:              0.01,
		GlassBoxExplanation: "test",
		CompressionBaseline: 10,
		CompressionWindow:   20,
		CompressionCombined: 30,
		ConfidenceLevel:     0.9,
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/anomalies", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateAnomaly(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, publisher.called, "expected anomaly-created event to be published")
	if publisher.anomaly == nil {
		t.Fatal("expected anomaly payload in event publisher")
	}
	if publisher.anomaly.ID == uuid.Nil {
		t.Fatal("expected anomaly to have ID populated")
	}
}

func TestCreateAnomaly_InvalidPayload(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	req := httptest.NewRequest("POST", "/v1/anomalies", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.CreateAnomaly(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAnomaly_Success(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	// Create a test anomaly
	id := uuid.New()
	now := time.Now()
	testAnomaly := &models.Anomaly{
		ID:                  id,
		Timestamp:           now,
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.75,
		PValue:              0.01,
		Status:              models.StatusPending,
		GlassBoxExplanation: "Test anomaly",
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	mockStore.anomalies[id] = testAnomaly

	req := httptest.NewRequest("GET", "/v1/anomalies/"+id.String(), nil)
	w := httptest.NewRecorder()

	handler.GetAnomaly(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Anomaly
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, id, response.ID)
	assert.Equal(t, testAnomaly.NCDScore, response.NCDScore)
}

func TestGetAnomaly_InvalidID(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	req := httptest.NewRequest("GET", "/v1/anomalies/invalid-id", nil)
	w := httptest.NewRecorder()

	handler.GetAnomaly(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAnomaly_NotFound(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	nonExistentID := uuid.New()
	req := httptest.NewRequest("GET", "/v1/anomalies/"+nonExistentID.String(), nil)
	w := httptest.NewRecorder()

	handler.GetAnomaly(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListAnomalies_Success(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	// Create test anomalies
	now := time.Now()
	for i := 0; i < 5; i++ {
		id := uuid.New()
		mockStore.anomalies[id] = &models.Anomaly{
			ID:                         id,
			Timestamp:                  now.Add(time.Duration(i) * time.Hour),
			StreamType:                 models.StreamTypeMetrics,
			NCDScore:                   float64(i) * 0.1,
			PValue:                     0.01,
			Status:                     models.StatusPending,
			IsStatisticallySignificant: true,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
	}

	req := httptest.NewRequest("GET", "/v1/anomalies", nil)
	w := httptest.NewRecorder()

	handler.ListAnomalies(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AnomalyListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 5, response.Total)
	assert.Len(t, response.Anomalies, 5)
	assert.False(t, response.HasMore)
}

func TestListAnomalies_WithFilters(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	// Create mixed anomalies
	now := time.Now()
	for i := 0; i < 10; i++ {
		id := uuid.New()
		streamType := models.StreamTypeMetrics
		if i%2 == 0 {
			streamType = models.StreamTypeLogs
		}
		mockStore.anomalies[id] = &models.Anomaly{
			ID:                         id,
			Timestamp:                  now.Add(time.Duration(i) * time.Hour),
			StreamType:                 streamType,
			NCDScore:                   float64(i) * 0.1,
			PValue:                     0.01,
			Status:                     models.StatusPending,
			IsStatisticallySignificant: true,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
	}

	req := httptest.NewRequest("GET", "/v1/anomalies?stream_type=metrics", nil)
	w := httptest.NewRecorder()

	handler.ListAnomalies(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AnomalyListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 5, response.Total)
	for _, a := range response.Anomalies {
		assert.Equal(t, models.StreamTypeMetrics, a.StreamType)
	}
}

func TestListAnomalies_Pagination(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	// Create 10 anomalies
	now := time.Now()
	for i := 0; i < 10; i++ {
		id := uuid.New()
		mockStore.anomalies[id] = &models.Anomaly{
			ID:        id,
			Timestamp: now.Add(time.Duration(i) * time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	req := httptest.NewRequest("GET", "/v1/anomalies?limit=5&offset=0", nil)
	w := httptest.NewRecorder()

	handler.ListAnomalies(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AnomalyListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 10, response.Total)
	assert.Len(t, response.Anomalies, 5)
	assert.Equal(t, 5, response.Limit)
	assert.Equal(t, 0, response.Offset)
	assert.True(t, response.HasMore)
}

func TestUpdateAnomalyStatus_Success(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	// Create test anomaly
	id := uuid.New()
	now := time.Now()
	mockStore.anomalies[id] = &models.Anomaly{
		ID:        id,
		Timestamp: now,
		Status:    models.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	updateReq := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
	}
	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PATCH", "/v1/anomalies/"+id.String()+"/status", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdateAnomalyStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Anomaly
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, models.StatusAcknowledged, response.Status)
	assert.NotNil(t, response.AcknowledgedAt)
	assert.NotNil(t, response.AcknowledgedBy)
}

func TestUpdateAnomalyStatus_InvalidID(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	updateReq := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
	}
	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PATCH", "/v1/anomalies/invalid-id/status", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdateAnomalyStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAnomalyStatus_InvalidPayload(t *testing.T) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(10)
	handler := NewAnomaliesHandler(mockStore, mockStream, nil)

	id := uuid.New()
	req := httptest.NewRequest("PATCH", "/v1/anomalies/"+id.String()+"/status", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.UpdateAnomalyStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test helper to verify the anomaly model methods
func TestAnomalyModel_IsAnomaly(t *testing.T) {
	tests := []struct {
		name     string
		anomaly  *models.Anomaly
		expected bool
	}{
		{
			name: "significant anomaly",
			anomaly: &models.Anomaly{
				IsStatisticallySignificant: true,
				PValue:                     0.01,
			},
			expected: true,
		},
		{
			name: "not significant",
			anomaly: &models.Anomaly{
				IsStatisticallySignificant: false,
				PValue:                     0.01,
			},
			expected: false,
		},
		{
			name: "high p-value",
			anomaly: &models.Anomaly{
				IsStatisticallySignificant: true,
				PValue:                     0.1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.anomaly.IsAnomaly()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnomalyModel_GetSeverity(t *testing.T) {
	tests := []struct {
		name     string
		anomaly  *models.Anomaly
		expected string
	}{
		{
			name: "critical severity",
			anomaly: &models.Anomaly{
				PValue:   0.0001,
				NCDScore: 0.8,
			},
			expected: "critical",
		},
		{
			name: "high severity",
			anomaly: &models.Anomaly{
				PValue:   0.005,
				NCDScore: 0.4,
			},
			expected: "high",
		},
		{
			name: "medium severity",
			anomaly: &models.Anomaly{
				PValue:   0.03,
				NCDScore: 0.2,
			},
			expected: "medium",
		},
		{
			name: "low severity",
			anomaly: &models.Anomaly{
				PValue:   0.1,
				NCDScore: 0.1,
			},
			expected: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.anomaly.GetSeverity()
			assert.Equal(t, tt.expected, result)
		})
	}
}
