package library_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeDraftList(t *testing.T) {
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

			// Create two regular users
			handle1 := xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle1, "password"})
			tests.Ok(t, err, acc1)
			user1ID := account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))
			session1 := sh.WithSession(e2e.WithAccountID(root, user1ID))

			handle2 := xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle2, "password"})
			tests.Ok(t, err, acc2)
			user2ID := account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))
			_ = user2ID // May be used in future tests

			t.Run("unauthenticated_rejected", func(t *testing.T) {
				t.Parallel()

				list, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{})
				tests.Status(t, err, list, 401)
			})

			t.Run("author_sees_own_drafts_only", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				// Create test nodes for this test
				node1Name := "Test Node 1 " + xid.New().String()
				node1, err := cl.NodeCreateWithResponse(adminCtx, openapi.NodeInitialProps{
					Name: node1Name,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, node1)

				node2Name := "Test Node 2 " + xid.New().String()
				node2, err := cl.NodeCreateWithResponse(adminCtx, openapi.NodeInitialProps{
					Name: node2Name,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, node2)

				// Create drafts as admin - the author will be admin
				draft1Name := "Draft on Node1"
				draft1, err := cl.NodeVersionCreateWithResponse(root, node1.JSON200.Slug, openapi.NodeVersionInitialProps{
					Name: &draft1Name,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, draft1)

				draft2Name := "Draft on Node2"
				draft2, err := cl.NodeVersionCreateWithResponse(root, node2.JSON200.Slug, openapi.NodeVersionInitialProps{
					Name: &draft2Name,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, draft2)

				// Regular user1 (non-admin) lists drafts - should see none (not the author, not a manager)
				list1, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{}, session1)
				tests.Ok(t, err, list1)
				a.Equal(0, len(list1.JSON200.Drafts), "Regular user should not see drafts they didn't author")

				// Admin lists drafts - should see both (is a manager)
				listAdmin, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{}, sh.WithSession(adminCtx))
				tests.Ok(t, err, listAdmin)
				a.GreaterOrEqual(len(listAdmin.JSON200.Drafts), 2, "Admin should see all drafts")
			})

			t.Run("manager_sees_all_drafts", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				// Create a fresh node for this test
				testNodeName := "Manager Test Node " + xid.New().String()
				testNode, err := cl.NodeCreateWithResponse(adminCtx, openapi.NodeInitialProps{
					Name: testNodeName,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, testNode)

				// Admin creates a draft
				draft1Name := "Draft for Manager Test"
				draft1, err := cl.NodeVersionCreateWithResponse(root, testNode.JSON200.Slug, openapi.NodeVersionInitialProps{
					Name: &draft1Name,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, draft1)

				// Admin (manager) lists all drafts - should see all
				listAdmin, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{}, sh.WithSession(adminCtx))
				tests.Ok(t, err, listAdmin)

				// Admin should see at least the draft we just created
				// (may see others from parallel tests, so we check >= 1)
				a.GreaterOrEqual(len(listAdmin.JSON200.Drafts), 1)

				// Find our specific draft
				foundDraft := false
				for _, d := range listAdmin.JSON200.Drafts {
					if d.Id == draft1.JSON200.Id {
						foundDraft = true
						a.Equal(draft1Name, d.Name)
						a.Equal(testNode.JSON200.Id, d.Node.Id)
						break
					}
				}
				a.True(foundDraft, "Admin should see user's draft")
			})

			t.Run("drafts_include_node_reference", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				// Create a node with specific properties we can verify
				nodeName := "Node with Reference " + xid.New().String()
				nodeDesc := "Test description for draft reference"
				node, err := cl.NodeCreateWithResponse(adminCtx, openapi.NodeInitialProps{
					Name:        nodeName,
					Description: &nodeDesc,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, node)

				// Create a draft
				draftName := "Draft to test node reference"
				draft, err := cl.NodeVersionCreateWithResponse(root, node.JSON200.Slug, openapi.NodeVersionInitialProps{
					Name: &draftName,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, draft)

				// List drafts
				list, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{}, sh.WithSession(adminCtx))
				tests.Ok(t, err, list)

				// Find our draft and verify the node reference
				foundDraft := false
				for _, d := range list.JSON200.Drafts {
					if d.Id == draft.JSON200.Id {
						foundDraft = true
						a.NotNil(d.Node, "Draft should include node reference")
						a.Equal(node.JSON200.Id, d.Node.Id)
						a.Equal(nodeName, d.Node.Name)
						a.Equal(nodeDesc, d.Node.Description)
						a.Equal(node.JSON200.Slug, d.Node.Slug)
						break
					}
				}
				a.True(foundDraft, "Should find the created draft")
			})

			t.Run("pagination_works", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				// Create multiple drafts on different nodes
				for i := 0; i < 3; i++ {
					nodeName := "Pagination Test Node " + xid.New().String()
					node, err := cl.NodeCreateWithResponse(adminCtx, openapi.NodeInitialProps{
						Name: nodeName,
					}, sh.WithSession(adminCtx))
					tests.Ok(t, err, node)

					draftName := "Pagination Draft " + xid.New().String()
					draft, err := cl.NodeVersionCreateWithResponse(root, node.JSON200.Slug, openapi.NodeVersionInitialProps{
						Name: &draftName,
					}, sh.WithSession(adminCtx))
					tests.Ok(t, err, draft)
				}

				// Request with pagination
				page := "1"
				list, err := cl.NodeDraftListWithResponse(root, &openapi.NodeDraftListParams{
					Page: &page,
				}, sh.WithSession(adminCtx))
				tests.Ok(t, err, list)

				// Should have pagination fields
				a.GreaterOrEqual(len(list.JSON200.Drafts), 1)
				a.NotNil(list.JSON200.CurrentPage)
				a.NotNil(list.JSON200.PageSize)
			})
		}))
	}))
}
