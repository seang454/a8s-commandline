package features

import (
	"testing"

	"github.com/yourname/a8s/internal/cli/catalog"
)

func TestBackendFeatureInventory(t *testing.T) {
	if len(Names) != 21 {
		t.Fatalf("expected 21 backend features, got %d", len(Names))
	}
	known := map[string]bool{}
	for _, name := range Names {
		if known[name] {
			t.Fatalf("duplicate feature name %q", name)
		}
		known[name] = true
	}
	for _, route := range catalog.Routes {
		if !known[route.Feature] {
			t.Fatalf("route %s %s has unknown feature %q", route.Method, route.Endpoint, route.Feature)
		}
	}
}
