package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/yourname/a8s/internal/cli/catalog"
)

func TestDatabaseDeployDryRunWithManifestAndOverride(t *testing.T) {
	dir := t.TempDir()
	manifest := filepath.Join(dir, "database.yaml")
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(manifest, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseDeployment
spec:
  projectName: payments
  engine: postgresql
  databaseName: payments
  version: "16"
  storageSize: 20Gi
`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: Config
currentContext: default
contexts:
  default:
    server: http://localhost:8080
`), 0o600); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, &errOut)
	root.SetArgs([]string{
		"--config", configPath,
		"--output", "json",
		"database", "deploy",
		"--file", manifest,
		"--storage-size", "50Gi",
		"--dry-run",
	})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v\nstderr: %s", err, errOut.String())
	}
	if !strings.Contains(out.String(), `"storageSize": "50Gi"`) {
		t.Fatalf("expected override in dry-run output: %s", out.String())
	}
	if !strings.Contains(out.String(), `"deploymentMode": "single"`) {
		t.Fatalf("expected default deployment mode: %s", out.String())
	}
}

func TestManifestInitAndValidateDatabaseDeployment(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	manifestPath := filepath.Join(dir, "database.yaml")

	var initOut bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &initOut, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--output", "json", "manifest", "init", "DatabaseDeployment", "--output-file", manifestPath})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "kind: DatabaseDeployment") {
		t.Fatalf("unexpected manifest template:\n%s", data)
	}

	var validateOut bytes.Buffer
	root = NewRootCommand(strings.NewReader(""), &validateOut, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--output", "json", "manifest", "validate", "--file", manifestPath})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(validateOut.String(), `"valid": true`) || !strings.Contains(validateOut.String(), `"validatedBy": "strict"`) {
		t.Fatalf("unexpected validation output: %s", validateOut.String())
	}
}

func TestManifestValidateRejectsUnknownFieldForStrictKind(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	manifestPath := filepath.Join(dir, "database.yaml")
	if err := os.WriteFile(manifestPath, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: DatabaseDeployment
spec:
  projectName: payments
  engine: postgresql
  databaseName: payments
  version: "16"
  unexpectedField: nope
`), 0o600); err != nil {
		t.Fatal(err)
	}

	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"--config", configPath, "manifest", "validate", "--file", manifestPath})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "field unexpectedField not found") {
		t.Fatalf("expected strict unknown-field validation error, got %v", err)
	}
}

func TestManifestValidateGenericKind(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	manifestPath := filepath.Join(dir, "domain.yaml")
	if err := os.WriteFile(manifestPath, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: ProjectDomain
spec:
  customDomain: api.example.com
`), 0o600); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--output", "json", "manifest", "validate", "--file", manifestPath})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"kind": "ProjectDomain"`) || !strings.Contains(out.String(), `"validatedBy": "generic-envelope"`) {
		t.Fatalf("unexpected validation output: %s", out.String())
	}
}

