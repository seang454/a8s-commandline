# a8s admin project

Manage project

## Usage

```text
a8s admin project
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

- [a8s admin project deactivate](a8s_admin_project_deactivate.md) - DELETE /api/v1/admin/projects/{projectId}
- [a8s admin project list](a8s_admin_project_list.md) - GET /api/v1/admin/projects
- [a8s admin project restore](a8s_admin_project_restore.md) - POST /api/v1/admin/projects/{projectId}/restore
- [a8s admin project update](a8s_admin_project_update.md) - PATCH /api/v1/admin/projects/{projectId}

