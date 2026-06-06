# A8S CLI Configuration

## Purpose

This document defines the CLI configuration schema, named contexts, environment variables, precedence, validation, secret handling, and migration from the current `.a8s.yaml`.

## Configuration Locations

Default locations:

| Platform | Config path |
|---|---|
| Windows | `%APPDATA%\a8s\config.yaml` |
| macOS | `~/Library/Application Support/a8s/config.yaml` |
| Linux | `${XDG_CONFIG_HOME:-~/.config}/a8s/config.yaml` |

Support `--config <path>` and `A8S_CONFIG` for explicit overrides.

## Versioned Schema

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: Config

currentContext: development

preferences:
  output: table
  color: auto
  timeout: 30s
  requestTimeout: 20s
  pollingInterval: 3s

contexts:
  development:
    server: http://localhost:8080
    namespace: ns-local
    targetCluster: local
    tls:
      insecureSkipVerify: false
      caFile: ""
    auth:
      issuer: https://keycloak.autonomous-istad.com/realms/a8s
      clientId: a8s-cli
      credentialKey: context:development
```

`apiVersion` enables future migrations. Reject unsupported major schema versions with a useful error.

## Context Commands

```bash
a8s context create development --server http://localhost:8080
a8s context list
a8s context get development
a8s context use development
a8s context update development --namespace ns-team-a
a8s context rename development local
a8s context delete local --yes
```

Context output must never include stored tokens.

## Resolution Precedence

Highest priority wins:

1. command flags
2. `--context <name>`
3. active context
4. environment variables
5. configuration preferences
6. built-in defaults

Authentication precedence:

1. `--token`
2. `A8S_TOKEN`
3. credential store for selected context
4. legacy `api_token`, temporarily during migration

## Global Environment Variables

| Variable | Purpose |
|---|---|
| `A8S_CONFIG` | Explicit config path. |
| `A8S_CONTEXT` | Select a named context. |
| `A8S_SERVER` | Override backend URL. |
| `A8S_TOKEN` | Ephemeral bearer token; never persisted. |
| `A8S_NAMESPACE` | Override namespace. |
| `A8S_TARGET_CLUSTER` | Override cluster alias. |
| `A8S_OUTPUT` | Default output: `table`, `json`, or `yaml`. |
| `A8S_TIMEOUT` | Complete command timeout. |
| `A8S_REQUEST_TIMEOUT` | Individual HTTP request timeout. |
| `A8S_COLOR` | `auto`, `always`, or `never`. |
| `NO_COLOR` | Disable colored output when present. |

Viper must bind environment names explicitly or use a replacer so nested keys behave predictably.

## Global Flags

```text
--config
--context
--server
--token
--namespace
--target-cluster
--output, -o
--timeout
--request-timeout
--verbose
--no-color
```

Use `--token` mainly for debugging because shell history may retain it.

## Validation Rules

- `server` must be an absolute HTTP or HTTPS URL.
- HTTPS is required for non-local production contexts.
- context names use lowercase letters, digits, `-`, and `_`.
- durations must be positive.
- `output` must be `table`, `json`, or `yaml`.
- `insecureSkipVerify` requires an explicit warning.
- `credentialKey` must reference a credential-store record, not contain a secret.
- unknown fields should warn initially and become errors after schema stabilization.

## Writes and Concurrency

Configuration writes must:

- use a temporary file and atomic rename
- preserve restrictive permissions
- lock or detect concurrent modification
- create a backup before schema migration
- avoid rewriting the file for read-only commands

## Legacy Migration

The current configuration is:

```yaml
api_url: http://localhost:8080
api_token: your-api-token-here
```

Migration behavior:

1. Detect legacy keys.
2. Create a `default` context using `api_url`.
3. Move a real `api_token` into the credential store.
4. Remove the token from the new YAML.
5. Back up the old file.
6. Print a migration summary to stderr.

Example:

```text
Migrated legacy configuration to context "default".
Stored credentials in the operating-system credential manager.
Backup: ~/.a8s.yaml.bak
```

Never migrate the example placeholder token as a real credential.

## Config Inspection

```bash
a8s config view
a8s config path
a8s config validate
```

`config view` must redact secrets and show resolved values only with `--resolved`.

## Acceptance Criteria

- contexts resolve according to documented precedence
- credentials never appear in config output
- legacy configuration migrates without losing server settings
- malformed configuration returns exit code `2`
- concurrent writes do not corrupt configuration
- Windows, macOS, and Linux default paths work

