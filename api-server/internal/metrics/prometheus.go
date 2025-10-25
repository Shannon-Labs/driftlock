package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// AnomaliesDetected tracks total anomalies detected by stream type
	AnomaliesDetected = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "driftlock_anomalies_detected_total",
			Help: "Total number of anomalies detected",
		},
		[]string{"stream_type", "severity"},
	)

	// EventsProcessed tracks total events processed
	EventsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "driftlock_events_processed_total",
			Help: "Total number of events processed",
		},
		[]string{"stream_type"},
	)

	// CompressionRatio tracks compression ratios
	CompressionRatio = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "driftlock_compression_ratio",
			Help:    "Compression ratio distribution",
			Buckets: prometheus.LinearBuckets(1, 0.5, 10), // 1.0 to 5.5
		},
		[]string{"stream_type", "window_type"},
	)

	// NCDScore tracks NCD score distribution
	NCDScore = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "driftlock_ncd_score",
			Help:    "NCD score distribution",
			Buckets: prometheus.LinearBuckets(0, 0.1, 11), // 0.0 to 1.0
		},
	)

	// PValue tracks p-value distribution
	PValue = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "driftlock_p_value",
			Help:    "P-value distribution",
			Buckets: []float64{0.001, 0.01, 0.05, 0.1, 0.25, 0.5, 0.75, 1.0},
		},
	)

	// APIRequestDuration tracks API request latencies
	APIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "driftlock_api_request_duration_seconds",
			Help:    "API request duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10), // 1ms to ~1s
		},
		[]string{"endpoint", "method", "status"},
	)

	// ActiveSSEConnections tracks active SSE connections
	ActiveSSEConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "driftlock_active_sse_connections",
			Help: "Number of active SSE connections",
		},
	)

	// DatabaseConnections tracks database connection pool stats
	DatabaseConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "driftlock_database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"}, // open, in_use, idle
	)

	// CBADComputationDuration tracks CBAD computation time
	CBADComputationDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "driftlock_cbad_computation_duration_seconds",
			Help:    "CBAD computation duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		},
	)

	// DatabaseQueryDuration tracks database query latencies
	DatabaseQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "driftlock_database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		},
		[]string{"query_type"},
	)
)

// Register all metrics with Prometheus
func init() {
	prometheus.MustRegister(AnomaliesDetected)
	prometheus.MustRegister(EventsProcessed)
	prometheus.MustRegister(CompressionRatio)
	prometheus.MustRegister(NCDScore)
	prometheus.MustRegister(PValue)
	prometheus.MustRegister(APIRequestDuration)
	prometheus.MustRegister(ActiveSSEConnections)
	prometheus.MustRegister(DatabaseConnections)
	prometheus.MustRegister(CBADComputationDuration)
	prometheus.MustRegister(DatabaseQueryDuration)
}

// Handler returns the Prometheus metrics HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}

// RecordAnomaly records an anomaly detection
func RecordAnomaly(streamType, severity string) {
	AnomaliesDetected.WithLabelValues(streamType, severity).Inc()
}

// RecordEvent records an event processed
func RecordEvent(streamType string) {
	EventsProcessed.WithLabelValues(streamType).Inc()
}

// RecordCompressionRatio records a compression ratio
func RecordCompressionRatio(streamType, windowType string, ratio float64) {
	CompressionRatio.WithLabelValues(streamType, windowType).Observe(ratio)
}

// RecordNCDScore records an NCD score
func RecordNCDScore(score float64) {
	NCDScore.Observe(score)
}

// RecordPValue records a p-value
func RecordPValue(pvalue float64) {
	PValue.Observe(pvalue)
}

// InstrumentHandler wraps an HTTP handler with request duration metrics
func InstrumentHandler(endpoint string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		handler(rw, r)

		duration := time.Since(start).Seconds()
		APIRequestDuration.WithLabelValues(endpoint, r.Method, http.StatusText(rw.statusCode)).Observe(duration)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// UpdateDatabaseConnectionStats updates database connection pool metrics
func UpdateDatabaseConnectionStats(open, inUse, idle int) {
	DatabaseConnections.WithLabelValues("open").Set(float64(open))
	DatabaseConnections.WithLabelValues("in_use").Set(float64(inUse))
	DatabaseConnections.WithLabelValues("idle").Set(float64(idle))
}

// RecordCBADComputation records CBAD computation duration
func RecordCBADComputation(duration time.Duration) {
	CBADComputationDuration.Observe(duration.Seconds())
}

// RecordDatabaseQuery records database query duration
func RecordDatabaseQuery(queryType string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
}
