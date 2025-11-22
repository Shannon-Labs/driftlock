package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Rate limiter for signup endpoint (5 per hour per IP)
var signupLimiter = struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}{
	limiters: make(map[string]*rate.Limiter),
}

type onboardSignupRequest struct {
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Plan        string `json:"plan"`
	Source      string `json:"source,omitempty"`
}

type onboardSignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	APIKey  string `json:"api_key,omitempty"`
	Tenant  struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Slug      string    `json:"slug"`
		Plan      string    `json:"plan"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"tenant"`
}

func onboardSignupHandler(cfg config, store *store, emailer *emailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST"))
			return
		}
		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 64*1024) // defensive bound for signup payloads

		// Rate limit by IP
		ip := getClientIP(r)
		if !getSignupLimiter(ip).Allow() {
			writeError(w, r, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded: max 5 signups per hour per IP"))
			return
		}

		// Check for Firebase Auth token (optional for backward compatibility)
		var firebaseUID string
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") && firebaseAuth != nil {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := verifyFirebaseToken(r.Context(), tokenString)
			if err == nil {
				firebaseUID = user.UID
				// If Firebase token is valid, use the email from the token (more secure)
				// But we'll still accept email from body for backward compatibility
			}
		}

		var req onboardSignupRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
			return
		}

		// Validate input
		if err := validateSignup(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Check if email already exists
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		exists, err := store.checkTenantEmail(ctx, req.Email)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("database error: %w", err))
			return
		}
		if exists {
			writeError(w, r, http.StatusConflict, fmt.Errorf("email already registered"))
			return
		}

		// If Firebase UID provided, also check if that UID is already linked
		if firebaseUID != "" {
			exists, err := store.checkTenantFirebaseUID(ctx, firebaseUID)
			if err != nil {
				writeError(w, r, http.StatusInternalServerError, fmt.Errorf("database error: %w", err))
				return
			}
			if exists {
				writeError(w, r, http.StatusConflict, fmt.Errorf("firebase account already registered"))
				return
			}
		}

		// Create tenant with active status and API key immediately
		plan := req.Plan
		if plan == "" {
			plan = "trial"
		}

		result, err := store.createTenantWithKey(ctx, tenantCreateParams{
			Name:                req.CompanyName,
			Slug:                slugify(req.CompanyName),
			Plan:                plan,
			StreamSlug:          "default",
			StreamType:          "logs",
			StreamDescription:   "Default stream created during onboarding",
			StreamRetentionDays: 14,
			KeyRole:             "admin",
			KeyName:             "default-key",
			KeyRateLimit:        cfg.DefaultRateLimit(),
			TenantRateLimit:     cfg.DefaultRateLimit(),
			DefaultBaseline:     cfg.DefaultBaseline,
			DefaultWindow:       cfg.DefaultWindow,
			DefaultHop:          cfg.DefaultHop,
			NCDThreshold:        cfg.NCDThreshold,
			PValueThreshold:     cfg.PValueThreshold,
			PermutationCount:    cfg.PermutationCount,
			DefaultCompressor:   cfg.DefaultAlgo,
			Email:               req.Email,
			SignupIP:            ip,
			Source:              req.Source,
			Seed:                int64(cfg.Seed),
			Status:              "active",    // Explicitly active
			FirebaseUID:         firebaseUID, // Link to Firebase Auth user if provided
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create tenant: %w", err))
			return
		}

		// Send welcome email (async)
		if emailer != nil {
			go emailer.sendWelcomeEmail(req.Email, req.CompanyName, result.APIKey)
		}

		// Build response
		resp := onboardSignupResponse{
			Success: true,
			Message: "Signup successful! Your API key is included below.",
			APIKey:  result.APIKey,
		}
		resp.Tenant.ID = result.Tenant.ID.String()
		resp.Tenant.Name = result.Tenant.Name
		resp.Tenant.Slug = result.Tenant.Slug
		resp.Tenant.Plan = result.Tenant.Plan
		resp.Tenant.CreatedAt = result.Tenant.CreatedAt

		writeJSON(w, r, http.StatusCreated, resp)
	}
}

func verifyHandler(store *store, emailer *emailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verification is now optional/post-signup since we give access immediately.
		// We can keep this handler for future email verification flows if needed,
		// or just return a "Already Verified" page.
		// For now, leaving it as is, but it won't be the primary flow.
		writeError(w, r, http.StatusNotImplemented, fmt.Errorf("verification flow deprecated in favor of instant access"))
	}
}

func validateSignup(req *onboardSignupRequest) error {
	// Validate email
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}

	// Validate company name
	if req.CompanyName == "" {
		return fmt.Errorf("company_name is required")
	}
	if len(req.CompanyName) < 2 {
		return fmt.Errorf("company_name must be at least 2 characters")
	}
	if len(req.CompanyName) > 100 {
		return fmt.Errorf("company_name must be at most 100 characters")
	}

	// Validate plan
	if req.Plan != "" {
		validPlans := map[string]bool{
			"trial":   true,
			"starter": true,
			"growth":  true,
			"pilot":   true,
		}
		if !validPlans[req.Plan] {
			return fmt.Errorf("invalid plan: must be one of trial, starter, growth, pilot")
		}
	}

	return nil
}

func getSignupLimiter(ip string) *rate.Limiter {
	signupLimiter.mu.Lock()
	defer signupLimiter.mu.Unlock()

	if limiter, exists := signupLimiter.limiters[ip]; exists {
		return limiter
	}

	// 5 requests per hour
	limiter := rate.NewLimiter(rate.Every(time.Hour/5), 5)
	signupLimiter.limiters[ip] = limiter
	return limiter
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by load balancers/proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fall back to RemoteAddr
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return r.RemoteAddr
}
