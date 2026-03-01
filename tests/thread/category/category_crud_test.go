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

func TestCategoryCRUD_CreateAndList(t *testing.T) {
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

			name := "Category " + uuid.NewString()
			create := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#abc123",
				Description: "category testing",
				Name:        name,
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, create.JSON200)

			t.Run("create_sets_name_and_slug", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				r.NotNil(create.JSON200)
				a.Equal(name, create.JSON200.Name)
				a.Equal(mark.Slugify(name), create.JSON200.Slug)
			})

			t.Run("get_returns_created_category", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				get := tests.AssertRequest(cl.CategoryGetWithResponse(root, create.JSON200.Slug, adminSession))(t, http.StatusOK)
				r.NotNil(get.JSON200)
				a.Equal(create.JSON200.Id, get.JSON200.Id)
				a.Equal(create.JSON200.Slug, get.JSON200.Slug)
			})

			t.Run("list_includes_created_category", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				list := tests.AssertRequest(cl.CategoryListWithResponse(root, adminSession))(t, http.StatusOK)
				r.NotNil(list.JSON200)
				found, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == create.JSON200.Id })
				r.True(ok)
				a.Equal(create.JSON200.Slug, found.Slug)
				a.Equal(create.JSON200.Name, found.Name)
			})
		}))
	}))
}

func TestCategoryCRUD_RenameCategoryUpdatesSlug(t *testing.T) {
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

			originalName := "Category " + uuid.NewString()
			newName := "Category " + uuid.NewString()
			newSlug := mark.Slugify(newName)
			create := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#def456",
				Description: "category testing",
				Name:        originalName,
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, create.JSON200)
			oldSlug := create.JSON200.Slug

			update := tests.AssertRequest(cl.CategoryUpdateWithResponse(root, oldSlug, openapi.CategoryUpdateJSONRequestBody{
				Name: lo.ToPtr(newName),
				Slug: lo.ToPtr(newSlug),
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, update.JSON200)

			t.Run("update_returns_new_name_and_slug", func(t *testing.T) {
				a := assert.New(t)
				a.Equal(newSlug, update.JSON200.Slug)
				a.Equal(newName, update.JSON200.Name)
			})

			t.Run("new_slug_resolves_and_old_slug_does_not", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				get := tests.AssertRequest(cl.CategoryGetWithResponse(root, newSlug, adminSession))(t, http.StatusOK)
				r.NotNil(get.JSON200)
				a.Equal(newSlug, get.JSON200.Slug)
				a.Equal(newName, get.JSON200.Name)

				oldGet := tests.AssertRequest(cl.CategoryGetWithResponse(root, oldSlug, adminSession))(t, http.StatusNotFound)
				r.NotNil(oldGet)
			})

			t.Run("list_reflects_updated_category", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				list := tests.AssertRequest(cl.CategoryListWithResponse(root, adminSession))(t, http.StatusOK)
				r.NotNil(list.JSON200)
				found, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == update.JSON200.Id })
				r.True(ok)
				a.Equal(newSlug, found.Slug)
				a.Equal(newName, found.Name)
			})
		}))
	}))
}

func TestCategoryCRUD_DeleteCategory(t *testing.T) {
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

			originName := "Category " + uuid.NewString()
			targetName := "Category " + uuid.NewString()
			origin := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#123123",
				Description: "category testing",
				Name:        originName,
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, origin.JSON200)

			target := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#321321",
				Description: "category testing",
				Name:        targetName,
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, target.JSON200)

			deleted := tests.AssertRequest(cl.CategoryDeleteWithResponse(root, origin.JSON200.Slug, openapi.CategoryDeleteJSONRequestBody{
				MoveTo: target.JSON200.Id,
			}, adminSession))(t, http.StatusOK)
			require.NotNil(t, deleted.JSON200)

			t.Run("delete_returns_deleted_category", func(t *testing.T) {
				a := assert.New(t)
				a.Equal(origin.JSON200.Id, deleted.JSON200.Id)
			})

			t.Run("deleted_slug_not_found", func(t *testing.T) {
				getDeleted := tests.AssertRequest(cl.CategoryGetWithResponse(root, origin.JSON200.Slug, adminSession))(t, http.StatusNotFound)
				require.NotNil(t, getDeleted)
			})

			t.Run("list_excludes_deleted_category", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				list := tests.AssertRequest(cl.CategoryListWithResponse(root, adminSession))(t, http.StatusOK)
				r.NotNil(list.JSON200)
				_, ok := lo.Find(list.JSON200.Categories, func(c openapi.Category) bool { return c.Id == origin.JSON200.Id })
				a.False(ok)
			})
		}))
	}))
}
