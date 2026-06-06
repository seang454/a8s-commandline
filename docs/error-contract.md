# A8S CLI Error Contract

## Purpose

This document defines how the CLI interprets backend errors, reports them to users, produces machine-readable output, retries requests, and selects process exit codes.

## Current Backend Reality

The backend currently returns more than one error shape.

Simple API error:

```json
{
  "status": 404,
  "message": "Resource not found",
  "timestamp": "2026-06-06T10:00:00Z"
}
```

Error with details:

```json
{
  "message": "Validation failed",
  "status": 400,
  "timestamp": "2026-06-06T10:00:00",
  "details": {
    "name": "Name is required"
  }
}
```

Spring Security, proxies, streaming endpoints, and some controllers may return empty bodies, plain text, or different JSON. The CLI must normalize these responses.

## Normalized CLI Error

All commands should return or wrap one shared error type:

```go
type Error struct {
    Code       string
    Message    string
    ExitCode   int
    HTTPStatus int
    RequestID  string
    Details    any
    RetryAfter time.Duration
    Cause      error
}
```

Recommended stable error codes:

```text
invalid_usage
validation_failed
authentication_required
authentication_failed
permission_denied
not_found
conflict
rate_limited
timeout
network_error
backend_unavailable
operation_failed
unexpected_response
internal_error
```

## Exit Codes

| Exit code | Meaning | Typical causes |
|---|---|---|
| `0` | Success | Command completed successfully. |
| `1` | General failure | Unexpected response, operation failure, or unclassified error. |
| `2` | Invalid usage or validation | Invalid arguments, flags, manifest, or backend `400` validation error. |
| `3` | Authentication failure | Missing, invalid, expired, or unrefreshable credentials; backend `401`. |
| `4` | Authorization failure | Backend `403`. |
| `5` | Resource not found | Backend `404`. |
| `6` | Conflict or invalid state | Backend `409`, duplicate resource, or invalid operation state. |
| `7` | Timeout | Context deadline, wait timeout, or stream timeout. |
| `8` | Backend unavailable | Network failure, DNS failure, TLS failure, backend `502`, `503`, or `504`. |
| `9` | Rate limited | Backend `429` after retry policy is exhausted. |
| `130` | Interrupted | User pressed Ctrl+C. |

## HTTP Status Mapping

| HTTP status | CLI code | Exit code | Retry |
|---|---|---:|---|
| `400`, `422` | `validation_failed` | `2` | No |
| `401` | `authentication_required` | `3` | Refresh once, then stop |
| `403` | `permission_denied` | `4` | No |
| `404` | `not_found` | `5` | No |
| `409` | `conflict` | `6` | No |
| `408` | `timeout` | `7` | Only safe requests |
| `429` | `rate_limited` | `9` | Respect `Retry-After` |
| `500` | `internal_error` | `1` | Generally no |
| `502`, `503`, `504` | `backend_unavailable` | `8` | Safe requests only |

For `DELETE`, create, deploy, payment, restore, rollback, and other mutations, do not retry automatically unless the backend provides an idempotency guarantee.

## Error Decoding Order

When a non-success HTTP response is received:

1. Capture HTTP status and safe response headers.
2. Capture `X-Request-ID`, `Request-ID`, or trace ID when present.
3. Limit the response body size before reading it.
4. Attempt to decode the detailed backend error shape.
5. Attempt to decode the simple backend error shape.
6. Attempt to decode common Spring Security or problem-detail shapes.
7. Fall back to sanitized plain text.
8. If the body is empty, use the HTTP status text and command context.
9. Redact secrets before returning the error.

Never return raw HTML proxy pages or unbounded bodies directly to users.

## Human-Readable Error Output

Default output goes to stderr:

```text
Error: project "abc" was not found
Code: not_found
HTTP status: 404
Request ID: req-123

Run `a8s project list` to view available projects.
```

Rules:

- Start with a concise actionable message.
- Include request ID when available.
- Include validation field details.
- Do not print stack traces unless a dedicated debug mode is enabled.
- Never print tokens, credentials, or sensitive request bodies.

## Machine-Readable Error Output

When `--output json` is selected, errors should be emitted as JSON to stderr:

```json
{
  "error": {
    "code": "validation_failed",
    "message": "Validation failed",
    "exitCode": 2,
    "httpStatus": 400,
    "requestId": "req-123",
    "details": {
      "name": "Name is required"
    }
  }
}
```

YAML output should use the same fields. Successful data remains on stdout; errors remain on stderr.

## Validation Errors

Local validation should happen before network calls when possible:

- required arguments and flags
- mutually exclusive flags
- file existence and syntax
- supported output formats
- valid identifiers, namespaces, and providers
- positive durations, limits, and resource quantities

Backend validation details should be preserved and displayed by field.

## Retry Policy

Retry only when all conditions are true:

- the failure is transient
- the operation is safe or idempotent
- the request body can be replayed safely
- the context deadline allows another attempt

Recommended defaults:

```text
Maximum attempts: 3
Initial delay:    250ms
Maximum delay:    2s
Backoff:          exponential with jitter
```

Respect `Retry-After`. Do not retry authentication or permission failures.

## Timeout and Cancellation

- Every request must receive a `context.Context`.
- Global `--timeout` applies to the complete command unless a command documents otherwise.
- Streaming commands run until Ctrl+C or timeout.
- Ctrl+C maps to exit code `130`.
- Workflow timeout messages must include the operation and last known status.

Example:

```text
Error: timed out waiting for cluster "cluster-123"
Last status: DEPLOYING
Code: timeout
```

## Asynchronous Operation Failures

Polling and streaming workflows must distinguish:

- API request failure
- terminal operation failure
- timeout
- user interruption

Terminal backend operation failures should use code `operation_failed`, include the final status, and preserve safe backend details.

## Partial Success

For commands operating on multiple resources:

- print successful items normally
- report failed items
- return non-zero if any requested operation failed
- include a summary suitable for automation

Example:

```text
Updated: 8
Failed:  2
```

## Backend Error Improvements Recommended

The backend currently has overlapping exception handlers and inconsistent error shapes. Before production, standardize all endpoints on one response:

```json
{
  "code": "not_found",
  "message": "Project not found",
  "status": 404,
  "timestamp": "2026-06-06T10:00:00Z",
  "requestId": "req-123",
  "details": {}
}
```

Recommended backend changes:

- use one global exception handler
- add a stable machine-readable `code`
- include request or trace IDs
- avoid returning raw exception messages for `500`
- use consistent UTC timestamp formatting
- return validation details consistently
- define error events for SSE and WebSocket streams
- return `Retry-After` for rate limiting and temporary unavailability

## Error Contract Acceptance Tests

- Each documented HTTP status maps to the expected exit code.
- Both current backend error shapes decode correctly.
- Empty, plain-text, and malformed error responses are handled safely.
- Validation details render in table, JSON, and YAML modes.
- `401` refreshes once and never loops.
- `429` respects `Retry-After`.
- Unsafe mutations are not automatically retried.
- Secrets are redacted from all errors and verbose logs.
- Ctrl+C exits with code `130`.
- Machine-readable errors remain stable for scripts.

