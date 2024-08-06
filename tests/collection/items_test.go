package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/cookie"
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
		cj *cookie.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := e2e.WithSession(adminCtx, cj)

			acc1, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{xid.New().String(), "password"})
			tests.Ok(t, err, acc1)
			session1 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))), cj)

			acc2, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{xid.New().String(), "password"})
			tests.Ok(t, err, acc2)
			session2 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))), cj)

			collection1, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
				Name:        "c1",
				Description: "owned by acc1",
			}, session1)
			tests.Ok(t, err, collection1)

			cat1, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "",
				Description: "cat",
				Name:        xid.New().String(),
			}, adminSession)
			tests.Ok(t, err, cat1)

			threadCreateProps := openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat1.JSON200.Id,
				Visibility: openapi.Published,
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
				tests.Status(t, err, addPost1, http.StatusForbidden)
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
					Name:        "x1",
					Description: "owned by acc1",
				}, session1)
				tests.Ok(t, err, col)

				// Add the same post twice, should not error and be a no-op

				addPost1, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1)

				addPost1again, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addPost1again)
			})
		}))
	}))
}
