package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	"github.com/Shannon-Labs/driftlock/collector-processor/internal/ai"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type config struct {
	MaxBodyBytes     int64
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	MaxEvents        int
	DefaultBaseline  int
	DefaultWindow    int
	DefaultHop       int
	DefaultAlgo      string
	PValueThreshold  float64
	NCDThreshold     float64
	PermutationCount int
	Seed             uint64
	RateLimitRPS     int
	QueueCapacity    int
	PreferOpenZL     bool
}

func loadConfig() config {
	return config{
		MaxBodyBytes:     int64(envInt("MAX_BODY_MB", 10)) * 1024 * 1024,
		ReadTimeout:      time.Duration(envInt("READ_TIMEOUT_SEC", 15)) * time.Second,
		WriteTimeout:     time.Duration(envInt("WRITE_TIMEOUT_SEC", 30)) * time.Second,
		IdleTimeout:      time.Duration(envInt("IDLE_TIMEOUT_SEC", 60)) * time.Second,
		MaxEvents:        envInt("MAX_EVENTS", 1000),
		DefaultBaseline:  envInt("DEFAULT_BASELINE", 400),
		DefaultWindow:    envInt("DEFAULT_WINDOW", 50),
		DefaultHop:       envInt("DEFAULT_HOP", 10),
		DefaultAlgo:      env("DEFAULT_ALGO", "zstd"),
		PValueThreshold:  envFloat("PVALUE_THRESHOLD", 0.05),
		NCDThreshold:     envFloat("NCD_THRESHOLD", 0.3),
		PermutationCount: envInt("PERMUTATION_COUNT", 1000),
		Seed:             envInt64("SEED", 42),
		RateLimitRPS:     envInt("RATE_LIMIT_RPS", 60),
		QueueCapacity:    envInt("QUEUE_CAPACITY", 512),
		PreferOpenZL:     envBool("PREFER_OPENZL", false),
	}
}

func (c config) DefaultRateLimit() int {
	if c.RateLimitRPS > 0 {
		return c.RateLimitRPS
	}
	return 60
}

type detectRequest struct {
	StreamID       string            `json:"stream_id"`
	Events         []json.RawMessage `json:"events"`
	ConfigOverride *configOverride   `json:"config_override,omitempty"`
}

type configOverride struct {
	BaselineSize     *int     `json:"baseline_size,omitempty"`
	WindowSize       *int     `json:"window_size,omitempty"`
	HopSize          *int     `json:"hop_size,omitempty"`
	NCDThreshold     *float64 `json:"ncd_threshold,omitempty"`
	PValueThreshold  *float64 `json:"p_value_threshold,omitempty"`
	PermutationCount *int     `json:"permutation_count,omitempty"`
	Compressor       *string  `json:"compressor,omitempty"`
}

type detectResponse struct {
	Success         bool            `json:"success"`
	BatchID         string          `json:"batch_id"`
	StreamID        string          `json:"stream_id"`
	TotalEvents     int             `json:"total_events"`
	AnomalyCount    int             `json:"anomaly_count"`
	ProcessingTime  string          `json:"processing_time"`
	CompressionAlg  string          `json:"compression_algo"`
	FallbackFromAlg string          `json:"fallback_from_algo,omitempty"`
	Anomalies       []anomalyOutput `json:"anomalies"`
	RequestID       string          `json:"request_id"`
}

type anomalyOutput struct {
	ID       string                        `json:"id"`
	Index    int                           `json:"index"`
	Metrics  driftlockcbad.EnhancedMetrics `json:"metrics"`
	Event    json.RawMessage               `json:"event"`
	Why      string                        `json:"why"`
	Detected bool                          `json:"detected"`
}

type anomalyListResponse struct {
	Anomalies     []anomalyListItem `json:"anomalies"`
	NextPageToken string            `json:"next_page_token,omitempty"`
	Total         int               `json:"total"`
}

type anomalyListItem struct {
	ID               string    `json:"id"`
	StreamID         string    `json:"stream_id"`
	NCD              float64   `json:"ncd"`
	CompressionRatio float64   `json:"compression_ratio"`
	EntropyChange    float64   `json:"entropy_change"`
	PValue           float64   `json:"p_value"`
	Confidence       float64   `json:"confidence"`
	Status           string    `json:"status"`
	Explanation      string    `json:"explanation"`
	DetectedAt       time.Time `json:"detected_at"`
}

type anomalyDetailResponse struct {
	ID               string               `json:"id"`
	StreamID         string               `json:"stream_id"`
	BatchID          string               `json:"batch_id,omitempty"`
	Status           string               `json:"status"`
	DetectedAt       time.Time            `json:"detected_at"`
	Explanation      string               `json:"explanation"`
	Metrics          anomalyDetailMetrics `json:"metrics"`
	Details          json.RawMessage      `json:"details,omitempty"`
	BaselineSnapshot json.RawMessage      `json:"baseline_snapshot,omitempty"`
	WindowSnapshot   json.RawMessage      `json:"window_snapshot,omitempty"`
	Evidence         []anomalyEvidence    `json:"evidence"`
}

