// Package apiexec executes the API-backed commands declared in opencli.yaml.
package apiexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

type Operation struct {
	ID      string `json:"id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
	item    *openapi3.PathItem
	spec    *openapi3.Operation
}

type Executor struct {
	store      *config.Store
	spec       *openapi3.T
	operations map[string]Operation
}

func New(store *config.Store) (*Executor, error) {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return nil, err
	}
	e := &Executor{store: store, spec: spec, operations: map[string]Operation{}}
	for path, item := range spec.Paths.Map() {
		for method, op := range item.Operations() {
			e.operations[op.OperationID] = Operation{op.OperationID, method, path, op.Summary, item, op}
		}
	}
	return e, nil
}

type Request struct {
	Operation   string
	Path        map[string]string
	Query       url.Values
	Data        string
	File        string
	ContentType string
	Output      string
	OutputFile  string
	Timeout     string
}

func (e *Executor) Run(ctx context.Context, streams cligen.IO, input Request) error {
	op, ok := e.operations[input.Operation]
	if !ok {
		return fmt.Errorf("unknown API operation %q; use sd api operations", input.Operation)
	}
	if input.Output != "json" && input.Output != "raw" {
		return fmt.Errorf("--output must be json or raw")
	}
	duration, err := time.ParseDuration(input.Timeout)
	if err != nil || duration < 0 {
		return fmt.Errorf("--timeout must be a non-negative duration, such as 60s or 0s")
	}
	if duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}
	req, err := e.prepareRequest(ctx, op, streams.In, input)
	if err != nil {
		return err
	}
	client, err := api.NewAuthenticatedClient(ctx, e.store, api.WithRateLimitWarnings(streams.Err))
	if err != nil {
		return err
	}
	target, err := url.Parse(strings.TrimRight(client.BaseURL, "/") + req.URL.EscapedPath())
	if err != nil {
		return err
	}
	target.RawQuery = req.URL.RawQuery
	req.URL = target
	req.Host = ""
	response, err := client.Do(req, streams.Err)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 65536))
		return fmt.Errorf("%s failed (HTTP %d): %s", op.ID, response.StatusCode, strings.TrimSpace(string(message)))
	}
	return renderResponse(streams.Out, input, response)
}

func operationPath(op Operation, values map[string]string) (string, error) {
	path := op.Path
	for name, value := range values {
		placeholder := "{" + name + "}"
		if !strings.Contains(path, placeholder) {
			return "", fmt.Errorf("unknown path parameter %q", name)
		}
		if value == "" || value == "." || value == ".." {
			return "", fmt.Errorf("invalid path parameter %q", name)
		}
		path = strings.ReplaceAll(path, placeholder, url.PathEscape(value))
	}
	if strings.Contains(path, "{") {
		return "", fmt.Errorf("missing path parameter for %s", op.Path)
	}
	return path, nil
}

func validateQuery(op Operation, query url.Values) error {
	allowed := map[string]bool{}
	for _, p := range append(append(openapi3.Parameters{}, op.item.Parameters...), op.spec.Parameters...) {
		if p.Value.In == "query" {
			allowed[p.Value.Name] = true
		}
	}
	if value := query.Get("page"); value != "" {
		page, err := strconv.ParseUint(value, 10, 32)
		if err != nil || page == 0 {
			return fmt.Errorf("invalid request: page must be a positive integer")
		}
	}
	for key := range query {
		if !allowed[key] {
			return fmt.Errorf("unknown query parameter %q for %s", key, op.ID)
		}
	}
	return nil
}

func readBody(in io.Reader, input Request, op Operation) ([]byte, string, error) {
	if input.Data != "" && input.File != "" {
		return nil, "", fmt.Errorf("use only one of --data and --file")
	}
	var body []byte
	var err error
	switch {
	case input.File == "-":
		body, err = io.ReadAll(in)
	case input.File != "":
		body, err = os.ReadFile(input.File)
	case input.Data != "":
		body = []byte(input.Data)
	}
	if err != nil {
		return nil, "", err
	}
	if op.spec.RequestBody == nil {
		if len(body) > 0 {
			return nil, "", fmt.Errorf("%s does not accept a request body", op.ID)
		}
		return nil, "", nil
	}
	if len(body) == 0 {
		return nil, "", fmt.Errorf("%s requires a body: use --data or --file (- for stdin); inspect sd api schema %s", op.ID, op.ID)
	}
	ct := input.ContentType
	if ct == "" {
		if op.spec.RequestBody.Value.Content.Get("application/json") != nil {
			ct = "application/json"
		} else {
			ct = "application/octet-stream"
		}
	}
	return body, ct, nil
}

func validateBody(op Operation, body []byte, contentType string) error {
	if op.spec.RequestBody == nil {
		return nil
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	media := op.spec.RequestBody.Value.Content.Get(mediaType)
	if media == nil {
		return fmt.Errorf("unsupported content type %q for %s", mediaType, op.ID)
	}
	if mediaType != "application/json" {
		return nil
	}
	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if media.Schema != nil {
		// Avoid echoing sensitive settings or credentials in validation errors.
		if err := media.Schema.Value.VisitJSON(value, openapi3.VisitAsRequest(), openapi3.SetSchemaErrorMessageCustomizer(func(e *openapi3.SchemaError) string { return e.Reason })); err != nil {
			return fmt.Errorf("invalid body for %s: %w", op.ID, err)
		}
	}
	return nil
}

func (e *Executor) prepareRequest(ctx context.Context, op Operation, in io.Reader, input Request) (*http.Request, error) {
	body, contentType, err := readBody(in, input, op)
	if err != nil {
		return nil, err
	}
	path, err := operationPath(op, input.Path)
	if err != nil {
		return nil, err
	}
	// Validate locally before loading credentials or sending a mutation.
	req, err := http.NewRequestWithContext(ctx, op.Method, "http://localhost"+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = input.Query.Encode()
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}
	if err := validateQuery(op, input.Query); err != nil {
		return nil, err
	}
	validation := &openapi3filter.RequestValidationInput{Request: req, PathParams: input.Path, Route: &routers.Route{Spec: e.spec, Path: op.Path, PathItem: op.item, Method: op.Method, Operation: op.spec}, Options: &openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, ExcludeRequestBody: true}}
	if err := openapi3filter.ValidateRequest(ctx, validation); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if err := validateBody(op, body, contentType); err != nil {
		return nil, err
	}
	if input.Query.Get("live") != "" && input.Output != "raw" && input.OutputFile == "" {
		return nil, fmt.Errorf("live streams require --output raw or --output-file; use --timeout 0s to wait indefinitely")
	}

	return req, nil
}
