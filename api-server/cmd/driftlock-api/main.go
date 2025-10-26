package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Hmbown/driftlock/api-server/internal/api"
    "github.com/Hmbown/driftlock/api-server/internal/engine"
    "github.com/Hmbown/driftlock/api-server/internal/telemetry"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    ctx := context.Background()
    shutdown, err := telemetry.Setup(ctx)
    if err != nil {
        log.Fatalf("failed telemetry setup: %v", err)
    }
    defer func() {
        // Allow up to 5s to flush spans.
        c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := shutdown(c); err != nil {
            log.Printf("otel shutdown error: %v", err)
        }
    }()

    e := engine.New(nil) // Use minimal engine for simple version
    mux := api.NewMux(e)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    // Handle signals for graceful shutdown
    go func() {
        log.Printf("driftlockd listening on :%s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    <-sigc
    log.Printf("shutdown signal received")

    c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(c); err != nil {
        log.Printf("server shutdown error: %v", err)
    }
}
