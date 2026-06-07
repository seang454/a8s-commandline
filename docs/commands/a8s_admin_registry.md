# a8s admin registry

Manage registry

## Usage

```text
a8s admin registry
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

- [a8s admin registry artifact](a8s_admin_registry_artifact.md) - Manage artifact
- [a8s admin registry health](a8s_admin_registry_health.md) - GET /api/v1/admin/registry/health
- [a8s admin registry project](a8s_admin_registry_project.md) - Manage project
- [a8s admin registry repository](a8s_admin_registry_repository.md) - Manage repository

