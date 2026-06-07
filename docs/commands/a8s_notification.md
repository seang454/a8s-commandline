# a8s notification

Manage notification

## Usage

```text
a8s notification
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

- [a8s notification list](a8s_notification_list.md) - GET /api/notifications/history/{userId}
- [a8s notification preferences](a8s_notification_preferences.md) - Manage preferences
- [a8s notification read](a8s_notification_read.md) - POST /api/notifications/{notificationId}/read
- [a8s notification watch](a8s_notification_watch.md) - Watch notifications

