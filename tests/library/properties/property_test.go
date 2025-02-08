package properties_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesProperty(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := e2e.WithSession(ctx1, cj)

			ctx2, _ := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			randomUser := e2e.WithSession(ctx2, cj)

			parentname := "parent"
			parentslug := parentname + uuid.NewString()
			parent, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name: parentname,
				Slug: &parentslug,
			}, session)
			tests.Ok(t, err, parent)

			// add 3 child nodes to parent

			name1 := "child-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:   name1,
				Slug:   &slug1,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node1)

			name2 := "child-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:   name2,
				Slug:   &slug2,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node2)

			name3 := "child-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:   name3,
				Slug:   &slug3,
				Parent: &parent.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node3)

			// add two child nodes to child-3

			name34 := "child-3-4"
			slug34 := name34 + uuid.NewString()
			node34, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:   name34,
				Slug:   &slug34,
				Parent: &node3.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node34)

			name35 := "child-3-5"
			slug35 := name35 + uuid.NewString()
			node35, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:   name35,
				Slug:   &slug35,
				Parent: &node3.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node35)

			res, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(root, parentslug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
				{Name: "weight", Type: "number", Sort: "1"},
				{Name: "kind", Type: "string", Sort: "2"},
				{Name: "added", Type: "timestamp", Sort: "3"},
			}, session)
			tests.Ok(t, err, res)

			res, err = cl.NodeUpdateChildrenPropertySchemaWithResponse(root, slug3, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
				{Name: "size", Type: "number", Sort: "1"},
				{Name: "brand", Type: "string", Sort: "2"},
			}, session)
			tests.Ok(t, err, res)

			t.Run("fail_unauthenticated", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				update, err := cl.NodeUpdatePropertiesWithResponse(root, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Value: "5"},
						{Name: "kind", Value: "mythical"},
						{Name: "added", Value: "2025-01-01T12:59:21Z"},
						{Name: "new", Value: "prop"},
					},
				})
				r.NoError(err)
				a.Equal(403, update.StatusCode())
			})

			t.Run("fail_no_permission", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				update, err := cl.NodeUpdatePropertiesWithResponse(root, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Value: "5"},
						{Name: "kind", Value: "mythical"},
						{Name: "added", Value: "2025-01-01T12:59:21Z"},
						{Name: "new", Value: "prop"},
					},
				}, randomUser)
				r.NoError(err)
				a.Equal(401, update.StatusCode())
			})

			t.Run("fail_missing_new_property_type", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				update, err := cl.NodeUpdatePropertiesWithResponse(root, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Value: "5"},
						{Name: "kind", Value: "mythical"},
						{Name: "added", Value: "2025-01-01T12:59:21Z"},
						{Name: "new", Value: "prop"},
					},
				}, session)
				r.NoError(err)
				a.Equal(400, update.StatusCode())
			})

			t.Run("success_set_property_values", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				update, err := cl.NodeUpdatePropertiesWithResponse(root, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Value: "4"},
						{Name: "kind", Value: "legendary"},
						{Name: "added", Value: "2025-02-06T20:59:21Z"},
					},
				}, session)
				r.NoError(err)
				a.Equal(200, update.StatusCode())
			})

			t.Run("success_add_new_property", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				ptype := "string"

				update, err := cl.NodeUpdatePropertiesWithResponse(root, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Value: "5"},
						{Name: "kind", Value: "mythical"},
						{Name: "added", Value: "2025-01-01T12:59:21Z"},
						{Name: "new", Value: "prop", Type: &ptype},
					},
				}, session)
				r.NoError(err)
				a.Equal(200, update.StatusCode())
			})
		}))
	}))
}

