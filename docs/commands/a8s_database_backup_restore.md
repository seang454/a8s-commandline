# a8s database backup restore

POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore

## Usage

```text
a8s database backup restore <deployment-id> <run-id> [flags]
```

## Flags

- `--dry-run` `bool` - print the resolved request without sending it
- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)
- `--wait` `bool` - wait for the asynchronous operation to finish
- `--yes` `bool` - confirm a destructive operation

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

- [a8s database backup restore cancel](a8s_database_backup_restore_cancel.md) - POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore/cancel

## Backend Endpoint

- `method`: `POST`
- `endpoint`: `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore`
- `controller`: `DatabaseDeploymentBackupController`

