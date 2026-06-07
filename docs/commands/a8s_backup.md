# a8s backup

Manage backup

## Usage

```text
a8s backup
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

- [a8s backup delete](a8s_backup_delete.md) - DELETE /api/backups/{targetType}/{id}/{runId}
- [a8s backup download](a8s_backup_download.md) - GET /api/backups/download/{targetType}/{id}/{runId}
- [a8s backup restore](a8s_backup_restore.md) - POST /api/backups/restore/{targetType}/{id}/{runId}
- [a8s backup settings](a8s_backup_settings.md) - Manage settings
- [a8s backup trigger](a8s_backup_trigger.md) - POST /api/backups/trigger/{targetType}/{id}

