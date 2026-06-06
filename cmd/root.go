package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourname/a8s/internal/config"
)

var (
	cfgFile   string
	outputFmt string
	apiURL    string
	apiToken  string
)

var rootCmd = &cobra.Command{
	Use:   "a8s",
	Short: "A8S — Your Platform CLI",
	Long: `A8S is a command-line interface to interact with your platform API.
Built with Go + Cobra following the enterprise standard used by kubectl, gh, and helm.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.a8s.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "output format: table|json")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL (overrides config)")
	rootCmd.PersistentFlags().StringVar(&apiToken, "token", "", "API token (overrides config)")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	config.Load(cfgFile)

	if apiURL != "" {
		viper.Set("api_url", apiURL)
	}
	if apiToken != "" {
		viper.Set("api_token", apiToken)
	}
}
