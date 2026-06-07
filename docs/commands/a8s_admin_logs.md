# a8s admin logs

Manage logs

## Usage

```text
a8s admin logs
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

- [a8s admin logs clusters](a8s_admin_logs_clusters.md) - GET /api/v1/admin/logs/clusters
- [a8s admin logs namespaces](a8s_admin_logs_namespaces.md) - GET /api/v1/admin/logs/namespaces
- [a8s admin logs pods](a8s_admin_logs_pods.md) - GET /api/v1/admin/logs/pods
- [a8s admin logs query](a8s_admin_logs_query.md) - GET /api/v1/admin/logs/query
- [a8s admin logs workloads](a8s_admin_logs_workloads.md) - GET /api/v1/admin/logs/workloads

