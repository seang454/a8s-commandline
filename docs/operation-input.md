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

## Current Implementation Status

This document is the production input contract for the A8S CLI. The current
implementation already supports the foundation:

- generic backend mutation commands can send YAML or JSON with `--file`
- generic backend mutation commands can use `--set key=value` overrides
- `a8s database deploy` has a typed strict operation model and equivalent
  deployment flags
- `a8s context create` and `a8s context update` support `--file` plus explicit
  flag overrides
- selected high-value commands already have convenience flags and `--wait`
  workflows, including database deploy, cluster deploy, scan start, and
  workspace quota purchase
- payload-free backend actions must not accept request-body input such as
  `--file`, `--set`, `--form`, or `--upload`
- `a8s manifest kinds`, `a8s manifest schema`, `a8s manifest init`, and
  `a8s manifest validate` are available for operation-file discovery,
  starter generation, and local validation

The remaining production work is to replace the generic `--set` path with typed
models, full equivalent flags, strict field validation, and tests for every
operation kind listed below. Until a command has a typed Cobra implementation,
its YAML and flag examples define the required behavior that still needs to be
built.

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

Dense tables are only a quick command index. Detailed user-facing command
sections should use the clearer format demonstrated by `a8s cluster deploy`:

1. explain what the command does
2. show the minimal flags command
3. show minimal and production YAML
4. describe fields and whether they are required
5. group related flags by purpose
6. show file-plus-flag override behavior

Every command listed in a `Command and kind | YAML spec example | Equivalent
flags` table has a real copy-ready YAML document and complete command directly
below that table. Inline table values are summaries only.

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

### Real Context Examples

#### Create Context

```yaml
# context.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: Context
spec:
  server: https://api.a8s.example.com
  namespace: ns-team
  targetCluster: primary
```

```bash
a8s context create production --file context.yaml
a8s context create production --server https://api.a8s.example.com --namespace ns-team --target-cluster primary
```

#### Update Context

```yaml
# context-patch.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ContextPatch
spec:
  namespace: ns-prod
  targetCluster: prod-primary
```

```bash
a8s context update production --file context-patch.yaml
a8s context update production --namespace ns-prod --target-cluster prod-primary
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

### Real Workspace, Profile, and Integration Examples

#### Request Workspace Quota

```yaml
# quota-request.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: WorkspaceQuotaRequest
spec:
  requestedCpu: "4"
  requestedMemory: 8Gi
  requestedStorage: 100Gi
  reason: Production workload
  isPaid: false
  planName: Free
```

```bash
a8s workspace quota request --file quota-request.yaml
a8s workspace quota request --cpu 4 --memory 8Gi --storage 100Gi --reason "Production workload" --plan Free
```

#### Purchase Workspace Quota

```yaml
# quota-purchase.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: WorkspaceQuotaPurchase
spec:
  requestedCpu: "8"
  requestedMemory: 16Gi
  requestedStorage: 200Gi
  reason: Production upgrade
  isPaid: true
  planName: Premium
  paymentProvider: BAKONG
```

```bash
a8s workspace quota purchase --file quota-purchase.yaml --wait
a8s workspace quota purchase --cpu 8 --memory 16Gi --storage 200Gi --reason "Production upgrade" --plan Premium --payment-provider BAKONG --wait
```

#### Update Profile

```yaml
# profile.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProfileUpdate
spec:
  personal:
    firstName: Dara
    lastName: Sok
    displayName: Dara
    email: dara@example.com
    jobTitle: Engineer
    department: Platform
    bio: ""
  locale:
    city: Phnom Penh
    country: Cambodia
    timezone: Asia/Phnom_Penh
    language: en
    dateFormat: yyyy-MM-dd
```

```bash
a8s profile update --file profile.yaml
a8s profile update --first-name Dara --last-name Sok --display-name Dara --email dara@example.com --job-title Engineer --department Platform --city "Phnom Penh" --country Cambodia --timezone Asia/Phnom_Penh --language en --date-format yyyy-MM-dd
```

#### Connect Git Provider

```yaml
# git-connection.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: GitProviderConnection
spec:
  accessTokenFrom:
    env: GITHUB_TOKEN
  accessLevel: repository
  grantedScopes: repo,read:user
```

```bash
a8s git connect github --file git-connection.yaml
a8s git connect github --access-token-env GITHUB_TOKEN --access-level repository --granted-scopes repo,read:user
```

#### Sync DefectDojo Token

```yaml
# defectdojo-token.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DefectDojoToken
spec:
  apiTokenFrom:
    env: DEFECTDOJO_API_TOKEN
```

```bash
a8s defectdojo token sync project-123 --file defectdojo-token.yaml
a8s defectdojo token sync project-123 --api-token-env DEFECTDOJO_API_TOKEN
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

### Real Monolithic Project Examples

#### Deploy Project

Use the complete `project.yaml` example in the copy-ready examples section.

```bash
a8s project deploy --file project.yaml --wait
a8s project deploy --name shop-api --source-type git --repo-url https://github.com/acme/shop-api.git --repo-full-name acme/shop-api --branch main --app-port 8080 --architecture monolithic --auto-deploy --auto-deploy-trigger push --wait
```

The backend and frontend also support ZIP source deployment. The ZIP is a
domain-content file, while the YAML contains the deployment metadata:

```yaml
# project-zip.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectDeployment
spec:
  projectName: shop-api
  sourceType: zip
  architectureType: monolithic
  appPort: 8080
  envVars:
    - name: SPRING_PROFILES_ACTIVE
      value: production
      secret: false
```

```bash
a8s project deploy --file project-zip.yaml --source-archive shop-api.zip --wait
a8s project deploy --name shop-api --source-type zip --architecture monolithic --app-port 8080 --source-archive shop-api.zip --env SPRING_PROFILES_ACTIVE=production --wait
```

`--source-archive` is separate from `--file` because the archive is uploaded as
multipart content rather than decoded as an operation document.