func TestManifestCommandsAreLocalOnly(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: Config
currentContext: default
contexts:
  default:
    server: https://api.example.com
    tls:
      caFile: C:\missing\a8s-ca.crt
`), 0o600); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--output", "json", "manifest", "kinds"})
	if err := root.Execute(); err != nil {
		t.Fatalf("local manifest command should not configure backend TLS: %v", err)
	}
	if !strings.Contains(out.String(), "DatabaseDeployment") {
		t.Fatalf("unexpected kinds output: %s", out.String())
	}
}

func TestEveryCatalogCommandPathIsRegistered(t *testing.T) {
	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	for _, route := range catalog.Routes {
		command := root
		for _, name := range route.Command {
			command = findChild(command, name)
			if command == nil {
				t.Fatalf("catalog command is not registered: %s", strings.Join(route.Command, " "))
			}
		}
	}
}

func TestCatalogMutationExecutesWithSetAndQuery(t *testing.T) {
	var gotMethod, gotPath, gotToken string
	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		gotMethod = request.Method
		gotPath = request.URL.RequestURI()
		gotToken = request.Header.Get("Authorization")
		defer request.Body.Close()
		if err := json.NewDecoder(request.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode body: %v", err)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{
		"--server", server.URL,
		"--token", "test-token",
		"--output", "json",
		"project", "domain", "set", "project-1",
		"--set", "customDomain=api.example.com",
		"--query", "force=true",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPatch || gotPath != "/api/v1/projects/project-1/domain?force=true" {
		t.Fatalf("unexpected request: %s %s", gotMethod, gotPath)
	}
	if gotToken != "Bearer test-token" {
		t.Fatalf("unexpected authorization header: %q", gotToken)
	}
	if gotBody["customDomain"] != "api.example.com" {
		t.Fatalf("unexpected request body: %#v", gotBody)
	}
}

func TestContextCreateAndUse(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	for _, args := range [][]string{
		{"--config", configPath, "context", "create", "staging", "--server", "https://staging.example.com", "--namespace", "ns-staging"},
		{"--config", configPath, "context", "use", "staging"},
	} {
		root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
		root.SetArgs(args)
		if err := root.Execute(); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "currentContext: staging") || !strings.Contains(string(data), "ns-staging") {
		t.Fatalf("unexpected config:\n%s", data)
	}
}

func TestContextCreateAndUpdateFromFileWithFlagOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	contextFile := filepath.Join(dir, "context.yaml")
	patchFile := filepath.Join(dir, "context-patch.yaml")
	if err := os.WriteFile(contextFile, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: Context
spec:
  server: https://api.a8s.example.com
  namespace: ns-from-file
  targetCluster: primary
  auth:
    issuer: https://keycloak.example.com/realms/a8s
    clientId: a8s-custom
`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(patchFile, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: ContextPatch
spec:
  namespace: ns-prod
  targetCluster: prod-primary
  auth:
    clientId: a8s-prod
`), 0o600); err != nil {
		t.Fatal(err)
	}

	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"--config", configPath, "context", "create", "production", "--file", contextFile, "--namespace", "ns-override"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	root = NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"--config", configPath, "context", "update", "production", "--file", patchFile, "--target-cluster", "prod-override"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	configText := string(data)
	for _, wanted := range []string{
		"https://api.a8s.example.com",
		"namespace: ns-prod",
		"targetCluster: prod-override",
		"https://keycloak.example.com/realms/a8s",
		"clientId: a8s-prod",
		"credentialKey: context:production",
	} {
		if !strings.Contains(configText, wanted) {
			t.Fatalf("expected %q in config:\n%s", wanted, configText)
		}
	}
	if strings.Contains(configText, "ns-from-file") || strings.Contains(configText, "ns-override") {
		t.Fatalf("context patch or flag override was not applied correctly:\n%s", configText)
	}
}

func TestCatalogDryRunRedactsSecrets(t *testing.T) {
	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{
		"--output", "json",
		"database", "rotate-password", "database-1",
		"--set", "password=very-secret",
		"--dry-run",
		"--yes",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "very-secret") || !strings.Contains(out.String(), "[redacted]") {
		t.Fatalf("secret was not redacted: %s", out.String())
	}
}

func TestPayloadFreeCatalogCommandRejectsRequestBody(t *testing.T) {
	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"admin", "quota", "approve", "request-1", "--file", "payload.yaml"})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "does not accept --file") {
		t.Fatalf("expected payload-free body validation error, got %v", err)
	}
}

func TestMultipartUpload(t *testing.T) {
	upload := filepath.Join(t.TempDir(), "avatar.txt")
	if err := os.WriteFile(upload, []byte("avatar-content"), 0o600); err != nil {
		t.Fatal(err)
	}
	var gotContentType, gotFilename, gotContent string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		gotContentType = request.Header.Get("Content-Type")
		if err := request.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("parse multipart: %v", err)
		}
		file, header, err := request.FormFile("file")
		if err != nil {
			t.Errorf("read upload: %v", err)
		} else {
			defer file.Close()
			data, _ := io.ReadAll(file)
			gotFilename, gotContent = header.Filename, string(data)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"--server", server.URL, "profile", "avatar", "upload", "--upload", "file=" + upload})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(gotContentType, "multipart/form-data;") || gotFilename != "avatar.txt" || gotContent != "avatar-content" {
		t.Fatalf("unexpected multipart request: %q %q %q", gotContentType, gotFilename, gotContent)
	}
}

func TestStaticTokenUnauthorizedRequestIsNotRetried(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if request.Header.Get("Authorization") != "Bearer static-token" {
			t.Fatalf("unexpected authorization header: %q", request.Header.Get("Authorization"))
		}
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	root := NewRootCommand(strings.NewReader(""), io.Discard, io.Discard)
	root.SetArgs([]string{"--server", server.URL, "--token", "static-token", "project", "list"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected authentication error")
	}
	if requests != 1 {
		t.Fatalf("static-token request was retried %d time(s)", requests)
	}
}

func TestScanStartWaitsWithTypedImageFlag(t *testing.T) {
	configPath := writeTestConfig(t)
	var gotBody map[string]any
	var pollRequests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/v1/image-scanner/scans":
			if request.Method != http.MethodPost {
				t.Fatalf("unexpected scan start method: %s", request.Method)
			}
			if err := json.NewDecoder(request.Body).Decode(&gotBody); err != nil {
				t.Fatal(err)
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"scanId":"scan-1","status":"RUNNING"}`))
		case "/api/v1/image-scanner/scans/scan-1":
			pollRequests++
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"scanId":"scan-1","status":"COMPLETED"}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--server", server.URL, "--output", "json", "scan", "start", "--image", "nginx:1.27", "--wait"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if gotBody["image"] != "nginx:1.27" {
		t.Fatalf("typed image flag was not applied: %#v", gotBody)
	}
	if pollRequests != 1 || !strings.Contains(out.String(), `"COMPLETED"`) {
		t.Fatalf("scan wait did not poll to completion: polls=%d output=%s", pollRequests, out.String())
	}
}

