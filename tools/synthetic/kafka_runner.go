package main

import (
	"context"
	"log"

	driftlockkafka "github.com/shannon-labs/driftlock/collector-processor/driftlockcbad/kafka"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func runKafkaMode(config Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	publisher, err := driftlockkafka.NewPublisher(driftlockkafka.PublisherConfig{
		Brokers:      config.KafkaBrokers,
		ClientID:     config.KafkaClient,
		EventsTopic:  config.KafkaTopic,
		BatchSize:    config.NormalBatch,
		BatchTimeout: config.BatchInterval,
	}, logger)
	if err != nil {
		log.Fatalf("failed to create kafka publisher: %v", err)
	}
	defer func() {
		if cerr := publisher.Close(); cerr != nil {
			log.Printf("failed to close kafka publisher: %v", cerr)
		}
	}()

	generator := NewOTLPGenerator(config.CollectorURL)
	ctx := context.Background()

	var published int

	if config.NormalBatch > 0 {
		published += publishLogs(ctx, publisher, generator.GenerateNormalLogs(config.NormalBatch))
		published += publishMetrics(ctx, publisher, generator.GenerateNormalMetrics(config.NormalBatch))
	}

	if config.AnomalousBatch > 0 {
		published += publishLogs(ctx, publisher, generator.GenerateAnomalousLogs(config.AnomalousBatch))
		published += publishMetrics(ctx, publisher, generator.GenerateAnomalousMetrics(config.AnomalousBatch))
	}

	log.Printf("Published %d OTLP events to Kafka topic %s", published, config.KafkaTopic)
}

func publishLogs(ctx context.Context, publisher *driftlockkafka.Publisher, logs plog.Logs) int {
	count := 0
	resourceLogs := logs.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		scopeLogs := resourceLogs.At(i).ScopeLogs()
		for j := 0; j < scopeLogs.Len(); j++ {
			logRecords := scopeLogs.At(j).LogRecords()
			for k := 0; k < logRecords.Len(); k++ {
				if err := publisher.PublishLog(ctx, logRecords.At(k)); err != nil {
					log.Printf("failed to publish log record: %v", err)
					continue
				}
				count++
			}
		}
	}
	return count
}

func publishMetrics(ctx context.Context, publisher *driftlockkafka.Publisher, metrics pmetric.Metrics) int {
	count := 0
	resourceMetrics := metrics.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		scopeMetrics := resourceMetrics.At(i).ScopeMetrics()
		for j := 0; j < scopeMetrics.Len(); j++ {
			metricSlice := scopeMetrics.At(j).Metrics()
			for k := 0; k < metricSlice.Len(); k++ {
				if err := publisher.PublishMetric(ctx, metricSlice.At(k)); err != nil {
					log.Printf("failed to publish metric: %v", err)
					continue
				}
				count++
			}
		}
	}
	return count
}
