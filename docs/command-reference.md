# A8S CLI Command Reference Design

## Purpose

This file defines how the user-facing command reference should be generated and maintained. The complete endpoint-to-command inventory remains in `backend-api-cli-catalog.md`.

## Generation Policy

Generate the final command reference from the actual Cobra command tree so documentation cannot drift from implementation.

Recommended generator command:

```bash
go run ./scripts/generate-command-docs
```

Generated documentation should be written under:

```text
docs/commands/
|-- a8s.md
|-- a8s_project.md
|-- a8s_cluster.md
`-- ...
```

Do not manually edit generated command pages.

## Required Content Per Command

- command path
- short and long description
- usage
- positional arguments
- inherited and local flags
- examples
- output behavior
- destructive-operation warning
- authentication/role requirements
- related commands

## Command Conventions

```text
a8s <resource> <verb>
a8s admin <resource> <verb>
```

Preferred verbs:

```text
get, list, create, update, delete, status, watch, logs, apply, run,
restore, rollback, approve, reject, connect, disconnect
```

Avoid action-first top-level commands such as `a8s create user`.

## Mutation Input Policy

Every command that sends a configurable request payload to the backend must
support both:

- individual command flags
- a YAML or JSON operation document supplied with `--file <path>` or
  `--file -`

This rule applies even when the backend request contains only one configurable
field. Users may choose the most convenient form:

```bash
a8s database upgrade <deployment-id> --version 17
a8s database upgrade <deployment-id> --file upgrade.yaml

a8s project domain set <project-id> --domain api.example.com
a8s project domain set <project-id> --file domain.yaml

a8s scan start --image nginx:1.27
a8s scan start --file scan.yaml
```

Both input forms must resolve to the same typed internal request model and
produce the same backend payload.

Input precedence is:

```text
explicit flags > operation file > active-context defaults > backend defaults
```

Only flags explicitly supplied by the user override operation-file values.
Cobra default values must not accidentally replace values loaded from a file.

Commands without a configurable request payload do not accept operation
documents. This includes ordinary `get`, `list`, `status`, `watch`, `logs`,
`download`, and delete commands, plus payload-free actions such as restart or
sync. Their positional identifiers, output controls, confirmation flags, and
workflow controls remain normal arguments and flags.

File-content commands such as environment import, avatar upload, SQL query
files, and documentation upload use their domain file directly. That file is
not an operation YAML/JSON document.

## Required CLI-Only Commands

```bash
a8s auth login|status|logout
a8s context create|list|get|use|update|rename|delete
a8s config view|path|validate
a8s doctor
a8s completion bash|zsh|fish|powershell
a8s version --client --server
```

## Command Quality Checklist

- resource-first placement
- clear singular/plural behavior
- stable arguments
- safe destructive confirmation
- JSON/YAML output suitable for automation
- no secret values in help examples
- `--file` and equivalent flags for every configurable mutation payload,
  including small payloads
- `--file -` support for YAML/JSON operation input from stdin
- identical validation and backend payloads for file and flag input
- `--wait` and `--timeout` where asynchronous
- role requirements documented for admin commands

## Temporary Legacy Aliases

Legacy commands currently include:

```bash
a8s create user
a8s list users
a8s list projects
a8s delete user
```

If compatibility is needed, keep them as hidden or deprecated aliases that print migration guidance. Remove them in the next major CLI version.

## Reference Validation

CI should generate the Cobra reference and fail when committed generated files differ. It should also compare implemented command paths with `backend-api-cli-catalog.md`.
