package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const OrganizationIDKey contextKey = "organization_id"

import (
	"context"
	"net/http"
	"strconv"
)

// OrganizationMiddleware validates the organization ID from the API gateway
func OrganizationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get organization ID from the header set by Cloudflare Workers
			orgIDStr := r.Header.Get("X-Organization-ID")
			if orgIDStr == "" {
				http.Error(w, "Organization ID is required (from API gateway)", http.StatusUnauthorized)
				return
			}

			// Parse organization ID to uint for compatibility with existing code
			orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
			if err != nil {
				http.Error(w, "Invalid organization ID format", http.StatusBadRequest)
				return
			}

			// For backward compatibility with existing handler code, we'll set both the new key and the old "tenant_id"
			ctx := r.Context()
			ctx = context.WithValue(ctx, OrganizationIDKey, orgIDStr)
			// For backward compatibility with existing handler code
			ctx = context.WithValue(ctx, "tenant_id", uint(orgID)) 
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetOrganizationID extracts the organization ID from request context
func GetOrganizationID(r *http.Request) string {
	orgID, ok := r.Context().Value(OrganizationIDKey).(string)
	if !ok {
		return ""
	}
	return orgID
}

// WithOrganizationID adds organization ID to the request context
func WithOrganizationID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, OrganizationIDKey, orgID)
}