# a8s kubernetes events

GET /api/kubernetes/namespaces/{namespace}/events

## Usage

```text
a8s kubernetes events [flags]
```

## Flags

- `--limit` `int` - maximum number of events (default `0`)
- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)
- `--warnings-only` `bool` - show warning events only

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

- `method`: `GET`
- `endpoint`: `/api/kubernetes/namespaces/{namespace}/events`
- `controller`: `KubernetesController`

