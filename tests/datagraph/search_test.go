package search_test

import (
	"context"
	"encoding/json"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/samber/lo"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestSearchMultipleKinds(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{}, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		lc fx.Lifecycle,
		cfg config.Config,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw *account_writer.Writer,
	) {
		if cfg.SemdexEnabled {
			return
		}

		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := e2e.WithSession(adminCtx, cj)
			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := e2e.WithSession(ctx1, cj)
			ctx2, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			session2 := e2e.WithSession(ctx2, cj)

			cat1, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Name: uuid.NewString(), Colour: "#000"}, adminSession)
			tests.Ok(t, err, cat1)

			published := openapi.Published
			draft := openapi.Draft

			hot := "<p>this contains the keyword we want</p>"
			cold := "<p>this contains none of the words we want</p>"

			t1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       hot,
				Category:   cat1.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "thread",
			}, session1)
			tests.Ok(t, err, t1)
			t2, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       cold,
				Category:   cat1.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "thread",
			}, session2)
			tests.Ok(t, err, t2)
			n1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "node 1" + uuid.NewString(),
				Content:    &hot,
				Visibility: &published,
			}, adminSession)
			tests.Ok(t, err, n1)
			n2, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "node 2" + uuid.NewString(),
				Content:    &cold,
				Visibility: &draft,
			}, session2)
			tests.Ok(t, err, n2)

			t.Run("search_all_kinds", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				search1, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q: "keyword",
				}, session1)
				tests.Ok(t, err, search1)

				foundt1, err := findItem(search1.JSON200.Items, t1.JSON200.Id).AsDatagraphItemThread()
				r.NoError(err)
				a.Equal(t1.JSON200.Id, foundt1.Ref.Id)
				a.Equal(t1.JSON200.CreatedAt, foundt1.Ref.CreatedAt)
				a.Equal(t1.JSON200.DeletedAt, foundt1.Ref.DeletedAt)
				a.Equal(t1.JSON200.Title, foundt1.Ref.Title)
				a.Equal(t1.JSON200.Slug, foundt1.Ref.Slug)
				a.Equal(t1.JSON200.Body, foundt1.Ref.Body)
				a.Equal(t1.JSON200.Description, foundt1.Ref.Description)
				// a.Equal(t1.JSON200.Category, foundt1.Ref.Category) // NOTE: Not implemented yet because postsearcher doesn't distinguish properly between threads and posts.
				a.Equal(t1.JSON200.Author, foundt1.Ref.Author)

				foundn1, err := findItem(search1.JSON200.Items, n1.JSON200.Id).AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(n1.JSON200.Id, foundn1.Ref.Id)
				a.Equal(n1.JSON200.CreatedAt, foundn1.Ref.CreatedAt)
				a.Equal(n1.JSON200.DeletedAt, foundn1.Ref.DeletedAt)
				a.Equal(n1.JSON200.Name, foundn1.Ref.Name)
				a.Equal(n1.JSON200.Slug, foundn1.Ref.Slug)
				a.Equal(n1.JSON200.Content, foundn1.Ref.Content)
				a.Equal(n1.JSON200.Description, foundn1.Ref.Description)
				a.Equal(n1.JSON200.Owner, foundn1.Ref.Owner)

				foundt2 := findItem(search1.JSON200.Items, t2.JSON200.Id)
				r.Nil(foundt2)

				foundn2 := findItem(search1.JSON200.Items, n2.JSON200.Id)
				r.Nil(foundn2)
			})

			t.Run("search_only_threads", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				search1, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q: "keyword",
					Kind: &[]openapi.DatagraphItemKind{
						openapi.DatagraphItemKindThread,
					},
				}, session1)
				tests.Ok(t, err, search1)

				foundt1, err := findItem(search1.JSON200.Items, t1.JSON200.Id).AsDatagraphItemThread()
				r.NoError(err)
				a.Equal(t1.JSON200.Id, foundt1.Ref.Id)
				a.Equal(t1.JSON200.CreatedAt, foundt1.Ref.CreatedAt)
				a.Equal(t1.JSON200.DeletedAt, foundt1.Ref.DeletedAt)
				a.Equal(t1.JSON200.Title, foundt1.Ref.Title)
				a.Equal(t1.JSON200.Slug, foundt1.Ref.Slug)
				a.Equal(t1.JSON200.Body, foundt1.Ref.Body)
				a.Equal(t1.JSON200.Description, foundt1.Ref.Description)
				// a.Equal(t1.JSON200.Category, foundt1.Ref.Category) // NOTE: Not implemented yet because postsearcher doesn't distinguish properly between threads and posts.
				a.Equal(t1.JSON200.Author, foundt1.Ref.Author)

				foundn1 := findItem(search1.JSON200.Items, n1.JSON200.Id)
				r.Nil(foundn1)

				foundt2 := findItem(search1.JSON200.Items, t2.JSON200.Id)
				r.Nil(foundt2)

				foundn2 := findItem(search1.JSON200.Items, n2.JSON200.Id)
				r.Nil(foundn2)
			})

			t.Run("search_only_nodes", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				search1, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q: "keyword",
					Kind: &[]openapi.DatagraphItemKind{
						openapi.DatagraphItemKindNode,
					},
				}, session1)
				tests.Ok(t, err, search1)

				foundt1 := findItem(search1.JSON200.Items, t1.JSON200.Id)
				r.Nil(foundt1)

				foundn1, err := findItem(search1.JSON200.Items, n1.JSON200.Id).AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(n1.JSON200.Id, foundn1.Ref.Id)
				a.Equal(n1.JSON200.CreatedAt, foundn1.Ref.CreatedAt)
				a.Equal(n1.JSON200.DeletedAt, foundn1.Ref.DeletedAt)
				a.Equal(n1.JSON200.Name, foundn1.Ref.Name)
				a.Equal(n1.JSON200.Slug, foundn1.Ref.Slug)
				a.Equal(n1.JSON200.Content, foundn1.Ref.Content)
				a.Equal(n1.JSON200.Description, foundn1.Ref.Description)
				a.Equal(n1.JSON200.Owner, foundn1.Ref.Owner)

				foundt2 := findItem(search1.JSON200.Items, t2.JSON200.Id)
				r.Nil(foundt2)

				foundn2 := findItem(search1.JSON200.Items, n2.JSON200.Id)
				r.Nil(foundn2)
			})
		}))
	}))
}

