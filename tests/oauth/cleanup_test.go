package oauth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_querier"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthCleanup(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("OAuth cleanup is intentionally global")
	}

	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		oq *oauth_querier.Querier,
		ow *oauth_writer.Writer,
		cl *openapi.ClientWithResponses,
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
				rotatedParentToken := "rotated-parent-" + uuid.NewString()
				rotatedChildToken := "rotated-child-" + uuid.NewString()
				hashToken := func(token string) string {
					sum := sha256.Sum256([]byte(token))
					return hex.EncodeToString(sum[:])
				}
				now := time.Now()
				cleanupNow := now.Add(60 * 24 * time.Hour)

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

				create(staleHash, cleanupNow.Add(-30*24*time.Hour))
				create(recentlyExpiredHash, cleanupNow.Add(-time.Hour))
				create(liveHash, cleanupNow.Add(24*time.Hour))

				// An expired, rotated parent must survive the simulated cleanup while its
				// child is retained, so reuse detection can revoke the whole token family.
				parent, err := ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
					ClientID:  client.ID,
					AccountID: owner.ID,
					TokenHash: hashToken(rotatedParentToken),
					Scope:     "openid",
					ExpiresAt: now.Add(-time.Hour),
				})
				r.NoError(err)
				child, err := ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
					ClientID:  client.ID,
					AccountID: owner.ID,
					TokenHash: hashToken(rotatedChildToken),
					Scope:     "openid",
					ExpiresAt: cleanupNow.Add(24 * time.Hour),
				})
				r.NoError(err)
				rotated, err := ow.RevokeRefreshToken(root, parent.ID, now, opt.New(child.ID))
				r.NoError(err)
				r.True(rotated)

				deleted, err := ow.DeleteExpiredRefreshTokens(root, cleanupNow.Add(-7*24*time.Hour))
				r.NoError(err)
				a.Positive(deleted)

				_, err = oq.GetRefreshTokenByTokenHash(root, staleHash)
				a.Error(err)

				_, err = oq.GetRefreshTokenByTokenHash(root, recentlyExpiredHash)
				a.NoError(err, "a token still inside the window stays so the rotation chain can be walked")

				_, err = oq.GetRefreshTokenByTokenHash(root, liveHash)
				a.NoError(err)

				_, err = oq.GetRefreshTokenByTokenHash(root, hashToken(rotatedParentToken))
				a.NoError(err, "a rotated parent stays so reuse detection can walk to its retained child")

				reuse := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantRefreshToken,
					ClientId:     clientID,
					RefreshToken: &rotatedParentToken,
				}))(t, http.StatusBadRequest)
				r.NotNil(reuse.JSON400)
				a.Equal("invalid_grant", reuse.JSON400.Error)

				familyRevoked := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantRefreshToken,
					ClientId:     clientID,
					RefreshToken: &rotatedChildToken,
				}))(t, http.StatusBadRequest)
				r.NotNil(familyRevoked.JSON400)
				a.Equal("invalid_grant", familyRevoked.JSON400.Error)
			})
		}))
	}))
}

func TestRefreshTokenListIsPaginated(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		oq *oauth_querier.Querier,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			_, owner := e2e.WithAccount(root, aw, seed.Account_002_Frigg)

			clientID := "paged-refresh-tokens-" + uuid.NewString()
			client := createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauth.GrantTypeRefreshToken})

			const total = 12
			for range total {
				_, err := ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
					ClientID:  client.ID,
					AccountID: owner.ID,
					TokenHash: uuid.NewString(),
					Scope:     "openid",
					ExpiresAt: time.Now().Add(24 * time.Hour),
				})
				r.NoError(err)
			}

			first, err := oq.ListRefreshTokensByAccount(root, owner.ID, pagination.NewPageParams(1, 5))
			r.NoError(err)
			a.Len(first.Items, 5, "a page must not hand back the whole table")
			a.Equal(5, first.Results)
			a.Equal(3, first.TotalPages, "total_pages is how the client learns there are 12 rows")
			a.True(first.NextPage.Ok())

			last, err := oq.ListRefreshTokensByAccount(root, owner.ID, pagination.NewPageParams(3, 5))
			r.NoError(err)
			a.Len(last.Items, 2)
			a.False(last.NextPage.Ok(), "the final page must not advertise another")

			seen := map[string]struct{}{}
			for _, item := range append(first.Items, last.Items...) {
				seen[xid.ID(item.ID).String()] = struct{}{}
			}
			a.Len(seen, 7, "pages must not overlap")
		}))
	}))
}

func TestRefreshTokenPagesAreStableWhenTimestampsCollide(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		oq *oauth_querier.Querier,
		ow *oauth_writer.Writer,
		raw *sqlx.DB,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			_, owner := e2e.WithAccount(root, aw, seed.Account_003_Baldur)

			clientID := "collided-timestamps-" + uuid.NewString()
			client := createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauth.GrantTypeRefreshToken})

			const total = 9
			created := make([]oauthresource.RefreshTokenID, 0, total)
			for range total {
				rt, err := ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
					ClientID:  client.ID,
					AccountID: owner.ID,
					TokenHash: uuid.NewString(),
					Scope:     "openid",
					ExpiresAt: time.Now().Add(24 * time.Hour),
				})
				r.NoError(err)
				created = append(created, rt.ID)
			}

			// force every row onto the same instant, which is what an ordering
			// on created_at alone cannot separate
			shared := time.Now().Truncate(time.Second)
			for _, id := range created {
				_, err := raw.ExecContext(root,
					raw.Rebind("update oauth_refresh_tokens set created_at = ? where id = ?"),
					shared, xid.ID(id).String())
				r.NoError(err)
			}

			seen := map[string]struct{}{}
			for page := uint(1); page <= 3; page++ {
				result, err := oq.ListRefreshTokensByAccount(root, owner.ID, pagination.NewPageParams(page, 3))
				r.NoError(err)

				for _, item := range result.Items {
					id := xid.ID(item.ID).String()
					_, repeated := seen[id]
					a.False(repeated, "a row must not appear on two pages")
					seen[id] = struct{}{}
				}
			}

			a.Len(seen, total, "walking the pages must reach every row exactly once")
		}))
	}))
}
