package thread_tags_test

import (
	"context"
	"testing"

	"go.uber.org/fx"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadTags(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := e2e.WithSession(adminCtx, cj)

			cat, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: "Category " + uuid.NewString()}, adminSession)
			tests.Ok(t, err, cat)
			catID := cat.JSON200.Id

			t.Run("create_with_new_tags", func(t *testing.T) {
				t.Parallel()

				a := assert.New(t)
				r := require.New(t)

				t1 := xid.New().String()
				t2 := xid.New().String()
				t3 := xid.New().String()

				tags := []string{t1, t2, t3}
				create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      un("n1"),
					Category:   catID,
					Visibility: openapi.Published,
					Tags:       &tags,
				}, adminSession)
				tests.Ok(t, err, create)
				r.NotEmpty(create.JSON200.Tags)
				f := find(create.JSON200.Tags)
				a.True(f(t1))
				a.True(f(t2))
				a.True(f(t3))
			})

			t.Run("create_with_existing_tags", func(t *testing.T) {
				t.Parallel()

				a := assert.New(t)
				r := require.New(t)

				t1 := xid.New().String()
				t2 := xid.New().String()
				t3 := xid.New().String()

				n1tags := []string{t1, t2, t3}
				create1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      un("n1"),
					Category:   catID,
					Visibility: openapi.Published,
					Tags:       &n1tags,
				}, adminSession)
				tests.Ok(t, err, create1)

				t4 := xid.New().String()
				n2tags := []string{t2, t3, t4}
				create2, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      un("n1"),
					Category:   catID,
					Visibility: openapi.Published,
					Tags:       &n2tags,
				}, adminSession)
				tests.Ok(t, err, create2)
				r.NotEmpty(create2.JSON200.Tags)
				f := find(create2.JSON200.Tags)
				a.False(f(t1))
				a.True(f(t2))
				a.True(f(t3))
				a.True(f(t4))
			})

			t.Run("update_tags", func(t *testing.T) {
				t.Parallel()

				a := assert.New(t)
				r := require.New(t)

				t1 := xid.New().String()
				t2 := xid.New().String()
				t3 := xid.New().String()

				tags := []string{t1, t2}
				create1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      un("n1"),
					Category:   catID,
					Visibility: openapi.Published,
					Tags:       &tags,
				}, adminSession)
				tests.Ok(t, err, create1)
				r.NotEmpty(create1.JSON200.Tags)
				f := find(create1.JSON200.Tags)
				a.True(f(t1))
				a.True(f(t2))
				a.False(f(t3))

				newTags := []string{t1, t2, t3}
				create2, err := cl.ThreadUpdateWithResponse(root, create1.JSON200.Slug, openapi.ThreadMutableProps{
					Tags: &newTags,
				}, adminSession)
				tests.Ok(t, err, create2)
				r.NotEmpty(create2.JSON200.Tags)
				f = find(create2.JSON200.Tags)
				a.True(f(t1))
				a.True(f(t2))
				a.True(f(t3))
			})

			t.Run("remove_tags", func(t *testing.T) {
				t.Parallel()

				a := assert.New(t)
				r := require.New(t)

				t1 := xid.New().String()
				t2 := xid.New().String()
				t3 := xid.New().String()

				tags := []string{t1, t2, t3}
				create1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      un("n1"),
					Category:   catID,
					Visibility: openapi.Published,
					Tags:       &tags,
				}, adminSession)
				tests.Ok(t, err, create1)

				newTags := []string{t1, t2}
				create2, err := cl.ThreadUpdateWithResponse(root, create1.JSON200.Slug, openapi.ThreadMutableProps{
					Tags: &newTags,
				}, adminSession)
				tests.Ok(t, err, create2)
				r.NotEmpty(create2.JSON200.Tags)
				f := find(create2.JSON200.Tags)
				a.True(f(t1))
				a.True(f(t2))
				a.False(f(t3))
			})
		}))
	}))
}

func un(n string) string {
	return n + " " + xid.New().String()
}

func find(tags []openapi.TagReference) func(string) bool {
	return func(n string) bool {
		_, ok := lo.Find(tags, func(t openapi.TagReference) bool { return t.Name == n })
		return ok
	}
}
