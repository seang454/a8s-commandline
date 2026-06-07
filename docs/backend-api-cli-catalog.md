# A8S Backend API to CLI Catalog

Generated from controller annotations in `D:\CSTADPreUniversityTraining\ITP\finalProject\a8s-backend`.

- Feature folders: 21
- Controllers: 38
- HTTP route patterns: 248
- CLI-eligible route patterns mapped: 238
- Automation-only route patterns excluded: 10
- Unmapped CLI-eligible route patterns: 0
- WebSocket routes: 4

Global CLI flags should include `--server`, `--context`, `--namespace`, `--target-cluster`, `--output`, `--timeout`, and `--verbose`.

## Recommended CLI Command Tree

Use resource-first Cobra command groups. Avoid generic top-level commands such as `a8s create user` or `a8s list projects`.

```text
a8s
|-- auth
|-- context                 # CLI-local server, token, namespace, and cluster contexts
|-- workspace
|   `-- quota
|-- profile
|-- project
|-- microservice
|-- database
|   `-- backup
|-- cluster
|   |-- backup
|   `-- console
|-- backup
|-- kubernetes
|-- logs
|-- git
|-- scan
|-- monitoring
|-- benchmark
|-- sonarqube
|-- defectdojo
|-- alert
|-- notification
|-- doctor
|-- completion
|-- version
`-- admin
    |-- user
    |-- project
    |-- cluster
    |-- quota
    |-- gitops
    |-- registry
    |-- sonarqube
    |-- monitoring
    |-- logs
    |-- docs
    `-- events
