package main

import (
    "context"
    "database/sql"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "sort"
    "time"

    _ "github.com/lib/pq"
)

func main() {
    var dir string
    flag.StringVar(&dir, "dir", "api-server/internal/storage/migrations", "migrations directory")
    flag.Parse()

    dsn := getenv("DATABASE_URL", "")
    if dsn == "" {
        host := getenv("DB_HOST", "localhost")
        port := getenv("DB_PORT", "5432")
        user := getenv("DB_USER", "postgres")
        pass := getenv("DB_PASSWORD", "")
        dbname := getenv("DB_DATABASE", "driftlock")
        sslmode := getenv("DB_SSL_MODE", "disable")
        dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    defer db.Close()

    if err := db.PingContext(ctx); err != nil {
        log.Fatalf("ping db: %v", err)
    }

    files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
    if err != nil {
        log.Fatalf("list migrations: %v", err)
    }
    sort.Strings(files)
    if len(files) == 0 {
        log.Printf("no migration files found in %s", dir)
        return
    }

    for _, f := range files {
        log.Printf("applying migration: %s", filepath.Base(f))
        b, err := ioutil.ReadFile(f)
        if err != nil {
            log.Fatalf("read %s: %v", f, err)
        }
        if _, err := db.ExecContext(ctx, string(b)); err != nil {
            log.Fatalf("execute %s: %v", f, err)
        }
    }
    log.Printf("migrations applied successfully")
}

func getenv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}

