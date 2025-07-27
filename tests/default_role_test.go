package tests

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

// NOTE: The generated client takes an optional final parameter list of "request
// modifiers" these are used throughout tests primarily to add "session" objects
// to declare who is making the request. In these tests, we make use of these on
// both member and guest checks. Guest requests include no final argument so the
// requests are completely unauthenticated. With that said, mutative requests on
// resources will always yield an Unauthorised, while most read requests respond
// with a Forbidden. Certain read operations may return an Unauthorised response

// NOTE 2: These tests mutate global state, by modifying a role that's shared by
// the system as a whole (the Member role) so they cannoy be parallel and also a
// cleanup must be run after the tests finish to delete the modified member role

func TestGuestRolePermissions(t *testing.T) {
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			cat := AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name: "TestGuestRolePermissions" + uuid.NewString(),
			}, adminSession))(t, http.StatusOK)

			// Helper values for pointers
			published := openapi.Published
			review := openapi.Review
			content := "<body>This is a test node.</body>"

			// Remove all permissions from the guest role to test restrictions
			edit, err := cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{},
				}, adminSession)
			Ok(t, err, edit)

			t.Run("guest_cannot_create_post", func(t *testing.T) {
				// PermissionCreatePost
				AssertRequest(
					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_read_published_threads", func(t *testing.T) {
				// PermissionReadPublishedThreads
				AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_reaction", func(t *testing.T) {
				thread := AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "guest_cannot_create_reaction" + uuid.NewString(),
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				// PermissionCreateReaction
				AssertRequest(
					cl.PostReactAddWithResponse(root, thread.JSON200.Slug, openapi.PostReactAddJSONRequestBody{
						Emoji: "👍",
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_read_published_library", func(t *testing.T) {
				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "guest_cannot_read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)
				// PermissionReadPublishedLibrary
				AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}),
				)(t, http.StatusForbidden)
				AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_node", func(t *testing.T) {
				// PermissionSubmitLibraryNode
				AssertRequest(
					cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
						Name:       "guest_cannot_create_node" + uuid.NewString(),
						Content:    &content,
						Visibility: &review,
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_upload_asset", func(t *testing.T) {
				// PermissionUploadAsset
				AssertRequest(
					cl.AssetUploadWithBodyWithResponse(root, &openapi.AssetUploadParams{
						ContentLength: 69,
					}, "application/whatever", bytes.NewBuffer([]byte("test"))),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_list_profiles", func(t *testing.T) {
				// PermissionListProfiles
				AssertRequest(
					cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_read_profile", func(t *testing.T) {
				// PermissionReadProfile
				AssertRequest(
					cl.ProfileGetWithResponse(root, admin.Handle),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_collection", func(t *testing.T) {
				// PermissionCreateCollection
				AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_create_collection" + uuid.NewString(),
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_list_collections", func(t *testing.T) {
				// PermissionListCollections
				AssertRequest(
					cl.CollectionListWithResponse(root, &openapi.CollectionListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_read_collection", func(t *testing.T) {
				col := AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_read_collection" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				// PermissionReadCollection
				AssertRequest(
					cl.CollectionGetWithResponse(root, col.JSON200.Slug),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_collection_item", func(t *testing.T) {
				col := AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_create_collection_item" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				thread := AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "guest_cannot_create_collection_item",
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "guest_cannot_create_collection_item" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// PermissionCollectionSubmit
				AssertRequest(
					cl.CollectionAddPostWithResponse(root, col.JSON200.Slug, thread.JSON200.Id),
				)(t, http.StatusUnauthorized)
				AssertRequest(
					cl.CollectionAddNodeWithResponse(root, col.JSON200.Slug, node.JSON200.Id),
				)(t, http.StatusUnauthorized)
			})

			AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				adminSession),
			)(t, http.StatusOK)
		}))
	}))
}

func TestGuestRoleWithPermissions(t *testing.T) {
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			cat := AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name: "TestGuestRoleWithPermissions" + uuid.NewString(),
			}, adminSession))(t, http.StatusOK)

			// Helper values for pointers
			published := openapi.Published
			content := "<body>This is a test node.</body>"

			// Grant read permissions to guest role
			edit, err := cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{
						"READ_PUBLISHED_THREADS",
						"READ_PUBLISHED_LIBRARY",
						"LIST_PROFILES",
						"READ_PROFILE",
						"LIST_COLLECTIONS",
						"READ_COLLECTION",
					},
				}, adminSession)
			Ok(t, err, edit)

			t.Run("guest_can_read_published_threads", func(t *testing.T) {
				thread := AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "guest_can_read_published_threads" + uuid.NewString(),
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				// PermissionReadPublishedThreads
				AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}),
				)(t, http.StatusOK)
				AssertRequest(
					cl.ThreadGetWithResponse(root, thread.JSON200.Slug, &openapi.ThreadGetParams{}),
				)(t, http.StatusOK)
			})

			t.Run("guest_can_read_published_library", func(t *testing.T) {
				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "guest_can_read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// PermissionReadPublishedLibrary
				AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}),
				)(t, http.StatusOK)
				AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}),
				)(t, http.StatusOK)
			})

			t.Run("guest_can_list_profiles", func(t *testing.T) {
				// PermissionListProfiles
				AssertRequest(
					cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}),
				)(t, http.StatusOK)
			})

			t.Run("guest_can_read_profile", func(t *testing.T) {
				// PermissionReadProfile
				AssertRequest(
					cl.ProfileGetWithResponse(root, admin.Handle),
				)(t, http.StatusOK)
			})

			t.Run("guest_can_list_collections", func(t *testing.T) {
				// PermissionListCollections
				AssertRequest(
					cl.CollectionListWithResponse(root, &openapi.CollectionListParams{}),
				)(t, http.StatusOK)
			})

			t.Run("guest_can_read_collection", func(t *testing.T) {
				col := AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_can_read_collection" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				// PermissionReadCollection
				AssertRequest(
					cl.CollectionGetWithResponse(root, col.JSON200.Slug),
				)(t, http.StatusOK)
			})

			// Clean up by deleting the guest role customization
			AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				adminSession),
			)(t, http.StatusOK)
		}))
	}))
}

