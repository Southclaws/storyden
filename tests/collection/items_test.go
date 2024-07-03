package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
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
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, ar, seed.Account_001_Odin)
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

			collection2, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
				Name:        "c2",
				Description: "owned by acc2",
			}, session2)
			tests.Ok(t, err, collection2)

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

			thread1create, err := cl.ThreadCreateWithResponse(root, threadCreateProps, session1)
			tests.Ok(t, err, thread1create)

			thread2create, err := cl.ThreadCreateWithResponse(root, threadCreateProps, session2)
			tests.Ok(t, err, thread2create)

			node1create, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr()}, session1)
			tests.Ok(t, err, node1create)

			node2create, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr()}, session2)
			tests.Ok(t, err, node2create)

			t.Run("unauthorised", func(t *testing.T) {
				t.Parallel()

				addPost1, err := cl.CollectionAddPostWithResponse(root, collection1.JSON200.Id, thread1create.JSON200.Id, session2)
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

				matchThreadToItem(t, thread1create.JSON200, get1.JSON200.Items[0])
				matchThreadToItem(t, thread2create.JSON200, get1.JSON200.Items[1])
				matchNodeToItem(t, node1create.JSON200, get1.JSON200.Items[2])
				matchNodeToItem(t, node2create.JSON200, get1.JSON200.Items[3])
			})
		}))
	}))
}

func matchThreadToItem(t *testing.T, thread *openapi.Thread, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	a.Equal(openapi.DatagraphNodeKindPost, item.Kind)
	a.Equal(thread.Id, item.Id)
	// a.Equal(thread.CreatedAt, item.CreatedAt) // TODO
	a.Equal(thread.Title, item.Name)
	a.Contains(thread.Slug, item.Slug)
	a.Equal(thread.Short, item.Description)
	a.Equal(thread.Author, item.Owner)
}

func matchNodeToItem(t *testing.T, node *openapi.Node, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	a.Equal(openapi.DatagraphNodeKindNode, item.Kind)
	a.Equal(node.Id, item.Id)
	// a.Equal(node.CreatedAt, item.CreatedAt) // TODO
	a.Equal(node.Name, item.Name)
	a.Contains(node.Slug, item.Slug)
	a.Equal(node.Description, *item.Description)
	a.Equal(node.Owner, item.Owner)
}
