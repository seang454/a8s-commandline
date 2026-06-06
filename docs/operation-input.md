# A8S CLI Operation Input

## Purpose

This document defines how users provide request data to A8S CLI commands.
Operation input is separate from the persistent CLI configuration described in
`configuration.md`.

- CLI configuration selects the backend, context, namespace, authentication,
  timeouts, and default output behavior.
- Operation input describes the resource mutation the user wants to perform.

## Production-Wide Rule

Every command that sends a configurable request payload must support both:

1. YAML or JSON operation input using `--file`
2. equivalent command flags

This applies to large deployment requests and small one-field requests.

```bash
# File only
a8s database upgrade db-123 --file upgrade.yaml

# Flags only
a8s database upgrade db-123 --version 17

# File with an explicit flag override
a8s database upgrade db-123 --file upgrade.yaml --version 18

# YAML or JSON from stdin
a8s database upgrade db-123 --file -
```

All forms must produce the same typed internal request and backend API payload.

## Input Resolution

Use this precedence:

```text
explicit flags > operation file > active-context defaults > backend defaults
```

The CLI must distinguish explicitly supplied flags from Cobra defaults.
Defaults that the user did not specify must not overwrite operation-file
values.

The processing flow is:

```text
read file -> decode typed input -> apply explicit flags -> apply context
defaults -> validate -> optionally dry-run -> send request
```

Reject:

- unknown YAML or JSON fields
- unsupported `apiVersion` or `kind`
- missing required values after merging
- mutually exclusive or conflicting inputs
- invalid formats, enum values, durations, sizes, URLs, or identifiers

## How to Read the Command Catalog

Every operation file uses this envelope:

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectDomain
spec:
  customDomain: api.example.com
```

The tables below show only the `spec` fields to keep the examples readable.
Create a complete file by adding the envelope and using the listed kind.
The examples show commonly used fields. Each implemented kind must expose all
backend-supported request fields in its generated schema.

Equivalent repeatable flags use this convention:

```text
--env NAME=VALUE
--secret-env NAME
--header NAME=VALUE
--service-domain SERVICE_ID=DOMAIN
```

The CLI must document escaping rules and allow repeatable flags to be supplied
more than once.

## Complete Copy-Ready Examples

Small operation:

```yaml
# database-upgrade.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseUpgrade
spec:
  version: "17"
```

```bash
a8s database upgrade db-123 --file database-upgrade.yaml
a8s database upgrade db-123 --version 17
```

Monolithic project deployment:

```yaml
# project.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectDeployment
spec:
  projectName: shop-api
  sourceType: git
  repoUrl: https://github.com/acme/shop-api.git
  repoFullName: acme/shop-api
  branch: main
  appPort: 8080
  architectureType: monolithic
  autoDeployEnabled: true
  autoDeployTrigger: push
  envVars:
    - name: SPRING_PROFILES_ACTIVE
      value: production
      secret: false
```

```bash
a8s project deploy \
  --name shop-api \
  --source-type git \
  --repo-url https://github.com/acme/shop-api.git \
  --repo-full-name acme/shop-api \
  --branch main \
  --app-port 8080 \
  --architecture monolithic \
  --auto-deploy \
  --auto-deploy-trigger push \
  --env SPRING_PROFILES_ACTIVE=production \
  --wait
```

Database cluster deployment:

```yaml
# cluster.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterDeployment
spec:
  releaseName: orders-cluster
  projectName: orders
  cluster:
    name: orders
    environment: PRODUCTION
    platformConfig:
      targetClusterName: production-primary
  database:
    engine: POSTGRESQL
    enabled: true
    instances: 3
    storageSize: 100Gi
    version: "16"
    monitoringEnabled: true
    backup:
      enabled: true
      destinationPath: s3://backups/orders
      credentialSecret: backup-credentials
      retentionPolicy: 30d
      schedule: "0 0 * * *"
