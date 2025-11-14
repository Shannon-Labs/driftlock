package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type config struct {
	MaxBodyBytes     int64
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	DefaultBaseline  int
	DefaultWindow    int
	DefaultHop       int
	DefaultAlgo      string
	PValueThreshold  float64
	NCDThreshold     float64
	PermutationCount int
	Seed             uint64
}

func loadConfig() config {
	return config{
		MaxBodyBytes:     int64(envInt("MAX_BODY_MB", 10)) * 1024 * 1024,
		ReadTimeout:      time.Duration(envInt("READ_TIMEOUT_SEC", 15)) * time.Second,
		WriteTimeout:     time.Duration(envInt("WRITE_TIMEOUT_SEC", 30)) * time.Second,
		IdleTimeout:      time.Duration(envInt("IDLE_TIMEOUT_SEC", 60)) * time.Second,
		DefaultBaseline:  envInt("DEFAULT_BASELINE", 400),
		DefaultWindow:    envInt("DEFAULT_WINDOW", 1),
		DefaultHop:       envInt("DEFAULT_HOP", 1),
		DefaultAlgo:      env("DEFAULT_ALGO", "zstd"),
		PValueThreshold:  envFloat("PVALUE_THRESHOLD", 0.05),
		NCDThreshold:     envFloat("NCD_THRESHOLD", 0.3),
		PermutationCount: envInt("PERMUTATION_COUNT", 1000),
		Seed:             envInt64("SEED", 42),
	}
}

type detectResponse struct {
	Success         bool            `json:"success"`
	TotalEvents     int             `json:"total_events,omitempty"`
	AnomalyCount    int             `json:"anomaly_count,omitempty"`
	ProcessingTime  string          `json:"processing_time,omitempty"`
	CompressionAlg  string          `json:"compression_algo,omitempty"`
	FallbackFromAlg string          `json:"fallback_from_algo,omitempty"`
	Anomalies       []anomalyOutput `json:"anomalies,omitempty"`
	RequestID       string          `json:"request_id"`
	Error           string          `json:"error,omitempty"`
}

type anomalyOutput struct {
	Index    int                           `json:"index"`
	Metrics  driftlockcbad.EnhancedMetrics `json:"metrics"`
	Event    json.RawMessage               `json:"event"`
	Why      string                        `json:"why"`
	Detected bool                          `json:"detected"`
}

func main() {
	cfg := loadConfig()
	registerMetrics()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/v1/detect", func(w http.ResponseWriter, r *http.Request) {
		detectHandler(w, r, cfg)
	})
	mux.Handle("/metrics", promhttp.Handler())

	addr := env("PORT", "8080")
	log.Printf("driftlock-http listening on :%s", addr)

	srv := &http.Server{
		Addr:         ":" + addr,
		Handler:      withCommon(withRequestContext(mux)),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Channel to receive errors
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to receive system signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-shutdown:
		log.Printf("Server received signal %v, beginning graceful shutdown", sig)

		// Create a context with timeout for shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
			// Force shutdown if graceful fails
			if err := srv.Close(); err != nil {
				log.Printf("Server close error: %v", err)
			}
		} else {
			log.Printf("Server gracefully shutdown")
		}
	}
}

