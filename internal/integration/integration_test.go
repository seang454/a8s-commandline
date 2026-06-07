package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/yourname/a8s/internal/api"
)

func TestAuthenticatedBackendSmoke(t *testing.T) {
	if os.Getenv("A8S_RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("set A8S_RUN_INTEGRATION_TESTS=true to run backend integration tests")
	}
	server := os.Getenv("A8S_SERVER")
	token := os.Getenv("A8S_TOKEN")
	if server == "" || token == "" {
		t.Fatal("A8S_SERVER and A8S_TOKEN are required for integration tests")
	}
	client := api.New(server, token, 20*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := client.Do(ctx, http.MethodGet, "/actuator/health", nil)
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	_ = response.Body.Close()

	response, err = client.Do(ctx, http.MethodGet, "/api/v1/profile/me", nil)
	if err != nil {
		t.Fatalf("authenticated profile check failed: %v", err)
	}
	_ = response.Body.Close()
}
