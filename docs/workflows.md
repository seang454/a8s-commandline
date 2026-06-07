# A8S CLI Workflow Contracts

## Purpose

This document defines production behavior for asynchronous commands, including polling, streaming, terminal states, timeouts, cancellation, and fallback behavior.

## Common Wait Contract

Commands that support `--wait` must:

1. submit the operation
2. print or retain its resource/operation ID
3. prefer a stream or WebSocket when available
4. fall back to polling when streaming fails
5. stop on a terminal success or failure state
6. return the final resource with exit code `0` on success
7. return `operation_failed`, `timeout`, or `130` on failure, timeout, or interruption

Recommended flags:

```text
--wait
--timeout <duration>
--poll-interval <duration>
--watch
--output table|json|yaml
```

Unknown statuses are non-terminal unless the operation disappears or timeout expires.

## Default Timing

| Workflow | Poll interval | Default timeout |
|---|---:|---:|
| Project deployment | 3s | 20m |
| Microservice deployment | 3s | 30m |
| Database deployment/update | 5s | 30m |
| Cluster deployment/update | 5s | 45m |
| Backup | 5s | 60m |
| Restore | 5s | 90m |
| Image scan | 3s | 20m |
| Benchmark | 3s | 30m |
| Payment | 3s | 15m |

Users may override timeouts. Poll intervals below one second should be rejected.

## Status Normalization

The backend currently uses feature-specific strings. The CLI should normalize them internally:

```text
queued
running
succeeded
failed
cancelled
unknown
```

Preserve the raw backend status in structured output.

## Project and Microservice Deployments

Observed active statuses include:

```text
CREATED, PENDING, BUILDING, DEPLOYING
```

Observed successful statuses include:

```text
DEPLOYED, READY, SUCCEEDED, COMPLETED
```

Observed failure statuses include:

```text
FAILED, ERROR, CANCELLED, DELETED
```

Preferred sources:

1. Jenkins or deployment stream
2. project/release detail endpoint
3. project detail polling

On stream disconnect, reconnect with bounded backoff and then fall back to polling.

## Database Deployments

Observed active statuses include:

```text
PENDING, INSTALLING, UPDATING, UPGRADING, RESTARTING, RESTORING
```

Observed successful statuses include:

```text
DEPLOYED, READY, SUCCEEDED
```

Observed failure statuses include:

```text
FAILED, ERROR
```

`a8s database deploy --wait`, `update`, `upgrade`, `restart`, and restore workflows should poll the deployment detail endpoint. Success requires both a successful status and a usable final resource response.

## Cluster Deployments

Use the cluster deployment stream when available:

```text
/api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream
```

Fallback to cluster status/detail endpoints. A successful Helm release alone is insufficient if the final cluster record indicates failure.

## Backup and Restore

Backup and restore run statuses include:

```text
RUNNING, SUCCEEDED, FAILED
```

Rules:

- backup trigger returns immediately unless `--wait`
- only `SUCCEEDED` is terminal success
- `FAILED` is terminal failure
- restore cancellation should treat an already-terminal run as a conflict or no-op according to backend response
- download is allowed only for successful backup runs

## Image Scans

Use scan detail polling:

```bash
a8s scan start --image <image> --wait
```

Normalize queued/running/completed/failed variants. On success, optionally fetch the report with `--report` or a separate `scan report` command.

## Benchmarks

Known benchmark statuses:

```text
QUEUED, RUNNING, COMPLETED, FAILED
```

`COMPLETED` is success and `FAILED` is failure. Deleting a running benchmark should require confirmation and respect backend conflict behavior.

## Payments and Quota Purchases

Payment status endpoint results:

```text
PENDING, PAID, NO_PAYMENT_REQUIRED
```

Rules:

- `PAID` is terminal success
- `NO_PAYMENT_REQUIRED` is terminal success
- `PENDING` continues polling
- HTTP failure or timeout must not create a second payment automatically
- always retain and display the returned MD5 safely for status recovery
- refresh workspace entitlements after successful payment

The backend may approve the request during status polling, so status checks are not strictly read-only.

## Quota Approval

Admin quota request statuses:

```text
PENDING, APPROVED, REJECTED
```

Approval and rejection are terminal mutations. Require confirmation for rejection and display the final request state when available.

## Streaming and Reconnect

Recommended reconnect policy:

```text
Attempts before polling fallback: 5
Initial delay: 500ms
Maximum delay: 10s
Backoff: exponential with jitter
```

Never reconnect after explicit authentication or authorization failure. Refresh authentication once if appropriate.

## Progress Output

- human mode: progress to stderr, final result to stdout
- JSON/YAML mode: no spinner; emit final result only unless `--watch`
- `--watch --output json`: emit one JSON object per line
- never include secrets or complete tokenized WebSocket URLs

## Recovery and Resume

Every asynchronous create command should display identifiers immediately so users can resume:

```text
Cluster ID: cluster-123
Release: postgres-team-a

Resume:
  a8s cluster get cluster-123
  a8s cluster watch postgres-team-a
```

## Acceptance Criteria

- waits terminate correctly for every documented status
- unknown statuses do not cause false success
- stream failures fall back to polling
- Ctrl+C exits cleanly with code `130`
- timeouts include the last known status
- payment retries never create duplicate purchases
- workflow tests use deterministic fake clocks

## Current Implementation Status

Implemented:

- typed `a8s database deploy --wait`
- generic `--wait` on selected async generated commands
- polling via relative `Location`, `statusPath`, `statusUrl`, `operationPath`, or `operationUrl`
- built-in polling for image scan start, workspace quota payment status, cluster deployment status, and database backup restore fallback to deployment detail
- terminal status normalization for common success and failure strings
- command tests for scan, payment, and cluster wait behavior

Still required:

- workflow-specific typed clients for project, microservice, monolithic, database cluster, backup, restore, and admin approval workflows
- streaming-first behavior with reconnect before polling fallback
- richer resume output with operation IDs
- authenticated integration tests against the real backend for each critical workflow
