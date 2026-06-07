# a8s database

Manage single database deployments

## Usage

```text
a8s database
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

- [a8s database backup](a8s_database_backup.md) - Manage backup
- [a8s database clone-from-backup](a8s_database_clone-from-backup.md) - POST /api/v1/database-deployments/clone-from-backup
- [a8s database console](a8s_database_console.md) - Manage console
- [a8s database credentials](a8s_database_credentials.md) - GET /api/v1/database-deployments/{deploymentId}/credentials
- [a8s database delete](a8s_database_delete.md) - DELETE /api/v1/database-deployments/{deploymentId}
- [a8s database deploy](a8s_database_deploy.md) - Deploy a single database using flags or an operation file
- [a8s database get](a8s_database_get.md) - GET /api/v1/database-deployments/{deploymentId}
- [a8s database list](a8s_database_list.md) - GET /api/v1/database-deployments
- [a8s database metrics](a8s_database_metrics.md) - GET /api/v1/database-deployments/{deploymentId}/metrics
- [a8s database restart](a8s_database_restart.md) - POST /api/v1/database-deployments/{deploymentId}/restart
- [a8s database rotate-password](a8s_database_rotate-password.md) - POST /api/v1/database-deployments/{deploymentId}/rotate-password
- [a8s database settings](a8s_database_settings.md) - Manage settings
- [a8s database update](a8s_database_update.md) - PATCH /api/v1/database-deployments/{deploymentId}
- [a8s database upgrade](a8s_database_upgrade.md) - POST /api/v1/database-deployments/{deploymentId}/upgrade-version
- [a8s database verify-password](a8s_database_verify-password.md) - POST /api/v1/database-deployments/{deploymentId}/verify-password

