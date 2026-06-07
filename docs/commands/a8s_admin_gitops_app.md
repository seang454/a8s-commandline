# a8s admin gitops app

Manage app

## Usage

```text
a8s admin gitops app
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

- [a8s admin gitops app abort](a8s_admin_gitops_app_abort.md) - POST /api/v1/admin/gitops/apps/{appId}/abort
- [a8s admin gitops app create](a8s_admin_gitops_app_create.md) - POST /api/v1/admin/gitops/apps
- [a8s admin gitops app retry](a8s_admin_gitops_app_retry.md) - POST /api/v1/admin/gitops/apps/{appId}/retry
- [a8s admin gitops app sync](a8s_admin_gitops_app_sync.md) - POST /api/v1/admin/gitops/apps/{appId}/sync

