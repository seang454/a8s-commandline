# a8s alert project-config

Manage project config

## Usage

```text
a8s alert project-config
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

- [a8s alert project-config get](a8s_alert_project-config_get.md) - GET /api/v1/alerts/projects/{projectId}/config
- [a8s alert project-config list](a8s_alert_project-config_list.md) - GET /api/v1/alerts/projects/configs
- [a8s alert project-config set](a8s_alert_project-config_set.md) - PUT /api/v1/alerts/projects/{projectId}/config

