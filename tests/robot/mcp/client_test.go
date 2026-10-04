package mcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	robot_mcp "github.com/Southclaws/storyden/app/resources/robot/mcp"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremote"
	"github.com/Southclaws/storyden/app/services/semdex/robot/mcpclient"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_flow "github.com/Southclaws/storyden/internal/ent/oauthremoteauthorisationflow"
	ent_connection "github.com/Southclaws/storyden/internal/ent/oauthremoteconnection"
	ent_robot_mcp_server "github.com/Southclaws/storyden/internal/ent/robotmcpserver"
	ent_robot_mcp_tool "github.com/Southclaws/storyden/internal/ent/robotmcptool"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestRobotMCPServerCreateDiscoversBearerProtectedTools(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		withLoopbackMCPClient(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				endpoint, addRemoteTool := newMutableBearerMCPServer(t, "test-token")
				enabled := true
				token := "test-token"
				slug := "test-server-" + xid.New().String()

				created := tests.AssertRequest(cl.RobotMCPServerCreateWithResponse(root,
					openapi.RobotMCPServerCreateJSONRequestBody{
						Name:        "Test MCP Server",
						Slug:        &slug,
						Description: strPtr("Bearer protected test MCP server"),
						EndpointUrl: endpoint,
						Enabled:     &enabled,
						BearerToken: &token,
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, created.JSON200)

				assert.Equal(t, slug, created.JSON200.Slug)
				assert.Equal(t, mcpclient.ToolsetID(slug), created.JSON200.ToolsetId)
				assert.True(t, created.JSON200.HasBearerToken)
				assert.NotContains(t, string(created.Body), token)
				require.Len(t, created.JSON200.Tools, 1)
				assert.Equal(t, "mcp:"+slug+":echo", created.JSON200.Tools[0].Id)
				assert.Equal(t, mcpclient.CallableName("Test MCP Server", "echo"), created.JSON200.Tools[0].CallableName)
				assert.True(t, created.JSON200.Tools[0].Available)

				catalogue := tests.AssertRequest(cl.RobotToolsListWithResponse(root, adminSession))(t, http.StatusOK)
				require.NotNil(t, catalogue.JSON200)

				var found *openapi.RobotToolInfo
				for i := range catalogue.JSON200.Tools {
					tool := &catalogue.JSON200.Tools[i]
					if tool.Id == "mcp:"+slug+":echo" {
						found = tool
						break
					}
				}

				require.NotNil(t, found)
				assert.True(t, found.Available)
				assert.Equal(t, "mcp", string(found.Source))
				assert.Equal(t, mcpclient.CallableName("Test MCP Server", "echo"), found.CallableName)

				toolsetCatalogue := tests.AssertRequest(cl.RobotToolsetsListWithResponse(root, adminSession))(t, http.StatusOK)
				require.NotNil(t, toolsetCatalogue.JSON200)
				var serverToolset *openapi.RobotToolset
				for i := range toolsetCatalogue.JSON200.Toolsets {
					candidate := &toolsetCatalogue.JSON200.Toolsets[i]
					if candidate.Id == created.JSON200.ToolsetId {
						serverToolset = candidate
						break
					}
				}
				require.NotNil(t, serverToolset)
				assert.Equal(t, openapi.RobotToolsetSourceMcp, serverToolset.Source)
				require.NotNil(t, serverToolset.SourceRef)
				assert.Equal(t, string(created.JSON200.Id), *serverToolset.SourceRef)
				assert.False(t, serverToolset.Editable)
				assert.Equal(t, []string{"mcp:" + slug + ":echo"}, serverToolset.Tools)

				addRemoteTool("search", "Search events", "Search calendar events by title.")
				refreshed, err := cl.RobotMCPServerRefreshWithResponse(root, created.JSON200.Id, adminSession)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, refreshed.StatusCode(), "%s", refreshed.Body)
				require.NotNil(t, refreshed.JSON200)
				require.Len(t, refreshed.JSON200.Tools, 2)
				assert.Equal(t, "mcp:"+slug+":echo", refreshed.JSON200.Tools[0].Id)
				assert.Equal(t, "mcp:"+slug+":search", refreshed.JSON200.Tools[1].Id)

				t.Run("concurrent_refresh", func(t *testing.T) {
					const requests = 16
					responses := make([]*openapi.RobotMCPServerRefreshResponse, requests)
					errors := make([]error, requests)
					start := make(chan struct{})
					var wg sync.WaitGroup
					for i := range requests {
						wg.Go(func() {
							<-start
							responses[i], errors[i] = cl.RobotMCPServerRefreshWithResponse(root, created.JSON200.Id, adminSession)
						})
					}
					close(start)
					wg.Wait()

					for i, response := range responses {
						require.NoError(t, errors[i])
						require.Equal(t, http.StatusOK, response.StatusCode(), "%s", response.Body)
						require.NotNil(t, response.JSON200)
						require.Len(t, response.JSON200.Tools, 2)
					}
				})

				refreshedToolsets := tests.AssertRequest(cl.RobotToolsetsListWithResponse(root, adminSession))(t, http.StatusOK)
				require.NotNil(t, refreshedToolsets.JSON200)
				for _, candidate := range refreshedToolsets.JSON200.Toolsets {
					if candidate.Id == created.JSON200.ToolsetId {
						assert.Equal(t, []string{"mcp:" + slug + ":echo", "mcp:" + slug + ":search"}, candidate.Tools)
						return
					}
				}
				require.Fail(t, "refreshed MCP Toolset was not listed")
			}))
		}),
	)
}

