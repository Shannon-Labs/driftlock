package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	detectStdin bool
)

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

		req, err := http.NewRequest("POST", apiURL+"/detect?algo=zstd", bodyReader)
		if err != nil {
			fmt.Printf("❌ Error creating request: %v\n", err)
			return
		}

		// Detect content type based on extension
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
	detectCmd.Flags().BoolVar(&detectStdin, "stdin", false, "Read payload from STDIN instead of a file")
}
