# a8s project webhook

Manage webhook

## Usage

```text
a8s project webhook
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

- [a8s project webhook create](a8s_project_webhook_create.md) - POST /api/v1/projects/{projectId}/webhook
- [a8s project webhook delete](a8s_project_webhook_delete.md) - DELETE /api/v1/projects/{projectId}/webhook
- [a8s project webhook get](a8s_project_webhook_get.md) - GET /api/v1/projects/{projectId}/webhook
- [a8s project webhook rotate](a8s_project_webhook_rotate.md) - POST /api/v1/projects/{projectId}/webhook/rotate

