package operation

import (
	"strings"
	"testing"
)

type testSpec struct {
	Name string `yaml:"name"`
}

func TestLoadFileFromStdin(t *testing.T) {
	input := `apiVersion: cli.a8s.io/v1alpha1
kind: Test
spec:
  name: example
`
	spec, err := LoadFile[testSpec]("-", "Test", strings.NewReader(input))
	if err != nil {
		t.Fatalf("LoadFile returned error: %v", err)
	}
	if spec.Name != "example" {
		t.Fatalf("expected name example, got %q", spec.Name)
	}
}

func TestLoadFileRejectsUnknownField(t *testing.T) {
	input := `apiVersion: cli.a8s.io/v1alpha1
kind: Test
spec:
  name: example
  unknown: value
`
	_, err := LoadFile[testSpec]("-", "Test", strings.NewReader(input))
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}
