package mcpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/app/resources/robot/mcp"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremotetoken"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
	storydenmcp "github.com/Southclaws/storyden/lib/mcp"
	"github.com/google/jsonschema-go/jsonschema"
	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	toolIDPrefix   = "mcp:"
	refreshTimeout = 30 * time.Second
	callTimeout    = 60 * time.Second
	probeTimeout   = 10 * time.Second
	maxCardBytes   = 256 * 1024
)

type Manager struct {
	logger   *slog.Logger
	repo     *mcp.Repository
	registry *tools.Registry
	settings *settings.SettingsRepository
	tokens   *oauthremotetoken.Service
	config   config.Config
	client   *http.Client
}

func New(
	lc fx.Lifecycle,
	ctx context.Context,
	logger *slog.Logger,
	repo *mcp.Repository,
	registry *tools.Registry,
	settings *settings.SettingsRepository,
	tokens *oauthremotetoken.Service,
	cfg config.Config,
) *Manager {
	m := Manager{
		logger:   logger,
		repo:     repo,
		registry: registry,
		settings: settings,
		tokens:   tokens,
		config:   cfg,
		client:   httpsafe.NewClient(httpsafe.Config{DialTimeout: probeTimeout}),
	}

	lc.Append(fx.StartHook(func() error {
		return m.Start(ctx)
	}))

	return &m
}

func (m *Manager) Start(ctx context.Context) error {
	if err := m.SyncRegistry(ctx); err != nil {
		return err
	}

	servers, err := m.repo.ListEnabledServers(ctx)
	if err != nil {
		return err
	}

	for _, server := range servers {
		go func(server mcp.Server) {
			refreshCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), refreshTimeout)
			defer cancel()
			if _, err := m.Refresh(refreshCtx, server.ID); err != nil {
				m.logger.Warn("failed to refresh robot MCP server", slog.String("server", server.Slug), slog.String("error", err.Error()))
			}
		}(server)
	}

	return nil
}

func (m *Manager) ListServers(ctx context.Context) ([]mcp.Server, error) {
	return m.repo.ListServers(ctx)
}

func (m *Manager) GetServer(ctx context.Context, id mcp.ServerID) (mcp.Server, error) {
	return m.repo.GetServer(ctx, id)
}

func (m *Manager) CreateServer(ctx context.Context, in mcp.ServerCreate) (mcp.Server, error) {
	if strings.TrimSpace(in.Slug) == "" {
		in.Slug = Slugify(in.Name)
	}
	if in.Slug == "" {
		return mcp.Server{}, fault.New("MCP server slug is required")
	}
	server, err := m.repo.CreateServer(ctx, in)
	if err != nil {
		return mcp.Server{}, err
	}
	if server.Enabled {
		refreshed, err := m.Refresh(ctx, server.ID)
		if err != nil {
			return server, err
		}
		return refreshed, nil
	}
	return server, nil
}

func (m *Manager) UpdateServer(ctx context.Context, id mcp.ServerID, in mcp.ServerUpdate) (mcp.Server, error) {
	server, err := m.repo.UpdateServer(ctx, id, in)
	if err != nil {
		return mcp.Server{}, err
	}
	if !server.Enabled {
		if err := m.SyncRegistry(ctx); err != nil {
			return mcp.Server{}, err
		}
		return server, nil
	}
	if _, err := m.Refresh(ctx, server.ID); err != nil {
		return server, err
	}
	return m.repo.GetServer(ctx, id)
}

func (m *Manager) DeleteServer(ctx context.Context, id mcp.ServerID) error {
	if err := m.repo.DeleteServer(ctx, id); err != nil {
		return err
	}
	return m.SyncRegistry(ctx)
}

