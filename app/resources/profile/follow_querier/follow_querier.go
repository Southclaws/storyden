package follow_querier

import (
	"context"
	"math"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountfollow"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, roleQuerier *role_querier.Querier) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Profiles    []*profile.Ref
}

func (q *Querier) GetFollowers(ctx context.Context, id account.AccountID, page, size int) (*Result, error) {
	total, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowingAccountID(xid.ID(id))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowingAccountID(xid.ID(id))).
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Desc(accountfollow.FieldCreatedAt)).
		WithFollower().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(r))
	for _, follow := range r {
		if follow.Edges.Follower != nil {
			roleTargets = append(roleTargets, follow.Edges.Follower)
		}
	}
	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nextPage := opt.NewSafe(page+1, len(r) >= size)
	if len(r) > 1 {
		r = r[:len(r)-1]
	}

	profiles, err := dt.MapErr(r, func(in *ent.AccountFollow) (*profile.Ref, error) {
		return profile.MapRef(in.Edges.Follower)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Result{
		PageSize:    size,
		Results:     len(profiles),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Profiles:    profiles,
	}, nil
}

func (q *Querier) GetFollowing(ctx context.Context, id account.AccountID, page, size int) (*Result, error) {
	total, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowerAccountID(xid.ID(id))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowerAccountID(xid.ID(id))).
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Desc(accountfollow.FieldCreatedAt)).
		WithFollowing().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleTargets := make([]*ent.Account, 0, len(r))
	for _, follow := range r {
		if follow.Edges.Following != nil {
			roleTargets = append(roleTargets, follow.Edges.Following)
		}
	}
	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nextPage := opt.NewSafe(page+1, len(r) >= size)
	if len(r) > 1 {
		r = r[:len(r)-1]
	}

	profiles, err := dt.MapErr(r, func(in *ent.AccountFollow) (*profile.Ref, error) {
		return profile.MapRef(in.Edges.Following)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Result{
		PageSize:    size,
		Results:     len(profiles),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Profiles:    profiles,
	}, nil
}
