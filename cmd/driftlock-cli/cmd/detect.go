package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	detectStdin bool
)

// Response structs for pretty printing
type DetectResponse struct {
	Success        bool            `json:"success"`
	StreamID       string          `json:"stream_id"`
	TotalEvents    int             `json:"total_events"`
	AnomalyCount   int             `json:"anomaly_count"`
	ProcessingTime string          `json:"processing_time"`
	Anomalies      []AnomalyOutput `json:"anomalies"`
	Status         string          `json:"status"`
}

type AnomalyOutput struct {
	ID       string                 `json:"id"`
	Index    int                    `json:"index"`
	Metrics  map[string]interface{} `json:"metrics"`
	Why      string                 `json:"why"`
	Detected bool                   `json:"detected"`
}

var detectCmd = &cobra.Command{
	Use:   "detect [file]",
	Short: "Run anomaly detection on a file",
	Long: `Upload a file (JSON/NDJSON) to the Driftlock API for anomaly detection.
Example: driftlock detect logs.json`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !detectStdin {
			fmt.Println("❌ provide a file path or pass --stdin to read from STDIN")
			return
		}

		filePath := ""
		if len(args) > 0 {
			filePath = args[0]
		}
		apiKey := viper.GetString("api_key")
		apiURL := viper.GetString("api_url")

		if apiKey == "" {
			fmt.Println("❌ Not logged in. Run 'driftlock login' first.")
			return
		}

		var bodyReader io.ReadCloser
		if detectStdin || filePath == "-" {
			bodyReader = io.NopCloser(os.Stdin)
			filePath = "stdin"
		} else {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("❌ Error opening file: %v\n", err)
				return
			}
			bodyReader = file
			defer file.Close()
		}

		start := time.Now()
		req, err := http.NewRequest("POST", apiURL+"/detect?algo=zstd", bodyReader)
		if err != nil {
			fmt.Printf("❌ Error creating request: %v\n", err)
			return
		}

		// Detect content type
		ext := filepath.Ext(filePath)
		contentType := "application/json"
		if ext == ".ndjson" || ext == ".jsonl" {
			contentType = "application/x-ndjson"
		} else if filePath == "stdin" {
			contentType = "application/octet-stream"
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("X-Api-Key", apiKey)

		client := &http.Client{}
		fmt.Print("⏳ Analyzing... ")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("\n❌ Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("\n❌ API Error (%d): %s\n", resp.StatusCode, string(respBody))
			return
		}

		fmt.Printf("Done (%v)\n", time.Since(start).Round(time.Millisecond))

		var result DetectResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			fmt.Println("⚠️  Could not parse response JSON. Raw output:")
			fmt.Println(string(respBody))
			return
		}

		printSummary(result)
	},
}

func printSummary(r DetectResponse) {
	fmt.Println("\n📊 Detection Results")
	fmt.Println("------------------------------------------------")
	fmt.Printf("Stream ID:       %s\n", r.StreamID)
	fmt.Printf("Events:          %d\n", r.TotalEvents)
	fmt.Printf("Server Time:     %s\n", r.ProcessingTime)

	if r.Status == "calibrating" {
		fmt.Printf("Status:          🚧 Calibrating (Send more data)\n")
	} else {
		fmt.Printf("Status:          ✅ Active\n")
	}

	if r.AnomalyCount > 0 {
		fmt.Printf("\n🚨 Anomalies Detected: %d\n", r.AnomalyCount)
		for _, a := range r.Anomalies {
			fmt.Printf("   - [Idx %d] %s\n", a.Index, a.Why)
		}
	} else {
		fmt.Printf("\n✅ No anomalies detected. System nominal.\n")
	}
	fmt.Println("------------------------------------------------")
}

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().BoolVar(&detectStdin, "stdin", false, "Read payload from STDIN instead of a file")
}
