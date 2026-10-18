package asset_test

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"path"
	"testing"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	migration "github.com/Southclaws/storyden/app/resources/asset/asset_migration"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	job_resource "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/settings"
	migrator "github.com/Southclaws/storyden/app/services/asset/asset_migration"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/branding/banner"
	"github.com/Southclaws/storyden/app/services/branding/icon"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/mime"
	"github.com/Southclaws/storyden/tests"
)

func TestBootMigratesLegacyImages(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("changes global branding settings")
	}

	t.Parallel()

	var avatarPNG bytes.Buffer
	require.NoError(t, png.Encode(&avatarPNG, image.NewNRGBA(image.Rect(0, 0, 80, 40))))

	var iconPNG bytes.Buffer
	require.NoError(t, png.Encode(&iconPNG, image.NewNRGBA(image.Rect(0, 0, 512, 512))))

	var bannerPNG bytes.Buffer
	require.NoError(t, png.Encode(&bannerPNG, image.NewNRGBA(image.Rect(0, 0, 1200, 630))))

	ownerID := account.AccountID(xid.New())
	started := time.Now().Truncate(time.Second)

	generic := asset.NewFilename("existing.png")
	document := asset.NewFilename("existing.pdf")

	integration.Test(t, &config.Config{AssetStorageLocalPath: t.TempDir()}, e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			accounts *account_writer.Writer,
			accountRepo *account_repo.Repository,
			assets *asset_writer.Writer,
			objects object.Storer,
			db *ent.Client,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				configured, err := settingsRepo.Get(ctx)
				require.NoError(t, err)
				require.False(t, configured.IconAssetID.Ok())
				require.False(t, configured.BannerAssetID.Ok())

				owner, err := accounts.Create(ctx, "legacy-"+ownerID.String(),
					account_writer.WithID(ownerID), account_writer.WithAdmin(true))
				require.NoError(t, err)
				err = db.Account.UpdateOneID(xid.ID(owner.ID)).Modify(func(update *sql.UpdateBuilder) {
					update.Set(ent_account.FieldUpdatedAt, "2026-10-04 12:04:05.814050+00:00")
				}).Exec(ctx)
				require.NoError(t, err)
				_, err = accountRepo.GetByID(ctx, owner.ID)
				require.NoError(t, err)
				require.NoError(t, objects.Write(ctx, path.Join("avatar", owner.ID.String()),
					bytes.NewReader(avatarPNG.Bytes()), int64(avatarPNG.Len())))

				for _, key := range []string{
					"assets/icon-512x512-png",
					"assets/icon-180x180-png",
					"assets/icon-167x167-png",
					"assets/icon-152x152-png",
					"assets/icon-120x120-png",
					"assets/icon-32x32-png",
				} {
					require.NoError(t, objects.Write(ctx, key,
						bytes.NewReader(iconPNG.Bytes()), int64(iconPNG.Len())))
				}

				require.NoError(t, objects.Write(ctx, "assets/banner-1200x630-png",
					bytes.NewReader(bannerPNG.Bytes()), int64(bannerPNG.Len())))

				require.NoError(t, objects.Write(ctx, asset.BuildAssetPath(generic),
					bytes.NewReader(avatarPNG.Bytes()), int64(avatarPNG.Len())))
				_, err = assets.Add(ctx, xid.ID(owner.ID), generic, avatarPNG.Len(), mime.New("image/png"))
				require.NoError(t, err)

				pdf := []byte("%PDF-1.7\nexisting document")
				require.NoError(t, objects.Write(ctx, asset.BuildAssetPath(document), bytes.NewReader(pdf), int64(len(pdf))))
				_, err = assets.Add(ctx, xid.ID(owner.ID), document, len(pdf), mime.New("application/pdf"))
				require.NoError(t, err)
			}))
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			_ *migrator.Migrator,
			avatars avatar.Service,
			icons icon.Service,
			banners banner.Service,
			settings *settings.SettingsRepository,
			images *asset_image.Repository,
			objects object.Storer,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				keys := []string{"avatar/" + ownerID.String(), "icon", "banner"}
				require.EventuallyWithT(t, func(collect *assert.CollectT) {
					for _, key := range keys {
						row, err := db.BackgroundJob.Query().
							Where(backgroundjob.KindEQ(migration.Kind), backgroundjob.KeyEQ(key)).Only(ctx)
						if !assert.NoError(collect, err) {
							continue
						}

						assert.Equal(collect, backgroundjob.StatusComplete, row.Status)
						assert.False(collect, row.ID.Time().Before(started))
						assert.False(collect, row.ID.Time().After(time.Now()))
					}
				}, 10*time.Second, 20*time.Millisecond)

				metadata, reader, pending, err := avatars.GetVariant(ctx, ownerID, asset.Original)
				require.NoError(t, err)
				require.False(t, pending)
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.NoError(t, reader.(io.Closer).Close())
				require.Equal(t, avatarPNG.Bytes(), data)
				require.NotEmpty(t, metadata.ID)

				configured, err := settings.Get(ctx)
				require.NoError(t, err)
				require.True(t, configured.IconAssetID.Ok())
				require.True(t, configured.BannerAssetID.Ok())
				for _, sourceID := range []asset.AssetID{configured.IconAssetID.OrZero(), configured.BannerAssetID.OrZero()} {
					original, err := db.Asset.Get(ctx, sourceID)
					require.NoError(t, err)
					reader, _, err := objects.Read(ctx, asset.BuildAssetPath(asset.NewFilepathFilename(original.Filename)))
					require.NoError(t, err)
					data, err := io.ReadAll(reader)
					require.NoError(t, err)
					require.NoError(t, reader.(io.Closer).Close())
					expected := iconPNG.Bytes()
					if sourceID == configured.BannerAssetID.OrZero() {
						expected = bannerPNG.Bytes()
					}
					require.Equal(t, expected, data)
				}

				require.Eventually(t, func() bool {
					for _, id := range []asset.AssetID{configured.IconAssetID.OrZero(), configured.BannerAssetID.OrZero()} {
						job, err := images.Get(ctx, id)
						if err != nil || job.Status != job_resource.StatusComplete {
							return false
						}
					}

					return true
				}, 10*time.Second, 20*time.Millisecond)

				_, reader, err = icons.Get(ctx, "32x32")
				require.NoError(t, err)
				decoded, err := png.Decode(reader)
				require.NoError(t, err)
				require.NoError(t, reader.(io.Closer).Close())
				require.Equal(t, image.Pt(32, 32), decoded.Bounds().Size())
				_, reader, err = banners.Get(ctx)
				require.NoError(t, err)
				decoded, err = png.Decode(reader)
				require.NoError(t, err)
				require.NoError(t, reader.(io.Closer).Close())
				require.Equal(t, image.Pt(1200, 630), decoded.Bounds().Size())

				paths := []string{"avatar/" + ownerID.String()}
				for _, name := range []string{
					"icon-512x512-png", "icon-180x180-png", "icon-167x167-png",
					"icon-152x152-png", "icon-120x120-png", "icon-32x32-png", "banner-1200x630-png",
				} {
					paths = append(paths, "assets/"+name)
				}
				for _, key := range paths {
					exists, err := objects.Exists(ctx, key)
					require.NoError(t, err)
					require.False(t, exists, key)
				}
				_, err = images.Get(ctx, generic.GetID())
				require.True(t, ent.IsNotFound(err))

				reader, _, err = objects.Read(ctx, asset.BuildAssetPath(generic))
				require.NoError(t, err)
				data, err = io.ReadAll(reader)
				require.NoError(t, err)
				require.NoError(t, reader.(io.Closer).Close())
				require.Equal(t, avatarPNG.Bytes(), data)
				_, err = images.Get(ctx, document.GetID())
				require.True(t, ent.IsNotFound(err))
			}))
		}),
	)
}
