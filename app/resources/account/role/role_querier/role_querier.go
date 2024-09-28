package role_querier

import (
	"context"
	"sort"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
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

	mapped = append(mapped, &role.DefaultRoleAdmin, &role.DefaultRoleEveryone)

	sort.Sort(mapped)

	return mapped, nil
}
