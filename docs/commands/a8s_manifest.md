# a8s manifest

Generate and validate operation manifests

## Usage

```text
a8s manifest
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

- [a8s manifest init](a8s_manifest_init.md) - Generate a starter manifest for a kind
- [a8s manifest kinds](a8s_manifest_kinds.md) - List supported operation manifest kinds
- [a8s manifest schema](a8s_manifest_schema.md) - Show the manifest schema summary for a kind
- [a8s manifest validate](a8s_manifest_validate.md) - Validate an operation manifest without sending a backend request

