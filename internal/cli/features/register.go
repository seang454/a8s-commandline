package features

import (
	"github.com/spf13/cobra"

	"github.com/yourname/a8s/internal/cli/catalog"
	"github.com/yourname/a8s/internal/cli/commands/catalogcmd"
	"github.com/yourname/a8s/internal/cli/commands/watchcmd"
	"github.com/yourname/a8s/internal/cli/features/admin"
	"github.com/yourname/a8s/internal/cli/features/alerts"
	"github.com/yourname/a8s/internal/cli/features/auth"
	"github.com/yourname/a8s/internal/cli/features/databasebackup"
	"github.com/yourname/a8s/internal/cli/features/databaseconsole"
	"github.com/yourname/a8s/internal/cli/features/dbcluster"
	documentationfeature "github.com/yourname/a8s/internal/cli/features/documentation"
	"github.com/yourname/a8s/internal/cli/features/entitlements"
	"github.com/yourname/a8s/internal/cli/features/gitintegration"
	"github.com/yourname/a8s/internal/cli/features/imagescanner"
	"github.com/yourname/a8s/internal/cli/features/microservice"
	"github.com/yourname/a8s/internal/cli/features/monitoring"
	"github.com/yourname/a8s/internal/cli/features/monolithic"
	"github.com/yourname/a8s/internal/cli/features/notifications"
	"github.com/yourname/a8s/internal/cli/features/payments"
	"github.com/yourname/a8s/internal/cli/features/profile"
	"github.com/yourname/a8s/internal/cli/features/projects"
	"github.com/yourname/a8s/internal/cli/features/singledb"
	"github.com/yourname/a8s/internal/cli/features/sonarqube"
	"github.com/yourname/a8s/internal/cli/features/testingkit"
	"github.com/yourname/a8s/internal/cli/features/workspaces"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

var Names = []string{
	"admin",
	"alerts",
	"auth",
	"databasebackup",
	"databaseconsole",
	"dbcluster",
	"documentation",
	"entitlements",
	"gitintegration",
	"imagescanner",
	"microservice",
	"monitoring",
	"monolithic",
	"notifications",
	"payments",
	"profile",
	"projects",
	"singledb",
	"sonarqube",
	"testingkit",
	"workspaces",
}

// RegisterAll is the single inventory of backend feature packages exposed by the CLI.
func RegisterAll(root *cobra.Command, runtime *cliruntime.Runtime) {
	// singledb owns the specialized database root and must register first.
	singledb.Register(root, runtime)

	admin.Register(root, runtime)
	alerts.Register(root, runtime)
	auth.Register(root, runtime)
	databasebackup.Register(root, runtime)
	databaseconsole.Register(root, runtime)
	dbcluster.Register(root, runtime)
	documentationfeature.Register(root, runtime)
	entitlements.Register(root, runtime)
	gitintegration.Register(root, runtime)
	imagescanner.Register(root, runtime)
	microservice.Register(root, runtime)
	monitoring.Register(root, runtime)
	monolithic.Register(root, runtime)
	notifications.Register(root, runtime)
	payments.Register(root, runtime)
	profile.Register(root, runtime)
	projects.Register(root, runtime)
	sonarqube.Register(root, runtime)
	testingkit.Register(root, runtime)
	workspaces.Register(root, runtime)

	catalogcmd.RegisterUtilities(root, runtime)
	watchcmd.Register(root, runtime)
	root.AddCommand(newInventoryCommand(runtime))
}

func newInventoryCommand(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "features",
		Short: "List backend features exposed by the CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			counts := map[string]int{}
			for _, route := range catalog.Routes {
				counts[route.Feature]++
			}
			rows := make([]map[string]any, 0, len(Names))
			for _, name := range Names {
				note := ""
				switch name {
				case "databaseconsole":
					note = "routes owned by singledb and dbcluster"
				case "payments":
					note = "routes exposed through workspaces"
				}
				rows = append(rows, map[string]any{"feature": name, "routes": counts[name], "note": note})
			}
			return runtime.Printer.Print(rows)
		},
	}
}
