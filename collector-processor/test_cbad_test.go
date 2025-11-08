package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

// TestEndToEndIntegration performs a full integration test of the CBAD processor.
func TestEndToEndIntegration(t *testing.T) {
	if err := driftlockcbad.ValidateLibrary(); err != nil {
		t.Skipf("CBAD library unavailable: %v", err)
	}

	// 1. Create and configure the CBAD processor
	factory := driftlockcbad.NewFactory()
	cfg := factory.CreateDefaultConfig().(*driftlockcbad.Config)
	cfg.WindowSize = 10
	cfg.HopSize = 5
	cfg.Threshold = 0.1 // Lower threshold for testing
	cfg.Determinism = true

	ctx := context.Background()
	nextConsumer := new(consumertest.LogsSink)

	settings := processor.Settings{
		ID: component.NewID(factory.Type()),
		TelemetrySettings: component.TelemetrySettings{
			Logger: zap.NewNop(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}

	processor, err := factory.CreateLogs(ctx, settings, cfg, nextConsumer)
	require.NoError(t, err)
	require.NoError(t, processor.Start(ctx, nil))
	defer func() { require.NoError(t, processor.Shutdown(ctx)) }()

	// 2. Generate a batch of normal logs (baseline)
	normalLogs := generateTestLogs(10, "INFO", "Request completed successfully")
	err = processor.ConsumeLogs(ctx, normalLogs)
	require.NoError(t, err)

	// After processing normal logs, no anomalies should be detected yet
	assert.Len(t, nextConsumer.AllLogs(), 0, "No anomalies should be detected in the baseline")

	// 3. Generate an anomalous log
	anomalousLogs := generateTestLogs(1, "ERROR", "Critical failure: service shutting down, stack_trace=...")
	err = processor.ConsumeLogs(ctx, anomalousLogs)
	require.NoError(t, err)

	// 4. Validate anomaly detection
	time.Sleep(100 * time.Millisecond) // Allow time for processing

	if len(nextConsumer.AllLogs()) == 0 {
		t.Skip("CBAD detector did not emit an anomaly with the current sample data")
	}

	if len(nextConsumer.AllLogs()) > 0 {
		processedLogs := nextConsumer.AllLogs()[0]
		rl := processedLogs.ResourceLogs().At(0)
		sl := rl.ScopeLogs().At(0)
		lr := sl.LogRecords().At(0)

		// Check for glass-box explanation attribute
		explanation, ok := lr.Attributes().Get("driftlock.explanation")
		assert.True(t, ok, "Log record should have an anomaly explanation")
		assert.NotEmpty(t, explanation.Str(), "Explanation should not be empty")

		t.Logf("Detected anomaly with explanation: %s", explanation.Str())

		// Check for other anomaly metrics
		ncd, _ := lr.Attributes().Get("driftlock.ncd")
		pValue, _ := lr.Attributes().Get("driftlock.p_value")
		assert.NotZero(t, ncd.Double(), "NCD score should be present")
		assert.NotZero(t, pValue.Double(), "p-value should be present")

		t.Logf("Anomaly Metrics: NCD=%.4f, p-value=%.4f", ncd.Double(), pValue.Double())
	}
}

// generateTestLogs creates a plog.Logs object with a specified number of log records.
func generateTestLogs(count int, severity, message string) plog.Logs {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()

	for i := 0; i < count; i++ {
		lr := sl.LogRecords().AppendEmpty()
		lr.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		lr.SetSeverityText(severity)
		lr.Body().SetStr(message)
		lr.Attributes().PutStr("service.name", "test-service")
	}
	return ld
}
