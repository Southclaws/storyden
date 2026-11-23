package library_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeTreeRelationshipOwnership(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		nw *node_writer.Writer,
		roles *role_repo.Repository,
		assignments *role_assign.Assignment,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			_, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			managerCtx, manager := e2e.WithAccount(root, aw, seed.Account_005_Þórr)
			managerRole, err := roles.Create(root, "tree-manager-"+xid.New().String(), "blue", rbac.PermissionList{rbac.PermissionManageLibrary})
			require.NoError(t, err)
			require.False(t, manager.Admin)
			require.NoError(t, assignments.UpdateRoles(root, account_ref.ID(manager.ID), role_assign.Add(managerRole.ID)))

			adminSession := sh.WithSession(adminCtx)
			memberSession := sh.WithSession(memberCtx)
			managerSession := sh.WithSession(managerCtx)

			cases := []struct {
				name            string
				childOwner      account.AccountID
				parentOwner     account.AccountID
				childVisibility visibility.Visibility
				actor           openapi.RequestEditorFn
				moveStatus      int
				severStatus     int
			}{
				{
					name:            "owns_parent_only",
					childOwner:      other.ID,
					parentOwner:     member.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           memberSession,
					moveStatus:      http.StatusForbidden,
					severStatus:     http.StatusForbidden,
				},
				{
					name:            "owns_child_only",
					childOwner:      member.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           memberSession,
					moveStatus:      http.StatusForbidden,
					severStatus:     http.StatusOK,
				},
				{
					name:            "owns_neither",
					childOwner:      other.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           memberSession,
					moveStatus:      http.StatusForbidden,
					severStatus:     http.StatusForbidden,
				},
				{
					name:            "owns_both",
					childOwner:      member.ID,
					parentOwner:     member.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           memberSession,
					moveStatus:      http.StatusOK,
					severStatus:     http.StatusOK,
				},
				{
					name:            "draft_submission",
					childOwner:      member.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityDraft,
					actor:           memberSession,
					moveStatus:      http.StatusOK,
					severStatus:     http.StatusOK,
				},
				{
					name:            "review_submission",
					childOwner:      member.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityReview,
					actor:           memberSession,
					moveStatus:      http.StatusOK,
					severStatus:     http.StatusOK,
				},
				{
					name:            "library_manager",
					childOwner:      other.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           managerSession,
					moveStatus:      http.StatusOK,
					severStatus:     http.StatusOK,
				},
				{
					name:            "administrator",
					childOwner:      other.ID,
					parentOwner:     other.ID,
					childVisibility: visibility.VisibilityPublished,
					actor:           adminSession,
					moveStatus:      http.StatusOK,
					severStatus:     http.StatusOK,
				},
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					reader := sh.WithSession(e2e.WithAccountID(root, tc.childOwner))

					t.Run("move", func(t *testing.T) {
						original := createTreeOwnershipNode(t, root, nw, tc.childOwner, visibility.VisibilityPublished)
						parent := createTreeOwnershipNode(t, root, nw, tc.parentOwner, visibility.VisibilityPublished)
						child := createTreeOwnershipNode(t, root, nw, tc.childOwner, tc.childVisibility, node_writer.WithParent(library.NodeID(original.GetID())))
						assertTreeParent(t, root, cl, reader, child, original)

						response, err := cl.NodeAddNodeWithResponse(root, parent.GetSlug(), child.GetSlug(), tc.actor)
						tests.Status(t, err, response, tc.moveStatus)
						if tc.moveStatus == http.StatusOK {
							assertTreeParent(t, root, cl, reader, child, parent)
						} else {
							assertTreeParent(t, root, cl, reader, child, original)
						}
					})

					t.Run("detach", func(t *testing.T) {
						parent := createTreeOwnershipNode(t, root, nw, tc.parentOwner, visibility.VisibilityPublished)
						child := createTreeOwnershipNode(t, root, nw, tc.childOwner, tc.childVisibility, node_writer.WithParent(library.NodeID(parent.GetID())))
						assertTreeParent(t, root, cl, reader, child, parent)

						response, err := cl.NodeRemoveNodeWithResponse(root, parent.GetSlug(), child.GetSlug(), tc.actor)
						tests.Status(t, err, response, tc.severStatus)
						if tc.severStatus == http.StatusOK {
							assertTreeParent(t, root, cl, reader, child, nil)
						} else {
							assertTreeParent(t, root, cl, reader, child, parent)
						}
					})
				})
			}
		}))
	}))
}

