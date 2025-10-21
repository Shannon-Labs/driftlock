package driftlockcbad

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/processor"
    "go.opentelemetry.io/collector/processor/processorhelper"
)

const (
    typeStr   = "driftlock_cbad"
    stability = component.StabilityLevelDevelopment
)

func NewFactory() processor.Factory {
    return processor.NewFactory(
        typeStr,
        createDefaultConfig,
        processor.WithLogs(createLogsProcessor, stability),
        processor.WithMetrics(createMetricsProcessor, stability),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        WindowSize:  1024,
        HopSize:     256,
        Threshold:   0.9,
        Determinism: true,
    }
}

func createLogsProcessor(_ context.Context, set processor.Settings, cfg component.Config, next consumer.Logs) (processor.Logs, error) {
    c := cfg.(*Config)
    p := &cbadProcessor{cfg: *c, logger: set.Logger}
    return processorhelper.NewLogsProcessor(set, cfg, next, p.processLogs)
}

func createMetricsProcessor(_ context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
    c := cfg.(*Config)
    p := &cbadProcessor{cfg: *c, logger: set.Logger}
    return processorhelper.NewMetricsProcessor(set, cfg, next, p.processMetrics)
}

