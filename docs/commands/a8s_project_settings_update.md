# a8s project settings update

PATCH /api/v1/projects/{projectId}/settings

## Usage

```text
a8s project settings update <project-id> [flags]
```

## Flags

- `--dry-run` `bool` - print the resolved request without sending it
- `--file` `string` - YAML or JSON request body; operation envelopes use their spec
- `--form` `stringArray` - add multipart form field using key=value; repeatable (default `[]`)
- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)
- `--set` `stringArray` - set request field using dotted key=value; repeatable (default `[]`)
- `--upload` `stringArray` - upload a file using field=path; repeatable (default `[]`)

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

- `method`: `PATCH`
- `endpoint`: `/api/v1/projects/{projectId}/settings`
- `controller`: `ProjectController`

