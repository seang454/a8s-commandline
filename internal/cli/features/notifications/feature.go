package notifications

import (
	"github.com/spf13/cobra"
	"github.com/yourname/a8s/internal/cli/commands/catalogcmd"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

func Register(root *cobra.Command, runtime *cliruntime.Runtime) {
	catalogcmd.RegisterRoutes(root, runtime, Routes)
}
