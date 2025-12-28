package account_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
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

			accGet1, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet1)
			a.Len(accGet1.JSON200.EmailAddresses, 1, "account should have one email")

			etag1 := accGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")

			email2 := xid.New().String() + "second@example.com"

			addEmail, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
				EmailAddress: email2,
			}, session)
			tests.Ok(t, err, addEmail)

			accGet2, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet2)
			r.Len(accGet2.JSON200.EmailAddresses, 2, "account should now have two emails")

			etag2 := accGet2.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag2, "ETag header should be present")
			a.NotEqual(etag1, etag2, "ETag should change after adding email")

			// Make a conditional request - should return 304 since cache is current
			accGet304, err := cl.AccountGetWithResponse(root, session, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag2)
				return nil
			})
			tests.Status(t, err, accGet304, 304)

			// Remove the email
			removeEmail, err := cl.AccountEmailRemoveWithResponse(root, addEmail.JSON200.Id, session)
			tests.Ok(t, err, removeEmail)

			// Conditional GET with old ETag should now return 200 (not 304)
			accGet200, err := cl.AccountGetWithResponse(root, session, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag2)
				return nil
			})
			tests.Ok(t, err, accGet200)
			r.NotNil(accGet200.JSON200, "should return 200 with body after cache invalidation")
			r.Len(accGet200.JSON200.EmailAddresses, 1, "account should be back to one email")
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

			// Get initial account state
			accGet1, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet1)

			etag1 := accGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")

			// Update the profile
			newBio := "This is my new bio"
			updateResp, err := cl.AccountUpdateWithResponse(root, openapi.AccountMutableProps{
				Bio: &newBio,
			}, session)
			tests.Ok(t, err, updateResp)

			accGet2, err := cl.AccountGetWithResponse(root, session)
			tests.Ok(t, err, accGet2)
			a.Contains(accGet2.JSON200.Bio, newBio, "bio should contain the updated text")

			etag2 := accGet2.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag2, "ETag header should be present")
			a.NotEqual(etag1, etag2, "ETag should change after profile update")

			// Verify 304 with current ETag
			accGet304, err := cl.AccountGetWithResponse(root, session, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag2)
				return nil
			})
			tests.Status(t, err, accGet304, 304)

			// Update again
			newBio2 := "This is my second bio"
			updateResp2, err := cl.AccountUpdateWithResponse(root, openapi.AccountMutableProps{
				Bio: &newBio2,
			}, session)
			tests.Ok(t, err, updateResp2)

			// Conditional GET with old ETag should now return 200 (not 304)
			accGet200, err := cl.AccountGetWithResponse(root, session, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag2)
				return nil
			})
			tests.Ok(t, err, accGet200)
			r.NotNil(accGet200.JSON200, "should return 200 with body after cache invalidation")
			a.Contains(accGet200.JSON200.Bio, newBio2, "bio should contain the second updated text")
		}))
	}))
}
