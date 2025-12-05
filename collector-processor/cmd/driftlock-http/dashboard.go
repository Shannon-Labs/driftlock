package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handleListKeys(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		keys, err := store.listAPIKeys(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Mask keys for security (though these are hashes in DB mostly, except if we stored them differently?)
		// Actually api_keys table stores key_hash. We can't show the original key.
		// Users have to generate a new one if they lost it.

		type keyResponse struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Prefix    string `json:"prefix"` // To help identify
			CreatedAt string `json:"created_at"`
			Status    string `json:"status"`
		}

		resp := make([]keyResponse, len(keys))
		for i, k := range keys {
			status := "active"
			if k.RevokedAt != nil {
				status = "revoked"
			}
			resp[i] = keyResponse{
				ID:        k.ID.String(),
				Name:      k.Name,
				Prefix:    "dlk_...", // We don't store prefix separately in current schema, assume all are dlk_
				CreatedAt: k.CreatedAt.Format(time.RFC3339),
				Status:    status,
			}
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"keys": resp,
		})
	}
}

func handleGetUsage(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		usage, err := store.getMonthlyUsage(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		limit := 10000 // default Developer/Free
		switch tc.Tenant.Plan {
		case "starter", "radar", "signal", "basic":
			limit = 500000
		case "growth", "tensor", "pro", "sentinel", "lock", "transistor":
			limit = 5000000
		case "orbit", "enterprise", "horizon":
			limit = 25000000
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"current_period_usage": usage,
			"plan_limit":           limit,
			"plan":                 tc.Tenant.Plan,
		})
	}
}

func handleRegenerateKey(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		key, keyID, err := store.regenerateAPIKey(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "API key regenerated. Store this key securely - it won't be shown again.",
			"api_key": key,
			"key_id":  keyID.String(),
		})
	}
}

func handleCreateKey(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		var req struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
			return
		}

		if req.Name == "" {
			req.Name = "api-key"
		}
		if req.Role == "" {
			req.Role = "admin"
		}

		key, keyID, err := store.createAPIKey(r.Context(), tc.Tenant.ID, req.Name, req.Role)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "API key created. Store this key securely - it won't be shown again.",
			"api_key": key,
			"key_id":  keyID.String(),
			"name":    req.Name,
			"role":    req.Role,
		})
	}
}

func handleRevokeKey(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodDelete {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST or DELETE"))
			return
		}

		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		var req struct {
			KeyID string `json:"key_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
			return
		}

		keyID, err := uuid.Parse(req.KeyID)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid key_id"))
			return
		}

		if err := store.softRevokeKey(r.Context(), keyID, tc.Tenant.ID); err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "API key revoked",
		})
	}
}

func handleGetBillingStatus(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		status, err := store.getBillingStatus(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, status)
	}
}

type usageDetailsResponse struct {
	CurrentPeriodUsage int64               `json:"current_period_usage"`
	PlanLimit          int64               `json:"plan_limit"`
	Plan               string              `json:"plan"`
	UsagePercent       float64             `json:"usage_percent"`
	DailyUsage         []dailyUsageRecord  `json:"daily_usage"`
	StreamBreakdown    []streamUsageRecord `json:"stream_breakdown"`
	PeriodStart        time.Time           `json:"period_start"`
	PeriodEnd          time.Time           `json:"period_end"`
}

func handleGetUsageDetails(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		// Get monthly total
		monthlyUsage, err := store.getMonthlyUsage(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Get daily breakdown (last 30 days)
		dailyUsage, err := store.getDailyUsage(r.Context(), tc.Tenant.ID, 30)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Get stream breakdown
		streamUsage, err := store.getStreamUsage(r.Context(), tc.Tenant.ID)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Determine plan limit
		// Tier names: pilot (free), radar ($15), tensor ($100), orbit ($299)
		var limit int64 = 10000 // default trial/pilot (free tier)
		switch tc.Tenant.Plan {
		case "starter", "radar", "signal": // Standard tier ($15/mo)
			limit = 500_000
		case "growth", "lock", "tensor": // Pro tier ($100/mo)
			limit = 5_000_000
		case "enterprise", "orbit", "horizon": // Enterprise tier ($299/mo)
			limit = 25_000_000
		}

		// Calculate period dates
		now := time.Now()
		periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd := periodStart.AddDate(0, 1, -1)

		usagePercent := float64(monthlyUsage) / float64(limit) * 100
		if usagePercent > 100 {
			usagePercent = 100
		}

		resp := usageDetailsResponse{
			CurrentPeriodUsage: monthlyUsage,
			PlanLimit:          limit,
			Plan:               tc.Tenant.Plan,
			UsagePercent:       usagePercent,
			DailyUsage:         dailyUsage,
			StreamBreakdown:    streamUsage,
			PeriodStart:        periodStart,
			PeriodEnd:          periodEnd,
		}

		writeJSON(w, r, http.StatusOK, resp)
	}
}
