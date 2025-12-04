package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Waitlist rate limiter - 5 requests/hour per IP (generous for legitimate signups)
type waitlistRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newWaitlistRateLimiter() *waitlistRateLimiter {
	return &waitlistRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    5,
		window:   time.Hour,
	}
}

func (w *waitlistRateLimiter) allow(ip string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-w.window)

	// Clean old requests
	var recent []time.Time
	for _, t := range w.requests[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	w.requests[ip] = recent

	if len(recent) >= w.limit {
		return false
	}

	w.requests[ip] = append(w.requests[ip], now)
	return true
}

// Periodically clean up old entries
func (w *waitlistRateLimiter) cleanup() {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := time.Now().Add(-w.window * 2)
	for ip, times := range w.requests {
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}
		if len(recent) == 0 {
			delete(w.requests, ip)
		} else {
			w.requests[ip] = recent
		}
	}
}

// Simple email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type waitlistRequest struct {
	Email  string `json:"email"`
	Source string `json:"source,omitempty"` // Optional: "website", "api-docs", etc.
}

type waitlistResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func waitlistHandler(st *store, limiter *waitlistRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handlePreflight(w, r)
			return
		}
		if r.Method != http.MethodPost {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}

		clientIP := getClientIP(r)

		// Check rate limit
		if !limiter.allow(clientIP) {
			resp := apiErrorPayload{
				Error: apiError{
					Code:              "rate_limit_exceeded",
					Message:           "Too many waitlist requests. Please try again later.",
					RequestID:         requestIDFrom(r.Context()),
					RetryAfterSeconds: 3600, // 1 hour
				},
			}
			writeJSON(w, r, http.StatusTooManyRequests, resp)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, 1024)) // Small body limit
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("unable to read body: %w", err))
			return
		}
		defer r.Body.Close()

		var payload waitlistRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
			return
		}

		// Validate email
		email := strings.TrimSpace(strings.ToLower(payload.Email))
		if email == "" {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("email is required"))
			return
		}
		if !emailRegex.MatchString(email) {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid email format"))
			return
		}

		// Default source
		source := "website"
		if payload.Source != "" {
			source = strings.TrimSpace(payload.Source)
		}

		// Insert into waitlist (ignore duplicate)
		err = st.addToWaitlist(r.Context(), email, source, clientIP)
		if err != nil {
			// Check if it's a duplicate (unique constraint violation)
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				// Return success anyway - don't reveal if email exists
				writeJSON(w, r, http.StatusOK, waitlistResponse{
					Success: true,
					Message: "You're on the list! We'll notify you when we launch.",
				})
				return
			}
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("failed to add to waitlist"))
			return
		}

		writeJSON(w, r, http.StatusOK, waitlistResponse{
			Success: true,
			Message: "You're on the list! We'll notify you when we launch.",
		})
	}
}