func TestRobotMCPConcurrentRefreshAcrossInstances(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	var first *openapi.ClientWithResponses
	var second *openapi.ClientWithResponses
	var sh *e2e.SessionHelper
	var aw *account_writer.Writer
	var db *ent.Client
	integration.Test(t, nil, e2e.Setup(), withLoopbackMCPClient(), fx.Populate(&cfg, &first, &sh, &aw, &db))
	integration.Test(t, &cfg, e2e.Setup(), withLoopbackMCPClient(), fx.Populate(&second))

	ctx := context.Background()
	adminCtx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
	adminSession := sh.WithSession(adminCtx)
	endpoint, addRemoteTool := newMutableBearerMCPServer(t, "test-token")
	slug := "replicated-server-" + xid.New().String()
	created := tests.AssertRequest(first.RobotMCPServerCreateWithResponse(ctx, openapi.RobotMCPServerCreateJSONRequestBody{
		Name:        "Replicated MCP Server",
		Slug:        &slug,
		EndpointUrl: endpoint,
		Enabled:     new(true),
		BearerToken: new("test-token"),
	}, adminSession))(t, http.StatusOK)
	require.NotNil(t, created.JSON200)

	echoID := "mcp:" + slug + ":echo"
	echo, err := db.RobotMCPTool.Query().Where(ent_robot_mcp_tool.ToolID(echoID)).Only(ctx)
	require.NoError(t, err)
	addRemoteTool("echo", "Updated echo", "Echoes the supplied message.")
	addRemoteTool("search", "Search events", "Search calendar events by title.")

	const requests = 16
	responses := make([]*openapi.RobotMCPServerRefreshResponse, requests)
	errors := make([]error, requests)
	clients := []*openapi.ClientWithResponses{first, second}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range requests {
		wg.Go(func() {
			<-start
			responses[i], errors[i] = clients[i%len(clients)].RobotMCPServerRefreshWithResponse(ctx, created.JSON200.Id, adminSession)
		})
	}
	close(start)
	wg.Wait()

	for i, response := range responses {
		require.NoError(t, errors[i])
		require.Equal(t, http.StatusOK, response.StatusCode(), "%s", response.Body)
		require.NotNil(t, response.JSON200)
		require.Len(t, response.JSON200.Tools, 2)
	}

	updated, err := db.RobotMCPTool.Query().Where(ent_robot_mcp_tool.ToolID(echoID)).Only(ctx)
	require.NoError(t, err)
	require.Equal(t, echo.ID, updated.ID)
	require.True(t, echo.CreatedAt.Equal(updated.CreatedAt))
	require.Equal(t, "Updated echo", updated.Title)
	search, err := db.RobotMCPTool.Query().Where(ent_robot_mcp_tool.ToolID("mcp:" + slug + ":search")).Only(ctx)
	require.NoError(t, err)
	require.True(t, search.Enabled)
	require.Equal(t, "Search events", search.Title)

	t.Run("failed_refresh_preserves_cached_tools", func(t *testing.T) {
		addRemoteTool("collision-tool", "First collision", "First colliding callable name.")
		addRemoteTool("collision_tool", "Second collision", "Second colliding callable name.")
		tests.AssertRequest(first.RobotMCPServerRefreshWithResponse(ctx, created.JSON200.Id, adminSession))(t, http.StatusBadRequest)

		exists, err := db.RobotMCPTool.Query().Where(
			ent_robot_mcp_tool.ServerID(echo.ServerID),
			ent_robot_mcp_tool.RemoteNameIn("collision-tool", "collision_tool"),
		).Exist(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		for _, id := range []xid.ID{echo.ID, search.ID} {
			stored, err := db.RobotMCPTool.Get(ctx, id)
			require.NoError(t, err)
			require.True(t, stored.Enabled)
		}
	})
}

func TestRobotMCPRefreshRecordsBearerFailure(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		withLoopbackMCPClient(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			aw *account_writer.Writer,
			repo *robot_mcp.Repository,
			manager *mcpclient.Manager,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				endpoint := newBearerMCPServer(t, "test-token")

				server, err := repo.CreateServer(adminCtx, robot_mcp.ServerCreate{
					Name:        "Broken MCP Server",
					Slug:        "broken-server-" + xid.New().String(),
					EndpointURL: endpoint,
					Enabled:     true,
					BearerToken: "wrong-token",
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)

				_, err = manager.Refresh(root, server.ID)
				require.Error(t, err)

				stored, err := repo.GetServer(root, server.ID)
				require.NoError(t, err)
				require.NotNil(t, stored.LastError)
				assert.NotEmpty(t, strings.TrimSpace(*stored.LastError))
			}))
		}),
	)
}