```

```bash
a8s cluster deploy \
  --release-name orders-cluster \
  --project-name orders \
  --name orders \
  --environment PRODUCTION \
  --target-cluster production-primary \
  --engine POSTGRESQL \
  --instances 3 \
  --storage-size 100Gi \
  --version 16 \
  --monitoring \
  --backup-enabled \
  --backup-destination s3://backups/orders \
  --backup-credential-secret backup-credentials \
  --backup-retention 30d \
  --backup-schedule "0 0 * * *" \
  --wait
```

Use PowerShell backticks instead of backslashes when splitting a command
across lines on Windows.

These commands and flags define the intended production CLI contract. Until a
command is implemented in Cobra, its example documents the required behavior
rather than an already available executable command.

## CLI-Local Context Mutations

Context operations modify local CLI configuration instead of calling the
backend, but they follow the same file-or-flags rule.

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s context create production`, `Context` | `{server: https://api.a8s.example.com, namespace: ns-team, targetCluster: primary}` | `--server https://api.a8s.example.com --namespace ns-team --target-cluster primary` |
| `a8s context update production`, `ContextPatch` | `{namespace: ns-prod, targetCluster: prod-primary}` | `--namespace ns-prod --target-cluster prod-primary` |

```bash
a8s context create production --file context.yaml
a8s context update production --file context-patch.yaml
```

## Workspace, Profile, and Integration Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s workspace quota request`, `WorkspaceQuotaRequest` | `{requestedCpu: "4", requestedMemory: 8Gi, requestedStorage: 100Gi, reason: "Production workload", isPaid: false, planName: Free}` | `--cpu 4 --memory 8Gi --storage 100Gi --reason "Production workload" --plan Free` |
| `a8s workspace quota purchase`, `WorkspaceQuotaPurchase` | `{requestedCpu: "8", requestedMemory: 16Gi, requestedStorage: 200Gi, reason: "Upgrade", isPaid: true, planName: Premium, paymentProvider: BAKONG}` | `--cpu 8 --memory 16Gi --storage 200Gi --reason Upgrade --plan Premium --payment-provider BAKONG` |
| `a8s profile update`, `ProfileUpdate` | `{personal: {firstName: Dara, lastName: Sok, displayName: Dara, email: dara@example.com, jobTitle: Engineer, department: Platform, bio: ""}, locale: {city: Phnom Penh, country: Cambodia, timezone: Asia/Phnom_Penh, language: en, dateFormat: yyyy-MM-dd}}` | `--first-name Dara --last-name Sok --display-name Dara --email dara@example.com --job-title Engineer --department Platform --city "Phnom Penh" --country Cambodia --timezone Asia/Phnom_Penh --language en --date-format yyyy-MM-dd` |
| `a8s git connect github`, `GitProviderConnection` | `{accessTokenFrom: {env: GITHUB_TOKEN}, accessLevel: repository, grantedScopes: "repo,read:user"}` | `--access-token-stdin --access-level repository --granted-scopes repo,read:user` |
| `a8s defectdojo token sync <project-id>`, `DefectDojoToken` | `{apiTokenFrom: {env: DEFECTDOJO_API_TOKEN}}` | `--api-token-stdin` |

```bash
a8s workspace quota purchase --file quota-purchase.yaml --wait
a8s profile update --file profile.yaml
a8s git connect github --file git-connection.yaml
```

Attachment files for quota requests are separate from operation input:

```bash
a8s workspace quota request --file quota.yaml --attachment evidence.pdf
```

