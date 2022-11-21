// Package transport contains all the transport layers that facilitate
// interfacing with the application. The main transport method is OpenAPI which
// is implemented using HTTP and code generated from an OpenAPI specification.
package transport

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/graphql"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/http"
)

func Build() fx.Option {
	return fx.Options(
		http.Build(),

		openapi.Build(),
		graphql.Build(),
	)
}
