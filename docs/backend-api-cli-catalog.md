# A8S Backend API to CLI Catalog

Generated from controller annotations in `D:\CSTADPreUniversityTraining\ITP\finalProject\a8s-backend`.

- Feature folders: 21
- Controllers: 38
- HTTP route patterns: 248
- WebSocket routes: 4

Global CLI flags should include `--server`, `--context`, `--namespace`, `--target-cluster`, `--output`, `--timeout`, and `--verbose`.

## admin

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/admin/clusters` | `a8s admin cluster list` | `AdminClusterController` |
| `PATCH` | `/api/v1/admin/clusters/{clusterId}` | `a8s admin cluster update <cluster-id>` | `AdminClusterController` |
| `GET` | `/api/v1/admin/clusters/kubernetes` | `a8s admin cluster nodes` | `AdminClusterController` |
| `GET` | `/api/v1/admin/clusters/kubernetes/{alias}/health` | `a8s admin cluster health <alias>` | `AdminClusterController` |
| `GET` | `/api/v1/admin/clusters/kubernetes/{alias}/quotas` | `a8s admin cluster quota list <alias>` | `AdminClusterController` |
| `DELETE` | `/api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}` | `a8s admin cluster quota delete <alias> <namespace>` | `AdminClusterController` |
| `PUT` | `/api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace}` | `a8s admin cluster quota set <alias> <namespace>` | `AdminClusterController` |
| `GET` | `/api/v1/admin/projects` | `a8s admin project list` | `AdminController` |
| `DELETE` | `/api/v1/admin/projects/{projectId}` | `a8s admin project deactivate <project-id>` | `AdminController` |
| `PATCH` | `/api/v1/admin/projects/{projectId}` | `a8s admin project update <project-id>` | `AdminController` |
| `POST` | `/api/v1/admin/projects/{projectId}/restore` | `a8s admin project restore <project-id>` | `AdminController` |
| `GET` | `/api/v1/admin/users` | `a8s admin user list` | `AdminController` |
| `POST` | `/api/v1/admin/users` | `a8s admin user create` | `AdminController` |
| `DELETE` | `/api/v1/admin/users/{userId}` | `a8s admin user deactivate <user-id>` | `AdminController` |
| `PATCH` | `/api/v1/admin/users/{userId}` | `a8s admin user update <user-id>` | `AdminController` |
| `POST` | `/api/v1/admin/users/{userId}/reactivate` | `a8s admin user reactivate <user-id>` | `AdminController` |
| `POST` | `/api/v1/admin/gitops/apps` | `a8s admin gitops app create` | `AdminGitOpsController` |
| `POST` | `/api/v1/admin/gitops/apps/{appId}/abort` | `a8s admin gitops app abort <app-id>` | `AdminGitOpsController` |
| `POST` | `/api/v1/admin/gitops/apps/{appId}/retry` | `a8s admin gitops app retry <app-id>` | `AdminGitOpsController` |
| `POST` | `/api/v1/admin/gitops/apps/{appId}/sync` | `a8s admin gitops app sync <app-id>` | `AdminGitOpsController` |
| `GET` | `/api/v1/admin/gitops/overview` | `a8s admin gitops overview` | `AdminGitOpsController` |
| `GET` | `/api/v1/admin/logs/clusters` | `a8s admin logs clusters` | `AdminLogsController` |
| `GET` | `/api/v1/admin/logs/namespaces` | `a8s admin logs namespaces` | `AdminLogsController` |
| `GET` | `/api/v1/admin/logs/pods` | `a8s admin logs pods` | `AdminLogsController` |
| `GET` | `/api/v1/admin/logs/query` | `a8s admin logs query` | `AdminLogsController` |
| `GET` | `/api/v1/admin/logs/workloads` | `a8s admin logs workloads` | `AdminLogsController` |
| `GET` | `/api/v1/admin/monitoring/overview` | `a8s admin monitoring overview` | `AdminMonitoringController` |
| `GET` | `/api/v1/admin/quota-requests` | `a8s admin quota list` | `AdminQuotaRequestController` |
| `POST` | `/api/v1/admin/quota-requests/{id}/approve` | `a8s admin quota approve <request-id>` | `AdminQuotaRequestController` |
| `POST` | `/api/v1/admin/quota-requests/{id}/reject` | `a8s admin quota reject <request-id>` | `AdminQuotaRequestController` |
| `GET` | `/api/v1/admin/registry/health` | `a8s admin registry health` | `AdminRegistryController` |
| `GET` | `/api/v1/admin/registry/projects` | `a8s admin registry project list` | `AdminRegistryController` |
| `POST` | `/api/v1/admin/registry/projects` | `a8s admin registry project create` | `AdminRegistryController` |
| `DELETE` | `/api/v1/admin/registry/projects/{projectName}/artifacts` | `a8s admin registry artifact delete <project-name>` | `AdminRegistryController` |
| `GET` | `/api/v1/admin/registry/projects/{projectName}/artifacts` | `a8s admin registry artifact list <project-name>` | `AdminRegistryController` |
| `DELETE` | `/api/v1/admin/registry/projects/{projectName}/repositories` | `a8s admin registry repository delete <project-name>` | `AdminRegistryController` |
| `GET` | `/api/v1/admin/registry/projects/{projectName}/repositories` | `a8s admin registry repository list <project-name>` | `AdminRegistryController` |
| `GET` | `/api/v1/admin/sonarqube/projects` | `a8s admin sonarqube project list` | `AdminSonarQubeController` |
| `GET` | `/api/v1/admin/sonarqube/projects/{projectId}` | `a8s admin sonarqube project get <project-id>` | `AdminSonarQubeController` |
| `GET` | `/api/v1/admin/sonarqube/server-projects` | `a8s admin sonarqube server-project list` | `AdminSonarQubeController` |
| `POST` | `/api/v1/admin/sonarqube/server-projects` | `a8s admin sonarqube server-project create` | `AdminSonarQubeController` |
| `DELETE` | `/api/v1/admin/sonarqube/server-projects/{projectKey}` | `a8s admin sonarqube server-project delete <project-key>` | `AdminSonarQubeController` |
| `GET` | `/api/v1/admin/sonarqube/server-projects/{projectKey}` | `a8s admin sonarqube server-project get <project-key>` | `AdminSonarQubeController` |
| `PATCH` | `/api/v1/admin/sonarqube/server-projects/{projectKey}` | `a8s admin sonarqube server-project update <project-key>` | `AdminSonarQubeController` |

## alerts

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/alerts/channels` | `a8s alert channel list` | `AlertController` |
| `POST` | `/api/v1/alerts/channels` | `a8s alert channel create` | `AlertController` |
| `DELETE` | `/api/v1/alerts/channels/{channelId}` | `a8s alert channel delete <channel-id>` | `AlertController` |
| `PUT` | `/api/v1/alerts/channels/{channelId}` | `a8s alert channel update <channel-id>` | `AlertController` |
| `GET` | `/api/v1/alerts/projects/{projectId}/config` | `a8s alert project-config get <project-id>` | `AlertController` |
| `PUT` | `/api/v1/alerts/projects/{projectId}/config` | `a8s alert project-config set <project-id>` | `AlertController` |
| `GET` | `/api/v1/alerts/projects/configs` | `a8s alert project-config list` | `AlertController` |
| `GET` | `/api/v1/alerts/user-config` | `a8s alert user-config get` | `AlertController` |
| `PUT` | `/api/v1/alerts/user-config` | `a8s alert user-config set` | `AlertController` |

