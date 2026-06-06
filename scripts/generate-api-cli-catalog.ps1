param(
    [string]$BackendRoot = "D:\CSTADPreUniversityTraining\ITP\finalProject\a8s-backend",
    [string]$OutputPath = (Join-Path $PSScriptRoot "..\docs\backend-api-cli-catalog.md")
)

$ErrorActionPreference = "Stop"

$controllerRoot = Join-Path $BackendRoot "src\main\java\com\a8s\autonomous\features"
$controllerFiles = Get-ChildItem $controllerRoot -Recurse -Filter "*Controller.java" | Sort-Object FullName

function Get-QuotedValues([string]$Text) {
    return [regex]::Matches($Text, '"([^"]*)"') | ForEach-Object { $_.Groups[1].Value }
}

function Get-BasePaths([string]$Source) {
    $classIndex = $Source.IndexOf("public class ")
    if ($classIndex -lt 0) {
        $classIndex = $Source.IndexOf("public record ")
    }
    $header = if ($classIndex -ge 0) { $Source.Substring(0, $classIndex) } else { $Source }
    $match = [regex]::Match($header, '@RequestMapping\s*\((?<args>[\s\S]*?)\)')
    if (-not $match.Success) {
        return @("")
    }

    $paths = @(Get-QuotedValues $match.Groups["args"].Value | Where-Object { $_.StartsWith("/") })
    return $(if ($paths.Count -gt 0) { $paths } else { @("") })
}

function Get-MethodMappings([string]$Source) {
    $pattern = '@(?<annotation>GetMapping|PostMapping|PutMapping|PatchMapping|DeleteMapping)\s*(?:\((?<args>[\s\S]*?)\))?'
    $methodNames = @{
        GetMapping = "GET"
        PostMapping = "POST"
        PutMapping = "PUT"
        PatchMapping = "PATCH"
        DeleteMapping = "DELETE"
    }

    return [regex]::Matches($Source, $pattern) | ForEach-Object {
        $values = @(Get-QuotedValues $_.Groups["args"].Value)
        $path = $values | Where-Object { $_.StartsWith("/") } | Select-Object -First 1
        [pscustomobject]@{
            Method = $methodNames[$_.Groups["annotation"].Value]
            Path = if ($path) { $path } else { "" }
        }
    }
}

