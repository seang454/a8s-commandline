package catalogcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourname/a8s/internal/api"
	"github.com/yourname/a8s/internal/cli/catalog"
	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
	"github.com/yourname/a8s/internal/operation"
)

type requestFlags struct {
	file       string
	set        []string
	query      []string
	form       []string
	uploads    []string
	outputFile string
	yes        bool
	dryRun     bool
	wait       bool
}

type convenienceFlags struct {
	container    string
	tail         int
	warningsOnly bool
	limit        int
	job          string
	build        int
	queueItem    int
	full         bool
	image        string
	releaseName  string
	projectName  string
	name         string
	environment  string
}

// RegisterRoutes adds an explicit feature-owned route list.
func RegisterRoutes(root *cobra.Command, runtime *cliruntime.Runtime, routes []catalog.Route) {
	for _, route := range routes {
		registerRoute(root, runtime, route)
	}
}

// RegisterUtilities adds cross-feature compatibility and catalog commands.
func RegisterUtilities(root *cobra.Command, runtime *cliruntime.Runtime) {
	registerQuotaPurchase(root, runtime)
	root.AddCommand(newAPICommand(runtime))
}

func registerQuotaPurchase(root *cobra.Command, runtime *cliruntime.Runtime) {
	workspace := child(root, "workspace")
	if workspace == nil {
		return
	}
	quota := child(workspace, "quota")
	if quota == nil || child(quota, "purchase") != nil {
		return
	}
	var flags requestFlags
	var plan string
	command := &cobra.Command{
		Use:   "purchase",
		Short: "Purchase a workspace quota plan using Bakong KHQR",
		Annotations: map[string]string{
			"a8s.io/endpoint": "/api/v1/workspaces/quota-requests",
			"a8s.io/method":   http.MethodPost,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved := flags
			if cmd.Flags().Changed("plan") {
				resolved.set = append(resolved.set, "planName="+plan, "isPaid=true", "paymentProvider=BAKONG")
			}
			return execute(cmd.Context(), runtime, http.MethodPost, "/api/v1/workspaces/quota-requests", nil, nil, resolved)
		},
	}
	addFlags(command, http.MethodPost, &flags)
	command.Flags().StringVar(&plan, "plan", "", "quota plan name")
	quota.AddCommand(command)
}

func registerRoute(root *cobra.Command, runtime *cliruntime.Runtime, route catalog.Route) {
	parent := root
	for index, name := range route.Command {
		existing := child(parent, name)
		if existing != nil {
			parent = existing
			continue
		}
		if index == len(route.Command)-1 {
			leaf := newRouteCommand(runtime, route)
			parent.AddCommand(leaf)
			return
		}
		group := &cobra.Command{Use: name, Short: "Manage " + strings.ReplaceAll(name, "-", " ")}
		parent.AddCommand(group)
		parent = group
	}
	// A specialized command already owns this exact command path.
}

func child(parent *cobra.Command, name string) *cobra.Command {
	for _, command := range parent.Commands() {
		if command.Name() == name {
			return command
		}
	}
	return nil
}

func newRouteCommand(runtime *cliruntime.Runtime, route catalog.Route) *cobra.Command {
	var flags requestFlags
	var convenience convenienceFlags
	use := route.Command[len(route.Command)-1]
	for _, arg := range route.Args {
		use += " <" + arg + ">"
	}
	command := &cobra.Command{
		Use:   use,
		Short: route.Method + " " + route.Endpoint,
		Args:  cobra.ExactArgs(len(route.Args)),
		Annotations: map[string]string{
			"a8s.io/controller": route.Controller,
			"a8s.io/endpoint":   route.Endpoint,
			"a8s.io/method":     route.Method,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved := flags
			appendConvenienceQuery(cmd, route.Endpoint, &resolved, convenience)
			appendConvenienceSet(cmd, route.Endpoint, &resolved, convenience)
			endpoint := route.Endpoint
			if cmd.Flags().Changed("full") && convenience.full {
				endpoint = strings.TrimRight(endpoint, "/") + "/full"
			}
			return execute(cmd.Context(), runtime, route.Method, endpoint, route.Args, args, resolved)
		},
	}
	addFlags(command, route.Method, &flags)
	addConvenienceFlags(command, route.Endpoint, &convenience)
	return command
}

