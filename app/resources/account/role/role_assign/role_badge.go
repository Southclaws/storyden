package role_assign

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_accountroles "github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

func (w *Assignment) SetBadge(ctx context.Context, accountID account_ref.ID, roleID role.RoleID, badge bool) error {
	tx, err := w.db.Tx(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer func() { _ = tx.Rollback() }()

	predicate := []predicate.AccountRoles{
		accountroles.AccountIDEQ(xid.ID(accountID)),
		accountroles.RoleIDEQ(xid.ID(roleID)),
	}

	// Only one role can be set as a badge, clear all first, then set if true.

	err = tx.AccountRoles.Update().
		Where(accountroles.AccountIDEQ(xid.ID(accountID))).
		ClearBadge().
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if badge {
		err = tx.AccountRoles.Update().
			Where(predicate...).
			SetBadge(true).
			Exec(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
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
		_ = w.invalidateRoleIDsCache(ctx, xid.ID(accountID))
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
