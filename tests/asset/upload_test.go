package asset_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	authsession "github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestUploadRecordsBytesWrittenNotDeclaredSize(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		uploader *asset_upload.Uploader,
		objects object.Storer,
	) {
		lc.Append(fx.StartHook(func() {
			_, acc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			ctx := authsession.WithAccountPermissions(root, *acc, rbac.NewList(rbac.PermissionUploadAsset))

			payload := bytes.Repeat([]byte("storyden"), 512)
			require.Len(t, payload, 4096)

			name := asset.NewFilename("declared-size-lies.txt")

			a, err := uploader.Upload(ctx, bytes.NewReader(payload), 10, name, asset_upload.Options{})
			require.NoError(t, err)

			assert.Equal(t, len(payload), a.Size, "the recorded size must be the bytes actually stored")

			r, size, err := objects.Read(ctx, asset.BuildAssetPath(a.Name))
			require.NoError(t, err)

			stored, err := io.ReadAll(r)
			require.NoError(t, err)

			assert.Equal(t, len(payload), len(stored))
			assert.Equal(t, int64(a.Size), size, "AssetGet serves a.Size as Content-Length, it must match the object")
		}))
	}))
}

func TestUploadRecordsSizeForUnknownLength(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		uploader *asset_upload.Uploader,
		objects object.Storer,
	) {
		lc.Append(fx.StartHook(func() {
			_, acc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			ctx := authsession.WithAccountPermissions(root, *acc, rbac.NewList(rbac.PermissionUploadAsset))

			payload := bytes.Repeat([]byte("chunked"), 128)
			name := asset.NewFilename("unknown-length.txt")

			a, err := uploader.Upload(ctx, bytes.NewReader(payload), -1, name, asset_upload.Options{})
			require.NoError(t, err)

			assert.Equal(t, len(payload), a.Size, "a chunked upload reports -1 and must still record a real size")

			_, size, err := objects.Read(ctx, asset.BuildAssetPath(a.Name))
			require.NoError(t, err)
			assert.Equal(t, int64(a.Size), size)
		}))
	}))
}
