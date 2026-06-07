package database

import (
	"fmt"
	"strings"

	"github.com/yourname/a8s/internal/api/resources/databases"
	"github.com/yourname/a8s/internal/clierrors"
)

const DeployKind = "DatabaseDeployment"

type SecretReference struct {
	Env string `yaml:"env,omitempty" json:"env,omitempty"`
}

type TLS struct {
	Enabled            *bool  `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	RequireSSL         *bool  `yaml:"requireSsl,omitempty" json:"requireSsl,omitempty"`
	ExistingSecretName string `yaml:"existingSecretName,omitempty" json:"existingSecretName,omitempty"`
	IncludeCA          *bool  `yaml:"includeCa,omitempty" json:"includeCa,omitempty"`
}

type Deploy struct {
	ReleaseName            string           `yaml:"releaseName,omitempty" json:"releaseName,omitempty"`
	ProjectName            string           `yaml:"projectName" json:"projectName"`
	Engine                 string           `yaml:"engine" json:"engine"`
	DeploymentMode         string           `yaml:"deploymentMode,omitempty" json:"deploymentMode,omitempty"`
	DatabaseName           string           `yaml:"databaseName" json:"databaseName"`
	Username               string           `yaml:"username,omitempty" json:"username,omitempty"`
	PasswordFrom           *SecretReference `yaml:"passwordFrom,omitempty" json:"passwordFrom,omitempty"`
	Version                string           `yaml:"version" json:"version"`
	SizeProfile            string           `yaml:"sizeProfile,omitempty" json:"sizeProfile,omitempty"`
	StorageSize            string           `yaml:"storageSize,omitempty" json:"storageSize,omitempty"`
	StorageClassName       string           `yaml:"storageClassName,omitempty" json:"storageClassName,omitempty"`
	Environment            string           `yaml:"environment,omitempty" json:"environment,omitempty"`
	ExistingAuthSecretName string           `yaml:"existingAuthSecretName,omitempty" json:"existingAuthSecretName,omitempty"`
	NetworkPolicyEnabled   *bool            `yaml:"networkPolicyEnabled,omitempty" json:"networkPolicyEnabled,omitempty"`
	TLS                    *TLS             `yaml:"tls,omitempty" json:"tls,omitempty"`
}

type Overrides struct {
	Changed                map[string]bool
	ReleaseName            string
	ProjectName            string
	Engine                 string
	DeploymentMode         string
	DatabaseName           string
	Username               string
	Version                string
	SizeProfile            string
	StorageSize            string
	StorageClassName       string
	Environment            string
	ExistingAuthSecretName string
	NetworkPolicyEnabled   bool
	TLSEnabled             bool
	RequireSSL             bool
	TLSSecret              string
	IncludeCA              bool
}

func (d *Deploy) Apply(overrides Overrides) {
	setString := func(name string, target *string, value string) {
		if overrides.Changed[name] {
			*target = value
		}
	}
	setString("release-name", &d.ReleaseName, overrides.ReleaseName)
	setString("project-name", &d.ProjectName, overrides.ProjectName)
	setString("engine", &d.Engine, overrides.Engine)
	setString("deployment-mode", &d.DeploymentMode, overrides.DeploymentMode)
	setString("database-name", &d.DatabaseName, overrides.DatabaseName)
	setString("username", &d.Username, overrides.Username)
	setString("version", &d.Version, overrides.Version)
	setString("size-profile", &d.SizeProfile, overrides.SizeProfile)
	setString("storage-size", &d.StorageSize, overrides.StorageSize)
	setString("storage-class", &d.StorageClassName, overrides.StorageClassName)
	setString("environment", &d.Environment, overrides.Environment)
	setString("existing-auth-secret", &d.ExistingAuthSecretName, overrides.ExistingAuthSecretName)

	if overrides.Changed["network-policy"] {
		d.NetworkPolicyEnabled = boolPtr(overrides.NetworkPolicyEnabled)
	}
	if overrides.Changed["tls"] || overrides.Changed["require-ssl"] || overrides.Changed["tls-secret"] || overrides.Changed["include-ca"] {
		if d.TLS == nil {
			d.TLS = &TLS{}
		}
	}
	if overrides.Changed["tls"] {
		d.TLS.Enabled = boolPtr(overrides.TLSEnabled)
	}
	if overrides.Changed["require-ssl"] {
		d.TLS.RequireSSL = boolPtr(overrides.RequireSSL)
	}
	if overrides.Changed["tls-secret"] {
		d.TLS.ExistingSecretName = overrides.TLSSecret
	}
	if overrides.Changed["include-ca"] {
		d.TLS.IncludeCA = boolPtr(overrides.IncludeCA)
	}
}

func (d *Deploy) ApplyDefaults() {
	d.Engine = strings.ToLower(strings.TrimSpace(d.Engine))
	if d.DeploymentMode == "" {
		d.DeploymentMode = "single"
	}
	if d.ReleaseName == "" {
		d.ReleaseName = d.ProjectName
	}
}

func (d Deploy) Validate() error {
	required := map[string]string{
		"projectName":    d.ProjectName,
		"engine":         d.Engine,
		"deploymentMode": d.DeploymentMode,
		"databaseName":   d.DatabaseName,
		"version":        d.Version,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			return clierrors.Validation(fmt.Sprintf("%s is required", name))
		}
	}
	switch d.Engine {
	case "postgresql", "mongodb", "mysql", "redis", "cassandra", "oracle", "sqlserver":
	default:
		return clierrors.Validation(fmt.Sprintf("unsupported database engine %q", d.Engine))
	}
	if d.PasswordFrom != nil && d.ExistingAuthSecretName != "" {
		return clierrors.Validation("passwordFrom and existingAuthSecretName are mutually exclusive")
	}
	return nil
}

func (d Deploy) BackendRequest(password string) databases.DeployRequest {
	var tls *databases.TLSRequest
	if d.TLS != nil {
		tls = &databases.TLSRequest{
			Enabled:            d.TLS.Enabled,
			RequireSSL:         d.TLS.RequireSSL,
			ExistingSecretName: d.TLS.ExistingSecretName,
			IncludeCA:          d.TLS.IncludeCA,
		}
	}
	return databases.DeployRequest{
		ReleaseName:            d.ReleaseName,
		ProjectName:            d.ProjectName,
		Engine:                 d.Engine,
		DeploymentMode:         d.DeploymentMode,
		DatabaseName:           d.DatabaseName,
		Username:               d.Username,
		Password:               password,
		Version:                d.Version,
		SizeProfile:            d.SizeProfile,
		StorageSize:            d.StorageSize,
		StorageClassName:       d.StorageClassName,
		Environment:            d.Environment,
		ExistingAuthSecretName: d.ExistingAuthSecretName,
		NetworkPolicyEnabled:   d.NetworkPolicyEnabled,
		TLS:                    tls,
	}
}

func boolPtr(value bool) *bool {
	return &value
}
