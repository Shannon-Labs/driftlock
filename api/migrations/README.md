# Driftlock API Migrations

Migrations use [goose](https://github.com/pressly/goose) with SQL files. Run them via the Driftlock API binary (planned `./driftlock-api migrate up`) or the goose CLI:

```bash
# Example (local Postgres from docker-compose)
goose -dir api/migrations postgres "postgres://driftlock:driftlock@localhost:5432/driftlock?sslmode=disable" up
```

Use timestamped filenames (UTC) and include both `Up` and `Down` sections.
