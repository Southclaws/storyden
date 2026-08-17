package like_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/like/item_like"
	"github.com/Southclaws/storyden/app/resources/like/profile_like"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entlikepost "github.com/Southclaws/storyden/internal/ent/likepost"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
)

type LikeQuerier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *LikeQuerier {
	return &LikeQuerier{db: db, roleQuerier: roleQuerier}
}

func (l *LikeQuerier) GetPostLikes(ctx context.Context, postID post.ID) ([]*item_like.Like, error) {
	r, err := l.db.LikePost.
		Query().
		Where(entlikepost.HasPostWith(entpost.ID(xid.ID(postID)))).
		WithAccount().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(r))
	for _, liked := range r {
		if liked.Edges.Account != nil {
			roleTargets = append(roleTargets, liked.Edges.Account)
		}
	}
	if err := l.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := dt.MapErr(r, item_like.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return likes, nil
}

func (l *LikeQuerier) GetProfileLikes(ctx context.Context, accountID account.AccountID, params pagination.Parameters) (*pagination.Result[*profile_like.Like], error) {
	total, err := l.db.LikePost.Query().Where(entlikepost.HasAccountWith(entaccount.ID(xid.ID(accountID)))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	q := l.db.LikePost.Query().
		Limit(params.Limit()).
		Offset(params.Offset()).
		Order(ent.Desc(entlikepost.FieldCreatedAt), ent.Desc(entlikepost.FieldID)).
		Where(entlikepost.HasAccountWith(entaccount.ID(xid.ID(accountID)))).
		WithPost(func(pq *ent.PostQuery) {
			pq.WithAuthor()
			pq.WithCategory()
			pq.WithTags()
			pq.WithRoot()
		})

	r, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(r))
	for _, liked := range r {
		if liked.Edges.Post == nil {
			continue
		}

		if author := liked.Edges.Post.Edges.Author; author != nil {
			roleTargets = append(roleTargets, author)
		}
	}
	if err := l.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := dt.MapErr(r, profile_like.Map)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, likes)

	return &result, nil
}