## auth

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/auth/keycloak/users/{keycloakUserId}/verify-email` | `a8s auth verify-email status` | `AuthController` |
| `POST` | `/api/v1/auth/keycloak/users/{keycloakUserId}/verify-email` | `a8s auth verify-email start` | `AuthController` |
| `GET` | `/api/v1/auth/session/onboarding` | `a8s auth onboarding status` | `AuthController` |
| `POST` | `/api/v1/auth/session/onboarding` | `a8s auth onboarding start` | `AuthController` |

## databasebackup

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/database-deployments/{deploymentId}/backup` | `a8s database backup settings get <deployment-id>` | `DatabaseDeploymentBackupController` |
| `PATCH` | `/api/v1/database-deployments/{deploymentId}/backup` | `a8s database backup settings set <deployment-id>` | `DatabaseDeploymentBackupController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/backup/run` | `a8s database backup run <deployment-id>` | `DatabaseDeploymentBackupController` |
| `DELETE` | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}` | `a8s database backup delete <deployment-id> <run-id>` | `DatabaseDeploymentBackupController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/download` | `a8s database backup download <deployment-id> <run-id>` | `DatabaseDeploymentBackupController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore` | `a8s database backup restore <deployment-id> <run-id>` | `DatabaseDeploymentBackupController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore/cancel` | `a8s database backup restore cancel <deployment-id> <run-id>` | `DatabaseDeploymentBackupController` |
| `POST` | `/api/internal/backups/callback/{deploymentId}/{targetType}` | `(internal service callback; no user CLI command)` | `InternalBackupController` |
| `DELETE` | `/api/backups/{targetType}/{id}/{runId}` | `a8s backup delete <type> <id> <run-id>` | `UnifiedBackupController` |
| `GET` | `/api/backups/download/{targetType}/{id}/{runId}` | `a8s backup download <type> <id> <run-id>` | `UnifiedBackupController` |
| `POST` | `/api/backups/restore/{targetType}/{id}/{runId}` | `a8s backup restore <type> <id> <run-id>` | `UnifiedBackupController` |
| `POST` | `/api/backups/restore/{targetType}/{id}/{runId}/cancel` | `a8s backup restore cancel <type> <id> <run-id>` | `UnifiedBackupController` |
| `GET` | `/api/backups/settings/{targetType}/{id}` | `a8s backup settings get <type> <id>` | `UnifiedBackupController` |
| `POST` | `/api/backups/settings/{targetType}/{id}` | `a8s backup settings set <type> <id>` | `UnifiedBackupController` |
| `POST` | `/api/backups/trigger/{targetType}/{id}` | `a8s backup trigger <type> <id>` | `UnifiedBackupController` |

## databaseconsole

No standalone controller. Console APIs are exposed through `singledb` and `dbcluster` endpoints.

## dbcluster

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/namespaces/{namespace}/clusters` | `a8s cluster list` | `ClusterController` |
| `DELETE` | `/api/namespaces/{namespace}/clusters/{id}` | `a8s cluster delete <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}` | `a8s cluster get <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/namespaces/{namespace}/clusters/{id}` | `a8s cluster update <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/namespaces/{namespace}/clusters/{id}/backup` | `a8s cluster backup settings set <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/certificate` | `a8s cluster certificate <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/console/credentials` | `a8s cluster console credentials <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/console/data` | `a8s cluster console data <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/console/deployment` | `a8s cluster console deployment <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/console/namespaces` | `a8s cluster console namespaces <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/console/objects` | `a8s cluster console objects <cluster-id>` | `ClusterController` |
| `POST` | `/api/namespaces/{namespace}/clusters/{id}/console/query` | `a8s cluster console query <cluster-id>` | `ClusterController` |
| `POST` | `/api/namespaces/{namespace}/clusters/{id}/console/test` | `a8s cluster console test <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/deployments` | `a8s cluster history <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/metrics` | `a8s cluster metrics <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/namespaces/{namespace}/clusters/{id}/settings` | `a8s cluster settings update <cluster-id>` | `ClusterController` |
| `POST` | `/api/namespaces/{namespace}/clusters/{id}/upgrade-version` | `a8s cluster upgrade <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/values` | `a8s cluster values <cluster-id>` | `ClusterController` |
| `GET` | `/api/namespaces/{namespace}/clusters/{id}/values/full` | `a8s cluster values <cluster-id> --full` | `ClusterController` |
| `POST` | `/api/namespaces/{namespace}/clusters/clone-from-backup` | `a8s cluster clone-from-backup` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters` | `a8s cluster list` | `ClusterController` |
| `DELETE` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}` | `a8s cluster delete <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}` | `a8s cluster get <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}` | `a8s cluster update <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/backup` | `a8s cluster backup settings set <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/certificate` | `a8s cluster certificate <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/credentials` | `a8s cluster console credentials <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/data` | `a8s cluster console data <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/deployment` | `a8s cluster console deployment <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/namespaces` | `a8s cluster console namespaces <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/objects` | `a8s cluster console objects <cluster-id>` | `ClusterController` |
| `POST` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/query` | `a8s cluster console query <cluster-id>` | `ClusterController` |
| `POST` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/console/test` | `a8s cluster console test <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/deployments` | `a8s cluster history <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/metrics` | `a8s cluster metrics <cluster-id>` | `ClusterController` |
| `PATCH` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/settings` | `a8s cluster settings update <cluster-id>` | `ClusterController` |
| `POST` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/upgrade-version` | `a8s cluster upgrade <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/values` | `a8s cluster values <cluster-id>` | `ClusterController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/clusters/{id}/values/full` | `a8s cluster values <cluster-id> --full` | `ClusterController` |
| `POST` | `/api/v1/cluster/namespaces/{namespace}/clusters/clone-from-backup` | `a8s cluster clone-from-backup` | `ClusterController` |
| `POST` | `/api/namespaces/{namespace}/cluster-deployments` | `a8s cluster deploy` | `ClusterDeploymentController` |
| `GET` | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}` | `a8s cluster status <release-name>` | `ClusterDeploymentController` |
| `PATCH` | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}/backup` | `a8s cluster backup settings set --release <release-name>` | `ClusterDeploymentController` |
| `GET` | `/api/namespaces/{namespace}/cluster-deployments/{releaseName}/values` | `a8s cluster deployment values <release-name>` | `ClusterDeploymentController` |
| `POST` | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments` | `a8s cluster deploy` | `ClusterDeploymentController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}` | `a8s cluster status <release-name>` | `ClusterDeploymentController` |
| `PATCH` | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}/backup` | `a8s cluster backup settings set --release <release-name>` | `ClusterDeploymentController` |
| `GET` | `/api/v1/cluster/namespaces/{namespace}/cluster-deployments/{releaseName}/values` | `a8s cluster deployment values <release-name>` | `ClusterDeploymentController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/database-resources` | `a8s kubernetes database-resources` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/events` | `a8s kubernetes events` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/overview` | `a8s kubernetes overview` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/persistent-volume-claims` | `a8s kubernetes pvc` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/pods` | `a8s kubernetes pods` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/pods/{podName}/logs/stream` | `a8s logs <pod-name> --follow` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream` | `a8s cluster watch <release-name>` | `KubernetesController` |
| `GET` | `/api/kubernetes/namespaces/{namespace}/services` | `a8s kubernetes services` | `KubernetesController` |
| `GET` | `/api/kubernetes/test` | `a8s kubernetes test` | `KubernetesController` |

