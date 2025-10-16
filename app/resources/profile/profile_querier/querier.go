package profile_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, roleQuerier *role_querier.Querier) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

func (d *Querier) GetByID(ctx context.Context, id account.AccountID) (*profile.Public, error) {
	q := d.db.Account.
		Query().
		Where(account_ent.ID(xid.ID(id))).
		WithTags().
		WithEmails().
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		WithAuthentication()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	hr, err := d.roleQuerier.ListFor(ctx, result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := profile.Map(hr)(result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err = hydrateEdgeAggregations(ctx, result, acc)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (d *Querier) LookupByHandle(ctx context.Context, handle string) (*profile.Public, bool, error) {
	q := d.db.Account.
		Query().
		Where(account_ent.Handle(handle)).
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator()
		}).
		WithTags()

	result, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	hr, err := d.roleQuerier.ListFor(ctx, result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := profile.Map(hr)(result)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err = hydrateEdgeAggregations(ctx, result, acc)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, true, nil
}

func (d *Querier) GetMany(ctx context.Context, ids ...account.AccountID) ([]*profile.Public, error) {
	xids := dt.Map(ids, func(id account.AccountID) xid.ID { return xid.ID(id) })

	accounts, err := d.db.Account.
		Query().
		Where(account_ent.IDIn(xids...)).
		WithTags().
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
		WithInvitedBy(func(iq *ent.InvitationQuery) {
			iq.WithCreator(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Optimise: roleQuerier needs a mapping lookup, edge aggregations
	// are probably not needed downstream.
	profiles := make([]*profile.Public, 0, len(accounts))
	for _, a := range accounts {
		hr, err := d.roleQuerier.ListFor(ctx, a)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		acc, err := profile.Map(hr)(a)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		acc, err = hydrateEdgeAggregations(ctx, a, acc)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		profiles = append(profiles, acc)
	}

	return profiles, nil
}

func hydrateEdgeAggregations(ctx context.Context, a *ent.Account, acc *profile.Public) (*profile.Public, error) {
	following, err := a.QueryFollowing().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	followers, err := a.QueryFollowedBy().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likes, err := a.QueryPosts().QueryLikes().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc.Followers = followers
	acc.Following = following
	acc.LikeScore = likes

	return acc, nil
}
