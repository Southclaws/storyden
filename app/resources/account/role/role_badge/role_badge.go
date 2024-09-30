package role_badge

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	account_ent "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (w *Writer) Update(ctx context.Context, accountID account.AccountID, roleID role.RoleID, badge bool) (*account.Account, error) {
	predicate := []predicate.AccountRoles{
		accountroles.AccountIDEQ(xid.ID(accountID)),
		accountroles.RoleIDEQ(xid.ID(roleID)),
	}

	// Only one role can be set as a badge, clear all first, then set if true.

	err := w.db.AccountRoles.Update().
		Where(predicate...).
		SetBadge(false).
		Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if badge {
		err = w.db.AccountRoles.Update().
			Where(predicate...).
			SetBadge(true).
			Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
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