func TestMemberRolePermissions(t *testing.T) {
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			cat := AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name: "TestMemberRolePermissions" + uuid.NewString(),
			}, adminSession))(t, http.StatusOK)

			// Helper values for pointers
			published := openapi.Published
			review := openapi.Review
			content := "<body>This is a test node.</body>"

			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			member1Session := sh.WithSession(memberCtx)

			// Remove all permissions from the default member role.
			edit, err := cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleMemberID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{},
				}, adminSession)
			Ok(t, err, edit)

			t.Run("member_cannot_create_post", func(t *testing.T) {
				// PermissionCreatePost
				AssertRequest(
					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_published_threads", func(t *testing.T) {
				// PermissionReadPublishedThreads
				AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_reaction", func(t *testing.T) {
				thread := AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "member_cannot_create_reaction" + uuid.NewString(),
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				// PermissionCreateReaction
				AssertRequest(
					cl.PostReactAddWithResponse(root, thread.JSON200.Slug, openapi.PostReactAddJSONRequestBody{
						Emoji: "👍",
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_published_library", func(t *testing.T) {
				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "member_cannot_read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)
				// PermissionReadPublishedLibrary
				AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}, member1Session),
				)(t, http.StatusForbidden)
				AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_node", func(t *testing.T) {
				// PermissionSubmitLibraryNode
				AssertRequest(
					cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
						Name:       "member_cannot_create_node" + uuid.NewString(),
						Content:    &content,
						Visibility: &review,
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_upload_asset", func(t *testing.T) {
				// PermissionUploadAsset
				AssertRequest(
					cl.AssetUploadWithBodyWithResponse(root,
						&openapi.AssetUploadParams{
							ContentLength: 69,
						},
						"application/whatever",
						bytes.NewBuffer([]byte("test")),
						member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_list_profiles", func(t *testing.T) {
				// PermissionListProfiles
				AssertRequest(
					cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_profile", func(t *testing.T) {
				// PermissionReadProfile
				AssertRequest(
					cl.ProfileGetWithResponse(root, admin.Handle, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_collection", func(t *testing.T) {
				// PermissionCreateCollection
				AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_create_collection" + uuid.NewString(),
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_list_collections", func(t *testing.T) {
				// PermissionListCollections
				AssertRequest(
					cl.CollectionListWithResponse(root, &openapi.CollectionListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_collection", func(t *testing.T) {
				col := AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_read_collection" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				// PermissionReadCollection
				AssertRequest(
					cl.CollectionGetWithResponse(root, col.JSON200.Slug, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_collection_item", func(t *testing.T) {
				col := AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_create_collection_item" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				thread := AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "member_cannot_create_collection_item",
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "member_cannot_create_collection_item" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// PermissionCollectionSubmit
				AssertRequest(
					cl.CollectionAddPostWithResponse(root, col.JSON200.Slug, thread.JSON200.Id, member1Session),
				)(t, http.StatusForbidden)
				AssertRequest(
					cl.CollectionAddNodeWithResponse(root, col.JSON200.Slug, node.JSON200.Id, member1Session),
				)(t, http.StatusForbidden)
			})

			AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleMemberID.String(),
				adminSession),
			)(t, http.StatusOK)
		}))
	}))
}

func TestGuestVsMemberAccess(t *testing.T) {
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

			cat := AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name: "TestMemberRolePermissions" + uuid.NewString(),
			}, adminSession))(t, http.StatusOK)

			AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "guest_cannot_create_reaction" + uuid.NewString(),
				Body:       "<body>This is a test thread.</body>",
				Visibility: openapi.Published,
				Category:   cat.JSON200.Id,
			}, adminSession))(t, http.StatusOK)

			// Helper values for pointers
			published := openapi.Published
			// review := openapi.Review
			content := "<body>This is a test node.</body>"

			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			member1Session := sh.WithSession(memberCtx)

			// Set explicit permissions for Guest
			AssertRequest(cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{
						"READ_PUBLISHED_THREADS",
					},
				}, adminSession))(t, http.StatusOK)

			// Set explicit permissions for Member
			AssertRequest(cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleMemberID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{
						"READ_PUBLISHED_LIBRARY",
					},
				}, adminSession))(t, http.StatusOK)

			t.Run("read_published_library", func(t *testing.T) {
				node := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// Guest can read threads, but not library

				// PermissionReadPublishedThreads
				AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}),
				)(t, http.StatusOK)
				// PermissionReadPublishedLibrary
				AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}),
				)(t, http.StatusForbidden)
				AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}),
				)(t, http.StatusForbidden)

				// Member can read both

				AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, member1Session),
				)(t, http.StatusForbidden)
				// PermissionReadPublishedLibrary
				AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}, member1Session),
				)(t, http.StatusOK)
				AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}, member1Session),
				)(t, http.StatusOK)
			})

			AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				adminSession),
			)(t, http.StatusOK)

			AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleMemberID.String(),
				adminSession),
			)(t, http.StatusOK)
		}))
	}))
}

func TestCannotGrantWritePermissionsToGuest(t *testing.T) {
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

			// Set explicit permissions for Guest
			AssertRequest(cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{
						openapi.CREATEPOST,
					},
				}, adminSession))(t, http.StatusBadRequest)
		}))
	}))
}
