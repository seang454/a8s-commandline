# a8s project

Manage project

## Usage

```text
a8s project
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

- [a8s project auto-deploy](a8s_project_auto-deploy.md) - Manage auto deploy
- [a8s project branches](a8s_project_branches.md) - GET /api/v1/projects/{projectId}/branches
- [a8s project delete](a8s_project_delete.md) - DELETE /api/v1/projects/{projectId}
- [a8s project deploy](a8s_project_deploy.md) - POST /api/v1/projects
- [a8s project domain](a8s_project_domain.md) - Manage domain
- [a8s project env](a8s_project_env.md) - Manage env
- [a8s project get](a8s_project_get.md) - GET /api/v1/projects/{projectId}
- [a8s project list](a8s_project_list.md) - GET /api/v1/projects
- [a8s project live](a8s_project_live.md) - Manage live
- [a8s project logs](a8s_project_logs.md) - GET /api/v1/jenkins/logs/stream
- [a8s project redeploy](a8s_project_redeploy.md) - POST /api/v1/projects/{projectId}/sync
- [a8s project release](a8s_project_release.md) - Manage release
- [a8s project releases](a8s_project_releases.md) - GET /api/v1/projects/{projectId}/releases
- [a8s project repository](a8s_project_repository.md) - Manage repository
- [a8s project rollback](a8s_project_rollback.md) - POST /api/v1/projects/{projectId}/rollback
- [a8s project settings](a8s_project_settings.md) - Manage settings
- [a8s project webhook](a8s_project_webhook.md) - Manage webhook

