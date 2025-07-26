package role_querier

import (
	"context"
	"sort"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) Get(ctx context.Context, id role.RoleID) (*role.Role, error) {
	r, err := q.db.Role.Get(ctx, xid.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rl, err := role.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rl, nil
}

func (q *Querier) List(ctx context.Context) (role.Roles, error) {
	roles, err := q.db.Role.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := role.MapList(roles)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defaultRole, err := q.GetDefaultRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped = append(mapped, defaultRole)

	sort.Sort(mapped)

	return mapped, nil
}

func (q *Querier) ListFor(ctx context.Context, account *ent.Account) (held.Roles, error) {
	roles, err := q.db.AccountRoles.
		Query().
		Where(
			ent_account_role.AccountID(account.ID),
		).
		WithRole(func(rq *ent.RoleQuery) {
			rq.Order(ent.Asc(ent_role.FieldSortKey))
		}).
		Order(ent.Asc(ent_account_role.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := held.MapList(roles, account.Admin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	dr, drExists, err := q.lookupDefaultRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// If the default member role has not been modified (aka not added to the DB
	// with custom permissions) we add the default manually.
	if drExists {
		defaultRole, err := role.Map(dr)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		mapped = append(mapped, &held.Role{
			Role:     *defaultRole,
			Assigned: account.CreatedAt,
			Badge:    false,
			Default:  true,
		})
	} else {
		mapped = append(mapped, &held.Role{
			Role: role.DefaultRoleEveryone,
		})
	}

	// TODO: Implement sorting on API - currently it's pointless.
	// sort.Sort(mapped)

	return mapped, nil
}

func (q *Querier) GetDefaultRole(ctx context.Context) (*role.Role, error) {
	dr, drExists, err := q.lookupDefaultRole(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !drExists {
		return &role.DefaultRoleEveryone, nil
	}

	return role.Map(dr)
}

func (q *Querier) lookupDefaultRole(ctx context.Context) (*ent.Role, bool, error) {
	defaultRole, err := q.db.Role.Get(ctx, xid.ID(role.DefaultRoleEveryoneID))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return defaultRole, true, nil
}
