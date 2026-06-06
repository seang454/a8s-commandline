# A8S CLI Architecture

## Purpose

This document defines the production architecture for the Go/Cobra A8S CLI. The CLI must provide consistent access to all CLI-eligible backend endpoints documented in `backend-api-cli-catalog.md`, while keeping authentication, output, errors, retries, and asynchronous workflows consistent across commands.

## Architecture Goals

- Use resource-first commands such as `a8s project list` and `a8s admin user create`.
- Keep Cobra command files thin and move behavior into reusable internal packages.
- Support JSON, YAML, table, file-download, streaming, and interactive output.
- Support Keycloak authentication, secure token storage, and multiple contexts.
- Provide predictable errors and exit codes for humans, scripts, and CI/CD.
- Avoid exposing internal callbacks, provider webhooks, and Jenkins callbacks.
- Make commands testable without a live backend.

## Recommended Project Structure

```text
a8s-commandline/
|-- cmd/
|   |-- root.go
|   |-- auth.go
|   |-- context.go
|   |-- workspace.go
|   |-- project.go
|   |-- microservice.go
|   |-- database.go
|   |-- cluster.go
|   |-- backup.go
|   |-- kubernetes.go
|   |-- logs.go
|   |-- git.go
|   |-- scan.go
|   |-- monitoring.go
|   |-- benchmark.go
|   |-- sonarqube.go
|   |-- defectdojo.go
|   |-- alert.go
|   |-- notification.go
|   |-- admin.go
|   |-- doctor.go
|   |-- completion.go
|   `-- version.go
|-- internal/
|   |-- api/
|   |   |-- client.go
|   |   |-- request.go
|   |   |-- response.go
|   |   `-- services/
|   |-- auth/
|   |-- config/
|   |-- context/
|   |-- credentials/
|   |-- errors/
|   |-- output/
|   |-- confirm/
|   |-- workflow/
|   |-- stream/
|   |-- manifest/
|   `-- testutil/
|-- pkg/version/
|-- docs/
|-- scripts/
`-- main.go
```

## Package Responsibilities

| Package | Responsibility |
|---|---|
| `cmd` | Define Cobra commands, arguments, flags, help text, and dependency wiring. |
| `internal/api` | Execute HTTP requests, attach authentication, decode responses, normalize errors, and apply safe retries. |
| `internal/api/services` | Typed clients grouped by backend resource, such as projects, clusters, and admin users. |
| `internal/auth` | Login, logout, token refresh, identity inspection, and role discovery. |
| `internal/context` | Create, select, update, list, and delete named CLI contexts. |
| `internal/credentials` | Store tokens using the operating-system credential manager with a restricted-file fallback. |
| `internal/config` | Load configuration and resolve precedence without storing secrets directly. |
| `internal/errors` | Normalize API and local errors and map them to documented exit codes. |
| `internal/output` | Render table, JSON, YAML, raw text, and downloaded files. |
| `internal/confirm` | Confirm destructive operations and implement `--yes`. |
| `internal/workflow` | Poll, wait, and coordinate multi-endpoint operations. |
| `internal/stream` | Handle SSE, WebSocket, and log streams with cancellation and reconnect logic. |
| `internal/manifest` | Read, strictly decode, merge, and validate YAML/JSON operation input and stdin for every payload-bearing mutation command. |
| `internal/testutil` | Provide mock servers, fixtures, fake credentials, and output capture. |

## Dependency Direction

Dependencies should flow inward:

```text
main -> cmd -> workflow/services -> api/auth/context/output/errors
```

Rules:

- `cmd` may call services and workflows, but services must never import `cmd`.
- API services return typed data and typed errors; they do not print.
- Output packages print data but do not call the API.
- Workflows combine services, waiting, and streams but do not own global configuration.
- Context and credential packages remain independent of Cobra.

## Root Runtime

Create one runtime object after configuration is resolved:

```go
type Runtime struct {
    Context     context.Context
    Config      config.Resolved
    Credentials credentials.Store
    Auth        *auth.Manager
    API         *api.Client
    Output      output.Printer
    Confirm     confirm.Prompter
}
```

Cobra commands should receive this runtime through constructors instead of reading global Viper variables directly. This makes commands deterministic and testable.

## Shared API Client Contract

The current client only supports empty-body requests and returns status-only errors. Replace it with a shared client that supports:

- `context.Context` cancellation and timeouts
- JSON request and response bodies
- multipart uploads
- query parameters
- streaming responses and file downloads
- request IDs and user-agent headers
- authentication token refresh
- structured backend error decoding
- bounded retries for safe operations
- `Retry-After` handling for `429` and `503`

Recommended interface:

