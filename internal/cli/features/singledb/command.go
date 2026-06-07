package singledb

import (
	"github.com/spf13/cobra"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

func NewCommand(runtime *cliruntime.Runtime) *cobra.Command {
	command := &cobra.Command{
		Use:   "database",
		Short: "Manage single database deployments",
	}
	command.AddCommand(newDeployCommand(runtime))
	return command
}