func TestNodeParentMutationAuthorization(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		nw *node_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminSession := sh.WithSession(adminCtx)
			memberSession := sh.WithSession(memberCtx)
			otherSession := sh.WithSession(otherCtx)

			parents := []struct {
				name   string
				owner  account.AccountID
				status int
			}{
				{name: "own_destination", owner: member.ID, status: http.StatusOK},
				{name: "foreign_destination", owner: other.ID, status: http.StatusForbidden},
			}

			for _, destination := range parents {
				t.Run(destination.name, func(t *testing.T) {
					for _, operation := range []string{"update", "position", "delete_with_target"} {
						t.Run(operation, func(t *testing.T) {
							original := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished)
							parent := createTreeOwnershipNode(t, root, nw, destination.owner, visibility.VisibilityPublished)
							child := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished, node_writer.WithParent(library.NodeID(original.GetID())))
							grandchild := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished, node_writer.WithParent(library.NodeID(child.GetID())))
							assertTreeParent(t, root, cl, memberSession, child, original)
							parentSlug := parent.GetSlug()

							switch operation {
							case "update":
								response, err := cl.NodeUpdateWithResponse(root, child.GetSlug(), openapi.NodeUpdateJSONRequestBody{Parent: &parentSlug}, memberSession)
								tests.Status(t, err, response, destination.status)
							case "position":
								position := openapi.NodeUpdatePositionJSONRequestBody{}
								position.Parent.Set(parentSlug)
								response, err := cl.NodeUpdatePositionWithResponse(root, child.GetSlug(), position, memberSession)
								tests.Status(t, err, response, destination.status)
							case "delete_with_target":
								response, err := cl.NodeDeleteWithResponse(root, child.GetSlug(), &openapi.NodeDeleteParams{TargetNode: &parentSlug}, memberSession)
								tests.Status(t, err, response, destination.status)
							}

							if destination.status == http.StatusForbidden {
								assertTreeParent(t, root, cl, memberSession, child, original)
								assertTreeParent(t, root, cl, memberSession, grandchild, child)
								return
							}

							if operation == "delete_with_target" {
								deleted, err := cl.NodeGetWithResponse(root, child.GetSlug(), &openapi.NodeGetParams{}, memberSession)
								tests.Status(t, err, deleted, http.StatusNotFound)
								assertTreeParent(t, root, cl, memberSession, grandchild, parent)
								return
							}

							assertTreeParent(t, root, cl, memberSession, child, parent)
						})
					}
				})
			}

			t.Run("cannot_delete_parent_to_detach_another_members_page", func(t *testing.T) {
				for _, reparent := range []bool{false, true} {
					name := "to_root"
					if reparent {
						name = "to_owned_destination"
					}

					t.Run(name, func(t *testing.T) {
						parent := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished)
						child := createTreeOwnershipNode(t, root, nw, other.ID, visibility.VisibilityPublished, node_writer.WithParent(library.NodeID(parent.GetID())))
						params := &openapi.NodeDeleteParams{}
						if reparent {
							destination := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished)
							slug := destination.GetSlug()
							params.TargetNode = &slug
						}

						response, err := cl.NodeDeleteWithResponse(root, parent.GetSlug(), params, memberSession)
						tests.Status(t, err, response, http.StatusForbidden)
						assertTreeParent(t, root, cl, otherSession, child, parent)

						response, err = cl.NodeDeleteWithResponse(root, parent.GetSlug(), params, adminSession)
						tests.Status(t, err, response, http.StatusOK)
					})
				}
			})

			t.Run("draft_submission_requires_review_before_publication", func(t *testing.T) {
				parent := createTreeOwnershipNode(t, root, nw, other.ID, visibility.VisibilityPublished)
				parentSlug := parent.GetSlug()
				created := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name: "Submitted page " + xid.New().String(), Parent: &parentSlug,
				}, memberSession))(t, http.StatusOK)
				page := tests.AssertRequest(cl.NodeGetWithResponse(root, created.JSON200.Slug, &openapi.NodeGetParams{}, memberSession))(t, http.StatusOK)
				require.NotNil(t, page.JSON200.Parent)
				require.Equal(t, parent.GetID().String(), page.JSON200.Parent.Id)
				require.Equal(t, openapi.VisibilityDraft, page.JSON200.Visibility)

				published, err := cl.NodeUpdateVisibilityWithResponse(root, created.JSON200.Slug, openapi.NodeUpdateVisibilityJSONRequestBody{Visibility: openapi.VisibilityPublished}, memberSession)
				tests.Status(t, err, published, http.StatusForbidden)
				unchanged := tests.AssertRequest(cl.NodeGetWithResponse(root, created.JSON200.Slug, &openapi.NodeGetParams{}, memberSession))(t, http.StatusOK)
				require.Equal(t, openapi.VisibilityDraft, unchanged.JSON200.Visibility)

				reviewed, err := cl.NodeUpdateVisibilityWithResponse(root, created.JSON200.Slug, openapi.NodeUpdateVisibilityJSONRequestBody{Visibility: openapi.VisibilityReview}, memberSession)
				tests.Status(t, err, reviewed, http.StatusOK)
				published, err = cl.NodeUpdateVisibilityWithResponse(root, created.JSON200.Slug, openapi.NodeUpdateVisibilityJSONRequestBody{Visibility: openapi.VisibilityPublished}, adminSession)
				tests.Status(t, err, published, http.StatusOK)
			})

			t.Run("owner_can_edit_page_without_changing_managed_parent", func(t *testing.T) {
				parent := createTreeOwnershipNode(t, root, nw, other.ID, visibility.VisibilityPublished)
				child := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityPublished, node_writer.WithParent(library.NodeID(parent.GetID())))
				parentSlug := parent.GetSlug()
				description := "An owner can update their page without moving it."
				updated := tests.AssertRequest(cl.NodeUpdateWithResponse(root, child.GetSlug(), openapi.NodeUpdateJSONRequestBody{
					Parent: &parentSlug, Description: &description,
				}, memberSession))(t, http.StatusOK)
				require.Equal(t, description, updated.JSON200.Description)
				assertTreeParent(t, root, cl, memberSession, child, parent)

				unchanged, err := cl.NodeUpdateVisibilityWithResponse(root, child.GetSlug(), openapi.NodeUpdateVisibilityJSONRequestBody{Visibility: openapi.VisibilityPublished}, memberSession)
				tests.Status(t, err, unchanged, http.StatusOK)
			})

			t.Run("cannot_submit_to_another_members_private_page", func(t *testing.T) {
				parent := createTreeOwnershipNode(t, root, nw, other.ID, visibility.VisibilityDraft)
				parentSlug := parent.GetSlug()
				created, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name: "Private parent submission " + xid.New().String(), Parent: &parentSlug,
				}, memberSession)
				tests.Status(t, err, created, http.StatusForbidden)
			})

			t.Run("cannot_reverse_relationship_to_detach_another_members_page", func(t *testing.T) {
				draft := createTreeOwnershipNode(t, root, nw, member.ID, visibility.VisibilityDraft)
				published := createTreeOwnershipNode(t, root, nw, other.ID, visibility.VisibilityPublished, node_writer.WithParent(library.NodeID(draft.GetID())))

				response, err := cl.NodeAddNodeWithResponse(root, published.GetSlug(), draft.GetSlug(), memberSession)
				tests.Status(t, err, response, http.StatusForbidden)
				assertTreeParent(t, root, cl, memberSession, published, draft)
				assertTreeParent(t, root, cl, memberSession, draft, nil)

				response, err = cl.NodeAddNodeWithResponse(root, published.GetSlug(), draft.GetSlug(), adminSession)
				tests.Status(t, err, response, http.StatusOK)
				assertTreeParent(t, root, cl, memberSession, draft, published)
				assertTreeParent(t, root, cl, otherSession, published, nil)
			})
		}))
	}))
}

func createTreeOwnershipNode(t *testing.T, ctx context.Context, writer *node_writer.Writer, owner account.AccountID, vis visibility.Visibility, opts ...node_writer.Option) *library.Node {
	t.Helper()

	name := "tree-ownership-" + xid.New().String()
	opts = append(opts, node_writer.WithVisibility(vis))
	node, err := writer.Create(ctx, owner, name, mark.NewSlugFromName(name), opts...)
	require.NoError(t, err)
	return node
}

func assertTreeParent(t *testing.T, ctx context.Context, client *openapi.ClientWithResponses, actor openapi.RequestEditorFn, child, parent *library.Node) {
	t.Helper()

	response := tests.AssertRequest(client.NodeGetWithResponse(ctx, child.GetSlug(), &openapi.NodeGetParams{}, actor))(t, http.StatusOK)
	if parent == nil {
		require.Nil(t, response.JSON200.Parent)
		return
	}

	require.NotNil(t, response.JSON200.Parent)
	require.Equal(t, parent.GetID().String(), response.JSON200.Parent.Id)
}
