package engine

import (
    "context"
    "log"

    "go.opentelemetry.io/otel"
)

type Engine struct{}

func New() *Engine { return &Engine{} }

// Process consumes an opaque JSON payload and performs minimal work for now.
// Extend this to run actual detection/aggregation logic.
func (e *Engine) Process(ctx context.Context, payload []byte) error {
    tracer := otel.Tracer("driftlock/engine")
    ctx, span := tracer.Start(ctx, "engine.Process")
    defer span.End()

    // Stub behavior: log and return success.
    if len(payload) == 0 {
        log.Printf("engine: received empty payload")
        return nil
    }
    log.Printf("engine: processed %d bytes", len(payload))
    return nil
}
