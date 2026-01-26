package email_verification_test

// NOTE: These tests require AuthenticationMode to be set, however this is a
// global setting and cannot be easily reset between tests. Will figure out...

// import (
// 	"context"
// 	"net/http"
// 	"regexp"
// 	"testing"

// 	"github.com/Southclaws/opt"
// 	"github.com/google/uuid"
// 	"github.com/rs/xid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/fx"

// 	"github.com/Southclaws/storyden/app/resources/account"
// 	"github.com/Southclaws/storyden/app/resources/account/account_writer"
// 	"github.com/Southclaws/storyden/app/resources/seed"
// 	"github.com/Southclaws/storyden/app/transports/http/openapi"
// 	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
// 	"github.com/Southclaws/storyden/internal/integration"
// 	"github.com/Southclaws/storyden/internal/integration/e2e"
// 	"github.com/Southclaws/storyden/tests"
// )

// func TestUnverifiedUserPermissions(t *testing.T) {
// 	t.Parallel()

// 	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
// 		lc fx.Lifecycle,
// 		root context.Context,
// 		cl *openapi.ClientWithResponses,
// 		sh *e2e.SessionHelper,
// 		aw *account_writer.Writer,
// 		mail mailer.Sender,
// 	) {
// 		inbox := mail.(*mailer.Mock)

// 		lc.Append(fx.StartHook(func() {
// 			// Set authentication mode to Email for these tests
// 			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
// 			adminSession := sh.WithSession(adminCtx)

// 			emailMode := openapi.Email
// 			_, err := cl.AdminSettingsUpdateWithResponse(adminCtx, openapi.AdminSettingsUpdateJSONRequestBody{
// 				AuthenticationMode: &emailMode,
// 			}, adminSession)
// 			if err != nil {
// 				t.Fatalf("Failed to set authentication mode: %v", err)
// 			}

// 			t.Run("unverified_user_has_guest_permissions", func(t *testing.T) {
// 				r := require.New(t)
// 				a := assert.New(t)

// 				catName := "Category " + uuid.NewString()
// 				cat := tests.AssertRequest(
// 					cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
// 						Colour:      "#fe4efd",
// 						Description: "test category",
// 						Name:        catName,
// 					}, adminSession),
// 				)(t, http.StatusOK)

// 				// Sign up with email (unverified)
// 				address := xid.New().String() + "@storyden.org"
// 				signup, err := cl.AuthEmailSignupWithResponse(root, nil, openapi.AuthEmailSignupJSONRequestBody{Email: address})
// 				tests.Ok(t, err, signup)

// 				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
// 				unverifiedCtx := e2e.WithAccountID(root, accountID)
// 				unverifiedSession := sh.WithSession(unverifiedCtx)

// 				// Verify account is unverified
// 				unverifiedAccount, err := cl.AccountGetWithResponse(root, unverifiedSession)
// 				tests.Ok(t, err, unverifiedAccount)
// 				r.Equal(openapi.AccountVerifiedStatusNone, unverifiedAccount.JSON200.VerifiedStatus)

// 				// Try to create a thread - should fail with 403 (guest permissions)
// 				threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
// 					Body:       opt.New("<p>test thread from unverified user</p>").Ptr(),
// 					Category:   opt.New(cat.JSON200.Id).Ptr(),
// 					Visibility: opt.New(openapi.Published).Ptr(),
// 					Title:      "Unverified user thread",
// 				}, unverifiedSession)
// 				r.NoError(err)
// 				a.Equal(http.StatusForbidden, threadCreate.StatusCode(), "unverified user should not be able to create threads")

// 				// Verify email
// 				verification := inbox.GetLast()
// 				match := regexp.MustCompile(`verify your account: ([0-9]{6})`).FindStringSubmatch(verification.Plain)
// 				r.NotNil(match, "verification email should contain a 6-digit code")
// 				r.GreaterOrEqual(len(match), 2, "regex match should have at least 2 elements (full match and capture group)")
// 				code := match[1]
// 				verify, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{Email: address, Code: code}, unverifiedSession)
// 				tests.Ok(t, err, verify)

// 				// Verify account is now verified
// 				verifiedAccount, err := cl.AccountGetWithResponse(root, unverifiedSession)
// 				tests.Ok(t, err, verifiedAccount)
// 				a.Equal(openapi.AccountVerifiedStatusVerifiedEmail, verifiedAccount.JSON200.VerifiedStatus)

// 				// Now creating a thread should succeed
// 				threadCreateAfterVerify := tests.AssertRequest(
// 					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
// 						Body:       opt.New("<p>test thread from verified user</p>").Ptr(),
// 						Category:   opt.New(cat.JSON200.Id).Ptr(),
// 						Visibility: opt.New(openapi.Published).Ptr(),
// 						Title:      "Verified user thread",
// 					}, unverifiedSession),
// 				)(t, http.StatusOK)
// 				a.Equal("Verified user thread", threadCreateAfterVerify.JSON200.Title)
// 			})