func TestNodesPropertyFieldOrdering(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			session := e2e.WithSession(ctx, cj)

			parentname := "parent"
			parentslug := parentname + uuid.NewString()
			parent, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentname,
				Slug: &parentslug,
			}, session)
			tests.Ok(t, err, parent)

			// add 3 child nodes to parent

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

			res, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parentslug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
				{Name: "weight", Type: "number", Sort: "1"},
				{Name: "kind", Type: "string", Sort: "2"},
				{Name: "added", Type: "timestamp", Sort: "3"},
			}, session)
			tests.Ok(t, err, res)

			t.Run("sort_fields", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				update, err := cl.NodeUpdatePropertiesWithResponse(ctx, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
					Properties: openapi.PropertyMutationList{
						{Name: "added", Value: "2025-01-01T12:59:21Z"},
						{Name: "weight", Value: "5"},
						{Name: "kind", Value: "mythical"},
					},
				}, session)
				r.NoError(err)
				a.Equal(200, update.StatusCode())

				fields := dt.Map(update.JSON200.Properties, func(p openapi.Property) string {
					return p.Name
				})

				fieldIDs := dt.Map(update.JSON200.Properties, func(p openapi.Property) string {
					return p.Fid
				})

				a.Equal([]string{"weight", "kind", "added"}, fields)

				// re-order the fields
				// assuming the above assertion is correct, these field IDs are

				field1ID := &fieldIDs[0]
				field2ID := &fieldIDs[1]
				field3ID := &fieldIDs[2]

				schemaUpdate, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parentslug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
					{Fid: field1ID, Name: "weight", Type: "number", Sort: "3"},
					{Fid: field2ID, Name: "kind", Type: "string", Sort: "1"},
					{Fid: field3ID, Name: "added", Type: "timestamp", Sort: "2"},
				}, session)
				tests.Ok(t, err, schemaUpdate)

				updatedFields := dt.Map(schemaUpdate.JSON200.Properties, func(p openapi.PropertySchema) string {
					return p.Name
				})

				a.Equal([]string{"kind", "added", "weight"}, updatedFields)
			})
		}))
	}))
}

func TestNodesPropertySchemaOnParentAndChildNodes(t *testing.T) {
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
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			session := e2e.WithSession(ctx, cj)

			parentname := "parent"
			parentslug := parentname + uuid.NewString()
			parent, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: parentname,
				Slug: &parentslug,
			}, session)
			tests.Ok(t, err, parent)

			// add 3 child nodes to parent

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

			// add two child nodes to child-3

			name34 := "child-3-4"
			slug34 := name34 + uuid.NewString()
			node34, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:   name34,
				Slug:   &slug34,
				Parent: &node3.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node34)

			name35 := "child-3-5"
			slug35 := name35 + uuid.NewString()
			node35, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:   name35,
				Slug:   &slug35,
				Parent: &node3.JSON200.Slug,
			}, session)
			tests.Ok(t, err, node35)

			res, err := cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, parentslug, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
				{Name: "weight", Type: "number", Sort: "1"},
				{Name: "kind", Type: "string", Sort: "2"},
				{Name: "added", Type: "timestamp", Sort: "3"},
			}, session)
			tests.Ok(t, err, res)

			res, err = cl.NodeUpdateChildrenPropertySchemaWithResponse(ctx, slug3, openapi.NodeUpdateChildrenPropertySchemaJSONRequestBody{
				{Name: "size", Type: "number", Sort: "1"},
				{Name: "brand", Type: "string", Sort: "2"},
			}, session)
			tests.Ok(t, err, res)

			update, err := cl.NodeUpdatePropertiesWithResponse(ctx, slug1, openapi.NodeUpdatePropertiesJSONRequestBody{
				Properties: openapi.PropertyMutationList{
					{Name: "weight", Value: "4"},
					{Name: "kind", Value: "legendary"},
					{Name: "added", Value: "2025-02-06T20:59:21Z"},
				},
			}, session)
			r.NoError(err)
			a.Equal(200, update.StatusCode())

			t.Run("assert_fields_and_data_in_parent_children", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// get the parent
				parent, err := cl.NodeGetWithResponse(ctx, parentslug, session)
				r.NoError(err)
				r.NotNil(parent)
				// yield the children's schema
				a.Len(parent.JSON200.ChildPropertySchema, 3)

				// get the child
				child, err := cl.NodeGetWithResponse(ctx, slug3, session)
				r.NoError(err)
				r.NotNil(child)
				// yield the children's schema and its own schema
				a.Len(child.JSON200.ChildPropertySchema, 2)
				r.Equal(len(child.JSON200.Properties), len(parent.JSON200.ChildPropertySchema))
				matchSchemaToProperties(t, parent.JSON200.ChildPropertySchema, child.JSON200.Properties)
			})
		}))
	}))
}

func matchSchemaToProperties(t *testing.T, schema openapi.PropertySchemaList, properties []openapi.Property) {
	a := assert.New(t)

	for i := range schema {
		schemaField := schema[i]
		propertyField := properties[i]

		a.Equal(schemaField.Fid, propertyField.Fid)
		a.Equal(schemaField.Name, propertyField.Name)
		a.Equal(schemaField.Type, propertyField.Type)
		a.Equal(schemaField.Sort, propertyField.Sort)
	}
	return
}