## Monolithic Project Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s project deploy`, `ProjectDeployment` | `{projectName: shop-api, repoUrl: https://github.com/acme/shop-api.git, repoFullName: acme/shop-api, sourceType: git, branch: main, appPort: 8080, architectureType: monolithic, autoDeployEnabled: true, autoDeployTrigger: push}` | `--name shop-api --repo-url https://github.com/acme/shop-api.git --repo-full-name acme/shop-api --source-type git --branch main --app-port 8080 --architecture monolithic --auto-deploy --auto-deploy-trigger push` |
| `a8s project domain set <project-id>`, `ProjectDomain` | `{customDomain: api.example.com}` | `--domain api.example.com` |
| `a8s project repository connect <project-id>`, `ProjectRepositoryConnection` | `{repoProvider: github, repoUrl: https://github.com/acme/shop-api.git, repoFullName: acme/shop-api, branch: main, autoDeployEnabled: true, autoDeployTrigger: push}` | `--provider github --repo-url https://github.com/acme/shop-api.git --repo-full-name acme/shop-api --branch main --auto-deploy --auto-deploy-trigger push` |
| `a8s project settings update <project-id>`, `ProjectSettings` | `{alias: shop, operatorNote: "Critical service", failureAlerts: true, maintenanceMode: false, protectFromDelete: true}` | `--alias shop --operator-note "Critical service" --failure-alerts --maintenance-mode=false --protect-from-delete` |
| `a8s project env set <project-id>`, `ProjectEnvironment` | `{envVars: [{name: SPRING_PROFILES_ACTIVE, value: production, secret: false}]}` | `--env SPRING_PROFILES_ACTIVE=production` |
| `a8s project auto-deploy set <project-id>`, `ProjectAutoDeploy` | `{enabled: true, branch: main, autoDeployTrigger: push, releaseTagPattern: "v*"}` | `--enabled --branch main --trigger push --release-tag-pattern "v*"` |
| `a8s project webhook create <project-id>`, `ProjectWebhook` | `{name: shop-webhook, branch: main, autoDeployEnabled: true, autoDeployTrigger: push, createOnProvider: true}` | `--name shop-webhook --branch main --auto-deploy --trigger push --create-on-provider` |
| `a8s project rollback <project-id>`, `ProjectRollback` | `{releaseId: 11111111-1111-1111-1111-111111111111}` | `--release-id 11111111-1111-1111-1111-111111111111` |
| `a8s project release rollback <project-id> <release-id>`, `ProjectReleaseRollback` | `{buildNumber: 42, framework: spring, statusMessage: "Operator rollback"}` | `--build-number 42 --framework spring --status-message "Operator rollback"` |

```bash
a8s project deploy --file project.yaml --wait
a8s project deploy --file project.yaml --branch release --wait
a8s project domain set project-123 --file domain.yaml
```

The CLI must derive `userId` from the authenticated session. Users must not
provide backend ownership IDs for ordinary project deployment.

## Microservice Project Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s microservice deploy`, `MicroserviceDeployment` | `{projectName: shop, branch: main, services: [{name: api, repoUrl: https://github.com/acme/shop.git, repoFullName: acme/shop, path: services/api, appPort: 8080, serviceType: backend, exposePublic: true}]}` | `--project-name shop --branch main --service-file api-service.yaml` |
| `a8s microservice apply <project-id>`, `MicroserviceCanvas` | `{branch: main, services: [{name: api, repoUrl: https://github.com/acme/shop.git, repoFullName: acme/shop, path: services/api, appPort: 8080}]}` | `--branch main --service-file api-service.yaml` |
| `a8s microservice domains update <project-id>`, `MicroserviceDomains` | `{services: [{serviceId: 11111111-1111-1111-1111-111111111111, customDomain: api.example.com, platformSubdomain: api}]}` | `--service-domain 11111111-1111-1111-1111-111111111111=api.example.com --service-subdomain 11111111-1111-1111-1111-111111111111=api` |
| `a8s microservice rollback <project-id>`, `MicroserviceRollback` | `{snapshotId: snapshot-123}` | `--snapshot-id snapshot-123` |
| `a8s microservice env set <project-id> <service-id>`, `MicroserviceEnvironment` | `{envVars: [{name: SPRING_PROFILES_ACTIVE, value: production, secret: false}], runtimeConfigFile: {fileName: application.yaml, content: "server:\n  port: 8080"}}` | `--env SPRING_PROFILES_ACTIVE=production --runtime-config-file application.yaml` |
| `a8s microservice webhook update <project-id>`, `MicroserviceWebhook` | `{name: shop-webhook, branch: main, autoDeployEnabled: true, autoDeployTrigger: push, releaseTagPattern: "v*", releaseTriggerMode: tag}` | `--name shop-webhook --branch main --auto-deploy --trigger push --release-tag-pattern "v*" --release-trigger-mode tag` |
| `a8s microservice detect --repo`, `MicroserviceDetection` | `{repoUrl: https://github.com/acme/shop.git, branch: main, githubTokenFrom: {env: GITHUB_TOKEN}}` | `--repo https://github.com/acme/shop.git --branch main --github-token-stdin` |

