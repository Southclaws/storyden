package semdex_weaviate_test

import (
	"context"
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
)

//go:embed testdata/*.txt
var content embed.FS

const dir = "testdata"

func TestSemdexWeaviate(t *testing.T) {
	t.Parallel()

	os.Setenv("WEAVIATE_ENABLED", "true")
	os.Setenv("WEAVIATE_URL", "http://localhost:8080")

	integration.Test(t, &config.Config{}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)

			es, err := content.ReadDir(dir)
			r.NoError(err)

			cat1create, err := cl.CategoryCreateWithResponse(ctx, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "#fe4efd",
				Description: "TestSemdexWeaviate",
				Name:        uuid.NewString(),
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cat1create)
			r.Equal(http.StatusOK, cat1create.StatusCode())

			for _, e := range es {
				filename := filepath.Join(dir, e.Name())
				title := strings.TrimSuffix(e.Name(), ".txt")

				b, err := os.ReadFile(filename)
				r.NoError(err)

				response, err := cl.ThreadCreateWithResponse(ctx, openapi.ThreadInitialProps{
					Title:    title,
					Category: cat1create.JSON200.Id,
					Body:     string(b),
					Status:   openapi.Published,
				}, e2e.WithSession(ctx, cj))
				r.NoError(err)
				r.Equal(http.StatusOK, response.StatusCode())

				a.Equal(seed.Account_001_Odin.Name, response.JSON200.Author.Name)
			}

			query := "outage"

			search1, err := cl.DatagraphSearchWithResponse(ctx, &openapi.DatagraphSearchParams{
				Q: &query,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, search1.StatusCode())

			a.Len(search1.JSON200.Items, 6)
		}))
	}))
}
