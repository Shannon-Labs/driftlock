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
    "github.com/Hmbown/driftlock/api-server/internal/config"
    "github.com/Hmbown/driftlock/api-server/internal/cbad"
    "github.com/Hmbown/driftlock/api-server/internal/auth"
    "github.com/Hmbown/driftlock/api-server/internal/engine"
    "github.com/Hmbown/driftlock/api-server/internal/stream"
    "github.com/Hmbown/driftlock/api-server/internal/storage"
    "github.com/Hmbown/driftlock/api-server/internal/supabase"
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

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Set up storage (Postgres)
    var store *storage.Storage
    if dsn := cfg.GetDatabaseConnectionString(); dsn != "" {
        store, err = storage.NewPostgres(dsn)
        if err != nil {
            log.Printf("warning: failed to connect to Postgres (continuing without DB): %v", err)
        }
    }

    // Set up SSE streamer
    streamer := stream.NewStreamer(0)

    // Initialize Supabase client if configured
    var sbClient *supabase.Client
    if cfg.Supabase.BaseURL != "" {
        sbClient, err = supabase.NewClient(supabase.Config{
            ProjectID:      cfg.Supabase.ProjectID,
            AnonKey:        cfg.Supabase.AnonKey,
            ServiceRoleKey: cfg.Supabase.ServiceRoleKey,
            BaseURL:        cfg.Supabase.BaseURL,
            RedisAddr:      cfg.Cache.Redis.Addr,
            RedisPassword:  cfg.Cache.Redis.Password,
            RedisDB:        cfg.Cache.Redis.DB,
        })
        if err != nil {
            log.Printf("warning: failed to init Supabase client: %v", err)
        }
    }

    // Initialize engine; attempt to wire CBAD detector when storage ready
    var e *engine.Engine
    if store != nil {
        if det, derr := cbad.NewDetector(store, streamer, sbClient); derr == nil {
            e = engine.New(det)
        } else {
            log.Printf("warning: CBAD detector not initialized: %v", derr)
            e = engine.New(nil)
        }
    } else {
        e = engine.New(nil)
    }

    // Build HTTP mux with optional dependencies
    var deps *api.Deps
    if store != nil {
        deps = &api.Deps{Storage: store, Streamer: streamer, Supabase: sbClient}
    }
    // Optional API key auth for ingestion route
    if deps == nil {
        deps = &api.Deps{}
    }
    if key := os.Getenv("DEFAULT_API_KEY"); key != "" {
        org := os.Getenv("DEFAULT_ORG_ID")
        a := auth.NewAuthenticator()
        a.AddAPIKey(auth.HashAPIKey(key), auth.APIKeyInfo{
            Name:          "default",
            Role:          "ingest",
            OrganizationID: org,
        })
        deps.IngestMiddleware = a.OptionalMiddleware
    }
    mux := api.NewMuxWithDeps(e, deps)

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

    // Cleanup
    if streamer != nil {
        _ = streamer.Shutdown(c)
    }
    if store != nil {
        _ = store.Close()
    }
}
