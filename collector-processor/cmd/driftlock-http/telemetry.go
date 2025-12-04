package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	logger *zap.Logger

	// Global tracer instance
	tracer trace.Tracer

	// Tracer provider for shutdown
	tracerProvider *sdktrace.TracerProvider
)

// TelemetryConfig holds configuration for observability components.
type TelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	LogLevel       string
	LogFormat      string // "json" or "console"
	OTELEnabled    bool
	OTELEndpoint   string
	SampleRate     float64
	SentryDSN      string
	SentryEnabled  bool
}

// LoadTelemetryConfig loads telemetry configuration from environment variables.
func LoadTelemetryConfig() TelemetryConfig {
	sampleRate := 0.1 // 10% sampling by default in production
	if os.Getenv("DRIFTLOCK_DEV_MODE") == "true" {
		sampleRate = 1.0 // 100% sampling in dev
	}
	if v := os.Getenv("OTEL_SAMPLE_RATE"); v != "" {
		if f := envFloat("OTEL_SAMPLE_RATE", sampleRate); f >= 0 && f <= 1 {
			sampleRate = f
		}
	}

	sentryDSN := env("SENTRY_DSN", "")
	return TelemetryConfig{
		ServiceName:    env("OTEL_SERVICE_NAME", "driftlock-http"),
		ServiceVersion: env("DRIFTLOCK_VERSION", "1.0.0"),
		Environment:    env("OTEL_ENV", "production"),
		LogLevel:       env("LOG_LEVEL", "info"),
		LogFormat:      env("LOG_FORMAT", "json"),
		OTELEnabled:    envBool("OTEL_TRACES_ENABLED", false),
		OTELEndpoint:   env("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		SampleRate:     sampleRate,
		SentryDSN:      sentryDSN,
		SentryEnabled:  sentryDSN != "",
	}
}

// InitTelemetry initializes logging and tracing. Call ShutdownTelemetry on exit.
func InitTelemetry(ctx context.Context, cfg TelemetryConfig) error {
	// Initialize structured logging
	if err := initLogger(cfg); err != nil {
		return err
	}

	// Initialize Sentry error tracking if enabled
	if cfg.SentryEnabled {
		if err := initSentry(cfg); err != nil {
			logger.Warn("Failed to initialize Sentry, continuing without it", zap.Error(err))
		} else {
			logger.Info("Sentry error tracking initialized",
				zap.String("environment", cfg.Environment))
		}
	} else {
		logger.Info("Sentry error tracking disabled")
	}

	// Initialize OpenTelemetry tracing if enabled
	if cfg.OTELEnabled && cfg.OTELEndpoint != "" {
		if err := initTracing(ctx, cfg); err != nil {
			logger.Warn("Failed to initialize tracing, continuing without it", zap.Error(err))
		} else {
			logger.Info("OpenTelemetry tracing initialized",
				zap.String("endpoint", cfg.OTELEndpoint),
				zap.Float64("sample_rate", cfg.SampleRate))
		}
	} else {
		// Create a no-op tracer
		tracer = otel.Tracer(cfg.ServiceName)
		logger.Info("OpenTelemetry tracing disabled")
	}

	return nil
}

// initSentry initializes Sentry error tracking.
func initSentry(cfg TelemetryConfig) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.Environment,
		Release:          cfg.ServiceName + "@" + cfg.ServiceVersion,
		AttachStacktrace: true,
		// Sample rate for error events (1.0 = 100% of errors)
		SampleRate: 1.0,
		// BeforeSend can be used to filter/modify events
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Scrub sensitive data if needed
			return event
		},
	})
}

// initLogger initializes the zap logger based on configuration.
func initLogger(cfg TelemetryConfig) error {
	var level zapcore.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var zapCfg zap.Config
	if strings.ToLower(cfg.LogFormat) == "console" || os.Getenv("DRIFTLOCK_DEV_MODE") == "true" {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	var err error
	logger, err = zapCfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("service", cfg.ServiceName),
			zap.String("version", cfg.ServiceVersion),
			zap.String("env", cfg.Environment),
		),
	)
	if err != nil {
		return err
	}

	// Replace the global logger
	zap.ReplaceGlobals(logger)

	return nil
}

// initTracing initializes OpenTelemetry tracing with OTLP exporter.
func initTracing(ctx context.Context, cfg TelemetryConfig) error {
	// Create OTLP exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.OTELEndpoint),
		otlptracegrpc.WithInsecure(), // Use TLS in production via env config
	)
	if err != nil {
		return err
	}

	// Detect GCP resource attributes (project, region, etc.)
	gcpDetector := gcp.NewDetector()
	res, err := resource.New(ctx,
		resource.WithDetectors(gcpDetector),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		logger.Warn("Failed to create resource, using default", zap.Error(err))
		res = resource.Default()
	}

	// Create sampler based on configuration
	var sampler sdktrace.Sampler
	if cfg.SampleRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if cfg.SampleRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(cfg.SampleRate)
	}

	// Create trace provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global trace provider
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator for distributed tracing
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer for this service
	tracer = tracerProvider.Tracer(cfg.ServiceName)

	return nil
}

// ShutdownTelemetry gracefully shuts down telemetry components.
func ShutdownTelemetry(ctx context.Context) {
	// Flush Sentry events before shutdown
	if sentry.CurrentHub().Client() != nil {
		sentry.Flush(2 * time.Second)
	}

	if tracerProvider != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}
	if logger != nil {
		_ = logger.Sync()
	}
}

// Logger returns the global zap logger.
func Logger() *zap.Logger {
	if logger == nil {
		// Fallback to a basic logger if not initialized
		logger, _ = zap.NewProduction()
	}
	return logger
}

// Tracer returns the global tracer.
func Tracer() trace.Tracer {
	if tracer == nil {
		tracer = otel.Tracer("driftlock-http")
	}
	return tracer
}

// StartSpan starts a new span with the given name and returns a context and span.
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name, opts...)
}

// SpanFromContext returns the current span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanEvent adds an event to the current span.
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanError records an error on the current span.
func SetSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// WithRequestID returns a logger with the request ID field.
func WithRequestID(reqID string) *zap.Logger {
	return Logger().With(zap.String("request_id", reqID))
}

// WithTenant returns a logger with tenant context.
func WithTenant(tenantID, tenantName string) *zap.Logger {
	return Logger().With(
		zap.String("tenant_id", tenantID),
		zap.String("tenant_name", tenantName),
	)
}

// WithStream returns a logger with stream context.
func WithStream(streamID, streamSlug string) *zap.Logger {
	return Logger().With(
		zap.String("stream_id", streamID),
		zap.String("stream_slug", streamSlug),
	)
}

// Common span attribute helpers
func TenantAttr(id string) attribute.KeyValue {
	return attribute.String("tenant.id", id)
}

func StreamAttr(id string) attribute.KeyValue {
	return attribute.String("stream.id", id)
}

func RequestIDAttr(id string) attribute.KeyValue {
	return attribute.String("request.id", id)
}

func EventCountAttr(count int) attribute.KeyValue {
	return attribute.Int("event.count", count)
}

func AnomalyCountAttr(count int) attribute.KeyValue {
	return attribute.Int("anomaly.count", count)
}

// --- Sentry Error Capture Helpers ---

// CaptureError captures an error to Sentry with optional context.
func CaptureError(err error, ctx context.Context, tags map[string]string) {
	if sentry.CurrentHub().Client() == nil {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		// Add trace context if available
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			scope.SetTag("trace_id", span.SpanContext().TraceID().String())
			scope.SetTag("span_id", span.SpanContext().SpanID().String())
		}
		// Add custom tags
		for k, v := range tags {
			scope.SetTag(k, v)
		}
	})
	hub.CaptureException(err)
}

// CaptureMessage captures a message to Sentry.
func CaptureMessage(msg string, level sentry.Level, tags map[string]string) {
	if sentry.CurrentHub().Client() == nil {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		for k, v := range tags {
			scope.SetTag(k, v)
		}
	})
	hub.CaptureMessage(msg)
}

// SetSentryUser sets user context for the current scope (call per-request).
func SetSentryUser(tenantID, tenantName, email string) {
	if sentry.CurrentHub().Client() == nil {
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:       tenantID,
			Username: tenantName,
			Email:    email,
		})
	})
}

// SentryPanicRecoveryMiddleware wraps an http.Handler with Sentry panic recovery.
func SentryPanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetRequest(r)
			scope.SetTag("method", r.Method)
			scope.SetTag("path", r.URL.Path)
		})

		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				Logger().Error("Panic recovered",
					zap.Any("panic", err),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)

				// Capture to Sentry
				hub.RecoverWithContext(r.Context(), err)

				// Return 500 to client
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
