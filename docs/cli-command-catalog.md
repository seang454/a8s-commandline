| Command | Usage | Description | Input |
|---|---|---|---|
| a8s admin cluster health | a8s admin cluster health <alias> [flags] | GET /api/v1/admin/clusters/kubernetes/{alias}/health | flags |
| a8s admin cluster list | a8s admin cluster list [flags] | GET /api/v1/admin/clusters | flags |
| a8s admin cluster nodes | a8s admin cluster nodes [flags] | GET /api/v1/admin/clusters/kubernetes | flags |
| a8s admin cluster quota delete | a8s admin cluster quota delete <alias> <namespace> [flags] | DELETE /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace} | flags |
| a8s admin cluster quota list | a8s admin cluster quota list <alias> [flags] | GET /api/v1/admin/clusters/kubernetes/{alias}/quotas | flags |
| a8s admin cluster quota set | a8s admin cluster quota set <alias> <namespace> [flags] | PUT /api/v1/admin/clusters/kubernetes/{alias}/quotas/{namespace} | flags, file |
| a8s admin cluster update | a8s admin cluster update <cluster-id> [flags] | PATCH /api/v1/admin/clusters/{clusterId} | flags, file |
| a8s admin docs delete | a8s admin docs delete [flags] | DELETE /api/admin/documentation/content | flags |
| a8s admin docs files | a8s admin docs files [flags] | GET /api/admin/documentation/files | flags, file |
| a8s admin docs get | a8s admin docs get [flags] | GET /api/admin/documentation/content | flags |
| a8s admin docs publish | a8s admin docs publish [flags] | POST /api/admin/documentation/publish | flags |
| a8s admin docs publish-logs | a8s admin docs publish-logs [flags] | GET /api/admin/documentation/publish/logs | flags |
| a8s admin docs update | a8s admin docs update [flags] | PUT /api/admin/documentation/content | flags, file |
| a8s admin events watch | a8s admin events watch | Watch administrative events | flags |
| a8s admin gitops app abort | a8s admin gitops app abort <app-id> [flags] | POST /api/v1/admin/gitops/apps/{appId}/abort | flags |
| a8s admin gitops app create | a8s admin gitops app create [flags] | POST /api/v1/admin/gitops/apps | flags, file |
| a8s admin gitops app retry | a8s admin gitops app retry <app-id> [flags] | POST /api/v1/admin/gitops/apps/{appId}/retry | flags |
| a8s admin gitops app sync | a8s admin gitops app sync <app-id> [flags] | POST /api/v1/admin/gitops/apps/{appId}/sync | flags |
| a8s admin gitops overview | a8s admin gitops overview [flags] | GET /api/v1/admin/gitops/overview | flags |
| a8s admin logs clusters | a8s admin logs clusters [flags] | GET /api/v1/admin/logs/clusters | flags |
| a8s admin logs namespaces | a8s admin logs namespaces [flags] | GET /api/v1/admin/logs/namespaces | flags |
| a8s admin logs pods | a8s admin logs pods [flags] | GET /api/v1/admin/logs/pods | flags |
| a8s admin logs query | a8s admin logs query [flags] | GET /api/v1/admin/logs/query | flags |
| a8s admin logs workloads | a8s admin logs workloads [flags] | GET /api/v1/admin/logs/workloads | flags |
| a8s admin monitoring overview | a8s admin monitoring overview [flags] | GET /api/v1/admin/monitoring/overview | flags |
| a8s admin project deactivate | a8s admin project deactivate <project-id> [flags] | DELETE /api/v1/admin/projects/{projectId} | flags |
| a8s admin project list | a8s admin project list [flags] | GET /api/v1/admin/projects | flags |
| a8s admin project restore | a8s admin project restore <project-id> [flags] | POST /api/v1/admin/projects/{projectId}/restore | flags, file |
| a8s admin project update | a8s admin project update <project-id> [flags] | PATCH /api/v1/admin/projects/{projectId} | flags, file |
| a8s admin quota approve | a8s admin quota approve <request-id> [flags] | POST /api/v1/admin/quota-requests/{id}/approve | flags |
| a8s admin quota list | a8s admin quota list [flags] | GET /api/v1/admin/quota-requests | flags |
| a8s admin quota reject | a8s admin quota reject <request-id> [flags] | POST /api/v1/admin/quota-requests/{id}/reject | flags |
| a8s admin registry artifact delete | a8s admin registry artifact delete <project-name> [flags] | DELETE /api/v1/admin/registry/projects/{projectName}/artifacts | flags |
| a8s admin registry artifact list | a8s admin registry artifact list <project-name> [flags] | GET /api/v1/admin/registry/projects/{projectName}/artifacts | flags |
| a8s admin registry health | a8s admin registry health [flags] | GET /api/v1/admin/registry/health | flags |
| a8s admin registry project create | a8s admin registry project create [flags] | POST /api/v1/admin/registry/projects | flags, file |
| a8s admin registry project list | a8s admin registry project list [flags] | GET /api/v1/admin/registry/projects | flags |
| a8s admin registry repository delete | a8s admin registry repository delete <project-name> [flags] | DELETE /api/v1/admin/registry/projects/{projectName}/repositories | flags |
| a8s admin registry repository list | a8s admin registry repository list <project-name> [flags] | GET /api/v1/admin/registry/projects/{projectName}/repositories | flags |
| a8s admin sonarqube project get | a8s admin sonarqube project get <project-id> [flags] | GET /api/v1/admin/sonarqube/projects/{projectId} | flags |
| a8s admin sonarqube project list | a8s admin sonarqube project list [flags] | GET /api/v1/admin/sonarqube/projects | flags |
| a8s admin sonarqube server-project create | a8s admin sonarqube server-project create [flags] | POST /api/v1/admin/sonarqube/server-projects | flags, file |
| a8s admin sonarqube server-project delete | a8s admin sonarqube server-project delete <project-key> [flags] | DELETE /api/v1/admin/sonarqube/server-projects/{projectKey} | flags |
| a8s admin sonarqube server-project get | a8s admin sonarqube server-project get <project-key> [flags] | GET /api/v1/admin/sonarqube/server-projects/{projectKey} | flags |
| a8s admin sonarqube server-project list | a8s admin sonarqube server-project list [flags] | GET /api/v1/admin/sonarqube/server-projects | flags |
| a8s admin sonarqube server-project update | a8s admin sonarqube server-project update <project-key> [flags] | PATCH /api/v1/admin/sonarqube/server-projects/{projectKey} | flags, file |
| a8s admin user create | a8s admin user create [flags] | POST /api/v1/admin/users | flags, file |
| a8s admin user deactivate | a8s admin user deactivate <user-id> [flags] | DELETE /api/v1/admin/users/{userId} | flags |
| a8s admin user list | a8s admin user list [flags] | GET /api/v1/admin/users | flags |
| a8s admin user reactivate | a8s admin user reactivate <user-id> [flags] | POST /api/v1/admin/users/{userId}/reactivate | flags |
| a8s admin user update | a8s admin user update <user-id> [flags] | PATCH /api/v1/admin/users/{userId} | flags, file |
| a8s alert channel create | a8s alert channel create [flags] | POST /api/v1/alerts/channels | flags, file |
| a8s alert channel delete | a8s alert channel delete <channel-id> [flags] | DELETE /api/v1/alerts/channels/{channelId} | flags |
| a8s alert channel list | a8s alert channel list [flags] | GET /api/v1/alerts/channels | flags |
| a8s alert channel update | a8s alert channel update <channel-id> [flags] | PUT /api/v1/alerts/channels/{channelId} | flags, file |
| a8s alert project-config get | a8s alert project-config get <project-id> [flags] | GET /api/v1/alerts/projects/{projectId}/config | flags |
| a8s alert project-config list | a8s alert project-config list [flags] | GET /api/v1/alerts/projects/configs | flags |
| a8s alert project-config set | a8s alert project-config set <project-id> [flags] | PUT /api/v1/alerts/projects/{projectId}/config | flags, file |
| a8s alert user-config get | a8s alert user-config get [flags] | GET /api/v1/alerts/user-config | flags |
| a8s alert user-config set | a8s alert user-config set [flags] | PUT /api/v1/alerts/user-config | flags, file |
| a8s api catalog | a8s api catalog [flags] | List implemented backend route mappings | flags |
| a8s api request | a8s api request <method> <path> [flags] | Send an authenticated request to any backend route | flags |
| a8s auth login | a8s auth login [flags] | Authenticate through Keycloak using browser PKCE | flags |
| a8s auth logout | a8s auth logout [flags] | Clear stored credentials for the active context | flags |
| a8s auth onboarding start | a8s auth onboarding start [flags] | POST /api/v1/auth/session/onboarding | flags |
| a8s auth onboarding status | a8s auth onboarding status [flags] | GET /api/v1/auth/session/onboarding | flags |
| a8s auth status | a8s auth status | Show authentication status without displaying tokens | flags |
| a8s auth verify-email start | a8s auth verify-email start [keycloak-user-id] | Start email verification for the authenticated user | flags |
| a8s auth verify-email status | a8s auth verify-email status [keycloak-user-id] | Show email verification status | flags |
| a8s backup delete | a8s backup delete <type> <id> <run-id> [flags] | DELETE /api/backups/{targetType}/{id}/{runId} | flags |
| a8s backup download | a8s backup download <type> <id> <run-id> [flags] | GET /api/backups/download/{targetType}/{id}/{runId} | flags, file |
| a8s backup restore | a8s backup restore <type> <id> <run-id> [flags] | POST /api/backups/restore/{targetType}/{id}/{runId} | flags, file |
| a8s backup restore cancel | a8s backup restore cancel <type> <id> <run-id> [flags] | POST /api/backups/restore/{targetType}/{id}/{runId}/cancel | flags, file |
| a8s backup settings get | a8s backup settings get <type> <id> [flags] | GET /api/backups/settings/{targetType}/{id} | flags, file |
| a8s backup settings set | a8s backup settings set <type> <id> [flags] | POST /api/backups/settings/{targetType}/{id} | flags, file |
| a8s backup trigger | a8s backup trigger <type> <id> [flags] | POST /api/backups/trigger/{targetType}/{id} | flags |
| a8s benchmark delete | a8s benchmark delete <project-id> <run-id> [flags] | DELETE /api/v1/projects/live/{projectId}/benchmark/runs/{runId} | flags |
| a8s benchmark get | a8s benchmark get <project-id> <run-id> [flags] | GET /api/v1/projects/live/{projectId}/benchmark/runs/{runId} | flags |
| a8s benchmark list | a8s benchmark list <project-id> [flags] | GET /api/v1/projects/live/{projectId}/benchmark/runs | flags |
| a8s benchmark run | a8s benchmark run <project-id> [flags] | POST /api/v1/projects/live/{projectId}/benchmark/run | flags |
| a8s cluster backup settings set | a8s cluster backup settings set <release-name> [flags] | PATCH /api/namespaces/{namespace}/cluster-deployments/{releaseName}/backup | flags, file |
| a8s cluster certificate | a8s cluster certificate <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/certificate | flags |
| a8s cluster clone-from-backup | a8s cluster clone-from-backup [flags] | POST /api/namespaces/{namespace}/clusters/clone-from-backup | flags |
| a8s cluster console credentials | a8s cluster console credentials <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/console/credentials | flags |
| a8s cluster console data | a8s cluster console data <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/console/data | flags |
| a8s cluster console deployment | a8s cluster console deployment <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/console/deployment | flags, file |
| a8s cluster console namespaces | a8s cluster console namespaces <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/console/namespaces | flags |
| a8s cluster console objects | a8s cluster console objects <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/console/objects | flags |
| a8s cluster console query | a8s cluster console query <cluster-id> [flags] | POST /api/namespaces/{namespace}/clusters/{id}/console/query | flags |
| a8s cluster console test | a8s cluster console test <cluster-id> [flags] | POST /api/namespaces/{namespace}/clusters/{id}/console/test | flags |
| a8s cluster delete | a8s cluster delete <cluster-id> [flags] | DELETE /api/namespaces/{namespace}/clusters/{id} | flags |
| a8s cluster deploy | a8s cluster deploy [flags] | POST /api/namespaces/{namespace}/cluster-deployments | flags, file |
| a8s cluster deployment values | a8s cluster deployment values <release-name> [flags] | GET /api/namespaces/{namespace}/cluster-deployments/{releaseName}/values | flags, file |
| a8s cluster get | a8s cluster get <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id} | flags |
| a8s cluster history | a8s cluster history <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/deployments | flags, file |
| a8s cluster list | a8s cluster list [flags] | GET /api/namespaces/{namespace}/clusters | flags |
| a8s cluster metrics | a8s cluster metrics <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/metrics | flags |
| a8s cluster settings update | a8s cluster settings update <cluster-id> [flags] | PATCH /api/namespaces/{namespace}/clusters/{id}/settings | flags, file |
| a8s cluster status | a8s cluster status <release-name> [flags] | GET /api/namespaces/{namespace}/cluster-deployments/{releaseName} | flags, file |
| a8s cluster update | a8s cluster update <cluster-id> [flags] | PATCH /api/namespaces/{namespace}/clusters/{id} | flags, file |
| a8s cluster upgrade | a8s cluster upgrade <cluster-id> [flags] | POST /api/namespaces/{namespace}/clusters/{id}/upgrade-version | flags |
| a8s cluster values | a8s cluster values <cluster-id> [flags] | GET /api/namespaces/{namespace}/clusters/{id}/values | flags |
| a8s cluster watch | a8s cluster watch <release-name> [flags] | GET /api/kubernetes/namespaces/{namespace}/releases/{releaseName}/deployment-stream | flags, file |
| a8s completion bash | a8s completion bash | Generate the autocompletion script for bash | flags |
| a8s completion fish | a8s completion fish [flags] | Generate the autocompletion script for fish | flags |
| a8s completion powershell | a8s completion powershell [flags] | Generate the autocompletion script for powershell | flags |
| a8s completion zsh | a8s completion zsh [flags] | Generate the autocompletion script for zsh | flags |
| a8s config path | a8s config path | Print the active configuration path | flags |
| a8s config view | a8s config view | Print resolved non-secret configuration | flags |
| a8s context create | a8s context create <name> [flags] | Create a named context | flags, file |
| a8s context delete | a8s context delete <name> [flags] | Delete a named context | flags |
| a8s context get | a8s context get <name> | Get a configured context | flags |
| a8s context list | a8s context list | List configured contexts | flags |
| a8s context update | a8s context update <name> [flags] | Update a named context | flags, file |
| a8s context use | a8s context use <name> | Set the default context | flags, file |
| a8s database backup delete | a8s database backup delete <deployment-id> <run-id> [flags] | DELETE /api/v1/database-deployments/{deploymentId}/backup/runs/{runId} | flags, file |
| a8s database backup download | a8s database backup download <deployment-id> <run-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/download | flags, file |
| a8s database backup restore | a8s database backup restore <deployment-id> <run-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore | flags, file |
| a8s database backup restore cancel | a8s database backup restore cancel <deployment-id> <run-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/backup/runs/{runId}/restore/cancel | flags, file |
| a8s database backup run | a8s database backup run <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/backup/run | flags, file |
| a8s database backup settings get | a8s database backup settings get <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/backup | flags, file |
| a8s database backup settings set | a8s database backup settings set <deployment-id> [flags] | PATCH /api/v1/database-deployments/{deploymentId}/backup | flags, file |
| a8s database clone-from-backup | a8s database clone-from-backup [flags] | POST /api/v1/database-deployments/clone-from-backup | flags, file |
| a8s database console data | a8s database console data <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/console/data | flags, file |
| a8s database console namespaces | a8s database console namespaces <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/console/namespaces | flags, file |
| a8s database console objects | a8s database console objects <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/console/objects | flags, file |
| a8s database console query | a8s database console query <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/console/query | flags, file |
| a8s database console test | a8s database console test <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/console/test | flags, file |
| a8s database credentials | a8s database credentials <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/credentials | flags, file |
| a8s database delete | a8s database delete <deployment-id> [flags] | DELETE /api/v1/database-deployments/{deploymentId} | flags, file |
| a8s database deploy | a8s database deploy [flags] | Deploy a single database using flags or an operation file | flags, file |
| a8s database get | a8s database get <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId} | flags, file |
| a8s database list | a8s database list [flags] | GET /api/v1/database-deployments | flags, file |
| a8s database metrics | a8s database metrics <deployment-id> [flags] | GET /api/v1/database-deployments/{deploymentId}/metrics | flags, file |
| a8s database restart | a8s database restart <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/restart | flags, file |
| a8s database rotate-password | a8s database rotate-password <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/rotate-password | flags, file |
| a8s database settings update | a8s database settings update <deployment-id> [flags] | PATCH /api/v1/database-deployments/{deploymentId}/settings | flags, file |
| a8s database update | a8s database update <deployment-id> [flags] | PATCH /api/v1/database-deployments/{deploymentId} | flags, file |
| a8s database upgrade | a8s database upgrade <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/upgrade-version | flags, file |
| a8s database verify-password | a8s database verify-password <deployment-id> [flags] | POST /api/v1/database-deployments/{deploymentId}/verify-password | flags, file |
| a8s defectdojo access | a8s defectdojo access <project-id> [flags] | GET /api/v1/projects/{projectId}/defectdojo | flags |
| a8s defectdojo token sync | a8s defectdojo token sync <project-id> [flags] | PUT /api/v1/projects/{projectId}/defectdojo/token | flags |
| a8s doctor | a8s doctor | Check CLI configuration and backend connectivity | flags |
| a8s features | a8s features | List backend features exposed by the CLI | flags |
| a8s git account | a8s git account <provider> [flags] | GET /api/v1/git-integrations/{provider}/brokered-account | flags |
| a8s git connect | a8s git connect <provider> [flags] | POST /api/v1/git-integrations/{provider}/connect | flags |
| a8s git disconnect | a8s git disconnect <provider> [flags] | DELETE /api/v1/git-integrations/{provider} | flags |
| a8s git providers | a8s git providers [flags] | GET /api/v1/git-integrations/linked-providers | flags |
| a8s git repos | a8s git repos <provider> [flags] | GET /api/v1/git-integrations/{provider}/repos | flags |
| a8s git state | a8s git state <provider> [flags] | GET /api/v1/git-integrations/{provider}/state | flags |
| a8s git sync-token | a8s git sync-token <provider> [flags] | POST /api/v1/git-integrations/{provider}/sync-keycloak-token | flags |
| a8s help | a8s help [command] | Help about any command | flags |
| a8s kubernetes database-resources | a8s kubernetes database-resources [flags] | GET /api/kubernetes/namespaces/{namespace}/database-resources | flags |
| a8s kubernetes events | a8s kubernetes events [flags] | GET /api/kubernetes/namespaces/{namespace}/events | flags |
| a8s kubernetes overview | a8s kubernetes overview [flags] | GET /api/kubernetes/namespaces/{namespace}/overview | flags |
| a8s kubernetes pods | a8s kubernetes pods [flags] | GET /api/kubernetes/namespaces/{namespace}/pods | flags |
| a8s kubernetes pvc | a8s kubernetes pvc [flags] | GET /api/kubernetes/namespaces/{namespace}/persistent-volume-claims | flags |
| a8s kubernetes services | a8s kubernetes services [flags] | GET /api/kubernetes/namespaces/{namespace}/services | flags |
| a8s kubernetes test | a8s kubernetes test [flags] | GET /api/kubernetes/test | flags |
| a8s list all | a8s list all [flags] | List every available runnable command | flags |
| a8s list sections | a8s list sections [flags] | List commands grouped by top-level section | flags |
| a8s logs | a8s logs <pod-name> [flags] | GET /api/kubernetes/namespaces/{namespace}/pods/{podName}/logs/stream | flags |
| a8s manifest init | a8s manifest init <kind> [flags] | Generate a starter manifest for a kind | flags, file |
| a8s manifest kinds | a8s manifest kinds | List supported operation manifest kinds | flags, file |
| a8s manifest schema | a8s manifest schema <kind> | Show the manifest schema summary for a kind | flags, file |
| a8s manifest validate | a8s manifest validate [flags] | Validate an operation manifest without sending a backend request | flags, file |
| a8s microservice apply | a8s microservice apply <project-id> [flags] | PUT /api/v1/projects/microservices/{projectId}/canvas | flags, file |
| a8s microservice delete | a8s microservice delete <project-id> [flags] | DELETE /api/v1/projects/microservices/{projectId} | flags |
| a8s microservice deploy | a8s microservice deploy [flags] | POST /api/v1/projects/microservices | flags, file |
| a8s microservice detect | a8s microservice detect [flags] | POST /api/v1/projects/microservices/detect | flags |
| a8s microservice domains update | a8s microservice domains update <project-id> [flags] | PATCH /api/v1/projects/microservices/{projectId}/domains | flags, file |
| a8s microservice env clear | a8s microservice env clear <project-id> <service-id> [flags] | DELETE /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment | flags |
| a8s microservice env get | a8s microservice env get <project-id> <service-id> [flags] | GET /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment | flags |
| a8s microservice env import | a8s microservice env import <project-id> <service-id> [flags] | POST /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment/import | flags, file |
| a8s microservice env set | a8s microservice env set <project-id> <service-id> [flags] | PUT /api/v1/projects/microservices/{projectId}/services/{serviceId}/environment | flags, file |
| a8s microservice get | a8s microservice get <project-id> [flags] | GET /api/v1/projects/microservices/{projectId} | flags |
| a8s microservice history delete | a8s microservice history delete <project-id> <snapshot-id> [flags] | DELETE /api/v1/projects/microservices/{projectId}/history/{snapshotId} | flags |
| a8s microservice history list | a8s microservice history list <project-id> [flags] | GET /api/v1/projects/microservices/{projectId}/history | flags |
| a8s microservice pods | a8s microservice pods <project-id> [flags] | GET /api/v1/projects/microservices/{projectId}/runtime-pods | flags |
| a8s microservice readiness | a8s microservice readiness <project-id> [flags] | GET /api/v1/projects/microservices/{projectId}/readiness | flags |
| a8s microservice redeploy | a8s microservice redeploy <project-id> [flags] | POST /api/v1/projects/microservices/{projectId}/redeploy | flags, file |
| a8s microservice rollback | a8s microservice rollback <project-id> [flags] | POST /api/v1/projects/microservices/{projectId}/rollback | flags |
| a8s microservice webhook get | a8s microservice webhook get <project-id> [flags] | GET /api/v1/projects/microservices/{projectId}/webhook | flags |
| a8s microservice webhook update | a8s microservice webhook update <project-id> [flags] | POST /api/v1/projects/microservices/{projectId}/webhook | flags, file |
| a8s monitoring overview | a8s monitoring overview [flags] | GET /api/v1/monitoring/overview | flags |
| a8s monitoring watch | a8s monitoring watch | Watch monitoring updates | flags, file |
| a8s notification list | a8s notification list <user-id> [flags] | GET /api/notifications/history/{userId} | flags |
| a8s notification preferences get | a8s notification preferences get <user-id> [flags] | GET /api/notifications/preferences/{userId} | flags |
| a8s notification preferences set | a8s notification preferences set <user-id> [flags] | POST /api/notifications/preferences/{userId} | flags, file |
| a8s notification read | a8s notification read <notification-id> [flags] | POST /api/notifications/{notificationId}/read | flags |
| a8s notification watch | a8s notification watch | Watch notifications | flags |
| a8s profile account deactivate | a8s profile account deactivate [flags] | POST /api/v1/profile/me/deactivate | flags, file |
| a8s profile account delete | a8s profile account delete [flags] | DELETE /api/v1/profile/me | flags, file |
| a8s profile account reactivate | a8s profile account reactivate [flags] | POST /api/v1/profile/me/reactivate | flags, file |
| a8s profile account status | a8s profile account status [flags] | GET /api/v1/profile/me/account-status | flags, file |
| a8s profile avatar delete | a8s profile avatar delete [flags] | DELETE /api/v1/profile/me/avatar | flags, file |
| a8s profile avatar download | a8s profile avatar download [flags] | GET /api/v1/profile/me/avatar | flags, file |
| a8s profile avatar upload | a8s profile avatar upload [flags] | POST /api/v1/profile/me/avatar | flags, file |
| a8s profile get | a8s profile get [flags] | GET /api/v1/profile/me | flags, file |
| a8s profile update | a8s profile update [flags] | PATCH /api/v1/profile/me | flags, file |
| a8s project auto-deploy set | a8s project auto-deploy set <project-id> [flags] | PATCH /api/v1/projects/{projectId}/auto-deploy | flags, file |
| a8s project branches | a8s project branches <project-id> [flags] | GET /api/v1/projects/{projectId}/branches | flags |
| a8s project delete | a8s project delete <project-id> [flags] | DELETE /api/v1/projects/{projectId} | flags |
| a8s project deploy | a8s project deploy [flags] | POST /api/v1/projects | flags, file |
| a8s project domain set | a8s project domain set <project-id> [flags] | PATCH /api/v1/projects/{projectId}/domain | flags, file |
| a8s project domain sync | a8s project domain sync <project-id> [flags] | POST /api/v1/projects/{projectId}/domain/sync | flags |
| a8s project env get | a8s project env get <project-id> [flags] | GET /api/v1/projects/{projectId}/environment | flags |
| a8s project env import | a8s project env import <project-id> [flags] | POST /api/v1/projects/{projectId}/environment/import | flags, file |
| a8s project env set | a8s project env set <project-id> [flags] | PUT /api/v1/projects/{projectId}/environment | flags, file |
| a8s project get | a8s project get <project-id> [flags] | GET /api/v1/projects/{projectId} | flags |
| a8s project list | a8s project list [flags] | GET /api/v1/projects | flags |
| a8s project live list | a8s project live list [flags] | GET /api/v1/projects/live | flags |
| a8s project logs | a8s project logs [flags] | GET /api/v1/jenkins/logs/stream | flags |
| a8s project logs websocket | a8s project logs websocket | Watch Jenkins logs over WebSocket | flags |
| a8s project redeploy | a8s project redeploy <project-id> [flags] | POST /api/v1/projects/{projectId}/sync | flags, file |
| a8s project release delete | a8s project release delete <project-id> <release-id> [flags] | DELETE /api/v1/projects/{projectId}/releases/{releaseId} | flags |
| a8s project release rollback | a8s project release rollback <project-id> <release-id> [flags] | POST /api/v1/projects/{projectId}/releases/{releaseId}/rollback | flags |
| a8s project releases | a8s project releases <project-id> [flags] | GET /api/v1/projects/{projectId}/releases | flags |
| a8s project repository connect | a8s project repository connect <project-id> [flags] | POST /api/v1/projects/{projectId}/repository/connect | flags |
| a8s project rollback | a8s project rollback <project-id> [flags] | POST /api/v1/projects/{projectId}/rollback | flags |
| a8s project settings update | a8s project settings update <project-id> [flags] | PATCH /api/v1/projects/{projectId}/settings | flags, file |
| a8s project webhook create | a8s project webhook create <project-id> [flags] | POST /api/v1/projects/{projectId}/webhook | flags, file |
| a8s project webhook delete | a8s project webhook delete <project-id> [flags] | DELETE /api/v1/projects/{projectId}/webhook | flags |
| a8s project webhook get | a8s project webhook get <project-id> [flags] | GET /api/v1/projects/{projectId}/webhook | flags |
| a8s project webhook rotate | a8s project webhook rotate <project-id> [flags] | POST /api/v1/projects/{projectId}/webhook/rotate | flags |
| a8s scan get | a8s scan get <scan-id> [flags] | GET /api/v1/image-scanner/scans/{scanId} | flags |
| a8s scan images | a8s scan images [flags] | GET /api/v1/image-scanner/images | flags |
| a8s scan list | a8s scan list [flags] | GET /api/v1/image-scanner/scans | flags |
| a8s scan report | a8s scan report <scan-id> [flags] | GET /api/v1/image-scanner/scans/{scanId}/report | flags |
| a8s scan start | a8s scan start [flags] | POST /api/v1/image-scanner/scans | flags |
| a8s sonarqube access | a8s sonarqube access <project-id> [flags] | POST /api/v1/projects/{projectId}/sonarqube/access | flags |
| a8s sonarqube summary | a8s sonarqube summary <project-id> [flags] | GET /api/v1/projects/{projectId}/sonarqube | flags |
| a8s version | a8s version | Print the CLI version | flags |
| a8s workspace bootstrap | a8s workspace bootstrap [flags] | POST /api/v1/workspaces/bootstrap | flags |
| a8s workspace entitlements | a8s workspace entitlements [flags] | GET /api/v1/workspaces/entitlements | flags |
| a8s workspace quota payment-status | a8s workspace quota payment-status <md5> [flags] | GET /api/v1/workspaces/quota-requests/payment-status | flags |
| a8s workspace quota pricing | a8s workspace quota pricing [flags] | GET /api/v1/workspaces/quota-pricing | flags |
| a8s workspace quota purchase | a8s workspace quota purchase [flags] | Purchase a workspace quota plan using Bakong KHQR | flags |
| a8s workspace quota request | a8s workspace quota request [flags] | POST /api/v1/workspaces/quota-requests | flags |
| a8s workspace status | a8s workspace status [flags] | GET /api/v1/workspaces/bootstrap | flags |
