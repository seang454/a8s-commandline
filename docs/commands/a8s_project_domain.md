# a8s project domain

Manage domain

## Usage

```text
a8s project domain
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

- [a8s project domain set](a8s_project_domain_set.md) - PATCH /api/v1/projects/{projectId}/domain
- [a8s project domain sync](a8s_project_domain_sync.md) - POST /api/v1/projects/{projectId}/domain/sync

