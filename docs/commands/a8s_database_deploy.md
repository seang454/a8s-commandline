# a8s database deploy

Deploy a single database using flags or an operation file

## Usage

```text
a8s database deploy [flags]
```

## Examples

```bash
a8s database deploy --file database.yaml --wait
  a8s database deploy --project-name payments --engine postgresql --database-name payments --version 16 --password-env DATABASE_PASSWORD --wait
```

## Flags

- `--database-name` `string` - initial database name
- `--deployment-mode` `string` - deployment mode
- `--dry-run` `bool` - validate and print the final request without applying it
- `--engine` `string` - database engine
- `--environment` `string` - deployment environment
- `--existing-auth-secret` `string` - existing authentication secret
- `--file` `string` - YAML or JSON operation file; use - for stdin
- `--include-ca` `bool` - include CA certificate
- `--network-policy` `bool` - enable network policy
- `--password-env` `string` - read database password from an environment variable
- `--password-stdin` `bool` - read database password from stdin
- `--project-name` `string` - A8S project name
- `--release-name` `string` - deployment release name
- `--require-ssl` `bool` - require SSL connections
- `--size-profile` `string` - size profile
- `--storage-class` `string` - Kubernetes storage class
- `--storage-size` `string` - persistent storage size
- `--tls-secret` `string` - existing TLS secret
- `--tls` `bool` - enable TLS
- `--username` `string` - initial application username
- `--version` `string` - database version
- `--wait` `bool` - wait for the deployment to reach a terminal state

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

