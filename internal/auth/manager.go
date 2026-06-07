package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"

	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/config"
	"github.com/yourname/a8s/internal/credentials"
)

const refreshWindow = 2 * time.Minute

type Manager struct {
	Store       credentials.Store
	OpenBrowser func(string) error
	Now         func() time.Time
}

type LoginOptions struct {
	NoBrowser    bool
	CallbackPort int
	Out          io.Writer
}

type LogoutOptions struct {
	NoBrowser    bool
	CallbackPort int
	Out          io.Writer
}

type LogoutResult struct {
	RemoteAttempted bool
}

func NewManager(store credentials.Store) *Manager {
	return &Manager{Store: store, OpenBrowser: browser.OpenURL, Now: time.Now}
}

func (m *Manager) ResolveToken(ctx context.Context, resolved config.Resolved) (string, error) {
	if resolved.Token != "" {
		return resolved.Token, nil
	}
	if resolved.Auth.CredentialKey == "" {
		return "", nil
	}
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	record, err := m.Store.Get(key)
	if err != nil {
		if errors.Is(err, credentials.ErrNotFound) {
			return "", nil
		}
		return "", err
	}
	if record.AccessToken != "" && record.AccessTokenExpiry.After(m.Now().Add(refreshWindow)) {
		return record.AccessToken, nil
	}
	if record.RefreshToken == "" {
		return "", clierrors.New("authentication_required", "stored access token expired; run a8s auth login", 3)
	}
	refreshed, err := m.refresh(ctx, record)
	if err != nil {
		return "", clierrors.New("authentication_required", "token refresh failed; run a8s auth login", 3)
	}
	if err := m.Store.Set(key, refreshed); err != nil {
		return "", err
	}
	return refreshed.AccessToken, nil
}

func (m *Manager) RefreshToken(ctx context.Context, resolved config.Resolved) (string, error) {
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	record, err := m.Store.Get(key)
	if err != nil {
		if errors.Is(err, credentials.ErrNotFound) {
			return "", clierrors.New("authentication_required", "stored credentials not found; run a8s auth login", 3)
		}
		return "", err
	}
	if record.RefreshToken == "" {
		return "", clierrors.New("authentication_required", "stored credentials cannot be refreshed; run a8s auth login", 3)
	}
	record.AccessTokenExpiry = m.Now().Add(-time.Minute)
	refreshed, err := m.refresh(ctx, record)
	if err != nil {
		if invalidRefreshGrant(err) {
			_ = m.Store.Delete(key)
		}
		return "", clierrors.New("authentication_required", "token refresh failed; run a8s auth login", 3)
	}
	if err := m.Store.Set(key, refreshed); err != nil {
		return "", err
	}
	return refreshed.AccessToken, nil
}

