package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	CBAD          CBADConfig
	Auth          AuthConfig
	Observability ObservabilityConfig
	Streaming     StreamingConfig
	Cache         CacheConfig
	Storage       StorageConfig
	Tenant        TenantConfig
	Supabase      SupabaseConfig
}

// ServerConfig holds server settings
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database settings
type DatabaseConfig struct {
	Host           string
	Port           int
	Database       string
	User           string
	Password       string
	MaxConnections int
	SSLMode        string
}

// CBADConfig holds CBAD detection settings
type CBADConfig struct {
	NCDThreshold    float64
	PValueThreshold float64
	BaselineSize    int
	WindowSize      int
	HopSize         int
}

// AuthConfig holds authentication settings
type AuthConfig struct {
	Type     string // "apikey" or "oidc"
	Issuer   string
	ClientID string
}

// ObservabilityConfig holds observability settings
type ObservabilityConfig struct {
	PrometheusPort int
	LogLevel       string
	TraceSampling  float64
}

// RedisConfig holds Redis settings for distributed state management
type RedisConfig struct {
	Enabled  bool
	Addr     string
	Password string
	DB       int
	Prefix   string
}

// StreamingConfig holds stream processing settings
type StreamingConfig struct {
	Kafka KafkaConfig
}

// CompressionConfig holds storage compression configuration
type CompressionConfig struct {
	Enabled      bool
	Algorithm    string // openzl, zstd, lz4, gzip
	Level        int
	MinSizeBytes int // Only compress if data exceeds this size
}

// TieredStorageConfig holds hot/warm/cold storage configuration
type TieredStorageConfig struct {
	Enabled           bool
	HotRetentionDays  int
	WarmRetentionDays int
	ArchiveInterval   string // Duration string like "24h"
	Compression       CompressionConfig
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Tiered TieredStorageConfig
}

// CacheConfig holds cache settings
type CacheConfig struct {
	Redis RedisConfig
}

// KafkaConfig holds Apache Kafka settings
type KafkaConfig struct {
	Enabled        bool
	Brokers        []string
	ClientID       string
	GroupID        string
	EventsTopic    string
	AnomaliesTopic string
	TLSEnabled     bool
}

// TenantConfig holds multi-tenant settings
type TenantConfig struct {
	Enabled        bool
	IsolationMode  string // "schema" or "database"
	DefaultTenant  string
	ResourceQuotas TenantResourceQuotas
}

// TenantResourceQuotas holds per-tenant resource limits
type TenantResourceQuotas struct {
	MaxAnomaliesPerDay   int
	MaxEventsPerDay      int
	MaxStorageGB         int
	MaxAPIRequestsPerMin int
}

// SupabaseConfig holds Supabase settings
type SupabaseConfig struct {
	ProjectID      string
	AnonKey        string
	ServiceRoleKey string
	BaseURL        string
	WebhookURL     string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvInt("DB_PORT", 5432),
			Database:       getEnv("DB_DATABASE", "driftlock"),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", ""),
			MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 100),
			SSLMode:        getEnv("DB_SSL_MODE", "prefer"),
		},
		CBAD: CBADConfig{
			NCDThreshold:    getEnvFloat("CBAD_NCD_THRESHOLD", 0.3),
			PValueThreshold: getEnvFloat("CBAD_P_VALUE_THRESHOLD", 0.05),
			BaselineSize:    getEnvInt("CBAD_BASELINE_SIZE", 100),
			WindowSize:      getEnvInt("CBAD_WINDOW_SIZE", 50),
			HopSize:         getEnvInt("CBAD_HOP_SIZE", 10),
		},
		Auth: AuthConfig{
			Type:     getEnv("AUTH_TYPE", "apikey"),
			Issuer:   getEnv("AUTH_ISSUER", ""),
			ClientID: getEnv("AUTH_CLIENT_ID", ""),
		},
		Observability: ObservabilityConfig{
			PrometheusPort: getEnvInt("PROMETHEUS_PORT", 9090),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			TraceSampling:  getEnvFloat("TRACE_SAMPLING", 0.1),
		},
		Streaming: StreamingConfig{
			Kafka: KafkaConfig{
				Enabled:        getEnvBool("KAFKA_ENABLED", false),
				Brokers:        getEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
				ClientID:       getEnv("KAFKA_CLIENT_ID", "driftlock-api"),
				GroupID:        getEnv("KAFKA_GROUP_ID", "driftlock-api"),
				EventsTopic:    getEnv("KAFKA_EVENTS_TOPIC", "otlp-events"),
				AnomaliesTopic: getEnv("KAFKA_ANOMALIES_TOPIC", "anomaly-events"),
				TLSEnabled:     getEnvBool("KAFKA_TLS_ENABLED", false),
			},
		},
		Cache: CacheConfig{
			Redis: RedisConfig{
				Enabled:  getEnvBool("REDIS_ENABLED", false),
				Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       getEnvInt("REDIS_DB", 0),
				Prefix:   getEnv("REDIS_PREFIX", "driftlock"),
			},
		},
		Storage: StorageConfig{
			Tiered: TieredStorageConfig{
				Enabled:           getEnvBool("TIERED_STORAGE_ENABLED", false),
				HotRetentionDays:  getEnvInt("HOT_RETENTION_DAYS", 7),
				WarmRetentionDays: getEnvInt("WARM_RETENTION_DAYS", 83),
				ArchiveInterval:   getEnv("ARCHIVE_INTERVAL", "24h"),
				Compression: CompressionConfig{
					Enabled:      getEnvBool("COMPRESSION_ENABLED", true),
					Algorithm:    getEnv("COMPRESSION_ALGORITHM", "zstd"),
					Level:        getEnvInt("COMPRESSION_LEVEL", 3),
					MinSizeBytes: getEnvInt("COMPRESSION_MIN_SIZE", 1024),
				},
			},
		},
		Tenant: TenantConfig{
			Enabled:       getEnvBool("TENANT_ENABLED", false),
			IsolationMode: getEnv("TENANT_ISOLATION_MODE", "schema"),
			DefaultTenant: getEnv("TENANT_DEFAULT_TENANT", "default"),
			ResourceQuotas: TenantResourceQuotas{
				MaxAnomaliesPerDay:   getEnvInt("TENANT_MAX_ANOMALIES_PER_DAY", 1000),
				MaxEventsPerDay:      getEnvInt("TENANT_MAX_EVENTS_PER_DAY", 10000),
				MaxStorageGB:         getEnvInt("TENANT_MAX_STORAGE_GB", 100),
				MaxAPIRequestsPerMin: getEnvInt("TENANT_MAX_API_REQUESTS_PER_MIN", 100),
			},
		},
		Supabase: SupabaseConfig{
			ProjectID:      getEnv("SUPABASE_PROJECT_ID", ""),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			BaseURL:        getEnv("SUPABASE_BASE_URL", ""),
			WebhookURL:     getEnv("SUPABASE_WEBHOOK_URL", ""),
		},
	}

	return config, nil
}

// GetDatabaseConnectionString returns the PostgreSQL connection string
func (c *Config) GetDatabaseConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}

// Helper functions

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

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		var trimmed []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				trimmed = append(trimmed, part)
			}
		}
		if len(trimmed) > 0 {
			return trimmed
		}
	}
	return append([]string(nil), defaultValue...)
}
