package tree_visibility_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesTreeQueryingVisibilityFilters(t *testing.T) {
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

			member1Ctx, _ := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			member1Session := sh.WithSession(member1Ctx)

			// member2Ctx, _ := e2e.WithAccount(root, aw, seed.Account_008_Heimdallr)
			// member2Session := sh.WithSession(member2Ctx)

			published := openapi.Published
			draft := openapi.Draft

			// SETUP: 4 published library pages
			//
			// node1       <- root               has 2 children: node2 and node3
			// |- node2    <- child of node1     has no children
			// |- node3    <- child of node1     has 1 child: node4
			//    |- node4 <- child of node3     has no children

			node1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("n1"), Visibility: &published}, adminSession)
			tests.Ok(t, err, node1)

			node2, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("n2"), Visibility: &published, Parent: &node1.JSON200.Slug}, adminSession)
			tests.Ok(t, err, node2)

			node3, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("n3"), Visibility: &published, Parent: &node1.JSON200.Slug}, adminSession)
			tests.Ok(t, err, node3)

			node4, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("n4"), Visibility: &published, Parent: &node3.JSON200.Slug}, adminSession)
			tests.Ok(t, err, node4)

			rootNodeIDs := []string{node1.JSON200.Id}
			nonRootNodeIDs := []string{node2.JSON200.Id, node3.JSON200.Id, node4.JSON200.Id}

			t.Run("query_all_top_level", func(t *testing.T) {
				a := assert.New(t)

				// member 1 creates a draft under node3
				draft1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("m1draft"), Visibility: &draft, Parent: &node3.JSON200.Slug}, member1Session)
				tests.Ok(t, err, draft1)

				draftIDs := []string{draft1.JSON200.Id}

				// listing with no filters should return only published nodes
				list1, err := cl.NodeListWithResponse(root, &openapi.NodeListParams{}, member1Session)
				tests.Ok(t, err, list1)
				list1IDs := ids(list1.JSON200.Nodes)
				a.Subset(list1IDs, rootNodeIDs)
				a.NotSubset(list1IDs, nonRootNodeIDs)

				// listing only drafts yields nothing because there are no
				// drafts at the root meaning there are no fully draft trees.
				list2, err := cl.NodeListWithResponse(root, &openapi.NodeListParams{Visibility: &[]openapi.Visibility{draft}}, member1Session)
				tests.Ok(t, err, list2)
				list2IDs := ids(list2.JSON200.Nodes)
				a.NotSubset(list2IDs, rootNodeIDs)
				a.NotSubset(list2IDs, draftIDs)
				a.NotSubset(list2IDs, nonRootNodeIDs)
			})

			t.Run("query_published_and_drafts", func(t *testing.T) {
				a := assert.New(t)

				// member 1 creates a draft under node3
				draft1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: un("m1draft"), Visibility: &draft, Parent: &node3.JSON200.Slug}, member1Session)
				tests.Ok(t, err, draft1)

				draftIDs := []string{draft1.JSON200.Id}

				// listing published and drafts yields published nodes as well
				// as drafts from this member interspersed with the full tree.
				list3, err := cl.NodeListWithResponse(root, &openapi.NodeListParams{Visibility: &[]openapi.Visibility{published, draft}}, member1Session)
				tests.Ok(t, err, list3)
				list3IDs := ids(list3.JSON200.Nodes)
				a.Subset(list3IDs, rootNodeIDs)
				a.NotSubset(list3IDs, draftIDs)
				a.NotSubset(list3IDs, nonRootNodeIDs)

				n5 := find(t, list3.JSON200.Nodes, draft1.JSON200.Id)
				a.Equal(draft1.JSON200.Id, n5.Id, "draft node is in the list under node3")
			})
		}))
	}))
}

func un(n string) string {
	return n + " " + xid.New().String()
}

func ids(nodes []openapi.NodeWithChildren) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.Id
	}
	return ids
}

func find(t *testing.T, roots []openapi.NodeWithChildren, id string) *openapi.NodeWithChildren {
	t.Helper()

	// walk from root through children looking for node with id == id
	var walk func([]openapi.NodeWithChildren) *openapi.NodeWithChildren
	walk = func(nodes []openapi.NodeWithChildren) *openapi.NodeWithChildren {
		for _, n := range nodes {
			if n.Id == id {
				return &n
			}
			if n.Children != nil {
				if r := walk(n.Children); r != nil {
					return r
				}
			}
		}
		return nil
	}

	r := walk(roots)

	if r == nil {
		t.Fatalf("could not find node with id %s", id)
	}

	return r
}
