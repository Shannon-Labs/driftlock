package main

import (
	"fmt"
	"net/http"
	"time"
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
		// TODO: Add "Create Key" endpoint that returns the raw key once.
		
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

		limit := 10000 // default Developer
		switch tc.Tenant.Plan {
		case "starter":
			limit = 500000
		case "growth":
			limit = 5000000
		}

		writeJSON(w, r, http.StatusOK, map[string]interface{}{
			"current_period_usage": usage,
			"plan_limit":           limit,
			"plan":                 tc.Tenant.Plan,
		})
	}
}


