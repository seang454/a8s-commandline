package manifest

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/operation"
	databaseoperation "github.com/yourname/a8s/internal/operation/kinds/database"
)

const APIVersion = "cli.a8s.io/v1alpha1"

type Definition struct {
	Kind        string    `json:"kind" yaml:"kind"`
	Description string    `json:"description" yaml:"description"`
	Required    []string  `json:"required,omitempty" yaml:"required,omitempty"`
	Example     string    `json:"example" yaml:"example"`
	Strict      Validator `json:"-" yaml:"-"`
}

type Validator func(data []byte, source string) error

type ValidationResult struct {
	Valid       bool     `json:"valid" yaml:"valid"`
	APIVersion  string   `json:"apiVersion" yaml:"apiVersion"`
	Kind        string   `json:"kind" yaml:"kind"`
	ValidatedBy string   `json:"validatedBy" yaml:"validatedBy"`
	Warnings    []string `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

type rawEnvelope struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   map[string]any `yaml:"metadata,omitempty"`
	Spec       yaml.Node      `yaml:"spec"`
}

var definitions = map[string]Definition{}

func init() {
	register(entry("Context", "Create a local CLI context.", []string{"server"}, `
server: https://api.a8s.example.com
namespace: ns-team
targetCluster: primary
auth:
  issuer: https://keycloak.autonomous-istad.com/realms/a8s
  clientId: a8s-cli
`, strictContext))
	register(entry("ContextPatch", "Update a local CLI context.", nil, `
namespace: ns-prod
targetCluster: prod-primary
auth:
  clientId: a8s-cli
`, strictContextPatch))

	register(entry("WorkspaceQuotaRequest", "Request additional workspace quota.", []string{"requestedCpu", "requestedMemory", "requestedStorage", "reason", "planName"}, `
requestedCpu: "4"
requestedMemory: 8Gi
requestedStorage: 100Gi
reason: Production workload
isPaid: false
planName: Free
`, nil))
	register(entry("WorkspaceQuotaPurchase", "Purchase workspace quota through a payment provider.", []string{"requestedCpu", "requestedMemory", "requestedStorage", "planName", "paymentProvider"}, `
requestedCpu: "8"
requestedMemory: 16Gi
requestedStorage: 200Gi
reason: Production upgrade
isPaid: true
planName: Premium
paymentProvider: BAKONG
`, nil))
	register(entry("ProfileUpdate", "Update the authenticated user profile.", nil, `
personal:
  firstName: Dara
  lastName: Sok
  displayName: Dara
  email: dara@example.com
  jobTitle: Engineer
  department: Platform
locale:
  city: Phnom Penh
  country: Cambodia
  timezone: Asia/Phnom_Penh
  language: en
  dateFormat: yyyy-MM-dd
`, nil))
	register(entry("GitProviderConnection", "Connect a Git provider account.", []string{"accessTokenFrom.env"}, `
accessTokenFrom:
  env: GITHUB_TOKEN
accessLevel: repository
grantedScopes: repo,read:user
`, nil))
	register(entry("DefectDojoToken", "Sync a DefectDojo API token.", []string{"apiTokenFrom.env"}, `
apiTokenFrom:
  env: DEFECTDOJO_API_TOKEN
`, nil))

	register(entry("ProjectDeployment", "Deploy a monolithic project.", []string{"projectName", "sourceType", "architectureType"}, `
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
`, nil))
	register(entry("ProjectDomain", "Set a monolithic project domain.", []string{"customDomain"}, `
customDomain: api.example.com
`, nil))
	register(entry("ProjectRepositoryConnection", "Connect a project repository.", []string{"repoProvider", "repoUrl", "repoFullName", "branch"}, `
repoProvider: github
repoUrl: https://github.com/acme/shop-api.git
repoFullName: acme/shop-api
branch: main
autoDeployEnabled: true
autoDeployTrigger: push
`, nil))
	register(entry("ProjectSettings", "Update project operator settings.", nil, `
alias: shop
operatorNote: Critical service
failureAlerts: true
maintenanceMode: false
protectFromDelete: true
`, nil))
	register(entry("ProjectEnvironment", "Set project environment variables.", []string{"envVars"}, `
envVars:
  - name: SPRING_PROFILES_ACTIVE
    value: production
    secret: false
`, nil))
	register(entry("ProjectAutoDeploy", "Configure project auto deployment.", []string{"enabled"}, `
enabled: true
branch: main
autoDeployTrigger: push
releaseTagPattern: "v*"
`, nil))
	register(entry("ProjectWebhook", "Create a project webhook.", []string{"name", "branch"}, `
name: shop-webhook
branch: main
autoDeployEnabled: true
autoDeployTrigger: push
createOnProvider: true
`, nil))
	register(entry("ProjectRollback", "Roll back a project.", []string{"releaseId"}, `
releaseId: 11111111-1111-1111-1111-111111111111
`, nil))
	register(entry("ProjectReleaseRollback", "Roll back a specific project release.", nil, `
buildNumber: 42
framework: spring
statusMessage: Operator rollback
`, nil))

	register(entry("MicroserviceDeployment", "Deploy a microservice project.", []string{"projectName", "services"}, `
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
`, nil))
	register(entry("MicroserviceCanvas", "Apply a microservice canvas.", []string{"services"}, `
branch: main
services:
  - name: api
    repoUrl: https://github.com/acme/shop.git
    repoFullName: acme/shop
    path: services/api
    appPort: 8080
`, nil))
	register(entry("MicroserviceDomains", "Update microservice domains.", []string{"services"}, `
services:
  - serviceId: 11111111-1111-1111-1111-111111111111
    customDomain: api.example.com
    platformSubdomain: api
`, nil))
	register(entry("MicroserviceRollback", "Roll back a microservice project.", []string{"snapshotId"}, `
snapshotId: snapshot-123
`, nil))
	register(entry("MicroserviceEnvironment", "Set microservice environment configuration.", nil, `
envVars:
  - name: SPRING_PROFILES_ACTIVE
    value: production
    secret: false
runtimeConfigFile:
  fileName: application.yaml
  content: |
    server:
      port: 8080
`, nil))
	register(entry("MicroserviceWebhook", "Update a microservice webhook.", []string{"name", "branch"}, `
name: shop-webhook
branch: main
autoDeployEnabled: true
autoDeployTrigger: push
releaseTagPattern: "v*"
releaseTriggerMode: tag
`, nil))
	register(entry("MicroserviceDetection", "Detect microservices from a repository.", []string{"repoUrl"}, `
repoUrl: https://github.com/acme/shop.git
branch: main
githubTokenFrom:
  env: GITHUB_TOKEN
`, nil))
	register(entry("MicroserviceUploadDetection", "Detect microservices from uploaded source.", []string{"sourceName"}, `
sourceName: shop-source
`, nil))

	register(entry(databaseoperation.DeployKind, "Deploy a single database.", []string{"projectName", "engine", "databaseName", "version"}, `
projectName: payments
engine: postgresql
deploymentMode: single
databaseName: payments
username: app
version: "16"
sizeProfile: small
storageSize: 20Gi
tls:
  enabled: true
  requireSsl: true
`, strictDatabaseDeployment))
	register(entry("DatabaseDeploymentPatch", "Update a single database deployment.", nil, `
sizeProfile: medium
networkPolicyEnabled: true
tls:
  enabled: true
  requireSsl: true
`, nil))
	register(entry("DatabaseSettings", "Update database operator settings.", nil, `
alias: payments
operatorNote: Primary database
failureAlerts: true
maintenanceMode: false
protectFromDelete: true
`, nil))
	register(entry("DatabaseUpgrade", "Upgrade a database engine version.", []string{"version"}, `
version: "17"
`, nil))
	register(entry("DatabaseClone", "Clone a database from backup.", []string{"sourceDeploymentId", "backupRunId", "projectName"}, `
sourceDeploymentId: 11111111-1111-1111-1111-111111111111
backupRunId: 22222222-2222-2222-2222-222222222222
projectName: payments-clone
databaseName: payments
version: "16"
storageSize: 20Gi
`, nil))
	register(entry("DatabasePasswordRotation", "Rotate a database password.", []string{"passwordFrom.env"}, `
passwordFrom:
  env: A8S_NEW_DATABASE_PASSWORD
`, nil))
	register(entry("DatabasePasswordVerification", "Verify a database password.", []string{"passwordFrom.env"}, `
passwordFrom:
  env: A8S_DATABASE_PASSWORD
`, nil))
	register(entry("DatabaseQuery", "Run a database console query.", []string{"query"}, `
query: SELECT now()
`, nil))
	register(entry("DatabaseBackupSettings", "Configure single database backup settings.", []string{"enabled"}, `
enabled: true
destinationPath: s3://backups/payments
credentialSecret: backup-credentials
retentionPolicy: 30d
schedule: "0 0 * * *"
`, nil))

	register(entry("ClusterDeployment", "Deploy a managed database cluster.", []string{"projectName", "cluster.name", "database.engine"}, `
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
  version: "16"
  storageSize: 100Gi
`, nil))
	register(entry("ClusterDeploymentPatch", "Update a managed database cluster.", nil, `
database:
  instances: 5
  storageSize: 200Gi
  monitoringEnabled: true
`, nil))
	register(entry("ClusterSettings", "Update cluster operator settings.", nil, `
alias: orders
operatorNote: Primary production cluster
failureAlerts: true
maintenanceMode: false
protectFromDelete: true
`, nil))
	register(entry("ClusterUpgrade", "Upgrade a cluster database version.", []string{"version"}, `
version: "17"
`, nil))
	register(entry("ClusterPasswordRotation", "Rotate a cluster database password.", []string{"passwordFrom.env"}, `
passwordFrom:
  env: A8S_NEW_CLUSTER_PASSWORD
`, nil))
	register(entry("ClusterClone", "Clone a cluster from backup.", []string{"sourceClusterId", "backupRunId", "projectName"}, `
sourceClusterId: 11111111-1111-1111-1111-111111111111
backupRunId: 22222222-2222-2222-2222-222222222222
projectName: orders-clone
releaseName: orders-clone
version: "16"
instances: 3
storageSize: 100Gi
`, nil))
	register(entry("ClusterBackupSettings", "Configure cluster backup settings.", []string{"enabled"}, `
enabled: true
destinationPath: s3://backups/orders
credentialSecret: backup-credentials
retentionPolicy: 30d
schedule: "0 0 * * *"
`, nil))
	register(entry("ClusterQuery", "Run a cluster database console query.", []string{"query"}, `
query: SELECT now()
`, nil))

	register(entry("BackupSettings", "Configure backup settings for a resource.", []string{"enabled"}, `
enabled: true
destinationPath: s3://backups/resource
credentialSecret: backup-credentials
retentionPolicy: 30d
schedule: "0 0 * * *"
`, nil))
	register(entry("ImageScan", "Start an image security scan.", []string{"sourceKind"}, `
sourceKind: image
imageRef: nginx:1.27
forceRescan: false
`, nil))
	register(entry("BenchmarkRun", "Run a project benchmark.", []string{"concurrency", "totalRequests", "targetPath", "method"}, `
concurrency: 20
totalRequests: 1000
targetPath: /api/health
method: GET
headers:
  Accept: application/json
`, nil))
	register(entry("AlertChannel", "Create or update an alert channel.", []string{"name", "type"}, `
name: operations
type: telegram
credentialFrom:
  env: TELEGRAM_BOT_TOKEN
secondaryCredentialFrom:
  env: TELEGRAM_CHAT_ID
targetProject: shop-api
`, nil))
	register(entry("ProjectAlertConfig", "Configure project alerts.", nil, `
telegramEnabled: true
emailEnabled: true
backupAlertsEnabled: true
securityAlertsEnabled: true
telegramChannelName: operations
emailAddress: ops@example.com
`, nil))
	register(entry("UserAlertConfig", "Configure user alerts.", nil, `
quotaAlertsEnabled: true
globalSecurityAlertsEnabled: true
`, nil))
	register(entry("NotificationPreferences", "Configure notification preferences.", nil, `
buildFailures: true
rolloutReady: true
vulnerabilityFindings: true
weeklyDigest: false
`, nil))

	register(entry("AdminUserCreate", "Create an admin-managed user.", []string{"username", "email"}, `
username: dara
email: dara@example.com
firstName: Dara
lastName: Sok
passwordFrom:
  env: A8S_INITIAL_PASSWORD
`, nil))
	register(entry("AdminUserUpdate", "Update an admin-managed user.", nil, `
username: dara
email: dara@example.com
firstName: Dara
lastName: Sok
`, nil))
	register(entry("AdminProjectUpdate", "Update a project as an admin.", nil, `
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
`, nil))
	register(entry("AdminClusterUpdate", "Update a cluster as an admin.", nil, `
alias: primary
operatorNote: Production cluster
failureAlerts: true
maintenanceMode: false
protectFromDelete: true
`, nil))
	register(entry("AdminClusterQuota", "Set a cluster namespace quota as an admin.", []string{"cpuLimit", "memoryLimit", "pvcLimit"}, `
cpuLimit: 8
memoryLimit: 17179869184
pvcLimit: 20
`, nil))
	register(entry("AdminGitOpsApplication", "Create an admin GitOps application.", []string{"name", "repoUrl"}, `
name: shop-api
repoUrl: https://github.com/acme/shop-gitops.git
authType: token
tokenFrom:
  env: GITOPS_TOKEN
`, nil))
	register(entry("AdminRegistryProject", "Create a registry project as an admin.", []string{"name"}, `
name: shop
publicProject: false
`, nil))
	register(entry("AdminSonarQubeProject", "Create a SonarQube project as an admin.", []string{"key", "name"}, `
key: shop-api
name: Shop API
mainBranch: main
visibility: private
`, nil))
	register(entry("AdminSonarQubeProjectPatch", "Update a SonarQube project as an admin.", nil, `
key: shop-api-v2
visibility: private
`, nil))
	register(entry("AdminDocumentationUpdate", "Update documentation content as an admin.", []string{"path", "contentFile", "message"}, `
path: guides/deploy.md
contentFile: deploy.md
sha: abc123
message: Update deploy guide
`, nil))
}

func Definitions() []Definition {
	values := make([]Definition, 0, len(definitions))
	for _, definition := range definitions {
		values = append(values, definition)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Kind < values[j].Kind
	})
	return values
}

func Kinds() []string {
	defs := Definitions()
	kinds := make([]string, 0, len(defs))
	for _, definition := range defs {
		kinds = append(kinds, definition.Kind)
	}
	return kinds
}

func Get(kind string) (Definition, bool) {
	definition, ok := definitions[normalize(kind)]
	return definition, ok
}

func Template(kind string) (string, error) {
	definition, ok := Get(kind)
	if !ok {
		return "", unknownKind(kind)
	}
	return definition.Example, nil
}

func ValidateFile(path string, stdin io.Reader) (ValidationResult, error) {
	data, err := operation.ReadFile(path, stdin)
	if err != nil {
		return ValidationResult{}, err
	}
	source := path
	if path == "-" {
		source = "stdin"
	}
	return ValidateBytes(data, source)
}

func ValidateBytes(data []byte, source string) (ValidationResult, error) {
	raw, err := decodeEnvelope(data, source)
	if err != nil {
		return ValidationResult{}, err
	}
	if raw.APIVersion != APIVersion {
		return ValidationResult{}, clierrors.Validation(fmt.Sprintf("unsupported apiVersion %q", raw.APIVersion))
	}
	definition, ok := Get(raw.Kind)
	if !ok {
		return ValidationResult{}, unknownKind(raw.Kind)
	}
	if raw.Spec.Kind == 0 {
		return ValidationResult{}, clierrors.Validation("operation file requires spec")
	}
	if raw.Spec.Kind != yaml.MappingNode {
		return ValidationResult{}, clierrors.Validation("operation spec must be a YAML object")
	}
	if definition.Strict != nil {
		if err := definition.Strict(data, source); err != nil {
			return ValidationResult{}, err
		}
		return ValidationResult{Valid: true, APIVersion: raw.APIVersion, Kind: definition.Kind, ValidatedBy: "strict"}, nil
	}

	spec := map[string]any{}
	if err := raw.Spec.Decode(&spec); err != nil {
		return ValidationResult{}, clierrors.Validation(fmt.Sprintf("decode operation spec: %v", err))
	}
	if missing := missingRequired(spec, definition.Required); len(missing) > 0 {
		return ValidationResult{}, clierrors.Validation(fmt.Sprintf("missing required field(s): %s", strings.Join(missing, ", ")))
	}
	return ValidationResult{
		Valid:       true,
		APIVersion:  raw.APIVersion,
		Kind:        definition.Kind,
		ValidatedBy: "generic-envelope",
		Warnings:    []string{"strict field validation is not implemented for this kind yet"},
	}, nil
}

func register(definition Definition) {
	definitions[normalize(definition.Kind)] = definition
}

func entry(kind, description string, required []string, spec string, strict Validator) Definition {
	body := strings.Trim(spec, "\r\n")
	if strings.TrimSpace(body) == "" {
		body = "{}"
	}
	return Definition{
		Kind:        kind,
		Description: description,
		Required:    required,
		Example:     fmt.Sprintf("# Generated by: a8s manifest init %s\napiVersion: %s\nkind: %s\nspec:\n%s\n", kind, APIVersion, kind, indent(body, "  ")),
		Strict:      strict,
	}
}

func indent(value, prefix string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines[i] = line
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func decodeEnvelope(data []byte, source string) (rawEnvelope, error) {
	var raw rawEnvelope
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&raw); err != nil {
		return rawEnvelope{}, clierrors.Validation(fmt.Sprintf("decode operation file %q: %v", source, err))
	}
	return raw, nil
}

func missingRequired(spec map[string]any, required []string) []string {
	var missing []string
	for _, field := range required {
		if !hasDotted(spec, field) {
			missing = append(missing, field)
		}
	}
	return missing
}

func hasDotted(value map[string]any, path string) bool {
	parts := strings.Split(path, ".")
	var current any = value
	for _, part := range parts {
		object, ok := current.(map[string]any)
		if !ok {
			return false
		}
		current, ok = object[part]
		if !ok || current == nil {
			return false
		}
	}
	switch typed := current.(type) {
	case string:
		return strings.TrimSpace(typed) != ""
	case []any:
		return len(typed) > 0
	case map[string]any:
		return len(typed) > 0
	default:
		return true
	}
}

func strictDatabaseDeployment(data []byte, source string) error {
	spec, err := operation.LoadBytes[databaseoperation.Deploy](data, source, databaseoperation.DeployKind)
	if err != nil {
		return err
	}
	spec.ApplyDefaults()
	return spec.Validate()
}

type contextSpec struct {
	Server        string `yaml:"server,omitempty" json:"server,omitempty"`
	Namespace     string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	TargetCluster string `yaml:"targetCluster,omitempty" json:"targetCluster,omitempty"`
	TLS           struct {
		InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
		CAFile             string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
	} `yaml:"tls,omitempty" json:"tls,omitempty"`
	Auth struct {
		Issuer        string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
		ClientID      string `yaml:"clientId,omitempty" json:"clientId,omitempty"`
		CredentialKey string `yaml:"credentialKey,omitempty" json:"credentialKey,omitempty"`
	} `yaml:"auth,omitempty" json:"auth,omitempty"`
}

type contextPatchSpec struct {
	Server        *string `yaml:"server,omitempty" json:"server,omitempty"`
	Namespace     *string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	TargetCluster *string `yaml:"targetCluster,omitempty" json:"targetCluster,omitempty"`
	TLS           *struct {
		InsecureSkipVerify *bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
		CAFile             *string `yaml:"caFile,omitempty" json:"caFile,omitempty"`
	} `yaml:"tls,omitempty" json:"tls,omitempty"`
	Auth *struct {
		Issuer        *string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
		ClientID      *string `yaml:"clientId,omitempty" json:"clientId,omitempty"`
		CredentialKey *string `yaml:"credentialKey,omitempty" json:"credentialKey,omitempty"`
	} `yaml:"auth,omitempty" json:"auth,omitempty"`
}

func strictContext(data []byte, source string) error {
	spec, err := operation.LoadBytes[contextSpec](data, source, "Context")
	if err != nil {
		return err
	}
	if strings.TrimSpace(spec.Server) == "" {
		return clierrors.Validation("server is required")
	}
	return nil
}

func strictContextPatch(data []byte, source string) error {
	_, err := operation.LoadBytes[contextPatchSpec](data, source, "ContextPatch")
	return err
}

func unknownKind(kind string) error {
	return clierrors.Validation(fmt.Sprintf("unknown operation kind %q; run 'a8s manifest kinds' to list supported kinds", kind))
}

func normalize(kind string) string {
	return strings.ToLower(strings.TrimSpace(kind))
}