Complex service definitions should use `--service-file` rather than dozens of
flattened service flags. The top-level operation still supports both forms.

## Single Database Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s database deploy`, `DatabaseDeployment` | `{releaseName: payments-db, projectName: payments, engine: postgresql, deploymentMode: single, databaseName: payments, username: app, version: "16", sizeProfile: small, storageSize: 20Gi, networkPolicyEnabled: true, tls: {enabled: true, requireSsl: true}}` | `--release-name payments-db --project-name payments --engine postgresql --deployment-mode single --database-name payments --username app --version 16 --size-profile small --storage-size 20Gi --network-policy --tls --require-ssl` |
| `a8s database update <deployment-id>`, `DatabaseDeploymentPatch` | `{sizeProfile: medium, networkPolicyEnabled: true, tls: {enabled: true, requireSsl: true}}` | `--size-profile medium --network-policy --tls --require-ssl` |
| `a8s database settings update <deployment-id>`, `DatabaseSettings` | `{alias: payments, operatorNote: "Primary database", failureAlerts: true, maintenanceMode: false, protectFromDelete: true}` | `--alias payments --operator-note "Primary database" --failure-alerts --maintenance-mode=false --protect-from-delete` |
| `a8s database upgrade <deployment-id>`, `DatabaseUpgrade` | `{version: "17"}` | `--version 17` |
| `a8s database clone-from-backup`, `DatabaseClone` | `{sourceDeploymentId: 11111111-1111-1111-1111-111111111111, backupRunId: 22222222-2222-2222-2222-222222222222, projectName: payments-clone, databaseName: payments, version: "16", storageSize: 20Gi}` | `--source-deployment-id 11111111-1111-1111-1111-111111111111 --backup-run-id 22222222-2222-2222-2222-222222222222 --project-name payments-clone --database-name payments --version 16 --storage-size 20Gi` |
| `a8s database rotate-password <deployment-id>`, `DatabasePasswordRotation` | `{passwordFrom: {env: A8S_NEW_DATABASE_PASSWORD}}` | `--password-stdin` |
| `a8s database verify-password <deployment-id>`, `DatabasePasswordVerification` | `{passwordFrom: {env: A8S_DATABASE_PASSWORD}}` | `--password-stdin` |
| `a8s database console query <deployment-id>`, `DatabaseQuery` | `{query: "SELECT now()"}` | `--query "SELECT now()"` |
| `a8s database backup settings set <deployment-id>`, `DatabaseBackupSettings` | `{enabled: true, destinationPath: s3://backups/payments, credentialSecret: backup-credentials, retentionPolicy: 30d, schedule: "0 0 * * *", scheduleStartAt: "2026-06-06T00:00:00Z"}` | `--enabled --destination s3://backups/payments --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *" --schedule-start-at 2026-06-06T00:00:00Z` |

Secret values shown as `passwordFrom` are CLI-level secure references. The CLI
resolves them immediately before mapping to the current backend DTO.

