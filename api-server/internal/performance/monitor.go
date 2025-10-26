package performance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Monitor tracks various performance metrics for Driftlock
type Monitor struct {
	meter          metric.Meter
	requestsTotal  metric.Int64Counter
	requestsFailed metric.Int64Counter
	requestLatency metric.Float64Histogram
	activeWorkers  metric.Int64UpDownCounter
	queueSize      metric.Int64UpDownCounter
	mu             sync.RWMutex
}

// NewMonitor creates a new performance monitor
func NewMonitor() *Monitor {
	meter := otel.Meter("driftlock.performance")
	
	requestsTotal, _ := meter.Int64Counter(
		"driftlock_requests_total",
		metric.WithDescription("Total number of requests processed"),
		metric.WithUnit("1"),
	)
	
	requestsFailed, _ := meter.Int64Counter(
		"driftlock_requests_failed",
		metric.WithDescription("Number of failed requests"),
		metric.WithUnit("1"),
	)
	
	requestLatency, _ := meter.Float64Histogram(
		"driftlock_request_duration_seconds",
		metric.WithDescription("Request latency in seconds"),
		metric.WithUnit("s"),
	)
	
	activeWorkers, _ := meter.Int64UpDownCounter(
		"driftlock_active_workers",
		metric.WithDescription("Number of active worker goroutines"),
		metric.WithUnit("1"),
	)
	
	queueSize, _ := meter.Int64UpDownCounter(
		"driftlock_queue_size",
		metric.WithDescription("Current queue size for processing"),
		metric.WithUnit("1"),
	)

	return &Monitor{
		meter:          meter,
		requestsTotal:  requestsTotal,
		requestsFailed: requestsFailed,
		requestLatency: requestLatency,
		activeWorkers:  activeWorkers,
		queueSize:      queueSize,
	}
}

// RecordRequest records metrics for a completed request
func (m *Monitor) RecordRequest(ctx context.Context, duration time.Duration, path, method string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("path", path),
		attribute.String("method", method),
	}

	m.requestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	
	if !success {
		m.requestsFailed.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	
	m.requestLatency.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// UpdateActiveWorkers updates the counter for active workers
func (m *Monitor) UpdateActiveWorkers(delta int64) {
	m.activeWorkers.Add(context.Background(), delta)
}

// UpdateQueueSize updates the gauge for queue size
func (m *Monitor) UpdateQueueSize(size int64) {
	m.queueSize.Add(context.Background(), size)
}

// BatchProcessorMetrics tracks metrics for a batch processing system
type BatchProcessorMetrics struct {
	monitor        *Monitor
	batchSize      metric.Int64Histogram
	processingTime metric.Float64Histogram
	throughput     metric.Float64Gauge
	batchesTotal   metric.Int64Counter
	batchesSuccess metric.Int64Counter
	batchesFailed  metric.Int64Counter
}

// NewBatchProcessorMetrics creates metrics for a batch processor
func NewBatchProcessorMetrics() *BatchProcessorMetrics {
	meter := otel.Meter("driftlock.batch_processor")
	
	batchSize, _ := meter.Int64Histogram(
		"driftlock_batch_size",
		metric.WithDescription("Size of each batch processed"),
		metric.WithUnit("1"),
	)
	
	processingTime, _ := meter.Float64Histogram(
		"driftlock_batch_processing_duration_seconds",
		metric.WithDescription("Time to process a batch"),
		metric.WithUnit("s"),
	)
	
	throughput, _ := meter.Float64Gauge(
		"driftlock_events_per_second",
		metric.WithDescription("Current processing throughput"),
		metric.WithUnit("1/s"),
	)
	
	batchesTotal, _ := meter.Int64Counter(
		"driftlock_batches_total",
		metric.WithDescription("Total number of batches processed"),
		metric.WithUnit("1"),
	)
	
	batchesSuccess, _ := meter.Int64Counter(
		"driftlock_batches_success",
		metric.WithDescription("Number of successful batches"),
		metric.WithUnit("1"),
	)
	
	batchesFailed, _ := meter.Int64Counter(
		"driftlock_batches_failed",
		metric.WithDescription("Number of failed batches"),
		metric.WithUnit("1"),
	)
	
	return &BatchProcessorMetrics{
		batchSize:      batchSize,
		processingTime: processingTime,
		throughput:     throughput,
		batchesTotal:   batchesTotal,
		batchesSuccess: batchesSuccess,
		batchesFailed:  batchesFailed,
	}
}

// RecordBatch records metrics for a completed batch
func (b *BatchProcessorMetrics) RecordBatch(ctx context.Context, size int, duration time.Duration, success bool) {
	attrs := []attribute.KeyValue{
		attribute.Int("size", size),
	}

	b.batchesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	
	if success {
		b.batchesSuccess.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		b.batchesFailed.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	
	b.batchSize.Record(ctx, int64(size), metric.WithAttributes(attrs...))
	b.processingTime.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
	
	// Calculate and set throughput
	throughput := float64(size) / duration.Seconds()
	b.throughput.Set(ctx, throughput)
}

// BenchmarkRunner runs performance benchmarks
type BenchmarkRunner struct {
	monitor *Monitor
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner(monitor *Monitor) *BenchmarkRunner {
	return &BenchmarkRunner{
		monitor: monitor,
	}
}

// RunBenchmark runs a benchmark function and records metrics
func (b *BenchmarkRunner) RunBenchmark(ctx context.Context, name string, benchFn func() error) error {
	start := time.Now()
	err := benchFn()
	duration := time.Since(start)
	
	b.monitor.RecordRequest(ctx, duration, name, "BENCH", err == nil)
	
	if err != nil {
		return fmt.Errorf("benchmark %s failed after %v: %w", name, duration, err)
	}
	
	fmt.Printf("Benchmark %s completed in %v\n", name, duration)
	return nil
}

// RunBenchmarkWithIterations runs a benchmark for multiple iterations
func (b *BenchmarkRunner) RunBenchmarkWithIterations(ctx context.Context, name string, iterations int, benchFn func() error) error {
	start := time.Now()
	
	for i := 0; i < iterations; i++ {
		if err := benchFn(); err != nil {
			return fmt.Errorf("benchmark %s failed at iteration %d: %w", name, i, err)
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / time.Duration(iterations)
	
	b.monitor.RecordRequest(ctx, avgDuration, name, "BENCH_ITER", true)
	
	fmt.Printf("Benchmark %s (%d iterations) completed in %v (avg: %v)\n", 
		name, iterations, duration, avgDuration)
	return nil
}