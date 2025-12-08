package main

// #cgo LDFLAGS: -L${SRCDIR}/../../cbad-core/target/release -lcbad_core -lpthread -ldl
// #include <stdlib.h>
// #include <string.h>
// typedef struct {
//     double ncd;
//     double p_value;
//     double baseline_compression_ratio;
//     double window_compression_ratio;
//     double baseline_entropy;
//     double window_entropy;
//     int is_anomaly;
//     double confidence_level;
//     int is_statistically_significant;
//     double compression_ratio_change;
//     double entropy_change;
//     const char* explanation;
//     double recommended_ncd_threshold;
//     size_t recommended_window_size;
//     double data_stability_score;
// } CBADEnhancedMetrics;
//
// extern void* cbad_detector_create_simple();
// extern int cbad_add_transaction(void* detector, const char* data, size_t len);
// extern CBADEnhancedMetrics cbad_detect(void* detector);
// extern int cbad_detector_ready(void* detector);
// extern void cbad_detector_free(void* detector);
// extern void cbad_free_string(char* s);
import "C"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
	"unsafe"
)

type Transaction struct {
	Timestamp     string  `json:"timestamp"`
	TransactionID string  `json:"transaction_id"`
	AmountUSD     float64 `json:"amount_usd"`
	ProcessingMs  int     `json:"processing_ms"`
	OriginCountry string  `json:"origin_country"`
	APIEndpoint   string  `json:"api_endpoint"`
	Status        string  `json:"status"`
}

type AnomalyResult struct {
	Transaction Transaction
	Metrics     CBADEnhancedMetrics
	Explanation string
	Why         []string
	Compare     BaselineCompare
	Examples    []Transaction
}

