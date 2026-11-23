package moderation_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	mcptransport "github.com/Southclaws/storyden/app/transports/mcp"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestSuspensionToolsEnforceCurrentPermissions(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{MCPEnabled: true}, e2e.Setup(), mcptransport.Build(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		server *httptest.Server,
		sessions *e2e.SessionHelper,
		aw *account_writer.Writer,
		aq *account_querier.Querier,
		roles *role_repo.Repository,
		assignments *role_assign.Assignment,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			moderatorCtx, moderator := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)

			client, _ := connectMCP(t, root, server.URL, sessions.WithSession(moderatorCtx))
			adminClient, adminCookies := connectMCP(t, root, server.URL, sessions.WithSession(adminCtx))

			moderatorRole, err := roles.Create(root, "suspensions-"+xid.New().String(), "blue", rbac.PermissionList{rbac.PermissionManageSuspensions})
			require.NoError(t, err)
			administratorRole, err := roles.Create(root, "administrator-"+xid.New().String(), "red", rbac.PermissionList{rbac.PermissionAdministrator})
			require.NoError(t, err)

			call := func(t *testing.T, client *mcp.ClientSession, tool string, target account.AccountID, allowed, suspended bool) {
				t.Helper()
				result, err := client.CallTool(root, &mcp.CallToolParams{
					Name:      tool,
					Arguments: map[string]any{"account_id": target.String()},
				})
				require.NoError(t, err)
				content, err := json.Marshal(result.Content)
				require.NoError(t, err)
				require.Equal(t, !allowed, result.IsError, "%s", content)
				if !allowed {
					require.Contains(t, result.Content[0].(*mcp.TextContent).Text, "invalid permissions")
				}

				saved, err := aq.GetRefByID(root, target)
				require.NoError(t, err)
				require.Equal(t, suspended, saved.IsSuspended())
			}

			t.Run("member_cannot_suspend_or_reinstate", func(t *testing.T) {
				call(t, client, "member_suspend", member.ID, false, false)
				call(t, client, "member_reinstate", member.ID, false, false)
			})

			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(moderator.ID), role_assign.Add(moderatorRole.ID)))
			t.Run("granted_role_applies_to_existing_session", func(t *testing.T) {
				call(t, client, "member_suspend", member.ID, true, true)
			})

			client, _ = connectMCP(t, root, server.URL, sessions.WithSession(moderatorCtx))

			_, err = roles.Update(root, moderatorRole.ID, role_repo.WithPermissions(rbac.PermissionList{rbac.PermissionManageReports}))
			require.NoError(t, err)
			t.Run("removed_permission_blocks_existing_session", func(t *testing.T) {
				call(t, client, "member_reinstate", member.ID, false, true)
			})

			_, err = roles.Update(root, moderatorRole.ID, role_repo.WithPermissions(rbac.PermissionList{rbac.PermissionManageSuspensions}))
			require.NoError(t, err)
			t.Run("restored_permission_allows_reinstatement", func(t *testing.T) {
				call(t, client, "member_reinstate", member.ID, true, false)
			})

			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(moderator.ID), role_assign.Remove(moderatorRole.ID)))
			t.Run("removed_role_blocks_existing_session", func(t *testing.T) {
				call(t, client, "member_suspend", member.ID, false, false)
			})

			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(moderator.ID), role_assign.Add(moderatorRole.ID)))
			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(member.ID), role_assign.Add(administratorRole.ID)))
			t.Run("cached_member_promoted_to_administrator_is_protected", func(t *testing.T) {
				call(t, client, "member_suspend", member.ID, false, false)
				call(t, adminClient, "member_suspend", member.ID, true, true)
				call(t, client, "member_reinstate", member.ID, false, true)
				call(t, adminClient, "member_reinstate", member.ID, true, false)
			})

			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(member.ID), role_assign.Remove(administratorRole.ID)))
			t.Run("former_administrator_can_be_moderated", func(t *testing.T) {
				call(t, client, "member_suspend", member.ID, true, true)
				call(t, client, "member_reinstate", member.ID, true, false)
			})

			t.Run("connection_cannot_retain_previous_callers_privileges", func(t *testing.T) {
				request, err := http.NewRequestWithContext(root, http.MethodPost, server.URL+"/mcp", nil)
				require.NoError(t, err)
				require.NoError(t, sessions.WithSession(memberCtx)(root, request))
				adminCookies.SetCookies(request.URL, request.Cookies())

				call(t, adminClient, "member_suspend", moderator.ID, false, false)
			})
		}))
	}))
}

func connectMCP(t *testing.T, ctx context.Context, serverURL string, edit openapi.RequestEditorFn) (*mcp.ClientSession, *cookiejar.Jar) {
	t.Helper()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+"/mcp", nil)
	require.NoError(t, err)
	require.NoError(t, edit(ctx, request))
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	jar.SetCookies(request.URL, request.Cookies())

	client := mcp.NewClient(&mcp.Implementation{Name: "rbac-test", Version: "v1"}, nil)
	connection, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:             request.URL.String(),
		HTTPClient:           &http.Client{Jar: jar},
		DisableStandaloneSSE: true,
	}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = connection.Close() })
	return connection, jar
}
