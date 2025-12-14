package main

// Legacy onboarding implementation (Go).
// Superseded by the Rust handler in crates/driftlock-api/src/routes/onboarding.rs.
// Kept for historical reference only.

/*
import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// Rate limiter for signup endpoint (5 per hour per IP)
var signupLimiter = struct {
	mu      sync.Mutex
	limiter map[string]*rate.Limiter
}{
	limiter: make(map[string]*rate.Limiter),
}

func onboardSignupHandler(store *store) http.HandlerFunc {
	type request struct {
		Email      string `json:"email"`
		Company    string `json:"company_name"`
		Plan       string `json:"plan"`
		Source     string `json:"source"`
	}

	type response struct {
		Success bool                   `json:"success"`
		Tenant  map[string]interface{} `json:"tenant"`
		APIKey  string                 `json:"api_key"`
		Message string                 `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("use POST"))
			return
		}

		// Rate limit by IP
		ip := getClientIP(r)
		if !getSignupLimiter(ip).Allow() {
			writeError(w, r, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded"))
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Validate
		if err := validateSignup(req); err != nil {
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Check if email already exists
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		exists, err := store.checkTenantEmail(ctx, req.Email)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		if exists {
			writeError(w, r, http.StatusConflict, fmt.Errorf("email already registered"))
			return
		}

		// Create tenant (dev mode for now, add license check later)
		cfg := loadConfig()
		cfg.DevMode = true // Skip license for onboarding

		result, err := createTenant(ctx, cfg, tenantCreateParams{
			Name:                req.Company,
			Slug:                slugify(req.Company),
			Plan:                "free", // Set to "free" plan
			StreamSlug:          "default",
			StreamType:          "logs",
			StreamDescription:   "Default stream from onboarding",
			StreamRetentionDays: 14,
			KeyRole:             "admin",
			KeyName:             "onboarding-key",
			TenantRateLimit:     60, // Default rate limit for free tier
			DefaultBaseline:     cfg.DefaultBaseline,
			DefaultWindow:       cfg.DefaultWindow,
			DefaultHop:          cfg.DefaultHop,
			NCDThreshold:        cfg.NCDThreshold,
			PValueThreshold:     cfg.PValueThreshold,
			PermutationCount:    cfg.PermutationCount,
			DefaultCompressor:   cfg.DefaultAlgo,
			Email:               req.Email,
			SignupIP:            ip,
			Source:              "firebase_frontend", // Explicitly set the source
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Send welcome email (async)
		go sendWelcomeEmail(req.Email, req.Company, result.APIKey)

		// Return response
		resp := response{
			Success: true,
			Tenant: map[string]interface{}{
				"id":         result.Tenant.ID,
				"name":       result.Tenant.Name,
				"slug":       result.Tenant.Slug,
				"plan":       result.Tenant.Plan,
				"created_at": result.Tenant.CreatedAt,
			},
			APIKey:  result.APIKey,
			Message: "Welcome to Driftlock! Check your email for next steps.",
		}

		writeJSON(w, r, http.StatusCreated, resp)
	}
}

func validateSignup(req request) error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	if len(req.Company) < 2 {
		return fmt.Errorf("company name too short")
	}
	validPlans := []string{"trial", "pilot", "demo"}
	for _, p := range validPlans {
		if req.Plan == p {
			return nil
		}
	}
	return fmt.Errorf("invalid plan")
}

func getSignupLimiter(ip string) *rate.Limiter {
	signupLimiter.mu.Lock()
	defer signupLimiter.mu.Unlock()

	if limiter, exists := signupLimiter.limiter[ip]; exists {
		return limiter
	}

	// 5 requests per hour
	limiter := rate.NewLimiter(rate.Every(time.Hour/5), 5)
	signupLimiter.limiter[ip] = limiter
	return limiter
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "-")
	s = strings.ReplaceAll(s, "@", "-")
	// Remove non-alphanumeric except dash
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func sendWelcomeEmail(email, company, apiKey string) {
	// TODO: Integrate SendGrid or similar
	log.Printf("Welcome email to %s for %s (API key: %s...)", email, company, apiKey[:12])
}

// Add to buildHTTPHandler in main.go:
// mux.HandleFunc("/v1/onboard/signup", onboardSignupHandler(store))
*/