function Get-SuggestedCommand([string]$Method, [string]$Path) {
    $p = $Path

    if ($p -like "/api/internal/*") { return "(internal service callback; no user CLI command)" }
    if ($p -like "/api/v1/webhooks/*") { return "(provider webhook receiver; no user CLI command)" }
    if ($p -match "/releases/\{releaseId\}/(complete|failed)$|/delete/(complete|failed)$") {
        return "(Jenkins callback; no user CLI command)"
    }

    $rules = @(
        @{ Pattern = '^/api/v1/auth/.*/verify-email$'; GET = 'a8s auth verify-email status'; POST = 'a8s auth verify-email start' }
        @{ Pattern = '^/api/v1/auth/session/onboarding$'; GET = 'a8s auth onboarding status'; POST = 'a8s auth onboarding start' }
        @{ Pattern = '^/api/v1/workspaces/bootstrap$'; GET = 'a8s workspace status'; POST = 'a8s workspace bootstrap' }
        @{ Pattern = '^/api/v1/workspaces/entitlements$'; GET = 'a8s workspace entitlements' }
        @{ Pattern = '^/api/v1/workspaces/quota-pricing$'; GET = 'a8s workspace quota pricing' }
        @{ Pattern = '^/api/v1/workspaces/quota-requests$'; POST = 'a8s workspace quota request' }
        @{ Pattern = '^/api/v1/workspaces/quota-requests/payment-status$'; GET = 'a8s workspace quota payment-status <md5>' }
        @{ Pattern = '^/api/v1/profile/me/avatar$'; GET = 'a8s profile avatar download'; POST = 'a8s profile avatar upload'; DELETE = 'a8s profile avatar delete' }
        @{ Pattern = '^/api/v1/profile/me/account-status$'; GET = 'a8s profile account status' }
        @{ Pattern = '^/api/v1/profile/me/deactivate$'; POST = 'a8s profile account deactivate' }
        @{ Pattern = '^/api/v1/profile/me/reactivate$'; POST = 'a8s profile account reactivate' }
        @{ Pattern = '^/api/v1/profile/me/delete$'; POST = 'a8s profile account delete' }
        @{ Pattern = '^/api/v1/profile/me$'; GET = 'a8s profile get'; PATCH = 'a8s profile update'; DELETE = 'a8s profile account delete' }
        @{ Pattern = '^/api/v1/projects/live$'; GET = 'a8s project live list' }
        @{ Pattern = '^/api/v1/projects/microservices/detect/upload$'; POST = 'a8s microservice detect --upload' }
        @{ Pattern = '^/api/v1/projects/microservices/detect$'; POST = 'a8s microservice detect --repo' }
        @{ Pattern = '^/api/v1/projects/microservices$'; POST = 'a8s microservice deploy' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/canvas$'; PUT = 'a8s microservice apply <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/domains$'; PATCH = 'a8s microservice domains update <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/rollback$'; POST = 'a8s microservice rollback <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/redeploy$'; POST = 'a8s microservice redeploy <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/readiness$'; GET = 'a8s microservice readiness <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/history$'; GET = 'a8s microservice history list <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/history/\{snapshotId\}$'; DELETE = 'a8s microservice history delete <project-id> <snapshot-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/runtime-pods$'; GET = 'a8s microservice pods <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/services/\{serviceId\}/environment/import$'; POST = 'a8s microservice env import <project-id> <service-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/services/\{serviceId\}/environment$'; GET = 'a8s microservice env get <project-id> <service-id>'; PUT = 'a8s microservice env set <project-id> <service-id>'; DELETE = 'a8s microservice env clear <project-id> <service-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}/webhook$'; GET = 'a8s microservice webhook get <project-id>'; POST = 'a8s microservice webhook update <project-id>' }
        @{ Pattern = '^/api/v1/projects/microservices/\{projectId\}$'; GET = 'a8s microservice get <project-id>'; DELETE = 'a8s microservice delete <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/environment/import$'; POST = 'a8s project env import <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/environment$'; GET = 'a8s project env get <project-id>'; PUT = 'a8s project env set <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/auto-deploy$'; PATCH = 'a8s project auto-deploy set <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/branches$'; GET = 'a8s project branches <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/webhook/rotate$'; POST = 'a8s project webhook rotate <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/webhook$'; GET = 'a8s project webhook get <project-id>'; POST = 'a8s project webhook create <project-id>'; DELETE = 'a8s project webhook delete <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/repository/connect$'; POST = 'a8s project repository connect <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/settings$'; PATCH = 'a8s project settings update <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/domain/sync$'; POST = 'a8s project domain sync <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/domain$'; PATCH = 'a8s project domain set <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/sync$'; POST = 'a8s project redeploy <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/releases/\{releaseId\}/rollback$'; POST = 'a8s project release rollback <project-id> <release-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/releases/\{releaseId\}$'; DELETE = 'a8s project release delete <project-id> <release-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/releases$'; GET = 'a8s project releases <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/rollback$'; POST = 'a8s project rollback <project-id>' }
        @{ Pattern = '^/api/v1/projects/me$'; GET = 'a8s project list' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}$'; GET = 'a8s project get <project-id>'; DELETE = 'a8s project delete <project-id>' }
        @{ Pattern = '^/api/v1/projects$'; GET = 'a8s project list'; POST = 'a8s project deploy' }
        @{ Pattern = '^/api/v1/database-deployments/clone-from-backup$'; POST = 'a8s database clone-from-backup' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/credentials$'; GET = 'a8s database credentials <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/metrics$'; GET = 'a8s database metrics <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/settings$'; PATCH = 'a8s database settings update <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/restart$'; POST = 'a8s database restart <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/rotate-password$'; POST = 'a8s database rotate-password <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/verify-password$'; POST = 'a8s database verify-password <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/upgrade-version$'; POST = 'a8s database upgrade <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/console/test$'; POST = 'a8s database console test <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/console/namespaces$'; GET = 'a8s database console namespaces <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/console/objects$'; GET = 'a8s database console objects <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/console/data$'; GET = 'a8s database console data <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/console/query$'; POST = 'a8s database console query <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}$'; GET = 'a8s database get <deployment-id>'; PATCH = 'a8s database update <deployment-id>'; DELETE = 'a8s database delete <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments$'; GET = 'a8s database list'; POST = 'a8s database deploy' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup/runs/\{runId\}/download$'; GET = 'a8s database backup download <deployment-id> <run-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup/runs/\{runId\}/restore/cancel$'; POST = 'a8s database backup restore cancel <deployment-id> <run-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup/runs/\{runId\}/restore$'; POST = 'a8s database backup restore <deployment-id> <run-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup/runs/\{runId\}$'; DELETE = 'a8s database backup delete <deployment-id> <run-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup/run$'; POST = 'a8s database backup run <deployment-id>' }
        @{ Pattern = '^/api/v1/database-deployments/\{deploymentId\}/backup$'; GET = 'a8s database backup settings get <deployment-id>'; PATCH = 'a8s database backup settings set <deployment-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/cluster-deployments/\{releaseName\}/values$'; GET = 'a8s cluster deployment values <release-name>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/cluster-deployments/\{releaseName\}/backup$'; PATCH = 'a8s cluster backup settings set --release <release-name>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/cluster-deployments/\{releaseName\}$'; GET = 'a8s cluster status <release-name>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/cluster-deployments$'; POST = 'a8s cluster deploy' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/clone-from-backup$'; POST = 'a8s cluster clone-from-backup' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/deployments$'; GET = 'a8s cluster history <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/metrics$'; GET = 'a8s cluster metrics <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/settings$'; PATCH = 'a8s cluster settings update <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/upgrade-version$'; POST = 'a8s cluster upgrade <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/certificate$'; GET = 'a8s cluster certificate <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/values/full$'; GET = 'a8s cluster values <cluster-id> --full' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/values$'; GET = 'a8s cluster values <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/backup$'; PATCH = 'a8s cluster backup settings set <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/console/(?<op>deployment|credentials|namespaces|objects|data)$'; GET = 'a8s cluster console <operation> <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}/console/(?<op>test|query)$'; POST = 'a8s cluster console <operation> <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters/\{id\}$'; GET = 'a8s cluster get <cluster-id>'; PATCH = 'a8s cluster update <cluster-id>'; DELETE = 'a8s cluster delete <cluster-id>' }
        @{ Pattern = '^/api/(?:v1/cluster/)?namespaces/\{namespace\}/clusters$'; GET = 'a8s cluster list' }
        @{ Pattern = '^/api/kubernetes/test$'; GET = 'a8s kubernetes test' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/overview$'; GET = 'a8s kubernetes overview' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/pods$'; GET = 'a8s kubernetes pods' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/events$'; GET = 'a8s kubernetes events' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/services$'; GET = 'a8s kubernetes services' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/persistent-volume-claims$'; GET = 'a8s kubernetes pvc' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/database-resources$'; GET = 'a8s kubernetes database-resources' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/pods/\{podName\}/logs/stream$'; GET = 'a8s logs <pod-name> --follow' }
        @{ Pattern = '^/api/kubernetes/namespaces/\{namespace\}/releases/\{releaseName\}/deployment-stream$'; GET = 'a8s cluster watch <release-name>' }
        @{ Pattern = '^/api/backups/settings/\{targetType\}/\{id\}$'; GET = 'a8s backup settings get <type> <id>'; POST = 'a8s backup settings set <type> <id>' }
        @{ Pattern = '^/api/backups/trigger/\{targetType\}/\{id\}$'; POST = 'a8s backup trigger <type> <id>' }
        @{ Pattern = '^/api/backups/restore/\{targetType\}/\{id\}/\{runId\}/cancel$'; POST = 'a8s backup restore cancel <type> <id> <run-id>' }
        @{ Pattern = '^/api/backups/restore/\{targetType\}/\{id\}/\{runId\}$'; POST = 'a8s backup restore <type> <id> <run-id>' }
        @{ Pattern = '^/api/backups/download/\{targetType\}/\{id\}/\{runId\}$'; GET = 'a8s backup download <type> <id> <run-id>' }
        @{ Pattern = '^/api/backups/\{targetType\}/\{id\}/\{runId\}$'; DELETE = 'a8s backup delete <type> <id> <run-id>' }
        @{ Pattern = '^/api/v1/git-integrations/linked-providers$'; GET = 'a8s git providers' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}/connect$'; POST = 'a8s git connect <provider>' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}/sync-keycloak-token$'; POST = 'a8s git sync-token <provider>' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}/brokered-account$'; GET = 'a8s git account <provider>' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}/repos$'; GET = 'a8s git repos <provider>' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}/state$'; GET = 'a8s git state <provider>' }
        @{ Pattern = '^/api/v1/git-integrations/\{provider\}$'; DELETE = 'a8s git disconnect <provider>' }
        @{ Pattern = '^/api/v1/image-scanner/images$'; GET = 'a8s scan images' }
        @{ Pattern = '^/api/v1/image-scanner/scans$'; GET = 'a8s scan list'; POST = 'a8s scan start' }
        @{ Pattern = '^/api/v1/image-scanner/scans/\{scanId\}/report$'; GET = 'a8s scan report <scan-id>' }
        @{ Pattern = '^/api/v1/image-scanner/scans/\{scanId\}$'; GET = 'a8s scan get <scan-id>' }
        @{ Pattern = '^/api/v1/monitoring/overview$'; GET = 'a8s monitoring overview' }
        @{ Pattern = '^/api/v1/projects/live/\{projectId\}/benchmark/run$'; POST = 'a8s benchmark run <project-id>' }
        @{ Pattern = '^/api/v1/projects/live/\{projectId\}/benchmark/runs/\{runId\}$'; GET = 'a8s benchmark get <project-id> <run-id>'; DELETE = 'a8s benchmark delete <project-id> <run-id>' }
        @{ Pattern = '^/api/v1/projects/live/\{projectId\}/benchmark/runs$'; GET = 'a8s benchmark list <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/sonarqube/access$'; POST = 'a8s sonarqube access <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/sonarqube$'; GET = 'a8s sonarqube summary <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/defectdojo/token$'; PUT = 'a8s defectdojo token sync <project-id>' }
        @{ Pattern = '^/api/v1/projects/\{projectId\}/defectdojo$'; GET = 'a8s defectdojo access <project-id>' }
        @{ Pattern = '^/api/v1/alerts/channels/\{channelId\}$'; PUT = 'a8s alert channel update <channel-id>'; DELETE = 'a8s alert channel delete <channel-id>' }
        @{ Pattern = '^/api/v1/alerts/channels$'; GET = 'a8s alert channel list'; POST = 'a8s alert channel create' }
        @{ Pattern = '^/api/v1/alerts/projects/configs$'; GET = 'a8s alert project-config list' }
        @{ Pattern = '^/api/v1/alerts/projects/\{projectId\}/config$'; GET = 'a8s alert project-config get <project-id>'; PUT = 'a8s alert project-config set <project-id>' }
        @{ Pattern = '^/api/v1/alerts/user-config$'; GET = 'a8s alert user-config get'; PUT = 'a8s alert user-config set' }
        @{ Pattern = '^/api/notifications/history/\{userId\}$'; GET = 'a8s notification list' }
        @{ Pattern = '^/api/notifications/\{notificationId\}/read$'; POST = 'a8s notification read <notification-id>' }
        @{ Pattern = '^/api/notifications/preferences/\{userId\}$'; GET = 'a8s notification preferences get'; POST = 'a8s notification preferences set' }
        @{ Pattern = '^/api/v1/jenkins/logs/stream$'; GET = 'a8s project logs --follow' }
        @{ Pattern = '^/api/v1/admin/users$'; GET = 'a8s admin user list'; POST = 'a8s admin user create' }
        @{ Pattern = '^/api/v1/admin/users/\{userId\}/reactivate$'; POST = 'a8s admin user reactivate <user-id>' }
        @{ Pattern = '^/api/v1/admin/users/\{userId\}$'; PATCH = 'a8s admin user update <user-id>'; DELETE = 'a8s admin user deactivate <user-id>' }
        @{ Pattern = '^/api/v1/admin/projects$'; GET = 'a8s admin project list' }
        @{ Pattern = '^/api/v1/admin/projects/\{projectId\}/restore$'; POST = 'a8s admin project restore <project-id>' }
        @{ Pattern = '^/api/v1/admin/projects/\{projectId\}$'; PATCH = 'a8s admin project update <project-id>'; DELETE = 'a8s admin project deactivate <project-id>' }
        @{ Pattern = '^/api/v1/admin/clusters$'; GET = 'a8s admin cluster list' }
        @{ Pattern = '^/api/v1/admin/clusters/\{clusterId\}$'; PATCH = 'a8s admin cluster update <cluster-id>' }
        @{ Pattern = '^/api/v1/admin/clusters/kubernetes$'; GET = 'a8s admin cluster nodes' }
        @{ Pattern = '^/api/v1/admin/clusters/kubernetes/\{alias\}/health$'; GET = 'a8s admin cluster health <alias>' }
        @{ Pattern = '^/api/v1/admin/clusters/kubernetes/\{alias\}/quotas$'; GET = 'a8s admin cluster quota list <alias>' }
        @{ Pattern = '^/api/v1/admin/clusters/kubernetes/\{alias\}/quotas/\{namespace\}$'; PUT = 'a8s admin cluster quota set <alias> <namespace>'; DELETE = 'a8s admin cluster quota delete <alias> <namespace>' }
        @{ Pattern = '^/api/v1/admin/gitops/overview$'; GET = 'a8s admin gitops overview' }
        @{ Pattern = '^/api/v1/admin/gitops/apps$'; POST = 'a8s admin gitops app create' }
        @{ Pattern = '^/api/v1/admin/gitops/apps/\{appId\}/sync$'; POST = 'a8s admin gitops app sync <app-id>' }
        @{ Pattern = '^/api/v1/admin/gitops/apps/\{appId\}/retry$'; POST = 'a8s admin gitops app retry <app-id>' }
        @{ Pattern = '^/api/v1/admin/gitops/apps/\{appId\}/abort$'; POST = 'a8s admin gitops app abort <app-id>' }
        @{ Pattern = '^/api/v1/admin/logs/clusters$'; GET = 'a8s admin logs clusters' }
        @{ Pattern = '^/api/v1/admin/logs/namespaces$'; GET = 'a8s admin logs namespaces' }
        @{ Pattern = '^/api/v1/admin/logs/workloads$'; GET = 'a8s admin logs workloads' }
        @{ Pattern = '^/api/v1/admin/logs/pods$'; GET = 'a8s admin logs pods' }
        @{ Pattern = '^/api/v1/admin/logs/query$'; GET = 'a8s admin logs query' }
        @{ Pattern = '^/api/v1/admin/monitoring/overview$'; GET = 'a8s admin monitoring overview' }
        @{ Pattern = '^/api/v1/admin/quota-requests$'; GET = 'a8s admin quota list' }
        @{ Pattern = '^/api/v1/admin/quota-requests/\{id\}/approve$'; POST = 'a8s admin quota approve <request-id>' }
        @{ Pattern = '^/api/v1/admin/quota-requests/\{id\}/reject$'; POST = 'a8s admin quota reject <request-id>' }
        @{ Pattern = '^/api/v1/admin/registry/health$'; GET = 'a8s admin registry health' }
        @{ Pattern = '^/api/v1/admin/registry/projects$'; GET = 'a8s admin registry project list'; POST = 'a8s admin registry project create' }
        @{ Pattern = '^/api/v1/admin/registry/projects/\{projectName\}/repositories$'; GET = 'a8s admin registry repository list <project-name>'; DELETE = 'a8s admin registry repository delete <project-name>' }
        @{ Pattern = '^/api/v1/admin/registry/projects/\{projectName\}/artifacts$'; GET = 'a8s admin registry artifact list <project-name>'; DELETE = 'a8s admin registry artifact delete <project-name>' }
        @{ Pattern = '^/api/v1/admin/sonarqube/projects$'; GET = 'a8s admin sonarqube project list' }
        @{ Pattern = '^/api/v1/admin/sonarqube/projects/\{projectId\}$'; GET = 'a8s admin sonarqube project get <project-id>' }
        @{ Pattern = '^/api/v1/admin/sonarqube/server-projects$'; GET = 'a8s admin sonarqube server-project list'; POST = 'a8s admin sonarqube server-project create' }
        @{ Pattern = '^/api/v1/admin/sonarqube/server-projects/\{projectKey\}$'; GET = 'a8s admin sonarqube server-project get <project-key>'; PATCH = 'a8s admin sonarqube server-project update <project-key>'; DELETE = 'a8s admin sonarqube server-project delete <project-key>' }
        @{ Pattern = '^/api/admin/documentation/files$'; GET = 'a8s admin docs files' }
        @{ Pattern = '^/api/admin/documentation/content$'; GET = 'a8s admin docs get'; PUT = 'a8s admin docs update'; DELETE = 'a8s admin docs delete' }
        @{ Pattern = '^/api/admin/documentation/publish$'; POST = 'a8s admin docs publish' }
        @{ Pattern = '^/api/admin/documentation/publish/logs$'; GET = 'a8s admin docs publish-logs' }
    )

    foreach ($rule in $rules) {
        if ($p -match $rule.Pattern -and $rule.ContainsKey($Method)) {
            $command = $rule[$Method]
            if ($Matches.ContainsKey("op")) {
                $command = $command.Replace("<operation>", $Matches["op"])
            }
            return $command
        }
    }

    if ($p -like "/api/v1/admin/*" -or $p -like "/api/admin/*") {
        $suffix = ($p -replace '^/api/(v1/)?admin/?', '') -replace '[{}]', '' -replace '/', ' '
        return "a8s admin $suffix [$($Method.ToLowerInvariant())]"
    }

    return "(review command design)"
}

