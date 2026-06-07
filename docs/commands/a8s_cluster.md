# a8s cluster

Manage cluster

## Usage

```text
a8s cluster
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

- [a8s cluster backup](a8s_cluster_backup.md) - Manage backup
- [a8s cluster certificate](a8s_cluster_certificate.md) - GET /api/namespaces/{namespace}/clusters/{id}/certificate
- [a8s cluster clone-from-backup](a8s_cluster_clone-from-backup.md) - POST /api/namespaces/{namespace}/clusters/clone-from-backup
- [a8s cluster console](a8s_cluster_console.md) - Manage console
- [a8s cluster delete](a8s_cluster_delete.md) - DELETE /api/namespaces/{namespace}/clusters/{id}
- [a8s cluster deploy](a8s_cluster_deploy.md) - POST /api/namespaces/{namespace}/cluster-deployments
- [a8s cluster deployment](a8s_cluster_deployment.md) - Manage deployment
- [a8s cluster get](a8s_cluster_get.md) - GET /api/namespaces/{namespace}/clusters/{id}
- [a8s cluster history](a8s_cluster_history.md) - GET /api/namespaces/{namespace}/clusters/{id}/deployments
- [a8s cluster list](a8s_cluster_list.md) - GET /api/namespaces/{namespace}/clusters
- [a8s cluster metrics](a8s_cluster_metrics.md) - GET /api/namespaces/{namespace}/clusters/{id}/metrics
- [a8s cluster settings](a8s_cluster_settings.md) - Manage settings
- [a8s cluster status](a8s_cluster_status.md) - GET /api/namespaces/{namespace}/cluster-deployments/{releaseName}
- [a8s cluster update](a8s_cluster_update.md) - PATCH /api/namespaces/{namespace}/clusters/{id}
- [a8s cluster upgrade](a8s_cluster_upgrade.md) - POST /api/namespaces/{namespace}/clusters/{id}/upgrade-version
- [a8s cluster values](a8s_cluster_values.md) - GET /api/namespaces/{namespace}/clusters/{id}/values
- [a8s cluster watch](a8s_cluster_watch.md) - GET /api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream

