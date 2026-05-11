package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestCollectionItems(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{Identifier: xid.New().String(), Token: "password"})
			tests.Ok(t, err, acc1)
			session1 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))))

			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{Identifier: xid.New().String(), Token: "password"})
			tests.Ok(t, err, acc2)
			session2 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))))

			collection1, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
				Name: "c1",
			}, session1)
			tests.Ok(t, err, collection1)

			cat1, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "",
				Description: "cat",
				Name:        xid.New().String(),
			}, adminSession)
			tests.Ok(t, err, cat1)

			threadCreateProps := openapi.ThreadInitialProps{
				Body:       opt.New("<p>this is a thread</p>").Ptr(),
				Category:   opt.New(cat1.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "thread",
			}

			published := openapi.Published

			thread1create, err := cl.ThreadCreateWithResponse(root, threadCreateProps, session1)
			tests.Ok(t, err, thread1create)

			thread2create, err := cl.ThreadCreateWithResponse(root, threadCreateProps, session2)
			tests.Ok(t, err, thread2create)

			node1create, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr(), Visibility: &published}, adminSession)
			tests.Ok(t, err, node1create)

			node2create, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr(), Visibility: &published}, adminSession)
			tests.Ok(t, err, node2create)

			t.Run("unauthorised", func(t *testing.T) {
				t.Parallel()

				addPost1, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id)
				tests.Status(t, err, addPost1, http.StatusUnauthorized)
			})

			t.Run("add_remove_items", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)

				// Add 2 posts and 2 nodes to collection1
				addPost1, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1)
				addNode1, err := cl.CollectionAddNodeWithResponse(root, collection1.JSON200.Id, node1create.JSON200.Id, session1)
				tests.Ok(t, err, addNode1)
				addPost2, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread2create.JSON200.Id, session1)
				tests.Ok(t, err, addPost2)
				addNode2, err := cl.CollectionAddNodeWithResponse(root, collection1.JSON200.Id, node2create.JSON200.Id, session1)
				tests.Ok(t, err, addNode2)

				get1, err := cl.CollectionGetWithResponse(root, collection1.JSON200.Id)
				tests.Ok(t, err, get1)

				r.Len(get1.JSON200.Items, 4)

				// These must be in order of addition

				matchNodeToItem(t, node2create.JSON200, get1.JSON200.Items[0])
				matchThreadToItem(t, thread2create.JSON200, get1.JSON200.Items[1])
				matchNodeToItem(t, node1create.JSON200, get1.JSON200.Items[2])
				matchThreadToItem(t, thread1create.JSON200, get1.JSON200.Items[3])

				removePost1, err := cl.CollectionRemovePostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, removePost1)
				removePost2, err := cl.CollectionRemoveNodeWithResponse(root, collection1.JSON200.Id, node1create.JSON200.Id, session1)
				tests.Ok(t, err, removePost2)

				get2, err := cl.CollectionGetWithResponse(root, collection1.JSON200.Id)
				tests.Ok(t, err, get2)

				r.Len(get2.JSON200.Items, 2)

				matchNodeToItem(t, node2create.JSON200, get1.JSON200.Items[0])
				matchThreadToItem(t, thread2create.JSON200, get1.JSON200.Items[1])
			})

			t.Run("add_idempotent", func(t *testing.T) {
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "x1",
				}, session1)
				tests.Ok(t, err, col)

				// Add the same post twice, should not error and be a no-op

				addPost1, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1)

				addPost1again, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1again)
			})

			t.Run("collection_item_status", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "x1",
				}, session1)
				tests.Ok(t, err, col)

				thr, err := cl.ThreadCreateWithResponse(root, threadCreateProps, session1)
				tests.Ok(t, err, thr)

				addThr, err := cl.CollectionAddPostWithResponse(root, col.JSON200.Id, thr.JSON200.Id, session1)
				tests.Ok(t, err, addThr)

				// Get the thread in a list and directly as the owner of the collection

				ownerlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, session1)
				tests.Ok(t, err, ownerlist)
				fromList, found := lo.Find(ownerlist.JSON200.Threads, func(tr openapi.ThreadReference) bool { return tr.Id == thr.JSON200.Id })
				r.True(found)
				a.True(fromList.Collections.HasCollected)
				a.Equal(1, fromList.Collections.InCollections)

				ownerget, err := cl.ThreadGetWithResponse(root, thr.JSON200.Id, nil, session1)
				tests.Ok(t, err, ownerget)
				a.True(ownerget.JSON200.Collections.HasCollected)
				a.Equal(1, ownerget.JSON200.Collections.InCollections)

				// A different session this time, not the owner of the collection

				randolist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, session2)
				tests.Ok(t, err, randolist)
				fromList, found = lo.Find(randolist.JSON200.Threads, func(tr openapi.ThreadReference) bool { return tr.Id == thr.JSON200.Id })
				r.True(found)
				a.False(fromList.Collections.HasCollected)
				a.Equal(1, fromList.Collections.InCollections)

				randoget, err := cl.ThreadGetWithResponse(root, thr.JSON200.Id, nil, session2)
				tests.Ok(t, err, randoget)
				a.False(randoget.JSON200.Collections.HasCollected)
				a.Equal(1, randoget.JSON200.Collections.InCollections)

				// And finally, as a guest with no session at all

				guestlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
				tests.Ok(t, err, guestlist)
				fromList, found = lo.Find(guestlist.JSON200.Threads, func(tr openapi.ThreadReference) bool { return tr.Id == thr.JSON200.Id })
				r.True(found)
				a.False(fromList.Collections.HasCollected)
				a.Equal(1, fromList.Collections.InCollections)

				guestget, err := cl.ThreadGetWithResponse(root, thr.JSON200.Id, nil)
				tests.Ok(t, err, guestget)
				a.False(guestget.JSON200.Collections.HasCollected)
				a.Equal(1, guestget.JSON200.Collections.InCollections)
			})

			t.Run("query_having_item_by_id", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				acc1, err := cl.AccountGetWithResponse(root, session1)
				tests.Ok(t, err, acc1)

				colHavingItem, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "c1",
				}, session1)
				tests.Ok(t, err, colHavingItem)

				// add a post and node
				addPost1, err := cl.CollectionAddPostWithResponse(root, colHavingItem.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1)
				addNode1, err := cl.CollectionAddNodeWithResponse(root, colHavingItem.JSON200.Id, node1create.JSON200.Id, session1)
				tests.Ok(t, err, addNode1)

				get1, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &acc1.JSON200.Handle,
					HasItem:       &thread1create.JSON200.Id,
				})
				tests.Ok(t, err, get1)
				col1 := find(t, get1.JSON200.Collections, colHavingItem.JSON200.Id)
				a.True(col1.HasQueriedItem)

				get2, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &acc1.JSON200.Handle,
					HasItem:       &node1create.JSON200.Id,
				})
				tests.Ok(t, err, get2)
				col2 := find(t, get2.JSON200.Collections, colHavingItem.JSON200.Id)
				a.True(col2.HasQueriedItem)

				get3, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &acc1.JSON200.Handle,
					HasItem:       nil,
				})
				tests.Ok(t, err, get3)
				col3 := find(t, get3.JSON200.Collections, colHavingItem.JSON200.Id)
				a.False(col3.HasQueriedItem)

				get4, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &acc1.JSON200.Handle,
					HasItem:       &[]string{xid.New().String()}[0],
				})
				tests.Ok(t, err, get4)
				col4 := find(t, get4.JSON200.Collections, colHavingItem.JSON200.Id)
				a.False(col4.HasQueriedItem)
			})
		}))
	}))
}

func find(t *testing.T, collections []openapi.Collection, cid openapi.CollectionMarkParam) *openapi.Collection {
	for _, c := range collections {
		if c.Id == cid {
			return &c
		}
	}
	t.Fatalf("collection not found: %s", cid)
	return nil
}