func (m *Manager) Login(ctx context.Context, resolved config.Resolved, options LoginOptions) (credentials.Record, error) {
	if resolved.Auth.Issuer == "" || resolved.Auth.ClientID == "" {
		return credentials.Record{}, clierrors.Validation("active context requires auth.issuer and auth.clientId")
	}
	if options.CallbackPort < 0 || options.CallbackPort > 65535 {
		return credentials.Record{}, clierrors.Validation("callback port must be between 0 and 65535")
	}
	provider, err := oidc.NewProvider(ctx, resolved.Auth.Issuer)
	if err != nil {
		return credentials.Record{}, fmt.Errorf("discover OIDC provider: %w", err)
	}
	listenAddress := "127.0.0.1:0"
	if options.CallbackPort > 0 {
		listenAddress = fmt.Sprintf("127.0.0.1:%d", options.CallbackPort)
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return credentials.Record{}, fmt.Errorf("start login callback: %w", err)
	}
	defer listener.Close()
	redirectURL := "http://" + listener.Addr().String() + "/callback"
	oauthConfig := oauth2.Config{
		ClientID:    resolved.Auth.ClientID,
		Endpoint:    provider.Endpoint(),
		RedirectURL: redirectURL,
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
	}
	state, err := randomValue(32)
	if err != nil {
		return credentials.Record{}, err
	}
	nonce, err := randomValue(32)
	if err != nil {
		return credentials.Record{}, err
	}
	verifier, err := randomValue(64)
	if err != nil {
		return credentials.Record{}, err
	}
	challengeBytes := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(challengeBytes[:])
	authURL := oauthConfig.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	fmt.Fprintln(options.Out, "CLI callback redirect URI:")
	fmt.Fprintln(options.Out, redirectURL)
	fmt.Fprintln(options.Out, "Keycloak must allow this redirect URI, or a wildcard such as http://127.0.0.1:*")
	fmt.Fprintln(options.Out, "Open this URL to authenticate:")
	fmt.Fprintln(options.Out, authURL)
	if !options.NoBrowser {
		_ = m.OpenBrowser(authURL)
	}

	type result struct {
		code string
		err  error
	}
	resultChannel := make(chan result, 1)
	server := &http.Server{ReadHeaderTimeout: 5 * time.Second}
	server.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/callback" {
			http.NotFound(writer, request)
			return
		}
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "Invalid login state.", http.StatusBadRequest)
			resultChannel <- result{err: clierrors.New("authentication_required", "OIDC state validation failed", 3)}
			return
		}
		if remoteError := request.URL.Query().Get("error"); remoteError != "" {
			http.Error(writer, "Authentication failed.", http.StatusBadRequest)
			resultChannel <- result{err: clierrors.New("authentication_required", remoteError, 3)}
			return
		}
		fmt.Fprintln(writer, "A8S CLI authentication completed. You may close this window.")
		resultChannel <- result{code: request.URL.Query().Get("code")}
	})
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	var login result
	select {
	case login = <-resultChannel:
	case <-ctx.Done():
		return credentials.Record{}, clierrors.New("timeout", "authentication timed out", 7)
	}
	if login.err != nil {
		return credentials.Record{}, login.err
	}
	token, err := oauthConfig.Exchange(ctx, login.code, oauth2.SetAuthURLParam("code_verifier", verifier))
	if err != nil {
		return credentials.Record{}, fmt.Errorf("exchange authorization code: %w", err)
	}
	rawIDToken, _ := token.Extra("id_token").(string)
	if rawIDToken == "" {
		return credentials.Record{}, clierrors.New("authentication_required", "OIDC provider returned no ID token", 3)
	}
	idToken, err := provider.Verifier(&oidc.Config{ClientID: resolved.Auth.ClientID}).Verify(ctx, rawIDToken)
	if err != nil {
		return credentials.Record{}, clierrors.New("authentication_required", "ID token validation failed", 3)
	}
	var claims claims
	if err := idToken.Claims(&claims); err != nil {
		return credentials.Record{}, fmt.Errorf("decode ID token claims: %w", err)
	}
	if claims.Nonce != nonce {
		return credentials.Record{}, clierrors.New("authentication_required", "OIDC nonce validation failed", 3)
	}
	record := recordFromToken(token, rawIDToken, resolved.Auth.Issuer, resolved.Auth.ClientID, claims)
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	if err := m.Store.Set(key, record); err != nil {
		return credentials.Record{}, err
	}
	return record, nil
}

func (m *Manager) Status(resolved config.Resolved) (credentials.Record, error) {
	if resolved.Token != "" {
		record := credentials.Record{AccessToken: resolved.Token, Issuer: resolved.Auth.Issuer, ClientID: resolved.Auth.ClientID}
		_ = applyJWTClaims(&record, resolved.Token)
		return record, nil
	}
	return m.Store.Get(credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey))
}

func (m *Manager) Logout(resolved config.Resolved) error {
	return m.Store.Delete(credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey))
}

