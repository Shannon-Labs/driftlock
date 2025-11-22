package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "driftlock",
	Short: "Driftlock CLI for anomaly detection",
	Long: `Driftlock CLI - The standard interface for the Driftlock SaaS platform.

Use this tool to:
- Authenticate with the Driftlock API
- Stream data for anomaly detection
- Manage your account and usage`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.driftlock/config.json)")
	rootCmd.PersistentFlags().String("api-url", "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1", "Driftlock API URL")
	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDir := filepath.Join(home, ".driftlock")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			os.MkdirAll(configDir, 0755)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Config loaded
	}
}
