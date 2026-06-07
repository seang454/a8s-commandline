# a8s context update

Update a named context

## Usage

```text
a8s context update <name> [flags]
```

## Flags

- `--client-id` `string` - OIDC public client ID (default `a8s-cli`)
- `--file` `string` - YAML or JSON context operation file; use - for stdin
- `--issuer` `string` - OIDC issuer
- `--namespace` `string` - default namespace
- `--server` `string` - backend server URL
- `--target-cluster` `string` - default managed cluster alias

## Inherited Flags

- `--config` `string` - config file path
- `--context` `string` - named context to use
- `--request-timeout` `string` - single HTTP request timeout
- `--timeout` `string` - complete command timeout
- `--token` `string` - temporary bearer token; prefer A8S_TOKEN
- `-o, --output` `string` - output format: table|json|yaml

