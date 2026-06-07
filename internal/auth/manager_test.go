package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/config"
	"github.com/yourname/a8s/internal/credentials"
)

func TestResolveTokenPrefersExplicitToken(t *testing.T) {
	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	resolved := testResolved()
	resolved.Token = "explicit-token"

	token, err := manager.ResolveToken(context.Background(), resolved)
	if err != nil {
		t.Fatal(err)
	}
	if token != "explicit-token" {
		t.Fatalf("expected explicit token, got %q", token)
	}
}

func TestResolveTokenUsesValidStoredToken(t *testing.T) {
	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	now := time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC)
	manager.Now = func() time.Time { return now }
	resolved := testResolved()
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	if err := store.Set(key, credentials.Record{
		AccessToken:       "stored-token",
		AccessTokenExpiry: now.Add(10 * time.Minute),
	}); err != nil {
		t.Fatal(err)
	}

	token, err := manager.ResolveToken(context.Background(), resolved)
	if err != nil {
		t.Fatal(err)
	}
	if token != "stored-token" {
		t.Fatalf("expected stored token, got %q", token)
	}
}

func TestResolveTokenRefreshesAndPersistsRotatedToken(t *testing.T) {
	var issuer string
	var refreshGrant, refreshToken string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"issuer":                 issuer,
				"authorization_endpoint": issuer + "/authorize",
				"token_endpoint":         issuer + "/token",
				"jwks_uri":               issuer + "/keys",
			})
		case "/token":
			if err := request.ParseForm(); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			refreshGrant = request.Form.Get("grant_type")
			refreshToken = request.Form.Get("refresh_token")
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"access_token":  "new-access",
				"refresh_token": "new-refresh",
				"token_type":    "Bearer",
				"expires_in":    3600,
			})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	issuer = server.URL

	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	now := time.Now()
	manager.Now = func() time.Time { return now }
	resolved := testResolved()
	resolved.Auth.Issuer = issuer
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	if err := store.Set(key, credentials.Record{
		AccessToken: "old-access", RefreshToken: "old-refresh",
		AccessTokenExpiry: now.Add(-time.Minute), Issuer: issuer, ClientID: resolved.Auth.ClientID,
	}); err != nil {
		t.Fatal(err)
	}

	token, err := manager.ResolveToken(context.Background(), resolved)
	if err != nil {
		t.Fatal(err)
	}
	if token != "new-access" {
		t.Fatalf("expected refreshed token, got %q", token)
	}
	record, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if record.AccessToken != "new-access" || record.RefreshToken != "new-refresh" {
		t.Fatalf("rotated credentials were not persisted: %#v", record)
	}
	if refreshGrant != "refresh_token" || refreshToken != "old-refresh" {
		t.Fatalf("unexpected refresh request: grant=%q token=%q", refreshGrant, refreshToken)
	}
}

func TestResolveTokenExpiredWithoutRefreshRequiresLogin(t *testing.T) {
	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	now := time.Now()
	manager.Now = func() time.Time { return now }
	resolved := testResolved()
	if err := store.Set(credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey), credentials.Record{
		AccessToken: "expired", AccessTokenExpiry: now.Add(-time.Minute),
	}); err != nil {
		t.Fatal(err)
	}

	_, err := manager.ResolveToken(context.Background(), resolved)
	if clierrors.ExitCode(err) != 3 {
		t.Fatalf("expected authentication exit code 3, got %d: %v", clierrors.ExitCode(err), err)
	}
}

func TestRefreshTokenForcesRefreshAndPersistsResult(t *testing.T) {
	var issuer string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"issuer": issuer, "authorization_endpoint": issuer + "/authorize",
				"token_endpoint": issuer + "/token", "jwks_uri": issuer + "/keys",
			})
		case "/token":
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"access_token": "forced-access", "refresh_token": "rotated-refresh",
				"token_type": "Bearer", "expires_in": 3600,
			})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	issuer = server.URL

	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	resolved := testResolved()
	resolved.Auth.Issuer = issuer
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	if err := store.Set(key, credentials.Record{
		AccessToken: "apparently-valid", RefreshToken: "old-refresh",
		AccessTokenExpiry: time.Now().Add(time.Hour), Issuer: issuer, ClientID: resolved.Auth.ClientID,
	}); err != nil {
		t.Fatal(err)
	}

	token, err := manager.RefreshToken(context.Background(), resolved)
	if err != nil {
		t.Fatal(err)
	}
	if token != "forced-access" {
		t.Fatalf("expected forced refresh token, got %q", token)
	}
	record, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if record.AccessToken != "forced-access" || record.RefreshToken != "rotated-refresh" {
		t.Fatalf("forced refresh was not persisted: %#v", record)
	}
}

func TestRefreshTokenClearsCredentialOnInvalidGrant(t *testing.T) {
	var issuer string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(writer).Encode(map[string]any{
				"issuer": issuer, "authorization_endpoint": issuer + "/authorize",
				"token_endpoint": issuer + "/token", "jwks_uri": issuer + "/keys",
			})
		case "/token":
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{"error": "invalid_grant"})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	issuer = server.URL

	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	resolved := testResolved()
	resolved.Auth.Issuer = issuer
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	if err := store.Set(key, credentials.Record{
		AccessToken: "access", RefreshToken: "invalid-refresh",
		AccessTokenExpiry: time.Now().Add(time.Hour), Issuer: issuer, ClientID: resolved.Auth.ClientID,
	}); err != nil {
		t.Fatal(err)
	}

	_, err := manager.RefreshToken(context.Background(), resolved)
	if clierrors.ExitCode(err) != 3 {
		t.Fatalf("expected authentication exit code 3, got %d: %v", clierrors.ExitCode(err), err)
	}
	if _, err := store.Get(key); !errors.Is(err, credentials.ErrNotFound) {
		t.Fatalf("invalid credential was not cleared: %v", err)
	}
}

func TestLogoutOnlyClearsActiveContext(t *testing.T) {
	store := credentials.NewMemoryStore()
	manager := NewManager(store)
	development := testResolved()
	production := testResolved()
	production.ContextName = "production"
	production.Auth.CredentialKey = "context:production"
	for _, resolved := range []config.Resolved{development, production} {
		if err := store.Set(credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey), credentials.Record{AccessToken: resolved.ContextName}); err != nil {
			t.Fatal(err)
		}
	}

	if err := manager.Logout(development); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get("context:development"); err != credentials.ErrNotFound {
		t.Fatalf("development credentials still exist: %v", err)
	}
	if _, err := store.Get("context:production"); err != nil {
		t.Fatalf("production credentials were removed: %v", err)
	}
}

func testResolved() config.Resolved {
	return config.Resolved{
		ContextName: "development",
		Auth: config.Auth{
			Issuer:        "https://keycloak.example.com/realms/a8s",
			ClientID:      "a8s-cli",
			CredentialKey: "context:development",
		},
	}
}