```

### Implementation order

1. Foundation: `auth`, `context`, configuration, shared API client, output formats, confirmation prompts, and error handling.
2. Core workflow: `workspace`, `profile`, `project`, `microservice`, `database`, `cluster`, and `backup`.
3. Operations: `kubernetes`, `logs`, `git`, `scan`, `monitoring`, and `notification`.
4. Quality and security: `benchmark`, `sonarqube`, `defectdojo`, and `alert`.
5. Administration: all commands under `a8s admin` with backend `ROLE_ADMIN` enforcement.

### Command design rules

- Use `get`, `list`, `create`, `update`, and `delete` consistently under each resource group.
- Require `--yes` for destructive commands and support `--dry-run` where the API permits it.
- Support `--output table|json|yaml`.
- Every command with a configurable request payload must support both equivalent
  flags and YAML/JSON input through `--file`, even when the request payload is
  small.
- Resolve mutation input using `explicit flags > operation file > active
context defaults > backend defaults`; only explicitly supplied flags may
  override file values.
- Commands without a configurable request payload, such as reads, downloads,
  deletes, streams, and payload-free actions, do not require operation-file
  input.
- Keep payment commands under `a8s workspace quota` because payment currently exists only for quota and plan purchases.
- Never expose internal callbacks, provider webhook receivers, or Jenkins completion callbacks as ordinary CLI commands.

## Authentication and Session Management

Authentication commands are CLI workflows rather than direct one-to-one endpoint mappings. The CLI should use Keycloak/OIDC login, securely store the resulting credentials, refresh access tokens when possible, and clear credentials on logout.

| Command                        | Behavior                                                                           |
| ------------------------------ | ---------------------------------------------------------------------------------- |
| `a8s auth login`               | Start browser or device-code login and store credentials for the active context.   |
| `a8s auth status`              | Show the authenticated identity, token expiry, active context, and detected roles. |
| `a8s auth logout`              | Revoke or clear locally stored credentials for the active context.                 |
| `a8s auth verify-email status` | Check authenticated email verification status.                                     |
| `a8s auth verify-email start`  | Request the backend to start email verification.                                   |

The backend remains the authorization authority. A local admin-role check may improve error messages, but every `a8s admin` operation must still be authorized by the backend.

## Context Configuration

Contexts are CLI-local configuration records. They select the backend server, credentials, default namespace, and optional target Kubernetes cluster.

```bash
a8s context create production --server https://api.example.com --namespace team-a
a8s context list
a8s context get production
a8s context use production
a8s context update production --target-cluster primary
a8s context delete production --yes
```

Recommended context precedence: explicit command flags, active context, environment variables, then built-in defaults. Store context metadata in `~/.a8s/config.yaml` and store secrets in the operating-system credential manager rather than directly in YAML.

## Global Flags

| Flag                   | Purpose                                                                                                                                |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | ----- | ------------------------------------------------- |
| `--server`             | Override the backend base URL.                                                                                                         |
| `--context`            | Run using a named context without changing the active context.                                                                         |
| `--namespace`          | Override the workspace or Kubernetes namespace.                                                                                        |
| `--target-cluster`     | Select a configured Kubernetes cluster alias.                                                                                          |
| `--output table        | json                                                                                                                                   | yaml` | Select machine-readable or human-readable output. |
| `--output-file <path>` | Write downloaded certificates, backups, reports, or other binary content to a file.                                                    |
| `--file <path>`        | Read a mutation request body from YAML or JSON; supported by every command with configurable payload data, and supports `-` for stdin. |
| `--wait`               | Wait until an asynchronous operation reaches a terminal state.                                                                         |
| `--timeout <duration>` | Limit request, polling, or streaming duration.                                                                                         |
| `--yes`                | Skip destructive-operation confirmation.                                                                                               |
| `--dry-run`            | Validate and display a request without applying it when supported.                                                                     |
| `--verbose`            | Print diagnostic request and workflow information without exposing secrets.                                                            |

List commands should additionally support pagination, filtering, sorting, and `--all` where the backend supports those behaviors.

## Workflow Commands

Workflow commands should combine multiple endpoints, polling, or WebSocket streams into one operator-friendly action.

| Workflow command                                     | Expected behavior                                                                                                  |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `a8s cluster create --file cluster.yaml --wait`      | Submit cluster deployment, watch deployment status, and return the final cluster record.                           |
| `a8s project deploy --file project.yaml --wait`      | Submit deployment, stream build progress, and return the deployed project.                                         |
| `a8s workspace quota purchase --plan premium --wait` | Create KHQR payment, display payment data, poll payment status, and refresh entitlements.                          |
| `a8s backup restore <type> <id> <run-id> --wait`     | Start restore, monitor completion, and report the final result.                                                    |
| `a8s doctor`                                         | Check configuration, authentication, backend reachability, workspace readiness, and optional cluster connectivity. |

## Streaming Commands

| Command                            | Transport                                    |
| ---------------------------------- | -------------------------------------------- |
| `a8s logs <pod-name> --follow`     | Kubernetes log stream endpoint.              |
| `a8s project logs --follow`        | Jenkins log stream and/or Jenkins WebSocket. |
| `a8s monitoring watch`             | `/ws/monitoring/overview`.                   |
| `a8s notification watch`           | `/ws/notifications`.                         |
| `a8s admin events watch`           | `/ws/admin/events`.                          |
| `a8s cluster watch <release-name>` | Cluster deployment stream endpoint.          |

Streaming commands should reconnect with bounded backoff, respect `--timeout`, stop cleanly on Ctrl+C, and print structured records when `--output json` is selected.

## Exit Codes

| Code | Meaning                                          |
| ---- | ------------------------------------------------ |
| `0`  | Success.                                         |
| `1`  | General or unexpected failure.                   |
| `2`  | Invalid command usage or validation failure.     |
| `3`  | Authentication required or token refresh failed. |
| `4`  | Authenticated but not authorized.                |
| `5`  | Requested resource not found.                    |
| `6`  | Conflict or invalid resource state.              |
| `7`  | Operation timed out.                             |
| `8`  | Backend unavailable or network failure.          |

Machine-readable error output should include an error code, message, HTTP status when available, request ID, and actionable details.

## Security Requirements

- Never print access tokens, refresh tokens, passwords, database credentials, payment payload secrets, or sensitive headers.
- Prefer the operating-system credential manager for tokens; restrict permissions if a file fallback is required.
- Require confirmation or `--yes` for delete, deactivate, reject, restore, rollback, password rotation, and destructive admin operations.
- Validate TLS certificates by default. Any insecure-development override must be explicit and visibly warned.
- Redact secrets from verbose logs, diagnostic bundles, shell completion, command history guidance, and error messages.
- Treat backend authorization as mandatory; the CLI must never attempt to bypass role or ownership checks.

## Implementation Status

| Area                        | Status      | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Backend endpoint discovery  | Complete    | All controller route patterns are scanned by this generator.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| CLI endpoint command design | Complete    | All CLI-eligible HTTP routes have suggested commands.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| Automation-only exclusions  | Complete    | Internal callbacks, provider webhooks, and Jenkins callbacks remain excluded.                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| WebSocket command design    | Complete    | Four WebSocket routes have suggested watch commands.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| Go/Cobra implementation     | In progress | All 209 unique catalog command paths are registered through the shared endpoint executor, including JSON/YAML bodies, multipart uploads, downloads, SSE responses, and four WebSocket watch routes. Keycloak PKCE login, secure context credentials, pre-command token refresh, one-time backend-401 refresh/replay, invalid-refresh cleanup, auth status, and local logout are implemented. `a8s database deploy` is typed; scan start, cluster deploy, and quota purchase have typed convenience flags; selected async operations support `--wait`. |
| End-to-end command tests    | In progress | Catalog-to-Cobra coverage, shared executor tests, command wait tests, and gated backend smoke-test scaffolding are implemented. Authenticated backend integration tests remain required per command group and workflow.                                                                                                                                                                                                                                                                                                                               |

## Example Operator Workflows

### Deploy and inspect an application

```bash
a8s project deploy --file project.yaml --wait
a8s project list
a8s project get <project-id>
a8s project logs --follow
```

### Create and operate a database cluster

```bash
a8s cluster create --file cluster.yaml --wait
a8s cluster get <cluster-id>
a8s cluster metrics <cluster-id>
a8s cluster certificate <cluster-id> --output-file ca.crt
```

### Back up and restore a database

```bash
a8s database backup run <deployment-id>
a8s database backup download <deployment-id> <run-id> --output-file backup.tar.gz
a8s database backup restore <deployment-id> <run-id> --wait
```

### Purchase workspace quota

```bash
a8s workspace quota pricing
a8s workspace quota purchase --plan premium --wait
a8s workspace entitlements
```

### Diagnose an incident

```bash
a8s doctor
a8s monitoring overview
a8s kubernetes events --warnings-only
a8s logs <pod-name> --container <container-name> --follow --tail 100
```

## Known Backend and CLI Limitations

- Payments currently use Bakong KHQR and are exposed only through workspace quota and plan purchase endpoints.
- `paymentProvider` mentions Stripe fields, but the current controller flow generates Bakong KHQR; do not advertise Stripe until backend support is complete.
- Payment status is queried by MD5 and currently returns `PENDING`, `PAID`, or `NO_PAYMENT_REQUIRED`.
- Several deployment and restore operations are asynchronous and require polling or streaming for a complete CLI experience.
- The CLI must not expose the 10 automation-only callback and webhook routes listed in this catalog.
- `context`, `doctor`, shell completion, local token storage, and some authentication commands are CLI-only features without direct backend endpoint mappings.

## admin

| Method   | Endpoint                                                       | Suggested CLI command                                     | Controller                    |
| -------- | -------------------------------------------------------------- | --------------------------------------------------------- | ----------------------------- |
| `GET`    | `/api/v1/admin/clusters`                                       | `a8s admin cluster list`                                  | `AdminClusterController`      |
| `PATCH`  | `/api/v1/admin/clusters/{clusterId}`                           | `a8s admin cluster update <cluster-id>`                   | `AdminClusterController`      |
| `GET`    | `/api/v1/admin/clusters/kubernetes`                            | `a8s admin cluster nodes`                                 | `AdminClusterController`      |
| `GET`    | `/api/v1/admin/clusters/kubernetes/{alias}/health`             | `a8s admin cluster health <alias>`                        | `AdminClusterController`      |
| `GET`    | `/api/v1/admin/clusters/kubernetes/{alias}/quotas`             | `a8s admin cluster quota list <alias>`                    | `AdminClusterController`      |
| `DELETE` | `/api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}` | `a8s admin cluster quota delete <alias> <namespace>`      | `AdminClusterController`      |
| `PUT`    | `/api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}` | `a8s admin cluster quota set <alias> <namespace>`         | `AdminClusterController`      |
| `GET`    | `/api/v1/admin/projects`                                       | `a8s admin project list`                                  | `AdminController`             |
| `DELETE` | `/api/v1/admin/projects/{projectId}`                           | `a8s admin project deactivate <project-id>`               | `AdminController`             |
| `PATCH`  | `/api/v1/admin/projects/{projectId}`                           | `a8s admin project update <project-id>`                   | `AdminController`             |
| `POST`   | `/api/v1/admin/projects/{projectId}/restore`                   | `a8s admin project restore <project-id>`                  | `AdminController`             |
| `GET`    | `/api/v1/admin/users`                                          | `a8s admin user list`                                     | `AdminController`             |
| `POST`   | `/api/v1/admin/users`                                          | `a8s admin user create`                                   | `AdminController`             |
| `DELETE` | `/api/v1/admin/users/{userId}`                                 | `a8s admin user deactivate <user-id>`                     | `AdminController`             |
| `PATCH`  | `/api/v1/admin/users/{userId}`                                 | `a8s admin user update <user-id>`                         | `AdminController`             |
| `POST`   | `/api/v1/admin/users/{userId}/reactivate`                      | `a8s admin user reactivate <user-id>`                     | `AdminController`             |
| `POST`   | `/api/v1/admin/gitops/apps`                                    | `a8s admin gitops app create`                             | `AdminGitOpsController`       |
| `POST`   | `/api/v1/admin/gitops/apps/{appId}/abort`                      | `a8s admin gitops app abort <app-id>`                     | `AdminGitOpsController`       |
| `POST`   | `/api/v1/admin/gitops/apps/{appId}/retry`                      | `a8s admin gitops app retry <app-id>`                     | `AdminGitOpsController`       |
| `POST`   | `/api/v1/admin/gitops/apps/{appId}/sync`                       | `a8s admin gitops app sync <app-id>`                      | `AdminGitOpsController`       |
| `GET`    | `/api/v1/admin/gitops/overview`                                | `a8s admin gitops overview`                               | `AdminGitOpsController`       |
| `GET`    | `/api/v1/admin/logs/clusters`                                  | `a8s admin logs clusters`                                 | `AdminLogsController`         |
| `GET`    | `/api/v1/admin/logs/namespaces`                                | `a8s admin logs namespaces`                               | `AdminLogsController`         |
| `GET`    | `/api/v1/admin/logs/pods`                                      | `a8s admin logs pods`                                     | `AdminLogsController`         |
| `GET`    | `/api/v1/admin/logs/query`                                     | `a8s admin logs query`                                    | `AdminLogsController`         |
| `GET`    | `/api/v1/admin/logs/workloads`                                 | `a8s admin logs workloads`                                | `AdminLogsController`         |
| `GET`    | `/api/v1/admin/monitoring/overview`                            | `a8s admin monitoring overview`                           | `AdminMonitoringController`   |
| `GET`    | `/api/v1/admin/quota-requests`                                 | `a8s admin quota list`                                    | `AdminQuotaRequestController` |
| `POST`   | `/api/v1/admin/quota-requests/{id}/approve`                    | `a8s admin quota approve <request-id>`                    | `AdminQuotaRequestController` |
| `POST`   | `/api/v1/admin/quota-requests/{id}/reject`                     | `a8s admin quota reject <request-id>`                     | `AdminQuotaRequestController` |
| `GET`    | `/api/v1/admin/registry/health`                                | `a8s admin registry health`                               | `AdminRegistryController`     |
| `GET`    | `/api/v1/admin/registry/projects`                              | `a8s admin registry project list`                         | `AdminRegistryController`     |
| `POST`   | `/api/v1/admin/registry/projects`                              | `a8s admin registry project create`                       | `AdminRegistryController`     |
| `DELETE` | `/api/v1/admin/registry/projects/{projectName}/artifacts`      | `a8s admin registry artifact delete <project-name>`       | `AdminRegistryController`     |
| `GET`    | `/api/v1/admin/registry/projects/{projectName}/artifacts`      | `a8s admin registry artifact list <project-name>`         | `AdminRegistryController`     |
| `DELETE` | `/api/v1/admin/registry/projects/{projectName}/repositories`   | `a8s admin registry repository delete <project-name>`     | `AdminRegistryController`     |
| `GET`    | `/api/v1/admin/registry/projects/{projectName}/repositories`   | `a8s admin registry repository list <project-name>`       | `AdminRegistryController`     |
| `GET`    | `/api/v1/admin/sonarqube/projects`                             | `a8s admin sonarqube project list`                        | `AdminSonarQubeController`    |
| `GET`    | `/api/v1/admin/sonarqube/projects/{projectId}`                 | `a8s admin sonarqube project get <project-id>`            | `AdminSonarQubeController`    |
| `GET`    | `/api/v1/admin/sonarqube/server-projects`                      | `a8s admin sonarqube server-project list`                 | `AdminSonarQubeController`    |
| `POST`   | `/api/v1/admin/sonarqube/server-projects`                      | `a8s admin sonarqube server-project create`               | `AdminSonarQubeController`    |
| `DELETE` | `/api/v1/admin/sonarqube/server-projects/{projectKey}`         | `a8s admin sonarqube server-project delete <project-key>` | `AdminSonarQubeController`    |
| `GET`    | `/api/v1/admin/sonarqube/server-projects/{projectKey}`         | `a8s admin sonarqube server-project get <project-key>`    | `AdminSonarQubeController`    |
| `PATCH`  | `/api/v1/admin/sonarqube/server-projects/{projectKey}`         | `a8s admin sonarqube server-project update <project-key>` | `AdminSonarQubeController`    |

## alerts

| Method   | Endpoint                                     | Suggested CLI command                       | Controller        |
| -------- | -------------------------------------------- | ------------------------------------------- | ----------------- |
| `GET`    | `/api/v1/alerts/channels`                    | `a8s alert channel list`                    | `AlertController` |
| `POST`   | `/api/v1/alerts/channels`                    | `a8s alert channel create`                  | `AlertController` |
| `DELETE` | `/api/v1/alerts/channels/{channelId}`        | `a8s alert channel delete <channel-id>`     | `AlertController` |
| `PUT`    | `/api/v1/alerts/channels/{channelId}`        | `a8s alert channel update <channel-id>`     | `AlertController` |
| `GET`    | `/api/v1/alerts/projects/{projectId}/config` | `a8s alert project-config get <project-id>` | `AlertController` |
| `PUT`    | `/api/v1/alerts/projects/{projectId}/config` | `a8s alert project-config set <project-id>` | `AlertController` |
| `GET`    | `/api/v1/alerts/projects/configs`            | `a8s alert project-config list`             | `AlertController` |
| `GET`    | `/api/v1/alerts/user-config`                 | `a8s alert user-config get`                 | `AlertController` |
| `PUT`    | `/api/v1/alerts/user-config`                 | `a8s alert user-config set`                 | `AlertController` |

## auth

| Method | Endpoint                                                    | Suggested CLI command          | Controller       |
| ------ | ----------------------------------------------------------- | ------------------------------ | ---------------- |
| `GET`  | `/api/v1/auth/keycloak/users/{keycloakUserId}/verify-email` | `a8s auth verify-email status` | `AuthController` |
| `POST` | `/api/v1/auth/keycloak/users/{keycloakUserId}/verify-email` | `a8s auth verify-email start`  | `AuthController` |
| `GET`  | `/api/v1/auth/session/onboarding`                           | `a8s auth onboarding status`   | `AuthController` |
| `POST` | `/api/v1/auth/session/onboarding`                           | `a8s auth onboarding start`    | `AuthController` |

## databasebackup

| Method   | Endpoint                                                                         | Suggested CLI command                                         | Controller                           |
| -------- | -------------------------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------------ |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/backup`                             | `a8s database backup settings get <deployment-id>`            | `DatabaseDeploymentBackupController` |
| `PATCH`  | `/api/v1/database-deployments/{deploymentId}/backup`                             | `a8s database backup settings set <deployment-id>`            | `DatabaseDeploymentBackupController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/backup/run`                         | `a8s database backup run <deployment-id>`                     | `DatabaseDeploymentBackupController` |
| `DELETE` | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}`                | `a8s database backup delete <deployment-id> <run-id>`         | `DatabaseDeploymentBackupController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/download`       | `a8s database backup download <deployment-id> <run-id>`       | `DatabaseDeploymentBackupController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore`        | `a8s database backup restore <deployment-id> <run-id>`        | `DatabaseDeploymentBackupController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore/cancel` | `a8s database backup restore cancel <deployment-id> <run-id>` | `DatabaseDeploymentBackupController` |
| `POST`   | `/api/internal/backups/callback/{deploymentId}/{targetType}`                     | `(internal service callback; no user CLI command)`            | `InternalBackupController`           |
| `DELETE` | `/api/backups/{targetType}/{id}/{runId}`                                         | `a8s backup delete <type> <id> <run-id>`                      | `UnifiedBackupController`            |
| `GET`    | `/api/backups/download/{targetType}/{id}/{runId}`                                | `a8s backup download <type> <id> <run-id>`                    | `UnifiedBackupController`            |
| `POST`   | `/api/backups/restore/{targetType}/{id}/{runId}`                                 | `a8s backup restore <type> <id> <run-id>`                     | `UnifiedBackupController`            |
| `POST`   | `/api/backups/restore/{targetType}/{id}/{runId}/cancel`                          | `a8s backup restore cancel <type> <id> <run-id>`              | `UnifiedBackupController`            |
| `GET`    | `/api/backups/settings/{targetType}/{id}`                                        | `a8s backup settings get <type> <id>`                         | `UnifiedBackupController`            |
| `POST`   | `/api/backups/settings/{targetType}/{id}`                                        | `a8s backup settings set <type> <id>`                         | `UnifiedBackupController`            |
| `POST`   | `/api/backups/trigger/{targetType}/{id}`                                         | `a8s backup trigger <type> <id>`                              | `UnifiedBackupController`            |

## databaseconsole

No standalone controller. Console APIs are exposed through `singledb` and `dbcluster` endpoints.

## dbcluster

| Method   | Endpoint                                                                          | Suggested CLI command                                      | Controller                                                                           |
| -------- | --------------------------------------------------------------------------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `GET`    | `/api/namespaces/{namespace}/clusters`                                            | `a8s cluster list`                                         | `ClusterController`                                                                  |
| `DELETE` | `/api/namespaces/{namespace}/clusters/{id}`                                       | `a8s cluster delete <cluster-id>`                          | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}`                                       | `a8s cluster get <cluster-id>`                             | `ClusterController`                                                                  |
| `PATCH`  | `/api/namespaces/{namespace}/clusters/{id}`                                       | `a8s cluster update <cluster-id>`                          | `ClusterController`                                                                  |
| `PATCH`  | `/api/namespaces/{namespace}/clusters/{id}/backup`                                | `a8s cluster backup settings set <cluster-id>`             | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/certificate`                           | `a8s cluster certificate <cluster-id>`                     | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/console/credentials`                   | `a8s cluster console credentials <cluster-id>`             | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/console/data`                          | `a8s cluster console data <cluster-id>`                    | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/console/deployment`                    | `a8s cluster console deployment <cluster-id>`              | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/console/namespaces`                    | `a8s cluster console namespaces <cluster-id>`              | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/console/objects`                       | `a8s cluster console objects <cluster-id>`                 | `ClusterController`                                                                  |
| `POST`   | `/api/namespaces/{namespace}/clusters/{id}/console/query`                         | `a8s cluster console query <cluster-id>`                   | `ClusterController`                                                                  |
| `POST`   | `/api/namespaces/{namespace}/clusters/{id}/console/test`                          | `a8s cluster console test <cluster-id>`                    | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/deployments`                           | `a8s cluster history <cluster-id>`                         | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/metrics`                               | `a8s cluster metrics <cluster-id>`                         | `ClusterController`                                                                  |
| `PATCH`  | `/api/namespaces/{namespace}/clusters/{id}/settings`                              | `a8s cluster settings update <cluster-id>`                 | `ClusterController`                                                                  |
| `POST`   | `/api/namespaces/{namespace}/clusters/{id}/upgrade-version`                       | `a8s cluster upgrade <cluster-id>`                         | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/values`                                | `a8s cluster values <cluster-id>`                          | `ClusterController`                                                                  |
| `GET`    | `/api/namespaces/{namespace}/clusters/{id}/values/full`                           | `a8s cluster values <cluster-id> --full`                   | `ClusterController`                                                                  |
| `POST`   | `/api/namespaces/{namespace}/clusters/clone-from-backup`                          | `a8s cluster clone-from-backup`                            | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters`                                 | `a8s cluster list`                                         | `ClusterController`                                                                  |
| `DELETE` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}`                            | `a8s cluster delete <cluster-id>`                          | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}`                            | `a8s cluster get <cluster-id>`                             | `ClusterController`                                                                  |
| `PATCH`  | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}`                            | `a8s cluster update <cluster-id>`                          | `ClusterController`                                                                  |
| `PATCH`  | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}`                            | `a8s cluster rotate-password <cluster-id>`                 | `ClusterController` workflow alias that updates the engine-specific `secrets` field. |
| `PATCH`  | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/backup`                     | `a8s cluster backup settings set <cluster-id>`             | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/certificate`                | `a8s cluster certificate <cluster-id>`                     | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/credentials`        | `a8s cluster console credentials <cluster-id>`             | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/data`               | `a8s cluster console data <cluster-id>`                    | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/deployment`         | `a8s cluster console deployment <cluster-id>`              | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/namespaces`         | `a8s cluster console namespaces <cluster-id>`              | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/objects`            | `a8s cluster console objects <cluster-id>`                 | `ClusterController`                                                                  |
| `POST`   | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/query`              | `a8s cluster console query <cluster-id>`                   | `ClusterController`                                                                  |
| `POST`   | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/test`               | `a8s cluster console test <cluster-id>`                    | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/deployments`                | `a8s cluster history <cluster-id>`                         | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/metrics`                    | `a8s cluster metrics <cluster-id>`                         | `ClusterController`                                                                  |
| `PATCH`  | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/settings`                   | `a8s cluster settings update <cluster-id>`                 | `ClusterController`                                                                  |
| `POST`   | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/upgrade-version`            | `a8s cluster upgrade <cluster-id>`                         | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/values`                     | `a8s cluster values <cluster-id>`                          | `ClusterController`                                                                  |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/values/full`                | `a8s cluster values <cluster-id> --full`                   | `ClusterController`                                                                  |
| `POST`   | `/api/v1/cluster/namespaces/{namespace}/clusters/clone-from-backup`               | `a8s cluster clone-from-backup`                            | `ClusterController`                                                                  |
| `POST`   | `/api/namespaces/{namespace}/cluster-deployments`                                 | `a8s cluster deploy`                                       | `ClusterDeploymentController`                                                        |
| `GET`    | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}`                   | `a8s cluster status <release-name>`                        | `ClusterDeploymentController`                                                        |
| `PATCH`  | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}/backup`            | `a8s cluster backup settings set --release <release-name>` | `ClusterDeploymentController`                                                        |
| `GET`    | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}/values`            | `a8s cluster deployment values <release-name>`             | `ClusterDeploymentController`                                                        |
| `POST`   | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments`                      | `a8s cluster deploy`                                       | `ClusterDeploymentController`                                                        |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}`        | `a8s cluster status <release-name>`                        | `ClusterDeploymentController`                                                        |
| `PATCH`  | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}/backup` | `a8s cluster backup settings set --release <release-name>` | `ClusterDeploymentController`                                                        |
| `GET`    | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}/values` | `a8s cluster deployment values <release-name>`             | `ClusterDeploymentController`                                                        |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/database-resources`                       | `a8s kubernetes database-resources`                        | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/events`                                   | `a8s kubernetes events`                                    | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/overview`                                 | `a8s kubernetes overview`                                  | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/persistent-volume-claims`                 | `a8s kubernetes pvc`                                       | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/pods`                                     | `a8s kubernetes pods`                                      | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/pods/{podName}/logs/stream`               | `a8s logs <pod-name> --follow`                             | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream` | `a8s cluster watch <release-name>`                         | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/namespaces/{namespace}/services`                                 | `a8s kubernetes services`                                  | `KubernetesController`                                                               |
| `GET`    | `/api/kubernetes/test`                                                            | `a8s kubernetes test`                                      | `KubernetesController`                                                               |

