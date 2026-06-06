# A8S OpenAPI and Backend Compatibility

## Purpose

This document defines how the CLI tracks backend API changes, uses OpenAPI, detects breaking changes, and verifies client/server compatibility.

## OpenAPI Source

The backend exposes Springdoc routes such as:

```text
/v3/api-docs
/swagger-ui/
```

Generate and commit a sanitized OpenAPI snapshot:

```text
api/openapi.json
```

The snapshot must not contain environment-specific server URLs or secrets.

## Operation Requirements

Every CLI-eligible operation should have:

- stable `operationId`
- documented authentication requirement
- request and response schemas
- documented status codes
- error response schema
- tags matching backend features
- descriptions for asynchronous behavior

Internal callbacks and provider webhooks should be tagged as internal and excluded from CLI generation.

## Compatibility Policy

Use semantic versioning:

- patch: compatible fixes
- minor: additive endpoints, fields, and commands
- major: removals or incompatible behavior

The CLI must tolerate unknown additive JSON fields. It must not silently ignore missing required fields.

## Compatibility Handshake

Add a backend capabilities/version endpoint before production:

```http
GET /api/v1/meta
```

Recommended response:

```json
{
  "version": "1.0.0",
  "apiVersion": "v1",
  "capabilities": [
    "bakong-payments",
    "cluster-stream",
    "admin-events"
  ]
}
```

Commands should fail clearly when a required capability is absent.

## Drift Detection

CI should detect:

- removed operations
- method or path changes
- newly required request fields
- incompatible response changes
- removed enum values
- authentication changes
- new eligible operations without command mappings

Additive optional fields should not fail compatibility checks.

## Typed Client Strategy

OpenAPI may generate base request/response models, but handwritten service wrappers should own:

- CLI-specific arguments
- workflows
- retries
- output adaptation
- error normalization

Do not generate Cobra commands directly from raw OpenAPI without review.

## Backend Improvements Required

- standardize error schemas
- add stable operation IDs
- document all async terminal statuses
- document pagination and filtering
- document WebSocket/SSE payloads separately
- add version/capability endpoint
- correctly describe admin and ownership security

## Compatibility Tests

- current CLI supports the minimum documented backend version
- unsupported backend versions return actionable errors
- missing capabilities disable only affected commands
- endpoint catalog contains every CLI-eligible OpenAPI operation
- generated typed models compile after OpenAPI updates

