package driftlockcbad

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad/kafka"
)

type cbadProcessor struct {
	cfg            Config
	logger         *zap.Logger
	baselineReady  bool
	baseline       []logSample
	window         []logSample
	detector       *Detector    // New streaming detector
	detectorMu     sync.RWMutex // Protect detector access
	kafkaPublisher *kafka.Publisher
	redisClient    *redis.Client // For distributed state management
	baselineMu     sync.Mutex
	logBaseline    [][]byte
	metricBaseline [][]byte
	baselineCap    int
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

	// Build baseline from historical patterns (updated after each record)
	baseline := p.buildLogBaseline()

	rls := logs.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rl := rls.At(i)
		sls := rl.ScopeLogs()
		for j := 0; j < sls.Len(); j++ {
			sl := sls.At(j)
			lrs := sl.LogRecords()
			for k := 0; k < lrs.Len(); k++ {
				lr := lrs.At(k)

				// Publish the log record to Kafka if publisher is available
				if p.kafkaPublisher != nil {
					if err := p.kafkaPublisher.PublishLog(ctx, lr); err != nil {
						p.logger.Error("Failed to publish log to Kafka", zap.Error(err))
					}
				}

				logData := p.logRecordToBytes(lr)
				metrics, err := ComputeMetricsQuick(baseline, logData)
				if err != nil {
					p.logger.Warn("CBAD analysis failed", zap.Error(err))
				} else if metrics.IsAnomaly {
					p.addAnomalyAttributes(lr, metrics)
					p.logger.Info("Anomaly detected in log",
						zap.String("explanation", metrics.GetAnomalyExplanation()),
						zap.Float64("ncd", metrics.NCD),
						zap.Float64("p_value", metrics.PValue),
					)
				}

				// Continuously evolve the baseline with observed traffic
				p.recordLogBaseline(logData)
				baseline = p.buildLogBaseline()
			}
		}
	}
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

				// Publish the metric to Kafka if publisher is available
				if p.kafkaPublisher != nil {
					if err := p.kafkaPublisher.PublishMetric(ctx, metric); err != nil {
						p.logger.Error("Failed to publish metric to Kafka", zap.Error(err))
					}
				}

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

				// Evolve baseline for future analyses
				p.recordMetricBaseline(metricData)
				baseline = p.buildMetricBaseline()
			}
		}
	}

	return metrics, nil
}

// buildLogBaseline creates a baseline for log anomaly detection
func (p *cbadProcessor) buildLogBaseline() []byte {
	p.baselineMu.Lock()
	defer p.baselineMu.Unlock()

	if len(p.logBaseline) == 0 {
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

	return flattenBuffers(p.logBaseline)
}

// buildMetricBaseline creates a baseline for metric anomaly detection
func (p *cbadProcessor) buildMetricBaseline() []byte {
	p.baselineMu.Lock()
	defer p.baselineMu.Unlock()

	if len(p.metricBaseline) == 0 {
		return []byte(`{"name":"http_request_duration","type":"histogram","value":42.0,"labels":{"service":"api-gateway","method":"GET","status":"200"}}`)
	}

	return flattenBuffers(p.metricBaseline)
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

func (p *cbadProcessor) recordLogBaseline(data []byte) {
	if p.baselineCap == 0 {
		p.baselineCap = 512
	}
	p.baselineMu.Lock()
	defer p.baselineMu.Unlock()
	p.logBaseline = append(p.logBaseline, append([]byte(nil), data...))
	if len(p.logBaseline) > p.baselineCap {
		excess := len(p.logBaseline) - p.baselineCap
		p.logBaseline = p.logBaseline[excess:]
	}
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

	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		points := make([]map[string]interface{}, 0, metric.Gauge().DataPoints().Len())
		for i := 0; i < metric.Gauge().DataPoints().Len(); i++ {
			points = append(points, datapointToMap(metric.Gauge().DataPoints().At(i)))
		}
		metricData["datapoints"] = points
	case pmetric.MetricTypeSum:
		points := make([]map[string]interface{}, 0, metric.Sum().DataPoints().Len())
		for i := 0; i < metric.Sum().DataPoints().Len(); i++ {
			points = append(points, datapointToMap(metric.Sum().DataPoints().At(i)))
		}
		metricData["datapoints"] = points
	case pmetric.MetricTypeHistogram:
		points := make([]map[string]interface{}, 0, metric.Histogram().DataPoints().Len())
		for i := 0; i < metric.Histogram().DataPoints().Len(); i++ {
			dp := metric.Histogram().DataPoints().At(i)
			m := map[string]interface{}{
				"count":   dp.Count(),
				"sum":     dp.Sum(),
				"min":     dp.Min(),
				"max":     dp.Max(),
				"bounds":  dp.ExplicitBounds().AsRaw(),
				"bucket":  dp.BucketCounts().AsRaw(),
				"attrs":   attributesToMap(dp.Attributes()),
				"ts":      dp.Timestamp().String(),
				"startTs": dp.StartTimestamp().String(),
			}
			points = append(points, m)
		}
		metricData["datapoints"] = points
	default:
		// Fallback: note the type without discarding the metric
		metricData["note"] = "unsupported metric type payload preserved"
	}

	jsonData, _ := json.Marshal(metricData)
	return jsonData
}

func (p *cbadProcessor) recordMetricBaseline(data []byte) {
	if p.baselineCap == 0 {
		p.baselineCap = 512
	}
	p.baselineMu.Lock()
	defer p.baselineMu.Unlock()
	p.metricBaseline = append(p.metricBaseline, append([]byte(nil), data...))
	if len(p.metricBaseline) > p.baselineCap {
		excess := len(p.metricBaseline) - p.baselineCap
		p.metricBaseline = p.metricBaseline[excess:]
	}
}

func flattenBuffers(chunks [][]byte) []byte {
	total := 0
	for _, c := range chunks {
		total += len(c)
	}
	out := make([]byte, 0, total)
	for _, c := range chunks {
		out = append(out, c...)
	}
	return out
}

func attributesToMap(attrs pcommon.Map) map[string]interface{} {
	out := make(map[string]interface{}, attrs.Len())
	attrs.Range(func(k string, v pcommon.Value) bool {
		switch v.Type() {
		case pcommon.ValueTypeDouble:
			out[k] = v.Double()
		case pcommon.ValueTypeInt:
			out[k] = v.Int()
		case pcommon.ValueTypeBool:
			out[k] = v.Bool()
		default:
			out[k] = v.AsString()
		}
		return true
	})
	return out
}

func datapointToMap(dp pmetric.NumberDataPoint) map[string]interface{} {
	m := map[string]interface{}{
		"ts":      dp.Timestamp().String(),
		"startTs": dp.StartTimestamp().String(),
		"attrs":   attributesToMap(dp.Attributes()),
	}
	switch dp.ValueType() {
	case pmetric.NumberDataPointValueTypeDouble:
		m["value"] = dp.DoubleValue()
	case pmetric.NumberDataPointValueTypeInt:
		m["value"] = dp.IntValue()
	}
	return m
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
