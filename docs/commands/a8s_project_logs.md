# a8s project logs

GET /api/v1/jenkins/logs/stream

## Usage

```text
a8s project logs [flags]
```

## Flags

- `--build` `int` - Jenkins build number (default `0`)
- `--follow` `bool` - keep reading the event stream
- `--job` `string` - Jenkins job name
- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)
- `--queue-item` `int` - Jenkins queue item number (default `0`)

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

- [a8s project logs websocket](a8s_project_logs_websocket.md) - Watch Jenkins logs over WebSocket

## Backend Endpoint

- `method`: `GET`
- `endpoint`: `/api/v1/jenkins/logs/stream`
- `controller`: `JenkinsController`

