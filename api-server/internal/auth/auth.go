package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Shannon-Labs/driftlock/api-server/internal/errors"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
	"github.com/Shannon-Labs/driftlock/api-server/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	// UserContextKey for storing user info in request context
	UserContextKey = "user"

	// Default token expiration time
	DefaultTokenExpiration = 24 * time.Hour

	// Minimum API key length
	MinAPIKeyLength = 32
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles authentication and authorization
type AuthService struct {
	storage    *storage.Storage
	jwtSecret  []byte
	tokenExpiry time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(storage *storage.Storage, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	if tokenExpiry == 0 {
		tokenExpiry = DefaultTokenExpiration
	}
	
	return &AuthService{
		storage:     storage,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

// GenerateAPIKey creates a new API key
func (s *AuthService) GenerateAPIKey(createReq *models.APIKeyCreate) (*models.APIKeyResponse, error) {
	// Generate a random API key
	rawKey, err := s.generateSecureKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash API key: %w", err)
	}

	// Create API key record
	apiKey := &models.APIKey{
		ID:                 s.generateSecureID(),
		KeyHash:            string(hashedKey),
		Name:               createReq.Name,
		Description:        createReq.Description,
		Role:               createReq.Role,
		Scopes:             createReq.Scopes,
		CreatedBy:          createReq.CreatedBy,
		ExpiresAt:          createReq.ExpiresAt,
		IsActive:           true,
		RateLimitPerMinute: createReq.RateLimitPerMinute,
	}

	// Store in database (we'll need to add this to storage interface)
	// For now, we'll just return the key since we don't have storage methods for API keys yet

	return &models.APIKeyResponse{
		APIKey: *apiKey,
		Key:    rawKey,
	}, nil
}

// ValidateAPIKey validates an API key and returns user info if valid
func (s *AuthService) ValidateAPIKey(key string) (*models.APIKey, error) {
	// In a real implementation, we would query the database to get the API key by its ID
	// and then verify the hash
	
	// For now, we'll implement a basic validation
	if len(key) < MinAPIKeyLength {
		return nil, fmt.Errorf("invalid API key format")
	}

	// This is where we'd query the database to find the key by ID and verify
	// Since we don't have the storage methods for this yet, we'll return a mock response
	// In a real implementation, this would look something like:
	//
	// apiKey, err := s.storage.GetAPIKeyByID(ctx, keyID)
	// if err != nil {
	//     return nil, err
	// }
	// 
	// if !apiKey.IsActive {
	//     return nil, fmt.Errorf("API key is inactive")
	// }
	// 
	// if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
	//     return nil, fmt.Errorf("API key has expired")
	// }
	// 
	// err = bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(key))
	// if err != nil {
	//     return nil, fmt.Errorf("invalid API key")
	// }

	return &models.APIKey{
		ID:       key[:8], // Use first 8 chars as ID for demo
		Name:     "Demo Key",
		Role:     "admin",
		IsActive: true,
	}, nil
}

// GenerateJWT creates a new JWT token
func (s *AuthService) GenerateJWT(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseJWT parses and validates a JWT token
func (s *AuthService) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// generateSecureKey generates a cryptographically secure random key
func (s *AuthService) generateSecureKey() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateSecureID generates a unique ID for API keys
func (s *AuthService) generateSecureID() string {
	bytes := make([]byte, 16) // 128 bits
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// APIKeyMiddleware provides API key authentication middleware
func (s *AuthService) APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for API key in header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Authorization header required"))
			return
		}

		// Support both "Bearer <token>" and "Api-Key <key>" formats
		var apiKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "Api-Key ") {
			apiKey = strings.TrimPrefix(authHeader, "Api-Key ")
		} else {
			// Also check for a custom header
			apiKey = r.Header.Get("X-API-Key")
			if apiKey == "" {
				errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Valid API key required"))
				return
			}
		}

		// Validate the API key
		keyInfo, err := s.ValidateAPIKey(apiKey)
		if err != nil {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Invalid API key"))
			return
		}

		if !keyInfo.IsActive {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("API key is inactive"))
			return
		}

		if keyInfo.ExpiresAt != nil && time.Now().After(*keyInfo.ExpiresAt) {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("API key has expired"))
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, keyInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWTMiddleware provides JWT authentication middleware
func (s *AuthService) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Authorization header required"))
			return
		}

		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Authorization header must use Bearer scheme"))
			return
		}

		claims, err := s.ParseJWT(tokenString)
		if err != nil {
			errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Invalid or expired token"))
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves user info from request context
func GetUserFromContext(ctx context.Context) *models.APIKey {
	if user, ok := ctx.Value(UserContextKey).(*models.APIKey); ok {
		return user
	}
	if claims, ok := ctx.Value(UserContextKey).(*Claims); ok {
		return &models.APIKey{
			ID:   claims.UserID,
			Role: claims.Role,
		}
	}
	return nil
}

// RequireRole middleware checks if the authenticated user has the required role
func (s *AuthService) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Authentication required"))
				return
			}

			if user.Role != requiredRole && requiredRole != "admin" && user.Role != "admin" {
				errors.WriteJSON(w, errors.ErrForbidden.WithDetails("Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole middleware checks if the authenticated user has any of the required roles
func (s *AuthService) RequireAnyRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				errors.WriteJSON(w, errors.ErrUnauthorized.WithDetails("Authentication required"))
				return
			}

			hasRole := false
			for _, role := range requiredRoles {
				if user.Role == role || user.Role == "admin" {
					hasRole = true
					break
				}
			}

			if !hasRole {
				errors.WriteJSON(w, errors.ErrForbidden.WithDetails("Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}