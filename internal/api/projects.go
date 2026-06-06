package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourname/a8s/internal/models"
)

func (c *Client) ListProjects() ([]models.Project, error) {
	req, err := c.newRequest(http.MethodGet, "/api/projects")
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var projects []models.Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return projects, nil
}
