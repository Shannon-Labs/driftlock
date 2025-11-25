package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const tenantContextKey ctxKey = "tenant"

func tenantFromContext(ctx context.Context) (tenantContext, bool) {
	tc, ok := ctx.Value(tenantContextKey).(tenantContext)
	return tc, ok
}

var firebaseAuth *auth.Client

func initFirebaseAuth() error {
	if os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY") == "" {
		return fmt.Errorf("FIREBASE_SERVICE_ACCOUNT_KEY not set")
	}

	// Decode if base64 encoded (common for env vars)
	keyJSON := []byte(os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"))
	if decoded, err := base64.StdEncoding.DecodeString(os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")); err == nil {
		keyJSON = decoded
	}

	opt := option.WithCredentialsJSON(keyJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error getting auth client: %v", err)
	}

	firebaseAuth = client
	return nil
}

type authUser struct {
	UID   string
	Email string
}

func verifyFirebaseToken(ctx context.Context, tokenString string) (*authUser, error) {
	if firebaseAuth == nil {
		return nil, fmt.Errorf("firebase auth not initialized")
	}

	token, err := firebaseAuth.VerifyIDToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	email, _ := token.Claims["email"].(string)
	return &authUser{
		UID:   token.UID,
		Email: email,
	}, nil
}

// withFirebaseAuth middleware verifies the Firebase ID token
// It attaches the user to the context
func withFirebaseAuth(store *store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("missing bearer token"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		user, err := verifyFirebaseToken(r.Context(), tokenString)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid token: %v", err))
			return
		}

		// Find tenant by email
		// Note: This assumes 1:1 mapping for now, or 1st tenant found
		// In a real multi-user system, we'd have a users table linking to tenants.
		// For MVP, we look up tenant by email.
		tenantID, err := store.tenantIDByEmail(r.Context(), user.Email)
		if err != nil {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("no tenant found for email %s", user.Email))
			return
		}

		// Inject tenant context (reusing existing context structure if possible, or creating new)
		// We need to fetch the API key to fully populate tenantContext if downstream needs it,
		// but for dashboard actions, just TenantID is often enough.
		// Let's fetch the full record to be safe.

		// Fetch a valid API key for this tenant to populate the context fully
		apiKey, err := store.primaryAPIKey(r.Context(), tenantID)
		if err != nil {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("tenant has no active api keys"))
			return
		}

		// Populate context
		tc := tenantContext{
			Tenant: tenantRecord{ID: tenantID, Name: user.Email, Email: user.Email}, // Partial record
			Key:    apiKey,
		}

		ctx := context.WithValue(r.Context(), tenantContextKey, tc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withAuth(store *store, limiter *tenantRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey == "" {
			// Try Bearer token
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("missing api key"))
			return
		}

		// Parse Key ID from "dlk_<uuid>.<secret>"
		if !strings.HasPrefix(apiKey, "dlk_") {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key format"))
			return
		}
		parts := strings.Split(strings.TrimPrefix(apiKey, "dlk_"), ".")
		if len(parts) != 2 {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key format"))
			return
		}
		keyID, err := uuid.Parse(parts[0])
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key id"))
			return
		}

		// Lookup key
		rec, err := store.resolveAPIKey(r.Context(), keyID)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key"))
			return
		}

		// Verify secret (hash stored over the full key string)
		candidateKey := fmt.Sprintf("dlk_%s.%s", keyID.String(), parts[1])
		if !verifyAPIKey(rec.KeyHash, candidateKey) {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key secret"))
			return
		}

		// Populate context first so we can use it for rate limiting if needed
		tc := tenantContext{
			Tenant: tenantRecord{
				ID:                rec.TenantID,
				Name:              rec.TenantName,
				Slug:              rec.TenantSlug,
				DefaultCompressor: rec.DefaultCompressor,
				RateLimitRPS:      rec.TenantRateLimit,
			},
			Key: *rec,
		}

		// Rate limit check
		if limiter != nil {
			decision := limiter.Allow(tc)
			decorateRateLimitHeaders(w, decision)
			if !decision.Allowed {
				writeRateLimitExceeded(w, r, decision)
				return
			}
		} else if tc.Tenant.RateLimitRPS > 0 {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(tc.Tenant.RateLimitRPS))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(tc.Tenant.RateLimitRPS))
		}

		ctx := context.WithValue(r.Context(), tenantContextKey, tc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decorateRateLimitHeaders(w http.ResponseWriter, decision rateLimitDecision) {
	if decision.TenantLimit > 0 {
		remaining := decision.TenantRemaining
		if decision.KeyLimit > 0 && decision.KeyRemaining < remaining {
			remaining = decision.KeyRemaining
		}
		if remaining < 0 {
			remaining = 0
		}
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(decision.TenantLimit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	}
	if decision.RetryAfter > 0 {
		seconds := int(math.Ceil(decision.RetryAfter.Seconds()))
		if seconds < 1 {
			seconds = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(seconds))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(decision.RetryAfter).Unix(), 10))
	}
}

func writeRateLimitExceeded(w http.ResponseWriter, r *http.Request, decision rateLimitDecision) {
	seconds := int(math.Ceil(decision.RetryAfter.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	logError(
		r,
		requestIDFrom(r.Context()),
		"rate_limit_exceeded",
		fmt.Sprintf("http_status=%d", http.StatusTooManyRequests),
		fmt.Errorf("rate limit exceeded"),
	)
	decorateRateLimitHeaders(w, decision)
	payload := apiErrorPayload{
		Error: apiError{
			Code:              "rate_limit_exceeded",
			Message:           "rate limit exceeded",
			RequestID:         requestIDFrom(r.Context()),
			RetryAfterSeconds: seconds,
		},
	}
	writeJSON(w, r, http.StatusTooManyRequests, payload)
}
