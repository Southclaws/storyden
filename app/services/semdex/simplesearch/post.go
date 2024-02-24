package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
)

type postSearcher struct {
	ec *ent.Client
}

func (s *postSearcher) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	pq := s.ec.Post.Query().Where(
		post.Or(
			post.TitleContainsFold(query),
			post.BodyContainsFold(query),
		),
	).WithRoot()

	rs, err := pq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results, err := dt.MapErr(rs, func(p *ent.Post) (*datagraph.NodeReference, error) {
		if p.Edges.Root == nil {
			return &datagraph.NodeReference{
				ID:          p.ID,
				Kind:        datagraph.KindThread,
				Name:        p.Title,
				Description: p.Short,
				Slug:        p.Slug,
			}, nil
		} else {
			return &datagraph.NodeReference{
				ID:          p.ID,
				Kind:        datagraph.KindReply,
				Name:        p.Edges.Root.Title,
				Description: p.Short,
				Slug:        p.Edges.Root.Slug,
			}, nil
		}
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return results, nil
}
