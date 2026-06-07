# a8s manifest init

Generate a starter manifest for a kind

## Usage

```text
a8s manifest init <kind> [flags]
```

## Flags

- `--output-file` `string` - write the starter manifest to a file instead of stdout
- `--overwrite` `bool` - replace output file if it already exists

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

