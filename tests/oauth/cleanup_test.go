package oauth_test

import (
	"context"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_querier"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestOAuthCleanup(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		oq *oauth_querier.Querier,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			_, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)

			t.Run("expired_authorization_requests_are_deleted", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "cleanup-auth-request-" + uuid.NewString()
				client := createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, "cleanup-secret")), standardScopes(), []string{oauth.GrantTypeAuthorizationCode})

				expiredHash := "expired-" + uuid.NewString()
				activeHash := "active-" + uuid.NewString()
				_, err := ow.CreateAuthorisationRequest(root, oauth_writer.AuthorisationRequestCreate{
					ClientID:            client.ID,
					AccountID:           owner.ID,
					RequestIDHash:       expiredHash,
					RedirectURI:         "https://client.example/callback",
					CodeChallenge:       codeChallenge("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
					CodeChallengeMethod: "S256",
					ExpiresAt:           time.Now().Add(-time.Minute),
				})
				r.NoError(err)
				_, err = ow.CreateAuthorisationRequest(root, oauth_writer.AuthorisationRequestCreate{
					ClientID:            client.ID,
					AccountID:           owner.ID,
					RequestIDHash:       activeHash,
					RedirectURI:         "https://client.example/callback",
					CodeChallenge:       codeChallenge("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
					CodeChallengeMethod: "S256",
					ExpiresAt:           time.Now().Add(time.Minute),
				})
				r.NoError(err)

				deleted, err := ow.DeleteExpiredAuthorisationRequests(root, time.Now())
				r.NoError(err)
				a.Equal(1, deleted)

				_, err = oq.GetAuthorisationRequestByRequestIDHash(root, expiredHash)
				a.Error(err)
				_, err = oq.GetAuthorisationRequestByRequestIDHash(root, activeHash)
				a.NoError(err)
			})

			t.Run("refresh_tokens_are_dropped_once_they_are_long_expired", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "cleanup-refresh-token-" + uuid.NewString()
				client := createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauth.GrantTypeRefreshToken})

				staleHash := "stale-" + uuid.NewString()
				recentlyExpiredHash := "recent-" + uuid.NewString()
				liveHash := "live-" + uuid.NewString()

				create := func(hash string, expiresAt time.Time) {
					_, err := ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
						ClientID:  client.ID,
						AccountID: owner.ID,
						TokenHash: hash,
						Scope:     "openid",
						ExpiresAt: expiresAt,
					})
					r.NoError(err)
				}

				create(staleHash, time.Now().Add(-30*24*time.Hour))
				create(recentlyExpiredHash, time.Now().Add(-time.Hour))
				create(liveHash, time.Now().Add(24*time.Hour))

				deleted, err := ow.DeleteExpiredRefreshTokens(root, time.Now().Add(-7*24*time.Hour))
				r.NoError(err)
				a.Equal(1, deleted, "only rows past the retention window go")

				_, err = oq.GetRefreshTokenByTokenHash(root, staleHash)
				a.Error(err)

				_, err = oq.GetRefreshTokenByTokenHash(root, recentlyExpiredHash)
				a.NoError(err, "a token still inside the window stays so the rotation chain can be walked")

				_, err = oq.GetRefreshTokenByTokenHash(root, liveHash)
				a.NoError(err)
			})
		}))
	}))
}
