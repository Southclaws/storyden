package asset_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"io"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestImageUploadLimits(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global media settings mutations on shared postgres database")
	}

	t.Parallel()

	integration.Test(t, &config.Config{
		ImageWorkers:      1,
		ImageMemoryBudget: 128 << 20,
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		accounts *account_writer.Writer,
		repo *settings.SettingsRepository,
		uploader *asset_upload.Uploader,
		objects object.Storer,
		client *openapi.ClientWithResponses,
		sessions *e2e.SessionHelper,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			auth := sessions.WithSession(ctx)

			var encoded bytes.Buffer
			require.NoError(t, png.Encode(&encoded, image.NewNRGBA(image.Rect(0, 0, 100, 80))))
			payload := encoded.Bytes()

			setLimits := func(size, pixels int64) {
				t.Helper()
				_, err := repo.Set(ctx, settings.Settings{
					Services: opt.New(settings.ServiceSettings{
						Media: opt.New(settings.MediaServiceSettings{
							ImageMaxSize:   opt.New(size),
							ImageMaxPixels: opt.New(pixels),
						}),
					}),
				})
				require.NoError(t, err)
			}

			t.Run("exact_limits_preserve_original_bytes", func(t *testing.T) {
				setLimits(int64(len(payload)), 8000)
				name := "within-limits.png"
				response := tests.AssertRequest(client.AssetUploadWithBodyWithResponse(ctx,
					&openapi.AssetUploadParams{Filename: &name, ContentLength: int64(len(payload))},
					"image/png", bytes.NewReader(payload), auth,
				))(t, http.StatusOK)

				original := tests.AssertRequest(client.AssetGetWithResponse(ctx,
					response.JSON200.Filename, nil, auth,
				))(t, http.StatusOK)
				require.Equal(t, payload, original.Body)
			})

			t.Run("oversized_dimensions_rejected_by_http", func(t *testing.T) {
				setLimits(1048576, 7999)
				name := "too-many-pixels.png"
				tests.AssertRequest(client.AssetUploadWithBodyWithResponse(ctx,
					&openapi.AssetUploadParams{Filename: &name, ContentLength: int64(len(payload))},
					"image/png", bytes.NewReader(payload), auth,
				))(t, http.StatusBadRequest)
			})

			t.Run("tiny_image_with_huge_pixel_surface_is_rejected", func(t *testing.T) {
				setLimits(1048576, 12000000)
				crafted := bytes.Clone(payload)
				binary.BigEndian.PutUint32(crafted[16:20], 65535)
				binary.BigEndian.PutUint32(crafted[20:24], 65535)
				binary.BigEndian.PutUint32(crafted[29:33], crc32.ChecksumIEEE(crafted[12:29]))

				name := "oversized-surface.png"
				response := tests.AssertRequest(client.AssetUploadWithBodyWithResponse(ctx,
					&openapi.AssetUploadParams{Filename: &name, ContentLength: int64(len(crafted))},
					"image/png", bytes.NewReader(crafted), auth,
				))(t, http.StatusBadRequest)
				require.Contains(t, string(response.Body), "maximum resolution")
			})

			t.Run("actual_bytes_are_checked_and_candidate_is_deleted", func(t *testing.T) {
				setLimits(int64(len(payload)-1), 8000)
				for _, declared := range []int64{1, -1} {
					name := asset.NewFilename("encoded-limit.png")
					_, err := uploader.Upload(ctx, account.ID, bytes.NewReader(payload), declared, name, asset_upload.Options{})
					require.Error(t, err)

					_, _, err = objects.Read(ctx, asset.BuildAssetPath(name))
					require.Error(t, err, "rejected candidate must be removed")
				}
			})

			t.Run("generic_files_are_unaffected", func(t *testing.T) {
				setLimits(1, 1)
				data := []byte("%PDF-1.7\nA document larger than the image byte limit.")
				name := asset.NewFilename("document.pdf")
				uploaded, err := uploader.Upload(ctx, account.ID, bytes.NewReader(data), int64(len(data)), name, asset_upload.Options{})
				require.NoError(t, err)

				reader, _, err := objects.Read(ctx, asset.BuildAssetPath(uploaded.Name))
				require.NoError(t, err)
				if closer, ok := reader.(io.Closer); ok {
					defer closer.Close()
				}
				stored, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, data, stored)
			})
		}))
	}))
}
