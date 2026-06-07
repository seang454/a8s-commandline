# a8s admin quota

Manage quota

## Usage

```text
a8s admin quota
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

- [a8s admin quota approve](a8s_admin_quota_approve.md) - POST /api/v1/admin/quota-requests/{id}/approve
- [a8s admin quota list](a8s_admin_quota_list.md) - GET /api/v1/admin/quota-requests
- [a8s admin quota reject](a8s_admin_quota_reject.md) - POST /api/v1/admin/quota-requests/{id}/reject

