package library_test

import (
	"context"
	"net/http"
	"sync"
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

			visibility := openapi.VisibilityPublished
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
			r.Empty(nodeGet1.HTTPResponse.Header.Get("Last-Modified"), "ETag is the sole conditional validator")

			nodeGet304 := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(nodeGet304.JSON200, "304 response should have no body")

			legacyValidator := tests.AssertRequest(cl.NodeGetWithResponse(ctx, slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", "Sun, 06 Nov 2094 08:49:37 GMT")
				return nil
			}))(t, http.StatusOK)
			r.NotNil(legacyValidator.JSON200, "If-Modified-Since is intentionally ignored")

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

func TestNodeCacheCanonicalisesIDQueryForms(t *testing.T) {
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
			visibility := openapi.VisibilityPublished
			name := "cache-alias-node-" + uuid.NewString()

			created := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name, Slug: &name, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200
			mark := created.Id + "-" + created.Slug

			byID := tests.AssertRequest(cl.NodeGetWithResponse(ctx, created.Id, &openapi.NodeGetParams{}))(t, http.StatusOK)
			etag := byID.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag)

			byMark := tests.AssertRequest(cl.NodeGetWithResponse(ctx, mark, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(byMark.JSON200)

			newName := "updated-cache-alias-node-" + uuid.NewString()
			tests.AssertRequest(cl.NodeUpdateWithResponse(ctx, created.Slug, openapi.NodeMutableProps{
				Name: &newName,
			}, session))(t, http.StatusOK)

			for _, key := range []string{created.Id, mark} {
				got := tests.AssertRequest(cl.NodeGetWithResponse(ctx, key, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
					req.Header.Set("If-None-Match", etag)
					return nil
				}))(t, http.StatusOK)
				r.NotNil(got.JSON200)
				a.Equal(newName, got.JSON200.Name)
			}
		}))
	}))
}

func TestNodeRenameAndDeleteRejectCurrentValidatorForRetiredSlugs(t *testing.T) {
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

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)
			visibility := openapi.VisibilityPublished
			oldSlug := "cache-retired-node-" + uuid.NewString()
			newSlug := "cache-current-node-" + uuid.NewString()

			created := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: oldSlug, Slug: &oldSlug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200

			// Populate the old slug validator before retiring it.
			tests.AssertRequest(cl.NodeGetWithResponse(ctx, oldSlug, &openapi.NodeGetParams{}))(t, http.StatusOK)

			tests.AssertRequest(cl.NodeUpdateWithResponse(ctx, created.Slug, openapi.NodeMutableProps{
				Slug: &newSlug,
			}, session))(t, http.StatusOK)

			current := tests.AssertRequest(cl.NodeGetWithResponse(ctx, newSlug, &openapi.NodeGetParams{}))(t, http.StatusOK)
			currentETag := current.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(currentETag)

			retired, err := cl.NodeGetWithResponse(ctx, oldSlug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", currentETag)
				return nil
			})
			r.NoError(err)
			r.Equal(http.StatusNotFound, retired.StatusCode(), "a live validator must not make a retired slug return 304")

			tests.AssertRequest(cl.NodeDeleteWithResponse(ctx, newSlug, &openapi.NodeDeleteParams{}, session))(t, http.StatusOK)

			deleted, err := cl.NodeGetWithResponse(ctx, newSlug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", currentETag)
				return nil
			})
			r.NoError(err)
			r.Equal(http.StatusNotFound, deleted.StatusCode(), "a deleted slug must not return 304")
		}))
	}))
}

func TestNodeDeleteInvalidatesEveryConditionalRequestAlias(t *testing.T) {
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

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)
			visibility := openapi.VisibilityPublished
			slug := "cache-delete-node-" + uuid.NewString()

			created := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: slug, Slug: &slug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200
			mark := created.Id + "-" + created.Slug

			aliases := map[string]string{}
			for _, key := range []string{created.Id, created.Slug, mark} {
				initial := tests.AssertRequest(cl.NodeGetWithResponse(ctx, key, &openapi.NodeGetParams{}))(t, http.StatusOK)
				etag := initial.HTTPResponse.Header.Get("ETag")
				r.NotEmpty(etag)
				aliases[key] = etag

				tests.AssertRequest(cl.NodeGetWithResponse(ctx, key, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
					req.Header.Set("If-None-Match", etag)
					return nil
				}))(t, http.StatusNotModified)
			}

			tests.AssertRequest(cl.NodeDeleteWithResponse(ctx, created.Slug, &openapi.NodeDeleteParams{}, session))(t, http.StatusOK)

			for key, etag := range aliases {
				response, err := cl.NodeGetWithResponse(ctx, key, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
					req.Header.Set("If-None-Match", etag)
					return nil
				})
				r.NoError(err)
				r.Equal(http.StatusNotFound, response.StatusCode(), "deleted alias %q must never return 304", key)
			}
		}))
	}))
}

