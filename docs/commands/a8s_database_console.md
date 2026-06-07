# a8s database console

Manage console

## Usage

```text
a8s database console
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

- [a8s database console data](a8s_database_console_data.md) - GET /api/v1/database-deployments/{deploymentId}/console/data
- [a8s database console namespaces](a8s_database_console_namespaces.md) - GET /api/v1/database-deployments/{deploymentId}/console/namespaces
- [a8s database console objects](a8s_database_console_objects.md) - GET /api/v1/database-deployments/{deploymentId}/console/objects
- [a8s database console query](a8s_database_console_query.md) - POST /api/v1/database-deployments/{deploymentId}/console/query
- [a8s database console test](a8s_database_console_test.md) - POST /api/v1/database-deployments/{deploymentId}/console/test