```go
type Client interface {
    Do(ctx context.Context, method, path string, requestBody, responseBody any, options ...RequestOption) error
    Download(ctx context.Context, path, outputPath string, options ...RequestOption) error
    Stream(ctx context.Context, method, path string, options ...RequestOption) (io.ReadCloser, error)
}
```

Do not retry unsafe create, deploy, payment, restore, or delete requests unless the backend supports an idempotency key.

## Commands and Services

Each resource group should have:

1. A Cobra command constructor in `cmd`.
2. A typed service interface under `internal/api/services`.
3. Request and response models owned by the service package.
4. Focused unit tests for arguments, flags, and service calls.

Example:

```text
cmd/cluster.go
internal/api/services/clusters.go
internal/workflow/cluster_create.go
```

`a8s cluster create --file cluster.yaml --wait` should call the cluster service, then delegate waiting or streaming to a workflow.

## Output Contract

All read commands should support:

- `table`: default human-readable output
- `json`: stable machine-readable output
- `yaml`: stable machine-readable output
- `raw`: only where the endpoint naturally returns text
- `--output-file`: required for binary downloads

Rules:

- Data goes to stdout.
- Progress, warnings, and diagnostics go to stderr.
- Machine-readable output must not contain spinners, color codes, or extra prose.
- Commands should return typed values to a common printer instead of implementing custom JSON output repeatedly.

## Context and Configuration Resolution

Recommended precedence:

1. Explicit command flags
2. Named context selected with `--context`
3. Active context
4. Environment variables
5. Built-in defaults

Configuration metadata belongs in `~/.a8s/config.yaml`. Access and refresh tokens belong in a credential store, referenced by context name.

## Asynchronous Workflows

Deployments, scans, backups, restores, payments, and some administrative actions are asynchronous. Workflow implementations must define:

- terminal success states
- terminal failure states
- polling interval
- maximum timeout
- preferred WebSocket or stream endpoint
- fallback polling endpoint
- Ctrl+C cancellation behavior

Every waiting command must work without `--wait`; without it, return the accepted operation or created resource immediately.

## Streaming Architecture

Use one stream package for:

- Kubernetes logs
- Jenkins logs
- monitoring WebSocket
- notification WebSocket
- admin event WebSocket
- cluster deployment stream

The backend WebSocket interceptor currently accepts JWTs through a `token` query parameter. Treat this as sensitive: never log the complete URL and prefer a header-based or short-lived WebSocket token backend design in the future.

## Production Safety

- Require `--yes` for destructive commands.
- Do not expose internal, webhook, or callback endpoints.
- Redact secrets from verbose output and errors.
- Validate TLS by default.
- Keep admin commands under `a8s admin`.
- Enforce authorization in the backend, not only in the CLI.
- Add `a8s doctor` checks for configuration, auth, backend health, workspace readiness, and optional cluster access.

## Backend Risks to Address

The current backend security configuration marks several powerful route groups as `permitAll`, including cluster, Kubernetes, Git integration, internal routes, and WebSockets. The CLI must not treat that as intentional authorization. Before production release:

- require authentication and ownership checks for cluster and Kubernetes routes
- protect `/api/internal/**` with service authentication
- review Git integration authorization
- protect WebSocket handshakes consistently
- ensure `/api/admin/documentation/**` requires an admin role

## Migration Plan

### Phase 1: Foundation

- Replace global Viper access with a runtime and resolved configuration.
- Implement contexts and credential storage.
- Implement the shared API client, normalized errors, output printer, and confirmation prompts.
- Add `auth`, `context`, `doctor`, `completion`, and improved `version`.

### Phase 2: Core Commands

- Implement workspace, profile, project, microservice, database, cluster, and backup services.
- Replace legacy action-first commands with resource-first commands.
- Keep temporary aliases only if users already depend on them.

### Phase 3: Operations and Streams

- Implement Kubernetes, logs, monitoring, notifications, Git, scanner, and workflows.
- Add polling, WebSocket, SSE, and download support.

### Phase 4: Quality and Administration

- Implement benchmark, SonarQube, DefectDojo, alerts, and all admin commands.
- Add audit-aware confirmation text for admin mutations.

### Phase 5: Release Hardening

- Add contract and integration tests.
- Add signed cross-platform releases, checksums, completions, and compatibility checks.

## Definition of Done

The CLI architecture is production-ready when:

- every eligible endpoint is reachable through a typed service and command
- every command uses common auth, errors, output, and context behavior
- destructive commands require confirmation
- streams and waits cancel cleanly
- tokens are securely stored and refreshed
- all commands have unit or contract tests
- critical workflows have authenticated integration tests
- no internal callback or webhook endpoint is exposed