#### Set Project Domain

```yaml
# project-domain.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectDomain
spec:
  customDomain: api.example.com
```

```bash
a8s project domain set project-123 --file project-domain.yaml
a8s project domain set project-123 --domain api.example.com
```

#### Connect Project Repository

```yaml
# repository-connection.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectRepositoryConnection
spec:
  repoProvider: github
  repoUrl: https://github.com/acme/shop-api.git
  repoFullName: acme/shop-api
  branch: main
  autoDeployEnabled: true
  autoDeployTrigger: push
```

```bash
a8s project repository connect project-123 --file repository-connection.yaml
a8s project repository connect project-123 --provider github --repo-url https://github.com/acme/shop-api.git --repo-full-name acme/shop-api --branch main --auto-deploy --auto-deploy-trigger push
```

#### Update Project Settings

```yaml
# project-settings.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectSettings
spec:
  alias: shop
  operatorNote: Critical service
  failureAlerts: true
  maintenanceMode: false
  protectFromDelete: true
```

```bash
a8s project settings update project-123 --file project-settings.yaml
a8s project settings update project-123 --alias shop --operator-note "Critical service" --failure-alerts --maintenance-mode=false --protect-from-delete
```

#### Set Project Environment

```yaml
# project-environment.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectEnvironment
spec:
  envVars:
    - name: SPRING_PROFILES_ACTIVE
      value: production
      secret: false
```

```bash
a8s project env set project-123 --file project-environment.yaml
a8s project env set project-123 --env SPRING_PROFILES_ACTIVE=production
```

#### Configure Auto Deploy

```yaml
# project-auto-deploy.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectAutoDeploy
spec:
  enabled: true
  branch: main
  autoDeployTrigger: push
  releaseTagPattern: "v*"
```

```bash
a8s project auto-deploy set project-123 --file project-auto-deploy.yaml
a8s project auto-deploy set project-123 --enabled --branch main --trigger push --release-tag-pattern "v*"
```

#### Create Project Webhook

```yaml
# project-webhook.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectWebhook
spec:
  name: shop-webhook
  branch: main
  autoDeployEnabled: true
  autoDeployTrigger: push
  createOnProvider: true
```

```bash
a8s project webhook create project-123 --file project-webhook.yaml
a8s project webhook create project-123 --name shop-webhook --branch main --auto-deploy --trigger push --create-on-provider
```

#### Roll Back Project

```yaml
# project-rollback.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectRollback
spec:
  releaseId: 11111111-1111-1111-1111-111111111111
```

```bash
a8s project rollback project-123 --file project-rollback.yaml --yes --wait
a8s project rollback project-123 --release-id 11111111-1111-1111-1111-111111111111 --yes --wait
```

#### Roll Back Project Release

```yaml
# release-rollback.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectReleaseRollback
spec:
  buildNumber: 42
  framework: spring
  statusMessage: Operator rollback
```

```bash
a8s project release rollback project-123 release-123 --file release-rollback.yaml --yes --wait
a8s project release rollback project-123 release-123 --build-number 42 --framework spring --status-message "Operator rollback" --yes --wait
```

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
| `a8s microservice detect --source-archive <path>`, `MicroserviceUploadDetection` | `{sourceName: shop-source}` | `--source-name shop-source --source-archive shop-source.zip` |

Complex service definitions should use `--service-file` rather than dozens of
flattened service flags. The top-level operation still supports both forms.

### Real Microservice Examples

#### Deploy Microservice Project

```yaml
# microservices.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceDeployment
spec:
  projectName: shop
  branch: main
  services:
    - name: api
      repoUrl: https://github.com/acme/shop.git
      repoFullName: acme/shop
      path: services/api
      appPort: 8080
      serviceType: backend
      exposePublic: true
```

```bash
a8s microservice deploy --file microservices.yaml --wait
a8s microservice deploy --project-name shop --branch main --service-file api-service.yaml --wait
```

For flag-based deployment, `api-service.yaml` contains one service object.

#### Apply Microservice Canvas

```yaml
# canvas.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceCanvas
spec:
  branch: main
  services:
    - name: api
      repoUrl: https://github.com/acme/shop.git
      repoFullName: acme/shop
      path: services/api
      appPort: 8080
```

```bash
a8s microservice apply project-123 --file canvas.yaml --wait
a8s microservice apply project-123 --branch main --service-file api-service.yaml --wait
```

#### Update Microservice Domains

```yaml
# microservice-domains.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceDomains
spec:
  services:
    - serviceId: 11111111-1111-1111-1111-111111111111
      customDomain: api.example.com
      platformSubdomain: api
```

```bash
a8s microservice domains update project-123 --file microservice-domains.yaml
a8s microservice domains update project-123 --service-domain 11111111-1111-1111-1111-111111111111=api.example.com --service-subdomain 11111111-1111-1111-1111-111111111111=api
```

#### Roll Back Microservice Project

```yaml
# microservice-rollback.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceRollback
spec:
  snapshotId: snapshot-123
```

```bash
a8s microservice rollback project-123 --file microservice-rollback.yaml --yes --wait
a8s microservice rollback project-123 --snapshot-id snapshot-123 --yes --wait
```

#### Set Microservice Environment

```yaml
# microservice-environment.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceEnvironment
spec:
  envVars:
    - name: SPRING_PROFILES_ACTIVE
      value: production
      secret: false
  runtimeConfigFile:
    fileName: application.yaml
    content: |
      server:
        port: 8080
```

```bash
a8s microservice env set project-123 service-123 --file microservice-environment.yaml
a8s microservice env set project-123 service-123 --env SPRING_PROFILES_ACTIVE=production --runtime-config-file application.yaml
```

#### Update Microservice Webhook

```yaml
# microservice-webhook.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceWebhook
spec:
  name: shop-webhook
  branch: main
  autoDeployEnabled: true
  autoDeployTrigger: push
  releaseTagPattern: "v*"
  releaseTriggerMode: tag
```

