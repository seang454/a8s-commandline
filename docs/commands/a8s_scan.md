# a8s scan

Manage scan

## Usage

```text
a8s scan
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

- [a8s scan get](a8s_scan_get.md) - GET /api/v1/image-scanner/scans/{scanId}
- [a8s scan images](a8s_scan_images.md) - GET /api/v1/image-scanner/images
- [a8s scan list](a8s_scan_list.md) - GET /api/v1/image-scanner/scans
- [a8s scan report](a8s_scan_report.md) - GET /api/v1/image-scanner/scans/{scanId}/report
- [a8s scan start](a8s_scan_start.md) - POST /api/v1/image-scanner/scans

