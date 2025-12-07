package driftlockcbad

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad/kafka"
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

func createLogsProcessor(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Logs) (processor.Logs, error) {
	c := cfg.(*Config)
	p := &cbadProcessor{cfg: *c, logger: set.Logger}
	p.baselineCap = 512

	// Initialize Redis client for distributed state if enabled
	if c.Redis.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
		})
		p.redisClient = redisClient
		set.Logger.Info("Redis client initialized for distributed state", 
			zap.String("addr", c.Redis.Addr))
	} else {
		set.Logger.Info("Redis disabled, using local state only")
	}

	// Initialize Kafka publisher if enabled
	if c.Kafka.Enabled {
		kafkaConfig := kafka.PublisherConfig{
			Brokers:  c.Kafka.Brokers,
			ClientID: c.Kafka.ClientID,
			EventsTopic: c.Kafka.EventsTopic,
			BatchSize: c.Kafka.BatchSize,
			BatchTimeout: time.Duration(c.Kafka.BatchTimeoutMs) * time.Millisecond,
		}

		// Configure TLS if enabled
		if c.Kafka.TLSEnabled {
			kafkaConfig.TLSConfig = &tls.Config{}
		}

		kafkaPublisher, err := kafka.NewPublisher(kafkaConfig, set.Logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kafka publisher: %w", err)
		}
		p.kafkaPublisher = kafkaPublisher
		set.Logger.Info("Kafka publisher initialized for OTLP events", 
			zap.String("topic", c.Kafka.EventsTopic), 
			zap.Strings("brokers", c.Kafka.Brokers))
	} else {
		set.Logger.Info("Kafka publisher disabled for OTLP events")
	}

	return &logProcessor{processor: p, nextConsumer: next}, nil
}

func createMetricsProcessor(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
	c := cfg.(*Config)
	p := &cbadProcessor{cfg: *c, logger: set.Logger}
	p.baselineCap = 512

	// Initialize Redis client for distributed state if enabled
	if c.Redis.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
		})
		p.redisClient = redisClient
		set.Logger.Info("Redis client initialized for distributed state", 
			zap.String("addr", c.Redis.Addr))
	} else {
		set.Logger.Info("Redis disabled, using local state only")
	}

	// Initialize Kafka publisher if enabled (reuse the same publisher for logs and metrics)
	if c.Kafka.Enabled {
		// The publisher should already be created in the logs processor
		// For now, we'll create a new one, but in production you might want to share
		kafkaConfig := kafka.PublisherConfig{
			Brokers:  c.Kafka.Brokers,
			ClientID: c.Kafka.ClientID,
			EventsTopic: c.Kafka.EventsTopic,
			BatchSize: c.Kafka.BatchSize,
			BatchTimeout: time.Duration(c.Kafka.BatchTimeoutMs) * time.Millisecond,
		}

		// Configure TLS if enabled
		if c.Kafka.TLSEnabled {
			kafkaConfig.TLSConfig = &tls.Config{}
		}

		kafkaPublisher, err := kafka.NewPublisher(kafkaConfig, set.Logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kafka publisher: %w", err)
		}
		p.kafkaPublisher = kafkaPublisher
		set.Logger.Info("Kafka publisher initialized for OTLP metrics", 
			zap.String("topic", c.Kafka.EventsTopic), 
			zap.Strings("brokers", c.Kafka.Brokers))
	}

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
	// Close the Kafka publisher if it exists
	if lp.processor.kafkaPublisher != nil {
		if err := lp.processor.kafkaPublisher.Close(); err != nil {
			return err
		}
	}
	
	// Close the Redis client if it exists
	if lp.processor.redisClient != nil {
		return lp.processor.redisClient.Close()
	}
	
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
	// Close the Kafka publisher if it exists
	if mp.processor.kafkaPublisher != nil {
		if err := mp.processor.kafkaPublisher.Close(); err != nil {
			return err
		}
	}
	
	// Close the Redis client if it exists
	if mp.processor.redisClient != nil {
		return mp.processor.redisClient.Close()
	}
	
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
