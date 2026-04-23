package admin_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/comms/mailqueue"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	ent_emailqueue "github.com/Southclaws/storyden/internal/ent/emailqueue"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

const (
	welcomeSubject       = "Welcome to Storyden!"
	passwordResetSubject = "Reset your password on Storyden!"
	testPassword         = "mysupersecretpasswordwhichissosecretiforgotwhatitwas"
	emailQueueWait       = 15 * time.Second
	emailQueuePoll       = 100 * time.Millisecond
)

func TestEmailQueueList(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		mailSender mailer.Sender,
		db *ent.Client,
		_ *mailqueue.Queuer,
	) {
		inbox := mailSender.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			t.Run("lists_sent_signup_email", func(t *testing.T) {
				inbox.Reset()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				address := uniqueEmailAddress()

				signupEmailOnly(t, root, cl, address)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.RecipientAddress == address &&
						row.Subject == welcomeSubject &&
						row.Status == ent_emailqueue.StatusSent
				})

				email := waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == address &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusSent
				})

				require.NotNil(t, email.ProcessedAt)
				require.Len(t, email.Attempts, 1)
				require.Equal(t, openapi.EmailQueueAttemptStatusSent, email.Attempts[0].Status)
				require.Nil(t, email.Attempts[0].Error)
				require.Equal(t, 1, inbox.Count())

				last := inbox.GetLast()
				require.Equal(t, address, last.Address.Address)
				require.Equal(t, email.RecipientName, last.Name)
				require.Equal(t, welcomeSubject, last.Subject)
			})

			t.Run("lists_failed_password_reset_email", func(t *testing.T) {
				inbox.Reset()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				address := uniqueEmailAddress()

				signupEmailPassword(t, root, cl, address)

				welcomeEmail := waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == address &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusSent
				})
				require.NotNil(t, welcomeEmail.ProcessedAt)
				require.Equal(t, 1, inbox.Count())

				sendErr := errors.New("simulated mail transport failure")
				inbox.SetSendError(sendErr)
				t.Cleanup(inbox.ClearSendError)

				requestPasswordReset(t, root, cl, address)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.RecipientAddress == address &&
						row.Subject == passwordResetSubject &&
						row.Status == ent_emailqueue.StatusFailed
				})

				email := waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == address &&
						item.Subject == passwordResetSubject &&
						item.Status == openapi.EmailQueueStatusFailed
				})

				require.Nil(t, email.ProcessedAt)
				require.NotEmpty(t, email.Attempts)

				attempt := email.Attempts[len(email.Attempts)-1]
				require.Equal(t, openapi.EmailQueueAttemptStatusFailed, attempt.Status)
				require.NotNil(t, attempt.Error)
				require.Contains(t, *attempt.Error, sendErr.Error())
				require.Equal(t, 1, inbox.Count())
			})

			t.Run("lists_newest_first", func(t *testing.T) {
				inbox.Reset()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				olderAddress := uniqueEmailAddress()
				newerAddress := uniqueEmailAddress()

				signupEmailOnly(t, root, cl, olderAddress)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.RecipientAddress == olderAddress &&
						row.Subject == welcomeSubject &&
						row.Status == ent_emailqueue.StatusSent
				})
				waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == olderAddress &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusSent
				})

				sendErr := errors.New("simulated mail transport failure for ordering test")
				inbox.SetSendError(sendErr)
				t.Cleanup(inbox.ClearSendError)

				signupEmailOnly(t, root, cl, newerAddress)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.RecipientAddress == newerAddress &&
						row.Subject == welcomeSubject &&
						row.Status == ent_emailqueue.StatusFailed
				})
				waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == newerAddress &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusFailed
				})

				list := waitForEmailQueueList(t, adminCtx, cl, adminSession, func(items []openapi.EmailQueueItem) bool {
					return indexOfEmail(items, func(item openapi.EmailQueueItem) bool {
						return string(item.RecipientAddress) == olderAddress &&
							item.Subject == welcomeSubject &&
							item.Status == openapi.EmailQueueStatusSent
					}) >= 0 &&
						indexOfEmail(items, func(item openapi.EmailQueueItem) bool {
							return string(item.RecipientAddress) == newerAddress &&
								item.Subject == welcomeSubject &&
								item.Status == openapi.EmailQueueStatusFailed
						}) >= 0
				})

				olderIndex := indexOfEmail(list, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == olderAddress &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusSent
				})
				newerIndex := indexOfEmail(list, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == newerAddress &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusFailed
				})

				require.GreaterOrEqual(t, olderIndex, 0)
				require.GreaterOrEqual(t, newerIndex, 0)
				require.Less(t, newerIndex, olderIndex)
			})

			t.Run("requires_admin_permission", func(t *testing.T) {
				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				list, err := cl.EmailQueueListWithResponse(memberCtx, &openapi.EmailQueueListParams{}, memberSession)
				require.NoError(t, err)
				require.Equal(t, http.StatusForbidden, list.StatusCode())
			})

			t.Run("retries_failed_email_now", func(t *testing.T) {
				inbox.Reset()

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				address := uniqueEmailAddress()

				inbox.SetSendError(errors.New("simulated retryable failure"))
				t.Cleanup(inbox.ClearSendError)

				signupEmailOnly(t, root, cl, address)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.RecipientAddress == address &&
						row.Subject == welcomeSubject &&
						row.Status == ent_emailqueue.StatusFailed
				})

				email := waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return string(item.RecipientAddress) == address &&
						item.Subject == welcomeSubject &&
						item.Status == openapi.EmailQueueStatusFailed
				})

				emailXID, err := xid.FromString(string(email.Id))
				require.NoError(t, err)

				rowBefore, err := db.EmailQueue.Get(adminCtx, emailXID)
				require.NoError(t, err)
				require.True(t, rowBefore.AvailableAt.After(time.Now()))
				require.Len(t, rowBefore.Attempts, 1)

				retry, err := cl.EmailQueueRetryWithResponse(adminCtx, email.Id, adminSession)
				tests.Ok(t, err, retry)
				require.NotNil(t, retry.JSON200)
				require.Equal(t, email.Id, retry.JSON200.Id)
				require.Equal(t, openapi.EmailQueueStatusPending, retry.JSON200.Status)
				waitForEmailQueueRow(t, adminCtx, db, func(row *ent.EmailQueue) bool {
					return row.ID == emailXID &&
						row.Status == ent_emailqueue.StatusFailed &&
						len(row.Attempts) == 2
				})

				retriedEmail := waitForEmailQueueItem(t, adminCtx, cl, adminSession, func(item openapi.EmailQueueItem) bool {
					return item.Id == email.Id &&
						item.Status == openapi.EmailQueueStatusFailed &&
						len(item.Attempts) == 2
				})
				require.Len(t, retriedEmail.Attempts, 2)

				rowAfter, err := db.EmailQueue.Get(adminCtx, emailXID)
				require.NoError(t, err)
				require.Len(t, rowAfter.Attempts, 2)
				require.True(t, rowAfter.AvailableAt.After(time.Now()))
			})
		}))
	}))
}

