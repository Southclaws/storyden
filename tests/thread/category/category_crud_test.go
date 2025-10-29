package category_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestCategoryCRUD(t *testing.T) {
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

			t.Run("create_and_list", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				name := "Category " + uuid.NewString()
				create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#abc123", Description: "category testing", Name: name}, adminSession)
				tests.Ok(t, err, create)
				r.NotNil(create.JSON200)
				r.Equal(name, create.JSON200.Name)
				r.Equal(mark.Slugify(name), create.JSON200.Slug)

				get, err := cl.CategoryGetWithResponse(root, create.JSON200.Slug, adminSession)
				tests.Ok(t, err, get)
				r.NotNil(get.JSON200)
				a.Equal(create.JSON200.Id, get.JSON200.Id)
				a.Equal(create.JSON200.Slug, get.JSON200.Slug)

				list, err := cl.CategoryListWithResponse(root, adminSession)
				tests.Ok(t, err, list)
				r.NotNil(list.JSON200)
				found, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == create.JSON200.Id })
				r.True(ok)
				a.Equal(create.JSON200.Slug, found.Slug)
				a.Equal(create.JSON200.Name, found.Name)
			})

			t.Run("rename_category_updates_slug", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				name := "Category " + uuid.NewString()
				create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#def456", Description: "category testing", Name: name}, adminSession)
				tests.Ok(t, err, create)
				r.NotNil(create.JSON200)
				oldSlug := create.JSON200.Slug

				newName := "Category " + uuid.NewString()
				newSlug := mark.Slugify(newName)
				update, err := cl.CategoryUpdateWithResponse(root, oldSlug, openapi.CategoryUpdateJSONRequestBody{Name: lo.ToPtr(newName), Slug: lo.ToPtr(newSlug)}, adminSession)
				tests.Ok(t, err, update)
				r.NotNil(update.JSON200)
				a.Equal(newSlug, update.JSON200.Slug)
				a.Equal(newName, update.JSON200.Name)

				get, err := cl.CategoryGetWithResponse(root, newSlug, adminSession)
				tests.Ok(t, err, get)
				r.NotNil(get.JSON200)
				a.Equal(newSlug, get.JSON200.Slug)
				a.Equal(newName, get.JSON200.Name)

				oldGet, err := cl.CategoryGetWithResponse(root, oldSlug, adminSession)
				tests.Status(t, err, oldGet, http.StatusNotFound)

				list, err := cl.CategoryListWithResponse(root, adminSession)
				tests.Ok(t, err, list)
				r.NotNil(list.JSON200)
				found, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == update.JSON200.Id })
				r.True(ok)
				a.Equal(newSlug, found.Slug)
				a.Equal(newName, found.Name)
			})

			t.Run("delete_category", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				originName := "Category " + uuid.NewString()
				origin, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#123123", Description: "category testing", Name: originName}, adminSession)
				tests.Ok(t, err, origin)
				r.NotNil(origin.JSON200)

				targetName := "Category " + uuid.NewString()
				target, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#321321", Description: "category testing", Name: targetName}, adminSession)
				tests.Ok(t, err, target)
				r.NotNil(target.JSON200)

				deleted, err := cl.CategoryDeleteWithResponse(root, origin.JSON200.Slug, openapi.CategoryDeleteJSONRequestBody{MoveTo: target.JSON200.Id}, adminSession)
				tests.Ok(t, err, deleted)
				r.NotNil(deleted.JSON200)
				a.Equal(origin.JSON200.Id, deleted.JSON200.Id)

				getDeleted, err := cl.CategoryGetWithResponse(root, origin.JSON200.Slug, adminSession)
				tests.Status(t, err, getDeleted, http.StatusNotFound)

				list, err := cl.CategoryListWithResponse(root, adminSession)
				tests.Ok(t, err, list)
				r.NotNil(list.JSON200)
				_, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == origin.JSON200.Id })
				a.False(ok)
			})
		}))
	}))
}
