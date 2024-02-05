package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/item"
)

type itemSearcher struct {
	ec *ent.Client
}

func (s *itemSearcher) Search(ctx context.Context, query string) ([]*semdex.Result, error) {
	iq := s.ec.Item.Query().Where(
		item.Or(
			item.NameContainsFold(query),
			item.DescriptionContainsFold(query),
			item.ContentContainsFold(query),
		),
	)

	rs, err := iq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := dt.MapErr(rs, datagraph.ItemFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := dt.Map(items, indexableToResult)

	return results, nil
}