## Database Cluster Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s cluster deploy`, `ClusterDeployment` | `{releaseName: orders-cluster, projectName: orders, cluster: {name: orders, environment: PRODUCTION}, database: {engine: POSTGRESQL, enabled: true, instances: 3, storageSize: 100Gi, version: "16", monitoringEnabled: true}}` | `--release-name orders-cluster --project-name orders --name orders --environment PRODUCTION --engine POSTGRESQL --instances 3 --storage-size 100Gi --version 16 --monitoring` |
| `a8s cluster update <cluster-id>`, `ClusterDeploymentPatch` | `{database: {instances: 5, storageSize: 200Gi, monitoringEnabled: true}}` | `--instances 5 --storage-size 200Gi --monitoring` |
| `a8s cluster settings update <cluster-id>`, `ClusterSettings` | `{alias: orders, operatorNote: "Primary cluster", failureAlerts: true, maintenanceMode: false, protectFromDelete: true}` | `--alias orders --operator-note "Primary cluster" --failure-alerts --maintenance-mode=false --protect-from-delete` |
| `a8s cluster upgrade <cluster-id>`, `ClusterUpgrade` | `{version: "17"}` | `--version 17` |
| `a8s cluster clone-from-backup`, `ClusterClone` | `{sourceClusterId: 11111111-1111-1111-1111-111111111111, backupRunId: 22222222-2222-2222-2222-222222222222, projectName: orders-clone, version: "16", instances: 3, storageSize: 100Gi, monitoringEnabled: true}` | `--source-cluster-id 11111111-1111-1111-1111-111111111111 --backup-run-id 22222222-2222-2222-2222-222222222222 --project-name orders-clone --version 16 --instances 3 --storage-size 100Gi --monitoring` |
| `a8s cluster backup settings set <cluster-id>`, `ClusterBackupSettings` | `{enabled: true, destinationPath: s3://backups/orders, credentialSecret: backup-credentials, retentionPolicy: 30d, schedule: "0 0 * * *"}` | `--enabled --destination s3://backups/orders --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *"` |
| `a8s cluster backup settings set --release <release-name>`, `ClusterBackupSettings` | `{enabled: true, destinationPath: s3://backups/orders, credentialSecret: backup-credentials, retentionPolicy: 30d, schedule: "0 0 * * *"}` | `--release orders-cluster --enabled --destination s3://backups/orders --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *"` |
| `a8s cluster console query <cluster-id>`, `ClusterQuery` | `{query: "SELECT now()"}` | `--query "SELECT now()"` |

Cluster secrets must use secret references or secret stdin flags. Do not place
database passwords or Cloudflare tokens directly in cluster operation files.

