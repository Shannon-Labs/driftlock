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
	var args []interface{}

	// Add request ID if present
	if reqID := ctx.Value("request_id"); reqID != nil {
		args = append(args, "request_id", reqID.(string))
	}

	// Add user ID if present
	if userID := ctx.Value("user_id"); userID != nil {
		args = append(args, "user_id", userID.(string))
	}

	if len(args) > 0 {
		return &Logger{Logger: l.With(args...)}
	}

	return l
}

// HTTP logs HTTP request/response information
func (l *Logger) HTTP(method, path string, status int, duration time.Duration, attrs ...slog.Attr) {
	args := []interface{}{
		"method", method,
		"path", path,
		"status", status,
		"duration", duration.String(),
		"duration_ms", duration.Milliseconds(),
	}

	// Convert attrs to args
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}

	if status >= 500 {
		l.Logger.Error("HTTP request", args...)
	} else if status >= 400 {
		l.Logger.Warn("HTTP request", args...)
	} else {
		l.Logger.Info("HTTP request", args...)
	}
}

// Database logs database operations
func (l *Logger) Database(operation string, duration time.Duration, err error, attrs ...slog.Attr) {
	args := []interface{}{
		"operation", operation,
		"duration", duration.String(),
	}

	// Convert attrs to args
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}

	if err != nil {
		args = append(args, "error", err.Error())
		l.Logger.Error("Database operation failed", args...)
	} else {
		l.Logger.Debug("Database operation", args...)
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
	args := []interface{}{"event", event}
	
	// Convert attrs to args
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	
	l.Logger.Warn("Security event", args...)
}

// Startup logs application startup information
func (l *Logger) Startup(version, port string, attrs ...slog.Attr) {
	args := []interface{}{
		"version", version,
		"port", port,
	}
	
	// Convert attrs to args
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	
	l.Logger.Info("Application starting", args...)
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
	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	Global().Logger.Info(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, attrs ...slog.Attr) {
	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	Global().Logger.Debug(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, attrs ...slog.Attr) {
	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	Global().Logger.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, attrs ...slog.Attr) {
	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	Global().Logger.Error(msg, args...)
}

// HTTP logs HTTP request information
func HTTP(method, path string, status int, duration time.Duration, attrs ...slog.Attr) {
	Global().HTTP(method, path, status, duration, attrs...)
}
