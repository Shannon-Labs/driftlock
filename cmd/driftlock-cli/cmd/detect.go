package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var detectCmd = &cobra.Command{
	Use:   "detect [file]",
	Short: "Run anomaly detection on a file",
	Long: `Upload a file (JSON/NDJSON) to the Driftlock API for anomaly detection.
Example: driftlock detect logs.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		apiKey := viper.GetString("api_key")
		apiURL := viper.GetString("api_url")

		if apiKey == "" {
			fmt.Println("❌ Not logged in. Run 'driftlock login' first.")
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("❌ Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		body := &bytes.Buffer{}
		// For /v1/detect endpoint, we assume it takes raw body or multipart. 
		// Based on previous context, /v1/detect usually accepts raw JSON/NDJSON body.
		// Let's stream the file content directly.
		
		// Ideally we stream directly to Request, but here we read to memory for simplicity in CLI MVP.
		// For production CLI with large files, we should use io.Pipe or pass file directly as Body.
		
		req, err := http.NewRequest("POST", apiURL+"/detect?algo=zstd", file)
		if err != nil {
			fmt.Printf("❌ Error creating request: %v\n", err)
			return
		}

		// Detect content type based on extension
		ext := filepath.Ext(filePath)
		contentType := "application/json"
		if ext == ".ndjson" || ext == ".jsonl" {
			contentType = "application/x-ndjson"
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("X-Api-Key", apiKey)

		client := &http.Client{}
		fmt.Println("⏳ Analyzing...")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ API Error (%d): %s\n", resp.StatusCode, string(respBody))
			return
		}

		fmt.Println("✅ Detection Complete:")
		fmt.Println(string(respBody))
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
