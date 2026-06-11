package node_versions_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
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

func TestNodeVersionAuthAndVisibility(t *testing.T) {
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

			authorCtx, _ := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			authorSession := sh.WithSession(authorCtx)

			otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			otherSession := sh.WithSession(otherCtx)

			node := createPublishedNode(t, root, cl, adminSession, "auth-target")
			version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Private draft "+uuid.NewString())

			t.Run("unauthenticated member cannot create draft", func(t *testing.T) {
				t.Parallel()
				updatedName := "Unauthenticated proposal " + uuid.NewString()

				create, err := cl.NodeVersionCreateWithResponse(root, node.Slug, openapi.NodeVersionCreateJSONRequestBody{
					Name: &updatedName,
				})
				tests.Status(t, err, create, http.StatusForbidden)
			})

			t.Run("public cannot list draft", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				list, err := cl.NodeVersionListWithResponse(root, node.Slug, nil)
				tests.Ok(t, err, list)
				a.Empty(list.JSON200.Versions)
			})

			t.Run("public cannot get draft", func(t *testing.T) {
				t.Parallel()

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, version.Id)
				tests.Status(t, err, get, http.StatusNotFound)
			})

			t.Run("other member cannot list draft", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				list, err := cl.NodeVersionListWithResponse(root, node.Slug, nil, otherSession)
				tests.Ok(t, err, list)
				a.Empty(list.JSON200.Versions)
			})

			t.Run("other member cannot get draft", func(t *testing.T) {
				t.Parallel()

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, version.Id, otherSession)
				tests.Status(t, err, get, http.StatusNotFound)
			})

			t.Run("other member cannot update draft", func(t *testing.T) {
				t.Parallel()
				updatedName := "Other update " + uuid.NewString()

				update, err := cl.NodeVersionUpdateWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateJSONRequestBody{
					Name: &updatedName,
				}, otherSession)
				tests.Status(t, err, update, http.StatusNotFound)
			})

			t.Run("other member cannot discard draft", func(t *testing.T) {
				t.Parallel()

				deleteAsOther, err := cl.NodeVersionDeleteWithResponse(root, node.Slug, version.Id, otherSession)
				tests.Status(t, err, deleteAsOther, http.StatusNotFound)
			})

			t.Run("author can discard their own draft", func(t *testing.T) {
				t.Parallel()

				node := createPublishedNode(t, root, cl, adminSession, "author-discard-"+uuid.NewString())
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Draft by author "+uuid.NewString())

				deleteOwn, err := cl.NodeVersionDeleteWithResponse(root, node.Slug, version.Id, authorSession)
				tests.Ok(t, err, deleteOwn)

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, version.Id, authorSession)
				tests.Status(t, err, get, http.StatusNotFound)
			})

			t.Run("manager can discard any draft", func(t *testing.T) {
				t.Parallel()

				node := createPublishedNode(t, root, cl, adminSession, "manager-discard-"+uuid.NewString())
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Draft by author "+uuid.NewString())

				deleteAny, err := cl.NodeVersionDeleteWithResponse(root, node.Slug, version.Id, adminSession)
				tests.Ok(t, err, deleteAny)

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, version.Id, adminSession)
				tests.Status(t, err, get, http.StatusNotFound)
			})

			t.Run("manager can list draft", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)
				r := require.New(t)

				list, err := cl.NodeVersionListWithResponse(root, node.Slug, nil, adminSession)
				tests.Ok(t, err, list)
				r.Len(list.JSON200.Versions, 1)
				a.Equal(version.Id, list.JSON200.Versions[0].Id)
			})

			t.Run("manager can get draft", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)
				r := require.New(t)

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, version.Id, adminSession)
				tests.Ok(t, err, get)
				r.NotNil(get.JSON200)
				a.Equal(version.Id, get.JSON200.Id)
				a.Equal(openapi.NodeVersionStatusDraft, get.JSON200.Status)
			})

			t.Run("author cannot apply draft", func(t *testing.T) {
				t.Parallel()

				apply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, authorSession)
				tests.Status(t, err, apply, http.StatusForbidden)
			})
		}))
	}))
}
