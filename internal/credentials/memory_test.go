package credentials

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestMemoryStoreAndContextKey(t *testing.T) {
	store := NewMemoryStore()
	key := Key("development", "")
	if key != "context:development" {
		t.Fatalf("unexpected default key: %q", key)
	}
	if configured := Key("development", " custom-key "); configured != "custom-key" {
		t.Fatalf("unexpected configured key: %q", configured)
	}
	if err := store.Set(key, Record{AccessToken: "token"}); err != nil {
		t.Fatal(err)
	}
	record, err := store.Get(key)
	if err != nil || record.AccessToken != "token" {
		t.Fatalf("unexpected stored record: %#v, %v", record, err)
	}
	if err := store.Delete(key); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(key); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestNativeStoreFallbackWarningOnlyPrintsOnce(t *testing.T) {
	var output bytes.Buffer
	store := &NativeStore{fallbackPath: "credentials.json"}
	store.WarnFallbackTo(&output)

	store.warnFallback()
	store.warnFallback()

	if strings.Count(output.String(), "Warning:") != 1 || !strings.Contains(output.String(), "credentials.json") {
		t.Fatalf("unexpected fallback warning: %q", output.String())
	}
}
