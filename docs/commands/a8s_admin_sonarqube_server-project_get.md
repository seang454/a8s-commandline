# a8s admin sonarqube server-project get

GET /api/v1/admin/sonarqube/server-projects/{projectKey}

## Usage

```text
a8s admin sonarqube server-project get <project-key> [flags]
```

## Flags

- `--output-file` `string` - write the response body to a file
- `--query` `stringArray` - add query parameter using key=value; repeatable (default `[]`)

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
- `endpoint`: `/api/v1/admin/sonarqube/server-projects/{projectKey}`
- `controller`: `AdminSonarQubeController`

