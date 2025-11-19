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
