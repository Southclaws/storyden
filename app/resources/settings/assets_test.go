package settings_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/mime"
	"github.com/Southclaws/storyden/tests"
)

func TestSettingsAssetSnapshots(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("global branding changes require isolated storage")
	}
	t.Parallel()
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, sr *settings.SettingsRepository, aw *account_writer.Writer, assets *asset_writer.Writer, db *ent.Client, logger *slog.Logger, cfg config.Config) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			ctx, owner := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			source, err := assets.Add(ctx, xid.ID(owner.ID), asset.NewFilename("banner-original"), 100, mime.New("image/png"))
			require.NoError(t, err)
			crop, err := assets.AddVersion(ctx, xid.ID(owner.ID), asset.NewFilename("banner-crop"), 50, mime.New("image/png"), source.ID)
			require.NoError(t, err)

			_, err = sr.Get(ctx)
			require.NoError(t, err)
			_, err = sr.Set(ctx, settings.Settings{BannerAssetID: opt.New(crop.ID)})
			require.NoError(t, err)
			cropped, err := sr.GetWithAssets(ctx)
			require.NoError(t, err)
			require.Equal(t, crop.ID, cropped.BannerAsset.OrZero().ID)
			require.Equal(t, source.ID, cropped.BannerAsset.OrZero().Parent.OrZero().ID)

			restarted, err := settings.New(ctx, fxtest.NewLifecycle(t), logger, db, cfg)
			require.NoError(t, err)
			cold, err := restarted.GetWithAssets(ctx)
			require.NoError(t, err)
			require.Equal(t, crop.ID, cold.BannerAsset.OrZero().ID)
			require.Equal(t, source.ID, cold.BannerAsset.OrZero().Parent.OrZero().ID)

			_, err = sr.Set(ctx, settings.Settings{BannerAssetID: opt.New(source.ID)})
			require.NoError(t, err)
			current, err := sr.GetWithAssets(ctx)
			require.NoError(t, err)
			require.Equal(t, source.ID, current.BannerAsset.OrZero().ID)
			require.False(t, current.BannerAsset.OrZero().Parent.Ok())
			require.Equal(t, crop.ID, cropped.BannerAssetID.OrZero())
			require.Equal(t, crop.ID, cropped.BannerAsset.OrZero().ID)

			var wg sync.WaitGroup
			wg.Go(func() {
				for range 20 {
					for _, id := range []asset.AssetID{crop.ID, source.ID} {
						_, err := sr.Set(ctx, settings.Settings{BannerAssetID: opt.New(id)})
						if !assert.NoError(t, err) {
							return
						}
					}
				}
			})
			for range 2 {
				wg.Go(func() {
					for range 40 {
						snapshot, err := sr.GetWithAssets(ctx)
						if !assert.NoError(t, err) {
							return
						}
						assert.Equal(t, snapshot.BannerAssetID.OrZero(), snapshot.BannerAsset.OrZero().ID)
					}
				})
			}
			wg.Wait()
		}))
	}))
}