$rows = @()
foreach ($file in $controllerFiles) {
    $source = Get-Content $file.FullName -Raw
    $feature = $file.Directory.Parent.Name
    $basePaths = @(Get-BasePaths $source)

    foreach ($mapping in @(Get-MethodMappings $source)) {
        foreach ($basePath in $basePaths) {
            $fullPath = "$basePath$($mapping.Path)"
            if (-not $fullPath) { $fullPath = "/" }
            $rows += [pscustomobject]@{
                Feature = $feature
                Controller = $file.BaseName
                Method = $mapping.Method
                Path = $fullPath
                Command = Get-SuggestedCommand $mapping.Method $fullPath
            }
        }
    }
}

$rows = $rows | Sort-Object Feature, Controller, Path, Method
$uniqueHandlers = ($rows | Group-Object Controller, Method, Path).Count
$featureFolders = Get-ChildItem $controllerRoot -Directory | Sort-Object Name
$excludedRows = @($rows | Where-Object Command -like "(*no user CLI command*)")
$mappedRows = @($rows | Where-Object {
    $_.Command -notlike "(*no user CLI command*)" -and
    $_.Command -ne "(review command design)"
})
$unmappedRows = @($rows | Where-Object Command -eq "(review command design)")

$lines = [System.Collections.Generic.List[string]]::new()
$lines.Add("# A8S Backend API to CLI Catalog")
$lines.Add("")
$lines.Add("Generated from controller annotations in ``$BackendRoot``.")
$lines.Add("")
$lines.Add("- Feature folders: $($featureFolders.Count)")
$lines.Add("- Controllers: $($controllerFiles.Count)")
$lines.Add("- HTTP route patterns: $($rows.Count)")
$lines.Add("- CLI-eligible route patterns mapped: $($mappedRows.Count)")
$lines.Add("- Automation-only route patterns excluded: $($excludedRows.Count)")
$lines.Add("- Unmapped CLI-eligible route patterns: $($unmappedRows.Count)")
$lines.Add("- WebSocket routes: 4")
$lines.Add("")
$lines.Add("Global CLI flags should include ``--server``, ``--context``, ``--namespace``, ``--target-cluster``, ``--output``, ``--timeout``, and ``--verbose``.")
$lines.Add("")
$lines.Add("## Recommended CLI Command Tree")
$lines.Add("")
$lines.Add("Use resource-first Cobra command groups. Avoid generic top-level commands such as ``a8s create user`` or ``a8s list projects``.")
$lines.Add("")
$lines.Add('```text')
$lines.Add("a8s")
$lines.Add("|-- auth")
$lines.Add("|-- context                 # CLI-local server, token, namespace, and cluster contexts")
$lines.Add("|-- workspace")
$lines.Add('|   `-- quota')
$lines.Add("|-- profile")
$lines.Add("|-- project")
$lines.Add("|-- microservice")
$lines.Add("|-- database")
$lines.Add('|   `-- backup')
$lines.Add("|-- cluster")
$lines.Add("|   |-- backup")
$lines.Add('|   `-- console')
$lines.Add("|-- backup")
$lines.Add("|-- kubernetes")
$lines.Add("|-- logs")
$lines.Add("|-- git")
$lines.Add("|-- scan")
$lines.Add("|-- monitoring")
$lines.Add("|-- benchmark")
$lines.Add("|-- sonarqube")
$lines.Add("|-- defectdojo")
$lines.Add("|-- alert")
$lines.Add("|-- notification")
$lines.Add('`-- admin')
$lines.Add("    |-- user")
$lines.Add("    |-- project")
$lines.Add("    |-- cluster")
$lines.Add("    |-- quota")
$lines.Add("    |-- gitops")
$lines.Add("    |-- registry")
$lines.Add("    |-- sonarqube")
$lines.Add("    |-- monitoring")
$lines.Add("    |-- logs")
$lines.Add("    |-- docs")
$lines.Add('    `-- events')
$lines.Add('```')
$lines.Add("")
$lines.Add("### Implementation order")
$lines.Add("")
$lines.Add("1. Foundation: ``auth``, ``context``, configuration, shared API client, output formats, confirmation prompts, and error handling.")
$lines.Add("2. Core workflow: ``workspace``, ``profile``, ``project``, ``microservice``, ``database``, ``cluster``, and ``backup``.")
$lines.Add("3. Operations: ``kubernetes``, ``logs``, ``git``, ``scan``, ``monitoring``, and ``notification``.")
$lines.Add("4. Quality and security: ``benchmark``, ``sonarqube``, ``defectdojo``, and ``alert``.")
$lines.Add("5. Administration: all commands under ``a8s admin`` with backend ``ROLE_ADMIN`` enforcement.")
$lines.Add("")
$lines.Add("### Command design rules")
$lines.Add("")
$lines.Add("- Use ``get``, ``list``, ``create``, ``update``, and ``delete`` consistently under each resource group.")
$lines.Add("- Require ``--yes`` for destructive commands and support ``--dry-run`` where the API permits it.")
$lines.Add("- Support ``--output table|json|yaml``, plus ``--file`` for complex request bodies.")
$lines.Add("- Keep payment commands under ``a8s workspace quota`` because payment currently exists only for quota and plan purchases.")
$lines.Add("- Never expose internal callbacks, provider webhook receivers, or Jenkins completion callbacks as ordinary CLI commands.")
$lines.Add("")

