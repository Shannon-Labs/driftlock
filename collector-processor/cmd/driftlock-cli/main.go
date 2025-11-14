package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/driftlockcbad"
)

type anomalyOutput struct {
	Index    int                           `json:"index"`
	Metrics  driftlockcbad.EnhancedMetrics `json:"metrics"`
	Event    json.RawMessage               `json:"event"`
	Why      string                        `json:"why"`
	Detected bool                          `json:"detected"`
}

func main() {
	input := flag.String("input", "", "Path to input file (JSON array or NDJSON)")
	format := flag.String("format", "ndjson", "Input format: json|ndjson")
	output := flag.String("output", "-", "Output path (file or '-' for stdout)")
	baselineSize := flag.Int("baseline", 400, "Baseline size (events)")
	windowSize := flag.Int("window", 1, "Window size (events)")
	hopSize := flag.Int("hop", 1, "Hop size (events)")
	algo := flag.String("algo", "zstd", "Compression algorithm: zstd|lz4|gzip|openzl")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: --input is required")
		os.Exit(2)
	}

	in, err := os.Open(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open input: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// Initialize streaming detector
	usedAlgo := strings.ToLower(*algo)
	detector, err := driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
		BaselineSize:         *baselineSize,
		WindowSize:           *windowSize,
		HopSize:              *hopSize,
		MaxCapacity:          *baselineSize + 4**windowSize + 1024,
		PValueThreshold:      0.05,
		NCDThreshold:         0.3,
		PermutationCount:     1000,
		Seed:                 42,
		CompressionAlgorithm: usedAlgo,
	})
	if err != nil {
		if usedAlgo == "openzl" {
			usedAlgo = "zstd"
			detector, err = driftlockcbad.NewDetector(driftlockcbad.DetectorConfig{
				BaselineSize:         *baselineSize,
				WindowSize:           *windowSize,
				HopSize:              *hopSize,
				MaxCapacity:          *baselineSize + 4**windowSize + 1024,
				PValueThreshold:      0.05,
				NCDThreshold:         0.3,
				PermutationCount:     1000,
				Seed:                 42,
				CompressionAlgorithm: usedAlgo,
			})
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create detector: %v\n", err)
			os.Exit(1)
		}
	}
	defer detector.Close()

	start := time.Now()
	var anomalies []anomalyOutput
	var idx int

	switch strings.ToLower(*format) {
	case "ndjson":
		rd := bufio.NewReader(in)
		for {
			line, err := rd.ReadBytes('\n')
			if len(line) > 0 {
				_ = processEvent(detector, line, idx, &anomalies)
				idx++
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "read error at index %d: %v\n", idx, err)
				os.Exit(1)
			}
		}
	case "json":
		var arr []json.RawMessage
		dec := json.NewDecoder(in)
		if err := dec.Decode(&arr); err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode JSON array: %v\n", err)
			os.Exit(1)
		}
		for i, ev := range arr {
			_ = processEvent(detector, ev, i, &anomalies)
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported --format %q (use json|ndjson)\n", *format)
		os.Exit(2)
	}

	duration := time.Since(start)
	out := map[string]interface{}{
		"total_events":     idx,
		"anomaly_count":    len(anomalies),
		"processing_time":  duration.String(),
		"compression_algo": usedAlgo,
		"anomalies":        anomalies,
	}

	var outWriter io.Writer = os.Stdout
	if *output != "-" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		outWriter = f
	}

	enc := json.NewEncoder(outWriter)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
		os.Exit(1)
	}
}

func processEvent(detector *driftlockcbad.Detector, ev []byte, index int, sink *[]anomalyOutput) error {
	added, err := detector.AddData(ev)
	if err != nil {
		return err
	}
	if !added {
		// Dropped (e.g., privacy rules); skip silently
		return nil
	}

	ready, err := detector.IsReady()
	if err != nil || !ready {
		return err
	}

	detected, metrics, err := detector.DetectAnomaly()
	if err != nil {
		return err
	}
	if detected {
		why := metrics.GetDetailedExplanation()
		*sink = append(*sink, anomalyOutput{
			Index:    index,
			Metrics:  *metrics,
			Event:    json.RawMessage(append([]byte{}, ev...)),
			Why:      why,
			Detected: true,
		})
	}
	return nil
}


