package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAPIKeyRegistration(t *testing.T) {
	// Save original env vars
	originalDefaultKey := os.Getenv("DEFAULT_API_KEY")
	originalDevKey := os.Getenv("DRIFTLOCK_DEV_API_KEY")
	originalOrgID := os.Getenv("DEFAULT_ORG_ID")
	
	defer func() {
		os.Setenv("DEFAULT_API_KEY", originalDefaultKey)
		os.Setenv("DRIFTLOCK_DEV_API_KEY", originalDevKey)
		os.Setenv("DEFAULT_ORG_ID", originalOrgID)
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectKeys     int
		expectAdminKey string
	}{
		{
			name: "DEFAULT_API_KEY only",
			envVars: map[string]string{
				"DEFAULT_API_KEY": "test-default-key-123",
				"DEFAULT_ORG_ID":  "test-org",
			},
			expectKeys:     1,
			expectAdminKey: "test-default-key-123",
		},
		{
			name: "Both DEFAULT_API_KEY and DRIFTLOCK_DEV_API_KEY",
			envVars: map[string]string{
				"DEFAULT_API_KEY":      "prod-key-456",
				"DRIFTLOCK_DEV_API_KEY": "dev-key-789",
				"DEFAULT_ORG_ID":       "multi-org",
			},
			expectKeys:     2,
			expectAdminKey: "prod-key-456",
		},
		{
			name: "DRIFTLOCK_DEV_API_KEY only",
			envVars: map[string]string{
				"DRIFTLOCK_DEV_API_KEY": "dev-only-key",
				"DEFAULT_ORG_ID":        "dev-org",
			},
			expectKeys:     1,
			expectAdminKey: "dev-only-key",
		},
		{
			name: "No keys set",
			envVars: map[string]string{
				"DEFAULT_ORG_ID": "default-org",
			},
			expectKeys:     0,
			expectAdminKey: "",
		},
		{
			name: "Placeholder key ignored",
			envVars: map[string]string{
				"DEFAULT_API_KEY": "your_api_key_here_for_dashboard_access",
				"DEFAULT_ORG_ID":  "test-org",
			},
			expectKeys:     0,
			expectAdminKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env
			os.Unsetenv("DEFAULT_API_KEY")
			os.Unsetenv("DRIFTLOCK_DEV_API_KEY")
			os.Unsetenv("DEFAULT_ORG_ID")

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Create authenticator and register keys
			auth := NewAuthenticator()
			defaultOrgID := strings.TrimSpace(os.Getenv("DEFAULT_ORG_ID"))
			if defaultOrgID == "" {
				defaultOrgID = "default"
			}

			commonScopes := []string{"read:anomalies", "write:anomalies", "admin:config"}
			registered := 0

			registerAPIKey := func(rawKey, label string, info APIKeyInfo) {
				rawKey = strings.TrimSpace(rawKey)
				placeholderKey := "your_api_key_here_for_dashboard_access"
				if rawKey == "" || rawKey == placeholderKey {
					return
				}
				auth.AddAPIKey(HashAPIKey(rawKey), info)
				registered++
			}

			registerAPIKey(
				os.Getenv("DEFAULT_API_KEY"),
				"DEFAULT_API_KEY",
				APIKeyInfo{
					Name:           "default",
					Role:           "admin",
					Scopes:         commonScopes,
					OrganizationID: defaultOrgID,
				},
			)

			registerAPIKey(
				os.Getenv("DRIFTLOCK_DEV_API_KEY"),
				"DRIFTLOCK_DEV_API_KEY",
				APIKeyInfo{
					Name:           "development",
					Role:           "admin",
					Scopes:         commonScopes,
					OrganizationID: defaultOrgID,
				},
			)

			if registered != tt.expectKeys {
				t.Errorf("Expected %d keys registered, got %d", tt.expectKeys, registered)
			}

			// Verify the admin key works if expected
			if tt.expectAdminKey != "" {
				info, ok := auth.apiKeys[HashAPIKey(tt.expectAdminKey)]
				if !ok {
					t.Error("Expected admin key to be valid")
				}
				if info.Role != "admin" {
					t.Errorf("Expected admin role, got %s", info.Role)
				}
			}
		})
	}
}

func TestMiddlewareAuthentication(t *testing.T) {
	auth := NewAuthenticator()
	
	// Register a test key
	testKey := "test-api-key-12345"
	auth.AddAPIKey(HashAPIKey(testKey), APIKeyInfo{
		Name:   "test-key",
		Role:   "admin",
		Scopes: []string{"read:anomalies", "write:anomalies"},
	})

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectStatus   int
		expectAuthenticated bool
	}{
		{
			name: "Valid API key in Authorization header",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+testKey)
			},
			expectStatus:   http.StatusOK,
			expectAuthenticated: true,
		},
		{
			name: "Invalid API key",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-key")
			},
			expectStatus:   http.StatusUnauthorized,
			expectAuthenticated: false,
		},
		{
			name: "No API key",
			setupRequest: func(r *http.Request) {
				// Don't set any auth header
			},
			expectStatus:   http.StatusUnauthorized,
			expectAuthenticated: false,
		},
		{
			name: "Malformed Authorization header",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat")
			},
			expectStatus:   http.StatusUnauthorized,
			expectAuthenticated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that checks authentication
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if authentication succeeded by looking for context values
				username := GetUsernameFromContext(r.Context())
				if username == "anonymous" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with auth middleware
			middleware := auth.Middleware(handler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)
			
			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, rec.Code)
			}
		})
	}
}

func TestOptionalMiddleware(t *testing.T) {
	auth := NewAuthenticator()
	
	// Register a test key
	testKey := "optional-test-key"
	auth.AddAPIKey(HashAPIKey(testKey), APIKeyInfo{
		Name:   "optional-test",
		Role:   "viewer",
		Scopes: []string{"read:anomalies"},
	})

	tests := []struct {
		name         string
		setupRequest func(*http.Request)
		expectStatus int
	}{
		{
			name: "Valid API key (authenticated)",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+testKey)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "No API key (anonymous, allowed)",
			setupRequest: func(r *http.Request) {
				// No auth header
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "Invalid API key (anonymous, allowed)",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid")
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Optional auth should always allow the request
				w.WriteHeader(http.StatusOK)
			})

			middleware := auth.OptionalMiddleware(handler)

			req := httptest.NewRequest("GET", "/stream", nil)
			tt.setupRequest(req)
			
			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, rec.Code)
			}
		})
	}
}

func TestHashAPIKey(t *testing.T) {
	key1 := "test-api-key-12345"
	key2 := "test-api-key-12345"
	key3 := "different-key"

	hash1 := HashAPIKey(key1)
	hash2 := HashAPIKey(key2)
	hash3 := HashAPIKey(key3)

	if hash1 != hash2 {
		t.Error("Same keys should produce same hash")
	}

	if hash1 == hash3 {
		t.Error("Different keys should produce different hashes")
	}

	if len(hash1) != 64 { // SHA-256 hex length
		t.Errorf("Expected 64 character hash, got %d", len(hash1))
	}
}