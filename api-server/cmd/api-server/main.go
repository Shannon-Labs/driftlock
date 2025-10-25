package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/your-org/driftlock/api-server/internal/auth"
	"github.com/your-org/driftlock/api-server/internal/cbad"
	"github.com/your-org/driftlock/api-server/internal/config"
	"github.com/your-org/driftlock/api-server/internal/export"
	"github.com/your-org/driftlock/api-server/internal/handlers"
	"github.com/your-org/driftlock/api-server/internal/metrics"
	"github.com/your-org/driftlock/api-server/internal/storage"
	"github.com/your-org/driftlock/api-server/internal/stream"
	"github.com/your-org/driftlock/api-server/internal/telemetry"
	"github.com/your-org/driftlock/pkg/version"
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

	// Initialize SSE streamer
	streamer := stream.NewStreamer(1000) // Max 1000 concurrent connections
	log.Printf("SSE streamer initialized")

	// Initialize CBAD detector
	detector, err := cbad.NewDetector(db, streamer)
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

	// Create handlers
	anomaliesHandler := handlers.NewAnomaliesHandler(db, streamer)
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

	// Metrics endpoint (Prometheus)
	mux.Handle("/metrics", promhttp.Handler())

	// SSE streaming endpoint (optional auth)
	mux.Handle("/v1/stream/anomalies", authenticator.OptionalMiddleware(
		otelhttp.NewHandler(streamer, "stream-anomalies"),
	))

	// Authenticated endpoints
	authMux := authenticator.Middleware(http.NewServeMux())

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
			metrics.UpdateDatabaseConnectionStats(stats.OpenConnections, stats.InUse, stats.Idle)

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

	// Mount authenticated routes
	mux.Handle("/v1/", authMux)

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
		http.ListenAndServe(metricsAddr, metrics.Handler())
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

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Printf("Driftlock API Server stopped")
}
