package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
)

type linkSearcher struct {
	ec *ent.Client
}

func (s *linkSearcher) Search(ctx context.Context, query string) ([]*semdex.Result, error) {
	lq := s.ec.Link.Query().Where(
		link.Or(
			link.TitleContainsFold(query),
			link.DescriptionContainsFold(query),
			link.URLContainsFold(query),
		),
	)

	r, err := lq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := dt.Map(datagraph.LinksFromModel(r), indexableToResult)

	return results, nil
}
