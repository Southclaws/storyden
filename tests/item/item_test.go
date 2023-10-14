package item_test

import (
	"context"
	"testing"

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

			name1 := "test-item-1"
			slug1 := name1 + uuid.NewString()
			item1, err := cl.ItemCreateWithResponse(ctx, openapi.ItemInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing items api",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(item1)
			r.Equal(200, item1.StatusCode())

			a.Equal(name1, item1.JSON200.Name)
			a.Equal(slug1, item1.JSON200.Slug)
			a.Equal("testing items api", item1.JSON200.Description)
			a.Equal(acc.ID.String(), string(item1.JSON200.Owner.Id))
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

			create400, err := cl.ItemCreateWithResponse(ctx, openapi.ItemInitialProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(create400)
			a.Equal(400, create400.StatusCode())

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
