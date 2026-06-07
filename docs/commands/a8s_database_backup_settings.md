# a8s database backup settings

Manage settings

## Usage

```text
a8s database backup settings
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

- [a8s database backup settings get](a8s_database_backup_settings_get.md) - GET /api/v1/database-deployments/{deploymentId}/backup
- [a8s database backup settings set](a8s_database_backup_settings_set.md) - PATCH /api/v1/database-deployments/{deploymentId}/backup