```bash
a8s microservice webhook update project-123 --file microservice-webhook.yaml
a8s microservice webhook update project-123 --name shop-webhook --branch main --auto-deploy --trigger push --release-tag-pattern "v*" --release-trigger-mode tag
```

#### Detect Microservices from Repository

```yaml
# microservice-detection.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceDetection
spec:
  repoUrl: https://github.com/acme/shop.git
  branch: main
  githubTokenFrom:
    env: GITHUB_TOKEN
```

```bash
a8s microservice detect --file microservice-detection.yaml
a8s microservice detect --repo https://github.com/acme/shop.git --branch main --github-token-env GITHUB_TOKEN
```

The backend also supports uploaded source detection:

```yaml
# microservice-upload-detection.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: MicroserviceUploadDetection
spec:
  sourceName: shop-source
```

```bash
a8s microservice detect --file microservice-upload-detection.yaml --source-archive shop-source.zip
a8s microservice detect --source-name shop-source --source-archive shop-source.zip
```

For multiple uploaded files, support repeatable `--source-file` and
`--source-path` flags while preserving their order:

```bash
a8s microservice detect \
  --source-name shop-source \
  --source-file services/api/pom.xml \
  --source-path services/api/pom.xml \
  --source-file services/web/package.json \
  --source-path services/web/package.json
```

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

### Real Single Database Examples

#### Deploy Database

```yaml
# database.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseDeployment
spec:
  releaseName: payments-db
  projectName: payments
  engine: postgresql
  deploymentMode: single
  databaseName: payments
  username: app
  version: "16"
  sizeProfile: small
  storageSize: 20Gi
  networkPolicyEnabled: true
  tls:
    enabled: true
    requireSsl: true
```

```bash
a8s database deploy --file database.yaml --wait
a8s database deploy --release-name payments-db --project-name payments --engine postgresql --deployment-mode single --database-name payments --username app --version 16 --size-profile small --storage-size 20Gi --network-policy --tls --require-ssl --wait
```

Provide database credentials securely. The backend accepts either a resolved
password or an existing authentication secret:

```yaml
# database-with-secret.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseDeployment
spec:
  projectName: payments
  engine: postgresql
  deploymentMode: single
  databaseName: payments
  version: "16"
  existingAuthSecretName: payments-database-credentials
  tls:
    enabled: true
    requireSsl: true
    existingSecretName: payments-database-tls
    includeCa: true
```

```bash
a8s database deploy --file database-with-secret.yaml --wait
a8s database deploy --project-name payments --engine postgresql --deployment-mode single --database-name payments --version 16 --existing-auth-secret payments-database-credentials --tls --require-ssl --tls-secret payments-database-tls --include-ca --wait
```

For a new password, use a secure reference or secret input flag:

```yaml
spec:
  passwordFrom:
    env: A8S_DATABASE_PASSWORD
```

```bash
a8s database deploy --file database.yaml --password-env A8S_DATABASE_PASSWORD --wait
a8s database deploy --file database.yaml --password-stdin --wait
```

#### Update Database

```yaml
# database-patch.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseDeploymentPatch
spec:
  sizeProfile: medium
  networkPolicyEnabled: true
  tls:
    enabled: true
    requireSsl: true
```

```bash
a8s database update db-123 --file database-patch.yaml --wait
a8s database update db-123 --size-profile medium --network-policy --tls --require-ssl --wait
```

#### Update Database Settings

```yaml
# database-settings.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseSettings
spec:
  alias: payments
  operatorNote: Primary database
  failureAlerts: true
  maintenanceMode: false
  protectFromDelete: true
```

```bash
a8s database settings update db-123 --file database-settings.yaml
a8s database settings update db-123 --alias payments --operator-note "Primary database" --failure-alerts --maintenance-mode=false --protect-from-delete
```

#### Upgrade Database

```yaml
# database-upgrade.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseUpgrade
spec:
  version: "17"
```

```bash
a8s database upgrade db-123 --file database-upgrade.yaml --wait
a8s database upgrade db-123 --version 17 --wait
```

#### Clone Database from Backup

```yaml
# database-clone.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseClone
spec:
  sourceDeploymentId: 11111111-1111-1111-1111-111111111111
  backupRunId: 22222222-2222-2222-2222-222222222222
  projectName: payments-clone
  databaseName: payments
  version: "16"
  storageSize: 20Gi
```

```bash
a8s database clone-from-backup --file database-clone.yaml --wait
a8s database clone-from-backup --source-deployment-id 11111111-1111-1111-1111-111111111111 --backup-run-id 22222222-2222-2222-2222-222222222222 --project-name payments-clone --database-name payments --version 16 --storage-size 20Gi --wait
```

#### Rotate and Verify Database Password

```yaml
# database-password.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabasePasswordRotation
spec:
  passwordFrom:
    env: A8S_NEW_DATABASE_PASSWORD
```

```bash
a8s database rotate-password db-123 --file database-password.yaml --wait
a8s database rotate-password db-123 --password-env A8S_NEW_DATABASE_PASSWORD --wait
```

```yaml
# database-password-verification.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabasePasswordVerification
spec:
  passwordFrom:
    env: A8S_DATABASE_PASSWORD
```

```bash
a8s database verify-password db-123 --file database-password-verification.yaml
a8s database verify-password db-123 --password-env A8S_DATABASE_PASSWORD
```

#### Run Database Query

```yaml
# database-query.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseQuery
spec:
  query: SELECT now()
```

```bash
a8s database console query db-123 --file database-query.yaml
a8s database console query db-123 --query "SELECT now()"
```

#### Configure Database Backup

```yaml
# database-backup.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseBackupSettings
spec:
  enabled: true
  destinationPath: s3://backups/payments
  credentialSecret: backup-credentials
  retentionPolicy: 30d
  schedule: "0 0 * * *"
  scheduleStartAt: "2026-06-06T00:00:00Z"
```

