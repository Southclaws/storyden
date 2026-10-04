package oauth_test

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthClientAssertionReplay(t *testing.T) {
	for _, provider := range []string{"local", "redis"} {
		t.Run(provider, func(t *testing.T) {
			cfg := oauthConfig(t)
			if provider == "redis" {
				address := os.Getenv("REDIS_URL")
				if address == "" {
					t.Skip("REDIS_URL is not set")
				}

				u, err := url.Parse(address)
				require.NoError(t, err)
				cfg.RedisURL = *u
				cfg.CacheProvider = "redis"
			}

			var first *openapi.ClientWithResponses
			var database *ent.Client
			var store cache.Store
			var redis rueidis.Client
			integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, "open"), fx.Populate(&first, &database, &store, &redis))

			ctx := context.Background()
			body, key := agentRegistration(t)
			registered := tests.AssertRequest(first.OAuthClientRegisterWithResponse(ctx, body))(t, http.StatusCreated).JSON201

			t.Run("invalid_signature_does_not_consume", func(t *testing.T) {
				jti := uuid.NewString()
				_, wrongKey := agentRegistration(t)
				invalid := signedAssertionForm(t, registered.ClientId, wrongKey, jti)
				response := tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(invalid.Encode())))(t, http.StatusBadRequest)
				require.Equal(t, "invalid_client", response.JSON400.Error)

				form := signedAssertionForm(t, registered.ClientId, key, jti)
				tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
			})

			replace := fx.Replace(database)
			if provider == "local" {
				replace = fx.Options(replace, fx.Replace(fx.Annotate(store, fx.As(new(cache.Store)))))
			}

			t.Run("concurrent_instances", func(t *testing.T) {
				form := signedAssertionForm(t, registered.ClientId, key, uuid.NewString())
				var second *openapi.ClientWithResponses
				integration.Test(t, cfg, e2e.Setup(), replace, fx.Populate(&second))

				const requests = 16
				responses := make([]*openapi.OAuthTokenResponse, requests)
				errors := make([]error, requests)
				start := make(chan struct{})
				var wg sync.WaitGroup
				for i := range requests {
					wg.Go(func() {
						<-start
						client := []*openapi.ClientWithResponses{first, second}[i%2]
						responses[i], errors[i] = client.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
					})
				}
				close(start)
				wg.Wait()

				accepted := 0
				for i, response := range responses {
					require.NoError(t, errors[i])
					if response.StatusCode() == http.StatusOK {
						accepted++
						require.NotNil(t, response.JSON200.AccessToken)
					} else {
						require.Equal(t, http.StatusBadRequest, response.StatusCode(), string(response.Body))
						require.Equal(t, "invalid_client", response.JSON400.Error)
					}
				}
				require.Equal(t, 1, accepted)
			})

			t.Run("recreated_instance", func(t *testing.T) {
				form := signedAssertionForm(t, registered.ClientId, key, uuid.NewString())
				tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)

				var restarted *openapi.ClientWithResponses
				integration.Test(t, cfg, e2e.Setup(), replace, fx.Populate(&restarted))

				response := tests.AssertRequest(restarted.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusBadRequest)
				require.Equal(t, "invalid_client", response.JSON400.Error)
			})

			t.Run("independent_clients", func(t *testing.T) {
				jti := uuid.NewString()
				firstForm := signedAssertionForm(t, registered.ClientId, key, jti)
				tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(firstForm.Encode())))(t, http.StatusOK)

				body, key := agentRegistration(t)
				registered := tests.AssertRequest(first.OAuthClientRegisterWithResponse(ctx, body))(t, http.StatusCreated).JSON201
				form := signedAssertionForm(t, registered.ClientId, key, jti)
				tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
			})

			if provider == "redis" {
				t.Run("cache_failure_denies_token", func(t *testing.T) {
					form := signedAssertionForm(t, registered.ClientId, key, uuid.NewString())
					redis.Close()
					response := tests.AssertRequest(first.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusInternalServerError)
					require.Nil(t, response.JSON200)
				})
			}
		})
	}
}

func signedAssertionForm(t *testing.T, clientID string, key *rsa.PrivateKey, id string) url.Values {
	t.Helper()

	assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  jwt.ClaimStrings{"http://localhost:8000/api/oauth/token"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		ID:        id,
	})
	assertion.Header["kid"] = "agent-key"
	signed, err := assertion.SignedString(key)
	require.NoError(t, err)

	return url.Values{
		"grant_type":            {"client_credentials"},
		"client_id":             {clientID},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"client_assertion":      {signed},
	}
}
