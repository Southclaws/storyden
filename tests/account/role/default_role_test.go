package role_test

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
	"github.com/Southclaws/storyden/tests"
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
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name: "TestGuestRolePermissions" + uuid.NewString(),
			}, adminSession))(t, http.StatusOK)

			// Helper values for pointers
			published := openapi.Published
			review := openapi.Review
			content := "<body>This is a test node.</body>"

			edit, err := cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleEveryoneID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{},
				}, adminSession)
			tests.Ok(t, err, edit)

			t.Run("guest_cannot_create_post", func(t *testing.T) {
				// PermissionCreatePost
				tests.AssertRequest(
					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_read_published_threads", func(t *testing.T) {
				// PermissionReadPublishedThreads
				tests.AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_reaction", func(t *testing.T) {
				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "guest_cannot_create_reaction" + uuid.NewString(),
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				// PermissionCreateReaction
				tests.AssertRequest(
					cl.PostReactAddWithResponse(root, thread.JSON200.Slug, openapi.PostReactAddJSONRequestBody{
						Emoji: "👍",
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_read_published_library", func(t *testing.T) {
				node := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "guest_cannot_read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)
				// PermissionReadPublishedLibrary
				tests.AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}),
				)(t, http.StatusForbidden)
				tests.AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_node", func(t *testing.T) {
				// PermissionSubmitLibraryNode
				tests.AssertRequest(
					cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
						Name:       "guest_cannot_create_node" + uuid.NewString(),
						Content:    &content,
						Visibility: &review,
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_upload_asset", func(t *testing.T) {
				// PermissionUploadAsset
				tests.AssertRequest(
					cl.AssetUploadWithBodyWithResponse(root, &openapi.AssetUploadParams{
						ContentLength: 69,
					}, "application/whatever", bytes.NewBuffer([]byte("test"))),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_list_profiles", func(t *testing.T) {
				// PermissionListProfiles
				tests.AssertRequest(
					cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_read_profile", func(t *testing.T) {
				// PermissionReadProfile
				tests.AssertRequest(
					cl.ProfileGetWithResponse(root, seed.Account_001_Odin.Handle),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_collection", func(t *testing.T) {
				// PermissionCreateCollection
				tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_create_collection" + uuid.NewString(),
					}),
				)(t, http.StatusUnauthorized)
			})

			t.Run("guest_cannot_list_collections", func(t *testing.T) {
				// PermissionListCollections
				tests.AssertRequest(
					cl.CollectionListWithResponse(root, &openapi.CollectionListParams{}),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_read_collection", func(t *testing.T) {
				col := tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_read_collection" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				// PermissionReadCollection
				tests.AssertRequest(
					cl.CollectionGetWithResponse(root, col.JSON200.Slug),
				)(t, http.StatusForbidden)
			})

			t.Run("guest_cannot_create_collection_item", func(t *testing.T) {
				col := tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "guest_cannot_create_collection_item" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "guest_cannot_create_collection_item",
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				node := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "guest_cannot_create_collection_item" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// PermissionCollectionSubmit
				tests.AssertRequest(
					cl.CollectionAddPostWithResponse(root, col.JSON200.Slug, thread.JSON200.Id),
				)(t, http.StatusUnauthorized)
				tests.AssertRequest(
					cl.CollectionAddNodeWithResponse(root, col.JSON200.Slug, node.JSON200.Id),
				)(t, http.StatusUnauthorized)
			})

			tests.AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleEveryoneID.String(),
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
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
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
				role.DefaultRoleEveryoneID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{},
				}, adminSession)
			tests.Ok(t, err, edit)

			t.Run("member_cannot_create_post", func(t *testing.T) {
				// PermissionCreatePost
				tests.AssertRequest(
					cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_published_threads", func(t *testing.T) {
				// PermissionReadPublishedThreads
				tests.AssertRequest(
					cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_reaction", func(t *testing.T) {
				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "member_cannot_create_reaction" + uuid.NewString(),
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				// PermissionCreateReaction
				tests.AssertRequest(
					cl.PostReactAddWithResponse(root, thread.JSON200.Slug, openapi.PostReactAddJSONRequestBody{
						Emoji: "👍",
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_published_library", func(t *testing.T) {
				node := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "member_cannot_read_published_library" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)
				// PermissionReadPublishedLibrary
				tests.AssertRequest(
					cl.NodeListWithResponse(root, &openapi.NodeListParams{}, member1Session),
				)(t, http.StatusForbidden)
				tests.AssertRequest(
					cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_node", func(t *testing.T) {
				// PermissionSubmitLibraryNode
				tests.AssertRequest(
					cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
						Name:       "member_cannot_create_node" + uuid.NewString(),
						Content:    &content,
						Visibility: &review,
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_upload_asset", func(t *testing.T) {
				// PermissionUploadAsset
				tests.AssertRequest(
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
				tests.AssertRequest(
					cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_profile", func(t *testing.T) {
				// PermissionReadProfile
				tests.AssertRequest(
					cl.ProfileGetWithResponse(root, seed.Account_001_Odin.Handle, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_collection", func(t *testing.T) {
				// PermissionCreateCollection
				tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_create_collection" + uuid.NewString(),
					}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_list_collections", func(t *testing.T) {
				// PermissionListCollections
				tests.AssertRequest(
					cl.CollectionListWithResponse(root, &openapi.CollectionListParams{}, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_read_collection", func(t *testing.T) {
				col := tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_read_collection" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				// PermissionReadCollection
				tests.AssertRequest(
					cl.CollectionGetWithResponse(root, col.JSON200.Slug, member1Session),
				)(t, http.StatusForbidden)
			})

			t.Run("member_cannot_create_collection_item", func(t *testing.T) {
				col := tests.AssertRequest(
					cl.CollectionCreateWithResponse(root, openapi.CollectionInitialProps{
						Name: "member_cannot_create_collection_item" + uuid.NewString(),
					}, adminSession),
				)(t, http.StatusOK)

				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "member_cannot_create_collection_item",
					Body:       "<body>This is a test thread.</body>",
					Visibility: openapi.Published,
					Category:   cat.JSON200.Id,
				}, adminSession))(t, http.StatusOK)

				node := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "member_cannot_create_collection_item" + uuid.NewString(),
					Content:    &content,
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// PermissionCollectionSubmit
				tests.AssertRequest(
					cl.CollectionAddPostWithResponse(root, col.JSON200.Slug, thread.JSON200.Id, member1Session),
				)(t, http.StatusForbidden)
				tests.AssertRequest(
					cl.CollectionAddNodeWithResponse(root, col.JSON200.Slug, node.JSON200.Id, member1Session),
				)(t, http.StatusForbidden)
			})

			tests.AssertRequest(cl.RoleDeleteWithResponse(adminCtx,
				role.DefaultRoleEveryoneID.String(),
				adminSession),
			)(t, http.StatusOK)
		}))
	}))
}
