package properties_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
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

func TestNodesPropertySorting(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx1)

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
			node1 := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name1,
					Slug:   &slug1,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			res := tests.AssertRequest(
				cl.NodeUpdatePropertiesWithResponse(root, node1.JSON200.Slug, openapi.PropertyMutableProps{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Type: opt.New(openapi.PropertyTypeNumber).Ptr(), Value: "4", Sort: opt.New("1").Ptr()},
						{Name: "height", Type: opt.New(openapi.PropertyTypeNumber).Ptr(), Value: "6", Sort: opt.New("2").Ptr()},
						{Name: "nickname", Type: opt.New(openapi.PropertyTypeText).Ptr(), Value: "ahmed", Sort: opt.New("3").Ptr()},
					},
				}, session),
			)(t, http.StatusOK)

			// first mutation creates all the properties in the schema, so the
			// following 2 mutations need to use field IDs to actually set the
			// values because the API cannot discern between create and update.
			s1fieldIDs := dt.Map(res.JSON200.Properties, func(p openapi.Property) string {
				return p.Fid
			})

			s1field1ID := &s1fieldIDs[0]
			s1field2ID := &s1fieldIDs[1]
			s1field3ID := &s1fieldIDs[2]

			name2 := "child-2"
			slug2 := name2 + uuid.NewString()
			node2 := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name2,
					Slug:   &slug2,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			tests.AssertRequest(
				cl.NodeUpdatePropertiesWithResponse(root, node2.JSON200.Slug, openapi.PropertyMutableProps{
					Properties: openapi.PropertyMutationList{
						{Fid: s1field1ID, Value: "19"},
						{Fid: s1field2ID, Value: "1"},
						{Fid: s1field3ID, Value: "john"},
					},
				}, session),
			)(t, http.StatusOK)

			name3 := "child-3"
			slug3 := name3 + uuid.NewString()
			node3 := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name3,
					Slug:   &slug3,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			tests.AssertRequest(
				cl.NodeUpdatePropertiesWithResponse(root, node3.JSON200.Slug, openapi.PropertyMutableProps{
					Properties: openapi.PropertyMutationList{
						{Fid: s1field1ID, Value: "5"},
						{Fid: s1field2ID, Value: "69"},
						{Fid: s1field3ID, Value: "zenith"},
					},
				}, session),
			)(t, http.StatusOK)

			t.Run("sort_by_weight", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				n := tests.AssertRequest(cl.NodeGetWithResponse(root, parentslug, &openapi.NodeGetParams{
					ChildrenSort: opt.New("-weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Children, 3)
				slugs := dt.Map(n.JSON200.Children, bySlug)
				wantSlugs := []string{slug2, slug3, slug1}
				a.Equal(wantSlugs, slugs)

				n = tests.AssertRequest(cl.NodeGetWithResponse(root, parentslug, &openapi.NodeGetParams{
					ChildrenSort: opt.New("weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Children, 3)
				slugs = dt.Map(n.JSON200.Children, bySlug)
				wantSlugs = []string{slug1, slug3, slug2}
				a.Equal(wantSlugs, slugs)
			})

			t.Run("sort_by_weight_direct_children", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				n := tests.AssertRequest(cl.NodeListChildrenWithResponse(root, parentslug, &openapi.NodeListChildrenParams{
					ChildrenSort: opt.New("-weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Nodes, 3)
				slugs := dt.Map(n.JSON200.Nodes, bySlug)
				wantSlugs := []string{slug2, slug3, slug1}
				a.Equal(wantSlugs, slugs)

				n = tests.AssertRequest(cl.NodeListChildrenWithResponse(root, parentslug, &openapi.NodeListChildrenParams{
					ChildrenSort: opt.New("weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Nodes, 3)
				slugs = dt.Map(n.JSON200.Nodes, bySlug)
				wantSlugs = []string{slug1, slug3, slug2}
				a.Equal(wantSlugs, slugs)
			})
		}))
	}))
}

func TestNodesPropertySorting_WithEmptyValues(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx1)

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
			node1 := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name1,
					Slug:   &slug1,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			res := tests.AssertRequest(
				cl.NodeUpdatePropertiesWithResponse(root, node1.JSON200.Slug, openapi.PropertyMutableProps{
					Properties: openapi.PropertyMutationList{
						{Name: "weight", Type: opt.New(openapi.PropertyTypeNumber).Ptr(), Value: "4", Sort: opt.New("1").Ptr()},
						{Name: "height", Type: opt.New(openapi.PropertyTypeNumber).Ptr(), Value: "6", Sort: opt.New("2").Ptr()},
						{Name: "nickname", Type: opt.New(openapi.PropertyTypeText).Ptr(), Value: "ahmed", Sort: opt.New("3").Ptr()},
					},
				}, session),
			)(t, http.StatusOK)

			// first mutation creates all the properties in the schema, so the
			// following 2 mutations need to use field IDs to actually set the
			// values because the API cannot discern between create and update.
			s1fieldIDs := dt.Map(res.JSON200.Properties, func(p openapi.Property) string {
				return p.Fid
			})

			s1field1ID := &s1fieldIDs[0]
			s1field2ID := &s1fieldIDs[1]
			s1field3ID := &s1fieldIDs[2]

			name2 := "child-2"
			slug2 := name2 + uuid.NewString()
			node2 := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name2,
					Slug:   &slug2,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			tests.AssertRequest(
				cl.NodeUpdatePropertiesWithResponse(root, node2.JSON200.Slug, openapi.PropertyMutableProps{
					Properties: openapi.PropertyMutationList{
						{Fid: s1field1ID, Value: "19"},
						{Fid: s1field2ID, Value: "1"},
						{Fid: s1field3ID, Value: "john"},
					},
				}, session),
			)(t, http.StatusOK)

			name3 := "child-3"
			slug3 := name3 + uuid.NewString()
			tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:   name3,
					Slug:   &slug3,
					Parent: &parent.JSON200.Slug,
				}, session),
			)(t, http.StatusOK)

			// node 3 has no values, should be sorted to the end
			// tests.AssertRequest(
			// 	cl.NodeUpdatePropertiesWithResponse(root, node3.JSON200.Slug, openapi.PropertyMutableProps{
			// 		Properties: openapi.PropertyMutationList{
			// 			{Fid: s1field1ID, Value: "5"},
			// 			{Fid: s1field2ID, Value: "69"},
			// 			{Fid: s1field3ID, Value: "zenith"},
			// 		},
			// 	}, session),
			// )(t, http.StatusOK)

			t.Run("sort_by_weight_with_empties", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				n := tests.AssertRequest(cl.NodeGetWithResponse(root, parentslug, &openapi.NodeGetParams{
					ChildrenSort: opt.New("-weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Children, 3)
				slugs := dt.Map(n.JSON200.Children, bySlug)
				wantSlugs := []string{slug2, slug1, slug3}
				a.Equal(wantSlugs, slugs)

				n = tests.AssertRequest(cl.NodeGetWithResponse(root, parentslug, &openapi.NodeGetParams{
					ChildrenSort: opt.New("weight").Ptr(),
				}, session))(t, http.StatusOK)
				r.Len(n.JSON200.Children, 3)
				slugs = dt.Map(n.JSON200.Children, bySlug)
				wantSlugs = []string{slug1, slug2, slug3}
				a.Equal(wantSlugs, slugs)
			})
		}))
	}))
}

func bySlug(n openapi.NodeWithChildren) string { return n.Slug }