## documentation

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `DELETE` | `/api/admin/documentation/content` | `a8s admin docs delete` | `DocumentationController` |
| `GET` | `/api/admin/documentation/content` | `a8s admin docs get` | `DocumentationController` |
| `PUT` | `/api/admin/documentation/content` | `a8s admin docs update` | `DocumentationController` |
| `GET` | `/api/admin/documentation/files` | `a8s admin docs files` | `DocumentationController` |
| `POST` | `/api/admin/documentation/publish` | `a8s admin docs publish` | `DocumentationController` |
| `GET` | `/api/admin/documentation/publish/logs` | `a8s admin docs publish-logs` | `DocumentationController` |

## entitlements

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/workspaces/entitlements` | `a8s workspace entitlements` | `WorkspaceEntitlementController` |

## gitintegration

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `DELETE` | `/api/v1/git-integrations/{provider}` | `a8s git disconnect <provider>` | `GitIntegrationController` |
| `GET` | `/api/v1/git-integrations/{provider}/brokered-account` | `a8s git account <provider>` | `GitIntegrationController` |
| `POST` | `/api/v1/git-integrations/{provider}/connect` | `a8s git connect <provider>` | `GitIntegrationController` |
| `GET` | `/api/v1/git-integrations/{provider}/repos` | `a8s git repos <provider>` | `GitIntegrationController` |
| `GET` | `/api/v1/git-integrations/{provider}/state` | `a8s git state <provider>` | `GitIntegrationController` |
| `POST` | `/api/v1/git-integrations/{provider}/sync-keycloak-token` | `a8s git sync-token <provider>` | `GitIntegrationController` |
| `GET` | `/api/v1/git-integrations/linked-providers` | `a8s git providers` | `GitIntegrationController` |

## imagescanner

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/image-scanner/images` | `a8s scan images` | `ImageScannerController` |
| `GET` | `/api/v1/image-scanner/scans` | `a8s scan list` | `ImageScannerController` |
| `POST` | `/api/v1/image-scanner/scans` | `a8s scan start` | `ImageScannerController` |
| `GET` | `/api/v1/image-scanner/scans/{scanId}` | `a8s scan get <scan-id>` | `ImageScannerController` |
| `GET` | `/api/v1/image-scanner/scans/{scanId}/report` | `a8s scan report <scan-id>` | `ImageScannerController` |
| `POST` | `/api/internal/image-scanner/callback/{scanId}` | `(internal service callback; no user CLI command)` | `InternalImageScannerCallbackController` |

