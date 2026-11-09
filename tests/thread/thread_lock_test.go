package thread_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadLock(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			authorCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			managerCtx, manager := e2e.WithAccount(root, aw, seed.Account_006_Freyja)
			authorSession := sh.WithSession(authorCtx)
			otherSession := sh.WithSession(otherCtx)
			adminSession := sh.WithSession(adminCtx)
			managerSession := sh.WithSession(managerCtx)

			role := tests.AssertRequest(cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
				Name:        "Post manager " + xid.New().String(),
				Colour:      "green",
				Permissions: openapi.PermissionList{openapi.MANAGEPOSTS},
			}, adminSession))(t, http.StatusOK)
			tests.AssertRequest(cl.AccountAddRoleWithResponse(root, manager.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)

			createThread := func(t *testing.T) *openapi.Thread {
				t.Helper()
				return tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Lock " + xid.New().String(),
					Body:       opt.New("<p>A discussion</p>").Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, authorSession))(t, http.StatusOK).JSON200
			}

			t.Run("only_manage_posts_can_create_lock_or_unlock", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				r.Nil(thread.LockedAt)

				for _, locked := range []bool{true, false} {
					partial := openapi.ThreadMutableProps{Locked: &locked}
					tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial, authorSession))(t, http.StatusForbidden)
					tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial, otherSession))(t, http.StatusForbidden)
					tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial))(t, http.StatusUnauthorized)
					tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
						Title:  "Restricted lock " + xid.New().String(),
						Locked: &locked,
					}, authorSession))(t, http.StatusForbidden)
				}

				locked := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					Locked: opt.New(true).Ptr(),
				}, managerSession))(t, http.StatusOK).JSON200
				r.NotNil(locked.LockedAt)
				r.False(locked.LockedAt.IsZero())

				edited := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					Title: opt.New(thread.Title + " edited").Ptr(),
				}, authorSession))(t, http.StatusOK).JSON200
				r.Equal(locked.LockedAt, edited.LockedAt)

				unlocked := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					Locked: opt.New(false).Ptr(),
				}, adminSession))(t, http.StatusOK).JSON200
				r.Nil(unlocked.LockedAt)

				createdLocked := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:  "Announcement " + xid.New().String(),
					Locked: opt.New(true).Ptr(),
				}, managerSession))(t, http.StatusOK).JSON200
				r.NotNil(createdLocked.LockedAt)
			})

			t.Run("lock_is_independent_of_resolution_and_enforced_for_all_replies", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				first := tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, openapi.ReplyInitialProps{
					Body: "<p>An answer</p>",
				}, otherSession))(t, http.StatusOK).JSON200

				locked := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					Locked:         opt.New(true).Ptr(),
					LifecycleState: opt.New(openapi.ThreadLifecycleStateResolved).Ptr(),
				}, managerSession))(t, http.StatusOK).JSON200
				r.NotNil(locked.LockedAt)
				r.Equal(openapi.ThreadLifecycleStateResolved, locked.LifecycleState)

				for _, replyTo := range []*string{nil, &first.Id} {
					partial := openapi.ReplyInitialProps{Body: "<p>A reply</p>", ReplyTo: replyTo}
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, partial, authorSession))(t, http.StatusForbidden)
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, partial, otherSession))(t, http.StatusForbidden)
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, partial))(t, http.StatusUnauthorized)
				}

				afterRejected := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil))(t, http.StatusOK).JSON200
				r.Equal(1, afterRejected.ReplyStatus.Replies)
				r.Equal(locked.LastReplyAt, afterRejected.LastReplyAt)
				r.Equal(locked.LockedAt, afterRejected.LockedAt)

				for _, replyTo := range []*string{nil, &first.Id} {
					partial := openapi.ReplyInitialProps{Body: "<p>A moderator update</p>", ReplyTo: replyTo}
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, partial, managerSession))(t, http.StatusOK)
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, partial, adminSession))(t, http.StatusOK)
				}

				listed := tests.AssertRequest(cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Q: &thread.Title,
				}))(t, http.StatusOK).JSON200
				r.Len(listed.Threads, 1)
				r.Equal(locked.LockedAt, listed.Threads[0].LockedAt)

				reopened := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					LifecycleState: opt.New(openapi.ThreadLifecycleStateOpen).Ptr(),
				}, authorSession))(t, http.StatusOK).JSON200
				r.Equal(locked.LockedAt, reopened.LockedAt)

				tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					Locked: opt.New(false).Ptr(),
				}, managerSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, openapi.ReplyInitialProps{
					Body: "<p>Discussion continues</p>",
				}, authorSession))(t, http.StatusOK)
			})
		}))
	}))
}
