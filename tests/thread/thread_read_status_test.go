package thread_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_read_state"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/postread"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadReadStatusUpdate(t *testing.T) {
	t.Parallel()

	var run func()
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		rw *post_read_state.Writer,
		db *ent.Client,
	) {
		run = func() {
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)
			otherSession := sh.WithSession(otherCtx)

			createThread := func(t *testing.T, visibility openapi.Visibility) *openapi.Thread {
				t.Helper()
				response := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Explicit read status " + xid.New().String(),
					Body:       opt.New("<p>A discussion</p>").Ptr(),
					Visibility: &visibility,
				}, memberSession))(t, http.StatusOK)
				return response.JSON200
			}

			updateStatus := func(t *testing.T, slug string, status openapi.ThreadReadStatusMutablePropsStatus) *openapi.ReadStatus {
				t.Helper()
				response := tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, slug, openapi.ThreadReadStatusMutableProps{
					Status: status,
				}, memberSession))(t, http.StatusOK)
				return response.JSON200
			}

			t.Run("preserves_baseline_and_reply_count", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t, openapi.VisibilityPublished)
				id := post.ID(openapi.ParseID(thread.Id))
				baseline := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
				r.NoError(rw.MarkRead(root, member.ID, id, baseline))

				for range 5 {
					tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, openapi.ReplyInitialProps{
						Body: "A reply after the baseline",
					}, otherSession))(t, http.StatusOK)
				}

				before := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil, memberSession))(t, http.StatusOK)
				r.Equal(5, before.JSON200.ReadStatus.RepliesSince)
				etag := before.HTTPResponse.Header.Get("ETag")

				unread := updateStatus(t, thread.Slug, openapi.ThreadReadStatusMutablePropsStatusUnread)
				r.True(unread.LastReadAt.Equal(baseline))
				r.Equal(5, unread.RepliesSince)
				r.NotNil(unread.MarkedUnreadAt)

				again := updateStatus(t, thread.Slug, openapi.ThreadReadStatusMutablePropsStatusUnread)
				r.True(again.LastReadAt.Equal(baseline))
				r.Equal(5, again.RepliesSince)

				get := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil, memberSession, func(ctx context.Context, request *http.Request) error {
					request.Header.Set("If-None-Match", etag)
					return nil
				}))(t, http.StatusOK)
				r.NotNil(get.JSON200.ReadStatus.MarkedUnreadAt)
				r.Equal(5, get.JSON200.ReadStatus.RepliesSince)

				list := tests.AssertRequest(cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Q: new(thread.Title),
				}, memberSession))(t, http.StatusOK)
				found := false
				for _, item := range list.JSON200.Threads {
					if item.Id == thread.Id {
						found = true
						r.NotNil(item.ReadStatus.MarkedUnreadAt)
						r.Equal(5, item.ReadStatus.RepliesSince)
					}
				}
				r.True(found)

				otherGet := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil, otherSession))(t, http.StatusOK)
				r.Nil(otherGet.JSON200.ReadStatus)
				r.False(db.PostRead.Query().Where(postread.AccountID(xid.ID(other.ID)), postread.RootPostID(xid.ID(id))).ExistX(root))

				read := updateStatus(t, thread.Slug, openapi.ThreadReadStatusMutablePropsStatusRead)
				r.Nil(read.MarkedUnreadAt)
				r.Zero(read.RepliesSince)
				r.True(read.LastReadAt.After(baseline))
			})

			t.Run("initialises_untracked_threads", func(t *testing.T) {
				for _, status := range []openapi.ThreadReadStatusMutablePropsStatus{openapi.ThreadReadStatusMutablePropsStatusUnread, openapi.ThreadReadStatusMutablePropsStatusRead} {
					t.Run(string(status), func(t *testing.T) {
						r := require.New(t)
						thread := createThread(t, openapi.VisibilityPublished)
						result := updateStatus(t, thread.Slug, status)
						r.Zero(result.RepliesSince)
						r.False(result.LastReadAt.IsZero())

						if status == openapi.ThreadReadStatusMutablePropsStatusUnread {
							r.NotNil(result.MarkedUnreadAt)
							r.True(result.MarkedUnreadAt.Equal(result.LastReadAt))
						} else {
							r.Nil(result.MarkedUnreadAt)
						}
					})
				}
			})

			t.Run("orders_beacons_and_explicit_actions", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t, openapi.VisibilityPublished)
				id := post.ID(openapi.ParseID(thread.Id))
				baseline := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
				marked := baseline.Add(10 * time.Minute)
				visited := baseline.Add(20 * time.Minute)

				r.NoError(rw.MarkRead(root, member.ID, id, baseline))
				r.NoError(rw.MarkUnread(root, member.ID, id, marked))
				r.NoError(rw.MarkRead(root, member.ID, id, marked.Add(-time.Second)))
				r.NoError(rw.MarkRead(root, member.ID, id, marked))

				row := db.PostRead.Query().Where(postread.RootPostID(xid.ID(id)), postread.AccountID(xid.ID(member.ID))).OnlyX(root)
				r.True(row.LastSeenAt.Equal(baseline))
				r.NotNil(row.MarkedUnreadAt)
				r.True(row.MarkedUnreadAt.Equal(marked))

				r.NoError(rw.MarkRead(root, member.ID, id, visited))
				r.NoError(rw.MarkRead(root, member.ID, id, baseline))
				r.NoError(rw.MarkUnread(root, member.ID, id, marked))

				row = db.PostRead.Query().Where(postread.RootPostID(xid.ID(id)), postread.AccountID(xid.ID(member.ID))).OnlyX(root)
				r.True(row.LastSeenAt.Equal(visited))
				r.Nil(row.MarkedUnreadAt)

				unread := updateStatus(t, thread.Slug, openapi.ThreadReadStatusMutablePropsStatusUnread)
				body := "{\"k\":\"thread\",\"id\":\"" + thread.Id + "\"}"
				tests.AssertRequest(cl.SendBeaconWithBodyWithResponse(root, "text/plain", strings.NewReader(body), memberSession))(t, http.StatusAccepted)
				r.Eventually(func() bool {
					row, err := db.PostRead.Query().Where(postread.RootPostID(xid.ID(id)), postread.AccountID(xid.ID(member.ID))).Only(root)
					return err == nil && row.MarkedUnreadAt == nil && row.LastSeenAt.After(*unread.MarkedUnreadAt)
				}, 5*time.Second, 10*time.Millisecond)
			})

			t.Run("rejects_unauthorised_or_invalid_updates", func(t *testing.T) {
				thread := createThread(t, openapi.VisibilityPublished)
				body := openapi.ThreadReadStatusMutableProps{Status: openapi.ThreadReadStatusMutablePropsStatusUnread}
				tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, thread.Slug, body))(t, http.StatusUnauthorized)
				tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, xid.New().String(), body, memberSession))(t, http.StatusNotFound)

				invalid := openapi.ThreadReadStatusMutableProps{Status: "invalid"}
				tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, thread.Slug, invalid, memberSession))(t, http.StatusBadRequest)

				draft := createThread(t, openapi.VisibilityDraft)
				tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, draft.Slug, body, otherSession))(t, http.StatusForbidden)
			})
		}
	}))
	run()
}
