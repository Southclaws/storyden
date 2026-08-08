package oauth_writer

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"
)

func TestCreateDeviceAuthorisationRejectsDuplicateUserCode(t *testing.T) {
	sqlDB, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	db := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, sqlDB))))
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	ctx := context.Background()
	writer := New(db)
	client, err := writer.CreateClient(ctx, ClientCreate{
		AccountID:        opt.NewEmpty[account.AccountID](),
		ClientID:         "device-code-collision-test",
		Name:             "Device code collision test",
		Type:             oauth.ClientTypePublic,
		RedirectURIs:     []string{},
		AllowedScopes:    []string{},
		AllowedGrants:    []string{},
		ClientSecretHash: opt.NewEmpty[string](),
	})
	require.NoError(t, err)

	input := DeviceAuthorisationCreate{
		ClientID:            client.ID,
		DeviceCodeHash:      "first-device-code",
		UserCodeHash:        "same-user-code",
		UserCodeDisplay:     "0123-4567",
		ExpiresAt:           time.Now().Add(time.Minute),
		PollIntervalSeconds: 5,
	}
	_, err = writer.CreateDeviceAuthorisation(ctx, input)
	require.NoError(t, err)

	input.DeviceCodeHash = "second-device-code"
	_, err = writer.CreateDeviceAuthorisation(ctx, input)
	require.Error(t, err)
	assert.Equal(t, ftag.AlreadyExists, ftag.Get(err))
}
