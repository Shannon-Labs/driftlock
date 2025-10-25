package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	CBAD      CBADConfig
	Auth      AuthConfig
	Observability ObservabilityConfig
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
