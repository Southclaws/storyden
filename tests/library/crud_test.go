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
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
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
		cj *session.Jar,
		aw account_writer.Writer,
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
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node1)

			a.Equal(name1, node1.JSON200.Name)
			a.Equal(slug1, node1.JSON200.Slug)
			a.Equal("Rich text content.", node1.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Rich text content.</p></body>", *node1.JSON200.Content)
			r.NotNil(node1.JSON200.Link)
			a.Equal(url1, node1.JSON200.Link.Url)
			a.Equal(acc.ID.String(), string(node1.JSON200.Owner.Id))

			// Get the one just created

			node1get, err := cl.NodeGetWithResponse(ctx, slug1)
			tests.Ok(t, err, node1get)

			a.Equal(name1, node1get.JSON200.Name)
			a.Equal(slug1, node1get.JSON200.Slug)
			a.Equal("Rich text content.", node1get.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Rich text content.</p></body>", *node1get.JSON200.Content)
			a.Equal(acc.ID.String(), string(node1get.JSON200.Owner.Id))

			// Update the one just created

			name1 = "test-node-1-UPDATED"
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
				Url:     &url1,
				Meta:    &prop1,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node1update)

			a.Equal(name1, node1update.JSON200.Name)
			a.Equal(slug1, node1update.JSON200.Slug)
			a.Equal("Newly changed content.", node1update.JSON200.Description)
			a.Equal("<body><h1>Nodes</h1><p>Newly changed content.</p></body>", *node1update.JSON200.Content)
			r.NotNil(node1update.JSON200.Link)
			a.Equal(url1, node1update.JSON200.Link.Url)
			a.Equal(prop1, node1update.JSON200.Meta)
		}))
	}))
}

func TestNodesErrors(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			get404, err := cl.NodeGetWithResponse(ctx, "nonexistent")
			r.NoError(err)
			r.NotNil(get404)
			a.Equal(http.StatusNotFound, get404.StatusCode())

			update403, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{})
			r.NoError(err)
			r.NotNil(update403)
			a.Equal(http.StatusForbidden, update403.StatusCode())

			update404, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(http.StatusNotFound, update404.StatusCode())
		}))
	}))
}
