package auth

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"

	internalauth "github.com/yourname/a8s/internal/auth"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/credentials"
)

func addSessionCommands(group *cobra.Command, runtime *cliruntime.Runtime) {
	group.AddCommand(newLoginCommand(runtime), newStatusCommand(runtime), newLogoutCommand(runtime))
}

func newLoginCommand(runtime *cliruntime.Runtime) *cobra.Command {
	var noBrowser bool
	var callbackPort int
	var loginTimeout time.Duration
	command := &cobra.Command{
		Use:   "login",
		Short: "Authenticate through Keycloak using browser PKCE",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.Auth == nil {
				return clierrors.New("authentication_required", "credential storage is unavailable", 3)
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), loginTimeout)
			defer cancel()
			record, err := runtime.Auth.Login(ctx, runtime.Config, internalauth.LoginOptions{NoBrowser: noBrowser, CallbackPort: callbackPort, Out: runtime.Out})
			if err != nil {
				return err
			}
			return runtime.Printer.Print(statusOutput(runtime, record))
		},
	}
	command.Flags().BoolVar(&noBrowser, "no-browser", false, "print the login URL without opening a browser")
	command.Flags().IntVar(&callbackPort, "callback-port", 0, "fixed local callback port; Keycloak redirect URI must allow http://127.0.0.1:<port>/callback")
	command.Flags().DurationVar(&loginTimeout, "login-timeout", 5*time.Minute, "maximum time to complete browser authentication")
	return command
}

func newStatusCommand(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status without displaying tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.Auth == nil {
				return runtime.Printer.Print(map[string]any{"status": "not-authenticated", "context": runtime.Config.ContextName})
			}
			record, err := runtime.Auth.Status(runtime.Config)
			if errors.Is(err, credentials.ErrNotFound) {
				return runtime.Printer.Print(map[string]any{"status": "not-authenticated", "context": runtime.Config.ContextName, "server": runtime.Config.Server})
			}
			if err != nil {
				return err
			}
			return runtime.Printer.Print(statusOutput(runtime, record))
		},
	}
}

func newLogoutCommand(runtime *cliruntime.Runtime) *cobra.Command {
	var noBrowser bool
	var callbackPort int
	var keycloak bool
	var logoutTimeout time.Duration
	command := &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials for the active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			remoteLogout := keycloak || callbackPort > 0 || noBrowser
			remoteStatus := "skipped"
			remoteError := ""
			if runtime.Auth != nil {
				if remoteLogout {
					ctx, cancel := context.WithTimeout(cmd.Context(), logoutTimeout)
					defer cancel()
					result, err := runtime.Auth.EndSession(ctx, runtime.Config, internalauth.LogoutOptions{NoBrowser: noBrowser, CallbackPort: callbackPort, Out: runtime.Out})
					if result.RemoteAttempted {
						remoteStatus = "completed"
					}
					if err != nil {
						remoteStatus = "failed"
						remoteError = err.Error()
					}
				}
				if err := runtime.Auth.Logout(runtime.Config); err != nil {
					return err
				}
			}
			output := map[string]any{"status": "logged-out", "context": runtime.Config.ContextName, "keycloakLogout": remoteStatus}
			if remoteError != "" {
				output["keycloakLogoutError"] = remoteError
			}
			return runtime.Printer.Print(output)
		},
	}
	command.Flags().BoolVar(&keycloak, "keycloak", false, "also end the Keycloak browser session")
	command.Flags().BoolVar(&noBrowser, "no-browser", false, "print the Keycloak logout URL without opening a browser; implies --keycloak")
	command.Flags().IntVar(&callbackPort, "callback-port", 0, "fixed local logout callback port; Keycloak post logout redirect URI must allow http://127.0.0.1:<port>/callback")
	command.Flags().DurationVar(&logoutTimeout, "logout-timeout", 2*time.Minute, "maximum time to complete browser logout")
	return command
}

func statusOutput(runtime *cliruntime.Runtime, record credentials.Record) map[string]any {
	status := "authenticated"
	if !record.AccessTokenExpiry.IsZero() && record.AccessTokenExpiry.Before(time.Now()) {
		status = "expired"
	}
	return map[string]any{
		"status": status, "context": runtime.Config.ContextName, "server": runtime.Config.Server,
		"issuer": record.Issuer, "subject": record.Subject, "username": record.Username,
		"email": record.Email, "roles": record.Roles, "accessTokenExpiry": record.AccessTokenExpiry,
	}
}
