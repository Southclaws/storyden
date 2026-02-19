package role_repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_account_role "github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
)

func (h *Repository) syncAll(ctx context.Context) error {
	roles, err := h.db.Role.Query().Order(ent.Asc(ent_role.FieldSortKey)).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	mappedRoles, err := role.MapList(roles)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	stored := map[role.RoleID]struct{}{}
	customIDs := make([]role.RoleID, 0, len(mappedRoles))

	for _, rl := range mappedRoles {
		if err := h.storeRole(ctx, rl, true); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		stored[rl.ID] = struct{}{}
		if !isDefaultRoleID(rl.ID) {
			customIDs = append(customIDs, rl.ID)
		}
	}

	for _, defaultID := range []role.RoleID{
		role.DefaultRoleGuestID,
		role.DefaultRoleMemberID,
		role.DefaultRoleAdminID,
	} {
		if _, ok := stored[defaultID]; ok {
			continue
		}

		defaultRole, ok := defaultRole(defaultID)
		if !ok {
			continue
		}

		defaultCopy := defaultRole
		if err := h.storeRole(ctx, &defaultCopy, false); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := h.storeCustomRoleOrdering(ctx, customIDs); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	accountRows, err := h.db.Account.Query().Select(ent_account.FieldID).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	relationships, err := h.db.AccountRoles.Query().Where(
		ent_account_role.RoleIDNotIn(
			xid.ID(role.DefaultRoleGuestID),
			xid.ID(role.DefaultRoleMemberID),
			xid.ID(role.DefaultRoleAdminID),
		),
	).Order(
		ent.Asc(ent_account_role.FieldAccountID),
		ent.Asc(ent_account_role.FieldCreatedAt),
	).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	assignmentsByAccount := map[xid.ID][]cachedAssignment{}
	for _, relationship := range relationships {
		accountID := xid.ID(relationship.AccountID)
		assignmentsByAccount[accountID] = append(assignmentsByAccount[accountID], cachedAssignment{
			RoleID:     relationship.RoleID.String(),
			AssignedAt: relationship.CreatedAt,
			Badge:      relationship.Badge != nil && *relationship.Badge,
		})
	}

	for _, accountRow := range accountRows {
		accountID := xid.ID(accountRow.ID)
		assignments := assignmentsByAccount[accountID]
		if assignments == nil {
			assignments = []cachedAssignment{}
		}

		if err := h.storeAssignments(ctx, accountID, assignments); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (h *Repository) recoverFromCacheWriteFailure(ctx context.Context, operation string, cacheWriteErr error) error {
	slog.WarnContext(ctx, "role repository cache write failed, running full sync",
		slog.String("operation", operation),
		slog.String("error", cacheWriteErr.Error()))

	if err := h.syncAll(ctx); err != nil {
		return fault.Wrap(errors.Join(cacheWriteErr, err), fctx.With(ctx))
	}

	slog.WarnContext(ctx, "role repository full sync recovered cache state",
		slog.String("operation", operation))

	return nil
}

func isDefaultRoleID(id role.RoleID) bool {
	switch id {
	case role.DefaultRoleGuestID, role.DefaultRoleMemberID, role.DefaultRoleAdminID:
		return true
	default:
		return false
	}
}
