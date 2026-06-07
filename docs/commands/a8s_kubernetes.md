# a8s kubernetes

Manage kubernetes

## Usage

```text
a8s kubernetes
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

- [a8s kubernetes database-resources](a8s_kubernetes_database-resources.md) - GET /api/kubernetes/namespaces/{namespace}/database-resources
- [a8s kubernetes events](a8s_kubernetes_events.md) - GET /api/kubernetes/namespaces/{namespace}/events
- [a8s kubernetes overview](a8s_kubernetes_overview.md) - GET /api/kubernetes/namespaces/{namespace}/overview
- [a8s kubernetes pods](a8s_kubernetes_pods.md) - GET /api/kubernetes/namespaces/{namespace}/pods
- [a8s kubernetes pvc](a8s_kubernetes_pvc.md) - GET /api/kubernetes/namespaces/{namespace}/persistent-volume-claims
- [a8s kubernetes services](a8s_kubernetes_services.md) - GET /api/kubernetes/namespaces/{namespace}/services
- [a8s kubernetes test](a8s_kubernetes_test.md) - GET /api/kubernetes/test

