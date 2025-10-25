package logging

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// Logger wraps slog.Logger with additional convenience methods
type Logger struct {
	*slog.Logger
}

// Config holds logger configuration
type Config struct {
	Level  string // debug, info, warn, error
	Format string // json or text
}

// New creates a new structured logger
func New(config Config) *Logger {
	// Parse log level
	var level slog.Level
	switch config.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.String("timestamp", a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	// Select handler based on format
	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Default creates a logger with default settings (JSON, Info level)
func Default() *Logger {
	return New(Config{
		Level:  "info",
		Format: "json",
	})
}

// WithContext adds context values to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract common context values if present
	attrs := []slog.Attr{}

	// Add request ID if present
	if reqID := ctx.Value("request_id"); reqID != nil {
		attrs = append(attrs, slog.String("request_id", reqID.(string)))
	}

	// Add user ID if present
	if userID := ctx.Value("user_id"); userID != nil {
		attrs = append(attrs, slog.String("user_id", userID.(string)))
	}

	if len(attrs) > 0 {
		return &Logger{Logger: l.With(attrs...)}
	}

	return l
}

// HTTP logs HTTP request/response information
func (l *Logger) HTTP(method, path string, status int, duration time.Duration, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.String("duration", duration.String()),
		slog.Int64("duration_ms", duration.Milliseconds()),
	}

	allAttrs := append(baseAttrs, attrs...)

	if status >= 500 {
		l.Error("HTTP request", allAttrs...)
	} else if status >= 400 {
		l.Warn("HTTP request", allAttrs...)
	} else {
		l.Info("HTTP request", allAttrs...)
	}
}

// Database logs database operations
func (l *Logger) Database(operation string, duration time.Duration, err error, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("duration", duration.String()),
	}

	if err != nil {
		baseAttrs = append(baseAttrs, slog.String("error", err.Error()))
		l.Error("Database operation failed", append(baseAttrs, attrs...)...)
	} else {
		l.Debug("Database operation", append(baseAttrs, attrs...)...)
	}
}

// Anomaly logs anomaly detection events
func (l *Logger) Anomaly(id, streamType string, ncdScore, pValue float64, significant bool) {
	l.Info("Anomaly detected",
		slog.String("anomaly_id", id),
		slog.String("stream_type", streamType),
		slog.Float64("ncd_score", ncdScore),
		slog.Float64("p_value", pValue),
		slog.Bool("significant", significant),
	)
}

// Security logs security-related events
func (l *Logger) Security(event string, attrs ...slog.Attr) {
	l.Warn("Security event",
		append([]slog.Attr{slog.String("event", event)}, attrs...)...,
	)
}

// Startup logs application startup information
func (l *Logger) Startup(version, port string, attrs ...slog.Attr) {
	l.Info("Application starting",
		append([]slog.Attr{
			slog.String("version", version),
			slog.String("port", port),
		}, attrs...)...,
	)
}

// Shutdown logs application shutdown
func (l *Logger) Shutdown(reason string) {
	l.Info("Application shutting down", slog.String("reason", reason))
}

// Global logger instance
var global *Logger

// SetGlobal sets the global logger instance
func SetGlobal(l *Logger) {
	global = l
}

// Global returns the global logger instance
func Global() *Logger {
	if global == nil {
		global = Default()
	}
	return global
}

// Helper functions that use the global logger

// Info logs an info message
func Info(msg string, attrs ...slog.Attr) {
	Global().Info(msg, attrs...)
}

// Debug logs a debug message
func Debug(msg string, attrs ...slog.Attr) {
	Global().Debug(msg, attrs...)
}

// Warn logs a warning message
func Warn(msg string, attrs ...slog.Attr) {
	Global().Warn(msg, attrs...)
}

// Error logs an error message
func Error(msg string, attrs ...slog.Attr) {
	Global().Error(msg, attrs...)
}

// HTTP logs HTTP request information
func HTTP(method, path string, status int, duration time.Duration, attrs ...slog.Attr) {
	Global().HTTP(method, path, status, duration, attrs...)
}
