package driftlockcbad

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type cbadProcessor struct {
	cfg           Config
	logger        *zap.Logger
	baselineReady bool
	baseline      []logSample
	window        []logSample
}

type logSample struct {
	resource pcommon.Resource
	scope    pcommon.InstrumentationScope
	record   plog.LogRecord
	data     []byte
}

type anomalyBatch struct {
	samples []logSample
	metrics *Metrics
}

// processLogs processes OTLP log data and detects anomalies using CBAD
func (p *cbadProcessor) processLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	_ = ctx
	p.logger.Debug("driftlock_cbad.processLogs invoked", zap.Int("resource_logs", logs.ResourceLogs().Len()))

	// Build baseline from historical patterns (simplified - in production, this would come from a database)
	baseline := p.buildLogBaseline()

	// Process each log record
	logs.ResourceLogs().RemoveIf(func(rl plog.ResourceLogs) bool {
		rl.ScopeLogs().RemoveIf(func(sl plog.ScopeLogs) bool {
			sl.LogRecords().RemoveIf(func(lr plog.LogRecord) bool {
				logData := p.logRecordToBytes(lr)
				metrics, err := ComputeMetricsQuick(baseline, logData)
				if err != nil {
					p.logger.Warn("CBAD analysis failed", zap.Error(err))
					return false // Keep the log record if analysis fails
				}

				if metrics.IsAnomaly {
					p.addAnomalyAttributes(lr, metrics)
					p.logger.Info("Anomaly detected in log",
						zap.String("explanation", metrics.GetAnomalyExplanation()),
						zap.Float64("ncd", metrics.NCD),
						zap.Float64("p_value", metrics.PValue),
					)
					return false // Keep the anomalous log record
				}

				return true // Remove the non-anomalous log record
			})
			return sl.LogRecords().Len() == 0
		})
		return rl.ScopeLogs().Len() == 0
	})

	return logs, nil
}

// processMetrics processes OTLP metric data and detects anomalies using CBAD
func (p *cbadProcessor) processMetrics(ctx context.Context, metrics pmetric.Metrics) (pmetric.Metrics, error) {
	_ = ctx
	p.logger.Debug("driftlock_cbad.processMetrics invoked", zap.Int("resource_metrics", metrics.ResourceMetrics().Len()))

	// Build baseline from historical patterns
	baseline := p.buildMetricBaseline()

	// Process each metric
	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		resourceMetrics := metrics.ResourceMetrics().At(i)
		for j := 0; j < resourceMetrics.ScopeMetrics().Len(); j++ {
			scopeMetrics := resourceMetrics.ScopeMetrics().At(j)
			for k := 0; k < scopeMetrics.Metrics().Len(); k++ {
				metric := scopeMetrics.Metrics().At(k)

				// Convert metric to byte array for CBAD analysis
				metricData := p.metricToBytes(metric)

				// Compute CBAD metrics
				cbadMetrics, err := ComputeMetricsQuick(baseline, metricData)
				if err != nil {
					p.logger.Warn("CBAD analysis failed for metric", zap.Error(err))
					continue
				}

				// Add CBAD attributes to metric if anomaly detected
				if cbadMetrics.IsAnomaly {
					p.addMetricAnomalyAttributes(metric, cbadMetrics)
					p.logger.Info("Anomaly detected in metric",
						zap.String("name", metric.Name()),
						zap.String("explanation", cbadMetrics.GetAnomalyExplanation()),
						zap.Float64("ncd", cbadMetrics.NCD),
						zap.Float64("p_value", cbadMetrics.PValue),
					)
				}
			}
		}
	}

	return metrics, nil
}

// buildLogBaseline creates a baseline for log anomaly detection
func (p *cbadProcessor) buildLogBaseline() []byte {
	logData := map[string]interface{}{
		"timestamp": "2025-10-24T00:00:00Z",
		"severity":  "INFO",
		"body":      "Request completed successfully",
		"attributes": map[string]interface{}{
			"service.name": "test-service",
		},
	}
	jsonData, _ := json.Marshal(logData)
	return jsonData
}

// buildMetricBaseline creates a baseline for metric anomaly detection
func (p *cbadProcessor) buildMetricBaseline() []byte {
	// In a real implementation, this would come from historical data
	// For now, use a representative normal metric pattern
	return []byte(`{"name":"http_request_duration","type":"histogram","value":42.0,"labels":{"service":"api-gateway","method":"GET","status":"200"}}`)
}

// logRecordToBytes converts a log record to bytes for CBAD analysis
func (p *cbadProcessor) logRecordToBytes(logRecord plog.LogRecord) []byte {
	// Create a JSON representation of the log record
	logData := map[string]interface{}{
		"timestamp": logRecord.Timestamp().String(),
		"severity":  logRecord.SeverityText(),
		"body":      logRecord.Body().AsString(),
	}

	// Add attributes
	attrs := make(map[string]interface{})
	logRecord.Attributes().Range(func(k string, v pcommon.Value) bool {
		attrs[k] = v.AsString()
		return true
	})
	if len(attrs) > 0 {
		logData["attributes"] = attrs
	}

	jsonData, _ := json.Marshal(logData)
	return jsonData
}

// metricToBytes converts a metric to bytes for CBAD analysis
func (p *cbadProcessor) metricToBytes(metric pmetric.Metric) []byte {
	// Create a JSON representation of the metric
	metricData := map[string]interface{}{
		"name": metric.Name(),
		"type": metric.Type().String(),
	}

	// Add description and unit if available
	if desc := metric.Description(); desc != "" {
		metricData["description"] = desc
	}
	if unit := metric.Unit(); unit != "" {
		metricData["unit"] = unit
	}

	jsonData, _ := json.Marshal(metricData)
	return jsonData
}

// addAnomalyAttributes adds CBAD anomaly attributes to a log record
func (p *cbadProcessor) addAnomalyAttributes(logRecord plog.LogRecord, metrics *Metrics) {
	// Add CBAD anomaly detection attributes
	logRecord.Attributes().PutStr("driftlock.anomaly_detected", "true")
	logRecord.Attributes().PutDouble("driftlock.ncd", metrics.NCD)
	logRecord.Attributes().PutDouble("driftlock.p_value", metrics.PValue)
	logRecord.Attributes().PutDouble("driftlock.confidence_level", metrics.ConfidenceLevel)
	logRecord.Attributes().PutStr("driftlock.explanation", metrics.GetAnomalyExplanation())

	// Add compression metrics
	logRecord.Attributes().PutDouble("driftlock.baseline_compression_ratio", metrics.BaselineCompressionRatio)
	logRecord.Attributes().PutDouble("driftlock.window_compression_ratio", metrics.WindowCompressionRatio)
	logRecord.Attributes().PutDouble("driftlock.baseline_entropy", metrics.BaselineEntropy)
	logRecord.Attributes().PutDouble("driftlock.window_entropy", metrics.WindowEntropy)

	// Add statistical significance flag
	if metrics.IsStatisticallySignificant() {
		logRecord.Attributes().PutStr("driftlock.statistically_significant", "true")
	}
}

// addMetricAnomalyAttributes adds CBAD anomaly attributes to a metric
func (p *cbadProcessor) addMetricAnomalyAttributes(metric pmetric.Metric, metrics *Metrics) {
	// Add CBAD anomaly detection as a description
	explanation := metrics.GetDetailedExplanation()
	if len(explanation) > 1024 { // Truncate if too long
		explanation = explanation[:1024] + "..."
	}
	metric.SetDescription(explanation)
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
