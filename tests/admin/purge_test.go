package admin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
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

func TestPurgeAccountContent(t *testing.T) {
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
			t.Run("purges_account_threads", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				vis := openapi.Published
				thread1, err := cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
					Title:      "Test Thread 1",
					Body:       opt.New("<p>Content 1</p>").Ptr(),
					Visibility: &vis,
				}, memberSession)
				tests.Ok(t, err, thread1)

				thread2, err := cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
					Title:      "Test Thread 2",
					Body:       opt.New("<p>Content 2</p>").Ptr(),
					Visibility: &vis,
				}, memberSession)
				tests.Ok(t, err, thread2)

				threads, err := cl.ThreadListWithResponse(adminCtx, &openapi.ThreadListParams{}, adminSession)
				tests.Ok(t, err, threads)
				initialCount := countUserThreads(threads.JSON200.Threads, member.Handle)
				a.GreaterOrEqual(initialCount, 2, "Member should have at least 2 threads")

				var purgeAction openapi.ModerationActionCreate
				err = purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Threads,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				threadsAfter, err := cl.ThreadListWithResponse(adminCtx, &openapi.ThreadListParams{}, adminSession)
				tests.Ok(t, err, threadsAfter)
				afterCount := countUserThreads(threadsAfter.JSON200.Threads, member.Handle)
				a.Equal(0, afterCount, "Member should have no threads after purge")

				get1, err := cl.ThreadGetWithResponse(adminCtx, thread1.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				r.NoError(err)
				a.Equal(http.StatusNotFound, get1.StatusCode())

				get2, err := cl.ThreadGetWithResponse(adminCtx, thread2.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				r.NoError(err)
				a.Equal(http.StatusNotFound, get2.StatusCode())
			})

			t.Run("purges_account_replies", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				otherMemberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				otherMemberSession := sh.WithSession(otherMemberCtx)

				vis := openapi.Published
				createdThread, err := cl.ThreadCreateWithResponse(otherMemberCtx, openapi.ThreadInitialProps{
					Title:      "Test Thread",
					Body:       opt.New("<p>Thread content</p>").Ptr(),
					Visibility: &vis,
				}, otherMemberSession)
				tests.Ok(t, err, createdThread)

				reply1, err := cl.ReplyCreateWithResponse(memberCtx, createdThread.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "Test reply 1",
				}, memberSession)
				tests.Ok(t, err, reply1)

				reply2, err := cl.ReplyCreateWithResponse(memberCtx, createdThread.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "Test reply 2",
				}, memberSession)
				tests.Ok(t, err, reply2)

				threadGet, err := cl.ThreadGetWithResponse(adminCtx, createdThread.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				tests.Ok(t, err, threadGet)
				initialReplyCount := threadGet.JSON200.Replies.Results
				a.GreaterOrEqual(initialReplyCount, 2)

				var purgeAction openapi.ModerationActionCreate
				err = purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Replies,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				threadAfter, err := cl.ThreadGetWithResponse(adminCtx, createdThread.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				tests.Ok(t, err, threadAfter)

				memberReplyCount := 0
				for _, post := range threadAfter.JSON200.Replies.Replies {
					if post.Author.Handle == member.Handle {
						memberReplyCount++
					}
				}
				a.Equal(0, memberReplyCount, "Member should have no replies after purge")
			})

			t.Run("purges_reply_preserves_nested_replies", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				threadAuthorCtx, _ := e2e.WithAccount(root, aw, seed.Account_005_Þórr)
				threadAuthorSession := sh.WithSession(threadAuthorCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				nestedReplyAuthorCtx, nestedReplyAuthor := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				nestedReplyAuthorSession := sh.WithSession(nestedReplyAuthorCtx)

				// Create a thread
				vis := openapi.Published
				threadResp := tests.AssertRequest(
					cl.ThreadCreateWithResponse(threadAuthorCtx, openapi.ThreadInitialProps{
						Title:      "Thread for Nested Replies",
						Body:       opt.New("<p>Thread content</p>").Ptr(),
						Visibility: &vis,
					}, threadAuthorSession),
				)(t, http.StatusOK)

				// Member creates a reply to the thread
				parentReplyResp := tests.AssertRequest(
					cl.ReplyCreateWithResponse(memberCtx, threadResp.JSON200.Slug, openapi.ReplyInitialProps{
						Body: "Parent reply from member",
					}, memberSession),
				)(t, http.StatusOK)

				// Another user creates a nested reply to the member's reply
				nestedReplyResp := tests.AssertRequest(
					cl.ReplyCreateWithResponse(nestedReplyAuthorCtx, threadResp.JSON200.Slug, openapi.ReplyInitialProps{
						Body:    "Nested reply to parent",
						ReplyTo: &parentReplyResp.JSON200.Id,
					}, nestedReplyAuthorSession),
				)(t, http.StatusOK)

				// Verify the nested reply has a replyTo reference
				threadBefore, err := cl.ThreadGetWithResponse(adminCtx, threadResp.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				tests.Ok(t, err, threadBefore)

				var foundNested bool
				for _, reply := range threadBefore.JSON200.Replies.Replies {
					if reply.Id == nestedReplyResp.JSON200.Id {
						foundNested = true
						a.NotNil(reply.ReplyTo, "Nested reply should have replyTo set before purge")
						if reply.ReplyTo != nil {
							a.Equal(parentReplyResp.JSON200.Id, reply.ReplyTo.Id)
							a.Nil(reply.ReplyTo.DeletedAt, "Parent reply should not be deleted yet")
						}
						break
					}
				}
				a.True(foundNested, "Should find nested reply before purge")

				// Purge the member's replies (deleting the parent reply)
				var purgeAction openapi.ModerationActionCreate
				err = purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Replies,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				// Verify the parent reply is deleted but nested reply is preserved
				threadAfter, err := cl.ThreadGetWithResponse(adminCtx, threadResp.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				tests.Ok(t, err, threadAfter)

				// Verify parent reply is gone
				var foundParent bool
				for _, reply := range threadAfter.JSON200.Replies.Replies {
					if reply.Author.Handle == member.Handle {
						foundParent = true
						break
					}
				}
				a.False(foundParent, "Parent reply should be deleted")

				// Verify nested reply still exists but ReplyTo is nil (soft-deleted parent filtered out)
				var foundNestedAfter bool
				for _, reply := range threadAfter.JSON200.Replies.Replies {
					if reply.Id == nestedReplyResp.JSON200.Id {
						foundNestedAfter = true
						a.Equal(nestedReplyAuthor.Handle, reply.Author.Handle)
						// The ReplyTo should be nil because soft-deleted parents are filtered from the query
						a.Nil(reply.ReplyTo, "Nested reply's replyTo should be nil (soft-deleted parent filtered)")
						break
					}
				}
				a.True(foundNestedAfter, "Nested reply should still exist after parent deletion")
			})

			t.Run("purges_multiple_content_types", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				otherMemberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				otherMemberSession := sh.WithSession(otherMemberCtx)

				vis := openapi.Published
				memberThread, err := cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
					Title:      "Member Thread",
					Body:       opt.New("<p>Member content</p>").Ptr(),
					Visibility: &vis,
				}, memberSession)
				tests.Ok(t, err, memberThread)

				otherThread, err := cl.ThreadCreateWithResponse(otherMemberCtx, openapi.ThreadInitialProps{
					Title:      "Other Thread",
					Body:       opt.New("<p>Other content</p>").Ptr(),
					Visibility: &vis,
				}, otherMemberSession)
				tests.Ok(t, err, otherThread)

				reply, err := cl.ReplyCreateWithResponse(memberCtx, otherThread.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "Test reply",
				}, memberSession)
				tests.Ok(t, err, reply)

				var purgeAction openapi.ModerationActionCreate
				err = purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Threads,
						openapi.Replies,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				get, err := cl.ThreadGetWithResponse(adminCtx, memberThread.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				r.NoError(err)
				a.Equal(http.StatusNotFound, get.StatusCode())

				threadAfter, err := cl.ThreadGetWithResponse(adminCtx, otherThread.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession)
				tests.Ok(t, err, threadAfter)

				memberReplyCount := 0
				for _, post := range threadAfter.JSON200.Replies.Replies {
					if post.Author.Handle == member.Handle {
						memberReplyCount++
					}
				}
				a.Equal(0, memberReplyCount, "Member should have no replies after purge")
			})

			t.Run("purges_profile_bio", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				_, err := cl.AccountUpdateWithResponse(memberCtx, openapi.AccountUpdateJSONRequestBody{
					Bio: toPtr("<body>This is my bio</body>"),
				}, memberSession)
				r.NoError(err)

				profile, err := cl.ProfileGetWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, profile)
				a.NotEmpty(profile.JSON200.Bio)
				a.Equal("<body>This is my bio</body>", profile.JSON200.Bio)

				var purgeAction openapi.ModerationActionCreate
				err = purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.ProfileBio,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				profileAfter, err := cl.ProfileGetWithResponse(adminCtx, member.Handle, adminSession)
				tests.Ok(t, err, profileAfter)
				a.Equal("<body></body>", profileAfter.JSON200.Bio, "Bio should be empty after purge")
			})

			t.Run("requires_admin_permission", func(t *testing.T) {
				r := require.New(t)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				otherMemberCtx, otherMember := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				_ = sh.WithSession(otherMemberCtx)

				var purgeAction openapi.ModerationActionCreate
				err := purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: otherMember.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Threads,
					},
				})
				r.NoError(err)

				create, err := cl.ModerationActionCreateWithResponse(memberCtx, purgeAction, memberSession)
				r.NoError(err)
				r.Equal(http.StatusForbidden, create.StatusCode())
			})

			t.Run("creates_audit_event", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				_ = sh.WithSession(memberCtx)

				var purgeAction openapi.ModerationActionCreate
				err := purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Threads,
					},
				})
				r.NoError(err)

				create := tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				a.Equal(openapi.AccountContentPurged, create.JSON201.Type)

				time.Sleep(100 * time.Millisecond)

				list, err := cl.AuditEventListWithResponse(adminCtx, &openapi.AuditEventListParams{}, adminSession)
				tests.Ok(t, err, list)

				var found bool
				for _, event := range *list.JSON200.Events {
					if event.Type == openapi.AccountContentPurged {
						purgedEvent, err := event.AsAuditEventAccountContentPurged()
						r.NoError(err)
						a.Equal(openapi.Identifier(member.ID.String()), purgedEvent.AccountId)
						a.NotNil(purgedEvent.Included)
						found = true
						break
					}
				}
				a.True(found, "Should find account_content_purged event in audit log")
			})

			t.Run("purges_collections_and_cascade_deletes_items", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				// Create a thread to add to collection
				vis := openapi.Published
				threadResp := tests.AssertRequest(
					cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
						Title:      "Thread for Collection",
						Body:       opt.New("<p>Content</p>").Ptr(),
						Visibility: &vis,
					}, memberSession),
				)(t, http.StatusOK)

				// Create a collection owned by the member
				collectionResp := tests.AssertRequest(
					cl.CollectionCreateWithResponse(memberCtx, openapi.CollectionInitialProps{
						Name: "Test Collection",
					}, memberSession),
				)(t, http.StatusOK)

				// Add the thread to the collection
				tests.AssertRequest(
					cl.CollectionAddPostWithResponse(memberCtx, collectionResp.JSON200.Slug, threadResp.JSON200.Id, memberSession),
				)(t, http.StatusOK)

				// Verify the collection has items
				collectionGet := tests.AssertRequest(
					cl.CollectionGetWithResponse(memberCtx, collectionResp.JSON200.Slug, memberSession),
				)(t, http.StatusOK)
				a.GreaterOrEqual(len(collectionGet.JSON200.Items), 1, "Collection should have at least 1 item")

				// Purge member's collections
				var purgeAction openapi.ModerationActionCreate
				err := purgeAction.FromModerationActionCreatePurgeAccount(openapi.ModerationActionCreatePurgeAccount{
					Action:    "purge_account",
					AccountId: member.ID.String(),
					Include: []openapi.ModerationActionPurgeAccountContentType{
						openapi.Collections,
					},
				})
				r.NoError(err)

				tests.AssertRequest(
					cl.ModerationActionCreateWithResponse(adminCtx, purgeAction, adminSession),
				)(t, http.StatusCreated)

				// Verify collection is deleted
				tests.AssertRequest(
					cl.CollectionGetWithResponse(adminCtx, collectionResp.JSON200.Slug, adminSession),
				)(t, http.StatusNotFound)

				// Verify the thread itself still exists (wasn't cascade deleted)
				threadAfter := tests.AssertRequest(
					cl.ThreadGetWithResponse(adminCtx, threadResp.JSON200.Slug, &openapi.ThreadGetParams{}, adminSession),
				)(t, http.StatusOK)
				a.Equal(threadResp.JSON200.Id, threadAfter.JSON200.Id, "Thread should still exist after collection deletion")
			})
		}))
	}))
}

func countUserThreads(threads []openapi.ThreadReference, handle string) int {
	count := 0
	for _, thread := range threads {
		if thread.Author.Handle == handle {
			count++
		}
	}
	return count
}

func toPtr[T any](v T) *T {
	return &v
}