```bash
a8s database backup settings set db-123 --file database-backup.yaml
a8s database backup settings set db-123 --enabled --destination s3://backups/payments --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *" --schedule-start-at 2026-06-06T00:00:00Z
```

## Database Cluster Mutations

### Deploy a Database Cluster

Creates and deploys a managed database cluster.

```text
a8s cluster deploy
```

The input has three main sections:

```text
spec
|-- releaseName and projectName     Identifies the A8S deployment
|-- cluster                         Selects and describes the Kubernetes cluster
`-- database                        Configures the database running in the cluster
    |-- resource                    Configures CPU and memory
    `-- backup                      Configures automated backups
```

Users do not need to provide every available field. Start with the minimal
input and add production options only when required.

Use YAML for a repeatable production deployment:

```bash
a8s cluster deploy --file cluster.yaml --wait
```

Use flags for a quick deployment:

```bash
a8s cluster deploy \
  --release-name orders-cluster \
  --project-name orders \
  --name orders \
  --engine POSTGRESQL \
  --target-cluster production-primary \
  --wait
```

Use YAML with flags when only a few values need to change:

```bash
a8s cluster deploy \
  --file cluster.yaml \
  --target-cluster staging-primary \
  --wait
```

In this example, the CLI loads everything from `cluster.yaml`, then replaces
only `targetClusterName` with `staging-primary`.

#### Minimal Cluster YAML

This contains only the main values needed to identify and deploy the cluster:

```yaml
# cluster.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterDeployment

spec:
  releaseName: orders-cluster
  projectName: orders

  cluster:
    name: orders
    platformConfig:
      targetClusterName: production-primary

  database:
    engine: POSTGRESQL
```

#### Production Cluster YAML

Add capacity, networking, monitoring, and backup configuration when needed:

```yaml
# cluster-production.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterDeployment

spec:
  releaseName: orders-cluster
  projectName: orders

  cluster:
    name: orders
    environment: PRODUCTION
    notes: Primary orders database cluster

    platformConfig:
      targetClusterName: production-primary

  database:
    engine: POSTGRESQL
    enabled: true
    instances: 3
    version: "16"
    databaseName: appdb
    username: app_postgres

    storageSize: 100Gi
    storageClass: longhorn

    externalAccessEnabled: false
    port: 5432
    publicHostnames:
      - postgres.orders.example.com
    tlsEnabled: true
    tlsMode: CERT_MANAGER
    monitoringEnabled: true

    resource:
      resourceProfile: MEDIUM
      cpuRequest: "1"
      memRequest: 2Gi
      cpuLimit: "2"
      memLimit: 4Gi

    backup:
      enabled: true
      destinationPath: s3://backups/orders
      credentialSecret: backup-credentials
      retentionPolicy: 30d
      schedule: "0 0 * * *"

    postgresql:
      walEnabled: true
      walSize: 10Gi
      bootstrapDatabase: appdb
      bootstrapOwner: app_postgres
```

#### Cluster Deployment Fields

| Field | Required | Description |
|---|---|---|
| `spec.releaseName` | Recommended | Kubernetes or Helm release name used to track deployment status. |
| `spec.projectName` | Yes | A8S project name associated with the cluster. |
| `spec.cluster.name` | Yes | Human-readable cluster name. |
| `spec.cluster.clusterKey` | No | Stable internal cluster key. Normally generated by the backend. |
| `spec.cluster.workspaceId` | No | Backend workspace identifier. Ordinary users should not provide it. |
| `spec.cluster.environment` | No | Environment classification, such as `DEVELOPMENT`, `STAGING`, or `PRODUCTION`. |
| `spec.cluster.domain` | No | Base domain associated with the cluster. |
| `spec.cluster.externalIp` | No | Explicit external IP when required by the deployment environment. |
| `spec.cluster.notes` | No | Operator notes describing the cluster. |
| `spec.cluster.platformConfig.targetClusterName` | No | Kubernetes target-cluster alias. The active context supplies it when omitted. |
| `spec.database.engine` | Yes | Database engine, such as `POSTGRESQL`, `MONGODB`, `MYSQL`, `REDIS`, or `CASSANDRA`. |
| `spec.database.enabled` | No | Enables or disables the database deployment. Disabling may uninstall resources. |
| `spec.database.instances` | No | Number of database instances. Support depends on the selected engine. |
| `spec.database.version` | No | Requested database version. |
| `spec.database.databaseName` | No | Initial database, keyspace, or Redis database identifier. |
| `spec.database.username` | No | Initial application username. |
| `spec.database.storageSize` | No | Persistent storage size, such as `100Gi`. |
| `spec.database.storageClass` | No | Kubernetes storage class, such as `longhorn`. |
| `spec.database.externalAccessEnabled` | No | Exposes the database outside its Kubernetes cluster when supported. |
| `spec.database.port` | No | Requested external port. The backend allocates one when omitted. |
| `spec.database.publicHostnames` | No | Public DNS hostnames. MongoDB and Redis may require multiple hostnames. |
| `spec.database.tlsEnabled` | No | Enables TLS for database connections. |
| `spec.database.tlsMode` | No | TLS strategy: `DISABLED`, `OPERATOR`, `CERT_MANAGER`, or `EXISTING_SECRET`. |
| `spec.database.tlsSecretName` | No | Existing TLS certificate secret when required by the TLS mode. |
| `spec.database.tlsCaSecretName` | No | Existing certificate-authority secret. |
| `spec.database.monitoringEnabled` | No | Enables database monitoring integration. |
| `spec.database.notes` | No | Notes specifically associated with the database instance. |
| `spec.database.resource` | No | CPU and memory request, limit, and resource-profile settings. |
| `spec.database.backup` | No | Automated backup destination, credentials reference, retention, and schedule. |
| `spec.database.postgresql` | No | PostgreSQL-specific WAL and bootstrap settings. |
| `spec.database.mongo` | No | MongoDB replica-set horizon and DNS settings. |
| `spec.database.mysql` | No | MySQL HAProxy settings. |
| `spec.database.redis` | No | Redis exporter, follower, and ACL settings. |
| `spec.database.cassandra` | No | Cassandra cluster, datacenter, and client-auth settings. |

