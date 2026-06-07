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
|   `-- a8s/
|       `-- main.go
|-- internal/
|   |-- cli/
|   |   |-- root.go
|   |   |-- runtime.go
|   |   |-- global_flags.go
|   |   |-- features/             # Mirrors every Spring backend feature folder
|   |   |   |-- admin/
|   |   |   |-- alerts/
|   |   |   |-- auth/
|   |   |   |-- databasebackup/
|   |   |   |-- databaseconsole/
|   |   |   |-- dbcluster/
|   |   |   |-- documentation/
|   |   |   |-- entitlements/
|   |   |   |-- gitintegration/
|   |   |   |-- imagescanner/
|   |   |   |-- microservice/
|   |   |   |-- monitoring/
|   |   |   |-- monolithic/
|   |   |   |-- notifications/
|   |   |   |-- payments/
|   |   |   |-- profile/
|   |   |   |-- projects/
|   |   |   |-- singledb/
|   |   |   |-- sonarqube/
|   |   |   |-- testingkit/
|   |   |   `-- workspaces/
|   |   `-- commands/             # Shared cross-feature command machinery
|   |       |-- catalogcmd/
|   |       |-- contextcmd/
|   |       |-- doctorcmd/
|   |       `-- watchcmd/
|   |-- api/
|   |   |-- client.go
|   |   |-- options.go
|   |   |-- multipart.go
|   |   |-- download.go
|   |   |-- stream.go
|   |   `-- resources/
|   |       |-- projects/
|   |       |-- microservices/
|   |       |-- databases/
|   |       |-- clusters/
|   |       |-- backups/
|   |       |-- workspaces/
|   |       `-- admin/
|   |-- operation/
|   |   |-- envelope.go
|   |   |-- loader.go
|   |   |-- merge.go
|   |   |-- validate.go
|   |   |-- registry.go
|   |   |-- secrets.go
|   |   `-- kinds/
|   |       |-- project/
|   |       |-- microservice/
|   |       |-- database/
|   |       |-- cluster/
|   |       |-- workspace/
|   |       |-- quality/
|   |       `-- admin/
|   |-- workflow/
|   |   |-- wait.go
|   |   |-- poll.go
|   |   |-- deployment/
|   |   |-- backup/
|   |   |-- payment/
|   |   `-- scan/
|   |-- auth/
|   |-- config/
|   |-- contexts/
|   |-- credentials/
|   |-- clierrors/
|   |-- output/
|   |-- confirm/
|   |-- stream/
|   |-- files/
|   `-- testutil/
|-- pkg/
|   `-- version/
|-- examples/
|   |-- project/
|   |-- microservice/
|   |-- database/
|   |-- cluster/
|   |-- backup/
|   |-- scan/
|   `-- admin/
|-- docs/
|   `-- commands/              # Generated Cobra reference
|-- scripts/
|-- .github/workflows/
|-- go.mod
|-- go.sum
|-- Makefile
`-- README.md

internal/cli/features/singledb/deploy.go       Cobra arguments and flags
internal/operation/kinds/database/deploy.go    YAML model, merge, validation
internal/api/resources/databases/deploy.go     Backend request and response
internal/workflow/deployment/database.go       Wait and polling behavior

```

This structure is intentionally more modular than a flat `cmd/` directory.
With hundreds of mapped routes, each backend feature has a matching package
under `internal/cli/features`. Developers can locate a feature using the same
name they see in the Spring Boot monolith.

### Command Package Shape

Each command group should follow a consistent structure:

```text
internal/cli/features/singledb/
|-- command.go                 # Creates `a8s database`
|-- deploy.go                  # Cobra wiring for deploy
|-- routes_gen.go              # Generated endpoints and CLI paths owned by singledb
|-- update.go
|-- upgrade.go
|-- backup.go
|-- flags.go                   # Shared database flags
|-- examples.go                # Cobra help examples
`-- command_test.go
```

Cobra files only:

- declare arguments and flags
- load operation input
- call a typed resource service or workflow
- print the returned value

They must not construct HTTP requests, implement polling loops, or contain
backend-specific normalization logic.

### API Resource Package Shape

Each backend resource package owns its transport request and response models:

```text
internal/api/resources/databases/
|-- client.go
|-- models.go
|-- deploy.go
|-- backup.go
|-- console.go
`-- client_test.go
```

