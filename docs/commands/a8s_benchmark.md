# a8s benchmark

Manage benchmark

## Usage

```text
a8s benchmark
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

- [a8s benchmark delete](a8s_benchmark_delete.md) - DELETE /api/v1/projects/live/{projectId}/benchmark/runs/{runId}
- [a8s benchmark get](a8s_benchmark_get.md) - GET /api/v1/projects/live/{projectId}/benchmark/runs/{runId}
- [a8s benchmark list](a8s_benchmark_list.md) - GET /api/v1/projects/live/{projectId}/benchmark/runs
- [a8s benchmark run](a8s_benchmark_run.md) - POST /api/v1/projects/live/{projectId}/benchmark/run

