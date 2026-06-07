package contextcmd

import (
	"sort"

	"github.com/spf13/cobra"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/config"
	"github.com/yourname/a8s/internal/operation"
)

type contextFlags struct {
	file          string
	server        string
	namespace     string
	targetCluster string
	issuer        string
	clientID      string
}

type contextSpec struct {
	Server        string   `yaml:"server,omitempty" json:"server,omitempty"`
	Namespace     string   `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	TargetCluster string   `yaml:"targetCluster,omitempty" json:"targetCluster,omitempty"`
	TLS           tlsSpec  `yaml:"tls,omitempty" json:"tls,omitempty"`
	Auth          authSpec `yaml:"auth,omitempty" json:"auth,omitempty"`
}

type contextPatchSpec struct {
	Server        *string        `yaml:"server,omitempty" json:"server,omitempty"`
	Namespace     *string        `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	TargetCluster *string        `yaml:"targetCluster,omitempty" json:"targetCluster,omitempty"`
	TLS           *tlsPatchSpec  `yaml:"tls,omitempty" json:"tls,omitempty"`
	Auth          *authPatchSpec `yaml:"auth,omitempty" json:"auth,omitempty"`
}

type tlsSpec struct {
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	CAFile             string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
}

type tlsPatchSpec struct {
	InsecureSkipVerify *bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	CAFile             *string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
}

type authSpec struct {
	Issuer        string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	ClientID      string `yaml:"clientId,omitempty" json:"clientId,omitempty"`
	CredentialKey string `yaml:"credentialKey,omitempty" json:"credentialKey,omitempty"`
}

type authPatchSpec struct {
	Issuer        *string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	ClientID      *string `yaml:"clientId,omitempty" json:"clientId,omitempty"`
	CredentialKey *string `yaml:"credentialKey,omitempty" json:"credentialKey,omitempty"`
}

func NewCommand(runtime *cliruntime.Runtime) *cobra.Command {
	command := &cobra.Command{Use: "context", Short: "Manage named backend environments"}
	command.AddCommand(newList(runtime), newGet(runtime), newUse(runtime), newCreate(runtime), newUpdate(runtime), newDelete(runtime))
	return command
}

func newList(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			names := make([]string, 0, len(file.Contexts))
			for name := range file.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)
			rows := make([]map[string]any, 0, len(names))
			for _, name := range names {
				value := file.Contexts[name]
				rows = append(rows, map[string]any{"name": name, "current": name == file.CurrentContext, "server": value.Server, "namespace": value.Namespace, "targetCluster": value.TargetCluster})
			}
			return runtime.Printer.Print(rows)
		},
	}
}

func newGet(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get a configured context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			value, ok := file.Contexts[args[0]]
			if !ok {
				return clierrors.New("not_found", "context does not exist", 5)
			}
			return runtime.Printer.Print(value)
		},
	}
}

func newUse(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Set the default context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			if _, ok := file.Contexts[args[0]]; !ok {
				return clierrors.New("not_found", "context does not exist", 5)
			}
			file.CurrentContext = args[0]
			if err := config.Save(runtime.Config.ConfigPath, file); err != nil {
				return err
			}
			return runtime.Printer.Print(map[string]string{"currentContext": args[0]})
		},
	}
}

func newCreate(runtime *cliruntime.Runtime) *cobra.Command {
	var flags contextFlags
	command := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a named context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := config.Context{}
			if flags.file != "" {
				spec, err := operation.LoadFile[contextSpec](flags.file, "Context", runtime.In)
				if err != nil {
					return err
				}
				value = contextFromSpec(spec)
			}
			applyContextFlagOverrides(cmd, &value, flags)
			defaultContextAuth(args[0], &value)
			if value.Server == "" {
				return clierrors.Validation("--server or spec.server is required")
			}

			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			if _, exists := file.Contexts[args[0]]; exists {
				return clierrors.New("conflict", "context already exists", 6)
			}
			file.Contexts[args[0]] = value
			if len(file.Contexts) == 1 {
				file.CurrentContext = args[0]
			}
			if err := config.Save(runtime.Config.ConfigPath, file); err != nil {
				return err
			}
			return runtime.Printer.Print(file.Contexts[args[0]])
		},
	}
	addContextFlags(command, &flags)
	return command
}

