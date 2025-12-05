package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
	"github.com/google/uuid"
)

const (
	// Demo runs with a much smaller baseline/window so anomaly detection works with the 50-event cap.
	// This is intentionally lower than production defaults to keep demo latency low while still showing signal.
	demoBaselineSize = 40
	demoWindowSize   = 10
	demoHopSize      = 5
)

// Demo endpoint rate limiter - 10 requests/min per IP
type demoRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newDemoRateLimiter() *demoRateLimiter {
	return &demoRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    10,
		window:   time.Minute,
	}
}

func (d *demoRateLimiter) allow(ip string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-d.window)

	// Clean old requests
	var recent []time.Time
	for _, t := range d.requests[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	d.requests[ip] = recent

	if len(recent) >= d.limit {
		return false
	}

	d.requests[ip] = append(d.requests[ip], now)
	return true
}

func (d *demoRateLimiter) remaining(ip string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-d.window)

	var count int
	for _, t := range d.requests[ip] {
		if t.After(cutoff) {
			count++
		}
	}

	remaining := d.limit - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Periodically clean up old entries
func (d *demoRateLimiter) cleanup() {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := time.Now().Add(-d.window * 2)
	for ip, times := range d.requests {
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}
		if len(recent) == 0 {
			delete(d.requests, ip)
		} else {
			d.requests[ip] = recent
		}
	}
}

const maxDemoEvents = 50

type demoDetectRequest struct {
	Events         []json.RawMessage `json:"events"`
	ConfigOverride *configOverride   `json:"config_override,omitempty"`
}

type demoDetectResponse struct {
	Success        bool            `json:"success"`
	TotalEvents    int             `json:"total_events"`
	AnomalyCount   int             `json:"anomaly_count"`
	ProcessingTime string          `json:"processing_time"`
	CompressionAlg string          `json:"compression_algo"`
	Anomalies      []anomalyOutput `json:"anomalies"`
	RequestID      string          `json:"request_id"`
	Demo           demoInfo        `json:"demo"`
}

type demoInfo struct {
	Message         string `json:"message"`
	RemainingCalls  int    `json:"remaining_calls"`
	LimitPerMinute  int    `json:"limit_per_minute"`
	MaxEventsPerReq int    `json:"max_events_per_request"`
	SignupURL       string `json:"signup_url"`
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for Cloud Run, load balancers, etc.)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the chain
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	addr := r.RemoteAddr
	// Strip port if present
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

func demoDetectHandler(cfg config, limiter *demoRateLimiter) http.HandlerFunc {
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
					Message:           "Demo rate limit exceeded. Sign up for unlimited access.",
					RequestID:         requestIDFrom(r.Context()),
					RetryAfterSeconds: 60,
				},
			}
			writeJSON(w, r, http.StatusTooManyRequests, resp)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, cfg.MaxBodyBytes))
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("unable to read body: %w", err))
			return
		}
		defer r.Body.Close()

		var payload demoDetectRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
			return
		}

		if len(payload.Events) == 0 {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("events required"))
			return
		}

		if len(payload.Events) > maxDemoEvents {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("demo limited to %d events per request (got %d). Sign up for unlimited access", maxDemoEvents, len(payload.Events)))
			return
		}

		for idx, ev := range payload.Events {
			if len(bytes.TrimSpace(ev)) == 0 || bytes.Equal(bytes.TrimSpace(ev), []byte("null")) {
				writeError(w, r, http.StatusBadRequest, fmt.Errorf("event %d is empty", idx))
				return
			}
		}

		// Build detection settings with defaults
	plan := buildDemoDetectionSettings(cfg, payload.ConfigOverride)

	usedAlgo := plan.CompressionAlgorithm
	if usedAlgo == "openzl" && !driftlockcbad.HasOpenZL() {
		usedAlgo = "zstd"
		}

		detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
			BaselineSize:         plan.BaselineSize,
			WindowSize:           plan.WindowSize,
			HopSize:              plan.HopSize,
			MaxCapacity:          plan.BaselineSize + 4*plan.WindowSize + 1024,
			PValueThreshold:      plan.PValueThreshold,
			NCDThreshold:         plan.NCDThreshold,
			PermutationCount:     plan.PermutationCount,
			Seed:                 plan.Seed,
			CompressionAlgorithm: usedAlgo,
		})
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		defer detector.Close()

		requestCounter.Inc()
		start := time.Now()

		// SHA-141: Use default tokenizer for demo (all patterns enabled)
		tokenizer := driftlockcbad.GetTokenizer(driftlockcbad.DefaultTokenizerConfig())

		anomalies, _, err := runDetectionWithRecovery(r.Context(), detector, payload.Events, tokenizer)
		if err != nil {
			// Distinguish between user errors and internal FFI errors
			if errors.Is(err, ErrCBADPanic) || errors.Is(err, ErrCBADTimeout) {
				log.Printf("CBAD FFI error in demo: %v", err)
				writeError(w, r, http.StatusInternalServerError, fmt.Errorf("detection service temporarily unavailable"))
				return
			}
			writeError(w, r, http.StatusBadRequest, err)
			return
		}

		// Generate temporary IDs for demo anomalies (not persisted)
		for i := range anomalies {
			anomalies[i].ID = uuid.New().String()
		}

		resp := demoDetectResponse{
			Success:        true,
			TotalEvents:    len(payload.Events),
			AnomalyCount:   len(anomalies),
			ProcessingTime: time.Since(start).String(),
			CompressionAlg: usedAlgo,
			Anomalies:      anomalies,
			RequestID:      requestIDFrom(r.Context()),
			Demo: demoInfo{
				Message:         "This is a demo response. Sign up for full access with persistence, history, and evidence bundles.",
				RemainingCalls:  limiter.remaining(clientIP),
				LimitPerMinute:  10,
				MaxEventsPerReq: maxDemoEvents,
				SignupURL:       "https://driftlock.net/#signup",
			},
		}

		writeJSON(w, r, http.StatusOK, resp)
		requestDuration.Observe(time.Since(start).Seconds())
	}
}

func buildDemoDetectionSettings(cfg config, override *configOverride) detectionPlan {
	plan := detectionPlan{
		BaselineSize:         demoBaselineSize,
		WindowSize:           demoWindowSize,
		HopSize:              demoHopSize,
		NCDThreshold:         cfg.NCDThreshold,
		PValueThreshold:      cfg.PValueThreshold,
		PermutationCount:     cfg.PermutationCount,
		CompressionAlgorithm: strings.ToLower(cfg.DefaultAlgo),
		Seed:                 cfg.Seed,
	}

	if override != nil {
		if override.BaselineSize != nil {
			plan.BaselineSize = *override.BaselineSize
		}
		if override.WindowSize != nil {
			plan.WindowSize = *override.WindowSize
		}
		if override.HopSize != nil {
			plan.HopSize = *override.HopSize
		}
		if override.NCDThreshold != nil {
			plan.NCDThreshold = *override.NCDThreshold
		}
		if override.PValueThreshold != nil {
			plan.PValueThreshold = *override.PValueThreshold
		}
		if override.PermutationCount != nil {
			plan.PermutationCount = *override.PermutationCount
		}
		if override.Compressor != nil && *override.Compressor != "" {
			plan.CompressionAlgorithm = strings.ToLower(*override.Compressor)
		}
	}

	return plan
}