foreach ($folder in $featureFolders) {
    $feature = $folder.Name
    $featureRows = @($rows | Where-Object Feature -eq $feature)
    $lines.Add("## $feature")
    $lines.Add("")

    if ($featureRows.Count -eq 0) {
        if ($feature -eq "databaseconsole") {
            $lines.Add("No standalone controller. Console APIs are exposed through ``singledb`` and ``dbcluster`` endpoints.")
        } elseif ($feature -eq "payments") {
            $lines.Add("No standalone payment controller. Payments currently support Bakong KHQR purchases for workspace quota and plan upgrades.")
            $lines.Add("")
            $lines.Add("### Payment and quota-purchase endpoints")
            $lines.Add("")
            $lines.Add("| Method | Endpoint | Suggested CLI command | Purpose |")
            $lines.Add("|---|---|---|---|")
            $lines.Add("| ``GET`` | ``/api/v1/workspaces/quota-pricing`` | ``a8s workspace quota pricing`` | Get unit prices and plan prices. |")
            $lines.Add("| ``POST`` | ``/api/v1/workspaces/quota-requests`` | ``a8s workspace quota purchase --plan <plan>`` | Submit a paid quota request and generate a Bakong KHQR payload. |")
            $lines.Add("| ``GET`` | ``/api/v1/workspaces/quota-requests/payment-status?md5=<md5>`` | ``a8s workspace quota payment-status <md5>`` | Check payment status and apply the quota upgrade after payment. |")
            $lines.Add("")
            $lines.Add("The purchase request accepts ``requestedCpu``, ``requestedMemory``, ``requestedStorage``, ``reason``, ``isPaid``, ``planName``, and ``paymentProvider``. Set ``isPaid`` to ``true`` and ``paymentProvider`` to ``BAKONG`` to generate KHQR.")
            $lines.Add("")
            $lines.Add("The purchase response contains ``qrString`` and ``md5``. Use the returned ``md5`` when polling payment status. Status responses currently include ``PENDING``, ``PAID``, and ``NO_PAYMENT_REQUIRED``.")
            $lines.Add("")
            $lines.Add("When payment is confirmed, the backend approves the quota request, applies the workspace quota, activates the subscription for 30 days, and sends a payment receipt notification.")
            $lines.Add("")
            $lines.Add("### Related admin endpoints")
            $lines.Add("")
            $lines.Add("| Method | Endpoint | Suggested CLI command | Purpose |")
            $lines.Add("|---|---|---|---|")
            $lines.Add("| ``GET`` | ``/api/v1/admin/quota-requests`` | ``a8s admin quota list`` | List pending quota and payment-related requests. |")
            $lines.Add("| ``POST`` | ``/api/v1/admin/quota-requests/{id}/approve`` | ``a8s admin quota approve <request-id>`` | Approve a pending request and apply its quota. |")
            $lines.Add("| ``POST`` | ``/api/v1/admin/quota-requests/{id}/reject`` | ``a8s admin quota reject <request-id>`` | Reject a pending request. |")
        } else {
            $lines.Add("No HTTP controller endpoints.")
        }
        $lines.Add("")
        continue
    }

    $lines.Add("| Method | Endpoint | Suggested CLI command | Controller |")
    $lines.Add("|---|---|---|---|")
    foreach ($row in $featureRows) {
        $lines.Add("| ``$($row.Method)`` | ``$($row.Path)`` | ``$($row.Command)`` | ``$($row.Controller)`` |")
    }
    $lines.Add("")
}