func (m *Manager) Refresh(ctx context.Context, id mcp.ServerID) (mcp.Server, error) {
	server, err := m.repo.GetServer(ctx, id)
	if err != nil {
		return mcp.Server{}, err
	}
	if !server.Enabled {
		if err := m.SyncRegistry(ctx); err != nil {
			return mcp.Server{}, err
		}
		return server, nil
	}

	discovered, err := m.discoverTools(ctx, server)
	if err != nil {
		_ = m.repo.MarkRefreshError(context.WithoutCancel(ctx), server.ID, err.Error())
		return server, err
	}

	now := time.Now()
	if err := m.repo.UpsertTools(ctx, server, discovered); err != nil {
		return mcp.Server{}, err
	}
	if err := m.repo.MarkRefreshSuccess(ctx, server.ID, now); err != nil {
		return mcp.Server{}, err
	}
	if err := m.SyncRegistry(ctx); err != nil {
		return mcp.Server{}, err
	}

	return m.repo.GetServer(ctx, id)
}

func (m *Manager) RefreshOAuthConnectionServers(ctx context.Context, id oauth_remote.ConnectionID) error {
	servers, err := m.repo.ListServersByOAuthRemoteConnection(ctx, xid.ID(id))
	if err != nil {
		return err
	}

	for _, server := range servers {
		if !server.Enabled {
			enabled := true
			server, err = m.repo.UpdateServer(ctx, server.ID, mcp.ServerUpdate{Enabled: &enabled})
			if err != nil {
				return err
			}
		}
		if _, err := m.Refresh(ctx, server.ID); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) SyncRegistry(ctx context.Context) error {
	m.registry.UnregisterPrefix(toolIDPrefix)

	servers, err := m.repo.ListEnabledServers(ctx)
	if err != nil {
		return err
	}
	for _, server := range servers {
		if !server.Enabled {
			continue
		}
		for _, cachedTool := range server.Tools {
			if !cachedTool.Enabled {
				continue
			}
			m.registry.Register(m.makeTool(server, cachedTool))
		}
	}

	return nil
}

func (m *Manager) ListToolCatalogue(ctx context.Context) ([]tools.CatalogueTool, error) {
	catalogue := m.registry.ListCatalogue(ctx)
	registered := map[string]struct{}{}
	for _, tool := range catalogue {
		registered[tool.ID] = struct{}{}
	}

	cached, err := m.repo.ListTools(ctx)
	if err != nil {
		return nil, err
	}
	for _, tool := range cached {
		if _, ok := registered[tool.ID]; ok {
			continue
		}
		catalogue = append(catalogue, tools.CatalogueTool{
			ID:                   tool.ID,
			CallableName:         tool.CallableName,
			Name:                 tool.Title,
			Description:          tool.Description,
			Source:               "mcp",
			Available:            false,
			RequiresConfirmation: false,
		})
	}

	return catalogue, nil
}

type ServerCard struct {
	Name        string             `json:"name"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Title       string             `json:"title,omitempty"`
	WebsiteURL  string             `json:"websiteUrl,omitempty"`
	Remotes     []ServerCardRemote `json:"remotes,omitempty"`
}

type ServerCardRemote struct {
	Type                      string   `json:"type"`
	URL                       string   `json:"url"`
	SupportedProtocolVersions []string `json:"supportedProtocolVersions,omitempty"`
}

type ProbeResult struct {
	InputURL      string
	EndpointURL   string
	ServerCardURL string
	ServerCard    *ServerCard
	RemoteType    string
	Active        bool
	ProbeError    string
}

func (m *Manager) Probe(ctx context.Context, rawURL string, bearerToken string) (ProbeResult, error) {
	probeCtx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	input, err := normaliseHTTPURL(rawURL)
	if err != nil {
		return ProbeResult{}, err
	}

	result := ProbeResult{
		InputURL:    input.String(),
		EndpointURL: input.String(),
	}

	if card, cardURL, err := fetchServerCard(probeCtx, m.client, input); err == nil {
		result.ServerCard = &card
		result.ServerCardURL = cardURL
		if input.Path == "" || input.Path == "/" {
			if remote, ok := selectStreamableHTTPRemote(card); ok {
				result.EndpointURL = remote.URL
				result.RemoteType = remote.Type
			}
		}
	}

	session, err := m.connect(probeCtx, mcp.Server{
		EndpointURL:    result.EndpointURL,
		BearerToken:    bearerToken,
		Description:    "",
		Enabled:        true,
		HasBearerToken: bearerToken != "",
	})
	if err != nil {
		if input.Path == "" || input.Path == "/" {
			mcpURL := *input
			mcpURL.Path = "/mcp"
			if mcpURL.String() != result.EndpointURL {
				session, mcpErr := m.connect(probeCtx, mcp.Server{
					EndpointURL:    mcpURL.String(),
					BearerToken:    bearerToken,
					Description:    "",
					Enabled:        true,
					HasBearerToken: bearerToken != "",
				})
				if mcpErr == nil {
					defer session.Close()
					result.EndpointURL = mcpURL.String()
					result.Active = true
					return result, nil
				}
			}
		}
		result.ProbeError = err.Error()
		return result, nil
	}
	defer session.Close()

	result.Active = true
	return result, nil
}

func (m *Manager) discoverTools(ctx context.Context, server mcp.Server) ([]mcp.Tool, error) {
	session, err := m.connect(ctx, server)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var out []mcp.Tool
	for remoteTool, err := range session.Tools(ctx, &sdkmcp.ListToolsParams{}) {
		if err != nil {
			return nil, err
		}
		toolID := ToolID(server.Slug, remoteTool.Name)
		out = append(out, mcp.Tool{
			ID:           toolID,
			RemoteName:   remoteTool.Name,
			CallableName: CallableName(server.Name, remoteTool.Name),
			Title:        remoteTool.Title,
			Description:  remoteTool.Description,
			InputSchema:  schemaToMap(remoteTool.InputSchema),
			OutputSchema: schemaToMap(remoteTool.OutputSchema),
			Annotations:  annotationsToMap(remoteTool.Annotations),
			Enabled:      true,
			ServerID:     server.ID,
			ServerSlug:   server.Slug,
		})
	}

	return out, nil
}

func (m *Manager) makeTool(server mcp.Server, cachedTool mcp.Tool) *tools.Tool {
	def := &storydenmcp.ToolDefinition{
		Name:         cachedTool.ID,
		Title:        firstNonEmpty(cachedTool.Title, cachedTool.RemoteName),
		Description:  cachedTool.Description,
		InputSchema:  mapToSchema(cachedTool.InputSchema),
		OutputSchema: mapToSchema(cachedTool.OutputSchema),
		Annotations: storydenmcp.ToolAnnotations{
			ReadOnlyHint:    boolAnnotation(cachedTool.Annotations, "readOnlyHint"),
			DestructiveHint: boolAnnotation(cachedTool.Annotations, "destructiveHint"),
			IdempotentHint:  boolAnnotation(cachedTool.Annotations, "idempotentHint"),
			OpenWorldHint:   boolAnnotation(cachedTool.Annotations, "openWorldHint"),
		},
	}

	return &tools.Tool{
		Definition:   def,
		Source:       "mcp",
		CallableName: cachedTool.CallableName,
		Handler: func(ctx context.Context, raw json.RawMessage) (json.RawMessage, error) {
			var args map[string]any
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &args); err != nil {
					return nil, err
				}
			}
			out := m.callTool(ctx, server, cachedTool.RemoteName, args)
			return json.Marshal(out)
		},
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        cachedTool.CallableName,
					Description: cachedTool.Description,
					InputSchema: def.InputSchema,
				},
				func(ctx tool.Context, args map[string]any) (map[string]any, error) {
					return m.callTool(ctx, server, cachedTool.RemoteName, args), nil
				},
			)
		},
	}
}

func (m *Manager) callTool(ctx context.Context, server mcp.Server, remoteName string, args map[string]any) map[string]any {
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	session, err := m.connect(callCtx, server)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	defer session.Close()

	res, err := session.CallTool(callCtx, &sdkmcp.CallToolParams{
		Name:      remoteName,
		Arguments: args,
	})
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	if res.IsError {
		return map[string]any{"error": callToolErrorText(res)}
	}
	if res.StructuredContent != nil {
		return map[string]any{"output": res.StructuredContent}
	}
	text := callToolText(res)
	if text == "" {
		return map[string]any{"output": nil}
	}
	return map[string]any{"output": text}
}

func (m *Manager) connect(ctx context.Context, server mcp.Server) (*sdkmcp.ClientSession, error) {
	set, err := m.settings.Get(ctx)
	if err != nil {
		return nil, err
	}

	client := sdkmcp.NewClient(&sdkmcp.Implementation{
		Name:       "storyden",
		Title:      set.Title.Or(settings.DefaultTitle),
		Version:    config.Version,
		WebsiteURL: m.config.PublicWebAddress.String(),
		Icons: []sdkmcp.Icon{
			{
				Source:   m.config.PublicAPIAddress.ResolveReference(&url.URL{Path: "/api/info/icon/512x512"}).String(),
				MIMEType: "image/png",
				Sizes:    []string{"512x512"},
				Theme:    sdkmcp.IconThemeLight,
			},
			// TODO: Validate and add more icon sizes.
		},
	}, nil)

	bearerToken := server.BearerToken
	if bearerToken == "" && server.OAuthRemoteConnectionID != nil {
		token, err := m.tokens.AccessToken(ctx, *server.OAuthRemoteConnectionID)
		if err != nil {
			return nil, err
		}
		bearerToken = token
	}

	var transport sdkmcp.Transport

	// try with streamable (modern default)
	transport = &sdkmcp.StreamableClientTransport{
		Endpoint:             server.EndpointURL,
		DisableStandaloneSSE: true,
		OAuthHandler:         staticBearer(bearerToken),
		HTTPClient:           m.client,
	}
	s, err1 := client.Connect(ctx, transport, nil)
	if err1 == nil {
		return s, nil
	}

	// fall back to sse transport for older servers
	transport = &sdkmcp.SSEClientTransport{
		Endpoint:   server.EndpointURL,
		HTTPClient: m.client,
	}
	s, err2 := client.Connect(ctx, transport, nil)
	if err2 == nil {
		return s, nil
	}

	return nil, fmt.Errorf("failed to connect to MCP server with available transports: %v, %v", err1, err2)
}

type staticBearerHandler struct {
	token string
}

func staticBearer(token string) sdkauth.OAuthHandler {
	if token == "" {
		return nil
	}
	return &staticBearerHandler{token: token}
}

func (h *staticBearerHandler) TokenSource(context.Context) (oauth2.TokenSource, error) {
	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: h.token, TokenType: "Bearer"}), nil
}

func (h *staticBearerHandler) Authorize(context.Context, *http.Request, *http.Response) error {
	return errors.New("bearer token rejected by MCP server")
}

func ToolID(serverSlug, remoteName string) string {
	return toolIDPrefix + serverSlug + ":" + remoteName
}

func CallableName(providerName, remoteName string) string {
	provider := providerKey(providerName)
	tool := sanitizeIdentifier(strings.ToLower(remoteName))
	for _, part := range strings.Split(provider, "_") {
		if part == "" || part == "mcp" {
			continue
		}
		if strings.HasPrefix(tool, part+"_") {
			tool = strings.TrimPrefix(tool, part+"_")
			break
		}
	}
	if tool == "" {
		tool = "tool"
	}
	return provider + "_" + tool
}

var nonIdentifierChar = regexp.MustCompile(`[^A-Za-z0-9_]+`)
var bracketedText = regexp.MustCompile(`\s*\([^)]*\)`)

func sanitizeIdentifier(in string) string {
	out := strings.Trim(nonIdentifierChar.ReplaceAllString(in, "_"), "_")
	if out == "" {
		return "tool"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "x_" + out
	}
	return out
}

func providerKey(name string) string {
	name = bracketedText.ReplaceAllString(name, "")
	identifier := sanitizeIdentifier(strings.ToLower(name))
	parts := strings.Split(identifier, "_")
	parts = trimProviderNoise(parts)
	if len(parts) == 0 {
		return identifier
	}
	return strings.Join(parts, "_")
}

func trimProviderNoise(parts []string) []string {
	for len(parts) > 1 && parts[0] == "mcp" {
		parts = parts[1:]
	}
	for len(parts) > 0 {
		last := parts[len(parts)-1]
		switch last {
		case "com", "net", "org", "io", "so", "app", "dev", "cloud":
			parts = parts[:len(parts)-1]
		default:
			return parts
		}
	}
	return parts
}

var nonSlugChar = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(in string) string {
	slug := strings.Trim(nonSlugChar.ReplaceAllString(strings.ToLower(in), "-"), "-")
	if slug == "" {
		return ""
	}
	return slug
}

func schemaToMap(schema any) map[string]any {
	if schema == nil {
		return map[string]any{}
	}
	var out map[string]any
	body, err := json.Marshal(schema)
	if err == nil {
		_ = json.Unmarshal(body, &out)
	}
	if out == nil {
		return map[string]any{}
	}
	return out
}

func mapToSchema(in map[string]any) *jsonschema.Schema {
	if len(in) == 0 {
		return nil
	}
	body, err := json.Marshal(in)
	if err != nil {
		return nil
	}
	var out jsonschema.Schema
	if err := json.Unmarshal(body, &out); err != nil {
		return nil
	}
	return &out
}

func annotationsToMap(in *sdkmcp.ToolAnnotations) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	out := map[string]any{}
	if in.DestructiveHint != nil {
		out["destructiveHint"] = *in.DestructiveHint
	}
	out["idempotentHint"] = in.IdempotentHint
	if in.OpenWorldHint != nil {
		out["openWorldHint"] = *in.OpenWorldHint
	}
	out["readOnlyHint"] = in.ReadOnlyHint
	return out
}

func boolAnnotation(in map[string]any, key string) bool {
	v, _ := in[key].(bool)
	return v
}

func callToolText(res *sdkmcp.CallToolResult) string {
	var b strings.Builder
	for _, content := range res.Content {
		text, ok := content.(*sdkmcp.TextContent)
		if !ok {
			continue
		}
		b.WriteString(text.Text)
	}
	return b.String()
}

func callToolErrorText(res *sdkmcp.CallToolResult) string {
	if text := callToolText(res); text != "" {
		return text
	}
	return "Tool execution failed."
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func normaliseHTTPURL(raw string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, fault.New("MCP endpoint URL must use http or https")
	}
	if u.Host == "" {
		return nil, fault.New("MCP endpoint URL host is required")
	}
	if u.User != nil || u.Fragment != "" {
		return nil, fault.New("MCP endpoint URL must not include user info or fragment")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	if u.Path == "" {
		u.Path = "/"
	}
	return u, nil
}

func serverCardURL(input *url.URL) string {
	u := *input
	u.Path = "/.well-known/mcp-server-card"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func fetchServerCard(ctx context.Context, client *http.Client, input *url.URL) (ServerCard, string, error) {
	cardURL := serverCardURL(input)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cardURL, nil)
	if err != nil {
		return ServerCard{}, cardURL, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return ServerCard{}, cardURL, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return ServerCard{}, cardURL, fmt.Errorf("MCP server card discovery failed: %s", res.Status)
	}

	var card ServerCard
	if err := json.NewDecoder(io.LimitReader(res.Body, maxCardBytes)).Decode(&card); err != nil {
		return ServerCard{}, cardURL, err
	}

	return card, cardURL, nil
}

func selectStreamableHTTPRemote(card ServerCard) (ServerCardRemote, bool) {
	for _, remote := range card.Remotes {
		if remote.URL == "" {
			continue
		}
		if remote.Type == "streamable-http" || remote.Type == "streamable_http" {
			return remote, true
		}
	}
	return ServerCardRemote{}, false
}
