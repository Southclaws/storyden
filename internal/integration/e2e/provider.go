package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Southclaws/fault"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/config"
	internal_http "github.com/Southclaws/storyden/internal/http"
	internal_openapi "github.com/Southclaws/storyden/internal/openapi"
)

func newHttpTestServer(lc fx.Lifecycle, l *zap.Logger, cfg config.Config, router *http.ServeMux) *httptest.Server {
	server := httptest.NewServer(router)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			server.Close()
			return nil
		},
	})

	return server
}

func newClient(ts *httptest.Server) (*internal_openapi.ClientWithResponses, error) {
	server := fmt.Sprintf("%s/api", ts.URL)

	cl, err := internal_openapi.NewClientWithResponses(server)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return cl, nil
}

func Setup() fx.Option {
	return fx.Options(
		openapi.Build(),
		fx.Provide(internal_http.NewRouter, newHttpTestServer, newClient),
	)
}
