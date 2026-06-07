# a8s admin user

Manage user

## Usage

```text
a8s admin user
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

- [a8s admin user create](a8s_admin_user_create.md) - POST /api/v1/admin/users
- [a8s admin user deactivate](a8s_admin_user_deactivate.md) - DELETE /api/v1/admin/users/{userId}
- [a8s admin user list](a8s_admin_user_list.md) - GET /api/v1/admin/users
- [a8s admin user reactivate](a8s_admin_user_reactivate.md) - POST /api/v1/admin/users/{userId}/reactivate
- [a8s admin user update](a8s_admin_user_update.md) - PATCH /api/v1/admin/users/{userId}

