package asset_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestImageUploadDispatchesCommittedIntent(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			accounts *account_writer.Writer,
			uploader *asset_upload.Uploader,
			jobs *asset_image.Repository,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				ctx, owner := e2e.WithAccount(ctx, accounts, seed.Account_001_Odin)
				ctx = session.WithAccountPermissions(ctx, *owner, rbac.NewList(rbac.PermissionUploadAsset))
				image := imagePayload(t)
				started := time.Now()
				source, err := uploader.Upload(ctx, owner.ID, bytes.NewReader(image), int64(len(image)),
					asset.NewFilename("committed.png"), asset_upload.Options{})
				require.NoError(t, err)

				stored, err := db.Asset.Get(ctx, source.ID)
				require.NoError(t, err)
				require.Equal(t, source.ID, stored.ID)

				job, err := jobs.Get(ctx, source.ID)
				require.NoError(t, err)
				require.Equal(t, asset_image.Kind, job.Kind)
				require.Equal(t, asset.EagerImagePresets, job.Payload.Requested)
				require.True(t, job.DispatchUntil.After(started),
					"upload returned without dispatching its committed image job")

				document := []byte("%PDF-1.7\nexample document")
				file, err := uploader.Upload(ctx, owner.ID, bytes.NewReader(document), int64(len(document)),
					asset.NewFilename("document.pdf"), asset_upload.Options{})
				require.NoError(t, err)

				_, err = jobs.Get(ctx, file.ID)
				require.True(t, ent.IsNotFound(err))
			}))
		}),
	)
}