The field-to-flag mapping follows predictable names:

```text
spec.releaseName                            -> --release-name
spec.projectName                            -> --project-name
spec.cluster.name                           -> --name
spec.cluster.environment                    -> --environment
spec.cluster.platformConfig.targetClusterName -> --target-cluster
spec.database.engine                        -> --engine
spec.database.instances                     -> --instances
spec.database.storageSize                   -> --storage-size
spec.database.monitoringEnabled             -> --monitoring
spec.database.backup.schedule               -> --backup-schedule
```

#### Frontend Defaults and CLI Parity

The frontend currently accepts a simplified cluster form and expands it into
the structured backend request. The CLI must perform equivalent normalization
when users provide flags or a simplified operation file.

| Input | Current frontend behavior | Recommended CLI behavior |
|---|---|---|
| `releaseName` | Defaults to `projectName`. | Use the same default when omitted. |
| `cluster.name` | Uses `projectName`. | Use `projectName` unless `--name` is supplied. |
| `environment` | Defaults to `DEVELOPMENT`. | Use `DEVELOPMENT` unless explicitly configured. |
| `engine` | Defaults to `POSTGRESQL`. | Prefer requiring it; otherwise clearly document the default. |
| `instances` | Defaults to `3`. | Apply an engine-aware default. |
| `sizeProfile` | Required by the frontend and expanded to resources. | Support `--size-profile SMALL|MEDIUM|LARGE` and expand it before the API call. |
| `databaseName` | Uses an engine-specific default. | Use the same engine-specific defaults when omitted. |
| `username` | Uses an engine-specific default. | Use the same engine-specific defaults when omitted. |
| `targetClusterName` | Uses the configured frontend default. | Resolve from `--target-cluster`, then the active context. |
| `password` | Collected from the UI and mapped to the engine-specific secret field. | Read securely using `--password-stdin`, a hidden prompt, or a secure reference. |
| Redis ACL permissions | Applies a default ACL permission set. | Apply and document the same allowed default permissions. |
| Engine-specific settings | Generates PostgreSQL, MongoDB, MySQL, Redis, or Cassandra defaults. | Generate the same defaults or require an explicit engine-specific block. |

The CLI should send the structured backend request directly. It must not call
the frontend Next.js proxy merely to obtain these defaults.

#### Engine-Specific Configuration

Only configure the block matching `spec.database.engine`.

PostgreSQL:

```yaml
database:
  engine: POSTGRESQL
  postgresql:
    walEnabled: true
    walSize: 10Gi
    bootstrapDatabase: appdb
    bootstrapOwner: app_postgres
```

MongoDB:

```yaml
database:
  engine: MONGODB
  mongo:
    replicaSetHorizonsEnabled: true
    replicaSetHorizonsBasePort: 30000
    clusterServiceDNSSuffix: cluster.local
```

MySQL:

```yaml
database:
  engine: MYSQL
  mysql:
    haproxySize: 2
```

Redis:

```yaml
database:
  engine: REDIS
  redis:
    exporterEnabled: true
    followersEnabled: true
    aclPermissions:
      - +@read
      - +@write
      - +@keyspace
      - +@connection
```

Cassandra:

```yaml
database:
  engine: CASSANDRA
  cassandra:
    clusterName: orders
    datacenter: dc1
    requireClientAuth: true
```

Reject operation files containing an engine-specific block that does not match
the selected engine.

#### Platform Networking Configuration

Advanced deployments may configure ingress, a shared TCP gateway, an external
TCP proxy, and Cloudflare DNS:

```yaml
cluster:
  platformConfig:
    targetClusterName: production-primary
    ingressEnabled: true
    ingressClassName: nginx
    sharedGatewayEnabled: true
    externalTcpProxyEnabled: false
    externalTcpProxyNodeSelectorKey: node-role.kubernetes.io/gateway
    externalTcpProxyNodeSelectorValue: "true"
    externalTcpProxyClusterDomain: cluster.local
    cloudflareEnabled: true
    cloudflareZoneName: example.com
    cloudflareZoneId: zone-id
```

Equivalent flags:

```text
--ingress-enabled
--ingress-class <name>
--shared-gateway
--external-tcp-proxy
--external-tcp-proxy-node-selector-key <key>
--external-tcp-proxy-node-selector-value <value>
--external-tcp-proxy-cluster-domain <domain>
--cloudflare-enabled
--cloudflare-zone-name <name>
--cloudflare-zone-id <id>
```

Cloudflare API tokens must use the secure secret input described below.

#### Cluster Deployment Flags

Identity and target:

```text
--release-name <name>
--project-name <name>
--name <cluster-name>
--environment DEVELOPMENT|STAGING|PRODUCTION
--target-cluster <alias>
--notes <text>
```

Database capacity:

```text
--engine POSTGRESQL|MONGODB|MYSQL|REDIS|CASSANDRA
--instances <number>
--version <version>
--storage-size <quantity>
--storage-class <name>
--database-name <name>
--username <name>
--port <number>
--public-hostname <hostname>       Repeatable
```

Features:

```text
--external-access
--tls
--tls-mode DISABLED|OPERATOR|CERT_MANAGER|EXISTING_SECRET
--tls-secret <secret-name>
--tls-ca-secret <secret-name>
--monitoring
--size-profile SMALL|MEDIUM|LARGE|CUSTOM
--resource-profile <profile>
--cpu-request <quantity>
--memory-request <quantity>
--cpu-limit <quantity>
--memory-limit <quantity>
```

Backup:

```text
--backup-enabled
--backup-destination <path>
--backup-credential-secret <secret-name>
--backup-retention <duration>
--backup-schedule <cron>
```

