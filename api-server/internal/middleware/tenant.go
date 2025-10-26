package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/services"
)

// TenantMiddleware extracts tenant information from requests
type TenantMiddleware struct {
	tenantService *services.TenantService
}

// NewTenantMiddleware creates a new tenant middleware
func NewTenantMiddleware(tenantService *services.TenantService) *TenantMiddleware {
	return &TenantMiddleware{
		tenantService: tenantService,
	}
}

// Middleware returns a middleware function that extracts tenant information
func (m *TenantMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract tenant from subdomain or header
		tenantID := m.extractTenantID(r)

		// If no tenant ID found, use default tenant
		if tenantID == "" {
			tenantID = "default"
		}

		// Get tenant information
		tenant, err := m.tenantService.GetTenant(r.Context(), tenantID)
		if err != nil {
			// If tenant not found, return error
			http.Error(w, "Tenant not found", http.StatusNotFound)
			return
		}

		// Check if tenant is active
		if tenant.Status != models.TenantStatusActive {
			http.Error(w, "Tenant is not active", http.StatusForbidden)
			return
		}

		// Add tenant to context
		ctx := context.WithValue(r.Context(), "tenant", tenant)

		// Call next handler with tenant context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractTenantID extracts tenant ID from request
func (m *TenantMiddleware) extractTenantID(r *http.Request) string {
	// Try to extract from subdomain
	host := r.Host
	if host != "" {
		// Extract subdomain (e.g., tenant1.driftlock.com -> tenant1)
		parts := strings.Split(host, ".")
		if len(parts) >= 2 && parts[1] == "driftlock" {
			return parts[0]
		}
	}

	// Try to extract from header
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID != "" {
		return tenantID
	}

	// Try to extract from query parameter
	tenantID = r.URL.Query().Get("tenant")
	if tenantID != "" {
		return tenantID
	}

	return ""
}

// GetTenantFromContext retrieves tenant from request context
func GetTenantFromContext(ctx context.Context) *models.Tenant {
	if tenant, ok := ctx.Value("tenant").(*models.Tenant); ok {
		return tenant
	}
	return nil
}

// RequireTenant middleware ensures a tenant is present in the request
func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := GetTenantFromContext(r.Context())
		if tenant == nil {
			http.Error(w, "Tenant required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireActiveTenant middleware ensures a tenant is active
func RequireActiveTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := GetTenantFromContext(r.Context())
		if tenant == nil || tenant.Status != models.TenantStatusActive {
			http.Error(w, "Active tenant required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CheckTenantQuota middleware checks if tenant has exceeded quota
func CheckTenantQuota(resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenant := GetTenantFromContext(r.Context())
			if tenant == nil {
				http.Error(w, "Tenant required", http.StatusBadRequest)
				return
			}

			// Get tenant service from context (this is a simplified approach)
			// In a real implementation, you might inject this differently
			tenantService, ok := r.Context().Value("tenantService").(*services.TenantService)
			if !ok {
				http.Error(w, "Tenant service not available", http.StatusInternalServerError)
				return
			}

			exceeded, err := tenantService.CheckTenantQuota(r.Context(), tenant.ID, resourceType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if exceeded {
				http.Error(w, "Tenant quota exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
