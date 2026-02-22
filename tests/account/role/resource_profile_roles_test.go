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
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestResourceProfileReferencesIncludeRoles(t *testing.T) {
	// Intentionally not parallel: this suite mutates default roles and can race
	// with other role tests that run against the shared CI database.
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
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)

			roleMeta := openapi.Metadata{
				"style": "bold",
			}
			roleName := "resource-role-" + xid.New().String()
			roleColour := "#f97316"

			roleCreate := tests.AssertRequest(cl.RoleCreateWithResponse(authorCtx, openapi.RoleCreateJSONRequestBody{
				Name:        roleName,
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

			memberThreadCreate := tests.AssertRequest(cl.ThreadCreateWithResponse(memberCtx, openapi.ThreadInitialProps{
				Title:      "member-thread-with-default-role-" + xid.New().String(),
				Body:       opt.New("<p>member thread body</p>").Ptr(),
				Category:   opt.New(threadCategory.JSON200.Id).Ptr(),
				Visibility: &vis,
			}, memberSession))(t, http.StatusOK)

			nodeCreate := tests.AssertRequest(cl.NodeCreateWithResponse(authorCtx, openapi.NodeCreateJSONRequestBody{
				Name:       "node-with-roles-" + xid.New().String(),
				Visibility: &vis,
			}, authorSession))(t, http.StatusOK)

			collectionCreate := tests.AssertRequest(cl.CollectionCreateWithResponse(authorCtx, openapi.CollectionCreateJSONRequestBody{
				Name: "collection-with-roles-" + xid.New().String(),
			}, authorSession))(t, http.StatusOK)
			tests.AssertRequest(cl.CollectionAddPostWithResponse(authorCtx, collectionCreate.JSON200.Id, threadCreate.JSON200.Id, authorSession))(t, http.StatusOK)
			tests.AssertRequest(cl.CollectionAddNodeWithResponse(authorCtx, collectionCreate.JSON200.Id, nodeCreate.JSON200.Id, authorSession))(t, http.StatusOK)

			assertRoleRef := func(
				t *testing.T,
				where string,
				roles []openapi.AccountRoleRef,
				id openapi.Identifier,
				name string,
				colour string,
				meta *openapi.Metadata,
				badge bool,
				isDefault bool,
			) *openapi.AccountRoleRef {
				t.Helper()
				r := require.New(t)
				a := assert.New(t)

				found := findRoleRef(roles, id)
				r.NotNil(found, where)
				a.Equal(id, found.Id, where)
				a.Equal(name, found.Name, where)
				a.Equal(colour, found.Colour, where)
				a.Equal(badge, found.Badge, where)
				a.Equal(isDefault, found.Default, where)
				if meta == nil {
					a.Nil(found.Meta, where)
				} else {
					r.NotNil(found.Meta, where)
					a.Equal(*meta, *found.Meta, where)
				}

				return found
			}

			assertProfileRefHasCustomRole := func(t *testing.T, where string, roles []openapi.AccountRoleRef) {
				t.Helper()
				assertRoleRef(
					t,
					where,
					roles,
					roleCreate.JSON200.Id,
					roleName,
					roleColour,
					&roleMeta,
					false,
					false,
				)
			}

			t.Run("thread_get_author", func(t *testing.T) {
				threadGet := tests.AssertRequest(cl.ThreadGetWithResponse(authorCtx, threadCreate.JSON200.Slug, nil, authorSession))(t, http.StatusOK)
				assertProfileRefHasCustomRole(t, "thread.get author should include assigned custom role", threadGet.JSON200.Author.Roles)
			})

			t.Run("thread_get_replies_author", func(t *testing.T) {
				r := require.New(t)
				threadGet := tests.AssertRequest(cl.ThreadGetWithResponse(authorCtx, threadCreate.JSON200.Slug, nil, authorSession))(t, http.StatusOK)
				replyItem, found := lo.Find(threadGet.JSON200.Replies.Replies, func(in openapi.Reply) bool {
					return in.Id == replyCreate.JSON200.Id
				})
				r.True(found, "created reply should appear in thread replies")
				assertProfileRefHasCustomRole(t, "thread.get reply author should include assigned custom role", replyItem.Author.Roles)
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
				assertProfileRefHasCustomRole(t, "thread.list author should include assigned custom role", threadListItem.Author.Roles)
			})

			t.Run("thread_list_admin_default_role", func(t *testing.T) {
				r := require.New(t)

				authorFilter := openapi.AccountHandle(author.Handle)
				threadList := tests.AssertRequest(cl.ThreadListWithResponse(authorCtx, &openapi.ThreadListParams{
					Author: &authorFilter,
				}, authorSession))(t, http.StatusOK)
				threadListItem, found := lo.Find(threadList.JSON200.Threads, func(in openapi.ThreadReference) bool {
					return in.Id == threadCreate.JSON200.Id
				})
				r.True(found, "created admin thread should appear in thread list")

				assertRoleRef(
					t,
					"thread.list author should include admin default role",
					threadListItem.Author.Roles,
					openapi.Identifier(role.DefaultRoleAdminID.String()),
					role.DefaultRoleAdmin.Name,
					role.DefaultRoleAdmin.Colour,
					nil,
					false,
					true,
				)
			})

			t.Run("thread_list_member_default_role", func(t *testing.T) {
				r := require.New(t)

				authorFilter := openapi.AccountHandle(member.Handle)
				threadList := tests.AssertRequest(cl.ThreadListWithResponse(authorCtx, &openapi.ThreadListParams{
					Author: &authorFilter,
				}, authorSession))(t, http.StatusOK)
				threadListItem, found := lo.Find(threadList.JSON200.Threads, func(in openapi.ThreadReference) bool {
					return in.Id == memberThreadCreate.JSON200.Id
				})
				r.True(found, "created member thread should appear in thread list")

				assertRoleRef(
					t,
					"member default role should be present on profile reference",
					threadListItem.Author.Roles,
					openapi.Identifier(role.DefaultRoleMemberID.String()),
					role.DefaultRoleMember.Name,
					role.DefaultRoleMember.Colour,
					nil,
					false,
					true,
				)
			})

			t.Run("reply_create_author", func(t *testing.T) {
				assertProfileRefHasCustomRole(t, "reply.create author should include assigned custom role", replyCreate.JSON200.Author.Roles)
			})

			t.Run("node_get_owner", func(t *testing.T) {
				nodeGet := tests.AssertRequest(cl.NodeGetWithResponse(authorCtx, nodeCreate.JSON200.Slug, &openapi.NodeGetParams{}, authorSession))(t, http.StatusOK)
				assertProfileRefHasCustomRole(t, "node.get owner should include assigned custom role", nodeGet.JSON200.Owner.Roles)
			})

			t.Run("node_list_owner", func(t *testing.T) {
				r := require.New(t)
				nodeList := tests.AssertRequest(cl.NodeListWithResponse(authorCtx, &openapi.NodeListParams{}, authorSession))(t, http.StatusOK)
				nodeItem, found := lo.Find(nodeList.JSON200.Nodes, func(in openapi.NodeWithChildren) bool {
					return in.Id == nodeCreate.JSON200.Id
				})
				r.True(found, "created node should appear in node list")
				assertProfileRefHasCustomRole(t, "node.list owner should include assigned custom role", nodeItem.Owner.Roles)
			})

			t.Run("collection_get_owner", func(t *testing.T) {
				collectionGet := tests.AssertRequest(cl.CollectionGetWithResponse(authorCtx, collectionCreate.JSON200.Id, authorSession))(t, http.StatusOK)
				assertProfileRefHasCustomRole(t, "collection.get owner should include assigned custom role", collectionGet.JSON200.Owner.Roles)
			})

			t.Run("collection_list_owner", func(t *testing.T) {
				r := require.New(t)
				collectionList := tests.AssertRequest(cl.CollectionListWithResponse(authorCtx, &openapi.CollectionListParams{}, authorSession))(t, http.StatusOK)
				collectionItem, found := lo.Find(collectionList.JSON200.Collections, func(in openapi.Collection) bool {
					return in.Id == collectionCreate.JSON200.Id
				})
				r.True(found, "created collection should appear in collection list")
				assertProfileRefHasCustomRole(t, "collection.list owner should include assigned custom role", collectionItem.Owner.Roles)
			})

			t.Run("collection_get_items_owner_and_item_refs", func(t *testing.T) {
				r := require.New(t)
				collectionGet := tests.AssertRequest(cl.CollectionGetWithResponse(authorCtx, collectionCreate.JSON200.Id, authorSession))(t, http.StatusOK)

				var foundThreadItem bool
				var foundNodeItem bool

				for _, item := range collectionGet.JSON200.Items {
					assertProfileRefHasCustomRole(t, "collection item owner should include assigned custom role", item.Owner.Roles)

					postItem, postErr := item.Item.AsDatagraphItemPost()
					if postErr == nil && postItem.Ref.Id == threadCreate.JSON200.Id {
						foundThreadItem = true
						assertProfileRefHasCustomRole(t, "collection item post author should include assigned custom role", postItem.Ref.Author.Roles)
					}

					nodeItem, nodeErr := item.Item.AsDatagraphItemNode()
					if nodeErr == nil && nodeItem.Ref.Id == nodeCreate.JSON200.Id {
						foundNodeItem = true
						assertProfileRefHasCustomRole(t, "collection item node owner should include assigned custom role", nodeItem.Ref.Owner.Roles)
					}
				}

				r.True(foundThreadItem, "expected thread item in collection")
				r.True(foundNodeItem, "expected node item in collection")
			})

			t.Run("default_role_overrides_are_reflected", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)
				toPerms := func(perms []rbac.Permission) openapi.PermissionList {
					out := make(openapi.PermissionList, 0, len(perms))
					for _, perm := range perms {
						out = append(out, openapi.Permission(perm.String()))
					}
					return out
				}

				memberName := "member-override-" + xid.New().String()
				memberColour := "#22c55e"
				memberMeta := openapi.Metadata{"style": "member-gradient"}
				memberPerms := toPerms(role.DefaultRoleMember.Permissions.List())

				adminName := "admin-override-" + xid.New().String()
				adminColour := "#ef4444"
				adminMeta := openapi.Metadata{"style": "admin-glow"}

				guestName := "guest-override-" + xid.New().String()
				guestColour := "#64748b"
				guestMeta := openapi.Metadata{"style": "guest-muted"}
				guestPerms := toPerms(role.DefaultRoleGuest.Permissions.List())

				t.Cleanup(func() {
					assertDelete := func(roleID string) {
						resp, err := cl.RoleDeleteWithResponse(authorCtx, roleID, authorSession)
						if err != nil {
							t.Errorf("cleanup failed deleting default role %s: %v", roleID, err)
							return
						}
						if resp == nil {
							t.Errorf("cleanup failed deleting default role %s: nil response", roleID)
							return
						}
						if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNotFound {
							t.Errorf("cleanup failed deleting default role %s: unexpected status %d", roleID, resp.StatusCode())
						}
					}
					assertDelete(role.DefaultRoleGuestID.String())
					assertDelete(role.DefaultRoleMemberID.String())
					assertDelete(role.DefaultRoleAdminID.String())
				})

				tests.AssertRequest(cl.RoleUpdateWithResponse(authorCtx, role.DefaultRoleMemberID.String(), openapi.RoleUpdateJSONRequestBody{
					Name:        &memberName,
					Colour:      &memberColour,
					Meta:        &memberMeta,
					Permissions: &memberPerms,
				}, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RoleUpdateWithResponse(authorCtx, role.DefaultRoleAdminID.String(), openapi.RoleUpdateJSONRequestBody{
					Name:        &adminName,
					Colour:      &adminColour,
					Meta:        &adminMeta,
					Permissions: lo.ToPtr(toPerms(role.DefaultRoleAdmin.Permissions.List())),
				}, authorSession))(t, http.StatusBadRequest)
				tests.AssertRequest(cl.RoleUpdateWithResponse(authorCtx, role.DefaultRoleAdminID.String(), openapi.RoleUpdateJSONRequestBody{
					Name:   &adminName,
					Colour: &adminColour,
					Meta:   &adminMeta,
				}, authorSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RoleUpdateWithResponse(authorCtx, role.DefaultRoleGuestID.String(), openapi.RoleUpdateJSONRequestBody{
					Name:        &guestName,
					Colour:      &guestColour,
					Meta:        &guestMeta,
					Permissions: &guestPerms,
				}, authorSession))(t, http.StatusOK)

				memberFilter := openapi.AccountHandle(member.Handle)
				memberThreadList := tests.AssertRequest(cl.ThreadListWithResponse(authorCtx, &openapi.ThreadListParams{
					Author: &memberFilter,
				}, authorSession))(t, http.StatusOK)
				memberThreadListItem, found := lo.Find(memberThreadList.JSON200.Threads, func(in openapi.ThreadReference) bool {
					return in.Id == memberThreadCreate.JSON200.Id
				})
				r.True(found, "created member thread should appear in thread list")
				assertRoleRef(
					t,
					"member default override should appear on profile reference",
					memberThreadListItem.Author.Roles,
					openapi.Identifier(role.DefaultRoleMemberID.String()),
					memberName,
					memberColour,
					&memberMeta,
					false,
					true,
				)

				adminFilter := openapi.AccountHandle(author.Handle)
				adminThreadList := tests.AssertRequest(cl.ThreadListWithResponse(authorCtx, &openapi.ThreadListParams{
					Author: &adminFilter,
				}, authorSession))(t, http.StatusOK)
				adminThreadListItem, found := lo.Find(adminThreadList.JSON200.Threads, func(in openapi.ThreadReference) bool {
					return in.Id == threadCreate.JSON200.Id
				})
				r.True(found, "created admin thread should appear in thread list")
				assertRoleRef(
					t,
					"admin default override should appear on profile reference",
					adminThreadListItem.Author.Roles,
					openapi.Identifier(role.DefaultRoleAdminID.String()),
					adminName,
					adminColour,
					&adminMeta,
					false,
					true,
				)

				roleList := tests.AssertRequest(cl.RoleListWithResponse(authorCtx, authorSession))(t, http.StatusOK)
				adminRole, found := lo.Find(roleList.JSON200.Roles, func(in openapi.Role) bool {
					return in.Id == role.DefaultRoleAdminID.String()
				})
				r.True(found, "admin default role should be present in role list")
				a.Equal(openapi.PermissionList{openapi.ADMINISTRATOR}, adminRole.Permissions)

				guestRole, found := lo.Find(roleList.JSON200.Roles, func(in openapi.Role) bool {
					return in.Id == role.DefaultRoleGuestID.String()
				})
				r.True(found, "guest default role should be present in role list")
				a.Equal(guestName, guestRole.Name)
				a.Equal(guestColour, guestRole.Colour)
				r.NotNil(guestRole.Meta)
				a.Equal(guestMeta, *guestRole.Meta)
			})
		}))
	}))
}
