package role_assign

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
	ent_accountroles "github.com/Southclaws/storyden/internal/ent/accountroles"
)

// ResolveRoleIDs resolves account<>role assignments using internal cache with
// DB fallback for misses.
func (w *Assignment) ResolveRoleIDs(ctx context.Context, accountIDs []xid.ID) (map[xid.ID][]xid.ID, map[xid.ID][]*ent.AccountRoles, error) {
	idsByAccount := make(map[xid.ID][]xid.ID, len(accountIDs))
	byAccountRows := map[xid.ID][]*ent.AccountRoles{}
	missing := make([]xid.ID, 0, len(accountIDs))

	for _, accountID := range accountIDs {
		if rows, ok := w.getRoleAssignmentsCache(ctx, accountID); ok {
			roleIDs := dt.Map(rows, func(ar *ent.AccountRoles) xid.ID { return ar.RoleID })
			idsByAccount[accountID] = roleIDs
			byAccountRows[accountID] = rows
			continue
		}

		missing = append(missing, accountID)
	}

	if len(missing) == 0 {
		return idsByAccount, byAccountRows, nil
	}

	loadedIDs, loadedRows, err := w.ResolveRoleIDsFresh(ctx, missing)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	for accountID, roleIDs := range loadedIDs {
		idsByAccount[accountID] = roleIDs
	}
	for accountID, rows := range loadedRows {
		byAccountRows[accountID] = rows
	}

	return idsByAccount, byAccountRows, nil
}

func (w *Assignment) ResolveRoleIDsFresh(ctx context.Context, accountIDs []xid.ID) (map[xid.ID][]xid.ID, map[xid.ID][]*ent.AccountRoles, error) {
	idsByAccount := make(map[xid.ID][]xid.ID, len(accountIDs))
	if len(accountIDs) == 0 {
		return idsByAccount, map[xid.ID][]*ent.AccountRoles{}, nil
	}

	accountRoles, err := w.db.AccountRoles.Query().
		Where(ent_accountroles.AccountIDIn(accountIDs...)).
		Order(ent.Asc(ent_accountroles.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	byAccountRows := lo.GroupBy(accountRoles, func(ar *ent.AccountRoles) xid.ID { return ar.AccountID })
	for _, accountID := range accountIDs {
		rows := byAccountRows[accountID]
		roleIDs := dt.Map(rows, func(ar *ent.AccountRoles) xid.ID { return ar.RoleID })
		idsByAccount[accountID] = roleIDs
		if err := w.storeRoleAssignmentsCache(ctx, accountID, rows); err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return idsByAccount, byAccountRows, nil
}
