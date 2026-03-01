package library_test

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

func TestNodeCacheWithUpdate(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			visibility := openapi.Published
			name := "cache-test-node-" + uuid.NewString()
			slug := name

			nodeCreate := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name,
				Slug:       &slug,
				Visibility: &visibility,
			}, session))(t, http.StatusOK)
			a.Equal(name, nodeCreate.JSON200.Name)

			nodeGet1 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}))(t, http.StatusOK)

			etag1 := nodeGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")
			lastModified1Header := nodeGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present for backward compatibility")

			nodeGet304 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(nodeGet304.JSON200, "304 response should have no body")

			newName := "updated-cache-test-node-" + uuid.NewString()
			nodeUpdate := tests.AssertRequest(cl.NodeUpdateWithResponse(ctx, slug, openapi.NodeMutableProps{
				Name: &newName,
			}, session))(t, http.StatusOK)
			a.Equal(newName, nodeUpdate.JSON200.Name)

			nodeGet2 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}))(t, http.StatusOK)
			a.Equal(newName, nodeGet2.JSON200.Name)

			etag2 := nodeGet2.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag2, "ETag header should be present")
			a.NotEqual(etag1, etag2, "ETag should change after update")

			nodeGetAfterUpdate := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusOK)
			r.NotNil(nodeGetAfterUpdate.JSON200, "should return 200 with body after cache invalidation")
			a.Equal(newName, nodeGetAfterUpdate.JSON200.Name)
		}))
	}))
}

func TestNodeCacheWithPropertySchemaUpdate(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			visibility := openapi.Published
			name := "cache-test-schema-" + uuid.NewString()
			slug := name
			ptype := openapi.Text

			nodeCreate := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name,
				Slug:       &slug,
				Visibility: &visibility,
				Properties: &openapi.PropertyMutationList{
					{
						Name:  "test_field",
						Type:  &ptype,
						Value: "initial value",
					},
				},
			}, session))(t, http.StatusOK)
			r.NotNil(nodeCreate.JSON200)

			nodeGet1 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}))(t, http.StatusOK)

			etag1 := nodeGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")

			nodeGet304 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(nodeGet304.JSON200, "304 response should have no body")

			schemaUpdate := tests.AssertRequest(cl.NodeUpdatePropertySchemaWithResponse(ctx, slug, openapi.NodeUpdatePropertySchemaJSONRequestBody{
				{
					Name: "new_field",
					Type: openapi.Text,
					Sort: "b",
				},
			}, session))(t, http.StatusOK)
			r.NotNil(schemaUpdate.JSON200)

			nodeGetAfterSchema := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusOK)
			r.NotNil(nodeGetAfterSchema.JSON200, "should return 200 with body after schema update invalidates cache")

			etag2 := nodeGetAfterSchema.HTTPResponse.Header.Get("ETag")
			a.NotEqual(etag1, etag2, "ETag should change after property schema update")
		}))
	}))
}
