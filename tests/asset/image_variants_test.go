package asset_test

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"golang.org/x/image/webp"

	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	image_jobs "github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	image_processor "github.com/Southclaws/storyden/app/services/asset/asset_image"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/avatar"
	job_runner "github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/app/services/branding/banner"
	"github.com/Southclaws/storyden/app/services/branding/icon"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	ea "github.com/Southclaws/storyden/internal/ent/asset"
	ej "github.com/Southclaws/storyden/internal/ent/backgroundjob"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func imagePayload(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 1600, 800))
	for y := 0; y < 800; y++ {
		for x := 0; x < 1600; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x), G: uint8(y), B: 100, A: 255})
		}
	}
	var out bytes.Buffer
	require.NoError(t, png.Encode(&out, img))
	return out.Bytes()
}

func TestImageVariantsHTTP(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("isolated file storage is incompatible with shared job workers")
	}

	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		uploader *asset_upload.Uploader,
		db *ent.Client,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		jobs *image_jobs.Repository,
		accounts *account_repo.Repository,
		settingsRepo *settings.SettingsRepository,
		avatars avatar.Service,
		icons icon.Service,
		banners banner.Service,
		_ *image_processor.Processor,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, acc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *acc, rbac.NewList(rbac.PermissionUploadAsset))
			auth := sh.WithSession(ctx)
			payload := imagePayload(t)
			a, err := uploader.Upload(
				ctx,
				acc.ID,
				bytes.NewReader(payload),
				int64(len(payload)),
				asset.NewFilename("cover.png"),
				asset_upload.Options{},
			)
			require.NoError(t, err)

			original, err := cl.AssetGetWithResponse(root, a.Name.String(), nil, auth)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, original.StatusCode())
			require.Equal(t, payload, original.Body)
			require.Contains(t, original.HTTPResponse.Header.Get("Cache-Control"), "31536000")

			medium := openapi.M
			mediumParams := &openapi.AssetGetParams{Variant: &medium}
			pending, err := cl.AssetGetWithResponse(root, a.Name.String(), mediumParams, auth)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, pending.StatusCode())
			require.Contains(t, pending.HTTPResponse.Header.Get("Cache-Control"), "no-cache")

			require.Eventually(t, func() bool {
				job, err := jobs.Get(ctx, a.ID)
				return err == nil && job.Status == background_job.StatusComplete
			}, 20*time.Second, 50*time.Millisecond)

			rows, err := jobs.Variants(ctx, a.ID)
			require.NoError(t, err)
			require.Len(t, rows, 5)
			expected := map[asset.Variant]image.Point{
				asset.Small:           {320, 160},
				asset.Medium:          {960, 480},
				asset.ThumbnailSmall:  {64, 64},
				asset.ThumbnailMedium: {128, 128},
				asset.ThumbnailLarge:  {256, 256},
			}
			for _, child := range rows {
				response, err := cl.AssetGetWithResponse(root, child.Name.String(), nil, auth)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, response.StatusCode())
				decoded, err := webp.Decode(bytes.NewReader(response.Body))
				require.NoError(t, err)
				metadata, err := child.ImageMetadata()
				require.NoError(t, err)
				require.NotNil(t, metadata)
				require.Equal(t, expected[metadata.Preset], decoded.Bounds().Size())
				stored, err := db.Asset.Get(ctx, child.ID)
				require.NoError(t, err)
				require.Equal(t, a.ID, *stored.ParentAssetID)
				require.Equal(t, xid.ID(acc.ID), stored.AccountID)
				var storedMetadata asset.ImageMetadata
				require.NoError(t, json.Unmarshal(stored.Metadata, &storedMetadata))
				require.Equal(t, metadata.Preset, storedMetadata.Preset)
			}

			ready, err := cl.AssetGetWithResponse(root, a.Name.String(), mediumParams, auth)
			require.NoError(t, err)
			require.Equal(t, "image/webp", ready.HTTPResponse.Header.Get("Content-Type"))
			require.Equal(t, len(ready.Body), int(ready.HTTPResponse.ContentLength))
			require.Contains(t, ready.HTTPResponse.Header.Get("Cache-Control"), "86400")

			require.NoError(t, jobs.Request(ctx, a.ID, asset.EagerImagePresets|2))
			count, err := db.Asset.Query().Where(ea.ParentAssetIDEQ(a.ID)).Count(ctx)
			require.NoError(t, err)
			require.Equal(t, 5, count)

			bad := openapi.AssetVariantQuery("unknown")
			invalid, err := cl.AssetGetWithResponse(root, a.Name.String(), &openapi.AssetGetParams{Variant: &bad}, auth)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, invalid.StatusCode())

			document := []byte("%PDF-1.7 storyden attachment")
			doc, err := uploader.Upload(
				ctx,
				acc.ID,
				bytes.NewReader(document),
				int64(len(document)),
				asset.NewFilename("guide.pdf"),
				asset_upload.Options{},
			)
			require.NoError(t, err)
			docResponse, err := cl.AssetGetWithResponse(root, doc.Name.String(), mediumParams, auth)
			require.NoError(t, err)
			require.Equal(t, document, docResponse.Body)
			_, err = jobs.Get(ctx, doc.ID)
			require.True(t, ent.IsNotFound(err))

			require.NoError(t, avatars.Set(ctx, acc.ID, bytes.NewReader(payload), int64(len(payload))))
			require.Eventually(t, func() bool {
				owner, err := accounts.GetRefByID(ctx, acc.ID)
				if err != nil || !owner.AvatarAssetID.Ok() {
					return false
				}
				j, err := jobs.Get(ctx, owner.AvatarAssetID.OrZero())
				return err == nil && j.Status == background_job.StatusComplete
			}, 20*time.Second, 50*time.Millisecond)
			thumb := openapi.ThumbM
			avatarParams := &openapi.AccountGetAvatarParams{Variant: &thumb}
			avatarResponse, err := cl.AccountGetAvatarWithResponse(root, acc.Handle, avatarParams, auth)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, avatarResponse.StatusCode())
			require.Equal(t, "image/webp", avatarResponse.HTTPResponse.Header.Get("Content-Type"))
			require.Equal(t, "default-src 'none'; sandbox", avatarResponse.HTTPResponse.Header.Get("Content-Security-Policy"))
			avatarImage, err := webp.Decode(bytes.NewReader(avatarResponse.Body))
			require.NoError(t, err)
			require.Equal(t, image.Pt(128, 128), avatarImage.Bounds().Size())

			require.NoError(t, icons.Upload(ctx, bytes.NewReader(payload)))
			require.Eventually(t, func() bool {
				configured, err := settingsRepo.Get(ctx)
				if err != nil || !configured.IconAssetID.Ok() {
					return false
				}
				j, err := jobs.Get(ctx, configured.IconAssetID.OrZero())
				return err == nil && j.Status == background_job.StatusComplete
			}, 20*time.Second, 50*time.Millisecond)
			iconAsset, iconReader, err := icons.Get(ctx, "512x512")
			require.NoError(t, err)
			require.Equal(t, "image/png", iconAsset.MIME.String())
			iconImage, err := png.Decode(iconReader)
			require.NoError(t, err)
			require.Equal(t, image.Pt(512, 512), iconImage.Bounds().Size())
			if c, ok := iconReader.(io.Closer); ok {
				c.Close()
			}

			svg := bytes.NewBufferString(`<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32"/>`)
			require.NoError(t, icons.Upload(ctx, svg))
			fallbackIcon, fallbackReader, err := icons.Get(ctx, "512x512")
			require.NoError(t, err)
			require.Equal(t, "image/png", fallbackIcon.MIME.String())
			_, err = png.Decode(fallbackReader)
			require.NoError(t, err)
			if c, ok := fallbackReader.(io.Closer); ok {
				c.Close()
			}

			require.NoError(t, banners.Upload(ctx, bytes.NewReader(payload)))
			require.Eventually(t, func() bool {
				configured, err := settingsRepo.Get(ctx)
				if err != nil || !configured.BannerAssetID.Ok() {
					return false
				}
				j, err := jobs.Get(ctx, configured.BannerAssetID.OrZero())
				return err == nil && j.Status == background_job.StatusComplete
			}, 20*time.Second, 50*time.Millisecond)
			bannerAsset, bannerReader, err := banners.Get(ctx)
			require.NoError(t, err)
			require.Equal(t, "image/png", bannerAsset.MIME.String())
			bannerImage, err := png.Decode(bannerReader)
			require.NoError(t, err)
			require.Equal(t, image.Pt(1200, 630), bannerImage.Bounds().Size())
			if c, ok := bannerReader.(io.Closer); ok {
				c.Close()
			}

			configured, err := settingsRepo.Get(ctx)
			require.NoError(t, err)
			storedSettings, err := db.Setting.Get(ctx, settings.StorydenPrimarySettingsKey)
			require.NoError(t, err)
			var persisted settings.Settings
			require.NoError(t, json.Unmarshal([]byte(storedSettings.Value), &persisted))
			require.Equal(t, configured.IconAssetID, persisted.IconAssetID)
			require.Equal(t, configured.BannerAssetID, persisted.BannerAssetID)

		}))
	}))
}

