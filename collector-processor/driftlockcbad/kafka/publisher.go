package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	segment "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// PublisherConfig describes how to connect to Kafka for publishing OTLP data.
type PublisherConfig struct {
	Brokers      []string
	ClientID     string
	TLSConfig    *tls.Config
	EventsTopic  string
	BatchSize    int
	BatchTimeout time.Duration
}

// OTLPEvent represents an OTLP event to be published to Kafka.
type OTLPEvent struct {
	Type      string      `json:"type"` // "log" or "metric"
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
}

// Publisher publishes OTLP events to Kafka.
type Publisher struct {
	dialer      *segment.Dialer
	cfg         PublisherConfig
	eventsTopic string
	writer      *segment.Writer
	writerMu    sync.Mutex
	logger      *zap.Logger
}

// NewPublisher creates a new Kafka publisher for OTLP events.
func NewPublisher(cfg PublisherConfig, logger *zap.Logger) (*Publisher, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka publisher: at least one broker is required")
	}
	if cfg.EventsTopic == "" {
		return nil, errors.New("kafka publisher: events topic is required")
	}
	if logger == nil {
		return nil, errors.New("kafka publisher: logger is required")
	}

	dialer := &segment.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  cfg.ClientID,
		TLS:       cfg.TLSConfig,
	}

	p := &Publisher{
		dialer:      dialer,
		cfg:         cfg,
		eventsTopic: cfg.EventsTopic,
		logger:      logger,
	}

	// Create the writer with default settings
	writerConfig := segment.WriterConfig{
		Brokers:      cfg.Brokers,
		Topic:        cfg.EventsTopic,
		Dialer:       dialer,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
	}

	if writerConfig.BatchSize == 0 {
		writerConfig.BatchSize = 100
	}
	if writerConfig.BatchTimeout == 0 {
		writerConfig.BatchTimeout = 5 * time.Millisecond
	}

	p.writer = segment.NewWriter(writerConfig)
	return p, nil
}

// PublishLog publishes an OTLP log record to Kafka.
func (p *Publisher) PublishLog(ctx context.Context, lr plog.LogRecord) error {
	logData := map[string]interface{}{
		"timestamp": lr.Timestamp().AsTime().Format(time.RFC3339Nano),
		"severity":  lr.SeverityText(),
		"body":      lr.Body().AsString(),
		"severity_number": int(lr.SeverityNumber()),
		"flags":     uint32(lr.Flags()),
	}

	// Add attributes
	attrs := make(map[string]interface{})
	lr.Attributes().Range(func(k string, v pcommon.Value) bool {
		attrs[k] = v.AsString()
		return true
	})
	if len(attrs) > 0 {
		logData["attributes"] = attrs
	}

	event := OTLPEvent{
		Type:      "log",
		Data:      logData,
		Timestamp: time.Now(),
		Source:    "collector-processor",
	}

	return p.publish(ctx, event)
}