Resource clients know backend paths and DTOs, but do not know Cobra, terminal
output, or operation-file envelopes.

### Operation Kind Package Shape

Each manifest domain owns user-facing operation kinds, defaults, validation,
and mapping to backend DTOs:

```text
internal/operation/kinds/database/
|-- deploy.go
|-- update.go
|-- upgrade.go
|-- backup.go
|-- defaults.go
|-- validate.go
|-- map_backend.go
`-- operation_test.go
```

This is the core of the YAML-and-flags design:

```text
YAML or explicit flags
-> operation kind
-> defaults and validation
-> backend request DTO
-> API resource client
```

Operation kinds must not perform HTTP requests or print terminal output.

## Package Responsibilities

| Package | Responsibility |
|---|---|
| `cmd/a8s` | Minimal executable entry point that constructs and executes the root command. |
| `internal/cli` | Construct the root command, runtime, global flags, and command groups. |
| `internal/cli/features/*` | Mirror backend feature folders and own feature routes, Cobra commands, friendly flags, examples, and tests. |
| `internal/cli/commands/*` | Provide shared cross-feature Cobra machinery such as generic endpoint execution, contexts, diagnostics, and watches. |
| `internal/api` | Execute HTTP requests, attach authentication, decode responses, normalize errors, and apply safe retries. |
| `internal/api/resources/*` | Typed backend clients and transport DTOs grouped by backend resource. |
| `internal/operation` | Load strict YAML/JSON, merge explicit flags, resolve secrets, validate operation kinds, and register schemas. |
| `internal/operation/kinds/*` | User-facing operation models, defaults, validation, and mapping to backend request DTOs. |
| `internal/auth` | Login, logout, token refresh, identity inspection, and role discovery. |
| `internal/contexts` | Create, select, update, list, and delete named CLI contexts. |
| `internal/credentials` | Store tokens using the operating-system credential manager with a restricted-file fallback. |
| `internal/config` | Load configuration and resolve precedence without storing secrets directly. |
| `internal/clierrors` | Normalize API and local errors and map them to documented exit codes. |
| `internal/output` | Render table, JSON, YAML, raw text, and downloaded files. |
| `internal/confirm` | Confirm destructive operations and implement `--yes`. |
| `internal/workflow` | Poll, wait, and coordinate multi-endpoint operations. |
| `internal/stream` | Handle SSE, WebSocket, and log streams with cancellation and reconnect logic. |
| `internal/files` | Read domain-content files such as source archives, dotenv files, avatars, documentation, and query files. |
| `internal/testutil` | Provide mock servers, fixtures, fake credentials, and output capture. |

## Dependency Direction

Dependencies should flow inward:

```text
cmd/a8s -> cli/commands -> operation/workflow/api resources
                          -> auth/contexts/output/clierrors
```

Rules:

- CLI command packages may call operation kinds, resource clients, and
  workflows, but those packages must never import CLI command packages.
- API resource clients return typed data and typed errors; they do not print.
- Output packages print data but do not call the API.
- Workflows combine resource clients, waiting, and streams but do not own global configuration.
- Context and credential packages remain independent of Cobra.
- Operation kinds may map to API resource DTOs but do not execute HTTP calls.
- API resource packages must not apply user-facing defaults; defaults belong to
  operation kinds.

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

Use plural or specific package names such as `contexts` and `clierrors` to
avoid confusing them with Go's standard-library `context` and `errors`
packages.

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

1. Cobra command constructors under `internal/cli/commands/<resource>`.
2. Typed backend clients under `internal/api/resources/<resource>`.
3. User-facing YAML/flag kinds under `internal/operation/kinds/<resource>` for
   configurable mutations.
4. Workflows under `internal/workflow/<resource>` for asynchronous or
   multi-endpoint operations.
5. Focused unit tests for command wiring, operation mapping, and API calls.

Example:

```text
internal/cli/commands/cluster/deploy.go
internal/operation/kinds/cluster/deploy.go
internal/api/resources/clusters/deploy.go
internal/workflow/deployment/cluster.go
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

- Implement workspace, profile, project, microservice, database, cluster, and backup resource clients and operation kinds.
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