func TestSearchVisibilityRules(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{}, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		lc fx.Lifecycle,
		cfg config.Config,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw *account_writer.Writer,
	) {
		if cfg.SemdexEnabled {
			return
		}

		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := e2e.WithSession(adminCtx, cj)
			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := e2e.WithSession(ctx1, cj)
			ctx2, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			session2 := e2e.WithSession(ctx2, cj)

			cat1, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Name: uuid.NewString(), Colour: "#000"}, adminSession)
			tests.Ok(t, err, cat1)

			published := openapi.Published
			draft := openapi.Draft

			hot := "<p>this contains the keyword we want</p>"
			cold := "<p>this contains the keyword we want but it's not published</p>"

			t1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       hot,
				Category:   cat1.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "thread",
			}, session1)
			tests.Ok(t, err, t1)
			t2, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       cold,
				Category:   cat1.JSON200.Id,
				Visibility: openapi.Draft,
				Title:      "thread",
			}, session2)
			tests.Ok(t, err, t2)
			n1, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "node 1" + uuid.NewString(),
				Content:    &hot,
				Visibility: &published,
			}, adminSession)
			tests.Ok(t, err, n1)
			n2, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "node 2" + uuid.NewString(),
				Content:    &cold,
				Visibility: &draft,
			}, session2)
			tests.Ok(t, err, n2)

			t.Run("only_published", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				search1, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q: "keyword",
				}, session1)
				tests.Ok(t, err, search1)

				foundt1, err := findItem(search1.JSON200.Items, t1.JSON200.Id).AsDatagraphItemThread()
				r.NoError(err)
				a.Equal(t1.JSON200.Id, foundt1.Ref.Id)

				foundn1, err := findItem(search1.JSON200.Items, n1.JSON200.Id).AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(n1.JSON200.Id, foundn1.Ref.Id)

				foundt2 := findItem(search1.JSON200.Items, t2.JSON200.Id)
				r.Nil(foundt2)

				foundn2 := findItem(search1.JSON200.Items, n2.JSON200.Id)
				r.Nil(foundn2)
			})
		}))
	}))
}

func findItem(items []openapi.DatagraphItem, id openapi.Identifier) *openapi.DatagraphItem {
	found, ok := lo.Find(items, func(i openapi.DatagraphItem) bool {
		iid := coerceDatagraphItem(i)
		return iid == id
	})
	if !ok {
		return nil
	}
	return &found
}

func coerceDatagraphItem(i openapi.DatagraphItem) string {
	b, err := i.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var v struct {
		Ref struct {
			ID string `json:"id"`
		} `json:"ref"`
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		panic(err)
	}

	return v.Ref.ID
}