## microservice

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/internal/microservices/source-archive` | `(internal service callback; no user CLI command)` | `InternalMicroserviceSourceArchiveController` |
| `GET` | `/api/internal/microservices/defectdojo-token` | `(internal service callback; no user CLI command)` | `InternalWorkspaceDefectDojoTokenController` |
| `POST` | `/api/v1/projects/microservices` | `a8s microservice deploy` | `MicroserviceProjectController` |
| `DELETE` | `/api/v1/projects/microservices/{projectId}` | `a8s microservice delete <project-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}` | `a8s microservice get <project-id>` | `MicroserviceProjectController` |
| `PUT` | `/api/v1/projects/microservices/{projectId}/canvas` | `a8s microservice apply <project-id>` | `MicroserviceProjectController` |
| `PATCH` | `/api/v1/projects/microservices/{projectId}/domains` | `a8s microservice domains update <project-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}/history` | `a8s microservice history list <project-id>` | `MicroserviceProjectController` |
| `DELETE` | `/api/v1/projects/microservices/{projectId}/history/{snapshotId}` | `a8s microservice history delete <project-id> <snapshot-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}/readiness` | `a8s microservice readiness <project-id>` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/{projectId}/redeploy` | `a8s microservice redeploy <project-id>` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/{projectId}/rollback` | `a8s microservice rollback <project-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}/runtime-pods` | `a8s microservice pods <project-id>` | `MicroserviceProjectController` |
| `DELETE` | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment` | `a8s microservice env clear <project-id> <service-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment` | `a8s microservice env get <project-id> <service-id>` | `MicroserviceProjectController` |
| `PUT` | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment` | `a8s microservice env set <project-id> <service-id>` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/{projectId}/services/{serviceId}/environment/import` | `a8s microservice env import <project-id> <service-id>` | `MicroserviceProjectController` |
| `GET` | `/api/v1/projects/microservices/{projectId}/webhook` | `a8s microservice webhook get <project-id>` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/{projectId}/webhook` | `a8s microservice webhook update <project-id>` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/detect` | `a8s microservice detect --repo` | `MicroserviceProjectController` |
| `POST` | `/api/v1/projects/microservices/detect/upload` | `a8s microservice detect --upload` | `MicroserviceProjectController` |

