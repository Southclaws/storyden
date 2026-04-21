package account_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/tests"
)

func TestMemberCanViewOwnWarningsButCannotModify(t *testing.T) {
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

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)

			reason := "please cool it down"
			createResp, err := cl.AccountWarningCreateWithResponse(
				root,
				member.ID.String(),
				openapi.AccountWarningCreateJSONRequestBody{Reason: reason},
				adminSession,
			)
			tests.Ok(t, err, createResp)
			warningID := createResp.JSON200.Id

			listResp, err := cl.AccountWarningListWithResponse(root, member.ID.String(), memberSession)
			tests.Ok(t, err, listResp)
			r.Equal(1, listResp.JSON200.Total)
			r.Len(listResp.JSON200.Warnings, 1)
			a.Equal(reason, listResp.JSON200.Warnings[0].Reason)
			a.Equal(createResp.JSON200.IssuedAt, listResp.JSON200.Warnings[0].IssuedAt)

			updateResp, err := cl.AccountWarningUpdateWithResponse(
				root,
				member.ID.String(),
				warningID,
				openapi.AccountWarningUpdateJSONRequestBody{Reason: "not allowed"},
				memberSession,
			)
			tests.Status(t, err, updateResp, http.StatusForbidden)

			deleteResp, err := cl.AccountWarningDeleteWithResponse(root, member.ID.String(), warningID, memberSession)
			tests.Status(t, err, deleteResp, http.StatusForbidden)

			listAfterResp, err := cl.AccountWarningListWithResponse(root, member.ID.String(), memberSession)
			tests.Ok(t, err, listAfterResp)
			r.Equal(1, listAfterResp.JSON200.Total)
			a.Equal(reason, listAfterResp.JSON200.Warnings[0].Reason)
		}))
	}))
}

func TestWarningUpdateRequiresMatchingAccountAndWarning(t *testing.T) {
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

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			firstMemberCtx, firstMember := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			firstMemberSession := sh.WithSession(firstMemberCtx)

			_, secondMember := e2e.WithAccount(root, aw, seed.Account_005_Þórr)

			createResp, err := cl.AccountWarningCreateWithResponse(
				root,
				firstMember.ID.String(),
				openapi.AccountWarningCreateJSONRequestBody{Reason: "original warning reason"},
				adminSession,
			)
			tests.Ok(t, err, createResp)
			warningID := createResp.JSON200.Id

			updateResp, err := cl.AccountWarningUpdateWithResponse(
				root,
				secondMember.ID.String(),
				warningID,
				openapi.AccountWarningUpdateJSONRequestBody{Reason: "wrong account update"},
				adminSession,
			)
			tests.Status(t, err, updateResp, http.StatusNotFound)

			listResp, err := cl.AccountWarningListWithResponse(root, firstMember.ID.String(), firstMemberSession)
			tests.Ok(t, err, listResp)
			r.Equal(1, listResp.JSON200.Total)
			r.Len(listResp.JSON200.Warnings, 1)
			a.Equal("original warning reason", listResp.JSON200.Warnings[0].Reason)
			a.Equal(warningID, listResp.JSON200.Warnings[0].Id)
		}))
	}))
}

func TestWarningIssueCreatesNotification(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)

			createResp, err := cl.AccountWarningCreateWithResponse(
				root,
				member.ID.String(),
				openapi.AccountWarningCreateJSONRequestBody{Reason: "notification check"},
				adminSession,
			)
			tests.Ok(t, err, createResp)

			var warningNotification openapi.Notification
			require.Eventually(t, func() bool {
				notList, listErr := cl.NotificationListWithResponse(
					root,
					&openapi.NotificationListParams{},
					memberSession,
				)
				if listErr != nil || notList == nil || notList.JSON200 == nil {
					return false
				}

				for _, n := range notList.JSON200.Notifications {
					if n.Event == openapi.WarningIssued {
						warningNotification = n
						return true
					}
				}

				return false
			}, 5*time.Second, 100*time.Millisecond, "expected warning notification for member")

			a.Equal(openapi.WarningIssued, warningNotification.Event)
			if a.NotNil(warningNotification.Source) {
				a.Equal(admin.ID.String(), warningNotification.Source.Id)
			}
			a.Equal(openapi.Unread, warningNotification.Status)
		}))
	}))
}

func TestWarningDeleteEventIncludesActor(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) error {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			_ = sh.WithSession(memberCtx)

			events := make(chan rpc.EventAccountWarningDeleted, 1)
			sub, err := pubsub.Subscribe(ctx, bus, "tests.warning_delete_actor", func(ctx context.Context, event *rpc.EventAccountWarningDeleted) error {
				events <- *event
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(sub.Close())
			})

			createResp, err := cl.AccountWarningCreateWithResponse(
				root,
				member.ID.String(),
				openapi.AccountWarningCreateJSONRequestBody{Reason: "delete actor attribution"},
				adminSession,
			)
			tests.Ok(t, err, createResp)

			deleteResp, err := cl.AccountWarningDeleteWithResponse(
				root,
				member.ID.String(),
				createResp.JSON200.Id,
				adminSession,
			)
			tests.Status(t, err, deleteResp, http.StatusNoContent)

			select {
			case event := <-events:
				a.Equal(member.ID, event.AccountID)
				a.Equal(string(createResp.JSON200.Id), event.WarningID)
				a.Equal(admin.ID, event.AuthorID)
			case <-time.After(2 * time.Second):
				t.Fatal("timed out waiting for warning delete event")
			}
			return nil
		}))
	}))
}
