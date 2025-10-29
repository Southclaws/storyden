package account_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestHandleValidation(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("fails_leading_hyphen", func(t *testing.T) {
				leadingHyphen, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: "-abc123",
					Token:      "password",
				})
				tests.Status(t, err, leadingHyphen, http.StatusBadRequest)
			})
			t.Run("fails_trailing_hyphen", func(t *testing.T) {
				trailingHyphen, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: "abc123-",
					Token:      "password",
				})
				tests.Status(t, err, trailingHyphen, http.StatusBadRequest)
			})

			t.Run("fails_uppercase_handle", func(t *testing.T) {
				uppercaseHandle, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: "Hello123",
					Token:      "password",
				})
				tests.Status(t, err, uppercaseHandle, http.StatusBadRequest)
			})

			t.Run("fails_long_handle", func(t *testing.T) {
				longHandle := "this-is-a-very-long-username-over-30-chars"
				tooLong, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: longHandle,
					Token:      "password",
				})
				tests.Status(t, err, tooLong, http.StatusBadRequest)
			})

			t.Run("valid_hyphen", func(t *testing.T) {
				validHyphen, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: "abc-123",
					Token:      "password",
				})
				tests.Ok(t, err, validHyphen)
			})

			t.Run("valid_length_exactly_at_limit", func(t *testing.T) {
				validLength := "abc12345678901234567890123456" // exactly 30 chars
				validLengthResp, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: validLength,
					Token:      "password",
				})
				tests.Ok(t, err, validLengthResp)
			})
		}))
	}))
}

func TestInternationalHandles(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
	) {
		lc.Append(fx.StartHook(func() {
			testCases := []struct {
				name   string
				handle string
			}{
				{"greek", "ελ"},
				{"cyrillic", "док"},
				{"japanese", "日本"},
				{"chinese", "中文"},
				{"hebrew", "עב"},
				{"arabic", "عر"},
				{"hindi", "हि"},
				{"korean", "한국"},
				{"thai", "ไท"},
				{"mixed_with_underscore", "hi_世"},
				{"mixed_with_hyphen", "hi-世"},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					a := assert.New(t)

					uniqueHandle := tc.handle + "-" + xid.New().String()[8:]

					acc, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
						Identifier: uniqueHandle,
						Token:      "password",
					})
					tests.Ok(t, err, acc)
					a.Equal(http.StatusOK, acc.StatusCode())

					profile, err := cl.ProfileGetWithResponse(root, uniqueHandle)
					tests.Ok(t, err, profile)
					a.Equal(uniqueHandle, profile.JSON200.Handle)

					session := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc.JSON200.Id)))))
					newHandle := tc.handle + "x-" + xid.New().String() // Short suffix to stay under 30 chars

					upd, err := cl.AccountUpdateWithResponse(root, openapi.AccountUpdateJSONRequestBody{
						Handle: &newHandle,
					}, session)
					tests.Ok(t, err, upd)
					a.Equal(newHandle, upd.JSON200.Handle)

					oldProfile, err := cl.ProfileGetWithResponse(root, uniqueHandle)
					tests.Status(t, err, oldProfile, http.StatusNotFound)

					newProfile, err := cl.ProfileGetWithResponse(root, newHandle)
					tests.Ok(t, err, newProfile)
					a.Equal(newHandle, newProfile.JSON200.Handle)
				})
			}
		}))
	}))
}
