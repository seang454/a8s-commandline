# a8s admin cluster quota

Manage quota

## Usage

```text
a8s admin cluster quota
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

- [a8s admin cluster quota delete](a8s_admin_cluster_quota_delete.md) - DELETE /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}
- [a8s admin cluster quota list](a8s_admin_cluster_quota_list.md) - GET /api/v1/admin/clusters/kubernetes/{alias}/quotas
- [a8s admin cluster quota set](a8s_admin_cluster_quota_set.md) - PUT /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}

