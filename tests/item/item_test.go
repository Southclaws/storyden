package item_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
)

func TestItemsHappyPath(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, acc := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)

			name1 := "test-item-1-" + uuid.NewString()
			slug1 := name1
			cont1 := "# Item content\n\nRich text"
			iurl1 := "https://picsum.photos/200/200"
			url1 := "https://southcla.ws"
			item1, err := cl.ItemCreateWithResponse(ctx, openapi.ItemInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing items api",
				Content:     &cont1,
				ImageUrl:    &iurl1,
				Url:         &url1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(item1)
			r.Equal(200, item1.StatusCode())

			a.Equal(name1, item1.JSON200.Name)
			a.Equal(slug1, item1.JSON200.Slug)
			a.Equal("testing items api", item1.JSON200.Description)
			a.Equal(iurl1, *item1.JSON200.ImageUrl)
			a.Equal(url1, *item1.JSON200.Url)
			a.Equal(cont1, *item1.JSON200.Content)
			a.Equal(acc.ID.String(), string(item1.JSON200.Owner.Id))

			// Get the item just created

			item1get, err := cl.ItemGetWithResponse(ctx, item1.JSON200.Slug)
			r.NoError(err)
			r.NotNil(item1)
			r.Equal(200, item1.StatusCode())

			a.Equal(name1, item1get.JSON200.Name)
			a.Equal(slug1, item1get.JSON200.Slug)
			a.Equal("testing items api", item1get.JSON200.Description)
			a.Equal(cont1, *item1get.JSON200.Content)
			a.Equal(acc.ID.String(), string(item1get.JSON200.Owner.Id))

			// Update the item just created

			name1 = "test-item-1-UPDATED"
			slug1 = name1 + uuid.NewString()
			desc1 := "a new description"
			cont1 = "# New content"
			iurl1 = "https://picsum.photos/500/500"
			url1 = "https://cla.ws"
			prop1 := any(map[string]any{
				"key": "value",
			})
			item1update, err := cl.ItemUpdateWithResponse(ctx, item1.JSON200.Slug, openapi.ItemMutableProps{
				Name:        &name1,
				Slug:        &slug1,
				Description: &desc1,
				Content:     &cont1,
				ImageUrl:    &iurl1,
				Url:         &url1,
				Properties:  &prop1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(item1update)
			r.Equal(200, item1update.StatusCode())

			a.Equal(name1, item1update.JSON200.Name)
			a.Equal(slug1, item1update.JSON200.Slug)
			a.Equal(desc1, item1update.JSON200.Description)
			a.Equal(cont1, *item1update.JSON200.Content)
			a.Equal(iurl1, *item1update.JSON200.ImageUrl)
			a.Equal(url1, *item1update.JSON200.Url)
			a.Equal(prop1, item1update.JSON200.Properties)

			// Query for the exact item

			q := "test-item-1"
			items1, err := cl.ItemListWithResponse(ctx, &openapi.ItemListParams{
				Q: &q,
			})
			r.NoError(err)
			r.NotNil(items1)
			r.Equal(200, items1.StatusCode())

			ids := dt.Map(items1.JSON200.Items, func(c openapi.ItemWithParents) string { return c.Id })
			a.Contains(ids, item1.JSON200.Id)

			// Query for all items

			items2, err := cl.ItemListWithResponse(ctx, &openapi.ItemListParams{})
			r.NoError(err)
			r.NotNil(items2)
			r.Equal(200, items2.StatusCode())

			ids = dt.Map(items2.JSON200.Items, func(c openapi.ItemWithParents) string { return c.Id })
			a.Contains(ids, item1.JSON200.Id)
		}))
	}))
}

func TestItemsErrors(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)

			create401, err := cl.ItemCreateWithResponse(ctx, openapi.ItemInitialProps{})
			r.NoError(err)
			r.NotNil(create401)
			a.Equal(401, create401.StatusCode())

			get404, err := cl.ItemGetWithResponse(ctx, "nonexistent")
			r.NoError(err)
			r.NotNil(get404)
			a.Equal(404, get404.StatusCode())

			update401, err := cl.ItemUpdateWithResponse(ctx, "nonexistent", openapi.ItemMutableProps{})
			r.NoError(err)
			r.NotNil(update401)
			a.Equal(401, update401.StatusCode())

			update404, err := cl.ItemUpdateWithResponse(ctx, "nonexistent", openapi.ItemMutableProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(404, update404.StatusCode())
		}))
	}))
}
