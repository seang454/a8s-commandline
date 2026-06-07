# Backend Feature Packages

This directory mirrors the 21 folders under the Spring Boot backend's
`features/` package. It is the first place developers should look when adding
or changing feature-specific CLI behavior.

Each feature package owns:

- registration of the backend routes belonging to that feature
- generated `routes_gen.go` inventory showing its endpoints and CLI paths
- feature-specific Cobra commands and friendly flags
- feature-specific command tests
- links to typed operation, API-resource, and workflow packages when needed

Shared HTTP execution, output, multipart uploads, downloads, and request-body
handling remain in `internal/cli/commands/catalogcmd`.

Some backend folders do not own standalone controllers:

- `databaseconsole` routes are registered by `singledb` and `dbcluster`
- `payments` routes are exposed through `workspaces`

The complete feature inventory and registration order live in `register.go`.
