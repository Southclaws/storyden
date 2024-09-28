package role_assign

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
)

type Assignment struct {
	db *ent.Client
}

func New(db *ent.Client) *Assignment {
	return &Assignment{db: db}
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

func (w *Assignment) UpdateRoles(ctx context.Context, accountID account.AccountID, roles ...Mutation) (*account.Account, error) {
	update := w.db.Account.UpdateOneID(xid.ID(accountID))
	mutation := update.Mutation()

	adds, removes, admin := split(roles...)

	mutation.AddRoleIDs(adds...)
	mutation.RemoveRoleIDs(removes...)

	if a, ok := admin.Get(); ok {
		mutation.SetAdmin(a)
	}

	_, err := update.Save(ctx)
	if err != nil && !ent.IsConstraintError(err) {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := w.db.Account.
		Query().
		Where(account_ent.ID(xid.ID(accountID))).
		WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := account.MapAccount(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
