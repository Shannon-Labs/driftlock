package driftlockcbad

import (
	"go.opentelemetry.io/collector/component"
)

// KafkaConfig holds Apache Kafka settings for the processor
type KafkaConfig struct {
	Enabled        bool     `mapstructure:"enabled"`
	Brokers        []string `mapstructure:"brokers"`
	ClientID       string   `mapstructure:"client_id"`
	EventsTopic    string   `mapstructure:"events_topic"`
	TLSEnabled     bool     `mapstructure:"tls_enabled"`
	BatchSize      int      `mapstructure:"batch_size"`
	BatchTimeoutMs int      `mapstructure:"batch_timeout_ms"`
}

// RedisConfig holds Redis settings for distributed state management
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
}

// Config defines the configuration for the driftlock CBAD processor
type Config struct {
	// Embed the component.Config interface to satisfy the component configuration contract
	component.Config `mapstructure:",squash"`

	WindowSize  int     `mapstructure:"window_size"`
	HopSize     int     `mapstructure:"hop_size"`
	Threshold   float64 `mapstructure:"threshold"`
	Determinism bool    `mapstructure:"determinism"`

	// Kafka configuration for publishing OTLP events
	Kafka KafkaConfig `mapstructure:"kafka"`

	// Redis configuration for distributed state management
	Redis RedisConfig `mapstructure:"redis"`
}
