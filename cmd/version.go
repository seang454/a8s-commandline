package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourname/a8s/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("a8s version %s (build: %s)\n", version.Version, version.BuildDate)
	},
}
