package main

// #cgo LDFLAGS: -L${SRCDIR}/../../cbad-core/target/release -lcbad_core -lpthread -ldl
// #include <stdlib.h>
// #include <string.h>
// typedef struct {
//     const char* data;
//     size_t len;
// } CBADData;
//
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
// } CBADEnhancedMetrics;
//
// extern void* cbad_detector_create_simple();
// extern CBADEnhancedMetrics cbad_detector_process(void* detector, const char* data, size_t len);
// extern void cbad_detector_free(void* detector);
// extern void cbad_free_string(char* s);
import "C"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

type Transaction struct {
	Timestamp    string  `json:"timestamp"`
	TransactionID string `json:"transaction_id"`
	AmountUSD    float64 `json:"amount_usd"`
	ProcessingMs int     `json:"processing_ms"`
	OriginCountry string `json:"origin_country"`
	APIEndpoint  string  `json:"api_endpoint"`
	Status       string  `json:"status"`
}

type AnomalyResult struct {
	Transaction Transaction
	Metrics     CBADEnhancedMetrics
	Explanation string
}

type CBADEnhancedMetrics struct {
	NCD                         float64 `json:"ncd"`
	PValue                      float64 `json:"p_value"`
	BaselineCompressionRatio    float64 `json:"baseline_compression_ratio"`
	WindowCompressionRatio      float64 `json:"window_compression_ratio"`
	BaselineEntropy             float64 `json:"baseline_entropy"`
	WindowEntropy               float64 `json:"window_entropy"`
	IsAnomaly                   int     `json:"is_anomaly"`
	ConfidenceLevel             float64 `json:"confidence_level"`
	IsStatisticallySignificant  int     `json:"is_statistically_significant"`
	CompressionRatioChange      float64 `json:"compression_ratio_change"`
	EntropyChange               float64 `json:"entropy_change"`
	Explanation                 string  `json:"explanation"`
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

	// Process transactions and detect anomalies
	var anomalies []AnomalyResult
	startTime := time.Now()
	
	// Process first 1000 transactions for demo (speeds things up)
	processingLimit := 1000
	if len(transactions) < processingLimit {
		processingLimit = len(transactions)
	}
	
	fmt.Printf("Processing %d transactions for demo...\n\n", processingLimit)
	
	for i := 0; i < processingLimit; i++ {
		tx := transactions[i]
		
		// Convert transaction to string for analysis
		txJSON, _ := json.Marshal(tx)
		dataStr := C.CString(string(txJSON))
		
		// Process with CBAD
		metrics := C.cbad_detector_process(detector, dataStr, C.size_t(len(txJSON)))
		C.free(unsafe.Pointer(dataStr))
		
		if metrics.is_anomaly != 0 {
			// Get explanation
			explanation := C.GoString(metrics.explanation)
			if metrics.explanation != nil {
				defer C.cbad_free_string(C.CString(explanation))
			}

			// Convert C struct to Go struct
			goMetrics := CBADEnhancedMetrics{
				NCD:                         float64(metrics.ncd),
				PValue:                      float64(metrics.p_value),
				BaselineCompressionRatio:    float64(metrics.baseline_compression_ratio),
				WindowCompressionRatio:      float64(metrics.window_compression_ratio),
				BaselineEntropy:             float64(metrics.baseline_entropy),
				WindowEntropy:               float64(metrics.window_entropy),
				IsAnomaly:                   int(metrics.is_anomaly),
				ConfidenceLevel:             float64(metrics.confidence_level),
				IsStatisticallySignificant:  int(metrics.is_statistically_significant),
				CompressionRatioChange:      float64(metrics.compression_ratio_change),
				EntropyChange:               float64(metrics.entropy_change),
				Explanation:                 explanation,
			}

			anomalies = append(anomalies, AnomalyResult{
				Transaction: tx,
				Metrics:     goMetrics,
				Explanation: explanation,
			})
		}

		if i%100 == 0 && i > 0 {
			fmt.Printf("✅ Processed %d transactions...\n", i)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\n🎉 Processing complete!\n")
	fmt.Printf("📊 Found %d anomalies out of %d transactions\n", len(anomalies), processingLimit)
	fmt.Printf("⚡ Processing time: %v\n\n", duration)

	// Generate HTML report
	outputPath := "demo-output.html"
	if err := generateHTMLReport(anomalies, transactions[:processingLimit], outputPath); err != nil {
		log.Fatalf("❌ Failed to generate HTML report: %v", err)
	}

	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("✅ Demo output written to: %s\n", absPath)
	fmt.Printf("🌐 Open this file in your browser to view results\n")
	fmt.Printf("\n💡 Tip: Use 'open %s' on macOS or 'xdg-open %s' on Linux\n", absPath, absPath)
}

func generateHTMLReport(anomalies []AnomalyResult, allTransactions []Transaction, outputPath string) error {
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
            justify-content: between;
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
                    <span class="number">{{.DetectionRate}}%</span>
                    <span class="label">Detection Rate</span>
                </div>
                <div class="stat-card">
                    <span class="number">{{.ProcessingTime}}</span>
                    <span class="label">Processing Time</span>
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
                        <div class="anomaly-badge">ANOMALY</div>
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
                            <span class="metric-value">{{printf "%.1f%%" .Metrics.ConfidenceLevel}}</span>
                        </div>
                        <div class="metric">
                            <span class="metric-label">Compression Δ</span>
                            <span class="metric-value">{{printf "%.3f" .Metrics.CompressionRatioChange}}</span>
                        </div>
                    </div>
                    
                    <div class="explanation">
                        <strong>🧠 Explanation:</strong> {{.Explanation}}
                    </div>
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

	// Parse template
	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("template parse error: %v", err)
	}

	// Prepare data
	data := struct {
		TotalTransactions int
		AnomalyCount      int
		DetectionRate     float64
		ProcessingTime    string
		Anomalies         []AnomalyResult
	}{
		TotalTransactions: len(allTransactions),
		AnomalyCount:      len(anomalies),
		DetectionRate:     float64(len(anomalies)) / float64(len(allTransactions)) * 100,
		ProcessingTime:    duration.String(),
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