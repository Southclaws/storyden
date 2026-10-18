package asset_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"golang.org/x/image/webp"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestUploadedCropsHaveIndependentVariants(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("isolated file storage is incompatible with shared job workers")
	}

	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		accounts *account_writer.Writer,
		client *openapi.ClientWithResponses,
		sessions *e2e.SessionHelper,
		jobs *asset_image.Repository,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, accounts, seed.Account_001_Odin)
			auth := sessions.WithSession(ctx)
			original := image.NewNRGBA(image.Rect(0, 0, 1800, 900))
			blue := image.NewUniform(color.NRGBA{B: 255, A: 255})
			red := image.NewUniform(color.NRGBA{R: 255, A: 255})
			draw.Draw(original, image.Rect(0, 0, 1800, 450), blue, image.Point{}, draw.Src)
			draw.Draw(original, image.Rect(0, 450, 1800, 900), red, image.Point{}, draw.Src)
			var originalPNG bytes.Buffer
			require.NoError(t, png.Encode(&originalPNG, original))

			upload := func(filename string, data []byte, parent *string) openapi.Asset {
				response, err := client.AssetUploadWithBodyWithResponse(ctx, &openapi.AssetUploadParams{
					Filename:      &filename,
					ContentLength: int64(len(data)),
					ParentAssetId: parent,
				}, "application/octet-stream", bytes.NewReader(data), auth)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, response.StatusCode(), string(response.Body))
				require.NotNil(t, response.JSON200)
				return *response.JSON200
			}
			source := upload("original.png", originalPNG.Bytes(), nil)

			for index, rectangle := range []image.Rectangle{
				image.Rect(0, 0, 1280, 320),
				image.Rect(0, 500, 1280, 820),
			} {
				var cropPNG bytes.Buffer
				require.NoError(t, png.Encode(&cropPNG, original.SubImage(rectangle)))
				crop := upload("cropped-cover.png", cropPNG.Bytes(), &source.Id)
				require.NotNil(t, crop.Parent)
				require.Equal(t, source.Id, crop.Parent.Id)

				for _, input := range []struct {
					asset openapi.Asset
					data  []byte
				}{
					{asset: source, data: originalPNG.Bytes()},
					{asset: crop, data: cropPNG.Bytes()},
				} {
					response, err := client.AssetGetWithResponse(ctx, input.asset.Filename, nil, auth)
					require.NoError(t, err)
					require.Equal(t, input.data, response.Body)
				}

				cropID := openapi.ParseID(crop.Id)
				require.Eventually(t, func() bool {
					job, err := jobs.Get(ctx, cropID)
					return err == nil && job.Status == background_job.StatusComplete
				}, 10*time.Second, 20*time.Millisecond)

				variants, err := jobs.Variants(ctx, cropID)
				require.NoError(t, err)
				require.Len(t, variants, 4)
				expected := map[asset.Variant]image.Point{
					asset.Small:           {320, 80},
					asset.ThumbnailSmall:  {64, 64},
					asset.ThumbnailMedium: {128, 128},
					asset.ThumbnailLarge:  {256, 256},
				}
				for _, variant := range variants {
					row, err := db.Asset.Get(ctx, variant.ID)
					require.NoError(t, err)
					require.Equal(t, cropID, *row.ParentAssetID)
					metadata, err := variant.ImageMetadata()
					require.NoError(t, err)
					require.Equal(t, expected[metadata.Preset], image.Pt(metadata.Width, metadata.Height))
				}

				medium := openapi.M
				params := &openapi.AssetGetParams{Variant: &medium}
				require.Eventually(t, func() bool {
					response, err := client.AssetGetWithResponse(ctx, crop.Filename, params, auth)
					if err != nil {
						return false
					}
					return strings.Contains(response.HTTPResponse.Header.Get("Cache-Control"), "86400")
				}, 10*time.Second, 20*time.Millisecond)
				response, err := client.AssetGetWithResponse(ctx, crop.Filename, params, auth)
				require.NoError(t, err)
				decoded, err := webp.Decode(bytes.NewReader(response.Body))
				require.NoError(t, err)
				require.Equal(t, image.Pt(960, 240), decoded.Bounds().Size())
				pixel := color.NRGBAModel.Convert(decoded.At(0, 0)).(color.NRGBA)
				if index == 0 {
					require.Greater(t, int(pixel.B)-int(pixel.R), 200)
				} else {
					require.Greater(t, int(pixel.R)-int(pixel.B), 200)
				}
			}

			sourceID := openapi.ParseID(source.Id)
			require.Eventually(t, func() bool {
				job, err := jobs.Get(ctx, sourceID)
				return err == nil && job.Status == background_job.StatusComplete
			}, 10*time.Second, 20*time.Millisecond)

			variants, err := jobs.Variants(ctx, sourceID)
			require.NoError(t, err)
			require.Len(t, variants, 4)
			expected := map[asset.Variant]image.Point{
				asset.Small:           {320, 160},
				asset.ThumbnailSmall:  {64, 64},
				asset.ThumbnailMedium: {128, 128},
				asset.ThumbnailLarge:  {256, 256},
			}
			for _, variant := range variants {
				metadata, err := variant.ImageMetadata()
				require.NoError(t, err)
				require.Equal(t, expected[metadata.Preset], image.Pt(metadata.Width, metadata.Height))
			}
		}))
	}))
}
