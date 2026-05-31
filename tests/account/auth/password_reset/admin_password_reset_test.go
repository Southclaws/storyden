package password_reset_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/audit/audit_logger"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAdminPasswordResetTokenGet(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("admin_generates_token_for_user", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				email := xid.New().String() + "@storyden.org"
				password := "originalsecretpassword"

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				tokenResp, err := cl.AccountPasswordResetTokenGetWithResponse(root, userAccount.JSON200.Id, adminSession)
				tests.Ok(t, err, tokenResp)

				a.NotEmpty(tokenResp.JSON200.Token)

				newPassword := "newsuperpassword"
				resetResp, err := cl.AuthPasswordResetWithResponse(root, openapi.AuthPasswordResetJSONRequestBody{
					Token: tokenResp.JSON200.Token,
					New:   newPassword,
				})
				tests.Ok(t, err, resetResp)

				r.Equal(userAccount.JSON200.Id, resetResp.JSON200.Id)

				loginResp, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email,
					Password: newPassword,
				})
				tests.Ok(t, err, loginResp)

				r.Equal(userAccount.JSON200.Id, loginResp.JSON200.Id)
			})

			t.Run("non_admin_cannot_get_token", func(t *testing.T) {
				r := require.New(t)

				email := xid.New().String() + "@storyden.org"
				password := "somepassword"

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				otherUserEmail := xid.New().String() + "@storyden.org"
				otherUserSignup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    otherUserEmail,
					Password: "password",
				})
				tests.Ok(t, err, otherUserSignup)
				otherUserSession := e2e.WithSessionFromHeader(t, root, otherUserSignup.HTTPResponse.Header)

				tokenResp, err := cl.AccountPasswordResetTokenGetWithResponse(root, userAccount.JSON200.Id, otherUserSession)
				r.NoError(err)
				r.Equal(http.StatusForbidden, tokenResp.StatusCode())
			})

			t.Run("account_not_found", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				nonexistentID := xid.New().String()
				tokenResp, err := cl.AccountPasswordResetTokenGetWithResponse(root, nonexistentID, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNotFound, tokenResp.StatusCode())
			})

		}))
	}))
}

func TestAdminPasswordResetEmail(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			t.Run("admin_sends_reset_email_for_user", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				email := xid.New().String() + "@storyden.org"
				password := "originalsecretpassword"
				signupEmailCount := inbox.Count()

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				r.Len(userAccount.JSON200.EmailAddresses, 1)
				emailAddressID := userAccount.JSON200.EmailAddresses[0].Id

				tests.WaitForNextEmail(t, inbox, signupEmailCount)
				resetEmailCount := inbox.Count()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, userAccount.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: emailAddressID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNoContent, resetReq.StatusCode())

				resetEmail := tests.WaitForNextEmail(t, inbox, resetEmailCount)
				a.Equal(userAccount.JSON200.Name, resetEmail.Name)
				a.Equal(email, resetEmail.Address.Address)
				token := regexp.MustCompile(`\?token=(.+)`).FindStringSubmatch(resetEmail.Plain)[1]

				newPassword := "newsuperpassword"
				resetResp, err := cl.AuthPasswordResetWithResponse(root, openapi.AuthPasswordResetJSONRequestBody{
					Token: token,
					New:   newPassword,
				})
				tests.Ok(t, err, resetResp)

				r.Equal(userAccount.JSON200.Id, resetResp.JSON200.Id)

				loginResp, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email,
					Password: newPassword,
				})
				tests.Ok(t, err, loginResp)

				r.Equal(userAccount.JSON200.Id, loginResp.JSON200.Id)
			})

			t.Run("admin_sends_to_specific_email_when_multiple_exist", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				email1 := xid.New().String() + "@storyden.org"
				email2 := xid.New().String() + "@example.com"
				password := "password123"
				signupEmailCount := inbox.Count()

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email1,
					Password: password,
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				tests.WaitForNextEmail(t, inbox, signupEmailCount)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				addEmailResp, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailAddJSONRequestBody{
					EmailAddress: email2,
				}, userSession)
				tests.Ok(t, err, addEmailResp)

				userAccount, err = cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				r.Len(userAccount.JSON200.EmailAddresses, 2)

				var email2AddressID openapi.Identifier
				for _, e := range userAccount.JSON200.EmailAddresses {
					if e.EmailAddress == email2 {
						email2AddressID = e.Id
						break
					}
				}
				r.NotEmpty(email2AddressID)

				resetEmailCount := inbox.Count()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, userAccount.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: email2AddressID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNoContent, resetReq.StatusCode())

				resetEmail := tests.WaitForNextEmail(t, inbox, resetEmailCount)
				a.Equal(email2, resetEmail.Address.Address)
			})

			t.Run("email_address_not_found", func(t *testing.T) {
				r := require.New(t)

				email := xid.New().String() + "@storyden.org"
				password := "password"

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				nonexistentEmailID := xid.New().String()
				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, userAccount.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: nonexistentEmailID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNotFound, resetReq.StatusCode())
			})

			t.Run("email_belongs_to_different_account", func(t *testing.T) {
				r := require.New(t)

				user1Email := xid.New().String() + "@storyden.org"
				user2Email := xid.New().String() + "@example.com"

				signup1, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    user1Email,
					Password: "password",
				})
				tests.Ok(t, err, signup1)
				user1Session := e2e.WithSessionFromHeader(t, root, signup1.HTTPResponse.Header)

				user1Account, err := cl.AccountGetWithResponse(root, user1Session)
				tests.Ok(t, err, user1Account)

				signup2, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    user2Email,
					Password: "password",
				})
				tests.Ok(t, err, signup2)
				user2Session := e2e.WithSessionFromHeader(t, root, signup2.HTTPResponse.Header)

				user2Account, err := cl.AccountGetWithResponse(root, user2Session)
				tests.Ok(t, err, user2Account)

				user2EmailID := user2Account.JSON200.EmailAddresses[0].Id

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, user1Account.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: user2EmailID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNotFound, resetReq.StatusCode())
			})

			t.Run("non_admin_cannot_send_email", func(t *testing.T) {
				r := require.New(t)

				user1Email := xid.New().String() + "@storyden.org"
				user2Email := xid.New().String() + "@example.com"

				signup1, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    user1Email,
					Password: "password",
				})
				tests.Ok(t, err, signup1)
				user1Session := e2e.WithSessionFromHeader(t, root, signup1.HTTPResponse.Header)

				user1Account, err := cl.AccountGetWithResponse(root, user1Session)
				tests.Ok(t, err, user1Account)

				signup2, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    user2Email,
					Password: "password",
				})
				tests.Ok(t, err, signup2)
				user2Session := e2e.WithSessionFromHeader(t, root, signup2.HTTPResponse.Header)

				user2Account, err := cl.AccountGetWithResponse(root, user2Session)
				tests.Ok(t, err, user2Account)

				user2EmailID := user2Account.JSON200.EmailAddresses[0].Id

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, user2Account.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: user2EmailID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, user1Session)
				r.NoError(err)
				r.Equal(http.StatusForbidden, resetReq.StatusCode())
			})

			t.Run("account_not_found", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				nonexistentAccountID := xid.New().String()
				fakeEmailID := xid.New().String()

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, nonexistentAccountID, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: fakeEmailID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				r.NoError(err)
				r.Equal(http.StatusNotFound, resetReq.StatusCode())
			})
		}))
	}))
}