// PublishMetric publishes an OTLP metric to Kafka.
func (p *Publisher) PublishMetric(ctx context.Context, metric pmetric.Metric) error {
	metricData := map[string]interface{}{
		"name":        metric.Name(),
		"description": metric.Description(),
		"unit":        metric.Unit(),
		"type":        metric.Type().String(),
	}

	// Add the metric data based on the metric type
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		gauge := metric.Gauge()
		pointsData := make([]map[string]interface{}, gauge.DataPoints().Len())
		for i := 0; i < gauge.DataPoints().Len(); i++ {
			dp := gauge.DataPoints().At(i)
			point := map[string]interface{}{
				"timestamp": dp.Timestamp().AsTime().Format(time.RFC3339Nano),
				"value":     dp.DoubleValue(),
			}
			// Add attributes for this data point
			dpAttrs := make(map[string]interface{})
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				dpAttrs[k] = v.AsString()
				return true
			})
			if len(dpAttrs) > 0 {
				point["attributes"] = dpAttrs
			}
			pointsData[i] = point
		}
		metricData["gauge"] = pointsData
	case pmetric.MetricTypeSum:
		sum := metric.Sum()
		pointsData := make([]map[string]interface{}, sum.DataPoints().Len())
		for i := 0; i < sum.DataPoints().Len(); i++ {
			dp := sum.DataPoints().At(i)
			point := map[string]interface{}{
				"timestamp": dp.Timestamp().AsTime().Format(time.RFC3339Nano),
				"value":     dp.DoubleValue(),
			}
			// Add attributes for this data point
			dpAttrs := make(map[string]interface{})
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				dpAttrs[k] = v.AsString()
				return true
			})
			if len(dpAttrs) > 0 {
				point["attributes"] = dpAttrs
			}
			pointsData[i] = point
		}
		metricData["sum"] = pointsData
	case pmetric.MetricTypeHistogram:
		hist := metric.Histogram()
		pointsData := make([]map[string]interface{}, hist.DataPoints().Len())
		for i := 0; i < hist.DataPoints().Len(); i++ {
			dp := hist.DataPoints().At(i)
			point := map[string]interface{}{
				"timestamp": dp.Timestamp().AsTime().Format(time.RFC3339Nano),
				"count":     dp.Count(),
				"sum":       dp.Sum(),
				"min":       dp.Min(),
				"max":       dp.Max(),
			}
			// Add attributes for this data point
			dpAttrs := make(map[string]interface{})
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				dpAttrs[k] = v.AsString()
				return true
			})
			if len(dpAttrs) > 0 {
				point["attributes"] = dpAttrs
			}
			pointsData[i] = point
		}
		metricData["histogram"] = pointsData
	case pmetric.MetricTypeSummary:
		summary := metric.Summary()
		pointsData := make([]map[string]interface{}, summary.DataPoints().Len())
		for i := 0; i < summary.DataPoints().Len(); i++ {
			dp := summary.DataPoints().At(i)
			point := map[string]interface{}{
				"timestamp": dp.Timestamp().AsTime().Format(time.RFC3339Nano),
				"count":     dp.Count(),
				"sum":       dp.Sum(),
			}
			// Add quantile values
			quantilesData := make([]map[string]interface{}, dp.QuantileValues().Len())
			for j := 0; j < dp.QuantileValues().Len(); j++ {
				qv := dp.QuantileValues().At(j)
				quantilesData[j] = map[string]interface{}{
					"quantile": qv.Quantile(),
					"value":    qv.Value(),
				}
			}
			if len(quantilesData) > 0 {
				point["quantiles"] = quantilesData
			}
			// Add attributes for this data point
			dpAttrs := make(map[string]interface{})
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				dpAttrs[k] = v.AsString()
				return true
			})
			if len(dpAttrs) > 0 {
				point["attributes"] = dpAttrs
			}
			pointsData[i] = point
		}
		metricData["summary"] = pointsData
	}

	event := OTLPEvent{
		Type:      "metric",
		Data:      metricData,
		Timestamp: time.Now(),
		Source:    "collector-processor",
	}

	return p.publish(ctx, event)
}

// publish serializes and sends the OTLP event to Kafka.
func (p *Publisher) publish(ctx context.Context, event OTLPEvent) error {
	p.writerMu.Lock()
	defer p.writerMu.Unlock()

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal OTLP event: %w", err)
	}

	// Create a unique key for this event to ensure consistent partitioning
	key := fmt.Sprintf("%s-%d", event.Type, time.Now().UnixNano())

	msg := segment.Message{
		Topic: p.eventsTopic,
		Key:   []byte(key),
		Value: jsonData,
		Time:  time.Now(),
		Headers: []segment.Header{
			{
				Key:   "driftlock-event-type",
				Value: []byte(event.Type),
			},
			{
				Key:   "driftlock-source",
				Value: []byte(event.Source),
			},
		},
	}

	return p.writer.WriteMessages(ctx, msg)
}

// Close closes the Kafka publisher and releases resources.
func (p *Publisher) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}