func newAPICommand(runtime *cliruntime.Runtime) *cobra.Command {
	var flags requestFlags
	api := &cobra.Command{Use: "api", Short: "Access backend API routes directly"}
	request := &cobra.Command{
		Use:   "request <method> <path>",
		Short: "Send an authenticated request to any backend route",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(cmd.Context(), runtime, strings.ToUpper(args[0]), args[1], nil, nil, flags)
		},
	}
	addFlags(request, "*", &flags)
	var search string
	catalogCommand := &cobra.Command{
		Use:   "catalog",
		Short: "List implemented backend route mappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			var routes []catalog.Route
			for _, route := range catalog.Routes {
				value := route.Method + " " + route.Endpoint + " " + strings.Join(route.Command, " ") + " " + route.Controller
				if search == "" || strings.Contains(strings.ToLower(value), strings.ToLower(search)) {
					routes = append(routes, route)
				}
			}
			return runtime.Printer.Print(routes)
		},
	}
	catalogCommand.Flags().StringVar(&search, "search", "", "filter by method, endpoint, command, or controller")
	api.AddCommand(request, catalogCommand)
	return api
}

func addFlags(command *cobra.Command, method string, flags *requestFlags) {
	endpoint := command.Annotations["a8s.io/endpoint"]
	command.Flags().StringVar(&flags.file, "file", "", "YAML or JSON request body; operation envelopes use their spec")
	command.Flags().StringArrayVar(&flags.set, "set", nil, "set request field using dotted key=value; repeatable")
	command.Flags().StringArrayVar(&flags.query, "query", nil, "add query parameter using key=value; repeatable")
	command.Flags().StringArrayVar(&flags.form, "form", nil, "add multipart form field using key=value; repeatable")
	command.Flags().StringArrayVar(&flags.uploads, "upload", nil, "upload a file using field=path; repeatable")
	command.Flags().StringVar(&flags.outputFile, "output-file", "", "write the response body to a file")
	command.Flags().BoolVar(&flags.yes, "yes", false, "confirm a destructive operation")
	command.Flags().BoolVar(&flags.dryRun, "dry-run", false, "print the resolved request without sending it")
	if supportsWait(method, endpoint) {
		command.Flags().BoolVar(&flags.wait, "wait", false, "wait for the asynchronous operation to finish")
	}
	if !supportsRequestBody(method, endpoint) {
		_ = command.Flags().MarkHidden("file")
		_ = command.Flags().MarkHidden("set")
		_ = command.Flags().MarkHidden("form")
		_ = command.Flags().MarkHidden("upload")
	}
	if method == http.MethodGet || (method != "*" && !destructive(method, endpoint)) {
		_ = command.Flags().MarkHidden("yes")
	}
	if method == http.MethodGet {
		_ = command.Flags().MarkHidden("dry-run")
	}
}

func execute(ctx context.Context, runtime *cliruntime.Runtime, method, endpoint string, argNames, args []string, flags requestFlags) error {
	if destructive(method, endpoint) && !flags.yes {
		return clierrors.Validation("destructive operation requires --yes")
	}
	if !supportsRequestBody(method, endpoint) && hasRequestBodyInput(flags) {
		return clierrors.Validation("this command does not accept --file, --set, --form, or --upload")
	}
	path, remaining, err := resolvePath(endpoint, runtime.Config.Namespace, argNames, args)
	if err != nil {
		return err
	}
	query := url.Values{}
	for index, value := range remaining {
		query.Add(argNames[len(argNames)-len(remaining)+index], value)
	}
	if err := addPairs(query, flags.query); err != nil {
		return err
	}
	if strings.HasPrefix(endpoint, "/api/kubernetes") && runtime.Config.TargetCluster != "" && !query.Has("targetClusterName") {
		query.Set("targetClusterName", runtime.Config.TargetCluster)
	}
	if encoded := query.Encode(); encoded != "" {
		separator := "?"
		if strings.Contains(path, "?") {
			separator = "&"
		}
		path += separator + encoded
	}

	var body any
	if len(flags.form) > 0 || len(flags.uploads) > 0 {
		if flags.file != "" || len(flags.set) > 0 {
			return clierrors.Validation("--form/--upload cannot be combined with --file/--set")
		}
		body, err = multipartBody(flags.form, flags.uploads)
		if err != nil {
			return err
		}
	}
	if flags.file != "" || len(flags.set) > 0 {
		payload, err := operation.LoadGeneric(flags.file, runtime.In)
		if err != nil {
			return err
		}
		if err := operation.ApplySet(payload, flags.set); err != nil {
			return err
		}
		body = payload
	}
	if flags.dryRun {
		return runtime.Printer.Print(map[string]any{"method": method, "path": path, "body": redact(body)})
	}

	requestContext, cancel := context.WithTimeout(ctx, runtime.Config.Timeout)
	defer cancel()
	response, err := runtime.API.Do(requestContext, method, path, body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if flags.outputFile != "" {
		return writeResponseFile(flags.outputFile, response.Body)
	}
	if flags.wait {
		return waitForResponse(requestContext, runtime, method, endpoint, path, response)
	}
	if response.StatusCode == http.StatusNoContent {
		return runtime.Printer.Print(map[string]any{"status": "success", "httpStatus": response.StatusCode})
	}
	return printResponse(runtime, response)
}

func redact(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			lower := strings.ToLower(key)
			if strings.Contains(lower, "password") || strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "credential") {
				result[key] = "[redacted]"
			} else {
				result[key] = redact(item)
			}
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			result[index] = redact(item)
		}
		return result
	default:
		return value
	}
}