type anomalyDetailMetrics struct {
	NCD              float64 `json:"ncd"`
	CompressionRatio float64 `json:"compression_ratio"`
	EntropyChange    float64 `json:"entropy_change"`
	PValue           float64 `json:"p_value"`
	Confidence       float64 `json:"confidence"`
}

type anomalyEvidence struct {
	Format    string    `json:"format"`
	URI       string    `json:"uri"`
	Checksum  string    `json:"checksum,omitempty"`
	SizeBytes int64     `json:"size_bytes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type exportRequest struct {
	Format   string          `json:"format"`
	Filters  json.RawMessage `json:"filters"`
	Delivery json.RawMessage `json:"delivery"`
}

type exportResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type apiErrorPayload struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	RequestID         string `json:"request_id"`
	RetryAfterSeconds int    `json:"retry_after_seconds,omitempty"`
}

type healthResponse struct {
	Success         bool          `json:"success"`
	RequestID       string        `json:"request_id"`
	Error           string        `json:"error,omitempty"`
	LibraryStatus   string        `json:"library_status"`
	Version         string        `json:"version,omitempty"`
	AvailableAlgos  []string      `json:"available_algos,omitempty"`
	OpenZLAvailable bool          `json:"openzl_available"`
	License         licenseStatus `json:"license"`
	Database        string        `json:"database"`
	Queue           queueStatus   `json:"queue"`
}

type queueStatus struct {
	Mode     string `json:"mode"`
	Pending  int    `json:"pending"`
	Capacity int    `json:"capacity"`
}

func main() {
	// Initialize telemetry (logging + tracing) first
	telemetryCfg := LoadTelemetryConfig()
	if err := InitTelemetry(context.Background(), telemetryCfg); err != nil {
		log.Fatalf("telemetry init failed: %v", err)
	}
	defer ShutdownTelemetry(context.Background())

	cfg := loadConfig()
	var err error
	licenseInfo, err = loadLicense(time.Now())
	if err != nil {
		Logger().Fatal("license validation failed", zap.Error(err))
	}
	if os.Getenv("DRIFTLOCK_DEV_MODE") == "true" {
		Logger().Warn("Running in development mode - license validation bypassed")
	}

	if handleCLI(os.Args[1:], cfg) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := connectDB(ctx)
	if err != nil {
		Logger().Fatal("db connection failed", zap.Error(err))
	}
	store := newStore(pool)
	if err := store.loadCache(ctx); err != nil {
		Logger().Fatal("load configs failed", zap.Error(err))
	}

	// Initialize Firebase Auth
	if os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY") != "" {
		if err := initFirebaseAuth(); err != nil {
			Logger().Warn("Failed to init Firebase Auth - Dashboard endpoints will fail", zap.Error(err))
		} else {
			Logger().Info("Firebase Auth initialized")
		}
	} else {
		Logger().Warn("FIREBASE_SERVICE_ACCOUNT_KEY not set - Dashboard endpoints disabled")
	}

	defer store.Close()

	registerMetrics()

	queue := newMemoryQueue(cfg.QueueCapacity)
	limiter := newTenantRateLimiter(cfg.DefaultRateLimit())
	emailer := newEmailService()
	tracker := newUsageTracker(store, emailer)

	// Initialize webhook event store and retry worker
	// Use a long-lived context for the worker (not the 10-second startup context)
	webhookStore := NewWebhookEventStore(pool, DefaultRetryConfig())
	workerCtx, workerCancel := context.WithCancel(context.Background())
	webhookRetryWorker := NewWebhookRetryWorker(store, webhookStore, emailer, 30*time.Second, 50)
	webhookRetryWorker.Start(workerCtx)

	// Initialize billing cron worker for grace period downgrades and trial reminders
	billingCronWorker := NewBillingCronWorker(store, emailer)
	go billingCronWorker.Start(workerCtx)

	defer func() {
		workerCancel()
		webhookRetryWorker.Stop()
	}()

	// Initialize AI components
	configRepo := ai.NewConfigRepository(store.pool)
	costLimiter := ai.NewCostLimiter(configRepo)
	router := ai.NewRouter(costLimiter)

	var aiClient ai.AIClient
	aiClient, err = ai.NewAIClientFromEnv()
	if err != nil {
		Logger().Warn("Failed to initialize AI client - AI analysis disabled", zap.Error(err))
	} else {
		Logger().Info("AI client initialized", zap.String("provider", aiClient.Provider()))
	}

	handler := buildHTTPHandler(cfg, store, queue, limiter, emailer, tracker, router, aiClient, webhookStore)

	// Strip /api prefix for Firebase hosting rewrite compatibility
	// Firebase routes /api/** to Cloud Run with the full path, but our routes don't have the /api prefix
	apiPrefixHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If path starts with /api, strip it for routing
		if strings.HasPrefix(r.URL.Path, "/api/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		} else if r.URL.Path == "/api" {
			r.URL.Path = "/"
		}
		handler.ServeHTTP(w, r)
	})

	// Wrap handler with OpenTelemetry HTTP instrumentation for automatic tracing
	otelHandler := otelhttp.NewHandler(apiPrefixHandler, "driftlock-http",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	// Wrap with Sentry panic recovery (outermost middleware)
	finalHandler := SentryPanicRecoveryMiddleware(otelHandler)

	addr := env("PORT", "8080")
	Logger().Info("driftlock-http starting", zap.String("port", addr))

	srv := &http.Server{
		Addr:         ":" + addr,
		Handler:      finalHandler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		Logger().Fatal("server failed", zap.Error(err))
	case sig := <-shutdown:
		Logger().Info("signal received, shutting down", zap.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			Logger().Error("graceful shutdown failed", zap.Error(err))
			_ = srv.Close()
		}
		Logger().Info("shutdown complete")
	}
}

func buildHTTPHandler(cfg config, store *store, queue jobQueue, limiter *tenantRateLimiter, emailer *emailService, tracker *usageTracker, router *ai.Router, aiClient ai.AIClient, webhookStore *WebhookEventStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler(store, queue))
	mux.HandleFunc("/readyz", readinessHandler(store))
	mux.Handle("/metrics", promhttp.Handler())

	// Onboarding endpoints
	mux.HandleFunc("/v1/onboard/signup", onboardSignupHandler(cfg, store, emailer))
	mux.HandleFunc("/v1/onboard/verify", verifyHandler(store, emailer))
	mux.HandleFunc("/v1/onboard/resend-verification", resendVerificationHandler(store, emailer))

	// Demo endpoint (no auth, rate limited by IP)
	demoLimiter := newDemoRateLimiter()
	waitlistLimiter := newWaitlistRateLimiter()
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			demoLimiter.cleanup()
			waitlistLimiter.cleanup()
			cleanupSignupLimiters() // Prevent memory leak from signup rate limiters
		}
	}()
	mux.HandleFunc("/v1/demo/detect", demoDetectHandler(cfg, demoLimiter))

	// Pre-launch waitlist (no auth, rate limited)
	mux.HandleFunc("/v1/waitlist", waitlistHandler(store, waitlistLimiter))

	mux.Handle("/v1/detect", withAuth(store, limiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		detectHandler(w, r, cfg, store, tracker, router, aiClient)
	})))
	mux.Handle("/v1/anomalies", withAuth(store, limiter, anomaliesHandler(store)))
	mux.Handle("/v1/anomalies/", withAuth(store, limiter, anomalyRouter(cfg, store, queue)))

	// Billing endpoints
	mux.Handle("/v1/billing/checkout", withAuth(store, limiter, billingCheckoutHandler(store)))
	mux.Handle("/v1/billing/portal", withFirebaseAuth(store, billingPortalHandler(store)))
	mux.HandleFunc("/v1/billing/webhook", billingWebhookHandler(store, emailer, webhookStore))

	// Dashboard/User endpoints (Firebase Auth)
	mux.Handle("/v1/me/keys", withFirebaseAuth(store, http.HandlerFunc(handleListKeys(store))))
	mux.Handle("/v1/me/keys/create", withFirebaseAuth(store, http.HandlerFunc(handleCreateKey(store))))
	mux.Handle("/v1/me/keys/regenerate", withFirebaseAuth(store, http.HandlerFunc(handleRegenerateKey(store))))
	mux.Handle("/v1/me/keys/revoke", withFirebaseAuth(store, http.HandlerFunc(handleRevokeKey(store))))
	mux.Handle("/v1/me/usage", withFirebaseAuth(store, http.HandlerFunc(handleGetUsage(store))))
	mux.Handle("/v1/me/usage/details", withFirebaseAuth(store, http.HandlerFunc(handleGetUsageDetails(store))))
	mux.Handle("/v1/me/billing", withFirebaseAuth(store, http.HandlerFunc(handleGetBillingStatus(store))))
	mux.Handle("/v1/me/usage/ai", withFirebaseAuth(store, aiUsageHandler(store)))
	mux.Handle("/v1/me/ai/config", withFirebaseAuth(store, aiConfigHandler(store)))

	return withCommon(withRequestContext(mux))
}

func healthHandler(store *store, queue jobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handlePreflight(w, r)
			return
		}
		resp := healthResponse{
			Success:         true,
			RequestID:       requestIDFrom(r.Context()),
			LibraryStatus:   "healthy",
			Version:         "1.0.0",
			AvailableAlgos:  []string{"zstd", "lz4", "gzip"},
			OpenZLAvailable: driftlockcbad.HasOpenZL(),
			License:         currentLicenseStatus(time.Now()),
			Database:        "connected",
		}
		if queue != nil {
			stats := queue.Stats()
			resp.Queue = queueStatus{
				Mode:     stats.Mode,
				Pending:  stats.Pending,
				Capacity: stats.Capacity,
			}
		}
		if !resp.License.ExpiresAt.IsZero() && resp.License.Status != "valid" {
			resp.Success = false
			resp.Error = resp.License.Message
			writeJSON(w, r, http.StatusServiceUnavailable, resp)
			return
		}
		if resp.OpenZLAvailable {
			resp.AvailableAlgos = append(resp.AvailableAlgos, "openzl")
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := store.pool.Ping(ctx); err != nil {
			resp.Success = false
			resp.Database = "error"
			resp.Error = err.Error()
			writeJSON(w, r, http.StatusServiceUnavailable, resp)
			return
		}
		if err := driftlockcbad.ValidateLibrary(); err != nil {
			resp.Success = false
			resp.LibraryStatus = "unhealthy"
			resp.Error = err.Error()
			writeJSON(w, r, http.StatusServiceUnavailable, resp)
			return
		}
		writeJSON(w, r, http.StatusOK, resp)
	}
}

// readinessResponse represents the response from the /readyz endpoint.
type readinessResponse struct {
	Ready     bool              `json:"ready"`
	RequestID string            `json:"request_id"`
	Checks    map[string]string `json:"checks"`
	Error     string            `json:"error,omitempty"`
}

// readinessHandler checks if the service is ready to accept traffic.
// Unlike /healthz (liveness), this checks that all dependencies are ready.
func readinessHandler(store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handlePreflight(w, r)
			return
		}

		resp := readinessResponse{
			Ready:     true,
			RequestID: requestIDFrom(r.Context()),
			Checks:    make(map[string]string),
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Check 1: Database connectivity
		if err := store.pool.Ping(ctx); err != nil {
			resp.Ready = false
			resp.Checks["database"] = "failed: " + err.Error()
			Logger().Warn("readiness check failed: database", zap.Error(err))
		} else {
			resp.Checks["database"] = "ok"
		}

		// Check 2: CBAD library validation
		if err := driftlockcbad.ValidateLibrary(); err != nil {
			resp.Ready = false
			resp.Checks["cbad_library"] = "failed: " + err.Error()
			Logger().Warn("readiness check failed: cbad_library", zap.Error(err))
		} else {
			resp.Checks["cbad_library"] = "ok"
		}

		// Check 3: Config cache loaded (check cache has some entries or is accessible)
		store.cache.mu.RLock()
		tenantCount := len(store.cache.tenants)
		streamCount := len(store.cache.streams)
		store.cache.mu.RUnlock()
		resp.Checks["cache"] = fmt.Sprintf("ok (tenants: %d, streams: %d)", tenantCount, streamCount)

		if !resp.Ready {
			resp.Error = "one or more readiness checks failed"
			writeJSON(w, r, http.StatusServiceUnavailable, resp)
			return
		}

		writeJSON(w, r, http.StatusOK, resp)
	}
}

func anomaliesHandler(store *store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		tc, ok := tenantFromContext(r.Context())
		if !ok {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		limit := 50
		if v := r.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
				limit = n
			}
		}
		streamFilter := r.URL.Query().Get("stream_id")
		args := []any{tc.Tenant.ID}
		baseQuery := `FROM anomalies WHERE tenant_id=$1`
		countArgs := []any{tc.Tenant.ID}
		if streamFilter != "" {
			var streamID uuid.UUID
			if id, err := uuid.Parse(streamFilter); err == nil {
				streamID = id
			} else {
				stream, _, ok := store.streamBySlugOrID(tc.Tenant.ID, streamFilter)
				if !ok {
					writeJSON(w, r, http.StatusOK, anomalyListResponse{Anomalies: []anomalyListItem{}, Total: 0})
					return
				}
				streamID = stream.ID
			}
			baseQuery += " AND stream_id = $2"
			args = append(args, streamID)
			countArgs = append(countArgs, streamID)
		}

		// Get total count
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		countQuery := `SELECT COUNT(*) ` + baseQuery
		var total int
		if err := store.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		// Get paginated results
		query := `SELECT id, stream_id, ncd, compression_ratio, entropy_change, p_value, confidence, explanation, status, detected_at ` + baseQuery
		query += " ORDER BY detected_at DESC, id DESC LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, limit)

		rows, err := store.pool.Query(ctx, query, args...)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		var items []anomalyListItem
		for rows.Next() {
			var it anomalyListItem
			var streamID uuid.UUID
			if err := rows.Scan(&it.ID, &streamID, &it.NCD, &it.CompressionRatio, &it.EntropyChange, &it.PValue, &it.Confidence, &it.Explanation, &it.Status, &it.DetectedAt); err != nil {
				writeError(w, r, http.StatusInternalServerError, err)
				return
			}
			it.StreamID = streamID.String()
			items = append(items, it)
		}
		if err := rows.Err(); err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, r, http.StatusOK, anomalyListResponse{Anomalies: items, Total: total})
	})
}

func anomalyRouter(cfg config, store *store, queue jobQueue) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/anomalies/")
		path = strings.Trim(path, "/")
		if path == "" {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("resource not found"))
			return
		}
		parts := strings.Split(path, "/")
		if len(parts) == 1 && strings.EqualFold(parts[0], "export") {
			handleBulkExport(w, r, cfg, store, queue)
			return
		}
		id, err := uuid.Parse(parts[0])
		if err != nil {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("invalid anomaly id"))
			return
		}
		if len(parts) == 1 {
			handleAnomalyDetail(w, r, store, id)
			return
		}
		if len(parts) == 2 && strings.EqualFold(parts[1], "export") {
			handleSingleExport(w, r, cfg, store, queue, id)
			return
		}
		writeError(w, r, http.StatusNotFound, fmt.Errorf("resource not found"))
	})
}

func handleAnomalyDetail(w http.ResponseWriter, r *http.Request, store *store, anomalyID uuid.UUID) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	tc, ok := tenantFromContext(r.Context())
	if !ok {
		writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	rec, evidence, err := store.fetchAnomaly(ctx, tc.Tenant.ID, anomalyID)
	if err != nil {
		if errors.Is(err, errNotFound) {
			writeError(w, r, http.StatusNotFound, fmt.Errorf("anomaly not found"))
			return
		}
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}
	resp := anomalyDetailResponse{
		ID:          rec.ID,
		StreamID:    rec.StreamID,
		BatchID:     rec.BatchID,
		Status:      rec.Status,
		DetectedAt:  rec.DetectedAt,
		Explanation: rec.Explanation,
		Metrics: anomalyDetailMetrics{
			NCD:              rec.NCD,
			CompressionRatio: rec.CompressionRatio,
			EntropyChange:    rec.EntropyChange,
			PValue:           rec.PValue,
			Confidence:       rec.Confidence,
		},
		Details:          rec.Details,
		BaselineSnapshot: rec.BaselineSnapshot,
		WindowSnapshot:   rec.WindowSnapshot,
		Evidence:         make([]anomalyEvidence, 0, len(evidence)),
	}
	for _, ev := range evidence {
		resp.Evidence = append(resp.Evidence, anomalyEvidence{
			Format:    ev.Format,
			URI:       ev.URI,
			Checksum:  ev.Checksum,
			SizeBytes: ev.SizeBytes,
			CreatedAt: ev.CreatedAt,
		})
	}
	writeJSON(w, r, http.StatusOK, resp)
}

func handleBulkExport(w http.ResponseWriter, r *http.Request, cfg config, store *store, queue jobQueue) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	tc, ok := tenantFromContext(r.Context())
	if !ok {
		writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	req, err := decodeExportRequest(r, cfg.MaxBodyBytes)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	jobID, err := store.createExportJob(ctx, tc.Tenant.ID, req.Format, req.Filters, req.Delivery)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}
	if queue != nil {
		payload := append([]byte(nil), req.Filters...)
		if err := queue.Enqueue(r.Context(), queueJob{
			ID:       jobID,
			TenantID: tc.Tenant.ID,
			Type:     "anomaly_export",
			Payload:  payload,
		}); err != nil {
			log.Printf("queue enqueue failed: %v", err)
		}
	}
	writeJSON(w, r, http.StatusAccepted, exportResponse{
		JobID:   jobID.String(),
		Status:  "not_implemented",
		Message: "export worker queue stubbed; payload recorded for future worker",
	})
}

func handleSingleExport(w http.ResponseWriter, r *http.Request, cfg config, store *store, queue jobQueue, anomalyID uuid.UUID) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	tc, ok := tenantFromContext(r.Context())
	if !ok {
		writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	req, err := decodeExportRequest(r, cfg.MaxBodyBytes)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}
	if len(req.Filters) == 0 {
		filterPayload, _ := json.Marshal(map[string]string{"anomaly_id": anomalyID.String()})
		req.Filters = filterPayload
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	jobID, err := store.createExportJob(ctx, tc.Tenant.ID, req.Format, req.Filters, req.Delivery)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}
	if queue != nil {
		payload := append([]byte(nil), req.Filters...)
		if err := queue.Enqueue(r.Context(), queueJob{
			ID:       jobID,
			TenantID: tc.Tenant.ID,
			Type:     "anomaly_export_single",
			Payload:  payload,
		}); err != nil {
			log.Printf("queue enqueue failed: %v", err)
		}
	}
	writeJSON(w, r, http.StatusAccepted, exportResponse{
		JobID:   jobID.String(),
		Status:  "not_implemented",
		Message: "per-anomaly export job recorded; worker implementation pending",
	})
}

func decodeExportRequest(r *http.Request, maxBody int64) (exportRequest, error) {
	var req exportRequest
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		return req, fmt.Errorf("unable to read body: %w", err)
	}
	if len(bytes.TrimSpace(body)) == 0 {
		req.Format = "json"
		req.Filters = []byte("{}")
		req.Delivery = []byte(`{"type":"inline"}`)
		return req, nil
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return req, fmt.Errorf("invalid json: %w", err)
	}
	if req.Format == "" {
		req.Format = "json"
	}
	if len(req.Filters) == 0 {
		req.Filters = []byte("{}")
	} else if !json.Valid(req.Filters) {
		return req, fmt.Errorf("filters must be valid json")
	}
	if len(req.Delivery) == 0 {
		req.Delivery = []byte(`{"type":"inline"}`)
	} else if !json.Valid(req.Delivery) {
		return req, fmt.Errorf("delivery must be valid json")
	}
	req.Format = strings.ToLower(req.Format)
	switch req.Format {
	case "json", "markdown", "html", "pdf":
	default:
		return req, fmt.Errorf("unsupported format %q", req.Format)
	}
	return req, nil
}

func detectHandler(w http.ResponseWriter, r *http.Request, cfg config, store *store, tracker *usageTracker, router *ai.Router, aiClient ai.AIClient) {
	if r.Method == http.MethodOptions {
		handlePreflight(w, r)
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}
	tc, ok := tenantFromContext(r.Context())
	if !ok {
		writeError(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, cfg.MaxBodyBytes))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("unable to read body: %w", err))
		return
	}
	defer r.Body.Close()

	var payload detectRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("invalid json: %w", err))
		return
	}
	if len(payload.Events) == 0 {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("events required"))
		return
	}
	if cfg.MaxEvents > 0 && len(payload.Events) > cfg.MaxEvents {
		writeError(w, r, http.StatusBadRequest, fmt.Errorf("too many events: max %d per request", cfg.MaxEvents))
		return
	}
	for idx, ev := range payload.Events {
		if len(bytes.TrimSpace(ev)) == 0 || bytes.Equal(bytes.TrimSpace(ev), []byte("null")) {
			writeError(w, r, http.StatusBadRequest, fmt.Errorf("event %d is empty", idx))
			return
		}
	}

	stream, settings, err := resolveStream(store, tc, payload.StreamID)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err)
		return
	}
	plan := buildDetectionSettings(cfg, stream, settings, payload.ConfigOverride)

	usedAlgo := plan.CompressionAlgorithm
	fallbackFrom := ""
	if usedAlgo == "openzl" && !driftlockcbad.HasOpenZL() {
		fallbackFrom = usedAlgo
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

	anomalies, records, err := runDetectionWithRecovery(r.Context(), detector, payload.Events)
	if err != nil {
		// Distinguish between user errors and internal FFI errors
		if errors.Is(err, ErrCBADPanic) || errors.Is(err, ErrCBADTimeout) {
			log.Printf("CBAD FFI error: %v", err)
			writeError(w, r, http.StatusInternalServerError, fmt.Errorf("detection service temporarily unavailable"))
			return
		}
		writeError(w, r, http.StatusBadRequest, err)
		return
	}

	batchID, anomalyIDs, err := persistDetection(r.Context(), store, tc, stream, body, records)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, err)
		return
	}
	for i := range anomalies {
		if i < len(anomalyIDs) {
			anomalies[i].ID = anomalyIDs[i]
		}
	}

	// Track usage asynchronously
	go tracker.track(context.Background(), tc.Tenant.ID, stream.ID, tc.Tenant.Plan, len(payload.Events), len(anomalies))

	// AI Analysis (Async)
	if router != nil && aiClient != nil {
		go func() {
			ctx := context.Background()
			for _, anomaly := range anomalies {
				// Create event object for router
				event := ai.Event{
					ID:           anomaly.ID,
					TenantID:     tc.Tenant.ID.String(),
					StreamID:     stream.ID.String(),
					AnomalyScore: anomaly.Metrics.ConfidenceLevel, // Use confidence as score
					CreatedAt:    time.Now(),
				}

				decision, err := router.RouteEvent(ctx, event)
				if err != nil {
					log.Printf("Failed to route event %s: %v", anomaly.ID, err)
					continue
				}

				if decision.ShouldAnalyze {
					// Construct prompt
					prompt := fmt.Sprintf("Analyze this anomaly: %+v", anomaly) // Simplified prompt

					explanation, inputTokens, outputTokens, err := aiClient.AnalyzeAnomaly(ctx, decision.Model, prompt)
					if err != nil {
						log.Printf("Failed to analyze anomaly %s: %v", anomaly.ID, err)
						continue
					}

					// Update anomaly with explanation
					err = store.updateAnomalyExplanation(ctx, anomaly.ID, explanation)
					if err != nil {
						log.Printf("Failed to update anomaly explanation %s: %v", anomaly.ID, err)
					}

					// Track AI usage
					cost := router.CalculateCost(decision.Model, inputTokens, outputTokens)
					err = trackAIUsage(store, tc.Tenant.ID.String(), stream.ID.String(), decision.Model, inputTokens, outputTokens, cost)
					if err != nil {
						log.Printf("Failed to track AI usage for %s: %v", anomaly.ID, err)
					}
				}
			}
		}()
	}

	resp := detectResponse{
		Success:         true,
		BatchID:         batchID,
		StreamID:        stream.ID.String(),
		TotalEvents:     len(payload.Events),
		AnomalyCount:    len(anomalies), // Count of anomalies detected in this batch only
		ProcessingTime:  time.Since(start).String(),
		CompressionAlg:  usedAlgo,
		FallbackFromAlg: fallbackFrom,
		Anomalies:       anomalies,
		RequestID:       requestIDFrom(r.Context()),
	}
	writeJSON(w, r, http.StatusOK, resp)
	requestDuration.Observe(time.Since(start).Seconds())
}

func resolveStream(store *store, tc tenantContext, hint string) (streamRecord, streamSettings, error) {
	if hint != "" {
		if s, cfg, ok := store.streamBySlugOrID(tc.Tenant.ID, hint); ok {
			return s, cfg, nil
		}
		return streamRecord{}, streamSettings{}, fmt.Errorf("unknown stream %s", hint)
	}
	if tc.Key.StreamID != uuid.Nil {
		if s, cfg, ok := store.streamByID(tc.Key.StreamID); ok {
			return s, cfg, nil
		}
	}
	if s, cfg, ok := store.defaultStream(tc.Tenant.ID); ok {
		return s, cfg, nil
	}
	return streamRecord{}, streamSettings{}, fmt.Errorf("no stream configured")
}

type detectionPlan struct {
	BaselineSize         int
	WindowSize           int
	HopSize              int
	NCDThreshold         float64
	PValueThreshold      float64
	PermutationCount     int
	CompressionAlgorithm string
	Seed                 uint64
}

func buildDetectionSettings(cfg config, stream streamRecord, settings streamSettings, override *configOverride) detectionPlan {
	plan := detectionPlan{
		BaselineSize:         settings.BaselineSize,
		WindowSize:           settings.WindowSize,
		HopSize:              settings.HopSize,
		NCDThreshold:         settings.NCDThreshold,
		PValueThreshold:      settings.PValueThreshold,
		PermutationCount:     settings.PermutationCount,
		CompressionAlgorithm: strings.ToLower(settings.Compressor),
		Seed:                 uint64(stream.Seed),
	}
	if plan.Seed == 0 {
		plan.Seed = cfg.Seed
	}
	if plan.CompressionAlgorithm == "" {
		plan.CompressionAlgorithm = strings.ToLower(cfg.DefaultAlgo)
	}

	overrodeCompressor := false
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
			overrodeCompressor = true
		}
	}
	if cfg.PreferOpenZL && driftlockcbad.HasOpenZL() && plan.CompressionAlgorithm != "openzl" && !overrodeCompressor {
		plan.CompressionAlgorithm = "openzl"
	}
	return plan
}

// cbadTimeout is the maximum time allowed for CBAD FFI operations
const cbadTimeout = 15 * time.Second

// ErrCBADPanic is returned when the CBAD detector panics
var ErrCBADPanic = errors.New("CBAD detector panic recovered")

// ErrCBADTimeout is returned when a CBAD operation times out
var ErrCBADTimeout = errors.New("CBAD detector operation timed out")

// runDetectionWithRecovery wraps runDetection with panic recovery and timeout
func runDetectionWithRecovery(ctx context.Context, detector *driftlockcbad.Detector, events []json.RawMessage) (outputs []anomalyOutput, records []persistedAnomaly, err error) {
	// Create a context with timeout for the entire detection operation
	ctx, cancel := context.WithTimeout(ctx, cbadTimeout)
	defer cancel()

	// Channel to receive results
	type result struct {
		outputs []anomalyOutput
		records []persistedAnomaly
		err     error
	}
	resultCh := make(chan result, 1)

	go func() {
		// Recover from any panics in the CBAD FFI code
		defer func() {
			if r := recover(); r != nil {
				log.Printf("CBAD FFI panic recovered: %v", r)
				resultCh <- result{nil, nil, fmt.Errorf("%w: %v", ErrCBADPanic, r)}
			}
		}()

		outputs, records, err := runDetection(detector, events)
		resultCh <- result{outputs, records, err}
	}()

	select {
	case <-ctx.Done():
		log.Printf("CBAD detection timed out after %v", cbadTimeout)
		return nil, nil, ErrCBADTimeout
	case res := <-resultCh:
		return res.outputs, res.records, res.err
	}
}

func runDetection(detector *driftlockcbad.Detector, events []json.RawMessage) ([]anomalyOutput, []persistedAnomaly, error) {
	outputs := make([]anomalyOutput, 0)
	records := make([]persistedAnomaly, 0)
	for idx, ev := range events {
		added, err := detector.AddData(ev)
		if err != nil {
			return nil, nil, err
		}
		if !added {
			continue
		}
		ready, err := detector.IsReady()
		if err != nil {
			return nil, nil, err
		}
		if !ready {
			continue
		}
		detected, metrics, err := detector.DetectAnomaly()
		if err != nil {
			return nil, nil, err
		}
		if !detected {
			continue
		}
		explanation := metrics.GetDetailedExplanation()
		snapshot := append([]byte(nil), ev...)
		outputs = append(outputs, anomalyOutput{
			Index:    idx,
			Metrics:  *metrics,
			Event:    json.RawMessage(snapshot),
			Why:      explanation,
			Detected: true,
		})
		detail := map[string]any{
			"event":   json.RawMessage(snapshot),
			"metrics": metricsToMap(metrics),
		}
		detailJSON, _ := json.Marshal(detail)
		records = append(records, persistedAnomaly{
			NCD:              metrics.NCD,
			CompressionRatio: metrics.WindowCompressionRatio,
			EntropyChange:    metrics.EntropyChange,
			PValue:           metrics.PValue,
			Confidence:       metrics.ConfidenceLevel,
			Explanation:      explanation,
			Details:          detailJSON,
			EvidenceFormat:   "markdown",
		})
	}
	return outputs, records, nil
}

func metricsToMap(m *driftlockcbad.EnhancedMetrics) map[string]any {
	return map[string]any{
		"ncd":                        m.NCD,
		"p_value":                    m.PValue,
		"baseline_compression_ratio": m.BaselineCompressionRatio,
		"window_compression_ratio":   m.WindowCompressionRatio,
		"baseline_entropy":           m.BaselineEntropy,
		"window_entropy":             m.WindowEntropy,
		"confidence":                 m.ConfidenceLevel,
		"statistically_significant":  m.IsStatisticallySignificant,
		"compression_ratio_change":   m.CompressionRatioChange,
		"entropy_change":             m.EntropyChange,
	}
}

func persistDetection(ctx context.Context, store *store, tc tenantContext, stream streamRecord, body []byte, records []persistedAnomaly) (string, []string, error) {
	batchHashVal := batchHash(tc.Tenant.ID, time.Now().UTC(), body)
	batchID, err := store.insertBatch(ctx, tc.Tenant.ID, stream.ID, batchHashVal, "sync")
	if err != nil {
		return "", nil, err
	}
	idsUUID, err := store.insertAnomalies(ctx, batchID, tc.Tenant.ID, stream.ID, records)
	if err != nil {
		return "", nil, err
	}
	ids := make([]string, len(idsUUID))
	for i, id := range idsUUID {
		ids[i] = id.String()
	}
	return batchID.String(), ids, nil
}

// Middleware and helpers remain largely unchanged from the earlier version.

type ctxKey string

const requestIDKey ctxKey = "reqid"

func withRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		logRequest(r, reqID, "request_start", "")
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
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
	logEntry := map[string]any{
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
	logEntry := map[string]any{
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
	code := "internal"
	switch status {
	case http.StatusBadRequest:
		code = "invalid_argument"
	case http.StatusUnauthorized:
		code = "unauthorized"
	case http.StatusForbidden:
		code = "forbidden"
	case http.StatusNotFound:
		code = "not_found"
	case http.StatusTooManyRequests:
		code = "rate_limit_exceeded"
	case http.StatusServiceUnavailable:
		code = "service_unavailable"
	case http.StatusMethodNotAllowed:
		code = "method_not_allowed"
	}
	logError(r, reqID, code, fmt.Sprintf("http_status=%d", status), err)
	payload := apiErrorPayload{Error: apiError{Code: code, Message: err.Error(), RequestID: reqID}}
	writeJSON(w, r, status, payload)
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

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "y", "on":
			return true
		case "0", "false", "no", "n", "off":
			return false
		}
	}
	return def
}

func parseAllowedOrigins(v string) []string {
	v = strings.ReplaceAll(v, "%2C", ",")
	v = strings.ReplaceAll(v, "%2c", ",")
	v = strings.ReplaceAll(v, "|", ",")
	v = strings.ReplaceAll(v, ";", ",")
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

func withCommon(next http.Handler) http.Handler {
	allowed := parseAllowedOrigins(env("CORS_ALLOW_ORIGINS", "*"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin, allowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-Id, X-Api-Key")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Content-Type", "application/json")
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
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
	openZLAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "driftlock_openzl_available",
		Help: "Whether OpenZL symbols are present in the CBAD core",
	})
	registerMetricsOnce sync.Once
)

func registerMetrics() {
	registerMetricsOnce.Do(func() {
		prometheus.MustRegister(requestCounter, requestDuration, openZLAvailable)
		if driftlockcbad.HasOpenZL() {
			openZLAvailable.Set(1)
		} else {
			openZLAvailable.Set(0)
		}
	})
}
