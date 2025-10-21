package driftlockcbad

import (
    "context"
    "fmt"

    "go.uber.org/zap"
)

type cbadProcessor struct {
    cfg    Config
    logger *zap.Logger
}

// TODO: Replace `any` with pdata types when wiring fully.
func (p *cbadProcessor) processLogs(ctx context.Context, logs any) (any, error) {
    _ = ctx
    p.logger.Debug("driftlock_cbad.processLogs invoked")
    // TODO: Extract records, compute CBAD metrics, attach explanations.
    return logs, nil
}

func (p *cbadProcessor) processMetrics(ctx context.Context, metrics any) (any, error) {
    _ = ctx
    p.logger.Debug("driftlock_cbad.processMetrics invoked")
    // TODO: Extract series, compute CBAD metrics, attach explanations.
    return metrics, nil
}

// Example router rationale (to be recorded per stream)
func (p *cbadProcessor) route(streamKind string) string {
    switch streamKind {
    case "logs":
        return "entropy->zstd_ratio->ncd"
    case "metrics":
        return "delta_bits->lz4_ratio->ncd"
    default:
        return fmt.Sprintf("default-router(%s)", streamKind)
    }
}

