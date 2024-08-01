package node_read

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/services/account/session"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type HydratedQuerier struct {
	session    session.SessionProvider
	nodereader library.Repository
	scorer     semdex.RelevanceScorer
}

func New(
	session session.SessionProvider,
	nodereader library.Repository,
	scorer semdex.RelevanceScorer,
) *HydratedQuerier {
	return &HydratedQuerier{
		session:    session,
		nodereader: nodereader,
		scorer:     scorer,
	}
}

func (q *HydratedQuerier) GetBySlug(ctx context.Context, slug library.NodeSlug) (*library.Node, error) {
	session := q.session.AccountOpt(ctx)

	n, err := q.nodereader.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if acc, ok := session.Get(); ok && q.scorer != nil {
		pro := profile.ProfileFromAccount(&acc)
		nid := xid.ID(n.ID)

		scores, err := q.scorer.ScoreRelevance(ctx, pro, nid)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		score, ok := scores[nid]
		if !ok {
			return n, nil
		}

		n.RelevanceScore = opt.New(score)

		// TODO: Hydrate recommendations
	}

	return n, nil
}
