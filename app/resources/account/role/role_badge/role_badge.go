package role_badge

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type Writer struct {
	db             *ent.Client
	accountQuerier *account_querier.Querier
}

func New(db *ent.Client, accountQuerier *account_querier.Querier) *Writer {
	return &Writer{db: db, accountQuerier: accountQuerier}
}

func (w *Writer) Update(ctx context.Context, accountID account.AccountID, roleID role.RoleID, badge bool) (*account.AccountWithEdges, error) {
	predicate := []predicate.AccountRoles{
		accountroles.AccountIDEQ(xid.ID(accountID)),
		accountroles.RoleIDEQ(xid.ID(roleID)),
	}

	// Only one role can be set as a badge, clear all first, then set if true.

	err := w.db.AccountRoles.Update().
		Where(accountroles.AccountIDEQ(xid.ID(accountID))).
		ClearBadge().
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

	return w.accountQuerier.GetByID(ctx, accountID)
}