func newUpdate(runtime *cliruntime.Runtime) *cobra.Command {
	var flags contextFlags
	command := &cobra.Command{
		Use:   "update <name>",
		Short: "Update a named context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			value, exists := file.Contexts[args[0]]
			if !exists {
				return clierrors.New("not_found", "context does not exist", 5)
			}
			if flags.file != "" {
				patch, err := operation.LoadFile[contextPatchSpec](flags.file, "ContextPatch", runtime.In)
				if err != nil {
					return err
				}
				applyContextPatch(&value, patch)
			}
			applyContextFlagOverrides(cmd, &value, flags)
			defaultContextAuth(args[0], &value)
			file.Contexts[args[0]] = value
			if err := config.Save(runtime.Config.ConfigPath, file); err != nil {
				return err
			}
			return runtime.Printer.Print(value)
		},
	}
	addContextFlags(command, &flags)
	return command
}

func newDelete(runtime *cliruntime.Runtime) *cobra.Command {
	var yes bool
	command := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return clierrors.Validation("context deletion requires --yes")
			}
			file, err := config.LoadFile(runtime.Config.ConfigPath)
			if err != nil {
				return err
			}
			if _, exists := file.Contexts[args[0]]; !exists {
				return clierrors.New("not_found", "context does not exist", 5)
			}
			delete(file.Contexts, args[0])
			if file.CurrentContext == args[0] {
				file.CurrentContext = ""
			}
			return config.Save(runtime.Config.ConfigPath, file)
		},
	}
	command.Flags().BoolVar(&yes, "yes", false, "confirm context deletion")
	return command
}

func addContextFlags(command *cobra.Command, flags *contextFlags) {
	command.Flags().StringVar(&flags.file, "file", "", "YAML or JSON context operation file; use - for stdin")
	command.Flags().StringVar(&flags.server, "server", "", "backend server URL")
	command.Flags().StringVar(&flags.namespace, "namespace", "", "default namespace")
	command.Flags().StringVar(&flags.targetCluster, "target-cluster", "", "default managed cluster alias")
	command.Flags().StringVar(&flags.issuer, "issuer", "", "OIDC issuer")
	command.Flags().StringVar(&flags.clientID, "client-id", "a8s-cli", "OIDC public client ID")
}

func contextFromSpec(spec contextSpec) config.Context {
	return config.Context{
		Server:        spec.Server,
		Namespace:     spec.Namespace,
		TargetCluster: spec.TargetCluster,
		TLS: config.TLS{
			InsecureSkipVerify: spec.TLS.InsecureSkipVerify,
			CAFile:             spec.TLS.CAFile,
		},
		Auth: config.Auth{
			Issuer:        spec.Auth.Issuer,
			ClientID:      spec.Auth.ClientID,
			CredentialKey: spec.Auth.CredentialKey,
		},
	}
}

func applyContextPatch(value *config.Context, patch contextPatchSpec) {
	if patch.Server != nil {
		value.Server = *patch.Server
	}
	if patch.Namespace != nil {
		value.Namespace = *patch.Namespace
	}
	if patch.TargetCluster != nil {
		value.TargetCluster = *patch.TargetCluster
	}
	if patch.TLS != nil {
		if patch.TLS.InsecureSkipVerify != nil {
			value.TLS.InsecureSkipVerify = *patch.TLS.InsecureSkipVerify
		}
		if patch.TLS.CAFile != nil {
			value.TLS.CAFile = *patch.TLS.CAFile
		}
	}
	if patch.Auth != nil {
		if patch.Auth.Issuer != nil {
			value.Auth.Issuer = *patch.Auth.Issuer
		}
		if patch.Auth.ClientID != nil {
			value.Auth.ClientID = *patch.Auth.ClientID
		}
		if patch.Auth.CredentialKey != nil {
			value.Auth.CredentialKey = *patch.Auth.CredentialKey
		}
	}
}

func applyContextFlagOverrides(cmd *cobra.Command, value *config.Context, flags contextFlags) {
	if cmd.Flags().Changed("server") {
		value.Server = flags.server
	}
	if cmd.Flags().Changed("namespace") {
		value.Namespace = flags.namespace
	}
	if cmd.Flags().Changed("target-cluster") {
		value.TargetCluster = flags.targetCluster
	}
	if cmd.Flags().Changed("issuer") {
		value.Auth.Issuer = flags.issuer
	}
	if cmd.Flags().Changed("client-id") {
		value.Auth.ClientID = flags.clientID
	}
}

func defaultContextAuth(name string, value *config.Context) {
	if value.Auth.ClientID == "" {
		value.Auth.ClientID = "a8s-cli"
	}
	if value.Auth.CredentialKey == "" {
		value.Auth.CredentialKey = "context:" + name
	}
}
