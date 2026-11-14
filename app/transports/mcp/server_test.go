package mcp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	storydenmcp "github.com/Southclaws/storyden/lib/mcp"
)

func TestBindToolDefaultsMissingSchemas(t *testing.T) {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "test",
		Version: "v1",
	}, nil)

	tool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{
			Name:        "test_tool",
			Title:       "Test tool",
			Description: "Tool used for schema default testing.",
		},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{"ok":true}`), nil
		},
	}

	require.NotPanics(t, func() {
		bindTool(server, tool)
	})
}

func TestBindToolSetsMCPAudience(t *testing.T) {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "test",
		Version: "v1",
	}, nil)

	audience := make(chan robot_tools.ToolAudience, 1)
	tool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{
			Name:        "test_tool",
			Title:       "Test tool",
			Description: "Tool used for audience testing.",
		},
		Handler: func(ctx context.Context, _ json.RawMessage) (json.RawMessage, error) {
			audience <- robot_tools.ToolAudienceFromContext(ctx)
			return json.RawMessage(`{"ok":true}`), nil
		},
	}
	bindTool(server, tool)

	clientTransport, serverTransport := sdkmcp.NewInMemoryTransports()
	serverSession, err := server.Connect(t.Context(), serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, serverSession.Close()) })

	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "test-client", Version: "v1"}, nil)
	clientSession, err := client.Connect(t.Context(), clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientSession.Close()) })

	_, err = clientSession.CallTool(t.Context(), &sdkmcp.CallToolParams{
		Name:      "test_tool",
		Arguments: map[string]any{},
	})
	require.NoError(t, err)
	assert.Equal(t, robot_tools.ToolAudienceMCP, <-audience)
}

func TestNormaliseToolSchemaRequiresObjectSchema(t *testing.T) {
	assert.JSONEq(t,
		`{"type":"object","additionalProperties":true}`,
		string(normaliseToolSchema(map[string]any{"type": "string"})),
	)

	assert.JSONEq(t,
		`{"type":"object","properties":{"query":{"type":"string"}}}`,
		string(normaliseToolSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
		})),
	)
}

func TestShouldExportToolSkipsExternalMCPTools(t *testing.T) {
	nativeTool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{Name: "library_search"},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}

	externalTool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{Name: "mcp:notion-mcp-beta:notion-create-comment"},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}
	internalTool := &robot_tools.Tool{
		Definition: storydenmcp.GetDocumentGetTool(),
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}
	webFetchTool := &robot_tools.Tool{
		Definition: storydenmcp.GetWebFetchTool(),
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}

	assert.True(t, shouldExportTool(nativeTool))
	assert.True(t, shouldExportTool(webFetchTool))
	assert.False(t, shouldExportTool(externalTool))
	assert.False(t, shouldExportTool(internalTool))
}

func TestBindToolEnforcesRequiredPermission(t *testing.T) {
	member := account.Account{ID: account.AccountID(xid.New()), Handle: "member"}
	suspended := account.Account{ID: account.AccountID(xid.New()), Handle: "suspended", DeletedAt: opt.New(time.Now())}

	cases := map[string]struct {
		ctx     context.Context
		allowed bool
	}{
		"member_without_permission":        {session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionCreatePost)), false},
		"member_with_permission":           {session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionManageSuspensions)), true},
		"administrator":                    {session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionAdministrator)), true},
		"suspended_member_with_permission": {session.WithAccountPermissions(t.Context(), suspended, rbac.NewList(rbac.PermissionManageSuspensions)), false},
		"no_session":                       {t.Context(), false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			called := false
			tool := &robot_tools.Tool{
				Definition: &storydenmcp.ToolDefinition{
					Name:               "privileged_tool",
					RequiredPermission: opt.New(rbac.PermissionManageSuspensions),
				},
				Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
					called = true
					return json.RawMessage(`{"ok":true}`), nil
				},
			}

			res := callBoundTool(t, tc.ctx, tool)
			assert.Equal(t, tc.allowed, called)
			assert.Equal(t, !tc.allowed, res.IsError)
		})
	}
}

func TestBindToolRejectsSuspendedAccountWithoutRequiredPermission(t *testing.T) {
	suspended := account.Account{ID: account.AccountID(xid.New()), Handle: "suspended", DeletedAt: opt.New(time.Now())}
	called := false
	tool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{Name: "open_tool"},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			called = true
			return json.RawMessage(`{"ok":true}`), nil
		},
	}

	res := callBoundTool(t, session.WithAccountPermissions(t.Context(), suspended, rbac.NewList(rbac.PermissionCreatePost)), tool)
	assert.False(t, called)
	assert.True(t, res.IsError)
}

func callBoundTool(t *testing.T, ctx context.Context, tool *robot_tools.Tool) *sdkmcp.CallToolResult {
	t.Helper()

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "v1"}, nil)
	bindTool(server, tool)

	clientTransport, serverTransport := sdkmcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = serverSession.Close() })

	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "test-client", Version: "v1"}, nil)
	clientSession, err := client.Connect(t.Context(), clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = clientSession.Close() })

	res, err := clientSession.CallTool(t.Context(), &sdkmcp.CallToolParams{Name: tool.Definition.Name, Arguments: map[string]any{}})
	require.NoError(t, err)
	return res
}