Engine-specific flags should use an engine prefix so their purpose remains
clear:

```text
--postgres-wal-enabled
--postgres-wal-size <quantity>
--postgres-bootstrap-database <name>
--postgres-bootstrap-owner <name>

--mongo-replica-set-horizons
--mongo-replica-set-horizons-base-port <number>
--mongo-cluster-service-dns-suffix <suffix>

--mysql-haproxy-size <number>

--redis-exporter
--redis-followers
--redis-acl-permission <permission>     Repeatable

--cassandra-cluster-name <name>
--cassandra-datacenter <name>
--cassandra-require-client-auth
```

The active context supplies `namespace` and may supply `targetClusterName`.
Users can override them with `--namespace` and `--target-cluster`.

#### Cluster Secrets

The backend request accepts engine-specific password fields and a Cloudflare
API token. Do not store their literal values in YAML:

```yaml
spec:
  secretRefs:
    databasePassword:
      env: A8S_CLUSTER_DATABASE_PASSWORD
    cloudflareApiToken:
      env: CLOUDFLARE_API_TOKEN
```

Equivalent secure flags:

```text
--password-stdin
--password-env <environment-variable>
--cloudflare-api-token-stdin
--cloudflare-api-token-env <environment-variable>
```

The CLI resolves these references locally and maps the database password to
exactly one backend field: `pgPassword`, `mongoPassword`, `mysqlPassword`,
`redisPassword`, or `cassandraPassword`.

#### Backend Validation Gap

The backend currently requires a request containing `database.engine`, and
requires a cluster name only when both release name and namespace are missing.
Many other fields are normalized by the frontend rather than validated by the
backend.

For a consistent production experience, the CLI should validate its documented
contract before sending requests. The backend should eventually share or
enforce equivalent validation so frontend, CLI, and direct API callers behave
consistently.

#### Engine Compatibility Note

Some frontend TypeScript types also contain `oracle` and `sqlserver`, but the
cluster deployment proxy and backend currently support only:

```text
POSTGRESQL, MONGODB, MYSQL, REDIS, CASSANDRA
```

The CLI must reject unsupported engines before sending a request. Add Oracle or
SQL Server only after the backend cluster API supports them.

### Update Cluster Deployment

Changes configurable deployment values for an existing cluster.

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterDeploymentPatch
spec:
  database:
    instances: 5
    storageSize: 200Gi
    monitoringEnabled: true
```

```bash
a8s cluster update cluster-123 --file cluster-patch.yaml --wait

a8s cluster update cluster-123 \
  --instances 5 \
  --storage-size 200Gi \
  --monitoring \
  --wait
```

### Update Cluster Operator Settings

Changes A8S operator metadata and protection settings. It does not change
database capacity.

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterSettings
spec:
  alias: orders
  operatorNote: Primary production cluster
  failureAlerts: true
  maintenanceMode: false
  protectFromDelete: true
```

```bash
a8s cluster settings update cluster-123 --file cluster-settings.yaml

a8s cluster settings update cluster-123 \
  --alias orders \
  --failure-alerts \
  --protect-from-delete
```

### Upgrade Cluster Version

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterUpgrade
spec:
  version: "17"
```

```bash
a8s cluster upgrade cluster-123 --file cluster-upgrade.yaml --wait
a8s cluster upgrade cluster-123 --version 17 --wait
```

### Rotate Cluster Password

The frontend rotates a cluster password by sending an engine-specific secret
through the general cluster update endpoint. The CLI should expose a dedicated,
safe workflow command even though the backend has no dedicated rotate-password
route.

Do not place the password value directly in YAML. Use a secure reference:

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterPasswordRotation
spec:
  passwordFrom:
    env: A8S_NEW_CLUSTER_PASSWORD
```

```bash
a8s cluster rotate-password cluster-123 --file cluster-password.yaml --wait
a8s cluster rotate-password cluster-123 --password-stdin --wait
```

The CLI must inspect the cluster engine and map the secret to exactly one
backend update field:

```text
POSTGRESQL -> secrets.pgPassword
MONGODB    -> secrets.mongoPassword
MYSQL      -> secrets.mysqlPassword
REDIS      -> secrets.redisPassword
CASSANDRA  -> secrets.cassandraPassword
```

Never print the supplied or generated password. If the CLI generates a
password, display it once only when explicitly requested with a secure
one-time-output option.

### Clone Cluster from Backup

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterClone
spec:
  sourceClusterId: 11111111-1111-1111-1111-111111111111
  backupRunId: 22222222-2222-2222-2222-222222222222
  projectName: orders-clone
  releaseName: orders-clone
  version: "16"
  instances: 3
  storageSize: 100Gi
  monitoringEnabled: true
```

```bash
a8s cluster clone-from-backup --file cluster-clone.yaml --wait

a8s cluster clone-from-backup \
  --source-cluster-id 11111111-1111-1111-1111-111111111111 \
  --backup-run-id 22222222-2222-2222-2222-222222222222 \
  --project-name orders-clone \
  --wait
```

### Configure Cluster Backup

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterBackupSettings
spec:
  enabled: true
  destinationPath: s3://backups/orders
  credentialSecret: backup-credentials
  retentionPolicy: 30d
  schedule: "0 0 * * *"
```

```bash
a8s cluster backup settings set cluster-123 --file cluster-backup.yaml

a8s cluster backup settings set cluster-123 \
  --enabled \
  --destination s3://backups/orders \
  --credential-secret backup-credentials \
  --retention 30d \
  --schedule "0 0 * * *"
```

Use `--release <release-name>` instead of a cluster ID when configuring a
deployment that has not yet produced a persistent cluster record.

### Run Cluster Console Query

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ClusterQuery
spec:
  query: SELECT now()
