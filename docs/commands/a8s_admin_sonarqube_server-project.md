# a8s admin sonarqube server-project

Manage server project

## Usage

```text
a8s admin sonarqube server-project
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

- [a8s admin sonarqube server-project create](a8s_admin_sonarqube_server-project_create.md) - POST /api/v1/admin/sonarqube/server-projects
- [a8s admin sonarqube server-project delete](a8s_admin_sonarqube_server-project_delete.md) - DELETE /api/v1/admin/sonarqube/server-projects/{projectKey}
- [a8s admin sonarqube server-project get](a8s_admin_sonarqube_server-project_get.md) - GET /api/v1/admin/sonarqube/server-projects/{projectKey}
- [a8s admin sonarqube server-project list](a8s_admin_sonarqube_server-project_list.md) - GET /api/v1/admin/sonarqube/server-projects
- [a8s admin sonarqube server-project update](a8s_admin_sonarqube_server-project_update.md) - PATCH /api/v1/admin/sonarqube/server-projects/{projectKey}

