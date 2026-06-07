package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type File struct {
	APIVersion     string             `yaml:"apiVersion"`
	Kind           string             `yaml:"kind"`
	CurrentContext string             `yaml:"currentContext"`
	Preferences    Preferences        `yaml:"preferences"`
	Contexts       map[string]Context `yaml:"contexts"`
}

type Preferences struct {
	Output          string `yaml:"output"`
	Color           string `yaml:"color"`
	Timeout         string `yaml:"timeout"`
	RequestTimeout  string `yaml:"requestTimeout"`
	PollingInterval string `yaml:"pollingInterval"`
}

type Context struct {
	Server        string `yaml:"server"`
	Namespace     string `yaml:"namespace"`
	TargetCluster string `yaml:"targetCluster"`
	TLS           TLS    `yaml:"tls"`
	Auth          Auth   `yaml:"auth"`
}

type TLS struct {
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	CAFile             string `yaml:"caFile"`
}

type Auth struct {
	Issuer        string `yaml:"issuer"`
	ClientID      string `yaml:"clientId"`
	CredentialKey string `yaml:"credentialKey"`
}

type Overrides struct {
	ConfigPath     string
	ContextName    string
	Server         string
	Namespace      string
	TargetCluster  string
	Token          string
	Output         string
	Timeout        string
	RequestTimeout string
}

type Resolved struct {
	ConfigPath      string
	ContextName     string
	Server          string
	Namespace       string
	TargetCluster   string
	Token           string
	Output          string
	Timeout         time.Duration
	RequestTimeout  time.Duration
	PollingInterval time.Duration
	TLS             TLS
	Auth            Auth
}

func DefaultPath() (string, error) {
	if explicit := os.Getenv("A8S_CONFIG"); explicit != "" {
		return explicit, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory: %w", err)
	}
	return filepath.Join(dir, "a8s", "config.yaml"), nil
}

func LoadFile(path string) (File, error) {
	file := defaultFile()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return file, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("read config %q: %w", path, err)
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&file); err != nil {
		return File{}, fmt.Errorf("decode config %q: %w", path, err)
	}
	if file.Contexts == nil {
		file.Contexts = map[string]Context{}
	}
	return file, nil
}

func Save(path string, file File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := yaml.Marshal(file)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}
	return nil
}

func Resolve(overrides Overrides) (Resolved, error) {
	path := overrides.ConfigPath
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return Resolved{}, err
		}
	}

	file, err := LoadFile(path)
	if err != nil {
		return Resolved{}, err
	}

	contextName := firstNonEmpty(overrides.ContextName, os.Getenv("A8S_CONTEXT"), file.CurrentContext, "default")
	selected := file.Contexts[contextName]
	resolved := Resolved{
		ConfigPath:    path,
		ContextName:   contextName,
		Server:        firstNonEmpty(overrides.Server, os.Getenv("A8S_SERVER"), selected.Server, "http://localhost:8080"),
		Namespace:     firstNonEmpty(overrides.Namespace, os.Getenv("A8S_NAMESPACE"), selected.Namespace),
		TargetCluster: firstNonEmpty(overrides.TargetCluster, os.Getenv("A8S_TARGET_CLUSTER"), selected.TargetCluster),
		Token:         firstNonEmpty(overrides.Token, os.Getenv("A8S_TOKEN"), os.Getenv("A8S_API_TOKEN")),
		Output:        firstNonEmpty(overrides.Output, os.Getenv("A8S_OUTPUT"), file.Preferences.Output, "table"),
		TLS:           selected.TLS,
		Auth:          selected.Auth,
	}

	if resolved.Timeout, err = parseDuration(firstNonEmpty(overrides.Timeout, os.Getenv("A8S_TIMEOUT"), file.Preferences.Timeout, "30s")); err != nil {
		return Resolved{}, fmt.Errorf("invalid timeout: %w", err)
	}
	if resolved.RequestTimeout, err = parseDuration(firstNonEmpty(overrides.RequestTimeout, os.Getenv("A8S_REQUEST_TIMEOUT"), file.Preferences.RequestTimeout, "20s")); err != nil {
		return Resolved{}, fmt.Errorf("invalid request timeout: %w", err)
	}
	if resolved.PollingInterval, err = parseDuration(firstNonEmpty(os.Getenv("A8S_POLLING_INTERVAL"), file.Preferences.PollingInterval, "3s")); err != nil {
		return Resolved{}, fmt.Errorf("invalid polling interval: %w", err)
	}
	return resolved, nil
}

func defaultFile() File {
	return File{
		APIVersion:     "cli.a8s.io/v1alpha1",
		Kind:           "Config",
		CurrentContext: "default",
		Preferences: Preferences{
			Output:          "table",
			Color:           "auto",
			Timeout:         "30s",
			RequestTimeout:  "20s",
			PollingInterval: "3s",
		},
		Contexts: map[string]Context{},
	}
}

func parseDuration(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
