package watchcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
)

type route struct {
	path    []string
	socket  string
	summary string
}

var routes = []route{
	{[]string{"notification", "watch"}, "/ws/notifications", "Watch notifications"},
	{[]string{"monitoring", "watch"}, "/ws/monitoring/overview", "Watch monitoring updates"},
	{[]string{"admin", "events", "watch"}, "/ws/admin/events", "Watch administrative events"},
	{[]string{"project", "logs", "websocket"}, "/ws/jenkins/logs", "Watch Jenkins logs over WebSocket"},
}

func Register(root *cobra.Command, runtime *cliruntime.Runtime) {
	for _, value := range routes {
		register(root, runtime, value)
	}
	// The SSE-backed log commands accept --follow for familiar CLI ergonomics.
	for _, path := range [][]string{{"logs"}, {"project", "logs"}} {
		if command := find(root, path); command != nil && command.Flags().Lookup("follow") == nil {
			command.Flags().Bool("follow", false, "keep reading the event stream")
		}
	}
}

func register(root *cobra.Command, runtime *cliruntime.Runtime, value route) {
	parent := root
	for index, name := range value.path {
		existing := child(parent, name)
		if existing != nil {
			parent = existing
			continue
		}
		if index == len(value.path)-1 {
			parent.AddCommand(&cobra.Command{
				Use:   name,
				Short: value.summary,
				RunE: func(cmd *cobra.Command, args []string) error {
					return watch(cmd.Context(), runtime, value.socket)
				},
			})
			return
		}
		group := &cobra.Command{Use: name, Short: "Manage " + name}
		parent.AddCommand(group)
		parent = group
	}
}

func watch(parent context.Context, runtime *cliruntime.Runtime, path string) error {
	if runtime.Config.Token == "" {
		return clierrors.New("authentication_required", "WebSocket watch requires authentication", 3)
	}
	target, err := url.Parse(runtime.Config.Server)
	if err != nil {
		return clierrors.Validation("invalid server URL")
	}
	switch target.Scheme {
	case "http":
		target.Scheme = "ws"
	case "https":
		target.Scheme = "wss"
	default:
		return clierrors.Validation("WebSocket server must use http or https")
	}
	target.Path = strings.TrimRight(target.Path, "/") + path
	query := target.Query()
	query.Set("token", runtime.Config.Token)
	target.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(parent, runtime.Config.Timeout)
	defer cancel()
	connection, response, err := websocket.DefaultDialer.DialContext(ctx, target.String(), http.Header{"User-Agent": []string{"a8s-cli"}})
	if err != nil {
		if response != nil {
			return clierrors.FromHTTP(response.StatusCode, "WebSocket handshake failed", response.Header.Get("X-Request-ID"))
		}
		return &clierrors.Error{Code: "backend_unavailable", Message: err.Error(), Exit: 8, Cause: err}
	}
	defer connection.Close()

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return clierrors.New("timeout", "watch timed out", 7)
			}
			return fmt.Errorf("read WebSocket message: %w", err)
		}
		if messageType != websocket.TextMessage {
			continue
		}
		if runtime.Config.Output == "json" || runtime.Config.Output == "yaml" {
			var value any
			if json.Unmarshal(message, &value) == nil {
				if err := runtime.Printer.Print(value); err != nil {
					return err
				}
				continue
			}
		}
		if _, err := fmt.Fprintln(runtime.Out, string(message)); err != nil {
			return err
		}
	}
}

func find(root *cobra.Command, path []string) *cobra.Command {
	current := root
	for _, name := range path {
		current = child(current, name)
		if current == nil {
			return nil
		}
	}
	return current
}

func child(parent *cobra.Command, name string) *cobra.Command {
	for _, command := range parent.Commands() {
		if command.Name() == name {
			return command
		}
	}
	return nil
}
