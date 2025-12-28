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

			nodeCreate, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name,
				Slug:       &slug,
				Visibility: &visibility,
			}, session)
			tests.Ok(t, err, nodeCreate)
			a.Equal(name, nodeCreate.JSON200.Name)

			nodeGet1, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
			tests.Ok(t, err, nodeGet1)

			etag1 := nodeGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")
			lastModified1Header := nodeGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present for backward compatibility")

			nodeGet304, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			})
			tests.Status(t, err, nodeGet304, 304)
			a.Nil(nodeGet304.JSON200, "304 response should have no body")

			newName := "updated-cache-test-node-" + uuid.NewString()
			nodeUpdate, err := cl.NodeUpdateWithResponse(ctx, slug, openapi.NodeMutableProps{
				Name: &newName,
			}, session)
			tests.Ok(t, err, nodeUpdate)
			a.Equal(newName, nodeUpdate.JSON200.Name)

			nodeGet2, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
			tests.Ok(t, err, nodeGet2)
			a.Equal(newName, nodeGet2.JSON200.Name)

			etag2 := nodeGet2.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag2, "ETag header should be present")
			a.NotEqual(etag1, etag2, "ETag should change after update")

			nodeGetAfterUpdate, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			})
			tests.Ok(t, err, nodeGetAfterUpdate)
			r.NotNil(nodeGetAfterUpdate.JSON200, "should return 200 with body after cache invalidation")
			a.Equal(newName, nodeGetAfterUpdate.JSON200.Name)
		}))
	}))
}
