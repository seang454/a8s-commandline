# a8s auth logout

Clear stored credentials for the active context

## Usage

```text
a8s auth logout [flags]
```

## Flags

- `--callback-port` `int` - fixed local logout callback port; Keycloak post logout redirect URI must allow http://127.0.0.1:<port>/callback (default `0`)
- `--keycloak` `bool` - also end the Keycloak browser session
- `--logout-timeout` `duration` - maximum time to complete browser logout (default `2m0s`)
- `--no-browser` `bool` - print the Keycloak logout URL without opening a browser; implies --keycloak

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

