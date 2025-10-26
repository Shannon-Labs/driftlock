package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad_DefaultStreamingConfig(t *testing.T) {
	// Ensure no interfering env vars
	os.Unsetenv("KAFKA_ENABLED")
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_CLIENT_ID")
	os.Unsetenv("KAFKA_GROUP_ID")
	os.Unsetenv("KAFKA_EVENTS_TOPIC")
	os.Unsetenv("KAFKA_ANOMALIES_TOPIC")
	os.Unsetenv("KAFKA_TLS_ENABLED")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Streaming.Kafka.Enabled {
		t.Fatalf("expected Kafka to be disabled by default")
	}

	expectedBrokers := []string{"localhost:9092"}
	if !reflect.DeepEqual(cfg.Streaming.Kafka.Brokers, expectedBrokers) {
		t.Fatalf("expected brokers %v, got %v", expectedBrokers, cfg.Streaming.Kafka.Brokers)
	}

	if cfg.Streaming.Kafka.ClientID != "driftlock-api" {
		t.Fatalf("unexpected client id: %s", cfg.Streaming.Kafka.ClientID)
	}

	if cfg.Streaming.Kafka.EventsTopic != "otlp-events" {
		t.Fatalf("unexpected events topic: %s", cfg.Streaming.Kafka.EventsTopic)
	}
}

func TestLoad_StreamingOverrides(t *testing.T) {
	t.Setenv("KAFKA_ENABLED", "true")
	t.Setenv("KAFKA_BROKERS", "broker-1:9092, broker-2:9093")
	t.Setenv("KAFKA_CLIENT_ID", "driftlock-worker")
	t.Setenv("KAFKA_GROUP_ID", "cbad-consumers")
	t.Setenv("KAFKA_EVENTS_TOPIC", "telemetry.events")
	t.Setenv("KAFKA_ANOMALIES_TOPIC", "telemetry.anomalies")
	t.Setenv("KAFKA_TLS_ENABLED", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if !cfg.Streaming.Kafka.Enabled {
		t.Fatalf("expected Kafka to be enabled")
	}

	expectedBrokers := []string{"broker-1:9092", "broker-2:9093"}
	if !reflect.DeepEqual(cfg.Streaming.Kafka.Brokers, expectedBrokers) {
		t.Fatalf("expected brokers %v, got %v", expectedBrokers, cfg.Streaming.Kafka.Brokers)
	}

	if cfg.Streaming.Kafka.ClientID != "driftlock-worker" {
		t.Fatalf("unexpected client id: %s", cfg.Streaming.Kafka.ClientID)
	}

	if cfg.Streaming.Kafka.GroupID != "cbad-consumers" {
		t.Fatalf("unexpected group id: %s", cfg.Streaming.Kafka.GroupID)
	}

	if cfg.Streaming.Kafka.EventsTopic != "telemetry.events" {
		t.Fatalf("unexpected events topic: %s", cfg.Streaming.Kafka.EventsTopic)
	}

	if cfg.Streaming.Kafka.AnomaliesTopic != "telemetry.anomalies" {
		t.Fatalf("unexpected anomalies topic: %s", cfg.Streaming.Kafka.AnomaliesTopic)
	}

	if !cfg.Streaming.Kafka.TLSEnabled {
		t.Fatalf("expected TLS to be enabled")
	}
}
