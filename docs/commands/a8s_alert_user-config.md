# a8s alert user-config

Manage user config

## Usage

```text
a8s alert user-config
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

- [a8s alert user-config get](a8s_alert_user-config_get.md) - GET /api/v1/alerts/user-config
- [a8s alert user-config set](a8s_alert_user-config_set.md) - PUT /api/v1/alerts/user-config

