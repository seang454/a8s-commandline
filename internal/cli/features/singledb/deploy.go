package singledb

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/yourname/a8s/internal/api/resources/databases"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/operation"
	databaseoperation "github.com/yourname/a8s/internal/operation/kinds/database"
	"github.com/yourname/a8s/internal/workflow/deployment"
)

type deployFlags struct {
	file               string
	releaseName        string
	projectName        string
	engine             string
	deploymentMode     string
	databaseName       string
	username           string
	version            string
	sizeProfile        string
	storageSize        string
	storageClass       string
	environment        string
	existingAuthSecret string
	networkPolicy      bool
	tls                bool
	requireSSL         bool
	tlsSecret          string
	includeCA          bool
	passwordEnv        string
	passwordStdin      bool
	dryRun             bool
	wait               bool
}

func newDeployCommand(runtime *cliruntime.Runtime) *cobra.Command {
	var flags deployFlags
	command := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a single database using flags or an operation file",
		Example: `  a8s database deploy --file database.yaml --wait
  a8s database deploy --project-name payments --engine postgresql --database-name payments --version 16 --password-env DATABASE_PASSWORD --wait`,
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := operation.LoadFile[databaseoperation.Deploy](flags.file, databaseoperation.DeployKind, runtime.In)
			if err != nil {
				return err
			}
			spec.Apply(databaseoperation.Overrides{
				Changed:                changedFlags(cmd),
				ReleaseName:            flags.releaseName,
				ProjectName:            flags.projectName,
				Engine:                 flags.engine,
				DeploymentMode:         flags.deploymentMode,
				DatabaseName:           flags.databaseName,
				Username:               flags.username,
				Version:                flags.version,
				SizeProfile:            flags.sizeProfile,
				StorageSize:            flags.storageSize,
				StorageClassName:       flags.storageClass,
				Environment:            flags.environment,
				ExistingAuthSecretName: flags.existingAuthSecret,
				NetworkPolicyEnabled:   flags.networkPolicy,
				TLSEnabled:             flags.tls,
				RequireSSL:             flags.requireSSL,
				TLSSecret:              flags.tlsSecret,
				IncludeCA:              flags.includeCA,
			})
			spec.ApplyDefaults()
			if err := spec.Validate(); err != nil {
				return err
			}

			passwordEnv := flags.passwordEnv
			if passwordEnv == "" && spec.PasswordFrom != nil {
				passwordEnv = spec.PasswordFrom.Env
			}
			if flags.dryRun {
				request := spec.BackendRequest("")
				if passwordEnv != "" || flags.passwordStdin {
					request.Password = "[redacted]"
				}
				return runtime.Printer.Print(request)
			}
			password, err := operation.ResolveSecret(passwordEnv, flags.passwordStdin, runtime.In)
			if err != nil {
				return err
			}
			request := spec.BackendRequest(password)

			ctx, cancel := context.WithTimeout(cmd.Context(), runtime.Config.Timeout)
			defer cancel()
			client := databases.New(runtime.API)
			result, err := client.Deploy(ctx, request)
			if err != nil {
				return err
			}
			if flags.wait {
				id := firstNonEmpty(result.DeploymentID, result.ID)
				if id == "" {
					return clierrors.New("unexpected_response", "backend accepted deployment without returning a deployment ID", 1)
				}
				result, err = deployment.WaitForDatabase(ctx, client, id, runtime.Config.PollingInterval, runtime.Printer.Progress)
				if err != nil {
					return err
				}
			}
			return runtime.Printer.Print(result)
		},
	}

	f := command.Flags()
	f.StringVar(&flags.file, "file", "", "YAML or JSON operation file; use - for stdin")
	f.StringVar(&flags.releaseName, "release-name", "", "deployment release name")
	f.StringVar(&flags.projectName, "project-name", "", "A8S project name")
	f.StringVar(&flags.engine, "engine", "", "database engine")
	f.StringVar(&flags.deploymentMode, "deployment-mode", "", "deployment mode")
	f.StringVar(&flags.databaseName, "database-name", "", "initial database name")
	f.StringVar(&flags.username, "username", "", "initial application username")
	f.StringVar(&flags.version, "version", "", "database version")
	f.StringVar(&flags.sizeProfile, "size-profile", "", "size profile")
	f.StringVar(&flags.storageSize, "storage-size", "", "persistent storage size")
	f.StringVar(&flags.storageClass, "storage-class", "", "Kubernetes storage class")
	f.StringVar(&flags.environment, "environment", "", "deployment environment")
	f.StringVar(&flags.existingAuthSecret, "existing-auth-secret", "", "existing authentication secret")
	f.BoolVar(&flags.networkPolicy, "network-policy", false, "enable network policy")
	f.BoolVar(&flags.tls, "tls", false, "enable TLS")
	f.BoolVar(&flags.requireSSL, "require-ssl", false, "require SSL connections")
	f.StringVar(&flags.tlsSecret, "tls-secret", "", "existing TLS secret")
	f.BoolVar(&flags.includeCA, "include-ca", false, "include CA certificate")
	f.StringVar(&flags.passwordEnv, "password-env", "", "read database password from an environment variable")
	f.BoolVar(&flags.passwordStdin, "password-stdin", false, "read database password from stdin")
	f.BoolVar(&flags.dryRun, "dry-run", false, "validate and print the final request without applying it")
	f.BoolVar(&flags.wait, "wait", false, "wait for the deployment to reach a terminal state")
	return command
}

func changedFlags(command *cobra.Command) map[string]bool {
	changed := map[string]bool{}
	command.Flags().Visit(func(flag *pflag.Flag) {
		changed[flag.Name] = true
	})
	return changed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
