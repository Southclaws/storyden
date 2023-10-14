package item

import (
	"context"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type ItemManager interface {
	Create(ctx context.Context) (*datagraph.Item, error)
	Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
	Update(ctx context.Context, slug datagraph.ItemSlug, p Partial) (*datagraph.Item, error)
	Archive(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	ImageURL    opt.Optional[string]
	Description opt.Optional[string]
}
