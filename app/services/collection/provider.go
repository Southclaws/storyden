package collection

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/collection/collection_item_manager"
	"github.com/Southclaws/storyden/app/services/collection/collection_manager"
	"github.com/Southclaws/storyden/app/services/collection/collection_read"
)

func Build() fx.Option {
	return fx.Provide(
		collection_item_manager.New,
		collection_manager.New,
		collection_read.New,
	)
}
