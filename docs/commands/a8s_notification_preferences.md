# a8s notification preferences

Manage preferences

## Usage

```text
a8s notification preferences
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

- [a8s notification preferences get](a8s_notification_preferences_get.md) - GET /api/notifications/preferences/{userId}
- [a8s notification preferences set](a8s_notification_preferences_set.md) - POST /api/notifications/preferences/{userId}