func TestRobotMCPProbeFallsBackToCommonMCPPath(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		withLoopbackMCPClient(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				endpoint := newPathMCPServer(t, "/mcp")

				probe := tests.AssertRequest(cl.RobotMCPServerProbeWithResponse(root,
					openapi.RobotMCPServerProbeJSONRequestBody{
						Url: endpoint,
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, probe.JSON200)

				assert.True(t, probe.JSON200.Active)
				assert.Equal(t, endpoint+"/mcp", probe.JSON200.EndpointUrl)
			}))
		}),
	)
}

func TestRobotMCPServerDeleteRemovesLinkedOAuthConnection(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			db *ent.Client,
			aw *account_writer.Writer,
			oauth *oauthremote.Service,
			repo *robot_mcp.Repository,
			manager *mcpclient.Manager,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				connectionInput := oauthremote.CreateConnectionInput{
					ResourceURL: "https://mcp-delete-test.example/mcp",
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "delete-test-client",
						AuthorizationEndpoint: "https://auth-delete-test.example/authorize",
						TokenEndpoint:         "https://auth-delete-test.example/token",
						AuthorizationServer:   "https://auth-delete-test.example",
					},
					AddedBy: admin.ID,
				}

				connection, err := oauth.CreateConnection(adminCtx, connectionInput)
				require.NoError(t, err)
				_, err = oauth.StartAuthorization(adminCtx, connection.ID)
				require.NoError(t, err)

				server, err := repo.CreateServer(adminCtx, robot_mcp.ServerCreate{
					Name:                    "OAuth Delete Test",
					Slug:                    "oauth-delete-test-" + xid.New().String(),
					EndpointURL:             "https://mcp-delete-test.example/mcp",
					Enabled:                 false,
					OAuthRemoteConnectionID: &connection.ID,
					AddedBy:                 admin.ID,
				})
				require.NoError(t, err)

				err = repo.UpsertTools(adminCtx, server, []robot_mcp.Tool{{
					ID:           mcpclient.ToolID(server.Slug, "echo"),
					RemoteName:   "echo",
					CallableName: mcpclient.CallableName(server.Name, "echo"),
					Title:        "Echo",
					Description:  "Echoes input.",
					Enabled:      true,
					ServerID:     server.ID,
					ServerSlug:   server.Slug,
				}})
				require.NoError(t, err)

				err = manager.DeleteServer(adminCtx, server.ID)
				require.NoError(t, err)

				count, err := db.OAuthRemoteConnection.Query().Where(ent_connection.IDEQ(xid.ID(connection.ID))).Count(adminCtx)
				assertZeroCount(t, count, err)
				count, err = db.OAuthRemoteAuthorisationFlow.Query().Where(ent_flow.ConnectionIDEQ(xid.ID(connection.ID))).Count(adminCtx)
				assertZeroCount(t, count, err)
				count, err = db.RobotMCPServer.Query().Where(ent_robot_mcp_server.IDEQ(xid.ID(server.ID))).Count(adminCtx)
				assertZeroCount(t, count, err)
				count, err = db.RobotMCPTool.Query().Where(ent_robot_mcp_tool.ServerIDEQ(xid.ID(server.ID))).Count(adminCtx)
				assertZeroCount(t, count, err)

				_, err = oauth.CreateConnection(adminCtx, connectionInput)
				require.NoError(t, err)
			}))
		}),
	)
}

