package model_media_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/model_media"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestResolverAcceptsOnlyExistingStorydenImageAssets(t *testing.T) {
	t.Parallel()

	publicAPI, err := url.Parse("https://storyden.example/base")
	require.NoError(t, err)
	integration.Test(t, &config.Config{PublicAPIAddress: *publicAPI}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		resolver *model_media.Resolver,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("image-resolver-owner").SetName("Image Resolver Owner").Save(ctx)
			require.NoError(t, err)

			pngID := asset.AssetID(xid.New())
			pngName := asset.NewExistingFilename(pngID, "reference image.png").String()
			_, err = db.Asset.Create().SetID(pngID).SetAccountID(owner.ID).SetFilename(pngName).SetMimeType("image/png").SetSize(10).Save(ctx)
			require.NoError(t, err)

			image, err := resolver.ResolveImage(ctx, pngID)
			require.NoError(t, err)
			assert.Equal(t, "image/png", image.MIME)
			assert.Equal(t, "https://storyden.example/base/api/assets/"+pngName, image.URL)

			svgID := asset.AssetID(xid.New())
			_, err = db.Asset.Create().SetID(svgID).SetAccountID(owner.ID).SetFilename(asset.NewExistingFilename(svgID, "vector.svg").String()).SetMimeType("image/svg+xml").SetSize(10).Save(ctx)
			require.NoError(t, err)
			err = resolver.ValidateContents(ctx, []*genai.Content{{Role: genai.RoleUser, Parts: []*genai.Part{robotresource.NewImageAssetPart(svgID)}}})
			require.Error(t, err)
			assert.EqualError(t, err, "Robot image attachments must be GIF, JPEG, PNG, or WebP assets.")

			missingID := asset.AssetID(xid.New())
			err = resolver.ValidateContents(ctx, []*genai.Content{{Role: genai.RoleUser, Parts: []*genai.Part{robotresource.NewImageAssetPart(missingID)}}})
			require.Error(t, err)
			assert.ErrorContains(t, err, "is not available")
		}))
	}))
}
