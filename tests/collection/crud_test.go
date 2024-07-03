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
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
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
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, ar, seed.Account_001_Odin)

			acc1, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{xid.New().String(), "password"})
			tests.Ok(t, err, acc1)
			session1 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))), cj)

			acc2, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{xid.New().String(), "password"})
			tests.Ok(t, err, acc2)
			session2 := e2e.WithSession(session.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))), cj)

			cat1, err := cl.CategoryCreateWithResponse(adminCtx, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "",
				Description: "cat",
				Name:        xid.New().String(),
			}, e2e.WithSession(adminCtx, cj))
			tests.Ok(t, err, cat1)

			t.Run("unauthenticated", func(t *testing.T) {
				t.Parallel()

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: "c1 desc",
				})
				tests.Status(t, err, col, http.StatusForbidden)
			})

			t.Run("create", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: "c1 desc",
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", col.JSON200.Description)

				// owner
				get1, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, get1)
				a.Equal("c1", get1.JSON200.Name)
				a.Equal("c1 desc", get1.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get1.JSON200.Owner.Id)

				// another user
				get2, err := cl.CollectionGetWithResponse(root, col.JSON200.Id, session2)
				tests.Ok(t, err, get2)
				a.Equal("c1", get2.JSON200.Name)
				a.Equal("c1 desc", get2.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get2.JSON200.Owner.Id)

				// anonymous/guest
				get3, err := cl.CollectionGetWithResponse(root, col.JSON200.Id)
				tests.Ok(t, err, get3)
				a.Equal("c1", get3.JSON200.Name)
				a.Equal("c1 desc", get3.JSON200.Description)
				a.Equal(acc1.JSON200.Id, get3.JSON200.Owner.Id)
			})

			t.Run("update", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: "c1 desc",
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", col.JSON200.Description)

				id := col.JSON200.Id

				update1, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				}, session1)
				tests.Ok(t, err, update1)
				a.Equal("new name", update1.JSON200.Name)
				a.Equal("new desc", update1.JSON200.Description)

				col1get, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Ok(t, err, col1get)
				a.Equal("new name", col1get.JSON200.Name)
				a.Equal("new desc", col1get.JSON200.Description)

				updateUnauthorised, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				}, session2)
				tests.Status(t, err, updateUnauthorised, http.StatusUnauthorized)

				updateUnauthenticated, err := cl.CollectionUpdateWithResponse(root, id, openapi.CollectionUpdateJSONRequestBody{
					Description: opt.New("new desc").Ptr(),
					Name:        opt.New("new name").Ptr(),
				})
				tests.Status(t, err, updateUnauthenticated, http.StatusForbidden)
			})

			t.Run("delete", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: "c1 desc",
				}, session1)
				tests.Ok(t, err, col)

				id := col.JSON200.Id

				get1, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Ok(t, err, get1)
				a.Equal("c1", get1.JSON200.Name)
				a.Equal("c1 desc", get1.JSON200.Description)

				del, err := cl.CollectionDeleteWithResponse(root, col.JSON200.Id, session1)
				tests.Ok(t, err, del)

				get2, err := cl.CollectionGetWithResponse(root, id, session1)
				tests.Status(t, err, get2, http.StatusNotFound)
			})

			t.Run("list", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)

				col, err := cl.CollectionCreateWithResponse(root, openapi.CollectionCreateJSONRequestBody{
					Name:        "c1",
					Description: "c1 desc",
				}, session1)
				tests.Ok(t, err, col)
				a.Equal("c1", col.JSON200.Name)
				a.Equal("c1 desc", col.JSON200.Description)

				id := col.JSON200.Id

				list1, err := cl.CollectionListWithResponse(root)
				tests.Ok(t, err, list1)

				foundCol, foundOk := lo.Find(list1.JSON200.Collections, func(c openapi.Collection) bool { return c.Id == id })
				a.True(foundOk)
				a.Equal("c1", foundCol.Name)
				a.Equal("c1 desc", foundCol.Description)
				a.Equal(acc1.JSON200.Id, foundCol.Owner.Id)
			})
		}))
	}))
}
