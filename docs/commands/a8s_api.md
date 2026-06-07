# a8s api

Access backend API routes directly

## Usage

```text
a8s api
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

- [a8s api catalog](a8s_api_catalog.md) - List implemented backend route mappings
- [a8s api request](a8s_api_request.md) - Send an authenticated request to any backend route

