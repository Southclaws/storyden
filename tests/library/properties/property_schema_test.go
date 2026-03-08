package properties_test

import (
	"context"
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

func TestNodesPropertySchemas_Create(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			parentname := "parent"
			parentslug := parentname + uuid.NewString()
			parent, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentname,
				Slug: &parentslug,
			}, session)
			tests.Ok(t, err, parent)

			// Add a child node

			name1 := "child-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:   name1,
				Slug:   &slug1,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node1)

			name2 := "child-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:   name2,
				Slug:   &slug2,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node2)

			name3 := "child-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:   name3,
				Slug:   &slug3,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node3)

			t.Run("schema", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				res, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parent.JSON200.Slug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
					{Name: "weight", Type: openapi.PropertyTypeNumber, Sort: "1"},
				}, session)
				tests.Ok(t, err, res)

				parent, err := cl.NodeGetWithResponse(ctx, parent.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, parent)
				r.Equal(1, len(parent.JSON200.ChildPropertySchema))
				matchSchema(t, openapi.PropertySchemaList{
					{Name: "weight", Type: openapi.PropertyTypeNumber, Sort: "1"},
				}, parent.JSON200.ChildPropertySchema)

				// Update the schema

				weightField := parent.JSON200.ChildPropertySchema[0]

				res, err = cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parent.JSON200.Slug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
					{Fid: &weightField.Fid, Name: "weight", Type: openapi.PropertyTypeNumber, Sort: "1"},
					{Name: "kind", Type: openapi.PropertyTypeText, Sort: "2"},
					{Name: "added", Type: openapi.PropertyTypeTimestamp, Sort: "3"},
				}, session)
				tests.Ok(t, err, res)

				parent, err = cl.NodeGetWithResponse(ctx, parent.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, parent)
				r.Equal(3, len(parent.JSON200.ChildPropertySchema))
				matchSchema(t, openapi.PropertySchemaList{
					{Name: "weight", Type: openapi.PropertyTypeNumber, Sort: "1"},
					{Name: "kind", Type: openapi.PropertyTypeText, Sort: "2"},
					{Name: "added", Type: openapi.PropertyTypeTimestamp, Sort: "3"},
				}, parent.JSON200.ChildPropertySchema)

				// Delete the schema

				res, err = cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parent.JSON200.Slug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{}, session)
				tests.Ok(t, err, res)

				parent, err = cl.NodeGetWithResponse(ctx, parent.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, parent)
				r.Equal(0, len(parent.JSON200.ChildPropertySchema))
				a.Empty(parent.JSON200.ChildPropertySchema)
			})
		}))
	}))
}

func TestNodesPropertySchemas_EmptyNode(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			parentname := "parent"
			parentslug := parentname + uuid.NewString()
			parent, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentname,
				Slug: &parentslug,
			}, session)
			tests.Ok(t, err, parent)

			t.Run("empty_node_schema_update", func(t *testing.T) {
				r := require.New(t)

				res, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parent.JSON200.Slug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
					{Name: "weight", Type: openapi.PropertyTypeNumber, Sort: "1"},
				}, session)
				tests.Ok(t, err, res)

				r.Empty(res.JSON200.Properties, "updating property schema on a node with no children should return empty properties")

				parentRefreshed, err := cl.NodeGetWithResponse(ctx, parent.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, parentRefreshed)
				r.Empty(parentRefreshed.JSON200.ChildPropertySchema, "child property schema should remain empty when node has no children")
			})
		}))
	}))
}

func matchSchema(t *testing.T, want openapi.PropertySchemaList, got openapi.PropertySchemaList) {
	t.Helper()
	a := assert.New(t)
	r := require.New(t)

	r.Equal(len(want), len(got))
	for i, s := range got {
		a.Equal(want[i].Name, s.Name)
		a.Equal(want[i].Type, s.Type)
		a.Equal(want[i].Sort, s.Sort)
	}
}
