package thread_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/oapi-codegen/nullable"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadLifecycle(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		db *ent.Client,
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
				response := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Lifecycle " + xid.New().String(),
					Body:       opt.New("<p>A discussion</p>").Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, authorSession))(t, http.StatusOK)
				return response.JSON200
			}

			createReply := func(t *testing.T, thread *openapi.Thread, replyTo *string) *openapi.Reply {
				t.Helper()
				response := tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, openapi.ReplyInitialProps{
					Body:    "<p>A possible conclusion</p>",
					ReplyTo: replyTo,
				}, otherSession))(t, http.StatusOK)
				return response.JSON200
			}

			updateThread := func(t *testing.T, thread *openapi.Thread, partial openapi.ThreadMutableProps) *openapi.Thread {
				t.Helper()
				response := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial, authorSession))(t, http.StatusOK)
				return response.JSON200
			}

			getThread := func(t *testing.T, thread *openapi.Thread) *openapi.Thread {
				t.Helper()
				response := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil))(t, http.StatusOK)
				return response.JSON200
			}

			t.Run("read_single_reply_with_thread_and_reply_visibility", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				parent := createReply(t, thread, nil)
				nested := createReply(t, thread, &parent.Id)

				tests.AssertRequest(cl.LikePostAddWithResponse(root, nested.Id, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.LikePostAddWithResponse(root, parent.Id, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.PostReactAddWithResponse(root, parent.Id, openapi.ReactInitialProps{Emoji: "👍"}, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.PostReactAddWithResponse(root, nested.Id, openapi.ReactInitialProps{Emoji: "🎉"}, authorSession))(t, http.StatusOK)
				read := tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id, authorSession))(t, http.StatusOK).JSON200
				r.Equal(nested.Body, read.Body)
				r.Equal(thread.Id, read.RootId)
				r.Equal(nested.Author.Id, read.Author.Id)
				r.Equal(1, read.Likes.Likes)
				r.True(read.Likes.Liked)
				r.Len(read.Reacts, 1)
				r.Equal("🎉", read.Reacts[0].Emoji)
				r.Equal(parent.Id, read.ReplyTo.Id)
				r.Equal(thread.Id, read.ReplyTo.RootId)
				r.Equal(1, read.ReplyTo.Likes.Likes)
				r.Len(read.ReplyTo.Reacts, 1)
				listed := getThread(t, thread)
				for _, item := range listed.Replies.Replies {
					r.Equal(thread.Id, item.RootId)
					r.Equal(thread.Slug, item.RootSlug)
				}
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id))(t, http.StatusOK)

				foreign := createThread(t)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, foreign.Slug, nested.Id))(t, http.StatusNotFound)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, thread.Id))(t, http.StatusNotFound)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, xid.New().String()))(t, http.StatusNotFound)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, "invalid"))(t, http.StatusBadRequest)

				updateThread(t, thread, openapi.ThreadMutableProps{Visibility: opt.New(openapi.VisibilityDraft).Ptr()})
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id))(t, http.StatusForbidden)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id, otherSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id, managerSession))(t, http.StatusOK)
				updateThread(t, thread, openapi.ThreadMutableProps{Visibility: opt.New(openapi.VisibilityPublished).Ptr()})

				tests.AssertRequest(cl.PostUpdateWithResponse(root, parent.Id, openapi.PostMutableProps{Visibility: opt.New(openapi.VisibilityReview).Ptr()}, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, parent.Id))(t, http.StatusForbidden)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, parent.Id, authorSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, parent.Id, otherSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, parent.Id, managerSession))(t, http.StatusOK)
				publicNested := tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id))(t, http.StatusOK).JSON200
				r.Nil(publicNested.ReplyTo)

				tests.AssertRequest(cl.PostDeleteWithResponse(root, nested.Id, otherSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, nested.Id, managerSession))(t, http.StatusNotFound)
				tests.AssertRequest(cl.ThreadDeleteWithResponse(root, thread.Slug, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, parent.Id, managerSession))(t, http.StatusNotFound)
			})

			t.Run("resolve_without_conclusion_and_continue_replying", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				r.Equal(openapi.ThreadLifecycleStateOpen, thread.LifecycleState)
				r.Nil(thread.ConclusionReplyId)

				for range 2 {
					resolved := updateThread(t, thread, openapi.ThreadMutableProps{
						LifecycleState: opt.New(openapi.ThreadLifecycleStateResolved).Ptr(),
					})
					r.Equal(openapi.ThreadLifecycleStateResolved, resolved.LifecycleState)
					r.Nil(resolved.ConclusionReplyId)
				}

				createReply(t, thread, nil)
				r.Equal(openapi.ThreadLifecycleStateResolved, getThread(t, thread).LifecycleState)

				reopened := updateThread(t, thread, openapi.ThreadMutableProps{
					LifecycleState: opt.New(openapi.ThreadLifecycleStateOpen).Ptr(),
				})
				r.Equal(openapi.ThreadLifecycleStateOpen, reopened.LifecycleState)
			})

			t.Run("nominate_replace_remove_and_reopen", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				first := createReply(t, thread, nil)
				nested := createReply(t, thread, &first.Id)

				resolved := updateThread(t, thread, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullableWithValue(nested.Id),
				})
				r.Equal(openapi.ThreadLifecycleStateResolved, resolved.LifecycleState)
				r.Equal(nested.Id, *resolved.ConclusionReplyId)

				edited := updateThread(t, thread, openapi.ThreadMutableProps{Title: opt.New(thread.Title + " updated").Ptr()})
				r.Equal(nested.Id, *edited.ConclusionReplyId)
				r.Equal(openapi.ThreadLifecycleStateResolved, edited.LifecycleState)

				listed := tests.AssertRequest(cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Q: &edited.Title,
				}))(t, http.StatusOK)
				r.Len(listed.JSON200.Threads, 1)
				r.Equal(nested.Id, *listed.JSON200.Threads[0].ConclusionReplyId)
				r.Equal(openapi.ThreadLifecycleStateResolved, listed.JSON200.Threads[0].LifecycleState)

				replaced := updateThread(t, thread, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullableWithValue(first.Id),
				})
				r.Equal(first.Id, *replaced.ConclusionReplyId)

				removed := updateThread(t, thread, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullNullable[string](),
				})
				r.Nil(removed.ConclusionReplyId)
				r.Equal(openapi.ThreadLifecycleStateResolved, removed.LifecycleState)

				updateThread(t, thread, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullableWithValue(nested.Id),
				})
				reopened := updateThread(t, thread, openapi.ThreadMutableProps{
					LifecycleState: opt.New(openapi.ThreadLifecycleStateOpen).Ptr(),
				})
				r.Nil(reopened.ConclusionReplyId)
				r.Equal(openapi.ThreadLifecycleStateOpen, reopened.LifecycleState)
				r.Nil(getThread(t, thread).ConclusionReplyId)
			})

			t.Run("reject_invalid_targets_without_mutating_thread", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				other := createThread(t)
				foreign := createReply(t, other, nil)
				deleted := createReply(t, thread, nil)
				review := createReply(t, thread, nil)
				tests.AssertRequest(cl.PostDeleteWithResponse(root, deleted.Id, otherSession))(t, http.StatusOK)
				tests.AssertRequest(cl.PostUpdateWithResponse(root, review.Id, openapi.PostMutableProps{
					Visibility: opt.New(openapi.VisibilityReview).Ptr(),
				}, managerSession))(t, http.StatusOK)

				for name, target := range map[string]string{
					"foreign_reply": foreign.Id,
					"thread_root":   thread.Id,
					"missing":       xid.New().String(),
					"malformed":     "invalid-id",
					"deleted":       deleted.Id,
					"review":        review.Id,
				} {
					t.Run(name, func(t *testing.T) {
						tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
							Title:             opt.New("Must not be saved").Ptr(),
							ConclusionReplyId: nullable.NewNullableWithValue(target),
						}, authorSession))(t, http.StatusBadRequest)
					})
				}

				valid := createReply(t, thread, nil)
				tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					LifecycleState:    opt.New(openapi.ThreadLifecycleStateOpen).Ptr(),
					ConclusionReplyId: nullable.NewNullableWithValue(valid.Id),
				}, authorSession))(t, http.StatusBadRequest)

				tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					LifecycleState: opt.New(openapi.ThreadLifecycleState("closed")).Ptr(),
				}, authorSession))(t, http.StatusBadRequest)

				unchanged := getThread(t, thread)
				r.Equal(thread.Title, unchanged.Title)
				r.Equal(openapi.ThreadLifecycleStateOpen, unchanged.LifecycleState)
				r.Nil(unchanged.ConclusionReplyId)
			})

			t.Run("author_or_manage_posts_permission_required", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				reply := createReply(t, thread, nil)

				for _, partial := range []openapi.ThreadMutableProps{
					{LifecycleState: opt.New(openapi.ThreadLifecycleStateResolved).Ptr()},
					{ConclusionReplyId: nullable.NewNullableWithValue(reply.Id)},
					{ConclusionReplyId: nullable.NewNullNullable[string]()},
				} {
					tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial, otherSession))(t, http.StatusForbidden)
					tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, partial))(t, http.StatusUnauthorized)
				}

				resolved := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullableWithValue(reply.Id),
				}, managerSession))(t, http.StatusOK)
				r.Equal(openapi.ThreadLifecycleStateResolved, resolved.JSON200.LifecycleState)
				r.Equal(reply.Id, *resolved.JSON200.ConclusionReplyId)

				reopened := tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
					LifecycleState: opt.New(openapi.ThreadLifecycleStateOpen).Ptr(),
				}, managerSession))(t, http.StatusOK)
				r.Equal(openapi.ThreadLifecycleStateOpen, reopened.JSON200.LifecycleState)
				r.Nil(reopened.JSON200.ConclusionReplyId)
			})

			t.Run("removing_or_unpublishing_conclusion_preserves_resolution", func(t *testing.T) {
				for _, action := range []string{"delete", "unpublish", "purge"} {
					t.Run(action, func(t *testing.T) {
						r := require.New(t)
						thread := createThread(t)
						reply := createReply(t, thread, nil)
						updateThread(t, thread, openapi.ThreadMutableProps{
							ConclusionReplyId: nullable.NewNullableWithValue(reply.Id),
						})

						before := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil))(t, http.StatusOK)
						etag := before.HTTPResponse.Header.Get("ETag")
						r.NotEmpty(etag)

						if action == "delete" {
							tests.AssertRequest(cl.PostDeleteWithResponse(root, reply.Id, otherSession))(t, http.StatusOK)
						} else if action == "purge" {
							id, err := xid.FromString(reply.Id)
							r.NoError(err)
							r.NoError(db.Post.DeleteOneID(id).Exec(root))
						} else {
							tests.AssertRequest(cl.PostUpdateWithResponse(root, reply.Id, openapi.PostMutableProps{
								Visibility: opt.New(openapi.VisibilityReview).Ptr(),
							}, managerSession))(t, http.StatusOK)
						}

						after := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, nil, func(ctx context.Context, req *http.Request) error {
							if action != "purge" {
								req.Header.Set("If-None-Match", etag)
							}
							return nil
						}))(t, http.StatusOK)
						r.Nil(after.JSON200.ConclusionReplyId)
						r.Equal(openapi.ThreadLifecycleStateResolved, after.JSON200.LifecycleState)
					})
				}
			})

			t.Run("conclusion_available_on_every_reply_page", func(t *testing.T) {
				r := require.New(t)
				thread := createThread(t)
				first := createReply(t, thread, nil)
				for range 49 {
					createReply(t, thread, nil)
				}
				conclusion := createReply(t, thread, &first.Id)
				updateThread(t, thread, openapi.ThreadMutableProps{
					ConclusionReplyId: nullable.NewNullableWithValue(conclusion.Id),
				})

				for _, page := range []string{"1", "2"} {
					response := tests.AssertRequest(cl.ThreadGetWithResponse(root, thread.Slug, &openapi.ThreadGetParams{
						Page: &page,
					}))(t, http.StatusOK)
					r.Equal(conclusion.Id, *response.JSON200.ConclusionReplyId)
					r.Equal(openapi.ThreadLifecycleStateResolved, response.JSON200.LifecycleState)
					r.Equal(51, response.JSON200.ReplyStatus.Replies)
				}

				location := tests.AssertRequest(cl.PostLocationGetWithResponse(root, &openapi.PostLocationGetParams{
					Id: conclusion.Id,
				}))(t, http.StatusOK)
				r.Equal(2, *location.JSON200.Page)
				read := tests.AssertRequest(cl.ReplyGetWithResponse(root, thread.Slug, conclusion.Id))(t, http.StatusOK)
				r.Equal(conclusion.Body, read.JSON200.Body)
				r.Equal(conclusion.Id, read.JSON200.Id)
			})

			t.Run("concurrent_nomination_and_removal", func(t *testing.T) {
				for _, action := range []string{"delete", "unpublish"} {
					t.Run(action, func(t *testing.T) {
						r := require.New(t)
						for range 5 {
							thread := createThread(t)
							reply := createReply(t, thread, nil)
							var nomination *openapi.ThreadUpdateResponse
							var nominationErr error
							var removalErr error
							var removalStatus int
							var workers sync.WaitGroup
							start := make(chan struct{})
							workers.Add(2)
							go func() {
								defer workers.Done()
								<-start
								nomination, nominationErr = cl.ThreadUpdateWithResponse(root, thread.Slug, openapi.ThreadMutableProps{
									ConclusionReplyId: nullable.NewNullableWithValue(reply.Id),
								}, authorSession)
							}()
							go func() {
								defer workers.Done()
								<-start
								if action == "delete" {
									response, err := cl.PostDeleteWithResponse(root, reply.Id, otherSession)
									removalErr = err
									if response != nil {
										removalStatus = response.StatusCode()
									}
								} else {
									response, err := cl.PostUpdateWithResponse(root, reply.Id, openapi.PostMutableProps{
										Visibility: opt.New(openapi.VisibilityReview).Ptr(),
									}, managerSession)
									removalErr = err
									if response != nil {
										removalStatus = response.StatusCode()
									}
								}
							}()
							close(start)
							workers.Wait()
							r.NoError(nominationErr)
							r.NoError(removalErr)
							r.Equal(http.StatusOK, removalStatus)
							r.Contains([]int{http.StatusOK, http.StatusBadRequest}, nomination.StatusCode())
							r.Nil(getThread(t, thread).ConclusionReplyId)
						}
					})
				}
			})

			t.Run("filter_before_pagination_and_combine_with_category", func(t *testing.T) {
				r := require.New(t)
				category := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name:   "Lifecycle filters " + xid.New().String(),
					Colour: "green",
				}, adminSession))(t, http.StatusOK)

				var resolvedIDs []string
				var openIDs []string
				for index := range 54 {
					thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
						Title:      fmt.Sprintf("Filter discussion %d", index),
						Category:   &category.JSON200.Id,
						Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
					}, authorSession))(t, http.StatusOK).JSON200
					if index%26 == 0 {
						openIDs = append(openIDs, thread.Id)
					} else {
						updateThread(t, thread, openapi.ThreadMutableProps{
							LifecycleState: opt.New(openapi.ThreadLifecycleStateResolved).Ptr(),
						})
						resolvedIDs = append(resolvedIDs, thread.Id)
					}
				}

				for state, expected := range map[openapi.ThreadLifecycleState][]string{
					openapi.ThreadLifecycleStateOpen:     openIDs,
					openapi.ThreadLifecycleStateResolved: resolvedIDs,
				} {
					var actual []string
					for page := 1; ; page++ {
						response := tests.AssertRequest(cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
							Categories:     &openapi.CategorySlugListQuery{category.JSON200.Slug},
							LifecycleState: &state,
							Page:           opt.New(fmt.Sprint(page)).Ptr(),
						}))(t, http.StatusOK)
						r.Equal((len(expected)+49)/50, response.JSON200.TotalPages)
						for _, thread := range response.JSON200.Threads {
							r.Equal(state, thread.LifecycleState)
							actual = append(actual, thread.Id)
						}

						if response.JSON200.NextPage == nil {
							break
						}
					}
					r.ElementsMatch(expected, actual)
				}

				tests.AssertRequest(cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					LifecycleState: opt.New(openapi.ThreadLifecycleState("closed")).Ptr(),
				}))(t, http.StatusBadRequest)
			})
		}))
	}))
}