func addConvenienceFlags(command *cobra.Command, endpoint string, flags *convenienceFlags) {
	switch {
	case strings.Contains(endpoint, "/pods/{podName}/logs/stream"):
		command.Flags().StringVar(&flags.container, "container", "", "pod container name")
		command.Flags().IntVar(&flags.tail, "tail", 0, "number of previous log lines")
	case strings.HasSuffix(endpoint, "/events"):
		command.Flags().BoolVar(&flags.warningsOnly, "warnings-only", false, "show warning events only")
		command.Flags().IntVar(&flags.limit, "limit", 0, "maximum number of events")
	case endpoint == "/api/v1/jenkins/logs/stream":
		command.Flags().StringVar(&flags.job, "job", "", "Jenkins job name")
		command.Flags().IntVar(&flags.build, "build", 0, "Jenkins build number")
		command.Flags().IntVar(&flags.queueItem, "queue-item", 0, "Jenkins queue item number")
	case strings.Contains(endpoint, "/clusters/{id}/values"):
		command.Flags().BoolVar(&flags.full, "full", false, "return complete deployment values")
	case endpoint == "/api/v1/image-scanner/scans":
		command.Flags().StringVar(&flags.image, "image", "", "container image reference to scan")
	case strings.HasSuffix(endpoint, "/cluster-deployments"):
		command.Flags().StringVar(&flags.releaseName, "release-name", "", "cluster deployment release name")
		command.Flags().StringVar(&flags.projectName, "project-name", "", "owning project name")
		command.Flags().StringVar(&flags.name, "name", "", "cluster name")
		command.Flags().StringVar(&flags.environment, "environment", "", "deployment environment")
	}
}

func appendConvenienceQuery(command *cobra.Command, endpoint string, flags *requestFlags, values convenienceFlags) {
	add := func(name, value string) {
		if command.Flags().Changed(name) {
			flags.query = append(flags.query, value)
		}
	}
	switch {
	case strings.Contains(endpoint, "/pods/{podName}/logs/stream"):
		add("container", "container="+values.container)
		add("tail", fmt.Sprintf("tailLines=%d", values.tail))
	case strings.HasSuffix(endpoint, "/events"):
		add("warnings-only", fmt.Sprintf("warningsOnly=%t", values.warningsOnly))
		add("limit", fmt.Sprintf("limit=%d", values.limit))
	case endpoint == "/api/v1/jenkins/logs/stream":
		add("job", "job="+values.job)
		add("build", fmt.Sprintf("build=%d", values.build))
		add("queue-item", fmt.Sprintf("queueItem=%d", values.queueItem))
	}
}

func appendConvenienceSet(command *cobra.Command, endpoint string, flags *requestFlags, values convenienceFlags) {
	add := func(flag, field, value string) {
		if command.Flags().Changed(flag) {
			flags.set = append(flags.set, field+"="+value)
		}
	}
	switch {
	case endpoint == "/api/v1/image-scanner/scans":
		add("image", "image", values.image)
	case strings.HasSuffix(endpoint, "/cluster-deployments"):
		add("release-name", "releaseName", values.releaseName)
		add("project-name", "projectName", values.projectName)
		add("name", "cluster.name", values.name)
		add("environment", "cluster.environment", values.environment)
	}
}

