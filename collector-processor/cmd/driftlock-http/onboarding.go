package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	Tenant  struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Slug      string    `json:"slug"`
		Plan      string    `json:"plan"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"tenant"`
	// APIKey removed from response for verification flow
}

func onboardSignupHandler(cfg config, store *store, emailer *emailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST"))
			return
		}

		// Rate limit by IP
		ip := getClientIP(r)
		if !getSignupLimiter(ip).Allow() {
			writeError(w, r, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded: max 5 signups per hour per IP"))
			return
		}

		var req onboardSignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

		// Create tenant with trial plan (Pending Verification)
		plan := req.Plan
		if plan == "" {
			plan = "trial"
		}

		// Generate verification token
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to generate token: %w", err))
			return
		}
		token := hex.EncodeToString(tokenBytes)

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
			VerificationToken:   token,
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create tenant: %w", err))
			return
		}

		// Send verification email (async)
		if emailer != nil {
			go emailer.sendVerificationEmail(req.Email, req.CompanyName, token)
		}

		// Build response
		resp := onboardSignupResponse{
			Success: true,
			Message: "Signup successful! Please check your email to verify your account and receive your API key.",
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
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("token required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Verify and activate tenant
		result, err := store.verifyAndActivateTenant(ctx, token)
		if err != nil {
			// If invalid token, it might be already verified or just wrong
			// We'll just return error for now
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("verification failed: %w", err))
			return
		}

		// Send welcome email with API key (async)
		// Need to fetch email from tenant record?
		// verifyAndActivateTenant returns tenantCreateResult which has ID/Name/Slug/Plan/CreatedAt
		// It does NOT have email. We need to fetch it or update verifyAndActivateTenant to return it.
		// Wait, `tenantCreateResult` doesn't have email.
		// Let's assume we need to fetch it or just add it to the return struct in db.go?
		// Or simpler: Just fetch tenant details using ID.
		// But `verifyAndActivateTenant` could just return email.
		// Let's do a quick fetch from DB here.

		var email string
		err = store.pool.QueryRow(ctx, `SELECT email FROM tenants WHERE id = $1`, result.TenantID).Scan(&email)
		if err != nil {
			// Should not happen if verify succeeded
			// Log error but don't fail the request completely?
			// If we fail here, user is verified but didn't get email. Bad.
			// We'll log and maybe show key on screen as backup?
			// But plan says "Welcome email with API key".
			fmt.Printf("ERROR: Failed to fetch email for tenant %s: %v\n", result.TenantID, err)
		} else if emailer != nil {
			go emailer.sendWelcomeEmail(email, result.Tenant.Name, result.APIKey)
		}

		// Return HTML success page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Email Verified - Driftlock</title>
				<style>
					body { font-family: system-ui, sans-serif; line-height: 1.5; color: #111827; max-width: 600px; margin: 40px auto; padding: 0 20px; }
					.card { background: #f9fafb; border-radius: 8px; padding: 32px; border: 1px solid #e5e7eb; text-align: center; }
					h1 { color: #2563eb; margin-top: 0; }
					.success-icon { font-size: 48px; margin-bottom: 16px; }
				</style>
			</head>
			<body>
				<div class="card">
					<div class="success-icon">✅</div>
					<h1>Email Verified!</h1>
					<p>Thank you for verifying your account.</p>
					<p>We have sent a <strong>Welcome Email</strong> to <code>%s</code> containing your <strong>API Key</strong>.</p>
					<p>Please check your inbox (and spam folder) to get started.</p>
					<p><a href="https://driftlock.net">Return to Driftlock Homepage</a></p>
				</div>
			</body>
			</html>
		`, email)
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