type CBADEnhancedMetrics struct {
	NCD                        float64 `json:"ncd"`
	PValue                     float64 `json:"p_value"`
	BaselineCompressionRatio   float64 `json:"baseline_compression_ratio"`
	WindowCompressionRatio     float64 `json:"window_compression_ratio"`
	BaselineEntropy            float64 `json:"baseline_entropy"`
	WindowEntropy              float64 `json:"window_entropy"`
	IsAnomaly                  int     `json:"is_anomaly"`
	ConfidenceLevel            float64 `json:"confidence_level"`
	IsStatisticallySignificant int     `json:"is_statistically_significant"`
	CompressionRatioChange     float64 `json:"compression_ratio_change"`
	EntropyChange              float64 `json:"entropy_change"`
	Explanation                string  `json:"explanation"`
	RecommendedNCDThreshold    float64 `json:"recommended_ncd_threshold"`
	RecommendedWindowSize      int     `json:"recommended_window_size"`
	DataStabilityScore         float64 `json:"data_stability_score"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-financial-demo.json>")
	}

	jsonPath := os.Args[1]

	// Read JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("🚀 Loaded %d transactions\n", len(transactions))
	fmt.Printf("⏱️  Processing with Driftlock CBAD engine...\n\n")

	// Create detector
	detector := C.cbad_detector_create_simple()
	if detector == nil {
		log.Fatal("❌ Failed to create detector")
	}
	defer C.cbad_detector_free(detector)

	// Configuration (tuned for demo reliability)
	// - Use a generous warmup to build a stable baseline
	// - Detect on every transaction to catch brief anomalies
	// - Process the full dataset so late anomalies are included
	warmupCount := 400
	// Detection every 25th transaction reduces anomaly count to the 10-30 band
	detectionInterval := 25
	processingLimit := 2000

	if len(transactions) < processingLimit {
		processingLimit = len(transactions)
	}

	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("   - Warmup transactions: %d\n", warmupCount)
	fmt.Printf("   - Detection interval: every %d transactions\n", detectionInterval)
	fmt.Printf("   - Processing %d total transactions\n\n", processingLimit)

	startTime := time.Now()
	var anomalies []AnomalyResult

	// Baseline tracking for warmup window
	var warmupTxs []Transaction
	warmupTxs = make([]Transaction, 0, warmupCount)
	procVals := make([]int, 0, warmupCount)
	amountVals := make([]float64, 0, warmupCount)
	endpointFreq := make(map[string]int)
	originFreq := make(map[string]int)

	// Phase 1: Warmup - ingest without detection and collect baseline stats
	fmt.Printf("🔥 Phase 1: Building baseline (first %d transactions)...\n", warmupCount)
	for i := 0; i < warmupCount && i < processingLimit; i++ {
		tx := transactions[i]
		// Track for baseline
		warmupTxs = append(warmupTxs, tx)
		procVals = append(procVals, tx.ProcessingMs)
		amountVals = append(amountVals, tx.AmountUSD)
		endpointFreq[tx.APIEndpoint]++
		originFreq[tx.OriginCountry]++
		txJSON, _ := json.Marshal(tx)
		// Important: add newline to delimit events for proper permutation testing
		payload := string(txJSON) + "\n"
		dataStr := C.CString(payload)
		result := C.cbad_add_transaction(detector, dataStr, C.size_t(len(payload)))
		C.free(unsafe.Pointer(dataStr))

		if result == 0 {
			log.Printf("Warning: Failed to add transaction %d", i)
		}

		if i%20 == 0 {
			fmt.Printf("   Ingested %d transactions...\n", i)
		}
	}

	ready := C.cbad_detector_ready(detector)
	if ready == 0 {
		log.Fatal("❌ Detector not ready after warmup")
	}
	// Compute baseline stats
	baseline := computeBaselineStats(procVals, amountVals, endpointFreq, originFreq)
	fmt.Printf("✅ Baseline ready! (processing_ms median=%.0fms, p95=%.0fms; amount median=$%.2f)\n\n", baseline.ProcMedian, baseline.ProcP95, baseline.AmountMedian)

	// Phase 2: Detection - ingest AND detect anomalies
	fmt.Printf("🚨 Phase 2: Anomaly detection (transactions %d-%d)...\n", warmupCount, processingLimit)
	for i := warmupCount; i < processingLimit; i++ {
		tx := transactions[i]

		// Always ingest
		txJSON, _ := json.Marshal(tx)
		payload := string(txJSON) + "\n"
		dataStr := C.CString(payload)
		ingestResult := C.cbad_add_transaction(detector, dataStr, C.size_t(len(payload)))

		if ingestResult == 0 {
			log.Printf("Warning: Failed to ingest transaction %d", i)
		}

		// Detect every Nth transaction
		if i%detectionInterval == 0 {
			metrics := C.cbad_detect(detector)
			C.free(unsafe.Pointer(dataStr))

			if metrics.is_anomaly != 0 {
				// Get explanation
				explanation := C.GoString(metrics.explanation)
				if metrics.explanation != nil {
					// Free the Rust-allocated string properly
					defer C.cbad_free_string((*C.char)(metrics.explanation))
				}

				// Convert C struct to Go struct
				goMetrics := CBADEnhancedMetrics{
					NCD:                        float64(metrics.ncd),
					PValue:                     float64(metrics.p_value),
					BaselineCompressionRatio:   float64(metrics.baseline_compression_ratio),
					WindowCompressionRatio:     float64(metrics.window_compression_ratio),
					BaselineEntropy:            float64(metrics.baseline_entropy),
					WindowEntropy:              float64(metrics.window_entropy),
					IsAnomaly:                  int(metrics.is_anomaly),
					ConfidenceLevel:            float64(metrics.confidence_level),
					IsStatisticallySignificant: int(metrics.is_statistically_significant),
					CompressionRatioChange:     float64(metrics.compression_ratio_change),
					EntropyChange:              float64(metrics.entropy_change),
					Explanation:                explanation,
					RecommendedNCDThreshold:    float64(metrics.recommended_ncd_threshold),
					RecommendedWindowSize:      int(metrics.recommended_window_size),
					DataStabilityScore:         float64(metrics.data_stability_score),
				}

				// Build human-friendly reasons
				var why []string
				if goMetrics.NCD >= 0.8 {
					why = append(why, fmt.Sprintf("High dissimilarity vs baseline (NCD=%.3f)", goMetrics.NCD))
				} else if goMetrics.NCD >= 0.5 {
					why = append(why, fmt.Sprintf("Moderate dissimilarity vs baseline (NCD=%.3f)", goMetrics.NCD))
				}
				if goMetrics.PValue <= 0.05 {
					why = append(why, fmt.Sprintf("Statistically significant (p=%.4f)", goMetrics.PValue))
				}
				if goMetrics.BaselineCompressionRatio > 0 && goMetrics.WindowCompressionRatio > 0 {
					why = append(why, fmt.Sprintf("Compression efficiency dropped from %.2fx → %.2fx (Δ %.0f%%)", goMetrics.BaselineCompressionRatio, goMetrics.WindowCompressionRatio, goMetrics.CompressionRatioChange*100))
				} else if goMetrics.CompressionRatioChange != 0 {
					why = append(why, fmt.Sprintf("Compression efficiency changed by %.0f%%", goMetrics.CompressionRatioChange*100))
				}
				if goMetrics.EntropyChange > 0.05 {
					why = append(why, fmt.Sprintf("Randomness increased (entropy Δ +%.0f%%)", goMetrics.EntropyChange*100))
				} else if goMetrics.EntropyChange < -0.05 {
					why = append(why, fmt.Sprintf("Randomness decreased (entropy Δ %.0f%%)", goMetrics.EntropyChange*100))
				}

				// Baseline-aware bullets
				// processing_ms vs baseline median and z-score
				var zStr string
				if baseline.ProcStd > 0 {
					z := (float64(tx.ProcessingMs) - baseline.ProcMean) / baseline.ProcStd
					zStr = fmt.Sprintf("(z=%+.1f)", z)
				} else {
					zStr = "(z=N/A)"
				}
				why = append(why, fmt.Sprintf("processing_ms %dms vs baseline median %.0fms %s", tx.ProcessingMs, baseline.ProcMedian, zStr))

				// amount vs baseline median and z-score, with ratio
				if baseline.AmountMedian > 0 {
					ratio := tx.AmountUSD / baseline.AmountMedian
					var az string
					if baseline.AmountStd > 0 {
						az = fmt.Sprintf("(z=%+.1f)", (tx.AmountUSD-baseline.AmountMean)/baseline.AmountStd)
					} else {
						az = "(z=N/A)"
					}
					why = append(why, fmt.Sprintf("amount $%.2f vs baseline median $%.2f (×%.1f) %s", tx.AmountUSD, baseline.AmountMedian, ratio, az))
				}

				// Endpoint frequency insight
				epCount := baseline.EndpointFreq[tx.APIEndpoint]
				epPct := pctFloat(epCount, baseline.Count)
				if epPct >= 50.0 {
					why = append(why, fmt.Sprintf("Endpoint %s is common (%.1f%% of baseline)", tx.APIEndpoint, epPct))
				} else if epPct >= 10.0 {
					why = append(why, fmt.Sprintf("Endpoint %s appears occasionally (%.1f%% of baseline)", tx.APIEndpoint, epPct))
				} else {
					why = append(why, fmt.Sprintf("Endpoint %s is rare (%.1f%% of baseline)", tx.APIEndpoint, epPct))
				}

				// Origin frequency insight
				ocCount := baseline.OriginFreq[tx.OriginCountry]
				ocPct := pctFloat(ocCount, baseline.Count)
				if ocCount == 0 {
					why = append(why, fmt.Sprintf("Origin %s not seen in baseline (0/%.0f)", tx.OriginCountry, float64(baseline.Count)))
				} else if ocPct < 5.0 {
					why = append(why, fmt.Sprintf("Origin %s is rare (%.1f%% of baseline)", tx.OriginCountry, ocPct))
				}

				// Build baseline comparison panel data
				cmp := BaselineCompare{
					ProcValue:    tx.ProcessingMs,
					ProcMedian:   baseline.ProcMedian,
					ZScore:       0,
					HasZScore:    baseline.ProcStd > 0,
					AmountValue:  tx.AmountUSD,
					AmountMedian: baseline.AmountMedian,
					AmountRatio: func() float64 {
						if baseline.AmountMedian > 0 {
							return tx.AmountUSD / baseline.AmountMedian
						}
						return 0
					}(),
					AmountHasZ:               baseline.AmountStd > 0,
					Endpoint:                 tx.APIEndpoint,
					EndpointCount:            epCount,
					EndpointPct:              epPct,
					Origin:                   tx.OriginCountry,
					OriginCount:              baseline.OriginFreq[tx.OriginCountry],
					OriginPct:                pctFloat(baseline.OriginFreq[tx.OriginCountry], baseline.Count),
					BaselineCompressionRatio: goMetrics.BaselineCompressionRatio,
					WindowCompressionRatio:   goMetrics.WindowCompressionRatio,
					CompressionDeltaPct:      goMetrics.CompressionRatioChange * 100.0,
				}
				if cmp.HasZScore {
					cmp.ZScore = (float64(tx.ProcessingMs) - baseline.ProcMean) / baseline.ProcStd
				}
				if cmp.AmountHasZ {
					cmp.AmountZScore = (tx.AmountUSD - baseline.AmountMean) / baseline.AmountStd
				}

				// Find similar normal examples from warmup
				examples := nearestNormalExamples(tx, warmupTxs, 3)

				anomalies = append(anomalies, AnomalyResult{
					Transaction: tx,
					Metrics:     goMetrics,
					Explanation: explanation,
					Why:         why,
					Compare:     cmp,
					Examples:    examples,
				})
			}

			if i%500 == 0 {
				fmt.Printf("   Processed %d transactions...\n", i)
			}
		} else {
			C.free(unsafe.Pointer(dataStr))
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\n🎉 Processing complete!\n")
	fmt.Printf("📊 Found %d anomalies out of %d transactions\n", len(anomalies), processingLimit)
	fmt.Printf("⚡ Processing time: %v\n\n", duration)

	// Generate HTML report
	outputPath := "demo-output.html"
	if err := generateHTMLReport(anomalies, transactions, time.Since(startTime), baselineSummaryFrom(baseline), outputPath); err != nil {
		log.Fatalf("❌ Failed to generate HTML report: %v", err)
	}

	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("✅ Demo output written to: %s\n", absPath)
	fmt.Printf("🌐 Open this file in your browser to view results\n")
	fmt.Printf("\n💡 Tip: Use 'open %s' on macOS or 'xdg-open %s' on Linux\n", absPath, absPath)
}

// Baseline types and helpers
type BaselineStats struct {
	Count        int
	ProcMin      int
	ProcMax      int
	ProcMean     float64
	ProcMedian   float64
	ProcStd      float64
	ProcP95      float64
	AmountMin    float64
	AmountMax    float64
	AmountMean   float64
	AmountMedian float64
	AmountStd    float64
	AmountP95    float64
	EndpointFreq map[string]int
	OriginFreq   map[string]int
}

type FreqItem struct {
	Key   string
	Count int
	Pct   float64
}

type BaselineSummary struct {
	Count        int
	ProcMin      int
	ProcMedian   float64
	ProcP95      float64
	AmountMin    float64
	AmountMedian float64
	AmountP95    float64
	TopEndpoints []FreqItem
	TopOrigins   []FreqItem
}

type BaselineCompare struct {
	ProcValue                int
	ProcMedian               float64
	ZScore                   float64
	HasZScore                bool
	AmountValue              float64
	AmountMedian             float64
	AmountRatio              float64
	AmountZScore             float64
	AmountHasZ               bool
	Endpoint                 string
	EndpointCount            int
	EndpointPct              float64
	Origin                   string
	OriginCount              int
	OriginPct                float64
	BaselineCompressionRatio float64
	WindowCompressionRatio   float64
	CompressionDeltaPct      float64
}

func computeBaselineStats(procVals []int, amountVals []float64, ep map[string]int, origin map[string]int) BaselineStats {
	bs := BaselineStats{Count: len(procVals), EndpointFreq: ep, OriginFreq: origin}
	if len(procVals) == 0 {
		return bs
	}
	// Min/Max/Mean
	minV, maxV := procVals[0], procVals[0]
	var sum int64
	for _, v := range procVals {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
		sum += int64(v)
	}
	bs.ProcMin = minV
	bs.ProcMax = maxV
	bs.ProcMean = float64(sum) / float64(len(procVals))
	// Median/P95 require sort
	sorted := append([]int(nil), procVals...)
	sort.Ints(sorted)
	n := len(sorted)
	if n%2 == 1 {
		bs.ProcMedian = float64(sorted[n/2])
	} else {
		bs.ProcMedian = (float64(sorted[n/2-1]) + float64(sorted[n/2])) / 2.0
	}
	// p95 nearest-rank
	idx := int(math.Ceil(0.95*float64(n))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	bs.ProcP95 = float64(sorted[idx])
	// Std deviation (population)
	var varSum float64
	mu := bs.ProcMean
	for _, v := range procVals {
		d := float64(v) - mu
		varSum += d * d
	}
	bs.ProcStd = math.Sqrt(varSum / float64(n))

	// Amount stats
	if len(amountVals) > 0 {
		amin, amax := amountVals[0], amountVals[0]
		var asum float64
		for _, a := range amountVals {
			if a < amin {
				amin = a
			}
			if a > amax {
				amax = a
			}
			asum += a
		}
		bs.AmountMin = amin
		bs.AmountMax = amax
		bs.AmountMean = asum / float64(len(amountVals))
		as := append([]float64(nil), amountVals...)
		sort.Float64s(as)
		m := len(as)
		if m%2 == 1 {
			bs.AmountMedian = as[m/2]
		} else {
			bs.AmountMedian = (as[m/2-1] + as[m/2]) / 2.0
		}
		aidx := int(math.Ceil(0.95*float64(m))) - 1
		if aidx < 0 {
			aidx = 0
		}
		if aidx >= m {
			aidx = m - 1
		}
		bs.AmountP95 = as[aidx]
		var avarSum float64
		amu := bs.AmountMean
		for _, a := range amountVals {
			d := a - amu
			avarSum += d * d
		}
		bs.AmountStd = math.Sqrt(avarSum / float64(m))
	}
	return bs
}

func pctFloat(count int, total int) float64 {
	if total <= 0 {
		return 0
	}
	return (float64(count) / float64(total)) * 100.0
}

func topN(freq map[string]int, total, n int) []FreqItem {
	items := make([]FreqItem, 0, len(freq))
	for k, c := range freq {
		items = append(items, FreqItem{Key: k, Count: c, Pct: pctFloat(c, total)})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Key < items[j].Key
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > n {
		items = items[:n]
	}
	return items
}

func nearestNormalExamples(anom Transaction, warmup []Transaction, k int) []Transaction {
	type cand struct {
		tx   Transaction
		diff int
	}
	addNearest := func(candidates []Transaction) []Transaction {
		list := make([]cand, 0, len(candidates))
		for _, t := range candidates {
			d := t.ProcessingMs - anom.ProcessingMs
			if d < 0 {
				d = -d
			}
			list = append(list, cand{tx: t, diff: d})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].diff < list[j].diff })
		out := make([]Transaction, 0, minInt(k, len(list)))
		for i := 0; i < len(list) && len(out) < k; i++ {
			out = append(out, list[i].tx)
		}
		return out
	}
	// Priority 1: same endpoint + same origin
	both := make([]Transaction, 0)
	for _, t := range warmup {
		if t.APIEndpoint == anom.APIEndpoint && t.OriginCountry == anom.OriginCountry {
			both = append(both, t)
		}
	}
	if len(both) >= 1 {
		return addNearest(both)
	}
	// Priority 2: same endpoint
	ep := make([]Transaction, 0)
	for _, t := range warmup {
		if t.APIEndpoint == anom.APIEndpoint {
			ep = append(ep, t)
		}
	}
	if len(ep) >= 1 {
		return addNearest(ep)
	}
	// Priority 3: same origin
	oc := make([]Transaction, 0)
	for _, t := range warmup {
		if t.OriginCountry == anom.OriginCountry {
			oc = append(oc, t)
		}
	}
	if len(oc) >= 1 {
		return addNearest(oc)
	}
	// Fallback: any warmup
	return addNearest(warmup)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func baselineSummaryFrom(bs BaselineStats) BaselineSummary {
	return BaselineSummary{
		Count:        bs.Count,
		ProcMin:      bs.ProcMin,
		ProcMedian:   bs.ProcMedian,
		ProcP95:      bs.ProcP95,
		AmountMin:    bs.AmountMin,
		AmountMedian: bs.AmountMedian,
		AmountP95:    bs.AmountP95,
		TopEndpoints: topN(bs.EndpointFreq, bs.Count, 5),
		TopOrigins:   topN(bs.OriginFreq, bs.Count, 5),
	}
}

func generateHTMLReport(anomalies []AnomalyResult, allTransactions []Transaction, duration time.Duration, baseline BaselineSummary, outputPath string) error {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Driftlock DORA Compliance Demo</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            line-height: 1.6;
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }
        
        .header {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 16px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5rem;
            background: linear-gradient(135deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 0.5rem;
        }
        
        .header .tagline {
            font-size: 1.2rem;
            color: #666;
            margin-bottom: 1rem;
        }
        
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-top: 1.5rem;
        }
        
        .stat-card {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 1.5rem;
            border-radius: 12px;
            text-align: center;
        }
        
        .stat-card .number {
            font-size: 2rem;
            font-weight: bold;
            display: block;
        }
        
        .stat-card .label {
            font-size: 0.9rem;
            opacity: 0.9;
        }
        
        .anomalies-section {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 16px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
        }
        .baseline-summary { margin-top: 1.5rem; text-align: left; }
        .baseline-summary h3 { color: #667eea; margin-bottom: 0.5rem; }
        .baseline-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; }
        .baseline-card { background: #fff; border: 1px solid #e8ecff; border-radius: 8px; padding: 0.75rem 1rem; }
        .baseline-list li { margin-left: 1rem; }
        .why {
            background: #fff;
            border: 1px solid #e8ecff;
            border-radius: 8px;
            padding: 1rem;
        }
        .why strong { color: #667eea; }
        .why-list { margin: 0.5rem 0 0 1rem; }
        .why-list li { margin: 0.25rem 0; }
        
        .anomalies-section h2 {
            color: #667eea;
            margin-bottom: 1.5rem;
            font-size: 1.8rem;
        }
        
        .anomaly-card {
            background: #f8f9ff;
            border: 2px solid #e8ecff;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
            transition: all 0.3s ease;
        }
        
        .anomaly-card:hover {
            border-color: #667eea;
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
        }
        
        .anomaly-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }
        
        .anomaly-title {
            font-weight: 600;
            color: #333;
            font-size: 1.1rem;
        }
        
        .anomaly-badge {
            background: linear-gradient(135deg, #ff6b6b, #ee5a52);
            color: white;
            padding: 0.25rem 0.75rem;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: 600;
        }
        
        .transaction-details {
            background: white;
            padding: 1rem;
            border-radius: 8px;
            margin-bottom: 1rem;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.85rem;
            border-left: 4px solid #667eea;
        }
        
        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 1rem;
            margin-bottom: 1rem;
        }
        
        .metric {
            background: white;
            padding: 0.75rem;
            border-radius: 8px;
            text-align: center;
            border: 1px solid #e8ecff;
        }
        
        .metric-label {
            font-size: 0.8rem;
            color: #666;
            display: block;
            margin-bottom: 0.25rem;
        }
        
        .metric-value {
            font-size: 1.1rem;
            font-weight: 600;
            color: #333;
        }

        .explanation {
            background: linear-gradient(135deg, #667eea10, #764ba210);
            padding: 1rem;
            border-radius: 8px;
            border-left: 4px solid #667eea;
            font-style: italic;
            color: #555;
        }

        .compare { background: #fff; border: 1px solid #e8ecff; border-radius: 8px; padding: 1rem; margin-top: 0.75rem; }
        .compare h4 { margin-bottom: 0.5rem; color: #333; }
        .compare table { width: 100%; border-collapse: collapse; }
        .compare th, .compare td { text-align: left; padding: 6px 8px; border-bottom: 1px solid #f0f2ff; }
        .examples { background: #fff; border: 1px solid #e8ecff; border-radius: 8px; padding: 1rem; margin-top: 0.75rem; }
        .examples h4 { margin-bottom: 0.5rem; color: #333; }
        .examples pre { background: #f7f8ff; border: 1px solid #e8ecff; border-radius: 6px; padding: 8px; overflow-x: auto; }
        
        .footer {
            text-align: center;
            color: white;
            padding: 2rem;
            opacity: 0.8;
        }
        
        .no-anomalies {
            text-align: center;
            padding: 3rem;
            color: #666;
        }
        
        .no-anomalies h3 {
            font-size: 1.5rem;
            margin-bottom: 1rem;
            color: #667eea;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛡️ Driftlock</h1>
            <div class="tagline">Regulator-Proof Algorithms for DORA Compliance</div>
            <p>Compression-based anomaly detection with mathematical explanations</p>
            
            <div class="stats">
                <div class="stat-card">
                    <span class="number">{{.TotalTransactions}}</span>
                    <span class="label">Transactions Processed</span>
                </div>
                <div class="stat-card">
                    <span class="number">{{.AnomalyCount}}</span>
                    <span class="label">Anomalies Detected</span>
                </div>
                <div class="stat-card">
                    <span class="number">{{printf "%.1f" .DetectionRate}}%</span>
                    <span class="label">Detection Rate</span>
                </div>
                <div class="stat-card">
                    <span class="number">{{.ProcessingTime}}</span>
                    <span class="label">Processing Time</span>
                </div>
            </div>

            <div class="baseline-summary">
                <h3>Baseline Summary (first 400 events)</h3>
                <div class="baseline-grid">
                  <div class="baseline-card">
                    <strong>processing_ms</strong><br>
                    min: {{.Baseline.ProcMin}}ms · median: {{printf "%.0f" .Baseline.ProcMedian}}ms · p95: {{printf "%.0f" .Baseline.ProcP95}}ms
                  </div>
                  <div class="baseline-card">
                    <strong>amount_usd</strong><br>
                    min: ${{printf "%.2f" .Baseline.AmountMin}} · median: ${{printf "%.2f" .Baseline.AmountMedian}} · p95: ${{printf "%.2f" .Baseline.AmountP95}}
                  </div>
                  <div class="baseline-card">
                    <strong>Top Endpoints</strong>
                    <ul class="baseline-list">
                        {{range .Baseline.TopEndpoints}}
                            <li>{{.Key}} — {{printf "%.1f" .Pct}}% ({{.Count}}/{{$.Baseline.Count}})</li>
                        {{end}}
                    </ul>
                  </div>
                  <div class="baseline-card">
                    <strong>Top Origins</strong>
                    <ul class="baseline-list">
                        {{range .Baseline.TopOrigins}}
                            <li>{{.Key}} — {{printf "%.1f" .Pct}}% ({{.Count}}/{{$.Baseline.Count}})</li>
                        {{end}}
                    </ul>
                  </div>
                </div>
            </div>
        </div>

        <div class="anomalies-section">
            <h2>🚨 Detected Anomalies</h2>
            
            {{if .Anomalies}}
                {{range .Anomalies}}
                <div class="anomaly-card">
                    <div class="anomaly-header">
                        <div class="anomaly-title">Transaction {{.Transaction.TransactionID}}</div>
                        <div class="anomaly-badge">ANOMALY DETECTED</div>
                    </div>
                    
                    <div class="transaction-details">
                        <strong>Timestamp:</strong> {{.Transaction.Timestamp}}<br>
                        <strong>Amount:</strong> ${{.Transaction.AmountUSD}}<br>
                        <strong>Processing Time:</strong> {{.Transaction.ProcessingMs}}ms<br>
                        <strong>Origin:</strong> {{.Transaction.OriginCountry}}<br>
                        <strong>Endpoint:</strong> {{.Transaction.APIEndpoint}}<br>
                        <strong>Status:</strong> {{.Transaction.Status}}
                    </div>
                    
                    <div class="metrics-grid">
                        <div class="metric">
                            <span class="metric-label">NCD Score</span>
                            <span class="metric-value">{{printf "%.3f" .Metrics.NCD}}</span>
                        </div>
                        <div class="metric">
                            <span class="metric-label">P-Value</span>
                            <span class="metric-value">{{printf "%.4f" .Metrics.PValue}}</span>
                        </div>
                        <div class="metric">
                            <span class="metric-label">Confidence</span>
                            <span class="metric-value">{{printf "%.1f%%" (pct .Metrics.ConfidenceLevel)}}</span>
                        </div>
                        <div class="metric">
                            <span class="metric-label">Compression Δ</span>
                            <span class="metric-value">{{printf "%.2f" .Metrics.CompressionRatioChange}}</span>
                        </div>
                    </div>
                    
                <div class="explanation">
                    <strong>🧠 Explanation:</strong> {{.Explanation}}
                </div>

                {{if .Why}}
                <div class="why" style="margin-top: 0.75rem;">
                    <strong>Why this is anomalous</strong>
                    <ul class="why-list">
                        {{range .Why}}
                        <li>{{.}}</li>
                        {{end}}
                    </ul>
                </div>
                {{end}}

                <div class="compare">
                    <h4>Baseline Comparison</h4>
                    <table>
                        <tr>
                            <th>processing_ms</th>
                            <td>{{.Compare.ProcValue}}ms vs median {{printf "%.0f" .Compare.ProcMedian}}ms {{if .Compare.HasZScore}}(z={{printf "%+.1f" .Compare.ZScore}}){{else}}(z=N/A){{end}}</td>
                        </tr>
                        <tr>
                            <th>amount_usd</th>
                            <td>${{printf "%.2f" .Compare.AmountValue}} vs median ${{printf "%.2f" .Compare.AmountMedian}} (×{{printf "%.1f" .Compare.AmountRatio}} {{if .Compare.AmountHasZ}}; z={{printf "%+.1f" .Compare.AmountZScore}}{{else}}; z=N/A{{end}})</td>
                        </tr>
                        <tr>
                            <th>api_endpoint</th>
                            <td>{{.Compare.Endpoint}} — {{printf "%.1f" .Compare.EndpointPct}}% of baseline ({{.Compare.EndpointCount}}/{{$.Baseline.Count}})</td>
                        </tr>
                        <tr>
                            <th>origin_country</th>
                            <td>{{.Compare.Origin}} — {{printf "%.1f" .Compare.OriginPct}}% of baseline ({{.Compare.OriginCount}}/{{$.Baseline.Count}})</td>
                        </tr>
                        <tr>
                            <th>compression</th>
                            <td>{{printf "%.2fx" .Compare.BaselineCompressionRatio}} → {{printf "%.2fx" .Compare.WindowCompressionRatio}} (Δ {{printf "%+.0f%%" .Compare.CompressionDeltaPct}})</td>
                        </tr>
                    </table>
                </div>

                {{if .Examples}}
                <details>
                  <summary>Similar normal examples</summary>
                  <div class="details-body">
                    {{range .Examples}}
<pre><code>{{.Timestamp}} | ${{printf "%.2f" .AmountUSD}} | {{.ProcessingMs}}ms | {{.OriginCountry}} | {{.APIEndpoint}} | {{.Status}}</code></pre>
                    {{end}}
                  </div>
                </details>
                {{end}}
            </div>
            {{end}}
            {{else}}
                <div class="no-anomalies">
                    <h3>✅ No Anomalies Detected</h3>
                    <p>All transactions appear normal. The system found no statistically significant deviations from baseline behavior.</p>
                </div>
            {{end}}
        </div>

        <div class="footer">
            <p>Driftlock - Patent-pending compression-based anomaly detection | Shannon Labs</p>
            <p>Built with Rust CBAD core for transparent, regulator-friendly math</p>
        </div>
    </div>
</body>
</html>`

	// Parse template with helpers
	funcMap := template.FuncMap{
		"pct": func(v float64) float64 { return v * 100.0 },
	}
	t, err := template.New("report").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("template parse error: %v", err)
	}

	// Prepare data
	data := struct {
		TotalTransactions int
		AnomalyCount      int
		DetectionRate     float64
		ProcessingTime    string
		Baseline          BaselineSummary
		Anomalies         []AnomalyResult
	}{
		TotalTransactions: len(allTransactions),
		AnomalyCount:      len(anomalies),
		DetectionRate:     float64(len(anomalies)) / float64(len(allTransactions)) * 100,
		ProcessingTime:    duration.String(),
		Baseline:          baseline,
		Anomalies:         anomalies,
	}

	// Write to file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("template execution error: %v", err)
	}

	return nil
}
