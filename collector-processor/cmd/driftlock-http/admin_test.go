package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAdminStore struct {
	mock.Mock
}

func (m *mockAdminStore) ListTenants(ctx context.Context, limit, offset int) ([]adminTenant, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]adminTenant), args.Int(1), args.Error(2)
}

func (m *mockAdminStore) GetTenant(ctx context.Context, id uuid.UUID) (adminTenantDetail, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(adminTenantDetail), args.Error(1)
}

func (m *mockAdminStore) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockAdminStore) GetStats(ctx context.Context) (adminStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(adminStats), args.Error(1)
}

func TestAdminAuthMiddleware(t *testing.T) {
	handler := withAdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Setenv("ADMIN_API_KEY", "secret")

	// Case 1: Valid Key
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Admin-Key", "secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Invalid Key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Admin-Key", "wrong")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Case 3: Missing Key
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Case 4: Key not configured
	t.Setenv("ADMIN_API_KEY", "")
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestListTenants(t *testing.T) {
	mockStore := new(mockAdminStore)
	handler := adminHandler(mockStore)

	expectedTenants := []adminTenant{
		{ID: "1", Name: "T1", Email: "t1@example.com"},
	}
	mockStore.On("ListTenants", mock.Anything, 20, 0).Return(expectedTenants, 1, nil)

	req := httptest.NewRequest("GET", "/tenants", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp adminTenantListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "T1", resp.Tenants[0].Name)
}

func TestGetTenant(t *testing.T) {
	mockStore := new(mockAdminStore)
	handler := adminHandler(mockStore)

	uid := uuid.New()
	expectedTenant := adminTenantDetail{
		adminTenant: adminTenant{ID: uid.String(), Name: "T1"},
	}
	mockStore.On("GetTenant", mock.Anything, uid).Return(expectedTenant, nil)

	req := httptest.NewRequest("GET", fmt.Sprintf("/tenants/%s", uid), nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp adminTenantDetail
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "T1", resp.Name)
}

func TestUpdateTenantStatus(t *testing.T) {
	mockStore := new(mockAdminStore)
	handler := adminHandler(mockStore)

	uid := uuid.New()
	mockStore.On("UpdateTenantStatus", mock.Anything, uid, "suspended").Return(nil)

	payload := map[string]string{"status": "suspended"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/tenants/%s/status", uid), bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockStore.AssertExpectations(t)
}

func TestGetStats(t *testing.T) {
	mockStore := new(mockAdminStore)
	handler := adminHandler(mockStore)

	stats := adminStats{TotalTenants: 10}
	mockStore.On("GetStats", mock.Anything).Return(stats, nil)

	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp adminStats
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 10, resp.TotalTenants)
}
