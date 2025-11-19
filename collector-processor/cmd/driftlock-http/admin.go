package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func withAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminKey := os.Getenv("ADMIN_API_KEY")
		if adminKey == "" {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("admin access disabled"))
			return
		}
		if r.Header.Get("X-Admin-Key") == adminKey {
			next.ServeHTTP(w, r)
			return
		}
		writeError(w, r, http.StatusForbidden, fmt.Errorf("forbidden"))
	})
}

type adminStore interface {
	ListTenants(ctx context.Context, limit, offset int) ([]adminTenant, int, error)
	GetTenant(ctx context.Context, id uuid.UUID) (adminTenantDetail, error)
	UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error
	GetStats(ctx context.Context) (adminStats, error)
}

func adminHandler(store adminStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/tenants", handleListTenants(store))
	mux.HandleFunc("/tenants/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/tenants/")
		parts := strings.Split(path, "/")

		if len(parts) == 1 && parts[0] != "" {
			id, err := uuid.Parse(parts[0])
			if err != nil {
				writeError(w, r, http.StatusNotFound, fmt.Errorf("invalid tenant id"))
				return
			}
			handleGetTenant(store, id).ServeHTTP(w, r)
			return
		} else if len(parts) == 2 && parts[1] == "status" {
			id, err := uuid.Parse(parts[0])
			if err != nil {
				writeError(w, r, http.StatusNotFound, fmt.Errorf("invalid tenant id"))
				return
			}
			handleUpdateTenantStatus(store, id).ServeHTTP(w, r)
			return
		}
		writeError(w, r, http.StatusNotFound, fmt.Errorf("resource not found"))
	})
	mux.HandleFunc("/stats", handleGetStats(store))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
}

type adminTenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type adminTenantListResponse struct {
	Tenants []adminTenant `json:"tenants"`
	Total   int           `json:"total"`
}

func handleListTenants(store adminStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		limit := 20
		offset := 0
		if v := r.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				limit = n
			}
		}
		if v := r.URL.Query().Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		tenants, total, err := store.ListTenants(ctx, limit, offset)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, adminTenantListResponse{
			Tenants: tenants,
			Total:   total,
		})
	}
}

type dailyUsage struct {
	Date            string `json:"date"`
	EventCount      int    `json:"event_count"`
	ApiRequestCount int    `json:"api_request_count"`
	AnomalyCount    int    `json:"anomaly_count"`
}

type adminTenantDetail struct {
	adminTenant
	Usage []dailyUsage `json:"usage_metrics"`
}

func handleGetTenant(store adminStore, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		t, err := store.GetTenant(ctx, id)
		if err != nil {
			if err.Error() == "tenant not found" {
				writeError(w, r, http.StatusNotFound, err)
				return
			}
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, t)
	}
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func handleUpdateTenantStatus(store adminStore, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		var req updateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid json"))
			return
		}

		if req.Status == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("status required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := store.UpdateTenantStatus(ctx, id, req.Status); err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]string{"status": "updated"})
	}
}

type adminStats struct {
	TotalTenants        int `json:"total_tenants"`
	ActiveSubscriptions int `json:"active_subscriptions"`
	TotalAnomalies24h   int `json:"total_anomalies_24h"`
}

func handleGetStats(store adminStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		stats, err := store.GetStats(ctx)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, stats)
	}
}

// Implementation on *store

func (s *store) ListTenants(ctx context.Context, limit, offset int) ([]adminTenant, int, error) {
	var total int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, email, plan, status, created_at
		FROM tenants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tenants := make([]adminTenant, 0)
	for rows.Next() {
		var t adminTenant
		var id uuid.UUID
		if err := rows.Scan(&id, &t.Name, &t.Email, &t.Plan, &t.Status, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		t.ID = id.String()
		tenants = append(tenants, t)
	}
	return tenants, total, nil
}

func (s *store) GetTenant(ctx context.Context, id uuid.UUID) (adminTenantDetail, error) {
	var t adminTenantDetail
	var uid uuid.UUID
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, email, plan, status, created_at
		FROM tenants WHERE id=$1`, id).Scan(&uid, &t.Name, &t.Email, &t.Plan, &t.Status, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return t, fmt.Errorf("tenant not found")
		}
		return t, err
	}
	t.ID = uid.String()

	// Get usage metrics (last 30 days)
	rows, err := s.pool.Query(ctx, `
		SELECT date::text, SUM(event_count), SUM(api_request_count), SUM(anomaly_count)
		FROM usage_metrics
		WHERE tenant_id = $1 AND date >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY date
		ORDER BY date ASC`, id)
	if err != nil {
		return t, err
	}
	defer rows.Close()

	t.Usage = make([]dailyUsage, 0)
	for rows.Next() {
		var u dailyUsage
		if err := rows.Scan(&u.Date, &u.EventCount, &u.ApiRequestCount, &u.AnomalyCount); err != nil {
			return t, err
		}
		t.Usage = append(t.Usage, u)
	}
	return t, nil
}

func (s *store) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := s.pool.Exec(ctx, `UPDATE tenants SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	if err != nil {
		return err
	}
	_ = s.loadCache(ctx)
	return nil
}

func (s *store) GetStats(ctx context.Context) (adminStats, error) {
	var stats adminStats
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants`).Scan(&stats.TotalTenants); err != nil {
		return stats, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE stripe_status='active'`).Scan(&stats.ActiveSubscriptions); err != nil {
		return stats, err
	}
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM anomalies WHERE detected_at >= NOW() - INTERVAL '24 hours'`).Scan(&stats.TotalAnomalies24h); err != nil {
		return stats, err
	}
	return stats, nil
}
