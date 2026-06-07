# A8S CLI Documentation

## Start Here

| Document | Purpose |
|---|---|
| [backend-api-cli-catalog.md](backend-api-cli-catalog.md) | Complete backend endpoint inventory and endpoint-to-command mapping. |
| [architecture.md](architecture.md) | Production Go/Cobra architecture, package boundaries, and migration plan. |
| [authentication.md](authentication.md) | Keycloak login, token refresh, credential storage, logout, and authorization. |
| [configuration.md](configuration.md) | Context schema, flags, environment variables, precedence, and legacy migration. |
| [operation-input.md](operation-input.md) | Production-wide YAML/JSON and equivalent-flag policy for mutation commands. |
| [error-contract.md](error-contract.md) | Normalized errors, retries, timeouts, and stable CLI exit codes. |
| [workflows.md](workflows.md) | Waiting, polling, streaming, and terminal-state behavior. |
| [testing-strategy.md](testing-strategy.md) | Unit, contract, integration, security, and release testing. |
| [command-reference.md](command-reference.md) | Command conventions and generated reference requirements. |
| [build-and-install.md](build-and-install.md) | Build outputs, per-OS installation, PATH setup, and first-run flow. |
| [openapi-compatibility.md](openapi-compatibility.md) | OpenAPI generation, API compatibility, and drift detection. |
| [release-process.md](release-process.md) | Versioning, CI gates, signing, packaging, and release procedure. |

## Current Status

| Area | Status |
|---|---|
| Backend endpoint discovery | Complete |
| Endpoint-to-command mapping | Complete |
| Architecture specification | Complete |
| Authentication specification | Complete |
| Configuration specification | Complete |
| Operation-input specification | Complete |
| Error contract | Complete |
| Workflow contract | Complete |
| Testing strategy | Complete |
| Command reference | Generated from Cobra into `docs/command-reference.md` and `docs/commands/` |
| OpenAPI compatibility design | Complete |
| Release-process design | Complete; CI/release workflows added |
| Production Go/Cobra implementation | In progress: every unique catalog command path is registered; `database deploy` is typed; selected scan/cluster/payment typed flags and `--wait` polling are implemented; Keycloak PKCE login, secure context credentials, pre-command refresh, one-time backend-401 refresh/replay, auth status, and local logout are implemented |
| Backend security hardening | Required before production |

## Production Gate

Do not declare the CLI production-ready until:

- powerful backend routes are protected consistently
- the shared CLI foundation is implemented
- critical workflows pass authenticated integration tests
- OpenAPI drift checks pass
- release binaries are signed and checksummed
