package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current user info",
	Long:  `Shows the currently authenticated tenant/user and plan status.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("api_key")
		apiURL := viper.GetString("api_url")

		if apiKey == "" {
			fmt.Println("❌ Not logged in. Run 'driftlock login' first.")
			return
		}

		// Assuming endpoint /v1/me/usage or similar exists for user info
		// Based on backend checks, we might check /v1/me/usage or /healthz (authenticated)
		req, _ := http.NewRequest("GET", apiURL+"/me/usage", nil)
		req.Header.Set("X-Api-Key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error checking status: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("❌ Failed to fetch info (%d)\n", resp.StatusCode)
			return
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
