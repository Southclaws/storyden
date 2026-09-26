package node_read

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type HydratedQuerier struct {
	nodereader *node_querier.Querier
}

func New(
	nodereader *node_querier.Querier,
) *HydratedQuerier {
	return &HydratedQuerier{
		nodereader: nodereader,
	}
}

func (q *HydratedQuerier) GetBySlug(ctx context.Context, qk library.QueryKey, sortChildrenBy opt.Optional[node_querier.ChildSortRule]) (*library.Node, error) {
	requestingAccount := session.GetOptAccount(ctx)
	visibilityRules := node_querier.WithVisibilityRulesApplied(requestingAccount)
	opts := []node_querier.Option{visibilityRules}

	sortChildrenBy.Call(func(v node_querier.ChildSortRule) {
		opts = append(opts, node_querier.WithSortChildrenBy(v))
	})

	n, err := q.nodereader.Get(ctx, qk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ancestors, err := q.nodereader.Ancestry(ctx, library.NodeID(n.Mark.ID()), visibilityRules)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n.Ancestors = ancestors

	return n, nil
}

func (q *HydratedQuerier) ListChildren(ctx context.Context, qk library.QueryKey, pp pagination.Parameters, opts ...node_querier.Option) (*pagination.Result[*library.Node], error) {
	requestingAccount := session.GetOptAccount(ctx)
	opts = append(opts, node_querier.WithVisibilityRulesApplied(requestingAccount))

	r, err := q.nodereader.ListChildren(ctx, qk, pp, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}