## documentation

| Method   | Endpoint                                | Suggested CLI command         | Controller                |
| -------- | --------------------------------------- | ----------------------------- | ------------------------- |
| `DELETE` | `/api/admin/documentation/content`      | `a8s admin docs delete`       | `DocumentationController` |
| `GET`    | `/api/admin/documentation/content`      | `a8s admin docs get`          | `DocumentationController` |
| `PUT`    | `/api/admin/documentation/content`      | `a8s admin docs update`       | `DocumentationController` |
| `GET`    | `/api/admin/documentation/files`        | `a8s admin docs files`        | `DocumentationController` |
| `POST`   | `/api/admin/documentation/publish`      | `a8s admin docs publish`      | `DocumentationController` |
| `GET`    | `/api/admin/documentation/publish/logs` | `a8s admin docs publish-logs` | `DocumentationController` |

## entitlements

| Method | Endpoint                          | Suggested CLI command        | Controller                       |
| ------ | --------------------------------- | ---------------------------- | -------------------------------- |
| `GET`  | `/api/v1/workspaces/entitlements` | `a8s workspace entitlements` | `WorkspaceEntitlementController` |

## gitintegration

| Method   | Endpoint                                                  | Suggested CLI command           | Controller                 |
| -------- | --------------------------------------------------------- | ------------------------------- | -------------------------- |
| `DELETE` | `/api/v1/git-integrations/{provider}`                     | `a8s git disconnect <provider>` | `GitIntegrationController` |
| `GET`    | `/api/v1/git-integrations/{provider}/brokered-account`    | `a8s git account <provider>`    | `GitIntegrationController` |
| `POST`   | `/api/v1/git-integrations/{provider}/connect`             | `a8s git connect <provider>`    | `GitIntegrationController` |
| `GET`    | `/api/v1/git-integrations/{provider}/repos`               | `a8s git repos <provider>`      | `GitIntegrationController` |
| `GET`    | `/api/v1/git-integrations/{provider}/state`               | `a8s git state <provider>`      | `GitIntegrationController` |
| `POST`   | `/api/v1/git-integrations/{provider}/sync-keycloak-token` | `a8s git sync-token <provider>` | `GitIntegrationController` |
| `GET`    | `/api/v1/git-integrations/linked-providers`               | `a8s git providers`             | `GitIntegrationController` |

## imagescanner

| Method | Endpoint                                        | Suggested CLI command                              | Controller                               |
| ------ | ----------------------------------------------- | -------------------------------------------------- | ---------------------------------------- |
| `GET`  | `/api/v1/image-scanner/images`                  | `a8s scan images`                                  | `ImageScannerController`                 |
| `GET`  | `/api/v1/image-scanner/scans`                   | `a8s scan list`                                    | `ImageScannerController`                 |
| `POST` | `/api/v1/image-scanner/scans`                   | `a8s scan start`                                   | `ImageScannerController`                 |
| `GET`  | `/api/v1/image-scanner/scans/{scanId}`          | `a8s scan get <scan-id>`                           | `ImageScannerController`                 |
| `GET`  | `/api/v1/image-scanner/scans/{scanId}/report`   | `a8s scan report <scan-id>`                        | `ImageScannerController`                 |
| `POST` | `/api/internal/image-scanner/callback/{scanId}` | `(internal service callback; no user CLI command)` | `InternalImageScannerCallbackController` |

