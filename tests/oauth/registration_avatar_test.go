package oauth_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

type unavailableAvatarStorage struct {
	object.Storer
}

func (s *unavailableAvatarStorage) Write(context.Context, string, io.Reader, int64) error {
	return errors.New("avatar storage unavailable")
}

func TestAgentRegistrationStoresAnOrdinaryAvatar(t *testing.T) {
	for _, unavailable := range []bool{false, true} {
		name := "stored avatar"
		if unavailable {
			name = "storage unavailable"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			integration.Test(t, oauthConfig(t), e2e.Setup(), withOAuthRegistration(t, "open"),
				fx.Decorate(func(store object.Storer) object.Storer {
					if unavailable {
						return &unavailableAvatarStorage{Storer: store}
					}
					return store
				}),
				fx.Invoke(func(
					lc fx.Lifecycle,
					ctx context.Context,
					cl *openapi.ClientWithResponses,
					db *ent.Client,
					avatars avatar.Service,
					jobs *asset_image.Repository,
				) {
					lc.Append(fx.StartHook(func() {
						body, _ := agentRegistration(t)
						response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(ctx, body))(t, http.StatusCreated)
						client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(response.JSON201.ClientId)).Only(ctx)
						require.NoError(t, err)
						require.NotNil(t, client.AccountID)
						bot, err := db.Account.Get(ctx, *client.AccountID)
						require.NoError(t, err)
						id := account.AccountID(bot.ID)
						if unavailable {
							require.Nil(t, bot.AvatarAssetID)
							require.False(t, avatars.Exists(ctx, id))
							return
						}

						require.NotNil(t, bot.AvatarAssetID)
						source, err := db.Asset.Get(ctx, *bot.AvatarAssetID)
						require.NoError(t, err)
						require.Equal(t, bot.ID, source.AccountID)
						require.Equal(t, "image/png", source.MimeType)
						require.Nil(t, source.ParentAssetID)
						job, err := jobs.Get(ctx, source.ID)
						require.NoError(t, err)
						require.Equal(t, asset.EagerImagePresets, job.Payload.Requested)

						metadata, stream, pending, err := avatars.GetVariant(ctx, id, asset.Original)
						require.NoError(t, err)
						require.False(t, pending)
						require.Equal(t, asset.AssetID(source.ID), metadata.ID)
						decoded, err := png.Decode(stream)
						require.NoError(t, err)
						require.NoError(t, stream.(io.Closer).Close())
						require.Equal(t, image.Pt(256, 256), decoded.Bounds().Size())

						var replacement bytes.Buffer
						require.NoError(t, png.Encode(&replacement, image.NewNRGBA(image.Rect(0, 0, 13, 27))))
						require.NoError(t, avatars.Set(ctx, id, bytes.NewReader(replacement.Bytes()), int64(replacement.Len())))
						metadata, stream, _, err = avatars.GetVariant(ctx, id, asset.Original)
						require.NoError(t, err)
						require.NotEqual(t, asset.AssetID(source.ID), metadata.ID)
						data, err := io.ReadAll(stream)
						require.NoError(t, err)
						require.NoError(t, stream.(io.Closer).Close())
						require.Equal(t, replacement.Bytes(), data)

						require.NoError(t, avatars.Delete(ctx, id))
						require.False(t, avatars.Exists(ctx, id))
						_, stream, pending, err = avatars.GetVariant(ctx, id, asset.Original)
						require.NoError(t, err)
						require.False(t, pending)
						decoded, err = png.Decode(stream)
						require.NoError(t, err)
						require.NoError(t, stream.(io.Closer).Close())
						require.Equal(t, image.Pt(128, 128), decoded.Bounds().Size())
					}))
				}),
			)
		})
	}
}
