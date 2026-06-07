# a8s database backup

Manage backup

## Usage

```text
a8s database backup
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

- [a8s database backup delete](a8s_database_backup_delete.md) - DELETE /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}
- [a8s database backup download](a8s_database_backup_download.md) - GET /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/download
- [a8s database backup restore](a8s_database_backup_restore.md) - POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore
- [a8s database backup run](a8s_database_backup_run.md) - POST /api/v1/database-deployments/{deploymentId}/backup/run
- [a8s database backup settings](a8s_database_backup_settings.md) - Manage settings