func TestAdminPasswordResetAuditLogging(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		_ *audit_logger.Service,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		var adminCtx context.Context
		var adminSession openapi.RequestEditorFn

		lc.Append(fx.StartHook(func() {
			adminCtx, _ = e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession = sh.WithSession(adminCtx)
		}))

		lc.Append(fx.StartHook(func() {
			t.Run("audit_token_issued_event", func(t *testing.T) {
				a := assert.New(t)

				email := xid.New().String() + "@storyden.org"
				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: "password",
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				tokenResp, err := cl.AccountPasswordResetTokenGetWithResponse(root, userAccount.JSON200.Id, adminSession)
				tests.Ok(t, err, tokenResp)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				event, found := lo.Find(*list.JSON200.Events, func(e openapi.AuditEvent) bool {
					if e.Type != openapi.AccountPasswordResetTokenIssued {
						return false
					}
					eventIssued, err := e.AsAuditEventAccountPasswordResetTokenIssued()
					if err != nil {
						return false
					}
					return eventIssued.AccountId == userAccount.JSON200.Id
				})

				a.True(found, "Should find account_password_reset_token_issued event")
				eventIssued, err := event.AsAuditEventAccountPasswordResetTokenIssued()
				a.NoError(err)
				a.Equal(openapi.AccountPasswordResetTokenIssued, eventIssued.Type)
				a.Equal(userAccount.JSON200.Id, eventIssued.AccountId)
			})

			t.Run("audit_email_sent_event", func(t *testing.T) {
				a := assert.New(t)

				email := xid.New().String() + "@storyden.org"
				signupEmailCount := inbox.Count()

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: "password",
				})
				tests.Ok(t, err, signup)
				userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				tests.WaitForNextEmail(t, inbox, signupEmailCount)

				userAccount, err := cl.AccountGetWithResponse(root, userSession)
				tests.Ok(t, err, userAccount)

				emailAddressID := userAccount.JSON200.EmailAddresses[0].Id

				resetReq, err := cl.AccountEmailPasswordResetWithResponse(root, userAccount.JSON200.Id, openapi.AccountEmailPasswordResetJSONRequestBody{
					EmailAddressId: emailAddressID,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				}, adminSession)
				a.NoError(err)
				a.Equal(http.StatusNoContent, resetReq.StatusCode())

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				event, found := lo.Find(*list.JSON200.Events, func(e openapi.AuditEvent) bool {
					if e.Type != openapi.AccountPasswordResetEmailSent {
						return false
					}
					eventSent, err := e.AsAuditEventAccountPasswordResetEmailSent()
					if err != nil {
						return false
					}
					return eventSent.AccountId == userAccount.JSON200.Id &&
						eventSent.EmailAddressId == emailAddressID
				})

				a.True(found, "Should find account_password_reset_email_sent event")
				eventSent, err := event.AsAuditEventAccountPasswordResetEmailSent()
				a.NoError(err)
				a.Equal(openapi.AccountPasswordResetEmailSent, eventSent.Type)
				a.Equal(userAccount.JSON200.Id, eventSent.AccountId)
				a.Equal(emailAddressID, eventSent.EmailAddressId)
			})
		}))
	}))
}
