# a8s alert channel

Manage channel

## Usage

```text
a8s alert channel
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

- [a8s alert channel create](a8s_alert_channel_create.md) - POST /api/v1/alerts/channels
- [a8s alert channel delete](a8s_alert_channel_delete.md) - DELETE /api/v1/alerts/channels/{channelId}
- [a8s alert channel list](a8s_alert_channel_list.md) - GET /api/v1/alerts/channels
- [a8s alert channel update](a8s_alert_channel_update.md) - PUT /api/v1/alerts/channels/{channelId}

