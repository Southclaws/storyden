package like_querier

import (
	"context"
	"math"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/like/item_like"
	"github.com/Southclaws/storyden/app/resources/like/profile_like"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entlikepost "github.com/Southclaws/storyden/internal/ent/likepost"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
)

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Likes       []*profile_like.Like
}

type LikeQuerier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *LikeQuerier {
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

	accountEdges := dt.Map(r, func(l *ent.LikePost) *ent.Account { return l.Edges.Account })
	roleHydrator, err := l.roleQuerier.BuildMultiHydrator(ctx, accountEdges)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := dt.MapErr(r, func(in *ent.LikePost) (*item_like.Like, error) {
		return item_like.Map(in, roleHydrator.Hydrate)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return likes, nil
}

func (l *LikeQuerier) GetProfileLikes(ctx context.Context, accountID account.AccountID, page int, size int) (*Result, error) {
	total, err := l.db.LikePost.Query().Where(entlikepost.HasAccountWith(entaccount.ID(xid.ID(accountID)))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	q := l.db.LikePost.Query().
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Desc(entlikepost.FieldCreatedAt)).
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

	nextPage := opt.NewSafe(page+1, len(r) >= size)
	if len(r) > 1 {
		r = r[:len(r)-1]
	}

	postAuthorEdges := dt.Map(r, func(l *ent.LikePost) *ent.Account {
		if l.Edges.Post == nil {
			return nil
		}
		return l.Edges.Post.Edges.Author
	})
	roleHydrator, err := l.roleQuerier.BuildMultiHydrator(ctx, postAuthorEdges)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := dt.MapErr(r, func(in *ent.LikePost) (*profile_like.Like, error) {
		return profile_like.Map(in, roleHydrator.Hydrate)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Result{
		PageSize:    size,
		Results:     len(likes),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Likes:       likes,
	}, nil
}