func TestWorkspaceQuotaPurchaseWaitsForPayment(t *testing.T) {
	configPath := writeTestConfig(t)
	var gotBody map[string]any
	var statusRequests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/v1/workspaces/quota-requests":
			if err := json.NewDecoder(request.Body).Decode(&gotBody); err != nil {
				t.Fatal(err)
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"md5":"payment-md5","paymentStatus":"PENDING"}`))
		case "/api/v1/workspaces/quota-requests/payment-status":
			statusRequests++
			if request.URL.Query().Get("md5") != "payment-md5" {
				t.Fatalf("unexpected md5 query: %s", request.URL.RawQuery)
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"md5":"payment-md5","paymentStatus":"PAID"}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"--config", configPath, "--server", server.URL, "--output", "json", "workspace", "quota", "purchase", "--plan", "premium", "--wait"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if gotBody["planName"] != "premium" || gotBody["paymentProvider"] != "BAKONG" {
		t.Fatalf("purchase flags were not applied: %#v", gotBody)
	}
	if statusRequests != 1 || !strings.Contains(out.String(), `"PAID"`) {
		t.Fatalf("payment wait did not poll to paid: polls=%d output=%s", statusRequests, out.String())
	}
}

func TestClusterDeployWaitsWithTypedFlags(t *testing.T) {
	configPath := writeTestConfig(t)
	var gotBody map[string]any
	var statusRequests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/namespaces/ns-test/cluster-deployments":
			if err := json.NewDecoder(request.Body).Decode(&gotBody); err != nil {
				t.Fatal(err)
			}
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"releaseName":"orders-cluster","status":"INSTALLING"}`))
		case "/api/namespaces/ns-test/cluster-deployments/orders-cluster":
			statusRequests++
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(`{"releaseName":"orders-cluster","status":"DEPLOYED"}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	var out bytes.Buffer
	root := NewRootCommand(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{
		"--config", configPath, "--server", server.URL, "--output", "json",
		"cluster", "deploy", "--release-name", "orders-cluster", "--project-name", "orders",
		"--name", "orders", "--environment", "PRODUCTION", "--wait",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if gotBody["releaseName"] != "orders-cluster" || gotBody["projectName"] != "orders" {
		t.Fatalf("cluster flags were not applied: %#v", gotBody)
	}
	cluster, ok := gotBody["cluster"].(map[string]any)
	if !ok || cluster["name"] != "orders" || cluster["environment"] != "PRODUCTION" {
		t.Fatalf("cluster nested flags were not applied: %#v", gotBody)
	}
	if statusRequests != 1 || !strings.Contains(out.String(), `"DEPLOYED"`) {
		t.Fatalf("cluster wait did not poll to deployed: polls=%d output=%s", statusRequests, out.String())
	}
}

func writeTestConfig(t *testing.T) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(`apiVersion: cli.a8s.io/v1alpha1
kind: Config
currentContext: default
preferences:
  timeout: 2s
  requestTimeout: 1s
  pollingInterval: 1ms
contexts:
  default:
    server: http://localhost:8080
    namespace: ns-test
`), 0o600); err != nil {
		t.Fatal(err)
	}
	return configPath
}

func findChild(parent *cobra.Command, name string) *cobra.Command {
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}
