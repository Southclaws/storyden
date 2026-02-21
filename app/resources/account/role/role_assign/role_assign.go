package role_assign

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Assignment struct {
	db             *ent.Client
	accountQuerier *account_querier.Querier
	profileCache   *profile_cache.Cache
	bus            *pubsub.Bus
}

func New(db *ent.Client, accountQuerier *account_querier.Querier, profileCache *profile_cache.Cache, bus *pubsub.Bus) *Assignment {
	return &Assignment{db: db, accountQuerier: accountQuerier, profileCache: profileCache, bus: bus}
}

type Mutation struct {
	id     role.RoleID
	delete bool
}

func Add(id role.RoleID) Mutation {
	return Mutation{id: id}
}

func Remove(id role.RoleID) Mutation {
	return Mutation{id: id, delete: true}
}

func (m Mutation) xid() xid.ID { return xid.ID(m.id) }

func split(mutations ...Mutation) (adds, removes []xid.ID, admin opt.Optional[bool]) {
	for _, m := range mutations {
		if m.delete {
			if m.id == role.DefaultRoleAdminID {
				admin = opt.New(false)
			} else {
				removes = append(removes, m.xid())
			}
		} else {
			if m.id == role.DefaultRoleAdminID {
				admin = opt.New(true)
			} else {
				adds = append(adds, m.xid())
			}
		}
	}
	return
}

func (w *Assignment) UpdateRoles(ctx context.Context, accountID account.AccountID, roles ...Mutation) (*account.AccountWithEdges, error) {
	update := w.db.Account.UpdateOneID(xid.ID(accountID))
	mutation := update.Mutation()

	roles = dt.Filter(roles, func(m Mutation) bool {
		return m.id != role.DefaultRoleMemberID
	})

	adds, removes, admin := split(roles...)

	mutation.AddRoleIDs(adds...)
	mutation.RemoveRoleIDs(removes...)

	if a, ok := admin.Get(); ok {
		mutation.SetAdmin(a)
	}

	err := w.profileCache.Invalidate(ctx, xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = update.Save(ctx)
	if err != nil && !ent.IsConstraintError(err) {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	w.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: accountID,
	})

	return w.accountQuerier.GetByID(ctx, accountID)
}
