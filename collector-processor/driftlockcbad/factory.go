package driftlockcbad

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor"
)

var (
	typeStr   = component.MustNewType("driftlock_cbad")
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
	return &logProcessor{processor: p, nextConsumer: next}, nil
}

func createMetricsProcessor(_ context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
	c := cfg.(*Config)
	p := &cbadProcessor{cfg: *c, logger: set.Logger}
	return &metricProcessor{processor: p, nextConsumer: next}, nil
}

// logProcessor wraps the cbadProcessor to implement the processor.Logs interface
type logProcessor struct {
	processor    *cbadProcessor
	nextConsumer consumer.Logs
}

// Capabilities returns the capabilities of the processor
func (lp *logProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start starts the processor
func (lp *logProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

// Shutdown stops the processor
func (lp *logProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ConsumeLogs processes the logs and passes them to the next consumer
func (lp *logProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	processedLogs, err := lp.processor.processLogs(ctx, ld)
	if err != nil {
		return err
	}
	if processedLogs.ResourceLogs().Len() == 0 {
		return nil // Nothing to pass to the next consumer
	}
	return lp.nextConsumer.ConsumeLogs(ctx, processedLogs)
}

// metricProcessor wraps the cbadProcessor to implement the processor.Metrics interface
type metricProcessor struct {
	processor    *cbadProcessor
	nextConsumer consumer.Metrics
}

// Capabilities returns the capabilities of the processor
func (mp *metricProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start starts the processor
func (mp *metricProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

// Shutdown stops the processor
func (mp *metricProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ConsumeMetrics processes the metrics and passes them to the next consumer
func (mp *metricProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	processedMetrics, err := mp.processor.processMetrics(ctx, md)
	if err != nil {
		return err
	}
	return mp.nextConsumer.ConsumeMetrics(ctx, processedMetrics)
}
