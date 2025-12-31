package admin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/audit/audit_logger"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAuditEventList(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		_ *audit_logger.Service,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("lists_audit_events_for_admin", func(t *testing.T) {
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				_ = sh.WithSession(memberCtx)

				suspend, err := cl.AdminAccountBanCreateWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, suspend)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)
				a.NotEmpty(list.JSON200.Events)

				// Find the account_suspended event for this specific member
				event, found := lo.Find(*list.JSON200.Events, func(e openapi.AuditEvent) bool {
					if e.Type != openapi.AccountSuspended {
						return false
					}
					eventSuspended, err := e.AsAuditEventAccountSuspended()
					if err != nil {
						return false
					}
					return eventSuspended.AccountId == openapi.Identifier(member.ID.String())
				})

				a.True(found, "Should find account_suspended event for member")
				eventSuspended, err := event.AsAuditEventAccountSuspended()
				a.NoError(err)
				a.Equal(openapi.AccountSuspended, eventSuspended.Type)
				a.Equal(openapi.Identifier(member.ID.String()), eventSuspended.AccountId)
			})

			t.Run("filters_by_event_type", func(t *testing.T) {
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				eventType := openapi.AccountSuspended
				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{
					Types: &[]openapi.AuditEventType{eventType},
				}, adminSession)
				tests.Ok(t, err, list)

				for _, event := range *list.JSON200.Events {
					a.Equal(openapi.AccountSuspended, event.Type)
				}
			})

			t.Run("requires_admin_permission", func(t *testing.T) {
				r := require.New(t)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(memberCtx, &openapi.AuditEventListParams{}, memberSession)
				r.NoError(err)
				r.Equal(http.StatusForbidden, list.StatusCode())
			})
		}))
	}))
}

func TestAuditEventGet(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		_ *audit_logger.Service,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("gets_audit_event_by_id", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				_ = sh.WithSession(memberCtx)

				suspend, err := cl.AdminAccountBanCreateWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, suspend)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)
				r.NotEmpty(list.JSON200.Events)

				// Find the account_suspended event for this specific member
				event, found := lo.Find(*list.JSON200.Events, func(e openapi.AuditEvent) bool {
					if e.Type != openapi.AccountSuspended {
						return false
					}
					eventSuspended, err := e.AsAuditEventAccountSuspended()
					if err != nil {
						return false
					}
					return eventSuspended.AccountId == openapi.Identifier(member.ID.String())
				})
				r.True(found, "Should find account_suspended event for member")

				get, err := cl.AuditEventGetWithResponse(adminCtx, string(event.Id), adminSession)
				tests.Ok(t, err, get)
				a.Equal(event.Id, get.JSON200.Id)
			})

			t.Run("requires_admin_permission", func(t *testing.T) {
				r := require.New(t)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				get, err := cl.AuditEventGetWithResponse(memberCtx, xid.New().String(), memberSession)
				r.NoError(err)
				r.Equal(http.StatusForbidden, get.StatusCode())
			})
		}))
	}))
}

func TestAuditLogging(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		_ *audit_logger.Service,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("logs_thread_deletion", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				vis := openapi.Published
				createThread, err := cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
					Title:      "Test Thread",
					Body:       opt.New("<p>Test content</p>").Ptr(),
					Visibility: &vis,
				}, memberSession)
				tests.Ok(t, err, createThread)

				deleteResp, err := cl.ThreadDeleteWithResponse(adminCtx, createThread.JSON200.Slug, adminSession)
				r.NoError(err)
				a.Equal(http.StatusOK, deleteResp.StatusCode())

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				var found bool
				for _, event := range *list.JSON200.Events {
					if event.Type == openapi.ThreadDeleted {
						deletedEvent, err := event.AsAuditEventThreadDeleted()
						r.NoError(err)
						a.Equal(createThread.JSON200.Id, deletedEvent.ThreadId)
						found = true
						break
					}
				}
				a.True(found, "Should find thread_deleted event in audit log")
			})

			t.Run("logs_reply_deletion", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				vis := openapi.Published
				createThread, err := cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
					Title:      "Test Thread for Replies",
					Body:       opt.New("<p>Test content</p>").Ptr(),
					Visibility: &vis,
				}, memberSession)
				tests.Ok(t, err, createThread)

				replyResp, err := cl.ReplyCreateWithResponse(memberCtx, createThread.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "Test reply",
				}, memberSession)
				tests.Ok(t, err, replyResp)

				deleteResp, err := cl.PostDeleteWithResponse(adminCtx, replyResp.JSON200.Id, adminSession)
				r.NoError(err)
				a.Equal(http.StatusOK, deleteResp.StatusCode())

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				var found bool
				for _, event := range *list.JSON200.Events {
					if event.Type == openapi.ThreadReplyDeleted {
						deletedEvent, err := event.AsAuditEventThreadReplyDeleted()
						r.NoError(err)
						a.Equal(replyResp.JSON200.Id, deletedEvent.ReplyId)
						found = true
						break
					}
				}
				a.True(found, "Should find thread_reply_deleted event in audit log")
			})

			t.Run("logs_account_suspension", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				_ = sh.WithSession(memberCtx)

				suspend, err := cl.AdminAccountBanCreateWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, suspend)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				var found bool
				for _, event := range *list.JSON200.Events {
					if event.Type == openapi.AccountSuspended {
						suspendedEvent, err := event.AsAuditEventAccountSuspended()
						r.NoError(err)
						a.Equal(openapi.Identifier(member.ID.String()), suspendedEvent.AccountId)
						found = true
						break
					}
				}
				a.True(found, "Should find account_suspended event in audit log")
			})

			t.Run("logs_account_unsuspension", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				_ = sh.WithSession(memberCtx)

				suspend, err := cl.AdminAccountBanCreateWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, suspend)

				reinstate, err := cl.AdminAccountBanRemoveWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, reinstate)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				var found bool
				for _, event := range *list.JSON200.Events {
					if event.Type == openapi.AccountUnsuspended {
						unsuspendedEvent, err := event.AsAuditEventAccountUnsuspended()
						r.NoError(err)
						a.Equal(openapi.Identifier(member.ID.String()), unsuspendedEvent.AccountId)
						found = true
						break
					}
				}
				a.True(found, "Should find account_unsuspended event in audit log")
			})
		}))
	}))
}
