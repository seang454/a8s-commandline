# a8s auth onboarding

Manage onboarding

## Usage

```text
a8s auth onboarding
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

- [a8s auth onboarding start](a8s_auth_onboarding_start.md) - POST /api/v1/auth/session/onboarding
- [a8s auth onboarding status](a8s_auth_onboarding_status.md) - GET /api/v1/auth/session/onboarding

