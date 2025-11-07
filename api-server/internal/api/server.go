package api

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "context"
    "strings"
    "time"

    "github.com/shannon-labs/driftlock/api-server/internal/ctxutil"
    "github.com/shannon-labs/driftlock/api-server/internal/engine"
    "github.com/shannon-labs/driftlock/api-server/internal/handlers"
    "github.com/shannon-labs/driftlock/api-server/internal/stream"
    "github.com/shannon-labs/driftlock/api-server/internal/streaming"
    "github.com/shannon-labs/driftlock/api-server/internal/storage"
    "github.com/shannon-labs/driftlock/api-server/internal/supabase"
    "github.com/shannon-labs/driftlock/pkg/version"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Deps are optional dependencies for wiring feature routes.
type Deps struct {
    Storage  storage.StorageInterface
    Streamer *stream.Streamer
    Events   streaming.EventPublisher
    Supabase *supabase.Client
    // IngestMiddleware wraps the /v1/events route (e.g., API key auth)
    IngestMiddleware func(http.Handler) http.Handler
}

// NewMux builds the http mux with instrumented handlers.
func NewMux(e *engine.Engine) http.Handler {
    return NewMuxWithDeps(e, nil)
}

// NewMuxWithDeps builds mux with optional anomalies + SSE routes.
func NewMuxWithDeps(e *engine.Engine, d *Deps) http.Handler {
    mux := http.NewServeMux()

    mux.Handle("/healthz", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    }), "healthz"))

    mux.Handle("/readyz", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // If storage provided, ensure DB is reachable
        if d != nil && d.Storage != nil {
            ctx, cancel := context.WithTimeout(r.Context(), time.Second)
            defer cancel()
            if err := d.Storage.Ping(ctx); err != nil {
                http.Error(w, "not ready", http.StatusServiceUnavailable)
                return
            }
        }
        // Best-effort Supabase health (non-blocking)
        if d != nil && d.Supabase != nil {
            ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
            defer cancel()
            if err := d.Supabase.HealthCheck(ctx); err != nil {
                log.Printf("readyz: supabase health failed: %v", err)
                // Do not fail readiness solely due to Supabase
            }
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ready"))
    }), "readyz"))

    mux.Handle("/v1/version", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]string{
            "version": version.Version(),
        })
    }), "version"))

    // Events ingestion
    eventsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "failed to read body", http.StatusBadRequest)
            return
        }
        defer func() { _ = r.Body.Close() }()
        // Extract org/event metadata if JSON
        var raw map[string]interface{}
        if err := json.Unmarshal(body, &raw); err == nil {
            var orgID, evType string
            if v, ok := raw["organization_id"].(string); ok {
                orgID = v
            }
            if v, ok := raw["event_type"].(string); ok {
                evType = v
            }
            if orgID != "" || evType != "" {
                r = r.Clone(ctxutil.WithEventContext(r.Context(), orgID, evType))
            }
        }
        start := time.Now()
        if err := e.Process(r.Context(), body); err != nil {
            log.Printf("events: processing error: %v", err)
            http.Error(w, "processing error", http.StatusInternalServerError)
            return
        }
        took := time.Since(start)
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]any{
            "status":          "ok",
            "processed_bytes": len(body),
            "latency":         took.String(),
        })
    })
    var wrapped http.Handler = eventsHandler
    if d != nil && d.IngestMiddleware != nil {
        wrapped = d.IngestMiddleware(wrapped)
    }
    mux.Handle("/v1/events", otelhttp.NewHandler(wrapped, "events"))

    // Optional: anomalies + SSE routes when deps provided
    if d != nil && d.Storage != nil {
        var ah *handlers.AnomaliesHandler
        if d.Supabase != nil {
            ah = handlers.NewAnomaliesHandlerWithSupabase(d.Storage, d.Streamer, d.Events, d.Supabase)
        } else {
            ah = handlers.NewAnomaliesHandler(d.Storage, d.Streamer, d.Events)
        }

        mux.Handle("/v1/anomalies", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            switch r.Method {
            case http.MethodGet:
                ah.ListAnomalies(w, r)
            case http.MethodPost:
                ah.CreateAnomaly(w, r)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        }), "anomalies"))

        mux.Handle("/v1/anomalies/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !strings.HasPrefix(r.URL.Path, "/v1/anomalies/") {
                http.NotFound(w, r)
                return
            }
            switch r.Method {
            case http.MethodGet:
                ah.GetAnomaly(w, r)
            case http.MethodPatch:
                ah.UpdateAnomalyStatus(w, r)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        }), "anomaly_by_id"))

        if d.Streamer != nil {
            mux.Handle("/v1/stream/anomalies", otelhttp.NewHandler(d.Streamer, "sse_anomalies"))
        }
    }

    return mux
}
