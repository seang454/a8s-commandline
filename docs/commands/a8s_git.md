# a8s git

Manage git

## Usage

```text
a8s git
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

- [a8s git account](a8s_git_account.md) - GET /api/v1/git-integrations/{provider}/brokered-account
- [a8s git connect](a8s_git_connect.md) - POST /api/v1/git-integrations/{provider}/connect
- [a8s git disconnect](a8s_git_disconnect.md) - DELETE /api/v1/git-integrations/{provider}
- [a8s git providers](a8s_git_providers.md) - GET /api/v1/git-integrations/linked-providers
- [a8s git repos](a8s_git_repos.md) - GET /api/v1/git-integrations/{provider}/repos
- [a8s git state](a8s_git_state.md) - GET /api/v1/git-integrations/{provider}/state
- [a8s git sync-token](a8s_git_sync-token.md) - POST /api/v1/git-integrations/{provider}/sync-keycloak-token

