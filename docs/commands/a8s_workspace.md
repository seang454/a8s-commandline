# a8s workspace

Manage workspace

## Usage

```text
a8s workspace
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

- [a8s workspace bootstrap](a8s_workspace_bootstrap.md) - POST /api/v1/workspaces/bootstrap
- [a8s workspace entitlements](a8s_workspace_entitlements.md) - GET /api/v1/workspaces/entitlements
- [a8s workspace quota](a8s_workspace_quota.md) - Manage quota
- [a8s workspace status](a8s_workspace_status.md) - GET /api/v1/workspaces/bootstrap

