package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress      string
	Debug              bool
	DatabaseURL        string
	KafkaBrokers       []string
	JWTSecret          string
	AllowedOrigins     []string
	SessionSecret      string
	AnomalyThreshold   float64
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set default values
	config := &Config{
		ServerAddress:    getEnv("SERVER_ADDRESS", ":8080"),
		Debug:            getEnvAsBool("DEBUG", false),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost:5432/driftlock?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-change-in-production"),
		SessionSecret:    getEnv("SESSION_SECRET", "default-session-secret-change-in-production"),
		AnomalyThreshold: getEnvAsFloat64("ANOMALY_THRESHOLD", 1000.0),
		KafkaBrokers:     getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		AllowedOrigins:   getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),
	}

	return config
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return fallback
}

func getEnvAsFloat64(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return fallback
}