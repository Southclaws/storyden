package library_test

import (
	"context"
	"net/http"
	"testing"
	"time"

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

			stableETag, nodeGet304 := waitForNodeNotModified(t, ctx, cl, slug, etag1)
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

			nodeGetAfterUpdate := waitForNodeModified(t, ctx, cl, slug, stableETag)
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

			nodeCreate, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
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
			}, session)
			tests.Ok(t, err, nodeCreate)

			nodeGet1, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
			tests.Ok(t, err, nodeGet1)

			etag1 := nodeGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")

			stableETag, _ := waitForNodeNotModified(t, ctx, cl, slug, etag1)

			schemaUpdate, err := cl.NodeUpdatePropertySchemaWithResponse(ctx, slug, openapi.NodeUpdatePropertySchemaJSONRequestBody{
				{
					Name: "new_field",
					Type: openapi.Text,
					Sort: "b",
				},
			}, session)
			tests.Ok(t, err, schemaUpdate)

			nodeGetAfterSchema := waitForNodeModified(t, ctx, cl, slug, stableETag)
			r.NotNil(nodeGetAfterSchema.JSON200, "should return 200 with body after schema update invalidates cache")

			etag2 := nodeGetAfterSchema.HTTPResponse.Header.Get("ETag")
			a.NotEqual(stableETag, etag2, "ETag should change after property schema update")
		}))
	}))
}

func waitForNodeNotModified(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	slug string,
	initialETag string,
) (string, *openapi.NodeGetResponse) {
	t.Helper()

	const maxAttempts = 12
	deadline := time.Now().Add(2 * time.Second)
	etag := initialETag

	for attempt := 0; attempt < maxAttempts && time.Now().Before(deadline); attempt++ {
		resp, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("If-None-Match", etag)
			return nil
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		switch resp.StatusCode() {
		case http.StatusNotModified:
			return etag, resp
		case http.StatusOK:
			nextETag := resp.HTTPResponse.Header.Get("ETag")
			if nextETag != "" {
				etag = nextETag
			}
		default:
			tests.Status(t, err, resp, http.StatusNotModified)
		}

		time.Sleep(40 * time.Millisecond)
	}

	t.Fatalf("did not observe 304 for node %q within retry window", slug)
	return "", nil
}

func waitForNodeModified(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	slug string,
	ifNoneMatch string,
) *openapi.NodeGetResponse {
	t.Helper()

	const maxAttempts = 12
	deadline := time.Now().Add(2 * time.Second)

	for attempt := 0; attempt < maxAttempts && time.Now().Before(deadline); attempt++ {
		resp, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("If-None-Match", ifNoneMatch)
			return nil
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		switch resp.StatusCode() {
		case http.StatusOK:
			return resp
		case http.StatusNotModified:
			time.Sleep(40 * time.Millisecond)
			continue
		default:
			tests.Status(t, err, resp, http.StatusOK)
		}
	}

	t.Fatalf("did not observe 200 cache miss for node %q within retry window", slug)
	return nil
}
