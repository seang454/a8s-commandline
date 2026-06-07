# a8s microservice

Manage microservice

## Usage

```text
a8s microservice
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

- [a8s microservice apply](a8s_microservice_apply.md) - PUT /api/v1/projects/microservices/{projectId}/canvas
- [a8s microservice delete](a8s_microservice_delete.md) - DELETE /api/v1/projects/microservices/{projectId}
- [a8s microservice deploy](a8s_microservice_deploy.md) - POST /api/v1/projects/microservices
- [a8s microservice detect](a8s_microservice_detect.md) - POST /api/v1/projects/microservices/detect
- [a8s microservice domains](a8s_microservice_domains.md) - Manage domains
- [a8s microservice env](a8s_microservice_env.md) - Manage env
- [a8s microservice get](a8s_microservice_get.md) - GET /api/v1/projects/microservices/{projectId}
- [a8s microservice history](a8s_microservice_history.md) - Manage history
- [a8s microservice pods](a8s_microservice_pods.md) - GET /api/v1/projects/microservices/{projectId}/runtime-pods
- [a8s microservice readiness](a8s_microservice_readiness.md) - GET /api/v1/projects/microservices/{projectId}/readiness
- [a8s microservice redeploy](a8s_microservice_redeploy.md) - POST /api/v1/projects/microservices/{projectId}/redeploy
- [a8s microservice rollback](a8s_microservice_rollback.md) - POST /api/v1/projects/microservices/{projectId}/rollback
- [a8s microservice webhook](a8s_microservice_webhook.md) - Manage webhook

