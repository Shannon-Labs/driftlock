package driftlockcbad

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap/zaptest"
)

// Integration-style checks for baseline evolution and data preservation.
func TestOTLPProcessorIntegration(t *testing.T) {
	if err := ValidateLibrary(); err != nil {
		t.Skipf("CBAD core unavailable in test environment: %v", err)
	}

	tests := []struct {
		name         string
		events       []string
		expectedAnom bool
	}{
		{
			name: "normal_logs_no_anomaly",
			events: []string{
				"User login successful",
				"User login successful",
				"User login successful",
				"User login successful",
			},
			expectedAnom: false,
		},
		{
			name: "pattern_change_anomaly",
			events: []string{
				"User login successful",
				"User login successful",
				"PANIC: database connection failed! stacktrace=AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNNOOOOPPPPQQQQRRRR",
				"User login successful",
			},
			expectedAnom: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestProcessor(t)
			logs := buildLogs(tt.events)

			out, err := p.processLogs(context.Background(), logs)
			if err != nil {
				t.Fatalf("processLogs error: %v", err)
			}

			// Ensure data flows through unchanged in count.
			if got, want := countLogRecords(out), len(tt.events); got != want {
				t.Fatalf("expected %d log records to flow through, got %d", want, got)
			}

			anomalyCount := countAnomalies(out)
			if tt.expectedAnom && anomalyCount == 0 {
				t.Fatalf("expected anomaly detection, found none")
			}
			if !tt.expectedAnom && anomalyCount > 0 {
				t.Fatalf("unexpected anomaly detection: %d", anomalyCount)
			}

			// Baseline should evolve with observed events.
			p.baselineMu.Lock()
			defer p.baselineMu.Unlock()
			if len(p.logBaseline) == 0 {
				t.Fatalf("baseline did not record any entries")
			}
		})
	}
}

func TestMetricPayloadsPreserved(t *testing.T) {
	if err := ValidateLibrary(); err != nil {
		t.Skipf("CBAD core unavailable in test environment: %v", err)
	}

	p := newTestProcessor(t)

	metrics := pmetric.NewMetrics()
	rm := metrics.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("test_gauge")
	m.SetDescription("desc")
	m.SetUnit("ms")
	g := m.SetEmptyGauge()

	dp := g.DataPoints().AppendEmpty()
	dp.SetDoubleValue(1.5)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(testNow()))
	dp.Attributes().PutStr("service", "svc-a")

	dp2 := g.DataPoints().AppendEmpty()
	dp2.SetDoubleValue(2.5)
	dp2.Attributes().PutStr("service", "svc-b")

	out, err := p.processMetrics(context.Background(), metrics)
	if err != nil {
		t.Fatalf("processMetrics error: %v", err)
	}

	// Verify datapoints and attributes are preserved.
	rmOut := out.ResourceMetrics().At(0)
	smOut := rmOut.ScopeMetrics().At(0)
	mOut := smOut.Metrics().At(0)
	if mOut.Gauge().DataPoints().Len() != 2 {
		t.Fatalf("expected 2 datapoints, got %d", mOut.Gauge().DataPoints().Len())
	}
	seen := map[string]bool{}
	for i := 0; i < mOut.Gauge().DataPoints().Len(); i++ {
		dp := mOut.Gauge().DataPoints().At(i)
		if v, ok := dp.Attributes().Get("service"); ok {
			seen[v.Str()] = true
		}
	}
	if !seen["svc-a"] || !seen["svc-b"] {
		t.Fatalf("missing preserved attributes in datapoints: %v", seen)
	}
}

// Helpers

func newTestProcessor(t *testing.T) *cbadProcessor {
	t.Helper()
	cfg := Config{
		WindowSize:  32,
		HopSize:     8,
		Threshold:   0.3,
		Determinism: true,
	}
	p := &cbadProcessor{
		cfg:         cfg,
		logger:      zaptest.NewLogger(t),
		baselineCap: 256,
	}
	return p
}

func buildLogs(messages []string) plog.Logs {
	logs := plog.NewLogs()
	rl := logs.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	for _, msg := range messages {
		lr := sl.LogRecords().AppendEmpty()
		lr.Body().SetStr(msg)
		lr.SetSeverityText("INFO")
		lr.Attributes().PutStr("svc", "svc-A")
	}
	return logs
}

func countLogRecords(logs plog.Logs) int {
	count := 0
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rl := logs.ResourceLogs().At(i)
		for j := 0; j < rl.ScopeLogs().Len(); j++ {
			sl := rl.ScopeLogs().At(j)
			count += sl.LogRecords().Len()
		}
	}
	return count
}

func countAnomalies(logs plog.Logs) int {
	count := 0
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rl := logs.ResourceLogs().At(i)
		for j := 0; j < rl.ScopeLogs().Len(); j++ {
			sl := rl.ScopeLogs().At(j)
			for k := 0; k < sl.LogRecords().Len(); k++ {
				lr := sl.LogRecords().At(k)
				if v, ok := lr.Attributes().Get("driftlock.anomaly_detected"); ok && v.AsString() == "true" {
					count++
				}
			}
		}
	}
	return count
}

func testNow() time.Time { return time.Unix(0, 0) }
