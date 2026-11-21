package moderation_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestSuspensionToolsProtectAdministrators(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		aq *account_querier.Querier,
		registry *robot_tools.Registry,
	) {
		lc.Append(fx.StartHook(func() {
			_, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			_, otherAdmin := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			_, moderator := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			_, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)

			moderatorCtx := session.WithAccountPermissions(root, *moderator, rbac.NewList(rbac.PermissionManageSuspensions))
			adminCtx := session.WithAccountPermissions(root, *admin, rbac.NewList(rbac.PermissionAdministrator))

			suspend, err := registry.GetTool(root, "member_suspend")
			require.NoError(t, err)
			reinstate, err := registry.GetTool(root, "member_reinstate")
			require.NoError(t, err)

			call := func(ctx context.Context, tool *robot_tools.Tool, target account.AccountID) error {
				args, err := json.Marshal(map[string]string{"account_id": target.String()})
				require.NoError(t, err)
				_, err = tool.Handler(ctx, args)
				return err
			}

			suspended := func(id account.AccountID) bool {
				acc, err := aq.GetByID(root, id)
				require.NoError(t, err)
				return acc.DeletedAt.Ok()
			}

			t.Run("moderator_cannot_suspend_administrator", func(t *testing.T) {
				require.Error(t, call(moderatorCtx, suspend, otherAdmin.ID))
				require.False(t, suspended(otherAdmin.ID))
			})

			t.Run("moderator_can_suspend_and_reinstate_member", func(t *testing.T) {
				require.NoError(t, call(moderatorCtx, suspend, member.ID))
				require.True(t, suspended(member.ID))
				require.NoError(t, call(moderatorCtx, reinstate, member.ID))
				require.False(t, suspended(member.ID))
			})

			t.Run("administrator_can_suspend_administrator", func(t *testing.T) {
				require.NoError(t, call(adminCtx, suspend, otherAdmin.ID))
				require.True(t, suspended(otherAdmin.ID))
			})

			t.Run("moderator_cannot_reinstate_administrator", func(t *testing.T) {
				require.Error(t, call(moderatorCtx, reinstate, otherAdmin.ID))
				require.True(t, suspended(otherAdmin.ID))
			})

			t.Run("administrator_can_reinstate_administrator", func(t *testing.T) {
				require.NoError(t, call(adminCtx, reinstate, otherAdmin.ID))
				require.False(t, suspended(otherAdmin.ID))
			})
		}))
	}))
}
