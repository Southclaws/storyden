package asset_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/ent"
	ea "github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestImageCheckpointPublishesAssetAndProgressTogether(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		accounts *account_writer.Writer,
		assets *asset_writer.Writer,
		images *asset_image.Repository,
		jobs *background_job.Repository,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, owner := e2e.WithAccount(ctx, accounts, seed.Account_001_Odin)
			source, err := assets.Add(ctx, xid.ID(owner.ID), asset.NewFilename("source.png"), 100, mime.New("image/png"))
			require.NoError(t, err)
			small, ok := asset.Preset(asset.Small)
			require.True(t, ok)
			medium, ok := asset.Preset(asset.Medium)
			require.True(t, ok)
			err = jobs.Ensure(ctx, source.ID, asset_image.Kind, source.ID.String(),
				asset_image.Payload{Requested: small.Bit})
			require.NoError(t, err)
			claimed, err := jobs.Claim(ctx, source.ID, 3)
			require.NoError(t, err)
			require.NotNil(t, claimed)
			job, err := jobs.Get[asset_image.Payload](ctx, claimed.ID)
			require.NoError(t, err)
			require.NoError(t, images.Request(ctx, source.ID, medium.Bit))

			metadata, err := asset.NewMetadata(mime.New("image/webp"), asset.ImageMetadata{
				Width:  320,
				Height: 160,
				Preset: asset.Small,
				Recipe: asset.ImageRecipe,
			})
			require.NoError(t, err)
			variant := &asset.Asset{
				Name:     asset.NewFilename("s.webp"),
				Size:     42,
				MIME:     mime.New("image/webp"),
				Metadata: metadata,
			}
			err = images.Checkpoint(ctx, job, small, func(tx *ent.Tx) error {
				return asset_writer.AddVariant(ctx, tx.Client(), source, variant)
			})
			require.NoError(t, err)
			stored, err := images.Get(ctx, source.ID)
			require.NoError(t, err)
			require.Equal(t, small.Bit|medium.Bit, stored.Payload.Requested)
			require.Equal(t, small.Bit, stored.Payload.Completed)
			require.Equal(t, background_job.StatusProcessing, stored.Status)
			child, err := db.Asset.Get(ctx, variant.Name.GetID())
			require.NoError(t, err)
			require.Equal(t, source.ID, *child.ParentAssetID)
			require.Equal(t, xid.ID(owner.ID), child.AccountID)
			require.JSONEq(t, string(metadata.RawMessage), string(child.Metadata))

			err = images.Checkpoint(ctx, job, small, func(tx *ent.Tx) error {
				t.Error("completed presets must not publish another asset")
				return nil
			})
			require.Error(t, err)

			err = images.Checkpoint(ctx, job, medium, func(tx *ent.Tx) error {
				return asset_writer.AddVariant(ctx, tx.Client(), source, variant)
			})
			require.True(t, ent.IsConstraintError(err))
			rolledBack, err := images.Get(ctx, source.ID)
			require.NoError(t, err)
			require.Equal(t, stored.Payload, rolledBack.Payload)
			require.Equal(t, stored.Revision, rolledBack.Revision)
			children, err := db.Asset.Query().Where(ea.ParentAssetIDEQ(source.ID)).Count(ctx)
			require.NoError(t, err)
			require.Equal(t, 1, children)
		}))
	}))
}
