package asset_test

import (
	"bytes"
	"context"
	"io"
	"path"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	migration "github.com/Southclaws/storyden/app/resources/asset/asset_migration"
	job_resource "github.com/Southclaws/storyden/app/resources/background_job"
	migrator "github.com/Southclaws/storyden/app/services/asset/asset_migration"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/tests"
)

func TestLegacyAvatarImportPreservesNewerAvatar(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("isolated file storage is incompatible with shared job workers")
	}

	t.Parallel()

	ownerID := account.AccountID(xid.New())
	key := path.Join("avatar", ownerID.String())
	data := imagePayload(t)
	jobID := xid.New()

	integration.Test(t, &config.Config{AssetStorageLocalPath: t.TempDir()},
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			accounts *account_writer.Writer,
			avatars avatar.Service,
			jobs *job_resource.Repository,
			objects object.Storer,
		) {
			lc.Append(fx.StartHook(func() {
				owner, err := accounts.Create(ctx, "legacy-"+ownerID.String(), account_writer.WithID(ownerID))
				require.NoError(t, err)

				require.NoError(t, objects.Write(ctx, key, bytes.NewReader(data), int64(len(data))))
				require.NoError(t, avatars.Set(ctx, ownerID, bytes.NewReader(data), int64(len(data))))

				intent := migration.Payload{
					Target:           migration.Avatar,
					Path:             key,
					AccountID:        ownerID,
					AccountUpdatedAt: owner.UpdatedAt,
				}

				err = jobs.Ensure(ctx, jobID, migration.Kind, key, intent)
				require.NoError(t, err)
			}))
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			_ *migrator.Migrator,
			avatars avatar.Service,
			jobs *job_resource.Repository,
			objects object.Storer,
		) {
			lc.Append(fx.StartHook(func() {
				require.Eventually(t, func() bool {
					job, err := jobs.Get[migration.Payload](ctx, jobID)
					return err == nil && job.Status == job_resource.StatusSkipped
				}, 5*time.Second, 20*time.Millisecond)

				job, err := jobs.Get[migration.Payload](ctx, jobID)
				require.NoError(t, err)
				require.Equal(t, "image reference changed during import", job.Reason)
				require.Zero(t, job.Payload.SourceID)

				_, reader, _, err := avatars.GetVariant(ctx, ownerID, asset.Original)
				require.NoError(t, err)

				stored, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.NoError(t, reader.(io.Closer).Close())
				require.Equal(t, data, stored)

				exists, err := objects.Exists(ctx, key)
				require.NoError(t, err)
				require.True(t, exists)
			}))
		}),
	)
}
