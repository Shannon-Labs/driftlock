package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// OTLPGenerator generates synthetic telemetry data for testing CBAD anomaly detection
type OTLPGenerator struct {
	collectorURL string
	client       *http.Client
	rng          *rand.Rand
}

// NewOTLPGenerator creates a new synthetic telemetry generator
func NewOTLPGenerator(collectorURL string) *OTLPGenerator {
	return &OTLPGenerator{
		collectorURL: collectorURL,
		client:       &http.Client{Timeout: 10 * time.Second},
		rng:          rand.New(rand.NewSource(42)), // Deterministic for reproducible tests
	}
}

// GenerateNormalLogs generates normal application logs
func (g *OTLPGenerator) GenerateNormalLogs(count int) plog.Logs {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	for i := 0; i < count; i++ {
		logRecord := scopeLogs.LogRecords().AppendEmpty()
		
		// Set timestamp
		logRecord.SetTimestamp(time.Now().Add(-time.Duration(i) * time.Second))
		
		// Set severity
		severities := []string{"INFO", "DEBUG", "TRACE"}
		severity := severities[g.rng.Intn(len(severities))]
		logRecord.SetSeverityText(severity)
		
		// Set body with structured log message
		message := fmt.Sprintf("Request completed successfully method=GET path=/api/users status=200 duration_ms=%d request_id=req-%d user_id=user-%d",
			40+g.rng.Intn(20), i, g.rng.Intn(1000))
		logRecord.Body().SetStr(message)
		
		// Add attributes
		logRecord.Attributes().PutStr("service.name", "api-gateway")
		logRecord.Attributes().PutStr("service.version", "1.0.0")
		logRecord.Attributes().PutStr("environment", "production")
		logRecord.Attributes().PutStr("region", "us-east-1")
		logRecord.Attributes().PutInt("request.size", g.rng.Intn(1024))
		logRecord.Attributes().PutDouble("response.time", float64(40+g.rng.Intn(20)))
	}

	return logs
}

// GenerateAnomalousLogs generates anomalous logs (errors, panics, unusual patterns)
func (g *OTLPGenerator) GenerateAnomalousLogs(count int) plog.Logs {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	for i := 0; i < count; i++ {
		logRecord := scopeLogs.LogRecords().AppendEmpty()
		
		// Set timestamp
		logRecord.SetTimestamp(time.Now().Add(-time.Duration(i) * time.Second))
		
		// Set severity - mostly ERROR for anomalous logs
		logRecord.SetSeverityText("ERROR")
		
		// Generate different types of anomalous messages
		var message string
		anomalyType := g.rng.Intn(4)
		switch anomalyType {
		case 0: // Stack trace
			message = fmt.Sprintf("PANIC: runtime error: index out of range [42] with length 10 goroutine 123 [running]: main.processRequest(0xc000123456, 0x10, 0xc000789012) /app/main.go:42 +0x123 github.com/gin-gonic/gin.(*Context).Next(0xc000345678) /go/pkg/mod/github.com/gin-gonic/gin@v1.7.7/context.go:123 +0x456")
		case 1: // Database error
			message = fmt.Sprintf("Database connection failed: connection refused to postgres://user:pass@db:5432/app?sslmode=disable dial tcp 10.0.0.1:5432: connect: connection refused request_id=req-%d", i)
		case 2: // Memory issue
			message = fmt.Sprintf("Out of memory: cannot allocate 1073741824 bytes (out of swap space) request_id=req-%d user_id=user-%d", i, g.rng.Intn(1000))
		case 3: // Timeout
			message = fmt.Sprintf("Request timeout: context deadline exceeded after 30.0s method=POST path=/api/expensive-operation request_id=req-%d", i)
		}
		
		logRecord.Body().SetStr(message)
		
		// Add attributes that might indicate the anomaly
		logRecord.Attributes().PutStr("service.name", "api-gateway")
		logRecord.Attributes().PutStr("error.type", "runtime_error")
		logRecord.Attributes().PutStr("error.severity", "critical")
		logRecord.Attributes().PutStr("environment", "production")
		logRecord.Attributes().PutStr("region", "us-east-1")
		
		// Add stack trace or error details
		if anomalyType == 0 {
			logRecord.Attributes().PutStr("stack_trace", message)
			logRecord.Attributes().PutStr("error.class", "panic")
		} else if anomalyType == 1 {
			logRecord.Attributes().PutStr("database.error", "connection_refused")
			logRecord.Attributes().PutStr("database.host", "10.0.0.1:5432")
		} else if anomalyType == 2 {
			logRecord.Attributes().PutStr("memory.error", "out_of_memory")
			logRecord.Attributes().PutInt("memory.requested", 1073741824)
		} else if anomalyType == 3 {
			logRecord.Attributes().PutStr("timeout.error", "deadline_exceeded")
			logRecord.Attributes().PutInt("timeout.duration_ms", 30000)
		}
	}

	return logs
}

// GenerateNormalMetrics generates normal application metrics
func (g *OTLPGenerator) GenerateNormalMetrics(count int) pmetric.Metrics {
	metrics := pmetric.NewMetrics()
	resourceMetrics := metrics.ResourceMetrics().AppendEmpty()
	scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()

	for i := 0; i < count; i++ {
		metric := scopeMetrics.Metrics().AppendEmpty()
		metric.SetName("http_request_duration")
		metric.SetDescription("HTTP request duration in milliseconds")
		metric.SetUnit("ms")
		
		// Create histogram
		histogram := metric.SetEmptyHistogram()
		histogram.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
		
		// Add data points with normal distribution
		dp := histogram.DataPoints().AppendEmpty()
		dp.SetTimestamp(time.Now().Add(-time.Duration(i) * time.Second))
		
		// Normal request duration around 40-60ms
		duration := float64(40 + g.rng.Intn(
