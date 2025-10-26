package metrics

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// Namespace for all Driftlock metrics
	Namespace = "driftlock"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:    Namespace,
			Name:         "http_request_duration_seconds",
			Help:         "HTTP request duration",
			Buckets:      prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	HTTPRequestSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  Namespace,
			Name:       "http_request_size_bytes",
			Help:       "Size of HTTP requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "endpoint"},
	)

	HTTPResponseSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  Namespace,
			Name:       "http_response_size_bytes",
			Help:       "Size of HTTP responses",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "endpoint"},
	)

	// Anomaly detection metrics
	AnomaliesDetected = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "anomalies_detected_total",
			Help:      "Total anomalies detected",
		},
		[]string{"stream_type", "severity"},
	)

	AnomalyProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:    Namespace,
			Name:         "anomaly_processing_duration_seconds",
			Help:         "Time taken to process anomalies",
			Buckets:      []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"stream_type"},
	)

	AnomalyNCDThreshold = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "anomaly_ncd_threshold",
			Help:      "Current NCD threshold for anomaly detection",
		},
		[]string{"stream_type"},
	)

	AnomalyPValueThreshold = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "anomaly_p_value_threshold",
			Help:      "Current p-value threshold for anomaly detection",
		},
		[]string{"stream_type"},
	)

	// Database metrics
	DBConnectionsTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "db_connections_total",
			Help:      "Total number of database connections",
		},
		[]string{"db_name", "connection_state"},
	)

	DBConnectionsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:    Namespace,
			Name:         "db_connection_duration_seconds",
			Help:         "Duration of database connections",
			Buckets:      prometheus.DefBuckets,
		},
		[]string{"db_name"},
	)

	// Performance metrics
	GoRoutines = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "goroutines",
			Help:      "Number of goroutines",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)

	GoMemoryAllocated = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "memory_allocated_bytes",
			Help:      "Memory allocated in bytes",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.Alloc)
		},
	)

	GoMemoryHeap = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "memory_heap_bytes",
			Help:      "Heap memory usage in bytes",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.HeapAlloc)
		},
	)

	// CBAD engine metrics
	CBADCompressionRatio = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:    Namespace,
			Name:         "cbad_compression_ratio",
			Help:         "Compression ratios produced by CBAD engine",
			Buckets:      prometheus.LinearBuckets(0, 0.1, 20),
		},
		[]string{"stream_type", "data_type"},
	)

	CBADProcessingRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "cbad_events_per_second",
			Help:      "Rate of events processed by CBAD engine",
		},
		[]string{"stream_type"},
	)

	// API Performance metrics
	APIResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:    Namespace,
			Name:         "api_response_time_seconds",
			Help:         "API response time",
			Buckets:      []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 15.0},
		},
		[]string{"handler", "method"},
	)
)

// Middleware returns a Prometheus middleware for HTTP request metrics
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		endpoint := r.URL.Path

		// Wrap response writer to capture status code and response size
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Continue with the next handler
		next.ServeHTTP(wrapped, r)

		// Calculate request duration
		duration := time.Since(start).Seconds()

		// Get request size
		requestSize := float64(r.ContentLength)
		if requestSize < 0 {
			requestSize = 0
		}

		// Record metrics
		HTTPRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(wrapped.status)).Inc()
		HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		HTTPRequestSize.WithLabelValues(method, endpoint).Observe(requestSize)
		HTTPResponseSize.WithLabelValues(method, endpoint).Observe(float64(wrapped.size))

		// Record API-specific response time
		APIResponseTime.WithLabelValues(endpoint, method).Observe(duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// UpdateAnomalyThresholds updates the current threshold values in Prometheus
func UpdateAnomalyThresholds(ncdThreshold, pValueThreshold float64, streamType string) {
	AnomalyNCDThreshold.WithLabelValues(streamType).Set(ncdThreshold)
	AnomalyPValueThreshold.WithLabelValues(streamType).Set(pValueThreshold)
}

// RecordAnomalyDetection records metrics for a detected anomaly
func RecordAnomalyDetection(streamType, severity string, processingTime time.Duration) {
	AnomaliesDetected.WithLabelValues(streamType, severity).Inc()
	AnomalyProcessingDuration.WithLabelValues(streamType).Observe(processingTime.Seconds())
}

// RecordCBADMetrics records metrics from the CBAD engine
func RecordCBADMetrics(streamType, dataType string, compressionRatio float64, eventsPerSecond float64) {
	CBADCompressionRatio.WithLabelValues(streamType, dataType).Observe(compressionRatio)
	CBADProcessingRate.WithLabelValues(streamType).Set(eventsPerSecond)
}

// RegisterMetricsHandler registers the Prometheus metrics handler
func RegisterMetricsHandler() http.Handler {
	return promhttp.Handler()
}