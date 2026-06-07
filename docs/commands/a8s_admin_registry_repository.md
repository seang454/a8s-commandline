# a8s admin registry repository

Manage repository

## Usage

```text
a8s admin registry repository
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

- [a8s admin registry repository delete](a8s_admin_registry_repository_delete.md) - DELETE /api/v1/admin/registry/projects/{projectName}/repositories
- [a8s admin registry repository list](a8s_admin_registry_repository_list.md) - GET /api/v1/admin/registry/projects/{projectName}/repositories

