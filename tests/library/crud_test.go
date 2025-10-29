package library_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesHappyPath(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			ctx, acc := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			visibility := openapi.Published

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			content1 := "<h1>Nodes</h1><p>Rich text content.</p>"
			url1 := "https://southcla.ws"
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name1,
				Slug:       &slug1,
				Content:    &content1,
				Url:        &url1,
				Visibility: &visibility, // Admin account can post directly to published
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node1)

			a.Equal(name1, node1.JSON200.Name)
			a.Equal(slug1, node1.JSON200.Slug)
			a.Equal("Rich text content.", node1.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Rich text content.</p></body>", *node1.JSON200.Content)
			r.NotNil(node1.JSON200.Link)
			a.Equal(url1, node1.JSON200.Link.Url)
			a.Equal(acc.ID.String(), string(node1.JSON200.Owner.Id))

			// Get the one just created

			node1get, err := cl.NodeGetWithResponse(ctx, slug1, &openapi.NodeGetParams{})
			tests.Ok(t, err, node1get)

			a.Equal(name1, node1get.JSON200.Name)
			a.Equal(slug1, node1get.JSON200.Slug)
			a.Equal("Rich text content.", node1get.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Rich text content.</p></body>", *node1get.JSON200.Content)
			a.Equal(acc.ID.String(), string(node1get.JSON200.Owner.Id))

			// Update the one just created

			name1 = "test-node-1-updated"
			slug1 = name1 + uuid.NewString()
			cont1 := "<h1>Nodes</h1><p>Newly changed content.</p>"
			url1 = "https://cla.ws"
			prop1 := openapi.Metadata(map[string]any{
				"key": "value",
			})
			node1update, err := cl.NodeUpdateWithResponse(ctx, node1.JSON200.Slug, openapi.NodeMutableProps{
				Name:    &name1,
				Slug:    &slug1,
				Content: &cont1,
				Url:     nullable.NewNullableWithValue(url1),
				Meta:    &prop1,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node1update)

			a.Equal(name1, node1update.JSON200.Name)
			a.Equal(slug1, node1update.JSON200.Slug)
			a.Equal("Newly changed content.", node1update.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Newly changed content.</p></body>", *node1update.JSON200.Content)
			r.NotNil(node1update.JSON200.Link)
			a.Equal(url1, node1update.JSON200.Link.Url)
			a.Equal(prop1, node1update.JSON200.Meta)

			t.Run("empty_slug", func(t *testing.T) {
				name2 := "Testing Node Number Two" + uuid.NewString()
				slug2 := ""
				node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       name2,
					Slug:       &slug2,
					Visibility: &visibility,
				}, sh.WithSession(ctx))
				tests.Ok(t, err, node2)

				a.Equal(name2, node2.JSON200.Name)
				a.Contains(node2.JSON200.Slug, "testing-node-number-two")
				a.Equal("", node2.JSON200.Description)
				a.Nil(node2.JSON200.Content)
				a.Nil(node2.JSON200.Link)
				a.Equal(acc.ID.String(), string(node2.JSON200.Owner.Id))
			})

			t.Run("remove_link", func(t *testing.T) {
				name3 := "test-node-with-link" + uuid.NewString()
				slug3 := name3
				url3 := "https://example.com"
				node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       name3,
					Slug:       &slug3,
					Url:        &url3,
					Visibility: &visibility,
				}, sh.WithSession(ctx))
				tests.Ok(t, err, node3)

				r.NotNil(node3.JSON200.Link)
				a.Equal(url3, node3.JSON200.Link.Url)

				node3nolink, err := cl.NodeUpdateWithResponse(ctx, node3.JSON200.Slug, openapi.NodeMutableProps{
					Url: nullable.NewNullNullable[string](),
				}, sh.WithSession(ctx))
				tests.Ok(t, err, node3nolink)

				a.Nil(node3nolink.JSON200.Link)

				node3get, err := cl.NodeGetWithResponse(ctx, node3.JSON200.Slug, &openapi.NodeGetParams{})
				tests.Ok(t, err, node3get)

				a.Nil(node3get.JSON200.Link)
			})

			t.Run("with_russian_slug", func(t *testing.T) {
				a := assert.New(t)

				nodeCreate, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Slug: opt.New("бабочки").Ptr(),
					Name: "Бабочки",
				}, sh.WithSession(ctx))
				tests.Ok(t, err, nodeCreate)
				a.Contains(nodeCreate.JSON200.Slug, "бабочки")

				nodeGet, err := cl.NodeGetWithResponse(ctx, nodeCreate.JSON200.Slug, nil, sh.WithSession(ctx))
				tests.Ok(t, err, nodeGet)
				a.Equal("Бабочки", nodeGet.JSON200.Name)
				a.Contains(nodeGet.JSON200.Slug, "бабочки")
			})

			t.Run("with_russian_title_slugified", func(t *testing.T) {
				a := assert.New(t)

				nodeCreate, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: "БАБОЧКИ Example",
				}, sh.WithSession(ctx))
				tests.Ok(t, err, nodeCreate)
				a.Contains(nodeCreate.JSON200.Slug, "бабочки-example")

				nodeGet, err := cl.NodeGetWithResponse(ctx, nodeCreate.JSON200.Slug, nil, sh.WithSession(ctx))
				tests.Ok(t, err, nodeGet)
				a.Equal("БАБОЧКИ Example", nodeGet.JSON200.Name)
			})
		}))
	}))
}

func TestNodesErrors(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			get404, err := cl.NodeGetWithResponse(ctx, "nonexistent", &openapi.NodeGetParams{})
			r.NoError(err)
			r.NotNil(get404)
			a.Equal(http.StatusNotFound, get404.StatusCode())

			update403, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{})
			r.NoError(err)
			r.NotNil(update403)
			a.Equal(http.StatusUnauthorized, update403.StatusCode())

			update404, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{}, sh.WithSession(ctx))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(http.StatusNotFound, update404.StatusCode())

			t.Run("invalid_slug", func(t *testing.T) {
				name := "Testing Node Bad Slug" + uuid.NewString()
				slug := "not@a/good'slug]"
				create, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: name,
					Slug: &slug,
				}, sh.WithSession(ctx))
				tests.Status(t, err, create, http.StatusBadRequest)
			})
		}))
	}))
}
