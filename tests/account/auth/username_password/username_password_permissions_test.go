package username_password_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestUsernamePasswordAuthMemberPermissions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("username_password_user_has_member_permissions_without_email", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				catName := "Category " + uuid.NewString()
				cat := tests.AssertRequest(
					cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
						Colour:      "#fe4efd",
						Description: "test category",
						Name:        catName,
					}, adminSession),
				)(t, http.StatusOK)

				handle := xid.New().String()
				password := "password123"

				signup, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{Identifier: handle, Token: password})
				tests.Ok(t, err, signup)

				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
				userCtx := e2e.WithAccountID(root, accountID)
				userSession := sh.WithSession(userCtx)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)
				r.Equal(openapi.AccountVerifiedStatusNone, userAccount.JSON200.VerifiedStatus)
				r.Len(userAccount.JSON200.EmailAddresses, 0)

				threadCreate := tests.AssertRequest(
					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
						Body:       opt.New("<p>test thread from username user</p>").Ptr(),
						Category:   opt.New(cat.JSON200.Id).Ptr(),
						Visibility: opt.New(openapi.Published).Ptr(),
						Title:      "Username user thread",
					}, userSession),
				)(t, http.StatusOK)
				a.Equal("Username user thread", threadCreate.JSON200.Title)

				reactCreate := tests.AssertRequest(
					cl.PostReactAddWithResponse(root, threadCreate.JSON200.Id, openapi.PostReactAddJSONRequestBody{Emoji: "👍"}, userSession),
				)(t, http.StatusOK)
				a.Equal("👍", reactCreate.JSON200.Emoji)
			})
		}))
	}))
}
