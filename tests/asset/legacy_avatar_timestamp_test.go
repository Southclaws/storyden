package asset_test

import (
	"context"
	"testing"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestLegacyAvatarTimestampGuard(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name       string
		updatedAt  string
		replaced   bool
		deleted    bool
		wantImport bool
	}{
		{name: "same instant", updatedAt: "2026-10-04 12:04:05.814050+00:00", wantImport: true},
		{name: "changed timestamp", updatedAt: "2026-10-04 12:04:05.814051+00:00"},
		{name: "replacement", updatedAt: "2026-10-04 12:04:05.814050+00:00", replaced: true},
		{name: "deleted account", updatedAt: "2026-10-04 12:04:05.814050+00:00", deleted: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			integration.Test(t, nil, fx.Invoke(func(
				lc fx.Lifecycle,
				ctx context.Context,
				db *ent.Client,
				accounts *account_writer.Writer,
			) {
				lc.Append(fx.StartHook(func() {
					owner, err := accounts.Create(ctx, "legacy-"+xid.New().String())
					require.NoError(t, err)
					ownerID := xid.ID(owner.ID)
					replacementID := xid.New()

					update := db.Account.UpdateOneID(ownerID).Modify(func(update *sql.UpdateBuilder) {
						update.Set(account.FieldUpdatedAt, test.updatedAt)
					})
					if test.replaced {
						update.SetAvatarAssetID(replacementID)
					}
					if test.deleted {
						update.SetDeletedAt(time.Now())
					}
					require.NoError(t, update.Exec(ctx))

					snapshot, err := time.Parse(time.RFC3339Nano, "2026-10-04T12:04:05.81405Z")
					require.NoError(t, err)
					assetID := xid.New()
					tx, err := db.Tx(ctx)
					require.NoError(t, err)
					defer tx.Rollback()

					updated, err := account_repo.SetAvatarAssetIfUnchanged(ctx, tx, ownerID, snapshot, assetID)
					require.NoError(t, err)
					require.Equal(t, test.wantImport, updated)
					require.NoError(t, tx.Commit())

					stored, err := db.Account.Get(ctx, ownerID)
					require.NoError(t, err)
					if test.wantImport {
						require.Equal(t, assetID, *stored.AvatarAssetID)
					} else if test.replaced {
						require.Equal(t, replacementID, *stored.AvatarAssetID)
					} else {
						require.Nil(t, stored.AvatarAssetID)
					}
					if !test.wantImport {
						expected, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", test.updatedAt)
						require.NoError(t, err)
						require.True(t, stored.UpdatedAt.Equal(expected))
					}
				}))
			}))
		})
	}
}

func TestLegacyAvatarImportSerializesReplacement(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		accounts *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := accounts.Create(ctx, "legacy-"+xid.New().String())
			require.NoError(t, err)
			ownerID := xid.ID(owner.ID)
			err = db.Account.UpdateOneID(ownerID).Modify(func(update *sql.UpdateBuilder) {
				update.Set(account.FieldUpdatedAt, "2026-10-04 12:04:05.814050+00:00")
			}).Exec(ctx)
			require.NoError(t, err)
			snapshot, err := time.Parse(time.RFC3339Nano, "2026-10-04T12:04:05.81405Z")
			require.NoError(t, err)

			tx, err := db.Tx(ctx)
			require.NoError(t, err)
			defer tx.Rollback()
			updated, err := account_repo.SetAvatarAssetIfUnchanged(ctx, tx, ownerID, snapshot, xid.New())
			require.NoError(t, err)
			require.True(t, updated)

			replacementID := xid.New()
			started := make(chan struct{})
			done := make(chan error, 1)
			go func() {
				close(started)
				done <- db.Account.UpdateOneID(ownerID).SetAvatarAssetID(replacementID).Exec(ctx)
			}()
			<-started
			require.NoError(t, tx.Commit())
			select {
			case err := <-done:
				require.NoError(t, err)
			case <-time.After(5 * time.Second):
				t.Fatal("avatar replacement did not finish")
			}

			stored, err := db.Account.Get(ctx, ownerID)
			require.NoError(t, err)
			require.Equal(t, replacementID, *stored.AvatarAssetID)
		}))
	}))
}
