package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Driftlock",
	Long: `Log in to your Driftlock account by providing your API Key.
You can find your API Key in the dashboard at https://driftlock.net/dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("🔑 Driftlock CLI Authentication")
		fmt.Println("--------------------------------")
		fmt.Println("Please enter your API Key from the dashboard.")
		fmt.Print("API Key: ")

		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		if apiKey == "" {
			fmt.Println("❌ API Key cannot be empty.")
			return
		}

		viper.Set("api_key", apiKey)

		// Ensure config dir exists
		home, _ := os.UserHomeDir()
		configDir := fmt.Sprintf("%s/.driftlock", home)
		os.MkdirAll(configDir, 0755)

		err := viper.WriteConfigAs(fmt.Sprintf("%s/config.json", configDir))
		if err != nil {
			fmt.Printf("❌ Error saving config: %v\n", err)
			return
		}

		fmt.Println("\n✅ Successfully logged in! Configuration saved to", configDir)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
