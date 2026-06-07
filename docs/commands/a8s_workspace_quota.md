# a8s workspace quota

Manage quota

## Usage

```text
a8s workspace quota
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

- [a8s workspace quota payment-status](a8s_workspace_quota_payment-status.md) - GET /api/v1/workspaces/quota-requests/payment-status
- [a8s workspace quota pricing](a8s_workspace_quota_pricing.md) - GET /api/v1/workspaces/quota-pricing
- [a8s workspace quota purchase](a8s_workspace_quota_purchase.md) - Purchase a workspace quota plan using Bakong KHQR
- [a8s workspace quota request](a8s_workspace_quota_request.md) - POST /api/v1/workspaces/quota-requests