type testPayload struct {
	Value string `json:"value"`
}

func TestInvalidImageJobPayloadDoesNotBlockPendingWork(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		runner *job_runner.Runner,
		jobs *background_job.Repository,
		db *ent.Client,
		_ *image_processor.Processor,
	) {
		err := runner.Register("test.valid", 1, job_runner.Handler[testPayload]{
			Run: func(context.Context, *background_job.Job[testPayload]) (job_runner.Result, error) {
				return job_runner.Result{}, nil
			},
		})
		require.NoError(t, err)

		lc.Append(fx.StartHook(func() {
			invalidID := xid.New()
			invalid := json.RawMessage(`{"requested":"invalid"}`)
			err := jobs.Ensure(ctx, invalidID, image_jobs.Kind, invalidID.String(), invalid)
			require.NoError(t, err)

			validID := xid.New()
			err = jobs.Ensure(ctx, validID, "test.valid", validID.String(), testPayload{Value: "valid"})
			require.NoError(t, err)

			runner.Dispatch(ctx, invalidID)
			runner.Dispatch(ctx, validID)

			require.Eventually(t, func() bool {
				invalidJob, err := db.BackgroundJob.Get(ctx, invalidID)
				if err != nil || invalidJob.Status != ej.StatusSkipped {
					return false
				}
				returned, err := db.BackgroundJob.Get(ctx, validID)
				return err == nil && returned.Status == ej.StatusComplete
			}, 5*time.Second, 20*time.Millisecond)

			invalidJob, err := db.BackgroundJob.Get(ctx, invalidID)
			require.NoError(t, err)
			require.Contains(t, invalidJob.Reason, "requested")
			require.Equal(t, 1, invalidJob.Attempts)
		}))
	}))
}

