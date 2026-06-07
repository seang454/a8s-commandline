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

## What Each Field Does

| Field | Purpose |
|---|---|
| `apiVersion` | Identifies the configuration schema version so future CLI releases can migrate it safely. |
| `kind` | Identifies this YAML document as an A8S CLI configuration. |
| `currentContext` | Selects the context used when the user does not pass `--context`. |
| `preferences.output` | Default command output: `table`, `json`, or `yaml`. |
| `preferences.color` | Controls terminal colors: `auto`, `always`, or `never`. |
| `preferences.timeout` | Maximum duration for a complete command or waiting workflow. |
| `preferences.requestTimeout` | Maximum duration for one HTTP request. |
| `preferences.pollingInterval` | Delay between status checks for asynchronous workflows. |
| `contexts` | Contains named A8S environments such as development, staging, and production. |
| `contexts.<name>.server` | Base URL of the A8S Spring Boot backend. |
| `contexts.<name>.namespace` | Default workspace or Kubernetes namespace used by namespace-based endpoints. |
| `contexts.<name>.targetCluster` | Default managed Kubernetes cluster alias when a command supports cluster selection. |
| `contexts.<name>.tls.insecureSkipVerify` | Disables TLS certificate verification. Keep `false` outside temporary local development. |
| `contexts.<name>.tls.caFile` | Optional custom CA certificate used to verify the backend HTTPS certificate. |
| `contexts.<name>.auth.issuer` | Keycloak realm issuer used for login, discovery, and token refresh. |
| `contexts.<name>.auth.clientId` | Public Keycloak client ID registered for the CLI. |
| `contexts.<name>.auth.credentialKey` | Reference to tokens stored in the operating-system credential manager. |

The configuration file stores non-secret metadata. It must not contain access tokens, refresh tokens, passwords, or database credentials.

## How Commands Use the Active Context

When a user runs:

```bash
a8s kubernetes pods
```

the CLI:

1. reads `currentContext`
2. loads that context's `server`, `namespace`, TLS, and authentication settings
3. loads the token referenced by `credentialKey`
4. applies `requestTimeout`
5. sends the authenticated request
6. prints the response using `preferences.output`

For the development example, the resulting request is:

```http
GET http://localhost:8080/api/kubernetes/namespaces/ns-local/pods
Authorization: Bearer <token-from-credential-store>
```

Another example:

```bash
a8s workspace quota pricing
```

uses `server` and authentication but does not require the configured namespace:

```http
GET http://localhost:8080/api/v1/workspaces/quota-pricing
Authorization: Bearer <token-from-credential-store>
```

## Development, Staging, and Production

One user can configure multiple environments in the same file:

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: Config

currentContext: development

preferences:
  output: table
  color: auto
  timeout: 30m
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

  staging:
    server: https://staging-api.a8s.example.com
    namespace: ns-user-staging
    targetCluster: staging-primary
    tls:
      insecureSkipVerify: false
      caFile: ""
    auth:
      issuer: https://staging-keycloak.a8s.example.com/realms/a8s
      clientId: a8s-cli
      credentialKey: context:staging

  production:
    server: https://api.a8s.example.com
    namespace: ns-user-production
    targetCluster: production-primary
    tls:
      insecureSkipVerify: false
      caFile: ""
    auth:
      issuer: https://keycloak.autonomous-istad.com/realms/a8s
      clientId: a8s-cli
      credentialKey: context:production
```

Recommended environment rules:

| Environment | Recommended behavior |
|---|---|
| Development | May use localhost; keep TLS verification enabled when HTTPS is used. |
| Staging | Use HTTPS and production-like authentication, permissions, and workflows. |
| Production | Require HTTPS, TLS verification, confirmation for destructive commands, and securely stored credentials. |

Shared production defaults such as the backend URL, Keycloak issuer, and CLI client ID may be built into the CLI or distributed by administrators. User-specific values such as namespace and credentials should be discovered or created after login.

## First-Use Configuration

Users should not normally create the YAML manually. The CLI should create and update it automatically.

Recommended first-use flow:

```bash
a8s auth login
```

The CLI should:

1. determine the operating-system config path
2. create an initial `default` or `production` context
3. use built-in or supplied backend and Keycloak defaults
4. authenticate the user
5. store tokens in the operating-system credential manager
6. call workspace bootstrap/status and entitlements endpoints
7. discover and save the user's namespace when available
8. save the non-secret context metadata

If production defaults are not built into the CLI, initialize explicitly:

```bash
a8s context create production \
  --server https://api.a8s.example.com \
  --issuer https://keycloak.autonomous-istad.com/realms/a8s \
  --client-id a8s-cli

a8s context use production
a8s auth login
```

After login, the CLI should fill user-specific values such as `namespace` and `credentialKey`.

If Keycloak shows `Invalid parameter: redirect_uri`, the configured
`auth.clientId` does not allow the CLI callback URL. The frontend client uses a
web callback such as `https://autonomous-istad.com/api/auth/callback/keycloak`,
but the CLI uses a loopback callback such as
`http://127.0.0.1:<port>/callback`. Use a dedicated public `a8s-cli` client with
`http://127.0.0.1:*` as an allowed redirect URI, or run:

```bash
a8s auth login --callback-port 64239
```

and allow exactly:

```text
http://127.0.0.1:64239/callback
```

## Context Commands

```bash
a8s context create development \
  --server http://localhost:8080 \
  --issuer https://keycloak.autonomous-istad.com/realms/a8s \
  --client-id a8s-cli

a8s context create staging \
  --server https://staging-api.a8s.example.com \
  --issuer https://staging-keycloak.a8s.example.com/realms/a8s \
  --client-id a8s-cli

a8s context create production \
  --server https://api.a8s.example.com \
  --issuer https://keycloak.autonomous-istad.com/realms/a8s \
  --client-id a8s-cli

a8s context list
a8s context get production
a8s context use staging
a8s context update production --namespace ns-team-production
a8s context update production --target-cluster production-primary
a8s context rename staging pre-production
a8s context delete pre-production --yes
```

Context output must never include stored tokens.

## Switching and Overriding Environments

Change the active context:

```bash
a8s context use production
a8s project list
```

Use another context for only one command:

```bash
a8s project list --context staging
```

Temporarily override one context value:

```bash
a8s kubernetes pods --context production --namespace ns-other
```

These overrides do not modify the configuration file.

Before destructive production operations, the CLI should clearly display the selected context and require confirmation:

```text
Context: production
Server:  https://api.a8s.example.com
Action:  delete cluster cluster-123

Continue? [y/N]
```

## Manual Configuration

Manual YAML editing should be supported for administrators and advanced users, but it is not the normal setup path.

After manual changes, validate the file:

```bash
a8s config validate
a8s context list
a8s auth status --context production
```

Do not manually place tokens in the YAML. Use `a8s auth login --context <name>` to create the credential-store record.

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
