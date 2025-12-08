package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mocks -------------------------------------------------------------------

type mockProfileStore struct {
	settings *StreamTuningSettings
	verifyOK bool
	updated  struct {
		profile  string
		auto     *bool
		adaptive *bool
	}
	updateErr error
}

func (m *mockProfileStore) getStreamTuningSettings(ctx context.Context, streamID uuid.UUID) (*StreamTuningSettings, error) {
	return m.settings, nil
}

func (m *mockProfileStore) verifyStreamOwnership(ctx context.Context, tenantID, streamID uuid.UUID) bool {
	return m.verifyOK
}

func (m *mockProfileStore) updateStreamProfile(ctx context.Context, streamID uuid.UUID, profile string, autoTune, adaptiveWindow *bool) error {
	m.updated.profile = profile
	m.updated.auto = autoTune
	m.updated.adaptive = adaptiveWindow
	return m.updateErr
}

type mockTuningStore struct {
	mockProfileStore
	history []TuneHistoryRecord
	stats   *FeedbackStats
}

func (m *mockTuningStore) getTuneHistory(ctx context.Context, streamID uuid.UUID, limit int) ([]TuneHistoryRecord, error) {
	return m.history, nil
}

func (m *mockTuningStore) getFeedbackStats(ctx context.Context, streamID uuid.UUID, since time.Time) (*FeedbackStats, error) {
	if m.stats != nil {
		return m.stats, nil
	}
	return &FeedbackStats{}, nil
}

type mockFeedbackStore struct {
	anomaly        *AnomalyRecord
	recorded       *FeedbackRecord
	autoTuneCalled bool
}

func (m *mockFeedbackStore) getAnomalyByID(ctx context.Context, tenantID, anomalyID uuid.UUID) (*AnomalyRecord, error) {
	return m.anomaly, nil
}

func (m *mockFeedbackStore) recordFeedback(ctx context.Context, feedback FeedbackRecord) error {
	m.recorded = &feedback
	return nil
}

func (m *mockFeedbackStore) applyAutoTune(ctx context.Context, streamID uuid.UUID) error {
	m.autoTuneCalled = true
	return nil
}

// --- Tests -------------------------------------------------------------------

func withTenant(r *http.Request, tenantID uuid.UUID) *http.Request {
	tc := tenantContext{Tenant: tenantRecord{ID: tenantID, Plan: "free"}}
	return r.WithContext(context.WithValue(r.Context(), tenantContextKey, tc))
}

func TestStreamProfileHandler_Get(t *testing.T) {
	streamID := uuid.New()
	store := &mockProfileStore{
		settings: &StreamTuningSettings{
			DetectionProfile:      "balanced",
			AutoTuneEnabled:       true,
			AdaptiveWindowEnabled: true,
			NCDThreshold:          0.3,
			PValueThreshold:       0.05,
			BaselineSize:          400,
			WindowSize:            50,
		},
		verifyOK: true,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/streams/"+streamID.String()+"/profile", nil)
	req = withTenant(req, uuid.New())
	rr := httptest.NewRecorder()

	streamProfileHandler(store)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp["profile"] != "balanced" {
		t.Fatalf("expected profile balanced, got %v", resp["profile"])
	}
}

func TestStreamProfileHandler_Patch(t *testing.T) {
	streamID := uuid.New()
	store := &mockProfileStore{
		settings: &StreamTuningSettings{
			DetectionProfile: "balanced",
		},
		verifyOK: true,
	}

	payload := map[string]interface{}{
		"profile":                 "strict",
		"auto_tune_enabled":       true,
		"adaptive_window_enabled": true,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/v1/streams/"+streamID.String()+"/profile", bytes.NewReader(body))
	req = withTenant(req, uuid.New())
	rr := httptest.NewRecorder()

	streamProfileHandler(store)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if store.updated.profile != "strict" {
		t.Fatalf("expected profile update to strict, got %s", store.updated.profile)
	}
	if store.updated.auto == nil || !*store.updated.auto {
		t.Fatalf("expected auto_tune_enabled to be true")
	}
	if store.updated.adaptive == nil || !*store.updated.adaptive {
		t.Fatalf("expected adaptive_window_enabled to be true")
	}
}

func TestAnomalyFeedbackHandler(t *testing.T) {
	anomalyID := uuid.New()
	streamID := uuid.New()
	store := &mockFeedbackStore{
		anomaly: &AnomalyRecord{
			ID:         anomalyID,
			StreamID:   streamID,
			TenantID:   uuid.New(),
			NCD:        0.4,
			PValue:     0.02,
			Confidence: 0.9,
		},
	}

	payload := map[string]string{
		"feedback_type": "confirmed",
		"reason":        "looks real",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/anomalies/"+anomalyID.String()+"/feedback", bytes.NewReader(body))
	req = withTenant(req, store.anomaly.TenantID)
	rr := httptest.NewRecorder()

	anomalyFeedbackHandler(store)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	time.Sleep(10 * time.Millisecond)
	if store.recorded == nil {
		t.Fatalf("expected feedback to be recorded")
	}
	if !store.autoTuneCalled {
		t.Fatalf("expected auto-tune to be triggered")
	}
}
