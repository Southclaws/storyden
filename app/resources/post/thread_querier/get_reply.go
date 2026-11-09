package thread_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (q *Querier) GetReply(ctx context.Context, threadID post.ID, replyID post.ID, accountID opt.Optional[account.AccountID]) (*reply.Reply, error) {
	p, err := q.db.Post.Query().
		Where(
			ent_post.ID(xid.ID(replyID)),
			ent_post.RootPostID(xid.ID(threadID)),
			ent_post.DeletedAtIsNil(),
			ent_post.HasRootWith(ent_post.DeletedAtIsNil()),
		).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets().
		WithReplyTo(func(pq *ent.PostQuery) {
			pq.Where(
				ent_post.DeletedAtIsNil(),
				ent_post.RootPostID(xid.ID(threadID)),
			).
				WithAuthor().
				WithReacts(func(rq *ent.ReactQuery) {
					rq.WithAccount()
				})
		}).
		WithReacts(func(rq *ent.ReactQuery) {
			rq.WithAccount()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	posts := []*ent.Post{p}
	if target := p.Edges.ReplyTo; target != nil {
		target.Edges.Root = p.Edges.Root
		posts = append(posts, target)
	}

	ids := make([]xid.ID, 0, len(posts))
	authors := []*ent.Account{p.Edges.Root.Edges.Author}
	var reactEdges []*ent.React
	for _, item := range posts {
		ids = append(ids, item.ID)
		authors = append(authors, item.Edges.Author)
		reactEdges = append(reactEdges, item.Edges.Reacts...)

		for _, react := range item.Edges.Reacts {
			authors = append(authors, react.Edges.Account)
		}
	}

	if err := q.roleQuerier.HydrateRoleEdges(ctx, authors...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := q.getLikesStatus(ctx, ids, accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reacts, err := reaction.MapList(reactEdges)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return reply.Mapper(account.NewAccountLookup(authors), likes, reaction.Reacts(reacts).Map())(p)
}
