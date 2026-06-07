# a8s sonarqube

Manage sonarqube

## Usage

```text
a8s sonarqube
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

- [a8s sonarqube access](a8s_sonarqube_access.md) - POST /api/v1/projects/{projectId}/sonarqube/access
- [a8s sonarqube summary](a8s_sonarqube_summary.md) - GET /api/v1/projects/{projectId}/sonarqube