func signupEmailOnly(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, address string) {
	t.Helper()

	signup, err := cl.AuthEmailSignupWithResponse(ctx, nil, openapi.AuthEmailSignupJSONRequestBody{
		Email: address,
	})
	tests.Ok(t, err, signup)
}

func signupEmailPassword(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, address string) {
	t.Helper()

	signup, err := cl.AuthEmailPasswordSignupWithResponse(ctx, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
		Email:    address,
		Password: testPassword,
	})
	tests.Ok(t, err, signup)
}

func requestPasswordReset(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, address string) {
	t.Helper()

	reset, err := cl.AuthPasswordResetRequestEmailWithResponse(ctx, openapi.AuthEmailPasswordReset{
		Email: address,
		TokenUrl: struct {
			Query string `json:"query"`
			Url   string `json:"url"`
		}{
			Url:   "http://localhost:3000/reset",
			Query: "token",
		},
	})
	tests.Ok(t, err, reset)
}

func waitForEmailQueueItem(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	session openapi.RequestEditorFn,
	match func(openapi.EmailQueueItem) bool,
) openapi.EmailQueueItem {
	t.Helper()

	var found openapi.EmailQueueItem

	waitForEmailQueueList(t, ctx, cl, session, func(items []openapi.EmailQueueItem) bool {
		index := indexOfEmail(items, match)
		if index < 0 {
			return false
		}

		found = items[index]
		return true
	})

	return found
}

func waitForEmailQueueList(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	session openapi.RequestEditorFn,
	match func([]openapi.EmailQueueItem) bool,
) []openapi.EmailQueueItem {
	t.Helper()

	var items []openapi.EmailQueueItem

	require.Eventually(t, func() bool {
		list, err := cl.EmailQueueListWithResponse(ctx, &openapi.EmailQueueListParams{}, session)
		if err != nil || list.StatusCode() != http.StatusOK || list.JSON200 == nil || list.JSON200.Emails == nil {
			return false
		}

		items = *list.JSON200.Emails
		return match(items)
	}, emailQueueWait, emailQueuePoll)

	return items
}

func waitForEmailQueueRow(
	t *testing.T,
	ctx context.Context,
	db *ent.Client,
	match func(*ent.EmailQueue) bool,
) *ent.EmailQueue {
	t.Helper()

	var found *ent.EmailQueue

	require.Eventually(t, func() bool {
		rows, err := db.EmailQueue.Query().All(ctx)
		if err != nil {
			return false
		}

		for _, row := range rows {
			if match(row) {
				found = row
				return true
			}
		}

		return false
	}, emailQueueWait, emailQueuePoll)

	return found
}

func indexOfEmail(items []openapi.EmailQueueItem, match func(openapi.EmailQueueItem) bool) int {
	for i, item := range items {
		if match(item) {
			return i
		}
	}

	return -1
}

func uniqueEmailAddress() string {
	return xid.New().String() + "@storyden.org"
}
