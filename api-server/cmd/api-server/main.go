package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/Hmbown/driftlock/api-server/internal/auth"
	"github.com/Hmbown/driftlock/api-server/internal/cbad"
	"github.com/Hmbown/driftlock/api-server/internal/config"
	"github.com/Hmbown/driftlock/api-server/internal/engine"
	"github.com/Hmbown/driftlock/api-server/internal/export"
	"github.com/Hmbown/driftlock/api-server/internal/handlers"
	"github.com/Hmbown/driftlock/api-server/internal/metrics"
	"github.com/Hmbown/driftlock/api-server/internal/storage"
	"github.com/Hmbown/driftlock/api-server/internal/storage/redis"
	"github.com/Hmbown/driftlock/api-server/internal/stream"
	"github.com/Hmbown/driftlock/api-server/internal/streaming"
	"github.com/Hmbown/driftlock/api-server/internal/telemetry"
	"github.com/Hmbown/driftlock/api-server/internal/streaming/kafka"
	"github.com/Hmbown/driftlock/pkg/version"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	log.Printf("Starting Driftlock API Server %s", version.Version())

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup OpenTelemetry
	ctx := context.Background()
	shutdownOtel, err := telemetry.Setup(ctx)
	if err != nil {
		log.Fatalf("Failed telemetry setup: %v", err)
	}
	defer func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownOtel(c); err != nil {
			log.Printf("OTel shutdown error: %v", err)
		}
	}()

	// Connect to PostgreSQL
	connString := cfg.GetDatabaseConnectionString()
	db, err := storage.NewPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to PostgreSQL at %s:%d", cfg.Database.Host, cfg.Database.Port)

	// Initialize tiered storage if configured
	if cfg.Storage.Tiered.Enabled {
		// Parse the archive interval
		archiveInterval, err := time.ParseDuration(cfg.Storage.Tiered.ArchiveInterval)
		if err != nil {
			log.Printf("Invalid archive interval '%s', using default 24h: %v", cfg.Storage.Tiered.ArchiveInterval, err)
			archiveInterval = 24 * time.Hour
		}

		// In a real implementation, we'd have separate connections for warm/cold storage
		// For now, using the single PostgreSQL instance for all tiers in the tiered storage
		tieredStorage := storage.NewTieredStorage(db, db, db, storage.TierConfig{
			HotRetentionDays:  cfg.Storage.Tiered.HotRetentionDays,
			WarmRetentionDays: cfg.Storage.Tiered.WarmRetentionDays,
			ArchiveInterval:   archiveInterval,
		})
		
		// Start the archive worker in the background
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tieredStorage.StartArchiveWorker(ctx)
		
		log.Printf("Tiered storage initialized with hot/warm/cold tiers: hot=%d days, warm=%d days", 
			cfg.Storage.Tiered.HotRetentionDays, cfg.Storage.Tiered.WarmRetentionDays)
	} else {
		log.Printf("Using single storage backend (no tiering)")
	}

	// Initialize SSE streamer
	streamer := stream.NewStreamer(1000) // Max 1000 concurrent connections
	log.Printf("SSE streamer initialized")

	// Initialize Redis state manager for distributed CBAD processing
	var redisStateMgr *redis.StateManager
	if cfg.Cache.Redis.Enabled {
		log.Printf("Initializing Redis state manager at %s", cfg.Cache.Redis.Addr)
		redisStateMgr = redis.NewStateManager(
			cfg.Cache.Redis.Addr,
			cfg.Cache.Redis.Password,
			cfg.Cache.Redis.Prefix,
			cfg.Cache.Redis.DB,
		)
		log.Printf("Redis state manager initialized")
	} else {
		log.Printf("Redis disabled, using local state only")
	}

	// Initialize CBAD detector
	cbadDetector, err := cbad.NewDetector(db, streamer)
	if err != nil {
		log.Fatalf("Failed to initialize CBAD detector: %v", err)
	}
	log.Printf("CBAD detector initialized")

	// Initialize authenticator
	authenticator := auth.NewAuthenticator()

	// Add default API key for development (CHANGE IN PRODUCTION!)
	devKey := os.Getenv("DRIFTLOCK_DEV_API_KEY")
	if devKey != "" {
		authenticator.AddAPIKey(
			auth.HashAPIKey(devKey),
			auth.APIKeyInfo{
				Name:   "development",
				Role:   "admin",
				Scopes: []string{"read:anomalies", "write:anomalies", "admin:config"},
			},
		)
		log.Printf("Authentication initialized (WARNING: using development API key)")
	} else {
		log.Printf("Authentication initialized (no development API key provided)")
	}

	// Initialize exporter
	exporter := export.NewExporter(true) // Enable signature

	// Initialize event publisher (Kafka or in-memory based on config)
	var eventPublisher streaming.EventPublisher
	var kafkaProducer *kafka.Producer // Keep track for proper cleanup

	if cfg.Streaming.Kafka.Enabled {
		log.Printf("Initializing Kafka event publisher with brokers: %v", cfg.Streaming.Kafka.Brokers)
		
		// Create Kafka producer configuration
		producerConfig := kafka.ProducerConfig{
			Brokers:  cfg.Streaming.Kafka.Brokers,
			ClientID: cfg.Streaming.Kafka.ClientID,
		}
		
		// Configure TLS if enabled
		if cfg.Streaming.Kafka.TLSEnabled {
			producerConfig.TLSConfig = &tls.Config{}
		}
		
		// Create Kafka producer
		var err error
		kafkaProducer, err = kafka.NewProducer(producerConfig)
		if err != nil {
			log.Fatalf("Failed to initialize Kafka producer: %v", err)
		}
		
		// Create Kafka event publisher
		eventPublisher = streaming.NewKafkaEventPublisher(kafkaProducer, cfg.Streaming.Kafka.AnomaliesTopic)
		log.Printf("Kafka event publisher initialized for topic: %s", cfg.Streaming.Kafka.AnomaliesTopic)
	} else {
		log.Printf("Kafka disabled, using in-memory publisher for testing")
		// Use in-memory broker for local testing
		memoryBroker := kafka.NewInMemoryBroker()
		eventPublisher = streaming.NewInMemoryPublisher(memoryBroker, "anomaly-events")
	}

	// Create handlers
	anomaliesHandler := handlers.NewAnomaliesHandler(db, streamer, eventPublisher)
	configHandler := handlers.NewConfigHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	exportHandler := handlers.NewExportHandler(db, exporter)

	// Create HTTP mux
	mux := http.NewServeMux()

	// Public endpoints (no auth required)
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check database connectivity
		if err := db.Ping(r.Context()); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	}))

	mux.Handle("/v1/version", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"` + version.Version() + `"}`))
	}), "version"))

	// Initialize the engine with the CBAD detector
	engine := engine.New(cbadDetector)

	// Metrics endpoint (Prometheus)
	mux.Handle("/metrics", promhttp.Handler())

	// Events endpoint (ingestion and anomaly detection)
	mux.Handle("/v1/events", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MiB cap
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		start := time.Now()
		if err := engine.Process(r.Context(), body); err != nil {
			log.Printf("events: processing error: %v", err)
			http.Error(w, "processing error", http.StatusInternalServerError)
			return
		}
		took := time.Since(start)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"processed_bytes": len(body),
			"latency": took.String(),
		})
	}), "events"))

	// SSE streaming endpoint (optional auth)
	mux.Handle("/v1/stream/anomalies", authenticator.OptionalMiddleware(
		otelhttp.NewHandler(streamer, "stream-anomalies"),
	))

	// Create authenticated routes mux
	authMux := http.NewServeMux()

	// Anomaly endpoints
	authMux.Handle("/v1/anomalies", otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				anomaliesHandler.ListAnomalies(w, r)
			case http.MethodPost:
				anomaliesHandler.CreateAnomaly(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		"anomalies",
	))

	authMux.Handle("/v1/anomalies/", otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > len("/v1/anomalies/") {
				if len(path) > 50 && path[len(path)-7:] == "/status" {
					// PATCH /v1/anomalies/:id/status
					if r.Method == http.MethodPatch {
						anomaliesHandler.UpdateAnomalyStatus(w, r)
						return
					}
				} else if len(path) > 50 && path[len(path)-7:] == "/export" {
					// GET /v1/anomalies/:id/export
					if r.Method == http.MethodGet {
						exportHandler.ExportAnomaly(w, r)
						return
					}
				} else if r.Method == http.MethodGet {
					// GET /v1/anomalies/:id
					anomaliesHandler.GetAnomaly(w, r)
					return
				}
			}
			http.Error(w, "Not found", http.StatusNotFound)
		}),
		"anomaly-detail",
	))

	// Configuration endpoints
	authMux.Handle("/v1/config", otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				configHandler.GetConfig(w, r)
			case http.MethodPatch:
				configHandler.UpdateConfig(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		"config",
	))

	// Analytics endpoints
	authMux.Handle("/v1/analytics/summary", otelhttp.NewHandler(
		http.HandlerFunc(analyticsHandler.GetSummary),
		"analytics-summary",
	))

	authMux.Handle("/v1/analytics/compression-timeline", otelhttp.NewHandler(
		http.HandlerFunc(analyticsHandler.GetCompressionTimeline),
		"analytics-timeline",
	))

	authMux.Handle("/v1/analytics/ncd-heatmap", otelhttp.NewHandler(
		http.HandlerFunc(analyticsHandler.GetNCDHeatmap),
		"analytics-heatmap",
	))

	// Performance metrics endpoint
	authMux.Handle("/v1/metrics/performance", otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stats := db.DB().Stats()

			// Update database metrics using available functions
			metrics.DBConnectionsTotal.WithLabelValues("driftlock_db", "total").Set(float64(stats.OpenConnections))
			metrics.DBConnectionsTotal.WithLabelValues("driftlock_db", "in_use").Set(float64(stats.InUse))
			metrics.DBConnectionsTotal.WithLabelValues("driftlock_db", "idle").Set(float64(stats.Idle))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"sse_connections": ` + string(rune(streamer.GetClientCount())) + `,
				"database": {
					"open": ` + string(rune(stats.OpenConnections)) + `,
					"in_use": ` + string(rune(stats.InUse)) + `,
					"idle": ` + string(rune(stats.Idle)) + `
				}
			}`))
		}),
		"performance-metrics",
	))

	// Mount authenticated routes with middleware
	mux.Handle("/v1/", authenticator.Middleware(authMux))

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + string(rune(cfg.Server.Port)),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		log.Printf("API server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Start Prometheus metrics server
	go func() {
		metricsAddr := ":" + string(rune(cfg.Observability.PrometheusPort))
		log.Printf("Prometheus metrics server listening on %s", metricsAddr)
		http.ListenAndServe(metricsAddr, metrics.RegisterMetricsHandler())
	}()

	// Graceful shutdown
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	log.Printf("Shutdown signal received")

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := streamer.Shutdown(shutdownCtx); err != nil {
		log.Printf("SSE shutdown error: %v", err)
	}

	// Close Kafka producer if it was created
	if kafkaProducer != nil {
		log.Printf("Closing Kafka producer...")
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("Kafka producer close error: %v", err)
		}
		log.Printf("Kafka producer closed")
	}

	// Close Redis state manager if it was created
	if redisStateMgr != nil {
		log.Printf("Closing Redis state manager...")
		if err := redisStateMgr.Close(); err != nil {
			log.Printf("Redis state manager close error: %v", err)
		}
		log.Printf("Redis state manager closed")
	}

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Printf("Driftlock API Server stopped")
}
