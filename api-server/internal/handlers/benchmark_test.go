package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/driftlock/api-server/internal/models"
	"github.com/your-org/driftlock/api-server/internal/stream"
)

// BenchmarkCreateAnomaly benchmarks the CreateAnomaly handler
func BenchmarkCreateAnomaly(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	createReq := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.85,
		PValue:              0.001,
		GlassBoxExplanation: "Test anomaly for benchmarking",
		CompressionBaseline: 1000.0,
		CompressionWindow:   1500.0,
		CompressionCombined: 2500.0,
		ConfidenceLevel:     0.99,
		BaselineData:        map[string]interface{}{"test": "data"},
		WindowData:          map[string]interface{}{"test": "window"},
		Metadata:            map[string]interface{}{"benchmark": true},
		Tags:                []string{"benchmark", "test"},
	}

	body, _ := json.Marshal(createReq)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/v1/anomalies", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.CreateAnomaly(w, req)
	}
}

// BenchmarkGetAnomaly benchmarks the GetAnomaly handler
func BenchmarkGetAnomaly(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	// Pre-populate with test anomaly
	id := uuid.New()
	now := time.Now()
	mockStore.anomalies[id] = &models.Anomaly{
		ID:                  id,
		Timestamp:           now,
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.75,
		PValue:              0.01,
		Status:              models.StatusPending,
		GlassBoxExplanation: "Benchmark test anomaly",
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/v1/anomalies/"+id.String(), nil)
		w := httptest.NewRecorder()
		handler.GetAnomaly(w, req)
	}
}

// BenchmarkListAnomalies benchmarks the ListAnomalies handler
func BenchmarkListAnomalies(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	// Pre-populate with test anomalies
	now := time.Now()
	for i := 0; i < 100; i++ {
		id := uuid.New()
		mockStore.anomalies[id] = &models.Anomaly{
			ID:                         id,
			Timestamp:                  now.Add(time.Duration(i) * time.Minute),
			StreamType:                 models.StreamTypeMetrics,
			NCDScore:                   float64(i) * 0.01,
			PValue:                     0.01,
			Status:                     models.StatusPending,
			IsStatisticallySignificant: true,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/v1/anomalies?limit=50", nil)
		w := httptest.NewRecorder()
		handler.ListAnomalies(w, req)
	}
}

// BenchmarkListAnomalies_WithFilters benchmarks filtered list queries
func BenchmarkListAnomalies_WithFilters(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	// Pre-populate with mixed anomalies
	now := time.Now()
	for i := 0; i < 200; i++ {
		id := uuid.New()
		streamType := models.StreamTypeMetrics
		if i%2 == 0 {
			streamType = models.StreamTypeLogs
		}
		mockStore.anomalies[id] = &models.Anomaly{
			ID:                         id,
			Timestamp:                  now.Add(time.Duration(i) * time.Minute),
			StreamType:                 streamType,
			NCDScore:                   float64(i) * 0.005,
			PValue:                     0.01,
			Status:                     models.StatusPending,
			IsStatisticallySignificant: i%3 == 0,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/v1/anomalies?stream_type=metrics&only_significant=true&limit=25", nil)
		w := httptest.NewRecorder()
		handler.ListAnomalies(w, req)
	}
}

// BenchmarkUpdateAnomalyStatus benchmarks status updates
func BenchmarkUpdateAnomalyStatus(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	// Pre-populate with test anomaly
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
	body, _ := json.Marshal(updateReq)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("PATCH", "/v1/anomalies/"+id.String()+"/status", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.UpdateAnomalyStatus(w, req)

		// Reset status for next iteration
		mockStore.anomalies[id].Status = models.StatusPending
	}
}

// BenchmarkListAnomalies_Parallel benchmarks concurrent list requests
func BenchmarkListAnomalies_Parallel(b *testing.B) {
	mockStore := newMockStorage()
	mockStream := stream.NewStreamer(1000)
	handler := NewAnomaliesHandler(mockStore, mockStream)

	// Pre-populate
	now := time.Now()
	for i := 0; i < 100; i++ {
		id := uuid.New()
		mockStore.anomalies[id] = &models.Anomaly{
			ID:        id,
			Timestamp: now.Add(time.Duration(i) * time.Minute),
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/v1/anomalies?limit=25", nil)
			w := httptest.NewRecorder()
			handler.ListAnomalies(w, req)
		}
	})
}

// BenchmarkJSON_Marshal benchmarks JSON marshaling of anomalies
func BenchmarkJSON_Marshal(b *testing.B) {
	now := time.Now()
	anomaly := &models.Anomaly{
		ID:                  uuid.New(),
		Timestamp:           now,
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.85,
		PValue:              0.001,
		Status:              models.StatusPending,
		GlassBoxExplanation: "Benchmark test",
		CompressionBaseline: 1000.0,
		CompressionWindow:   1500.0,
		CompressionCombined: 2500.0,
		ConfidenceLevel:     0.99,
		BaselineData:        map[string]interface{}{"test": "data"},
		WindowData:          map[string]interface{}{"test": "window"},
		Metadata:            map[string]interface{}{"benchmark": true},
		Tags:                []string{"benchmark", "test"},
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(anomaly)
	}
}

// BenchmarkJSON_Unmarshal benchmarks JSON unmarshaling
func BenchmarkJSON_Unmarshal(b *testing.B) {
	now := time.Now()
	anomaly := &models.Anomaly{
		ID:                  uuid.New(),
		Timestamp:           now,
		StreamType:          models.StreamTypeMetrics,
		NCDScore:            0.85,
		PValue:              0.001,
		GlassBoxExplanation: "Benchmark test",
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	data, _ := json.Marshal(anomaly)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var a models.Anomaly
		_ = json.Unmarshal(data, &a)
	}
}
