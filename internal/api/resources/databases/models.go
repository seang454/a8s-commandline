package databases

type TLSRequest struct {
	Enabled            *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	RequireSSL         *bool  `json:"requireSsl,omitempty" yaml:"requireSsl,omitempty"`
	ExistingSecretName string `json:"existingSecretName,omitempty" yaml:"existingSecretName,omitempty"`
	IncludeCA          *bool  `json:"includeCa,omitempty" yaml:"includeCa,omitempty"`
}

type DeployRequest struct {
	ReleaseName            string      `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	ProjectName            string      `json:"projectName" yaml:"projectName"`
	Engine                 string      `json:"engine" yaml:"engine"`
	DeploymentMode         string      `json:"deploymentMode" yaml:"deploymentMode"`
	DatabaseName           string      `json:"databaseName" yaml:"databaseName"`
	Username               string      `json:"username,omitempty" yaml:"username,omitempty"`
	Password               string      `json:"password,omitempty" yaml:"password,omitempty"`
	Version                string      `json:"version" yaml:"version"`
	SizeProfile            string      `json:"sizeProfile,omitempty" yaml:"sizeProfile,omitempty"`
	StorageSize            string      `json:"storageSize,omitempty" yaml:"storageSize,omitempty"`
	StorageClassName       string      `json:"storageClassName,omitempty" yaml:"storageClassName,omitempty"`
	Environment            string      `json:"environment,omitempty" yaml:"environment,omitempty"`
	ExistingAuthSecretName string      `json:"existingAuthSecretName,omitempty" yaml:"existingAuthSecretName,omitempty"`
	NetworkPolicyEnabled   *bool       `json:"networkPolicyEnabled,omitempty" yaml:"networkPolicyEnabled,omitempty"`
	TLS                    *TLSRequest `json:"tls,omitempty" yaml:"tls,omitempty"`
}

type Deployment struct {
	ID            string `json:"id,omitempty" yaml:"id,omitempty"`
	DeploymentID  string `json:"deploymentId,omitempty" yaml:"deploymentId,omitempty"`
	ReleaseName   string `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	ProjectName   string `json:"projectName,omitempty" yaml:"projectName,omitempty"`
	Engine        string `json:"engine,omitempty" yaml:"engine,omitempty"`
	Status        string `json:"status,omitempty" yaml:"status,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty" yaml:"statusMessage,omitempty"`
	Namespace     string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}
