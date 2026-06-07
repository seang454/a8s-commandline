# A8S CLI

A Go + Cobra CLI for the A8S platform.

The endpoint mapping and resource-first Cobra command tree are generated from
the backend catalog. Specialized typed workflows are added on top of that
shared endpoint executor.

## Documentation

- [Documentation index](docs/README.md)
- [Backend API to CLI catalog](docs/backend-api-cli-catalog.md)
- [Architecture](docs/architecture.md)
- [Authentication](docs/authentication.md)
- [Configuration](docs/configuration.md)
- [Error contract](docs/error-contract.md)
- [Workflow contracts](docs/workflows.md)
- [Testing strategy](docs/testing-strategy.md)
- [Command-reference design](docs/command-reference.md)
- [Build and install](docs/build-and-install.md)
- [OpenAPI compatibility](docs/openapi-compatibility.md)
- [Release process](docs/release-process.md)

## Quick Start

### 1. Configure

Use `.a8s.yaml` as an example for context configuration. Authenticate through
the configured Keycloak issuer, or supply a temporary token for automation.

```bash
export A8S_CONFIG="$PWD/.a8s.yaml"
./a8s auth login
./a8s auth status
```

For non-interactive automation:

```bash
export A8S_TOKEN="<temporary-access-token>"
```

### 2. Build
```bash
make build
```

### 3. Call a backend feature

Every unique user-facing catalog command path is registered:

```bash
./a8s project list
./a8s cluster get <cluster-id>
./a8s workspace quota pricing
./a8s admin monitoring overview
```

Inspect the backend-feature inventory or one feature's generated route file:

```bash
./a8s features
./a8s api catalog --search dbcluster
```

Feature-specific CLI ownership mirrors the Spring backend under
`internal/cli/features/<backend-feature>/`.

Mutation commands accept YAML/JSON request bodies and explicit field overrides:

```bash
./a8s project domain set <project-id> --set customDomain=api.example.com
./a8s alert channel create --file alert-channel.yaml
./a8s profile avatar upload --upload file=avatar.png
```

New backend routes can be reached before the catalog is regenerated:

```bash
./a8s api request GET /api/v1/new-resource
```

### 4. Deploy a database

Using a YAML operation:

```bash
export A8S_DATABASE_PASSWORD="<database-password>"
./a8s database deploy --file examples/database/deployment.yaml --wait
```

Using flags:

```bash
./a8s database deploy \
  --project-name payments \
  --engine postgresql \
  --database-name payments \
  --version 16 \
  --password-env A8S_DATABASE_PASSWORD \
  --wait
```

Preview the final backend request without sending it:

```bash
./a8s --output yaml database deploy \
  --file examples/database/deployment.yaml \
  --storage-size 50Gi \
  --dry-run
```

## Current Implementation

- production root runtime and global flags
- versioned named-context configuration loading
- normalized API errors and exit codes
- shared JSON/YAML output
- strict YAML operation loading with unknown-field rejection
- explicit flag-over-manifest merging
- secure secret input from environment variables or stdin
- typed `database deploy` backend client
- database deployment dry-run and wait workflow
- generated registry for all unique user-facing catalog command paths
- authenticated generic JSON/YAML mutation execution with `--file` and `--set`
- multipart uploads with `--upload` and `--form`
- raw downloads with `--output-file`
- SSE/raw response streaming and four authenticated WebSocket watch routes
- shared `--wait` polling for known async operations using status URLs, operation IDs, or built-in scan/cluster/payment poll paths
- typed convenience flags for `scan start`, `cluster deploy`, and `workspace quota purchase`
- local context create/list/get/use/update/delete commands
- Keycloak browser login using Authorization Code Flow with PKCE
- context-scoped OS keyring credentials with restricted-file fallback
- automatic access-token refresh before command execution
- one forced token refresh and exact HTTP request replay after backend `401`
- unusable refresh-credential cleanup after Keycloak `invalid_grant`
- cross-origin redirect and absolute API URL token-leak protection
- authentication status and local logout commands
- direct `api request` compatibility command
- generated Cobra command reference under `docs/command-reference.md` and `docs/commands/`
- gated authenticated backend smoke test scaffold
- GitHub Actions CI/release workflows with cross-platform builds and checksums
- local installer scripts for Windows and Unix-like systems
- catalog-to-Cobra coverage test
- focused unit and command tests

The generated commands provide broad backend endpoint coverage. High-value
operations still need deeper domain validation, more typed payload models,
workflow-specific integration tests, and backend security hardening before the
CLI should be called production-ready.

Legacy commands remain in the repository during migration but are no longer
registered by the production root command.

## Build for All Platforms
```bash
make build-all
# outputs: dist/a8s-linux-amd64, dist/a8s-darwin-arm64, dist/a8s-windows-amd64.exe
```
