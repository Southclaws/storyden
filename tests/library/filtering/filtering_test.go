package library_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesFiltering(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *cookie.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			ctx1, acc1 := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx2, acc2 := e2e.WithAccount(ctx, aw, seed.Account_002_Frigg)

			visibility := openapi.Published

			name1 := "test-node-owned-by-1"
			slug1 := name1 + uuid.NewString()
			content1 := "# Nodes\n\nOwned by Odin."
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name1,
				Slug:       &slug1,
				Content:    &content1,
				Visibility: &visibility,
			}, e2e.WithSession(ctx1, cj))
			tests.Ok(t, err, node1)

			name2 := "test-node-owned-by-2"
			slug2 := name2 + uuid.NewString()
			content2 := "# Nodes\n\nOwned by Frigg."
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name2,
				Slug:       &slug2,
				Content:    &content2,
				Visibility: &visibility,
			}, e2e.WithSession(ctx2, cj))
			tests.Ok(t, err, node2)

			clist, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &acc1.Handle,
			})
			tests.Ok(t, err, clist)

			ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

			a.Contains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)

			clist2, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &acc2.Handle,
			})
			tests.Ok(t, err, clist2)

			ids2 := dt.Map(clist2.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

			a.NotContains(ids2, node1.JSON200.Id)
			a.Contains(ids2, node2.JSON200.Id)
		}))
	}))
}
