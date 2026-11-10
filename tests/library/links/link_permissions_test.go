package links_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_writer"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entlink "github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestLinkPermissions(t *testing.T) {
	// This test changes built-in role permissions, so it must not run in parallel.
	var (
		root     context.Context
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
		roles    *role_writer.Writer
		links    *link_writer.LinkWriter
		nodes    *node_mutate.Manager
		db       *ent.Client
	)
	script := &scriptedScraper{}
	script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
		return &scrape.WebContent{Title: "Saved bookmark"}, nil
	})
	integration.Test(t, nil, e2e.Setup(),
		fx.Decorate(func(scrape.Scraper) scrape.Scraper { return script }),
		fx.Populate(&root, &client, &sessions, &accounts, &roles, &links, &nodes, &db),
	)

	memberBefore, err := roles.Get(root, role.DefaultRoleMemberID)
	require.NoError(t, err)
	guestBefore, err := roles.Get(root, role.DefaultRoleGuestID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := roles.Update(root, role.DefaultRoleMemberID, role_writer.WithPermissions(memberBefore.Permissions.List()))
		require.NoError(t, err)
		_, err = roles.Update(root, role.DefaultRoleGuestID, role_writer.WithPermissions(guestBefore.Permissions.List()))
		require.NoError(t, err)
	})

	ctx, owner := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
	member := sessions.WithSession(ctx)
	raw := "https://" + uuid.NewString() + ".example.com/"
	u := parseLinkURL(t, raw)
	bookmark, err := links.Store(ctx, u, "A saved bookmark", "")
	require.NoError(t, err)

	t.Run("guest_can_view_but_cannot_save_or_delete", func(t *testing.T) {
		tests.AssertRequest(client.LinkListWithResponse(root, nil))(t, http.StatusOK)
		tests.AssertRequest(client.LinkGetWithResponse(root, bookmark.Slug))(t, http.StatusOK)
		tests.AssertRequest(client.LinkCreateWithResponse(root, openapi.LinkCreateJSONRequestBody{Url: raw}))(t, http.StatusForbidden)
		tests.AssertRequest(client.LinkDeleteWithResponse(root, bookmark.Slug))(t, http.StatusUnauthorized)
	})

	t.Run("default_member_use_includes_view_but_not_delete", func(t *testing.T) {
		tests.AssertRequest(client.LinkListWithResponse(ctx, nil, member))(t, http.StatusOK)
		tests.AssertRequest(client.LinkGetWithResponse(ctx, bookmark.Slug, member))(t, http.StatusOK)
		tests.AssertRequest(client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, member))(t, http.StatusOK)
		tests.AssertRequest(client.LinkDeleteWithResponse(ctx, bookmark.Slug, member))(t, http.StatusForbidden)
	})

	linkedPage := tests.AssertRequest(client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString(), Url: &raw}, member))(t, http.StatusOK)
	linkedThread := tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString()}, member))(t, http.StatusOK)
	updatedThread := tests.AssertRequest(client.ThreadUpdateWithResponse(ctx, linkedThread.JSON200.Slug, openapi.ThreadUpdateJSONRequestBody{Url: &raw}, member))(t, http.StatusOK)
	require.NotNil(t, updatedThread.JSON200.Link)
	require.Equal(t, raw, updatedThread.JSON200.Link.Url)

	base := rbac.PermissionList{rbac.PermissionCreatePost, rbac.PermissionReadPublishedThreads, rbac.PermissionReadPublishedLibrary, rbac.PermissionManageLibrary}
	_, err = roles.Update(root, role.DefaultRoleMemberID, role_writer.WithPermissions(base))
	require.NoError(t, err)
	_, err = roles.Update(root, role.DefaultRoleGuestID, role_writer.WithPermissions(rbac.PermissionList{}))
	require.NoError(t, err)

	t.Run("removing_view_blocks_listing_but_preserves_direct_link_details", func(t *testing.T) {
		for _, editors := range [][]openapi.RequestEditorFn{nil, {member}} {
			tests.AssertRequest(client.LinkListWithResponse(ctx, nil, editors...))(t, http.StatusForbidden)
			tests.AssertRequest(client.LinkGetWithResponse(ctx, bookmark.Slug, editors...))(t, http.StatusOK)
		}
		page := tests.AssertRequest(client.NodeGetWithResponse(ctx, linkedPage.JSON200.Slug, nil, member))(t, http.StatusOK)
		require.NotNil(t, page.JSON200.Link)
		require.Equal(t, raw, page.JSON200.Link.Url)
		thread := tests.AssertRequest(client.ThreadGetWithResponse(ctx, linkedThread.JSON200.Slug, nil, member))(t, http.StatusOK)
		require.NotNil(t, thread.JSON200.Link)
		require.Equal(t, raw, thread.JSON200.Link.Url)

		tests.AssertRequest(client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, member))(t, http.StatusForbidden)
		tests.AssertRequest(client.LinkDeleteWithResponse(ctx, bookmark.Slug, member))(t, http.StatusForbidden)
	})

	t.Run("content_still_saves_and_indexes_mentioned_URLs_without_use", func(t *testing.T) {
		mentionedURL := raw + "mentioned"
		content := `<p>A normal hyperlink: <a href="` + mentionedURL + `">bookmark</a></p>`
		tests.AssertRequest(client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString(), Content: &content}, member))(t, http.StatusOK)
		tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString(), Body: &content}, member))(t, http.StatusOK)
		require.Eventually(t, func() bool {
			exists, err := db.Link.Query().Where(entlink.URLEQ(mentionedURL)).Exist(ctx)
			return err == nil && exists
		}, 5*time.Second, 10*time.Millisecond)
	})

	t.Run("primary_links_require_use_in_threads_nodes_and_batches", func(t *testing.T) {
		tests.AssertRequest(client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString(), Url: &raw}, member))(t, http.StatusForbidden)
		tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString(), Url: &raw}, member))(t, http.StatusForbidden)

		page := tests.AssertRequest(client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString()}, member))(t, http.StatusOK)
		thread := tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString()}, member))(t, http.StatusOK)
		tests.AssertRequest(client.NodeUpdateWithResponse(ctx, page.JSON200.Slug, openapi.NodeUpdateJSONRequestBody{Url: nullable.NewNullableWithValue(raw)}, member))(t, http.StatusForbidden)
		tests.AssertRequest(client.ThreadUpdateWithResponse(ctx, thread.JSON200.Slug, openapi.ThreadUpdateJSONRequestBody{Url: &raw}, member))(t, http.StatusForbidden)

		limited := session.WithAccountPermissions(ctx, *owner, rbac.NewList(rbac.PermissionManageLibrary))
		_, err := nodes.CreateMany(limited, []node_mutate.BatchCreate{{
			ID:     library.NodeID(xid.New()),
			Fields: node_mutate.BatchFields{Name: opt.New(uuid.NewString()), URL: opt.New(u)},
		}})
		require.Error(t, err)
		require.Equal(t, ftag.PermissionDenied, ftag.Get(err))

		_, err = nodes.UpdateMany(limited, []node_mutate.BatchUpdate{{
			Key:    library.NewKey(page.JSON200.Id),
			Fields: node_mutate.BatchFields{URL: opt.New(u)},
		}})
		require.Error(t, err)
		require.Equal(t, ftag.PermissionDenied, ftag.Get(err))
	})

	t.Run("manage_alone_grants_all_link_operations_and_deletion_preserves_content", func(t *testing.T) {
		_, err := roles.Update(root, role.DefaultRoleMemberID, role_writer.WithPermissions(rbac.PermissionList{rbac.PermissionManageLinks}))
		require.NoError(t, err)
		tests.AssertRequest(client.LinkListWithResponse(ctx, nil, member))(t, http.StatusOK)
		tests.AssertRequest(client.LinkGetWithResponse(ctx, bookmark.Slug, member))(t, http.StatusOK)
		tests.AssertRequest(client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, member))(t, http.StatusOK)

		node, err := db.Node.Create().SetName("Saved page").SetSlug(uuid.NewString()).SetAccountID(xid.ID(owner.ID)).SetLinkID(bookmark.ID).Save(ctx)
		require.NoError(t, err)
		post, err := db.Post.Create().SetTitle("Saved thread").SetSlug(uuid.NewString()).SetAccountPosts(xid.ID(owner.ID)).
			SetUpdatedAt(time.Now()).SetLastReplyAt(time.Now()).SetBody("Saved content").SetShort("Saved content").SetLinkID(bookmark.ID).AddContentLinkIDs(bookmark.ID).Save(ctx)
		require.NoError(t, err)
		image, err := db.Asset.Create().SetAccountID(xid.ID(owner.ID)).SetFilename("saved.png").SetSize(1).SetMimeType("image/png").Save(ctx)
		require.NoError(t, err)
		require.NoError(t, db.Link.UpdateOneID(bookmark.ID).SetPrimaryImageID(image.ID).Exec(ctx))
		tests.AssertRequest(client.LinkDeleteWithResponse(ctx, bookmark.Slug, member))(t, http.StatusNoContent)
		tests.AssertRequest(client.LinkGetWithResponse(ctx, bookmark.Slug, member))(t, http.StatusNotFound)
		tests.AssertRequest(client.LinkDeleteWithResponse(ctx, bookmark.Slug, member))(t, http.StatusNotFound)
		kept, err := db.Node.Get(ctx, node.ID)
		require.NoError(t, err)
		require.Equal(t, "Saved page", kept.Name)
		require.Equal(t, xid.NilID(), kept.LinkID)
		keptPost, err := db.Post.Get(ctx, post.ID)
		require.NoError(t, err)
		require.Equal(t, "Saved content", keptPost.Body)
		require.Equal(t, xid.NilID(), keptPost.LinkID)
		_, err = db.Asset.Get(ctx, image.ID)
		require.NoError(t, err)
	})
}
