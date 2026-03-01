package role_assign

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	ent_accountroles "github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Assignment struct {
	ins   spanner.Instrumentation
	db    *ent.Client
	store cache.Store
}

func New(
	ins spanner.Builder,
	db *ent.Client,
	store cache.Store,
) *Assignment {
	return &Assignment{
		ins:   ins.Build(),
		db:    db,
		store: store,
	}
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

func (w *Assignment) UpdateRoles(ctx context.Context, accountID account_ref.ID, roles ...Mutation) error {
	ctx, span := w.ins.Instrument(ctx,
		kv.String("account_id", xid.ID(accountID).String()),
		kv.Int("mutations_count", len(roles)),
	)
	defer span.End()

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer func() { _ = tx.Rollback() }()

	update := tx.Account.UpdateOneID(xid.ID(accountID))
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

	_, err = update.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		return fault.Wrap(err, fctx.With(ctx))
	}

	roleIDs, err := tx.AccountRoles.Query().
		Where(ent_accountroles.AccountIDEQ(xid.ID(accountID))).
		All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := w.storeRoleAssignmentsCache(ctx, xid.ID(accountID), roleIDs); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		_ = w.deleteRoleIDsCache(ctx, xid.ID(accountID))
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
