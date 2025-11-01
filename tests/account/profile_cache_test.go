package account_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAccountCacheWithEmailOperations(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		cacheStore cache.Store,
		profileCache *profile_cache.Cache,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			email1 := xid.New().String() + "first@example.com"
			password := "password"

			signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
				Email:    email1,
				Password: password,
			})
			tests.Ok(t, err, signup)
			session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

			accountID := openapi.ParseID(signup.JSON200.Id)

			accGet1, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet1)
			a.Len(accGet1.JSON200.EmailAddresses, 1, "account should have one email")

			lastModified1Header := accGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present")

			cacheKey := "profile:last-modified:" + accountID.String()

			email2 := xid.New().String() + "second@example.com"

			addEmail, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
				EmailAddress: email2,
			}, session)
			tests.Ok(t, err, addEmail)

			// Wait for the cache update event to be processed via pubsub
			time.Sleep(200 * time.Millisecond)

			// Now cache should be populated with timestamp from when email was added
			cachedValue1, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should exist after email added")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			accGet2, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet2)
			r.Len(accGet2.JSON200.EmailAddresses, 2, "account should now have two emails")

			lastModified2Header := accGet2.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified2Header, "Last-Modified header should be present")

			// Make a conditional request - should return 304 since cache timestamp is newer
			accGet304, err := cl.AccountGetWithResponse(root, session, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified2Header)
				return nil
			})
			tests.Status(t, err, accGet304, 304)

			// Remove the email
			removeEmail, err := cl.AccountEmailRemoveWithResponse(root, addEmail.JSON200.Id, session)
			tests.Ok(t, err, removeEmail)

			// Wait for the cache update event to be processed
			time.Sleep(200 * time.Millisecond)

			// Cache should be updated with a newer timestamp
			cachedValue2, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should still exist after email removed")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			a.True(cachedTime2.After(cachedTime1), "cached time should be updated after removing email (cachedTime1: %v, cachedTime2: %v)", cachedTime1, cachedTime2)

			accGet3, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet3)
			r.Len(accGet3.JSON200.EmailAddresses, 1, "account should be back to one email")
		}))
	}))
}

func TestAccountCacheWithProfileUpdate(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		cacheStore cache.Store,
		profileCache *profile_cache.Cache,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			email := uuid.NewString() + "@example.com"
			password := "password"

			signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
				Email:    email,
				Password: password,
			})
			tests.Ok(t, err, signup)
			session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

			accountID := openapi.ParseID(signup.JSON200.Id)

			cacheKey := "profile:last-modified:" + accountID.String()

			// Update the profile
			newBio := "This is my new bio"
			updateResp, err := cl.AccountUpdateWithResponse(root, openapi.AccountMutableProps{
				Bio: &newBio,
			}, session)
			tests.Ok(t, err, updateResp)

			// Wait for the cache update event to be processed
			time.Sleep(200 * time.Millisecond)

			// Cache should now be populated
			cachedValue1, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should exist after profile update")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			accGet2, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet2)
			a.Contains(accGet2.JSON200.Bio, newBio, "bio should contain the updated text")

			// Update again
			newBio2 := "This is my second bio"
			updateResp2, err := cl.AccountUpdateWithResponse(root, openapi.AccountMutableProps{
				Bio: &newBio2,
			}, session)
			tests.Ok(t, err, updateResp2)

			time.Sleep(200 * time.Millisecond)

			cachedValue2, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should still exist after second profile update")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			a.True(cachedTime2.After(cachedTime1), "cached time should be updated after profile update (cachedTime1: %v, cachedTime2: %v)", cachedTime1, cachedTime2)
		}))
	}))
}
