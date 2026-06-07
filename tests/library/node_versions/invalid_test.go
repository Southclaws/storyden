package node_versions_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeVersionInvalidRequests(t *testing.T) {
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

			node := createPublishedNode(t, root, cl, adminSession, "invalid-target")
			version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Invalid request draft "+uuid.NewString())

			t.Run("get rejects invalid version id", func(t *testing.T) {
				t.Parallel()

				get, err := cl.NodeVersionGetWithResponse(root, node.Slug, "not-a-version-id", authorSession)
				tests.Status(t, err, get, http.StatusBadRequest)
			})

			t.Run("update status rejects invalid version id", func(t *testing.T) {
				t.Parallel()

				update, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, "not-a-version-id", openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Status(t, err, update, http.StatusBadRequest)
			})

			t.Run("status transition only supports applied", func(t *testing.T) {
				t.Parallel()

				update, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusDraft,
				}, adminSession)
				tests.Status(t, err, update, http.StatusBadRequest)
			})

			t.Run("create rejects missing target node", func(t *testing.T) {
				t.Parallel()
				updatedName := "Missing target proposal " + uuid.NewString()

				create, err := cl.NodeVersionCreateWithResponse(root, "missing-node-"+uuid.NewString(), openapi.NodeVersionCreateJSONRequestBody{
					Name: &updatedName,
				}, authorSession)
				tests.Status(t, err, create, http.StatusNotFound)
			})

			t.Run("create rejects invalid property type", func(t *testing.T) {
				t.Parallel()
				invalidType := openapi.PropertyType("not-a-property-type")
				properties := openapi.PropertyMutationList{
					{Name: "Release year", Value: "1992", Type: &invalidType},
				}

				create, err := cl.NodeVersionCreateWithResponse(root, node.Slug, openapi.NodeVersionCreateJSONRequestBody{
					Properties: &properties,
				}, authorSession)
				tests.Status(t, err, create, http.StatusBadRequest)
			})
		}))
	}))
}
