package auth

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "strings"
    "sync"
    "time"

    "github.com/shannon-labs/driftlock/api-server/internal/ctxutil"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UsernameContextKey stores the authenticated username in context
	UsernameContextKey ContextKey = "username"
	// RoleContextKey stores the user role in context
	RoleContextKey ContextKey = "role"
)

// Authenticator handles authentication
type Authenticator struct {
	apiKeys map[string]APIKeyInfo
	mu      sync.RWMutex
}

// APIKeyInfo stores information about an API key
type APIKeyInfo struct {
    Name     string
    Role     string
    Scopes   []string
    LastUsed time.Time
    // Optional: associated organization ID for metering and scoping
    OrganizationID string
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator() *Authenticator {
	return &Authenticator{
		apiKeys: make(map[string]APIKeyInfo),
	}
}

// AddAPIKey registers an API key (key should already be hashed)
func (a *Authenticator) AddAPIKey(keyHash string, info APIKeyInfo) {
	a.apiKeys[keyHash] = info
}

// HashAPIKey hashes an API key for storage
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// Middleware returns an HTTP middleware for authentication
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract API key from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Parse "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		apiKey := parts[1]
		keyHash := HashAPIKey(apiKey)

		// Validate API key
		info, ok := a.apiKeys[keyHash]
		if !ok {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		// Update last used time
		info.LastUsed = time.Now()
		a.apiKeys[keyHash] = info

        // Add authentication info to context
        ctx := context.WithValue(r.Context(), UsernameContextKey, info.Name)
        ctx = context.WithValue(ctx, string(UsernameContextKey), info.Name)
        ctx = context.WithValue(ctx, RoleContextKey, info.Role)
        ctx = context.WithValue(ctx, string(RoleContextKey), info.Role)

        // If API key ties to an organization, propagate into event context
        if info.OrganizationID != "" {
            ctx = ctxutil.WithEventContext(ctx, info.OrganizationID, "")
        }

		// Call next handler with authenticated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalMiddleware allows requests with or without authentication
func (a *Authenticator) OptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				apiKey := parts[1]
				keyHash := HashAPIKey(apiKey)

            if info, ok := a.apiKeys[keyHash]; ok {
                info.LastUsed = time.Now()
                a.apiKeys[keyHash] = info

                ctx := context.WithValue(r.Context(), UsernameContextKey, info.Name)
                ctx = context.WithValue(ctx, string(UsernameContextKey), info.Name)
                ctx = context.WithValue(ctx, RoleContextKey, info.Role)
                ctx = context.WithValue(ctx, string(RoleContextKey), info.Role)
                if info.OrganizationID != "" {
                    ctx = ctxutil.WithEventContext(ctx, info.OrganizationID, "")
                }
                r = r.WithContext(ctx)
            }
        }
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole returns a middleware that checks for a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleContextKey).(string)
			if !ok || userRole != role {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole returns a middleware that checks for any of the specified roles
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	roleMap := make(map[string]bool)
	for _, r := range roles {
		roleMap[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleContextKey).(string)
			if !ok || !roleMap[userRole] {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUsernameFromContext extracts username from context
func GetUsernameFromContext(ctx context.Context) string {
	if username, ok := ctx.Value(UsernameContextKey).(string); ok {
		return username
	}

	if username, ok := ctx.Value(string(UsernameContextKey)).(string); ok {
		return username
	}

	return "anonymous"
}

// GetRoleFromContext extracts role from context
func GetRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(RoleContextKey).(string); ok {
		return role
	}

	if role, ok := ctx.Value(string(RoleContextKey)).(string); ok {
		return role
	}

	return "viewer"
}
