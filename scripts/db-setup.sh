#!/bin/bash
set -e

# Driftlock Database Setup Script
# Usage: ./scripts/db-setup.sh [database_url]

DB_URL=${1:-${DATABASE_URL}}

if [ -z "$DB_URL" ]; then
  echo "Error: DATABASE_URL is not set."
  echo "Usage: ./scripts/db-setup.sh <postgres_connection_string>"
  echo "Example: ./scripts/db-setup.sh postgres://user:pass@localhost:5432/driftlock?sslmode=disable"
  exit 1
fi

echo "Checking database connection..."
if ! command -v psql &> /dev/null; then
    echo "Warning: psql not found. Skipping connection check."
else
    if ! psql "$DB_URL" -c "\q" 2>/dev/null; then
        echo "Error: Could not connect to database."
        exit 1
    fi
fi

echo "Applying migrations..."

# Check for goose
if ! command -v goose &> /dev/null; then
    echo "Installing goose..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Navigate to migrations directory
cd api/migrations

echo "Running migrations against $DB_URL"
goose postgres "$DB_URL" up

echo "Database setup complete!"
