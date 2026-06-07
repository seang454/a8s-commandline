package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourname/a8s/internal/clierrors"
)

func TestDoJSONMapsBackendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-123")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"already exists"}`))
	}))
	defer server.Close()

	client := New(server.URL, "", time.Second)
	err := client.DoJSON(context.Background(), http.MethodPost, "/resource", map[string]string{"name": "x"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if clierrors.ExitCode(err) != 6 {
		t.Fatalf("expected conflict exit code, got %d: %v", clierrors.ExitCode(err), err)
	}
}

func TestDoRefreshesAndRetriesUnauthorizedRequestOnce(t *testing.T) {
	var requests int
	var bodies []string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		data, _ := io.ReadAll(request.Body)
		bodies = append(bodies, string(data))
		if requests == 1 {
			if request.Header.Get("Authorization") != "Bearer old-token" {
				t.Fatalf("unexpected first token: %q", request.Header.Get("Authorization"))
			}
			http.Error(writer, "expired", http.StatusUnauthorized)
			return
		}
		if request.Header.Get("Authorization") != "Bearer new-token" {
			t.Fatalf("unexpected retry token: %q", request.Header.Get("Authorization"))
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := New(server.URL, "old-token", time.Second)
	var refreshes int
	client.RefreshToken = func(context.Context) (string, error) {
		refreshes++
		return "new-token", nil
	}
	var result map[string]any
	err := client.DoJSON(context.Background(), http.MethodPost, "/resource", map[string]string{"name": "orders"}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 || refreshes != 1 {
		t.Fatalf("expected one refresh and one retry, requests=%d refreshes=%d", requests, refreshes)
	}
	if len(bodies) != 2 || bodies[0] != bodies[1] || !strings.Contains(bodies[0], "orders") {
		t.Fatalf("request body was not replayed exactly: %#v", bodies)
	}
	if result["status"] != "ok" {
		t.Fatalf("unexpected response: %#v", result)
	}
}

func TestDoStopsAfterOneUnauthorizedRetry(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New(server.URL, "old-token", time.Second)
	var refreshes int
	client.RefreshToken = func(context.Context) (string, error) {
		refreshes++
		return "new-token", nil
	}
	_, err := client.Do(context.Background(), http.MethodGet, "/resource", nil)
	if clierrors.ExitCode(err) != 3 {
		t.Fatalf("expected authentication exit code 3, got %d: %v", clierrors.ExitCode(err), err)
	}
	if requests != 2 || refreshes != 1 {
		t.Fatalf("request retried more than once: requests=%d refreshes=%d", requests, refreshes)
	}
}

func TestDoDoesNotRetryWhenRefreshFails(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New(server.URL, "old-token", time.Second)
	client.RefreshToken = func(context.Context) (string, error) {
		return "", clierrors.New("authentication_required", "refresh failed", 3)
	}
	_, err := client.Do(context.Background(), http.MethodGet, "/resource", nil)
	if clierrors.ExitCode(err) != 3 {
		t.Fatalf("expected authentication exit code 3, got %d: %v", clierrors.ExitCode(err), err)
	}
	if requests != 1 {
		t.Fatalf("request retried after refresh failure: requests=%d", requests)
	}
}

func TestDoReplaysMultipartRequestBodyAfterRefresh(t *testing.T) {
	var bodies [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		data, _ := io.ReadAll(request.Body)
		bodies = append(bodies, data)
		if len(bodies) == 1 {
			http.Error(writer, "expired", http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "old-token", time.Second)
	client.RefreshToken = func(context.Context) (string, error) { return "new-token", nil }
	body := RequestBody{Reader: strings.NewReader("multipart-payload"), ContentType: "multipart/form-data; boundary=test"}
	response, err := client.Do(context.Background(), http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if len(bodies) != 2 || string(bodies[0]) != string(bodies[1]) {
		encoded, _ := json.Marshal(bodies)
		t.Fatalf("multipart body was not replayed: %s", encoded)
	}
}

func TestDoRejectsAbsoluteRequestURL(t *testing.T) {
	client := New("https://api.example.com", "secret-token", time.Second)
	_, err := client.Do(context.Background(), http.MethodGet, "https://other.example.com/resource", nil)
	if clierrors.ExitCode(err) != 2 {
		t.Fatalf("expected validation exit code 2, got %d: %v", clierrors.ExitCode(err), err)
	}
}

func TestDoDoesNotFollowCrossOriginRedirect(t *testing.T) {
	var redirectedRequests int
	redirected := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		redirectedRequests++
		if request.Header.Get("Authorization") != "" {
			t.Errorf("authorization header leaked to redirected host")
		}
		writer.WriteHeader(http.StatusOK)
	}))
	defer redirected.Close()
	backend := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, redirected.URL+"/resource", http.StatusFound)
	}))
	defer backend.Close()

	client := New(backend.URL, "secret-token", time.Second)
	_, err := client.Do(context.Background(), http.MethodGet, "/redirect", nil)
	if err == nil {
		t.Fatal("expected cross-origin redirect response to fail")
	}
	if redirectedRequests != 0 {
		t.Fatalf("cross-origin redirect was followed %d time(s)", redirectedRequests)
	}
}

func TestDoFollowsSameOriginRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/redirect" {
			http.Redirect(writer, request, "/resource", http.StatusFound)
			return
		}
		if request.Header.Get("Authorization") != "Bearer secret-token" {
			t.Fatalf("same-origin request lost authorization header: %q", request.Header.Get("Authorization"))
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "secret-token", time.Second)
	response, err := client.Do(context.Background(), http.MethodGet, "/redirect", nil)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
}
