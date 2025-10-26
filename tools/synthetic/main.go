package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	defaultAPIURL          = "http://localhost:8080"
	defaultEventsPerSecond = 10
	defaultAnomalyRate     = 0.05
	defaultTenantID        = "demo"
)

type Config struct {
	APIURL          string
	EventsPerSecond int
	AnomalyRate     float64
	TenantID        string
}

func main() {
	// Parse configuration
	config := loadConfig()

	// Initialize OpenTelemetry
	initTracer(config.TenantID)

	// Start data generation
	log.Printf("Starting synthetic data generation for tenant: %s", config.TenantID)
	log.Printf("API URL: %s", config.APIURL)
	log.Printf("Events per second: %d", config.EventsPerSecond)
	log.Printf("Anomaly rate: %.2f", config.AnomalyRate)

	// Generate data continuously
	ticker := time.NewTicker(time.Second / time.Duration(config.EventsPerSecond))
	defer ticker.Stop()

	for range ticker.C {
		generateAndSendEvent(config)
	}
}

func loadConfig() Config {
	config := Config{
		APIURL:          getEnv("API_URL", defaultAPIURL),
		EventsPerSecond: getEnvInt("EVENTS_PER_SECOND", defaultEventsPerSecond),
		AnomalyRate:     getEnvFloat("ANOMALY_RATE", defaultAnomalyRate),
		TenantID:        getEnv("TENANT_ID", defaultTenantID),
	}

	return config
}

func generateAndSendEvent(config Config) {
	// Generate event data
	event := generateEvent(config.AnomalyRate)

	// Send event to API
	sendEvent(config.APIURL, event)
}

func generateEvent(anomalyRate float64) map[string]interface{} {
	// Determine if this should be an anomaly
	isAnomaly := rand.Float64() < anomalyRate

	// Generate base event
	event := map[string]interface{}{
		"timestamp":  time.Now().Unix(),
		"service":    generateRandomService(),
		"operation":  generateRandomOperation(),
		"duration":   rand.Intn(1000) + 1, // 1-1000ms
		"status":     generateRandomStatus(),
		"user_id":    fmt.Sprintf("user-%d", rand.Intn(1000)+1),
		"session_id": fmt.Sprintf("session-%d", rand.Intn(1000)+1),
		"ip_address": generateRandomIP(),
		"user_agent": generateRandomUserAgent(),
	}

	// Add anomaly data if needed
	if isAnomaly {
		event["anomaly"] = map[string]interface{}{
			"type":        generateRandomAnomalyType(),
			"severity":    generateRandomSeverity(),
			"confidence":  rand.Float64(),
			"description": generateRandomAnomalyDescription(),
			"indicators":  generateRandomIndicators(),
		}
	}

	return event
}

func sendEvent(apiURL string, event map[string]interface{}) {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL+"/v1/events", bytes.NewBuffer(eventJSON))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "demo")

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send event: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
	}
}

func initTracer(tenantID string) {
	// Create exporter
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint("http://localhost:4317"))
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	// Create a new tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			"",
			attribute.String("service.name", "synthetic-data-generator"),
			attribute.String("service.version", "1.0.0"),
			attribute.String("tenant.id", tenantID),
		)),
	)

	// Register tracer provider globally
	otel.SetTracerProvider(tp)
}

// Helper functions for generating random data

func generateRandomService() string {
	services := []string{
		"api-gateway",
		"user-service",
		"payment-service",
		"inventory-service",
		"order-service",
		"auth-service",
		"notification-service",
		"analytics-service",
		"search-service",
	}
	return services[rand.Intn(len(services))]
}

func generateRandomOperation() string {
	operations := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
	}
	return operations[rand.Intn(len(operations))]
}

func generateRandomStatus() string {
	statuses := []string{
		"success",
		"error",
		"timeout",
	}
	return statuses[rand.Intn(len(statuses))]
}

func generateRandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func generateRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"curl/7.68.0",
		"PostmanRuntime/7.28.4",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

func generateRandomAnomalyType() string {
	anomalyTypes := []string{
		"statistical",
		"behavioral",
		"contextual",
		"collective",
		"temporal",
		"spatial",
		"graph",
	}
	return anomalyTypes[rand.Intn(len(anomalyTypes))]
}

func generateRandomSeverity() string {
	severities := []string{
		"low",
		"medium",
		"high",
		"critical",
	}
	return severities[rand.Intn(len(severities))]
}

func generateRandomAnomalyDescription() string {
	descriptions := []string{
		"Unusual spike in request rate",
		"Abnormal response time pattern",
		"Unexpected error rate increase",
		"Irregular access pattern detected",
		"Statistical deviation from baseline",
		"Unusual geographic distribution",
		"Atypical user behavior pattern",
		"Anomalous resource utilization",
		"Unusual sequence of operations",
	}
	return descriptions[rand.Intn(len(descriptions))]
}

func generateRandomIndicators() []map[string]interface{} {
	indicators := make([]map[string]interface{}, 0)

	// Generate 2-5 indicators
	numIndicators := rand.Intn(4) + 2
	for i := 0; i < numIndicators; i++ {
		indicator := map[string]interface{}{
			"name":      generateRandomIndicatorName(),
			"value":     rand.Float64() * 100,
			"threshold": rand.Float64() * 100,
			"unit":      generateRandomIndicatorUnit(),
		}
		indicators = append(indicators, indicator)
	}

	return indicators
}

func generateRandomIndicatorName() string {
	indicatorNames := []string{
		"request_rate",
		"response_time",
		"error_rate",
		"cpu_usage",
		"memory_usage",
		"network_io",
		"disk_io",
		"cache_hit_rate",
		"queue_depth",
	}
	return indicatorNames[rand.Intn(len(indicatorNames))]
}

func generateRandomIndicatorUnit() string {
	units := []string{
		"requests/sec",
		"ms",
		"percent",
		"MB/s",
		"IOPS",
		"count",
	}
	return units[rand.Intn(len(units))]
}

// Environment variable helpers

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
