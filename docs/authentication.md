# A8S CLI Authentication

## Purpose

This document defines authentication, token storage, refresh, logout, and authorization behavior for the A8S CLI.

The backend is a stateless Spring Security OAuth2 resource server. It validates Keycloak JWT bearer tokens and maps roles from the JWT `realm_access.roles` claim to authorities such as `ROLE_USER` and `ROLE_ADMIN`.

## Required Commands

```bash
a8s auth login
a8s auth status
a8s auth logout
a8s auth verify-email status
a8s auth verify-email start
```

Authentication is scoped to a CLI context. A user may be logged into multiple A8S environments at the same time.

## Recommended Login Flow

Use OAuth 2.0 Authorization Code Flow with PKCE for an installed CLI application:

1. Discover Keycloak endpoints from the configured issuer.
2. Generate a PKCE verifier, challenge, state, and nonce.
3. Start a temporary loopback callback listener on `127.0.0.1`.
4. Open the authorization URL in the default browser.
5. Validate callback state and exchange the authorization code.
6. Validate token issuer, audience, nonce, and expiry.
7. Store tokens using the operating-system credential manager.
8. Call a lightweight authenticated backend endpoint to verify access.

If browser login is unavailable, support Keycloak Device Authorization Grant only when the configured Keycloak client enables it.

Do not use Resource Owner Password Credentials for ordinary CLI users.

## Keycloak Client Requirements

Create a dedicated public Keycloak client for the CLI:

```text
Client ID: a8s-cli
Client authentication: Off
Standard flow: On
PKCE method: S256
Valid redirect URIs: http://127.0.0.1:*
Device authorization grant: Optional
Direct access grants: Off
```

The CLI client must not contain a client secret. A secret embedded in a distributed binary is not secret.

Required scopes should include:

- `openid`
- `profile`
- `email`
- `offline_access` only if refresh-token policy permits it

## Context Authentication Configuration

Context metadata may include:

```yaml
current-context: production

contexts:
  production:
    server: https://api.example.com
    namespace: ns-team-a
    target-cluster: primary
    auth:
      issuer: https://keycloak.example.com/realms/a8s
      client-id: a8s-cli
      credential-key: context:production
```

Do not store access tokens or refresh tokens directly in this YAML.

## Credential Storage

Preferred storage:

| Platform | Storage |
|---|---|
| Windows | Windows Credential Manager |
| macOS | Keychain |
| Linux | Secret Service-compatible keyring |

If no credential manager is available, allow a restricted file fallback only after warning the user. The fallback file must:

- be separate from normal configuration
- use restrictive file permissions
- never be printed by `a8s context get`
- never be included in diagnostics unless explicitly requested and redacted

Recommended credential record:

```json
{
  "accessToken": "...",
  "refreshToken": "...",
  "idToken": "...",
  "accessTokenExpiry": "2026-06-06T12:00:00Z",
  "refreshTokenExpiry": "2026-07-06T12:00:00Z",
  "issuer": "https://keycloak.example.com/realms/a8s",
  "clientId": "a8s-cli"
}
```

## Token Refresh

Before an authenticated request:

1. Load the active context credentials.
2. If the access token remains valid beyond a small safety window, use it.
3. If it is near expiry, refresh it once.
4. Persist rotated refresh tokens.
5. Retry the original request once after successful refresh.
6. If refresh fails, clear unusable credentials and return authentication exit code `3`.

Do not repeatedly retry `401` responses. One refresh-and-retry attempt is the maximum.

## Request Authentication

HTTP requests use:

```http
Authorization: Bearer <access-token>
```

The CLI must not:

- send bearer tokens to a different host after redirects
- log authorization headers
- include tokens in ordinary query parameters
- expose tokens in shell completion or command examples

The current backend WebSocket interceptor expects `?token=<jwt>`. Until the backend supports a safer handshake:

- never print full WebSocket URLs
- avoid persistent debug logs containing the query string
- use short-lived access tokens
- clear reconnect state after logout

## Authorization and Roles

The backend reads Keycloak realm roles and converts them to Spring authorities:

```text
realm_access.roles: ["admin"] -> ROLE_ADMIN
```

Rules:

- `a8s admin ...` commands may inspect the token and warn when `ADMIN` is absent.
- The CLI must still send the request and rely on backend authorization where appropriate.
- A backend `403` is authoritative and maps to exit code `4`.
- Never add a CLI option that bypasses ownership or role checks.

## Authentication Status

`a8s auth status` should display:

```text
Context:       production
Server:        https://api.example.com
Issuer:        https://keycloak.example.com/realms/a8s
Subject:       <subject-id>
Username:      user@example.com
Roles:         USER, ADMIN
Token expires: 2026-06-06T12:00:00Z
Status:        authenticated
```

Do not display tokens. With JSON or YAML output, use stable field names.

## Logout

`a8s auth logout` should:

1. Attempt Keycloak end-session or token revocation when supported.
2. Delete the context credential record even if remote logout fails.
3. Clear cached identity and WebSocket reconnect state.
4. Preserve non-secret context metadata.

Support `--all-contexts` only with confirmation.

## Static Token Compatibility

For automation and CI, support an explicit environment variable or flag:

```bash
A8S_TOKEN=<token> a8s project list
a8s project list --token <token>
```

Rules:

- flags override stored credentials
- environment tokens are never persisted
- warn that command-line token flags may appear in shell history
- prefer workload identity or short-lived service-account tokens for CI

The current `api_token` YAML setting should be treated as legacy and deprecated.

## Git Provider Authentication

GitHub and GitLab integration endpoints have special backend behavior and may accept provider tokens in some flows. Keep provider authentication separate from the primary Keycloak session:

- Keycloak access tokens authenticate the A8S user.
- Git provider tokens authorize repository-provider operations.
- Never overwrite the A8S bearer token with a provider token.
- Store provider credentials separately if the CLI later supports direct provider authentication.

## Security Requirements

- Validate issuer, audience, nonce, expiry, and state.
- Use PKCE S256.
- Never ship a Keycloak client secret in the CLI.
- Redact tokens and sensitive claims.
- Validate TLS by default.
- Do not accept tokens through insecure configuration without an explicit warning.
- Do not store admin service-account credentials in user contexts.
- Keep refresh tokens scoped to the minimum required privileges.

## Backend Changes Required Before Production

- Register a dedicated public `a8s-cli` Keycloak client with PKCE.
- Confirm whether device authorization is enabled.
- Define supported issuer and audience validation.
- Add or confirm token revocation/logout support.
- Protect currently public cluster, Kubernetes, Git integration, internal, and WebSocket routes appropriately.
- Ensure `/api/admin/documentation/**` requires `ROLE_ADMIN`.

## Authentication Acceptance Tests

- Login succeeds through browser PKCE flow.
- Invalid state and nonce are rejected.
- Access tokens refresh before expiry.
- Rotated refresh tokens are persisted.
- Failed refresh returns exit code `3`.
- Logout clears local credentials.
- Tokens never appear in stdout, stderr, or verbose logs.
- Admin commands return exit code `4` for authenticated non-admin users.
- Multiple contexts retain independent credentials.

