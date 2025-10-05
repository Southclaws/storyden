package category_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/samber/lo"
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

func TestCategoryMovement(t *testing.T) {
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

			t.Run("move_category_to_new_parent", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				parentName := "Category " + uuid.NewString()
				parent, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#aaaaaa", Description: "category testing", Name: parentName}, adminSession)
				tests.Ok(t, err, parent)
				r.NotNil(parent.JSON200)

				targetParentName := "Category " + uuid.NewString()
				targetParent, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#bbbbbb", Description: "category testing", Name: targetParentName}, adminSession)
				tests.Ok(t, err, targetParent)
				r.NotNil(targetParent.JSON200)

				childName := "Category " + uuid.NewString()
				child, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#cccccc", Description: "category testing", Name: childName, Parent: lo.ToPtr(parent.JSON200.Id)}, adminSession)
				tests.Ok(t, err, child)
				r.NotNil(child.JSON200)

				moveBody := openapi.CategoryUpdatePositionJSONRequestBody{Parent: nullable.NewNullableWithValue(openapi.NullableIdentifier(targetParent.JSON200.Id))}
				move, err := cl.CategoryUpdatePositionWithResponse(root, child.JSON200.Slug, moveBody, adminSession)
				tests.Ok(t, err, move)
				r.NotNil(move.JSON200)

				updated, err := cl.CategoryGetWithResponse(root, child.JSON200.Slug, adminSession)
				tests.Ok(t, err, updated)
				r.NotNil(updated.JSON200)
				r.NotNil(updated.JSON200.Parent)
				a.Equal(openapi.Identifier(targetParent.JSON200.Id), *updated.JSON200.Parent)
			})

			t.Run("move_category_to_root", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				parentName := "Category " + uuid.NewString()
				parent, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#aaaa11", Description: "category testing", Name: parentName}, adminSession)
				tests.Ok(t, err, parent)
				r.NotNil(parent.JSON200)

				childName := "Category " + uuid.NewString()
				child, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#bbbb22", Description: "category testing", Name: childName, Parent: lo.ToPtr(parent.JSON200.Id)}, adminSession)
				tests.Ok(t, err, child)
				r.NotNil(child.JSON200)

				moveBody := openapi.CategoryUpdatePositionJSONRequestBody{Parent: nullable.NewNullNullable[openapi.NullableIdentifier]()}
				move, err := cl.CategoryUpdatePositionWithResponse(root, child.JSON200.Slug, moveBody, adminSession)
				tests.Ok(t, err, move)
				r.NotNil(move.JSON200)

				updated, err := cl.CategoryGetWithResponse(root, child.JSON200.Slug, adminSession)
				tests.Ok(t, err, updated)
				r.NotNil(updated.JSON200)
				a.Nil(updated.JSON200.Parent)
			})

			t.Run("resort_categories", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				firstName := "Category " + uuid.NewString()
				first, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#101010", Description: "category testing", Name: firstName}, adminSession)
				tests.Ok(t, err, first)
				r.NotNil(first.JSON200)

				secondName := "Category " + uuid.NewString()
				second, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#202020", Description: "category testing", Name: secondName}, adminSession)
				tests.Ok(t, err, second)
				r.NotNil(second.JSON200)

				thirdName := "Category " + uuid.NewString()
				third, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Colour: "#303030", Description: "category testing", Name: thirdName}, adminSession)
				tests.Ok(t, err, third)
				r.NotNil(third.JSON200)

				moveBody := openapi.CategoryUpdatePositionJSONRequestBody{Before: lo.ToPtr(openapi.Identifier(first.JSON200.Id))}
				move, err := cl.CategoryUpdatePositionWithResponse(root, third.JSON200.Slug, moveBody, adminSession)
				tests.Ok(t, err, move)
				r.NotNil(move.JSON200)

				ordered := lo.Filter(move.JSON200.Categories, func(c openapi.Category, _ int) bool {
					return c.Id == third.JSON200.Id || c.Id == first.JSON200.Id || c.Id == second.JSON200.Id
				})
				r.Len(ordered, 3)
				ids := lo.Map(ordered, func(c openapi.Category, _ int) openapi.Identifier { return c.Id })
				a.Equal([]openapi.Identifier{third.JSON200.Id, first.JSON200.Id, second.JSON200.Id}, ids)
			})
		}))
	}))
}
