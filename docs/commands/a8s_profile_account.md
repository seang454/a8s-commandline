# a8s profile account

Manage account

## Usage

```text
a8s profile account
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

- [a8s profile account deactivate](a8s_profile_account_deactivate.md) - POST /api/v1/profile/me/deactivate
- [a8s profile account delete](a8s_profile_account_delete.md) - DELETE /api/v1/profile/me
- [a8s profile account reactivate](a8s_profile_account_reactivate.md) - POST /api/v1/profile/me/reactivate
- [a8s profile account status](a8s_profile_account_status.md) - GET /api/v1/profile/me/account-status

