# a8s context

Manage named backend environments

## Usage

```text
a8s context
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

- [a8s context create](a8s_context_create.md) - Create a named context
- [a8s context delete](a8s_context_delete.md) - Delete a named context
- [a8s context get](a8s_context_get.md) - Get a configured context
- [a8s context list](a8s_context_list.md) - List configured contexts
- [a8s context update](a8s_context_update.md) - Update a named context
- [a8s context use](a8s_context_use.md) - Set the default context

