package collection_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
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

func TestCollectionSubmissions(t *testing.T) {
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

			acc3, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{Identifier: xid.New().String(), Token: "password"})
			tests.Ok(t, err, acc3)
			session3 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc3.JSON200.Id)))))

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

			t.Run("submit_published_to_someone_elses_collection", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "c1",
				}, session1)
				tests.Ok(t, err, col)

				// Owner adds a thread and a node to their own collection
				addpost, err := cl.CollectionAddPostWithResponse(root, col.JSON200.Id, thread1create.JSON200.Id, session1)
				tests.Ok(t, err, addpost)
				addnode, err := cl.CollectionAddNodeWithResponse(root, col.JSON200.Id, node1create.JSON200.Id, session1)
				tests.Ok(t, err, addnode)

				// Non-owner submits a thread and a node to Owner's collection
				submitnode, err := cl.CollectionAddNodeWithResponse(root, col.JSON200.Id, node2create.JSON200.Id, session2)
				tests.Ok(t, err, submitnode)
				submitpost, err := cl.CollectionAddPostWithResponse(root, col.JSON200.Id, thread2create.JSON200.Id, session2)
				tests.Ok(t, err, submitpost)

				// guest cannot see the node in the collection
				get1, err := cl.CollectionGetWithResponse(root, col.JSON200.Id)
				tests.Ok(t, err, get1)
				r.Len(get1.JSON200.Items, 2)

				// unrelated user, acc 3, cannot see the node in the collection
				gets3, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session3)
				tests.Ok(t, err, gets3)
				r.Len(gets3.JSON200.Items, 2)

				// acc1 can see the node in the collection
				gets1, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, gets1)
				r.Len(gets1.JSON200.Items, 4)

				a.Equal(openapi.SubmissionReview, gets1.JSON200.Items[0].MembershipType)
				a.Equal(openapi.SubmissionReview, gets1.JSON200.Items[1].MembershipType)

				// acc2 can see the node in the collection
				gets2, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session2)
				tests.Ok(t, err, gets2)
				r.Len(gets2.JSON200.Items, 3)
				// NOTE: This is 3 because currently the owner of the submission
				// is not stored and filtering is performed against the owner of
				// the resource itself. Which removes both node2 and thread2
				// from the list, however posts are not yet properly filtered
				// and are always treated as published so the post is present.

				a.Equal(openapi.SubmissionReview, gets1.JSON200.Items[0].MembershipType)
				a.Equal(openapi.SubmissionReview, gets1.JSON200.Items[1].MembershipType)
			})

			t.Run("submit_unlisted_node_to_someone_elses_collection", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				unlisted := openapi.Unlisted

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "c2",
				}, session1)
				tests.Ok(t, err, col)

				unlistedNode, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr(), Visibility: &unlisted}, session2)
				tests.Ok(t, err, unlistedNode)

				submitnode, err := cl.CollectionAddNodeWithResponse(root, col.JSON200.Id, unlistedNode.JSON200.Id, session2)
				tests.Ok(t, err, submitnode)

				// guest cannot see items in review
				getGuest, err := cl.CollectionGetWithResponse(root, col.JSON200.Id)
				tests.Ok(t, err, getGuest)
				a.Len(getGuest.JSON200.Items, 0)

				// unrelated cannot see items in review
				getUnrelated, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session3)
				tests.Ok(t, err, getUnrelated)
				r.Len(getUnrelated.JSON200.Items, 0)

				// acc1 can the node that's submitted and in review
				getOwner, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, getOwner)
				r.Len(getOwner.JSON200.Items, 1)
				a.Equal(unlistedNode.JSON200.Id, getOwner.JSON200.Items[0].Id)
				a.Equal(openapi.SubmissionReview, getOwner.JSON200.Items[0].MembershipType)

				// acc2 can the node that's submitted and in review
				getSubmitter, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session2)
				tests.Ok(t, err, getSubmitter)
				r.Len(getSubmitter.JSON200.Items, 1)
				a.Equal(unlistedNode.JSON200.Id, getSubmitter.JSON200.Items[0].Id)
				a.Equal(openapi.SubmissionReview, getSubmitter.JSON200.Items[0].MembershipType)
			})

			t.Run("accept_submission", func(t *testing.T) {
				t.Parallel()
				r := require.New(t)
				a := assert.New(t)

				unlisted := openapi.Unlisted

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "c2",
				}, session1)
				tests.Ok(t, err, col)

				unlistedNode, err := cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Content: opt.New("<p>hi</p>").Ptr(), Visibility: &unlisted}, session2)
				tests.Ok(t, err, unlistedNode)

				submitnode, err := cl.CollectionAddNodeWithResponse(root, col.JSON200.Id, unlistedNode.JSON200.Id, session2)
				tests.Ok(t, err, submitnode)

				getOwner, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, getOwner)
				r.Len(getOwner.JSON200.Items, 1)
				a.Equal(openapi.SubmissionReview, getOwner.JSON200.Items[0].MembershipType)

				acceptNode, err := cl.CollectionAddNodeWithResponse(root, col.JSON200.Id, unlistedNode.JSON200.Id, session1)
				tests.Ok(t, err, acceptNode)

				getOwner2, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, getOwner2)
				r.Len(getOwner2.JSON200.Items, 1)
				a.Equal(openapi.SubmissionAccepted, getOwner2.JSON200.Items[0].MembershipType)
			})
		}))
	}))
}
