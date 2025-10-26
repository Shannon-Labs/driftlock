package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// API request metrics
	requestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
		[]string{"endpoint", "method"},
	)

	// Anomaly detection metrics
	anomalyCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "anomalies_detected_total",
			Help: "Total number of anomalies detected",
		},
		[]string{"type", "severity", "tenant_id"},
	)

	anomalyProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "anomaly_processing_duration_seconds",
			Help:    "Duration of anomaly processing in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		},
		[]string{"query_type", "table"},
	)

	// Kafka metrics
	kafkaMessageCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_processed_total",
			Help: "Total number of Kafka messages processed",
		},
		[]string{"topic", "partition"},
	)
)

// MetricsMiddleware adds Prometheus metrics to HTTP requests
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(c.FullPath(), c.Request.Method, c.Writer.Status()).Inc()
		requestDuration.WithLabelValues(c.FullPath(), c.Request.Method).Observe(duration)
	}
}

// RecordAnomalyDetected records a detected anomaly
func RecordAnomalyDetected(anomalyType, severity, tenantID string) {
	anomalyCount.WithLabelValues(anomalyType, severity, tenantID).Inc()
}

// RecordAnomalyProcessingDuration records the duration of anomaly processing
func RecordAnomalyProcessingDuration(duration float64) {
	anomalyProcessingDuration.Observe(duration)
}

// RecordDBQuery records database query metrics
func RecordDBQuery(queryType, table string, duration float64) {
	dbQueryDuration.WithLabelValues(queryType, table).Observe(duration)
}

// RecordKafkaMessage records a processed Kafka message
func RecordKafkaMessage(topic, partition string) {
	kafkaMessageCount.WithLabelValues(topic, partition).Inc()
}

// SetupMetricsEndpoint sets up the Prometheus metrics endpoint
func SetupMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}