```

```bash
a8s cluster console query cluster-123 --file query.yaml
a8s cluster console query cluster-123 --query "SELECT now()"
```

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

### Real Backup, Quality, and Alert Examples

#### Configure Unified Backup

```yaml
# backup-settings.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: BackupSettings
spec:
  enabled: true
  destinationPath: s3://backups/resource
  credentialSecret: backup-credentials
  retentionPolicy: 30d
  schedule: "0 0 * * *"
```

```bash
a8s backup settings set database db-123 --file backup-settings.yaml
a8s backup settings set database db-123 --enabled --destination s3://backups/resource --credential-secret backup-credentials --retention 30d --schedule "0 0 * * *"
```

#### Start Image Scan

```yaml
# image-scan.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ImageScan
spec:
  sourceKind: image
  imageRef: nginx:1.27
  forceRescan: false
```

```bash
a8s scan start --file image-scan.yaml --wait
a8s scan start --source-kind image --image nginx:1.27 --force-rescan=false --wait
```

Private-registry scan:

```yaml
# private-image-scan.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ImageScan
spec:
  sourceKind: registry
  registryUrl: registry.example.com
  imageName: platform/shop-api
  imageTag: "1.0.0"
  privateRegistry: true
  username: scanner
  passwordFrom:
    env: REGISTRY_PASSWORD
  forceRescan: true
```

```bash
a8s scan start --file private-image-scan.yaml --wait
a8s scan start --source-kind registry --registry-url registry.example.com --image-name platform/shop-api --image-tag 1.0.0 --private-registry --username scanner --password-env REGISTRY_PASSWORD --force-rescan --wait
```

Build and scan from a Git repository:

```yaml
# git-image-scan.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ImageScan
spec:
  sourceKind: repository
  repositoryUrl: https://github.com/acme/shop-api.git
  branchOrTag: main
  dockerfilePath: Dockerfile
  buildContext: .
  targetImageName: shop-api:scan
  privateRepository: false
```

```bash
a8s scan start --file git-image-scan.yaml --wait
a8s scan start --source-kind repository --repository-url https://github.com/acme/shop-api.git --branch-or-tag main --dockerfile Dockerfile --build-context . --target-image shop-api:scan --wait
```

#### Run Benchmark

```yaml
# benchmark.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: BenchmarkRun
spec:
  concurrency: 20
  totalRequests: 1000
  targetPath: /api/health
  method: GET
  headers:
    Accept: application/json
```

```bash
a8s benchmark run project-123 --file benchmark.yaml --wait
a8s benchmark run project-123 --concurrency 20 --total-requests 1000 --target-path /api/health --method GET --header Accept=application/json --wait
```

#### Create or Update Alert Channel

```yaml
# alert-channel.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AlertChannel
spec:
  name: operations
  type: telegram
  credentialFrom:
    env: TELEGRAM_BOT_TOKEN
  secondaryCredentialFrom:
    env: TELEGRAM_CHAT_ID
  targetProject: shop-api
```

```bash
a8s alert channel create --file alert-channel.yaml
a8s alert channel create --name operations --type telegram --credential-env TELEGRAM_BOT_TOKEN --secondary-credential-env TELEGRAM_CHAT_ID --target-project shop-api

a8s alert channel update channel-123 --file alert-channel.yaml
a8s alert channel update channel-123 --name operations --type telegram --credential-env TELEGRAM_BOT_TOKEN --secondary-credential-env TELEGRAM_CHAT_ID --target-project shop-api
```

#### Configure Project and User Alerts

```yaml
# project-alerts.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: ProjectAlertConfig
spec:
  telegramEnabled: true
  emailEnabled: true
  backupAlertsEnabled: true
  securityAlertsEnabled: true
  telegramChannelName: operations
  emailAddress: ops@example.com
```

```bash
a8s alert project-config set project-123 --file project-alerts.yaml
a8s alert project-config set project-123 --telegram-enabled --email-enabled --backup-alerts --security-alerts --telegram-channel operations --email-address ops@example.com
```

```yaml
# user-alerts.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: UserAlertConfig
spec:
  quotaAlertsEnabled: true
  globalSecurityAlertsEnabled: true
```

```bash
a8s alert user-config set --file user-alerts.yaml
a8s alert user-config set --quota-alerts --global-security-alerts
```

#### Configure Notification Preferences

```yaml
# notification-preferences.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: NotificationPreferences
spec:
  buildFailures: true
  rolloutReady: true
  vulnerabilityFindings: true
  weeklyDigest: false
```

```bash
a8s notification preferences set --file notification-preferences.yaml
a8s notification preferences set --build-failures --rollout-ready --vulnerability-findings --weekly-digest=false
```

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

### Real Administrative Examples

#### Create or Update Admin User

```yaml
# admin-user-create.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminUserCreate
spec:
  username: dara
  email: dara@example.com
  firstName: Dara
  lastName: Sok
  passwordFrom:
    env: A8S_INITIAL_PASSWORD
```

```bash
a8s admin user create --file admin-user-create.yaml
a8s admin user create --username dara --email dara@example.com --first-name Dara --last-name Sok --password-env A8S_INITIAL_PASSWORD
```

```yaml
# admin-user-update.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminUserUpdate
spec:
  username: dara
  email: dara@example.com
  firstName: Dara
  lastName: Sok
```

```bash
a8s admin user update user-123 --file admin-user-update.yaml
a8s admin user update user-123 --username dara --email dara@example.com --first-name Dara --last-name Sok
```

#### Update Admin Project or Cluster

```yaml
# admin-project-update.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminProjectUpdate
spec:
  name: shop-api
  userId: 11111111-1111-1111-1111-111111111111
  repoProvider: github
  repoFullName: acme/shop-api
  repoUrl: https://github.com/acme/shop-api.git
  branch: main
  framework: spring
  appPort: 8080
  status: ACTIVE
  autoDeployEnabled: true
