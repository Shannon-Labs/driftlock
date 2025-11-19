package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Signup request/response types
type signupRequest struct {
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Plan        string `json:"plan,omitempty"`
}

type signupResponse struct {
	Success    bool   `json:"success"`
	TenantID   string `json:"tenant_id"`
	TenantSlug string `json:"tenant_slug"`
	StreamID   string `json:"stream_id"`
	APIKey     string `json:"api_key"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
}

// Rate limiter for signup requests
type signupRateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newSignupRateLimiter(limit int, window time.Duration) *signupRateLimiter {
	return &signupRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (l *signupRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	// Get existing requests and filter to window
	existing := l.requests[ip]
	var valid []time.Time
	for _, t := range existing {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.limit {
		l.requests[ip] = valid
		return false
	}

	l.requests[ip] = append(valid, now)
	return true
}

// Global signup rate limiter: 5 signups per hour per IP
var signupLimiter = newSignupRateLimiter(5, time.Hour)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (from proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func onboardSignupHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handlePreflight(w, r)
			return
		}
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// Rate limiting
		clientIP := getClientIP(r)
		if !signupLimiter.Allow(clientIP) {
			writeJSON(w, r, http.StatusTooManyRequests, apiErrorPayload{
				Error: apiError{
					Code:              "rate_limit_exceeded",
					Message:           "Too many signup requests. Please try again later.",
					RequestID:         requestIDFrom(r.Context()),
					RetryAfterSeconds: 3600,
				},
			})
			return
		}

		// Parse request body
		var req signupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
			return
		}
		defer r.Body.Close()

		// Validate request
		if err := validateSignup(req); err != nil {
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Normalize email
		email := strings.ToLower(strings.TrimSpace(req.Email))
		companyName := strings.TrimSpace(req.CompanyName)

		// Check for duplicate email
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		exists, err := store.checkTenantEmail(ctx, email)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}
		if exists {
			writeError(w, r, http.StatusConflict, fmt.Errorf("email already registered"))
			return
		}

		// Determine plan (default to trial)
		plan := "trial"
		if req.Plan != "" {
			plan = strings.ToLower(req.Plan)
		}
		if plan != "trial" && plan != "starter" && plan != "growth" && plan != "enterprise" {
			plan = "trial"
		}

		// Create tenant with API key
		result, err := store.createTenantForSignup(ctx, tenantSignupParams{
			Email:       email,
			CompanyName: companyName,
			Plan:        plan,
			SignupIP:    clientIP,
			Source:      "web_signup",
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create account"))
			return
		}

		// Return success response with API key
		resp := signupResponse{
			Success:    true,
			TenantID:   result.TenantID.String(),
			TenantSlug: result.TenantSlug,
			StreamID:   result.StreamID.String(),
			APIKey:     result.APIKey,
			Message:    "Account created successfully. Save your API key - it won't be shown again.",
			RequestID:  requestIDFrom(r.Context()),
		}
		writeJSON(w, r, http.StatusCreated, resp)
	}
}

func validateSignup(req signupRequest) error {
	email := strings.TrimSpace(req.Email)
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !validateEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	companyName := strings.TrimSpace(req.CompanyName)
	if companyName == "" {
		return fmt.Errorf("company_name is required")
	}
	if len(companyName) < 2 {
		return fmt.Errorf("company_name must be at least 2 characters")
	}
	if len(companyName) > 100 {
		return fmt.Errorf("company_name must be less than 100 characters")
	}

	return nil
}

// Tenant signup params
type tenantSignupParams struct {
	Email       string
	CompanyName string
	Plan        string
	SignupIP    string
	Source      string
}

// Store methods for onboarding

func (s *store) checkTenantEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tenants WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *store) createTenantForSignup(ctx context.Context, params tenantSignupParams) (*tenantCreateResult, error) {
	// Use the existing createTenantWithKey method with signup-specific defaults
	return s.createTenantWithKey(ctx, tenantCreateParams{
		Name:                params.CompanyName,
		Slug:                slugify(params.CompanyName),
		Plan:                params.Plan,
		StreamSlug:          "default",
		StreamType:          "logs",
		StreamDescription:   "Default stream for anomaly detection",
		StreamRetentionDays: 90,
		KeyRole:             "admin",
		KeyName:             "default-key",
		KeyRateLimit:        60,
		TenantRateLimit:     120,
		DefaultBaseline:     400,
		DefaultWindow:       50,
		DefaultHop:          10,
		NCDThreshold:        0.3,
		PValueThreshold:     0.05,
		PermutationCount:    1000,
		DefaultCompressor:   "zstd",
		Seed:                time.Now().UnixNano(),
	})
}

// Admin endpoints for tenant management

type tenantListItem struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Email     string     `json:"email"`
	Plan      string     `json:"plan"`
	CreatedAt time.Time  `json:"created_at"`
	Status    string     `json:"status"`
}

type tenantListResponse struct {
	Tenants []tenantListItem `json:"tenants"`
	Total   int              `json:"total"`
}

func adminTenantsHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// Check admin key
		adminKey := r.Header.Get("X-Admin-Key")
		expectedKey := env("ADMIN_KEY", "")
		if expectedKey == "" || adminKey != expectedKey {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid admin key"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		rows, err := store.pool.Query(ctx, `
			SELECT id, name, slug, COALESCE(email, ''), plan, created_at,
			       CASE WHEN verified_at IS NOT NULL THEN 'verified' ELSE 'pending' END as status
			FROM tenants
			ORDER BY created_at DESC
			LIMIT 100`)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		var tenants []tenantListItem
		for rows.Next() {
			var t tenantListItem
			if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Email, &t.Plan, &t.CreatedAt, &t.Status); err != nil {
				writeError(w, r, http.StatusInternalServerError, err)
				return
			}
			tenants = append(tenants, t)
		}
		if err := rows.Err(); err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, tenantListResponse{
			Tenants: tenants,
			Total:   len(tenants),
		})
	}
}

type usageMetricsResponse struct {
	TenantID     string `json:"tenant_id"`
	EventCount   int64  `json:"event_count"`
	AnomalyCount int64  `json:"anomaly_count"`
	APIRequests  int64  `json:"api_requests"`
	Period       string `json:"period"`
}

func adminTenantUsageHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		// Check admin key
		adminKey := r.Header.Get("X-Admin-Key")
		expectedKey := env("ADMIN_KEY", "")
		if expectedKey == "" || adminKey != expectedKey {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid admin key"))
			return
		}

		// Extract tenant ID from path
		path := strings.TrimPrefix(r.URL.Path, "/v1/admin/tenants/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 || parts[1] != "usage" {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("not found"))
			return
		}

		tenantID, err := uuid.Parse(parts[0])
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid tenant id"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Get anomaly count for last 30 days
		var anomalyCount int64
		err = store.pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM anomalies
			WHERE tenant_id = $1 AND detected_at > NOW() - INTERVAL '30 days'`,
			tenantID).Scan(&anomalyCount)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Get batch count (proxy for API requests)
		var batchCount int64
		err = store.pool.QueryRow(ctx, `
			SELECT COUNT(*) FROM ingest_batches
			WHERE tenant_id = $1 AND created_at > NOW() - INTERVAL '30 days'`,
			tenantID).Scan(&batchCount)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, r, http.StatusOK, usageMetricsResponse{
			TenantID:     tenantID.String(),
			EventCount:   batchCount * 100, // Rough estimate
			AnomalyCount: anomalyCount,
			APIRequests:  batchCount,
			Period:       "last_30_days",
		})
	}
}
