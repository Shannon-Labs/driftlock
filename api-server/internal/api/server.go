package api

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "time"

    "github.com/your-org/driftlock/api-server/internal/engine"
    "github.com/your-org/driftlock/pkg/version"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewMux builds the http mux with instrumented handlers.
func NewMux(e *engine.Engine) http.Handler {
    mux := http.NewServeMux()

    mux.Handle("/healthz", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    }), "healthz"))

    mux.Handle("/readyz", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ready"))
    }), "readyz"))

    mux.Handle("/v1/version", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]string{
            "version": version.Version(),
        })
    }), "version"))

    mux.Handle("/v1/events", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MiB cap
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "failed to read body", http.StatusBadRequest)
            return
        }
        defer func() { _ = r.Body.Close() }()

        start := time.Now()
        if err := e.Process(r.Context(), body); err != nil {
            log.Printf("events: processing error: %v", err)
            http.Error(w, "processing error", http.StatusInternalServerError)
            return
        }
        took := time.Since(start)
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]any{
            "status":  "ok",
            "latency": took.String(),
        })
    }), "events"))

    return mux
}