$lines.Add("## WebSockets")
$lines.Add("")
$lines.Add("| Endpoint | Suggested CLI use |")
$lines.Add("|---|---|")
$lines.Add("| ``/ws/jenkins/logs`` | ``a8s project logs --follow`` |")
$lines.Add("| ``/ws/notifications`` | ``a8s notification watch`` |")
$lines.Add("| ``/ws/monitoring/overview`` | ``a8s monitoring watch`` |")
$lines.Add("| ``/ws/admin/events`` | ``a8s admin events watch`` |")
$lines.Add("")
$lines.Add("## CLI Exclusions")
$lines.Add("")
$lines.Add("Do not expose provider webhook receivers, Jenkins completion/failure callbacks, or ``/api/internal/**`` endpoints as ordinary user commands. They are automation-to-backend contracts.")

$resolvedOutput = [System.IO.Path]::GetFullPath($OutputPath)
$outputDirectory = Split-Path $resolvedOutput -Parent
New-Item -ItemType Directory -Path $outputDirectory -Force | Out-Null
[System.IO.File]::WriteAllLines($resolvedOutput, $lines)

Write-Output "Generated $resolvedOutput"
Write-Output "Controllers: $($controllerFiles.Count)"
Write-Output "HTTP route patterns: $($rows.Count)"