## Unified Backup, Quality, and Alert Mutations

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s backup settings set <type> <id>`, `BackupSettings` | `{enabled: true, destinationPath: s3://backups/resource, credentialSecret: backup-credentials, retentionPolicy: 30d, schedule: "0 0 * * *"}` | `--enabled --destination s3://backups/resource --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *"` |
| `a8s scan start`, `ImageScan` | `{sourceKind: image, imageRef: nginx:1.27, forceRescan: false}` | `--source-kind image --image nginx:1.27 --force-rescan=false` |
| `a8s benchmark run <project-id>`, `BenchmarkRun` | `{concurrency: 20, totalRequests: 1000, targetPath: /api/health, method: GET, headers: {Accept: application/json}}` | `--concurrency 20 --total-requests 1000 --target-path /api/health --method GET --header Accept=application/json` |
| `a8s alert channel create`, `AlertChannel` | `{name: operations, type: telegram, credentialFrom: {env: TELEGRAM_BOT_TOKEN}, secondaryCredentialFrom: {env: TELEGRAM_CHAT_ID}, targetProject: shop-api}` | `--name operations --type telegram --credential-env TELEGRAM_BOT_TOKEN --secondary-credential-env TELEGRAM_CHAT_ID --target-project shop-api` |
| `a8s alert channel update <channel-id>`, `AlertChannel` | `{name: operations, type: email, credential: ops@example.com, targetProject: shop-api}` | `--name operations --type email --credential ops@example.com --target-project shop-api` |
| `a8s alert project-config set <project-id>`, `ProjectAlertConfig` | `{telegramEnabled: true, emailEnabled: true, backupAlertsEnabled: true, securityAlertsEnabled: true, telegramChannelName: operations, emailAddress: ops@example.com}` | `--telegram-enabled --email-enabled --backup-alerts --security-alerts --telegram-channel operations --email-address ops@example.com` |
| `a8s alert user-config set`, `UserAlertConfig` | `{quotaAlertsEnabled: true, globalSecurityAlertsEnabled: true}` | `--quota-alerts --global-security-alerts` |
| `a8s notification preferences set`, `NotificationPreferences` | `{buildFailures: true, rolloutReady: true, vulnerabilityFindings: true, weeklyDigest: false}` | `--build-failures --rollout-ready --vulnerability-findings --weekly-digest=false` |

## Administrative Mutations

All commands in this section require backend admin authorization.

| Command and kind | YAML `spec` example | Equivalent flags |
|---|---|---|
| `a8s admin user create`, `AdminUserCreate` | `{username: dara, email: dara@example.com, firstName: Dara, lastName: Sok, passwordFrom: {env: A8S_INITIAL_PASSWORD}}` | `--username dara --email dara@example.com --first-name Dara --last-name Sok --password-stdin` |
| `a8s admin user update <user-id>`, `AdminUserUpdate` | `{username: dara, email: dara@example.com, firstName: Dara, lastName: Sok}` | `--username dara --email dara@example.com --first-name Dara --last-name Sok` |
| `a8s admin project update <project-id>`, `AdminProjectUpdate` | `{name: shop-api, userId: 11111111-1111-1111-1111-111111111111, repoProvider: github, repoFullName: acme/shop-api, repoUrl: https://github.com/acme/shop-api.git, branch: main, framework: spring, appPort: 8080, status: ACTIVE, autoDeployEnabled: true}` | `--name shop-api --user-id 11111111-1111-1111-1111-111111111111 --repo-provider github --repo-full-name acme/shop-api --repo-url https://github.com/acme/shop-api.git --branch main --framework spring --app-port 8080 --status ACTIVE --auto-deploy` |
| `a8s admin cluster update <cluster-id>`, `AdminClusterUpdate` | `{alias: primary, operatorNote: "Production cluster", failureAlerts: true, maintenanceMode: false, protectFromDelete: true}` | `--alias primary --operator-note "Production cluster" --failure-alerts --maintenance-mode=false --protect-from-delete` |
| `a8s admin cluster quota set <alias> <namespace>`, `AdminClusterQuota` | `{cpuLimit: 8, memoryLimit: 17179869184, pvcLimit: 20}` | `--cpu-limit 8 --memory-limit 17179869184 --pvc-limit 20` |
| `a8s admin gitops app create`, `AdminGitOpsApplication` | `{name: shop-api, repoUrl: https://github.com/acme/shop-gitops.git, authType: token, tokenFrom: {env: GITOPS_TOKEN}}` | `--name shop-api --repo-url https://github.com/acme/shop-gitops.git --auth-type token --token-stdin` |
| `a8s admin registry project create`, `AdminRegistryProject` | `{name: shop, publicProject: false}` | `--name shop --public=false` |
| `a8s admin sonarqube server-project create`, `AdminSonarQubeProject` | `{key: shop-api, name: "Shop API", mainBranch: main, visibility: private}` | `--key shop-api --name "Shop API" --main-branch main --visibility private` |
| `a8s admin sonarqube server-project update <project-key>`, `AdminSonarQubeProjectPatch` | `{key: shop-api-v2, visibility: private}` | `--key shop-api-v2 --visibility private` |
| `a8s admin docs update`, `AdminDocumentationUpdate` | `{path: guides/deploy.md, contentFile: deploy.md, sha: abc123, message: "Update deploy guide"}` | `--path guides/deploy.md --content-file deploy.md --sha abc123 --message "Update deploy guide"` |

```bash
a8s admin user create --file user.yaml
a8s admin registry project create --file registry-project.yaml
a8s admin sonarqube server-project update shop-api --file sonar-patch.yaml
```

The current admin quota approve and reject endpoints have no configurable
request payload. They remain payload-free actions:

```bash
a8s admin quota approve request-123
a8s admin quota reject request-123
```

## Commands Without Operation Documents

Do not add `--file` when there is no configurable backend request payload.
Examples include:

- reads: `get`, `list`, `status`, `history`, `metrics`, and `overview`
- streams: `watch`, `logs`, and events
- downloads: certificates, reports, backups, and avatars
- deletion commands whose only input is the resource identifier and
  confirmation
- payload-free actions such as restart, sync, retry, abort, or mark-read
- CLI-local commands such as auth status, context use, doctor, completion, and
  version

Arguments and operational controls such as resource IDs, `--yes`, `--wait`,
`--timeout`, `--output`, and `--verbose` remain flags or positional arguments.
They are not fields in an operation document.

## Domain File Inputs

Some commands consume files that are not YAML/JSON operation documents:

```bash
a8s project env import project-123 --file .env
a8s profile avatar upload --file avatar.png
a8s admin docs update --file documentation.md
a8s cluster console query cluster-123 --query-file query.sql
```

Use a distinct flag such as `--request-file` if a command must accept both an
operation document and a domain-content file, avoiding ambiguity.

## Operation Document Shape

Use versioned, typed documents:

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseUpgrade
spec:
  version: "17"
```

Resource identifiers normally remain positional arguments:

```bash
a8s database upgrade db-123 --file upgrade.yaml
```

Do not require the identifier to be duplicated inside the file. If an
identifier appears in both places, reject mismatches.

## Schema Discovery and Starter Files

Users should not have to search source code to discover supported fields.
Generate schemas from the CLI's typed operation-input models and provide:

```bash
# List supported operation kinds
a8s manifest kinds

# Print the complete schema and descriptions
a8s manifest schema DatabaseDeployment

# Generate a commented starter YAML file
a8s manifest init DatabaseDeployment --output-file database.yaml

# Validate without sending a backend request
a8s manifest validate --file database.yaml
```

`manifest init` should include every supported field as a commented example,
mark required fields, identify secret-reference fields, and show allowed enum
values. CI must detect drift between these schemas, the Go request models, and
the backend OpenAPI contract.

## Secrets

Do not place access tokens, refresh tokens, passwords, or secret values in
operation documents or ordinary flags because flags may be stored in shell
history.

Support secure sources such as:

- operating-system credential storage
- existing Kubernetes secret references
- environment-variable references
- dedicated `--password-stdin` or equivalent secret stdin flags
- interactive hidden prompts when running in a terminal

Redact secret values from dry-run, verbose logs, errors, and generated command
documentation.

## Required Shared Flags

Payload-bearing mutation commands should support:

```text
--file <path|->
--dry-run
--output table|json|yaml
--wait
--timeout <duration>
```

Use `--wait` only for asynchronous operations. Destructive mutations must also
support confirmation and `--yes`.

## Implementation Contract

Each payload-bearing command should define:

- one typed operation-input structure
- one strict YAML/JSON decoder
- equivalent Cobra flags
- a merge function that applies only explicitly changed flags
- one validator used by both input forms
- one mapper from operation input to the typed backend request
- tests proving that file-only and flags-only input produce identical payloads

Required tests include:

- flags only
- YAML file only
- JSON file only
- stdin input
- explicit flag override
- omitted flag does not override file value
- invalid and unknown fields
- conflicting inputs
- secret redaction
- dry-run performs no mutation
