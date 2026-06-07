package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/yourname/a8s/internal/clierrors"
)

type Client struct {
	BaseURL      string
	Token        string
	HTTPClient   *http.Client
	RefreshToken func(context.Context) (string, error)
}

type RequestBody struct {
	Reader      io.Reader
	ContentType string
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: newHTTPClient(baseURL, 15*time.Second),
	}
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func New(baseURL, token string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Token:      token,
		HTTPClient: newHTTPClient(baseURL, timeout),
	}
}

func newHTTPClient(baseURL string, timeout time.Duration) *http.Client {
	origin, _ := url.Parse(baseURL)
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if origin == nil || !sameOrigin(origin, request.URL) {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

func (c *Client) ConfigureTLS(insecureSkipVerify bool, caFile string) error {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: insecureSkipVerify} //nolint:gosec -- explicit development-only config
	if caFile != "" {
		data, err := os.ReadFile(caFile)
		if err != nil {
			return fmt.Errorf("read CA file: %w", err)
		}
		pool, err := x509.SystemCertPool()
		if err != nil {
			return fmt.Errorf("load system CA pool: %w", err)
		}
		if !pool.AppendCertsFromPEM(data) {
			return fmt.Errorf("CA file contains no valid PEM certificates")
		}
		transport.TLSClientConfig.RootCAs = pool
	}
	c.HTTPClient.Transport = transport
	return nil
}

func (c *Client) DoJSON(ctx context.Context, method, path string, requestBody, responseBody any) error {
	resp, err := c.Do(ctx, method, path, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if responseBody == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// Do performs an authenticated request and returns successful responses without
// interpreting their body. Callers must close the returned response body.
func (c *Client) Do(ctx context.Context, method, path string, requestBody any) (*http.Response, error) {
	body, contentType, err := encodeRequestBody(requestBody)
	if err != nil {
		return nil, err
	}
	target, err := c.resolve(path)
	if err != nil {
		return nil, err
	}

	resp, err := c.execute(ctx, method, target, body, contentType)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && c.RefreshToken != nil {
		_ = resp.Body.Close()
		token, refreshErr := c.RefreshToken(ctx)
		if refreshErr != nil {
			return nil, refreshErr
		}
		if token != "" {
			c.Token = token
			resp, err = c.execute(ctx, method, target, body, contentType)
			if err != nil {
				return nil, err
			}
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, decodeError(resp)
	}
	return resp, nil
}

func (c *Client) execute(ctx context.Context, method, target string, body []byte, contentType string) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, reader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "a8s-cli")
	if body != nil {
		if contentType == "" {
			contentType = "application/json"
		}
		req.Header.Set("Content-Type", contentType)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, &clierrors.Error{Code: "backend_unavailable", Message: err.Error(), Exit: 8, Cause: err}
	}
	return resp, nil
}

func encodeRequestBody(requestBody any) ([]byte, string, error) {
	if requestBody == nil {
		return nil, "", nil
	}
	switch value := requestBody.(type) {
	case RequestBody:
		data, err := io.ReadAll(value.Reader)
		if err != nil {
			return nil, "", fmt.Errorf("read request body: %w", err)
		}
		return data, value.ContentType, nil
	case io.Reader:
		data, err := io.ReadAll(value)
		if err != nil {
			return nil, "", fmt.Errorf("read request body: %w", err)
		}
		return data, "", nil
	case []byte:
		return value, "", nil
	default:
		data, err := json.Marshal(requestBody)
		if err != nil {
			return nil, "", fmt.Errorf("encode request: %w", err)
		}
		return data, "", nil
	}
}

func (c *Client) resolve(path string) (string, error) {
	base, err := url.Parse(strings.TrimRight(c.BaseURL, "/") + "/")
	if err != nil {
		return "", fmt.Errorf("invalid server URL: %w", err)
	}
	relative, err := url.Parse(strings.TrimLeft(path, "/"))
	if err != nil {
		return "", fmt.Errorf("invalid request path: %w", err)
	}
	if relative.IsAbs() || relative.Host != "" {
		return "", clierrors.Validation("request path must be relative to the configured backend server")
	}
	return base.ResolveReference(relative).String(), nil
}

func sameOrigin(left, right *url.URL) bool {
	return strings.EqualFold(left.Scheme, right.Scheme) && strings.EqualFold(left.Host, right.Host)
}

func decodeError(resp *http.Response) error {
	const maxBody = 1 << 20
	data, _ := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	message := strings.TrimSpace(string(data))
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
		Detail  string `json:"detail"`
	}
	if json.Unmarshal(data, &payload) == nil {
		message = firstMessage(payload.Message, payload.Detail, payload.Error, message)
	}
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}
	requestID := firstMessage(resp.Header.Get("X-Request-ID"), resp.Header.Get("Request-ID"))
	return clierrors.FromHTTP(resp.StatusCode, message, requestID)
}

func firstMessage(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
