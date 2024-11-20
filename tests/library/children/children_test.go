package datagraph_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestLibraryNodeChildren(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			session := e2e.WithSession(ctx, cj)

			// iurl := "https://picsum.photos/500/500"

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name1,
				Slug: &slug1,
			}, session)
			tests.Ok(t, err, node1)

			// Add a child node

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name2,
				Slug: &slug2,
			}, session)
			tests.Ok(t, err, node2)

			cadd, err := cl.NodeAddNodeWithResponse(ctx, slug1, slug2, session)
			tests.Ok(t, err, cadd)

			r.Equal(node1.JSON200.Id, cadd.JSON200.Id)

			// Add another child to this child
			// node1
			// |- node2
			//    |- node3

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name3,
				Slug: &slug3,
			}, session)
			tests.Ok(t, err, node3)

			cadd, err = cl.NodeAddNodeWithResponse(ctx, slug2, slug3, session)
			tests.Ok(t, err, cadd)

			r.Equal(node2.JSON200.Id, cadd.JSON200.Id)
		}))
	}))
}
