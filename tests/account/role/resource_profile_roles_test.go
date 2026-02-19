package role_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestResourceProfileReferencesIncludeRoles(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			authorCtx, author := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			authorSession := sh.WithSession(authorCtx)

			roleMeta := openapi.Metadata{
				"style": "bold",
			}
			roleColour := "#f97316"

			roleCreate := tests.AssertRequest(cl.RoleCreateWithResponse(authorCtx, openapi.RoleCreateJSONRequestBody{
				Name:        "resource-role-" + xid.New().String(),
				Colour:      roleColour,
				Permissions: openapi.PermissionList{},
				Meta:        &roleMeta,
			}, authorSession))(t, http.StatusOK)

			tests.AssertRequest(cl.AccountAddRoleWithResponse(authorCtx, author.Handle, roleCreate.JSON200.Id, authorSession))(t, http.StatusOK)

			threadCategory := tests.AssertRequest(cl.CategoryCreateWithResponse(authorCtx, openapi.CategoryCreateJSONRequestBody{
				Name:        "category-" + xid.New().String(),
				Description: "test",
				Colour:      "#f97316",
			}, authorSession))(t, http.StatusOK)

			vis := openapi.Published
			threadCreate := tests.AssertRequest(cl.ThreadCreateWithResponse(authorCtx, openapi.ThreadInitialProps{
				Title:      "thread-with-roles-" + xid.New().String(),
				Body:       opt.New("<p>thread body</p>").Ptr(),
				Category:   opt.New(threadCategory.JSON200.Id).Ptr(),
				Visibility: &vis,
			}, authorSession))(t, http.StatusOK)

			replyCreate := tests.AssertRequest(cl.ReplyCreateWithResponse(authorCtx, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>reply body</p>",
			}, authorSession))(t, http.StatusOK)

			nodeCreate := tests.AssertRequest(cl.NodeCreateWithResponse(authorCtx, openapi.NodeCreateJSONRequestBody{
				Name:       "node-with-roles-" + xid.New().String(),
				Visibility: &vis,
			}, authorSession))(t, http.StatusOK)

			collectionCreate := tests.AssertRequest(cl.CollectionCreateWithResponse(authorCtx, openapi.CollectionCreateJSONRequestBody{
				Name: "collection-with-roles-" + xid.New().String(),
			}, authorSession))(t, http.StatusOK)
			tests.AssertRequest(cl.CollectionAddPostWithResponse(authorCtx, collectionCreate.JSON200.Id, threadCreate.JSON200.Id, authorSession))(t, http.StatusOK)
			tests.AssertRequest(cl.CollectionAddNodeWithResponse(authorCtx, collectionCreate.JSON200.Id, nodeCreate.JSON200.Id, authorSession))(t, http.StatusOK)

			assertProfileRefHasRole := func(t *testing.T, where string, roles []openapi.AccountRoleRef) {
				t.Helper()
				r := require.New(t)
				a := assert.New(t)

				found := findRoleRef(roles, roleCreate.JSON200.Id)
				r.NotNil(found, where)
				a.Equal(roleColour, found.Colour, where)
				r.NotNil(found.Meta, where)
				a.Equal("bold", (*found.Meta)["style"], where)
				a.False(found.Default, where)
			}

			t.Run("thread_get_author", func(t *testing.T) {
				threadGet := tests.AssertRequest(cl.ThreadGetWithResponse(authorCtx, threadCreate.JSON200.Slug, nil, authorSession))(t, http.StatusOK)
				assertProfileRefHasRole(t, "thread.get author should include assigned custom role", threadGet.JSON200.Author.Roles)
			})

			t.Run("thread_get_replies_author", func(t *testing.T) {
				r := require.New(t)
				threadGet := tests.AssertRequest(cl.ThreadGetWithResponse(authorCtx, threadCreate.JSON200.Slug, nil, authorSession))(t, http.StatusOK)
				replyItem, found := lo.Find(threadGet.JSON200.Replies.Replies, func(in openapi.Reply) bool {
					return in.Id == replyCreate.JSON200.Id
				})
				r.True(found, "created reply should appear in thread replies")
				assertProfileRefHasRole(t, "thread.get reply author should include assigned custom role", replyItem.Author.Roles)
			})

			t.Run("thread_list_author", func(t *testing.T) {
				r := require.New(t)
				authorFilter := openapi.AccountHandle(author.Handle)
				threadList := tests.AssertRequest(cl.ThreadListWithResponse(authorCtx, &openapi.ThreadListParams{
					Author: &authorFilter,
				}, authorSession))(t, http.StatusOK)
				threadListItem, found := lo.Find(threadList.JSON200.Threads, func(in openapi.ThreadReference) bool {
					return in.Id == threadCreate.JSON200.Id
				})
				r.True(found, "created thread should appear in thread list")
				assertProfileRefHasRole(t, "thread.list author should include assigned custom role", threadListItem.Author.Roles)
			})

			t.Run("reply_create_author", func(t *testing.T) {
				assertProfileRefHasRole(t, "reply.create author should include assigned custom role", replyCreate.JSON200.Author.Roles)
			})

			t.Run("node_get_owner", func(t *testing.T) {
				nodeGet := tests.AssertRequest(cl.NodeGetWithResponse(authorCtx, nodeCreate.JSON200.Slug, &openapi.NodeGetParams{}, authorSession))(t, http.StatusOK)
				assertProfileRefHasRole(t, "node.get owner should include assigned custom role", nodeGet.JSON200.Owner.Roles)
			})

			t.Run("node_list_owner", func(t *testing.T) {
				r := require.New(t)
				nodeList := tests.AssertRequest(cl.NodeListWithResponse(authorCtx, &openapi.NodeListParams{}, authorSession))(t, http.StatusOK)
				nodeItem, found := lo.Find(nodeList.JSON200.Nodes, func(in openapi.NodeWithChildren) bool {
					return in.Id == nodeCreate.JSON200.Id
				})
				r.True(found, "created node should appear in node list")
				assertProfileRefHasRole(t, "node.list owner should include assigned custom role", nodeItem.Owner.Roles)
			})

			t.Run("collection_get_owner", func(t *testing.T) {
				collectionGet := tests.AssertRequest(cl.CollectionGetWithResponse(authorCtx, collectionCreate.JSON200.Id, authorSession))(t, http.StatusOK)
				assertProfileRefHasRole(t, "collection.get owner should include assigned custom role", collectionGet.JSON200.Owner.Roles)
			})

			t.Run("collection_list_owner", func(t *testing.T) {
				r := require.New(t)
				collectionList := tests.AssertRequest(cl.CollectionListWithResponse(authorCtx, &openapi.CollectionListParams{}, authorSession))(t, http.StatusOK)
				collectionItem, found := lo.Find(collectionList.JSON200.Collections, func(in openapi.Collection) bool {
					return in.Id == collectionCreate.JSON200.Id
				})
				r.True(found, "created collection should appear in collection list")
				assertProfileRefHasRole(t, "collection.list owner should include assigned custom role", collectionItem.Owner.Roles)
			})

			t.Run("collection_get_items_owner_and_item_refs", func(t *testing.T) {
				r := require.New(t)
				collectionGet := tests.AssertRequest(cl.CollectionGetWithResponse(authorCtx, collectionCreate.JSON200.Id, authorSession))(t, http.StatusOK)

				var foundThreadItem bool
				var foundNodeItem bool

				for _, item := range collectionGet.JSON200.Items {
					assertProfileRefHasRole(t, "collection item owner should include assigned custom role", item.Owner.Roles)

					postItem, postErr := item.Item.AsDatagraphItemPost()
					if postErr == nil && postItem.Ref.Id == threadCreate.JSON200.Id {
						foundThreadItem = true
						assertProfileRefHasRole(t, "collection item post author should include assigned custom role", postItem.Ref.Author.Roles)
					}

					nodeItem, nodeErr := item.Item.AsDatagraphItemNode()
					if nodeErr == nil && nodeItem.Ref.Id == nodeCreate.JSON200.Id {
						foundNodeItem = true
						assertProfileRefHasRole(t, "collection item node owner should include assigned custom role", nodeItem.Ref.Owner.Roles)
					}
				}

				r.True(foundThreadItem, "expected thread item in collection")
				r.True(foundNodeItem, "expected node item in collection")
			})
		}))
	}))
}
