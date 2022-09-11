// Package transport contains all the transport layers that facilitate
// interfacing with the application. The main transport method is HTTP which is
// implemented using code generated from an OpenAPI specification.
package transport

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/graphql"
	"github.com/Southclaws/storyden/app/transports/http"
)

func Build() fx.Option {
	return fx.Options(
		http.Build(),
		graphql.Build(),
	)
}
