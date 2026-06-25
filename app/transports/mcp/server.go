package mcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/middleware/headers"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

func MountMCP(
	lc fx.Lifecycle,
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,

	settings *settings.SettingsRepository,
	registry *robot_tools.Registry,

	mux *http.ServeMux,

	// NOTE: This is duplicated from the OpenAPI router because there's an issue
	// in the OpenAPI codegen that makes mounting it on a sub-router difficult.
	// Eventually, when that's fixed, middleware can be declared once at root.
	ri *headers.Middleware,
	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	if !cfg.MCPEnabled {
		return
	}

	lc.Append(fx.StartHook(func() error {
		set, err := settings.Get(ctx)
		if err != nil {
			return err
		}

		title := set.Title.Or("Storyden")

		mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
			Name:       "storyden-mcp", // TODO: make this configurable?
			Title:      title,
			Version:    config.Version,
			WebsiteURL: cfg.PublicWebAddress.String(),
		}, nil)

		allTools, err := registry.GetTools(ctx)
		if err != nil {
			return err
		}

		for _, t := range allTools {
			if !shouldExportTool(t) {
				continue
			}

			bindTool(mcpServer, t)
		}

		handler := sdkmcp.NewStreamableHTTPHandler(func(r *http.Request) *sdkmcp.Server {
			return mcpServer
		}, nil)

		applied := httpserver.Apply(handler,
			ri.WithHeaderContext(),
			co.WithCORS(),
			lo.WithLogger(),
			cj.WithAuth(),
			rl.WithRequestSizeLimiter(),
			rl.WithRateLimit(),
			withStrictAuthMCP(),
		)

		mux.Handle("/mcp", applied)
		mux.Handle("/mcp/", applied)

		return nil
	}))
}

// bindTool registers a robot tool with the MCP server using the tool's schema
// definition and Handler function.
func bindTool(s *sdkmcp.Server, t *robot_tools.Tool) {
	def := t.Definition
	handler := t.Handler

	s.AddTool(&sdkmcp.Tool{
		Name:         def.Name,
		Title:        def.Title,
		Description:  def.Description,
		InputSchema:  normaliseToolSchema(def.InputSchema),
		OutputSchema: normaliseToolSchema(def.OutputSchema),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           def.Title,
			ReadOnlyHint:    def.Annotations.ReadOnlyHint,
			DestructiveHint: &def.Annotations.DestructiveHint,
			IdempotentHint:  def.Annotations.IdempotentHint,
			OpenWorldHint:   &def.Annotations.OpenWorldHint,
		},
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest) (*sdkmcp.CallToolResult, error) {
		args := req.Params.Arguments
		if args == nil {
			args = json.RawMessage(`{}`)
		}

		result, err := handler(ctx, args)
		if err != nil {
			res := &sdkmcp.CallToolResult{}
			res.SetError(err)
			return res, nil
		}

		// Unmarshal the JSON result into structured content.
		var structured any
		if err := json.Unmarshal(result, &structured); err != nil {
			// If unmarshal fails, return as text (shouldn't happen if Handler is correct).
			return &sdkmcp.CallToolResult{
				Content: []sdkmcp.Content{&sdkmcp.TextContent{Text: string(result)}},
			}, nil
		}

		return &sdkmcp.CallToolResult{
			StructuredContent: structured,
		}, nil
	})
}

func shouldExportTool(t *robot_tools.Tool) bool {
	if t == nil || t.Definition == nil {
		return false
	}

	if t.IsClientTool || t.Handler == nil {
		return false
	}

	// External MCP tools are available to Robots through the tool registry, but
	// Storyden's own MCP server should not re-export third-party MCP servers.
	return !strings.HasPrefix(t.Definition.Name, "mcp:")
}

func normaliseToolSchema(schema any) json.RawMessage {
	if schema == nil {
		return defaultObjectSchema()
	}

	b, err := json.Marshal(schema)
	if err != nil {
		return defaultObjectSchema()
	}

	var object map[string]any
	if err := json.Unmarshal(b, &object); err != nil {
		return defaultObjectSchema()
	}

	if schemaType, ok := object["type"].(string); ok && schemaType == "object" {
		return b
	}

	return defaultObjectSchema()
}

func defaultObjectSchema() json.RawMessage {
	return json.RawMessage(`{"type":"object","additionalProperties":true}`)
}

// withStrictAuthMCP is middleware for MCP-specific authentication checks. MCP
// is fully behind auth so any requests require either a session cookie or an
// access key.
func withStrictAuthMCP() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := session.GetAccountID(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
