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
- `--file -` support for complex stdin input
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

