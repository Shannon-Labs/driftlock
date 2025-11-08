package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// migrations contains the SQL migrations in order
var migrations = []struct {
	up   string
	down string
}{
	{
		up: `CREATE TABLE IF NOT EXISTS anomalies (
			id UUID PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			event_type VARCHAR(50) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			message TEXT NOT NULL,
			explanation TEXT,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		down: "DROP TABLE IF EXISTS anomalies",
	},
	{
		up: `CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			event_type VARCHAR(50) NOT NULL,
			data JSONB NOT NULL,
			processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		down: "DROP TABLE IF EXISTS events",
	},
	{
		up: `CREATE TABLE IF NOT EXISTS config (
			key VARCHAR(100) PRIMARY KEY,
			value JSONB NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		down: "DROP TABLE IF EXISTS config",
	},
	{
		up: `CREATE TABLE IF NOT EXISTS compression_metrics (
			id UUID PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			baseline_compression_ratio DOUBLE PRECISION,
			window_compression_ratio DOUBLE PRECISION,
			ncd DOUBLE PRECISION,
			p_value DOUBLE PRECISION,
			is_anomaly BOOLEAN,
			confidence_level DOUBLE PRECISION,
			metadata JSONB
		)`,
		down: "DROP TABLE IF EXISTS compression_metrics",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <up|down> [step]\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]
	step := -1
	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &step)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_DATABASE", "driftlock"),
		)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Ensure migration tracking table exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`); err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	switch command {
	case "up":
		runMigrationsUp(db, step)
	case "down":
		runMigrationsDown(db, step)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runMigrationsUp(db *sql.DB, step int) {
	currentVersion := getCurrentVersion(db)
	fmt.Printf("Current schema version: %d\n", currentVersion)

	for i := currentVersion; i < len(migrations) && (step < 0 || i < currentVersion+step); i++ {
		fmt.Printf("Applying migration %d...\n", i+1)
		if _, err := db.Exec(migrations[i].up); err != nil {
			log.Fatalf("Failed to apply migration %d: %v", i+1, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", i+1); err != nil {
			log.Fatalf("Failed to record migration %d: %v", i+1, err)
		}
		fmt.Printf("✓ Migration %d applied\n", i+1)
	}

	fmt.Printf("Schema version is now: %d\n", getCurrentVersion(db))
}

func runMigrationsDown(db *sql.DB, step int) {
	currentVersion := getCurrentVersion(db)
	fmt.Printf("Current schema version: %d\n", currentVersion)

	if currentVersion == 0 {
		fmt.Println("No migrations to rollback")
		return
	}

	end := currentVersion - step
	if step < 0 {
		end = 0
	}

	for i := currentVersion - 1; i >= end && i >= 0; i-- {
		fmt.Printf("Rolling back migration %d...\n", i+1)
		if _, err := db.Exec(migrations[i].down); err != nil {
			log.Fatalf("Failed to rollback migration %d: %v", i+1, err)
		}
		if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", i+1); err != nil {
			log.Fatalf("Failed to remove migration record %d: %v", i+1, err)
		}
		fmt.Printf("✓ Migration %d rolled back\n", i+1)
	}

	fmt.Printf("Schema version is now: %d\n", getCurrentVersion(db))
}

func getCurrentVersion(db *sql.DB) int {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		log.Printf("Warning: could not get current version: %v", err)
		return 0
	}
	return version
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}