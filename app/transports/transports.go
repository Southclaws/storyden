// Package transports contains all the transport layers that facilitate
// interfacing with the application. The main transport method is OpenAPI which
// is implemented using HTTP and code generated from an OpenAPI specification.
package transports

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http"
	"github.com/Southclaws/storyden/app/transports/mcp"
	"github.com/Southclaws/storyden/app/transports/sse"
)

func Build() fx.Option {
	return fx.Options(
		http.Build(),
		mcp.Build(),
		sse.Build(),
	)
}
