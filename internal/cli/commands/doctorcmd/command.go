package doctorcmd

import (
	"net/http"

	"github.com/spf13/cobra"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
)

func NewCommand(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check CLI configuration and backend connectivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			checks := []map[string]any{
				{"check": "configuration", "status": "ok", "path": runtime.Config.ConfigPath},
				{"check": "context", "status": status(runtime.Config.ContextName != ""), "value": runtime.Config.ContextName},
				{"check": "server", "status": status(runtime.Config.Server != ""), "value": runtime.Config.Server},
				{"check": "namespace", "status": status(runtime.Config.Namespace != ""), "value": runtime.Config.Namespace},
				{"check": "authentication", "status": status(runtime.Config.Token != "")},
			}
			response, err := runtime.API.Do(cmd.Context(), http.MethodGet, "/actuator/health", nil)
			if err != nil {
				checks = append(checks, map[string]any{"check": "backend", "status": "failed", "error": err.Error()})
			} else {
				response.Body.Close()
				checks = append(checks, map[string]any{"check": "backend", "status": "ok", "httpStatus": response.StatusCode})
			}
			return runtime.Printer.Print(checks)
		},
	}
}

func status(ok bool) string {
	if ok {
		return "ok"
	}
	return "missing"
}
