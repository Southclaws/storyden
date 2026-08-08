package oauth

import (
	"context"
	"database/sql"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/opt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
)

func TestDeviceEndpointsPropagateDatabaseFailures(t *testing.T) {
	db, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared")
	require.NoError(t, err)

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, db)))
	service := &Service{
		cfg:     config.Config{OAuthEnabled: true},
		clients: oauth_querier.New(client),
	}

	require.NoError(t, db.Close())

	ctx := context.Background()
	permissions := rbac.NewList(rbac.PermissionUseOauthClients)

	t.Run("start", func(t *testing.T) {
		result, oauthErr, err := service.StartDeviceAuthorization(ctx, "missing-client", opt.NewEmpty[string]())
		require.Nil(t, result)
		require.Nil(t, oauthErr)
		require.Error(t, err)
	})

	t.Run("read consent", func(t *testing.T) {
		result, oauthErr, err := service.GetDeviceConsent(ctx, account.AccountID{}, permissions, "ABCD-EFGH")
		require.Nil(t, result)
		require.Nil(t, oauthErr)
		require.Error(t, err)
	})

	t.Run("submit consent", func(t *testing.T) {
		oauthErr, err := service.ApproveDeviceAuthorization(ctx, account.AccountID{}, permissions, "ABCD-EFGH", true)
		require.Nil(t, oauthErr)
		require.Error(t, err)
	})
}