## monitoring

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/monitoring/overview` | `a8s monitoring overview` | `MonitoringController` |

## monolithic

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/projects` | `a8s project list` | `ProjectController` |
| `POST` | `/api/v1/projects` | `a8s project deploy` | `ProjectController` |
| `POST` | `/api/v1/projects` | `a8s project deploy` | `ProjectController` |
| `DELETE` | `/api/v1/projects/{projectId}` | `a8s project delete <project-id>` | `ProjectController` |
| `GET` | `/api/v1/projects/{projectId}` | `a8s project get <project-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/delete/complete` | `(Jenkins callback; no user CLI command)` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/delete/failed` | `(Jenkins callback; no user CLI command)` | `ProjectController` |
| `PATCH` | `/api/v1/projects/{projectId}/domain` | `a8s project domain set <project-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/domain/sync` | `a8s project domain sync <project-id>` | `ProjectController` |
| `GET` | `/api/v1/projects/{projectId}/releases` | `a8s project releases <project-id>` | `ProjectController` |
| `DELETE` | `/api/v1/projects/{projectId}/releases/{releaseId}` | `a8s project release delete <project-id> <release-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/releases/{releaseId}/complete` | `(Jenkins callback; no user CLI command)` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/releases/{releaseId}/failed` | `(Jenkins callback; no user CLI command)` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/releases/{releaseId}/rollback` | `a8s project release rollback <project-id> <release-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/repository/connect` | `a8s project repository connect <project-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/rollback` | `a8s project rollback <project-id>` | `ProjectController` |
| `PATCH` | `/api/v1/projects/{projectId}/settings` | `a8s project settings update <project-id>` | `ProjectController` |
| `POST` | `/api/v1/projects/{projectId}/sync` | `a8s project redeploy <project-id>` | `ProjectController` |
| `GET` | `/api/v1/projects/me` | `a8s project list` | `ProjectController` |
| `GET` | `/api/v1/projects/{projectId}/environment` | `a8s project env get <project-id>` | `ProjectEnvironmentController` |
| `PUT` | `/api/v1/projects/{projectId}/environment` | `a8s project env set <project-id>` | `ProjectEnvironmentController` |
| `POST` | `/api/v1/projects/{projectId}/environment/import` | `a8s project env import <project-id>` | `ProjectEnvironmentController` |
| `PATCH` | `/api/v1/projects/{projectId}/auto-deploy` | `a8s project auto-deploy set <project-id>` | `ProjectWebhookController` |
| `GET` | `/api/v1/projects/{projectId}/branches` | `a8s project branches <project-id>` | `ProjectWebhookController` |
| `DELETE` | `/api/v1/projects/{projectId}/webhook` | `a8s project webhook delete <project-id>` | `ProjectWebhookController` |
| `GET` | `/api/v1/projects/{projectId}/webhook` | `a8s project webhook get <project-id>` | `ProjectWebhookController` |
| `POST` | `/api/v1/projects/{projectId}/webhook` | `a8s project webhook create <project-id>` | `ProjectWebhookController` |
| `POST` | `/api/v1/projects/{projectId}/webhook/rotate` | `a8s project webhook rotate <project-id>` | `ProjectWebhookController` |
| `POST` | `/api/v1/webhooks/github` | `(provider webhook receiver; no user CLI command)` | `WebhookController` |
| `POST` | `/api/v1/webhooks/gitlab` | `(provider webhook receiver; no user CLI command)` | `WebhookController` |

## notifications

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `POST` | `/api/notifications/{notificationId}/read` | `a8s notification read <notification-id>` | `NotificationController` |
| `GET` | `/api/notifications/history/{userId}` | `a8s notification list` | `NotificationController` |
| `GET` | `/api/notifications/preferences/{userId}` | `a8s notification preferences get` | `NotificationController` |
| `POST` | `/api/notifications/preferences/{userId}` | `a8s notification preferences set` | `NotificationController` |

## payments

No standalone payment controller. Payments currently support Bakong KHQR purchases for workspace quota and plan upgrades.

### Payment and quota-purchase endpoints

| Method | Endpoint | Suggested CLI command | Purpose |
|---|---|---|---|
| `GET` | `/api/v1/workspaces/quota-pricing` | `a8s workspace quota pricing` | Get unit prices and plan prices. |
| `POST` | `/api/v1/workspaces/quota-requests` | `a8s workspace quota purchase --plan <plan>` | Submit a paid quota request and generate a Bakong KHQR payload. |
| `GET` | `/api/v1/workspaces/quota-requests/payment-status?md5=<md5>` | `a8s workspace quota payment-status <md5>` | Check payment status and apply the quota upgrade after payment. |

The purchase request accepts `requestedCpu`, `requestedMemory`, `requestedStorage`, `reason`, `isPaid`, `planName`, and `paymentProvider`. Set `isPaid` to `true` and `paymentProvider` to `BAKONG` to generate KHQR.

The purchase response contains `qrString` and `md5`. Use the returned `md5` when polling payment status. Status responses currently include `PENDING`, `PAID`, and `NO_PAYMENT_REQUIRED`.

When payment is confirmed, the backend approves the quota request, applies the workspace quota, activates the subscription for 30 days, and sends a payment receipt notification.

### Related admin endpoints

| Method | Endpoint | Suggested CLI command | Purpose |
|---|---|---|---|
| `GET` | `/api/v1/admin/quota-requests` | `a8s admin quota list` | List pending quota and payment-related requests. |
| `POST` | `/api/v1/admin/quota-requests/{id}/approve` | `a8s admin quota approve <request-id>` | Approve a pending request and apply its quota. |
| `POST` | `/api/v1/admin/quota-requests/{id}/reject` | `a8s admin quota reject <request-id>` | Reject a pending request. |

## profile

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `DELETE` | `/api/v1/profile/me` | `a8s profile account delete` | `ProfileController` |
| `GET` | `/api/v1/profile/me` | `a8s profile get` | `ProfileController` |
| `PATCH` | `/api/v1/profile/me` | `a8s profile update` | `ProfileController` |
| `GET` | `/api/v1/profile/me/account-status` | `a8s profile account status` | `ProfileController` |
| `DELETE` | `/api/v1/profile/me/avatar` | `a8s profile avatar delete` | `ProfileController` |
| `GET` | `/api/v1/profile/me/avatar` | `a8s profile avatar download` | `ProfileController` |
| `POST` | `/api/v1/profile/me/avatar` | `a8s profile avatar upload` | `ProfileController` |
| `POST` | `/api/v1/profile/me/deactivate` | `a8s profile account deactivate` | `ProfileController` |
| `POST` | `/api/v1/profile/me/delete` | `a8s profile account delete` | `ProfileController` |
| `POST` | `/api/v1/profile/me/reactivate` | `a8s profile account reactivate` | `ProfileController` |

## projects

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/jenkins/logs/stream` | `a8s project logs --follow` | `JenkinsController` |
| `GET` | `/api/v1/projects/live` | `a8s project live list` | `LiveProjectController` |
| `GET` | `/api/v1/projects/{projectId}/defectdojo` | `a8s defectdojo access <project-id>` | `ProjectDefectDojoController` |
| `PUT` | `/api/v1/projects/{projectId}/defectdojo/token` | `a8s defectdojo token sync <project-id>` | `ProjectDefectDojoController` |

