package deployment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yourname/a8s/internal/api/resources/databases"
)

type DatabaseGetter interface {
	Get(ctx context.Context, id string) (databases.Deployment, error)
}

func WaitForDatabase(ctx context.Context, getter DatabaseGetter, id string, interval time.Duration, progress func(string, ...any)) (databases.Deployment, error) {
	if interval <= 0 {
		interval = 3 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		deployment, err := getter.Get(ctx, id)
		if err != nil {
			return databases.Deployment{}, err
		}
		status := strings.ToUpper(deployment.Status)
		progress("database %s status: %s", id, status)
		switch status {
		case "READY", "RUNNING", "DEPLOYED", "SUCCEEDED", "SUCCESS":
			return deployment, nil
		case "FAILED", "ERROR", "CANCELLED", "CANCELED":
			return deployment, fmt.Errorf("database deployment reached terminal status %s: %s", status, deployment.StatusMessage)
		}
		select {
		case <-ctx.Done():
			return databases.Deployment{}, ctx.Err()
		case <-ticker.C:
		}
	}
}
