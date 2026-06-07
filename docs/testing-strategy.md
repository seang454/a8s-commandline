# A8S CLI Testing Strategy

## Goals

- prevent command and API contract regressions
- verify authentication, contexts, errors, output, waiting, and streaming centrally
- test destructive and admin commands safely
- detect backend endpoint and OpenAPI drift
- validate release binaries on supported platforms

## Test Layers

| Layer | Scope | Network |
|---|---|---|
| Unit | flags, validation, status mapping, config, output, errors | none |
| Command | Cobra command wiring with fake services | none |
| Contract | real HTTP client against mock server fixtures | local mock |
| Integration | CLI against deployed backend and Keycloak | test environment |
| End-to-end | critical user/admin workflows | isolated environment |
| Release smoke | packaged binary behavior | controlled environment |

## Unit Tests

Required focus:

- configuration precedence and migration
- context CRUD
- credential-store abstraction
- token expiry and refresh decisions
- error decoding and exit-code mapping
- output rendering and secret redaction
- manifest validation
- confirmation behavior
- workflow terminal-state mapping
- retry eligibility

Use table-driven Go tests. Avoid real sleeps by injecting clocks and backoff functions.

## Command Tests

Construct commands with fake services and in-memory stdout/stderr:

```go
func NewProjectCommand(runtime Runtime, service ProjectService) *cobra.Command
```

Test:

- required and conflicting flags
- positional argument validation
- correct service calls
- `--output` behavior
- destructive confirmation and `--yes`
- returned exit codes
- help and examples

## Mock-Server Contract Tests

Use `httptest.Server` to verify:

- request paths, methods, query parameters, and headers
- JSON and multipart bodies
- file downloads
- current backend error shapes
- empty and malformed responses
- `401` refresh-once behavior
- `429 Retry-After`
- safe retry behavior
- cancellation and timeouts

Maintain fixtures by backend feature.

## Integration Tests

Run against an isolated backend, Keycloak realm, database, and Kubernetes test namespace.

The repository includes a gated smoke-test scaffold:

```bash
A8S_RUN_INTEGRATION_TESTS=true \
A8S_SERVER=https://api.example.com \
A8S_TOKEN=<access-token> \
go test ./internal/integration
```

Required workflows:

- login, status, refresh, and logout
- workspace bootstrap and entitlements
- project deploy/get/delete
- database deploy/backup/restore/delete
- cluster deploy/status/delete
- image scan and report
- monitoring and log access
- quota purchase in payment test mode
- admin user/project/quota operations

Every test creates uniquely named resources and performs cleanup.

## Destructive Command Tests

Destructive tests must:

- use isolated test resources
- verify prompt cancellation
- verify `--yes`
- verify non-admin rejection
- verify ownership enforcement
- verify already-deleted/not-found behavior
- never target shared production resources

Require an explicit environment gate:

```text
A8S_RUN_DESTRUCTIVE_TESTS=true
```

## Security Tests

- tokens never appear in output or logs
- redirect handling never leaks authorization headers
- invalid TLS is rejected by default
- non-admin users cannot use admin endpoints
- internal callbacks are not exposed by Cobra
- shell completion does not expose secrets
- config output redacts credential references appropriately
- WebSocket token URLs are redacted

## Output Golden Tests

Use golden files for stable table, JSON, YAML, error, and help output. Normalize timestamps, IDs, and terminal width before comparison.

Machine-readable output is a compatibility contract; changes require review.

## OpenAPI Drift Tests

CI should:

1. fetch or generate the backend OpenAPI document
2. compare it with the committed snapshot
3. detect removed or changed operations
4. verify every eligible operation has a command mapping
5. fail on undocumented breaking changes

## Coverage Expectations

| Area | Target |
|---|---:|
| Shared foundation packages | 85%+ |
| Auth, errors, config, workflow state mapping | 90%+ |
| Individual thin command files | meaningful behavior coverage |
| Critical end-to-end workflows | at least one happy and one failure path |

Coverage numbers do not replace behavior tests.

## CI Matrix

Run on:

```text
windows-latest
ubuntu-latest
macos-latest
```

Required CI stages:

```text
format -> vet -> unit -> contract -> build -> OpenAPI drift -> security scan
```

Integration and destructive suites may run on protected branches or scheduled environments.

## Tooling Recommendations

- `go test ./...`
- `go test -race ./...` on supported runners
- `go vet ./...`
- `staticcheck`
- `govulncheck`
- `gofumpt` or `gofmt`
- `golangci-lint`
- mock HTTP servers using `httptest`

## Test Data Rules

- never commit real tokens or credentials
- use deterministic fake IDs and timestamps
- redact backend snapshots
- use payment test pricing and isolated Bakong behavior
- keep fixtures small and feature-specific

## Definition of Done for a Command

A command is complete when:

- help, arguments, flags, examples, and validation exist
- it uses shared auth, API, error, and output behavior
- unit and contract tests pass
- destructive behavior is confirmed where applicable
- endpoint mapping and command reference are updated
- an integration test exists for critical workflows