## microservice

| Method   | Endpoint                                                                             | Suggested CLI command                                        | Controller                                    |
| -------- | ------------------------------------------------------------------------------------ | ------------------------------------------------------------ | --------------------------------------------- |
| `GET`    | `/api/internal/microservices/source-archive`                                         | `(internal service callback; no user CLI command)`           | `InternalMicroserviceSourceArchiveController` |
| `GET`    | `/api/internal/microservices/defectdojo-token`                                       | `(internal service callback; no user CLI command)`           | `InternalWorkspaceDefectDojoTokenController`  |
| `POST`   | `/api/v1/projects/microservices`                                                     | `a8s microservice deploy`                                    | `MicroserviceProjectController`               |
| `DELETE` | `/api/v1/projects/microservices/{projectId}`                                         | `a8s microservice delete <project-id>`                       | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}`                                         | `a8s microservice get <project-id>`                          | `MicroserviceProjectController`               |
| `PUT`    | `/api/v1/projects/microservices/{projectId}/canvas`                                  | `a8s microservice apply <project-id>`                        | `MicroserviceProjectController`               |
| `PATCH`  | `/api/v1/projects/microservices/{projectId}/domains`                                 | `a8s microservice domains update <project-id>`               | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}/history`                                 | `a8s microservice history list <project-id>`                 | `MicroserviceProjectController`               |
| `DELETE` | `/api/v1/projects/microservices/{projectId}/history/{snapshotId}`                    | `a8s microservice history delete <project-id> <snapshot-id>` | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}/readiness`                               | `a8s microservice readiness <project-id>`                    | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/{projectId}/redeploy`                                | `a8s microservice redeploy <project-id>`                     | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/{projectId}/rollback`                                | `a8s microservice rollback <project-id>`                     | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}/runtime-pods`                            | `a8s microservice pods <project-id>`                         | `MicroserviceProjectController`               |
| `DELETE` | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment`        | `a8s microservice env clear <project-id> <service-id>`       | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment`        | `a8s microservice env get <project-id> <service-id>`         | `MicroserviceProjectController`               |
| `PUT`    | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment`        | `a8s microservice env set <project-id> <service-id>`         | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment/import` | `a8s microservice env import <project-id> <service-id>`      | `MicroserviceProjectController`               |
| `GET`    | `/api/v1/projects/microservices/{projectId}/webhook`                                 | `a8s microservice webhook get <project-id>`                  | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/{projectId}/webhook`                                 | `a8s microservice webhook update <project-id>`               | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/detect`                                              | `a8s microservice detect --repo`                             | `MicroserviceProjectController`               |
| `POST`   | `/api/v1/projects/microservices/detect/upload`                                       | `a8s microservice detect --source-archive <path>`            | `MicroserviceProjectController`               |

## monitoring

| Method | Endpoint                      | Suggested CLI command     | Controller             |
| ------ | ----------------------------- | ------------------------- | ---------------------- |
| `GET`  | `/api/v1/monitoring/overview` | `a8s monitoring overview` | `MonitoringController` |

## monolithic

| Method   | Endpoint                                                     | Suggested CLI command                                    | Controller                     |
| -------- | ------------------------------------------------------------ | -------------------------------------------------------- | ------------------------------ |
| `GET`    | `/api/v1/projects`                                           | `a8s project list`                                       | `ProjectController`            |
| `POST`   | `/api/v1/projects`                                           | `a8s project deploy`                                     | `ProjectController`            |
| `POST`   | `/api/v1/projects`                                           | `a8s project deploy`                                     | `ProjectController`            |
| `DELETE` | `/api/v1/projects/{projectId}`                               | `a8s project delete <project-id>`                        | `ProjectController`            |
| `GET`    | `/api/v1/projects/{projectId}`                               | `a8s project get <project-id>`                           | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/delete/complete`               | `(Jenkins callback; no user CLI command)`                | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/delete/failed`                 | `(Jenkins callback; no user CLI command)`                | `ProjectController`            |
| `PATCH`  | `/api/v1/projects/{projectId}/domain`                        | `a8s project domain set <project-id>`                    | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/domain/sync`                   | `a8s project domain sync <project-id>`                   | `ProjectController`            |
| `GET`    | `/api/v1/projects/{projectId}/releases`                      | `a8s project releases <project-id>`                      | `ProjectController`            |
| `DELETE` | `/api/v1/projects/{projectId}/releases/{releaseId}`          | `a8s project release delete <project-id> <release-id>`   | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/releases/{releaseId}/complete` | `(Jenkins callback; no user CLI command)`                | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/releases/{releaseId}/failed`   | `(Jenkins callback; no user CLI command)`                | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/releases/{releaseId}/rollback` | `a8s project release rollback <project-id> <release-id>` | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/repository/connect`            | `a8s project repository connect <project-id>`            | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/rollback`                      | `a8s project rollback <project-id>`                      | `ProjectController`            |
| `PATCH`  | `/api/v1/projects/{projectId}/settings`                      | `a8s project settings update <project-id>`               | `ProjectController`            |
| `POST`   | `/api/v1/projects/{projectId}/sync`                          | `a8s project redeploy <project-id>`                      | `ProjectController`            |
| `GET`    | `/api/v1/projects/me`                                        | `a8s project list`                                       | `ProjectController`            |
| `GET`    | `/api/v1/projects/{projectId}/environment`                   | `a8s project env get <project-id>`                       | `ProjectEnvironmentController` |
| `PUT`    | `/api/v1/projects/{projectId}/environment`                   | `a8s project env set <project-id>`                       | `ProjectEnvironmentController` |
| `POST`   | `/api/v1/projects/{projectId}/environment/import`            | `a8s project env import <project-id>`                    | `ProjectEnvironmentController` |
| `PATCH`  | `/api/v1/projects/{projectId}/auto-deploy`                   | `a8s project auto-deploy set <project-id>`               | `ProjectWebhookController`     |
| `GET`    | `/api/v1/projects/{projectId}/branches`                      | `a8s project branches <project-id>`                      | `ProjectWebhookController`     |
| `DELETE` | `/api/v1/projects/{projectId}/webhook`                       | `a8s project webhook delete <project-id>`                | `ProjectWebhookController`     |
| `GET`    | `/api/v1/projects/{projectId}/webhook`                       | `a8s project webhook get <project-id>`                   | `ProjectWebhookController`     |
| `POST`   | `/api/v1/projects/{projectId}/webhook`                       | `a8s project webhook create <project-id>`                | `ProjectWebhookController`     |
| `POST`   | `/api/v1/projects/{projectId}/webhook/rotate`                | `a8s project webhook rotate <project-id>`                | `ProjectWebhookController`     |
| `POST`   | `/api/v1/webhooks/github`                                    | `(provider webhook receiver; no user CLI command)`       | `WebhookController`            |
| `POST`   | `/api/v1/webhooks/gitlab`                                    | `(provider webhook receiver; no user CLI command)`       | `WebhookController`            |

## notifications

| Method | Endpoint                                   | Suggested CLI command                     | Controller               |
| ------ | ------------------------------------------ | ----------------------------------------- | ------------------------ |
| `POST` | `/api/notifications/{notificationId}/read` | `a8s notification read <notification-id>` | `NotificationController` |
| `GET`  | `/api/notifications/history/{userId}`      | `a8s notification list`                   | `NotificationController` |
| `GET`  | `/api/notifications/preferences/{userId}`  | `a8s notification preferences get`        | `NotificationController` |
| `POST` | `/api/notifications/preferences/{userId}`  | `a8s notification preferences set`        | `NotificationController` |

## payments

No standalone payment controller. Payments currently support Bakong KHQR purchases for workspace quota and plan upgrades.

### Payment and quota-purchase endpoints

| Method | Endpoint                                                     | Suggested CLI command                        | Purpose                                                         |
| ------ | ------------------------------------------------------------ | -------------------------------------------- | --------------------------------------------------------------- |
| `GET`  | `/api/v1/workspaces/quota-pricing`                           | `a8s workspace quota pricing`                | Get unit prices and plan prices.                                |
| `POST` | `/api/v1/workspaces/quota-requests`                          | `a8s workspace quota purchase --plan <plan>` | Submit a paid quota request and generate a Bakong KHQR payload. |
| `GET`  | `/api/v1/workspaces/quota-requests/payment-status?md5=<md5>` | `a8s workspace quota payment-status <md5>`   | Check payment status and apply the quota upgrade after payment. |

The purchase request accepts `requestedCpu`, `requestedMemory`, `requestedStorage`, `reason`, `isPaid`, `planName`, and `paymentProvider`. Set `isPaid` to `true` and `paymentProvider` to `BAKONG` to generate KHQR.

The purchase response contains `qrString` and `md5`. Use the returned `md5` when polling payment status. Status responses currently include `PENDING`, `PAID`, and `NO_PAYMENT_REQUIRED`.

When payment is confirmed, the backend approves the quota request, applies the workspace quota, activates the subscription for 30 days, and sends a payment receipt notification.

### Related admin endpoints

| Method | Endpoint                                    | Suggested CLI command                  | Purpose                                          |
| ------ | ------------------------------------------- | -------------------------------------- | ------------------------------------------------ |
| `GET`  | `/api/v1/admin/quota-requests`              | `a8s admin quota list`                 | List pending quota and payment-related requests. |
| `POST` | `/api/v1/admin/quota-requests/{id}/approve` | `a8s admin quota approve <request-id>` | Approve a pending request and apply its quota.   |
| `POST` | `/api/v1/admin/quota-requests/{id}/reject`  | `a8s admin quota reject <request-id>`  | Reject a pending request.                        |

## profile

