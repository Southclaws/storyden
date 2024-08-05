package account_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestPublicProfiles(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *cookie.Jar,
		ar account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			// Create 5 fresh accounts
			newAccount(t, root, cl, ar, "odin")
			newAccount(t, root, cl, ar, "frigg")
			newAccount(t, root, cl, ar, "baldur")
			newAccount(t, root, cl, ar, "odin2")
			newAccount(t, root, cl, ar, "þórr")

			// Get them all, default params.

			list1, err := cl.ProfileListWithResponse(root, nil)
			r.NoError(err)
			r.NotNil(list1)
			r.Equal(http.StatusOK, list1.StatusCode())

			a.Equal(1, list1.JSON200.CurrentPage)
			a.GreaterOrEqual(len(list1.JSON200.Profiles), 5)

			// Get one specific name - there are two odins

			list2, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Q: opt.New("odin").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list2)
			r.Equal(http.StatusOK, list2.StatusCode())

			a.Equal(1, list1.JSON200.CurrentPage)
			a.Equal(50, list1.JSON200.PageSize)
			a.GreaterOrEqual(len(list2.JSON200.Profiles), 2)

			// Query an invalid page.

			list3, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Page: opt.New("2147483647").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list3)
			r.Equal(http.StatusOK, list3.StatusCode())

			a.Equal(2147483647, list3.JSON200.CurrentPage)
			a.Nil(list3.JSON200.NextPage)
			a.Empty(list3.JSON200.Profiles)
		}))
	}))
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *cookie.Jar,
		accountQuery account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			handle1 := "user-" + xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{handle1, "password"})
			tests.Ok(t, err, acc1)
			session1 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))), cj)

			handle2 := "user-" + xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{handle2, "password"})
			tests.Ok(t, err, acc2)
			session2 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))), cj)

			t.Run("update_profile", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// as guest
				get1, err := cl.ProfileGetWithResponse(root, handle1)
				tests.Ok(t, err, get1)
				r.Equal(handle1, get1.JSON200.Handle)
				r.Equal(handle1, get1.JSON200.Name)

				// as another user
				get2, err := cl.ProfileGetWithResponse(root, handle1, session2)
				tests.Ok(t, err, get2)
				r.Equal(handle1, get2.JSON200.Handle)
				r.Equal(handle1, get2.JSON200.Name)

				// update account profile
				newbio := "new bio"
				newhandle := "newhandle-" + xid.New().String()
				newname := "newname"
				newlinks := []openapi.ProfileExternalLink{
					{Text: "link1", Url: "https://example.com"},
				}
				newmeta := openapi.Metadata{
					"some": "data",
				}

				upd1, err := cl.AccountUpdateWithResponse(root, openapi.AccountUpdateJSONRequestBody{
					Bio:    &newbio,
					Handle: &newhandle,
					Name:   &newname,
					Links:  &newlinks,
					Meta:   &newmeta,
				}, session1)
				tests.Ok(t, err, upd1)

				a.Contains(upd1.JSON200.Bio, newbio)
				a.Equal(newhandle, upd1.JSON200.Handle)
				a.Equal(newname, upd1.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, upd1.JSON200.Links)
				a.Equal(newmeta, upd1.JSON200.Meta)

				// old handle should not work
				getold, err := cl.ProfileGetWithResponse(root, handle1)
				tests.Status(t, err, getold, http.StatusNotFound)

				// as guest
				getAfterUpdateAsGuest, err := cl.ProfileGetWithResponse(root, newhandle)
				tests.Ok(t, err, getAfterUpdateAsGuest)
				a.Contains(getAfterUpdateAsGuest.JSON200.Bio, newbio)
				a.Equal(newhandle, getAfterUpdateAsGuest.JSON200.Handle)
				a.Equal(newname, getAfterUpdateAsGuest.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, getAfterUpdateAsGuest.JSON200.Links)
				a.Equal(newmeta, getAfterUpdateAsGuest.JSON200.Meta)

				// as another user
				getAfterUpdateAsUser2, err := cl.ProfileGetWithResponse(root, newhandle, session2)
				tests.Ok(t, err, getAfterUpdateAsUser2)
				a.Contains(getAfterUpdateAsUser2.JSON200.Bio, newbio)
				a.Equal(newhandle, getAfterUpdateAsUser2.JSON200.Handle)
				a.Equal(newname, getAfterUpdateAsUser2.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, getAfterUpdateAsUser2.JSON200.Links)
				a.Equal(newmeta, getAfterUpdateAsUser2.JSON200.Meta)
			})
		}))
	}))
}

func newAccount(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, ar account_querier.Querier, handle string) account.Account {
	r := require.New(t)

	hand1 := handle + "-" + xid.New().String()

	response, err := cl.AuthPasswordSignupWithResponse(ctx, openapi.AuthPair{
		Identifier: hand1,
		Token:      "password",
	})
	r.NoError(err)
	r.NotNil(response)
	r.Equal(http.StatusOK, response.StatusCode())

	acc, err := ar.GetByID(ctx, account.AccountID(utils.Must(xid.FromString(response.JSON200.Id))))
	r.NoError(err)
	r.NotNil(acc)

	return *acc
}
