# a8s admin docs

Manage docs

## Usage

```text
a8s admin docs
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

- [a8s admin docs delete](a8s_admin_docs_delete.md) - DELETE /api/admin/documentation/content
- [a8s admin docs files](a8s_admin_docs_files.md) - GET /api/admin/documentation/files
- [a8s admin docs get](a8s_admin_docs_get.md) - GET /api/admin/documentation/content
- [a8s admin docs publish-logs](a8s_admin_docs_publish-logs.md) - GET /api/admin/documentation/publish/logs
- [a8s admin docs publish](a8s_admin_docs_publish.md) - POST /api/admin/documentation/publish
- [a8s admin docs update](a8s_admin_docs_update.md) - PUT /api/admin/documentation/content

