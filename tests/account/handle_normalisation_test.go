package account_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestHandleNormalisation(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("uppercase_handle_normalised_to_lowercase_on_signup", func(t *testing.T) {
				suffix := xid.New().String()
				// Submit handle with mixed case — should be silently lowercased.
				resp := tests.AssertRequest(
					cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
						Identifier: "UpperCase-" + suffix,
						Token:      "password1",
					}),
				)(t, http.StatusOK)
				require.NotNil(t, resp.JSON200)

				accID := account.AccountID(utils.Must(xid.FromString(resp.JSON200.Id)))
				session := sh.WithSession(e2e.WithAccountID(root, accID))

				got := tests.AssertRequest(
					cl.AccountGetWithResponse(root, session),
				)(t, http.StatusOK)
				require.NotNil(t, got.JSON200)
				assert.Equal(t, "uppercase-"+suffix, got.JSON200.Handle)
			})

			t.Run("duplicate_handle_case_insensitive_rejected_on_signup", func(t *testing.T) {
				suffix := xid.New().String()
				lowercase := "lower-" + suffix
				uppercase := "Lower-" + suffix

				// First registration with lowercase handle — must succeed.
				first := tests.AssertRequest(
					cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
						Identifier: lowercase,
						Token:      "password1",
					}),
				)(t, http.StatusOK)
				require.NotNil(t, first.JSON200)

				// Second registration with same handle but mixed case — must be rejected.
				second, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: uppercase,
					Token:      "password2",
				})
				tests.Status(t, err, second, http.StatusConflict)
			})

			t.Run("invalid_handle_characters_rejected_not_silently_mangled", func(t *testing.T) {
				for _, bad := range []string{"-leading", "trailing-", "has space", "!!!"} {
					resp, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
						Identifier: bad,
						Token:      "password1",
					})
					tests.Status(t, err, resp, http.StatusBadRequest)
				}
			})

			t.Run("uppercase_handle_normalised_on_profile_update", func(t *testing.T) {
				suffix := xid.New().String()
				resp := tests.AssertRequest(
					cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
						Identifier: "before-" + suffix,
						Token:      "password1",
					}),
				)(t, http.StatusOK)
				require.NotNil(t, resp.JSON200)

				accID := account.AccountID(utils.Must(xid.FromString(resp.JSON200.Id)))
				session := sh.WithSession(e2e.WithAccountID(root, accID))

				newHandle := openapi.AccountHandle("After-" + suffix)
				updated := tests.AssertRequest(
					cl.AccountUpdateWithResponse(root, openapi.AccountMutableProps{
						Handle: &newHandle,
					}, session),
				)(t, http.StatusOK)
				require.NotNil(t, updated.JSON200)
				assert.Equal(t, "after-"+suffix, updated.JSON200.Handle)
			})
		}))
	}))
}
