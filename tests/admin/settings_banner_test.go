package admin_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"testing"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestBannerAssetSettings(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("global branding changes require isolated storage")
	}
	t.Parallel()
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			auth := sh.WithSession(ctx)
			makePNG := func(width, height int, fill color.NRGBA) []byte {
				image := image.NewNRGBA(image.Rect(0, 0, width, height))
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						image.SetNRGBA(x, y, fill)
					}
				}
				var b bytes.Buffer
				require.NoError(t, png.Encode(&b, image))
				return b.Bytes()
			}
			sourcePNG := makePNG(100, 100, color.NRGBA{R: 255, A: 255})
			croppedPNG := makePNG(50, 25, color.NRGBA{B: 255, A: 255})
			upload := func(name string, data []byte, parent *string) openapi.Asset {
				r := tests.AssertRequest(cl.AssetUploadWithBodyWithResponse(ctx, &openapi.AssetUploadParams{Filename: &name, ContentLength: int64(len(data)), ParentAssetId: parent}, "application/octet-stream", bytes.NewReader(data), auth))(t, http.StatusOK)
				return *r.JSON200
			}
			source := upload("banner-source.png", sourcePNG, nil)
			child := upload("banner-crop.png", croppedPNG, &source.Id)
			metadata := openapi.Metadata{
				"branding": map[string]any{
					"banner": map[string]any{
						"asset_id": child.Id,
						"crop":     map[string]any{"left": .25, "top": .25, "width": .5, "height": .25},
					},
				},
			}
			patch := openapi.AdminSettingsUpdateJSONRequestBody{BannerAssetId: &child.Id, Metadata: &metadata}
			updated := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, patch, auth))(t, http.StatusOK)
			require.Equal(t, child.Id, updated.JSON200.BannerAsset.Id)
			require.Equal(t, metadata, *updated.JSON200.Metadata)
			require.Equal(t, source.Id, updated.JSON200.BannerAsset.Parent.Id)
			original := openapi.Original
			banner := tests.AssertRequest(cl.BannerGetWithResponse(ctx, &openapi.BannerGetParams{Variant: &original}, auth))(t, http.StatusOK)
			require.Equal(t, croppedPNG, banner.Body)

			title := "Banner settings retained"
			tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, openapi.AdminSettingsUpdateJSONRequestBody{Title: &title}, auth))(t, http.StatusOK)
			current := tests.AssertRequest(cl.AdminSettingsGetWithResponse(ctx, auth))(t, http.StatusOK)
			require.Equal(t, metadata, *current.JSON200.Metadata)
			require.Equal(t, child.Id, current.JSON200.BannerAsset.Id)
			require.Equal(t, source.Id, current.JSON200.BannerAsset.Parent.Id)

			text := upload("not-an-image.txt", []byte("not a banner"), nil)
			tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, openapi.AdminSettingsUpdateJSONRequestBody{BannerAssetId: &text.Id}, auth))(t, http.StatusBadRequest)

			updated = tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, openapi.AdminSettingsUpdateJSONRequestBody{BannerAssetId: &source.Id}, auth))(t, http.StatusOK)
			require.Equal(t, source.Id, updated.JSON200.BannerAsset.Id)
			require.Nil(t, updated.JSON200.BannerAsset.Parent)

			legacyPNG := makePNG(60, 30, color.NRGBA{G: 255, A: 255})
			tests.AssertRequest(cl.BannerUploadWithBodyWithResponse(ctx, "application/octet-stream", bytes.NewReader(legacyPNG), auth))(t, http.StatusOK)
			current = tests.AssertRequest(cl.AdminSettingsGetWithResponse(ctx, auth))(t, http.StatusOK)
			require.NotEqual(t, child.Id, current.JSON200.BannerAsset.Id)
			require.Nil(t, current.JSON200.BannerAsset.Parent)
			require.Equal(t, metadata, *current.JSON200.Metadata)
			banner = tests.AssertRequest(cl.BannerGetWithResponse(ctx, &openapi.BannerGetParams{Variant: &original}, auth))(t, http.StatusOK)
			require.Equal(t, legacyPNG, banner.Body)
		}))
	}))
}