func TestBackgroundJobLeaseAndRecovery(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		jobs *background_job.Repository,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			id := xid.New()
			require.NoError(t, jobs.Ensure(ctx, id, "test", id.String(), testPayload{Value: "original"}))
			first, err := jobs.Claim(ctx, id, 3)
			require.NoError(t, err)
			require.NotNil(t, first)
			duplicate, err := jobs.Claim(ctx, id, 3)
			require.NoError(t, err)
			require.Nil(t, duplicate)
			require.NoError(t, db.BackgroundJob.UpdateOneID(id).SetLeaseUntil(time.Now().Add(-time.Minute)).Exec(ctx))
			second, err := jobs.Claim(ctx, id, 3)
			require.NoError(t, err)
			require.NotNil(t, second)
			require.NotEqual(t, first.LeaseToken, second.LeaseToken)
			validateLease := func(*testPayload) error {
				return nil
			}
			err = jobs.UpdatePayload(ctx, id, first.LeaseToken, validateLease)
			require.ErrorIs(t, err, background_job.ErrLeaseLost)
			require.NoError(t, jobs.Heartbeat(ctx, second.State))
			require.NoError(t, jobs.Finish(ctx, second.State, "", false, 3, nil))
			finished, err := jobs.Get[json.RawMessage](ctx, id)
			require.NoError(t, err)
			require.Equal(t, background_job.StatusComplete, finished.Status)
			var finishedPayload testPayload
			require.NoError(t, json.Unmarshal(finished.Payload, &finishedPayload))
			require.Equal(t, "original", finishedPayload.Value)

			concurrentID := xid.New()
			require.NoError(t, jobs.Ensure(ctx, concurrentID, "test", concurrentID.String(), testPayload{}))
			type claimResult struct {
				job *background_job.Job[json.RawMessage]
				err error
			}
			results := make(chan claimResult, 2)
			for range 2 {
				go func() {
					job, err := background_job.New(db).Claim(ctx, concurrentID, 3)
					results <- claimResult{
						job: job,
						err: err,
					}
				}()
			}
			claimed := 0
			for range 2 {
				result := <-results
				require.NoError(t, result.err)
				if result.job != nil {
					claimed++
				}
			}
			require.Equal(t, 1, claimed)
			failedID := xid.New()
			require.NoError(t, jobs.Ensure(ctx, failedID, "test", failedID.String(), testPayload{}))
			attempt, err := jobs.Claim(ctx, failedID, 1)
			require.NoError(t, err)
			require.NotNil(t, attempt)
			require.NoError(t, jobs.Finish(ctx, attempt.State, "storage failed", false, 1, nil))
			failed, err := jobs.Get[json.RawMessage](ctx, failedID)
			require.NoError(t, err)
			require.Equal(t, background_job.StatusFailed, failed.Status)
			require.Equal(t, "storage failed", failed.Reason)
		}))
	}))
}
