package driftlockcbad

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestProcessorWithoutKafka(t *testing.T) {
	// Create a config with Kafka disabled
	cfg := &Config{
		ComponentConfig: component.NewDefaultConfig(),
		WindowSize:      1024,
		HopSize:         256,
		Threshold:       0.9,
		Determinism:     true,
		Kafka: KafkaConfig{
			Enabled: false, // Kafka disabled
		},
	}

	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a consumer to pass to the processor
	nextConsumer := &consumertest.LogsSink{}

	// Create the processor using the factory function
	ctx := context.Background()
	set := componenttest.NewNopProcessorCreateSettings()
	set.Logger = logger

	logProcessor, err := createLogsProcessor(ctx, set, cfg, nextConsumer)
	if err != nil {
		t.Fatalf("Failed to create logs processor: %v", err)
	}

	// Verify that the processor was created successfully
	if logProcessor == nil {
		t.Fatal("Expected log processor to be created, got nil")
	}

	// Test that the processor can process logs without Kafka
	logs := plog.NewLogs()
	rl := logs.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	lr := sl.LogRecords().AppendEmpty()
	lr.Body().SetStr("Test log message")

	// Process the logs
	err = logProcessor.ConsumeLogs(ctx, logs)
	if err != nil {
		t.Fatalf("Failed to consume logs: %v", err)
	}

	// Verify that logs were processed
	if nextConsumer.LogRecordsCount() != 1 {
		t.Errorf("Expected 1 log record, got %d", nextConsumer.LogRecordsCount())
	}

	// Now test metrics processor as well
	metricConsumer := &consumertest.MetricsSink{}
	metricProcessor, err := createMetricsProcessor(ctx, set, cfg, metricConsumer)
	if err != nil {
		t.Fatalf("Failed to create metrics processor: %v", err)
	}

	// Test that the processor can process metrics without Kafka
	metrics := pmetric.NewMetrics()
	rm := metrics.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	metric := sm.Metrics().AppendEmpty()
	metric.SetName("test.metric")
	gauge := metric.SetEmptyGauge()
	dp := gauge.DataPoints().AppendEmpty()
	dp.SetDoubleValue(42.0)

	// Process the metrics
	err = metricProcessor.ConsumeMetrics(ctx, metrics)
	if err != nil {
		t.Fatalf("Failed to consume metrics: %v", err)
	}

	// Verify that metrics were processed
	if metricConsumer.MetricsCount() != 1 {
		t.Errorf("Expected 1 metric, got %d", metricConsumer.MetricsCount())
	}
}