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
	"github.com/Southclaws/storyden/app/resources/library/node_cache"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
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
		cacheStore cache.Store,
		nodeCache *node_cache.Cache,
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

			lastModified1Header := nodeGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present")

			lastModified1, err := time.Parse(time.RFC1123, lastModified1Header)
			r.NoError(err, "Last-Modified header should be parseable")

			cacheKey := "node:last-modified:" + slug
			cachedValue1, err := cacheStore.Get(ctx, cacheKey)
			r.NoError(err, "cache value should exist")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			time.Sleep(10 * time.Millisecond)

			nodeGet304, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
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

			cachedValue2, err := cacheStore.Get(ctx, cacheKey)
			r.NoError(err, "cache value should still exist after update")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			a.True(cachedTime2.After(cachedTime1),
				"cache should be updated IMMEDIATELY after node update (cachedTime1: %v, cachedTime2: %v)", cachedTime1, cachedTime2)

			nodeGet2, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{})
			tests.Ok(t, err, nodeGet2)
			a.Equal(newName, nodeGet2.JSON200.Name)

			lastModified2Header := nodeGet2.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified2Header, "Last-Modified header should be present")

			lastModified2, err := time.Parse(time.RFC1123, lastModified2Header)
			r.NoError(err, "Last-Modified header should be parseable")

			a.True(lastModified2.After(lastModified1) || lastModified2.Equal(lastModified1),
				"Last-Modified header should be updated or equal after update")

			lastModifiedFromCache := nodeCache.LastModified(ctx, slug)
			r.NotNil(lastModifiedFromCache, "node cache should return last modified time")
			a.True(lastModifiedFromCache.Equal(cachedTime2) || lastModifiedFromCache.After(cachedTime2),
				"cache last modified should match the cached value")

			nodeGetAfterUpdate, err := cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			if nodeGetAfterUpdate.HTTPResponse.StatusCode == 304 {
				t.Logf("Got 304 - cache timestamps: before=%v, after=%v, last-modified header=%s",
					cachedTime1, cachedTime2, lastModified1Header)
				t.Logf("New Last-Modified header: %s", nodeGetAfterUpdate.HTTPResponse.Header.Get("Last-Modified"))
			}
			tests.Ok(t, err, nodeGetAfterUpdate)
			r.NotNil(nodeGetAfterUpdate.JSON200, "should return 200 with body after cache invalidation")
			a.Equal(newName, nodeGetAfterUpdate.JSON200.Name)
		}))
	}))
}
