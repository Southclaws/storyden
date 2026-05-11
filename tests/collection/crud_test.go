package collection_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestCollectionCRUD(t *testing.T) {
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

			handle1 := xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{Identifier: handle1, Token: "password"})
			tests.Ok(t, err, acc1)
			session1 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))))

			handle2 := xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{Identifier: handle2, Token: "password"})
			tests.Ok(t, err, acc2)
			session2 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))))

			cat1, err := cl.CategoryCreateWithResponse(adminCtx, openapi.CategoryInitialProps{
				Colour:      "",
				Description: "cat",
				Name:        xid.New().String(),
			}, sh.WithSession(adminCtx))
			tests.Ok(t, err, cat1)

			t.Run("unauthenticated", func(t *testing.T) {
				t.Parallel()

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name: "c1",
				})
				tests.Status(t, err, col, http.StatusUnauthorized)
			})

			t.Run("create", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				desc := "c1 desc"
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: &desc,
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", *col.JSON200.Description)

				// owner
				get1, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, get1)
				a.Equal("c1", get1.JSON200.Name)
				a.Equal("c1 desc", *get1.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get1.JSON200.Owner.Id)

				// another user
				get2, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session2)
				tests.Ok(t, err, get2)
				a.Equal("c1", get2.JSON200.Name)
				a.Equal("c1 desc", *get2.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get2.JSON200.Owner.Id)

				// anonymous/guest
				get3, err := cl.CollectionGetWithResponse(root, col.JSON200.Id)
				tests.Ok(t, err, get3)
				a.Equal("c1", get3.JSON200.Name)
				a.Equal("c1 desc", *get3.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get3.JSON200.Owner.Id)
			})

			t.Run("update", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				desc := "c1 desc"
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: &desc,
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", *col.JSON200.Description)

				id := col.JSON200.Id

				update1, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				}, session1)
				tests.Ok(t, err, update1)
				a.Equal("new name", update1.JSON200.Name)
				a.Equal("new desc", *update1.JSON200.Description)

				col1get, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Ok(t, err, col1get)
				a.Equal("new name", col1get.JSON200.Name)
				a.Equal("new desc", *col1get.JSON200.Description)

				updateUnauthorised, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				}, session2)
				tests.Status(t, err, updateUnauthorised, http.StatusForbidden)

				updateUnauthenticated, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				})
				tests.Status(t, err, updateUnauthenticated, http.StatusUnauthorized)
			})

			t.Run("delete", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				desc := "c1 desc"
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: &desc,
				}, session1)
				tests.Ok(t, err, col)

				id := col.JSON200.Id

				get1, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Ok(t, err, get1)
				a.Equal("c1", get1.JSON200.Name)
				a.Equal("c1 desc", *get1.JSON200.Description)

				del, err := cl.CollectionDeleteWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, del)

				get2, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Status(t, err, get2, http.StatusNotFound)
			})

			t.Run("list", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				desc := "c1 desc"
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: &desc,
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", *col.JSON200.Description)

				id := col.JSON200.Id

				list1, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{})
				tests.Ok(t, err, list1)

				foundCol, foundOk := lo.Find(list1.JSON200.Collections, func(c openapi.Collection) bool { return c.Id == id })
				a.True(foundOk)
				a.Equal("c1", foundCol.Name)
				a.Equal("c1 desc", *foundCol.Description)
				a.Equal(acc1.JSON200.Id, foundCol.Owner.Id)
			})

			t.Run("list_filter_by_owner", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				desc := "c1 desc"
				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: &desc,
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", *col.JSON200.Description)

				id := col.JSON200.Id

				list1, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &handle1,
				})
				tests.Ok(t, err, list1)

				foundCol, foundOk := lo.Find(list1.JSON200.Collections, func(c openapi.Collection) bool { return c.Id == id })
				a.True(foundOk)
				a.Equal("c1", foundCol.Name)
				a.Equal("c1 desc", *foundCol.Description)
				a.Equal(acc1.JSON200.Id, foundCol.Owner.Id)

				list2, err := cl.CollectionListWithResponse(root, &openapi.CollectionListParams{
					AccountHandle: &handle2,
				})
				tests.Ok(t, err, list2)

				foundCol, foundOk = lo.Find(list2.JSON200.Collections, func(c openapi.Collection) bool { return c.Id == id })
				a.False(foundOk)
			})
		}))
	}))
}