func TestRobotMCPOAuthRefreshesTokenForRefreshAndToolCall(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		withLoopbackMCPClient(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			aw *account_writer.Writer,
			oauth *oauthremote.Service,
			oauthRepo *oauth_remote.Repository,
			repo *robot_mcp.Repository,
			manager *mcpclient.Manager,
			registry *robot_tools.Registry,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				tokenEndpoint, refreshCount := newRefreshTokenEndpoint(t, "refreshed-mcp-access-token")
				mcpEndpoint := newBearerMCPServer(t, "refreshed-mcp-access-token")

				connection, err := oauth.CreateConnection(adminCtx, oauthremote.CreateConnectionInput{
					ResourceURL: mcpEndpoint,
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "mcp-refresh-client",
						AuthorizationEndpoint: tokenEndpoint + "/authorize",
						TokenEndpoint:         tokenEndpoint + "/token",
						AuthorizationServer:   tokenEndpoint,
					},
					AddedBy: admin.ID,
				})
				require.NoError(t, err)

				expired := time.Now().Add(-time.Hour)
				connection, err = oauthRepo.StoreTokens(adminCtx, connection.ID, oauth_remote.TokenUpdate{
					AccessToken:  "expired-mcp-access-token",
					RefreshToken: "mcp-refresh-token",
					TokenType:    "Bearer",
					TokenExpiry:  &expired,
				})
				require.NoError(t, err)

				server, err := repo.CreateServer(adminCtx, robot_mcp.ServerCreate{
					Name:                    "OAuth Refresh MCP",
					Slug:                    "oauth-refresh-mcp-" + xid.New().String(),
					EndpointURL:             mcpEndpoint,
					Enabled:                 true,
					OAuthRemoteConnectionID: &connection.ID,
					AddedBy:                 admin.ID,
				})
				require.NoError(t, err)

				refreshed, err := manager.Refresh(adminCtx, server.ID)
				require.NoError(t, err)
				require.Len(t, refreshed.Tools, 1)
				assert.Equal(t, 1, refreshCount())

				tool, err := registry.GetTool(adminCtx, mcpclient.ToolID(server.Slug, "echo"))
				require.NoError(t, err)
				require.NotNil(t, tool.Handler)

				raw, err := tool.Handler(adminCtx, json.RawMessage(`{"message":"hello"}`))
				require.NoError(t, err)
				assert.NotContains(t, string(raw), "error")
				assert.Contains(t, string(raw), "hello")
				assert.Equal(t, 1, refreshCount(), "fresh stored token should be reused for the tool call")
			}))
		}),
	)
}

func newBearerMCPServer(t *testing.T, token string) string {
	t.Helper()
	endpoint, _ := newMutableBearerMCPServer(t, token)
	return endpoint
}

func newMutableBearerMCPServer(t *testing.T, token string) (string, func(name, title, description string)) {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{Name: "test-mcp", Version: "v1"}, nil)
	addTool := func(name, title, description string) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        name,
			Title:       title,
			Description: description,
		}, func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, map[string]any, error) {
			return nil, map[string]any{"message": args["message"]}, nil
		})
	}
	addTool("echo", "Echo", "Echoes a message.")

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(httpServer.Close)

	return httpServer.URL, addTool
}

func withLoopbackMCPClient() fx.Option {
	return fx.Replace(mcpclient.HTTPClient{Client: http.DefaultClient})
}

func newRefreshTokenEndpoint(t *testing.T, accessToken string) (string, func() int) {
	t.Helper()

	var mu sync.Mutex
	refreshes := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			http.NotFound(w, r)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.PostForm.Get("grant_type") != "refresh_token" || r.PostForm.Get("refresh_token") != "mcp-refresh-token" {
			http.Error(w, "invalid refresh request", http.StatusBadRequest)
			return
		}

		mu.Lock()
		refreshes++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"access_token":  accessToken,
			"refresh_token": "mcp-refresh-token-rotated",
			"token_type":    "Bearer",
			"expires_in":    3600,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)

	return server.URL, func() int {
		mu.Lock()
		defer mu.Unlock()
		return refreshes
	}
}

func newPathMCPServer(t *testing.T, path string) string {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{Name: "test-mcp", Version: "v1"}, nil)
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(httpServer.Close)

	return httpServer.URL
}

func strPtr(s string) *string {
	return &s
}

func assertZeroCount(t *testing.T, count int, err error) {
	t.Helper()

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