```

```bash
a8s admin project update project-123 --file admin-project-update.yaml
a8s admin project update project-123 --name shop-api --user-id 11111111-1111-1111-1111-111111111111 --repo-provider github --repo-full-name acme/shop-api --repo-url https://github.com/acme/shop-api.git --branch main --framework spring --app-port 8080 --status ACTIVE --auto-deploy
```

```yaml
# admin-cluster-update.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminClusterUpdate
spec:
  alias: primary
  operatorNote: Production cluster
  failureAlerts: true
  maintenanceMode: false
  protectFromDelete: true
```

```bash
a8s admin cluster update cluster-123 --file admin-cluster-update.yaml
a8s admin cluster update cluster-123 --alias primary --operator-note "Production cluster" --failure-alerts --maintenance-mode=false --protect-from-delete
```

#### Set Admin Cluster Quota

```yaml
# admin-cluster-quota.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminClusterQuota
spec:
  cpuLimit: 8
  memoryLimit: 17179869184
  pvcLimit: 20
```

```bash
a8s admin cluster quota set production-primary ns-team --file admin-cluster-quota.yaml
a8s admin cluster quota set production-primary ns-team --cpu-limit 8 --memory-limit 17179869184 --pvc-limit 20
```

#### Create Admin GitOps or Registry Project

```yaml
# gitops-app.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminGitOpsApplication
spec:
  name: shop-api
  repoUrl: https://github.com/acme/shop-gitops.git
  authType: token
  tokenFrom:
    env: GITOPS_TOKEN
```

```bash
a8s admin gitops app create --file gitops-app.yaml
a8s admin gitops app create --name shop-api --repo-url https://github.com/acme/shop-gitops.git --auth-type token --token-env GITOPS_TOKEN
```

```yaml
# registry-project.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminRegistryProject
spec:
  name: shop
  publicProject: false
```

```bash
a8s admin registry project create --file registry-project.yaml
a8s admin registry project create --name shop --public=false
```

#### Create or Update SonarQube Project

```yaml
# sonarqube-project.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminSonarQubeProject
spec:
  key: shop-api
  name: Shop API
  mainBranch: main
  visibility: private
```

```bash
a8s admin sonarqube server-project create --file sonarqube-project.yaml
a8s admin sonarqube server-project create --key shop-api --name "Shop API" --main-branch main --visibility private
```

```yaml
# sonarqube-project-patch.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminSonarQubeProjectPatch
spec:
  key: shop-api-v2
  visibility: private
```

```bash
a8s admin sonarqube server-project update shop-api --file sonarqube-project-patch.yaml
a8s admin sonarqube server-project update shop-api --key shop-api-v2 --visibility private
```

#### Update Admin Documentation

```yaml
# documentation-update.yaml
apiVersion: cli.a8s.io/v1alpha1
kind: AdminDocumentationUpdate
spec:
  path: guides/deploy.md
  contentFile: deploy.md
  sha: abc123
  message: Update deploy guide
```

```bash
a8s admin docs update --file documentation-update.yaml
a8s admin docs update --path guides/deploy.md --content-file deploy.md --sha abc123 --message "Update deploy guide"
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

The implementation should also reject body-input flags on these commands. A
hidden flag is not enough; commands such as `a8s admin quota approve` and
`a8s project redeploy` should fail fast if the user supplies `--file`, `--set`,
`--form`, or `--upload`.

## Domain File Inputs

Some commands consume files that are not YAML/JSON operation documents:

```bash
a8s project env import project-123 --env-file .env
a8s microservice env import project-123 service-123 --env-file .env
a8s profile avatar upload --avatar-file avatar.png
a8s admin docs update --content-file documentation.md --path guides/deploy.md --message "Update deploy guide"
a8s microservice detect --source-archive source.zip
a8s cluster console query cluster-123 --query-file query.sql
```

Use a distinct flag such as `--request-file` if a command must accept both an
operation document and a domain-content file, avoiding ambiguity.

Recommended file-flag meanings:

| Flag | Meaning |
|---|---|
| `--file <yaml-or-json>` | Typed operation document. |
| `--env-file <path>` | Dotenv content imported into a project or service. |
| `--source-archive <path>` | ZIP archive containing application source. |
| `--source-file <path>` | Individual uploaded source file; repeatable. |
| `--service-file <path>` | One microservice service-definition document. |
| `--content-file <path>` | Documentation or runtime configuration content. |
| `--avatar-file <path>` | Profile image upload. |
| `--query-file <path>` | SQL or database-console query content. |
| `--output-file <path>` | Destination for downloaded content. |

Do not overload `--file` for raw content because users must always be able to
tell whether the CLI will parse a file as an operation document or upload its
bytes unchanged.

## Final Input Coverage Audit

The backend and frontend mutation flows were compared against this document.
All currently configurable user-facing request payloads now have YAML and
equivalent-flag examples.

The following commands intentionally do not use operation YAML because their
backend operations have no configurable request payload:

```text
auth verify-email start
auth onboarding start
workspace bootstrap
admin project restore
admin user reactivate
admin gitops app abort
admin gitops app retry
admin gitops app sync
admin quota approve
admin quota reject
admin docs publish
git sync-token
database backup run
database backup restore
database backup restore cancel
backup trigger
backup restore
backup restore cancel
cluster console test
database console test
microservice redeploy
project domain sync
project redeploy
project webhook rotate
notification read
profile account deactivate
profile account reactivate
sonarqube access
```

These commands continue to use positional identifiers and operational flags
such as `--wait`, `--timeout`, and `--yes`.

Internal callbacks, Jenkins callbacks, and provider webhook receivers remain
excluded from the user CLI entirely.

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

Users should not have to search source code to discover supported fields. For
production, generate schemas from the CLI's typed operation-input models and
provide:

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

The current `a8s manifest` implementation provides schema summaries, starter
files, and local validation for known operation kinds. Strict field validation
is available for typed operation models such as `DatabaseDeployment`, `Context`,
and `ContextPatch`; other known kinds receive envelope, kind, and required-field
validation until their typed request models are implemented.

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