// 			t.Run("unverified_user_can_read_public_content", func(t *testing.T) {
// 				r := require.New(t)

// 				// Create content using an admin account
// 				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
// 				adminSession := sh.WithSession(adminCtx)

// 				catName := "ReadTest Category " + uuid.NewString()
// 				cat := tests.AssertRequest(
// 					cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
// 						Colour:      "#123456",
// 						Description: "read test category",
// 						Name:        catName,
// 					}, adminSession),
// 				)(t, http.StatusOK)

// 				thread := tests.AssertRequest(
// 					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
// 						Body:       opt.New("<p>public thread content</p>").Ptr(),
// 						Category:   opt.New(cat.JSON200.Id).Ptr(),
// 						Visibility: opt.New(openapi.Published).Ptr(),
// 						Title:      "Public Thread",
// 					}, adminSession),
// 				)(t, http.StatusOK)

// 				// Sign up with email (unverified)
// 				address := xid.New().String() + "@storyden.org"
// 				signup, err := cl.AuthEmailSignupWithResponse(root, nil, openapi.AuthEmailSignupJSONRequestBody{Email: address})
// 				tests.Ok(t, err, signup)

// 				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
// 				unverifiedCtx := e2e.WithAccountID(root, accountID)
// 				unverifiedSession := sh.WithSession(unverifiedCtx)

// 				// Unverified user should be able to read public threads
// 				threadGet := tests.AssertRequest(
// 					cl.ThreadGetWithResponse(root, thread.JSON200.Slug, nil, unverifiedSession),
// 				)(t, http.StatusOK)
// 				r.Equal("Public Thread", threadGet.JSON200.Title)

// 				// Unverified user should be able to list threads
// 				categoryFilter := []string{cat.JSON200.Slug}
// 				threadList := tests.AssertRequest(
// 					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
// 						Categories: &categoryFilter,
// 					}, unverifiedSession),
// 				)(t, http.StatusOK)
// 				r.NotEmpty(threadList.JSON200.Threads)
// 			})

// 			t.Run("unverified_user_cannot_add_reactions", func(t *testing.T) {
// 				r := require.New(t)
// 				a := assert.New(t)

// 				// Create a thread using an admin account
// 				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
// 				adminSession := sh.WithSession(adminCtx)

// 				catName := "ReactTest Category " + uuid.NewString()
// 				cat := tests.AssertRequest(
// 					cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
// 						Colour:      "#abcdef",
// 						Description: "react test category",
// 						Name:        catName,
// 					}, adminSession),
// 				)(t, http.StatusOK)

// 				thread := tests.AssertRequest(
// 					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
// 						Body:       opt.New("<p>thread for reaction test</p>").Ptr(),
// 						Category:   opt.New(cat.JSON200.Id).Ptr(),
// 						Visibility: opt.New(openapi.Published).Ptr(),
// 						Title:      "Reaction Test Thread",
// 					}, adminSession),
// 				)(t, http.StatusOK)

// 				// Sign up with email (unverified)
// 				address := xid.New().String() + "@storyden.org"
// 				signup, err := cl.AuthEmailSignupWithResponse(root, nil, openapi.AuthEmailSignupJSONRequestBody{Email: address})
// 				tests.Ok(t, err, signup)

// 				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
// 				unverifiedCtx := e2e.WithAccountID(root, accountID)
// 				unverifiedSession := sh.WithSession(unverifiedCtx)

// 				// Unverified user should NOT be able to add reactions (guest permissions)
// 				reactCreate, err := cl.PostReactAddWithResponse(root, thread.JSON200.Id, openapi.PostReactAddJSONRequestBody{Emoji: "❤️"}, unverifiedSession)
// 				r.NoError(err)
// 				a.Equal(http.StatusForbidden, reactCreate.StatusCode(), "unverified user should not be able to add reactions")

// 				// Verify email
// 				verification := inbox.GetLast()
// 				match := regexp.MustCompile(`verify your account: ([0-9]{6})`).FindStringSubmatch(verification.Plain)
// 				r.NotNil(match, "verification email should contain a 6-digit code")
// 				r.GreaterOrEqual(len(match), 2, "regex match should have at least 2 elements (full match and capture group)")
// 				code := match[1]
// 				verify, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{Email: address, Code: code}, unverifiedSession)
// 				tests.Ok(t, err, verify)

// 				// Now adding a reaction should succeed
// 				reactCreateAfterVerify := tests.AssertRequest(
// 					cl.PostReactAddWithResponse(root, thread.JSON200.Id, openapi.PostReactAddJSONRequestBody{Emoji: "❤️"}, unverifiedSession),
// 				)(t, http.StatusOK)
// 				a.Equal("❤️", reactCreateAfterVerify.JSON200.Emoji)
// 			})
// 		}))
// 	}))
// }
