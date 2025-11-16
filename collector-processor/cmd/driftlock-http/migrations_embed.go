package main

import (
	"os"
	"path/filepath"
)

func migrationsPath() string {
	if v := env("MIGRATIONS_DIR", ""); v != "" {
		return v
	}
	candidates := []string{
		// Relative to current working directory (for local dev)
		filepath.Join("api", "migrations"),
		// Relative to collector-processor directory (for Docker build)
		filepath.Join("..", "api", "migrations"),
		// From repo root when running from collector-processor
		filepath.Join("../..", "api", "migrations"),
		// Docker runtime location (set via MIGRATIONS_DIR env var, but also check directly)
		"/usr/local/share/driftlock/migrations",
	}
	// Add executable-relative paths
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(base, "api", "migrations"),
			filepath.Join(base, "..", "api", "migrations"),
			filepath.Join(base, "../..", "api", "migrations"))
		// Also check relative to executable's parent (for /usr/local/bin/driftlock-http)
		if filepath.Base(base) == "bin" {
			candidates = append(candidates,
				filepath.Join(base, "..", "api", "migrations"))
		}
	}
	// Check each candidate path
	for _, path := range candidates {
		// Try to resolve to absolute path for more reliable checking
		absPath, err := filepath.Abs(path)
		if err == nil {
			path = absPath
		}
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	// Fallback: return default path (may not exist, but caller can handle error)
	return filepath.Join("api", "migrations")
}