func (m *Manager) EndSession(ctx context.Context, resolved config.Resolved, options LogoutOptions) (LogoutResult, error) {
	if resolved.Auth.Issuer == "" || resolved.Auth.ClientID == "" {
		return LogoutResult{}, clierrors.Validation("active context requires auth.issuer and auth.clientId")
	}
	if options.CallbackPort < 0 || options.CallbackPort > 65535 {
		return LogoutResult{}, clierrors.Validation("callback port must be between 0 and 65535")
	}
	key := credentials.Key(resolved.ContextName, resolved.Auth.CredentialKey)
	record, err := m.Store.Get(key)
	if err != nil {
		if errors.Is(err, credentials.ErrNotFound) {
			return LogoutResult{}, nil
		}
		return LogoutResult{}, err
	}
	endpoint, err := discoverEndSessionEndpoint(ctx, resolved.Auth.Issuer)
	if err != nil {
		return LogoutResult{RemoteAttempted: true}, err
	}
	listenAddress := "127.0.0.1:0"
	if options.CallbackPort > 0 {
		listenAddress = fmt.Sprintf("127.0.0.1:%d", options.CallbackPort)
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return LogoutResult{RemoteAttempted: true}, fmt.Errorf("start logout callback: %w", err)
	}
	defer listener.Close()
	redirectURL := "http://" + listener.Addr().String() + "/callback"
	state, err := randomValue(32)
	if err != nil {
		return LogoutResult{RemoteAttempted: true}, err
	}
	values := url.Values{}
	values.Set("client_id", resolved.Auth.ClientID)
	values.Set("post_logout_redirect_uri", redirectURL)
	values.Set("state", state)
	if record.IDToken != "" {
		values.Set("id_token_hint", record.IDToken)
	}
	logoutURL := endpoint + "?" + values.Encode()

	fmt.Fprintln(options.Out, "CLI logout callback redirect URI:")
	fmt.Fprintln(options.Out, redirectURL)
	fmt.Fprintln(options.Out, "Keycloak must allow this post logout redirect URI, or a wildcard such as http://127.0.0.1:*")
	fmt.Fprintln(options.Out, "Open this URL to sign out:")
	fmt.Fprintln(options.Out, logoutURL)
	if !options.NoBrowser {
		_ = m.OpenBrowser(logoutURL)
	}

	type result struct {
		err error
	}
	resultChannel := make(chan result, 1)
	server := &http.Server{ReadHeaderTimeout: 5 * time.Second}
	server.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/callback" {
			http.NotFound(writer, request)
			return
		}
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "Invalid logout state.", http.StatusBadRequest)
			resultChannel <- result{err: clierrors.New("authentication_required", "OIDC logout state validation failed", 3)}
			return
		}
		if remoteError := request.URL.Query().Get("error"); remoteError != "" {
			http.Error(writer, "Logout failed.", http.StatusBadRequest)
			resultChannel <- result{err: clierrors.New("authentication_required", remoteError, 3)}
			return
		}
		fmt.Fprintln(writer, "A8S CLI logout completed. You may close this window.")
		resultChannel <- result{}
	})
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	select {
	case logout := <-resultChannel:
		if logout.err != nil {
			return LogoutResult{RemoteAttempted: true}, logout.err
		}
	case <-ctx.Done():
		return LogoutResult{RemoteAttempted: true}, clierrors.New("timeout", "logout timed out", 7)
	}
	return LogoutResult{RemoteAttempted: true}, nil
}

func (m *Manager) refresh(ctx context.Context, record credentials.Record) (credentials.Record, error) {
	provider, err := oidc.NewProvider(ctx, record.Issuer)
	if err != nil {
		return credentials.Record{}, err
	}
	config := oauth2.Config{ClientID: record.ClientID, Endpoint: provider.Endpoint()}
	token := &oauth2.Token{AccessToken: record.AccessToken, RefreshToken: record.RefreshToken, Expiry: record.AccessTokenExpiry}
	refreshed, err := config.TokenSource(ctx, token).Token()
	if err != nil {
		return credentials.Record{}, err
	}
	record.AccessToken = refreshed.AccessToken
	record.AccessTokenExpiry = refreshed.Expiry
	if refreshed.RefreshToken != "" {
		record.RefreshToken = refreshed.RefreshToken
	}
	if rawIDToken, ok := refreshed.Extra("id_token").(string); ok && rawIDToken != "" {
		record.IDToken = rawIDToken
		_ = applyJWTClaims(&record, rawIDToken)
	}
	return record, nil
}

func invalidRefreshGrant(err error) bool {
	var retrieveError *oauth2.RetrieveError
	return errors.As(err, &retrieveError) && retrieveError.ErrorCode == "invalid_grant"
}

func discoverEndSessionEndpoint(ctx context.Context, issuer string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(issuer, "/")+"/.well-known/openid-configuration", nil)
	if err != nil {
		return "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("discover OIDC logout endpoint: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf("discover OIDC logout endpoint: unexpected HTTP %d", response.StatusCode)
	}
	var metadata struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := json.NewDecoder(response.Body).Decode(&metadata); err != nil {
		return "", fmt.Errorf("decode OIDC provider metadata: %w", err)
	}
	if strings.TrimSpace(metadata.EndSessionEndpoint) == "" {
		return "", clierrors.New("authentication_required", "OIDC provider does not advertise an end-session endpoint", 3)
	}
	return metadata.EndSessionEndpoint, nil
}

type claims struct {
	Subject  string `json:"sub"`
	Username string `json:"preferred_username"`
	Email    string `json:"email"`
	Nonce    string `json:"nonce"`
	Realm    struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func recordFromToken(token *oauth2.Token, idToken, issuer, clientID string, values claims) credentials.Record {
	return credentials.Record{
		AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, IDToken: idToken,
		AccessTokenExpiry: token.Expiry, Issuer: issuer, ClientID: clientID,
		Subject: values.Subject, Username: values.Username, Email: values.Email, Roles: values.Realm.Roles,
	}
}

func applyJWTClaims(record *credentials.Record, token string) error {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return errors.New("not a JWT")
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	var values claims
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	record.Subject, record.Username, record.Email, record.Roles = values.Subject, values.Username, values.Email, values.Realm.Roles
	return nil
}

func randomValue(size int) (string, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return "", fmt.Errorf("generate secure random value: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
