package llmio

import (
    "context"

    "go.opentelemetry.io/collector/component"
    "go.uber.org/zap"
)

type Config struct{}

type llmReceiver struct {
    logger    *zap.Logger
    nextLogs  any // replace with consumer.Logs when wiring
    nextTraces any // replace with consumer.Traces when wiring
}

func (r *llmReceiver) Start(ctx context.Context, _ component.Host) error {
    r.logger.Info("driftlock_llmio receiver starting")
    return nil
}

func (r *llmReceiver) Shutdown(ctx context.Context) error {
    r.logger.Info("driftlock_llmio receiver shutting down")
    return nil
}
