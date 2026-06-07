package database

import "testing"

func TestDeployExplicitFlagsOverrideManifest(t *testing.T) {
	spec := Deploy{
		ProjectName:    "payments",
		Engine:         "postgresql",
		DeploymentMode: "single",
		DatabaseName:   "payments",
		Version:        "16",
		StorageSize:    "20Gi",
	}
	spec.Apply(Overrides{
		Changed:     map[string]bool{"storage-size": true},
		StorageSize: "50Gi",
		Engine:      "mysql",
	})
	if spec.StorageSize != "50Gi" {
		t.Fatalf("expected explicit storage override, got %q", spec.StorageSize)
	}
	if spec.Engine != "postgresql" {
		t.Fatalf("unchanged engine flag must not override manifest, got %q", spec.Engine)
	}
}

func TestDeployMapsToBackendRequest(t *testing.T) {
	enabled := true
	spec := Deploy{
		ProjectName:          "payments",
		Engine:               "postgresql",
		DatabaseName:         "payments",
		Version:              "16",
		NetworkPolicyEnabled: &enabled,
		TLS:                  &TLS{Enabled: &enabled, RequireSSL: &enabled},
	}
	spec.ApplyDefaults()
	if err := spec.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	request := spec.BackendRequest("secret")
	if request.ReleaseName != "payments" || request.DeploymentMode != "single" {
		t.Fatalf("defaults were not mapped: %#v", request)
	}
	if request.Password != "secret" || request.TLS == nil || request.TLS.RequireSSL == nil || !*request.TLS.RequireSSL {
		t.Fatalf("request mapping incomplete: %#v", request)
	}
}
