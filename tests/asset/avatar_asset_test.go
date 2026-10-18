package asset_test

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestAvatarUploadsReplaceAndClearTheAssetID(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		accounts *account_repo.Repository,
		writer *account_writer.Writer,
		avatars avatar.Service,
		download *asset_download.Downloader,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, owner := e2e.WithAccount(ctx, writer, seed.Account_001_Odin)
			require.False(t, avatars.Exists(ctx, owner.ID))
			firstImage := imagePayload(t)
			err := avatars.Set(ctx, owner.ID, bytes.NewReader(firstImage), int64(len(firstImage)))
			require.NoError(t, err)
			first, err := accounts.GetRefByID(ctx, owner.ID)
			require.NoError(t, err)
			firstID := first.AvatarAssetID
			require.True(t, firstID.Ok())
			require.True(t, avatars.Exists(ctx, owner.ID))

			var secondImage bytes.Buffer
			require.NoError(t, png.Encode(&secondImage, image.NewNRGBA(image.Rect(0, 0, 13, 27))))
			err = avatars.Set(ctx, owner.ID, bytes.NewReader(secondImage.Bytes()), int64(secondImage.Len()))
			require.NoError(t, err)
			second, err := accounts.GetRefByID(ctx, owner.ID)
			require.NoError(t, err)
			secondID := second.AvatarAssetID
			require.True(t, secondID.Ok())
			require.NotEqual(t, firstID, secondID)

			_, current, _, err := avatars.GetVariant(ctx, owner.ID, asset.Original)
			require.NoError(t, err)
			currentBytes, err := io.ReadAll(current)
			require.NoError(t, err)
			require.Equal(t, secondImage.Bytes(), currentBytes)
			if closer, ok := current.(io.Closer); ok {
				closer.Close()
			}

			require.NoError(t, avatars.Delete(ctx, owner.ID))
			cleared, err := accounts.GetRefByID(ctx, owner.ID)
			require.NoError(t, err)
			require.False(t, cleared.AvatarAssetID.Ok())
			require.False(t, avatars.Exists(ctx, owner.ID))

			_, generated, pending, err := avatars.GetVariant(ctx, owner.ID, asset.ThumbnailSmall)
			require.NoError(t, err)
			require.False(t, pending)
			decoded, err := png.Decode(generated)
			require.NoError(t, err)
			require.Equal(t, image.Pt(128, 128), decoded.Bounds().Size())
			if closer, ok := generated.(io.Closer); ok {
				closer.Close()
			}

			_, original, err := download.GetByID(ctx, firstID.OrZero())
			require.NoError(t, err)
			originalBytes, err := io.ReadAll(original)
			require.NoError(t, err)
			require.Equal(t, firstImage, originalBytes)
			if closer, ok := original.(io.Closer); ok {
				closer.Close()
			}
		}))
	}))
}

func TestBotWithoutAnAvatarAssetGetsTheDefaultAvatar(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		accounts *account_repo.Repository,
		writer *account_writer.Writer,
		avatars avatar.Service,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, bot := e2e.WithAccount(ctx, writer, seed.Account_001_Odin,
				account_writer.WithKind(account.AccountKindBot))
			stored, err := accounts.GetRefByID(ctx, bot.ID)
			require.NoError(t, err)
			require.False(t, stored.AvatarAssetID.Ok())

			metadata, stream, pending, err := avatars.GetVariant(ctx, bot.ID, asset.Original)
			require.NoError(t, err)
			defer stream.(io.Closer).Close()
			require.False(t, pending)
			require.Equal(t, "image/png", metadata.MIME.String())
			data, err := io.ReadAll(stream)
			require.NoError(t, err)
			decoded, err := png.Decode(bytes.NewReader(data))
			require.NoError(t, err)
			require.Equal(t, image.Pt(128, 128), decoded.Bounds().Size())
		}))
	}))
}
