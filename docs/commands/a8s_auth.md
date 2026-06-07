# a8s auth

Authenticate and manage the current session

## Usage

```text
a8s auth
```

## Inherited Flags

- `--config` `string` - config file path
- `--context` `string` - named context to use
- `--namespace` `string` - workspace or Kubernetes namespace
- `--request-timeout` `string` - single HTTP request timeout
- `--server` `string` - backend server URL
- `--target-cluster` `string` - managed Kubernetes cluster alias
- `--timeout` `string` - complete command timeout
- `--token` `string` - temporary bearer token; prefer A8S_TOKEN
- `-o, --output` `string` - output format: table|json|yaml

## Child Commands

- [a8s auth login](a8s_auth_login.md) - Authenticate through Keycloak using browser PKCE
- [a8s auth logout](a8s_auth_logout.md) - Clear stored credentials for the active context
- [a8s auth onboarding](a8s_auth_onboarding.md) - Manage onboarding
- [a8s auth status](a8s_auth_status.md) - Show authentication status without displaying tokens
- [a8s auth verify-email](a8s_auth_verify-email.md) - Manage email verification for the authenticated user

