package follow_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountfollow"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (q *Querier) GetFollowers(ctx context.Context, id account.AccountID, params pagination.Parameters) (*pagination.Result[*profile.Ref], error) {
	total, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowingAccountID(xid.ID(id))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowingAccountID(xid.ID(id))).
		Limit(params.Limit()).
		Offset(params.Offset()).
		Order(ent.Desc(accountfollow.FieldCreatedAt), ent.Desc(accountfollow.FieldID)).
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

	profiles, err := dt.MapErr(r, func(in *ent.AccountFollow) (*profile.Ref, error) {
		return profile.MapRef(in.Edges.Follower)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, profiles)

	return &result, nil
}

func (q *Querier) GetFollowing(ctx context.Context, id account.AccountID, params pagination.Parameters) (*pagination.Result[*profile.Ref], error) {
	total, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowerAccountID(xid.ID(id))).Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := q.db.AccountFollow.Query().
		Where(accountfollow.FollowerAccountID(xid.ID(id))).
		Limit(params.Limit()).
		Offset(params.Offset()).
		Order(ent.Desc(accountfollow.FieldCreatedAt), ent.Desc(accountfollow.FieldID)).
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

	profiles, err := dt.MapErr(r, func(in *ent.AccountFollow) (*profile.Ref, error) {
		return profile.MapRef(in.Edges.Following)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, profiles)

	return &result, nil
}
