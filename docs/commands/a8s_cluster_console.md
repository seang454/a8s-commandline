# a8s cluster console

Manage console

## Usage

```text
a8s cluster console
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

- [a8s cluster console credentials](a8s_cluster_console_credentials.md) - GET /api/namespaces/{namespace}/clusters/{id}/console/credentials
- [a8s cluster console data](a8s_cluster_console_data.md) - GET /api/namespaces/{namespace}/clusters/{id}/console/data
- [a8s cluster console deployment](a8s_cluster_console_deployment.md) - GET /api/namespaces/{namespace}/clusters/{id}/console/deployment
- [a8s cluster console namespaces](a8s_cluster_console_namespaces.md) - GET /api/namespaces/{namespace}/clusters/{id}/console/namespaces
- [a8s cluster console objects](a8s_cluster_console_objects.md) - GET /api/namespaces/{namespace}/clusters/{id}/console/objects
- [a8s cluster console query](a8s_cluster_console_query.md) - POST /api/namespaces/{namespace}/clusters/{id}/console/query
- [a8s cluster console test](a8s_cluster_console_test.md) - POST /api/namespaces/{namespace}/clusters/{id}/console/test