## singledb

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/database-deployments` | `a8s database list` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments` | `a8s database deploy` | `DatabaseDeploymentController` |
| `DELETE` | `/api/v1/database-deployments/{deploymentId}` | `a8s database delete <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}` | `a8s database get <deployment-id>` | `DatabaseDeploymentController` |
| `PATCH` | `/api/v1/database-deployments/{deploymentId}` | `a8s database update <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/console/data` | `a8s database console data <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/console/namespaces` | `a8s database console namespaces <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/console/objects` | `a8s database console objects <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/console/query` | `a8s database console query <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/console/test` | `a8s database console test <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/credentials` | `a8s database credentials <deployment-id>` | `DatabaseDeploymentController` |
| `GET` | `/api/v1/database-deployments/{deploymentId}/metrics` | `a8s database metrics <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/restart` | `a8s database restart <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/rotate-password` | `a8s database rotate-password <deployment-id>` | `DatabaseDeploymentController` |
| `PATCH` | `/api/v1/database-deployments/{deploymentId}/settings` | `a8s database settings update <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/upgrade-version` | `a8s database upgrade <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/{deploymentId}/verify-password` | `a8s database verify-password <deployment-id>` | `DatabaseDeploymentController` |
| `POST` | `/api/v1/database-deployments/clone-from-backup` | `a8s database clone-from-backup` | `DatabaseDeploymentController` |

