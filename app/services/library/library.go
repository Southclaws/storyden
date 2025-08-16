package library

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_property_schema"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_semdex"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(node_read.New, node_mutate.New, nodetree.New, node_visibility.New, node_property_schema.New),
		node_semdex.Build(),
	)
}
