# Database Migrations

This directory contains SQL migration files for the Driftlock database schema.

## Migration Tool

We recommend using [golang-migrate](https://github.com/golang-migrate/migrate) for running migrations.

### Installation

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Running Migrations

```bash
# Set your database URL
export DATABASE_URL="postgres://username:password@localhost:5432/driftlock?sslmode=disable"

# Run all pending migrations
migrate -path ./migrations -database "$DATABASE_URL" up

# Rollback the last migration
migrate -path ./migrations -database "$DATABASE_URL" down 1

# Check current version
migrate -path ./migrations -database "$DATABASE_URL" version

# Force a specific version (use with caution)
migrate -path ./migrations -database "$DATABASE_URL" force 1
```

### Local Development with Docker

```bash
# Start PostgreSQL container
docker run --name driftlock-postgres \
  -e POSTGRES_DB=driftlock \
  -e POSTGRES_USER=driftlock \
  -e POSTGRES_PASSWORD=driftlock \
  -p 5432:5432 \
  -d postgres:15

# Wait for PostgreSQL to be ready
sleep 5

# Run migrations
export DATABASE_URL="postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable"
migrate -path ./migrations -database "$DATABASE_URL" up
```

## Migration Files

- `001_initial_schema.up.sql` - Creates the core tables:
  - `anomalies` - Stores detected anomalies with CBAD metrics
  - `detection_config` - Configuration for anomaly detection thresholds
  - `performance_metrics` - API performance monitoring
  - `audit_log` - Change tracking and audit trail

## Schema Overview

### Anomalies Table

Stores all detected anomalies with comprehensive CBAD metrics:

- **Core fields**: ID, timestamp, stream type
- **CBAD metrics**: NCD score, p-value, compression ratios
- **Status tracking**: pending/acknowledged/dismissed/investigating
- **User interaction**: acknowledgements, dismissals, notes
- **Data storage**: baseline/window data as JSONB

### Indexes

Optimized for common query patterns:
- Timestamp (DESC) - for recent anomalies
- Status - for filtering by workflow state
- Stream type - for per-stream queries
- Statistical significance - for filtering true anomalies
- Tags (GIN) - for tag-based searches
- Metadata (GIN) - for flexible JSON queries

### Views

- `recent_significant_anomalies` - Last 100 statistically significant anomalies
- `anomaly_stats_by_stream` - Aggregated statistics per stream type

## Best Practices

1. **Always test migrations locally first**
2. **Back up production database before running migrations**
3. **Review migration SQL before applying**
4. **Use transactions where possible** (migrate does this automatically)
5. **Never modify existing migration files** - create new ones instead

## Creating New Migrations

```bash
# Create a new migration pair
migrate create -ext sql -dir ./migrations -seq add_new_feature
```

This creates:
- `002_add_new_feature.up.sql` - Forward migration
- `002_add_new_feature.down.sql` - Rollback migration
