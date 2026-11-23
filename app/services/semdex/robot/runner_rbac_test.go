package robot

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestRobotToolRBAC(t *testing.T) {
	member := account.Account{ID: account.AccountID(xid.New())}
	suspended := member
	suspended.DeletedAt = opt.New(time.Now())

	cases := []struct {
		name     string
		ctx      context.Context
		required opt.Optional[rbac.Permission]
		allowed  bool
	}{
		{
			name:     "missing_permission",
			ctx:      session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionManageReports)),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  false,
		},
		{
			name:     "permitted_member",
			ctx:      session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionManageSuspensions)),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  true,
		},
		{
			name:     "administrator_override",
			ctx:      session.WithAccountPermissions(t.Context(), member, rbac.NewList(rbac.PermissionAdministrator)),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  true,
		},
		{
			name:     "restricted_oauth",
			ctx:      session.WithOAuthToken(t.Context(), member, rbac.NewList(rbac.PermissionReadPublishedThreads), []string{"read"}),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  false,
		},
		{
			name:     "suspended_member",
			ctx:      session.WithAccountPermissions(t.Context(), suspended, rbac.NewList(rbac.PermissionManageSuspensions)),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  false,
		},
		{
			name:     "suspended_administrator",
			ctx:      session.WithAccountPermissions(t.Context(), suspended, rbac.NewList(rbac.PermissionAdministrator)),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  false,
		},
		{
			name:     "suspended_on_open_tool",
			ctx:      session.WithAccountPermissions(t.Context(), suspended, rbac.NewList()),
			required: opt.NewEmpty[rbac.Permission](),
			allowed:  false,
		},
		{
			name:     "missing_session",
			ctx:      t.Context(),
			required: opt.New(rbac.PermissionManageSuspensions),
			allowed:  false,
		},
		{
			name:     "open_tool",
			ctx:      session.WithAccountPermissions(t.Context(), member, rbac.NewList()),
			required: opt.NewEmpty[rbac.Permission](),
			allowed:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			registry := tools.NewRegistry(slog.Default())
			registered := &tools.Tool{
				Definition: &mcp.ToolDefinition{
					Name:               "test_tool",
					Title:              "Test tool",
					RequiredPermission: tc.required,
				},
			}
			require.NoError(t, registry.Register(registered))
			runner := &Agent{toolProvider: registry}
			ctx := &rbacToolContext{sessionContext: tc.ctx}

			_, err := runner.checkToolRBAC(ctx, registered, nil)
			if tc.allowed {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

type rbacToolContext struct {
	agent.Context
	sessionContext context.Context
}

func (c *rbacToolContext) Value(key any) any {
	return c.sessionContext.Value(key)
}
