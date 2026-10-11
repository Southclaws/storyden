package tools_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/ent"
	entnode "github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestLibraryPageListRejectsSQLInjectionInputs(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		registry *robottools.Registry,
		aw *account_writer.Writer,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, member := e2e.WithAccount(ctx, aw, seed.Account_004_Loki)
			ctx = session.WithAccountPermissions(ctx, *member, rbac.NewList(rbac.PermissionReadPublishedLibrary))
			list, err := registry.GetTool(ctx, "library_page_list")
			require.NoError(t, err)
			canary, err := db.Node.Create().
				SetAccountID(xid.ID(member.ID)).SetName("SQL injection canary").
				SetSlug(xid.New().String()).SetVisibility(entnode.VisibilityPublished).
				SetCreatedAt(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)).Save(ctx)
			require.NoError(t, err)

			inputs := []string{
				`{"created_since":"2040-01-01T00:00:00Z' OR 1=1 --"}`,
				`{"created_before":"2040-01-01T00:00:00Z'; DROP TABLE nodes; --"}`,
				`{"sort":"newest; DROP TABLE nodes; --"}`,
				`{"sort":"created_at DESC, (SELECT pg_sleep(0.01))"}`,
				`{"visibility":"published' OR 1=1 --"}`,
				`{"page":"1 UNION SELECT NULL"}`,
				`{"page_size":"1; DROP TABLE nodes; --"}`,
				`{"depth":"0) OR 1=1 --"}`,
			}
			for _, input := range inputs {
				t.Run(input, func(t *testing.T) {
					_, err := list.Handler(ctx, json.RawMessage(input))
					require.Error(t, err)
				})
			}

			saved, err := db.Node.Get(ctx, canary.ID)
			require.NoError(t, err)
			require.Equal(t, "SQL injection canary", saved.Name)
			require.Equal(t, entnode.VisibilityPublished, saved.Visibility)
		}))
	}))
}
