# a8s microservice history

Manage history

## Usage

```text
a8s microservice history
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

- [a8s microservice history delete](a8s_microservice_history_delete.md) - DELETE /api/v1/projects/microservices/{projectId}/history/{snapshotId}
- [a8s microservice history list](a8s_microservice_history_list.md) - GET /api/v1/projects/microservices/{projectId}/history

