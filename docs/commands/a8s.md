# a8s

A8S platform command-line interface

## Usage

```text
a8s
```

## Flags

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

- [a8s admin](a8s_admin.md) - Manage admin
- [a8s alert](a8s_alert.md) - Manage alert
- [a8s api](a8s_api.md) - Access backend API routes directly
- [a8s auth](a8s_auth.md) - Authenticate and manage the current session
- [a8s backup](a8s_backup.md) - Manage backup
- [a8s benchmark](a8s_benchmark.md) - Manage benchmark
- [a8s cluster](a8s_cluster.md) - Manage cluster
- [a8s config](a8s_config.md) - Inspect CLI configuration
- [a8s context](a8s_context.md) - Manage named backend environments
- [a8s database](a8s_database.md) - Manage single database deployments
- [a8s defectdojo](a8s_defectdojo.md) - Manage defectdojo
- [a8s doctor](a8s_doctor.md) - Check CLI configuration and backend connectivity
- [a8s features](a8s_features.md) - List backend features exposed by the CLI
- [a8s git](a8s_git.md) - Manage git
- [a8s kubernetes](a8s_kubernetes.md) - Manage kubernetes
- [a8s logs](a8s_logs.md) - GET /api/kubernetes/namespaces/{namespace}/pods/{podName}/logs/stream
- [a8s manifest](a8s_manifest.md) - Generate and validate operation manifests
- [a8s microservice](a8s_microservice.md) - Manage microservice
- [a8s monitoring](a8s_monitoring.md) - Manage monitoring
- [a8s notification](a8s_notification.md) - Manage notification
- [a8s profile](a8s_profile.md) - Manage profile
- [a8s project](a8s_project.md) - Manage project
- [a8s scan](a8s_scan.md) - Manage scan
- [a8s sonarqube](a8s_sonarqube.md) - Manage sonarqube
- [a8s version](a8s_version.md) - Print the CLI version
- [a8s workspace](a8s_workspace.md) - Manage workspace