func TestNodeDeleteInvalidatesParentRepresentation(t *testing.T) {
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
			visibility := openapi.VisibilityPublished
			parentSlug := "cache-delete-parent-" + uuid.NewString()
			childSlug := "cache-delete-child-" + uuid.NewString()

			parent := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentSlug, Slug: &parentSlug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200
			child := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: childSlug, Slug: &childSlug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200

			tests.AssertRequest(cl.NodeAddNodeWithResponse(ctx, parent.Slug, child.Slug, session))(t, http.StatusOK)

			before := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parent.Slug, &openapi.NodeGetParams{}))(t, http.StatusOK)
			parentETag := before.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(parentETag)
			r.Len(before.JSON200.Children, 1)

			tests.AssertRequest(cl.NodeDeleteWithResponse(ctx, child.Slug, &openapi.NodeDeleteParams{}, session))(t, http.StatusOK)

			after := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parent.Slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", parentETag)
				return nil
			}))(t, http.StatusOK)
			r.NotNil(after.JSON200)
			a.Empty(after.JSON200.Children)
		}))
	}))
}

func TestNodeCacheRejectsOldValidatorDuringConcurrentUpdates(t *testing.T) {
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

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)
			visibility := openapi.VisibilityPublished
			slug := "cache-concurrent-node-" + uuid.NewString()
			created := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: slug, Slug: &slug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200

			initial := tests.AssertRequest(cl.NodeGetWithResponse(ctx, created.Slug, &openapi.NodeGetParams{}))(t, http.StatusOK)
			oldETag := initial.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(oldETag)

			firstName := "cache-concurrent-first-" + uuid.NewString()
			tests.AssertRequest(cl.NodeUpdateWithResponse(ctx, created.Slug, openapi.NodeMutableProps{Name: &firstName}, session))(t, http.StatusOK)

			const writers = 8
			const readers = 64
			statuses := make(chan int, writers+readers)
			errors := make(chan error, writers+readers)
			var wg sync.WaitGroup
			wg.Add(writers + readers)

			for range writers {
				go func() {
					defer wg.Done()
					name := "cache-concurrent-update-" + uuid.NewString()
					response, err := cl.NodeUpdateWithResponse(ctx, created.Slug, openapi.NodeMutableProps{Name: &name}, session)
					if err != nil {
						errors <- err
						return
					}
					statuses <- response.StatusCode()
				}()
			}

			for range readers {
				go func() {
					defer wg.Done()
					response, err := cl.NodeGetWithResponse(ctx, created.Slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
						req.Header.Set("If-None-Match", oldETag)
						return nil
					})
					if err != nil {
						errors <- err
						return
					}
					statuses <- response.StatusCode()
				}()
			}

			wg.Wait()
			close(statuses)
			close(errors)

			for err := range errors {
				r.NoError(err)
			}
			for status := range statuses {
				r.Equal(http.StatusOK, status, "an old validator must never return 304 after the first completed update")
			}
		}))
	}))
}

func TestNodeTreeMoveInvalidatesParentRepresentation(t *testing.T) {
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
			visibility := openapi.VisibilityPublished
			parentSlug := "cache-parent-" + uuid.NewString()
			childSlug := "cache-child-" + uuid.NewString()

			parent := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentSlug, Slug: &parentSlug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200
			child := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: childSlug, Slug: &childSlug, Visibility: &visibility,
			}, session))(t, http.StatusOK).JSON200

			parentGet := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parent.Slug, &openapi.NodeGetParams{}))(t, http.StatusOK)
			etag := parentGet.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag)

			position := openapi.NodePositionMutableProps{}
			position.Parent.Set(parent.Slug)
			tests.AssertRequest(cl.NodeUpdatePositionWithResponse(ctx, child.Slug, position, session))(t, http.StatusOK)

			updatedParent := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parent.Slug, &openapi.NodeGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag)
				return nil
			}))(t, http.StatusOK)
			r.NotNil(updatedParent.JSON200)
			found := false
			for _, node := range updatedParent.JSON200.Children {
				if node.Id == child.Id {
					found = true
					break
				}
			}
			a.True(found, "moved child should appear in the invalidated parent representation")
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

			visibility := openapi.VisibilityPublished
			name := "cache-test-schema-" + uuid.NewString()
			slug := name
			ptype := openapi.PropertyTypeText

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
					Type: openapi.PropertyTypeText,
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
