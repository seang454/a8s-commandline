package operation

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/yourname/a8s/internal/clierrors"
)

// LoadGeneric reads either an operation envelope or a plain YAML/JSON object.
func LoadGeneric(path string, stdin io.Reader) (map[string]any, error) {
	if path == "" {
		return map[string]any{}, nil
	}
	var reader io.Reader
	if path == "-" {
		reader = stdin
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, clierrors.Validation(fmt.Sprintf("open request file %q: %v", path, err))
		}
		defer file.Close()
		reader = file
	}
	var document map[string]any
	if err := yaml.NewDecoder(reader).Decode(&document); err != nil {
		return nil, clierrors.Validation(fmt.Sprintf("decode request file %q: %v", path, err))
	}
	if spec, ok := document["spec"].(map[string]any); ok {
		return spec, nil
	}
	return document, nil
}

// ApplySet applies dotted key=value overrides to a request object.
func ApplySet(target map[string]any, values []string) error {
	for _, assignment := range values {
		key, raw, ok := strings.Cut(assignment, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return clierrors.Validation(fmt.Sprintf("invalid --set %q; expected key=value", assignment))
		}
		parts := strings.Split(key, ".")
		current := target
		for _, part := range parts[:len(parts)-1] {
			next, ok := current[part].(map[string]any)
			if !ok {
				next = map[string]any{}
				current[part] = next
			}
			current = next
		}
		current[parts[len(parts)-1]] = scalar(raw)
	}
	return nil
}

func scalar(value string) any {
	if parsed, err := strconv.ParseBool(value); err == nil {
		return parsed
	}
	if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsed
	}
	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return parsed
	}
	return value
}