type healthResponse struct {
	Success         bool     `json:"success"`
	RequestID       string   `json:"request_id"`
	Error           string   `json:"error,omitempty"`
	LibraryStatus   string   `json:"library_status"`
	Version         string   `json:"version,omitempty"`
	AvailableAlgos  []string `json:"available_algos,omitempty"`
	OpenZLAvailable bool     `json:"openzl_available"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		handlePreflight(w, r)
		return
	}

	openzlEnabled := driftlockcbad.HasOpenZL()

	resp := healthResponse{
		Success:         true,
		RequestID:       requestIDFrom(r.Context()),
		LibraryStatus:   "healthy",
		Version:         "1.0.0",
		AvailableAlgos:  []string{"zstd", "lz4", "gzip"},
		OpenZLAvailable: openzlEnabled,
	}

	if err := driftlockcbad.ValidateLibrary(); err != nil {
		resp.Success = false
		resp.LibraryStatus = "unhealthy"
		resp.Error = fmt.Sprintf("Library validation failed: %s", err.Error())
		writeJSON(w, r, http.StatusServiceUnavailable, resp)
		return
	}

	if openzlEnabled {
		resp.AvailableAlgos = append(resp.AvailableAlgos, "openzl")
	}

	writeJSON(w, r, http.StatusOK, resp)
}

func detectHandler(w http.ResponseWriter, r *http.Request, cfg config) {
	if r.Method == http.MethodOptions {
		handlePreflight(w, r)
		return
	}

	start := time.Now()
	// Body size cap
	r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxBodyBytes)

	algo := queryString(r, "algo", cfg.DefaultAlgo)
	fallbackFrom := ""
	requestedAlgo := strings.ToLower(algo)
	usedAlgo := requestedAlgo
	if usedAlgo == "openzl" && !driftlockcbad.HasOpenZL() {
		fallbackFrom = usedAlgo
		usedAlgo = "zstd"
	}
	format := queryString(r, "format", "")
	baseline := queryInt(r, "baseline", cfg.DefaultBaseline)
	window := queryInt(r, "window", cfg.DefaultWindow)
	hop := queryInt(r, "hop", cfg.DefaultHop)

	requestCounter.Inc()
	detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
		BaselineSize:         baseline,
		WindowSize:           window,
		HopSize:              hop,
		MaxCapacity:          baseline + 4*window + 1024,
		PValueThreshold:      cfg.PValueThreshold,
		NCDThreshold:         cfg.NCDThreshold,
		PermutationCount:     cfg.PermutationCount,
		Seed:                 cfg.Seed,
		CompressionAlgorithm: usedAlgo,
	})
	if err != nil {
		// Graceful fallback if openzl requested but not available
		if usedAlgo == "openzl" {
			fallbackFrom = usedAlgo
			usedAlgo = "zstd"
			var derr error
			detector, derr = driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
				BaselineSize:         baseline,
				WindowSize:           window,
				HopSize:              hop,
				MaxCapacity:          baseline + 4*window + 1024,
				PValueThreshold:      cfg.PValueThreshold,
				NCDThreshold:         cfg.NCDThreshold,
				PermutationCount:     cfg.PermutationCount,
				Seed:                 cfg.Seed,
				CompressionAlgorithm: usedAlgo,
			})
			if derr != nil {
				writeError(w, r, http.StatusInternalServerError, derr)
				return
			}
		} else {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
	}
	defer detector.Close()

	var anomalies []anomalyOutput
	var idx int

	body := r.Body
	defer body.Close()

	// Autodetect format if not specified
	if format == "" {
		// Peek first non-whitespace character
		buf := bufio.NewReader(body)
		first, _ := buf.Peek(1)
		if len(first) > 0 && first[0] == '[' {
			format = "json"
		} else {
			format = "ndjson"
		}
		body = io.NopCloser(buf)
	}

	switch strings.ToLower(format) {
	case "ndjson":
		rd := bufio.NewReader(body)
		for {
			line, err := rd.ReadBytes('\n')
			// Skip empty/whitespace-only lines
			trimmed := strings.TrimSpace(string(line))
			if len(trimmed) > 0 {
				// Validate that each NDJSON line is valid JSON
				var raw json.RawMessage
				if jerr := json.Unmarshal([]byte(trimmed), &raw); jerr != nil {
					writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid NDJSON at index %d: %s", idx, jerr.Error()))
					return
				}
				_ = processEvent(detector, raw, idx, &anomalies)
				idx++
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				writeError(w, r, http.StatusBadRequest, fmt.Errorf("read error at index %d: %s", idx, err.Error()))
				return
			}
		}
	case "json":
		var arr []json.RawMessage
		dec := json.NewDecoder(body)
		if err := dec.Decode(&arr); err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid JSON array: %s", err.Error()))
			return
		}
		for i, ev := range arr {
			_ = processEvent(detector, ev, i, &anomalies)
			idx++
		}
	default:
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("unsupported format %q", format))
		return
	}

	// Require at least one valid event
	if idx == 0 {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("empty request body"))
		return
	}

	resp := detectResponse{
		Success:         true,
		TotalEvents:     idx,
		AnomalyCount:    len(anomalies),
		ProcessingTime:  time.Since(start).String(),
		CompressionAlg:  usedAlgo,
		FallbackFromAlg: fallbackFrom,
		Anomalies:       anomalies,
		RequestID:       requestIDFrom(r.Context()),
	}
	writeJSON(w, r, http.StatusOK, resp)
	requestDuration.Observe(time.Since(start).Seconds())
}

func processEvent(detector *driftlockcbad.Detector, ev []byte, index int, sink *[]anomalyOutput) error {
	added, err := detector.AddData(ev)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}
	ready, err := detector.IsReady()
	if err != nil || !ready {
		return err
	}
	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		return err
	}
	if detected {
		why := metrics.GetDetailedExplanation()
		*sink = append(*sink, anomalyOutput{
			Index:    index,
			Metrics:  *metrics,
			Event:    json.RawMessage(append([]byte{}, ev...)),
			Why:      why,
			Detected: true,
		})
	}
	return nil
}

func withCommon(next http.Handler) http.Handler {
	allowed := parseAllowedOrigins(env("CORS_ALLOW_ORIGINS", "*"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin, allowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-Id")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}

		// Security headers
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// HSTS for HTTPS connections
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			handlePreflight(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

type ctxKey string

const requestIDKey ctxKey = "reqid"

func withRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// Log request start
		logRequest(r, reqID, "request_start", "")

		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Log request completion
		logRequest(r, reqID, "request_complete", fmt.Sprintf("status=%d", wrapped.statusCode))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logRequest(r *http.Request, requestID, event, details string) {
	logEntry := map[string]interface{}{
		"ts":         time.Now().Format(time.RFC3339Nano),
		"request_id": requestID,
		"method":     r.Method,
		"path":       r.URL.Path,
		"query":      r.URL.RawQuery,
		"remote":     r.RemoteAddr,
		"user_agent": r.Header.Get("User-Agent"),
		"event":      event,
	}
	if details != "" {
		logEntry["details"] = details
	}

	jsonLog, _ := json.Marshal(logEntry)
	log.Printf("%s", string(jsonLog))
}

func logError(r *http.Request, requestID, errType, details string, err error) {
	logEntry := map[string]interface{}{
		"ts":         time.Now().Format(time.RFC3339Nano),
		"request_id": requestID,
		"method":     r.Method,
		"path":       r.URL.Path,
		"remote":     r.RemoteAddr,
		"event":      "error",
		"error_type": errType,
		"details":    details,
	}
	if err != nil {
		logEntry["error"] = err.Error()
	}

	jsonLog, _ := json.Marshal(logEntry)
	log.Printf("%s", string(jsonLog))
}

func requestIDFrom(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, status int, err error) {
	reqID := requestIDFrom(r.Context())

	// Log the error with structured format
	errType := "unknown_error"
	switch status {
	case http.StatusBadRequest:
		errType = "bad_request"
	case http.StatusUnauthorized:
		errType = "unauthorized"
	case http.StatusForbidden:
		errType = "forbidden"
	case http.StatusNotFound:
		errType = "not_found"
	case http.StatusTooManyRequests:
		errType = "rate_limited"
	case http.StatusInternalServerError:
		errType = "internal_error"
	case http.StatusServiceUnavailable:
		errType = "service_unavailable"
	}

	logError(r, reqID, errType, fmt.Sprintf("http_status=%d", status), err)

	resp := detectResponse{
		Success:   false,
		RequestID: reqID,
		Error:     err.Error(),
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envInt64(key string, def int64) uint64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return uint64(n)
		}
	}
	return uint64(def)
}

func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return def
}

func queryString(r *http.Request, key, def string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return def
}

func queryInt(r *http.Request, key string, def int) int {
	if v := r.URL.Query().Get(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func parseAllowedOrigins(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func originAllowed(origin string, allow []string) bool {
	if origin == "" {
		return false
	}
	for _, a := range allow {
		if a == "*" || strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}

var (
	requestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "driftlock_http_requests_total",
		Help: "Total number of /v1/detect requests",
	})
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "driftlock_http_request_duration_seconds",
		Help:    "Duration of /v1/detect requests",
		Buckets: prometheus.DefBuckets,
	})
	registerMetricsOnce sync.Once
)

func registerMetrics() {
	registerMetricsOnce.Do(func() {
		prometheus.MustRegister(requestCounter, requestDuration)
	})
}
