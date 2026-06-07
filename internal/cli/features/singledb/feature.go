package singledb

import (
	"github.com/spf13/cobra"
	"github.com/yourname/a8s/internal/cli/commands/catalogcmd"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

func Register(root *cobra.Command, runtime *cliruntime.Runtime) {
	root.AddCommand(NewCommand(runtime))
	catalogcmd.RegisterRoutes(root, runtime, Routes)
}
