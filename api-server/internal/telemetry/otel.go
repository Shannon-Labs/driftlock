package telemetry

import (
    "context"
    "os"
    "strings"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Setup configures global OpenTelemetry and returns a shutdown function.
// If OTEL_EXPORTER_OTLP_ENDPOINT is not set, a no-op tracer provider is installed.
func Setup(ctx context.Context) (func(context.Context) error, error) {
    // Always set a propagator so downstream can parse incoming trace headers.
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{}, propagation.Baggage{},
    ))

    endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
    service := os.Getenv("OTEL_SERVICE_NAME")
    if service == "" {
        service = "driftlockd"
    }
    env := os.Getenv("OTEL_ENV")
    if env == "" {
        env = "dev"
    }
    version := os.Getenv("DRIFTLOCK_VERSION")

    res, err := resource.New(ctx,
        resource.WithSchemaURL(semconv.SchemaURL),
        resource.WithAttributes(
            semconv.ServiceName(service),
            semconv.ServiceVersion(version),
            semconv.DeploymentEnvironment(env),
        ),
    )
    if err != nil {
        return nil, err
    }

    // If no endpoint is provided, install a basic no-op provider; shutdown is a no-op.
    if endpoint == "" {
        tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res))
        otel.SetTracerProvider(tp)
        return func(context.Context) error { return nil }, nil
    }

    // Configure OTLP/HTTP exporter.
    opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
    if strings.HasPrefix(endpoint, "http://") || strings.Contains(endpoint, ":4318") {
        opts = append(opts, otlptracehttp.WithInsecure())
    }
    exp, err := otlptracehttp.New(ctx, opts...)
    if err != nil {
        return nil, err
    }

    bsp := sdktrace.NewBatchSpanProcessor(exp)
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithSpanProcessor(bsp),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)

    return tp.Shutdown, nil
}
