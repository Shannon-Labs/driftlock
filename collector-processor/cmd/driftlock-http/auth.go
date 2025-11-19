package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// Define context key
const tenantContextKey ctxKey = "tenantContext"

// Helper to get tenant from context
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
		tenantID, err := store.tenantIDByEmail(r.Context(), user.Email)
		if err != nil {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("no tenant found for email %s", user.Email))
			return
		}

		// Fetch a valid API key for this tenant to populate the context fully
		apiKey, err := store.primaryAPIKey(r.Context(), tenantID)
		if err != nil {
			writeError(w, r, http.StatusForbidden, fmt.Errorf("tenant has no active api keys"))
			return
		}

		// Populate context
		tc := tenantContext{
			Tenant: tenantRecord{ID: tenantID, Name: user.Email, Email: user.Email},
			Key:    apiKey,
		}
		
		ctx := context.WithValue(r.Context(), tenantContextKey, tc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withAuth verifies X-Api-Key
func withAuth(store *store, limiter *tenantRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Api-Key")
		if key == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				key = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if key == "" {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("missing api key"))
			return
		}

		if !strings.HasPrefix(key, "dlk_") {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key format"))
			return
		}

		parts := strings.Split(key, ".")
		if len(parts) != 2 {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key format"))
			return
		}

		idPart := strings.TrimPrefix(parts[0], "dlk_")
		keyID, err := uuid.Parse(idPart)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key id"))
			return
		}

		rec, err := store.resolveAPIKey(r.Context(), keyID)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key"))
			return
		}

		if !verifyAPIKey(rec.KeyHash, key) {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key"))
			return
		}

		tc := tenantContext{
			Tenant: tenantRecord{
				ID: rec.TenantID,
				Name: rec.TenantName,
				Slug: rec.TenantSlug,
				Plan: rec.Plan,
				DefaultCompressor: rec.DefaultCompressor,
				RateLimitRPS: rec.TenantRateLimit,
			},
			Key: *rec,
		}

		if allowed, wait := limiter.Allow(tc); !allowed {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(wait.Seconds())))
			writeError(w, r, http.StatusTooManyRequests, fmt.Errorf("rate limit exceeded"))
			return
		}

		ctx := context.WithValue(r.Context(), tenantContextKey, tc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
