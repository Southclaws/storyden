package node_read

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/services/account/session"
)

type HydratedQuerier struct {
	logger     *zap.Logger
	session    session.SessionProvider
	nodereader *node_querier.Querier
	scorer     semdex.RelevanceScorer
}

func New(
	logger *zap.Logger,
	session session.SessionProvider,
	nodereader *node_querier.Querier,
	scorer semdex.RelevanceScorer,
) *HydratedQuerier {
	return &HydratedQuerier{
		logger:     logger,
		session:    session,
		nodereader: nodereader,
		scorer:     scorer,
	}
}

func (q *HydratedQuerier) GetBySlug(ctx context.Context, qk library.QueryKey) (*library.Node, error) {
	session := q.session.AccountOpt(ctx)

	opts := []node_querier.Option{}

	if s, ok := session.Get(); ok {
		opts = append(opts, node_querier.WithVisibilityRulesApplied(&s.ID))
	} else {
		opts = append(opts, node_querier.WithVisibilityRulesApplied(nil))
	}

	n, err := q.nodereader.Get(ctx, qk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if acc, ok := session.Get(); ok && q.scorer != nil {
		pro := profile.ProfileFromAccount(&acc)
		nid := xid.ID(n.Mark.ID())

		scores, err := q.scorer.ScoreRelevance(ctx, pro, nid)
		if err != nil {
			q.logger.Warn("failed to score relevance", zap.Error(err))
		}

		score, ok := scores[nid]
		if ok {
			n.RelevanceScore = opt.New(score)
		}

		// TODO: Hydrate recommendations
	}

	return n, nil
}
