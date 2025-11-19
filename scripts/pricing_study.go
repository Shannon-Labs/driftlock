package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Cloud Run Pricing (Tier 1, us-central1, approx)
const (
	CostPerVCPUSecond = 0.00002400
	CostPerGBSecond   = 0.00000250
	CostPerRequest    = 0.40 / 1000000 // $0.40 per million
	AssumedVCPU       = 1.0
	AssumedMemoryGB   = 0.5 // 512MB
)

var (
	apiURL      string
	apiKey      string
	filePath    string
	batchSize   int
	concurrency int
)

func main() {
	flag.StringVar(&apiURL, "url", "http://localhost:8080", "API URL")
	flag.StringVar(&apiKey, "key", "", "API Key")
	flag.StringVar(&filePath, "file", "", "Path to JSONL file")
	flag.IntVar(&batchSize, "batch", 50, "Batch size")
	flag.IntVar(&concurrency, "concurrency", 4, "Concurrency level")
	flag.Parse()

	if filePath == "" || apiKey == "" {
		log.Fatal("File path and API key are required")
	}

	log.Printf("Starting pricing study on %s", filePath)
	log.Printf("Batch size: %d, Concurrency: %d", batchSize, concurrency)

	events := make(chan []json.RawMessage, 100)
	var totalEvents int64
	var totalRequests int64
	var totalDuration int64 // Microseconds
	var errors int64

	start := time.Now()

	// Worker pool
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range events {
				dur, err := sendBatch(batch)
				if err != nil {
					atomic.AddInt64(&errors, 1)
					log.Printf("Error: %v", err)
				} else {
					atomic.AddInt64(&totalDuration, dur.Microseconds())
					atomic.AddInt64(&totalRequests, 1)
					atomic.AddInt64(&totalEvents, int64(len(batch)))
				}
			}
		}()
	}

	// File reader
	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 10*1024*1024)

		var currentBatch []json.RawMessage
		for scanner.Scan() {
			line := scanner.Bytes()
			
			// Unmarshal to modify
			var event map[string]interface{}
			if err := json.Unmarshal(line, &event); err != nil {
				log.Printf("Skipping invalid JSON: %v", err)
				continue
			}

			// Update timestamp to now to avoid stale data handling
			event["timestamp"] = time.Now().Format(time.RFC3339)
			
			// Add nonce to ensure uniqueness and avoid dedup
			nonce := make([]byte, 8)
			rand.Read(nonce)
			event["_nonce"] = hex.EncodeToString(nonce)

			modifiedLine, _ := json.Marshal(event)
			currentBatch = append(currentBatch, json.RawMessage(modifiedLine))

			if len(currentBatch) >= batchSize {
				events <- currentBatch
				currentBatch = nil
			}
		}
		if len(currentBatch) > 0 {
			events <- currentBatch
		}
		close(events)
	}()

	wg.Wait()
	totalTime := time.Since(start)

	// Analysis
	var avgLatency time.Duration
	if totalRequests > 0 {
		avgLatency = time.Duration(totalDuration/totalRequests) * time.Microsecond
	}
	
	reqPerSec := float64(totalRequests) / totalTime.Seconds()
	eventsPerSec := float64(totalEvents) / totalTime.Seconds()

	// Cost Estimation
	billableSeconds := float64(totalDuration) / 1000000.0
	computeCost := billableSeconds * AssumedVCPU * CostPerVCPUSecond
	memoryCost := billableSeconds * AssumedMemoryGB * CostPerGBSecond
	requestCost := float64(totalRequests) * CostPerRequest
	totalCost := computeCost + memoryCost + requestCost

	costPer1k := 0.0
	if totalEvents > 0 {
		costPer1k = (totalCost / float64(totalEvents)) * 1000
	}

	fmt.Printf("\n=== Pricing Study Results ===\n")
	fmt.Printf("Input File:       %s\n", filePath)
	fmt.Printf("Total Events:     %d\n", totalEvents)
	fmt.Printf("Total Requests:   %d\n", totalRequests)
	fmt.Printf("Total Errors:     %d\n", errors)
	fmt.Printf("Wall Time:        %v\n", totalTime)
	fmt.Printf("Throughput:       %.2f events/sec (%.2f req/sec)\n", eventsPerSec, reqPerSec)
	fmt.Printf("Avg Latency:      %v\n", avgLatency)
	fmt.Printf("Billable Time:    %.4f seconds\n", billableSeconds)
	fmt.Printf("\n--- Estimated Cost (Cloud Run) ---\n")
	fmt.Printf("Compute Cost:     $%.6f\n", computeCost)
	fmt.Printf("Memory Cost:      $%.6f\n", memoryCost)
	fmt.Printf("Request Cost:     $%.6f\n", requestCost)
	fmt.Printf("TOTAL COST:       $%.6f\n", totalCost)
	fmt.Printf("Cost Per 1k Evts: $%.6f\n", costPer1k)
}

func sendBatch(events []json.RawMessage) (time.Duration, error) {
	payload := map[string]interface{}{
		"stream_id": "default",
		"events":    events,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL+"/v1/detect", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return duration, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return duration, nil
}
