# a8s admin cluster

Manage cluster

## Usage

```text
a8s admin cluster
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

- [a8s admin cluster health](a8s_admin_cluster_health.md) - GET /api/v1/admin/clusters/kubernetes/{alias}/health
- [a8s admin cluster list](a8s_admin_cluster_list.md) - GET /api/v1/admin/clusters
- [a8s admin cluster nodes](a8s_admin_cluster_nodes.md) - GET /api/v1/admin/clusters/kubernetes
- [a8s admin cluster quota](a8s_admin_cluster_quota.md) - Manage quota
- [a8s admin cluster update](a8s_admin_cluster_update.md) - PATCH /api/v1/admin/clusters/{clusterId}

