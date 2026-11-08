package thread_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

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

func TestThreadListReadStatusFilter(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberSession := sh.WithSession(memberCtx)
			otherSession := sh.WithSession(otherCtx)
			adminSession := sh.WithSession(adminCtx)

			createThread := func(t *testing.T, title string, visibility openapi.Visibility, session openapi.RequestEditorFn) *openapi.Thread {
				t.Helper()
				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      title,
					Body:       new("<p>A discussion</p>"),
					Visibility: new(visibility),
				}, session))(t, http.StatusOK).JSON200
				return thread
			}

			updateStatus := func(t *testing.T, thread *openapi.Thread, status openapi.ThreadReadStatusMutablePropsStatus, session openapi.RequestEditorFn) {
				t.Helper()
				tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, thread.Slug, openapi.ThreadReadStatusMutableProps{
					Status: status,
				}, session))(t, http.StatusOK)
			}

			listIDs := func(t *testing.T, params openapi.ThreadListParams, session openapi.RequestEditorFn) []string {
				t.Helper()
				response := tests.AssertRequest(cl.ThreadListWithResponse(root, &params, session))(t, http.StatusOK)
				ids := make([]string, 0, len(response.JSON200.Threads))
				for _, thread := range response.JSON200.Threads {
					ids = append(ids, thread.Id)
				}

				return ids
			}

			t.Run("membership_and_explicit_actions", func(t *testing.T) {
				r := require.New(t)
				prefix := "readmembership" + xid.New().String()
				unvisited := createThread(t, prefix+" unvisited", openapi.VisibilityPublished, otherSession)
				caughtUp := createThread(t, prefix+" caught up", openapi.VisibilityPublished, otherSession)
				marked := createThread(t, prefix+" marked", openapi.VisibilityPublished, otherSession)
				updateStatus(t, caughtUp, openapi.ThreadReadStatusMutablePropsStatusRead, memberSession)
				updateStatus(t, marked, openapi.ThreadReadStatusMutablePropsStatusUnread, memberSession)
				updateStatus(t, unvisited, openapi.ThreadReadStatusMutablePropsStatusRead, otherSession)
				updateStatus(t, caughtUp, openapi.ThreadReadStatusMutablePropsStatusUnread, otherSession)

				params := openapi.ThreadListParams{
					Q:          new(prefix),
					ReadStatus: new(openapi.ThreadReadStatusFilterUnread),
				}
				r.ElementsMatch([]string{unvisited.Id, marked.Id}, listIDs(t, params, memberSession))
				params.ReadStatus = new(openapi.ThreadReadStatusFilterRead)
				r.Equal([]string{caughtUp.Id}, listIDs(t, params, memberSession))
				r.Equal([]string{unvisited.Id}, listIDs(t, params, otherSession))
				params.ReadStatus = nil
				r.Len(listIDs(t, params, memberSession), 3)

				updateStatus(t, unvisited, openapi.ThreadReadStatusMutablePropsStatusRead, memberSession)
				updateStatus(t, caughtUp, openapi.ThreadReadStatusMutablePropsStatusUnread, memberSession)
				params.ReadStatus = new(openapi.ThreadReadStatusFilterUnread)
				r.ElementsMatch([]string{marked.Id, caughtUp.Id}, listIDs(t, params, memberSession))
			})

			t.Run("combines_with_author_category_and_tags", func(t *testing.T) {
				r := require.New(t)
				prefix := "readfilters" + xid.New().String()
				category := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name:        prefix,
					Colour:      "#abcdef",
					Description: "Read filter fixtures",
				}, adminSession))(t, http.StatusOK).JSON200

				createTagged := func(t *testing.T, title string, categoryID *openapi.Identifier, tags openapi.TagNameList, session openapi.RequestEditorFn) *openapi.Thread {
					t.Helper()
					thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
						Title:      title,
						Body:       new("<p>A discussion</p>"),
						Visibility: new(openapi.VisibilityPublished),
						Category:   categoryID,
						Tags:       new(tags),
					}, session))(t, http.StatusOK).JSON200
					return thread
				}

				matched := createTagged(t, prefix+" matched", new(category.Id), openapi.TagNameList{prefix}, otherSession)
				read := createTagged(t, prefix+" read", new(category.Id), openapi.TagNameList{prefix}, otherSession)
				createTagged(t, prefix+" wrong author", new(category.Id), openapi.TagNameList{prefix}, memberSession)
				createTagged(t, prefix+" wrong category", nil, openapi.TagNameList{prefix}, otherSession)
				createTagged(t, prefix+" wrong tag", new(category.Id), openapi.TagNameList{}, otherSession)
				updateStatus(t, read, openapi.ThreadReadStatusMutablePropsStatusRead, memberSession)
				updateStatus(t, matched, openapi.ThreadReadStatusMutablePropsStatusRead, otherSession)
				r.Len(matched.Tags, 1)
				tag := tests.AssertRequest(cl.TagGetWithResponse(root, prefix))(t, http.StatusOK).JSON200

				params := openapi.ThreadListParams{
					Q:          new(prefix),
					Author:     new(other.Handle),
					Categories: new([]string{category.Slug}),
					Tags:       new(openapi.TagListIDs{tag.Id}),
					ReadStatus: new(openapi.ThreadReadStatusFilterUnread),
				}
				r.Equal([]string{matched.Id}, listIDs(t, params, memberSession))
				params.ReadStatus = new(openapi.ThreadReadStatusFilterRead)
				r.Equal([]string{read.Id}, listIDs(t, params, memberSession))
			})

			t.Run("preserves_visibility_and_deletion_rules", func(t *testing.T) {
				r := require.New(t)
				prefix := "readvisibility" + xid.New().String()
				published := createThread(t, prefix+" published", openapi.VisibilityPublished, otherSession)
				ownReview := createThread(t, prefix+" own review", openapi.VisibilityReview, memberSession)
				otherReview := createThread(t, prefix+" other review", openapi.VisibilityReview, otherSession)
				ownDraft := createThread(t, prefix+" own draft", openapi.VisibilityDraft, memberSession)
				createThread(t, prefix+" other draft", openapi.VisibilityDraft, otherSession)
				deleted := createThread(t, prefix+" deleted", openapi.VisibilityPublished, otherSession)
				tests.AssertRequest(cl.ThreadDeleteWithResponse(root, deleted.Slug, otherSession))(t, http.StatusOK)

				params := openapi.ThreadListParams{
					Q:          new(prefix),
					ReadStatus: new(openapi.ThreadReadStatusFilterUnread),
				}
				r.ElementsMatch([]string{published.Id, ownReview.Id}, listIDs(t, params, memberSession))
				r.ElementsMatch([]string{published.Id, ownReview.Id, otherReview.Id}, listIDs(t, params, adminSession))
				params.Visibility = new([]openapi.Visibility{openapi.VisibilityDraft})
				params.Author = new(member.Handle)
				r.Equal([]string{ownDraft.Id}, listIDs(t, params, memberSession))
				r.Empty(listIDs(t, params, otherSession))
			})

			t.Run("filters_before_pagination_and_orders_by_activity", func(t *testing.T) {
				r := require.New(t)
				prefix := "readpages" + xid.New().String()
				var expected []string
				var oldest *openapi.Thread
				for index := range 3 {
					thread := createThread(t, fmt.Sprintf("%s %d", prefix, index), openapi.VisibilityPublished, otherSession)
					expected = append([]string{thread.Id}, expected...)
					if index == 0 {
						oldest = thread
						tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
							Pinned: new(openapi.PinnedRank(100)),
						}, adminSession))(t, http.StatusOK)
					}
				}

				for range 2 {
					thread := createThread(t, prefix+" read", openapi.VisibilityPublished, otherSession)
					updateStatus(t, thread, openapi.ThreadReadStatusMutablePropsStatusRead, memberSession)
				}

				params := openapi.ThreadListParams{
					Q:            new(prefix),
					ReadStatus:   new(openapi.ThreadReadStatusFilterUnread),
					IgnorePinned: new(true),
				}
				first := tests.AssertRequest(cl.ThreadListWithResponse(root, &params, memberSession))(t, http.StatusOK).JSON200
				r.Len(first.Threads, 3)
				r.Equal(3, first.Results)
				r.Equal(1, first.TotalPages)
				r.Nil(first.NextPage)
				params.Page = new("2")
				second := tests.AssertRequest(cl.ThreadListWithResponse(root, &params, memberSession))(t, http.StatusOK).JSON200
				r.Empty(second.Threads)
				r.Equal(2, second.CurrentPage)
				r.Equal(1, second.TotalPages)
				r.Nil(second.NextPage)
				var actual []string
				for _, thread := range append(first.Threads, second.Threads...) {
					actual = append(actual, thread.Id)
				}

				r.Equal(expected, actual)
				tests.AssertRequest(cl.ReplyCreateWithResponse(root, oldest.Slug, openapi.ReplyInitialProps{
					Body: "New activity on the oldest thread",
				}, otherSession))(t, http.StatusOK)
				params.Page = new("1")
				activity := listIDs(t, params, memberSession)
				r.Equal(oldest.Id, activity[0])
				params.ReadStatus = new(openapi.ThreadReadStatusFilterRead)
				read := tests.AssertRequest(cl.ThreadListWithResponse(root, &params, memberSession))(t, http.StatusOK).JSON200
				r.Len(read.Threads, 2)
				r.Equal(1, read.TotalPages)
				r.Nil(read.NextPage)
			})

			t.Run("requires_authentication_only_when_filtering", func(t *testing.T) {
				prefix := "readguest" + xid.New().String()
				createThread(t, prefix, openapi.VisibilityPublished, otherSession)
				params := openapi.ThreadListParams{Q: new(prefix)}
				response := tests.AssertRequest(cl.ThreadListWithResponse(root, &params))(t, http.StatusOK)
				require.Len(t, response.JSON200.Threads, 1)
				for _, status := range []openapi.ThreadReadStatusFilter{openapi.ThreadReadStatusFilterRead, openapi.ThreadReadStatusFilterUnread} {
					params.ReadStatus = new(status)
					tests.AssertRequest(cl.ThreadListWithResponse(root, &params))(t, http.StatusUnauthorized)
				}

				params.ReadStatus = new(openapi.ThreadReadStatusFilter("invalid"))
				tests.AssertRequest(cl.ThreadListWithResponse(root, &params, memberSession))(t, http.StatusBadRequest)
			})
		}))
	}))
}
