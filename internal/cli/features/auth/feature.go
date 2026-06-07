package auth

import (
	"github.com/spf13/cobra"
	"github.com/yourname/a8s/internal/cli/commands/catalogcmd"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

func Register(root *cobra.Command, runtime *cliruntime.Runtime) {
	group := &cobra.Command{Use: "auth", Short: "Authenticate and manage the current session"}
	root.AddCommand(group)
	addSessionCommands(group, runtime)
	catalogcmd.RegisterRoutes(root, runtime, Routes)
}