## sonarqube

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/projects/{projectId}/sonarqube` | `a8s sonarqube summary <project-id>` | `SonarQubeProjectController` |
| `POST` | `/api/v1/projects/{projectId}/sonarqube/access` | `a8s sonarqube access <project-id>` | `SonarQubeProjectController` |

## testingkit

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `POST` | `/api/v1/projects/live/{projectId}/benchmark/run` | `a8s benchmark run <project-id>` | `BenchmarkController` |
| `GET` | `/api/v1/projects/live/{projectId}/benchmark/runs` | `a8s benchmark list <project-id>` | `BenchmarkController` |
| `DELETE` | `/api/v1/projects/live/{projectId}/benchmark/runs/{runId}` | `a8s benchmark delete <project-id> <run-id>` | `BenchmarkController` |
| `GET` | `/api/v1/projects/live/{projectId}/benchmark/runs/{runId}` | `a8s benchmark get <project-id> <run-id>` | `BenchmarkController` |

## workspaces

| Method | Endpoint | Suggested CLI command | Controller |
|---|---|---|---|
| `GET` | `/api/v1/workspaces/bootstrap` | `a8s workspace status` | `WorkspaceController` |
| `POST` | `/api/v1/workspaces/bootstrap` | `a8s workspace bootstrap` | `WorkspaceController` |
| `GET` | `/api/v1/workspaces/quota-pricing` | `a8s workspace quota pricing` | `WorkspaceController` |
| `POST` | `/api/v1/workspaces/quota-requests` | `a8s workspace quota request` | `WorkspaceController` |
| `GET` | `/api/v1/workspaces/quota-requests/payment-status` | `a8s workspace quota payment-status` | `WorkspaceController` |

## WebSockets

| Endpoint | Suggested CLI use |
|---|---|
| `/ws/jenkins/logs` | `a8s project logs --follow` |
| `/ws/notifications` | `a8s notification watch` |
| `/ws/monitoring/overview` | `a8s monitoring watch` |
| `/ws/admin/events` | `a8s admin events watch` |

## CLI Exclusions

Do not expose provider webhook receivers, Jenkins completion/failure callbacks, or `/api/internal/**` endpoints as ordinary user commands. They are automation-to-backend contracts.
