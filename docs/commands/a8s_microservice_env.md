# a8s microservice env

Manage env

## Usage

```text
a8s microservice env
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

- [a8s microservice env clear](a8s_microservice_env_clear.md) - DELETE /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment
- [a8s microservice env get](a8s_microservice_env_get.md) - GET /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment
- [a8s microservice env import](a8s_microservice_env_import.md) - POST /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment/import
- [a8s microservice env set](a8s_microservice_env_set.md) - PUT /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment

