package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/shannon-labs/driftlock/api-server/internal/errors"
	"golang.org/x/time/rate"
)

// RateLimiter implements per-IP rate limiting using token bucket algorithm
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

// NewRateLimiter creates a new rate limiter
// rate: requests per second
// burst: maximum burst size
func NewRateLimiter(rps int, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	// Start background cleanup goroutine
	go rl.cleanupStale()

	return rl
}

// getLimiter gets or creates a rate limiter for the given key (IP address)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// cleanupStale removes stale rate limiters to prevent memory leak
func (rl *RateLimiter) cleanupStale() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// Simple cleanup: just clear all limiters periodically
		// A more sophisticated approach would track last access time
		rl.limiters = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// Middleware returns an HTTP middleware that enforces rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP (consider X-Forwarded-For in production)
		ip := getClientIP(r)

		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			errors.WriteJSON(w, errors.ErrTooManyRequests.WithDetails(
				"Rate limit exceeded. Please try again later.",
			))
			return
		}

		next.ServeHTTP(w, r)
	})
}


// GlobalRateLimiter implements a global rate limit (not per-IP)
type GlobalRateLimiter struct {
	limiter *rate.Limiter
}

// NewGlobalRateLimiter creates a new global rate limiter
func NewGlobalRateLimiter(rps int, burst int) *GlobalRateLimiter {
	return &GlobalRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Middleware returns an HTTP middleware that enforces global rate limiting
func (grl *GlobalRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !grl.limiter.Allow() {
			errors.WriteJSON(w, errors.ErrTooManyRequests.WithDetails(
				"Server is experiencing high load. Please try again later.",
			))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// EndpointRateLimiter implements per-endpoint rate limiting
type EndpointRateLimiter struct {
	limiters map[string]*RateLimiter
	mu       sync.RWMutex
}

// NewEndpointRateLimiter creates a new endpoint-specific rate limiter
func NewEndpointRateLimiter() *EndpointRateLimiter {
	return &EndpointRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

// AddEndpoint adds a rate limit for a specific endpoint
func (erl *EndpointRateLimiter) AddEndpoint(path string, rps int, burst int) {
	erl.mu.Lock()
	defer erl.mu.Unlock()
	erl.limiters[path] = NewRateLimiter(rps, burst)
}

// Middleware returns an HTTP middleware that enforces endpoint-specific rate limiting
func (erl *EndpointRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		erl.mu.RLock()
		limiter, exists := erl.limiters[r.URL.Path]
		erl.mu.RUnlock()

		if exists {
			ip := getClientIP(r)
			rl := limiter.getLimiter(ip)

			if !rl.Allow() {
				errors.WriteJSON(w, errors.ErrTooManyRequests.WithDetails(
					"Rate limit exceeded for this endpoint. Please try again later.",
				))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