func multipartBody(fields, uploads []string) (api.RequestBody, error) {
	var buffer strings.Builder
	writer := multipart.NewWriter(&buffer)
	for _, field := range fields {
		key, value, ok := strings.Cut(field, "=")
		if !ok || key == "" {
			return api.RequestBody{}, clierrors.Validation(fmt.Sprintf("invalid --form %q; expected field=value", field))
		}
		if err := writer.WriteField(key, value); err != nil {
			return api.RequestBody{}, fmt.Errorf("write multipart field: %w", err)
		}
	}
	for _, upload := range uploads {
		field, path, ok := strings.Cut(upload, "=")
		if !ok || field == "" || path == "" {
			return api.RequestBody{}, clierrors.Validation(fmt.Sprintf("invalid --upload %q; expected field=path", upload))
		}
		file, err := os.Open(path)
		if err != nil {
			return api.RequestBody{}, fmt.Errorf("open upload %q: %w", path, err)
		}
		part, err := writer.CreateFormFile(field, filepath.Base(path))
		if err == nil {
			_, err = io.Copy(part, file)
		}
		closeErr := file.Close()
		if err != nil {
			return api.RequestBody{}, fmt.Errorf("write upload %q: %w", path, err)
		}
		if closeErr != nil {
			return api.RequestBody{}, fmt.Errorf("close upload %q: %w", path, closeErr)
		}
	}
	if err := writer.Close(); err != nil {
		return api.RequestBody{}, fmt.Errorf("finish multipart request: %w", err)
	}
	return api.RequestBody{Reader: strings.NewReader(buffer.String()), ContentType: writer.FormDataContentType()}, nil
}

func resolvePath(endpoint, namespace string, argNames, args []string) (string, []string, error) {
	path := endpoint
	argIndex := 0
	for strings.Contains(path, "{") {
		start := strings.Index(path, "{")
		end := strings.Index(path[start:], "}")
		if end < 0 {
			return "", nil, clierrors.Validation("invalid endpoint path template")
		}
		end += start
		variable := path[start+1 : end]
		value := ""
		if variable == "namespace" && namespace != "" && !contains(argNames, "namespace") {
			value = namespace
		} else if argIndex < len(args) {
			value = args[argIndex]
			argIndex++
		}
		if value == "" {
			return "", nil, clierrors.Validation(fmt.Sprintf("endpoint requires %s", variable))
		}
		path = path[:start] + url.PathEscape(value) + path[end+1:]
	}
	return path, args[argIndex:], nil
}

func addPairs(values url.Values, pairs []string) error {
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return clierrors.Validation(fmt.Sprintf("invalid key=value pair %q", pair))
		}
		values.Add(key, value)
	}
	return nil
}

func destructive(method, endpoint string) bool {
	if method == http.MethodDelete {
		return true
	}
	for _, value := range []string{"/deactivate", "/reject", "/restore", "/rollback", "/rotate-password", "/abort", "/clear"} {
		if strings.Contains(endpoint, value) {
			return true
		}
	}
	return false
}

func hasRequestBodyInput(flags requestFlags) bool {
	return flags.file != "" || len(flags.set) > 0 || len(flags.form) > 0 || len(flags.uploads) > 0
}

func supportsRequestBody(method, endpoint string) bool {
	if method == "*" || method == "" {
		return true
	}
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodDelete:
		return false
	}
	for _, pattern := range []string{
		"/abort",
		"/approve",
		"/backup/run",
		"/bootstrap",
		"/console/test",
		"/deactivate",
		"/domain/sync",
		"/notifications/read",
		"/onboarding",
		"/publish",
		"/reactivate",
		"/redeploy",
		"/retry",
		"/restore",
		"/restore/cancel",
		"/sonarqube/access",
		"/sync",
		"/sync-keycloak-token",
		"/trigger",
		"/verify-email",
		"/webhook/rotate",
	} {
		if strings.Contains(endpoint, pattern) {
			return false
		}
	}
	return true
}

func supportsWait(method, endpoint string) bool {
	if method != http.MethodPost {
		return false
	}
	if endpoint == "/api/v1/workspaces/quota-requests" {
		return true
	}
	for _, value := range []string{
		"/cluster-deployments",
		"/image-scanner/scans",
		"/backup/run",
		"/backups/trigger/",
		"/restore/",
		"/runs/{runId}/restore",
	} {
		if strings.Contains(endpoint, value) {
			return true
		}
	}
	return false
}

func printResponse(runtime *cliruntime.Runtime, response *http.Response) error {
	contentType, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if contentType == "application/json" || strings.HasSuffix(contentType, "+json") {
		var value any
		if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		return runtime.Printer.Print(value)
	}
	_, err := io.Copy(runtime.Out, response.Body)
	return err
}

func writeResponseFile(path string, source io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, source); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
