package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
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
	root.AddCommand(newListCommand(root, runtime))
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
		"a8s list",
		"a8s manifest",
		"a8s version",
	} {
		if path == prefix || strings.HasPrefix(path, prefix+" ") {
			return true
		}
	}
	return false
}

type commandInfo struct {
	Command     string `json:"command" yaml:"command"`
	Usage       string `json:"usage" yaml:"usage"`
	Description string `json:"description" yaml:"description"`
	Help        string `json:"help" yaml:"help"`
}

func newListCommand(root *cobra.Command, runtime *cliruntime.Runtime) *cobra.Command {
	list := &cobra.Command{Use: "list", Short: "List CLI commands and local inventory"}
	list.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "List every available runnable command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runtime.Printer.Print(runnableCommands(root))
		},
	})

	// sections: group commands by their top-level section (e.g. 'cluster', 'project')
	sectionsCmd := &cobra.Command{
		Use:   "sections",
		Short: "List commands grouped by top-level section",
		RunE: func(cmd *cobra.Command, args []string) error {
			paginate, _ := cmd.Flags().GetBool("paginate")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			// If output is json or yaml, print structured grouped output
			if runtime.Printer.Format == "json" || runtime.Printer.Format == "yaml" {
				grouped := map[string][]commandInfo{}
				for _, c := range runnableCommands(root) {
					parts := strings.Fields(c.Command)
					section := "other"
					if len(parts) >= 2 {
						section = parts[1]
					}
					grouped[section] = append(grouped[section], c)
				}
				return runtime.Printer.Print(grouped)
			}

			// Human readable output with optional pagination
			sections := map[string][]commandInfo{}
			for _, c := range runnableCommands(root) {
				parts := strings.Fields(c.Command)
				section := "other"
				if len(parts) >= 2 {
					section = parts[1]
				}
				sections[section] = append(sections[section], c)
			}
			var names []string
			for n := range sections {
				names = append(names, n)
			}
			sort.Strings(names)

			reader := bufio.NewReader(runtime.In)
			for _, name := range names {
				fmt.Fprintf(runtime.Out, "%s:\n", name)
				cmds := sections[name]
				for i, ci := range cmds {
					fmt.Fprintf(runtime.Out, "  %s\n    %s\n", ci.Usage, ci.Description)
					if paginate && (i+1)%pageSize == 0 {
						fmt.Fprint(runtime.Out, "-- more -- press Enter to continue, 'q' to quit: ")
						input, err := reader.ReadString('\n')
						if err != nil {
							return nil
						}
						if strings.TrimSpace(input) == "q" {
							return nil
						}
					}
				}
				fmt.Fprintln(runtime.Out)
				if paginate {
					fmt.Fprint(runtime.Out, "-- press Enter to continue to next section, 'q' to quit: ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return nil
					}
					if strings.TrimSpace(input) == "q" {
						return nil
					}
				}
			}
			return nil
		},
	}
	sectionsCmd.Flags().BoolP("paginate", "p", false, "Paginate output interactively")
	sectionsCmd.Flags().IntP("page-size", "n", 20, "Number of command lines per page when paginating")
	list.AddCommand(sectionsCmd)
	return list
}

func runnableCommands(root *cobra.Command) []commandInfo {
	var result []commandInfo
	var walk func(command *cobra.Command)
	walk = func(command *cobra.Command) {
		if command.Hidden {
			return
		}
		if command.Run != nil || command.RunE != nil {
			path := command.CommandPath()
			result = append(result, commandInfo{
				Command:     path,
				Usage:       command.UseLine(),
				Description: command.Short,
				Help:        path + " --help",
			})
		}
		for _, child := range command.Commands() {
			walk(child)
		}
	}
	walk(root)
	sort.Slice(result, func(left, right int) bool {
		return result[left].Command < result[right].Command
	})
	return result
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
