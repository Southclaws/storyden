package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Southclaws/fault"
	"go.uber.org/fx"
	"go.uber.org/zap"

	http_transport "github.com/Southclaws/storyden/app/transports/http"
	"github.com/Southclaws/storyden/app/transports/http/bindings"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
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

func newClient(ts *httptest.Server) (*openapi.ClientWithResponses, error) {
	server := fmt.Sprintf("%s/api", ts.URL)

	cl, err := openapi.NewClientWithResponses(server)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return cl, nil
}

func Setup() fx.Option {
	return fx.Options(
		// In the normal app, we call http.Build() which constructs a production
		// HTTP server with the http.ServeMux router as the handler. In tests we
		// don't want this, instead we want the httptest.Server instead. So this
		// setup looks very similar to http.Build() but instead of calling the
		// httpserver.Build() provider, we provide the http.ServeMux and then we
		// mount it onto the httptest.Server instead of http.Server.
		//
		fx.Provide(httpserver.NewRouter, newHttpTestServer, newClient),

		fx.Provide(session.New),

		bindings.Build(),

		fx.Invoke(http_transport.MountOpenAPI),
	)
}
