# a8s profile

Manage profile

## Usage

```text
a8s profile
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

- [a8s profile account](a8s_profile_account.md) - Manage account
- [a8s profile avatar](a8s_profile_avatar.md) - Manage avatar
- [a8s profile get](a8s_profile_get.md) - GET /api/v1/profile/me
- [a8s profile update](a8s_profile_update.md) - PATCH /api/v1/profile/me

