package transport

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/transports/http"
)

func Build() fx.Option {
	return fx.Options(
		http.Build(),
	)
}
