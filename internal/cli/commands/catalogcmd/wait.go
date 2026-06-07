package catalogcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
)

const maxWaitBody = 1 << 20

func waitForResponse(ctx context.Context, runtime *cliruntime.Runtime, method, endpoint, requestPath string, response *http.Response) error {
	data, err := io.ReadAll(io.LimitReader(response.Body, maxWaitBody))
	if err != nil {
		return fmt.Errorf("read wait response: %w", err)
	}
	initial := map[string]any{}
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &initial); err != nil {
			return fmt.Errorf("decode wait response: %w", err)
		}
	}
	if terminal(statusOf(initial)) {
		return runtime.Printer.Print(initial)
	}
	pollPath, err := pollPath(endpoint, requestPath, response.Header, initial)
	if err != nil {
		return err
	}
	interval := runtime.Config.PollingInterval
	if interval <= 0 {
		interval = 3 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	current := initial
	for {
		status := statusOf(current)
		if status != "" {
			runtime.Printer.Progress("operation status: %s", strings.ToUpper(status))
		}
		if terminal(status) {
			if failed(status) {
				return clierrors.New("operation_failed", fmt.Sprintf("operation reached terminal status %s", status), 1)
			}
			return runtime.Printer.Print(current)
		}
		select {
		case <-ctx.Done():
			return clierrors.New("timeout", "operation wait timed out", 7)
		case <-ticker.C:
		}
		resp, err := runtime.API.Do(ctx, http.MethodGet, pollPath, nil)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(io.LimitReader(resp.Body, maxWaitBody))
		_ = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("read wait poll response: %w", err)
		}
		current = map[string]any{}
		if err := json.Unmarshal(data, &current); err != nil {
			return fmt.Errorf("decode wait poll response: %w", err)
		}
	}
}

func pollPath(endpoint, requestPath string, header http.Header, body map[string]any) (string, error) {
	if location := header.Get("Location"); location != "" {
		return relativeLocation(location)
	}
	for _, key := range []string{"statusPath", "statusUrl", "operationPath", "operationUrl"} {
		if value := stringField(body, key); value != "" {
			return relativeLocation(value)
		}
	}
	switch {
	case strings.Contains(endpoint, "/image-scanner/scans"):
		id := firstField(body, "scanId", "id")
		if id != "" {
			return "/api/v1/image-scanner/scans/" + url.PathEscape(id), nil
		}
	case strings.Contains(endpoint, "/quota-requests"):
		md5 := firstField(body, "md5", "paymentMd5", "paymentHash", "paymentStatusMd5")
		if md5 != "" {
			return "/api/v1/workspaces/quota-requests/payment-status?md5=" + url.QueryEscape(md5), nil
		}
	case strings.Contains(endpoint, "/cluster-deployments"):
		release := firstField(body, "releaseName", "name")
		if release != "" {
			base := strings.Split(requestPath, "?")[0]
			return strings.TrimRight(base, "/") + "/" + url.PathEscape(release), nil
		}
	case strings.Contains(endpoint, "/database-deployments/{deploymentId}/backup/runs"):
		deploymentID := pathValue(requestPath, "/api/v1/database-deployments/", "/backup/")
		if deploymentID != "" {
			return "/api/v1/database-deployments/" + url.PathEscape(deploymentID), nil
		}
	}
	return "", clierrors.Validation("backend response did not include an operation id or status URL required by --wait")
}

func relativeLocation(value string) (string, error) {
	parsed, err := url.Parse(value)
	if err != nil {
		return "", clierrors.Validation("invalid status URL returned by backend")
	}
	if parsed.IsAbs() || parsed.Host != "" {
		return "", clierrors.Validation("backend returned a non-relative status URL")
	}
	if parsed.Path == "" {
		return "", clierrors.Validation("backend returned an empty status URL")
	}
	if parsed.RawQuery != "" {
		return parsed.Path + "?" + parsed.RawQuery, nil
	}
	return parsed.Path, nil
}

func statusOf(value map[string]any) string {
	return strings.ToUpper(firstField(value, "status", "state", "phase", "paymentStatus"))
}

func terminal(status string) bool {
	return succeeded(status) || failed(status)
}

func succeeded(status string) bool {
	switch status {
	case "SUCCEEDED", "SUCCESS", "COMPLETED", "COMPLETE", "READY", "DEPLOYED", "PAID", "NO_PAYMENT_REQUIRED":
		return true
	default:
		return false
	}
}

func failed(status string) bool {
	switch status {
	case "FAILED", "ERROR", "CANCELLED", "CANCELED", "REJECTED":
		return true
	default:
		return false
	}
}

func firstField(value map[string]any, keys ...string) string {
	for _, key := range keys {
		if result := stringField(value, key); result != "" {
			return result
		}
	}
	return ""
}

func stringField(value map[string]any, key string) string {
	raw, ok := value[key]
	if !ok || raw == nil {
		return ""
	}
	switch typed := raw.(type) {
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func pathValue(path, prefix, suffix string) string {
	withoutQuery := strings.Split(path, "?")[0]
	start := strings.Index(withoutQuery, prefix)
	if start < 0 {
		return ""
	}
	remainder := withoutQuery[start+len(prefix):]
	end := strings.Index(remainder, suffix)
	if end < 0 {
		return ""
	}
	value, err := url.PathUnescape(remainder[:end])
	if err != nil {
		return ""
	}
	return value
}