| Method   | Endpoint                            | Suggested CLI command            | Controller          |
| -------- | ----------------------------------- | -------------------------------- | ------------------- |
| `DELETE` | `/api/v1/profile/me`                | `a8s profile account delete`     | `ProfileController` |
| `GET`    | `/api/v1/profile/me`                | `a8s profile get`                | `ProfileController` |
| `PATCH`  | `/api/v1/profile/me`                | `a8s profile update`             | `ProfileController` |
| `GET`    | `/api/v1/profile/me/account-status` | `a8s profile account status`     | `ProfileController` |
| `DELETE` | `/api/v1/profile/me/avatar`         | `a8s profile avatar delete`      | `ProfileController` |
| `GET`    | `/api/v1/profile/me/avatar`         | `a8s profile avatar download`    | `ProfileController` |
| `POST`   | `/api/v1/profile/me/avatar`         | `a8s profile avatar upload`      | `ProfileController` |
| `POST`   | `/api/v1/profile/me/deactivate`     | `a8s profile account deactivate` | `ProfileController` |
| `POST`   | `/api/v1/profile/me/delete`         | `a8s profile account delete`     | `ProfileController` |
| `POST`   | `/api/v1/profile/me/reactivate`     | `a8s profile account reactivate` | `ProfileController` |

## projects

| Method | Endpoint                                        | Suggested CLI command                    | Controller                    |
| ------ | ----------------------------------------------- | ---------------------------------------- | ----------------------------- |
| `GET`  | `/api/v1/jenkins/logs/stream`                   | `a8s project logs --follow`              | `JenkinsController`           |
| `GET`  | `/api/v1/projects/live`                         | `a8s project live list`                  | `LiveProjectController`       |
| `GET`  | `/api/v1/projects/{projectId}/defectdojo`       | `a8s defectdojo access <project-id>`     | `ProjectDefectDojoController` |
| `PUT`  | `/api/v1/projects/{projectId}/defectdojo/token` | `a8s defectdojo token sync <project-id>` | `ProjectDefectDojoController` |

## singledb

| Method   | Endpoint                                                         | Suggested CLI command                             | Controller                     |
| -------- | ---------------------------------------------------------------- | ------------------------------------------------- | ------------------------------ |
| `GET`    | `/api/v1/database-deployments`                                   | `a8s database list`                               | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments`                                   | `a8s database deploy`                             | `DatabaseDeploymentController` |
| `DELETE` | `/api/v1/database-deployments/{deploymentId}`                    | `a8s database delete <deployment-id>`             | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}`                    | `a8s database get <deployment-id>`                | `DatabaseDeploymentController` |
| `PATCH`  | `/api/v1/database-deployments/{deploymentId}`                    | `a8s database update <deployment-id>`             | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/console/data`       | `a8s database console data <deployment-id>`       | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/console/namespaces` | `a8s database console namespaces <deployment-id>` | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/console/objects`    | `a8s database console objects <deployment-id>`    | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/console/query`      | `a8s database console query <deployment-id>`      | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/console/test`       | `a8s database console test <deployment-id>`       | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/credentials`        | `a8s database credentials <deployment-id>`        | `DatabaseDeploymentController` |
| `GET`    | `/api/v1/database-deployments/{deploymentId}/metrics`            | `a8s database metrics <deployment-id>`            | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/restart`            | `a8s database restart <deployment-id>`            | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/rotate-password`    | `a8s database rotate-password <deployment-id>`    | `DatabaseDeploymentController` |
| `PATCH`  | `/api/v1/database-deployments/{deploymentId}/settings`           | `a8s database settings update <deployment-id>`    | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/upgrade-version`    | `a8s database upgrade <deployment-id>`            | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/{deploymentId}/verify-password`    | `a8s database verify-password <deployment-id>`    | `DatabaseDeploymentController` |
| `POST`   | `/api/v1/database-deployments/clone-from-backup`                 | `a8s database clone-from-backup`                  | `DatabaseDeploymentController` |

## sonarqube

| Method | Endpoint                                        | Suggested CLI command                | Controller                   |
| ------ | ----------------------------------------------- | ------------------------------------ | ---------------------------- |
| `GET`  | `/api/v1/projects/{projectId}/sonarqube`        | `a8s sonarqube summary <project-id>` | `SonarQubeProjectController` |
| `POST` | `/api/v1/projects/{projectId}/sonarqube/access` | `a8s sonarqube access <project-id>`  | `SonarQubeProjectController` |

## testingkit

| Method   | Endpoint                                                   | Suggested CLI command                        | Controller            |
| -------- | ---------------------------------------------------------- | -------------------------------------------- | --------------------- |
| `POST`   | `/api/v1/projects/live/{projectId}/benchmark/run`          | `a8s benchmark run <project-id>`             | `BenchmarkController` |
| `GET`    | `/api/v1/projects/live/{projectId}/benchmark/runs`         | `a8s benchmark list <project-id>`            | `BenchmarkController` |
| `DELETE` | `/api/v1/projects/live/{projectId}/benchmark/runs/{runId}` | `a8s benchmark delete <project-id> <run-id>` | `BenchmarkController` |
| `GET`    | `/api/v1/projects/live/{projectId}/benchmark/runs/{runId}` | `a8s benchmark get <project-id> <run-id>`    | `BenchmarkController` |

## workspaces

| Method | Endpoint                                           | Suggested CLI command                      | Controller            |
| ------ | -------------------------------------------------- | ------------------------------------------ | --------------------- |
| `GET`  | `/api/v1/workspaces/bootstrap`                     | `a8s workspace status`                     | `WorkspaceController` |
| `POST` | `/api/v1/workspaces/bootstrap`                     | `a8s workspace bootstrap`                  | `WorkspaceController` |
| `GET`  | `/api/v1/workspaces/quota-pricing`                 | `a8s workspace quota pricing`              | `WorkspaceController` |
| `POST` | `/api/v1/workspaces/quota-requests`                | `a8s workspace quota request`              | `WorkspaceController` |
| `GET`  | `/api/v1/workspaces/quota-requests/payment-status` | `a8s workspace quota payment-status <md5>` | `WorkspaceController` |

## WebSockets

| Endpoint                  | Suggested CLI use           |
| ------------------------- | --------------------------- |
| `/ws/jenkins/logs`        | `a8s project logs --follow` |
| `/ws/notifications`       | `a8s notification watch`    |
| `/ws/monitoring/overview` | `a8s monitoring watch`      |
| `/ws/admin/events`        | `a8s admin events watch`    |

## CLI Exclusions

Do not expose provider webhook receivers, Jenkins completion/failure callbacks, or `/api/internal/**` endpoints as ordinary user commands. They are automation-to-backend contracts.

## CLI Command Inventory (generated)

The following table is generated from the built CLI (`a8s list all -o json`) and shows every runnable command discovered in the current build. The `Input` column infers whether a command accepts an operation/manifest file via `--file` (marked `flags, file`) or is a read/flag-only command (`flags`). Use `a8s <command> --help` for exact accepted flags and file usage.

> Regenerate: `./dist/a8s list all -o json > docs/cli-command-catalog.json` then run the repo script to convert to Markdown.

<!-- BEGIN GENERATED COMMAND CATALOG -->

