//go:build main
// +build main

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

var (
	streams   = flag.Int("streams", 10, "Number of concurrent streams")
	events    = flag.Int("events", 1000, "Events per stream")
	verbose   = flag.Bool("v", false, "Verbose output")
)

type TestResult struct {
	StreamID       int
	TotalEvents    int
	DetectionTime  time.Duration
	AnomaliesFound int
	Errors         []string
}

func main() {
	flag.Parse()

	fmt.Printf("Starting CBAD load test\n")
	fmt.Printf("  Streams: %d\n", *streams)
	fmt.Printf("  Events per stream: %d\n", *events)
	fmt.Printf("  Total events: %d\n", *streams**events)

	if !driftlockcbad.IsAvailable() {
		log.Fatalf("CBAD library not available: %v", driftlockcbad.AvailabilityError())
	}

	// Track metrics
	var totalEvents int64
	var totalDetections int64
	var totalAnomalies int64
	var totalErrors int64
	var totalLatency int64

	// Create results channel
	results := make(chan TestResult, *streams)

	// Start load test
	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < *streams; i++ {
		wg.Add(1)
		go func(streamID int) {
			defer wg.Done()
			result := runStream(streamID, *events, *verbose)
			results <- result

			atomic.AddInt64(&totalEvents, int64(result.TotalEvents))
			atomic.AddInt64(&totalDetections, int64(1))
			atomic.AddInt64(&totalAnomalies, int64(result.AnomaliesFound))
			atomic.AddInt64(&totalErrors, int64(len(result.Errors)))
			atomic.AddInt64(&totalLatency, int64(result.DetectionTime.Nanoseconds()))
		}(i)
	}

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []TestResult
	for result := range results {
		allResults = append(allResults, result)
	}

	elapsed := time.Since(start)

	// Print summary
	fmt.Println()
	fmt.Println("=== Load Test Results ===")
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Events per second: %.2f\n", float64(totalEvents)/elapsed.Seconds())
	fmt.Printf("Average detection latency: %v\n", time.Duration(totalLatency/totalDetections))
	fmt.Printf("Anomalies detected: %d (%.2f%%)\n", totalAnomalies,
		float64(totalAnomalies)/float64(totalEvents)*100)
	fmt.Printf("Errors: %d\n", totalErrors)

	// Performance targets
	fmt.Println()
	fmt.Println("=== Performance Targets ===")
	targetLatency := 100 * time.Millisecond
	actualLatency := time.Duration(totalLatency / totalDetections)

	if actualLatency > targetLatency {
		fmt.Printf("❌ Latency target failed: %v > %v\n", actualLatency, targetLatency)
	} else {
		fmt.Printf("✅ Latency target met: %v ≤ %v\n", actualLatency, targetLatency)
	}

	eps := float64(totalEvents) / elapsed.Seconds()
	targetEPS := 1000.0

	if eps < targetEPS {
		fmt.Printf("❌ Throughput target failed: %.2f < %.2f events/sec\n", eps, targetEPS)
	} else {
		fmt.Printf("✅ Throughput target met: %.2f ≥ %.2f events/sec\n", eps, targetEPS)
	}

	// Detailed results if verbose
	if *verbose {
		fmt.Println()
		fmt.Println("=== Stream Details ===")
		for _, result := range allResults {
			fmt.Printf("Stream %d: events=%d, anomalies=%d, latency=%v",
				result.StreamID, result.TotalEvents, result.AnomaliesFound, result.DetectionTime)
			if len(result.Errors) > 0 {
				fmt.Printf(", errors=%d", len(result.Errors))
			}
			fmt.Println()
		}
	}
}

func runStream(streamID, eventCount int, verbose bool) TestResult {
	result := TestResult{
		StreamID: streamID,
		Errors:   []string{},
	}

	// Create detector with production-like config
	config := driftlockcbad.DefaultProductionConfig()
	config.BaselineSize = 500
	config.WindowSize = 100
	config.PermutationCount = 50 // Reduce for load testing

	detector, err := driftlockcbad.NewDetector(config)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to create detector: %v", err))
		return result
	}
	defer detector.Close()

	// Generate events
	normalEvent := []byte(`INFO service=api request_id=123 status=200 duration=12ms`)
	anomalyEvent := []byte(`ERROR service=db connection_failed stacktrace=at com.example.DB.connect(DB.java:42)`)

	events := make([][]byte, eventCount)
	for i := 0; i < eventCount; i++ {
		if i%100 == 99 {
			// Insert anomaly every 100 events
			events[i] = anomalyEvent
		} else {
			events[i] = normalEvent
		}
	}

	start := time.Now()

	// Feed events to detector
	for i, event := range events {
		_, err := detector.AddData(event)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Event %d: %v", i, err))
		}

		// Check if ready for detection
		if ready, _ := detector.IsReady(); ready {
			anomalous, _, err := detector.DetectAnomaly()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Detection %d: %v", i, err))
			}
			if anomalous {
				result.AnomaliesFound++
			}
		}
	}

	result.DetectionTime = time.Since(start)
	result.TotalEvents = eventCount

	if verbose {
		fmt.Printf("Stream %d completed: %d events, %d anomalies, %v\n",
			streamID, result.TotalEvents, result.AnomaliesFound, result.DetectionTime)
	}

	return result
}