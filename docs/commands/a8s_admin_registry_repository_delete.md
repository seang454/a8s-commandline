# a8s admin registry repository delete

DELETE /api/v1/admin/registry/projects/{projectName}/repositories

## Usage

```text
a8s admin registry repository delete <project-name> [flags]
```

## Flags

- `--dry-run` `bool` - print the resolved request without sending it
- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)
- `--yes` `bool` - confirm a destructive operation

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

## Backend Endpoint

- `method`: `DELETE`
- `endpoint`: `/api/v1/admin/registry/projects/{projectName}/repositories`
- `controller`: `AdminRegistryController`

