package cli

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourname/a8s/internal/api"
	internalauth "github.com/yourname/a8s/internal/auth"
	"github.com/yourname/a8s/internal/cli/commands/contextcmd"
	"github.com/yourname/a8s/internal/cli/commands/doctorcmd"
	"github.com/yourname/a8s/internal/cli/commands/manifestcmd"
	"github.com/yourname/a8s/internal/cli/features"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/config"
	"github.com/yourname/a8s/internal/credentials"
	"github.com/yourname/a8s/internal/output"
	"github.com/yourname/a8s/pkg/version"
)

type globalFlags struct {
	configPath     string
	contextName    string
	server         string
	namespace      string
	targetCluster  string
	token          string
	output         string
	timeout        string
	requestTimeout string
}

func Execute() error {
	return NewRootCommand(os.Stdin, os.Stdout, os.Stderr).Execute()
}

func NewRootCommand(in io.Reader, out, errOut io.Writer) *cobra.Command {
	runtime := &cliruntime.Runtime{In: in, Out: out, ErrOut: errOut}
	if store, err := credentials.NewNativeStore(); err == nil {
		store.WarnFallbackTo(errOut)
		runtime.Auth = internalauth.NewManager(store)
	}
	var flags globalFlags

	root := &cobra.Command{
		Use:           "a8s",
		Short:         "A8S platform command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := config.Resolve(config.Overrides{
				ConfigPath:     flags.configPath,
				ContextName:    flags.contextName,
				Server:         flags.server,
				Namespace:      flags.namespace,
				TargetCluster:  flags.targetCluster,
				Token:          flags.token,
				Output:         flags.output,
				Timeout:        flags.timeout,
				RequestTimeout: flags.requestTimeout,
			})
			if err != nil {
				return err
			}
			runtime.Config = resolved
			runtime.Printer = output.Printer{Out: out, ErrOut: errOut, Format: resolved.Output}
			if isLocalCommand(cmd) {
				return nil
			}
			staticToken := resolved.Token != ""
			if runtime.Auth != nil && shouldResolveCredentials(cmd) {
				token, err := runtime.Auth.ResolveToken(cmd.Context(), resolved)
				if err != nil {
					return err
				}
				resolved.Token = token
			}
			runtime.API = api.New(resolved.Server, resolved.Token, resolved.RequestTimeout)
			if runtime.Auth != nil && !staticToken && resolved.Auth.CredentialKey != "" && shouldResolveCredentials(cmd) {
				runtime.API.RefreshToken = func(ctx context.Context) (string, error) {
					token, err := runtime.Auth.RefreshToken(ctx, resolved)
					if err == nil {
						runtime.Config.Token = token
					}
					return token, err
				}
			}
			if err := runtime.API.ConfigureTLS(resolved.TLS.InsecureSkipVerify, resolved.TLS.CAFile); err != nil {
				return err
			}
			return nil
		},
	}
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)

	persistent := root.PersistentFlags()
	persistent.StringVar(&flags.configPath, "config", "", "config file path")
	persistent.StringVar(&flags.contextName, "context", "", "named context to use")
	persistent.StringVar(&flags.server, "server", "", "backend server URL")
	persistent.StringVar(&flags.namespace, "namespace", "", "workspace or Kubernetes namespace")
	persistent.StringVar(&flags.targetCluster, "target-cluster", "", "managed Kubernetes cluster alias")
	persistent.StringVar(&flags.token, "token", "", "temporary bearer token; prefer A8S_TOKEN")
	persistent.StringVarP(&flags.output, "output", "o", "", "output format: table|json|yaml")
	persistent.StringVar(&flags.timeout, "timeout", "", "complete command timeout")
	persistent.StringVar(&flags.requestTimeout, "request-timeout", "", "single HTTP request timeout")

	root.AddCommand(contextcmd.NewCommand(runtime))
	root.AddCommand(doctorcmd.NewCommand(runtime))
	root.AddCommand(manifestcmd.NewCommand(runtime))
	root.AddCommand(newConfigCommand(runtime))
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runtime.Printer.Print(map[string]string{"version": version.Version, "buildDate": version.BuildDate})
		},
	})
	features.RegisterAll(root, runtime)
	return root
}

func shouldResolveCredentials(cmd *cobra.Command) bool {
	path := cmd.CommandPath()
	return path != "a8s auth login" && path != "a8s auth status" && path != "a8s auth logout" && path != "a8s doctor"
}

func isLocalCommand(cmd *cobra.Command) bool {
	path := cmd.CommandPath()
	for _, prefix := range []string{
		"a8s auth login",
		"a8s auth logout",
		"a8s auth status",
		"a8s config",
		"a8s context",
		"a8s features",
		"a8s manifest",
		"a8s version",
	} {
		if path == prefix || strings.HasPrefix(path, prefix+" ") {
			return true
		}
	}
	return false
}

func newConfigCommand(runtime *cliruntime.Runtime) *cobra.Command {
	configCmd := &cobra.Command{Use: "config", Short: "Inspect CLI configuration"}
	configCmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print the active configuration path",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := runtime.Out.Write([]byte(runtime.Config.ConfigPath + "\n"))
			return err
		},
	})
	configCmd.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "Print resolved non-secret configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runtime.Printer.Print(struct {
				ContextName   string `json:"contextName" yaml:"contextName"`
				Server        string `json:"server" yaml:"server"`
				Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
				TargetCluster string `json:"targetCluster,omitempty" yaml:"targetCluster,omitempty"`
			}{
				ContextName:   runtime.Config.ContextName,
				Server:        runtime.Config.Server,
				Namespace:     runtime.Config.Namespace,
				TargetCluster: runtime.Config.TargetCluster,
			})
		},
	})
	return configCmd
}

func commandContext(parent context.Context, timeout config.Resolved) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout.Timeout)
}
