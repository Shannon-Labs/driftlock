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
		filepath.Join("api", "migrations"),
		filepath.Join("..", "api", "migrations"),
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(base, "api", "migrations"),
			filepath.Join(base, "..", "api", "migrations"))
	}
	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	return filepath.Join("api", "migrations")
}
