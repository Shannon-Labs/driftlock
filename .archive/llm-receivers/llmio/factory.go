package llmio

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver"
)

const (
    typeStr   = "driftlock_llmio"
    stability = component.StabilityLevelDevelopment
)

func NewFactory() receiver.Factory {
    return receiver.NewFactory(
        typeStr,
        createDefaultConfig,
        receiver.WithLogs(createLogsReceiver, stability),
        receiver.WithTraces(createTracesReceiver, stability),
    )
}

func createDefaultConfig() component.Config { return &Config{} }

func createLogsReceiver(_ context.Context, set receiver.Settings, cfg component.Config, next consumer.Logs) (receiver.Logs, error) {
    r := &llmReceiver{logger: set.Logger, nextLogs: next}
    return r, nil
}

func createTracesReceiver(_ context.Context, set receiver.Settings, cfg component.Config, next consumer.Traces) (receiver.Traces, error) {
    r := &llmReceiver{logger: set.Logger, nextTraces: next}
    return r, nil
}

