package operation

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/yourname/a8s/internal/clierrors"
)

type Envelope[T any] struct {
	APIVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       T        `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}

func LoadFile[T any](path, expectedKind string, stdin io.Reader) (T, error) {
	var zero T
	if path == "" {
		return zero, nil
	}
	data, err := ReadFile(path, stdin)
	if err != nil {
		return zero, err
	}
	return LoadBytes[T](data, sourceName(path), expectedKind)
}

func ReadFile(path string, stdin io.Reader) ([]byte, error) {
	var reader io.Reader
	if path == "-" {
		reader = stdin
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, clierrors.Validation(fmt.Sprintf("open operation file %q: %v", path, err))
		}
		defer file.Close()
		reader = file
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, clierrors.Validation(fmt.Sprintf("read operation file %q: %v", sourceName(path), err))
	}
	return data, nil
}

func LoadBytes[T any](data []byte, source, expectedKind string) (T, error) {
	var zero T
	var envelope Envelope[T]
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&envelope); err != nil {
		return zero, clierrors.Validation(fmt.Sprintf("decode operation file %q: %v", source, err))
	}
	if envelope.APIVersion != "cli.a8s.io/v1alpha1" {
		return zero, clierrors.Validation(fmt.Sprintf("unsupported apiVersion %q", envelope.APIVersion))
	}
	if !strings.EqualFold(envelope.Kind, expectedKind) {
		return zero, clierrors.Validation(fmt.Sprintf("expected kind %q, got %q", expectedKind, envelope.Kind))
	}
	return envelope.Spec, nil
}

func sourceName(path string) string {
	if path == "-" {
		return "stdin"
	}
	return path
}

func ResolveSecret(envName string, stdinEnabled bool, stdin io.Reader) (string, error) {
	if envName != "" && stdinEnabled {
		return "", clierrors.Validation("secret environment variable and stdin are mutually exclusive")
	}
	if envName != "" {
		value := os.Getenv(envName)
		if value == "" {
			return "", clierrors.Validation(fmt.Sprintf("environment variable %q is empty", envName))
		}
		return value, nil
	}
	if stdinEnabled {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return "", clierrors.Validation(fmt.Sprintf("read secret from stdin: %v", err))
		}
		value := strings.TrimRight(string(data), "\r\n")
		if value == "" {
			return "", clierrors.Validation("secret read from stdin is empty")
		}
		return value, nil
	}
	return "", nil
}
