# a8s auth login

Authenticate through Keycloak using browser PKCE

## Usage

```text
a8s auth login [flags]
```

## Flags

- `--callback-port` `int` - fixed local callback port; Keycloak redirect URI must allow http://127.0.0.1:<port>/callback (default `0`)
- `--login-timeout` `duration` - maximum time to complete browser authentication (default `5m0s`)
- `--no-browser` `bool` - print the login URL without opening a browser

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

