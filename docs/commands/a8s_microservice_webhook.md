# a8s microservice webhook

Manage webhook

## Usage

```text
a8s microservice webhook
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

- [a8s microservice webhook get](a8s_microservice_webhook_get.md) - GET /api/v1/projects/microservices/{projectId}/webhook
- [a8s microservice webhook update](a8s_microservice_webhook_update.md) - POST /api/v1/projects/microservices/{projectId}/webhook

