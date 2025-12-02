package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
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
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	APIKey           string `json:"api_key,omitempty"`
	PendingVerify    bool   `json:"pending_verification,omitempty"`
	VerificationSent bool   `json:"verification_sent,omitempty"`
	Tenant           struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Slug      string    `json:"slug"`
		Plan      string    `json:"plan"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"tenant"`
}

type verifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	APIKey  string `json:"api_key,omitempty"`
	Tenant  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		Plan string `json:"plan"`
	} `json:"tenant,omitempty"`
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

		// Generate verification token
		verificationToken, err := generateVerificationToken()
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to generate verification token"))
			return
		}

		// Create pending tenant (no API key until verified)
		plan := req.Plan
		if plan == "" {
			plan = "trial"
		}

		result, err := store.createPendingTenant(ctx, tenantCreateParams{
			Name:                req.CompanyName,
			Slug:                slugify(req.CompanyName),
			Plan:                plan,
			StreamSlug:          "default",
			StreamType:          "logs",
			StreamDescription:   "Default stream created during onboarding",
			StreamRetentionDays: 14,
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
			VerificationToken:   verificationToken,
			FirebaseUID:         firebaseUID, // Link to Firebase Auth user if provided
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create tenant: %w", err))
			return
		}

		// Send verification email (async)
		if emailer != nil {
			go emailer.sendVerificationEmail(req.Email, req.CompanyName, verificationToken)
		}

		// Build response - no API key yet, pending verification
		resp := onboardSignupResponse{
			Success:          true,
			Message:          "Please check your email to verify your account and receive your API key.",
			PendingVerify:    true,
			VerificationSent: true,
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
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use GET"))
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("verification token required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Verify and activate the tenant, generating the API key
		result, err := store.verifyAndActivateTenant(ctx, token)
		if err != nil {
			if strings.Contains(err.Error(), "invalid or expired") {
				writeError(w, r, http.StatusBadRequest, err)
				return
			}
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("verification failed: %w", err))
			return
		}

		// Send welcome email with API key
		if emailer != nil {
			go emailer.sendWelcomeEmail(result.Tenant.Name, result.Tenant.Name, result.APIKey)
		}

		// Return the API key
		resp := verifyResponse{
			Success: true,
			Message: "Email verified! Your API key is below. Store it securely - it won't be shown again.",
			APIKey:  result.APIKey,
		}
		resp.Tenant.ID = result.Tenant.ID.String()
		resp.Tenant.Name = result.Tenant.Name
		resp.Tenant.Slug = result.Tenant.Slug
		resp.Tenant.Plan = result.Tenant.Plan

		writeJSON(w, r, http.StatusOK, resp)
	}
}

// generateVerificationToken creates a secure random token for email verification
func generateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
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

	// 5 requests per hour (or 1000 in dev mode)
	limit := rate.Limit(5) / rate.Limit(time.Hour.Seconds())
	burst := 5
	if os.Getenv("DRIFTLOCK_DEV_MODE") == "true" {
		limit = rate.Inf
		burst = 1000
	}

	limiter := rate.NewLimiter(limit, burst)
	signupLimiter.limiters[ip] = limiter
	return limiter
}

// getClientIP is defined in demo.go