| Command                                   | Usage                                                               | Description                                                                             | Input       |
| ----------------------------------------- | ------------------------------------------------------------------- | --------------------------------------------------------------------------------------- | ----------- |
| a8s admin cluster health                  | a8s admin cluster health <alias> [flags]                            | GET /api/v1/admin/clusters/kubernetes/{alias}/health                                    | flags       |
| a8s admin cluster list                    | a8s admin cluster list [flags]                                      | GET /api/v1/admin/clusters                                                              | flags       |
| a8s admin cluster nodes                   | a8s admin cluster nodes [flags]                                     | GET /api/v1/admin/clusters/kubernetes                                                   | flags       |
| a8s admin cluster quota delete            | a8s admin cluster quota delete <alias> <namespace> [flags]          | DELETE /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}                     | flags       |
| a8s admin cluster quota list              | a8s admin cluster quota list <alias> [flags]                        | GET /api/v1/admin/clusters/kubernetes/{alias}/quotas                                    | flags       |
| a8s admin cluster quota set               | a8s admin cluster quota set <alias> <namespace> [flags]             | PUT /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}                        | flags, file |
| a8s admin cluster update                  | a8s admin cluster update <cluster-id> [flags]                       | PATCH /api/v1/admin/clusters/{clusterId}                                                | flags, file |
| a8s admin docs delete                     | a8s admin docs delete [flags]                                       | DELETE /api/admin/documentation/content                                                 | flags       |
| a8s admin docs files                      | a8s admin docs files [flags]                                        | GET /api/admin/documentation/files                                                      | flags, file |
| a8s admin docs get                        | a8s admin docs get [flags]                                          | GET /api/admin/documentation/content                                                    | flags       |
| a8s admin docs publish                    | a8s admin docs publish [flags]                                      | POST /api/admin/documentation/publish                                                   | flags       |
| a8s admin docs publish-logs               | a8s admin docs publish-logs [flags]                                 | GET /api/admin/documentation/publish/logs                                               | flags       |
| a8s admin docs update                     | a8s admin docs update [flags]                                       | PUT /api/admin/documentation/content                                                    | flags, file |
| a8s admin events watch                    | a8s admin events watch                                              | Watch administrative events                                                             | flags       |
| a8s admin gitops app abort                | a8s admin gitops app abort <app-id> [flags]                         | POST /api/v1/admin/gitops/apps/{appId}/abort                                            | flags       |
| a8s admin gitops app create               | a8s admin gitops app create [flags]                                 | POST /api/v1/admin/gitops/apps                                                          | flags, file |
| a8s admin gitops app retry                | a8s admin gitops app retry <app-id> [flags]                         | POST /api/v1/admin/gitops/apps/{appId}/retry                                            | flags       |
| a8s admin gitops app sync                 | a8s admin gitops app sync <app-id> [flags]                          | POST /api/v1/admin/gitops/apps/{appId}/sync                                             | flags       |
| a8s admin gitops overview                 | a8s admin gitops overview [flags]                                   | GET /api/v1/admin/gitops/overview                                                       | flags       |
| a8s admin logs clusters                   | a8s admin logs clusters [flags]                                     | GET /api/v1/admin/logs/clusters                                                         | flags       |
| a8s admin logs namespaces                 | a8s admin logs namespaces [flags]                                   | GET /api/v1/admin/logs/namespaces                                                       | flags       |
| a8s admin logs pods                       | a8s admin logs pods [flags]                                         | GET /api/v1/admin/logs/pods                                                             | flags       |
| a8s admin logs query                      | a8s admin logs query [flags]                                        | GET /api/v1/admin/logs/query                                                            | flags       |
| a8s admin logs workloads                  | a8s admin logs workloads [flags]                                    | GET /api/v1/admin/logs/workloads                                                        | flags       |
| a8s admin monitoring overview             | a8s admin monitoring overview [flags]                               | GET /api/v1/admin/monitoring/overview                                                   | flags       |
| a8s admin project deactivate              | a8s admin project deactivate <project-id> [flags]                   | DELETE /api/v1/admin/projects/{projectId}                                               | flags       |
| a8s admin project list                    | a8s admin project list [flags]                                      | GET /api/v1/admin/projects                                                              | flags       |
| a8s admin project restore                 | a8s admin project restore <project-id> [flags]                      | POST /api/v1/admin/projects/{projectId}/restore                                         | flags, file |
| a8s admin project update                  | a8s admin project update <project-id> [flags]                       | PATCH /api/v1/admin/projects/{projectId}                                                | flags, file |
| a8s admin quota approve                   | a8s admin quota approve <request-id> [flags]                        | POST /api/v1/admin/quota-requests/{id}/approve                                          | flags       |
| a8s admin quota list                      | a8s admin quota list [flags]                                        | GET /api/v1/admin/quota-requests                                                        | flags       |
| a8s admin quota reject                    | a8s admin quota reject <request-id> [flags]                         | POST /api/v1/admin/quota-requests/{id}/reject                                           | flags       |
| a8s admin registry artifact delete        | a8s admin registry artifact delete <project-name> [flags]           | DELETE /api/v1/admin/registry/projects/{projectName}/artifacts                          | flags       |
| a8s admin registry artifact list          | a8s admin registry artifact list <project-name> [flags]             | GET /api/v1/admin/registry/projects/{projectName}/artifacts                             | flags       |
| a8s admin registry health                 | a8s admin registry health [flags]                                   | GET /api/v1/admin/registry/health                                                       | flags       |
| a8s admin registry project create         | a8s admin registry project create [flags]                           | POST /api/v1/admin/registry/projects                                                    | flags, file |
| a8s admin registry project list           | a8s admin registry project list [flags]                             | GET /api/v1/admin/registry/projects                                                     | flags       |
| a8s admin registry repository delete      | a8s admin registry repository delete <project-name> [flags]         | DELETE /api/v1/admin/registry/projects/{projectName}/repositories                       | flags       |
| a8s admin registry repository list        | a8s admin registry repository list <project-name> [flags]           | GET /api/v1/admin/registry/projects/{projectName}/repositories                          | flags       |
| a8s admin sonarqube project get           | a8s admin sonarqube project get <project-id> [flags]                | GET /api/v1/admin/sonarqube/projects/{projectId}                                        | flags       |
| a8s admin sonarqube project list          | a8s admin sonarqube project list [flags]                            | GET /api/v1/admin/sonarqube/projects                                                    | flags       |
| a8s admin sonarqube server-project create | a8s admin sonarqube server-project create [flags]                   | POST /api/v1/admin/sonarqube/server-projects                                            | flags, file |
| a8s admin sonarqube server-project delete | a8s admin sonarqube server-project delete <project-key> [flags]     | DELETE /api/v1/admin/sonarqube/server-projects/{projectKey}                             | flags       |
| a8s admin sonarqube server-project get    | a8s admin sonarqube server-project get <project-key> [flags]        | GET /api/v1/admin/sonarqube/server-projects/{projectKey}                                | flags       |
| a8s admin sonarqube server-project list   | a8s admin sonarqube server-project list [flags]                     | GET /api/v1/admin/sonarqube/server-projects                                             | flags       |
| a8s admin sonarqube server-project update | a8s admin sonarqube server-project update <project-key> [flags]     | PATCH /api/v1/admin/sonarqube/server-projects/{projectKey}                              | flags, file |
| a8s admin user create                     | a8s admin user create [flags]                                       | POST /api/v1/admin/users                                                                | flags, file |
| a8s admin user deactivate                 | a8s admin user deactivate <user-id> [flags]                         | DELETE /api/v1/admin/users/{userId}                                                     | flags       |
| a8s admin user list                       | a8s admin user list [flags]                                         | GET /api/v1/admin/users                                                                 | flags       |
| a8s admin user reactivate                 | a8s admin user reactivate <user-id> [flags]                         | POST /api/v1/admin/users/{userId}/reactivate                                            | flags       |
| a8s admin user update                     | a8s admin user update <user-id> [flags]                             | PATCH /api/v1/admin/users/{userId}                                                      | flags, file |
| a8s alert channel create                  | a8s alert channel create [flags]                                    | POST /api/v1/alerts/channels                                                            | flags, file |
| a8s alert channel delete                  | a8s alert channel delete <channel-id> [flags]                       | DELETE /api/v1/alerts/channels/{channelId}                                              | flags       |
| a8s alert channel list                    | a8s alert channel list [flags]                                      | GET /api/v1/alerts/channels                                                             | flags       |
| a8s alert channel update                  | a8s alert channel update <channel-id> [flags]                       | PUT /api/v1/alerts/channels/{channelId}                                                 | flags, file |
| a8s alert project-config get              | a8s alert project-config get <project-id> [flags]                   | GET /api/v1/alerts/projects/{projectId}/config                                          | flags       |
| a8s alert project-config list             | a8s alert project-config list [flags]                               | GET /api/v1/alerts/projects/configs                                                     | flags       |
| a8s alert project-config set              | a8s alert project-config set <project-id> [flags]                   | PUT /api/v1/alerts/projects/{projectId}/config                                          | flags, file |
| a8s alert user-config get                 | a8s alert user-config get [flags]                                   | GET /api/v1/alerts/user-config                                                          | flags       |
| a8s alert user-config set                 | a8s alert user-config set [flags]                                   | PUT /api/v1/alerts/user-config                                                          | flags, file |
| a8s api catalog                           | a8s api catalog [flags]                                             | List implemented backend route mappings                                                 | flags       |
| a8s api request                           | a8s api request <method> <path> [flags]                             | Send an authenticated request to any backend route                                      | flags       |
| a8s auth login                            | a8s auth login [flags]                                              | Authenticate through Keycloak using browser PKCE                                        | flags       |
| a8s auth logout                           | a8s auth logout [flags]                                             | Clear stored credentials for the active context                                         | flags       |
| a8s auth onboarding start                 | a8s auth onboarding start [flags]                                   | POST /api/v1/auth/session/onboarding                                                    | flags       |
| a8s auth onboarding status                | a8s auth onboarding status [flags]                                  | GET /api/v1/auth/session/onboarding                                                     | flags       |
| a8s auth status                           | a8s auth status                                                     | Show authentication status without displaying tokens                                    | flags       |
| a8s auth verify-email start               | a8s auth verify-email start [keycloak-user-id]                      | Start email verification for the authenticated user                                     | flags       |
| a8s auth verify-email status              | a8s auth verify-email status [keycloak-user-id]                     | Show email verification status                                                          | flags       |
| a8s backup delete                         | a8s backup delete <type> <id> <run-id> [flags]                      | DELETE /api/backups/{targetType}/{id}/{runId}                                           | flags       |
| a8s backup download                       | a8s backup download <type> <id> <run-id> [flags]                    | GET /api/backups/download/{targetType}/{id}/{runId}                                     | flags, file |
| a8s backup restore                        | a8s backup restore <type> <id> <run-id> [flags]                     | POST /api/backups/restore/{targetType}/{id}/{runId}                                     | flags, file |
| a8s backup restore cancel                 | a8s backup restore cancel <type> <id> <run-id> [flags]              | POST /api/backups/restore/{targetType}/{id}/{runId}/cancel                              | flags, file |
| a8s backup settings get                   | a8s backup settings get <type> <id> [flags]                         | GET /api/backups/settings/{targetType}/{id}                                             | flags, file |
| a8s backup settings set                   | a8s backup settings set <type> <id> [flags]                         | POST /api/backups/settings/{targetType}/{id}                                            | flags, file |
| a8s backup trigger                        | a8s backup trigger <type> <id> [flags]                              | POST /api/backups/trigger/{targetType}/{id}                                             | flags       |
| a8s benchmark delete                      | a8s benchmark delete <project-id> <run-id> [flags]                  | DELETE /api/v1/projects/live/{projectId}/benchmark/runs/{runId}                         | flags       |
| a8s benchmark get                         | a8s benchmark get <project-id> <run-id> [flags]                     | GET /api/v1/projects/live/{projectId}/benchmark/runs/{runId}                            | flags       |
| a8s benchmark list                        | a8s benchmark list <project-id> [flags]                             | GET /api/v1/projects/live/{projectId}/benchmark/runs                                    | flags       |
| a8s benchmark run                         | a8s benchmark run <project-id> [flags]                              | POST /api/v1/projects/live/{projectId}/benchmark/run                                    | flags       |
| a8s cluster backup settings set           | a8s cluster backup settings set <release-name> [flags]              | PATCH /api/namespaces/{namespace}/cluster-deployments/{releaseName}/backup              | flags, file |
| a8s cluster certificate                   | a8s cluster certificate <cluster-id> [flags]                        | GET /api/namespaces/{namespace}/clusters/{id}/certificate                               | flags       |
| a8s cluster clone-from-backup             | a8s cluster clone-from-backup [flags]                               | POST /api/namespaces/{namespace}/clusters/clone-from-backup                             | flags       |
| a8s cluster console credentials           | a8s cluster console credentials <cluster-id> [flags]                | GET /api/namespaces/{namespace}/clusters/{id}/console/credentials                       | flags       |
| a8s cluster console data                  | a8s cluster console data <cluster-id> [flags]                       | GET /api/namespaces/{namespace}/clusters/{id}/console/data                              | flags       |
| a8s cluster console deployment            | a8s cluster console deployment <cluster-id> [flags]                 | GET /api/namespaces/{namespace}/clusters/{id}/console/deployment                        | flags, file |
| a8s cluster console namespaces            | a8s cluster console namespaces <cluster-id> [flags]                 | GET /api/namespaces/{namespace}/clusters/{id}/console/namespaces                        | flags       |
| a8s cluster console objects               | a8s cluster console objects <cluster-id> [flags]                    | GET /api/namespaces/{namespace}/clusters/{id}/console/objects                           | flags       |
| a8s cluster console query                 | a8s cluster console query <cluster-id> [flags]                      | POST /api/namespaces/{namespace}/clusters/{id}/console/query                            | flags       |
| a8s cluster console test                  | a8s cluster console test <cluster-id> [flags]                       | POST /api/namespaces/{namespace}/clusters/{id}/console/test                             | flags       |
| a8s cluster delete                        | a8s cluster delete <cluster-id> [flags]                             | DELETE /api/namespaces/{namespace}/clusters/{id}                                        | flags       |
| a8s cluster deploy                        | a8s cluster deploy [flags]                                          | POST /api/namespaces/{namespace}/cluster-deployments                                    | flags, file |
| a8s cluster deployment values             | a8s cluster deployment values <release-name> [flags]                | GET /api/namespaces/{namespace}/cluster-deployments/{releaseName}/values                | flags, file |
| a8s cluster get                           | a8s cluster get <cluster-id> [flags]                                | GET /api/namespaces/{namespace}/clusters/{id}                                           | flags       |
| a8s cluster history                       | a8s cluster history <cluster-id> [flags]                            | GET /api/namespaces/{namespace}/clusters/{id}/deployments                               | flags, file |
| a8s cluster list                          | a8s cluster list [flags]                                            | GET /api/namespaces/{namespace}/clusters                                                | flags       |
| a8s cluster metrics                       | a8s cluster metrics <cluster-id> [flags]                            | GET /api/namespaces/{namespace}/clusters/{id}/metrics                                   | flags       |
| a8s cluster settings update               | a8s cluster settings update <cluster-id> [flags]                    | PATCH /api/namespaces/{namespace}/clusters/{id}/settings                                | flags, file |
| a8s cluster status                        | a8s cluster status <release-name> [flags]                           | GET /api/namespaces/{namespace}/cluster-deployments/{releaseName}                       | flags, file |
| a8s cluster update                        | a8s cluster update <cluster-id> [flags]                             | PATCH /api/namespaces/{namespace}/clusters/{id}                                         | flags, file |
| a8s cluster upgrade                       | a8s cluster upgrade <cluster-id> [flags]                            | POST /api/namespaces/{namespace}/clusters/{id}/upgrade-version                          | flags       |
| a8s cluster values                        | a8s cluster values <cluster-id> [flags]                             | GET /api/namespaces/{namespace}/clusters/{id}/values                                    | flags       |
| a8s cluster watch                         | a8s cluster watch <release-name> [flags]                            | GET /api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream     | flags, file |
| a8s completion bash                       | a8s completion bash                                                 | Generate the autocompletion script for bash                                             | flags       |
| a8s completion fish                       | a8s completion fish [flags]                                         | Generate the autocompletion script for fish                                             | flags       |
| a8s completion powershell                 | a8s completion powershell [flags]                                   | Generate the autocompletion script for powershell                                       | flags       |
| a8s completion zsh                        | a8s completion zsh [flags]                                          | Generate the autocompletion script for zsh                                              | flags       |
| a8s config path                           | a8s config path                                                     | Print the active configuration path                                                     | flags       |
| a8s config view                           | a8s config view                                                     | Print resolved non-secret configuration                                                 | flags       |
| a8s context create                        | a8s context create <name> [flags]                                   | Create a named context                                                                  | flags, file |
| a8s context delete                        | a8s context delete <name> [flags]                                   | Delete a named context                                                                  | flags       |
| a8s context get                           | a8s context get <name>                                              | Get a configured context                                                                | flags       |
| a8s context list                          | a8s context list                                                    | List configured contexts                                                                | flags       |
| a8s context update                        | a8s context update <name> [flags]                                   | Update a named context                                                                  | flags, file |
| a8s context use                           | a8s context use <name>                                              | Set the default context                                                                 | flags, file |
| a8s database backup delete                | a8s database backup delete <deployment-id> <run-id> [flags]         | DELETE /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}                  | flags, file |
| a8s database backup download              | a8s database backup download <deployment-id> <run-id> [flags]       | GET /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/download            | flags, file |
| a8s database backup restore               | a8s database backup restore <deployment-id> <run-id> [flags]        | POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore            | flags, file |
| a8s database backup restore cancel        | a8s database backup restore cancel <deployment-id> <run-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore/cancel     | flags, file |
| a8s database backup run                   | a8s database backup run <deployment-id> [flags]                     | POST /api/v1/database-deployments/{deploymentId}/backup/run                             | flags, file |
| a8s database backup settings get          | a8s database backup settings get <deployment-id> [flags]            | GET /api/v1/database-deployments/{deploymentId}/backup                                  | flags, file |
| a8s database backup settings set          | a8s database backup settings set <deployment-id> [flags]            | PATCH /api/v1/database-deployments/{deploymentId}/backup                                | flags, file |
| a8s database clone-from-backup            | a8s database clone-from-backup [flags]                              | POST /api/v1/database-deployments/clone-from-backup                                     | flags, file |
| a8s database console data                 | a8s database console data <deployment-id> [flags]                   | GET /api/v1/database-deployments/{deploymentId}/console/data                            | flags, file |
| a8s database console namespaces           | a8s database console namespaces <deployment-id> [flags]             | GET /api/v1/database-deployments/{deploymentId}/console/namespaces                      | flags, file |
| a8s database console objects              | a8s database console objects <deployment-id> [flags]                | GET /api/v1/database-deployments/{deploymentId}/console/objects                         | flags, file |
| a8s database console query                | a8s database console query <deployment-id> [flags]                  | POST /api/v1/database-deployments/{deploymentId}/console/query                          | flags, file |
| a8s database console test                 | a8s database console test <deployment-id> [flags]                   | POST /api/v1/database-deployments/{deploymentId}/console/test                           | flags, file |
| a8s database credentials                  | a8s database credentials <deployment-id> [flags]                    | GET /api/v1/database-deployments/{deploymentId}/credentials                             | flags, file |
| a8s database delete                       | a8s database delete <deployment-id> [flags]                         | DELETE /api/v1/database-deployments/{deploymentId}                                      | flags, file |
| a8s database deploy                       | a8s database deploy [flags]                                         | Deploy a single database using flags or an operation file                               | flags, file |
| a8s database get                          | a8s database get <deployment-id> [flags]                            | GET /api/v1/database-deployments/{deploymentId}                                         | flags, file |
| a8s database list                         | a8s database list [flags]                                           | GET /api/v1/database-deployments                                                        | flags, file |
| a8s database metrics                      | a8s database metrics <deployment-id> [flags]                        | GET /api/v1/database-deployments/{deploymentId}/metrics                                 | flags, file |
| a8s database restart                      | a8s database restart <deployment-id> [flags]                        | POST /api/v1/database-deployments/{deploymentId}/restart                                | flags, file |
| a8s database rotate-password              | a8s database rotate-password <deployment-id> [flags]                | POST /api/v1/database-deployments/{deploymentId}/rotate-password                        | flags, file |
| a8s database settings update              | a8s database settings update <deployment-id> [flags]                | PATCH /api/v1/database-deployments/{deploymentId}/settings                              | flags, file |
| a8s database update                       | a8s database update <deployment-id> [flags]                         | PATCH /api/v1/database-deployments/{deploymentId}                                       | flags, file |
| a8s database upgrade                      | a8s database upgrade <deployment-id> [flags]                        | POST /api/v1/database-deployments/{deploymentId}/upgrade-version                        | flags, file |
| a8s database verify-password              | a8s database verify-password <deployment-id> [flags]                | POST /api/v1/database-deployments/{deploymentId}/verify-password                        | flags, file |
| a8s defectdojo access                     | a8s defectdojo access <project-id> [flags]                          | GET /api/v1/projects/{projectId}/defectdojo                                             | flags       |
| a8s defectdojo token sync                 | a8s defectdojo token sync <project-id> [flags]                      | PUT /api/v1/projects/{projectId}/defectdojo/token                                       | flags       |
| a8s doctor                                | a8s doctor                                                          | Check CLI configuration and backend connectivity                                        | flags       |
| a8s features                              | a8s features                                                        | List backend features exposed by the CLI                                                | flags       |
| a8s git account                           | a8s git account <provider> [flags]                                  | GET /api/v1/git-integrations/{provider}/brokered-account                                | flags       |
| a8s git connect                           | a8s git connect <provider> [flags]                                  | POST /api/v1/git-integrations/{provider}/connect                                        | flags       |
| a8s git disconnect                        | a8s git disconnect <provider> [flags]                               | DELETE /api/v1/git-integrations/{provider}                                              | flags       |
| a8s git providers                         | a8s git providers [flags]                                           | GET /api/v1/git-integrations/linked-providers                                           | flags       |
| a8s git repos                             | a8s git repos <provider> [flags]                                    | GET /api/v1/git-integrations/{provider}/repos                                           | flags       |
| a8s git state                             | a8s git state <provider> [flags]                                    | GET /api/v1/git-integrations/{provider}/state                                           | flags       |
| a8s git sync-token                        | a8s git sync-token <provider> [flags]                               | POST /api/v1/git-integrations/{provider}/sync-keycloak-token                            | flags       |
| a8s help                                  | a8s help [command]                                                  | Help about any command                                                                  | flags       |
| a8s kubernetes database-resources         | a8s kubernetes database-resources [flags]                           | GET /api/kubernetes/namespaces/{namespace}/database-resources                           | flags       |
| a8s kubernetes events                     | a8s kubernetes events [flags]                                       | GET /api/kubernetes/namespaces/{namespace}/events                                       | flags       |
| a8s kubernetes overview                   | a8s kubernetes overview [flags]                                     | GET /api/kubernetes/namespaces/{namespace}/overview                                     | flags       |
| a8s kubernetes pods                       | a8s kubernetes pods [flags]                                         | GET /api/kubernetes/namespaces/{namespace}/pods                                         | flags       |
| a8s kubernetes pvc                        | a8s kubernetes pvc [flags]                                          | GET /api/kubernetes/namespaces/{namespace}/persistent-volume-claims                     | flags       |
| a8s kubernetes services                   | a8s kubernetes services [flags]                                     | GET /api/kubernetes/namespaces/{namespace}/services                                     | flags       |
| a8s kubernetes test                       | a8s kubernetes test [flags]                                         | GET /api/kubernetes/test                                                                | flags       |
| a8s list all                              | a8s list all [flags]                                                | List every available runnable command                                                   | flags       |
| a8s list sections                         | a8s list sections [flags]                                           | List commands grouped by top-level section                                              | flags       |
| a8s logs                                  | a8s logs <pod-name> [flags]                                         | GET /api/kubernetes/namespaces/{namespace}/pods/{podName}/logs/stream                   | flags       |
| a8s manifest init                         | a8s manifest init <kind> [flags]                                    | Generate a starter manifest for a kind                                                  | flags, file |
| a8s manifest kinds                        | a8s manifest kinds                                                  | List supported operation manifest kinds                                                 | flags, file |
| a8s manifest schema                       | a8s manifest schema <kind>                                          | Show the manifest schema summary for a kind                                             | flags, file |
| a8s manifest validate                     | a8s manifest validate [flags]                                       | Validate an operation manifest without sending a backend request                        | flags, file |
| a8s microservice apply                    | a8s microservice apply <project-id> [flags]                         | PUT /api/v1/projects/microservices/{projectId}/canvas                                   | flags, file |
| a8s microservice delete                   | a8s microservice delete <project-id> [flags]                        | DELETE /api/v1/projects/microservices/{projectId}                                       | flags       |
| a8s microservice deploy                   | a8s microservice deploy [flags]                                     | POST /api/v1/projects/microservices                                                     | flags, file |
| a8s microservice detect                   | a8s microservice detect [flags]                                     | POST /api/v1/projects/microservices/detect                                              | flags       |
| a8s microservice domains update           | a8s microservice domains update <project-id> [flags]                | PATCH /api/v1/projects/microservices/{projectId}/domains                                | flags, file |
| a8s microservice env clear                | a8s microservice env clear <project-id> <service-id> [flags]        | DELETE /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment      | flags       |
| a8s microservice env get                  | a8s microservice env get <project-id> <service-id> [flags]          | GET /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment         | flags       |
| a8s microservice env import               | a8s microservice env import <project-id> <service-id> [flags]       | POST /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment/import | flags, file |
| a8s microservice env set                  | a8s microservice env set <project-id> <service-id> [flags]          | PUT /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment         | flags, file |
| a8s microservice get                      | a8s microservice get <project-id> [flags]                           | GET /api/v1/projects/microservices/{projectId}                                          | flags       |
| a8s microservice history delete           | a8s microservice history delete <project-id> <snapshot-id> [flags]  | DELETE /api/v1/projects/microservices/{projectId}/history/{snapshotId}                  | flags       |
| a8s microservice history list             | a8s microservice history list <project-id> [flags]                  | GET /api/v1/projects/microservices/{projectId}/history                                  | flags       |
| a8s microservice pods                     | a8s microservice pods <project-id> [flags]                          | GET /api/v1/projects/microservices/{projectId}/runtime-pods                             | flags       |
| a8s microservice readiness                | a8s microservice readiness <project-id> [flags]                     | GET /api/v1/projects/microservices/{projectId}/readiness                                | flags       |
| a8s microservice redeploy                 | a8s microservice redeploy <project-id> [flags]                      | POST /api/v1/projects/microservices/{projectId}/redeploy                                | flags, file |
| a8s microservice rollback                 | a8s microservice rollback <project-id> [flags]                      | POST /api/v1/projects/microservices/{projectId}/rollback                                | flags       |
| a8s microservice webhook get              | a8s microservice webhook get <project-id> [flags]                   | GET /api/v1/projects/microservices/{projectId}/webhook                                  | flags       |
| a8s microservice webhook update           | a8s microservice webhook update <project-id> [flags]                | POST /api/v1/projects/microservices/{projectId}/webhook                                 | flags, file |
| a8s monitoring overview                   | a8s monitoring overview [flags]                                     | GET /api/v1/monitoring/overview                                                         | flags       |
| a8s monitoring watch                      | a8s monitoring watch                                                | Watch monitoring updates                                                                | flags, file |
| a8s notification list                     | a8s notification list <user-id> [flags]                             | GET /api/notifications/history/{userId}                                                 | flags       |
| a8s notification preferences get          | a8s notification preferences get <user-id> [flags]                  | GET /api/notifications/preferences/{userId}                                             | flags       |
| a8s notification preferences set          | a8s notification preferences set <user-id> [flags]                  | POST /api/notifications/preferences/{userId}                                            | flags, file |
| a8s notification read                     | a8s notification read <notification-id> [flags]                     | POST /api/notifications/{notificationId}/read                                           | flags       |
| a8s notification watch                    | a8s notification watch                                              | Watch notifications                                                                     | flags       |
| a8s profile account deactivate            | a8s profile account deactivate [flags]                              | POST /api/v1/profile/me/deactivate                                                      | flags, file |
| a8s profile account delete                | a8s profile account delete [flags]                                  | DELETE /api/v1/profile/me                                                               | flags, file |
| a8s profile account reactivate            | a8s profile account reactivate [flags]                              | POST /api/v1/profile/me/reactivate                                                      | flags, file |
| a8s profile account status                | a8s profile account status [flags]                                  | GET /api/v1/profile/me/account-status                                                   | flags, file |
| a8s profile avatar delete                 | a8s profile avatar delete [flags]                                   | DELETE /api/v1/profile/me/avatar                                                        | flags, file |
| a8s profile avatar download               | a8s profile avatar download [flags]                                 | GET /api/v1/profile/me/avatar                                                           | flags, file |
| a8s profile avatar upload                 | a8s profile avatar upload [flags]                                   | POST /api/v1/profile/me/avatar                                                          | flags, file |
| a8s profile get                           | a8s profile get [flags]                                             | GET /api/v1/profile/me                                                                  | flags, file |
| a8s profile update                        | a8s profile update [flags]                                          | PATCH /api/v1/profile/me                                                                | flags, file |
| a8s project auto-deploy set               | a8s project auto-deploy set <project-id> [flags]                    | PATCH /api/v1/projects/{projectId}/auto-deploy                                          | flags, file |
| a8s project branches                      | a8s project branches <project-id> [flags]                           | GET /api/v1/projects/{projectId}/branches                                               | flags       |
| a8s project delete                        | a8s project delete <project-id> [flags]                             | DELETE /api/v1/projects/{projectId}                                                     | flags       |
| a8s project deploy                        | a8s project deploy [flags]                                          | POST /api/v1/projects                                                                   | flags, file |
| a8s project domain set                    | a8s project domain set <project-id> [flags]                         | PATCH /api/v1/projects/{projectId}/domain                                               | flags, file |
| a8s project domain sync                   | a8s project domain sync <project-id> [flags]                        | POST /api/v1/projects/{projectId}/domain/sync                                           | flags       |
| a8s project env get                       | a8s project env get <project-id> [flags]                            | GET /api/v1/projects/{projectId}/environment                                            | flags       |
| a8s project env import                    | a8s project env import <project-id> [flags]                         | POST /api/v1/projects/{projectId}/environment/import                                    | flags, file |
| a8s project env set                       | a8s project env set <project-id> [flags]                            | PUT /api/v1/projects/{projectId}/environment                                            | flags, file |
| a8s project get                           | a8s project get <project-id> [flags]                                | GET /api/v1/projects/{projectId}                                                        | flags       |
| a8s project list                          | a8s project list [flags]                                            | GET /api/v1/projects                                                                    | flags       |
| a8s project live list                     | a8s project live list [flags]                                       | GET /api/v1/projects/live                                                               | flags       |
| a8s project logs                          | a8s project logs [flags]                                            | GET /api/v1/jenkins/logs/stream                                                         | flags       |
| a8s project logs websocket                | a8s project logs websocket                                          | Watch Jenkins logs over WebSocket                                                       | flags       |
| a8s project redeploy                      | a8s project redeploy <project-id> [flags]                           | POST /api/v1/projects/{projectId}/sync                                                  | flags, file |
| a8s project release delete                | a8s project release delete <project-id> <release-id> [flags]        | DELETE /api/v1/projects/{projectId}/releases/{releaseId}                                | flags       |
| a8s project release rollback              | a8s project release rollback <project-id> <release-id> [flags]      | POST /api/v1/projects/{projectId}/releases/{releaseId}/rollback                         | flags       |
| a8s project releases                      | a8s project releases <project-id> [flags]                           | GET /api/v1/projects/{projectId}/releases                                               | flags       |
| a8s project repository connect            | a8s project repository connect <project-id> [flags]                 | POST /api/v1/projects/{projectId}/repository/connect                                    | flags       |
| a8s project rollback                      | a8s project rollback <project-id> [flags]                           | POST /api/v1/projects/{projectId}/rollback                                              | flags       |
| a8s project settings update               | a8s project settings update <project-id> [flags]                    | PATCH /api/v1/projects/{projectId}/settings                                             | flags, file |
| a8s project webhook create                | a8s project webhook create <project-id> [flags]                     | POST /api/v1/projects/{projectId}/webhook                                               | flags, file |
| a8s project webhook delete                | a8s project webhook delete <project-id> [flags]                     | DELETE /api/v1/projects/{projectId}/webhook                                             | flags       |
| a8s project webhook get                   | a8s project webhook get <project-id> [flags]                        | GET /api/v1/projects/{projectId}/webhook                                                | flags       |
| a8s project webhook rotate                | a8s project webhook rotate <project-id> [flags]                     | POST /api/v1/projects/{projectId}/webhook/rotate                                        | flags       |
| a8s scan get                              | a8s scan get <scan-id> [flags]                                      | GET /api/v1/image-scanner/scans/{scanId}                                                | flags       |
| a8s scan images                           | a8s scan images [flags]                                             | GET /api/v1/image-scanner/images                                                        | flags       |
| a8s scan list                             | a8s scan list [flags]                                               | GET /api/v1/image-scanner/scans                                                         | flags       |
| a8s scan report                           | a8s scan report <scan-id> [flags]                                   | GET /api/v1/image-scanner/scans/{scanId}/report                                         | flags       |
| a8s scan start                            | a8s scan start [flags]                                              | POST /api/v1/image-scanner/scans                                                        | flags       |
| a8s sonarqube access                      | a8s sonarqube access <project-id> [flags]                           | POST /api/v1/projects/{projectId}/sonarqube/access                                      | flags       |
| a8s sonarqube summary                     | a8s sonarqube summary <project-id> [flags]                          | GET /api/v1/projects/{projectId}/sonarqube                                              | flags       |
| a8s version                               | a8s version                                                         | Print the CLI version                                                                   | flags       |
| a8s workspace bootstrap                   | a8s workspace bootstrap [flags]                                     | POST /api/v1/workspaces/bootstrap                                                       | flags       |
| a8s workspace entitlements                | a8s workspace entitlements [flags]                                  | GET /api/v1/workspaces/entitlements                                                     | flags       |
| a8s workspace quota payment-status        | a8s workspace quota payment-status <md5> [flags]                    | GET /api/v1/workspaces/quota-requests/payment-status                                    | flags       |
| a8s workspace quota pricing               | a8s workspace quota pricing [flags]                                 | GET /api/v1/workspaces/quota-pricing                                                    | flags       |
| a8s workspace quota purchase              | a8s workspace quota purchase [flags]                                | Purchase a workspace quota plan using Bakong KHQR                                       | flags       |
| a8s workspace quota request               | a8s workspace quota request [flags]                                 | POST /api/v1/workspaces/quota-requests                                                  | flags       |
| a8s workspace status                      | a8s workspace status [flags]                                        | GET /api/v1/workspaces/bootstrap                                                        | flags       |

<!-- END GENERATED COMMAND CATALOG -->
