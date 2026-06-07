package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yourname/a8s/internal/api"
)

type Client struct {
	api *api.Client
}

func (c *Client) Get(ctx context.Context, id string) (Deployment, error) {
	var response Deployment
	err := c.api.DoJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/database-deployments/%s", id), nil, &response)
	return response, err
}

func New(client *api.Client) *Client {
	return &Client{api: client}
}

func (c *Client) Deploy(ctx context.Context, request DeployRequest) (Deployment, error) {
	var response Deployment
	err := c.api.DoJSON(ctx, http.MethodPost, "/api/v1/database-deployments", request, &response)
	return response, err
}
